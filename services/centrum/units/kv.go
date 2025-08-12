package units

import (
	"context"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"google.golang.org/protobuf/proto"
)

func (s *UnitDB) updateInKV(ctx context.Context, id uint64, unit *centrum.Unit) error {
	if err := s.store.ComputeUpdate(ctx, centrumutils.IdKey(id), func(key string, existing *centrum.Unit) (*centrum.Unit, bool, error) {
		if existing == nil {
			return unit, unit != nil, nil
		}

		if !proto.Equal(existing, unit) {
			existing.Merge(unit)
			return existing, true, nil
		}

		return existing, false, nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *UnitDB) deleteInKV(ctx context.Context, id uint64) error {
	return s.store.Delete(ctx, centrumutils.IdKey(id))
}

func (s *UnitDB) Get(ctx context.Context, id uint64) (*centrum.Unit, error) {
	return s.store.GetOrLoad(ctx, centrumutils.IdKey(id))
}

func (s *UnitDB) List(ctx context.Context, jobs []string) []*centrum.Unit {
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

	us := []*centrum.Unit{}
	for _, key := range keys {
		uid, err := centrumutils.ExtractIDString(key)
		if err != nil {
			continue
		}

		unit, err := s.store.Get(uid)
		if err != nil {
			continue
		}
		if unit == nil {
			continue
		}
		us = append(us, unit)
	}

	slices.SortFunc(us, func(a, b *centrum.Unit) int {
		return strings.Compare(a.GetName(), b.GetName())
	})

	return us
}

func (s *UnitDB) Filter(
	ctx context.Context,
	jobs []string,
	statuses []centrum.StatusUnit,
	notStatuses []centrum.StatusUnit,
	filterFn func(unit *centrum.Unit) bool,
) []*centrum.Unit {
	us := s.List(ctx, jobs)

	us = slices.DeleteFunc(us, func(unit *centrum.Unit) bool {
		// Include statuses that should be listed
		if len(statuses) > 0 && unit.GetStatus() != nil &&
			!slices.Contains(statuses, unit.GetStatus().GetStatus()) {
			return true
		} else if len(notStatuses) > 0 && unit.GetStatus() != nil {
			// Which statuses to ignore
			if slices.Contains(notStatuses, unit.GetStatus().GetStatus()) {
				return true
			}
		}

		if filterFn == nil || filterFn(unit) {
			return false
		}

		return true
	})

	return us
}

func (s *UnitDB) updateStatusInKV(
	ctx context.Context,
	id uint64,
	status *centrum.UnitStatus,
) error {
	if err := s.store.ComputeUpdate(ctx, centrumutils.IdKey(id), func(key string, existing *centrum.Unit) (*centrum.Unit, bool, error) {
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
