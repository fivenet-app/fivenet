package centrum

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

func (s *Server) sendHandshakre(ctx context.Context, srv pbcentrum.CentrumService_StreamServer, userJob string, jobs *pbcentrum.JobAccess) error {
	settings, err := s.state.GetSettings(ctx, userJob)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if err := srv.Send(&pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_Handshake{
			Handshake: &pbcentrum.StreamHandshake{
				ServerTime: timestamp.Now(),
				Settings:   settings,
				JobAccess:  jobs,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendLatestState(ctx context.Context, srv pbcentrum.CentrumService_StreamServer, userJob string, userId int32, jobs *pbcentrum.JobAccess) error {
	dispatchers := &pbcentrum.Dispatchers{}
	dispos, _ := s.state.GetDispatchers(ctx, userJob)
	dispatchers.Dispatchers = []*centrum.Dispatchers{dispos}
	for _, j := range jobs.Dispatches {
		dispos, err := s.state.GetDispatchers(ctx, j.Job.Name)
		if err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
		dispatchers.Dispatchers = append(dispatchers.Dispatchers, dispos)
	}

	ownUnitMapping, _ := s.tracker.GetUserMapping(userId)
	var pOwnUnitId *uint64
	if ownUnitMapping != nil && ownUnitMapping.UnitId != nil && *ownUnitMapping.UnitId > 0 {
		pOwnUnitId = ownUnitMapping.UnitId
	}

	// Retrieve units and dispatches
	units := s.state.ListUnits(ctx, userJob)
	for _, j := range jobs.Dispatches {
		units = append(units, s.state.ListUnits(ctx, j.Job.Name)...)
	}

	dispatchStatusFilter := []centrum.StatusDispatch{
		centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}
	dispatches := s.state.FilterDispatches(ctx, userJob, nil, dispatchStatusFilter)
	for _, j := range jobs.Dispatches {
		dispatches = append(dispatches, s.state.FilterDispatches(ctx, j.Job.Name, nil, dispatchStatusFilter)...)
	}

	// Send initial state to client
	if err := srv.Send(&pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_LatestState{
			LatestState: &pbcentrum.LatestState{
				Dispatchers: dispatchers,
				OwnUnitId:   pOwnUnitId,
				Units:       units,
				Dispatches:  dispatches,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stream(req *pbcentrum.StreamRequest, srv pbcentrum.CentrumService_StreamServer) error {
	userInfo := *auth.MustGetUserInfoFromContext(srv.Context())

	jobs := &pbcentrum.JobAccess{}
	additionalJobs := []string{}
	for _, job := range additionalJobs {
		j := s.enricher.GetJobByName(job)
		if j == nil {
			return errswrap.NewError(fmt.Errorf("job not found. %s", job), errorscentrum.ErrFailedQuery)
		}
		jobs.Dispatches = append(jobs.Dispatches, &pbcentrum.JobAccessEntry{
			Job:    j,
			Access: centrum.CentrumAccessLevel_ACCESS_LEVEL_DISPATCH,
		})
	}

	for {
		if err := s.sendHandshakre(srv.Context(), srv, userInfo.Job, jobs); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.sendLatestState(srv.Context(), srv, userInfo.Job, userInfo.UserId, jobs); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.stream(srv.Context(), srv, userInfo.Job, userInfo.UserId, []string{}); err != nil {
			return err
		}

		select {
		case <-srv.Context().Done():
			return nil

		case <-time.After(50 * time.Millisecond):
		}
	}
}

type feedCfg struct {
	Bucket     string
	Unmarshal  func([]byte) (proto.Message, error) // bucket → concrete proto
	WrapPut    func(proto.Message) *pbcentrum.StreamResponse
	WrapDelete func(key string) *pbcentrum.StreamResponse
}

var feeds = []feedCfg{
	{
		Bucket: "centrum_settings",
		Unmarshal: func(b []byte) (proto.Message, error) {
			var u centrum.Settings
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_Settings{Settings: m.(*centrum.Settings)}}
		},
	},
	{
		Bucket: "centrum_dispatchers",
		Unmarshal: func(b []byte) (proto.Message, error) {
			var u centrum.Dispatchers
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_Dispatchers{Dispatchers: m.(*centrum.Dispatchers)}}
		},
	},
	{
		Bucket: "centrum_units",
		Unmarshal: func(b []byte) (proto.Message, error) {
			var u centrum.Unit
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_UnitUpdated{UnitUpdated: m.(*centrum.Unit)}}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := extractID(key)
			if err != nil {
				return nil
			}

			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitDeleted{
					UnitDeleted: id,
				},
			}
		},
	},
	{
		Bucket: "centrum_dispatches",
		Unmarshal: func(b []byte) (proto.Message, error) {
			var d centrum.Dispatch
			return &d, proto.Unmarshal(b, &d)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_DispatchUpdated{DispatchUpdated: m.(*centrum.Dispatch)}}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := extractID(key)
			if err != nil {
				return nil
			}

			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchDeleted{
					DispatchDeleted: id,
				},
			}
		},
	},
}

func (s *Server) stream(ctx context.Context, srv pbcentrum.CentrumService_StreamServer, job string, userId int32, additionalJobs []string) error {
	s.logger.Debug("starting centrum stream", zap.String("job_main", job), zap.Int32("user_id", userId), zap.Strings("additional_jobs", additionalJobs))

	jobs := []string{job}
	jobs = append(jobs, additionalJobs...)
	jobs = utils.RemoveSliceDuplicates(jobs)

	out := make(chan *pbcentrum.StreamResponse, 256)
	g, gctx := errgroup.WithContext(ctx)

	// Centrum Events (e.g., dispatch and unit status updates)
	g.Go(func() error {
		// Create consumer with multi-filter
		consCfg := jetstream.ConsumerConfig{
			FilterSubjects: centrumSubjects(jobs),
			DeliverPolicy:  jetstream.DeliverNewPolicy,
			AckPolicy:      jetstream.AckNonePolicy,
		}
		consumer, err := s.js.CreateConsumer(gctx, "CENTRUM", consCfg)
		if err != nil {
			return fmt.Errorf("failed to create consumer. %w", err)
		}

		// Pull loop
		for {
			batch, err := consumer.Fetch(32,
				jetstream.FetchMaxWait(2*time.Second))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, jetstream.ErrNoMessages) {
					continue // idle
				}
				return err
			}

			for m := range batch.Messages() {
				_, topic, tType := eventscentrum.SplitSubject(m.Subject())

				var r *pbcentrum.StreamResponse

				switch topic {
				case eventscentrum.TopicDispatch:
					if tType != eventscentrum.TypeDispatchStatus {
						continue
					}

					var d centrum.DispatchStatus
					if err := proto.Unmarshal(m.Data(), &d); err != nil {
						s.logger.Error("failed to unmarshal dispatch status", zap.Error(err), zap.String("subject", m.Subject()))
					}

					r = &pbcentrum.StreamResponse{
						Change: &pbcentrum.StreamResponse_DispatchStatus{
							DispatchStatus: &d,
						},
					}

				case eventscentrum.TopicUnit:
					if tType != eventscentrum.TypeUnitStatus {
						continue
					}
					var u centrum.UnitStatus
					if err := proto.Unmarshal(m.Data(), &u); err != nil {
						s.logger.Error("failed to unmarshal unit status", zap.Error(err), zap.String("subject", m.Subject()))
					}

					r = &pbcentrum.StreamResponse{
						Change: &pbcentrum.StreamResponse_UnitStatus{
							UnitStatus: &u,
						},
					}
				}

				if r == nil {
					s.logger.Warn("received unknown centrum event", zap.String("subject", m.Subject()), zap.String("type", string(tType)))
					continue
				}

				select {
				case out <- r:

				case <-gctx.Done():
					return gctx.Err()
				}
			}
		}
	})

	// Setup feeds for each bucket
	for _, f := range feeds {
		g.Go(func() error {
			// Create consumer with multi-filter
			consCfg := jetstream.ConsumerConfig{
				FilterSubjects: kvSubjects(f.Bucket, jobs),
				DeliverPolicy:  jetstream.DeliverNewPolicy,
				AckPolicy:      jetstream.AckNonePolicy,
			}
			consumer, err := s.js.CreateConsumer(gctx, "KV_"+f.Bucket, consCfg)
			if err != nil {
				return fmt.Errorf("failed to create consumer. %w", err)
			}

			// Pull loop
			for {
				batch, err := consumer.Fetch(32,
					jetstream.FetchMaxWait(2*time.Second))
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) ||
						errors.Is(err, jetstream.ErrNoMessages) {
						continue // idle
					}
					return err
				}

				for m := range batch.Messages() {
					if op := m.Headers().Get("KV-Operation"); op == "DEL" || op == "PURGE" {
						key := strings.TrimPrefix(m.Subject(), "$KV."+f.Bucket+".")

						r := f.WrapDelete(key)
						if r == nil {
							continue
						}
						select {
						case out <- r:

						case <-gctx.Done():
							return gctx.Err()
						}
						continue
					}

					obj, err := f.Unmarshal(m.Data())
					if err != nil {
						// Bad payload – skip
						continue
					}

					r := f.WrapPut(obj)
					if r == nil {
						continue
					}
					select {
					case out <- r:

					case <-gctx.Done():
						return gctx.Err()
					}
				}
			}
		})
	}

	// Single writer
	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return gctx.Err()

			case resp := <-out:
				if err := srv.Send(resp); err != nil {
					return err
				}
			}
		}
	})

	return g.Wait()
}
