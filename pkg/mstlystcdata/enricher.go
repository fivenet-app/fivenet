package mstlystcdata

import (
	"strconv"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
)

const (
	NotAvailablePlaceholder = "N/A"
)

type Enricher struct {
	c *Cache
}

func NewEnricher(c *Cache) *Enricher {
	return &Enricher{
		c: c,
	}
}

func (e *Enricher) EnrichJobInfo(usr common.IJobInfo) {
	job, ok := e.c.jobs.Get(usr.GetJob())
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
		usr.SetJobLabel("N/A")
		usr.SetJobGradeLabel("N/A")
	}
}

func (e *Enricher) EnrichJobName(usr common.IJobName) {
	job, ok := e.c.jobs.Get(usr.GetJob())
	if ok {
		usr.SetJobLabel(job.Label)
	} else {
		usr.SetJobLabel(usr.GetJob())
	}
}

func (e *Enricher) EnrichCategory(doc common.ICategory) {
	cId := doc.GetCategoryId()

	// No category
	if cId == 0 {
		return
	}

	dc, ok := e.c.docCategories.Get(cId)
	if !ok {
		job := NotAvailablePlaceholder
		doc.SetCategory(&documents.Category{
			Id:   0,
			Name: "N/A",
			Job:  &job,
		})
	} else {
		doc.SetCategory(dc)
	}
}

func (e *Enricher) GetJobByName(job string) *jobs.Job {
	j, ok := e.c.jobs.Get(job)
	if !ok {
		return nil
	}

	return j
}

func (e *Enricher) GetJobGrade(job string, grade int32) (*jobs.Job, *jobs.JobGrade) {
	j, ok := e.c.jobs.Get(job)
	if !ok {
		return nil, nil
	}

	for i := 0; i < len(j.Grades); i++ {
		if j.Grades[i].Grade == grade {
			return j, j.Grades[i]
		}
	}

	return nil, nil
}
