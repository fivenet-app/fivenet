package croner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
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

	ctx       context.Context //nolint:containedctx // Executor lifecycle context is used for async job execution and publishing.
	js        *events.JSWrapper
	publisher events.IPublisher

	jsCons jetstream.ConsumeContext

	handlers *Handlers
	metrics  *executorMetrics
}

const cronExecutorAckWait = 40 * time.Minute

func NewExecutor(p ExecutorParams) (*Executor, error) {
	nodeName, err := getNodeName(p.Cfg.HTTP.AdminListen)
	if err != nil {
		return nil, fmt.Errorf("failed to get node name. %w", err)
	}

	ctxCancel, cancel := context.WithCancel(context.Background())

	logger := p.Logger.WithOptions(zap.IncreaseLevel(p.Cfg.Log.LevelOverrides.Get(config.LoggingComponentCron, p.Cfg.LogLevel))).
		Named("cron.executor")
	ag := &Executor{
		logger: logger,

		nodeName: nodeName,

		ctx:       ctxCancel,
		js:        p.JS,
		publisher: p.JS,

		handlers: p.Handlers,
		metrics:  getExecutorMetrics(),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := registerCronStreams(ctxCancel, ag.js); err != nil {
			return err
		}

		return ag.registerSubscriptions(ctxStartup, ctxCancel)
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		if ag.jsCons != nil {
			ag.jsCons.Stop()
			ag.jsCons = nil
		}
		cancel()

		return nil
	}))

	return ag, nil
}

