package documents

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/collab"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmldiffer"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

const housekeeperMinDays = 60

func init() {
	// Documents
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetDocuments,
		IDColumn:        table.FivenetDocuments.ID,
		DeletedAtColumn: table.FivenetDocuments.DeletedAt,

		MinDays: housekeeperMinDays,

		DependantTables: []*housekeeper.Table{
			// Activity
			{
				Table:      table.FivenetDocumentsActivity,
				IDColumn:   table.FivenetDocumentsActivity.ID,
				ForeignKey: table.FivenetDocumentsActivity.DocumentID,
			},

			// Comments
			{
				Table:           table.FivenetDocumentsComments,
				IDColumn:        table.FivenetDocumentsComments.ID,
				ForeignKey:      table.FivenetDocumentsComments.DocumentID,
				DeletedAtColumn: table.FivenetDocumentsComments.DeletedAt,

				MinDays: housekeeperMinDays,
			},

			// Document References and Relations
			{
				Table:           table.FivenetDocumentsReferences,
				IDColumn:        table.FivenetDocumentsReferences.ID,
				ForeignKey:      table.FivenetDocumentsReferences.SourceDocumentID,
				DeletedAtColumn: table.FivenetDocumentsReferences.DeletedAt,

				MinDays: housekeeperMinDays,
			},
			{
				Table:           table.FivenetDocumentsRelations,
				IDColumn:        table.FivenetDocumentsRelations.ID,
				ForeignKey:      table.FivenetDocumentsRelations.DocumentID,
				DeletedAtColumn: table.FivenetDocumentsRelations.DeletedAt,

				MinDays: housekeeperMinDays,
			},

			// Pins
			{
				Table:      table.FivenetDocumentsPins,
				ForeignKey: table.FivenetDocumentsPins.DocumentID,
			},
		},
	})

	// Categories
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetDocumentsCategories,
		IDColumn:        table.FivenetDocumentsCategories.ID,
		JobColumn:       table.FivenetDocumentsCategories.Job,
		DeletedAtColumn: table.FivenetDocumentsCategories.DeletedAt,

		MinDays: housekeeperMinDays,
	})
	// Templates
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetDocumentsTemplates,
		IDColumn:        table.FivenetDocumentsTemplates.ID,
		JobColumn:       table.FivenetDocumentsTemplates.CreatorJob,
		DeletedAtColumn: table.FivenetDocumentsTemplates.DeletedAt,

		MinDays: housekeeperMinDays,
	})
}

type Server struct {
	pbdocuments.DocumentsServiceServer
	pbdocuments.CollabServiceServer
	pbdocuments.ApprovalServiceServer
	pbdocuments.SigningServiceServer

	logger *zap.Logger
	db     *sql.DB

	js            *events.JSWrapper
	ps            perms.Permissions
	jobs          *mstlystcdata.Jobs
	docCategories *mstlystcdata.DocumentCategories
	enricher      *mstlystcdata.UserAwareEnricher
	aud           audit.IAuditer
	ui            userinfo.UserInfoRetriever
	notifi        notifi.INotifi
	htmlDiff      *htmldiffer.Differ

	access         *access.Grouped[documents.DocumentJobAccess, *documents.DocumentJobAccess, documents.DocumentUserAccess, *documents.DocumentUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel]
	templateAccess *access.Grouped[documents.TemplateJobAccess, *documents.TemplateJobAccess, documents.TemplateUserAccess, *documents.TemplateUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel]

	signingStampAccess *access.Grouped[documents.StampJobAccess, *documents.StampJobAccess, access.DummyUserAccess[documents.StampAccessLevel], *access.DummyUserAccess[documents.StampAccessLevel], access.DummyQualificationAccess[documents.StampAccessLevel], *access.DummyQualificationAccess[documents.StampAccessLevel], documents.StampAccessLevel]

	collabServer *collab.CollabServer
	fHandler     *filestore.Handler[int64]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger        *zap.Logger
	DB            *sql.DB
	Perms         perms.Permissions
	Storage       storage.IStorage
	Jobs          *mstlystcdata.Jobs
	DocCategories *mstlystcdata.DocumentCategories
	Enricher      *mstlystcdata.UserAwareEnricher
	Aud           audit.IAuditer
	Ui            userinfo.UserInfoRetriever
	Notif         notifi.INotifi
	HTMLDiffer    *htmldiffer.Differ
	JS            *events.JSWrapper
}

