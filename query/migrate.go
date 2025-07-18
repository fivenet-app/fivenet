package query

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/dsn"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MigrateLogger implements the migrate.Logger interface using zap.Logger for logging migration output.
type MigrateLogger struct {
	// logger is the zap logger instance used for migration logs.
	logger *zap.Logger
	// verbose indicates if verbose logging is enabled.
	verbose bool
}

// NewMigrateLogger creates a new MigrateLogger with the given zap.Logger.
// The logger is named "migrate" and verbosity is set based on the logger's level.
func NewMigrateLogger(logger *zap.Logger) *MigrateLogger {
	return &MigrateLogger{
		logger:  logger.Named("migrate"),
		verbose: logger.Level() == zapcore.DebugLevel,
	}
}

// Printf logs formatted migration output using zap at Info level.
func (l *MigrateLogger) Printf(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(strings.TrimRight(format, "\n"), v...))
}

// Verbose returns true if verbose logging output is wanted for migrations.
func (l *MigrateLogger) Verbose() bool {
	return false
}

// NewMigrate creates a new migrate.Migrate instance for the given DB and ESX compatibility flag.
// It sets up the migration source and driver, and injects template data for migration scripts.
func NewMigrate(db *sql.DB, esxCompat bool, disableLocking bool) (*migrate.Migrate, error) {
	// FiveNet's own users/chars table
	tableName := "fivenet_user"
	if esxCompat {
		// Use ESX's table
		tableName = "users"
	}

	// Setup migrate source and driver with template data for, e.g., ESX compatibility.
	source, err := iofs.New(&templateFS{
		data: map[string]any{
			"ESXCompat":      esxCompat,
			"UsersTableName": tableName,
		},
		FS: migrationsFS,
	}, "migrations")
	if err != nil {
		return nil, err
	}

	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{
		MigrationsTable: "fivenet_zschema_migrations",
		NoLock:          disableLocking,
	})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs", source,
		"mysql", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// MigrateDB runs database migrations using golang-migrate, logging progress and errors.
// It prepares the DSN, connects to the DB, runs migrations, and logs the result.
func MigrateDB(logger *zap.Logger, dbDSN string, esxCompat bool, disableLocking bool) error {
	logger.Info("starting database migrations")

	dsn, err := dsn.PrepareDSN(dbDSN, disableLocking, dsn.WithMultiStatements())
	if err != nil {
		return fmt.Errorf("failed to prepare DSN. %w", err)
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database. %w", err)
	}

	m, err := NewMigrate(db, esxCompat, disableLocking)
	if err != nil {
		return fmt.Errorf("failed to create migration instance. %w", err)
	}
	m.Log = NewMigrateLogger(logger)

	// Run migrations and handle "no change" as a non-error.
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		} else {
			logger.Info("database migrations have caused no changes")
		}
	} else {
		logger.Info("completed database migrations changes have been made")
	}

	return db.Close()
}
