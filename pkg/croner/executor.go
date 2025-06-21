package croner

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var ExecutorModule = fx.Module("executor",
	fx.Provide(
		NewExecutor,
	),
)

type ExecutorParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
	Cfg    *config.Config

	Handlers *Handlers
}

// Executor is responsible for executing cron jobs.
// Previously, it was called Agent.
type Executor struct {
	logger *zap.Logger

	nodeName string

	ctx context.Context
	js  *events.JSWrapper

	jsCons jetstream.ConsumeContext

	handlers *Handlers
}

func NewExecutor(p ExecutorParams) (*Executor, error) {
	nodeName, err := getNodeName(p.Cfg.HTTP.AdminListen)
	if err != nil {
		return nil, fmt.Errorf("failed to get node name. %w", err)
	}

	ctxCancel, cancel := context.WithCancel(context.Background())

	ag := &Executor{
		logger: p.Logger.Named("cron.executor"),

		nodeName: nodeName,

		ctx: ctxCancel,
		js:  p.JS,

		handlers: p.Handlers,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxCancel, ag.js); err != nil {
			return err
		}

		return ag.registerSubscriptions(ctxStartup, ctxCancel)
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return ag, nil
}

func (ag *Executor) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	consumer, err := ag.js.CreateOrUpdateConsumer(ctxStartup, CronScheduleStreamName, jetstream.ConsumerConfig{
		Durable:           instance.ID() + "_cron_executor",
		DeliverPolicy:     jetstream.DeliverNewPolicy,
		FilterSubject:     fmt.Sprintf("%s.%s", CronScheduleSubject, CronScheduleTopic),
		MaxDeliver:        3,
		InactiveThreshold: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	if ag.jsCons != nil {
		ag.jsCons.Stop()
		ag.jsCons = nil
	}

	if ag.jsCons, err = consumer.Consume(ag.watchForEvents,
		ag.js.ConsumeErrHandlerWithRestart(ctxCancel, ag.logger,
			ag.registerSubscriptions,
		)); err != nil {
		return err
	}

	return nil
}

func (ag *Executor) watchForEvents(msg jetstream.Msg) {
	job := &cron.CronjobSchedulerEvent{}
	if err := protojson.Unmarshal(msg.Data(), job); err != nil {
		ag.logger.Error("failed to unmarshal cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))

		if err := msg.NakWithDelay(100 * time.Millisecond); err != nil {
			ag.logger.Error("failed to nack unmarshal cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	fn := ag.handlers.getCronjobHandler(job.Cronjob.Name)
	if fn == nil {
		if err := msg.NakWithDelay(100 * time.Millisecond); err != nil {
			ag.logger.Error("failed to nack unmarshal cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	if err := msg.Ack(); err != nil {
		ag.logger.Error("failed to send in progress for cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
	}

	if job.Cronjob.Data == nil {
		job.Cronjob.Data = &cron.CronjobData{
			Data: &anypb.Any{},
		}
	}

	var timeout time.Duration
	if job.Cronjob.Timeout != nil {
		timeout = job.Cronjob.Timeout.AsDuration()
	} else {
		timeout = DefaultCronjobTimeout
	}

	ag.logger.Debug("running cron job", zap.String("name", job.Cronjob.Name), zap.Duration("timeout", timeout))

	var elapsed time.Duration

	var err error
	func() {
		ctx, cancel := context.WithTimeout(ag.ctx, timeout)
		defer cancel()

		start := time.Now()
		defer func() {
			elapsed = time.Since(start)

			// Recover from a panic and set err accordingly
			if e := recover(); e != nil {
				if er, ok := e.(error); ok {
					err = fmt.Errorf("recovered from panic. %w", er)
				} else {
					err = fmt.Errorf("recovered from panic. %v", er)
				}

				ag.logger.Error("cron job panic", zap.String("name", job.Cronjob.Name), zap.Error(err))
			}
		}()
		err = fn(ctx, job.Cronjob.Data)
	}()

	// Update timestamp in cronjob data
	now := timestamp.Now()
	job.Cronjob.Data.UpdatedAt = now
	var errMsg *string
	if err != nil {
		msg := err.Error()
		errMsg = &msg
	}

	if _, err := ag.js.PublishProto(ag.ctx, fmt.Sprintf("%s.%s", CronScheduleSubject, CronCompleteTopic), &cron.CronjobCompletedEvent{
		Name:      job.Cronjob.Name,
		Success:   err == nil,
		Cancelled: err != nil && err == context.Canceled,
		Elapsed:   durationpb.New(elapsed),
		EndDate:   now,

		NodeName: ag.nodeName,
		Data:     job.Cronjob.Data,

		ErrorMessage: errMsg,
	}); err != nil {
		ag.logger.Error("failed to publish cron schedule completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
		return
	}
}
