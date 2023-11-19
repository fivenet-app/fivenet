package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
)

func (s *State) GetDisponents(job string) ([]*users.UserShort, error) {
	disponents, err := s.disponents.GetOrLoad(job)
	if err != nil || disponents == nil {
		return nil, err
	}

	if disponents == nil {
		return nil, nil
	}

	return disponents.Disponents, nil
}

func (s *State) UpdateDisponents(job string, disponents []*users.UserShort) error {
	return s.disponents.Put(job, &centrum.Disponents{
		Job:        job,
		Disponents: disponents,
	})
}
