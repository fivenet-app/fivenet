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

func NewMigrate(db *sql.DB, esxCompat bool) (*migrate.Migrate, error) {
	// FiveNet's own users/chars table
	tableName := "fivenet_user"
	if esxCompat {
		// Use ESX's table
		tableName = "users"
	}

	// Setup migrate source and driver
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

func MigrateDB(logger *zap.Logger, dbDSN string, esxCompat bool) error {
	logger.Info("starting database migrations")

	dsn, err := dsn.PrepareDSN(dbDSN, dsn.WithMultiStatements())
	if err != nil {
		return err
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	m, err := NewMigrate(db, esxCompat)
	if err != nil {
		return err
	}
	m.Log = NewMigrateLogger(logger)

	// Run migrations
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
