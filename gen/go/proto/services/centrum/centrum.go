package centrum

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/store"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go"
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

	ctx     context.Context
	logger  *zap.Logger
	tracer  trace.Tracer
	db      *sql.DB
	p       perms.Permissions
	a       audit.IAuditer
	events  *events.Eventus
	tracker *tracker.Tracker

	visibleJobs []string

	settings   *store.Store[*dispatch.Settings]
	disponents *store.Store[*users.UserShort]
	units      *store.Store[*dispatch.Unit]
	dispatches *store.Store[*dispatch.Dispatch]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	TP      *tracesdk.TracerProvider
	DB      *sql.DB
	Perms   perms.Permissions
	Audit   audit.IAuditer
	Events  *events.Eventus
	Tracker *tracker.Tracker
	Config  *config.Config
}

func NewServer(p Params) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		logger: p.Logger,

		tracer: p.TP.Tracer("centrum-cache"),

		db:      p.DB,
		p:       p.Perms,
		a:       p.Audit,
		events:  p.Events,
		tracker: p.Tracker,

		visibleJobs: p.Config.Game.Livemap.Jobs,

		settings:   store.New[*dispatch.Settings]("centrum_settings"),
		disponents: store.New[*users.UserShort]("centrum_disponents"),
		units:      store.New[*dispatch.Unit]("centrum_units"),
		dispatches: store.New[*dispatch.Dispatch]("centrum_dispatches"),
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		if err := s.settings.Start(p.Events.JS); err != nil {
			return err
		}
		if err := s.disponents.Start(p.Events.JS); err != nil {
			return err
		}
		if err := s.units.Start(p.Events.JS); err != nil {
			return err
		}
		if err := s.dispatches.Start(p.Events.JS); err != nil {
			return err
		}

		if err := s.registerEvents(); err != nil {
			return fmt.Errorf("failed to register events: %w", err)
		}

		if err := s.loadInitialData(); err != nil {
			return err
		}

		go s.start()
		//go s.ConvertPhoneJobMsgToDispatch()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return s, nil
}

func (s *Server) start() {
	go func() {
		for {
			if err := s.refresh(); err != nil {
				s.logger.Error("failed to refresh centrum data", zap.Error(err))
			}

			select {
			case <-s.ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}()
}

func (s *Server) GetSettings(ctx context.Context, req *GetSettingsRequest) (*dispatch.Settings, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings := &dispatch.Settings{}
	if err := s.settings.Get(userInfo.Job, settings); err != nil {
		if !errors.Is(nats.ErrKeyNotFound, err) {
			return nil, err
		}
	}

	settings.Job = userInfo.Job

	return settings, nil
}

func (s *Server) UpdateSettings(ctx context.Context, req *dispatch.Settings) (*dispatch.Settings, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateSettings",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	stmt := tCentrumSettings.
		INSERT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
		).
		VALUES(
			userInfo.Job,
			req.Enabled,
			req.Mode,
			req.FallbackMode,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumSettings.Job.SET(jet.String(userInfo.Job)),
			tCentrumSettings.Enabled.SET(jet.Bool(req.Enabled)),
			tCentrumSettings.Mode.SET(jet.Int32(int32(req.Mode))),
			tCentrumSettings.FallbackMode.SET(jet.Int32(int32(req.FallbackMode))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	// Load settings from database so they are updated in the "cache"
	if err := s.loadSettings(ctx, userInfo.Job); err != nil {
		return nil, err
	}

	settings, err := s.getSettings(ctx, userInfo.Job)
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(settings)
	if err != nil {
		return nil, err
	}
	s.broadcastToAllUnits(TopicGeneral, TypeGeneralSettings, userInfo, data)

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return settings, nil
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

	disponents, err := s.getDisponents(ctx, userInfo.Job)
	if err != nil {
		return nil, err
	}

	change := &DisponentsChange{
		Disponents: disponents,
	}
	data, err := proto.Marshal(change)
	if err != nil {
		return nil, err
	}
	s.broadcastToAllUnits(TopicGeneral, TypeGeneralDisponents, userInfo, data)

	return &TakeControlResponse{}, nil
}

func (s *Server) waitForUnit(srv CentrumService_StreamServer, userInfo *userinfo.UserInfo) (uint64, error) {
	msgCh := make(chan *nats.Msg, 4)
	sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.%s.%s.%d", BaseSubject, userInfo.Job, TopicUnit, TypeUnitUserAssigned, 0), msgCh)
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

				resp := &StreamResponse{
					Change: &StreamResponse_UnitAssigned{
						UnitAssigned: &dest,
					},
				}

				// TODO check if user is in unit users list

				if err := srv.Send(resp); err != nil {
					return 0, err
				}

				return dest.Id, nil
			}
		}
	}
}

