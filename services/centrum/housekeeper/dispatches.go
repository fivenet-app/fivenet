package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	cancelOldDispatchesCancelledAttr = "dispatches_cancelled"
	cancelOldDispatchesTooOldAttr    = "too_old_flagged"
	cancelOldDispatchesBacklogAttr   = "backlog_remaining"

	deleteOldDispatchesDeletedAttr = "dispatches_deleted"

	// Keep recovery work bounded even when no entries qualify for deletion.
	maxKVDispatchKeysScannedPerRun = 500
	maxKVDispatchDeletesPerRun     = 100

	dispatchAssignmentExpiredAttr  = "expired_assignments"
	dispatchAssignmentDeletedAttr  = "assignments_deleted"
	dispatchAssignmentArchivedAttr = "archived_assignments_deleted"
	dispatchAssignmentAffectedAttr = "dispatches_affected"
	dispatchAssignmentJobsAttr     = "jobs_affected"
	dispatchAssignmentUnitsAttr    = "units_affected"
	dispatchAssignmentBacklogAttr  = "backlog_remaining"

	maxExpiredDispatchAssignmentsPerRun = 100

	deleteOldDispatchesKVKeysScannedAttr = "kv_keys_scanned"
	deleteOldDispatchesKVDeletedAttr     = "kv_dispatches_deleted"
	deleteOldDispatchesKVInvalidAttr     = "kv_invalid_keys"
	deleteOldDispatchesKVScanLimitAttr   = "kv_scan_limit_reached"
)

func (s *Housekeeper) runCancelOldDispatches(ctx context.Context, data *cron.CronjobData) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration(
			"cancel_old_dispatches",
			time.Since(startedAt).Seconds(),
		)
	}()

	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-cancel")
	defer span.End()

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal cancel old dispatches cron data", zap.Error(err))
	}

	cancelled, flaggedTooOld, backlogRemaining, err := s.cancelOldDispatches(ctx)
	if err != nil {
		s.logger.Error("failed to archive dispatches", zap.Error(err))
		return fmt.Errorf("failed to archive dispatches. %w", err)
	}
	s.metrics.SetHousekeeperWork("cancel_old_dispatches", "cancelled", cancelled)
	s.metrics.SetHousekeeperWork("cancel_old_dispatches", "too_old_flagged", flaggedTooOld)
	s.metrics.SetHousekeeperWork(
		"cancel_old_dispatches",
		"backlog_remaining",
		boolToInt(backlogRemaining),
	)

	dest.SetAttribute(cancelOldDispatchesCancelledAttr, strconv.Itoa(cancelled))
	dest.SetAttribute(cancelOldDispatchesTooOldAttr, strconv.Itoa(flaggedTooOld))
	setBacklogAttribute(dest, cancelOldDispatchesBacklogAttr, backlogRemaining)

	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal updated cancel old dispatches cron data. %w", err)
	}

	return nil
}

// Cancel dispatches that haven't been worked on for some time.
func (s *Housekeeper) cancelOldDispatches(ctx context.Context) (int, int, bool, error) {
	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus := table.FivenetCentrumDispatchesStatus

	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.DispatchID.AS("dispatch_id"),
			tDispatch.Jobs.AS("jobs"),
			tDispatchStatus.Status.AS("status"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.DispatchID),
				),
		).
		// Dispatches that are older than time X and are not in a completed/cancelled/archived state, or have no status at all
		WHERE(mysql.AND(
			tDispatchStatus.ID.EQ(
				mysql.RawInt(
					"SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`",
				),
			),
			tDispatchStatus.Status.NOT_IN(
				mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
				mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
				mysql.Int32(int32(centrumdispatches.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
			),
			tDispatch.CreatedAt.LT_EQ(
				mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(60, mysql.MINUTE)),
			),
		)).
		ORDER_BY(
			tDispatchStatus.DispatchID.ASC(),
		).
		LIMIT(MaxCancelledDispatchesPerRun + 1)

	var dest []*struct {
		DispatchID int64
		Jobs       []string
		Status     centrumdispatches.StatusDispatch
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, 0, false, err
	}
	dest, backlogRemaining := capWorkItems(dest, MaxCancelledDispatchesPerRun)

	s.logger.Debug("canceling expired dispatches", zap.Int("dispatch_count", len(dest)))
	cancelled := 0
	flaggedTooOld := 0
	for _, ds := range dest {
		// Ignore already cancelled dispatches
		if ds.Status == centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED {
			continue
		}

		// Add "too old" attribute when we are able to retrieve the dispatch
		if dsp, err := s.dispatches.Get(ctx, ds.DispatchID); err == nil && dsp != nil {
			if err := s.dispatches.AddAttributeToDispatch(
				ctx,
				dsp,
				centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_TOO_OLD,
			); err != nil {
				s.logger.Error(
					"failed to add too old attribute to cancelled dispatch",
					zap.Int64("dispatch_id", ds.DispatchID),
					zap.Error(err),
				)
			} else {
				flaggedTooOld++
			}
		}

		if _, err := s.dispatches.UpdateStatus(
			ctx,
			ds.DispatchID,
			&centrumdispatches.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: ds.DispatchID,
				Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			},
		); err != nil {
			s.logger.Error(
				"failed to cancel dispatch",
				zap.Int64("dispatch_id", ds.DispatchID),
				zap.Error(err),
			)
			continue
		}
		cancelled++

		// Remove dispatch from state and publish event so clients remove it
		if err := s.dispatches.Delete(ctx, ds.DispatchID, false); err != nil {
			s.logger.Error(
				"failed to delete cancelled dispatch",
				zap.Int64("dispatch_id", ds.DispatchID),
				zap.Error(err),
			)
			continue
		}
	}

	return cancelled, flaggedTooOld, backlogRemaining, nil
}

