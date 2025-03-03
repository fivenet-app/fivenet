package query

import (
	"fmt"
	"strings"

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

func (l *MigrateLogger) Printf(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(strings.TrimRight(format, "\n"), v...))
}

// Verbose should return true when verbose logging output is wanted
func (l *MigrateLogger) Verbose() bool {
	return false
}
