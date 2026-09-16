package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pkguserinfo "github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

const (
	userInfoReconcileConsumerName = "centrum_userinfo_reconcile"
	userInfoReconcileAckWait      = 30 * time.Second
	userInfoReconcileRetryDelay   = 5 * time.Second
	userInfoReconcileSlowAfter    = 10
	userInfoReconcileSlowRetry    = 5 * time.Minute
)

type userInfoReconcileConsume struct {
	ctx  jetstream.ConsumeContext
	errs <-chan error
}

// ensureUserInfoReconcileConsumer keeps user-info changes retained while no
// housekeeper leader is active. It intentionally has no inactivity timeout.
func (s *Housekeeper) ensureUserInfoReconcileConsumer(ctx context.Context) error {
	consumer, err := pkguserinfo.CreateOrUpdateUserInfoChangeConsumer(
		ctx,
		s.js,
		jetstream.ConsumerConfig{
			Durable:       userInfoReconcileConsumerName,
			DeliverPolicy: jetstream.DeliverNewPolicy,
			AckPolicy:     jetstream.AckExplicitPolicy,
			AckWait:       userInfoReconcileAckWait,
		},
	)
	if err != nil {
		return err
	}
	s.userInfoReconcileConsumer = consumer

	return nil
}

// startUserInfoReconcileConsumer starts durable delivery only while this
// process is leader. Acknowledging after correction preserves events across
// leader handoffs and retries transient persistence failures.
func (s *Housekeeper) startUserInfoReconcileConsumer(
	ctx context.Context,
) (*userInfoReconcileConsume, error) {
	consumeErr := make(chan error, 1)
	consumeCtx, err := s.userInfoReconcileConsumer.Consume(
		func(msg jetstream.Msg) {
			s.handleUserInfoReconcileMessage(ctx, msg)
		},
		jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
			if err != nil && ctx.Err() == nil {
				select {
				case consumeErr <- err:
				default:
				}
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	return &userInfoReconcileConsume{ctx: consumeCtx, errs: consumeErr}, nil
}

func (s *Housekeeper) maintainUserInfoReconcileConsumer(
	ctx context.Context,
	onStarted func(context.Context),
) {
	var consume *userInfoReconcileConsume
	started := false
	for {
		if consume == nil {
			next, err := s.startUserInfoReconcileConsumer(ctx)
			if err != nil {
				s.metrics.IncHousekeeperEvent("userinfo_reconcile", "consumer_start_failed")
				s.logger.Error("failed to start user info reconciliation consumer", zap.Error(err))
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
				}
				continue
			}
			consume = next
			if !started {
				started = true
				s.wg.Go(func() { onStarted(ctx) })
			} else {
				s.metrics.IncHousekeeperEvent("userinfo_reconcile", "consumer_restarted")
			}
		}

		select {
		case <-ctx.Done():
			consume.ctx.Stop()
			return
		case err := <-consume.errs:
			consume.ctx.Stop()
			consume = nil
			s.metrics.IncHousekeeperEvent("userinfo_reconcile", "consumer_failed")
			s.logger.Error("user info reconciliation consumer failed", zap.Error(err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}
}

func (s *Housekeeper) handleUserInfoReconcileMessage(ctx context.Context, msg jetstream.Msg) {
	event := &pbuserinfo.UserInfoChanged{}
	if err := protoutils.UnmarshalPartialJSON(msg.Data(), event); err != nil {
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", "decode_failed")
		s.logger.Error("failed to decode user info reconciliation event", zap.Error(err))
		if termErr := msg.Term(); termErr != nil {
			s.logger.Error("failed to terminate malformed user info event", zap.Error(termErr))
		}
		return
	}
	if event.GetUserId() <= 0 || event.GetNewJob() == "" {
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", "invalid_event")
		s.logger.Error("received invalid user info reconciliation event",
			zap.Int32("user_id", event.GetUserId()),
			zap.String("new_job", event.GetNewJob()),
		)
		if err := msg.Term(); err != nil {
			s.logger.Error("failed to terminate invalid user info event", zap.Error(err))
		}
		return
	}
	if ctx.Err() != nil {
		return
	}
	info, err := s.userinfo.GetUserInfo(ctx, event.GetUserId())
	if err != nil {
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", "userinfo_lookup_failed")
		s.logger.Error("failed to load current user info for reconciliation",
			zap.Int32("user_id", event.GetUserId()),
			zap.Error(err),
		)
		s.retryUserInfoReconcileMessage(msg, "userinfo_lookup")
		return
	}
	job := info.GetJob()
	if job == "" {
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", "invalid_authoritative_job")
		s.logger.Error(
			"current user info has no primary job",
			zap.Int32("user_id", event.GetUserId()),
		)
		if err := msg.Term(); err != nil {
			s.logger.Error("failed to terminate invalid user info event", zap.Error(err))
		}
		return
	}

	if err := s.reconcileUserJobChange(ctx, event.GetUserId(), job); err != nil {
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", "reconcile_failed")
		s.logger.Error("failed to reconcile user info event",
			zap.Int32("user_id", event.GetUserId()),
			zap.String("authoritative_job", job),
			zap.Error(err),
		)
		s.retryUserInfoReconcileMessage(msg, "reconcile")
		return
	}
	if ctx.Err() != nil {
		return
	}
	if err := msg.Ack(); err != nil {
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", "ack_failed")
		s.logger.Error("failed to acknowledge user info reconciliation event", zap.Error(err))
		return
	}
	s.metrics.IncHousekeeperEvent("userinfo_reconcile", "reconciled")
}

// retryUserInfoReconcileMessage backs off repeatedly failing messages without
// terminating them: account lookup errors cannot currently distinguish a
// transient database failure from an account that has become unavailable.
func (s *Housekeeper) retryUserInfoReconcileMessage(msg jetstream.Msg, operation string) {
	retryDelay := userInfoReconcileRetryDelay
	metadata, metadataErr := msg.Metadata()
	if metadataErr == nil && metadata != nil &&
		metadata.NumDelivered >= userInfoReconcileSlowAfter {
		retryDelay = userInfoReconcileSlowRetry
		s.metrics.IncHousekeeperEvent("userinfo_reconcile", operation+"_slow_retry")
		s.logger.Warn("slowing retries for repeatedly failing user info reconciliation event",
			zap.String("operation", operation),
			zap.Uint64("deliveries", metadata.NumDelivered),
		)
	}
	if err := msg.NakWithDelay(retryDelay); err != nil {
		s.logger.Error("failed to retry user info reconciliation event", zap.Error(err))
	}
}

// reconcileUserJobChange applies only authoritative job-mismatch corrections.
// Tracker-driven SyncUserUnitMapping remains responsible for duty state and
// repairing matching assignments.
func (s *Housekeeper) reconcileUserJobChange(ctx context.Context, userID int32, job string) error {
	return s.reconcileUserJobChangeWithDispatcherJobs(ctx, userID, job, nil)
}

// reconcileUserJobChangeWithDispatcherJobs applies job-mismatch corrections.
// A nil dispatcher job set discovers projections for an incremental event; a
// supplied set comes from the full sweep's existing dispatcher scan.
func (s *Housekeeper) reconcileUserJobChangeWithDispatcherJobs(
	ctx context.Context,
	userID int32,
	job string,
	dispatcherJobs map[string]struct{},
) error {
	removedUnit, err := s.unitUserState.ReconcileUserJobChange(ctx, userID, job)
	var errs error
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("failed to reconcile unit assignment: %w", err))
	} else if removedUnit {
		s.metrics.IncHousekeeperEvent("userinfo_changes", "unit_assignment_removed")
	}

	if dispatcherJobs == nil {
		dispatcherJobs = map[string]struct{}{}
		s.dispatcherUserState.Range(
			func(dispatcherJob string, value *centrumdispatchers.Dispatchers) bool {
				for _, dispatcher := range value.GetDispatchers() {
					if dispatcher.GetUserId() == userID {
						dispatcherJobs[dispatcherJob] = struct{}{}
						break
					}
				}
				return true
			},
		)
	}

	for dispatcherJob := range dispatcherJobs {
		if dispatcherJob == job {
			continue
		}
		if err := s.dispatcherUserState.SetUserState(
			ctx,
			dispatcherJob,
			userID,
			false,
		); err != nil {
			errs = errors.Join(errs, fmt.Errorf(
				"failed to remove user from %s dispatchers: %w",
				dispatcherJob,
				err,
			))
			continue
		}
		s.metrics.IncHousekeeperEvent("userinfo_changes", "dispatcher_removed")
	}

	return errs
}

// reconcileUserInfoState repairs current Centrum state from authoritative
// primary jobs. It is used after leadership acquisition and by the manual
// recovery cron; durable events handle the normal incremental path.
func (s *Housekeeper) reconcileUserInfoState(ctx context.Context) (int, int, error) {
	userIDs := map[int32]struct{}{}
	dispatcherJobsByUser := map[int32]map[string]struct{}{}

	s.unitUserState.Range(func(_ string, unit *centrumunits.Unit) bool {
		if unit == nil {
			return true
		}
		for _, assignment := range unit.GetUsers() {
			if assignment.GetUserId() > 0 {
				userIDs[assignment.GetUserId()] = struct{}{}
			}
		}
		return true
	})
	s.dispatcherUserState.Range(func(job string, value *centrumdispatchers.Dispatchers) bool {
		for _, dispatcher := range value.GetDispatchers() {
			if dispatcher.GetUserId() > 0 {
				userID := dispatcher.GetUserId()
				userIDs[userID] = struct{}{}
				if dispatcherJobsByUser[userID] == nil {
					dispatcherJobsByUser[userID] = map[string]struct{}{}
				}
				dispatcherJobsByUser[userID][job] = struct{}{}
			}
		}
		return true
	})

	processed := 0
	var errs error
	for userID := range userIDs {
		if ctx.Err() != nil {
			return processed, len(userIDs), errors.Join(errs, ctx.Err())
		}

		info, err := s.userinfo.GetUserInfo(ctx, userID)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to retrieve user %d: %w", userID, err))
			continue
		}
		processed++
		if err := s.reconcileUserJobChangeWithDispatcherJobs(
			ctx,
			userID,
			info.GetJob(),
			dispatcherJobsByUser[userID],
		); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to reconcile user %d: %w", userID, err))
		}
	}

	return processed, len(userIDs), errs
}

func (s *Housekeeper) runReconcileUserInfoState(ctx context.Context, _ *cron.CronjobData) error {
	startedAt := time.Now()
	defer func() {
		s.metrics.ObserveHousekeeperDuration(
			"reconcile_userinfo_state",
			time.Since(startedAt).Seconds(),
		)
	}()

	processed, discovered, err := s.reconcileUserInfoState(ctx)
	s.metrics.SetHousekeeperWork("reconcile_userinfo_state", "users_processed", processed)
	s.metrics.SetHousekeeperWork("reconcile_userinfo_state", "users_discovered", discovered)

	return err
}
