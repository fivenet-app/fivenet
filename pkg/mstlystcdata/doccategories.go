package mstlystcdata

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/cache"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// DocumentCategories manages document category data, including caching and enrichment.
type DocumentCategories struct {
	// Cache provides a cache for document categories keyed by string.
	*cache.Cache[documents.Category, *documents.Category]

	// logger is used for logging within the DocumentCategories service.
	logger *zap.Logger
	// db is the SQL database connection used for queries.
	db *sql.DB

	// tracer is used for distributed tracing of operations.
	tracer trace.Tracer
}

// DocumentCategoriesResult is the result struct for dependency injection, providing
// both the DocumentCategories service and a cron job registration.
type DocumentCategoriesResult struct {
	fx.Out

	// DocumentCategories is the main service for document category management.
	DocumentCategories *DocumentCategories
	// CronRegister registers cron jobs for document category refresh.
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

// NewDocumentCategories constructs a new DocumentCategories service, sets up cache,
// and registers lifecycle hooks for startup and shutdown.
func NewDocumentCategories(p Params) DocumentCategoriesResult {
	ctxCancel, cancel := context.WithCancel(context.Background())

	c := &DocumentCategories{
		logger: p.Logger,
		db:     p.DB,

		tracer: p.TP.Tracer("mstlystcdata.doccategories"),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		// Initialize the document categories cache with a key-value prefix.
		docCategories, err := cache.New(ctxStartup, p.Logger, p.JS, "cache",
			cache.WithKVPrefix[documents.Category]("doc_categories"),
		)
		if err != nil {
			return err
		}
		c.Cache = docCategories

		// Start the cache and load categories into it.
		if err := docCategories.Start(ctxCancel, false); err != nil {
			return err
		}

		if err := c.loadCategories(ctxStartup); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		// Cancel the context on shutdown.
		cancel()

		return nil
	}))

	return DocumentCategoriesResult{
		DocumentCategories: c,
		CronRegister:       c,
	}
}

// RegisterCronjobs registers the cron job for refreshing document categories.
func (c *DocumentCategories) RegisterCronjobs(
	ctx context.Context,
	registry croner.IRegistry,
) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "mstlystcdata.doccategories",
		Schedule: "* * * * *", // Every minute
	}); err != nil {
		return err
	}

	return nil
}

// RegisterCronjobHandlers registers the handler for the document categories cron job.
func (c *DocumentCategories) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add("mstlystcdata.doccategories", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := c.tracer.Start(ctx, "mstlystcdata-doccategories")
		defer span.End()

		if err := c.loadCategories(ctx); err != nil {
			c.logger.Error("failed to refresh doccategories cache", zap.Error(err))
			return err
		}

		return nil
	})

	return nil
}

// loadCategories loads all document categories from the database into the cache.
// It also builds a map of categories per job for potential future use.
func (c *DocumentCategories) loadCategories(ctx context.Context) error {
	tDCategory := table.FivenetDocumentsCategories.AS("category")

	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		FROM(tDCategory).
		WHERE(
			tDCategory.DeletedAt.IS_NULL(),
		).
		ORDER_BY(
			tDCategory.Job.ASC(),
			tDCategory.SortKey.ASC(),
		)

	var dest []*documents.Category
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	errs := multierr.Combine()
	categoriesPerJob := map[string][]*documents.Category{}
	for _, d := range dest {
		key := strconv.FormatInt(d.GetId(), 10)
		if err := c.Put(ctx, key, d); err != nil {
			errs = multierr.Append(errs, err)
		}

		if _, ok := categoriesPerJob[d.GetJob()]; !ok {
			categoriesPerJob[d.GetJob()] = []*documents.Category{}
		}
		categoriesPerJob[d.GetJob()] = append(categoriesPerJob[d.GetJob()], d)
	}

	return errs
}

// Enrich sets the full category object on the given ICategory if available in the cache.
// If the category is not found, it sets a placeholder category.
func (c *DocumentCategories) Enrich(doc common.ICategory) {
	cId := doc.GetCategoryId()

	// No category
	if cId == 0 {
		return
	}

	dc, err := c.Get(strconv.FormatInt(cId, 10))
	if err != nil {
		job := NotAvailablePlaceholder
		doc.SetCategory(&documents.Category{
			Id:   0,
			Name: NotAvailablePlaceholder,
			Job:  &job,
		})
	} else {
		doc.SetCategory(dc)
	}
}
