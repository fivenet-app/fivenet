package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"go.uber.org/zap"
)

const (
	cleanupUnitsDispatchesUnassignedAttr = "dispatches_unassigned"
	cleanupUnitsEmptyUnitsRemovedAttr    = "empty_units_removed"
	cleanupUnitsStatusesUpdatedAttr      = "unit_statuses_updated"
	cleanupUnitsOffDutyRemovedAttr       = "off_duty_users_removed"
	cleanupUnitsMappingsRepairedAttr     = "mappings_repaired"
)

func (s *Housekeeper) runCleanupUnits(ctx context.Context, data *cron.CronjobData) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration("cleanup_units", time.Since(startedAt).Seconds())
	}()

	ctx, span := s.tracer.Start(ctx, "centrum.units-cleanup")
	defer span.End()

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal cleanup units cron data", zap.Error(err))
	}

	unitStatusesUpdated, err := s.cleanupUnitStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to clean up unit status: %w", err)
	}
	s.metrics.SetHousekeeperWork("cleanup_units", "unit_statuses_updated", unitStatusesUpdated)

	dest.SetAttribute(cleanupUnitsStatusesUpdatedAttr, strconv.Itoa(unitStatusesUpdated))

	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal updated cleanup units cron data. %w", err)
	}

	return nil
}

func (s *Housekeeper) runAuditEmptyUnitDispatches(
	ctx context.Context,
	data *cron.CronjobData,
) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration(
			"audit_empty_unit_dispatches",
			time.Since(startedAt).Seconds(),
		)
	}()

	dest := &cron.GenericCronData{Attributes: map[string]string{}}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal empty unit dispatch audit cron data", zap.Error(err))
	}

	dispatchesUnassigned, emptyUnitsRemoved, err := s.removeDispatchesFromEmptyUnits(ctx)
	if err != nil {
		return fmt.Errorf("failed to audit empty units in dispatches. %w", err)
	}
	s.metrics.SetHousekeeperWork(
		"audit_empty_unit_dispatches",
		"dispatches_unassigned",
		dispatchesUnassigned,
	)
	s.metrics.SetHousekeeperWork(
		"audit_empty_unit_dispatches",
		"empty_units_removed",
		emptyUnitsRemoved,
	)

	dest.SetAttribute(cleanupUnitsDispatchesUnassignedAttr, strconv.Itoa(dispatchesUnassigned))
	dest.SetAttribute(cleanupUnitsEmptyUnitsRemovedAttr, strconv.Itoa(emptyUnitsRemoved))
	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal empty unit dispatch audit cron data. %w", err)
	}

	return nil
}

