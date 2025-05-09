package mstlystcdata

import (
	"slices"
	"strconv"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	permscitizenstore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizenstore/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

const (
	NotAvailablePlaceholder = "N/A"
)

type Enricher struct {
	jobs *Jobs

	appCfg        appconfig.IConfig
	jobStartIndex int32
}

func NewEnricher(jobs *Jobs, appCfg appconfig.IConfig, cfg *config.Config) *Enricher {
	return &Enricher{
		jobs: jobs,

		appCfg:        appCfg,
		jobStartIndex: cfg.Game.StartJobGrade,
	}
}

// Jobs

func (e *Enricher) EnrichJobInfo(usr common.IJobInfo) {
	job, ok := e.jobs.Get(usr.GetJob())
	if ok {
		usr.SetJobLabel(job.Label)

		gradeIndex := max(usr.GetJobGrade()-e.jobStartIndex, 0)

		if len(job.Grades) > int(gradeIndex) {
			usr.SetJobGradeLabel(job.Grades[gradeIndex].Label)
		} else {
			jg := strconv.Itoa(int(usr.GetJobGrade()))
			usr.SetJobGradeLabel(jg)
		}
	} else {
		appCfg := e.appCfg.Get()

		usr.SetJobLabel("N/A")
		usr.SetJob(appCfg.JobInfo.UnemployedJob.Name)
		usr.SetJobGradeLabel("N/A")
		usr.SetJobGrade(appCfg.JobInfo.UnemployedJob.Grade)
	}
}

func (e *Enricher) EnrichJobName(usr common.IJobName) {
	job, ok := e.jobs.Get(usr.GetJob())
	if ok {
		usr.SetJobLabel(job.Label)
	} else {
		usr.SetJobLabel(usr.GetJob())
	}
}

func (e *Enricher) GetJobByName(job string) *users.Job {
	j, ok := e.jobs.Get(job)
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

	for i := range j.Grades {
		if j.Grades[i].Grade == grade {
			return j, j.Grades[i]
		}
	}

	return nil, nil
}

type UserAwareEnricher struct {
	*Enricher

	ps perms.Permissions
}

func NewUserAwareEnricher(enricher *Enricher, ps perms.Permissions) *UserAwareEnricher {
	return &UserAwareEnricher{
		Enricher: enricher,
		ps:       ps,
	}
}

func (e *UserAwareEnricher) EnrichJobInfoSafe(userInfo *userinfo.UserInfo, usrs ...common.IJobInfo) {
	enrichFn := e.EnrichJobInfoSafeFunc(userInfo)

	for _, usr := range usrs {
		enrichFn(usr)
	}
}

func (e *UserAwareEnricher) EnrichJobInfoSafeFunc(userInfo *userinfo.UserInfo) func(usr common.IJobInfo) {
	jobGrades, _ := e.ps.AttrJobGradeList(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceGetUserPerm, permscitizenstore.CitizenStoreServiceGetUserJobsPermField)

	appCfg := e.appCfg.Get()
	publicJobs := appCfg.JobInfo.PublicJobs
	unemployedJob := appCfg.JobInfo.UnemployedJob

	return func(usr common.IJobInfo) {
		// Make sure user has permission to see that grade, otherwise "hide" the user's job grade
		ok := jobGrades.HasJobGrade(usr.GetJob(), usr.GetJobGrade())
		if !ok && !userInfo.SuperUser {
			if !slices.Contains(publicJobs, usr.GetJob()) {
				usr.SetJob(unemployedJob.Name)
				usr.SetJobGrade(unemployedJob.Grade)
			} else {
				usr.SetJobGrade(0)
			}
		} else {
			if !ok && !userInfo.SuperUser {
				usr.SetJobGrade(0)
			}
		}

		e.EnrichJobInfo(usr)
	}
}
