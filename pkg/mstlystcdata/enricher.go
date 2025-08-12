package mstlystcdata

import (
	"slices"
	"strconv"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

// NotAvailablePlaceholder is used as a fallback label when job info is not available.
const (
	NotAvailablePlaceholder = "N/A"
)

// Enricher provides methods to enrich job information for users based on job data and config.
type Enricher struct {
	// jobs is the job data source
	jobs *Jobs

	// appCfg is the application configuration provider
	appCfg appconfig.IConfig
	// jobStartIndex is the starting index for job grades
	jobStartIndex int32
}

// NewEnricher creates a new Enricher instance with the given job data and config.
func NewEnricher(jobs *Jobs, appCfg appconfig.IConfig, cfg *config.Config) *Enricher {
	return &Enricher{
		jobs: jobs,

		appCfg:        appCfg,
		jobStartIndex: cfg.Game.StartJobGrade,
	}
}

// EnrichJobInfo enriches the job information of an object that implements the common.IJobInfo interface.
// Sets job label and grade label, falling back to N/A and unemployed job if not found.
func (e *Enricher) EnrichJobInfo(usr common.IJobInfo) {
	job, err := e.jobs.Get(usr.GetJob())
	if err == nil {
		usr.SetJobLabel(job.GetLabel())

		gradeIndex := max(usr.GetJobGrade()-e.jobStartIndex, 0)

		if len(job.GetGrades()) > int(gradeIndex) {
			usr.SetJobGradeLabel(job.GetGrades()[gradeIndex].GetLabel())
		} else {
			jg := strconv.FormatInt(int64(usr.GetJobGrade()), 10)
			usr.SetJobGradeLabel(jg)
		}
	} else {
		appCfg := e.appCfg.Get()

		usr.SetJobLabel(NotAvailablePlaceholder)
		usr.SetJob(appCfg.JobInfo.GetUnemployedJob().GetName())
		usr.SetJobGradeLabel(NotAvailablePlaceholder)
		usr.SetJobGrade(appCfg.JobInfo.GetUnemployedJob().GetGrade())
	}
}

// EnrichJobInfoNoFallback enriches the job information of an object that implements the common.IJobInfo interface.
// Sets job label and grade label, in case the job is not found, will only set the labels to N/A and unemployed job.
func (e *Enricher) EnrichJobInfoNoFallback(usr common.IJobInfo) {
	job, err := e.jobs.Get(usr.GetJob())
	if err == nil {
		usr.SetJobLabel(job.GetLabel())

		gradeIndex := max(usr.GetJobGrade()-e.jobStartIndex, 0)

		if len(job.GetGrades()) > int(gradeIndex) {
			usr.SetJobGradeLabel(job.GetGrades()[gradeIndex].GetLabel())
		} else {
			jg := strconv.FormatInt(int64(usr.GetJobGrade()), 10)
			usr.SetJobGradeLabel(jg)
		}
	} else {
		usr.SetJobLabel(NotAvailablePlaceholder)
		usr.SetJobGradeLabel(NotAvailablePlaceholder)
	}
}

// EnrichJobName enriches the job label for an object that implements the common.IJobName interface.
func (e *Enricher) EnrichJobName(usr common.IJobName) {
	job, err := e.jobs.Get(usr.GetJob())
	if err == nil {
		usr.SetJobLabel(job.GetLabel())
	} else {
		usr.SetJobLabel(usr.GetJob())
	}
}

// GetJobByName returns the Job struct for a given job name, or nil if not found.
func (e *Enricher) GetJobByName(job string) *jobs.Job {
	j, err := e.jobs.Get(job)
	if err != nil {
		return nil
	}

	return j
}

// GetJobGrade returns the Job and JobGrade for a given job name and grade, or nil if not found.
func (e *Enricher) GetJobGrade(job string, grade int32) (*jobs.Job, *jobs.JobGrade) {
	j := e.GetJobByName(job)
	if j == nil {
		return nil, nil
	}

	for i := range j.GetGrades() {
		if j.GetGrades()[i].GetGrade() == grade {
			return j, j.GetGrades()[i]
		}
	}

	return nil, nil
}

// UserAwareEnricher extends Enricher with permission-aware enrichment for user job info.
type UserAwareEnricher struct {
	// Enricher is the embedded base enricher
	*Enricher

	// ps is the permissions provider
	ps perms.Permissions
}

// NewUserAwareEnricher creates a new UserAwareEnricher with the given enricher and permissions.
func NewUserAwareEnricher(enricher *Enricher, ps perms.Permissions) *UserAwareEnricher {
	return &UserAwareEnricher{
		Enricher: enricher,
		ps:       ps,
	}
}

// EnrichJobInfoSafe enriches job info for multiple users, applying permission checks for the given userInfo.
func (e *UserAwareEnricher) EnrichJobInfoSafe(
	userInfo *userinfo.UserInfo,
	usrs ...common.IJobInfo,
) {
	enrichFn := e.EnrichJobInfoSafeFunc(userInfo)

	for _, usr := range usrs {
		enrichFn(usr)
	}
}

// EnrichJobInfoSafeFunc returns a function that enriches job info for a user, applying permission checks.
func (e *UserAwareEnricher) EnrichJobInfoSafeFunc(
	userInfo *userinfo.UserInfo,
) func(usr common.IJobInfo) {
	jobGrades, _ := e.ps.AttrJobGradeList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceGetUserPerm,
		permscitizens.CitizensServiceGetUserJobsPermField,
	)

	appCfg := e.appCfg.Get()
	publicJobs := appCfg.JobInfo.GetPublicJobs()
	unemployedJob := appCfg.JobInfo.GetUnemployedJob()

	return func(usr common.IJobInfo) {
		// Make sure user has permission to see that grade, otherwise "hide" the user's job grade
		ok := jobGrades.HasJobGrade(usr.GetJob(), usr.GetJobGrade())
		if !ok && !userInfo.GetSuperuser() {
			if !slices.Contains(publicJobs, usr.GetJob()) {
				usr.SetJob(unemployedJob.GetName())
				usr.SetJobGrade(unemployedJob.GetGrade())
			} else {
				usr.SetJobGrade(0)
			}
		} else {
			if !ok && !userInfo.GetSuperuser() {
				usr.SetJobGrade(0)
			}
		}

		e.EnrichJobInfo(usr)
	}
}
