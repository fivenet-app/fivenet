package livemap

import (
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	permslivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	errorslivemap "github.com/fivenet-app/fivenet/v2025/services/livemap/errors"
	"go.uber.org/zap"
)

const (
	userMarkerChunkSize   = 30
	markerMarkerChunkSize = 75
)

func (s *Server) Stream(req *pblivemap.StreamRequest, srv pblivemap.LivemapService_StreamServer) error {
	origUI := auth.MustGetUserInfoFromContext(srv.Context()).Clone()
	userInfo := &origUI

	s.logger.Debug("starting livemap stream", zap.Int32("user_id", userInfo.UserId))
	markerJobs, err := s.ps.AttrJobList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceStreamPerm, permslivemap.LivemapServiceStreamMarkersPermField)
	if err != nil {
		return errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}
	usersJobs, err := s.ps.AttrJobGradeList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceStreamPerm, permslivemap.LivemapServiceStreamPlayersPermField)
	if err != nil {
		return errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}

	if userInfo.Superuser {
		s.markersCache.Range(func(job string, _ []*livemap.MarkerMarker) bool {
			markerJobs.Strings = append(markerJobs.Strings, job)
			return true
		})
		markerJobs.Strings = utils.RemoveSliceDuplicates(markerJobs.Strings)

		if usersJobs.Jobs == nil {
			usersJobs.Jobs = make(map[string]int32)
		}
		// Disable fine-grained permissions for superuser as the it is now just a list of jobs
		usersJobs.FineGrained = false
		for _, j := range s.tracker.ListTrackedJobs() {
			usersJobs.Jobs[j] = -1
		}
	}

	// Prepare js for client response
	js := &pblivemap.StreamResponse_Jobs{
		Jobs: &pblivemap.JobsList{
			Markers: make([]*jobs.Job, markerJobs.Len()),
			Users:   []*jobs.Job{},
		},
	}

	for i := range markerJobs.Strings {
		js.Jobs.Markers[i] = &jobs.Job{
			Name: markerJobs.Strings[i],
		}
		s.enricher.EnrichJobName(js.Jobs.Markers[i])
	}
	for job := range usersJobs.Iter() {
		j := &jobs.Job{
			Name: job,
		}
		s.enricher.EnrichJobName(j)
		js.Jobs.Users = append(js.Jobs.Users, j)
	}

	if err := srv.Send(&pblivemap.StreamResponse{
		Data: js,
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

			if event.UserRemoved != nil {
				if usersJobs.Len() == 0 {
					continue
				}

				// Send delete user marker event to client
				if event.UserRemoved.Hidden || (userInfo.Superuser || !usersJobs.HasJobGrade(event.UserRemoved.Job, event.UserRemoved.User.JobGrade)) {
					continue
				}

				resp := &pblivemap.StreamResponse{
					Data: &pblivemap.StreamResponse_Users{
						Users: &pblivemap.UserMarkersUpdates{
							Deleted: []int32{event.UserRemoved.UserId},
							Partial: true,
						},
					},
				}
				if err := srv.Send(resp); err != nil {
					return err
				}
			} else if event.MarkerUpdate != nil {
				// Send delete marker event to client
				resp := &pblivemap.StreamResponse{
					Data: &pblivemap.StreamResponse_Markers{
						Markers: &pblivemap.MarkerMarkersUpdates{
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
				resp := &pblivemap.StreamResponse{
					Data: &pblivemap.StreamResponse_Markers{
						Markers: &pblivemap.MarkerMarkersUpdates{
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
func (s *Server) sendChunkedUserMarkers(srv pblivemap.LivemapService_StreamServer, usersJobs *permissions.JobGradeList, userInfo *userinfo.UserInfo, updatedAt time.Time, lastOnDutyState bool) (bool, time.Time, bool, error) {
	// If the user was off duty and is now on duty, we need to send all user locations and not just updated ones
	// use the updatedAt time to force sending all users
	if onDutyState := s.tracker.IsUserOnDuty(userInfo.UserId); !userInfo.Superuser && onDutyState != lastOnDutyState {
		if !lastOnDutyState {
			updatedAt = time.Time{}
		} else {
			clear := true
			resp := &pblivemap.StreamResponse{
				Data: &pblivemap.StreamResponse_Users{
					Users: &pblivemap.UserMarkersUpdates{
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
		return true, lastUpdatedAt, onDutyState, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}

	// UpdatedAt is zero and no user updates or deletions? Early return
	if !updatedAt.IsZero() && len(updatedUsers) == 0 && len(deletedUsers) == 0 {
		return false, lastUpdatedAt, onDutyState, nil
	}

	// Less than chunk size or no markers, no need to chunk the response early return
	if len(updatedUsers) <= userMarkerChunkSize {
		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Users{
				Users: &pblivemap.UserMarkersUpdates{
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
		userUpdates := &pblivemap.UserMarkersUpdates{
			Updated: updatedUsers[0:userMarkerChunkSize:userMarkerChunkSize],
			Part:    currentPart,
			Partial: !updatedAt.IsZero(),
		}

		if totalParts == currentPart {
			userUpdates.Deleted = deletedUsers
		}

		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Users{
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
		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Users{
				Users: &pblivemap.UserMarkersUpdates{
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
func (s *Server) sendMarkerMarkers(srv pblivemap.LivemapService_StreamServer, jobs *permissions.StringList, updatedAt time.Time) (bool, error) {
	updatedMarkers, deletedMarkers, err := s.getMarkerMarkers(jobs, updatedAt)
	if err != nil {
		return true, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}

	// UpdatedAt is zero and no user updates or deletions? Early return
	if !updatedAt.IsZero() && len(updatedMarkers) == 0 && len(deletedMarkers) == 0 {
		return false, nil
	}

	// Less than chunk size or no markers, no need to chunk the response early return
	if len(updatedMarkers) <= markerMarkerChunkSize {
		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Markers{
				Markers: &pblivemap.MarkerMarkersUpdates{
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
		markerUpdates := &pblivemap.MarkerMarkersUpdates{
			Updated: updatedMarkers[0:markerMarkerChunkSize:markerMarkerChunkSize],
			Part:    currentPart,
			Partial: !updatedAt.IsZero(),
		}

		if totalParts == currentPart {
			markerUpdates.Deleted = deletedMarkers
		}

		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Markers{
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
		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Markers{
				Markers: &pblivemap.MarkerMarkersUpdates{
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
	if userInfo.Superuser {
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
			if key != userInfo.UserId && !jobs.HasJobGrade(job, marker.User.JobGrade) && !userInfo.Superuser {
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
