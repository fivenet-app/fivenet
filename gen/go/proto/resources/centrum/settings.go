package centrum

import (
	"database/sql/driver"
	"slices"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (x *Settings) Default(job string) {
	x.Job = job

	if x.Mode <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.Mode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.FallbackMode <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.FallbackMode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.PredefinedStatus == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	}

	if x.Timings == nil {
		x.Timings = &Timings{}
	}
	if x.Timings.DispatchMaxWait <= 0 {
		x.Timings.DispatchMaxWait = 900
	}
	if x.Timings.RequireUnitReminderSeconds <= 0 {
		x.Timings.RequireUnitReminderSeconds = 180
	}
}

func (x *Settings) Merge(in *Settings) *Settings {
	x.Job = in.Job
	x.Enabled = in.Enabled

	x.Mode = in.Mode
	x.FallbackMode = in.FallbackMode

	if in.PredefinedStatus == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	} else {
		x.PredefinedStatus = in.PredefinedStatus
	}

	if in.Timings == nil {
		x.Timings = &Timings{}
	} else {
		x.Timings = in.Timings
	}

	if in.Access == nil {
		x.Access = &CentrumAccess{}
	} else {
		x.Access = in.Access
	}

	x.Public = in.Public

	return x
}

func (x *Settings) JobHasAccess(job string, access CentrumAccessLevel) bool {
	if x.Access == nil || x.Access.Jobs == nil {
		return x.Job == job // No access restrictions defined, only the job itself is allowed
	}

	if x.Job == job {
		return true // Job is explicitly allowed
	}

	for _, j := range x.Access.Jobs {
		if j.Job == job && j.Access >= access {
			return true // Job is explicitly allowed in the access list
		}
	}

	return false // Default access is deny
}

// Scan implements driver.Valuer for protobuf JobList. Special case as only the job list is stored in the database, not the full protobuf message.
func (x *JobList) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.Unmarshal([]byte(t), &x.Jobs)
	case []byte:
		return json.Unmarshal(t, &x.Jobs)
	}
	return nil
}

// Value marshals the JobList value into driver.Valuer. Special case as only the job list is stored in the database, not the full protobuf message.
func (x *JobList) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	return json.MarshalToString(x.Jobs)
}

func (x *JobList) ContainsJob(job string) bool {
	return slices.Contains(x.Jobs, job)
}
