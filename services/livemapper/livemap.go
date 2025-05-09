package livemapper

import (
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
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
	userMarkerChunkSize   = 25
	markerMarkerChunkSize = 50
)

func (s *Server) Stream(req *pblivemapper.StreamRequest, srv pblivemapper.LivemapperService_StreamServer) error {
	origUI := auth.MustGetUserInfoFromContext(srv.Context()).Clone()
	userInfo := &origUI

	s.logger.Debug("starting livemap stream", zap.Int32("user_id", userInfo.UserId))
	markerJobs, err := s.ps.AttrJobList(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamMarkersPermField)
	if err != nil {
		return errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}
	usersJobs, err := s.ps.AttrJobGradeList(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamPlayersPermField)
	if err != nil {
		return errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}

	if userInfo.SuperUser {
		s.markersCache.Range(func(job string, _ []*livemap.MarkerMarker) bool {
			markerJobs.Strings = append(markerJobs.Strings, job)
			return true
		})
		markerJobs.Strings = utils.RemoveSliceDuplicates(markerJobs.Strings)

		// Disable fine-grained permissions for superuser as the it is now just a list of jobs
		usersJobs.FineGrained = false
		for _, j := range s.tracker.ListTrackedJobs() {
			usersJobs.Jobs[j] = -1
		}
	}

	// Prepare jobs for client response
	jobs := &pblivemapper.StreamResponse_Jobs{
		Jobs: &pblivemapper.JobsList{
			Markers: make([]*users.Job, markerJobs.Len()),
			Users:   []*users.Job{},
		},
	}

	for i := range markerJobs.Strings {
		jobs.Jobs.Markers[i] = &users.Job{
			Name: markerJobs.Strings[i],
		}
		s.enricher.EnrichJobName(jobs.Jobs.Markers[i])
	}
	for job := range usersJobs.Iter() {
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

	if end, err := s.sendMarkerMarkers(srv, markerJobs, time.Time{}); end || err != nil {
		return err
	}

	end, lastUpdatedAt, onDutyState, err := s.sendChunkedUserMarkers(srv, usersJobs, userInfo, time.Time{}, false)
	if end || err != nil {
		return err
	}
	updatedAt := lastUpdatedAt

	// Refresh Ticker
	refreshTime := s.appCfg.Get().UserTracker.RefreshTime.AsDuration()
	updateTicker := time.NewTicker(refreshTime)
	defer updateTicker.Stop()

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

			if event.Users != nil {
				if usersJobs.Len() == 0 {
					continue
				}

				// Send delete user markers event to client
				deleted := []int32{}
				for job := range usersJobs.Iter() {
					if _, ok := (*event.Users)[job]; !ok {
						continue
					}

					for _, um := range (*event.Users)[job] {
						if um.Hidden || (userInfo.SuperUser || !usersJobs.HasJobGrade(job, um.User.JobGrade)) {
							continue
						}

						deleted = append(deleted, um.UserId)
					}
				}
				if len(deleted) == 0 {
					continue
				}

				resp := &pblivemapper.StreamResponse{
					Data: &pblivemapper.StreamResponse_Users{
						Users: &pblivemapper.UserMarkersUpdates{
							Deleted: deleted,
							Partial: true,
						},
					},
				}
				if err := srv.Send(resp); err != nil {
					return err
				}
			} else if event.MarkerUpdate != nil {
				// Send delete marker event to client
				resp := &pblivemapper.StreamResponse{
					Data: &pblivemapper.StreamResponse_Markers{
						Markers: &pblivemapper.MarkerMarkersUpdates{
							Updated: []*livemap.MarkerMarker{event.MarkerUpdate},
							Partial: true,
						},
					},
				}

				if err := srv.Send(resp); err != nil {
					return err
				}
			} else if event.MarkerDelete != nil {
				// Send delete marker event to client
				resp := &pblivemapper.StreamResponse{
					Data: &pblivemapper.StreamResponse_Markers{
						Markers: &pblivemapper.MarkerMarkersUpdates{
							Deleted: []uint64{*event.MarkerDelete},
						},
					},
				}

				if err := srv.Send(resp); err != nil {
					return err
				}
			}

		case <-updateTicker.C:
			end, updatedAt, onDutyState, err = s.sendChunkedUserMarkers(srv, usersJobs, userInfo, updatedAt, onDutyState)
			if end || err != nil {
				return err
			}
		}
	}
}

