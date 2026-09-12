package units

import (
	"context"
	"errors"
	"fmt"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *UnitDB) SyncUserUnitMapping(ctx context.Context, userId int32) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user ID: %d", userId)
	}

	currentMapping, hasMapping, err := s.tracker.GetUserMapping(userId)
	if err != nil {
		return err
	}

	var currentUnitId int64
	if currentMapping != nil && currentMapping.UnitId != nil {
		currentUnitId = currentMapping.GetUnitId()
	}

	unitId, err := s.LoadUnitIDForUserID(ctx, userId)
	if err != nil {
		return err
	}
	if unitId > 0 {
		unit, err := s.Get(ctx, unitId)
		if err != nil {
			return err
		}

		marker, markerFound := s.tracker.GetUserMarkerById(userId)
		inJob, err := s.UserInJob(ctx, s.db, unit.GetJob(), userId)
		if err != nil {
			return err
		}

		// Unit membership is tied to the user's current duty context. Do not
		// retain an old-job or off-duty assignment until the periodic audit runs.
		if !markerFound || marker == nil || marker.GetHidden() ||
			!s.tracker.IsUserOnDuty(userId) || marker.GetJob() != unit.GetJob() || !inJob {
			if err := s.RemoveUnitAssignments(
				ctx,
				"",
				&userId,
				unitId,
				[]int32{userId},
			); err != nil {
				return err
			}
			unitId = 0
		}
	}

	var targetUnitId *int64
	if unitId > 0 && s.tracker.IsUserOnDuty(userId) {
		targetUnitId = &unitId
	}

	switch {
	case targetUnitId != nil:
		if err := s.tracker.SetUserMappingForUser(ctx, userId, targetUnitId); err != nil {
			return err
		}
	case s.tracker.IsUserOnDuty(userId):
		if err := s.tracker.UnsetUnitIDForUser(ctx, userId); err != nil {
			return err
		}
	case hasMapping:
		if err := s.tracker.DeleteUserMapping(ctx, userId); err != nil {
			return err
		}
	}

	var errs error
	if currentUnitId > 0 && currentUnitId != unitId {
		errs = errors.Join(errs, s.SyncUnitMembership(ctx, currentUnitId))
	}
	if unitId > 0 {
		errs = errors.Join(errs, s.SyncUnitMembership(ctx, unitId))
	}

	return errs
}

func (s *UnitDB) SyncUnitMembership(ctx context.Context, unitId int64) error {
	if unitId <= 0 {
		return fmt.Errorf("invalid unit ID: %d", unitId)
	}

	units, err := s.loadUnitsFromDB(ctx, unitId)
	if err != nil {
		return err
	}

	if len(units) == 0 {
		return s.syncMissingUnitMembership(ctx, unitId)
	}

	return s.syncLoadedUnitMembership(ctx, units[0])
}

func (s *UnitDB) syncLoadedUnitMembership(ctx context.Context, unit *centrumunits.Unit) error {
	if unit == nil {
		return nil
	}

	unitId := unit.GetId()
	previous, err := s.store.Get(centrumutils.IdKey(unitId))
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}
	if err := s.updateInKV(ctx, unitId, unit); err != nil {
		return err
	}

	userIds := make(map[int32]struct{}, len(unit.GetUsers()))
	var errs error
	for _, user := range unit.GetUsers() {
		userId := user.GetUserId()
		if userId <= 0 {
			continue
		}

		eligible, err := s.isEligibleUnitMember(ctx, unit.GetJob(), userId)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if !eligible {
			errs = errors.Join(errs, s.clearTrackerMappingForUnit(ctx, userId, unitId))
			continue
		}

		userIds[userId] = struct{}{}
		if err := s.tracker.SetUserMappingForUser(ctx, userId, &unit.Id); err != nil {
			s.logger.Error(
				"failed to set user's unit id",
				zap.Int64("unit_id", unitId),
				zap.Int32("user_id", userId),
				zap.Error(err),
			)
			errs = errors.Join(errs, err)
		}
	}

	return errors.Join(
		errs,
		s.clearRemovedTrackerMappingsForUnit(ctx, unitId, previous.GetUsers(), userIds),
	)
}

// isEligibleUnitMember keeps tracker mappings tied to the active duty
// context. Durable membership is reconciled separately by user changes.
func (s *UnitDB) isEligibleUnitMember(ctx context.Context, job string, userId int32) (bool, error) {
	marker, found := s.tracker.GetUserMarkerById(userId)
	if !found || marker == nil || marker.GetHidden() || !s.tracker.IsUserOnDuty(userId) ||
		marker.GetJob() != job {
		return false, nil
	}

	return s.UserInJob(ctx, s.db, job, userId)
}

func (s *UnitDB) syncMissingUnitMembership(ctx context.Context, unitId int64) error {
	previous, err := s.store.Get(centrumutils.IdKey(unitId))
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return err
	}

	var errs error
	if err := s.deleteInKV(ctx, unitId); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		errs = errors.Join(errs, err)
	}

	return errors.Join(
		errs,
		s.clearRemovedTrackerMappingsForUnit(ctx, unitId, previous.GetUsers(), nil),
	)
}

func (s *UnitDB) clearRemovedTrackerMappingsForUnit(
	ctx context.Context,
	unitId int64,
	previousUsers []*centrumunits.UnitAssignment,
	validUserIds map[int32]struct{},
) error {
	var errs error
	for _, user := range previousUsers {
		userId := user.GetUserId()
		if userId <= 0 {
			continue
		}

		if _, ok := validUserIds[userId]; ok {
			continue
		}

		errs = errors.Join(errs, s.clearTrackerMappingForUnit(ctx, userId, unitId))
	}

	return errs
}

func (s *UnitDB) clearTrackerMappingForUnit(ctx context.Context, userId int32, unitId int64) error {
	mapping, ok, err := s.tracker.GetUserMapping(userId)
	if err != nil {
		return err
	}
	if !ok || mapping == nil || mapping.UnitId == nil || mapping.GetUnitId() != unitId {
		return nil
	}

	return s.tracker.DeleteUserMapping(ctx, userId)
}
