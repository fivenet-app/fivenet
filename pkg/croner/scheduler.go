package croner

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"uuid"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/nats/leaderelection"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var SchedulerModule = fx.Module("cron_scheduler",
	fx.Provide(
		NewScheduler,
	),
)

const (
	OwnerKey = "_owner"

	jobNameLabel                    = "job_name"
	maxCronSchedulerClaimConcurrent = 32
)

var ErrCronjobAlreadyRunning = errors.New("cron job is already running")
var ErrCronjobRunNotReady = errors.New("cron job run state has not replicated yet")

const (
	cronCompletionAckWait    = 2 * time.Minute
	cronCompletionMaxDeliver = 10
	cronCompletionMinRetry   = 250 * time.Millisecond
	cronCompletionMaxRetry   = 30 * time.Second
)

type IScheduler interface {
	RunJob(ctx context.Context, name string) (*jetstream.PubAck, error)
}

type SchedulerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
	State  *Registry
	Cfg    *config.Config
}

type SchedulerResult struct {
	fx.Out

	Scheduler  *Scheduler
	IScheduler IScheduler
}

type Scheduler struct {
	IScheduler

	logger    *zap.Logger
	ctxCancel context.Context //nolint:containedctx // Scheduler keeps a root cancelable context for store operations and subscriptions.
	js        *events.JSWrapper
	publisher events.IPublisher
	registry  *Registry
	gron      *gronx.Gronx
	le        *leaderelection.LeaderElector
	metrics   *schedulerMetrics

	nodeName string
	jsCons   jetstream.ConsumeContext
	claimSem chan struct{}
}

func NewScheduler(p SchedulerParams) (SchedulerResult, error) {
	nodeName := instance.ID() + "_cron_scheduler"

	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.WithOptions(zap.IncreaseLevel(p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentCron, p.Cfg.LogLevel))).
		Named("cron.scheduler")
	s := &Scheduler{
		logger: logger,

		nodeName: nodeName,

		ctxCancel: ctxCancel,
		js:        p.JS,
		publisher: p.JS,
		registry:  p.State,
		gron:      gronx.New(),
		metrics:   getSchedulerMetrics(),
		claimSem:  make(chan struct{}, maxCronSchedulerClaimConcurrent),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, s.js); err != nil {
			return err
		}

		var err error
		s.le, err = leaderelection.New(
			ctxCancel, s.logger, s.js,
			"leader_election", // Bucket
			"cron_scheduler",  // Key
			12*time.Second,    // TTL for the lock
			6*time.Second,     // Heartbeat interval
			func(ctx context.Context) {
				s.logger.Info("scheduler started", zap.String("node_name", s.nodeName))

				s.start(ctx)
			},
			nil, // No on stopped function, context cancels the scheduler
		)
		if err != nil {
			return fmt.Errorf("failed to create leader elector. %w", err)
		}

		s.le.Start()

		return s.registerSubscriptions(ctxStartup, ctxCancel)
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		if s.jsCons != nil {
			s.jsCons.Stop()
			s.jsCons = nil
		}
		if s.le != nil {
			s.le.Stop()
		}
		cancel()

		return nil
	}))

	return SchedulerResult{
		Scheduler:  s,
		IScheduler: s,
	}, nil
}

func (s *Scheduler) start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

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

				var wg sync.WaitGroup

				s.registry.store.Range(func(key string, value *cron.Cronjob) bool {
					job, err := s.registry.store.GetOrLoad(ctx, key)
					if err != nil {
						s.logger.Error("failed to load cron job", zap.String(jobNameLabel, key))
						return true
					}

					// Check if the cron job is already/still running and under the timeout check
					if job.GetStartedTime() != nil &&
						job.GetState() == cron.CronjobState_CRONJOB_STATE_RUNNING {
						if time.Since(job.GetStartedTime().AsTime()) <= job.GetRunTimeout() {
							return true
						}
					}

					ok := false
					if nextSchedule := job.GetNextScheduleTime(); nextSchedule != nil {
						ok = !nextSchedule.AsTime().After(t)
					} else {
						var err error
						ok, err = s.gron.IsDue(job.GetSchedule(), t)
						if err != nil {
							s.logger.Error(
								"failed to check cron job due time",
								zap.String(jobNameLabel, key),
								zap.String("schedule", job.GetSchedule()),
							)
							return true
						}
					}
					if !ok {
						return true
					}

					select {
					case s.claimSem <- struct{}{}:
					case <-ctx.Done():
						return false
					}

					s.logger.Debug("scheduling cron job", zap.String("name", job.GetName()))
					wg.Go(func() {
						defer func() { <-s.claimSem }()
						claimed, err := s.claimJob(ctx, key)
						if err != nil {
							s.logger.Error(
								"failed to claim cron job",
								zap.String(jobNameLabel, job.GetName()),
								zap.Error(err),
							)
							return
						}
						if claimed == nil {
							return
						}

						if _, err := s.runCronjob(ctx, claimed); err != nil {
							s.logger.Error(
								"failed to trigger cron job run",
								zap.String(jobNameLabel, claimed.GetName()),
							)
							if err := s.releaseJob(ctx, claimed.GetName(), claimed.GetRunId(), true); err != nil {
								s.logger.Error(
									"failed to release cron job after publish failure",
									zap.String(jobNameLabel, claimed.GetName()),
									zap.Error(err),
								)
							}
						}
					})

					return true
				})

				wg.Wait()
			}()
		}
	}
}

