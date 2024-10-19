package centrumstate

import (
	"context"
	"slices"
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"google.golang.org/protobuf/proto"
)

func (s *State) GetDispatch(ctx context.Context, job string, id uint64) (*centrum.Dispatch, error) {
	return s.dispatches.GetOrLoad(ctx, JobIdKey(job, id))
}

func (s *State) ListDispatches(ctx context.Context, job string) ([]*centrum.Dispatch, bool) {
	ds := []*centrum.Dispatch{}

	ids, err := s.dispatches.Keys(ctx, job)
	if err != nil {
		return ds, false
	}

	for _, id := range ids {
		dsp, err := s.dispatches.GetOrLoad(ctx, id)
		if err != nil {
			continue
		}

		// Remove broken dispatches
		if dsp.Id == 0 || dsp.Job == "" {
			dId, err := strconv.Atoi(id)
			if err != nil {
				return ds, false
			}

			if err := s.DeleteDispatch(ctx, job, uint64(dId)); err != nil {
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

func (s *State) FilterDispatches(ctx context.Context, job string, statuses []centrum.StatusDispatch, notStatuses []centrum.StatusDispatch) []*centrum.Dispatch {
	dispatches, ok := s.ListDispatches(ctx, job)
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

func (s *State) DeleteDispatch(ctx context.Context, job string, id uint64) error {
	return s.dispatches.Delete(ctx, JobIdKey(job, id))
}

func (s *State) CreateDispatch(ctx context.Context, job string, id uint64, dsp *centrum.Dispatch) error {
	return s.UpdateDispatch(ctx, job, id, dsp)
}

func (s *State) UpdateDispatch(ctx context.Context, job string, id uint64, dsp *centrum.Dispatch) error {
	if err := s.dispatches.ComputeUpdate(ctx, JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
		if existing == nil {
			return dsp, dsp != nil, nil
		}

		if !proto.Equal(existing, dsp) {
			existing.Merge(dsp)
			return existing, true, nil
		}

		return existing, false, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *State) UpdateDispatchStatus(ctx context.Context, job string, id uint64, status *centrum.DispatchStatus) error {
	if err := s.dispatches.ComputeUpdate(ctx, JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
		if existing == nil {
			return existing, false, nil
		}

		existing.Status = status

		return existing, true, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *State) UpdateDispatchUnits(ctx context.Context, job string, id uint64, units []*centrum.DispatchAssignment) error {
	if err := s.dispatches.ComputeUpdate(ctx, JobIdKey(job, id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
		if existing == nil {
			return existing, false, nil
		}

		if units == nil {
			existing.Units = []*centrum.DispatchAssignment{}
		} else {
			existing.Units = units
		}

		return existing, true, nil
	}); err != nil {
		return err
	}

	return nil
}
