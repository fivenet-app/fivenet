package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type userDutyContext struct {
	job    string
	hidden bool
}

// userMarkerKeyContext extracts the stable user identity and job from a
// tracker marker key. KV delete and purge tombstones intentionally have no
// value, so their payload cannot be used for this purpose.
func userMarkerKeyContext(key string) (int32, userDutyContext, error) {
	userIDIndex := strings.LastIndexByte(key, '.')
	if userIDIndex < 1 || userIDIndex == len(key)-1 {
		return 0, userDutyContext{}, fmt.Errorf("invalid user marker key %q", key)
	}
	userID, err := strconv.ParseInt(key[userIDIndex+1:], 10, 32)
	if err != nil || userID <= 0 {
		return 0, userDutyContext{}, fmt.Errorf("invalid user marker key %q", key)
	}

	gradeIndex := strings.LastIndexByte(key[:userIDIndex], '.')
	if gradeIndex < 1 {
		return 0, userDutyContext{}, fmt.Errorf("invalid user marker key %q", key)
	}

	return int32(userID), userDutyContext{job: key[:gradeIndex], hidden: true}, nil
}

// dispatcherJobsToRemove returns jobs whose dispatcher state is no longer
// justified by a marker transition. Removal is idempotent in the KV store.
func dispatcherJobsToRemove(
	previous userDutyContext,
	seen bool,
	current userDutyContext,
) []string {
	var jobs []string
	if seen && previous.job != "" && (previous.job != current.job || current.hidden) {
		jobs = append(jobs, previous.job)
	}
	if current.hidden && current.job != "" && !slices.Contains(jobs, current.job) {
		jobs = append(jobs, current.job)
	}

	return jobs
}

func (s *Housekeeper) removeDispatchersForUser(
	ctx context.Context,
	userID int32,
	jobs []string,
) {
	for _, job := range jobs {
		s.metrics.IncHousekeeperEvent("user_changes", "dispatcher_removal_attempted")
		if err := s.setDispatcherState(ctx, job, userID, false); err != nil {
			s.metrics.IncHousekeeperEvent("user_changes", "dispatcher_removal_failed")
			s.logger.Error(
				"failed to remove stale dispatcher state for user marker change",
				zap.Int32("user_id", userID),
				zap.String("job", job),
				zap.Error(err),
			)
			continue
		}
		s.metrics.IncHousekeeperEvent("user_changes", "dispatcher_removal_succeeded")
	}
}

func (s *Housekeeper) runUserChangesWatch(ctx context.Context) {
	contexts := map[int32]userDutyContext{}
	for {
		if err := s.watchUserChanges(ctx, contexts); err != nil {
			if !errors.Is(err, context.Canceled) {
				s.recordWatcherRestart("user_changes", err)
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

func (s *Housekeeper) watchUserChanges(
	ctx context.Context,
	contexts map[int32]userDutyContext,
) error {
	watch, err := s.tracker.Subscribe(ctx)
	if err != nil {
		return err
	}
	defer watch.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case e, ok := <-watch.Updates():
			if !ok {
				return errWatcherUpdatesClosed
			}
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
					current := userDutyContext{
						job:    userMarker.GetJob(),
						hidden: userMarker.GetHidden(),
					}
					previous, seen := contexts[userMarker.GetUserId()]
					contexts[userMarker.GetUserId()] = current
					if seen && previous == current {
						return
					}
					s.removeDispatchersForUser(
						ctx,
						userMarker.GetUserId(),
						dispatcherJobsToRemove(previous, seen, current),
					)

					if err := s.syncUserUnitMapping(ctx, userMarker.GetUserId()); err != nil {
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

					userID, current, err := userMarkerKeyContext(e.Key())
					if err != nil {
						s.logger.Error(
							"failed to parse user marker key from usermarker delete event",
							zap.String("key", e.Key()),
							zap.Error(err),
						)
						return
					}

					// userLocStore contains one key per job/grade. A deletion can
					// therefore be cleanup of an obsolete key rather than an
					// off-duty transition. Consult the canonical marker before
					// signing the user out of every dispatcher state we have seen.
					if marker, found := s.tracker.GetUserMarkerById(userID); found &&
						marker != nil && !marker.GetHidden() {
						canonical := userDutyContext{
							job:    marker.GetJob(),
							hidden: marker.GetHidden(),
						}
						contexts[userID] = canonical

						// The user may have changed jobs, in which case the old
						// job's dispatcher state must still be removed. A same-job
						// deletion is normally an obsolete grade key, so retain the
						// active dispatcher state for the canonical marker.
						if current.job != "" && current.job != canonical.job {
							s.removeDispatchersForUser(ctx, userID, []string{current.job})
						}

						if err := s.syncUserUnitMapping(ctx, userID); err != nil {
							s.logger.Error(
								"failed to sync user unit mapping for stale usermarker delete event",
								zap.Int32("user_id", userID),
								zap.Error(err),
							)
						}
						return
					}

					previous, seen := contexts[userID]
					delete(contexts, userID)
					s.removeDispatchersForUser(
						ctx,
						userID,
						dispatcherJobsToRemove(previous, seen, current),
					)

					// A deleted marker means the user is no longer on duty. Reconcile
					// durable membership immediately instead of waiting for a unit ping
					// or the daily audit.
					if err := s.syncUserUnitMapping(ctx, userID); err != nil {
						s.logger.Error(
							"failed to sync user unit mapping for usermarker delete event",
							zap.Int32("user_id", userID),
							zap.Error(err),
						)
					}
				}()
			}
		}
	}
}
