package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"golang.org/x/exp/slices"
)

func (s *State) GetDispatch(job string, id uint64) (*dispatch.Dispatch, bool) {
	dispatches, ok := s.Dispatches.Load(job)
	if !ok {
		return nil, false
	}

	return dispatches.Load(id)
}

func (s *State) ListDispatches(job string) []*dispatch.Dispatch {
	ds := []*dispatch.Dispatch{}

	dispatches, ok := s.Dispatches.Load(job)
	if !ok {
		return nil
	}

	dispatches.Range(func(id uint64, dispatch *dispatch.Dispatch) bool {
		ds = append(ds, dispatch)
		return true
	})

	slices.SortFunc(ds, func(a, b *dispatch.Dispatch) int {
		return int(b.Id - a.Id)
	})

	return ds
}
