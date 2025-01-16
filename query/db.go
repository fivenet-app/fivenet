package query

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	Config *config.Config
}

func SetupDB(p Params) (*sql.DB, error) {
	if err := MigrateDB(p.Logger, p.Config.Database.DSN); err != nil {
		return nil, err
	}

	// Connect to database
	db, err := otelsql.Open("mysql", p.Config.Database.DSN,
		otelsql.WithAttributes(
			semconv.DBSystemMySQL,
		),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			DisableErrSkip: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	)); err != nil {
		return nil, err
	}

	// Setup tables "helper" vars to work with ESX directly
	if p.Config.Database.ESXCompat {
		tables.EnableESXCompat()
	}

	db.SetMaxOpenConns(p.Config.Database.MaxOpenConns)
	db.SetMaxIdleConns(p.Config.Database.MaxIdleConns)
	db.SetConnMaxIdleTime(p.Config.Database.ConnMaxIdleTime)
	db.SetConnMaxLifetime(p.Config.Database.ConnMaxLifetime)

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		return db.Close()
	}))

	// Setup SQL Prometheus metrics collector
	prometheus.MustRegister(collectors.NewDBStatsCollector(db, "fivenet"))

	return db, nil
}

func NewMigrate(db *sql.DB) (*migrate.Migrate, error) {
	// Setup migrate source and driver
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{
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

func MigrateDB(logger *zap.Logger, dsn string) error {
	logger.Info("starting database migrations")
	// Connect to database
	db, err := sql.Open("mysql", dsn+"&multiStatements=true")
	if err != nil {
		return err
	}

	m, err := NewMigrate(db)
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