func (s *Scheduler) runCronjob(ctx context.Context, job *cron.Cronjob) (*jetstream.PubAck, error) {
	startedAt := time.Now()
	pa, err := s.publisher.PublishProto(
		ctx,
		fmt.Sprintf("%s.%s", CronScheduleSubject, CronScheduleTopic),
		&cron.CronjobSchedulerEvent{
			Cronjob: job,
		},
	)
	s.metrics.handoffLatency.WithLabelValues(job.GetName()).Observe(time.Since(startedAt).Seconds())
	if err != nil {
		return nil, err
	}

	return pa, nil
}

func (s *Scheduler) RunJob(ctx context.Context, name string) (*jetstream.PubAck, error) {
	name = normalizeCronjobName(name)
	if _, err := s.registry.GetCronjob(ctx, name); err != nil {
		return nil, fmt.Errorf("failed to get cron job. %w", err)
	}

	job, err := s.claimJob(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to claim cron job. %w", err)
	}
	if job == nil {
		return nil, ErrCronjobAlreadyRunning
	}

	pa, err := s.runCronjob(ctx, job)
	if err != nil {
		releaseCtx, releaseCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer releaseCancel()
		if releaseErr := s.releaseJob(releaseCtx, job.GetName(), job.GetRunId(), false); releaseErr != nil {
			return nil, fmt.Errorf("failed to publish cron job and release claim. publish: %w; release: %v", err, releaseErr)
		}
		return nil, err
	}

	return pa, nil
}

func (s *Scheduler) claimJob(ctx context.Context, name string) (*cron.Cronjob, error) {
	name = normalizeCronjobName(name)
	runID := uuid.New().String()
	var claimed *cron.Cronjob

	if err := s.registry.store.ComputeUpdate(
		ctx,
		name,
		func(_ string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
			if existing == nil {
				return existing, false, nil
			}

			if existing.GetState() == cron.CronjobState_CRONJOB_STATE_RUNNING &&
				existing.GetStartedTime() != nil &&
				time.Since(existing.GetStartedTime().AsTime()) <= existing.GetRunTimeout() {
				return existing, false, nil
			}

			existing.SetStartedTime(timestamp.Now())
			existing.SetLastAttemptTime(timestamp.Now())
			existing.SetState(cron.CronjobState_CRONJOB_STATE_RUNNING)
			existing.SetRunId(runID)
			claimed = existing

			return existing, true, nil
		},
	); err != nil {
		return nil, err
	}

	return claimed, nil
}

func (s *Scheduler) releaseJob(ctx context.Context, name, runID string, retry bool) error {
	name = normalizeCronjobName(name)
	return s.registry.store.ComputeUpdate(
		ctx,
		name,
		func(_ string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
			if existing == nil || existing.GetRunId() != runID {
				return existing, false, nil
			}

			existing.SetState(cron.CronjobState_CRONJOB_STATE_WAITING)
			existing.ClearStartedTime()
			existing.SetRunId("")
			if retry {
				existing.SetNextScheduleTime(timestamp.Now())
			}

			return existing, true, nil
		},
	)
}

