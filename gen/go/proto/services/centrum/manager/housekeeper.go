package manager

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/paulmach/orb"
	"github.com/puzpuzpuz/xsync/v3"
	"go.uber.org/zap"
)

const (
	MaxCancelledDispatchesPerRun = 3
)

func (s *Manager) housekeeper() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runHandleDispatchAssignmentExpiration()
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runArchiveDispatches()
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runDispatchDeduplication()
	}()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runRemoveDispatchesFromEmptyUnits()
	}()
}

func (s *Manager) runHandleDispatchAssignmentExpiration() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(1 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-assignment-expiration")
				defer span.End()

				if err := s.handleDispatchAssignmentExpiration(ctx); err != nil {
					s.logger.Error("failed to handle expired dispatch assignments", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Manager) runArchiveDispatches() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(4 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-archival")
				defer span.End()

				if err := s.archiveDispatches(ctx); err != nil {
					s.logger.Error("failed to archive dispatches", zap.Error(err))
				}

				if err := s.cleanupDispatches(ctx); err != nil {
					s.logger.Error("failed to cleanup dispatches", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Manager) runDispatchDeduplication() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(2 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-dispatch-deduplicatation")
				defer span.End()

				if err := s.deduplicateDispatches(ctx); err != nil {
					s.logger.Error("failed to deduplicate dispatches", zap.Error(err))
				}
			}()
		}
	}
}

func (s *Manager) runRemoveDispatchesFromEmptyUnits() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(5 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-units-empty")
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
			}()
		}
	}
}

// Handle expired dispatch unit assignments
func (s *Manager) handleDispatchAssignmentExpiration(ctx context.Context) error {
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
			tDispatchUnit.ExpiresAt.LT_EQ(jet.NOW()),
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
		for dispatchId, units := range dsps {
			dsp, ok := s.GetDispatch(job, dispatchId)
			if !ok {
				continue
			}

			if err := s.UpdateDispatchAssignments(ctx, job, nil, dsp, nil, units, time.Time{}); err != nil {
				return err
			}
		}
	}

	return nil
}

// Set `COMPLETED`/`CANCELLED` dispatches to status `ARCHIVED` when the status is older than 5 minutes
func (s *Manager) archiveDispatches(ctx context.Context) error {
	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.DispatchID.AS("dispatch_id"),
			tDispatch.Job.AS("job"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.DispatchID),
				),
		).
		// Dispatches that are at 5 minutes or older, have completed/cancelled or no status set
		WHERE(jet.AND(
			tDispatchStatus.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(5, jet.MINUTE)),
			),
			tDispatchStatus.ID.IS_NULL().OR(
				jet.AND(
					tDispatchStatus.ID.EQ(
						jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
					),
					tDispatchStatus.Status.IN(
						jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED)),
						jet.Int16(int16(dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED)),
					),
				),
			),
		))

	var dest []*struct {
		DispatchID uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, ds := range dest {
		dsp, ok := s.GetDispatch(ds.Job, ds.DispatchID)
		if !ok {
			continue
		}

		// Ignore already archived dispatches
		if dsp.Status != nil && dsp.Status.Status == dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED {
			continue
		}

		if err := s.UpdateDispatchStatus(ctx, ds.Job, dsp, &dispatch.DispatchStatus{
			DispatchId: dsp.Id,
			Status:     dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			UserId:     dsp.Status.UserId,
		}); err != nil {
			return err
		}

		if err := s.DeleteDispatch(ctx, dsp.Job, dsp.Id); err != nil {
			return err
		}
	}

	return nil
}

func (s *Manager) cleanupDispatches(ctx context.Context) error {
	stmt := tDispatch.
		SELECT(
			tDispatch.ID.AS("dispatch_id"),
			tDispatch.Job.AS("job"),
		).
		FROM(
			tDispatch,
		).
		WHERE(jet.AND(
			tDispatch.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(90, jet.MINUTE)),
			),
		))

	var dest []*struct {
		DispatchID uint64
		Job        string
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for _, ds := range dest {
		dsp, ok := s.GetDispatch(ds.Job, ds.DispatchID)
		if !ok {
			continue
		}

		if err := s.DeleteDispatch(ctx, dsp.Job, dsp.Id); err != nil {
			return err
		}
	}

	return nil
}