func (s *Housekeeper) runDeleteOldDispatches(ctx context.Context, data *cron.CronjobData) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration(
			"delete_old_dispatches",
			time.Since(startedAt).Seconds(),
		)
	}()

	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-old-delete")
	defer span.End()

	if data == nil {
		data = &cron.CronjobData{}
	}

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal delete old dispatches cron data", zap.Error(err))
	}

	deleted, err := s.deleteOldDispatches(ctx)
	if err != nil {
		s.logger.Error("failed to remove old dispatches", zap.Error(err))
		return err
	}
	s.metrics.SetHousekeeperWork("delete_old_dispatches", "deleted", deleted)

	dest.SetAttribute(deleteOldDispatchesDeletedAttr, strconv.Itoa(deleted))
	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal updated delete old dispatches cron data. %w", err)
	}

	return nil
}

// deleteOldDispatches deletes dispatches that are older than a certain number of days.
// This can probably be moved into the general housekeeper service.
func (s *Housekeeper) deleteOldDispatches(ctx context.Context) (int, error) {
	tDispatch := table.FivenetCentrumDispatches

	stmt := tDispatch.
		SELECT(
			tDispatch.ID.AS("dispatch_id"),
		).
		FROM(
			tDispatch,
		).
		WHERE(mysql.AND(
			tDispatch.CreatedAt.LT_EQ(
				mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(DeleteDispatchDays, mysql.DAY)),
			),
		)).
		ORDER_BY(
			tDispatch.CreatedAt.ASC(),
			tDispatch.ID.ASC(),
		).
		// Get 75 at a time
		LIMIT(75)

	var dest []*struct {
		DispatchID int64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, err
	}

	errs := multierr.Combine()
	deleted := 0
	for _, ds := range dest {
		if err := s.dispatches.Delete(ctx, ds.DispatchID, true); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		deleted++
	}

	return deleted, errs
}

func (s *Housekeeper) runHandleDispatchAssignmentExpiration(
	ctx context.Context,
	data *cron.CronjobData,
) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-assignment-expiration")
	defer span.End()

	dest := &cron.GenericCronData{Attributes: map[string]string{}}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn(
			"failed to unmarshal dispatch assignment expiration cron data",
			zap.Error(err),
		)
	}

	expired, deleted, archivedDeleted, dispatchesAffected, jobsAffected, unitsAffected, backlog, err := s.handleDispatchAssignmentExpiration(
		ctx,
	)
	s.metrics.SetHousekeeperWork("dispatch_assignment_expiration", "expired", expired)
	s.metrics.SetHousekeeperWork("dispatch_assignment_expiration", "deleted", deleted)
	s.metrics.SetHousekeeperWork("dispatch_assignment_expiration", "archived_deleted", archivedDeleted)
	dest.SetAttribute(dispatchAssignmentExpiredAttr, strconv.Itoa(expired))
	dest.SetAttribute(dispatchAssignmentDeletedAttr, strconv.Itoa(deleted))
	dest.SetAttribute(dispatchAssignmentArchivedAttr, strconv.Itoa(archivedDeleted))
	dest.SetAttribute(dispatchAssignmentAffectedAttr, strconv.Itoa(dispatchesAffected))
	dest.SetAttribute(dispatchAssignmentJobsAttr, strconv.Itoa(jobsAffected))
	dest.SetAttribute(dispatchAssignmentUnitsAttr, strconv.Itoa(unitsAffected))
	s.metrics.SetHousekeeperWork(
		"dispatch_assignment_expiration",
		"backlog_remaining",
		boolToInt(backlog),
	)
	setBacklogAttribute(dest, dispatchAssignmentBacklogAttr, backlog)
	if err != nil {
		s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
		return err
	}
	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf(
			"failed to marshal updated dispatch assignment expiration cron data. %w",
			err,
		)
	}
	return nil
}

