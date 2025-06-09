package centrumstate

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
)

func (s *State) GetDispatchers(ctx context.Context, job string) (*centrum.Dispatchers, error) {
	dispatchers, err := s.dispatchers.GetOrLoad(ctx, job)
	if err != nil {
		return nil, err
	}

	if dispatchers == nil {
		return &centrum.Dispatchers{}, nil
	}

	return dispatchers, nil
}

func (s *State) UpdateDispatchers(ctx context.Context, job string, dispatchers []*jobs.Colleague) error {
	return s.dispatchers.Put(ctx, job, &centrum.Dispatchers{
		Job:         job,
		Dispatchers: dispatchers,
	})
}
