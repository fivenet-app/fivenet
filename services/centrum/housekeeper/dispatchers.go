package housekeeper

import (
	"context"
	"fmt"
	"strconv"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	dispatchersCheckedAttr = "dispatchers_checked"
	dispatchersRemovedAttr = "dispatchers_removed"
	dispatchersJobsAttr    = "jobs_affected"
	dispatchersUsersAttr   = "users_processed"
)

func (s *Housekeeper) runCleanupDispatchers(ctx context.Context, data *cron.CronjobData) error {
	ctx, span := s.tracer.Start(ctx, "centrum.dispatchers_cleanup")
	defer span.End()

	dest := &cron.GenericCronData{
		Attributes: map[string]string{},
	}
	if err := data.Unmarshal(dest); err != nil {
		s.logger.Warn("failed to unmarshal cleanup dispatchers cron data", zap.Error(err))
	}

	dispatchersChecked, dispatchersRemoved, usersProcessed, jobsAffected, err := s.cleanupDispatchers(
		ctx,
	)
	if err != nil {
		s.logger.Error("failed to remove old dispatchers", zap.Error(err))
		return err
	}

	dest.SetAttribute(dispatchersCheckedAttr, strconv.Itoa(dispatchersChecked))
	dest.SetAttribute(dispatchersRemovedAttr, strconv.Itoa(dispatchersRemoved))
	dest.SetAttribute(dispatchersUsersAttr, strconv.Itoa(usersProcessed))
	dest.SetAttribute(dispatchersJobsAttr, strconv.Itoa(jobsAffected))

	if err := data.MarshalFrom(dest); err != nil {
		return fmt.Errorf("failed to marshal updated cleanup dispatchers cron data. %w", err)
	}

	return nil
}

func (s *Housekeeper) cleanupDispatchers(ctx context.Context) (int, int, int, int, error) {
	var errs error
	dispatchersChecked := 0
	dispatchersRemoved := 0
	usersProcessed := 0
	jobsAffected := 0
	s.dispatchers.Range(func(job string, value *centrumdispatchers.Dispatchers) bool {
		dispatchersChecked++
		jobsAffected++
		for _, user := range value.GetDispatchers() {
			usersProcessed++
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
			dispatchersRemoved++
		}

		return true
	})

	return dispatchersChecked, dispatchersRemoved, usersProcessed, jobsAffected, errs
}