func (s *Server) Stream(req *StreamRequest, srv CentrumService_StreamServer) error {
	userInfo := auth.MustGetUserInfoFromContext(srv.Context())

	settings, err := s.getSettings(srv.Context(), userInfo.Job)
	if err != nil {
		return err
	}
	disponents, err := s.getDisponents(srv.Context(), userInfo.Job)
	if err != nil {
		return err
	}

	controller := utils.InSliceFunc(disponents, func(in *users.UserShort) bool {
		return in.UserId == userInfo.UserId
	})

	unitId, err := s.getUnitIDForUserID(srv.Context(), userInfo.UserId)
	if err != nil {
		return err
	}

	units, err := s.listUnits(srv.Context(), userInfo.Job)
	if err != nil {
		return err
	}

	unit, _ := s.getUnit(srv.Context(), userInfo, unitId)

	dispatches, err := s.ListDispatches(srv.Context(), &ListDispatchesRequest{})
	if err != nil {
		return err
	}

	// Send initial state message to client
	resp := &StreamResponse{
		Change: &StreamResponse_LatestState{
			LatestState: &LatestState{
				IsDisponent: true,
				Settings:    settings,
				Unit:        unit,
				Units:       units,
				Dispatches:  dispatches.Dispatches,
			},
		},
	}
	if err := srv.Send(resp); err != nil {
		return err
	}

	if !controller && unitId <= 0 {
		unitId, err = s.waitForUnit(srv, userInfo)
		if err != nil {
			return err
		}
	}

	msgCh := make(chan *nats.Msg, 48)
	if !controller {
		sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.*.*.%d", BaseSubject, userInfo.Job, unitId), msgCh)
		if err != nil {
			return err
		}
		defer sub.Unsubscribe()
	} else {
		sub, err := s.events.JS.ChanSubscribe(fmt.Sprintf("%s.%s.>", BaseSubject, userInfo.Job), msgCh)
		if err != nil {
			return err
		}
		defer sub.Unsubscribe()
	}

	// Watch for events from message queue
	for {
		select {
		case <-srv.Context().Done():
			return nil
		case msg := <-msgCh:
			msg.Ack()
			topic, tType := s.getEventTypeFromSubject(msg.Subject)

			switch topic {
			case TopicGeneral:
				switch tType {
				case TypeGeneralDisponents:
					var dest DisponentsChange
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
					}

					resp.Change = &StreamResponse_Disponents{
						Disponents: &dest,
					}

				case TypeGeneralSettings:
					var dest dispatch.Settings
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
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
						return err
					}

					resp.Change = &StreamResponse_DispatchCreated{
						DispatchCreated: &dest,
					}
				case TypeDispatchDeleted:
					var dest dispatch.Dispatch
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
					}

					resp.Change = &StreamResponse_DispatchDeleted{
						DispatchDeleted: &dest,
					}

				case TypeDispatchUpdated:
					var dest dispatch.Dispatch
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
					}

					resp.Change = &StreamResponse_DispatchUpdated{
						DispatchUpdated: &dest,
					}

				case TypeDispatchStatus:
					var dest dispatch.DispatchStatus
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
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
						return err
					}

					resp.Change = &StreamResponse_UnitAssigned{
						UnitAssigned: &dest,
					}

				case TypeUnitDeleted:
					var dest dispatch.Unit
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
					}

					resp.Change = &StreamResponse_UnitDeleted{
						UnitDeleted: &dest,
					}

				case TypeUnitUpdated:
					var dest dispatch.Unit
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
					}

					resp.Change = &StreamResponse_UnitUpdated{
						UnitUpdated: &dest,
					}

				case TypeUnitStatus:
					var dest dispatch.UnitStatus
					if err := proto.Unmarshal(msg.Data, &dest); err != nil {
						return err
					}

					resp.Change = &StreamResponse_UnitStatus{
						UnitStatus: &dest,
					}
				}
			}
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
	}
}
