package centrumstate

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/nats-io/nats.go/jetstream"
)

func (s *State) GetDisponents(ctx context.Context, job string) (*centrum.Disponents, error) {
	disponents, err := s.disponents.GetOrLoad(ctx, job)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, err
		}
	}

	if disponents == nil {
		return &centrum.Disponents{
			Job: job,
		}, nil
	}

	return disponents, nil
}

func (s *State) GetDisponentsJobs(ctx context.Context) []string {
	jobs := []string{}
	s.disponents.Range(ctx, func(key string, _ *centrum.Disponents) bool {
		jobs = append(jobs, key)
		return true
	})

	return jobs
}

func (s *State) UpdateDisponents(ctx context.Context, job string, disponents []*jobs.Colleague) error {
	d := &centrum.Disponents{
		Job:        job,
		Disponents: disponents,
	}
	s.enricher.EnrichJobName(d)

	return s.disponents.Put(ctx, job, d)
}
