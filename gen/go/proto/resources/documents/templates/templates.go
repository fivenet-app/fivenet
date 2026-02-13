package documentstemplates

import (
	"fmt"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
)

func (x *Template) GetJob() string {
	return x.GetCreatorJob()
}

func (x *Template) SetJobLabel(label string) {
	x.CreatorJobLabel = &label
}

func (x *TemplateShort) GetJob() string {
	return x.GetCreatorJob()
}

func (x *TemplateShort) SetJobLabel(label string) {
	x.CreatorJobLabel = &label
}

// pkg/access compatibility

func (x *TemplateUserAccess) GetAccess() documentsaccess.AccessLevel {
	return documentsaccess.AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *TemplateUserAccess) GetId() int64 {
	return 0
}

func (x *TemplateUserAccess) GetTargetId() int64 {
	return 0
}

func (x *TemplateUserAccess) SetAccess(access documentsaccess.AccessLevel) {}

func (x *TemplateUserAccess) GetUserId() int32 {
	return 0
}

func (x *TemplateUserAccess) SetUserId(userId int32) {}

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
