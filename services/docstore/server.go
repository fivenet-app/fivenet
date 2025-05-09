package docstore

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	pbdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/html/htmldiffer"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

const housekeeperMinDays = 60

func init() {
	housekeeper.AddTable(
		// Categories
		&housekeeper.Table{
			Table:           table.FivenetDocumentsCategories,
			TimestampColumn: table.FivenetDocumentsCategories.DeletedAt,
			MinDays:         housekeeperMinDays,
		},

		// Templates
		&housekeeper.Table{
			Table:           table.FivenetDocumentsTemplates,
			TimestampColumn: table.FivenetDocumentsTemplates.DeletedAt,
			MinDays:         housekeeperMinDays,
		},

		// Documents
		&housekeeper.Table{
			Table:           table.FivenetDocuments,
			TimestampColumn: table.FivenetDocuments.DeletedAt,
			MinDays:         housekeeperMinDays,
		},

		// Comments
		&housekeeper.Table{
			Table:           table.FivenetDocumentsComments,
			TimestampColumn: table.FivenetDocumentsComments.DeletedAt,
			MinDays:         housekeeperMinDays,
		},

		// Document References and Relations
		&housekeeper.Table{
			Table:           table.FivenetDocumentsReferences,
			TimestampColumn: table.FivenetDocumentsReferences.DeletedAt,
			MinDays:         housekeeperMinDays,
		},
		&housekeeper.Table{
			Table:           table.FivenetDocumentsRelations,
			TimestampColumn: table.FivenetDocumentsRelations.DeletedAt,
			MinDays:         housekeeperMinDays,
		},
	)
}

type Server struct {
	pbdocstore.DocStoreServiceServer

	db            *sql.DB
	ps            perms.Permissions
	jobs          *mstlystcdata.Jobs
	docCategories *mstlystcdata.DocumentCategories
	enricher      *mstlystcdata.UserAwareEnricher
	aud           audit.IAuditer
	ui            userinfo.UserInfoRetriever
	notif         notifi.INotifi
	htmlDiff      *htmldiffer.Differ

	access         *access.Grouped[documents.DocumentJobAccess, *documents.DocumentJobAccess, documents.DocumentUserAccess, *documents.DocumentUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel]
	templateAccess *access.Grouped[documents.TemplateJobAccess, *documents.TemplateJobAccess, documents.TemplateUserAccess, *documents.TemplateUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel]
}

type Params struct {
	fx.In

	DB            *sql.DB
	Perms         perms.Permissions
	Jobs          *mstlystcdata.Jobs
	DocCategories *mstlystcdata.DocumentCategories
	Enricher      *mstlystcdata.UserAwareEnricher
	Aud           audit.IAuditer
	Ui            userinfo.UserInfoRetriever
	Notif         notifi.INotifi
	HTMLDiffer    *htmldiffer.Differ
}

func NewServer(p Params) *Server {
	return &Server{
		db:            p.DB,
		ps:            p.Perms,
		jobs:          p.Jobs,
		docCategories: p.DocCategories,
		enricher:      p.Enricher,
		aud:           p.Aud,
		ui:            p.Ui,
		notif:         p.Notif,
		htmlDiff:      p.HTMLDiffer,

		access: newAccess(p.DB),

		templateAccess: access.NewGrouped[documents.TemplateJobAccess, *documents.TemplateJobAccess, documents.TemplateUserAccess, *documents.TemplateUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel](
			p.DB,
			table.FivenetDocumentsTemplates,
			&access.TargetTableColumns{
				ID:         table.FivenetDocumentsTemplates.ID,
				DeletedAt:  table.FivenetDocumentsTemplates.DeletedAt,
				CreatorID:  nil,
				CreatorJob: table.FivenetDocumentsTemplates.CreatorJob,
			},
			access.NewJobs[documents.TemplateJobAccess, *documents.TemplateJobAccess, documents.AccessLevel](
				table.FivenetDocumentsTemplatesAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetDocumentsTemplatesAccess.ID,
						TargetID: table.FivenetDocumentsTemplatesAccess.TargetID,
						Access:   table.FivenetDocumentsTemplatesAccess.Access,
					},
					Job:          table.FivenetDocumentsTemplatesAccess.Job,
					MinimumGrade: table.FivenetDocumentsTemplatesAccess.MinimumGrade,
				},
				table.FivenetDocumentsTemplatesAccess.AS("template_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetDocumentsTemplatesAccess.AS("template_job_access").ID,
						TargetID: table.FivenetDocumentsTemplatesAccess.AS("template_job_access").TargetID,
						Access:   table.FivenetDocumentsTemplatesAccess.AS("template_job_access").Access,
					},
					Job:          table.FivenetDocumentsTemplatesAccess.AS("template_job_access").Job,
					MinimumGrade: table.FivenetDocumentsTemplatesAccess.AS("template_job_access").MinimumGrade,
				},
			),
			nil,
			nil,
		),
	}
}

