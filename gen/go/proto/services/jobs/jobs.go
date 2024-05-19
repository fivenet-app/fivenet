package jobs

import (
	"context"
	"database/sql"
	"errors"
	sync "sync"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	tUser      = table.Users.AS("user")
	tUserProps = table.FivenetUserProps
	tJobProps  = table.FivenetJobProps
)

type Server struct {
	JobsConductServiceServer
	JobsServiceServer
	JobsTimeclockServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup

	tracer   trace.Tracer
	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer

	customDB config.CustomDB
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	TP                *tracesdk.TracerProvider
	DB                *sql.DB
	Perms             perms.Permissions
	UserAwareEnricher *mstlystcdata.UserAwareEnricher
	Audit             audit.IAuditer
	Config            *config.Config
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("jobs"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("jobs"),

		db:       p.DB,
		ps:       p.Perms,
		enricher: p.UserAwareEnricher,
		aud:      p.Audit,

		customDB: p.Config.Database.Custom,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterJobsConductServiceServer(srv, s)
	RegisterJobsServiceServer(srv, s)
	RegisterJobsTimeclockServiceServer(srv, s)
}

func (s *Server) GetMOTD(ctx context.Context, req *GetMOTDRequest) (*GetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tJobProps.
		SELECT(
			tJobProps.Motd.AS("getmotdresponse.motd"),
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(userInfo.Job)),
		).
		LIMIT(1)

	resp := &GetMOTDResponse{}
	if err := stmt.QueryContext(ctx, s.db, resp); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) SetMOTD(ctx context.Context, req *SetMOTDRequest) (*SetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "SetMOTD",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.Motd,
		).
		VALUES(
			userInfo.Job,
			req.Motd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Motd.SET(jet.String(req.Motd)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &SetMOTDResponse{
		Motd: req.Motd,
	}, nil
}
