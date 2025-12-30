package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/DeRuina/timberjack"
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

type LoggerResults struct {
	fx.Out

	Logger         *zap.Logger
	SuggaredLogger *zap.SugaredLogger
}

type LoggerParams struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

func NewLogger(p LoggerParams) (LoggerResults, error) {
	// Logger Setup
	level, err := zapcore.ParseLevel(p.Config.LogLevel)
	if err != nil {
		return LoggerResults{}, fmt.Errorf("failed to parse log level from config. %w", err)
	}

	var logger *zap.Logger
	if p.Config.Log.LogToFile {
		tl := &timberjack.Logger{
			Filename:         "/var/log/myapp/foo.log",
			MaxSize:          500,
			MaxBackups:       3,
			MaxAge:           28,
			Compress:         true,
			LocalTime:        true,
			RotationInterval: 24 * time.Hour,
		}
		w := zapcore.AddSync(tl)

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w, level,
		)

		logger = zap.New(core)
	} else {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Level.SetLevel(level)

		logger, err = loggerConfig.Build()
		if err != nil {
			return LoggerResults{}, fmt.Errorf("failed to configure logger. %w", err)
		}
	}

	zap.ReplaceGlobals(logger)

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		_ = logger.Sync()

		return nil
	}))

	return LoggerResults{
		Logger:         logger,
		SuggaredLogger: logger.Sugar(),
	}, nil
}
