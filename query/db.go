package query

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"strconv"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/cmd/envs"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/dbutils/dsn"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
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
	// Use jsoniter as a replacement for std json lib for jet qrm (in case it is used)
	qrm.GlobalConfig.JsonUnmarshalFunc = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal

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
