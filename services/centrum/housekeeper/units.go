package housekeeper

import (
	"context"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"go.uber.org/zap"
)

func (s *Housekeeper) runCleanupUnits(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.units-cleanup")
	defer span.End()

	if err := s.removeDispatchesFromEmptyUnits(ctx); err != nil {
		s.logger.Error("failed to clean empty units from dispatches", zap.Error(err))
	}

	if err := s.cleanupUnitStatus(ctx); err != nil {
		s.logger.Error("failed to clean up unit status", zap.Error(err))
	}

	if err := s.checkUnitUsers(ctx); err != nil {
		s.logger.Error("failed to check duty state of unit users", zap.Error(err))
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Housekeeper) removeDispatchesFromEmptyUnits(ctx context.Context) error {
	for _, settings := range s.settings.List(ctx) {
		job := settings.Job

		dsps := s.dispatches.Filter(ctx, []string{job}, nil, []centrum.StatusDispatch{
			centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			centrum.StatusDispatch_STATUS_DISPATCH_DELETED,
		})

		for _, dsp := range dsps {
			// Make sure unassigned dispatch has the unassigned status
			if len(dsp.Units) == 0 && dsp.Status != nil && !centrumutils.IsStatusDispatchUnassigned(dsp.Status.Status) {
				s.logger.Debug("updating dispatch status to unassigned because it has no assignments",
					zap.String("job", job), zap.Uint64("dispatch_id", dsp.Id))
				if _, err := s.dispatches.UpdateStatus(ctx, dsp.Id, &centrum.DispatchStatus{
					CreatedAt:  timestamp.Now(),
					DispatchId: dsp.Id,
					Status:     centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
					CreatorJob: &job,
				}); err != nil {
					return err
				}

				continue
			}

			for i := range slices.Backward(dsp.Units) {
				if i > (len(dsp.Units) - 1) {
					break
				}

				unitId := dsp.Units[i].UnitId
				// If unit isn't empty, continue with the loop
				if unitId <= 0 {
					continue
				}

				unit, err := s.units.Get(ctx, unitId)
				if err != nil {
					continue
				}

				if len(unit.Users) > 0 {
					continue
				}

				s.logger.Debug("removing empty unit from dispatch",
					zap.String("job", job), zap.Uint64("unit_id", unitId), zap.Uint64("dispatch_id", dsp.Id))

				if err := s.dispatches.UpdateAssignments(ctx, nil, dsp.Id, nil, []uint64{unitId}, time.Time{}); err != nil {
					s.logger.Error("failed to remove empty unit from dispatch",
						zap.String("job", job), zap.Uint64("unit_id", unitId), zap.Uint64("dispatch_id", dsp.Id), zap.Error(err))
					continue
				}
			}
		}
	}

	return nil
}

// Iterate over units to ensure that, e.g., an empty unit status is set to `unavailable`
func (s *Housekeeper) cleanupUnitStatus(ctx context.Context) error {
	for _, settings := range s.settings.List(ctx) {
		job := settings.Job

		units := s.units.List(ctx, []string{job})
		for _, unit := range units {
			// Either unit has users but is static and in a wrong status
			if len(unit.Users) > 0 {
				if unit.Attributes == nil || !unit.Attributes.Has(centrum.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
					continue
				}

				if unit.Status != nil &&
					(unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_BUSY ||
						unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_ON_BREAK ||
						unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE) {
					continue
				}
			} else {
				// Or the unit is not already set to be unavailable (because it is empty)
				if unit.Status != nil &&
					unit.Status.Status == centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE {
					continue
				}
			}

			var userId *int32
			if unit.Status != nil && unit.Status.UserId != nil {
				userId = unit.Status.UserId
			}

			s.logger.Debug("setting unit status to unavailable it is empty or static attribute (wrong status)",
				zap.String("job", job), zap.Uint64("unit_id", unit.Id), zap.Int32p("user_id", userId))
			if _, err := s.units.UpdateStatus(ctx, unit.Id, &centrum.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.Id,
				Status:     centrum.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:     userId,
				CreatorJob: &job,
			}); err != nil {
				s.logger.Error("failed to update empty unit status to unavailable",
					zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Error(err))
				continue
			}
		}
	}

	return nil
}

// Make sure that all users in units are still on duty
func (s *Housekeeper) checkUnitUsers(ctx context.Context) error {
	foundUserIds := []int32{}

	for _, settings := range s.settings.List(ctx) {
		job := settings.Job

		units := s.units.List(ctx, []string{job})
		for _, u := range units {
			unit, err := s.units.Get(ctx, u.Id)
			if err != nil {
				continue
			}

			if len(unit.Users) == 0 {
				continue
			}

			foundUids, _, err := s.checkAndUpdateUnitUsers(ctx, unit)
			if err != nil {
				s.logger.Error("failed to check users in unit", zap.Error(err))
			}
			foundUserIds = append(foundUserIds, foundUids...)
		}
	}

	userUnitIds, err := s.tracker.ListUserMappings(ctx)
	if err != nil {
		return err
	}

	var errs error
	for _, userUnit := range userUnitIds {
		// Check if user id is part of an unit
		if slices.Contains(foundUserIds, userUnit.UserId) {
			continue
		}

		s.logger.Warn("found user id with unit mapping that isn't in any unit anymore", zap.Int32("user_id", userUnit.UserId), zap.Int32s("users_in_units", foundUserIds), zap.Any("mapping", userUnit))

		// TODO this isn't working as intended at the moment..
		/*
			// Unset unit id for user when user is not in any unit
			if err := s.tracker.UnsetUnitIDForUser(ctx, userId); err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		*/
	}

	return errs
}

func (s *Housekeeper) checkAndUpdateUnitUsers(ctx context.Context, unit *centrum.Unit) ([]int32, bool, error) {
	if len(unit.Users) == 0 {
		return nil, false, nil
	}

	toRemove := []int32{}
	foundUserIds := []int32{}
	for i := range slices.Backward(unit.Users) {
		if i > (len(unit.Users) - 1) {
			break
		}

		userId := unit.Users[i].UserId
		if userId == 0 {
			s.logger.Warn("zero user id found during unit user checkup", zap.Uint64("unit_id", unit.Id))
			continue
		}

		unitMapping, err := s.tracker.GetUserMapping(userId)
		// If user is in that unit and still on duty, nothing to do, otherwise remove the user from the unit
		if err == nil && unitMapping.UnitId != nil && unit.Id == *unitMapping.UnitId && s.tracker.IsUserOnDuty(userId) {
			foundUserIds = append(foundUserIds, userId)
			continue
		}

		toRemove = append(toRemove, userId)
	}

	if len(toRemove) == 0 {
		return nil, false, nil
	}

	s.logger.Debug("removing off-duty users from unit",
		zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Int32s("to_remove", toRemove))

	if err := s.units.UpdateUnitAssignments(ctx, unit.Job, nil, unit.Id, nil, toRemove); err != nil {
		s.logger.Error("failed to remove off-duty users from unit",
			zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Int32s("user_ids", toRemove), zap.Error(err))
	}

	return foundUserIds, true, nil
}
