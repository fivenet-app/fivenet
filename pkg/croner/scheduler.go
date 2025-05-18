package croner

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/locks"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

var SchedulerModule = fx.Module("cron_scheduler",
	fx.Provide(
		NewScheduler,
	),
)

const OwnerKey = "_owner"

type SchedulerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
	State  *Registry
	Cfg    *config.Config
}

type Scheduler struct {
	logger *zap.Logger

	nodeName string

	ctxCancel context.Context
	js        *events.JSWrapper
	registry  *Registry
	gron      *gronx.Gronx
	lock      *locks.Locks

	jsCons jetstream.ConsumeContext
}

func NewScheduler(p SchedulerParams) (*Scheduler, error) {
	nodeName, err := getNodeName(p.Cfg.HTTP.AdminListen)
	if err != nil {
		return nil, fmt.Errorf("failed to get node name. %w", err)
	}

	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Scheduler{
		logger: p.Logger.Named("cron.scheduler"),

		nodeName: nodeName,

		ctxCancel: ctxCancel,
		js:        p.JS,
		registry:  p.State,
		gron:      gronx.New(),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, s.js); err != nil {
			return err
		}

		lock, err := locks.New(s.logger, s.registry.kv, s.registry.kv.Bucket(), 10*time.Second)
		if err != nil {
			return err
		}
		s.lock = lock

		go s.run(ctxCancel)

		return s.registerSubscriptions(ctxStartup, ctxCancel)
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return nil
	}))

	return s, nil
}

func (s *Scheduler) run(ctx context.Context) {
	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			s.unlock(ctx)
			return

		case <-time.After(2 * time.Second):
		}

		state, err := s.getLockState(ctx)
		if err != nil {
			s.logger.Error("failed to create lock", zap.Error(err))
			return
		}

		// Check if we are the owner of the lock or it is expired
		if state != nil && state.Hostname != s.nodeName {
			if state.UpdatedAt != nil {
				// Check if the lock is expired
				if time.Since(state.UpdatedAt.AsTime()) < 10*time.Second {
					s.logger.Debug("lock is still valid, sleeping for a few seconds")
					continue
				}
			}
		}

		s.logger.Info("no or expired lock, attempting to get lock")

		if err := s.createOrUpdateLock(ctx); err != nil {
			s.logger.Error("failed to create/update lock", zap.Error(err))
			continue
		}

		func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			// Keep lock owner state uptodate
			wg.Add(1)
			go func() {
				defer wg.Done()

				s.start(ctx)
			}()

			for {
				select {
				case <-ctx.Done():
					wg.Wait()
					return

				case <-time.After(5 * time.Second):
				}

				if err := s.createOrUpdateLock(ctx); err != nil {
					cancel()
					s.logger.Error("error creating/updating owner lock state, stopping scheduler", zap.Error(err))
					return
				}
			}
		}()
	}
}

func (s *Scheduler) getLockState(ctx context.Context) (*cron.CronjobLockOwnerState, error) {
	if err := s.lock.Lock(ctx, OwnerKey); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, fmt.Errorf("error while trying to get owner lock. %w", err)
	}
	defer func() {
		if err := s.lock.Unlock(ctx, OwnerKey); err != nil {
			s.logger.Error("failed to unlock owner lock", zap.Error(err))
		}
	}()

	entry, err := s.registry.kv.Get(ctx, OwnerKey)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, fmt.Errorf("error getting owner lock entry. %w", err)
	}

	if entry != nil {
		// Make sure the owner state is really expired and not just a "hiccup" of our nats lock logic
		state := &cron.CronjobLockOwnerState{}
		if err := protojson.Unmarshal(entry.Value(), state); err != nil {
			s.logger.Warn("failed to unmarshal owner lock entry", zap.Error(err))
		}

		return state, nil
	}

	return nil, nil
}

func (s *Scheduler) createOrUpdateLock(ctx context.Context) error {
	if err := s.lock.Lock(ctx, OwnerKey); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		return fmt.Errorf("error while trying to get owner lock. %w", err)
	}
	defer func() {
		if err := s.lock.Unlock(ctx, OwnerKey); err != nil {
			s.logger.Error("failed to unlock owner lock", zap.Error(err))
		}
	}()

	out, err := protoutils.Marshal(&cron.CronjobLockOwnerState{
		Hostname:  s.nodeName,
		UpdatedAt: timestamp.Now(),
	})
	if err != nil {
		s.logger.Error("error marshalling owner lock state", zap.Error(err))
		return err
	}

	if _, err := s.registry.kv.Put(ctx, OwnerKey, out); err != nil {
		s.logger.Error("failed to update owner lock state in kv", zap.Error(err))
		return err
	}

	return nil
}