func (s *Housekeeper) runAuditUnitMembership(ctx context.Context, data *cron.CronjobData) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration(
			"audit_unit_membership",
			time.Since(startedAt).Seconds(),
		)
	}()

	ctx, span := s.tracer.Start(ctx, "centrum.unit-membership-audit")
	defer span.End()

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal unit membership audit cron data", zap.Error(err))
	}

	offDutyUsersRemoved, mappingsRepaired, err := s.checkUnitUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to audit unit membership: %w", err)
	}
	s.metrics.SetHousekeeperWork(
		"audit_unit_membership",
		"off_duty_users_removed",
		offDutyUsersRemoved,
	)
	s.metrics.SetHousekeeperWork("audit_unit_membership", "mappings_repaired", mappingsRepaired)

	dest.SetAttribute(cleanupUnitsOffDutyRemovedAttr, strconv.Itoa(offDutyUsersRemoved))
	dest.SetAttribute(cleanupUnitsMappingsRepairedAttr, strconv.Itoa(mappingsRepaired))

	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal unit membership audit cron data. %w", err)
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty.
func (s *Housekeeper) removeDispatchesFromEmptyUnits(ctx context.Context) (int, int, error) {
	dispatchesUnassigned := 0
	emptyUnitsRemoved := 0
	for _, settings := range s.settings.List(ctx) {
		job := settings.GetJob()

		dsps := s.dispatches.Filter(ctx, []string{job}, nil, []centrumdispatches.StatusDispatch{
			centrumdispatches.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			centrumdispatches.StatusDispatch_STATUS_DISPATCH_DELETED,
		})

		for _, dsp := range dsps {
			// Make sure unassigned dispatch has the unassigned status
			if len(dsp.GetUnits()) == 0 && dsp.GetStatus() != nil &&
				!centrumutils.IsStatusDispatchUnassigned(dsp.GetStatus().GetStatus()) {
				s.logger.Debug(
					"updating dispatch status to unassigned because it has no assignments",
					zap.String("job", job),
					zap.Int64("dispatch_id", dsp.GetId()),
				)
				if _, err := s.dispatches.UpdateStatus(
					ctx,
					dsp.GetId(),
					&centrumdispatches.DispatchStatus{
						CreatedAt:  timestamp.Now(),
						DispatchId: dsp.GetId(),
						Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNASSIGNED,
						CreatorJob: &job,
					},
				); err != nil {
					return dispatchesUnassigned, emptyUnitsRemoved, err
				}
				dispatchesUnassigned++

				continue
			}

			for i := range slices.Backward(dsp.GetUnits()) {
				if i > (len(dsp.GetUnits()) - 1) {
					break
				}

				unitId := dsp.GetUnits()[i].GetUnitId()
				// If unit isn't empty, continue with the loop
				if unitId <= 0 {
					continue
				}

				unit, err := s.units.Get(ctx, unitId)
				if err != nil {
					continue
				}

				if len(unit.GetUsers()) > 0 {
					continue
				}

				s.logger.Debug(
					"removing empty unit from dispatch",
					zap.String(
						"job",
						job,
					),
					zap.Int64("unit_id", unitId),
					zap.Int64("dispatch_id", dsp.GetId()),
				)

				if err := s.dispatches.UpdateAssignments(
					ctx,
					new(job),
					nil,
					dsp.GetId(),
					nil,
					[]int64{unitId},
					time.Time{},
				); err != nil {
					s.logger.Error(
						"failed to remove empty unit from dispatch",
						zap.String(
							"job",
							job,
						),
						zap.Int64("unit_id", unitId),
						zap.Int64("dispatch_id", dsp.GetId()),
						zap.Error(err),
					)
					continue
				}
				emptyUnitsRemoved++
			}
		}
	}

	return dispatchesUnassigned, emptyUnitsRemoved, nil
}

// Iterate over units to ensure that, e.g., an empty unit status is set to `unavailable`.
func (s *Housekeeper) cleanupUnitStatus(ctx context.Context) (int, error) {
	updated := 0
	for _, settings := range s.settings.List(ctx) {
		job := settings.GetJob()

		units := s.units.List(ctx, []string{job})
		for _, unit := range units {
			// Either unit has users but is static and in a wrong status
			if len(unit.GetUsers()) > 0 {
				if unit.GetAttributes() == nil ||
					!unit.GetAttributes().Has(centrumunits.UnitAttribute_UNIT_ATTRIBUTE_STATIC) {
					continue
				}

				if unit.GetStatus() != nil &&
					(unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_BUSY ||
						unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_ON_BREAK ||
						unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE) {
					continue
				}
			} else if unit.GetStatus() != nil &&
				// Or the unit is not already set to be unavailable (because it is empty)
				unit.GetStatus().GetStatus() == centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE {
				continue
			}

			var userId *int32
			if unit.GetStatus() != nil && unit.Status.UserId != nil {
				userId = unit.GetStatus().UserId
			}

			s.logger.Debug(
				"setting unit status to unavailable it is empty or static attribute (wrong status)",
				zap.String(
					"job",
					job,
				),
				zap.Int64("unit_id", unit.GetId()),
				zap.Int32p("user_id", userId),
			)
			if _, _, err := s.units.UpdateStatus(ctx, unit.GetId(), &centrumunits.UnitStatus{
				CreatedAt:  timestamp.Now(),
				UnitId:     unit.GetId(),
				Status:     centrumunits.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId:     userId,
				CreatorJob: &job,
			}); err != nil {
				s.logger.Error(
					"failed to update empty unit status to unavailable",
					zap.String(
						"job",
						unit.GetJob(),
					),
					zap.Int64("unit_id", unit.GetId()),
					zap.Error(err),
				)
				continue
			}
			updated++
		}
	}

	return updated, nil
}

