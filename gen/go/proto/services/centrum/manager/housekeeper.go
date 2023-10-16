package manager

import (
	"context"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/puzpuzpuz/xsync/v2"
	"go.uber.org/zap"
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

		case <-time.After(5 * time.Second):
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

func (s *Manager) runRemoveDispatchesFromEmptyUnits() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case <-time.After(3 * time.Second):
			func() {
				ctx, span := s.tracer.Start(s.ctx, "centrum-units-empty")
				defer span.End()

				if err := s.removeDispatchesFromEmptyUnits(ctx); err != nil {
					s.logger.Error("failed to clean empty units from dispatches", zap.Error(err))
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

		s.GetDispatchesMap(ds.Job).Delete(ds.DispatchID)
	}

	return nil
}

func (s *Manager) cleanupDispatches(ctx context.Context) error {
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
		WHERE(jet.AND(
			tDispatchStatus.CreatedAt.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(2, jet.HOUR)),
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

		s.GetDispatchesMap(ds.Job).Delete(ds.DispatchID)

		if err := s.DeleteDispatch(ctx, dsp.Job, dsp.Id); err != nil {
			return err
		}
	}

	return nil
}

// Remove empty units from dispatches (if no other unit is assigned to dispatch update status to UNASSIGNED) by
// iterating over the dispatches and making sure the assigned units aren't empty
func (s *Manager) removeDispatchesFromEmptyUnits(ctx context.Context) error {
	s.Dispatches.Range(func(job string, value *xsync.MapOf[uint64, *dispatch.Dispatch]) bool {
		value.Range(func(id uint64, dsp *dispatch.Dispatch) bool {
			for i := len(dsp.Units) - 1; i >= 0; i-- {
				unit, _ := s.GetUnit(job, dsp.Units[i].UnitId)
				// If unit isn't empty, continue with the loop
				if unit != nil && len(unit.Users) > 0 {
					continue
				}

				if err := s.UpdateDispatchAssignments(ctx, job, nil, dsp, nil, []uint64{dsp.Units[i].UnitId}, time.Time{}); err != nil {
					s.logger.Error("failed to remove empty unit from dispatch", zap.Error(err))
					continue
				}
			}

			return true
		})

		return true
	})

	return nil
}
