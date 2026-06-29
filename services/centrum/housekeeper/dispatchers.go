package housekeeper

import (
	"context"
	"errors"
	"fmt"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func (s *Housekeeper) runCleanupDispatchers(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatchers_cleanup")
	defer span.End()

	if err := s.cleanupDispatchers(ctx); err != nil {
		s.logger.Error("failed to remove old dispatchers", zap.Error(err))
		return err
	}

	return nil
}

func (s *Housekeeper) cleanupDispatchers(ctx context.Context) error {
	var errs error
	s.dispatchers.Range(func(job string, value *centrumdispatchers.Dispatchers) bool {
		for _, user := range value.GetDispatchers() {
			um, err := s.tracker.GetUserMapping(user.GetUserId())
			if err != nil {
				if !errors.Is(err, jetstream.ErrKeyNotFound) {
					errs = multierr.Append(
						errs,
						fmt.Errorf(
							"unable to get user %d mapping for %s dispatchers. %w",
							user.GetUserId(),
							job,
							err,
						),
					)
				}
			}

			// If user mapping is not nil and not hidden, dispatcher should still be valid.
			if um != nil && !um.Hidden {
				continue
			}

			if err := s.dispatchers.SetUserState(ctx, job, user.GetUserId(), false); err != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf(
						"failed to remove user %d from %s dispatchers. %w",
						user.GetUserId(),
						job,
						err,
					),
				)
				continue
			}
		}

		return true
	})

	return errs
}
