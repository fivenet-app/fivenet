package dispatches

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/store"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

func (s *DispatchDB) Store() *store.Store[centrum.Dispatch, *centrum.Dispatch] {
	return s.store
}

func (s *DispatchDB) IdleStore() jetstream.KeyValue {
	return s.idleKV
}

func (s *DispatchDB) updateInKV(ctx context.Context, id uint64, dsp *centrum.Dispatch) error {
	if err := s.store.ComputeUpdate(ctx, centrumutils.IdKey(id), func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
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

// List returns all dispatches that match the given job prefixes.
// If jobs is nil, it returns all dispatches.
func (s *DispatchDB) List(ctx context.Context, jobs []string) []*centrum.Dispatch {
	if jobs == nil {
		jobs = []string{""}
	}

	keys := s.jobMapping.KeysFiltered("", func(key string) bool {
		for _, job := range jobs {
			if strings.HasPrefix(key, job+".") {
				return true
			}
		}

		return false
	})

	ds := []*centrum.Dispatch{}
	for _, key := range keys {
		uid, err := centrumutils.ExtractIDString(key)
		if err != nil {
			continue
		}

		dsp, err := s.store.GetOrLoad(ctx, uid)
		if err != nil {
			continue
		}
		// Skip any broken dispatches (e.g. missing job or ID)
		if dsp == nil || dsp.Id == 0 || len(dsp.Jobs.GetJobs()) == 0 {
			continue
		}
		ds = append(ds, dsp)
	}

	slices.SortFunc(ds, func(a, b *centrum.Dispatch) int {
		return int(a.Id - b.Id)
	})

	return ds
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
	if err := s.store.ComputeUpdate(ctx, centrumutils.IdKey(id), func(key string, existing *centrum.Dispatch) (*centrum.Dispatch, bool, error) {
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

const (
	InactiveTTL            = 60 * time.Minute // How long a Dispatch may stay idle/quiet
	InactiveLimitMarkerTTL = 2 * time.Hour    // How long tombstones live

)

func (d *DispatchDB) TouchActivity(ctx context.Context, id uint64) error {
	key := "idle." + centrumutils.IdKey(id)

	// First try to create the timer key with TTL.
	if _, err := d.idleKV.Create(
		ctx, key, nil,
		jetstream.KeyTTL(InactiveTTL),
	); err != nil {
		if !errors.Is(err, jetstream.ErrKeyExists) {
			return err // a real failure
		}

		// Key already exists → refresh its TTL with Update.
		// We need the current revision to comply with optimistic locking.
		entry, err := d.idleKV.Get(ctx, key)
		if err != nil {
			return err
		}

		// Update resets the previously-set TTL back to the full 60 min.
		if _, err = d.idleKV.Update(ctx, key, nil, entry.Revision()); err != nil {
			return err
		}
	}

	return nil
}
