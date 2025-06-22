package dispatches

import (
	"context"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"google.golang.org/protobuf/proto"
)

func (s *DispatchDB) updateInKV(ctx context.Context, id uint64, dsp *centrum.Dispatch) error {
	if err := s.store.ComputeUpdate(ctx, centrumutils.IdKey(id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
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

func (s *DispatchDB) deleteInKV(ctx context.Context, id uint64) error {
	return s.store.Delete(ctx, centrumutils.IdKey(id))
}

func (s *DispatchDB) Get(ctx context.Context, id uint64) (*centrum.Dispatch, error) {
	dsp, err := s.store.GetOrLoad(ctx, centrumutils.IdKey(id))
	if err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *DispatchDB) List(ctx context.Context, jobs []string) []*centrum.Dispatch {
	if jobs == nil {
		jobs = []string{""}
	}

	out := []*centrum.Dispatch{}

	for _, job := range jobs {
		keys := s.jobMapping.Keys(job)

		ds := []*centrum.Dispatch{}
		for _, key := range keys {
			did, err := centrumutils.ExtractIDString(key)
			if err != nil {
				continue
			}
			dsp, err := s.store.GetOrLoad(ctx, did)
			if err != nil {
				continue
			}

			// Skip any broken dispatches (e.g. missing job or ID)
			if dsp.Id == 0 || len(dsp.Jobs.GetJobs()) == 0 {
				continue
			}

			ds = append(ds, dsp)
		}

		out = append(out, ds...)
	}

	slices.SortFunc(out, func(a, b *centrum.Dispatch) int {
		return int(a.Id - b.Id)
	})

	return out
}

func (s *DispatchDB) Filter(ctx context.Context, jobs []string, statuses []centrum.StatusDispatch, notStatuses []centrum.StatusDispatch) []*centrum.Dispatch {
	ds := s.List(ctx, jobs)

	ds = slices.DeleteFunc(ds, func(dispatch *centrum.Dispatch) bool {
		// Hide user info when dispatch is anonymous
		if dispatch.Anon {
			dispatch.Creator = nil
		}

		// Include statuses that should be listed
		if len(statuses) > 0 && !slices.Contains(statuses, dispatch.Status.Status) {
			return true
		} else if len(notStatuses) > 0 && dispatch.Status != nil {
			// Which statuses to ignore
			if slices.Contains(notStatuses, dispatch.Status.Status) {
				return true
			}
		}

		return false
	})

	return ds
}

func (s *DispatchDB) updateStatusInKV(ctx context.Context, id uint64, status *centrum.DispatchStatus) error {
	if err := s.store.ComputeUpdate(ctx, centrumutils.IdKey(id), true, func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
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
