package mstlystcdata

import (
	"strconv"

	"github.com/galexrt/arpanet/proto/resources/common"
	"github.com/galexrt/arpanet/proto/resources/documents"
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
		usr.SetJobLabel(usr.GetJob())
		jg := strconv.Itoa(int(usr.GetJobGrade()))
		usr.SetJobGradeLabel(jg)
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