// Sends out chunked current user markers
func (s *Server) sendChunkedUserMarkers(srv pblivemapper.LivemapperService_StreamServer, usersJobs *permissions.JobGradeList, userInfo *userinfo.UserInfo, updatedAt time.Time, lastOnDutyState bool) (bool, time.Time, bool, error) {
	// If the user was off duty and is now on duty, we need to send all user locations and not just updated ones
	// use the updatedAt time to force sending all users
	if onDutyState := s.tracker.IsUserOnDuty(userInfo.UserId); onDutyState != lastOnDutyState {
		if !lastOnDutyState {
			updatedAt = time.Time{}
		} else {
			clear := true
			resp := &pblivemapper.StreamResponse{
				Data: &pblivemapper.StreamResponse_Users{
					Users: &pblivemapper.UserMarkersUpdates{
						Clear: &clear,
					},
				},
				UserOnDuty: &onDutyState,
			}
			if err := srv.Send(resp); err != nil {
				return true, updatedAt, onDutyState, err
			}

			return false, updatedAt, onDutyState, nil
		}
	}

	updatedUsers, deletedUsers, onDutyState, lastUpdatedAt, err := s.getUserLocations(usersJobs, userInfo, updatedAt)
	if err != nil {
		return true, lastUpdatedAt, onDutyState, errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}

	// UpdatedAt is zero and no user updates or deletions? Early return
	if !updatedAt.IsZero() && len(updatedUsers) == 0 && len(deletedUsers) == 0 {
		return false, lastUpdatedAt, onDutyState, nil
	}

	// Less than chunk size or no markers, no need to chunk the response early return
	if len(updatedUsers) <= userMarkerChunkSize {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Users{
				Users: &pblivemapper.UserMarkersUpdates{
					Updated: updatedUsers,
					Deleted: deletedUsers,
					Part:    0,
					Partial: !updatedAt.IsZero(),
				},
			},
			UserOnDuty: &onDutyState,
		}

		if err := srv.Send(resp); err != nil {
			return true, lastUpdatedAt, onDutyState, err
		}

		return false, lastUpdatedAt, onDutyState, nil
	}

	totalParts := int32(len(updatedUsers) / userMarkerChunkSize)
	currentPart := totalParts
	for userMarkerChunkSize < len(updatedUsers) {
		userUpdates := &pblivemapper.UserMarkersUpdates{
			Updated: updatedUsers[0:userMarkerChunkSize:userMarkerChunkSize],
			Part:    currentPart,
			Partial: !updatedAt.IsZero(),
		}

		if totalParts == currentPart {
			userUpdates.Deleted = deletedUsers
		}

		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Users{
				Users: userUpdates,
			},
			UserOnDuty: &onDutyState,
		}
		currentPart--

		if err := srv.Send(resp); err != nil {
			return true, lastUpdatedAt, onDutyState, err
		}

		updatedUsers = updatedUsers[userMarkerChunkSize:]

		select {
		case <-srv.Context().Done():
			return true, lastUpdatedAt, onDutyState, nil

		case <-time.After(25 * time.Millisecond):
		}
	}

	if len(updatedUsers) > 0 {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Users{
				Users: &pblivemapper.UserMarkersUpdates{
					Updated: updatedUsers,
					Part:    0,
					Partial: !updatedAt.IsZero(),
				},
			},
			UserOnDuty: &onDutyState,
		}
		if err := srv.Send(resp); err != nil {
			return true, lastUpdatedAt, onDutyState, err
		}
	}

	return false, lastUpdatedAt, onDutyState, nil
}

