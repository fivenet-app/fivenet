package housekeeper

import (
	"context"
	"errors"
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

func (s *Housekeeper) runIdleWatcher(ctx context.Context) {
	for {
		if err := s.idleWatcher(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.recordWatcherRestart("idle", err)
				s.logger.Error("idle watcher stopped", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) runDispatchCleanupWatcher(ctx context.Context) {
	for {
		if err := s.dispatchCleanupWatcher(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.recordWatcherRestart("dispatch_cleanup", err)
				s.logger.Error("dispatch cleanup watcher stopped", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) runProjectionCleanupWatcher(ctx context.Context) {
	for {
		if err := s.projectionCleanupWatcher(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.recordWatcherRestart("projection_cleanup", err)
				s.logger.Error("dispatch projection cleanup watcher stopped", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
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
				return errWatcherUpdatesClosed
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
			dsp, err := s.getDispatchProjection(ctx, id)
			if err != nil || dsp == nil ||
				centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) {
				continue // Already handled elsewhere
			}

			s.cancelDispatch(ctx, dsp)
		}
	}
}

func (s *Housekeeper) dispatchCleanupWatcher(ctx context.Context) error {
	watch, err := s.dispatches.IdleStore().Watch(ctx, "cleanup.*")
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
				return errWatcherUpdatesClosed
			}
			if e == nil ||
				(e.Operation() != jetstream.KeyValueDelete && e.Operation() != jetstream.KeyValuePurge) {
				continue
			}

			id, err := strconv.ParseInt(strings.TrimPrefix(e.Key(), "cleanup.id."), 10, 64)
			if err != nil || id <= 0 {
				s.logger.Warn(
					"received invalid dispatch cleanup timer key",
					zap.String("key", e.Key()),
					zap.Error(err),
				)
				continue
			}

			dsp, err := s.getDispatchProjection(ctx, id)
			if err != nil || dsp == nil ||
				!centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) {
				continue
			}

			if err := s.deleteDispatch(ctx, id, false); err != nil {
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
				return errWatcherUpdatesClosed
			}
			if e == nil ||
				(e.Operation() != jetstream.KeyValueDelete && e.Operation() != jetstream.KeyValuePurge) {
				continue
			}

			id, err := strconv.ParseInt(strings.TrimPrefix(e.Key(), "projection.id."), 10, 64)
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
	dsp, err := s.getDispatchProjection(ctx, id)
	if err != nil || dsp == nil {
		return
	}

	if dsp.GetCreatedAt() != nil &&
		time.Until(dsp.GetCreatedAt().AsTime().Add(dispatches.ProjectionTTL)) > 0 {
		if err := s.scheduleProjectionCleanup(ctx, id, dsp.GetCreatedAt()); err != nil {
			s.logger.Error(
				"failed to reschedule dispatch projection cleanup",
				zap.Int64("dispatch_id", id),
				zap.Error(err),
			)
		}
		return
	}

	if err := s.deleteDispatch(ctx, id, false); err != nil {
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
	if _, err := s.updateDispatchStatus(ctx, dsp.GetId(), &centrumdispatches.DispatchStatus{
		CreatedAt:  timestamp.Now(),
		DispatchId: dsp.GetId(),
		Status:     centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED,
	}); err != nil {
		s.logger.Error(
			"failed to update dispatch status to cancelled",
			zap.Int64("dispatch_id", dsp.GetId()),
			zap.Error(err),
		)
		return
	}
	if err := s.addDispatchAttribute(
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
	if err := s.deleteDispatch(ctx, dsp.GetId(), false); err != nil {
		s.logger.Error(
			"failed to delete idle dispatch",
			zap.Int64("dispatch_id", dsp.GetId()),
			zap.Error(err),
		)
	}
}
