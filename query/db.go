package query

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"os"
	"strconv"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/cmd/envs"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/dsn"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
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
	if skip, _ := strconv.ParseBool(os.Getenv(envs.SkipDBMigrationsEnv)); !skip {
		if err := MigrateDB(p.Logger, p.Config.Database.DSN, p.Config.Database.ESXCompat); err != nil {
			return nil, err
		}
	}

	dsn, err := dsn.PrepareDSN(p.Config.Database.DSN)
	if err != nil {
		return nil, err
	}

	// Open database connection
	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
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

func NewMigrate(db *sql.DB, esxCompat bool) (*migrate.Migrate, error) {
	// FiveNet's own `users` table
	tableName := "fivenet_users"
	if esxCompat {
		// Use ESX's table
		tableName = "users"
	}

	// Setup migrate source and driver
	source, err := iofs.New(&templateFS{
		data: map[string]any{
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
