package dataenricher

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

type Enricher struct {
	db                 *sql.DB
	cancel             context.CancelFunc
	Jobs               *cache.Cache[string, *jobs.Job]
	DocCategories      *cache.Cache[uint64, *documents.DocumentCategory]
	DocCategoriesByJob *cache.Cache[string, []*documents.DocumentCategory]
}

func New(db *sql.DB) *Enricher {
	ctx, cancel := context.WithCancel(context.Background())

	jobsCache := cache.NewContext(
		ctx,
		cache.AsLFU[string, *jobs.Job](lfu.WithCapacity(32)),
		cache.WithJanitorInterval[string, *jobs.Job](120*time.Second),
	)

	doccategoriesCache := cache.NewContext(
		ctx,
		cache.AsLRU[uint64, *documents.DocumentCategory](lru.WithCapacity(512)),
	)

	doccategoriesByJobCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*documents.DocumentCategory](lru.WithCapacity(32)),
	)

	c := &Enricher{
		db:                 db,
		cancel:             cancel,
		Jobs:               jobsCache,
		DocCategories:      doccategoriesCache,
		DocCategoriesByJob: doccategoriesByJobCache,
	}

	c.refreshCache()

	return c
}

func (c *Enricher) refreshCache() error {
	if err := c.refreshDocumentCategories(); err != nil {
		return err
	}

	if err := c.refreshJobsCache(); err != nil {
		return err
	}

	return nil
}

func (c *Enricher) refreshDocumentCategories() error {
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
	for _, d := range dest {
		c.DocCategories.Set(d.Id, d)

		if _, ok := categoriesPerJob[d.Job]; !ok {
			categoriesPerJob[d.Job] = []*documents.DocumentCategory{}
		}
		categoriesPerJob[d.Job] = append(categoriesPerJob[d.Job], d)
	}

	// Update cache
	for job, cs := range categoriesPerJob {
		c.DocCategoriesByJob.Set(job, cs)
	}

	return nil
}

func (c *Enricher) refreshJobsCache() error {
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
		c.Jobs.Set(strings.ToLower(job.Name), job)
	}

	return nil
}

func (c *Enricher) EnrichJobInfo(usr common.IJobInfo) {
	job, ok := c.Jobs.Get(usr.GetJob())
	if ok {
		usr.SetJobLabel(job.Label)

		jg := usr.GetJobGrade() - 1
		if jg < 0 {
			jg = 0
		}

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

func (c *Enricher) EnrichDocumentCategory(doc common.IDocumentCategory) {
	cId := doc.GetCategoryId()

	// No category
	if cId == 0 {
		return
	}

	dc, ok := c.DocCategories.Get(cId)
	if !ok {
		doc.SetCategory(&documents.DocumentCategory{
			Id:   0,
			Name: "N/A",
			Job:  "N/A",
		})
	} else {
		doc.SetCategory(dc)
	}
}
