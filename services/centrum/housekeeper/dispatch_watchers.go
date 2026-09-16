package housekeeper

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

// watcherFunc is a long-running watcher which returns when it needs to be
// restarted or when its context is cancelled.
type watcherFunc func(context.Context) error

// runWatcher owns the common watcher lifecycle and restart behavior. The
// underlying watcher implementations remain with their domain-specific
// handlers.
func (s *Housekeeper) runWatcher(ctx context.Context, name string, watch watcherFunc) {
	for {
		if err := watch(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.recordWatcherRestart(name, err)
			s.logger.Error(
				"housekeeper watcher stopped",
				zap.String("watcher", name),
				zap.Error(err),
			)
		}
		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) runIdleWatcher(ctx context.Context) {
	s.runWatcher(ctx, "idle", s.idleWatcher)
}

func (s *Housekeeper) runDispatchCleanupWatcher(ctx context.Context) {
	s.runWatcher(ctx, "dispatch_cleanup", s.dispatchCleanupWatcher)
}

func (s *Housekeeper) runProjectionCleanupWatcher(ctx context.Context) {
	s.runWatcher(ctx, "projection_cleanup", s.projectionCleanupWatcher)
}

func (s *Housekeeper) runDispatchAssignmentExpirationWatcher(ctx context.Context) {
	s.runWatcher(
		ctx,
		"dispatch_assignment_expiration",
		s.dispatchAssignmentExpirationWatcher,
	)
}
