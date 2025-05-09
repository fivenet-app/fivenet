package modules

import (
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerModule = fx.Module("logger",
	fx.Provide(
		NewLogger,
	),
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// Logger Setup
	loggerConfig := zap.NewProductionConfig()
	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level from config. %w", err)
	}
	loggerConfig.Level.SetLevel(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to configure logger. %w", err)
	}

	return logger, nil
}