// Send out chunked current marker markers
func (s *Server) sendMarkerMarkers(srv pblivemapper.LivemapperService_StreamServer, jobs *permissions.StringList, updatedAt time.Time) (bool, error) {
	updatedMarkers, deletedMarkers, err := s.getMarkerMarkers(jobs, updatedAt)
	if err != nil {
		return true, errswrap.NewError(err, errorslivemapper.ErrStreamFailed)
	}

	// UpdatedAt is zero and no user updates or deletions? Early return
	if !updatedAt.IsZero() && len(updatedMarkers) == 0 && len(deletedMarkers) == 0 {
		return false, nil
	}

	// Less than chunk size or no markers, no need to chunk the response early return
	if len(updatedMarkers) <= markerMarkerChunkSize {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Markers{
				Markers: &pblivemapper.MarkerMarkersUpdates{
					Updated: updatedMarkers,
					Deleted: deletedMarkers,
					Part:    0,
					Partial: !updatedAt.IsZero(),
				},
			},
		}

		if err := srv.Send(resp); err != nil {
			return true, err
		}

		return false, nil
	}

	totalParts := int32(len(updatedMarkers) / markerMarkerChunkSize)
	currentPart := totalParts
	for markerMarkerChunkSize < len(updatedMarkers) {
		markerUpdates := &pblivemapper.MarkerMarkersUpdates{
			Updated: updatedMarkers[0:markerMarkerChunkSize:markerMarkerChunkSize],
			Part:    currentPart,
			Partial: !updatedAt.IsZero(),
		}

		if totalParts == currentPart {
			markerUpdates.Deleted = deletedMarkers
		}

		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Markers{
				Markers: markerUpdates,
			},
		}
		currentPart--

		if err := srv.Send(resp); err != nil {
			return true, err
		}

		updatedMarkers = updatedMarkers[markerMarkerChunkSize:]

		select {
		case <-srv.Context().Done():
			return true, nil

		case <-time.After(25 * time.Millisecond):
		}
	}

	if len(updatedMarkers) > 0 {
		resp := &pblivemapper.StreamResponse{
			Data: &pblivemapper.StreamResponse_Markers{
				Markers: &pblivemapper.MarkerMarkersUpdates{
					Updated: updatedMarkers,
					Part:    0,
					Partial: !updatedAt.IsZero(),
				},
			},
		}
		if err := srv.Send(resp); err != nil {
			return true, err
		}
	}

	return false, nil
}

func (s *Server) getUserLocations(jobs *permissions.JobGradeList, userInfo *userinfo.UserInfo, updatedAt time.Time) ([]*livemap.UserMarker, []int32, bool, time.Time, error) {
	updated := []*livemap.UserMarker{}
	deleted := []int32{}

	found := false
	if userInfo.SuperUser {
		found = true
	}

	lastUpdatedAt := updatedAt

	for job, grades := range jobs.Iter() {
		if len(grades) == 0 {
			continue
		}

		markers, ok := s.tracker.GetUsersByJob(job)
		if !ok {
			continue
		}

		markers.Range(func(key int32, marker *livemap.UserMarker) bool {
			// If it isn't the user's own marker, user doesn't have access to grade, and is not a superuser, continue
			if key != userInfo.UserId && !jobs.HasJobGrade(job, marker.User.JobGrade) && !userInfo.SuperUser {
				return true
			}

			// Either no input updatedAt time set or the marker has been updated in the mean time
			if updatedAt.IsZero() || (marker.UpdatedAt != nil && updatedAt.Sub(marker.UpdatedAt.AsTime()) < 0) {
				if lastUpdatedAt.IsZero() || lastUpdatedAt.Sub(marker.UpdatedAt.AsTime()) < 0 {
					lastUpdatedAt = marker.UpdatedAt.AsTime()
				}

				if marker.Hidden {
					deleted = append(deleted, marker.UserId)
				} else {
					updated = append(updated, marker)
				}
			}

			// If the user is found in the list of user markers and not "off duty" (hidden), set found state
			if !found && !marker.Hidden && (userInfo.Job == job && userInfo.UserId == key) {
				found = true
			}

			return true
		})
	}

	if lastUpdatedAt.IsZero() {
		lastUpdatedAt = updatedAt
	}

	if found {
		return updated, deleted, true, lastUpdatedAt, nil
	}

	return nil, nil, false, lastUpdatedAt, nil
}
