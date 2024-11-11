package documents

import (
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

// Scan implements driver.Valuer for protobuf TemplateSchema.
func (x *TemplateSchema) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *TemplateSchema) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func (x *Template) GetJob() string {
	return x.CreatorJob
}

func (x *Template) SetJobLabel(label string) {
	x.CreatorJobLabel = &label
}

func (x *TemplateShort) GetJob() string {
	return x.CreatorJob
}

func (x *TemplateShort) SetJobLabel(label string) {
	x.CreatorJobLabel = &label
}

// pkg/access compatibility

func (x *TemplateJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *TemplateJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *TemplateUserAccess) GetAccess() AccessLevel {
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *TemplateUserAccess) GetId() uint64 {
	return 0
}

func (x *TemplateUserAccess) GetTargetId() uint64 {
	return 0
}

func (x *TemplateUserAccess) SetAccess(access AccessLevel) {}

func (x *TemplateUserAccess) GetUserId() int32 {
	return 0
}

func (x *TemplateUserAccess) SetUserId(userId int32) {}
