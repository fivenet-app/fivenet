package housekeeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func (s *Housekeeper) loadNewDispatches(ctx context.Context, data *cron.CronjobData) error {
	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")
	s.logger.Debug("loading new dispatches from DB")

	dest := &cron.GenericCronData{}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Error("failed to unmarshal centrum housekeeper cron data", zap.Error(err))
	}

	// Load dispatches with null postal field (they are considered "new")
	dspCount, err := s.dispatches.LoadFromDB(ctx, tDispatch.Postal.IS_NULL())
	if err != nil {
		s.logger.Error("failed loading new dispatches from DB", zap.Error(err))
	}

	count := int64(dspCount)

	if val := dest.GetAttribute("loaded_dispatches"); val != "" {
		if cc, err := strconv.ParseInt(val, 10, 64); err == nil {
			count += cc
		}
	}
	dest.SetAttribute("loaded_dispatches", strconv.FormatInt(int64(count), 10))

	if err := data.MarshalFrom(dest); err != nil {
		s.logger.Error("failed to marshal updated centrum housekeeper cron data", zap.Error(err))
	}

	return nil
}

func (s *Housekeeper) runHandleDispatchAssignmentExpiration(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-assignment-expiration")
	defer span.End()

	if err := s.handleDispatchAssignmentExpiration(ctx); err != nil {
		s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
	}

	return nil
}

// Handle expired dispatch unit assignments
func (s *Housekeeper) handleDispatchAssignmentExpiration(ctx context.Context) error {
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
		WHERE(jet.AND(
			tDispatchUnit.ExpiresAt.IS_NOT_NULL(),
			tDispatchUnit.ExpiresAt.LT_EQ(jet.CURRENT_TIMESTAMP()),
		))

	var dest []*struct {
		DispatchID uint64
		UnitID     uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	assignments := map[string]map[uint64][]uint64{}
	for _, ua := range dest {
		if _, ok := assignments[ua.Job]; !ok {
			assignments[ua.Job] = map[uint64][]uint64{}
		}
		if _, ok := assignments[ua.Job][ua.DispatchID]; !ok {
			assignments[ua.Job][ua.DispatchID] = []uint64{}
		}

		assignments[ua.Job][ua.DispatchID] = append(assignments[ua.Job][ua.DispatchID], ua.UnitID)
	}

	for job, dsps := range assignments {
		s.logger.Debug("handling dispatch assignment expiration", zap.String("job", job), zap.Int("expired_assignments", len(dsps)))
		for dispatchId, units := range dsps {
			if err := s.dispatches.UpdateAssignments(ctx, nil, dispatchId, nil, units, time.Time{}); err != nil {
				return fmt.Errorf("failed to update dispatch %d assignments. %w", dispatchId, err)
			}
		}
	}

	return nil
}

func (s *Housekeeper) runCancelOldDispatches(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-cancel")
	defer span.End()

	if err := s.cancelOldDispatches(ctx); err != nil {
		s.logger.Error("failed to archive dispatches", zap.Error(err))
	}

	return nil
}

// Cancel dispatches that haven't been worked on for some time
func (s *Housekeeper) cancelOldDispatches(ctx context.Context) error {
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
		WHERE(jet.AND(
			tDispatchStatus.ID.EQ(
				jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
			),
			tDispatchStatus.Status.NOT_IN(
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED)),
				jet.Int16(int16(centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
			),
			tDispatch.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE)),
			),
		)).
		GROUP_BY(
			tDispatchStatus.DispatchID,
		).
		ORDER_BY(
			tDispatchStatus.DispatchID.ASC(),
		).
		LIMIT(200)

	var dest []*struct {
		DispatchID uint64
		Jobs       []string
		Status     centrum.StatusDispatch
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	s.logger.Debug("canceling expired dispatches", zap.Int("dispatch_count", len(dest)))
	for _, ds := range dest {
		// Ignore already cancelled dispatches
		if ds.Status == centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED {
			continue
		}

		// Add "too old" attribute when we are able to retrieve the dispatch
		if dsp, err := s.dispatches.Get(ctx, ds.DispatchID); err == nil && dsp != nil {
			if err := s.dispatches.AddAttributeToDispatch(ctx, dsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_TOO_OLD); err != nil {
				s.logger.Error("failed to add too old attribute to cancelled dispatch", zap.Uint64("dispatch_id", ds.DispatchID), zap.Error(err))
			}
		}

		if _, err := s.dispatches.UpdateStatus(ctx, ds.DispatchID, &centrum.DispatchStatus{
			CreatedAt:  timestamp.Now(),
			DispatchId: ds.DispatchID,
			Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		}); err != nil {
			s.logger.Error("failed to cancel dispatch", zap.Uint64("dispatch_id", ds.DispatchID), zap.Error(err))
			continue
		}

		// Remove dispatch from state and publish event so clients remove it
		if err := s.dispatches.Delete(ctx, ds.DispatchID, false); err != nil {
			s.logger.Error("failed to delete cancelled dispatch", zap.Uint64("dispatch_id", ds.DispatchID), zap.Error(err))
			continue
		}
	}

	return nil
}

