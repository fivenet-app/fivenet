package housekeeper

import (
	"context"
	"fmt"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
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
			um, ok, err := s.tracker.GetUserMapping(user.GetUserId())
			if err != nil {
				errs = multierr.Append(
					errs,
					fmt.Errorf(
						"unable to get user %d mapping for %s dispatchers. %w",
						user.GetUserId(),
						job,
						err,
					),
				)
				continue
			}

			// Dispatcher is still valid when the user's job matches the dispatchers list job and the mapping is visible.
			if user.GetJob() == job && ok && um != nil && !um.Hidden {
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
