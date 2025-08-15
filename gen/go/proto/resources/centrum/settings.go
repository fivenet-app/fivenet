package centrum

import (
	"database/sql/driver"
	"encoding/json"
	"slices"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

func (x *Settings) Default(job string) {
	x.Job = job

	if x.GetMode() <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.Mode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.GetFallbackMode() <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.FallbackMode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.GetPredefinedStatus() == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	}

	if x.GetTimings() == nil {
		x.Timings = &Timings{}
	}
	if x.GetTimings().GetDispatchMaxWait() <= 0 {
		x.Timings.DispatchMaxWait = 900
	}
	if x.GetTimings().GetRequireUnitReminderSeconds() <= 0 {
		x.Timings.RequireUnitReminderSeconds = 180
	}

	if x.GetConfiguration() == nil {
		x.Configuration = &Configuration{
			DeduplicationEnabled:  true,
			DeduplicationRadius:   45,
			DeduplicationDuration: durationpb.New(3 * time.Minute),
		}
	}
}

func (x *Settings) Merge(in *Settings) *Settings {
	x.Job = in.GetJob()
	x.Enabled = in.GetEnabled()

	x.Mode = in.GetMode()
	x.FallbackMode = in.GetFallbackMode()

	x.Public = in.GetPublic()

	if in.GetPredefinedStatus() == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	} else {
		x.PredefinedStatus = in.GetPredefinedStatus()
	}

	if in.GetTimings() == nil {
		x.Timings = &Timings{}
	} else {
		x.Timings = in.GetTimings()
	}

	if in.GetAccess() == nil {
		x.Access = &CentrumAccess{}
	} else {
		x.Access = in.GetAccess()
	}

	if in.GetConfiguration() == nil {
		x.Configuration = &Configuration{
			DeduplicationEnabled:  true,
			DeduplicationRadius:   45,
			DeduplicationDuration: durationpb.New(3 * time.Minute),
		}
	} else {
		x.Configuration = in.GetConfiguration()
	}

	return x
}

func (x *Settings) JobHasAccess(job string, access CentrumAccessLevel) bool {
	if x.GetAccess() == nil || x.Access.Jobs == nil {
		return x.GetJob() == job // No access restrictions defined, only the job itself is allowed
	}

	if x.GetJob() == job {
		return true // Job is explicitly allowed
	}

	for _, j := range x.GetAccess().GetJobs() {
		if j.GetJob() == job && j.GetAccess() >= access {
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

	out, err := json.Marshal(x.GetJobStrings())
	if err != nil {
		return nil, err
	}
	return string(out), nil
}

func (x *JobList) IsEmpty() bool {
	return x == nil || len(x.GetJobs()) == 0
}

func (x *JobList) ContainsJob(job string) bool {
	return slices.ContainsFunc(x.GetJobs(), func(in *Job) bool {
		return in.GetName() == job
	})
}

func (x *JobList) GetJobStrings() []string {
	jobs := []string{}
	for _, job := range x.GetJobs() {
		jobs = append(jobs, job.GetName())
	}

	return jobs
}

func (x *Job) GetJob() string {
	return x.GetName()
}

func (x *Job) SetJob(job string) {
	x.Name = job
}

func (x *Job) SetJobLabel(label string) {
	x.Label = &label
}
