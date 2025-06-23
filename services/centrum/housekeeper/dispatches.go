package housekeeper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/paulmach/orb"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func (s *Housekeeper) loadNewDispatches(ctx context.Context, data *cron.CronjobData) error {
	tDispatch := table.FivenetCentrumDispatches.AS("dispatch")
	// Load dispatches with null postal field (they are considered "new")
	if err := s.dispatches.LoadFromDB(ctx, tDispatch.Postal.IS_NULL()); err != nil {
		s.logger.Error("failed loading new dispatches from DB", zap.Error(err))
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

	errs := multierr.Combine()

	if err := s.deleteOldDispatches(ctx); err != nil {
		s.logger.Error("failed to remove old dispatches", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	if err := s.deleteOldDispatchesFromKV(ctx); err != nil {
		s.logger.Error("failed to remove old dispatches from kv", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	if err := s.deleteOldUnitStatus(ctx); err != nil {
		s.logger.Error("failed to remove old unit status", zap.Error(err))
		errs = multierr.Append(errs, err)
	}

	return errs
}

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

func (s *Housekeeper) deleteOldDispatchesFromKV(ctx context.Context) error {
	errs := multierr.Combine()

	dsps := s.dispatches.List(ctx, nil)
	for _, dsp := range dsps {
		if dsp == nil {
			continue
		}

		if (
		// Created dispatches older than the delete dispatch days amount
		dsp.CreatedAt != nil && time.Since(dsp.CreatedAt.AsTime()) > DeleteDispatchDays*24*time.Hour) ||
			// Remove nil status dispatches
			dsp.Status == nil ||
			// "Completed" dispatches with their status being older than 15 minutes
			(centrumutils.IsStatusDispatchComplete(dsp.Status.Status) &&
				time.Since(dsp.Status.CreatedAt.AsTime()) > 15*time.Minute) {
			s.logger.Debug("old dispatch deleted from kv", zap.Uint64("dispatch_id", dsp.Id))

			if _, err := s.dispatches.UpdateStatus(ctx, dsp.Id, &centrum.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: dsp.Id,
				Status:     centrum.StatusDispatch_STATUS_DISPATCH_DELETED,
			}); err != nil {
				s.logger.Error("failed to update dispatch status to deleted", zap.Uint64("dispatch_id", dsp.Id), zap.Error(err))
			}

			if err := s.dispatches.Delete(ctx, dsp.Id, false); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("failed to delete dispatch from KV. %w", err))
				continue
			}
		}
	}

	return errs
}

func (s *Housekeeper) runDispatchDeduplication(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatch-deduplicatation")
	defer span.End()

	if err := s.deduplicateDispatches(ctx); err != nil {
		s.logger.Error("failed to deduplicate dispatches", zap.Error(err))
	}

	return nil
}

func (s *Housekeeper) deduplicateDispatches(ctx context.Context) error {
	wg := sync.WaitGroup{}

	for _, settings := range s.settings.List(ctx) {
		job := settings.Job
		locs, ok := s.dispatches.GetLocations(job)
		if locs == nil || !ok {
			continue
		}

		wg.Add(1)
		go func(job string) {
			defer wg.Done()

			dsps := s.dispatches.Filter(ctx, []string{job}, nil, []centrum.StatusDispatch{
				centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
				centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
				centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
				centrum.StatusDispatch_STATUS_DISPATCH_DELETED,
			})

			if len(dsps) <= 1 {
				return
			}

			removedCount := 0
			dispatchIds := map[uint64]any{}
			for _, dsp := range dsps {
				// Skip handled dispatches
				if _, ok := dispatchIds[dsp.Id]; ok {
					continue
				}

				// Add the handled dispatch to the list
				dispatchIds[dsp.Id] = nil

				if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
					continue
				}

				// Iterate over close by dispatches and collect the active ones (if locations are available)
				closestsDsp := locs.KNearest(dsp.Point(), 8, func(p orb.Pointer) bool {
					return p.(*centrum.Dispatch).Id != dsp.Id
				}, 45.0)
				s.logger.Debug("deduplicating dispatches", zap.Strings("job", dsp.Jobs.GetJobStrings()), zap.Uint64("dispatch_id", dsp.Id), zap.Int("closeby_dsps", len(closestsDsp)))

				var refs *centrum.DispatchReferences
				if dsp.References != nil {
					refs = dsp.References
				} else {
					refs = &centrum.DispatchReferences{}
				}

				activeDispatchesCloseBy := []*centrum.Dispatch{}
				for _, dest := range closestsDsp {
					if dest == nil {
						continue
					}

					closeByDsp := dest.(*centrum.Dispatch)
					if closeByDsp.Status != nil && centrumutils.IsStatusDispatchComplete(closeByDsp.Status.Status) {
						continue
					}

					if closeByDsp.CreatedAt != nil && time.Since(closeByDsp.CreatedAt.AsTime()) >= 3*time.Minute {
						continue
					}

					// Skip dispatches that are marked as multiple or duplicate dispatches they have already been handled
					if closeByDsp.Attributes != nil && (closeByDsp.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE) || closeByDsp.Attributes.Has(centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE)) {
						continue
					}

					// Add close by dispatch as a reference
					refs.Add(&centrum.DispatchReference{
						TargetDispatchId: closeByDsp.Id,
						ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATED_BY,
					})

					activeDispatchesCloseBy = append(activeDispatchesCloseBy, closeByDsp)
				}

				// Prevent unnecessary updates to the dispatch
				if len(activeDispatchesCloseBy) == 0 {
					continue
				}

				// Add "multiple" attribute when multiple dispatches close by
				if err := s.dispatches.AddAttributeToDispatch(ctx, dsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_MULTIPLE); err != nil {
					s.logger.Error("failed to update original dispatch attribute", zap.Error(err))
				}

				// Set dispatch references on dispatch
				if err := s.dispatches.AddReferencesToDispatch(ctx, dsp, refs.References...); err != nil {
					s.logger.Error("failed to update duplicate dispatch references", zap.Error(err))
				}

				sourceDspRef := &centrum.DispatchReference{
					TargetDispatchId: dsp.Id,
					ReferenceType:    centrum.DispatchReferenceType_DISPATCH_REFERENCE_TYPE_DUPLICATE_OF,
				}

				for _, closeByDsp := range activeDispatchesCloseBy {
					// Already took care of the dispatch
					if _, ok := dispatchIds[closeByDsp.Id]; ok {
						continue
					}
					dispatchIds[closeByDsp.Id] = nil

					if closeByDsp.Status != nil && centrumutils.IsStatusDispatchComplete(closeByDsp.Status.Status) {
						continue
					}

					if err := s.dispatches.AddAttributeToDispatch(ctx, closeByDsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_DUPLICATE); err != nil {
						s.logger.Error("failed to update duplicate dispatch attribute", zap.Error(err))
					}

					if err := s.dispatches.AddReferencesToDispatch(ctx, closeByDsp, sourceDspRef); err != nil {
						s.logger.Error("failed to update duplicate dispatch references", zap.Error(err))
					}

					if _, err := s.dispatches.UpdateStatus(ctx, closeByDsp.Id, &centrum.DispatchStatus{
						CreatedAt:  timestamp.Now(),
						DispatchId: closeByDsp.Id,
						Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
					}); err != nil {
						s.logger.Error("failed to update duplicate dispatch status", zap.Error(err))
						return
					}

					toRemove := []uint64{}
					for _, ua := range closeByDsp.Units {
						toRemove = append(toRemove, ua.UnitId)
					}
					if err := s.dispatches.UpdateAssignments(ctx, nil, closeByDsp.Id, nil, toRemove, time.Time{}); err != nil {
						s.logger.Error("failed to remove assigned units from duplicate dispatch", zap.Error(err))
						return
					}

					removedCount++

					if removedCount >= MaxCancelledDispatchesPerRun {
						break
					}
				}
			}
		}(job)
	}

	wg.Wait()

	return nil
}
