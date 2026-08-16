package jobs

import (
	"database/sql"
	sync "sync"

	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	colleaguehydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/fivenet-app/fivenet/v2026/stores/jobs/usersel"
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

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetJobGroups,
		IDColumn:        table.FivenetJobGroups.ID,
		JobColumn:       table.FivenetJobGroups.Job,
		DeletedAtColumn: table.FivenetJobGroups.DeletedAt,

		MinDays: 60,
	})
}

type Server struct {
	pbjobs.ConductServiceServer
	pbjobs.ColleaguesServiceServer
	pbjobs.JobsServiceServer
	pbjobs.TimeclockServiceServer
	pbjobs.StatsServiceServer
	pbjobs.UnimplementedGroupsServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup

	db       *sql.DB
	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	notifi   notifi.INotifi
	stats    *stats.Service

	customDB *config.CustomDB
	store    jobsstore.IStore

	groupAccess         access.JobGroupsAccess
	groupAccessResolver *access.SubjectResolver
	qualificationAccess qualificationAccessChecker

	colleagueHydrator colleaguehydrator.IHydrator

	fHandler             *filestore.Handler[int64]
	groupLogoFileHandler *filestore.Handler[int64]

	userSel usersel.IResolver
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger              *zap.Logger
	DB                  *sql.DB
	Config              *config.Config
	Perms               perms.Permissions
	UserAwareEnricher   mstlystcdata.IUserAwareEnricher
	Notifi              notifi.INotifi
	Storage             storage.IStorage
	Stats               *stats.Service
	Store               jobsstore.IStore
	ColleagueHydrator   colleaguehydrator.IHydrator
	UserSel             usersel.IResolver
	GroupAccess         *access.JobGroupsObjectAccess
	QualificationAccess *access.QualificationsObjectAccess
}

func NewServer(p Params) *Server {
	conductFileHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tConductFiles,
		tConductFiles.ConductID,
		tConductFiles.FileID,
		3<<20, // 3 MiB limit
		5,
		func(parentId int64) mysql.BoolExpression {
			return tConductFiles.ConductID.EQ(mysql.Int64(parentId))
		},
		filestore.InsertJoinRow,
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	tJobGroups := table.FivenetJobGroups
	groupLogoFileHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tJobGroups,
		tJobGroups.ID,
		tJobGroups.LogoFileID,
		2<<20,
		1,
		func(parentID int64) mysql.BoolExpression {
			return tJobGroups.ID.EQ(mysql.Int64(parentID))
		},
		filestore.UpdateJoinRow,
		true,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	s := &Server{
		logger: p.Logger.Named("jobs"),
		wg:     sync.WaitGroup{},

		db:       p.DB,
		perms:    p.Perms,
		enricher: p.UserAwareEnricher,
		notifi:   p.Notifi,
		stats:    p.Stats,

		customDB: &p.Config.Database.Custom,
		store:    p.Store,

		colleagueHydrator: p.ColleagueHydrator,

		fHandler:             conductFileHandler,
		groupLogoFileHandler: groupLogoFileHandler,
		groupAccess:          p.GroupAccess,
		groupAccessResolver:  access.NewSubjectResolver(p.DB),
		qualificationAccess:  p.QualificationAccess,

		userSel: p.UserSel,
	}
	if s.store == nil {
		r := jobsstore.New(p.DB, &p.Config.Database.Custom, p.GroupAccess)
		s.store = r.Store
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbjobs.RegisterConductServiceServer(srv, s)
	pbjobs.RegisterColleaguesServiceServer(srv, s)
	pbjobs.RegisterJobsServiceServer(srv, s)
	pbjobs.RegisterStatsServiceServer(srv, s)
	pbjobs.RegisterTimeclockServiceServer(srv, s)
	pbjobs.RegisterGroupsServiceServer(srv, s)
}
