package units

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/store"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
)

// WatchMembers provides read-only access to membership projection changes.
// Writes remain encapsulated so every component update also refreshes the
// job-scoped feed trigger and the unit supervision timer.
func (s *UnitDB) WatchMembers(
	ctx context.Context,
) (store.IKVWatcher[centrumunits.UnitAssignments, *centrumunits.UnitAssignments], error) {
	if s.membersStore == nil {
		return nil, errors.New("unit membership projection is not initialized")
	}
	return s.membersStore.WatchAll(ctx)
}

func unitSpecFromUnit(unit *centrumunits.Unit) *centrumunits.UnitSpec {
	if unit == nil {
		return nil
	}

	return &centrumunits.UnitSpec{
		Id: unit.Id, CreatedAt: unit.CreatedAt, UpdatedAt: unit.UpdatedAt,
		Job: unit.Job, JobLabel: unit.JobLabel, SortOrder: unit.SortOrder,
		Name: unit.Name, Initials: unit.Initials, Color: unit.Color, Icon: unit.Icon,
		Description: unit.Description, Attributes: unit.Attributes, HomePostal: unit.HomePostal,
		Access: unit.Access,
	}
}

func composeUnit(
	spec *centrumunits.UnitSpec,
	status *centrumunits.UnitStatus,
	members *centrumunits.UnitAssignments,
) *centrumunits.Unit {
	if spec == nil {
		return nil
	}

	unit := &centrumunits.Unit{
		Id: spec.Id, CreatedAt: spec.CreatedAt, UpdatedAt: spec.UpdatedAt,
		Job: spec.Job, JobLabel: spec.JobLabel, SortOrder: spec.SortOrder,
		Name: spec.Name, Initials: spec.Initials, Color: spec.Color, Icon: spec.Icon,
		Description: spec.Description, Attributes: spec.Attributes, HomePostal: spec.HomePostal,
		Access: spec.Access, Status: status,
	}
	if members != nil {
		unit.Users = members.Users
	}
	return unit
}

// Range iterates composed unit projections, including units whose job no
// longer has a Centrum settings entry.
func (s *UnitDB) Range(ctx context.Context, fn func(string, *centrumunits.Unit) bool) {
	s.specStore.Range(func(key string, spec *centrumunits.UnitSpec) bool {
		unit, err := s.Get(ctx, spec.GetId())
		if err != nil || unit == nil {
			return true
		}
		return fn(key, unit)
	})
}

// RangeMembers iterates membership projections without composing unrelated
// spec or status components. It is intended for membership-only maintenance.
func (s *UnitDB) RangeMembers(fn func(string, *centrumunits.UnitAssignments) bool) {
	s.membersStore.Range(fn)
}

func (s *UnitDB) updateSpecInKV(ctx context.Context, id int64, unit *centrumunits.Unit) error {
	if unit == nil {
		return s.deleteInKV(ctx, id)
	}

	key := centrumutils.IdKey(id)
	oldSpec, _ := s.specStore.Get(key)
	if _, err := s.specStore.PutIfChanged(ctx, key, unitSpecFromUnit(unit)); err != nil {
		return err
	}
	if oldSpec != nil && oldSpec.GetJob() != "" && oldSpec.GetJob() != unit.GetJob() {
		if err := s.jobMapping.Delete(
			ctx,
			centrumutils.JobIdKey(oldSpec.GetJob(), id),
		); err != nil &&
			!errors.Is(err, jetstream.ErrKeyNotFound) {
			return err
		}
	}
	return nil
}

func (s *UnitDB) updateInKV(ctx context.Context, id int64, unit *centrumunits.Unit) error {
	if err := s.updateSpecInKV(ctx, id, unit); err != nil {
		return err
	}
	if unit == nil {
		return nil
	}
	key := centrumutils.IdKey(id)
	if err := s.putStatusIfNewer(ctx, key, unit.Status); err != nil {
		return err
	}
	members := &centrumunits.UnitAssignments{UnitId: id, Job: unit.GetJob(), Users: unit.GetUsers()}
	if _, err := s.membersStore.PutIfChanged(ctx, key, members); err != nil {
		return err
	}

	return s.refreshProjectionSideEffects(ctx, id)
}

