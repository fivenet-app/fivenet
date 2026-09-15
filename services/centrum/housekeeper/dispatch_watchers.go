package housekeeper

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

// These functions own watcher lifecycle and restart behavior. The underlying
// watcher implementations remain with their domain-specific handlers.
func (s *Housekeeper) runIdleWatcher(ctx context.Context) {
	for {
		if err := s.idleWatcher(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.recordWatcherRestart("idle", err)
			s.logger.Error("idle watcher stopped", zap.Error(err))
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
		if err := s.dispatchCleanupWatcher(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.recordWatcherRestart("dispatch_cleanup", err)
			s.logger.Error("dispatch cleanup watcher stopped", zap.Error(err))
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
		if err := s.projectionCleanupWatcher(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.recordWatcherRestart("projection_cleanup", err)
			s.logger.Error("dispatch projection cleanup watcher stopped", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) runDispatchAssignmentExpirationWatcher(ctx context.Context) {
	for {
		if err := s.dispatchAssignmentExpirationWatcher(
			ctx,
		); err != nil &&
			!errors.Is(err, context.Canceled) {
			s.recordWatcherRestart("dispatch_assignment_expiration", err)
			s.logger.Error("dispatch assignment expiration watcher stopped", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}
