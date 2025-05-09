package centrumstate

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
)

func (s *State) GetDisponents(ctx context.Context, job string) ([]*jobs.Colleague, error) {
	disponents, err := s.disponents.GetOrLoad(ctx, job)
	if err != nil || disponents == nil {
		return nil, err
	}

	return disponents.Disponents, nil
}

func (s *State) UpdateDisponents(ctx context.Context, job string, disponents []*jobs.Colleague) error {
	return s.disponents.Put(ctx, job, &centrum.Disponents{
		Job:        job,
		Disponents: disponents,
	})
}
