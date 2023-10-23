package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/bot"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/gen/go/proto/services/centrum/manager"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/tracker/postals"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/nats-io/nats.go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const pingTickerTime = 30 * time.Second

type Server struct {
	CentrumServiceServer

	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer   trace.Tracer
	db       *sql.DB
	ps       perms.Permissions
	auditer  audit.IAuditer
	events   *events.Eventus
	enricher *mstlystcdata.Enricher
	tracker  *tracker.Tracker
	postals  *postals.Postals

	convertJobs []string

	state *manager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger     *zap.Logger
	TP         *tracesdk.TracerProvider
	DB         *sql.DB
	Perms      perms.Permissions
	Audit      audit.IAuditer
	Events     *events.Eventus
	Enricher   *mstlystcdata.Enricher
	Tracker    *tracker.Tracker
	Postals    *postals.Postals
	Config     *config.Config
	Manager    *manager.Manager
	BotManager *bot.Manager
}

func NewServer(p Params) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		logger: p.Logger.Named("centrum"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		db:       p.DB,
		ps:       p.Perms,
		auditer:  p.Audit,
		events:   p.Events,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		postals:  p.Postals,

		convertJobs: p.Config.Game.DispatchCenter.ConvertJobs,

		state: p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := eventscentrum.RegisterEvents(ctx, s.events); err != nil {
			return fmt.Errorf("failed to register events: %w", err)
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.ConvertPhoneJobMsgToDispatch()
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		return nil
	}))

	return s, nil
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
	disponents := s.state.GetDisponents(job)
	isDisponent := s.state.CheckIfUserIsDisponent(job, userId)
	unitId, _ := s.state.GetUnitIDForUserID(userId)
	units, _ := s.state.ListUnits(job)
	ownUnit, _ := s.state.GetUnit(job, unitId)

	dispatches, err := s.ListDispatches(srv.Context(), &ListDispatchesRequest{
		NotStatus: []dispatch.StatusDispatch{
			dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED,
		},
	})
	if err != nil {
		return 0, isDisponent, err
	}

	// Send initial state message to client
	resp := &StreamResponse{
		Change: &StreamResponse_LatestState{
			LatestState: &LatestState{
				Settings:    settings,
				Disponents:  disponents,
				IsDisponent: isDisponent,
				OwnUnit:     ownUnit,
				Units:       units,
				Dispatches:  dispatches.Dispatches,
			},
		},
	}
	if err := srv.Send(resp); err != nil {
		return 0, isDisponent, err
	}

	return unitId, isDisponent, nil
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

		time.Sleep(150 * time.Millisecond)
	}
}

func (s *Server) stream(srv CentrumService_StreamServer, isDisponent bool, job string, userId int32, unitId uint64) (bool, error) {
	msgCh := make(chan *nats.Msg, 48)
	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.>", eventscentrum.BaseSubject, job), msgCh, nats.DeliverNew())
	if err != nil {
		return true, err
	}
	defer sub.Unsubscribe()

	// Ping ticker to ensure better stream quality
	ticker := time.NewTicker(pingTickerTime * 2)
	defer ticker.Stop()

	// Watch for events from message queue
	for {
		resp := &StreamResponse{}

		select {
		case <-srv.Context().Done():
			return true, nil

		case t := <-ticker.C:
			resp.Change = &StreamResponse_Ping{
				Ping: t.String(),
			}

		case msg := <-msgCh:
			msg.Ack()

			topic, tType := eventscentrum.GetEventTypeFromSubject(msg.Subject)

			switch topic {
			case eventscentrum.TopicGeneral:
				switch tType {
				case eventscentrum.TypeGeneralDisponents:
					dest := &dispatch.DisponentsChange{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_Disponents{
						Disponents: dest,
					}

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

				case eventscentrum.TypeGeneralSettings:
					dest := &dispatch.Settings{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_Settings{
						Settings: dest,
					}
				}

			case eventscentrum.TopicDispatch:
				switch tType {
				case eventscentrum.TypeDispatchCreated:
					dest := &dispatch.Dispatch{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchCreated{
						DispatchCreated: dest,
					}

				case eventscentrum.TypeDispatchDeleted:
					dest := &dispatch.Dispatch{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchDeleted{
						DispatchDeleted: dest.Id,
					}

				case eventscentrum.TypeDispatchUpdated:
					dest := &dispatch.Dispatch{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchUpdated{
						DispatchUpdated: dest,
					}

				case eventscentrum.TypeDispatchStatus:
					dest := &dispatch.DispatchStatus{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchStatus{
						DispatchStatus: dest,
					}
				}

			case eventscentrum.TopicUnit:
				switch tType {
				case eventscentrum.TypeUnitDeleted:
					dest := &dispatch.Unit{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitDeleted{
						UnitDeleted: dest,
					}

				case eventscentrum.TypeUnitUpdated:
					dest := &dispatch.Unit{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitUpdated{
						UnitUpdated: dest,
					}

					// Either user is in that unit this update is about or they are not (yet) in an unit
					if dest.Id == unitId || unitId == 0 {
						if utils.InSliceFunc(dest.Users, func(a *dispatch.UnitAssignment) bool {
							return userId == a.UserId
						}) {
							// Seems that they got assigned to this unit, update the user's unitId here
							unitId = dest.Id
						}
					}

				case eventscentrum.TypeUnitStatus:
					dest := &dispatch.UnitStatus{}
					if err := proto.Unmarshal(msg.Data, dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitStatus{
						UnitStatus: dest,
					}
				}
			}
		}

		if err := srv.Send(resp); err != nil {
			return true, err
		}

		if resp.Restart != nil && *resp.Restart {
			return false, nil
		}

		// Reset ping ticker after every (successful) response
		ticker.Reset(pingTickerTime)
	}
}
