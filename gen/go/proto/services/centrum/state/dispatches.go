package state

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/paulmach/orb"
	"golang.org/x/exp/slices"
)

func (s *State) GetDispatch(job string, id uint64) *centrum.Dispatch {
	d, _ := s.dispatches.Get(JobIdKey(job, id))
	return d
}

func (s *State) ListDispatches(job string) ([]*centrum.Dispatch, bool) {
	ds := []*centrum.Dispatch{}

	ids, err := s.dispatches.Keys(job)
	if err != nil {
		return ds, false
	}

	for _, id := range ids {
		dsp, ok := s.dispatches.Get(id)
		if !ok || dsp == nil {
			continue
		}

		ds = append(ds, dsp)
	}

	slices.SortFunc(ds, func(a, b *centrum.Dispatch) int {
		return int(a.Id - b.Id)
	})

	return ds, true
}

func (s *State) FilterDispatches(job string, statuses []centrum.StatusDispatch, notStatuses []centrum.StatusDispatch) []*centrum.Dispatch {
	dispatches, ok := s.ListDispatches(job)
	if !ok {
		return nil
	}

	dsps := []*centrum.Dispatch{}
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

func (s *State) DeleteDispatch(job string, id uint64) error {
	dsp := s.GetDispatch(job, id)
	if dsp != nil {
		s.GetDispatchLocations(job).Remove(dsp, func(p orb.Pointer) bool {
			return p.(*centrum.Dispatch).Id == dsp.Id
		})
	}

	return s.dispatches.Delete(JobIdKey(job, id))
}

func (s *State) UpdateDispatch(ctx context.Context, job string, id uint64, dsp *centrum.Dispatch) error {
	s.dispatches.ComputeUpdate(JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, error) {
		if existing == nil {
			return dsp, nil
		}

		existing.Merge(dsp)

		return existing, nil
	})

	return nil
}

func (s *State) UpdateDispatchStatus(ctx context.Context, job string, id uint64, status *centrum.DispatchStatus) error {
	s.dispatches.ComputeUpdate(JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, error) {
		if existing == nil {
			return nil, nil
		}

		existing.Status = status

		return existing, nil
	})

	return nil
}
