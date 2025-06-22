package centrum

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

func (s *Server) sendHandshakre(ctx context.Context, srv pbcentrum.CentrumService_StreamServer, userJob string, aclJobs *pbcentrum.JobAccess) error {
	settings, err := s.settings.Get(ctx, userJob)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if err := srv.Send(&pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_Handshake{
			Handshake: &pbcentrum.StreamHandshake{
				ServerTime: timestamp.Now(),
				Settings:   settings,
				JobAccess:  aclJobs,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendLatestState(ctx context.Context, srv pbcentrum.CentrumService_StreamServer, userInfo *userinfo.UserInfo, aclJobs *pbcentrum.JobAccess, jobList []string) error {
	// Dispatchers
	dispatchers := &pbcentrum.Dispatchers{}
	for _, j := range aclJobs.Dispatches {
		dispos, err := s.dispatchers.Get(ctx, j.Job)
		if err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
		dispatchers.Dispatchers = append(dispatchers.Dispatchers, dispos)
	}

	// Own unit ID
	ownUnitMapping, _ := s.tracker.GetUserMapping(userInfo.UserId)
	var pOwnUnitId *uint64
	if ownUnitMapping != nil && ownUnitMapping.UnitId != nil && *ownUnitMapping.UnitId > 0 {
		pOwnUnitId = ownUnitMapping.UnitId
	}

	// Retrieve units and dispatches
	units := s.units.List(ctx, jobList)

	dispatches := s.dispatches.Filter(ctx, jobList, nil, []centrum.StatusDispatch{
		centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
		centrum.StatusDispatch_STATUS_DISPATCH_DELETED,
	})

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
	userInfo := auth.MustGetUserInfoFromContext(srv.Context()).Clone()

	// Check if user has access to other job's centrum
	jobList, jobAcls, err := s.settings.GetJobAccessList(srv.Context(), userInfo.Job, userInfo.JobGrade)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	for {
		if err := s.sendHandshakre(srv.Context(), srv, userInfo.Job, jobAcls); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.sendLatestState(srv.Context(), srv, userInfo, jobAcls, jobList); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.stream(srv.Context(), srv, userInfo, jobList); err != nil {
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
	StreamName string
	Bucket     string
	NoWildcard bool                                                                  // if true, the subjects are not a wildcard bucket (e.g., centrum_dispatchers.JOB)
	Unmarshal  func(ctx context.Context, s *Server, b []byte) (proto.Message, error) // bucket → concrete proto
	WrapPut    func(proto.Message) *pbcentrum.StreamResponse
	WrapDelete func(key string) *pbcentrum.StreamResponse
}

var feeds = []feedCfg{
	{
		StreamName: "centrum_settings",
		Bucket:     "centrum_settings",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var u centrum.Settings
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_Settings{Settings: m.(*centrum.Settings)}}
		},
	},
	{
		StreamName: "centrum_dispatchers",
		Bucket:     "centrum_dispatchers",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var u centrum.Dispatchers
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_Dispatchers{Dispatchers: m.(*centrum.Dispatchers)}}
		},
	},
	{
		StreamName: "centrum_units",
		Bucket:     "centrum_units.job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var d common.IDMapping
			if err := proto.Unmarshal(b, &d); err != nil {
				return nil, fmt.Errorf("failed to unmarshal unit id mapping. %w", err)
			}
			return s.units.Get(ctx, d.Id)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_UnitUpdated{UnitUpdated: m.(*centrum.Unit)}}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := centrumutils.ExtractID(key)
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
		StreamName: "centrum_dispatches",
		Bucket:     "centrum_dispatches.job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var d common.IDMapping
			if err := proto.Unmarshal(b, &d); err != nil {
				return nil, fmt.Errorf("failed to unmarshal dispatch id mapping. %w", err)
			}
			return s.dispatches.Get(ctx, d.Id)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_DispatchUpdated{DispatchUpdated: m.(*centrum.Dispatch)}}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := centrumutils.ExtractID(key)
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

func (s *Server) stream(ctx context.Context, srv pbcentrum.CentrumService_StreamServer, userInfo *userinfo.UserInfo, additionalJobs []string) error {
	s.logger.Debug("starting centrum stream", zap.String("job_main", userInfo.Job), zap.Int32("user_id", userInfo.UserId), zap.Strings("additional_jobs", additionalJobs))

	jobs := []string{userInfo.Job}
	jobs = append(jobs, additionalJobs...)
	jobs = utils.RemoveSliceDuplicates(jobs)

	out := make(chan *pbcentrum.StreamResponse, 256)
	g, gctx := errgroup.WithContext(ctx)

	// Centrum Events (e.g., dispatch and unit status updates)
	g.Go(func() error {
		// Create ephemeral consumer with multi-filter
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
				FilterSubjects: kvSubjects(f.Bucket, jobs, f.NoWildcard),
				DeliverPolicy:  jetstream.DeliverNewPolicy,
				AckPolicy:      jetstream.AckNonePolicy,
				MaxWaiting:     8,
			}
			consumer, err := s.js.CreateConsumer(gctx, "KV_"+f.StreamName, consCfg)
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

					obj, err := f.Unmarshal(gctx, s, m.Data())
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
