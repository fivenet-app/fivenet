package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/tracker"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/multierr"
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

				if event.Operation() == jetstream.KeyValuePut {
					userMarker, err := event.Value()
					if err != nil {
						s.logger.Error("failed to get user marker from event", zap.Error(err))
						return
					}

					if _, err := s.tracker.GetUserMapping(userMarker.UserId); err != nil {
						return
					}

					unitId, err := s.units.LoadUnitIDForUserID(ctx, userMarker.UserId)
					if err != nil {
						s.logger.Error("failed to load user unit id", zap.Error(err))
						return
					}

					if err := s.tracker.SetUserMappingForUser(ctx, userMarker.UserId, &unitId); err != nil {
						s.logger.Error("failed to update user unit id mapping in kv", zap.Error(err))
						return
					}
				} else if event.Operation() == jetstream.KeyValueDelete || event.Operation() == jetstream.KeyValuePurge {
					userId, job, _, err := tracker.DecodeUserMarkerKey(event.Key())
					if err != nil {
						s.logger.Error("failed to decode user marker key", zap.Error(err), zap.String("key", string(event.Key())))
						return
					}

					if err := s.handleRemovedUser(ctx, job, userId); err != nil {
						s.logger.Error("failed to handle removed user", zap.Int32("user_id", userId), zap.Error(err))
					}
				}
			}()
		}
	}
}

func (s *Housekeeper) handleRemovedUser(ctx context.Context, job string, userId int32) error {
	var errs error
	if s.helpers.CheckIfUserIsDispatcher(ctx, job, userId) {
		if err := s.dispatchers.SetUserState(ctx, job, userId, false); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to remove user from disponents. %w", err))
		}
	}

	um, err := s.tracker.GetUserMapping(userId)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to get user unit mapping. %w", err))
		// User not in any unit, nothing to do
		return errs
	}

	if um != nil && um.UnitId != nil && *um.UnitId > 0 {
		if err := s.units.UpdateUnitAssignments(ctx, job, &userId, *um.UnitId, nil, []int32{userId}); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to remove user from unit. %w", err))
		}
	}

	return errs
}
