package housekeeper

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	centrumutils "github.com/fivenet-app/fivenet/v2025/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *Housekeeper) runIdleWatcher(ctx context.Context) {
	for {
		if err := s.idleWatcher(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
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

func (s *Housekeeper) idleWatcher(ctx context.Context) error {
	watch, err := s.dispatches.IdleStore().Watch(ctx, "idle.*")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case e := <-watch.Updates():
			if e == nil || (e.Operation() != jetstream.KeyValueDelete &&
				e.Operation() != jetstream.KeyValuePurge) {
				continue // we only care about expiry events
			}

			idStr := strings.TrimPrefix(e.Key(), "idle.")
			id, _ := strconv.ParseInt(idStr, 10, 64)

			// double-check it is still open, then cancel & archive
			dsp, err := s.dispatches.Get(ctx, id)
			if err != nil || dsp == nil ||
				centrumutils.IsStatusDispatchComplete(dsp.GetStatus().GetStatus()) {
				continue // Already handled elsewhere
			}

			if _, err := s.dispatches.UpdateStatus(ctx, id, &centrum.DispatchStatus{
				CreatedAt:  timestamp.Now(),
				DispatchId: id,
				Status:     centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
			}); err != nil {
				s.logger.Error(
					"failed to update dispatch status to cancelled",
					zap.Int64("dispatch_id", id),
					zap.Error(err),
				)
			}
			if err := s.dispatches.AddAttributeToDispatch(ctx, dsp, centrum.DispatchAttribute_DISPATCH_ATTRIBUTE_TOO_OLD); err != nil {
				s.logger.Error(
					"failed to add too old attribute to cancelled dispatch",
					zap.Int64("dispatch_id", id),
					zap.Error(err),
				)
			}

			// Remove from kv so the UI gets the event
			if err := s.dispatches.Delete(ctx, id, false); err != nil {
				s.logger.Error(
					"failed to delete idle dispatch",
					zap.Int64("dispatch_id", id),
					zap.Error(err),
				)
			}
		}
	}
}
