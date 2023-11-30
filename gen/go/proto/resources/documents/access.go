package documents

import "database/sql/driver"

func (x *DocumentJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *DocumentJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
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

func (x *DocumentAccess) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf DocumentAccess.
func (x *DocumentAccess) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}
