package housekeeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	loadNewDispatchesLoadedAttr = "new_dispatches_loaded"
	loadNewDispatchesTotalAttr  = "loaded_dispatches"

	dispatchAssignmentExpiredAttr  = "expired_assignments"
	dispatchAssignmentAffectedAttr = "dispatches_affected"
	dispatchAssignmentJobsAttr     = "jobs_affected"
	dispatchAssignmentUnitsAttr    = "units_affected"
	dispatchAssignmentBacklogAttr  = "backlog_remaining"

	cancelOldDispatchesCancelledAttr = "dispatches_cancelled"
	cancelOldDispatchesTooOldAttr    = "too_old_flagged"
	cancelOldDispatchesBacklogAttr   = "backlog_remaining"

	deleteOldDispatchesDeletedAttr = "dispatches_deleted"

	deleteOldDispatchesKVKeysScannedAttr = "kv_keys_scanned"
	deleteOldDispatchesKVDeletedAttr     = "kv_dispatches_deleted"
	deleteOldDispatchesKVInvalidAttr     = "kv_invalid_keys"
	deleteOldDispatchesKVScanLimitAttr   = "kv_scan_limit_reached"

	// Keep recovery work bounded even when no entries qualify for deletion.
	maxKVDispatchKeysScannedPerRun = 500
	maxKVDispatchDeletesPerRun     = 100

	maxExpiredDispatchAssignmentsPerRun = 100
)

func (s *Housekeeper) loadNewDispatches(ctx context.Context, data *cron.CronjobData) error {
	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")

	s.logger.Debug("loading new dispatches from DB")

	dest := &cron.GenericCronData{}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Error("failed to unmarshal centrum housekeeper cron data", zap.Error(err))
	}

	// FIXME: Legacy FiveM plugins create dispatches without postal codes and rely
	// on this poller to load them. Remove this once they publish projections or
	// use the Centrum API instead.
	// Load dispatches with null postal field (they are considered "new").
	dspCount, err := s.dispatches.LoadFromDB(ctx, tDispatch.Postal.IS_NULL())
	if err != nil {
		return fmt.Errorf("failed loading new dispatches from DB: %w", err)
	}

	count := int64(dspCount)
	dest.SetAttribute(loadNewDispatchesLoadedAttr, strconv.FormatInt(count, 10))

	if val := dest.GetAttribute("loaded_dispatches"); val != "" {
		if cc, err := strconv.ParseInt(val, 10, 64); err == nil {
			count += cc
		}
	}
	dest.SetAttribute(loadNewDispatchesTotalAttr, strconv.FormatInt(count, 10))

	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal updated centrum housekeeper cron data: %w", err)
	}

	return nil
}

func (s *Housekeeper) runHandleDispatchAssignmentExpiration(
	ctx context.Context,
	data *cron.CronjobData,
) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-assignment-expiration")
	defer span.End()

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn(
			"failed to unmarshal dispatch assignment expiration cron data",
			zap.Error(err),
		)
	}

	expiredAssignments, dispatchesAffected, jobsAffected, unitsAffected, backlogRemaining, err := s.handleDispatchAssignmentExpiration(
		ctx,
	)
	if err != nil {
		s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
		return err
	}

	dest.SetAttribute(dispatchAssignmentExpiredAttr, strconv.Itoa(expiredAssignments))
	dest.SetAttribute(dispatchAssignmentAffectedAttr, strconv.Itoa(dispatchesAffected))
	dest.SetAttribute(dispatchAssignmentJobsAttr, strconv.Itoa(jobsAffected))
	dest.SetAttribute(dispatchAssignmentUnitsAttr, strconv.Itoa(unitsAffected))
	s.metrics.SetHousekeeperWork(
		"dispatch_assignment_expiration",
		"backlog_remaining",
		boolToInt(backlogRemaining),
	)
	setBacklogAttribute(dest, dispatchAssignmentBacklogAttr, backlogRemaining)

	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf(
			"failed to marshal updated dispatch assignment expiration cron data. %w",
			err,
		)
	}

	return nil
}

