package documentstemplates

import (
	"fmt"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
)

type TemplateJobAccess = resourcesaccess.JobAccess

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
