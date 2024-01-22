package mstlystcdata

import (
	"slices"
	"strconv"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	permscitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore/perms"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms"
)

const (
	NotAvailablePlaceholder = "N/A"
)

type Enricher struct {
	cache *Cache
	ps    perms.Permissions

	publicJobs         []string
	unemployedJob      string
	unemployedJobGrade int32
}

func NewEnricher(cache *Cache, ps perms.Permissions, cfg *config.Config) *Enricher {
	return &Enricher{
		cache: cache,
		ps:    ps,

		publicJobs:         cfg.Game.PublicJobs,
		unemployedJob:      cfg.Game.UnemployedJob.Name,
		unemployedJobGrade: cfg.Game.UnemployedJob.Grade,
	}
}

func (e *Enricher) EnrichJobInfo(usr common.IJobInfo) {
	job, ok := e.cache.jobs.Get(usr.GetJob())
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
	job, ok := e.cache.jobs.Get(usr.GetJob())
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

	dc, ok := e.cache.docCategories.Get(cId)
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

func (e *Enricher) GetJobByName(job string) *users.Job {
	j, ok := e.cache.jobs.Get(job)
	if !ok {
		return nil
	}

	return j
}

func (e *Enricher) GetJobGrade(job string, grade int32) (*users.Job, *users.JobGrade) {
	j := e.GetJobByName(job)
	if j == nil {
		return nil, nil
	}

	for i := 0; i < len(j.Grades); i++ {
		if j.Grades[i].Grade == grade {
			return j, j.Grades[i]
		}
	}

	return nil, nil
}

func (e *Enricher) EnrichJobInfoSafe(userInfo *userinfo.UserInfo, usrs ...common.IJobInfo) {
	enrichFn := e.EnrichJobInfoSafeFunc(userInfo)

	for _, usr := range usrs {
		enrichFn(usr)
	}
}

func (e *Enricher) EnrichJobInfoSafeFunc(userInfo *userinfo.UserInfo) func(usr common.IJobInfo) {
	jobGradesAttr, _ := e.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceGetUserPerm, permscitizenstore.CitizenStoreServiceGetUserJobsPermField)
	var jobGrades perms.JobGradeList
	if jobGradesAttr != nil {
		jobGrades = jobGradesAttr.(map[string]int32)
	}

	return func(usr common.IJobInfo) {
		// Make sure user has permission to see that grade, otherwise "hide" the user's job grade
		grade, ok := jobGrades[usr.GetJob()]
		if !ok && !userInfo.SuperUser {
			if !slices.Contains(e.publicJobs, usr.GetJob()) {
				usr.SetJob(e.unemployedJob)
				usr.SetJobGrade(e.unemployedJobGrade)
			} else {
				usr.SetJobGrade(0)
			}
		} else {
			if usr.GetJobGrade() > grade && !userInfo.SuperUser {
				usr.SetJobGrade(0)
			}
		}

		e.EnrichJobInfo(usr)
	}
}
