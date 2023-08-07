package audit

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("audit",
	fx.Provide(
		New,
	),
	fx.Decorate(
		wrapLogger,
	),
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("audit")
}
