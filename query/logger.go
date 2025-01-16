package query

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MigrateLogger struct {
	logger  *zap.Logger
	verbose bool
}

func NewMigrateLogger(logger *zap.Logger) *MigrateLogger {
	return &MigrateLogger{
		logger:  logger.Named("migrate"),
		verbose: logger.Level() == zapcore.DebugLevel,
	}
}

func (l *MigrateLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

// Verbose should return true when verbose logging output is wanted
func (l *MigrateLogger) Verbose() bool {
	return false
}