func (s *UnitDB) deleteInKV(ctx context.Context, id int64) error {
	key := centrumutils.IdKey(id)
	unit, _ := s.Get(ctx, id)
	var errs error
	for _, deleteFn := range []func(context.Context, string) error{
		s.specStore.Delete,
		s.statusStore.Delete,
		s.membersStore.Delete,
	} {
		if err := deleteFn(ctx, key); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			errs = errors.Join(errs, err)
		}
	}
	if unit != nil && unit.GetJob() != "" {
		if err := s.jobMapping.Delete(
			ctx,
			centrumutils.JobIdKey(unit.GetJob(), id),
		); err != nil &&
			!errors.Is(err, jetstream.ErrKeyNotFound) {
			errs = errors.Join(errs, err)
		}
	}
	if s.KVPing != nil {
		if err := s.KVPing.Delete(
			ctx,
			"ping."+strconv.FormatInt(id, 10),
		); err != nil &&
			!errors.Is(err, jetstream.ErrKeyNotFound) {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (s *UnitDB) Get(ctx context.Context, id int64) (*centrumunits.Unit, error) {
	key := centrumutils.IdKey(id)
	spec, err := s.specStore.GetOrLoad(ctx, key)
	if err != nil {
		return nil, err
	}
	status, err := s.statusStore.GetOrLoad(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, err
	}
	members, err := s.membersStore.GetOrLoad(ctx, key)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, err
	}
	return composeUnit(spec, status, members), nil
}

func (s *UnitDB) List(ctx context.Context, jobs []string) []*centrumunits.Unit {
	if len(jobs) == 0 {
		jobs = []string{""}
	}
	keys := s.jobMapping.KeysFiltered("", func(key string) bool {
		return slices.ContainsFunc(
			jobs,
			func(job string) bool { return strings.HasPrefix(key, job+".") },
		)
	})
	units := make([]*centrumunits.Unit, 0, len(keys))
	for _, key := range keys {
		id, err := centrumutils.ExtractID(key)
		if err != nil {
			continue
		}
		unit, err := s.Get(ctx, id)
		if err == nil && unit != nil {
			units = append(units, unit)
		}
	}
	slices.SortFunc(units, compareUnits)
	return units
}

func (s *UnitDB) Filter(
	ctx context.Context,
	jobs []string,
	statuses, notStatuses []centrumunits.StatusUnit,
	filterFn func(*centrumunits.Unit) bool,
) []*centrumunits.Unit {
	units := s.List(ctx, jobs)
	units = slices.DeleteFunc(units, func(unit *centrumunits.Unit) bool {
		if len(statuses) > 0 && unit.GetStatus() != nil &&
			!slices.Contains(statuses, unit.GetStatus().GetStatus()) {
			return true
		}
		if len(notStatuses) > 0 && unit.GetStatus() != nil &&
			slices.Contains(notStatuses, unit.GetStatus().GetStatus()) {
			return true
		}
		return filterFn != nil && !filterFn(unit)
	})
	slices.SortFunc(units, compareUnits)
	return units
}

func compareUnits(a, b *centrumunits.Unit) int {
	if a.GetSortOrder() < b.GetSortOrder() {
		return -1
	}
	if a.GetSortOrder() > b.GetSortOrder() {
		return 1
	}
	if c := strings.Compare(a.GetName(), b.GetName()); c != 0 {
		return c
	}
	if a.GetId() < b.GetId() {
		return -1
	}
	if a.GetId() > b.GetId() {
		return 1
	}
	return 0
}

func (s *UnitDB) updateStatusInKV(
	ctx context.Context,
	id int64,
	status *centrumunits.UnitStatus,
) error {
	if status == nil {
		return nil
	}
	if err := s.putStatusIfNewer(ctx, centrumutils.IdKey(id), status); err != nil {
		return err
	}
	return s.refreshProjectionSideEffects(ctx, id)
}

// putStatusIfNewer applies a database snapshot without regressing a status
// written by a concurrent runtime update. The comparison must happen inside
// ComputeUpdate so it is repeated after a KV revision conflict.
func (s *UnitDB) putStatusIfNewer(
	ctx context.Context,
	key string,
	status *centrumunits.UnitStatus,
) error {
	_, err := s.statusStore.ComputeUpdate(
		ctx,
		key,
		func(_ string, existing *centrumunits.UnitStatus) (*centrumunits.UnitStatus, bool, error) {
			if status == nil {
				return existing, false, nil
			}
			if existing != nil && existing.GetId() >= status.GetId() {
				return existing, false, nil
			}
			return status, true, nil
		},
	)
	return err
}

func (s *UnitDB) refreshProjectionSideEffects(ctx context.Context, id int64) error {
	unit, err := s.Get(ctx, id)
	if err != nil || unit == nil {
		return err
	}
	// The job mapping is the feed trigger. Rewrite it even when its value is
	// unchanged so status/member component updates reach subscribed clients.
	if err := s.jobMapping.Put(
		ctx,
		centrumutils.JobIdKey(unit.GetJob(), id),
		&common.IDMapping{Id: id},
	); err != nil {
		return err
	}
	return s.SyncUnitPing(ctx, unit)
}