func (s *Housekeeper) handleDispatchAssignmentExpiration(
	ctx context.Context,
) (int, int, int, int, int, int, bool, error) {
	assignmentsTable := table.FivenetCentrumDispatchesAsgmts
	unitsTable := table.FivenetCentrumUnits
	stmt := assignmentsTable.SELECT(
		assignmentsTable.DispatchID.AS(
			"dispatch_id",
		),
		assignmentsTable.UnitID.AS("unit_id"),
		unitsTable.Job.AS("job"),
	).FROM(assignmentsTable.INNER_JOIN(unitsTable, unitsTable.ID.EQ(assignmentsTable.UnitID))).WHERE(mysql.AND(
		assignmentsTable.ExpiresAt.IS_NOT_NULL(),
		assignmentsTable.ExpiresAt.LT_EQ(
			mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(2, mysql.SECOND)),
		),
	)).ORDER_BY(assignmentsTable.ExpiresAt.ASC(), assignmentsTable.DispatchID.ASC(), assignmentsTable.UnitID.ASC()).LIMIT(maxExpiredDispatchAssignmentsPerRun + 1)

	var rows []*struct {
		DispatchID int64
		UnitID     int64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil {
		return 0, 0, 0, 0, 0, 0, false, err
	}
	rows, backlog := capWorkItems(rows, maxExpiredDispatchAssignmentsPerRun)
	grouped := map[string]map[int64][]int64{}
	for _, row := range rows {
		if grouped[row.Job] == nil {
			grouped[row.Job] = map[int64][]int64{}
		}
		grouped[row.Job][row.DispatchID] = append(grouped[row.Job][row.DispatchID], row.UnitID)
	}
	dispatchIDs, jobs := map[int64]struct{}{}, map[string]struct{}{}
	errs := multierr.Combine()
	deleted := 0
	archivedDeleted := 0
	for job, dispatches := range grouped {
		jobs[job] = struct{}{}
		for id := range dispatches {
			dispatchIDs[id] = struct{}{}
		}
	}
	for job, dispatches := range grouped {
		for dispatchID, units := range dispatches {
			if err := s.assignmentExpirationWriter.UpdateAssignments(
				ctx,
				new(job),
				nil,
				dispatchID,
				nil,
				units,
				time.Time{},
			); err != nil {
				if errors.Is(err, jetstream.ErrKeyNotFound) {
					cleanupDeleted, cleanupErr := s.assignmentExpirationWriter.DeleteExpiredAssignments(
						ctx,
						dispatchID,
						units,
					)
					if cleanupErr != nil {
						errs = multierr.Append(errs, fmt.Errorf(
							"failed to delete expired assignments for archived dispatch %d. %w",
							dispatchID,
							cleanupErr,
						))
						continue
					}
					deleted += int(cleanupDeleted)
					archivedDeleted += int(cleanupDeleted)
					s.logger.Debug(
						"deleted expired assignments for archived dispatch",
						zap.Int64("dispatch_id", dispatchID),
						zap.Int("requested_assignments", len(units)),
						zap.Int64("deleted_assignments", cleanupDeleted),
					)
					continue
				}
				errs = multierr.Append(errs, fmt.Errorf(
					"failed to update dispatch %d assignments. %w",
					dispatchID,
					err,
				))
				continue
			}
			deleted += len(units)
		}
	}
	return len(rows), deleted, archivedDeleted, len(dispatchIDs), len(jobs), len(rows), backlog, errs
}