func (s *Manager) deduplicateDispatches(ctx context.Context) error {
	wg := sync.WaitGroup{}

	s.Dispatches.Range(func(job string, _ *xsync.MapOf[uint64, *dispatch.Dispatch]) bool {
		wg.Add(1)
		go func() {
			defer wg.Done()

			dsps := s.State.FilterDispatches(job, nil, []dispatch.StatusDispatch{
				dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
				dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED,
				dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED,
			})
			sort.Slice(dsps, func(i, j int) bool {
				return dsps[i].Id < dsps[j].Id
			})

			if len(dsps) <= 1 {
				return
			}

			removedCount := 0
			dispatchIds := map[uint64]interface{}{}
			for _, dsp := range dsps {
				// Add the handled dispatch to the list
				dispatchIds[dsp.Id] = nil

				if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
					continue
				}

				closestsDsp := s.State.DispatchLocations[dsp.Job].KNearest(dsp.Point(), 8, func(p orb.Pointer) bool {
					return p.(*dispatch.Dispatch).Id != dsp.Id
				}, 45.0)
				// Add "multiple" attribute when multiple dispatches close by
				if len(closestsDsp) > 0 {
					if err := s.addAttributeToDispatch(ctx, dsp, "multiple"); err != nil {
						s.logger.Error("failed to update original dispatch attribute", zap.Error(err))
					}
				}

				for _, dest := range closestsDsp {
					if dest == nil {
						continue
					}

					// Already took care of the dispatch
					closeByDsp := dest.(*dispatch.Dispatch)
					if _, ok := dispatchIds[closeByDsp.Id]; ok {
						continue
					}
					dispatchIds[closeByDsp.Id] = nil

					if closeByDsp.Status != nil && centrumutils.IsStatusDispatchComplete(closeByDsp.Status.Status) {
						continue
					}

					if closeByDsp.CreatedAt != nil && time.Since(closeByDsp.CreatedAt.AsTime()) >= 3*time.Minute {
						continue
					}

					s.State.DispatchLocations[closeByDsp.Job].Remove(closeByDsp, func(p orb.Pointer) bool {
						return p.(*dispatch.Dispatch).Id == closeByDsp.Id
					})

					if err := s.UpdateDispatchStatus(ctx, closeByDsp.Job, closeByDsp, &dispatch.DispatchStatus{
						DispatchId: closeByDsp.Id,
						Status:     dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED,
					}); err != nil {
						s.logger.Error("failed to update duplicate dispatch status", zap.Error(err))
						return
					}

					removedCount++

					if removedCount >= MaxCancelledDispatchesPerRun {
						break
					}
				}
			}
		}()

		return true
	})

	wg.Wait()

	return nil
}

func (s *Manager) addAttributeToDispatch(ctx context.Context, dsp *dispatch.Dispatch, attribute string) error {
	update := false
	if dsp.Attributes == nil {
		dsp.Attributes = &dispatch.Attributes{
			List: []string{attribute},
		}
		update = true
	} else {
		if !utils.InSlice(dsp.Attributes.List, attribute) {
			dsp.Attributes.List = append(dsp.Attributes.List, attribute)
			update = true
		}
	}

	if update {
		if err := s.UpdateDispatch(ctx, dsp.Job, nil, dsp); err != nil {
			return err
		}
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Manager) removeDispatchesFromEmptyUnits(ctx context.Context) error {
	s.Dispatches.Range(func(job string, _ *xsync.MapOf[uint64, *dispatch.Dispatch]) bool {
		dsps := s.State.FilterDispatches(job, nil, []dispatch.StatusDispatch{
			dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
			dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED,
		})

		for _, dsp := range dsps {
			for i := len(dsp.Units) - 1; i >= 0; i-- {
				if i > len(dsp.Units)-1 {
					break
				}

				unit, _ := s.GetUnit(job, dsp.Units[i].UnitId)
				// If unit isn't empty, continue with the loop
				if unit != nil && len(unit.Users) > 0 {
					continue
				}

				if err := s.UpdateDispatchAssignments(ctx, job, nil, dsp, nil, []uint64{unit.Id}, time.Time{}); err != nil {
					s.logger.Error("failed to remove empty unit from dispatch",
						zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Uint64("dispatch_id", dsp.Id), zap.Error(err))
					continue
				}
			}

			return true
		}

		return true
	})

	return nil
}

// Iterate over units to ensure that, e.g., an empty unit status is set to `unavailable`
func (s *Manager) cleanupUnitStatus(ctx context.Context) error {
	s.Units.Range(func(job string, units *xsync.MapOf[uint64, *dispatch.Unit]) bool {
		units.Range(func(id uint64, unit *dispatch.Unit) bool {
			if len(unit.Users) > 0 {
				return true
			}

			if unit.Status != nil && unit.Status.Status == dispatch.StatusUnit_STATUS_UNIT_UNAVAILABLE {
				return true
			}

			var userId *int32
			if unit.Status != nil && unit.Status.UserId != nil {
				userId = unit.Status.UserId
			}

			if err := s.UpdateUnitStatus(ctx, job, unit, &dispatch.UnitStatus{
				UnitId: unit.Id,
				Status: dispatch.StatusUnit_STATUS_UNIT_UNAVAILABLE,
				UserId: userId,
			}); err != nil {
				s.logger.Error("failed to update empty unit status to unavailable",
					zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Error(err))
				return true
			}

			return true
		})

		return true
	})

	return nil
}

func (s *Manager) checkUnitUsers(ctx context.Context) error {
	s.Units.Range(func(job string, _ *xsync.MapOf[uint64, *dispatch.Unit]) bool {
		units, ok := s.ListUnits(job)
		if !ok {
			return true
		}

		for _, unit := range units {
			if len(unit.Users) == 0 {
				return true
			}

			toRemove := []int32{}
			for i := len(unit.Users) - 1; i >= 0; i-- {
				if i > len(unit.Users)-1 {
					break
				}

				var userId int32
				if unit.Users[i].User != nil {
					userId = unit.Users[i].User.UserId
				} else {
					userId = unit.Users[i].UserId
				}

				if userId == 0 {
					continue
				}

				unitId, _ := s.UserIDToUnitID.Load(userId)
				// If user is in that unit and still on duty, nothing to do, otherwise remove the user from the unit
				if unit.Id == unitId && s.tracker.IsUserOnDuty(job, userId) {
					continue
				}

				toRemove = append(toRemove, userId)
			}

			if len(toRemove) > 0 {
				if err := s.UpdateUnitAssignments(ctx, job, nil, unit, nil, toRemove); err != nil {
					s.logger.Error("failed to remove off-duty users from unit",
						zap.String("job", unit.Job), zap.Uint64("unit_id", unit.Id), zap.Int32s("user_ids", toRemove), zap.Error(err))
				}
			}
		}

		return true
	})

	return nil
}
