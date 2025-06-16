package events

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const DescriptionPrefix = "FiveNet: "

// Ensures certain NATS config options are applied
type JSWrapper struct {
	jetstream.JetStream

	mu         sync.Mutex
	cfg        config.NATS
	shutdowner fx.Shutdowner
}

func NewJSWrapper(js jetstream.JetStream, cfg config.NATS, shutdowner fx.Shutdowner) *JSWrapper {
	return &JSWrapper{
		JetStream: js,

		mu:         sync.Mutex{},
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

type ConsumeErrRestartFn func(ctxTimeout context.Context, ctxCancel context.Context) error

func (j *JSWrapper) ConsumeErrHandlerWithRestart(ctxCancel context.Context, logger *zap.Logger, restartFn ConsumeErrRestartFn) jetstream.PullConsumeOpt {
	return jetstream.ConsumeErrHandler(func(ctxConsume jetstream.ConsumeContext, err error) {
		j.mu.Lock()
		defer j.mu.Unlock()

		if err != nil {
			logger.Error("error during jetstream consume, trying to restart...", zap.Error(err))

			if restartErr := j.consumeErrHandlerWithRestart(ctxCancel, logger, restartFn); restartErr != nil {
				logger.Error(fmt.Sprintf("failed to restart jetstream consumer after %d tries, attempting app shutdown", MaxRestartRetries), zap.Error(restartErr))

				if err := j.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		}
	})
}

func (j *JSWrapper) consumeErrHandlerWithRestart(ctxCancel context.Context, logger *zap.Logger, restartFn ConsumeErrRestartFn) error {
	var err error
	sleep := InitialRestartBackoffTime
	for try := range MaxRestartRetries {
		if func() bool {
			ctxTimeout, cancel := context.WithTimeout(ctxCancel, 10*time.Second)
			defer cancel()

			// Pass in a timeout context and the outer "passed in" context
			if err = restartFn(ctxTimeout, ctxCancel); err != nil {
				logger.Error(fmt.Sprintf("failed to restart jetstream consume, try %d of %d ...", try+1, MaxRestartRetries), zap.Error(err))

				if try < MaxRestartRetries {
					time.Sleep(sleep)
					sleep *= 2
				}
			} else {
				logger.Info(fmt.Sprintf("successfully restarted jetstream consume (try %d of %d)", try+1, MaxRestartRetries))
				return true
			}

			return false
		}() {
			break
		}
	}

	return err
}

func (j *JSWrapper) PublishProto(ctx context.Context, subject string, msg proto.Message, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	data, err := protoutils.MarshalToPJSON(msg)
	if err != nil {
		return nil, err
	}

	return j.Publish(ctx, subject, data, opts...)
}

func (j *JSWrapper) PublishAsyncProto(ctx context.Context, subject string, msg proto.Message, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	data, err := protoutils.MarshalToPJSON(msg)
	if err != nil {
		return nil, err
	}

	return j.PublishAsync(ctx, subject, data, opts...)
}

func (j *JSWrapper) BeginTx() *Transaction {
	return NewTransaction(j)
}
