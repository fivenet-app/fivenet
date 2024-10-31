package mstlystcdata

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/laws"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/croner"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/puzpuzpuz/xsync/v3"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var (
	tJobs      = table.Jobs.AS("job")
	tJGrades   = table.JobGrades.AS("jobgrade")
	tDCategory = table.FivenetDocumentsCategories.AS("category")
	tLawBooks  = table.FivenetLawbooks.AS("lawbook")
	tLaws      = table.FivenetLawbooksLaws.AS("law")
)

type Cache struct {
	logger *zap.Logger
	db     *sql.DB

	refreshTime time.Duration

	tracer             trace.Tracer
	jobs               *store.Store[users.Job, *users.Job]
	docCategories      *store.Store[documents.Category, *documents.Category]
	docCategoriesByJob *xsync.MapOf[string, []*documents.Category]
	lawBooks           *xsync.MapOf[uint64, *laws.LawBook]

	searcher *Searcher
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
	JS     *events.JSWrapper
	Config *config.Config

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

func NewCache(p Params) (*Cache, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	cc := &Cache{
		logger: p.Logger,
		db:     p.DB,

		refreshTime: p.Config.Cache.RefreshTime,

		tracer:             p.TP.Tracer("mstlystcdata-cache"),
		docCategoriesByJob: xsync.NewMapOf[string, []*documents.Category](),
		lawBooks:           xsync.NewMapOf[uint64, *laws.LawBook](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		jobs, err := store.New[users.Job, *users.Job](ctxStartup, p.Logger, p.JS, "cache",
			store.WithLocks[users.Job, *users.Job](nil),
			store.WithKVPrefix[users.Job, *users.Job]("jobs"),
		)
		if err != nil {
			return err
		}
		cc.jobs = jobs

		docCategories, err := store.New[documents.Category, *documents.Category](ctxStartup, p.Logger, p.JS, "cache",
			store.WithLocks[documents.Category, *documents.Category](nil),
			store.WithKVPrefix[documents.Category, *documents.Category]("doc_categories"),
		)
		if err != nil {
			return err
		}
		cc.docCategories = docCategories

		if err := jobs.Start(ctxCancel); err != nil {
			return err
		}

		if err := docCategories.Start(ctxCancel); err != nil {
			return err
		}

		cc.searcher, err = NewSearcher(cc)
		if err != nil {
			return err
		}
		// TODO we have to run addDataToIndex every now and then

		if err := cc.refreshCache(ctxStartup); err != nil {
			return err
		}

		p.CronHandlers.Add("mstlystcdata.cache", func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := cc.tracer.Start(ctx, "mstlystcdata-cache")
			defer span.End()

			if err := cc.refreshCache(ctx); err != nil {
				cc.logger.Error("failed to refresh mostly static data cache", zap.Error(err))
				return err
			}

			return nil
		})

		if err := p.Cron.RegisterCronjob(ctxStartup, &cron.Cronjob{
			Name:     "mstlystcdata.cache",
			Schedule: "@10minutes",
		}); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return cc, nil
}

func (c *Cache) GetSearcher() *Searcher {
	return c.searcher
}

func (c *Cache) refreshCache(ctx context.Context) error {
	errs := multierr.Combine()

	if err := c.refreshCategories(ctx); err != nil {
		errs = multierr.Append(errs, err)
	}

	if err := c.refreshJobs(ctx); err != nil {
		errs = multierr.Append(errs, err)
	}

	if err := c.RefreshLaws(ctx, 0); err != nil {
		errs = multierr.Append(errs, err)
	}

	if c.searcher != nil {
		if err := c.searcher.addDataToIndex(ctx); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (c *Cache) refreshCategories(ctx context.Context) error {
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
		if err := c.docCategories.Put(ctx, key, d); err != nil {
			errs = multierr.Append(errs, err)
		}

		if _, ok := categoriesPerJob[*d.Job]; !ok {
			categoriesPerJob[*d.Job] = []*documents.Category{}
		}
		categoriesPerJob[*d.Job] = append(categoriesPerJob[*d.Job], d)
	}

	// Update per jobs cache
	for job, cs := range categoriesPerJob {
		c.docCategoriesByJob.Store(job, cs)
	}

	return errs
}

func (c *Cache) refreshJobs(ctx context.Context) error {
	stmt := tJobs.
		SELECT(
			tJobs.Name,
			tJobs.Label,
			tJGrades.JobName.AS("job_grade.job_name"),
			tJGrades.Grade,
			tJGrades.Label,
		).
		FROM(tJobs.
			INNER_JOIN(tJGrades,
				tJGrades.JobName.EQ(tJobs.Name),
			),
		).
		ORDER_BY(
			tJobs.Name.ASC(),
			tJGrades.Grade.ASC(),
		)

	var dest []*users.Job
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		return err
	}

	// Update cache
	errs := multierr.Combine()
	if len(dest) == 0 {
		if err := c.jobs.Clear(ctx); err != nil {
			return err
		}
	} else {
		// Update cache
		found := []string{}
		for _, job := range dest {
			jobName := strings.ToLower(job.Name)
			if err := c.jobs.Put(ctx, jobName, job); err != nil {
				errs = multierr.Append(errs, err)
			}
			found = append(found, jobName)
		}

		// Delete non-existing jobs, based on which are in the database
		c.jobs.Range(ctx, func(key string, value *users.Job) bool {
			if !slices.ContainsFunc(found, func(in string) bool {
				return in == key
			}) {
				if err := c.jobs.Delete(ctx, key); err != nil {
					errs = multierr.Append(errs, err)
				}
			}
			return true
		})
	}

	return errs
}

func (c *Cache) RefreshLaws(ctx context.Context, lawBookId uint64) error {
	stmt := tLawBooks.
		SELECT(
			tLawBooks.ID,
			tLawBooks.CreatedAt,
			tLawBooks.UpdatedAt,
			tLawBooks.Name,
			tLawBooks.Description,
			tLaws.ID,
			tLaws.LawbookID,
			tLaws.CreatedAt,
			tLaws.UpdatedAt,
			tLaws.Name,
			tLaws.Description,
			tLaws.Hint,
			tLaws.Fine,
			tLaws.DetentionTime,
			tLaws.StvoPoints,
		).
		FROM(tLawBooks.
			LEFT_JOIN(tLaws,
				tLaws.LawbookID.EQ(tLawBooks.ID),
			),
		).
		ORDER_BY(
			tLawBooks.Name.ASC(),
			tLaws.Name.ASC(),
		)

	if lawBookId > 0 {
		stmt = stmt.WHERE(
			tLawBooks.ID.EQ(jet.Uint64(lawBookId)),
		)
	}

	var dest []*laws.LawBook
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// Single lawbook update or lawbook deleted => not found, need to remove it
	if len(dest) == 0 {
		c.lawBooks.Delete(lawBookId)
	} else if lawBookId > 0 {
		c.lawBooks.Store(lawBookId, dest[0])
	} else {
		// Update cache
		found := []uint64{}
		for _, lawbook := range dest {
			c.lawBooks.Store(lawbook.Id, lawbook)
			found = append(found, lawbook.Id)
		}

		// Delete non-existing law books, based on which are in the database
		c.lawBooks.Range(func(key uint64, value *laws.LawBook) bool {
			if !slices.ContainsFunc(found, func(in uint64) bool {
				return in == key
			}) {
				c.lawBooks.Delete(key)
			}
			return true
		})
	}

	return nil
}

func (c *Cache) GetLawBooks() []*laws.LawBook {
	lawBooks := []*laws.LawBook{}
	c.lawBooks.Range(func(key uint64, value *laws.LawBook) bool {
		lawBooks = append(lawBooks, value)
		return true
	})

	sort.Slice(lawBooks, func(i, j int) bool {
		return strings.EqualFold(lawBooks[i].Name, lawBooks[j].Name)
	})

	return lawBooks
}

func (c *Cache) GetHighestJobGrade(job string) *users.JobGrade {
	j, ok := c.jobs.Get(job)
	if !ok {
		return nil
	}

	if len(j.Grades) == 0 {
		return nil
	}

	return j.Grades[len(j.Grades)-1]
}
