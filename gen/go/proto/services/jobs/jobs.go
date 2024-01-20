package jobs

import (
	"context"
	"database/sql"
	sync "sync"

	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/query/fivenet/table"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	tUser = table.Users.AS("user")
)

type Server struct {
	ConductServiceServer
	JobsServiceServer
	RequestsServiceServer
	TimeclockServiceServer

	ctx    context.Context
	logger *zap.Logger
	wg     sync.WaitGroup

	tracer   trace.Tracer
	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.Enricher
	auditer  audit.IAuditer
	tracker  tracker.ITracker
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	Perms    perms.Permissions
	Enricher *mstlystcdata.Enricher
	Audit    audit.IAuditer
	Tracker  tracker.ITracker
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
		enricher: p.Enricher,
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
	RegisterConductServiceServer(srv, s)
	RegisterJobsServiceServer(srv, s)
	RegisterRequestsServiceServer(srv, s)
	RegisterTimeclockServiceServer(srv, s)
}