func newAccess(db *sql.DB) *access.Grouped[documents.DocumentJobAccess, *documents.DocumentJobAccess, documents.DocumentUserAccess, *documents.DocumentUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel] {
	return access.NewGrouped[documents.DocumentJobAccess, *documents.DocumentJobAccess, documents.DocumentUserAccess, *documents.DocumentUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel](
		db,
		table.FivenetDocuments,
		&access.TargetTableColumns{
			ID:         table.FivenetDocuments.ID,
			DeletedAt:  table.FivenetDocuments.DeletedAt,
			CreatorID:  table.FivenetDocuments.CreatorID,
			CreatorJob: table.FivenetDocuments.CreatorJob,
		},
		access.NewJobs[documents.DocumentJobAccess, *documents.DocumentJobAccess, documents.AccessLevel](
			table.FivenetDocumentsAccess,
			&access.JobAccessColumns{
				BaseAccessColumns: access.BaseAccessColumns{
					ID:       table.FivenetDocumentsAccess.ID,
					TargetID: table.FivenetDocumentsAccess.TargetID,
					Access:   table.FivenetDocumentsAccess.Access,
				},
				Job:          table.FivenetDocumentsAccess.Job,
				MinimumGrade: table.FivenetDocumentsAccess.MinimumGrade,
			},
			table.FivenetDocumentsAccess.AS("document_job_access"),
			&access.JobAccessColumns{
				BaseAccessColumns: access.BaseAccessColumns{
					ID:       table.FivenetDocumentsAccess.AS("document_job_access").ID,
					TargetID: table.FivenetDocumentsAccess.AS("document_job_access").TargetID,
					Access:   table.FivenetDocumentsAccess.AS("document_job_access").Access,
				},
				Job:          table.FivenetDocumentsAccess.AS("document_job_access").Job,
				MinimumGrade: table.FivenetDocumentsAccess.AS("document_job_access").MinimumGrade,
			},
		),
		access.NewUsers[documents.DocumentUserAccess, *documents.DocumentUserAccess, documents.AccessLevel](
			table.FivenetDocumentsAccess,
			&access.UserAccessColumns{
				BaseAccessColumns: access.BaseAccessColumns{
					ID:       table.FivenetDocumentsAccess.ID,
					TargetID: table.FivenetDocumentsAccess.TargetID,
					Access:   table.FivenetDocumentsAccess.Access,
				},
				UserId: table.FivenetDocumentsAccess.UserID,
			},
			table.FivenetDocumentsAccess.AS("document_user_access"),
			&access.UserAccessColumns{
				BaseAccessColumns: access.BaseAccessColumns{
					ID:       table.FivenetDocumentsAccess.AS("document_user_access").ID,
					TargetID: table.FivenetDocumentsAccess.AS("document_user_access").TargetID,
					Access:   table.FivenetDocumentsAccess.AS("document_user_access").Access,
				},
				UserId: table.FivenetDocumentsAccess.AS("document_user_access").UserID,
			},
		),
		nil,
	)
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbdocstore.RegisterDocStoreServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbdocstore.PermsRemap
}