// Handle expired dispatch unit assignments.
func (s *Housekeeper) handleDispatchAssignmentExpiration(
	ctx context.Context,
) (int, int, int, int, bool, error) {
	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts
	tUnits := table.FivenetCentrumUnits

	stmt := tDispatchUnit.
		SELECT(
			tDispatchUnit.DispatchID.AS("dispatch_id"),
			tDispatchUnit.UnitID.AS("unit_id"),
			tUnits.Job.AS("job"),
		).
		FROM(
			tDispatchUnit.
				INNER_JOIN(tUnits,
					tUnits.ID.EQ(tDispatchUnit.UnitID),
				),
		).
		WHERE(mysql.AND(
			tDispatchUnit.ExpiresAt.IS_NOT_NULL(),
			tDispatchUnit.ExpiresAt.LT_EQ(mysql.CURRENT_TIMESTAMP()),
		)).
		ORDER_BY(
			tDispatchUnit.ExpiresAt.ASC(),
			tDispatchUnit.DispatchID.ASC(),
			tDispatchUnit.UnitID.ASC(),
		).
		LIMIT(maxExpiredDispatchAssignmentsPerRun + 1)

	var dest []*struct {
		DispatchID int64
		UnitID     int64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, 0, 0, 0, false, err
	}
	dest, backlogRemaining := capWorkItems(dest, maxExpiredDispatchAssignmentsPerRun)

	assignments := map[string]map[int64][]int64{}
	for _, ua := range dest {
		if _, ok := assignments[ua.Job]; !ok {
			assignments[ua.Job] = map[int64][]int64{}
		}
		if _, ok := assignments[ua.Job][ua.DispatchID]; !ok {
			assignments[ua.Job][ua.DispatchID] = []int64{}
		}

		assignments[ua.Job][ua.DispatchID] = append(assignments[ua.Job][ua.DispatchID], ua.UnitID)
	}

	dispatchesSeen := map[int64]struct{}{}
	jobsSeen := map[string]struct{}{}
	for job, dsps := range assignments {
		jobsSeen[job] = struct{}{}
		for dispatchID := range dsps {
			dispatchesSeen[dispatchID] = struct{}{}
		}
	}

	dispatchesAffected := len(dispatchesSeen)
	jobsAffected := len(jobsSeen)
	unitsAffected := len(dest)

	for job, dsps := range assignments {
		s.logger.Debug(
			"handling dispatch assignment expiration",
			zap.String("job", job),
			zap.Int("expired_assignments", len(dsps)),
		)
		for dispatchId, units := range dsps {
			if err := s.dispatches.UpdateAssignments(
				ctx,
				new(job),
				nil,
				dispatchId,
				nil,
				units,
				time.Time{},
			); err != nil {
				return 0, 0, 0, 0, backlogRemaining, fmt.Errorf(
					"failed to update dispatch %d assignments. %w",
					dispatchId,
					err,
				)
			}
		}
	}

	return len(dest), dispatchesAffected, jobsAffected, unitsAffected, backlogRemaining, nil
}

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

func boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

// capWorkItems reserves one selected row to signal that another recovery run
// is needed, while keeping the work performed in this run bounded.
func capWorkItems[T any](items []T, limit int) ([]T, bool) {
	if len(items) <= limit {
		return items, false
	}

	return items[:limit], true
}

func setBacklogAttribute(data *cron.GenericCronData, key string, backlogRemaining bool) {
	data.SetAttribute(key, strconv.FormatBool(backlogRemaining))
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

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
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
	if deleted >= maxKVDispatchDeletesPerRun {
		s.metrics.SetHousekeeperWork("delete_old_dispatches_from_kv", "delete_limit_reached", 1)
	} else {
		s.metrics.SetHousekeeperWork("delete_old_dispatches_from_kv", "delete_limit_reached", 0)
	}

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
	keysScanned := 0
	deleted := 0
	invalid := 0

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

		dspId, err := centrumutils.ExtractIDString(key)
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

		dsp, err := s.dispatches.Store().Get(dspId)
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

		if (
		// Dispatches older than 3 hours will be removed from the KV store (not the database)
		dsp.GetCreatedAt() != nil && time.Since(dsp.GetCreatedAt().AsTime()) > 3*time.Hour) ||
			// Remove nil status dispatches
			dsp.GetStatus() == nil ||
			// "Completed" dispatches with their status being older than 15 minutes
			(centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) &&
				time.Since(dsp.GetStatus().GetCreatedAt().AsTime()) > 15*time.Minute) {
			s.logger.Debug("old dispatch deleted from kv", zap.Int64("dispatch_id", dsp.GetId()))

			if err := s.dispatches.Delete(ctx, dsp.GetId(), false); err != nil {
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
