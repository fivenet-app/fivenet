package documents

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsstamps "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/stamps"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/collab"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	pkggrpc "github.com/fivenet-app/fivenet/v2026/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	docstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
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
	pbdocuments.StampsServiceServer
	pbdocuments.StatsServiceServer

	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB

	js            *events.JSWrapper
	ps            perms.Permissions
	jobs          *mstlystcdata.Jobs
	docCategories *mstlystcdata.DocumentCategories
	enricher      *mstlystcdata.UserAwareEnricher
	ui            userinfo.UserInfoRetriever
	notifi        notifi.INotifi

	access         *access.Grouped[documentsaccess.DocumentJobAccess, *documentsaccess.DocumentJobAccess, documentsaccess.DocumentUserAccess, *documentsaccess.DocumentUserAccess, access.DummyQualificationAccess[documentsaccess.AccessLevel], *access.DummyQualificationAccess[documentsaccess.AccessLevel], documentsaccess.AccessLevel]
	templateAccess *access.Grouped[documentstemplates.TemplateJobAccess, *documentstemplates.TemplateJobAccess, documentstemplates.TemplateUserAccess, *documentstemplates.TemplateUserAccess, access.DummyQualificationAccess[documentsaccess.AccessLevel], *access.DummyQualificationAccess[documentsaccess.AccessLevel], documentsaccess.AccessLevel]

	signingStampAccess *access.Grouped[documentsstamps.StampJobAccess, *documentsstamps.StampJobAccess, access.DummyUserAccess[documentsstamps.StampAccessLevel], *access.DummyUserAccess[documentsstamps.StampAccessLevel], access.DummyQualificationAccess[documentsstamps.StampAccessLevel], *access.DummyQualificationAccess[documentsstamps.StampAccessLevel], documentsstamps.StampAccessLevel]

	collabServer *collab.CollabServer
	fHandler     *filestore.Handler[int64]
	stats        *docstats.Service
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger        *zap.Logger
	DB            *sql.DB
	TP            *tracesdk.TracerProvider
	Perms         perms.Permissions
	Storage       storage.IStorage
	Jobs          *mstlystcdata.Jobs
	DocCategories *mstlystcdata.DocumentCategories
	Enricher      *mstlystcdata.UserAwareEnricher
	Ui            userinfo.UserInfoRetriever
	Notif         notifi.INotifi
	JS            *events.JSWrapper
	Stats         *docstats.Service
}

