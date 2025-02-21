package mstlystcdata

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type DocumentCategories struct {
	logger *zap.Logger
	db     *sql.DB

	tracer trace.Tracer

	store *store.Store[documents.Category, *documents.Category]
	store.StoreRO[documents.Category, *documents.Category]
}

func NewDocumentCategories(p Params) *DocumentCategories {
	ctxCancel, cancel := context.WithCancel(context.Background())

	c := &DocumentCategories{
		logger: p.Logger,
		db:     p.DB,

		tracer: p.TP.Tracer("mstlystcdata-doccategories"),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		docCategories, err := store.New(ctxStartup, p.Logger, p.JS, "cache",
			store.WithLocks[documents.Category](nil),
			store.WithKVPrefix[documents.Category]("doc_categories"),
		)
		if err != nil {
			return err
		}
		c.store = docCategories
		c.StoreRO = docCategories

		if err := docCategories.Start(ctxCancel, false); err != nil {
			return err
		}

		if err := c.loadCategories(ctxStartup); err != nil {
			return err
		}

		p.CronHandlers.Add("mstlystcdata.doccategories", func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := c.tracer.Start(ctx, "mstlystcdata-doccategories")
			defer span.End()

			if err := c.loadCategories(ctx); err != nil {
				c.logger.Error("failed to refresh doccategories cache", zap.Error(err))
				return err
			}

			return nil
		})

		if err := p.Cron.RegisterCronjob(ctxStartup, &cron.Cronjob{
			Name:     "mstlystcdata.doccategories",
			Schedule: "@always", // Every minute
		}); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return c
}

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
		ORDER_BY(
			tDCategory.Job.ASC(),
			tDCategory.Name.ASC(),
		)

	var dest []*documents.Category
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		return err
	}

	errs := multierr.Combine()
	categoriesPerJob := map[string][]*documents.Category{}
	for _, d := range dest {
		key := strconv.FormatUint(d.Id, 10)
		if err := c.store.Put(ctx, key, d); err != nil {
			errs = multierr.Append(errs, err)
		}

		if _, ok := categoriesPerJob[*d.Job]; !ok {
			categoriesPerJob[*d.Job] = []*documents.Category{}
		}
		categoriesPerJob[*d.Job] = append(categoriesPerJob[*d.Job], d)
	}

	return errs
}
