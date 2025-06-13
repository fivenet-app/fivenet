package livemap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	permslivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/grpcws"
	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	errorslivemap "github.com/fivenet-app/fivenet/v2025/services/livemap/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/klauspost/compress/zstd"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

const (
	markerMarkerChunkSize = 75

	feedFetch = 128
	feedWait  = 2 * time.Second
)

func (s *Server) getAndSendACL(srv pblivemap.LivemapService_StreamServer, userInfo *userinfo.UserInfo) (*permissions.StringList, *permissions.JobGradeList, error) {
	markerJobs, err := s.ps.AttrJobList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceStreamPerm, permslivemap.LivemapServiceStreamMarkersPermField)
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
	}
	usersJobs, err := s.ps.AttrJobGradeList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceStreamPerm, permslivemap.LivemapServiceStreamPlayersPermField)
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorslivemap.ErrStreamFailed)
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
		return nil, nil, err
	}

	return markerJobs, usersJobs, nil
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

func (s *Server) Stream(req *pblivemap.StreamRequest, srv pblivemap.LivemapService_StreamServer) error {
	ctx := srv.Context()

	origUI := auth.MustGetUserInfoFromContext(ctx).Clone()
	userInfo := &origUI

	s.logger.Debug("starting livemap stream", zap.Int32("user_id", userInfo.UserId))

	markerJobs, usersJobs, err := s.getAndSendACL(srv, &origUI)
	if err != nil {
		return err
	}

	if end, err := s.sendMarkerMarkers(srv, markerJobs, time.Time{}); end || err != nil {
		return err
	}

	// Fetch latest snapshot (may be nil on a cold cluster)
	markers, snapTS, err := s.fetchSnapshot(ctx)
	if err != nil {
		return err
	}

	// Apply ACL from the request
	markers = tracker.FilterMarkers(markers, usersJobs, userInfo)

	// Send initial payload
	if err := srv.Send(&pblivemap.StreamResponse{
		Data: &pblivemap.StreamResponse_Snapshot{
			Snapshot: &pblivemap.Snapshot{
				Markers:     markers,
				GeneratedAt: snapTS,
			},
		},
	}); err != nil {
		return err
	}

	markerUpdateCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(markerUpdateCh)

	// Central pipe: all feeds push messages into outCh
	outCh := make(chan *pblivemap.StreamResponse, 1024)
	g, gctx := errgroup.WithContext(srv.Context())

	g.Go(func() error {
		for {
			select {
			case <-srv.Context().Done():
				return nil

			case event := <-markerUpdateCh:
				if event == nil {
					continue
				}

				if event.MarkerUpdate != nil {
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

	// Writer goroutine – single gRPC send loop
	g.Go(func() error {
		defer close(outCh)
		for {
			select {
			case <-gctx.Done():
				return gctx.Err()
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

	g.Go(func() error {
		meta := metadata.ExtractIncoming(ctx)
		connID := meta.Get(grpcws.ConnectionIdHeader)

		// Create / bind consumer
		cc := jetstream.ConsumerConfig{
			DeliverPolicy:  jetstream.DeliverNewPolicy,
			AckPolicy:      jetstream.AckNonePolicy,
			Durable:        "lm_userloc-deltas" + connID,
			FilterSubjects: buildFilters(usersJobs),
		}

		stream, err := s.js.Stream(ctx, "KV_"+tracker.BucketUserLoc)
		if err != nil {
			return err
		}

		cons, err := stream.CreateConsumer(ctx, cc)
		if err != nil {
			return err
		}

		for {
			msgs, err := cons.Fetch(feedFetch, jetstream.FetchMaxWait(feedWait))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, jetstream.ErrNoMessages) {
					continue // keep polling
				}
				return err
			}
			for m := range msgs.Messages() {
				if op := m.Headers().Get("KV-Operation"); op == "DEL" || op == "PURGE" {
					key := strings.TrimPrefix(m.Subject(), "$KV."+tracker.BucketUserLoc+".")

					userId, err := tracker.ExtractUserID(key)
					if err != nil {
						return errswrap.NewError(err, errorslivemap.ErrStreamFailed)
					}

					select {
					case outCh <- &pblivemap.StreamResponse{
						Data: &pblivemap.StreamResponse_UserDelete{
							UserDelete: userId,
						},
					}:

					case <-ctx.Done():
						return nil
					}
					continue
				}

				um := &livemap.UserMarker{}
				if err := proto.Unmarshal(m.Data(), um); err == nil {
					jg := um.User.JobGrade
					if um.JobGrade != nil {
						jg = *um.JobGrade
					}

					if !userInfo.Superuser && !usersJobs.HasJobGrade(um.Job, jg) {
						continue
					}

					select {
					case outCh <- &pblivemap.StreamResponse{
						Data: &pblivemap.StreamResponse_UserUpdate{
							UserUpdate: um,
						},
					}:

					case <-ctx.Done():
						return nil
					}
				}
			}
		}
	})

	return g.Wait()
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

func (s *Server) fetchSnapshot(ctx context.Context) ([]*livemap.UserMarker, int64, error) {
	stream, err := s.js.Stream(ctx, "KV_"+tracker.BucketUserLoc)
	if err != nil {
		return nil, 0, err
	}
	cons, err := stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		DeliverPolicy:  jetstream.DeliverLastPolicy, // Newest snapshot only
		AckPolicy:      jetstream.AckNonePolicy,
		FilterSubjects: []string{tracker.SnapshotSubject},
	})
	if err != nil {
		return nil, 0, err
	}

	msgs, err := cons.Fetch(1, jetstream.FetchMaxWait(2*time.Second))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, jetstream.ErrNoMessages) {
			return nil, 0, nil // no snapshot yet
		}
		return nil, 0, err
	}
	var msg jetstream.Msg
	for m := range msgs.Messages() {
		msg = m
	}

	if msg == nil {
		return nil, 0, nil
	}

	zr, err := zstd.NewReader(bytes.NewReader(msg.Data()))
	if err != nil {
		return nil, 0, err
	}
	buf, err := io.ReadAll(zr)
	if err != nil {
		return nil, 0, err
	}
	zr.Close()

	snap := &pblivemap.Snapshot{}
	if err := proto.Unmarshal(buf, snap); err != nil {
		return nil, 0, err
	}

	return snap.Markers, snap.GeneratedAt, nil
}
