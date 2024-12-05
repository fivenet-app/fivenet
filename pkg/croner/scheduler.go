package croner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/puzpuzpuz/xsync/v3"
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
}

type jobWrapper struct {
	ctx    context.Context
	cancel context.CancelFunc

	schedule string
}

type Scheduler struct {
	logger *zap.Logger

	ctx   context.Context
	js    *events.JSWrapper
	store *store.Store[cron.Cronjob, *cron.Cronjob]
	gron  *gronx.Gronx

	jsCons jetstream.ConsumeContext

	jobs *xsync.MapOf[string, *jobWrapper]
}

func NewScheduler(p SchedulerParams) (*Scheduler, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Scheduler{
		logger: p.Logger.Named("cron_scheduler"),
		ctx:    ctxCancel,
		js:     p.JS,
		gron:   gronx.New(),

		jobs: xsync.NewMapOf[string, *jobWrapper](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxStartup, s.js); err != nil {
			return err
		}

		st, err := store.New(ctxStartup, p.Logger, p.JS, "cron",
			store.WithOnUpdateFn(func(_ *store.Store[cron.Cronjob, *cron.Cronjob], cj *cron.Cronjob) (*cron.Cronjob, error) {
				if cj == nil {
					return cj, nil
				}

				jw, ok := s.jobs.Load(cj.Name)
				if !ok {
					ctx, cancel := context.WithCancel(ctxCancel)
					s.jobs.Store(cj.Name, &jobWrapper{
						ctx:    ctx,
						cancel: cancel,

						schedule: cj.Schedule,
					})
				} else {
					if cj.Schedule != jw.schedule {
						jw.schedule = cj.Schedule
					}
				}

				return cj, nil
			}),

			store.WithOnDeleteFn(func(_ *store.Store[cron.Cronjob, *cron.Cronjob], entry jetstream.KeyValueEntry, cj *cron.Cronjob) error {
				jw, ok := s.jobs.LoadAndDelete(entry.Key())
				if !ok {
					return nil
				}

				jw.cancel()
				jw = nil

				return nil
			}),
		)
		if err != nil {
			return err
		}
		s.store = st

		if err := st.Start(ctxCancel); err != nil {
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
			if s.jobs.Size() == 0 {
				continue
			}

			func() {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				wg := sync.WaitGroup{}

				s.jobs.Range(func(key string, value *jobWrapper) bool {
					job, err := s.store.GetOrLoad(ctx, key)
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

					wg.Add(1)
					go func() {
						defer wg.Done()

						if err := s.runCronjob(ctx, job); err != nil {
							s.logger.Error("failed to trigger cron job run", zap.String("job_name", job.Name))
						}

						if err := s.store.ComputeUpdate(ctx, key, true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
							if existing == nil {
								return existing, false, nil
							}

							existing.State = cron.CronjobState_CRONJOB_STATE_RUNNING

							return existing, true, nil
						}); err != nil {
							s.logger.Error("failed to update status of cron job", zap.String("job_name", job.Name))
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
	job := &cron.CronjobCompletedEvent{}
	if err := protojson.Unmarshal(msg.Data(), job); err != nil {
		s.logger.Error("failed to unmarshal cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))

		if err := msg.NakWithDelay(150 * time.Millisecond); err != nil {
			s.logger.Error("failed to nack unmarshal cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	if err := msg.InProgress(); err != nil {
		s.logger.Error("failed to send in progress for cron completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
	}

	if err := s.store.ComputeUpdate(s.ctx, job.Name, true, func(key string, existing *cron.Cronjob) (*cron.Cronjob, bool, error) {
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

		existing.Data.Merge(job.Data)

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
