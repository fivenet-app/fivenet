package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/tracker/postals"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/maps"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/nats-io/nats.go"
	"github.com/puzpuzpuz/xsync/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const pingTickerTime = 35 * time.Second

var (
	ErrFailedQuery       = status.Error(codes.Internal, "errors.CentrumService.ErrFailedQuery")
	ErrAlreadyInUnit     = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrAlreadyInUnit")
	ErrNotPartOfDispatch = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrNotPartOfDispatch")
	ErrNotOnDuty         = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrNotOnDuty")
)

var (
	tCentrumUsers = table.FivenetCentrumUsers
)

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

	visibleJobs []string
	convertJobs []string

	settings   *xsync.MapOf[string, *dispatch.Settings]
	disponents *xsync.MapOf[string, []*users.UserShort]
	units      *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]]
	dispatches *xsync.MapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]]

	userIDToUnitID *xsync.MapOf[int32, uint64]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	Perms    perms.Permissions
	Audit    audit.IAuditer
	Events   *events.Eventus
	Enricher *mstlystcdata.Enricher
	Tracker  *tracker.Tracker
	Postals  *postals.Postals
	Config   *config.Config
}

func NewServer(p Params) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		logger: p.Logger,
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		db:       p.DB,
		ps:       p.Perms,
		auditer:  p.Audit,
		events:   p.Events,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		postals:  p.Postals,

		visibleJobs: p.Config.Game.Livemap.Jobs,
		convertJobs: p.Config.Game.DispatchCenter.ConvertJobs,

		settings:   xsync.NewTypedMapOf[string, *dispatch.Settings](maps.HashString),
		disponents: xsync.NewTypedMapOf[string, []*users.UserShort](maps.HashString),
		units:      xsync.NewTypedMapOf[string, *xsync.MapOf[uint64, *dispatch.Unit]](maps.HashString),
		dispatches: xsync.NewTypedMapOf[string, *xsync.MapOf[uint64, *dispatch.Dispatch]](maps.HashString),

		userIDToUnitID: xsync.NewIntegerMapOf[int32, uint64](),
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		if err := s.registerEvents(); err != nil {
			return fmt.Errorf("failed to register events: %w", err)
		}

		if err := s.loadData(); err != nil {
			return err
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.loadDataLoop()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.ConvertPhoneJobMsgToDispatch()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchStateEvents()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchUserChanges()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.housekeeper()
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

func (s *Server) loadDataLoop() {
	for {
		if err := s.loadData(); err != nil {
			s.logger.Error("failed to refresh centrum data", zap.Error(err))
		}

		select {
		case <-s.ctx.Done():
			return
		case <-time.After(12 * time.Second):
		}
	}
}

func (s *Server) TakeControl(ctx context.Context, req *TakeControlRequest) (*TakeControlResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeControl",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	if err := s.dispatchCenterSignOn(ctx, userInfo.Job, userInfo.UserId, req.Signon); err != nil {
		return nil, err
	}

	if req.Signon {
		auditEntry.State = int16(rector.EVENT_TYPE_CREATED)
	} else {
		auditEntry.State = int16(rector.EVENT_TYPE_DELETED)
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &TakeControlResponse{}, nil
}

func (s *Server) sendLatestState(srv CentrumService_StreamServer, job string, userId int32) (uint64, bool, error) {
	settings := s.getSettings(job)
	disponents := s.getDisponents(job)
	isDisponent := s.checkIfUserIsDisponent(job, userId)

	unitId, _ := s.getUnitIDForUserID(userId)

	units, err := s.listUnits(job)
	if err != nil {
		return 0, isDisponent, err
	}

	ownUnit, _ := s.getUnit(job, unitId)

	dispatches, err := s.ListDispatches(srv.Context(), &ListDispatchesRequest{
		NotStatus: []dispatch.DISPATCH_STATUS{
			dispatch.DISPATCH_STATUS_ARCHIVED,
			dispatch.DISPATCH_STATUS_CANCELLED,
			dispatch.DISPATCH_STATUS_COMPLETED,
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
	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.>", BaseSubject, job), msgCh, nats.DeliverNew())
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
			topic, tType := s.getEventTypeFromSubject(msg.Subject)

			switch topic {
			case TopicGeneral:
				switch tType {
				case TypeGeneralDisponents:
					var dest DisponentsChange
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_Disponents{
						Disponents: &dest,
					}

					found := s.checkIfUserIsDisponent(job, userId)
					// Either user is a disponent currently and not anymore now,
					// or the user is not a disponent and joined as a disponent now
					if (isDisponent && !found) || (!isDisponent && found) {
						restart := true
						resp.Restart = &restart
					}

				case TypeGeneralSettings:
					var dest dispatch.Settings
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_Settings{
						Settings: &dest,
					}
				}

			case TopicDispatch:
				switch tType {
				case TypeDispatchCreated:
					var dest dispatch.Dispatch
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchCreated{
						DispatchCreated: &dest,
					}

				case TypeDispatchDeleted:
					var dest dispatch.Dispatch
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchDeleted{
						DispatchDeleted: dest.Id,
					}

				case TypeDispatchUpdated:
					var dest dispatch.Dispatch
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchUpdated{
						DispatchUpdated: &dest,
					}

				case TypeDispatchStatus:
					var dest dispatch.Dispatch
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_DispatchStatus{
						DispatchStatus: &dest,
					}

				}

			case TopicUnit:
				switch tType {
				case TypeUnitUserAssigned:
					var dest dispatch.Unit
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitAssigned{
						UnitAssigned: &dest,
					}

					// Either user is in the unit this update is about or they are not (yet) in an unit
					if dest.Id == unitId || unitId == 0 {
						// Restart user's stream if they got removed from their assigned unit
						inUnit := utils.InSliceFunc(dest.Users, func(a *dispatch.UnitAssignment) bool {
							return userId == a.UserId
						})

						if dest.Id == unitId && !inUnit {
							restart := true
							resp.Restart = &restart
						} else if inUnit {
							// Seems that they got assigned to this unit, update the user's unitId here
							unitId = dest.Id
						}
					}

				case TypeUnitDeleted:
					var dest dispatch.Unit
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitDeleted{
						UnitDeleted: &dest,
					}

				case TypeUnitUpdated:
					var dest dispatch.Unit
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitUpdated{
						UnitUpdated: &dest,
					}

				case TypeUnitStatus:
					var dest dispatch.Unit
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return true, err
					}

					resp.Change = &StreamResponse_UnitStatus{
						UnitStatus: &dest,
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