func NewServer(p Params) *Server {
	ctxCancel, cancel := context.WithCancel(context.Background())

	collabServer := collab.New(ctxCancel, p.Logger, p.JS, "documents")

	tDocFiles := table.FivenetDocumentsFiles

	// 3 MiB limit
	fHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tDocFiles,
		tDocFiles.DocumentID,
		tDocFiles.FileID,
		3<<20,
		func(parentId int64) mysql.BoolExpression {
			return tDocFiles.DocumentID.EQ(mysql.Int64(parentId))
		},
		filestore.InsertJoinRow,
		false,
	)

	docAccess := newAccess(p.DB)
	access.RegisterAccess("documents", &access.GroupedAccessAdapter{
		CanUserAccessTargetFn: func(ctx context.Context, targetId int64, userInfo *pbuserinfo.UserInfo, access int32) (bool, error) {
			// Type assert access to the correct enum type
			return docAccess.CanUserAccessTarget(
				ctx,
				targetId,
				userInfo,
				documents.AccessLevel(access),
			)
		},
	})

	s := &Server{
		logger: p.Logger.Named("documents"),
		db:     p.DB,

		js:            p.JS,
		ps:            p.Perms,
		jobs:          p.Jobs,
		docCategories: p.DocCategories,
		enricher:      p.Enricher,
		aud:           p.Aud,
		ui:            p.Ui,
		notifi:        p.Notif,
		htmlDiff:      p.HTMLDiffer,

		access: docAccess,
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
						ID: table.FivenetDocumentsTemplatesAccess.AS(
							"template_job_access",
						).ID,
						TargetID: table.FivenetDocumentsTemplatesAccess.AS(
							"template_job_access",
						).TargetID,
						Access: table.FivenetDocumentsTemplatesAccess.AS(
							"template_job_access",
						).Access,
					},
					Job: table.FivenetDocumentsTemplatesAccess.AS(
						"template_job_access",
					).Job,
					MinimumGrade: table.FivenetDocumentsTemplatesAccess.AS(
						"template_job_access",
					).MinimumGrade,
				},
			),
			nil,
			nil,
		),

		signingStampAccess: access.NewGrouped[documents.StampJobAccess, *documents.StampJobAccess, access.DummyUserAccess[documents.StampAccessLevel], *access.DummyUserAccess[documents.StampAccessLevel], access.DummyQualificationAccess[documents.StampAccessLevel], *access.DummyQualificationAccess[documents.StampAccessLevel], documents.StampAccessLevel](
			p.DB,
			table.FivenetDocumentsSignaturesStampsAccess,
			&access.TargetTableColumns{
				ID:         table.FivenetDocumentsSignaturesStamps.ID,
				DeletedAt:  table.FivenetDocumentsSignaturesStamps.DeletedAt,
				CreatorID:  table.FivenetDocumentsSignaturesStamps.OwnerID,
				CreatorJob: nil,
			},
			access.NewJobs[documents.StampJobAccess, *documents.StampJobAccess, documents.StampAccessLevel](
				table.FivenetDocumentsSignaturesStampsAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetDocumentsSignaturesStampsAccess.ID,
						TargetID: table.FivenetDocumentsSignaturesStampsAccess.TargetID,
						Access:   table.FivenetDocumentsSignaturesStampsAccess.Access,
					},
					Job:          table.FivenetDocumentsSignaturesStampsAccess.Job,
					MinimumGrade: table.FivenetDocumentsSignaturesStampsAccess.MinimumGrade,
				},
				table.FivenetDocumentsSignaturesStampsAccess.AS("stamp_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID: table.FivenetDocumentsSignaturesStampsAccess.AS(
							"stamp_job_access",
						).ID,
						TargetID: table.FivenetDocumentsSignaturesStampsAccess.AS(
							"stamp_job_access",
						).TargetID,
						Access: table.FivenetDocumentsSignaturesStampsAccess.AS(
							"stamp_job_access",
						).Access,
					},
					Job: table.FivenetDocumentsSignaturesStampsAccess.AS(
						"stamp_job_access",
					).Job,
					MinimumGrade: table.FivenetDocumentsSignaturesStampsAccess.AS(
						"stamp_job_access",
					).MinimumGrade,
				},
			),
			nil,
			nil,
		),

		collabServer: collabServer,
		fHandler:     fHandler,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		return s.collabServer.Start(ctxStartup)
	}))

	p.LC.Append(fx.StopHook(func(ctxStartup context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func newAccess(
	db *sql.DB,
) *access.Grouped[
	documents.DocumentJobAccess, *documents.DocumentJobAccess,
	documents.DocumentUserAccess, *documents.DocumentUserAccess,
	access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel],
	documents.AccessLevel,
] {
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
				UserID: table.FivenetDocumentsAccess.UserID,
			},
			table.FivenetDocumentsAccess.AS("document_user_access"),
			&access.UserAccessColumns{
				BaseAccessColumns: access.BaseAccessColumns{
					ID:       table.FivenetDocumentsAccess.AS("document_user_access").ID,
					TargetID: table.FivenetDocumentsAccess.AS("document_user_access").TargetID,
					Access:   table.FivenetDocumentsAccess.AS("document_user_access").Access,
				},
				UserID: table.FivenetDocumentsAccess.AS("document_user_access").UserID,
			},
		),
		nil,
	)
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbdocuments.RegisterDocumentsServiceServer(srv, s)
	pbdocuments.RegisterCollabServiceServer(srv, s)
	pbdocuments.RegisterApprovalServiceServer(srv, s)
	pbdocuments.RegisterSigningServiceServer(srv, s)
}

// GetPermsRemap returns the permissions re-mapping for the services.
func (s *Server) GetPermsRemap() map[string]string {
	return pbdocuments.PermsRemap
}
