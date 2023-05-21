package query

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/XSAM/otelsql"
	"github.com/galexrt/fivenet/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

var (
	//go:embed migrations/*.sql
	migrationsFS embed.FS
	db           *sql.DB
)

func SetupDB(logger *zap.Logger) (*sql.DB, error) {
	if err := MigrateDB(logger, config.C.Database.DSN); err != nil {
		return nil, err
	}

	// Connect to database
	var err error
	db, err = otelsql.Open("mysql", config.C.Database.DSN, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		return nil, err
	}

	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.C.Database.MaxOpenConns)
	db.SetMaxIdleConns(config.C.Database.MaxIdleConns)
	db.SetConnMaxIdleTime(config.C.Database.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.C.Database.ConnMaxLifetime)

	return db, nil
}

func MigrateDB(logger *zap.Logger, dsn string) error {
	logger.Info("starting database migrations")
	// Connect to database
	db, err := sql.Open("mysql", dsn+"&multiStatements=true")
	if err != nil {
		return err
	}

	// Setup migrate source and driver
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "fivenet_zschema_migrations",
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance(
		"iofs", source,
		"mysql", driver)
	if err != nil {
		return err
	}
	// Run migrations
	if err := m.Up(); err != nil {
		if !errors.Is(migrate.ErrNoChange, err) {
			return err
		} else {
			logger.Info("database migrations have caused no changes")
		}
	} else {
		logger.Info("completed database migrations changes have been made")
	}

	return db.Close()
}
