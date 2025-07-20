package query

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/dsn"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/reqs"
	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/multierr"
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

// Result holds the result of the database setup.
type Result struct {
	fx.Out

	// DB is the SQL database connection instance.
	DB *sql.DB
	// Reqs is the database requirements instance.
	Reqs *reqs.DBReqs
}

// SetupDB sets up the database connection, runs migrations (unless skipped),
// configures connection pool, metrics, and ESX compatibility, and registers a stop hook for cleanup.
// Returns a *sql.DB instance or an error.
func SetupDB(p Params) (Result, error) {
	res := Result{}

	// Use jsoniter as a replacement for std json lib for jet qrm (in case it is used)
	qrm.GlobalConfig.JsonUnmarshalFunc = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal

	var req *reqs.DBReqs
	// Run DB migrations unless explicitly skipped via environment variable.
	if !p.Config.Database.SkipMigrations {
		var err error
		if req, err = MigrateDB(p.Logger, p.Config.Database.DSN, p.Config.IgnoreRequirements, p.Config.Database.ESXCompat, p.Config.Database.DisableLocking); err != nil {
			return res, err
		}
	}

	// Prepare the DSN (Data Source Name) for the database connection.
	dsn, err := dsn.PrepareDSN(p.Config.Database.DSN, p.Config.Database.DisableLocking)
	if err != nil {
		return res, fmt.Errorf("failed to prepare DSN. %w", err)
	}

	// Open database connection with OpenTelemetry instrumentation.
	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			DisableErrSkip: true,
		}),
	)
	if err != nil {
		return res, err
	}

	// Register DB stats metrics for Prometheus monitoring.
	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	)); err != nil {
		return res, err
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

	if req == nil {
		req = reqs.NewDBReqs(db)
		// Collect requirement errors if any
		var errs error
		if err := req.ValidateVersion(); err != nil {
			if !p.Config.IgnoreRequirements {
				errs = multierr.Append(errs, fmt.Errorf("failed to validate database version requirements. %w", err))
			}
			p.Logger.Warn("ignoring failed database version requirements", zap.Error(err))
		}

		if err := req.ValidateTables(); err != nil {
			if !p.Config.IgnoreRequirements {
				errs = multierr.Append(errs, fmt.Errorf("failed to validate database tables requirements. %w", err))
			}
			p.Logger.Warn("ignoring failed database tables requirements", zap.Error(err))
		}

		if !p.Config.IgnoreRequirements && errs != nil {
			return res, errs
		}
	} else {
		// Replace db connection with the one we just created
		req.ReplaceDB(db)
	}

	res.DB = db
	res.Reqs = req

	return res, nil
}
