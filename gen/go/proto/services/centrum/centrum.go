package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/pkg/utils/maps"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
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

var (
	ErrFailedQuery   = status.Error(codes.Internal, "errors.CentrumService.ErrFailedQuery")
	ErrAlreadyInUnit = status.Error(codes.InvalidArgument, "errors.CentrumService.ErrAlreadyInUnit")
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

	visibleJobs []string

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

		visibleJobs: p.Config.Game.Livemap.Jobs,

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
			s.start()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.ConvertPhoneJobMsgToDispatch()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchForEvents()
		}()
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.watchForUserChanges()
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

func (s *Server) start() {
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

	if req.Signon {
		if _, ok := s.tracker.GetUserByJobAndID(userInfo.Job, userInfo.UserId); !ok {
			return nil, status.Error(codes.InvalidArgument, "You are not on duty!")
		}

		stmt := tCentrumUsers.
			INSERT(
				tCentrumUsers.Job,
				tCentrumUsers.UserID,
				tCentrumUsers.Identifier,
			).
			VALUES(
				userInfo.Job,
				userInfo.UserId,
				tUsers.
					SELECT(
						tUsers.Identifier.AS("identifier"),
					).
					FROM(tUsers).
					WHERE(
						tUsers.ID.EQ(jet.Int32(userInfo.UserId)),
					).
					LIMIT(1),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, err
			}
		}
	} else {
		stmt := tCentrumUsers.
			DELETE().
			WHERE(jet.AND(
				tCentrumUsers.Job.EQ(jet.String(userInfo.Job)),
				tCentrumUsers.UserID.EQ(jet.Int32(userInfo.UserId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}
	}

	// Load updated disponents into state
	if err := s.loadDisponents(ctx, userInfo.Job); err != nil {
		return nil, err
	}

	disponents := s.getDisponents(ctx, userInfo.Job)
	change := &DisponentsChange{
		Disponents: disponents,
	}
	data, err := proto.Marshal(change)
	if err != nil {
		return nil, err
	}
	s.broadcastToAllUnits(TopicGeneral, TypeGeneralDisponents, userInfo.Job, data)

	return &TakeControlResponse{}, nil
}

func (s *Server) waitForUserToBeAssignedUnit(srv CentrumService_StreamServer, job string, userId int32) (uint64, error) {
	msgCh := make(chan *nats.Msg, 8)
	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.%s.%s.%d", BaseSubject, job, TopicUnit, TypeUnitUserAssigned, 0), msgCh)
	if err != nil {
		return 0, err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-srv.Context().Done():
			return 0, nil

		case msg := <-msgCh:
			msg.Ack()
			topic, tType := s.getEventTypeFromSubject(msg.Subject)

			if topic == TopicUnit && tType == TypeUnitUserAssigned {
				var dest dispatch.Unit
				if err := proto.Unmarshal(msg.Data, &dest); err != nil {
					return 0, err
				}

				if !utils.InSliceFunc(dest.Users, func(a *dispatch.UnitAssignment) bool {
					return userId == a.UserId
				}) {
					continue
				}

				resp := &StreamResponse{
					Change: &StreamResponse_UnitAssigned{
						UnitAssigned: &dest,
					},
				}

				if err := srv.Send(resp); err != nil {
					return 0, err
				}

				return dest.Id, nil
			}
		}
	}
}

func (s *Server) sendLatestState(srv CentrumService_StreamServer, job string, userId int32) (uint64, bool, error) {
	settings := s.getSettings(srv.Context(), job)
	disponents := s.getDisponents(srv.Context(), job)

	isController := utils.InSliceFunc(disponents, func(in *users.UserShort) bool {
		return in.UserId == userId
	})

	unitId, _ := s.getUnitIDForUserID(userId)

	units, err := s.listUnits(job)
	if err != nil {
		return 0, isController, err
	}

	unit, _ := s.getUnit(job, unitId)

	ownOnly := !isController

	dispatches, err := s.ListDispatches(srv.Context(), &ListDispatchesRequest{
		OwnOnly: &ownOnly,
		NotStatus: []dispatch.DISPATCH_STATUS{
			dispatch.DISPATCH_STATUS_ARCHIVED,
		},
	})
	if err != nil {
		return 0, isController, err
	}

	// Send initial state message to client
	resp := &StreamResponse{
		Change: &StreamResponse_LatestState{
			LatestState: &LatestState{
				IsDisponent: isController,
				Settings:    settings,
				Unit:        unit,
				Units:       units,
				Dispatches:  dispatches.Dispatches,
			},
		},
	}
	if err := srv.Send(resp); err != nil {
		return 0, isController, err
	}

	return unitId, isController, nil
}

func (s *Server) Stream(req *StreamRequest, srv CentrumService_StreamServer) error {
	userInfo := *auth.MustGetUserInfoFromContext(srv.Context())

	for {
		unitId, isController, err := s.sendLatestState(srv, userInfo.Job, userInfo.UserId)
		if err != nil {
			return err
		}

		// If user is not in unit and not a controller, wait for user to be assigned/join an unit
		if unitId <= 0 && !isController {
			if _, err := s.waitForUserToBeAssignedUnit(srv, userInfo.Job, userInfo.UserId); err != nil {
				return err
			}

			unitId, isController, err = s.sendLatestState(srv, userInfo.Job, userInfo.UserId)
			if err != nil {
				return err
			}
		}

		end, err := s.stream(srv, isController, userInfo.Job, userInfo.UserId, unitId)
		if end {
			return err
		} else if err != nil {
			s.logger.Error("error during stream", zap.Error(err))
		}

		time.Sleep(150 * time.Millisecond)
	}
}

func (s *Server) stream(srv CentrumService_StreamServer, isController bool, job string, userId int32, unitId uint64) (bool, error) {
	msgCh := make(chan *nats.Msg, 48)
	if !isController {
		sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.*.*.%d", BaseSubject, job, unitId), msgCh)
		if err != nil {
			return true, err
		}
		defer sub.Unsubscribe()
	} else {
		sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.>", BaseSubject, job), msgCh)
		if err != nil {
			return true, err
		}
		defer sub.Unsubscribe()
	}

	resp := &StreamResponse{}
	// Watch for events from message queue
	for {
		select {
		case <-srv.Context().Done():
			return true, nil
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
						DispatchDeleted: &dest,
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

					// If not user's unit, skip
					if unitId != dest.Id {
						continue
					}

					resp.Change = &StreamResponse_UnitAssigned{
						UnitAssigned: &dest,
					}

					if utils.InSliceFunc(dest.Users, func(a *dispatch.UnitAssignment) bool {
						return userId == a.UserId
					}) {
					} else {
						restart := true
						resp.Restart = &restart
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
	}
}
