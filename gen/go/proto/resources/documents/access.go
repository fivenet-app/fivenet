package documents

import (
	"fmt"
)

func (x *DocumentJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *DocumentJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

// pkg/access compatibility

func (x *DocumentJobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *DocumentJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *DocumentJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *DocumentJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *DocumentJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *DocumentUserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *DocumentUserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func DocumentAccessHasDuplicates(access *DocumentAccess) bool {
	jobKeys := map[string]any{}
	for _, ja := range access.Jobs {
		key := fmt.Sprintf("%s-%d", ja.GetJob(), ja.GetMinimumGrade())
		if _, ok := jobKeys[key]; ok {
			return true
		}
		jobKeys[key] = nil
	}

	userKeys := map[int32]any{}
	for _, ja := range access.Users {
		if _, ok := userKeys[ja.GetUserId()]; ok {
			return true
		}
		userKeys[ja.GetUserId()] = nil
	}

	return false
}

func TemplateAccessHasDuplicates(jobs []*TemplateJobAccess) bool {
	jobKeys := map[string]any{}
	for _, ja := range jobs {
		key := fmt.Sprintf("%s-%d", ja.GetJob(), ja.GetMinimumGrade())
		if _, ok := jobKeys[key]; ok {
			return true
		}
		jobKeys[key] = nil
	}

	return false
}
