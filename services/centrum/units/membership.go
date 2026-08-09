package units

import (
	"context"
	"errors"
	"fmt"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
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

	var targetUnitId *int64
	if unitId > 0 {
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

	return errors.Join(errs, s.clearStaleTrackerMappingsForUnit(ctx, unitId, userIds))
}

func (s *UnitDB) syncMissingUnitMembership(ctx context.Context, unitId int64) error {
	var errs error
	if err := s.deleteInKV(ctx, unitId); err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		errs = errors.Join(errs, err)
	}

	return errors.Join(errs, s.clearStaleTrackerMappingsForUnit(ctx, unitId, nil))
}

func (s *UnitDB) clearStaleTrackerMappingsForUnit(
	ctx context.Context,
	unitId int64,
	validUserIds map[int32]struct{},
) error {
	mappings, err := s.tracker.ListUserMappings(ctx)
	if err != nil {
		return err
	}

	var errs error
	for userId, mapping := range mappings {
		if mapping == nil || mapping.UnitId == nil || mapping.GetUnitId() != unitId {
			continue
		}

		if _, ok := validUserIds[userId]; ok {
			continue
		}

		if err := s.tracker.DeleteUserMapping(ctx, userId); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
