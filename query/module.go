package query

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("database",
	fx.Provide(
		SetupDB,
	),
	fx.Decorate(wrapLogger),
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("db")
}
