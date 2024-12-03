package centrum

import (
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

func (x *UnitAccess) IsEmpty() bool {
	return len(x.Jobs) == 0 && len(x.Qualifications) == 0
}

func (x *UnitAccess) ClearQualificationResults() {
	for _, quali := range x.Qualifications {
		if quali.Qualification != nil && quali.Qualification.Result != nil {
			quali.Qualification.Result = nil
		}
	}
}

func (x *UnitJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *UnitJobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *UnitJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *UnitJobAccess) SetJobGradeLabel(jobLabel string) {
	x.JobGradeLabel = &jobLabel
}

// pkg/access compatibility

func (x *UnitJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *UnitJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *UnitJobAccess) SetAccess(access UnitAccessLevel) {
	x.Access = access
}

func (x *UnitUserAccess) GetId() uint64 {
	return 0
}

func (x *UnitUserAccess) GetTargetId() uint64 {
	return 0
}

func (x *UnitUserAccess) GetUserId() int32 {
	return 0
}

func (x *UnitUserAccess) GetAccess() UnitAccessLevel {
	return UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN
}

func (x *UnitUserAccess) SetUserId(id int32) {}

func (x *UnitUserAccess) SetAccess(access UnitAccessLevel) {}

func (x *UnitQualificationAccess) SetQualificationId(id uint64) {
	x.QualificationId = id
}

func (x *UnitQualificationAccess) SetAccess(access UnitAccessLevel) {
	x.Access = access
}

// Scan implements driver.Valuer for protobuf UnitAccess.
func (x *UnitAccess) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *UnitAccess) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
