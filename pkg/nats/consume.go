package nats

import (
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
