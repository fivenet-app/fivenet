package qualifications

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetQualifications,
		TimestampColumn: table.FivenetQualifications.DeletedAt,
		MinDays:         30,
	},
		&housekeeper.Table{
			Table:           table.FivenetQualificationsExamUsers,
			TimestampColumn: table.FivenetQualificationsExamUsers.EndsAt,
			MinDays:         30,
		},
		&housekeeper.Table{
			Table:           table.FivenetQualificationsRequests,
			TimestampColumn: table.FivenetQualificationsRequests.DeletedAt,
			MinDays:         30,
		},
		&housekeeper.Table{
			Table:           table.FivenetQualificationsResults,
			TimestampColumn: table.FivenetQualificationsResults.DeletedAt,
			MinDays:         30,
		},
	)

	housekeeper.AddJobTable(
		&housekeeper.JobTable{
			TargetTable:           table.FivenetQualifications,
			TargetJobColumn:       table.FivenetQualifications.Job,
			TargetSourceIDColumn:  table.FivenetQualifications.ID,
			TargetDeletedAtColumn: table.FivenetQualifications.DeletedAt,
		},
		&housekeeper.JobTable{
			Source: &housekeeper.JobTableSource{
				SourceTable:           table.FivenetQualifications,
				SourceJobColumn:       table.FivenetQualifications.Job,
				SourceDeletedAtColumn: table.FivenetQualifications.DeletedAt,
				SourceIDColumn:        table.FivenetQualifications.ID,
			},
			TargetTable:          table.FivenetQualificationsExamUsers,
			TargetSourceIDColumn: table.FivenetQualificationsExamUsers.QualificationID,
		},
		&housekeeper.JobTable{
			Source: &housekeeper.JobTableSource{
				SourceTable:           table.FivenetQualifications,
				SourceJobColumn:       table.FivenetQualifications.Job,
				SourceDeletedAtColumn: table.FivenetQualifications.DeletedAt,
				SourceIDColumn:        table.FivenetQualifications.ID,
			},
			TargetTable:           table.FivenetQualificationsRequests,
			TargetSourceIDColumn:  table.FivenetQualificationsRequests.QualificationID,
			TargetDeletedAtColumn: table.FivenetQualificationsRequests.DeletedAt,
		},
		&housekeeper.JobTable{
			Source: &housekeeper.JobTableSource{
				SourceTable:           table.FivenetQualifications,
				SourceJobColumn:       table.FivenetQualifications.Job,
				SourceDeletedAtColumn: table.FivenetQualifications.DeletedAt,
				SourceIDColumn:        table.FivenetQualifications.ID,
			},
			TargetTable:           table.FivenetQualificationsResults,
			TargetSourceIDColumn:  table.FivenetQualificationsResults.QualificationID,
			TargetDeletedAtColumn: table.FivenetQualificationsResults.DeletedAt,
		},
	)
}

var tQuali = table.FivenetQualifications.AS("qualification")

type Server struct {
	pbqualifications.QualificationsServiceServer

	logger   *zap.Logger
	db       *sql.DB
	perms    perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	notif    notifi.INotifi
	st       storage.IStorage

	access *access.Grouped[qualifications.QualificationJobAccess, *qualifications.QualificationJobAccess, qualifications.QualificationUserAccess, *qualifications.QualificationUserAccess, access.DummyQualificationAccess[qualifications.AccessLevel], *access.DummyQualificationAccess[qualifications.AccessLevel], qualifications.AccessLevel]
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
	Notif             notifi.INotifi
	Storage           storage.IStorage
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("jobs"),

		db:       p.DB,
		perms:    p.Perms,
		enricher: p.UserAwareEnricher,
		aud:      p.Audit,
		notif:    p.Notif,
		st:       p.Storage,

		access: access.NewGrouped[qualifications.QualificationJobAccess, *qualifications.QualificationJobAccess, qualifications.QualificationUserAccess, *qualifications.QualificationUserAccess, access.DummyQualificationAccess[qualifications.AccessLevel], *access.DummyQualificationAccess[qualifications.AccessLevel], qualifications.AccessLevel](
			p.DB,
			table.FivenetQualifications,
			&access.TargetTableColumns{
				ID:         table.FivenetQualifications.ID,
				DeletedAt:  table.FivenetQualifications.DeletedAt,
				CreatorID:  table.FivenetQualifications.CreatorID,
				CreatorJob: table.FivenetQualifications.CreatorJob,
			},
			access.NewJobs[qualifications.QualificationJobAccess, *qualifications.QualificationJobAccess, qualifications.AccessLevel](
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
						ID:       table.FivenetQualificationsAccess.AS("qualification_job_access").ID,
						TargetID: table.FivenetQualificationsAccess.AS("qualification_job_access").TargetID,
						Access:   table.FivenetQualificationsAccess.AS("qualification_job_access").Access,
					},
					Job:          table.FivenetQualificationsAccess.AS("qualification_job_access").Job,
					MinimumGrade: table.FivenetQualificationsAccess.AS("qualification_job_access").MinimumGrade,
				},
			),
			nil,
			nil,
		),
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbqualifications.RegisterQualificationsServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbqualifications.PermsRemap
}
