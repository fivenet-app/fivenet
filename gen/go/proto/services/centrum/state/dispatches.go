package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/utils"
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

func (s *State) FilterDispatches(job string, statuses []dispatch.StatusDispatch, notStatuses []dispatch.StatusDispatch) []*dispatch.Dispatch {
	dsps := []*dispatch.Dispatch{}
	dispatches := s.ListDispatches(job)
	for i := 0; i < len(dispatches); i++ {
		// Hide user info when dispatch is anonymous
		if dispatches[i].Anon {
			dispatches[i].Creator = nil
			dispatches[i].CreatorId = nil
		}

		include := true

		// Include statuses that should be listed
		if len(statuses) > 0 && !utils.InSlice(statuses, dispatches[i].Status.Status) {
			include = false
		} else if len(notStatuses) > 0 {
			// Which statuses to ignore
			for _, status := range notStatuses {
				if dispatches[i].Status != nil && dispatches[i].Status.Status == status {
					include = false
					break
				}
			}
		}

		if include {
			dsps = append(dsps, dispatches[i])
		}
	}

	return dsps
}
