package housekeeper

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

// dispatchAssignmentExpirationWatcher reacts only to MaxAge purge markers.
// Explicit Delete calls cancel a timer and must not be mistaken for expiry.
func (s *Housekeeper) dispatchAssignmentExpirationWatcher(ctx context.Context) error {
	watch, err := s.assignmentExpirationSource.IdleStore().Watch(ctx, "assignment.*.*")
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-watch.Updates():
			if !ok {
				return watcherUpdatesClosedError(ctx)
			}
			if event == nil || event.Operation() != jetstream.KeyValuePurge {
				continue
			}

			parts := strings.Split(event.Key(), ".")
			if len(parts) != 3 {
				s.logger.Warn(
					"received invalid dispatch assignment timer key",
					zap.String("key", event.Key()),
				)
				continue
			}
			dispatchID, dispatchErr := strconv.ParseInt(parts[1], 10, 64)
			unitID, unitErr := strconv.ParseInt(parts[2], 10, 64)
			if dispatchErr != nil || unitErr != nil || dispatchID <= 0 || unitID <= 0 {
				s.logger.Warn(
					"received invalid dispatch assignment timer key",
					zap.String("key", event.Key()),
					zap.Error(errors.Join(dispatchErr, unitErr)),
				)
				continue
			}

			// The timer is only a wake-up signal. Re-read the authoritative
			// projection so an acceptance, reassignment, or cancellation that
			// raced the expiry marker is never removed.
			dispatch, err := s.assignmentExpirationSource.Get(ctx, dispatchID)
			if err != nil || dispatch == nil {
				continue
			}
			assignment := slices.IndexFunc(
				dispatch.GetUnits(),
				func(assignment *centrumdispatches.DispatchAssignment) bool {
					return assignment.GetUnitId() == unitID
				},
			)
			if assignment < 0 || dispatch.GetUnits()[assignment].GetExpiresAt() == nil ||
				dispatch.GetUnits()[assignment].GetExpiresAt().
					AsTime().
					After(time.Now().Add(-2*time.Second)) {
				continue
			}

			if err := s.assignmentExpirationSource.UpdateAssignments(
				ctx,
				nil,
				nil,
				dispatchID,
				nil,
				[]int64{unitID},
				time.Time{},
			); err != nil {
				s.logger.Error(
					"failed to remove expired dispatch assignment",
					zap.Int64("dispatch_id", dispatchID),
					zap.Int64("unit_id", unitID),
					zap.Error(err),
				)
				s.metrics.IncHousekeeperEvent("dispatch_assignment_expiration", "failed")
				continue
			}
			s.metrics.IncHousekeeperEvent("dispatch_assignment_expiration", "expired")
		}
	}
}

func (s *Housekeeper) idleWatcher(ctx context.Context) error {
	watch, err := s.dispatches.IdleStore().Watch(ctx, "idle.*")
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e, ok := <-watch.Updates():
			if !ok {
				return watcherUpdatesClosedError(ctx)
			}
			// Ignore nil event
			if e == nil {
				continue
			}
			// We only care about expiry events
			if e.Operation() != jetstream.KeyValueDelete &&
				e.Operation() != jetstream.KeyValuePurge {
				continue
			}

			idStr := strings.TrimPrefix(e.Key(), "idle.")
			id, _ := strconv.ParseInt(idStr, 10, 64)

			// double-check it is still open, then cancel & archive
			dsp, err := s.dispatchLifecycle.Get(ctx, id)
			if err != nil || dsp == nil ||
				centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) {
				continue // Already handled elsewhere
			}

			s.cancelDispatch(ctx, dsp)
		}
	}
}

func (s *Housekeeper) dispatchCleanupWatcher(ctx context.Context) error {
	return s.watchDispatchCleanup(ctx, s.dispatches.IdleStore())
}

