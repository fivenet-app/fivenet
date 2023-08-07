package perms

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("perms",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("perms")
}
