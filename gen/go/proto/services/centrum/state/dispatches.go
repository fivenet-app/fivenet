package state

import (
	"context"
	"slices"
	"strconv"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/paulmach/orb"
	"go.uber.org/zap"
)

func (s *State) GetDispatch(job string, id uint64) (*centrum.Dispatch, error) {
	return s.dispatches.GetOrLoad(JobIdKey(job, id))
}

func (s *State) ListDispatches(job string) ([]*centrum.Dispatch, bool) {
	ds := []*centrum.Dispatch{}

	ids, err := s.dispatches.Keys(job)
	if err != nil {
		return ds, false
	}

	for _, id := range ids {
		dsp, err := s.dispatches.GetOrLoad(id)
		if err != nil {
			continue
		}

		// Remove broken dispatches
		if dsp.Id == 0 || dsp.Job == "" {
			dId, err := strconv.Atoi(id)
			if err != nil {
				return ds, false
			}

			if err := s.DeleteDispatch(job, uint64(dId)); err != nil {
				return ds, false
			}

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
		if len(statuses) > 0 && !slices.Contains(statuses, dispatches[i].Status.Status) {
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
	dsp, err := s.GetDispatch(job, id)
	if err != nil {
		return err
	}

	if locs := s.GetDispatchLocations(job); locs != nil {
		locs.Remove(dsp, func(p orb.Pointer) bool {
			return p.(*centrum.Dispatch).Id == dsp.Id
		})
	}

	return s.dispatches.Delete(JobIdKey(job, id))
}

func (s *State) CreateDispatch(ctx context.Context, job string, id uint64, dsp *centrum.Dispatch) error {
	return s.UpdateDispatch(ctx, job, id, dsp)
}

func (s *State) UpdateDispatch(ctx context.Context, job string, id uint64, dsp *centrum.Dispatch) error {
	if err := s.dispatches.ComputeUpdate(JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, error) {
		if existing == nil {
			// Dispatch must not be existant yet, so make sure to add to the
			// dispatch locations
			if locs := s.GetDispatchLocations(dsp.Job); locs != nil {
				if err := locs.Add(dsp); err != nil {
					s.logger.Error("failed to add non-existant dispatch to locations", zap.Uint64("dispatch_id", dsp.Id))
				}
			}

			return dsp, nil
		}

		locsReplace := false
		// Make sure to update the existing dispatch location if the dispatch
		// has its location changed
		if existing.X != dsp.X || existing.Y != dsp.Y {
			locsReplace = true
		}

		existing.Merge(dsp)

		if locsReplace {
			if locs := s.GetDispatchLocations(existing.Job); locs != nil {
				if err := locs.Replace(existing, func(p orb.Pointer) bool {
					return p.(*centrum.Dispatch).Id == existing.Id
				}); err != nil {
					s.logger.Error("failed to replace updated dispatch's in locations tree", zap.Error(err))
				}
			}
		}

		return existing, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *State) UpdateDispatchStatus(ctx context.Context, job string, id uint64, status *centrum.DispatchStatus) error {
	if err := s.dispatches.ComputeUpdate(JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, error) {
		if existing == nil {
			return nil, nil
		}

		existing.Status = status

		if existing.Status != nil {
			// Remove completed dispatches from locations
			if centrumutils.IsStatusDispatchComplete(existing.Status.Status) {
				if locs := s.GetDispatchLocations(job); locs != nil {
					locs.Remove(existing, func(p orb.Pointer) bool {
						return p.(*centrum.Dispatch).Id == existing.Id
					})
				}
			}
		}

		return existing, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *State) UpdateDispatchUnits(ctx context.Context, job string, id uint64, units []*centrum.DispatchAssignment) error {
	if err := s.dispatches.ComputeUpdate(JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, error) {
		if existing == nil {
			return nil, nil
		}

		if units == nil {
			existing.Units = []*centrum.DispatchAssignment{}
		} else {
			existing.Units = units
		}

		return existing, nil
	}); err != nil {
		return err
	}

	return nil
}
