package livemapper

import (
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	permslivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils"
	errorslivemapper "github.com/fivenet-app/fivenet/services/livemapper/errors"
	"go.uber.org/zap"
)

const (
	userMarkerChunkSize = 20
)

func (s *Server) Stream(req *pblivemapper.StreamRequest, srv pblivemapper.LivemapperService_StreamServer) error {
	userInfo := auth.MustGetUserInfoFromContext(srv.Context())

	s.logger.Debug("starting livemap stream", zap.Int32("user_id", userInfo.UserId))
	markerJobsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamMarkersPermField)
	if err != nil {
		return errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}
	userJobsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamPlayersPermField)
	if err != nil {
		return errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}

	var markersJobs []string
	if markerJobsAttr != nil {
		markersJobs = markerJobsAttr.([]string)
	}
	if userInfo.SuperUser {
		s.markersCache.Range(func(job string, _ []*livemap.MarkerMarker) bool {
			markersJobs = append(markersJobs, job)
			return true
		})
		markersJobs = utils.RemoveSliceDuplicates(markersJobs)
	}

	var usersJobs map[string]int32
	if userJobsAttr != nil {
		usersJobs, _ = userJobsAttr.(map[string]int32)
	}

	if userInfo.SuperUser {
		usersJobs = map[string]int32{}
		for _, j := range s.appCfg.Get().UserTracker.GetLivemapJobs() {
			usersJobs[j] = -1
		}
	}

	// Prepare jobs for client response
	jobs := &pblivemapper.StreamResponse_Jobs{
		Jobs: &pblivemapper.JobsList{
			Markers: make([]*users.Job, len(markersJobs)),
			Users:   []*users.Job{},
		},
	}

	for i := 0; i < len(markersJobs); i++ {
		jobs.Jobs.Markers[i] = &users.Job{
			Name: markersJobs[i],
		}
		s.enricher.EnrichJobName(jobs.Jobs.Markers[i])
	}
	for job := range usersJobs {
		j := &users.Job{
			Name: job,
		}
		s.enricher.EnrichJobName(j)
		jobs.Jobs.Users = append(jobs.Jobs.Users, j)
	}

	if err := srv.Send(&pblivemapper.StreamResponse{
		Data: jobs,
	}); err != nil {
		return err
	}

	if end, err := s.sendMarkerMarkers(srv, markersJobs); end || err != nil {
		return err
	}

	if end, err := s.sendChunkedUserMarkers(srv, usersJobs, userInfo); end || err != nil {
		return err
	}

	updateCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(updateCh)

	for {
		select {
		case <-srv.Context().Done():
			return nil

		case event := <-updateCh:
			if event == nil {
				continue
			}

			if event.Send == MarkerUpdate {
				if end, err := s.sendMarkerMarkers(srv, markersJobs); end || err != nil {
					return err
				}
			}

		case <-time.After(s.appCfg.Get().UserTracker.RefreshTime.AsDuration()):
			if end, err := s.sendChunkedUserMarkers(srv, usersJobs, userInfo); end || err != nil {
				return err
			}
		}
	}
}

// Sends out chunked current user markers
func (s *Server) sendChunkedUserMarkers(srv pblivemapper.LivemapperService_StreamServer, usersJobs map[string]int32, userInfo *userinfo.UserInfo) (bool, error) {
	userMarkers, onDutyState, err := s.getUserLocations(usersJobs, userInfo)
	if err != nil {
		return true, errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}

	// Less than chunk size or no markers, quick return here
	if len(userMarkers) <= userMarkerChunkSize {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Users{
				Users: &pblivemapper.UserMarkersUpdates{
					Users: userMarkers,
					Part:  0,
				},
			},
			UserOnDuty: &onDutyState,
		}

		if err := srv.Send(resp); err != nil {
			return true, err
		}

		return false, nil
	}

	parts := int32(len(userMarkers) / userMarkerChunkSize)
	for userMarkerChunkSize < len(userMarkers) {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Users{
				Users: &pblivemapper.UserMarkersUpdates{
					Users: userMarkers[0:userMarkerChunkSize:userMarkerChunkSize],
					Part:  parts,
				},
			},
			UserOnDuty: &onDutyState,
		}
		parts--

		if err := srv.Send(resp); err != nil {
			return true, err
		}

		userMarkers = userMarkers[userMarkerChunkSize:]

		select {
		case <-srv.Context().Done():
			return true, nil

		case <-time.After(125 * time.Millisecond):
		}
	}

	if len(userMarkers) > 0 {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Users{
				Users: &pblivemapper.UserMarkersUpdates{
					Users: userMarkers,
					Part:  0,
				},
			},
			UserOnDuty: &onDutyState,
		}
		if err := srv.Send(resp); err != nil {
			return true, err
		}
	}

	return false, nil
}

func (s *Server) sendMarkerMarkers(srv pblivemapper.LivemapperService_StreamServer, jobs []string) (bool, error) {
	markers, err := s.getMarkerMarkers(jobs)
	if err != nil {
		return true, errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}

	// Send current markers
	resp := &pblivemapper.StreamResponse{
		Data: &pblivemapper.StreamResponse_Markers{
			Markers: &pblivemapper.MarkerMarkersUpdates{
				Markers: markers,
			},
		},
	}
	if err := srv.Send(resp); err != nil {
		return true, err
	}

	return false, nil
}

func (s *Server) getUserLocations(jobs map[string]int32, userInfo *userinfo.UserInfo) ([]*livemap.UserMarker, bool, error) {
	ds := []*livemap.UserMarker{}

	found := false
	if userInfo.SuperUser {
		found = true
	}

	for job, grade := range jobs {
		markers, ok := s.tracker.GetUsersByJob(job)
		if !ok {
			continue
		}

		markers.Range(func(key int32, marker *livemap.UserMarker) bool {
			// SuperUser returns grade as `-1`, job has access to that grade or it is the user itself
			if (grade == -1 || (marker.User.JobGrade <= grade || key == userInfo.UserId)) && !marker.Hidden {
				ds = append(ds, marker)
			}

			// If the user is found in the list of user markers and not "off duty" (hidden), set found state
			if !found && !marker.Hidden && (userInfo.Job == job && userInfo.UserId == key) {
				found = true
			}

			return true
		})
	}

	if found {
		return ds, true, nil
	}

	return nil, false, nil
}
