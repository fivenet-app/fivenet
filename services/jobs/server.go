package jobs

import (
	"database/sql"
	sync "sync"

	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetJobColleagueProps,
		JobColumn:       table.FivenetJobColleagueProps.Job,
		IDColumn:        table.FivenetJobColleagueProps.UserID,
		DeletedAtColumn: table.FivenetJobColleagueProps.DeletedAt,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:      table.FivenetJobColleagueActivity,
				JobColumn:  table.FivenetJobColleagueActivity.Job,
				IDColumn:   table.FivenetJobColleagueActivity.ID,
				ForeignKey: table.FivenetJobColleagueActivity.TargetUserID,
			},
		},
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetJobLabels,
		JobColumn:       table.FivenetJobLabels.Job,
		IDColumn:        table.FivenetJobLabels.ID,
		DeletedAtColumn: table.FivenetJobLabels.DeletedAt,

		MinDays: 60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetJobConduct,
		JobColumn:       table.FivenetJobConduct.Job,
		IDColumn:        table.FivenetJobConduct.ID,
		DeletedAtColumn: table.FivenetJobConduct.DeletedAt,

		MinDays: 60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:      table.FivenetJobTimeclock,
		DateColumn: table.FivenetJobTimeclock.Date,
		JobColumn:  table.FivenetJobTimeclock.Job,

		MinDays: 365, // One year retention
	})
}

type Server struct {
	pbjobs.ConductServiceServer
	pbjobs.JobsServiceServer
	pbjobs.TimeclockServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	notifi   notifi.INotifi

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
	Notifi            notifi.INotifi
	Config            *config.Config
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("jobs"),
		wg:     sync.WaitGroup{},

		db:       p.DB,
		ps:       p.Perms,
		enricher: p.UserAwareEnricher,
		aud:      p.Audit,
		notifi:   p.Notifi,

		customDB: p.Config.Database.Custom,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbjobs.RegisterConductServiceServer(srv, s)
	pbjobs.RegisterJobsServiceServer(srv, s)
	pbjobs.RegisterTimeclockServiceServer(srv, s)
}

// GetPermsRemap returns the permissions re-mapping for the services.
func (s *Server) GetPermsRemap() map[string]string {
	return pbjobs.PermsRemap
}
