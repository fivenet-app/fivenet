package state

import "github.com/galexrt/fivenet/gen/go/proto/resources/users"

func (s *State) GetDisponents(job string) []*users.UserShort {
	disponents, ok := s.Disponents.Load(job)
	if !ok {
		return nil
	}

	return disponents
}
