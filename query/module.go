package query

import (
	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("database",
	fx.Provide(
		SetupDB,
	),
	fx.Decorate(wrapLogger),
)

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	Config *config.Config
}

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("db")
}
