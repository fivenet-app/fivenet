package mstlystcdata

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/nats-io/nats.go/jetstream"
)

type DummyJobs struct {
	Entries map[string]*jobs.Job
}

func NewDummyJobs(entries map[string]*jobs.Job) *DummyJobs {
	if entries == nil {
		entries = map[string]*jobs.Job{}
	}

	return &DummyJobs{Entries: entries}
}

func (j *DummyJobs) Get(job string) (*jobs.Job, error) {
	if j == nil {
		return nil, jetstream.ErrKeyNotFound
	}

	if out, ok := j.Entries[job]; ok {
		return out, nil
	}

	return nil, jetstream.ErrKeyNotFound
}

func (j *DummyJobs) List() []*jobs.Job {
	if j == nil || len(j.Entries) == 0 {
		return []*jobs.Job{}
	}

	out := make([]*jobs.Job, 0, len(j.Entries))
	for _, job := range j.Entries {
		out = append(out, job)
	}

	return out
}

func (j *DummyJobs) Range(fn func(key string, value *jobs.Job) bool) {
	if j == nil {
		return
	}

	for key, value := range j.Entries {
		if !fn(key, value) {
			return
		}
	}
}

func (j *DummyJobs) Has(job string) bool {
	if j == nil {
		return false
	}

	_, ok := j.Entries[job]
	return ok
}

func (j *DummyJobs) GetHighestJobGrade(job string) *jobs.JobGrade {
	out, err := j.Get(job)
	if err != nil || out == nil || len(out.GetGrades()) == 0 {
		return nil
	}

	return out.GetGrades()[len(out.GetGrades())-1]
}