func (s *Scheduler) unlock(ctx context.Context) {
	entry, err := s.registry.kv.Get(ctx, OwnerKey)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		s.logger.Error("error getting owner lock entry", zap.Error(err))
		return
	}

	// Make sure we are the owner of the lock
	if entry == nil || string(entry.Value()) != s.nodeName {
		return // We are not the owner, so we don't need to unlock
	}

	if err := s.lock.Unlock(ctx, OwnerKey); err != nil {
		s.logger.Error("failed to unlock owner lock on shutdown", zap.Error(err))
	}
}

func (s *Scheduler) start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return

		case t := <-ticker.C:
			if s.registry.store.Count() == 0 {
				continue
			}

			func() {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				wg := sync.WaitGroup{}

				s.registry.store.Range(ctx, func(key string, value *cron.Cronjob) bool {
					job, err := s.registry.store.GetOrLoad(ctx, key)
					if err != nil {
						s.logger.Error("failed to load cron job", zap.String("job_name", key))
						return true
					}

					// Check if the cron job is already/still running and under the timeout check
					if job.StartedTime != nil && job.State == cron.CronjobState_CRONJOB_STATE_RUNNING {
						if time.Since(job.StartedTime.AsTime()) <= job.GetRunTimeout() {
							return true
						}
					}

					ok, err := s.gron.IsDue(job.Schedule, t)
					if err != nil {
						s.logger.Error("failed to chek cron job due time", zap.String("job_name", key), zap.String("schedule", job.Schedule))
						return true
					}
					if !ok {
						return true
					}

					s.logger.Debug("scheduling cron job", zap.String("name", job.Name))
					wg.Add(1)
					go func() {
						defer wg.Done()

						if err := s.registry.store.ComputeUpdate(ctx, key, true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
							if existing == nil {
								return existing, false, nil
							}

							existing.StartedTime = timestamp.Now()
							existing.State = cron.CronjobState_CRONJOB_STATE_RUNNING

							job.StartedTime = existing.StartedTime
							job.State = existing.State

							return existing, true, nil
						}); err != nil {
							s.logger.Error("failed to update status of cron job", zap.String("job_name", job.Name))
						}

						if err := s.runCronjob(ctx, job); err != nil {
							s.logger.Error("failed to trigger cron job run", zap.String("job_name", job.Name))
						}
					}()

					return true
				})

				wg.Wait()
			}()
		}
	}
}

func (s *Scheduler) runCronjob(ctx context.Context, job *cron.Cronjob) error {
	msg := &cron.CronjobSchedulerEvent{
		Cronjob: job,
	}

	if _, err := s.js.PublishProto(ctx, fmt.Sprintf("%s.%s", CronScheduleSubject, CronScheduleTopic), msg); err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	consumer, err := s.js.CreateConsumer(ctxStartup, CronScheduleStreamName, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.%s", CronScheduleSubject, CronCompleteTopic),
		MaxDeliver:    3,
	})
	if err != nil {
		return err
	}

	if s.jsCons != nil {
		s.jsCons.Stop()
		s.jsCons = nil
	}

	s.jsCons, err = consumer.Consume(s.watchForCompletions,
		s.js.ConsumeErrHandlerWithRestart(ctxCancel, s.logger,
			s.registerSubscriptions,
		))
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) watchForCompletions(msg jetstream.Msg) {
	event := &cron.CronjobCompletedEvent{}
	if err := protojson.Unmarshal(msg.Data(), event); err != nil {
		s.logger.Error("failed to unmarshal cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))

		if err := msg.NakWithDelay(150 * time.Millisecond); err != nil {
			s.logger.Error("failed to nack unmarshal cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	if err := msg.InProgress(); err != nil {
		s.logger.Error("failed to send in progress for cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
	}

	if err := s.registry.store.ComputeUpdate(s.ctxCancel, event.Name, true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
		// No need to update the job, probably doesn't exist anymore
		if existing == nil {
			return existing, false, nil
		}

		existing.State = cron.CronjobState_CRONJOB_STATE_WAITING

		nextTime, err := gronx.NextTick(existing.Schedule, false)
		if err != nil {
			return existing, false, err
		}
		existing.NextScheduleTime = timestamp.New(nextTime)
		existing.LastAttemptTime = timestamp.New(time.Now())

		existing.Data = event.Data

		existing.LastCompletedEvent = event

		return existing, true, nil
	}); err != nil {
		s.logger.Error("failed to update cronjob state after completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
		return
	}

	if err := msg.Ack(); err != nil {
		s.logger.Error("failed to ack cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
		return
	}
}
