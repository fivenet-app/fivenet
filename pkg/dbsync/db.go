package dbsync

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils/dsn"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
)

type DBParams struct {
	fx.In

	LC fx.Lifecycle

	Config *Config
}

func NewDB(p DBParams) (*sql.DB, error) {
	dsn, err := dsn.PrepareDSN(p.Config.Load().Source.DSN, false)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to prepare source database connection (bad DSN). %w",
			err,
		)
	}

	// Connect to source database
	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(
			semconv.DBSystemMySQL,
		),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			DisableErrSkip: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to source database. %w", err)
	}

	reg, err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to register db stats metrics. %w", err)
	}

	// Setup SQL Prometheus metrics collector
	prometheus.MustRegister(collectors.NewDBStatsCollector(db, "fivenet"))

	p.LC.Append(fx.StopHook(func() error {
		_ = reg.Unregister()

		if err := db.Close(); err != nil {
			return err
		}

		return nil
	}))

	return db, nil
}