func (s *Housekeeper) runDeleteOldDispatches(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-old-delete")
	defer span.End()

	if err := s.deleteOldDispatches(ctx); err != nil {
		s.logger.Error("failed to remove old dispatches", zap.Error(err))
		return err
	}

	return nil
}

// deleteOldDispatches deletes dispatches that are older than a certain number of days.
// This can probably be moved into the general housekeeper service.
func (s *Housekeeper) deleteOldDispatches(ctx context.Context) error {
	tDispatch := table.FivenetCentrumDispatches

	stmt := tDispatch.
		SELECT(
			tDispatch.ID.AS("dispatch_id"),
		).
		FROM(
			tDispatch,
		).
		WHERE(jet.AND(
			tDispatch.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(DeleteDispatchDays, jet.DAY)),
			),
		)).
		// Get 75 at a time
		LIMIT(75)

	var dest []*struct {
		DispatchID uint64
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	errs := multierr.Combine()
	for _, ds := range dest {
		if err := s.dispatches.Delete(ctx, ds.DispatchID, true); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}

	return errs
}

func (s *Housekeeper) runDeleteOldDispatchesFromKV(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-old-delete-kv")
	defer span.End()

	if err := s.deleteOldDispatchesFromKV(ctx); err != nil {
		s.logger.Error("failed to remove old dispatches from kv", zap.Error(err))
		return err
	}

	return nil
}

func (s *Housekeeper) deleteOldDispatchesFromKV(ctx context.Context) error {
	errs := multierr.Combine()

	keyIter, err := s.dispatches.Store().KV().ListKeysFiltered(ctx, "id.*")
	if err != nil {
		s.logger.Error("failed to list dispatches from KV", zap.Error(err))
		return err
	}
	keysCh := keyIter.Keys()
	for key := range keysCh {
		if key == "" {
			continue
		}

		dspId, err := centrumutils.ExtractIDString(key)
		if err != nil {
			s.logger.Error("failed to extract dispatch ID from key", zap.String("key", key), zap.Error(err))
			errs = multierr.Append(errs, fmt.Errorf("failed to extract dispatch ID from key %q: %w", key, err))
			continue
		}

		dsp, err := s.dispatches.Store().Get(dspId)
		if err != nil {
			s.logger.Error("failed to get dispatch from KV", zap.String("key", key), zap.Error(err))

			if err := s.dispatches.Store().Delete(ctx, key); err != nil {
				s.logger.Error("failed to delete unavailable dispatch from KV", zap.String("key", key), zap.Error(err))
			}
			continue
		}

		if (
		// Dispatches older than 3 hours will be removed from the KV store (not the database)
		dsp.CreatedAt != nil && time.Since(dsp.CreatedAt.AsTime()) > 3*time.Hour) ||
			// Remove nil status dispatches
			dsp.Status == nil ||
			// "Completed" dispatches with their status being older than 15 minutes
			(centrumutils.IsStatusDispatchComplete(dsp.Status.Status) &&
				time.Since(dsp.Status.CreatedAt.AsTime()) > 15*time.Minute) {
			s.logger.Debug("old dispatch deleted from kv", zap.Uint64("dispatch_id", dsp.Id))

			if err := s.dispatches.Delete(ctx, dsp.Id, false); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to delete dispatch from KV. %w", err))
				continue
			}
		}
	}

	return errs
}
