package centrum

import (
	"database/sql/driver"
	"encoding/json"
	"slices"

	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
)

func (x *Dispatchers) SetJobLabel(jobLabel string) {
	x.JobLabel = &jobLabel
}

func (x *Dispatchers) Merge(in *Dispatchers) *Dispatchers {
	if len(in.GetDispatchers()) == 0 {
		x.Dispatchers = []*jobs.Colleague{}
	} else {
		x.Dispatchers = in.GetDispatchers()
	}

	return x
}

func (x *Dispatchers) IsEmpty() bool {
	return x == nil || len(x.GetDispatchers()) == 0
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
			x.Jobs = append(x.Jobs, &JobListEntry{Name: job})
		}
		return nil
	case []byte:
		var dest []string
		if err := json.Unmarshal(t, &dest); err != nil {
			return err
		}
		for _, job := range dest {
			x.Jobs = append(x.Jobs, &JobListEntry{Name: job})
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
	return slices.ContainsFunc(x.GetJobs(), func(in *JobListEntry) bool {
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

func (x *JobListEntry) GetJob() string {
	return x.GetName()
}

func (x *JobListEntry) SetJob(job string) {
	x.Name = job
}

func (x *JobListEntry) SetJobLabel(label string) {
	x.Label = &label
}
