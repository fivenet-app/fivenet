package qualifications

import (
	"database/sql"

	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
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
		Table:           table.FivenetQualifications,
		IDColumn:        table.FivenetQualifications.ID,
		JobColumn:       table.FivenetQualifications.Job,
		DeletedAtColumn: table.FivenetQualifications.DeletedAt,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:           table.FivenetQualificationsExamUsers,
				DeletedAtColumn: table.FivenetQualificationsExamUsers.EndsAt,
				ForeignKey:      table.FivenetQualificationsExamUsers.QualificationID,

				MinDays: 60,
			},
			{
				Table:           table.FivenetQualificationsRequests,
				DeletedAtColumn: table.FivenetQualificationsRequests.DeletedAt,
				ForeignKey:      table.FivenetQualificationsRequests.QualificationID,

				MinDays: 60,
			},
			{
				Table:           table.FivenetQualificationsResults,
				IDColumn:        table.FivenetQualificationsResults.ID,
				DeletedAtColumn: table.FivenetQualificationsResults.DeletedAt,
				ForeignKey:      table.FivenetQualificationsResults.QualificationID,

				MinDays: 60,
			},
		},
	})
}

var (
	tQuali      = table.FivenetQualifications.AS("qualification")
	tQualiFiles = table.FivenetQualificationsFiles
)

type Server struct {
	pbqualifications.QualificationsServiceServer

	logger   *zap.Logger
	db       *sql.DB
	perms    perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	notif    notifi.INotifi
	st       storage.IStorage

	access *access.Grouped[qualificationsaccess.QualificationJobAccess, *qualificationsaccess.QualificationJobAccess, qualificationsaccess.QualificationUserAccess, *qualificationsaccess.QualificationUserAccess, access.DummyQualificationAccess[qualificationsaccess.AccessLevel], *access.DummyQualificationAccess[qualificationsaccess.AccessLevel], qualificationsaccess.AccessLevel]

	fHandler *filestore.Handler[int64]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	DB                *sql.DB
	Perms             perms.Permissions
	UserAwareEnricher *mstlystcdata.UserAwareEnricher
	Config            *config.Config
	Notif             notifi.INotifi
	Storage           storage.IStorage
}

func NewServer(p Params) *Server {
	// 3 MiB limit
	qualiFileHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tQualiFiles,
		tQualiFiles.QualificationID,
		tQualiFiles.FileID,
		3<<20,
		func(parentId int64) mysql.BoolExpression {
			return tQualiFiles.QualificationID.EQ(mysql.Int64(parentId))
		},
		filestore.InsertJoinRow,
		false,
	)

	s := &Server{
		logger: p.Logger.Named("jobs"),

		db:       p.DB,
		perms:    p.Perms,
		enricher: p.UserAwareEnricher,
		notif:    p.Notif,
		st:       p.Storage,

		access: access.NewGrouped[qualificationsaccess.QualificationJobAccess, *qualificationsaccess.QualificationJobAccess, qualificationsaccess.QualificationUserAccess, *qualificationsaccess.QualificationUserAccess, access.DummyQualificationAccess[qualificationsaccess.AccessLevel], *access.DummyQualificationAccess[qualificationsaccess.AccessLevel], qualificationsaccess.AccessLevel](
			p.DB,
			table.FivenetQualifications,
			&access.TargetTableColumns{
				ID:         table.FivenetQualifications.ID,
				DeletedAt:  table.FivenetQualifications.DeletedAt,
				CreatorID:  table.FivenetQualifications.CreatorID,
				CreatorJob: table.FivenetQualifications.CreatorJob,
			},
			access.NewJobs[qualificationsaccess.QualificationJobAccess, *qualificationsaccess.QualificationJobAccess, qualificationsaccess.AccessLevel](
				table.FivenetQualificationsAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetQualificationsAccess.ID,
						TargetID: table.FivenetQualificationsAccess.TargetID,
						Access:   table.FivenetQualificationsAccess.Access,
					},
					Job:          table.FivenetQualificationsAccess.Job,
					MinimumGrade: table.FivenetQualificationsAccess.MinimumGrade,
				},
				table.FivenetQualificationsAccess.AS("qualification_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID: table.FivenetQualificationsAccess.AS(
							"qualification_job_access",
						).ID,
						TargetID: table.FivenetQualificationsAccess.AS(
							"qualification_job_access",
						).TargetID,
						Access: table.FivenetQualificationsAccess.AS(
							"qualification_job_access",
						).Access,
					},
					Job: table.FivenetQualificationsAccess.AS(
						"qualification_job_access",
					).Job,
					MinimumGrade: table.FivenetQualificationsAccess.AS(
						"qualification_job_access",
					).MinimumGrade,
				},
			),
			nil,
			nil,
		),

		fHandler: qualiFileHandler,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbqualifications.RegisterQualificationsServiceServer(srv, s)
}
