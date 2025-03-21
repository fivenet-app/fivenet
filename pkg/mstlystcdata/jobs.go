package mstlystcdata

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Jobs struct {
	logger *zap.Logger
	db     *sql.DB

	tracer trace.Tracer

	store *store.Store[users.Job, *users.Job]
	store.StoreRO[users.Job, *users.Job]

	updateCallbacks []updateCallbackFn
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
	JS     *events.JSWrapper
	Config *config.Config

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

func NewJobs(p Params) (*Jobs, error) {
	c := &Jobs{
		logger: p.Logger,
		db:     p.DB,

		tracer: p.TP.Tracer("mstlystcdata-cache"),

		updateCallbacks: []updateCallbackFn{},
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		jobs, err := store.New(ctxStartup, p.Logger, p.JS, "cache",
			store.WithLocks[users.Job](nil),
			store.WithKVPrefix[users.Job]("jobs"),
		)
		if err != nil {
			return err
		}
		c.store = jobs
		c.StoreRO = jobs

		if err := jobs.Start(ctxCancel, true); err != nil {
			return err
		}

		if err := c.loadJobs(ctxStartup); err != nil {
			return err
		}

		p.CronHandlers.Add("mstlystcdata.jobs", func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := c.tracer.Start(ctx, "mstlystcdata-jobs")
			defer span.End()

			if err := c.loadJobs(ctx); err != nil {
				c.logger.Error("failed to refresh jobs cache", zap.Error(err))
				return err
			}

			for _, fn := range c.updateCallbacks {
				if err := fn(ctx); err != nil {
					return err
				}
			}

			return nil
		})

		if err := p.Cron.RegisterCronjob(ctxStartup, &cron.Cronjob{
			Name:     "mstlystcdata.jobs",
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

	return c, nil
}

func (c *Jobs) loadJobs(ctx context.Context) error {
	tJobs := tables.Jobs().AS("job")
	tJobGrades := tables.JobGrades().AS("jobgrade")

	stmt := tJobs.
		SELECT(
			tJobs.Name,
			tJobs.Label,
			tJobGrades.JobName.AS("job_grade.job_name"),
			tJobGrades.Grade,
			tJobGrades.Label,
		).
		FROM(tJobs.
			INNER_JOIN(tJobGrades,
				tJobGrades.JobName.EQ(tJobs.Name),
			),
		).
		ORDER_BY(
			tJobs.Name.ASC(),
			tJobGrades.Grade.ASC(),
		)

	var dest []*users.Job
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// No jobs found in database, remove all from cache
	if len(dest) == 0 {
		if err := c.store.Clear(ctx); err != nil {
			return err
		}

		return nil
	}

	// Update cached jobs
	errs := multierr.Combine()
	// Check which jobs exist and which don't for deletion later
	found := []string{}
	for _, job := range dest {
		jobName := strings.ToLower(job.Name)
		if err := c.store.Put(ctx, jobName, job); err != nil {
			errs = multierr.Append(errs, err)
		}
		found = append(found, jobName)
	}

	// Delete non-existing jobs, based on which are in the database
	c.Range(ctx, func(key string, value *users.Job) bool {
		if !slices.ContainsFunc(found, func(in string) bool {
			return in == key
		}) {
			if err := c.store.Delete(ctx, key); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
		return true
	})

	return errs
}

type updateCallbackFn func(ctx context.Context) error

// Only call during init/fx startup hooks!
func (c *Jobs) addUpdateCallback(fn updateCallbackFn) {
	c.updateCallbacks = append(c.updateCallbacks, fn)
}

func (c *Jobs) GetHighestJobGrade(job string) *users.JobGrade {
	j, ok := c.Get(job)
	if !ok {
		return nil
	}

	if len(j.Grades) == 0 {
		return nil
	}

	return j.Grades[len(j.Grades)-1]
}

type JobsSearchParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger

	Jobs *Jobs
}

type JobsSearch struct {
	logger *zap.Logger

	mu    sync.Mutex
	index bleve.Index

	*Jobs
}

func NewJobsSearch(p JobsSearchParams) (*JobsSearch, error) {
	c := &JobsSearch{
		logger: p.Logger,

		mu:   sync.Mutex{},
		Jobs: p.Jobs,
	}

	index, err := c.newSearchIndex()
	if err != nil {
		return nil, err
	}
	c.index = index

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		var err error
		c.index, err = c.newSearchIndex()
		if err != nil {
			return err
		}

		c.Jobs.addUpdateCallback(c.loadDataIntoIndex)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		return nil
	}))

	return c, nil
}

func (c *JobsSearch) newSearchIndex() (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()

	jobMapping := bleve.NewDocumentMapping()
	gradesMapping := bleve.NewDocumentMapping()
	jobMapping.AddSubDocumentMapping("grades", gradesMapping)
	indexMapping.AddDocumentMapping("job", jobMapping)

	return bleve.NewMemOnly(indexMapping)
}

func (c *JobsSearch) loadDataIntoIndex(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	errs := multierr.Combine()

	batch := c.index.NewBatch()
	// Fill jobs search from cache
	c.Jobs.Range(ctx, func(key string, value *users.Job) bool {
		batch.Delete(key)

		if err := batch.Index(key, value); err != nil {
			errs = multierr.Append(errs, err)
		}

		return true
	})

	if err := c.index.Batch(batch); err != nil {
		errs = multierr.Append(errs, err)
	}

	return errs
}

func (c *JobsSearch) Search(ctx context.Context, search string, exactMatch bool) ([]*users.Job, error) {
	var searchQuery query.Query
	if search == "" {
		searchQuery = bleve.NewMatchAllQuery()
	} else {
		if exactMatch {
			searchQuery = bleve.NewMatchQuery(strings.ToLower(search))
		} else {
			searchQuery = bleve.NewQueryStringQuery(strings.ToLower(search) + "*")
		}
	}

	request := bleve.NewSearchRequest(searchQuery)
	if exactMatch {
		request.Size = 1
	} else {
		request.Size = 32
	}
	request.Fields = []string{"name", "label"}
	request.SortBy([]string{"_score", "label", "_id"})

	result, err := c.index.SearchInContext(ctx, request)
	if err != nil {
		return nil, err
	}

	jobs := []*users.Job{}
	for _, result := range result.Hits {
		job, ok := c.Jobs.Get(result.ID)
		if !ok {
			c.logger.Error("no job found for search result id", zap.String("job", result.ID))
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
