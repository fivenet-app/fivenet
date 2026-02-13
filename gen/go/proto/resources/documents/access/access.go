package documentsaccess

import (
	"fmt"
)

// pkg/access compatibility

func (x *DocumentJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *DocumentJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func DocumentAccessHasDuplicates(access *DocumentAccess) bool {
	jobKeys := map[string]any{}
	for _, ja := range access.GetJobs() {
		key := fmt.Sprintf("%s-%d", ja.GetJob(), ja.GetMinimumGrade())
		if _, ok := jobKeys[key]; ok {
			return true
		}
		jobKeys[key] = nil
	}

	userKeys := map[int32]any{}
	for _, ja := range access.GetUsers() {
		if _, ok := userKeys[ja.GetUserId()]; ok {
			return true
		}
		userKeys[ja.GetUserId()] = nil
	}

	return false
}
