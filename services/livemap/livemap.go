package livemap

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	permslivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	errorslivemap "github.com/fivenet-app/fivenet/v2025/services/livemap/errors"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

const (
	markerMarkerChunkSize = 75

	feedFetch = 32
)

func (s *Server) getAndSendACL(srv pblivemap.LivemapService_StreamServer, userInfo *userinfo.UserInfo) (*permissions.StringList, *permissions.JobGradeList, bool, error) {
	markerJobs, err := s.ps.AttrJobList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceStreamPerm, permslivemap.LivemapServiceStreamMarkersPermField)
	if err != nil {
		return nil, nil, false, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}
	usersJobs, err := s.ps.AttrJobGradeList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceStreamPerm, permslivemap.LivemapServiceStreamPlayersPermField)
	if err != nil {
		return nil, nil, false, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
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
			Markers: []*jobs.Job{},
			Users:   []*jobs.Job{},
		},
	}

	for i := range markerJobs.Strings {
		jm := &jobs.Job{
			Name: markerJobs.Strings[i],
		}
		s.enricher.EnrichJobName(jm)
		js.Jobs.Markers = append(js.Jobs.Markers, jm)
	}
	for job := range usersJobs.Iter() {
		j := &jobs.Job{
			Name: job,
		}
		s.enricher.EnrichJobName(j)
		js.Jobs.Users = append(js.Jobs.Users, j)
	}

	// Check if the user is on duty (superuser is always on duty)
	userOnDuty := false
	if userInfo.Superuser {
		userOnDuty = true
	} else if um, ok := s.tracker.GetUserMarkerById(userInfo.UserId); ok && !um.Hidden {
		userOnDuty = true
	}

	if err := srv.Send(&pblivemap.StreamResponse{
		UserOnDuty: &userOnDuty,
		Data:       js,
	}); err != nil {
		return nil, nil, false, err
	}

	return markerJobs, usersJobs, userOnDuty, nil
}

// buildFilters returns the FilterSubjects slice that encodes the caller’s ACL.
func buildFilters(jobs *permissions.JobGradeList) []string {
	var f []string
	for job, grades := range jobs.Iter() {
		if jobs.FineGrained {
			for _, g := range grades {
				f = append(f, fmt.Sprintf("$KV.%s.%s.%d", tracker.BucketUserLoc, job, g))
			}
		} else {
			// Non-fine-grained is expressed as a wildcard; higher grades are excluded in the stream's filter
			f = append(f, fmt.Sprintf("$KV.%s.%s.>", tracker.BucketUserLoc, job))
		}
	}
	return f
}

