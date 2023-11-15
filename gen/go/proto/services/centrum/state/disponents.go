package state

import "github.com/galexrt/fivenet/gen/go/proto/resources/users"

func (s *State) GetDisponents(job string) []*users.UserShort {
	disponents, ok := s.disponents.Load(job)
	if !ok {
		return nil
	}

	return disponents
}

func (s *State) UpdateDisponents(job string, disponents []*users.UserShort) {
	if job == "" {
		s.disponents.Clear()
		return
	}

	if len(disponents) == 0 {
		s.disponents.Delete(job)
	} else {
		s.disponents.Store(job, disponents)
	}
}