func (s *Scheduler) registerSubscriptions(
	ctxStartup context.Context,
	ctxCancel context.Context,
) error {
	consumer, err := s.js.CreateOrUpdateConsumer(
		ctxStartup,
		CronScheduleStreamName,
		jetstream.ConsumerConfig{
			Durable:       CronSchedulerConsumerName,
			DeliverPolicy: jetstream.DeliverAllPolicy,
			FilterSubject: fmt.Sprintf("%s.%s", CronScheduleSubject, CronCompleteTopic),
			AckPolicy:     jetstream.AckExplicitPolicy,
			AckWait:       cronCompletionAckWait,
			MaxDeliver:    cronCompletionMaxDeliver,
		},
	)
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
	heartbeatDone := make(chan struct{})
	go s.heartbeat(msg, heartbeatDone)
	defer close(heartbeatDone)

	event := &cron.CronjobCompletedEvent{}
	if err := protoutils.UnmarshalPartialJSON(msg.Data(), event); err != nil {
		s.logger.Error(
			"failed to unmarshal cron completion msg",
			zap.String("subject", msg.Subject()),
			zap.Error(err),
		)

		if err := msg.Term(); err != nil {
			s.logger.Error(
				"failed to terminate invalid cron completion msg",
				zap.String("subject", msg.Subject()),
				zap.Error(err),
			)
		}
		return
	}
	if event.GetName() == "" || event.GetRunId() == "" {
		if err := msg.Term(); err != nil {
			s.logger.Error(
				"failed to terminate incomplete cron completion msg",
				zap.String("subject", msg.Subject()),
				zap.Error(err),
			)
		}
		return
	}
	event.Name = normalizeCronjobName(event.GetName())

	if err := msg.InProgress(); err != nil {
		s.logger.Error(
			"failed to send in progress for cron completion msg",
			zap.String("subject", msg.Subject()),
			zap.String(jobNameLabel, event.GetName()),
			zap.Error(err),
		)
	}

	if err := s.registry.store.ComputeUpdate(
		s.ctxCancel,
		event.GetName(),
		func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
			// No need to update the job, probably doesn't exist anymore
			if existing == nil {
				return existing, false, nil
			}
			if existing.GetRunId() == "" || existing.GetRunId() != event.GetRunId() {
				if existing.GetRunId() != "" {
					if event.GetStartedTime() != nil && existing.GetStartedTime() != nil &&
						event.GetStartedTime().AsTime().Before(existing.GetStartedTime().AsTime()) {
						return existing, false, nil
					}
					return existing, false, ErrCronjobRunNotReady
				}
				return existing, false, nil
			}

			existing.State = cron.CronjobState_CRONJOB_STATE_WAITING

			nextTime, err := gronx.NextTick(existing.GetSchedule(), false)
			if err != nil {
				return existing, false, err
			}
			existing.NextScheduleTime = timestamp.New(nextTime)
			existing.ClearStartedTime()

			existing.Data = event.GetData()

			existing.LastCompletedEvent = event
			existing.SetRunId("")

			return existing, true, nil
		},
	); err != nil {
		s.logger.Error(
			"failed to update cronjob state after completion msg",
			zap.String("subject", msg.Subject()),
			zap.String(jobNameLabel, event.GetName()),
			zap.Error(err),
		)
		delay := cronCompletionMinRetry
		if errors.Is(err, ErrCronjobRunNotReady) {
			delay = cronCompletionRetryDelay(msg)
		}
		if nakErr := msg.NakWithDelay(delay); nakErr != nil {
			s.logger.Error("failed to nack cron completion msg", zap.Error(nakErr))
		}
		return
	}

	if err := msg.Ack(); err != nil {
		s.logger.Error(
			"failed to ack cron completion msg",
			zap.String("subject", msg.Subject()),
			zap.String(jobNameLabel, event.GetName()),
			zap.Error(err),
		)
		return
	}
}

func cronCompletionRetryDelay(msg jetstream.Msg) time.Duration {
	delay := cronCompletionMinRetry
	metadata, err := msg.Metadata()
	if err != nil || metadata == nil {
		return delay
	}

	for attempt := uint64(1); attempt < metadata.NumDelivered && delay < cronCompletionMaxRetry; attempt++ {
		delay *= 2
	}
	if delay > cronCompletionMaxRetry {
		return cronCompletionMaxRetry
	}

	return delay
}

func (s *Scheduler) heartbeat(msg jetstream.Msg, done <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := msg.InProgress(); err != nil {
				s.logger.Warn("failed to extend cron completion message", zap.Error(err))
			}
		}
	}
}
