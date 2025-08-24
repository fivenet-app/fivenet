package croner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/leaderelection"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	logger    *zap.Logger
	ctxCancel context.Context
	js        *events.JSWrapper
	registry  *Registry
	gron      *gronx.Gronx
	le        *leaderelection.LeaderElector

	nodeName string
	jsCons   jetstream.ConsumeContext
}

func NewScheduler(p SchedulerParams) (*Scheduler, error) {
	nodeName := instance.ID() + "_cron_scheduler"

	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.WithOptions(zap.IncreaseLevel(p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentCron, p.Cfg.LogLevel))).
		Named("cron.scheduler")
	s := &Scheduler{
		logger: logger,

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
		s.le.Stop()
		cancel()

		return nil
	}))

	return s, nil
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

				var wg sync.WaitGroup

				s.registry.store.Range(func(key string, value *cron.Cronjob) bool {
					job, err := s.registry.store.GetOrLoad(ctx, key)
					if err != nil {
						s.logger.Error("failed to load cron job", zap.String("job_name", key))
						return true
					}

					// Check if the cron job is already/still running and under the timeout check
					if job.GetStartedTime() != nil &&
						job.GetState() == cron.CronjobState_CRONJOB_STATE_RUNNING {
						if time.Since(job.GetStartedTime().AsTime()) <= job.GetRunTimeout() {
							return true
						}
					}

					ok, err := s.gron.IsDue(job.GetSchedule(), t)
					if err != nil {
						s.logger.Error(
							"failed to chek cron job due time",
							zap.String("job_name", key),
							zap.String("schedule", job.GetSchedule()),
						)
						return true
					}
					if !ok {
						return true
					}

					s.logger.Debug("scheduling cron job", zap.String("name", job.GetName()))
					wg.Add(1)
					go func() {
						defer wg.Done()

						if err := s.registry.store.ComputeUpdate(ctx, key, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
							if existing == nil {
								return existing, false, nil
							}

							existing.StartedTime = timestamp.Now()
							existing.State = cron.CronjobState_CRONJOB_STATE_RUNNING

							job.StartedTime = existing.GetStartedTime()
							job.State = existing.GetState()

							return existing, true, nil
						}); err != nil {
							s.logger.Error(
								"failed to update status of cron job",
								zap.String("job_name", job.GetName()),
							)
						}

						if err := s.runCronjob(ctx, job); err != nil {
							s.logger.Error(
								"failed to trigger cron job run",
								zap.String("job_name", job.GetName()),
							)
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
	if _, err := s.js.PublishProto(ctx, fmt.Sprintf("%s.%s", CronScheduleSubject, CronScheduleTopic), &cron.CronjobSchedulerEvent{
		Cronjob: job,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) registerSubscriptions(
	ctxStartup context.Context,
	ctxCancel context.Context,
) error {
	consumer, err := s.js.CreateOrUpdateConsumer(
		ctxStartup,
		CronScheduleStreamName,
		jetstream.ConsumerConfig{
			Durable:           instance.ID() + "_cron_scheduler",
			DeliverPolicy:     jetstream.DeliverNewPolicy,
			FilterSubject:     fmt.Sprintf("%s.%s", CronScheduleSubject, CronCompleteTopic),
			MaxDeliver:        3,
			InactiveThreshold: 1 * time.Minute, // Close consumer if inactive for 1 minute
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
	event := &cron.CronjobCompletedEvent{}
	if err := protoutils.UnmarshalPartialJSON(msg.Data(), event); err != nil {
		s.logger.Error(
			"failed to unmarshal cron completion msg",
			zap.String("subject", msg.Subject()),
			zap.Error(err),
		)

		if err := msg.NakWithDelay(150 * time.Millisecond); err != nil {
			s.logger.Error(
				"failed to nack unmarshal cron completion msg",
				zap.String("subject", msg.Subject()),
				zap.Error(err),
			)
		}
		return
	}

	if err := msg.InProgress(); err != nil {
		s.logger.Error(
			"failed to send in progress for cron completion msg",
			zap.String("subject", msg.Subject()),
			zap.String("job_name", event.GetName()),
			zap.Error(err),
		)
	}

	if err := s.registry.store.ComputeUpdate(s.ctxCancel, event.GetName(), func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
		// No need to update the job, probably doesn't exist anymore
		if existing == nil {
			return existing, false, nil
		}

		existing.State = cron.CronjobState_CRONJOB_STATE_WAITING

		nextTime, err := gronx.NextTick(existing.GetSchedule(), false)
		if err != nil {
			return existing, false, err
		}
		existing.NextScheduleTime = timestamp.New(nextTime)
		existing.LastAttemptTime = timestamp.New(time.Now())

		existing.Data = event.GetData()

		existing.LastCompletedEvent = event

		return existing, true, nil
	}); err != nil {
		s.logger.Error(
			"failed to update cronjob state after completion msg",
			zap.String("subject", msg.Subject()),
			zap.String("job_name", event.GetName()),
			zap.Error(err),
		)
		return
	}

	if err := msg.Ack(); err != nil {
		s.logger.Error(
			"failed to ack cron completion msg",
			zap.String("subject", msg.Subject()),
			zap.String("job_name", event.GetName()),
			zap.Error(err),
		)
		return
	}
}
