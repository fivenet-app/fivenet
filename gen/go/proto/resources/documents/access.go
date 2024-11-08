package documents

import (
	"database/sql/driver"
	"fmt"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

func (x *DocumentJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *DocumentJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

// pkg/access compatibility

func (x *DocumentJobAccess) GetTargetId() uint64 {
	return x.DocumentId
}

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

func (x *DocumentUserAccess) GetTargetId() uint64 {
	return x.DocumentId
}

func (x *DocumentUserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *DocumentUserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

// Scan implements driver.Valuer for protobuf DocumentAccess.
func (x *DocumentAccess) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *DocumentAccess) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func DocumentAccessHasDuplicates(access *DocumentAccess) bool {
	jobKeys := map[string]interface{}{}
	for _, ja := range access.Jobs {
		key := fmt.Sprintf("%s-%d", ja.GetJob(), ja.GetMinimumGrade())
		if _, ok := jobKeys[key]; ok {
			return true
		}
		jobKeys[key] = nil
	}

	userKeys := map[int32]interface{}{}
	for _, ja := range access.Users {
		if _, ok := userKeys[ja.GetUserId()]; ok {
			return true
		}
		userKeys[ja.GetUserId()] = nil
	}

	return false
}

func TemplateAccessHasDuplicates(jobs []*TemplateJobAccess) bool {
	jobKeys := map[string]interface{}{}
	for _, ja := range jobs {
		key := fmt.Sprintf("%s-%d", ja.GetJob(), ja.GetMinimumGrade())
		if _, ok := jobKeys[key]; ok {
			return true
		}
		jobKeys[key] = nil
	}

	return false
}
