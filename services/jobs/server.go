package jobs

import (
	"database/sql"
	sync "sync"

	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
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
	notifi   notifi.INotifi

	customDB config.CustomDB

	fHandler *filestore.Handler[int64]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	DB                *sql.DB
	Config            *config.Config
	Perms             perms.Permissions
	UserAwareEnricher *mstlystcdata.UserAwareEnricher
	Notifi            notifi.INotifi
	Storage           storage.IStorage
}

func NewServer(p Params) *Server {
	conductFileHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tConductFiles,
		tConductFiles.ConductID,
		tConductFiles.FileID,
		3<<20, // 3 MiB limit
		func(parentId int64) mysql.BoolExpression {
			return tConductFiles.ConductID.EQ(mysql.Int64(parentId))
		},
		filestore.InsertJoinRow,
		false,
	)

	s := &Server{
		logger: p.Logger.Named("jobs"),
		wg:     sync.WaitGroup{},

		db:       p.DB,
		ps:       p.Perms,
		enricher: p.UserAwareEnricher,
		notifi:   p.Notifi,

		customDB: p.Config.Database.Custom,

		fHandler: conductFileHandler,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbjobs.RegisterConductServiceServer(srv, s)
	pbjobs.RegisterJobsServiceServer(srv, s)
	pbjobs.RegisterTimeclockServiceServer(srv, s)
}
