package jobs

import (
	"database/sql/driver"

	"github.com/galexrt/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

func (x *Colleague) SetJob(job string) {
	x.Job = job
}

func (x *Colleague) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *Colleague) SetJobGrade(grade int32) {
	x.JobGrade = grade
}

func (x *Colleague) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

// Scan implements driver.Valuer for protobuf JobsUserActivityData.
func (x *JobsUserActivityData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *JobsUserActivityData) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
