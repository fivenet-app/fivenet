package documentsaccess

import (
	"fmt"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
)

type DocumentAccess = resourcesaccess.Access
type DocumentJobAccess = resourcesaccess.JobAccess
type DocumentUserAccess = resourcesaccess.UserAccess

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
