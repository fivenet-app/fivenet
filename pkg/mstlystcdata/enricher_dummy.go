package mstlystcdata

import (
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
)

type DummyEnricher struct {
	IEnricher
}

func NewDummyEnricher() IEnricher {
	return &DummyEnricher{}
}

func (e *DummyEnricher) EnrichJobInfo(user common.IJobInfo) {}

func (e *DummyEnricher) EnrichJobInfoNoFallback(usr common.IJobInfo) {}

func (e *DummyEnricher) EnrichJobName(usr common.IJobName) {}

func dummyJobGrades(job string) []*jobs.JobGrade {
	grades := make([]*jobs.JobGrade, 0, 3)
	for grade := int32(1); grade <= 3; grade++ {
		jobName := job
		grades = append(grades, &jobs.JobGrade{
			JobName: &jobName,
			Grade:   grade,
			Label:   fmt.Sprintf("Rank %d", grade),
		})
	}

	return grades
}

func (e *DummyEnricher) GetJobByName(job string) *jobs.Job {
	return &jobs.Job{
		Name:   job,
		Label:  job,
		Grades: dummyJobGrades(job),
	}
}

func (e *DummyEnricher) GetJobGrade(job string, grade int32) (*jobs.Job, *jobs.JobGrade) {
	return &jobs.Job{
			Name:   job,
			Label:  job,
			Grades: dummyJobGrades(job),
		}, &jobs.JobGrade{
			JobName: func() *string {
				jobName := job
				return &jobName
			}(),
			Grade: grade,
			Label: fmt.Sprintf("Rank %d", grade),
		}
}

type DummyUserAwareEnricher struct {
	IUserAwareEnricher
}

func NewDummyUserAwareEnricher() IUserAwareEnricher {
	return &DummyUserAwareEnricher{}
}

func (e *DummyUserAwareEnricher) EnrichJobInfoSafe(
	userInfo *userinfo.UserInfo,
	usrs ...common.IJobInfo,
) {
}

func (e *DummyUserAwareEnricher) EnrichJobInfoSafeFunc(
	userInfo *userinfo.UserInfo,
) func(usr common.IJobInfo) {
	return func(usr common.IJobInfo) {}
}
