package housekeeper

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func (s *Housekeeper) runUserChangesWatch(ctx context.Context) {
	for {
		if err := s.watchUserChanges(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.logger.Error("failed to watch user changes", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			s.logger.Info("stopping user changes watcher")
			return

		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Housekeeper) watchUserChanges(ctx context.Context) error {
	watch, err := s.tracker.Subscribe(ctx)
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case e := <-watch.Updates():
			if e == nil {
				s.logger.Error("received nil user changes event, skipping")
				continue
			}
			switch e.Operation() {
			case jetstream.KeyValuePut:
				func() {
					ctx, span := s.tracer.Start(ctx, "centrum.watch_users.put")
					defer span.End()

					userMarker, err := e.Value()
					if err != nil {
						s.logger.Error(
							"failed to get user marker from usermarker put event",
							zap.String("key", e.Key()),
							zap.Error(err),
						)
						return
					}

					if err := s.units.SyncUserUnitMapping(ctx, userMarker.GetUserId()); err != nil {
						s.logger.Error(
							"failed to sync user unit mapping for usermarker put event",
							zap.Int32("user_id", userMarker.GetUserId()),
							zap.Error(err),
						)
						return
					}
				}()

			case jetstream.KeyValueDelete, jetstream.KeyValuePurge:
				func() {
					ctx, span := s.tracer.Start(ctx, "centrum.watch_users.delete")
					defer span.End()

					userMarker, err := e.Value()
					if err != nil {
						s.logger.Error(
							"failed to get user marker from usermarker delete event",
							zap.String("key", e.Key()),
							zap.Error(err),
						)
						return
					}

					// Check if user is signed on as dispatcher
					if !s.helpers.CheckIfUserIsDispatcher(
						ctx,
						userMarker.GetJob(),
						userMarker.GetUserId(),
					) {
						return
					}

					// Remove user from dispatchers
					if err := s.dispatchers.SetUserState(
						ctx,
						userMarker.GetJob(),
						userMarker.GetUserId(),
						false,
					); err != nil {
						s.logger.Error(
							"failed to remove user as dispatcher for usermarker delete event",
							zap.Int32("user_id", userMarker.GetUserId()),
							zap.Error(err),
						)
						return
					}
				}()
			}
		}
	}
}
