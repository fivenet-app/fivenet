package mstlystcdata

import (
	"strconv"

	"github.com/galexrt/fivenet/proto/resources/common"
	"github.com/galexrt/fivenet/proto/resources/documents"
	"github.com/galexrt/fivenet/proto/resources/jobs"
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

func (e *Enricher) EnrichDocumentCategory(doc common.IDocumentCategory) {
	cId := doc.GetCategoryId()

	// No category
	if cId == 0 {
		return
	}

	dc, ok := e.c.docCategories.Get(cId)
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

func (e *Enricher) GetJobByName(name string) *jobs.Job {
	job, ok := e.c.jobs.Get(name)
	if !ok {
		return nil
	}

	return job
}