func (s *Server) sendUserMarkers(srv pblivemap.LivemapService_StreamServer, usersJobs *permissions.JobGradeList, userInfo *userinfo.UserInfo) error {
	// Get user markers
	markers := s.tracker.GetFilteredUserMarkers(usersJobs, userInfo)

	// Send initial payload
	if err := srv.Send(&pblivemap.StreamResponse{
		Data: &pblivemap.StreamResponse_Snapshot{
			Snapshot: &pblivemap.Snapshot{
				Markers: markers,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stream(req *pblivemap.StreamRequest, srv pblivemap.LivemapService_StreamServer) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx).Clone()

	s.logger.Debug("starting livemap stream", zap.Int32("user_id", userInfo.UserId))

	markerJobs, usersJobs, userOnDuty, err := s.getAndSendACL(srv, userInfo)
	if err != nil {
		return err
	}

	if end, err := s.sendMarkerMarkers(srv, markerJobs); end || err != nil {
		return err
	}

	if userOnDuty {
		if err := s.sendUserMarkers(srv, usersJobs, userInfo); err != nil {
			return errswrap.NewError(err, errorslivemap.ErrStreamFailed)
		}
	} else {
		// Send empty snapshot if the user is not on duty
		if err := srv.Send(&pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Snapshot{
				Snapshot: &pblivemap.Snapshot{},
			},
		}); err != nil {
			return err
		}
	}

	// Central pipe: all feeds push messages into outCh
	outCh := make(chan *pblivemap.StreamResponse, 256)
	defer close(outCh)
	g, gctx := errgroup.WithContext(ctx)

	// Writer goroutine – single gRPC send loop
	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil

			case msg := <-outCh:
				if msg == nil {
					return nil
				}

				if err := srv.Send(msg); err != nil {
					return err
				}
			}
		}
	})

	// Marker updates goroutine – listens for marker updates and sends them to outCh
	g.Go(func() error {
		markerUpdateCh := s.broker.Subscribe()
		defer s.broker.Unsubscribe(markerUpdateCh)

		for {
			select {
			case <-gctx.Done():
				return nil

			case event := <-markerUpdateCh:
				if event == nil {
					continue
				}

				if event.MarkerUpdate != nil {
					if event.MarkerUpdate.Job != userInfo.Job && !userInfo.Superuser {
						continue // Ignore updates for other jobs
					}

					// Send delete marker event to client
					outCh <- &pblivemap.StreamResponse{
						Data: &pblivemap.StreamResponse_Markers{
							Markers: &pblivemap.MarkerMarkersUpdates{
								Updated: []*livemap.MarkerMarker{event.MarkerUpdate},
								Partial: true,
							},
						},
					}
				} else if event.MarkerDelete != nil {
					// Send delete marker event to client
					outCh <- &pblivemap.StreamResponse{
						Data: &pblivemap.StreamResponse_Markers{
							Markers: &pblivemap.MarkerMarkersUpdates{
								Deleted: []uint64{*event.MarkerDelete},
							},
						},
					}
				} else {
					s.logger.Warn("received unknown event type in livemap stream", zap.Any("event", event))
				}
			}
		}
	})

	// User markers goroutine – listens for user marker updates and sends them to outCh
	g.Go(func() error {
		// Upsert pull consumer with multi-filter
		consCfg := jetstream.ConsumerConfig{
			FilterSubjects: buildFilters(usersJobs),
			DeliverPolicy:  jetstream.DeliverNewPolicy,
			AckPolicy:      jetstream.AckNonePolicy,
			MaxWaiting:     8,
		}
		consumer, err := s.js.CreateConsumer(gctx, "KV_"+tracker.BucketUserLoc, consCfg)
		if err != nil {
			return fmt.Errorf("failed to create consumer. %w", err)
		}

		for {
			batch, err := consumer.Fetch(feedFetch,
				jetstream.FetchMaxWait(2*time.Second))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, jetstream.ErrNoMessages) {
					continue // keep polling
				}
				return err
			}

			for m := range batch.Messages() {
				op := m.Headers().Get("KV-Operation")
				m.Ack()

				if op == "DEL" || op == "PURGE" {
					// Ignore delete and purge operations when not on duty
					if !userOnDuty {
						continue
					}

					key := strings.TrimPrefix(m.Subject(), "$KV."+tracker.BucketUserLoc+".")
					userId, job, jobGrade, err := tracker.DecodeUserMarkerKey(key)
					if err != nil {
						return errswrap.NewError(err, errorslivemap.ErrStreamFailed)
					}

					if userId == userInfo.UserId && job == userInfo.Job && jobGrade == userInfo.JobGrade {
						userOnDuty = false
					}

					select {
					case <-gctx.Done():
						return nil

					case outCh <- &pblivemap.StreamResponse{
						UserOnDuty: &userOnDuty,
						Data: &pblivemap.StreamResponse_UserDelete{
							UserDelete: &pblivemap.UserDelete{
								Id:  userId,
								Job: job,
							},
						},
					}:
					}
					continue
				}

				um := &livemap.UserMarker{}
				if err := proto.Unmarshal(m.Data(), um); err != nil {
					continue
				}

				// Marker is hidden, send delete event
				if um.Hidden {
					// If the user is hidden, we toggle the on duty state and "drop" any message not related to the user
					if um.UserId == userInfo.UserId && um.Job == userInfo.Job && (um.JobGrade == nil || *um.JobGrade == userInfo.JobGrade) {
						userOnDuty = false
					}

					select {
					case <-gctx.Done():
						return nil

					case outCh <- &pblivemap.StreamResponse{
						UserOnDuty: &userOnDuty,
						Data: &pblivemap.StreamResponse_UserDelete{
							UserDelete: &pblivemap.UserDelete{
								Id:  um.UserId,
								Job: um.Job,
							},
						},
					}:
					}
					continue
				}

				if !userOnDuty {
					if um.UserId == userInfo.UserId {
						userOnDuty = true
						// If the user is (back) on duty, we send the user markers snapshot
						if err := s.sendUserMarkers(srv, usersJobs, userInfo); err != nil {
							return errswrap.NewError(err, errorslivemap.ErrStreamFailed)
						}
					} else {
						// If the user is not on duty, we skip sending marker updates
						continue
					}
				}

				job := um.Job
				if um.Job == "" {
					job = um.User.Job
				}
				jg := um.User.JobGrade
				if um.JobGrade != nil {
					jg = *um.JobGrade
				}

				if !userInfo.Superuser && !usersJobs.HasJobGrade(job, jg) {
					continue
				}

				select {
				case <-gctx.Done():
					return nil

				case outCh <- &pblivemap.StreamResponse{
					UserOnDuty: &userOnDuty,
					Data: &pblivemap.StreamResponse_UserUpdate{
						UserUpdate: um,
					},
				}:
				}

			}
		}
	})

	return g.Wait()
}

// Send out chunked current marker markers
func (s *Server) sendMarkerMarkers(srv pblivemap.LivemapService_StreamServer, jobs *permissions.StringList) (bool, error) {
	updatedMarkers, deletedMarkers, err := s.getMarkerMarkers(jobs)
	if err != nil {
		return true, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}

	// Less than chunk size or no markers, no need to chunk the response early return
	if len(updatedMarkers) <= markerMarkerChunkSize {
		resp := &pblivemap.StreamResponse{
			Data: &pblivemap.StreamResponse_Markers{
				Markers: &pblivemap.MarkerMarkersUpdates{
					Updated: updatedMarkers,
					Deleted: deletedMarkers,
					Part:    0,
					Partial: false,
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
			Partial: false,
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
					Partial: false,
				},
			},
		}
		if err := srv.Send(resp); err != nil {
			return true, err
		}
	}

	return false, nil
}
