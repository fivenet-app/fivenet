package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
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

	Logger  *zap.Logger
	ZapCore zapcore.Core
}

type LoggerParams struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

func NewLogger(p LoggerParams) (LoggerResults, error) {
	// Logger Setup
	loggerConfig := zap.NewProductionConfig()
	level, err := zapcore.ParseLevel(p.Config.LogLevel)
	if err != nil {
		return LoggerResults{}, fmt.Errorf("failed to parse log level from config. %w", err)
	}
	loggerConfig.Level.SetLevel(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return LoggerResults{}, fmt.Errorf("failed to configure logger. %w", err)
	}

	res, err := newAttributes(p.Config.OTLP)
	if err != nil {
		return LoggerResults{}, fmt.Errorf(
			"failed to create attributes for tracer provider. %w",
			err,
		)
	}

	ctx := context.Background()
	logExporter, err := newLoggerExporter(
		ctx,
		p.Config.OTLP.Type,
		p.Config.OTLP.URL,
		p.Config.OTLP.Insecure,
		p.Config.OTLP.Timeout,
		p.Config.OTLP.Headers,
		p.Config.OTLP.Compression,
	)
	if err != nil {
		return LoggerResults{}, fmt.Errorf("failed to create logger exporter. %w", err)
	}

	if p.Config.OTLP.Enabled {
		loggerProvider := log.NewLoggerProvider(
			log.WithProcessor(log.NewBatchProcessor(logExporter)),
			log.WithResource(res),
		)
		global.SetLoggerProvider(loggerProvider)

		core := zapcore.NewTee(
			logger.Core(),
			otelzap.NewCore(
				"github.com/fivenet-app/fivenet",
				otelzap.WithLoggerProvider(loggerProvider),
			),
		)
		logger = zap.New(core)

		p.LC.Append(fx.StopHook(func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
			defer cancel()

			// Gracefully shutdown components
			if err := loggerProvider.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to cleanly shut down logger provider. %w", err)
			}

			return nil
		}))
	}

	zap.ReplaceGlobals(logger)

	return LoggerResults{
		Logger: logger,
	}, nil
}