func (s *Housekeeper) watchDispatchCleanup(
	ctx context.Context,
	idleStore jetstream.KeyValue,
) error {
	watch, err := idleStore.Watch(ctx, "cleanup.*")
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e, ok := <-watch.Updates():
			if !ok {
				return watcherUpdatesClosedError(ctx)
			}
			// A regular Delete cancels the timer when a completed dispatch becomes
			// active again. KeyTTL expiry emits a Purge marker, which is the only
			// operation that must remove the completed projection.
			if e == nil || e.Operation() != jetstream.KeyValuePurge {
				continue
			}

			id, err := strconv.ParseInt(strings.TrimPrefix(e.Key(), "cleanup."), 10, 64)
			if err != nil || id <= 0 {
				s.logger.Warn(
					"received invalid dispatch cleanup timer key",
					zap.String("key", e.Key()),
					zap.Error(err),
				)
				continue
			}

			dsp, err := s.dispatchLifecycle.Get(ctx, id)
			if err != nil || dsp == nil ||
				!centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) {
				continue
			}

			if err := s.dispatchLifecycle.Delete(ctx, id, false); err != nil {
				s.logger.Error(
					"failed to delete completed dispatch",
					zap.Int64("dispatch_id", id),
					zap.Error(err),
				)
				s.metrics.IncHousekeeperEvent("dispatch_cleanup", "failed")
			} else {
				s.metrics.IncHousekeeperEvent("dispatch_cleanup", "deleted")
			}
		}
	}
}

func (s *Housekeeper) projectionCleanupWatcher(ctx context.Context) error {
	watch, err := s.dispatches.IdleStore().Watch(ctx, "projection.*")
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e, ok := <-watch.Updates():
			if !ok {
				return watcherUpdatesClosedError(ctx)
			}
			if e == nil ||
				(e.Operation() != jetstream.KeyValueDelete && e.Operation() != jetstream.KeyValuePurge) {
				continue
			}

			id, err := strconv.ParseInt(strings.TrimPrefix(e.Key(), "projection."), 10, 64)
			if err != nil || id <= 0 {
				s.logger.Warn(
					"received invalid dispatch projection cleanup key",
					zap.String("key", e.Key()),
					zap.Error(err),
				)
				continue
			}

			s.handleProjectionCleanup(ctx, id)
		}
	}
}

func (s *Housekeeper) handleProjectionCleanup(ctx context.Context, id int64) {
	dsp, err := s.dispatchLifecycle.Get(ctx, id)
	if err != nil || dsp == nil {
		return
	}

	if dsp.GetCreatedAt() != nil &&
		time.Until(dsp.GetCreatedAt().AsTime().Add(dispatches.ProjectionTTL)) > 0 {
		if err := s.dispatchLifecycle.ScheduleProjectionCleanup(
			ctx,
			id,
			dsp.GetCreatedAt(),
		); err != nil {
			s.logger.Error(
				"failed to reschedule dispatch projection cleanup",
				zap.Int64("dispatch_id", id),
				zap.Error(err),
			)
		}
		return
	}

	if err := s.dispatchLifecycle.Delete(ctx, id, false); err != nil {
		s.logger.Error(
			"failed to delete expired dispatch projection",
			zap.Int64("dispatch_id", id),
			zap.Error(err),
		)
		s.metrics.IncHousekeeperEvent("projection_cleanup", "failed")
	} else {
		s.metrics.IncHousekeeperEvent("projection_cleanup", "deleted")
	}
}

func (s *Housekeeper) cancelDispatch(ctx context.Context, dsp *centrumdispatches.Dispatch) {
	if _, err := s.dispatchLifecycle.UpdateStatus(
		ctx,
		dsp.GetId(),
		&centrumdispatches.DispatchStatus{
			CreatedAt:  timestamp.Now(),
			DispatchId: dsp.GetId(),
			Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		},
	); err != nil {
		s.logger.Error(
			"failed to update dispatch status to cancelled",
			zap.Int64("dispatch_id", dsp.GetId()),
			zap.Error(err),
		)
		return
	}
	if err := s.dispatchLifecycle.AddAttributeToDispatch(
		ctx,
		dsp,
		centrumdispatches.DispatchAttribute_DISPATCH_ATTRIBUTE_TOO_OLD,
	); err != nil {
		s.logger.Error(
			"failed to add too old attribute to cancelled dispatch",
			zap.Int64("dispatch_id", dsp.GetId()),
			zap.Error(err),
		)
	}

	// Remove from kv so the UI gets the event
	if err := s.dispatchLifecycle.Delete(ctx, dsp.GetId(), false); err != nil {
		s.logger.Error(
			"failed to delete idle dispatch",
			zap.Int64("dispatch_id", dsp.GetId()),
			zap.Error(err),
		)
	}
}
