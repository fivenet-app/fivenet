package query

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"strconv"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/v2025/cmd/envs"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/dsn"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// migrationsFS embeds all SQL migration files from the migrations directory.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

var Module = fx.Module("database",
	fx.Provide(
		SetupDB,
	),
	fx.Decorate(wrapLogger),
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("db")
}

// Params contains dependencies for constructing the database connection and managing its lifecycle.
type Params struct {
	fx.In

	// LC is the Fx lifecycle for managing start/stop hooks.
	LC fx.Lifecycle

	// Logger is the zap logger instance for logging.
	Logger *zap.Logger
	// Config is the application configuration.
	Config *config.Config
}

// SetupDB sets up the database connection, runs migrations (unless skipped),
// configures connection pool, metrics, and ESX compatibility, and registers a stop hook for cleanup.
// Returns a *sql.DB instance or an error.
func SetupDB(p Params) (*sql.DB, error) {
	// Use jsoniter as a replacement for std json lib for jet qrm (in case it is used)
	qrm.GlobalConfig.JsonUnmarshalFunc = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal

	// Run DB migrations unless explicitly skipped via environment variable.
	if skip, _ := strconv.ParseBool(os.Getenv(envs.SkipDBMigrationsEnv)); !skip {
		if err := MigrateDB(p.Logger, p.Config.Database.DSN, p.Config.Database.ESXCompat, p.Config.Database.DisableLocking); err != nil {
			return nil, err
		}
	}

	// Prepare the DSN (Data Source Name) for the database connection.
	dsn, err := dsn.PrepareDSN(p.Config.Database.DSN, p.Config.Database.DisableLocking)
	if err != nil {
		return nil, err
	}

	// Open database connection with OpenTelemetry instrumentation.
	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			DisableErrSkip: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	// Register DB stats metrics for Prometheus monitoring.
	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	)); err != nil {
		return nil, err
	}

	// Setup tables "helper" vars to work with ESX directly if enabled in config.
	if p.Config.Database.ESXCompat {
		tables.EnableESXCompat()
	}

	// Configure connection pool settings.
	db.SetMaxOpenConns(p.Config.Database.MaxOpenConns)
	db.SetMaxIdleConns(p.Config.Database.MaxIdleConns)
	db.SetConnMaxIdleTime(p.Config.Database.ConnMaxIdleTime)
	db.SetConnMaxLifetime(p.Config.Database.ConnMaxLifetime)

	// Register a stop hook to close the DB connection on application shutdown.
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		return db.Close()
	}))

	// Setup SQL Prometheus metrics collector for DB stats.
	prometheus.MustRegister(collectors.NewDBStatsCollector(db, "fivenet"))

	return db, nil
}
