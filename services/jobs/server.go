package jobs

import (
	"context"
	"database/sql"
	sync "sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	pbjobs.JobsConductServiceServer
	pbjobs.JobsServiceServer
	pbjobs.JobsTimeclockServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup

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
	DB                *sql.DB
	Perms             perms.Permissions
	UserAwareEnricher *mstlystcdata.UserAwareEnricher
	Audit             audit.IAuditer
	Config            *config.Config

	Cron croner.ICron
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("jobs"),
		wg:     sync.WaitGroup{},

		db:       p.DB,
		ps:       p.Perms,
		enricher: p.UserAwareEnricher,
		aud:      p.Audit,

		customDB: p.Config.Database.Custom,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "jobs.timeclock_cleanup",
			Schedule: "@daily", // Daily
		}); err != nil {
			return err
		}

		p.Cron.UnregisterCronjob(ctx, "jobs.timeclock_handling")
		p.Cron.UnregisterCronjob(ctx, "jobs-timeclock-handling")

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbjobs.RegisterJobsConductServiceServer(srv, s)
	pbjobs.RegisterJobsServiceServer(srv, s)
	pbjobs.RegisterJobsTimeclockServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbjobs.PermsRemap
}