type Result struct {
	fx.Out

	Server       *Server
	Service      pkggrpc.Service     `group:"grpcservices"`
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewServer(p Params) Result {
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
				documentsaccess.AccessLevel(access),
			)
		},
	})

	s := &Server{
		logger: p.Logger.Named("documents"),
		tracer: p.TP.Tracer("documents"),
		db:     p.DB,

		js:            p.JS,
		ps:            p.Perms,
		jobs:          p.Jobs,
		docCategories: p.DocCategories,
		enricher:      p.Enricher,
		ui:            p.Ui,
		notifi:        p.Notif,

		access: docAccess,
		templateAccess: access.NewGrouped[documentstemplates.TemplateJobAccess, *documentstemplates.TemplateJobAccess, documentstemplates.TemplateUserAccess, *documentstemplates.TemplateUserAccess, access.DummyQualificationAccess[documentsaccess.AccessLevel], *access.DummyQualificationAccess[documentsaccess.AccessLevel], documentsaccess.AccessLevel](
			p.DB,
			table.FivenetDocumentsTemplates,
			&access.TargetTableColumns{
				ID:         table.FivenetDocumentsTemplates.ID,
				DeletedAt:  table.FivenetDocumentsTemplates.DeletedAt,
				CreatorID:  nil,
				CreatorJob: table.FivenetDocumentsTemplates.CreatorJob,
			},
			access.NewJobs[documentstemplates.TemplateJobAccess, *documentstemplates.TemplateJobAccess, documentsaccess.AccessLevel](
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

		signingStampAccess: access.NewGrouped[documentsstamps.StampJobAccess, *documentsstamps.StampJobAccess, access.DummyUserAccess[documentsstamps.StampAccessLevel], *access.DummyUserAccess[documentsstamps.StampAccessLevel], access.DummyQualificationAccess[documentsstamps.StampAccessLevel], *access.DummyQualificationAccess[documentsstamps.StampAccessLevel], documentsstamps.StampAccessLevel](
			p.DB,
			table.FivenetDocumentsStampsAccess,
			&access.TargetTableColumns{
				ID:         table.FivenetDocumentsStamps.ID,
				DeletedAt:  table.FivenetDocumentsStamps.DeletedAt,
				CreatorID:  nil,
				CreatorJob: nil,
			},
			access.NewJobs[documentsstamps.StampJobAccess, *documentsstamps.StampJobAccess, documentsstamps.StampAccessLevel](
				table.FivenetDocumentsStampsAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetDocumentsStampsAccess.ID,
						TargetID: table.FivenetDocumentsStampsAccess.TargetID,
						Access:   table.FivenetDocumentsStampsAccess.Access,
					},
					Job:          table.FivenetDocumentsStampsAccess.Job,
					MinimumGrade: table.FivenetDocumentsStampsAccess.MinimumGrade,
				},
				table.FivenetDocumentsStampsAccess.AS("stamp_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID: table.FivenetDocumentsStampsAccess.AS(
							"stamp_job_access",
						).ID,
						TargetID: table.FivenetDocumentsStampsAccess.AS(
							"stamp_job_access",
						).TargetID,
						Access: table.FivenetDocumentsStampsAccess.AS(
							"stamp_job_access",
						).Access,
					},
					Job: table.FivenetDocumentsStampsAccess.AS(
						"stamp_job_access",
					).Job,
					MinimumGrade: table.FivenetDocumentsStampsAccess.AS(
						"stamp_job_access",
					).MinimumGrade,
				},
			),
			nil,
			nil,
		),

		collabServer: collabServer,
		fHandler:     fHandler,
		stats:        p.Stats,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		return s.collabServer.Start(ctxStartup)
	}))

	p.LC.Append(fx.StopHook(func(ctxStartup context.Context) error {
		cancel()

		return nil
	}))

	return Result{
		Server:       s,
		Service:      s,
		CronRegister: s,
	}
}

