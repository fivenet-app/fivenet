package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/paulmach/orb"
	"github.com/puzpuzpuz/xsync/v3"
	"golang.org/x/exp/slices"
)

func (s *State) GetDispatch(job string, id uint64) (*dispatch.Dispatch, bool) {
	dispatches, ok := s.dispatches.Load(job)
	if !ok {
		return nil, false
	}

	return dispatches.Load(id)
}

func (s *State) ListDispatches(job string) []*dispatch.Dispatch {
	ds := []*dispatch.Dispatch{}

	dispatches, ok := s.dispatches.Load(job)
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

func (s *State) DeleteDispatch(job string, id uint64) {
	dsp, ok := s.GetDispatch(job, id)
	if !ok || dsp == nil {
		return
	}

	s.GetDispatchLocations(job).Remove(dsp, func(p orb.Pointer) bool {
		return p.(*dispatch.Dispatch).Id == dsp.Id
	})

	dispatches, ok := s.dispatches.Load(job)
	if ok {
		dispatches.Delete(id)
	}
}

func (s *State) UpdateDispatch(job string, dispatchId uint64, dsp *dispatch.Dispatch) error {
	if dispatch, ok := s.GetDispatchesMap(job).LoadOrStore(dispatchId, dsp); ok {
		dispatch.Merge(dsp)
	}

	return nil
}

func (s *State) GetDispatchesJobs() []string {
	list := []string{}

	s.dispatches.Range(func(job string, _ *xsync.MapOf[uint64, *dispatch.Dispatch]) bool {
		list = append(list, job)
		return true
	})

	return list
}
