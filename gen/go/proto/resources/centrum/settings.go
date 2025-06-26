package centrum

import (
	"database/sql/driver"
	"slices"
	"time"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/durationpb"
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

	if x.Configuration == nil {
		x.Configuration = &Configuration{
			DeduplicationEnabled:  true,
			DeduplicationRadius:   45,
			DeduplicationDuration: durationpb.New(3 * time.Minute),
		}
	}
}

func (x *Settings) Merge(in *Settings) *Settings {
	x.Job = in.Job
	x.Enabled = in.Enabled

	x.Mode = in.Mode
	x.FallbackMode = in.FallbackMode

	x.Public = in.Public

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

	if in.Configuration == nil {
		x.Configuration = &Configuration{
			DeduplicationEnabled:  true,
			DeduplicationRadius:   45,
			DeduplicationDuration: durationpb.New(3 * time.Minute),
		}
	} else {
		x.Configuration = in.Configuration
	}

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
		var dest []string
		if err := json.Unmarshal([]byte(t), &dest); err != nil {
			return err
		}
		for _, job := range dest {
			x.Jobs = append(x.Jobs, &Job{Name: job})
		}
		return nil
	case []byte:
		var dest []string
		if err := json.Unmarshal(t, &dest); err != nil {
			return err
		}
		for _, job := range dest {
			x.Jobs = append(x.Jobs, &Job{Name: job})
		}
		return nil
	}
	return nil
}

// Value marshals the JobList value into driver.Valuer. Special case as only the job list is stored in the database, not the full protobuf message.
func (x *JobList) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	return json.MarshalToString(x.GetJobStrings())
}

func (x *JobList) IsEmpty() bool {
	return x == nil || len(x.Jobs) == 0
}

func (x *JobList) ContainsJob(job string) bool {
	return slices.ContainsFunc(x.Jobs, func(in *Job) bool {
		return in.Name == job
	})
}

func (x *JobList) GetJobStrings() []string {
	jobs := []string{}
	for _, job := range x.Jobs {
		jobs = append(jobs, job.Name)
	}

	return jobs
}

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJob(job string) {
	x.Name = job
}

func (x *Job) SetJobLabel(label string) {
	x.Label = &label
}
