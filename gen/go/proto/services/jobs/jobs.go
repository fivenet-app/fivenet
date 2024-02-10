package jobs

import (
	"context"
	"database/sql"
	"errors"
	sync "sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	tUser     = table.Users.AS("user")
	tJobProps = table.FivenetJobProps
)

type Server struct {
	JobsConductServiceServer
	JobsServiceServer
	JobsRequestsServiceServer
	JobsTimeclockServiceServer

	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer   trace.Tracer
	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	auditer  audit.IAuditer
	tracker  tracker.ITracker
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
	Tracker           tracker.ITracker
}

func NewServer(p Params) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		logger: p.Logger.Named("jobs"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("jobs"),

		db:       p.DB,
		p:        p.Perms,
		enricher: p.UserAwareEnricher,
		auditer:  p.Audit,
		tracker:  p.Tracker,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runTimeclock()
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterJobsConductServiceServer(srv, s)
	RegisterJobsServiceServer(srv, s)
	RegisterJobsRequestsServiceServer(srv, s)
	RegisterJobsTimeclockServiceServer(srv, s)
}

func (s *Server) GetMOTD(ctx context.Context, req *GetMOTDRequest) (*GetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tJobProps.
		SELECT(
			tJobProps.JobsMotd.AS("getmotdresponse.motd"),
		).
		FROM(tJobProps).
		WHERE(tJobProps.Job.EQ(jet.String(userInfo.Job))).
		LIMIT(1)

	resp := &GetMOTDResponse{}
	if err := stmt.QueryContext(ctx, s.db, resp); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
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
	defer s.auditer.Log(auditEntry, req)

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.JobsMotd,
		).
		VALUES(
			userInfo.Job,
			req.Motd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.JobsMotd.SET(jet.String(req.Motd)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	return &SetMOTDResponse{
		Motd: req.Motd,
	}, nil
}
