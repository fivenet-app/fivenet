package mstlystcdata

import (
	"context"
	"database/sql"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/proto/resources/jobs"
	"github.com/galexrt/arpanet/query/arpanet/table"
	"go.uber.org/zap"
)

var (
	j   = table.Jobs.AS("job")
	jg  = table.JobGrades.AS("job_grade")
	adc = table.ArpanetDocumentsCategories.AS("documentcategory")
)

type Cache struct {
	logger *zap.Logger
	db     *sql.DB

	ctx                context.Context
	cancel             context.CancelFunc
	jobs               *cache.Cache[string, *jobs.Job]
	docCategories      *cache.Cache[uint64, *documents.DocumentCategory]
	docCategoriesByJob *cache.Cache[string, []*documents.DocumentCategory]

	searcher *Searcher
}

func NewCache(logger *zap.Logger, db *sql.DB) (*Cache, error) {
	ctx, cancel := context.WithCancel(context.Background())

	jobsCache := cache.NewContext(
		ctx,
		cache.AsLFU[string, *jobs.Job](lfu.WithCapacity(32)),
		cache.WithJanitorInterval[string, *jobs.Job](120*time.Second),
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

		ctx:                ctx,
		cancel:             cancel,
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
	if err := c.refreshCache(); err != nil {
		c.logger.Error("failed to refresh mostyl static data cache", zap.Error(err))
	}

	go func() {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(5 * time.Minute):
			if err := c.refreshCache(); err != nil {
				c.logger.Error("failed to refresh mostyl static data cache", zap.Error(err))
			}
		}
	}()
}

func (c *Cache) GetSearcher() *Searcher {
	return c.searcher
}

func (c *Cache) refreshCache() error {
	if err := c.refreshDocumentCategories(); err != nil {
		return err
	}

	if err := c.refreshJobsCache(); err != nil {
		return err
	}

	if c.searcher != nil {
		c.searcher.addDataToIndex()
	}

	return nil
}

func (c *Cache) refreshDocumentCategories() error {
	var dest []*documents.DocumentCategory

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

	if err := stmt.Query(c.db, &dest); err != nil {
		return err
	}

	categoriesPerJob := map[string][]*documents.DocumentCategory{}
	for _, d := range dest {
		c.docCategories.Set(d.Id, d)

		if _, ok := categoriesPerJob[d.Job]; !ok {
			categoriesPerJob[d.Job] = []*documents.DocumentCategory{}
		}
		categoriesPerJob[d.Job] = append(categoriesPerJob[d.Job], d)
	}

	// Update cache
	for job, cs := range categoriesPerJob {
		c.docCategoriesByJob.Set(job, cs)
	}

	return nil
}

func (c *Cache) refreshJobsCache() error {
	var dest []*jobs.Job

	stmt := j.
		SELECT(
			j.Name,
			j.Label,
			jg.JobName.AS("job_grade.job_name"),
			jg.Grade,
			jg.Label,
		).
		FROM(j.
			LEFT_JOIN(jg,
				jg.JobName.EQ(j.Name),
			),
		).
		ORDER_BY(
			j.Name.ASC(),
			jg.Grade.ASC(),
		)

	if err := stmt.Query(c.db, &dest); err != nil {
		return err
	}

	// Update cache
	for _, job := range dest {
		c.jobs.Set(strings.ToLower(job.Name), job)
	}

	return nil
}