// Make sure that all users in units are still on duty.
func (s *Housekeeper) checkUnitUsers(ctx context.Context) (int, int, error) {
	foundUserIds := map[int32]struct{}{}
	offDutyRemoved := 0

	for _, settings := range s.settings.List(ctx) {
		job := settings.GetJob()

		units := s.units.List(ctx, []string{job})
		for _, u := range units {
			unit, err := s.units.Get(ctx, u.GetId())
			if err != nil {
				continue
			}

			if len(unit.GetUsers()) == 0 {
				continue
			}

			foundUids, removed, err := s.checkAndUpdateUnitUsers(ctx, unit)
			if err != nil {
				s.logger.Error("failed to check users in unit", zap.Error(err))
			}
			for _, userId := range foundUids {
				foundUserIds[userId] = struct{}{}
			}
			offDutyRemoved += removed
		}
	}

	userUnitIds, err := s.tracker.ListUserMappings(ctx)
	if err != nil {
		return 0, 0, err
	}

	var errs error
	mappingsRepaired := 0
	for _, userUnit := range userUnitIds {
		if userUnit == nil || userUnit.UnitId == nil || userUnit.GetUnitId() <= 0 {
			continue
		}

		// Check if user id is part of an unit
		if _, ok := foundUserIds[userUnit.GetUserId()]; ok {
			continue
		}

		s.logger.Warn(
			"found user with unit mapping that isn't in any unit anymore",
			zap.Int32("user_id", userUnit.GetUserId()),
			zap.Int("users_in_units", len(foundUserIds)),
			zap.Any("mapping", userUnit),
		)

		if err := s.units.SyncUserUnitMapping(ctx, userUnit.GetUserId()); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		mappingsRepaired++
	}

	return offDutyRemoved, mappingsRepaired, errs
}

func (s *Housekeeper) checkAndUpdateUnitUsers(
	ctx context.Context,
	unit *centrumunits.Unit,
) ([]int32, int, error) {
	if len(unit.GetUsers()) == 0 {
		return nil, 0, nil
	}

	toRemove := []int32{}
	foundUserIds := []int32{}
	for i := range slices.Backward(unit.GetUsers()) {
		if i > (len(unit.GetUsers()) - 1) {
			break
		}

		userId := unit.GetUsers()[i].GetUserId()
		if userId == 0 {
			s.logger.Warn(
				"zero user id found during unit user checkup",
				zap.Int64("unit_id", unit.GetId()),
			)
			continue
		}

		inJob, err := s.unitAssignments.UserInJob(ctx, s.db, unit.GetJob(), userId)
		if err != nil {
			return foundUserIds, 0, fmt.Errorf("failed to check user job membership. %w", err)
		}

		marker, markerFound := s.tracker.GetUserMarkerById(userId)
		// Tracker mappings are a projection and may be temporarily absent. Only
		// duty/job facts may remove durable unit membership.
		if markerFound && marker != nil && !marker.GetHidden() &&
			s.tracker.IsUserOnDuty(userId) && marker.GetJob() == unit.GetJob() && inJob {
			foundUserIds = append(foundUserIds, userId)
			continue
		}

		toRemove = append(toRemove, userId)
	}

	if len(toRemove) == 0 {
		return foundUserIds, 0, nil
	}

	s.logger.Debug(
		"removing off-duty users from unit",
		zap.String(
			"job",
			unit.GetJob(),
		),
		zap.Int64("unit_id", unit.GetId()),
		zap.Int32s("to_remove", toRemove),
	)

	if err := s.unitAssignments.UpdateUnitAssignments(
		ctx,
		unit.GetJob(),
		nil,
		unit.GetId(),
		nil,
		toRemove,
	); err != nil {
		s.logger.Error(
			"failed to remove off-duty users from unit",
			zap.String(
				"job",
				unit.GetJob(),
			),
			zap.Int64("unit_id", unit.GetId()),
			zap.Int32s("user_ids", toRemove),
			zap.Error(err),
		)
		return foundUserIds, 0, fmt.Errorf("failed to update unit assignments. %w", err)
	}

	return foundUserIds, len(toRemove), nil
}
