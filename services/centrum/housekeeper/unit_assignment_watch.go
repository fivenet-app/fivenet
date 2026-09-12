package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *Housekeeper) runUnitAssignmentWatch(ctx context.Context) {
	for {
		if err := s.watchUnitAssignments(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.recordWatcherRestart("unit_assignments", err)
			s.logger.Error("unit assignment watcher stopped", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) watchUnitAssignments(ctx context.Context) error {
	watch, err := s.unitAssignmentWatchSource.WatchAll(ctx)
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-watch.Updates():
			if !ok {
				return errWatcherUpdatesClosed
			}
			if event == nil || event.Operation() != jetstream.KeyValuePut {
				continue
			}

			unit, err := event.Value()
			if err != nil {
				s.logger.Warn("cannot read unit assignment watch value", zap.Error(err))
				continue
			}
			if unit == nil || unit.GetId() <= 0 || len(unit.GetUsers()) > 0 {
				continue
			}

			s.metrics.IncHousekeeperEvent("unit_assignments", "empty_unit")
			removed, err := s.removeEmptyUnit(ctx, unit)
			if err != nil {
				s.metrics.IncHousekeeperEvent(
					"unit_assignments",
					"dispatch_assignment_removal_failed",
				)
				s.logger.Error(
					"failed to remove empty unit from assigned dispatches",
					zap.Int64("unit_id", unit.GetId()),
					zap.Error(err),
				)
				continue
			}
			s.metrics.AddHousekeeperEvents(
				"unit_assignments",
				"dispatch_assignments_removed",
				removed,
			)
		}
	}
}

func (s *Housekeeper) removeDispatchesFromEmptyUnit(
	ctx context.Context,
	unit *centrumunits.Unit,
) (int, error) {
	assignments := table.FivenetCentrumDispatchesAsgmts
	stmt := assignments.
		SELECT(assignments.DispatchID).
		FROM(assignments).
		WHERE(assignments.UnitID.EQ(mysql.Int64(unit.GetId())))

	rows := []*struct {
		DispatchID int64 `alias:"dispatch_id"`
	}{}
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return 0, fmt.Errorf("failed to find dispatches assigned to unit %d. %w", unit.GetId(), err)
	}

	removed := 0
	var errs error
	for _, row := range rows {
		if row == nil || row.DispatchID <= 0 {
			continue
		}
		if err := s.dispatches.UpdateAssignments(
			ctx,
			nil,
			nil,
			row.DispatchID,
			nil,
			[]int64{unit.GetId()},
			time.Time{},
		); err != nil {
			errs = errors.Join(errs, fmt.Errorf(
				"failed to remove unit %d from dispatch %d. %w",
				unit.GetId(),
				row.DispatchID,
				err,
			))
			continue
		}
		removed++
	}

	return removed, errs
}
