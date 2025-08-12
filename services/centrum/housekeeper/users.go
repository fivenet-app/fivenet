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
	userCh, err := s.tracker.Subscribe(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event := <-userCh:
			if event == nil {
				s.logger.Error("received nil user changes event, skipping")
				continue
			}

			func() {
				ctx, span := s.tracer.Start(ctx, "centrum.watch-users")
				defer span.End()

				if event.Operation() != jetstream.KeyValuePut {
					return
				}
				userMarker, err := event.Value()
				if err != nil {
					s.logger.Error("failed to get user marker from event", zap.Error(err))
					return
				}

				if _, err := s.tracker.GetUserMapping(userMarker.GetUserId()); err != nil {
					return
				}

				unitId, err := s.units.LoadUnitIDForUserID(ctx, userMarker.GetUserId())
				if err != nil {
					s.logger.Error("failed to load user unit id", zap.Error(err))
					return
				}

				if err := s.tracker.SetUserMappingForUser(ctx, userMarker.GetUserId(), &unitId); err != nil {
					s.logger.Error("failed to update user unit id mapping in kv", zap.Error(err))
					return
				}
			}()
		}
	}
}
