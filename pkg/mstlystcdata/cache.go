package mstlystcdata

import (
	"context"
	"database/sql"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/query/fivenet/table"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	j   = table.Jobs.AS("job")
	jg  = table.JobGrades.AS("jobgrade")
	adc = table.FivenetDocumentsCategories.AS("documentcategory")
)

type Cache struct {
	logger *zap.Logger
	db     *sql.DB

	tracer             trace.Tracer
	ctx                context.Context
	jobs               *cache.Cache[string, *jobs.Job]
	docCategories      *cache.Cache[uint64, *documents.DocumentCategory]
	docCategoriesByJob *cache.Cache[string, []*documents.DocumentCategory]

	searcher *Searcher
}

func NewCache(ctx context.Context, logger *zap.Logger, tp *tracesdk.TracerProvider, db *sql.DB) (*Cache, error) {
	jobsCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, *jobs.Job](lru.WithCapacity(32)),
		cache.WithJanitorInterval[string, *jobs.Job](5*time.Minute),
	)

	docCategoriesCache := cache.NewContext(
		ctx,
		cache.AsLRU[uint64, *documents.DocumentCategory](lru.WithCapacity(512)),
	)

	docCategoriesByJobCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*documents.DocumentCategory](lru.WithCapacity(32)),
	)

	c := &Cache{
		logger: logger,
		db:     db,

		tracer:             tp.Tracer("mstlystcdata-cache"),
		ctx:                ctx,
		jobs:               jobsCache,
		docCategories:      docCategoriesCache,
		docCategoriesByJob: docCategoriesByJobCache,
	}

	var err error
	c.searcher, err = NewSearcher(c)
	c.searcher.addDataToIndex()

	return c, err
}

func (c *Cache) Start() {
	for {
		if err := c.refreshCache(); err != nil {
			c.logger.Error("failed to refresh mostyl static data cache", zap.Error(err))
		}

		select {
		case <-c.ctx.Done():
			return
		case <-time.After(3 * time.Minute):
		}
	}
}

func (c *Cache) GetSearcher() *Searcher {
	return c.searcher
}

func (c *Cache) refreshCache() error {
	ctx, span := c.tracer.Start(c.ctx, "mstlystcdata-refresh-cache")
	defer span.End()

	if err := c.refreshDocumentCategories(ctx); err != nil {
		return err
	}

	if err := c.refreshJobs(ctx); err != nil {
		return err
	}

	if c.searcher != nil {
		c.searcher.addDataToIndex()
	}

	return nil
}

func (c *Cache) refreshDocumentCategories(ctx context.Context) error {
	stmt := adc.
		SELECT(
			adc.ID,
			adc.Name,
			adc.Description,
			adc.Job,
		).
		FROM(adc).
		ORDER_BY(
			adc.Job.ASC(),
			adc.Name.ASC(),
		)

	var dest []*documents.DocumentCategory
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		return err
	}

	categoriesPerJob := map[string][]*documents.DocumentCategory{}
	for _, d := range dest {
		c.docCategories.Set(d.Id, d)

		if _, ok := categoriesPerJob[*d.Job]; !ok {
			categoriesPerJob[*d.Job] = []*documents.DocumentCategory{}
		}
		categoriesPerJob[*d.Job] = append(categoriesPerJob[*d.Job], d)
	}

	// Update cache
	for job, cs := range categoriesPerJob {
		c.docCategoriesByJob.Set(job, cs)
	}

	return nil
}

func (c *Cache) refreshJobs(ctx context.Context) error {
	stmt := j.
		SELECT(
			j.Name,
			j.Label,
			jg.JobName.AS("job_grade.job_name"),
			jg.Grade,
			jg.Label,
		).
		FROM(j.
			INNER_JOIN(jg,
				jg.JobName.EQ(j.Name),
			),
		).
		ORDER_BY(
			j.Name.ASC(),
			jg.Grade.ASC(),
		)

	var dest []*jobs.Job
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		return err
	}

	// Update cache
	for _, job := range dest {
		c.jobs.Set(strings.ToLower(job.Name), job)
	}

	return nil
}