func (ag *Executor) registerSubscriptions(
	ctxStartup context.Context,
	ctxCancel context.Context,
) error {
	consumer, err := ag.js.CreateOrUpdateConsumer(
		ctxStartup,
		CronScheduleStreamName,
		jetstream.ConsumerConfig{
			Durable:       CronExecutorConsumerName,
			DeliverPolicy: jetstream.DeliverAllPolicy,
			FilterSubject: fmt.Sprintf("%s.%s", CronScheduleSubject, CronScheduleTopic),
			AckPolicy:     jetstream.AckExplicitPolicy,
			AckWait:       cronExecutorAckWait,
			MaxDeliver:    3,
		},
	)
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
		ag.logger.Error(
			"failed to unmarshal cron schedule msg",
			zap.String("subject", msg.Subject()),
			zap.Error(err),
		)

		if err := msg.Term(); err != nil {
			ag.logger.Error(
				"failed to terminate invalid cron schedule msg",
				zap.String("subject", msg.Subject()),
				zap.Error(err),
			)
		}
		return
	}
	if job.GetCronjob() == nil || job.GetCronjob().GetName() == "" || job.GetCronjob().GetRunId() == "" {
		if err := msg.Term(); err != nil {
			ag.logger.Error("failed to terminate invalid cron schedule msg", zap.Error(err))
		}
		return
	}

	fn := ag.handlers.getCronjobHandler(job.GetCronjob().GetName())
	if fn == nil {
		if job.GetCronjob().GetData() == nil {
			job.Cronjob.Data = &cron.CronjobData{Data: &anypb.Any{}}
		}
		job.Cronjob.Data.UpdatedAt = timestamp.Now()
		if err := ag.publishCompletion(job, errors.New("cronjob handler is not registered"), 0); err != nil {
			ag.logger.Error(
				"failed to publish missing cron handler completion msg",
				zap.String("subject", msg.Subject()),
				zap.Error(err),
			)
			if nakErr := msg.NakWithDelay(250 * time.Millisecond); nakErr != nil {
				ag.logger.Error("failed to nack cron schedule msg", zap.Error(nakErr))
			}
			return
		}
		if err := msg.Ack(); err != nil {
			ag.logger.Error("failed to acknowledge missing cron handler message", zap.Error(err))
		}
		return
	}

	if job.GetCronjob().GetData() == nil {
		job.Cronjob.Data = &cron.CronjobData{
			Data: &anypb.Any{},
		}
	}

	var timeout time.Duration
	if job.GetCronjob().GetTimeout() != nil {
		timeout = job.GetCronjob().GetTimeout().AsDuration()
	} else {
		timeout = DefaultCronjobTimeout
	}

	ag.logger.Debug(
		"running cron job",
		zap.String("name", job.GetCronjob().GetName()),
		zap.Duration("timeout", timeout),
	)

	jobName := job.GetCronjob().GetName()
	startedTime := job.GetCronjob().GetStartedTime()
	if startedTime != nil {
		ag.metrics.startLatency.WithLabelValues(jobName).
			Observe(time.Since(startedTime.AsTime()).Seconds())
	}

	var elapsed time.Duration

	var handlerErr error
	heartbeatDone := make(chan struct{})
	go ag.heartbeat(msg, heartbeatDone)
	start := time.Now()
	func() {
		defer close(heartbeatDone)
		ctx, cancel := context.WithTimeout(ag.ctx, timeout)
		defer cancel()
		defer func() {
			elapsed = time.Since(start)

			// Recover from a panic and set err accordingly
			if e := recover(); e != nil {
				if er, ok := e.(error); ok {
					handlerErr = fmt.Errorf("recovered from panic. %w", er)
				} else {
					//nolint:errorlint // `er` is not guaranteed to be an error type, so we want it to be treated as a "string" here.
					handlerErr = fmt.Errorf("recovered from panic. %v", er)
				}

				ag.logger.Error(
					"cron job panic",
					zap.String("name", job.GetCronjob().GetName()),
					zap.Error(handlerErr),
				)
			}
		}()
		handlerErr = fn(ctx, job.GetCronjob().GetData())
	}()
	<-heartbeatDone

	status := "success"
	if handlerErr != nil {
		status = "failure"
	}
	ag.metrics.handlerDuration.WithLabelValues(jobName).Observe(elapsed.Seconds())
	if status == "success" {
		ag.metrics.lastRunSuccess.WithLabelValues(jobName).Set(1)
	} else {
		ag.metrics.lastRunSuccess.WithLabelValues(jobName).Set(0)
	}

	// Update timestamp in cronjob data
	now := timestamp.Now()
	job.Cronjob.Data.UpdatedAt = now
	if err := ag.publishCompletion(job, handlerErr, elapsed); err != nil {
		ag.logger.Error(
			"failed to publish cron schedule completion msg",
			zap.String("subject", msg.Subject()),
			zap.Error(err),
		)
		if nakErr := msg.NakWithDelay(250 * time.Millisecond); nakErr != nil {
			ag.logger.Error("failed to nack cron schedule msg", zap.Error(nakErr))
		}
		return
	}

	if err := msg.Ack(); err != nil {
		ag.logger.Error(
			"failed to acknowledge completed cron schedule msg",
			zap.String("subject", msg.Subject()),
			zap.Error(err),
		)
	}
}

func (ag *Executor) publishCompletion(job *cron.CronjobSchedulerEvent, handlerErr error, elapsed time.Duration) error {
	var errMsg *string
	if handlerErr != nil {
		msg := handlerErr.Error()
		errMsg = &msg
	}

	publishCtx, publishCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer publishCancel()
	_, err := ag.publisher.PublishProto(
		publishCtx,
		fmt.Sprintf("%s.%s", CronScheduleSubject, CronCompleteTopic),
		&cron.CronjobCompletedEvent{
			Name:         job.GetCronjob().GetName(),
			RunId:        job.GetCronjob().GetRunId(),
			StartedTime:  job.GetCronjob().GetStartedTime(),
			Success:      handlerErr == nil,
			Cancelled:    handlerErr != nil && errors.Is(handlerErr, context.Canceled),
			Elapsed:      durationpb.New(elapsed),
			EndDate:      timestamp.Now(),
			NodeName:     ag.nodeName,
			Data:         job.GetCronjob().GetData(),
			ErrorMessage: errMsg,
		},
	)
	return err
}

func (ag *Executor) heartbeat(msg jetstream.Msg, done <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := msg.InProgress(); err != nil {
				ag.logger.Warn("failed to extend cron schedule message", zap.Error(err))
			}
		}
	}
}
