package complhelper

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/arpanet/proto/resources/common"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/proto/resources/jobs"
	"github.com/galexrt/arpanet/query/arpanet/table"
)

var (
	j   = table.Jobs.AS("job")
	jg  = table.JobGrades.AS("job_grade")
	adc = table.ArpanetDocumentsCategories.AS("documentcategory")
)

type Completor struct {
	db            *sql.DB
	cancel        context.CancelFunc
	Jobs          *cache.Cache[string, *jobs.Job]
	DocCategories *cache.Cache[string, []*documents.DocumentCategory]
}

func New(db *sql.DB) *Completor {
	ctx, cancel := context.WithCancel(context.Background())

	jobsCache := cache.NewContext(
		ctx,
		cache.AsLFU[string, *jobs.Job](lfu.WithCapacity(32)),
		cache.WithJanitorInterval[string, *jobs.Job](120*time.Second),
	)

	doccategoriesCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*documents.DocumentCategory](lru.WithCapacity(32)),
	)
	c := &Completor{
		db:            db,
		cancel:        cancel,
		Jobs:          jobsCache,
		DocCategories: doccategoriesCache,
	}

	c.refreshCache()

	return c
}

func (c *Completor) refreshCache() error {
	if err := c.refreshDocumentCategories(); err != nil {
		return err
	}

	if err := c.refreshJobsCache(); err != nil {
		return err
	}

	return nil
}

func (c *Completor) refreshDocumentCategories() error {
	var dest []*documents.DocumentCategory

	stmt := adc.
		SELECT(
			adc.ID,
			adc.Name,
			adc.Description,
			adc.Job,
		).
		FROM(adc).
		GROUP_BY(adc.Job).
		ORDER_BY(adc.Name.ASC())

	if err := stmt.Query(c.db, &dest); err != nil {
		return err
	}

	categoriesPerJob := map[string][]*documents.DocumentCategory{}
	for _, c := range dest {
		if _, ok := categoriesPerJob[c.Job]; !ok {
			categoriesPerJob[c.Job] = []*documents.DocumentCategory{}
		}
		categoriesPerJob[c.Job] = append(categoriesPerJob[c.Job], c)
	}

	// Update cache
	for job, cs := range categoriesPerJob {
		c.DocCategories.Set(job, cs)
	}

	return nil
}

func (c *Completor) refreshJobsCache() error {
	var dest []*jobs.Job

	stmt := j.
		SELECT(
			j.Name,
			j.Label,
			jg.JobName.AS("job_grade.job_name"),
			jg.Grade,
			jg.Label,
		).FROM(
		j.LEFT_JOIN(jg,
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
		c.Jobs.Set(strings.ToLower(job.Name), job)
	}

	return nil
}

func (c *Completor) ResolveJob(usr common.IJobInfo) {
	job, ok := c.Jobs.Get(usr.GetJob())
	if ok {
		usr.SetJobLabel(job.Label)

		jg := usr.GetJobGrade() - 1

		if len(job.Grades) >= int(jg) {
			usr.SetJobGradeLabel(job.Grades[jg].Label)
		} else {
			jg := strconv.Itoa(int(usr.GetJobGrade()))
			usr.SetJobGradeLabel(jg)
		}
	} else {
		usr.SetJobLabel(usr.GetJob())
		jg := strconv.Itoa(int(usr.GetJobGrade()))
		usr.SetJobGradeLabel(jg)
	}
}