func newAccess(
	db *sql.DB,
) *access.Grouped[
	documentsaccess.DocumentJobAccess, *documentsaccess.DocumentJobAccess,
	documentsaccess.DocumentUserAccess, *documentsaccess.DocumentUserAccess,
	access.DummyQualificationAccess[documentsaccess.AccessLevel], *access.DummyQualificationAccess[documentsaccess.AccessLevel],
	documentsaccess.AccessLevel,
] {
	return access.NewGrouped[documentsaccess.DocumentJobAccess, *documentsaccess.DocumentJobAccess, documentsaccess.DocumentUserAccess, *documentsaccess.DocumentUserAccess, access.DummyQualificationAccess[documentsaccess.AccessLevel], *access.DummyQualificationAccess[documentsaccess.AccessLevel], documentsaccess.AccessLevel](
		db,
		table.FivenetDocuments,
		&access.TargetTableColumns{
			ID:         table.FivenetDocuments.ID,
			DeletedAt:  table.FivenetDocuments.DeletedAt,
			CreatorID:  table.FivenetDocuments.CreatorID,
			CreatorJob: table.FivenetDocuments.CreatorJob,
		},
		access.NewJobs[documentsaccess.DocumentJobAccess, *documentsaccess.DocumentJobAccess, documentsaccess.AccessLevel](
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
		access.NewUsers[documentsaccess.DocumentUserAccess, *documentsaccess.DocumentUserAccess, documentsaccess.AccessLevel](
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
	pbdocuments.RegisterStampsServiceServer(srv, s)
	pbdocuments.RegisterStatsServiceServer(srv, s)
}

func (s *Server) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.approval.tasks.expire",
		Schedule: "* * * * *", // Every minute
	}); err != nil {
		return err
	}

	if err := registry.UnregisterCronjob(ctx, "documents.signature.tasks.expire"); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.stats.rollup.columns.recent",
		Schedule: "*/5 * * * *",
	}); err != nil {
		return err
	}
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.stats.rollup.columns.backfill",
		Schedule: "17 3 * * *",
	}); err != nil {
		return err
	}
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.stats.rollup.metrics.recent",
		Schedule: "*/5 * * * *",
	}); err != nil {
		return err
	}
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.stats.rollup.metrics.backfill",
		Schedule: "22 3 * * *",
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add(
		"documents.approval.tasks.expire",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "documents.approval.tasks.expire")
			defer span.End()

			dest := &cron.GenericCronData{}

			rowsAffected, err := s.expireApprovalTasks(ctx)
			if err != nil {
				return err
			}
			dest.SetAttribute("rows_affected", strconv.FormatInt(rowsAffected, 10))

			// Marshal the updated cron data
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal updated document workflow cron data. %w", err)
			}

			return nil
		},
	)

	h.Add(
		"documents.stats.rollup.columns.recent",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "documents.stats.rollup.columns.recent")
			defer span.End()

			startDay := time.Now().UTC().Truncate(24 * time.Hour)
			if err := s.stats.RebuildDocumentColumnRollups(ctx, startDay, startDay); err != nil {
				return err
			}

			dest := &cron.GenericCronData{}
			dest.SetAttribute("start_day", startDay.Format(time.DateOnly))
			dest.SetAttribute("kind", "recent")
			dest.SetAttribute("type", "document_column")

			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal updated document stats cron data. %w", err)
			}

			return nil
		},
	)

	h.Add(
		"documents.stats.rollup.columns.backfill",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "documents.stats.rollup.columns.backfill")
			defer span.End()

			endDay := time.Now().UTC().Truncate(24 * time.Hour)
			startDay := endDay.AddDate(0, 0, -30)
			if err := s.stats.RebuildDocumentColumnRollups(ctx, startDay, endDay); err != nil {
				return err
			}

			dest := &cron.GenericCronData{}
			dest.SetAttribute("start_day", startDay.Format(time.DateOnly))
			dest.SetAttribute("end_day", endDay.Format(time.DateOnly))
			dest.SetAttribute("kind", "backfill")
			dest.SetAttribute("type", "document_column")

			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal updated document stats cron data. %w", err)
			}

			return nil
		},
	)

	h.Add(
		"documents.stats.rollup.metrics.recent",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "documents.stats.rollup.metrics.recent")
			defer span.End()

			startDay := time.Now().UTC().Truncate(24 * time.Hour)
			if err := s.stats.RebuildDocumentMetricRollups(ctx, startDay, startDay); err != nil {
				return err
			}

			dest := &cron.GenericCronData{}
			dest.SetAttribute("start_day", startDay.Format(time.DateOnly))
			dest.SetAttribute("kind", "recent")
			dest.SetAttribute("type", "document_metric")

			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal updated document stats cron data. %w", err)
			}

			return nil
		},
	)

	h.Add(
		"documents.stats.rollup.metrics.backfill",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "documents.stats.rollup.metrics.backfill")
			defer span.End()

			endDay := time.Now().UTC().Truncate(24 * time.Hour)
			startDay := endDay.AddDate(0, 0, -30)
			if err := s.stats.RebuildDocumentMetricRollups(ctx, startDay, endDay); err != nil {
				return err
			}

			dest := &cron.GenericCronData{}
			dest.SetAttribute("start_day", startDay.Format(time.DateOnly))
			dest.SetAttribute("end_day", endDay.Format(time.DateOnly))
			dest.SetAttribute("kind", "backfill")
			dest.SetAttribute("type", "document_metric")

			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal updated document stats cron data. %w", err)
			}

			return nil
		},
	)

	return nil
}
