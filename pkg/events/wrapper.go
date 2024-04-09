package events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const DescriptionPrefix = "FiveNet: "

// Ensures certain NATS config options are applied
type JSWrapper struct {
	jetstream.JetStream

	cfg        config.NATS
	shutdowner fx.Shutdowner
}

func NewJSWrapper(js jetstream.JetStream, cfg config.NATS, shutdowner fx.Shutdowner) *JSWrapper {
	return &JSWrapper{
		JetStream:  js,
		cfg:        cfg,
		shutdowner: shutdowner,
	}
}

func (j *JSWrapper) CreateOrUpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	if cfg.Replicas == 0 || cfg.Replicas > j.cfg.Replicas {
		cfg.Replicas = j.cfg.Replicas
	}

	if !strings.HasPrefix(cfg.Description, DescriptionPrefix) {
		cfg.Description = DescriptionPrefix + cfg.Description
	}

	return j.JetStream.CreateOrUpdateStream(ctx, cfg)
}

func (j *JSWrapper) CreateOrUpdateKeyValue(ctx context.Context, cfg jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	if cfg.Replicas == 0 || cfg.Replicas > j.cfg.Replicas {
		cfg.Replicas = j.cfg.Replicas
	}

	if !strings.HasPrefix(cfg.Description, DescriptionPrefix) {
		cfg.Description = DescriptionPrefix + cfg.Description
	}

	return j.JetStream.CreateOrUpdateKeyValue(ctx, cfg)
}

const (
	MaxRestartRetries         = 5
	InitialRestartBackoffTime = 150 * time.Millisecond
)

func (j *JSWrapper) ConsumeErrHandler(logger *zap.Logger) jetstream.PullConsumeOpt {
	return jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		if err != nil {
			logger.Error("error during jetstream consume", zap.Error(err))
		}
	})
}

type ConsumeErrRestartFn func(ctx context.Context, c context.Context) error

func (j *JSWrapper) ConsumeErrHandlerWithRestart(c context.Context, logger *zap.Logger, restartFn ConsumeErrRestartFn) jetstream.PullConsumeOpt {
	return jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		if err != nil {
			logger.Error("error during jetstream consume, trying to restart...", zap.Error(err))

			sleep := InitialRestartBackoffTime
			var restartErr error
			for try := 0; try < MaxRestartRetries; try++ {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// Pass in a timeout context and the outer "passed in" context
				if restartErr = restartFn(ctx, c); restartErr != nil {
					logger.Error(fmt.Sprintf("failed to restart jetstream consume, try %d of %d ...", try+1, MaxRestartRetries), zap.Error(restartErr))

					if try < MaxRestartRetries {
						time.Sleep(sleep)
						sleep *= 2
					}
					continue
				} else {
					logger.Info(fmt.Sprintf("successfully restarted jetstream consume (try %d of %d)", try+1, MaxRestartRetries))
					break
				}
			}

			if restartErr != nil {
				logger.Error(fmt.Sprintf("failed to restart jetstream consume after %d tries, attempting app shutdown", MaxRestartRetries), zap.Error(restartErr))

				if err := j.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		}
	})
}

func (j *JSWrapper) addSpanInfoToMsg(ctx context.Context, msg *nats.Msg) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsSampled() {
		if msg.Header == nil {
			msg.Header = nats.Header{}
		}

		msg.Header.Set("X-Trace-Id", span.SpanContext().TraceID().String())
		msg.Header.Set("X-Span-Id", span.SpanContext().SpanID().String())
	}
}

func (j *JSWrapper) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return j.PublishMsg(ctx, &nats.Msg{Subject: subject, Data: payload}, opts...)
}

func (j *JSWrapper) PublishMsg(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	j.addSpanInfoToMsg(ctx, msg)
	return j.JetStream.PublishMsg(ctx, msg, opts...)
}

func (j *JSWrapper) PublishAsync(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	return j.PublishMsgAsync(ctx, &nats.Msg{Subject: subject, Data: payload}, opts...)
}

func (j *JSWrapper) PublishMsgAsync(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	j.addSpanInfoToMsg(ctx, msg)
	return j.JetStream.PublishMsgAsync(msg, opts...)
}

func GetJetstreamMsgContext(msg jetstream.Msg) (spanContext trace.SpanContext, err error) {
	headers := msg.Headers()

	var traceID trace.TraceID
	traceID, err = trace.TraceIDFromHex(headers.Get("X-Trace-Id"))
	if err != nil {
		return spanContext, err
	}
	var spanID trace.SpanID
	spanID, err = trace.SpanIDFromHex(headers.Get("X-Span-Id"))
	if err != nil {
		return spanContext, err
	}

	var spanContextConfig trace.SpanContextConfig
	spanContextConfig.TraceID = traceID
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = true
	spanContext = trace.NewSpanContext(spanContextConfig)

	return spanContext, nil
}

func GetNatsMsgContext(msg *nats.Msg) (spanContext trace.SpanContext, err error) {
	var traceID trace.TraceID
	traceID, err = trace.TraceIDFromHex(msg.Header.Get("X-Trace-Id"))
	if err != nil {
		return spanContext, err
	}
	var spanID trace.SpanID
	spanID, err = trace.SpanIDFromHex(msg.Header.Get("X-Span-Id"))
	if err != nil {
		return spanContext, err
	}

	var spanContextConfig trace.SpanContextConfig
	spanContextConfig.TraceID = traceID
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = true
	spanContext = trace.NewSpanContext(spanContextConfig)

	return spanContext, nil
}
