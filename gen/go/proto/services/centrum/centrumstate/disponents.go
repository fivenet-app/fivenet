package centrumstate

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
)

func (s *State) GetDisponents(ctx context.Context, job string) ([]*users.UserShort, error) {
	disponents, err := s.disponents.GetOrLoad(ctx, job)
	if err != nil || disponents == nil {
		return nil, err
	}

	return disponents.Disponents, nil
}

func (s *State) UpdateDisponents(ctx context.Context, job string, disponents []*users.UserShort) error {
	return s.disponents.Put(ctx, job, &centrum.Disponents{
		Job:        job,
		Disponents: disponents,
	})
}
