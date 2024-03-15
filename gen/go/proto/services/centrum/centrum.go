package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"sync"
	"time"

	centrum "github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/config/appconfig"
	"github.com/galexrt/fivenet/pkg/coords/postals"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/utils/broker"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	CentrumServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup
	jsCons jetstream.ConsumeContext

	tracer  trace.Tracer
	db      *sql.DB
	ps      perms.Permissions
	aud     audit.IAuditer
	js      jetstream.JetStream
	tracker tracker.ITracker
	postals postals.Postals
	appCfg  appconfig.IConfig

	brokersMutex sync.RWMutex
	brokers      map[string]*broker.Broker[*StreamResponse]

	state *manager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Perms     perms.Permissions
	Audit     audit.IAuditer
	JS        jetstream.JetStream
	Config    *config.Config
	AppConfig appconfig.IConfig
	Tracker   tracker.ITracker
	Postals   postals.Postals
	Manager   *manager.Manager
}

func NewServer(p Params) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	brokers := map[string]*broker.Broker[*StreamResponse]{}

	s := &Server{
		logger: p.Logger.Named("centrum"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		db:      p.DB,
		ps:      p.Perms,
		aud:     p.Audit,
		js:      p.JS,
		tracker: p.Tracker,
		postals: p.Postals,
		appCfg:  p.AppConfig,

		brokersMutex: sync.RWMutex{},
		brokers:      brokers,

		state: p.Manager,
	}

	p.LC.Append(fx.StartHook(func(c context.Context) error {
		s.handleAppConfigUpdate(ctx, p.AppConfig.Get())

		if err := s.registerSubscriptions(c); err != nil {
			return fmt.Errorf("failed to subscribe to events: %w", err)
		}

		// Handle app config updates
		go func() {
			configUpdateCh := p.AppConfig.Subscribe()
			for {
				select {
				case <-ctx.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					s.handleAppConfigUpdate(ctx, cfg)
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		s.jsCons.Stop()

		return nil
	}))

	return s, nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCentrumServiceServer(srv, s)
}

func (s *Server) registerSubscriptions(ctx context.Context) error {
	streamCfg, err := eventscentrum.RegisterStream(ctx, s.js)
	if err != nil {
		return fmt.Errorf("failed to register events: %w", err)
	}

	consumer, err := s.js.CreateConsumer(ctx, streamCfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		FilterSubject: fmt.Sprintf("%s.>", eventscentrum.BaseSubject),
	})
	if err != nil {
		return err
	}

	s.jsCons, err = consumer.Consume(s.watchForChanges)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleAppConfigUpdate(ctx context.Context, cfg *appconfig.Cfg) {
	s.brokersMutex.Lock()
	defer s.brokersMutex.Unlock()

	for _, job := range cfg.UserTracker.LivemapJobs {
		if _, ok := s.brokers[job]; !ok {
			s.brokers[job] = broker.New[*StreamResponse]()
			go s.brokers[job].Start(ctx)
		}
	}
}

func (s *Server) watchForChanges(msg jetstream.Msg) {
	if err := msg.Ack(); err != nil {
		s.logger.Error("failed to ack message", zap.Error(err))
	}

	job, topic, tType := eventscentrum.SplitSubject(msg.Subject())

	broker, ok := s.getJobBroker(job)
	if !ok {
		s.logger.Debug("no broker found for job", zap.String("job", job))
		return
	}

	resp := &StreamResponse{}
	switch topic {
	case eventscentrum.TopicGeneral:
		switch tType {
		case eventscentrum.TypeGeneralDisponents:
			dest := &centrum.Disponents{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_Disponents{
				Disponents: dest,
			}

		case eventscentrum.TypeGeneralSettings:
			dest := &centrum.Settings{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_Settings{
				Settings: dest,
			}
		}

	case eventscentrum.TopicUnit:
		switch tType {
		case eventscentrum.TypeUnitCreated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitCreated{
				UnitCreated: dest,
			}

		case eventscentrum.TypeUnitDeleted:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitDeleted{
				UnitDeleted: dest,
			}

		case eventscentrum.TypeUnitUpdated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitUpdated{
				UnitUpdated: dest,
			}

		case eventscentrum.TypeUnitStatus:
			dest := &centrum.UnitStatus{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitStatus{
				UnitStatus: dest,
			}
		}

	case eventscentrum.TopicDispatch:
		switch tType {
		case eventscentrum.TypeDispatchCreated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchCreated{
				DispatchCreated: dest,
			}

		case eventscentrum.TypeDispatchDeleted:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchDeleted{
				DispatchDeleted: dest,
			}

		case eventscentrum.TypeDispatchUpdated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchUpdated{
				DispatchUpdated: dest,
			}

		case eventscentrum.TypeDispatchStatus:
			dest := &centrum.DispatchStatus{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchStatus{
				DispatchStatus: dest,
			}
		}
	}

	broker.Publish(resp)

	meta, err := msg.Metadata()
	if err != nil {
		s.logger.Error("sent centrum message broker, but failed to get msg metadata ", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
			zap.String("job", job), zap.String("topic", string(topic)), zap.String("type", string(tType)), zap.Error(err))
		return
	}
	s.logger.Debug("sent centrum message broker", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
		zap.String("job", job), zap.String("topic", string(topic)), zap.String("type", string(tType)))
}

func (s *Server) TakeControl(ctx context.Context, req *TakeControlRequest) (*TakeControlResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeControl",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.state.DisponentSignOn(ctx, userInfo.Job, userInfo.UserId, req.Signon); err != nil {
		return nil, err
	}

	if req.Signon {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &TakeControlResponse{}, nil
}

func (s *Server) sendLatestState(srv CentrumService_StreamServer, job string, userId int32) (uint64, bool, error) {
	ctx := srv.Context()

	settings := s.state.GetSettings(ctx, job)
	disponents, _ := s.state.GetDisponents(ctx, job)
	isDisponent := s.state.CheckIfUserIsDisponent(ctx, job, userId)
	ownUnitId, _ := s.state.GetUserUnitID(ctx, userId)
	units, _ := s.state.ListUnits(ctx, job)
	ownUnit, _ := s.state.GetUnit(ctx, job, ownUnitId)

	dispatches := s.state.FilterDispatches(ctx, job, nil, []centrum.StatusDispatch{
		centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	})

	// Send initial state message to client
	resp := &StreamResponse{
		Change: &StreamResponse_LatestState{
			LatestState: &LatestState{
				Settings:    settings,
				Disponents:  disponents,
				IsDisponent: isDisponent,
				OwnUnit:     ownUnit,
				Units:       units,
				Dispatches:  dispatches,
				ServerTime:  timestamp.Now(),
			},
		},
	}
	if err := srv.Send(resp); err != nil {
		return 0, isDisponent, err
	}

	return ownUnitId, isDisponent, nil
}

func (s *Server) Stream(req *StreamRequest, srv CentrumService_StreamServer) error {
	userInfo := *auth.MustGetUserInfoFromContext(srv.Context())

	for {
		unitId, isDisponent, err := s.sendLatestState(srv, userInfo.Job, userInfo.UserId)
		if err != nil {
			return err
		}

		end, err := s.stream(srv, isDisponent, userInfo.Job, userInfo.UserId, unitId)
		if end {
			return err
		} else if err != nil {
			s.logger.Error("error during stream", zap.Error(err))
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func (s *Server) getJobBroker(job string) (*broker.Broker[*StreamResponse], bool) {
	s.brokersMutex.RLock()
	defer s.brokersMutex.RUnlock()

	broker, ok := s.brokers[job]
	return broker, ok
}

func (s *Server) stream(srv CentrumService_StreamServer, isDisponent bool, job string, userId int32, unitId uint64) (bool, error) {
	broker, ok := s.getJobBroker(job)
	if !ok {
		s.logger.Warn("no job broker found", zap.String("job", job), zap.Int32("user_id", userId))
		<-srv.Context().Done()
		return true, nil
	}

	stream := broker.Subscribe()
	defer broker.Unsubscribe(stream)

	// Watch for events from message queue
	for {
		resp := &StreamResponse{}

		select {
		case <-srv.Context().Done():
			return true, nil

		case msg := <-stream:
			resp.Change = msg.GetChange()

			if disponents := resp.GetDisponents(); disponents != nil {
				found := s.state.CheckIfUserIsDisponent(srv.Context(), job, userId)
				// Either user is a disponent currently and not anymore now,
				// or the user is not a disponent and joined as a disponent now
				if !isDisponent && found {
					restart := true
					resp.Restart = &restart
					isDisponent = true
				} else if isDisponent && !found {
					isDisponent = false
				}
			} else if unitUpdate := resp.GetUnitUpdated(); unitUpdate != nil {
				// Either user is in that unit this update is about or they are not (yet) in an unit
				if unitUpdate.Id == unitId || unitId == 0 {
					if slices.ContainsFunc(unitUpdate.Users, func(a *centrum.UnitAssignment) bool {
						return userId == a.UserId
					}) {
						// Seems that they got assigned to this unit, update the user's unitId here
						unitId = unitUpdate.Id
					} else {
						unitId = 0
						restart := true
						resp.Restart = &restart
					}
				}
			}

			if err := srv.Send(resp); err != nil {
				return true, err
			}

			if resp.Restart != nil && *resp.Restart {
				return false, nil
			}
		}
	}
}
