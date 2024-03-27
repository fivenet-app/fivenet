package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func ConsumeErrHandler(logger *zap.Logger) jetstream.PullConsumeOpt {
	return jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		if err != nil {
			logger.Error("error during jetstream consume", zap.Error(err))
		}
	})
}

type ConsumeErrRestartHandler func(ctx context.Context, c context.Context) error

func ConsumeErrHandlerWithRestart(c context.Context, logger *zap.Logger, restartHandler ConsumeErrRestartHandler) jetstream.PullConsumeOpt {
	return jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		if err != nil {
			logger.Error("error during jetstream consume, trying to restart...", zap.Error(err))

			if restartHandler != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// Pass in a timeout context and the outer "passed in" context
				if err := restartHandler(ctx, c); err != nil {
					logger.Error("failed to restart jetstream consume", zap.Error(err))
				}
			}
		}
	})
}
