package croner

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
)

var AgentModule = fx.Module("cron",
	fx.Provide(
		NewAgent,
	),
)

type AgentParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper

	Handlers *Handlers
}

type Agent struct {
	logger *zap.Logger
	ctx    context.Context
	js     *events.JSWrapper

	handlers *Handlers
}

func NewAgent(p AgentParams) (*Agent, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	ag := &Agent{
		logger: p.Logger.Named("cron_agent"),
		ctx:    ctxCancel,
		js:     p.JS,

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

func (ag *Agent) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	consumer, err := ag.js.CreateConsumer(ctxStartup, CronScheduleStreamName, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.%s", CronScheduleSubject, CronScheduleTopic),
		MaxDeliver:    3,
	})
	if err != nil {
		return err
	}

	if _, err := consumer.Consume(ag.watchForEvents,
		ag.js.ConsumeErrHandlerWithRestart(ctxCancel, ag.logger,
			ag.registerSubscriptions,
		)); err != nil {
		return err
	}

	return nil
}

func (ag *Agent) watchForEvents(msg jetstream.Msg) {
	job := &cron.CronjobSchedulerEvent{}
	if err := protojson.Unmarshal(msg.Data(), job); err != nil {
		ag.logger.Error("failed to unmarshal cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))

		if err := msg.NakWithDelay(150 * time.Millisecond); err != nil {
			ag.logger.Error("failed to nack unmarshal cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	fn := ag.handlers.getCronjobHandler(job.Cronjob.Name)
	if fn == nil {
		if err := msg.NakWithDelay(150 * time.Millisecond); err != nil {
			ag.logger.Error("failed to nack unmarshal cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	if err := msg.InProgress(); err != nil {
		ag.logger.Error("failed to send in progress for cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
	}

	start := time.Now()
	err := fn(ag.ctx, job.Cronjob.Data)
	elapsed := time.Since(start)

	// Update timestamp
	now := timestamp.Now()
	job.Cronjob.Data.UpdatedAt = now

	if _, err := ag.js.PublishProto(ag.ctx, fmt.Sprintf("%s.%s", CronScheduleSubject, CronCompleteTopic), &cron.CronjobCompletedEvent{
		Name:    job.Cronjob.Name,
		Sucess:  err == nil,
		Elapsed: durationpb.New(elapsed),
		EndDate: now,

		Data: job.Cronjob.Data,
	}); err != nil {
		ag.logger.Error("failed to publish cron schedule completion msg", zap.String("subject", msg.Subject()), zap.Error(err))
		return
	}

	if err := msg.Ack(); err != nil {
		ag.logger.Error("failed to ack cron schedule msg", zap.String("subject", msg.Subject()), zap.Error(err))
		return
	}
}
