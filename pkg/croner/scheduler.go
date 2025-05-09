package croner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
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

type SchedulerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
	State  *State
}

type Scheduler struct {
	logger *zap.Logger

	ctx   context.Context
	js    *events.JSWrapper
	state *State
	gron  *gronx.Gronx

	jsCons jetstream.ConsumeContext
}

func NewScheduler(p SchedulerParams) (*Scheduler, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Scheduler{
		logger: p.Logger.Named("cron_scheduler"),
		ctx:    ctxCancel,
		js:     p.JS,
		state:  p.State,
		gron:   gronx.New(),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, s.js); err != nil {
			return err
		}

		return s.registerSubscriptions(ctxStartup, ctxCancel)
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
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
			if s.state.cronjobs.Size() == 0 {
				continue
			}

			func() {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				wg := sync.WaitGroup{}

				s.state.cronjobs.Range(func(key string, value *jobWrapper) bool {
					job, err := s.state.store.GetOrLoad(ctx, key)
					if err != nil {
						s.logger.Error("failed to load cron job", zap.String("job_name", key))
						return true
					}

					// Check if the cron job is already/still running and under the timeout check
					if job.StartedTime != nil && (job.State == cron.CronjobState_CRONJOB_STATE_RUNNING &&
						time.Since(job.StartedTime.AsTime()) <= job.GetRunTimeout()) {
						return true
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

						if err := s.state.store.ComputeUpdate(ctx, key, true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
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

	if err := s.state.store.ComputeUpdate(s.ctx, event.Name, true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
		// No need to update the job, probably doesn't exist anymore
		if existing == nil {
			return existing, false, nil
		}

		existing.State = cron.CronjobState_CRONJOB_STATE_PENDING

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
