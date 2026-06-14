package documents

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
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
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
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
	pbdocuments.CategoriesServiceServer
	pbdocuments.CollabServiceServer
	pbdocuments.CommentsServiceServer
	pbdocuments.ApprovalServiceServer
	pbdocuments.StampsServiceServer
	pbdocuments.StatsServiceServer
	pbdocuments.TemplatesServiceServer

	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB

	js            *events.JSWrapper
	ps            perms.Permissions
	jobs          mstlystcdata.IJobs
	docCategories mstlystcdata.IDocumentCategories
	enricher      mstlystcdata.IUserAwareEnricher
	ui            userinfo.UserInfoRetriever
	notifi        notifi.INotifi
	store         documentsstore.IStore

	subjectAccess   *access.SubjectObjectAccess
	subjectResolver *access.SubjectResolver
	templateAccess  *access.SubjectObjectAccess

	signingStampAccess *access.SubjectObjectAccess

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
	Jobs          mstlystcdata.IJobs
	DocCategories mstlystcdata.IDocumentCategories
	Enricher      mstlystcdata.IUserAwareEnricher
	Ui            userinfo.UserInfoRetriever
	Notif         notifi.INotifi
	JS            *events.JSWrapper
	Stats         *docstats.Service
	Store         documentsstore.IStore
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
		5,
		func(parentId int64) mysql.BoolExpression {
			return tDocFiles.DocumentID.EQ(mysql.Int64(parentId))
		},
		filestore.InsertJoinRow,
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	docSubjectAccess := access.NewDocumentsSubjectObjectAccess(p.DB)
	docSubjectResolver := access.NewSubjectResolver(p.DB)
	access.RegisterAccess("documents", docSubjectAccess)

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
		store:         p.Store,

		subjectAccess:   docSubjectAccess,
		subjectResolver: docSubjectResolver,
		templateAccess:  access.NewDocumentTemplatesSubjectObjectAccess(p.DB),

		signingStampAccess: access.NewDocumentStampsSubjectObjectAccess(p.DB),

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

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbdocuments.RegisterDocumentsServiceServer(srv, s)
	pbdocuments.RegisterCategoriesServiceServer(srv, s)
	pbdocuments.RegisterCollabServiceServer(srv, s)
	pbdocuments.RegisterCommentsServiceServer(srv, s)
	pbdocuments.RegisterApprovalServiceServer(srv, s)
	pbdocuments.RegisterStampsServiceServer(srv, s)
	pbdocuments.RegisterStatsServiceServer(srv, s)
	pbdocuments.RegisterTemplatesServiceServer(srv, s)
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