func (s *Housekeeper) runDeleteOldDispatchesFromKV(
	ctx context.Context,
	data *cron.CronjobData,
) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration(
			"delete_old_dispatches_from_kv",
			time.Since(startedAt).Seconds(),
		)
	}()

	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-old-delete-kv")
	defer span.End()

	dest := &cron.GenericCronData{Attributes: map[string]string{}}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal delete old dispatches from kv cron data", zap.Error(err))
	}
	keysScanned, deleted, invalid, scanLimitReached, err := s.deleteOldDispatchesFromKV(ctx)
	if err != nil {
		s.logger.Error("failed to remove old dispatches from kv", zap.Error(err))
		return err
	}
	s.metrics.SetHousekeeperWork("delete_old_dispatches_from_kv", "keys_scanned", keysScanned)
	s.metrics.SetHousekeeperWork("delete_old_dispatches_from_kv", "deleted", deleted)
	s.metrics.SetHousekeeperWork("delete_old_dispatches_from_kv", "invalid", invalid)
	s.metrics.SetHousekeeperWork(
		"delete_old_dispatches_from_kv",
		"scan_limit_reached",
		boolToInt(scanLimitReached),
	)
	s.metrics.SetHousekeeperWork(
		"delete_old_dispatches_from_kv",
		"delete_limit_reached",
		boolToInt(deleted >= maxKVDispatchDeletesPerRun),
	)
	dest.SetAttribute(deleteOldDispatchesKVKeysScannedAttr, strconv.Itoa(keysScanned))
	dest.SetAttribute(deleteOldDispatchesKVDeletedAttr, strconv.Itoa(deleted))
	dest.SetAttribute(deleteOldDispatchesKVInvalidAttr, strconv.Itoa(invalid))
	setBacklogAttribute(dest, deleteOldDispatchesKVScanLimitAttr, scanLimitReached)
	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf(
			"failed to marshal updated delete old dispatches from kv cron data. %w",
			err,
		)
	}
	return nil
}

func (s *Housekeeper) deleteOldDispatchesFromKV(ctx context.Context) (int, int, int, bool, error) {
	errs := multierr.Combine()
	keysScanned, deleted, invalid := 0, 0, 0

	listCtx, cancelList := context.WithCancel(ctx)
	defer cancelList()

	keyIter, err := s.dispatches.Store().KV().ListKeysFiltered(listCtx, "id.*")
	if err != nil {
		s.logger.Error("failed to list dispatches from KV", zap.Error(err))
		return 0, 0, 0, false, err
	}

	scanLimitReached := false
	keysCh := keyIter.Keys()
	for key := range keysCh {
		if keysScanned >= maxKVDispatchKeysScannedPerRun {
			scanLimitReached = true
			cancelList()
			for range keysCh {
			}
			break
		}
		keysScanned++
		if key == "" {
			continue
		}

		dispatchID, err := centrumutils.ExtractIDString(key)
		if err != nil {
			s.logger.Error(
				"failed to extract dispatch ID from key",
				zap.String("key", key),
				zap.Error(err),
			)
			errs = multierr.Append(
				errs,
				fmt.Errorf("failed to extract dispatch ID from key %q. %w", key, err),
			)
			invalid++
			continue
		}

		dispatch, err := s.dispatches.Store().Get(dispatchID)
		if err != nil {
			s.logger.Error("failed to get dispatch from KV", zap.String("key", key), zap.Error(err))
			if err := s.dispatches.Store().Delete(ctx, key); err != nil {
				s.logger.Error(
					"failed to delete unavailable dispatch from KV",
					zap.String("key", key),
					zap.Error(err),
				)
			}
			deleted++
			if deleted >= maxKVDispatchDeletesPerRun {
				cancelList()
				for range keysCh {
				}
				break
			}
			continue
		}

		old := dispatch.GetCreatedAt() != nil &&
			time.Since(dispatch.GetCreatedAt().AsTime()) > 3*time.Hour
		completed := dispatch.GetStatus() != nil &&
			centrumutils.IsStatusDispatchComplete(dispatch.GetStatus().GetStatus()) &&
			time.Since(dispatch.GetStatus().GetCreatedAt().AsTime()) > 15*time.Minute
		if old || dispatch.GetStatus() == nil || completed {
			s.logger.Debug(
				"old dispatch deleted from kv",
				zap.Int64("dispatch_id", dispatch.GetId()),
			)
			if err := s.dispatches.Delete(ctx, dispatch.GetId(), false); err != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf("failed to delete dispatch from KV. %w", err),
				)
				continue
			}
			deleted++
			if deleted >= maxKVDispatchDeletesPerRun {
				cancelList()
				for range keysCh {
				}
				break
			}
		}
	}

	return keysScanned, deleted, invalid, scanLimitReached, errs
}
