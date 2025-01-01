package dbsync

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
)

var Module = fx.Module("dbsync",
	fx.Provide(
		New,
	),
)

type Sync struct {
	wg *sync.WaitGroup
	db *sql.DB

	cfg *config.DBSync
}

type Params struct {
	fx.In

	LC fx.Lifecycle
}

func New(p Params) (*Sync, error) {
	s := &Sync{
		wg: &sync.WaitGroup{},
	}

	if err := s.loadConfig(); err != nil {
		return nil, err
	}

	if !s.cfg.Enabled {
		return nil, fmt.Errorf("dbsync is disabled in config")
	}

	// Connect to source database
	db, err := otelsql.Open("mysql", s.cfg.Source.DSN,
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
	s.db = db

	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	)); err != nil {
		return nil, err
	}

	// Setup SQL Prometheus metrics collector
	prometheus.MustRegister(collectors.NewDBStatsCollector(db, "fivenet"))

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go s.Run(ctx)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return db.Close()
	}))

	return s, nil
}

func (s *Sync) Run(ctx context.Context) {
	s.syncUsers(ctx)
	// TODO run one loop per source table
}

func (s *Sync) syncUsers(ctx context.Context) error {
	if !s.cfg.Source.Tables.Users.Enabled {
		return nil
	}

	query := s.cfg.Source.Tables.Users.Queries[0]
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}

	_ = rows

	return nil
}
