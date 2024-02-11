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
	"github.com/galexrt/fivenet/pkg/coords/postals"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/nats-io/nats.go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	CentrumServiceServer

	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer  trace.Tracer
	db      *sql.DB
	ps      perms.Permissions
	auditer audit.IAuditer
	js      nats.JetStreamContext
	tracker tracker.ITracker
	postals postals.Postals

	brokers map[string]*utils.Broker[*StreamResponse]

	publicJobs []string

	state *manager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	DB      *sql.DB
	Perms   perms.Permissions
	Audit   audit.IAuditer
	JS      nats.JetStreamContext
	Tracker tracker.ITracker
	Postals postals.Postals
	Config  *config.Config
	Manager *manager.Manager
}

func NewServer(p Params) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	brokers := map[string]*utils.Broker[*StreamResponse]{}
	for _, job := range p.Config.Game.Livemap.Jobs {
		brokers[job] = utils.NewBroker[*StreamResponse](ctx)
	}

	s := &Server{
		ctx:    ctx,
		logger: p.Logger.Named("centrum"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		db:      p.DB,
		ps:      p.Perms,
		auditer: p.Audit,
		js:      p.JS,
		tracker: p.Tracker,
		postals: p.Postals,

		brokers: brokers,

		publicJobs: p.Config.Game.PublicJobs,

		state: p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		for _, broker := range s.brokers {
			go broker.Start()
		}

		if err := eventscentrum.RegisterStreams(ctx, s.js); err != nil {
			return fmt.Errorf("failed to register events: %w", err)
		}

		if err := s.RegisterSubscriptions(ctx); err != nil {
			return fmt.Errorf("failed to subscribe to events: %w", err)
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		return nil
	}))

	return s, nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCentrumServiceServer(srv, s)
}

func (s *Server) RegisterSubscriptions(ctx context.Context) error {
	if _, err := s.js.Subscribe(fmt.Sprintf("%s.>", eventscentrum.BaseSubject), s.watchForChanges, nats.DeliverLastPerSubject()); err != nil {
		return err
	}

	return nil
}

func (s *Server) watchForChanges(msg *nats.Msg) {
	job, topic, tType := eventscentrum.SplitSubject(msg.Subject)

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
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_Disponents{
				Disponents: dest,
			}

		case eventscentrum.TypeGeneralSettings:
			dest := &centrum.Settings{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
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
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitCreated{
				UnitCreated: dest,
			}

		case eventscentrum.TypeUnitDeleted:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitDeleted{
				UnitDeleted: dest,
			}

		case eventscentrum.TypeUnitUpdated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitUpdated{
				UnitUpdated: dest,
			}

		case eventscentrum.TypeUnitStatus:
			dest := &centrum.UnitStatus{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
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
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchCreated{
				DispatchCreated: dest,
			}

		case eventscentrum.TypeDispatchDeleted:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchDeleted{
				DispatchDeleted: dest,
			}

		case eventscentrum.TypeDispatchUpdated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchUpdated{
				DispatchUpdated: dest,
			}

		case eventscentrum.TypeDispatchStatus:
			dest := &centrum.DispatchStatus{}
			if err := proto.Unmarshal(msg.Data, dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchStatus{
				DispatchStatus: dest,
			}
		}
	}

	broker.Publish(resp)

	meta, _ := msg.Metadata()
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
	defer s.auditer.Log(auditEntry, req)

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
	settings := s.state.GetSettings(job)
	disponents, _ := s.state.GetDisponents(job)
	isDisponent := s.state.CheckIfUserIsDisponent(job, userId)
	ownUnitId, _ := s.state.GetUserUnitID(userId)
	units, _ := s.state.ListUnits(job)
	ownUnit, _ := s.state.GetUnit(job, ownUnitId)

	dispatches := s.state.FilterDispatches(job, nil, []centrum.StatusDispatch{
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

func (s *Server) getJobBroker(job string) (*utils.Broker[*StreamResponse], bool) {
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
				found := s.state.CheckIfUserIsDisponent(job, userId)
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
