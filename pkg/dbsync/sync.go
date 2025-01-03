package dbsync

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/XSAM/otelsql"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("dbsync",
	fx.Provide(
		New,
	),
)

type Sync struct {
	wg *sync.WaitGroup

	logger *zap.Logger
	db     *sql.DB

	cfg   *config.DBSync
	state *DBSyncState
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
}

func New(p Params) (*Sync, error) {
	s := &Sync{
		wg: &sync.WaitGroup{},

		logger: p.Logger.Named("dbsync"),
	}

	if err := s.loadConfig(); err != nil {
		return nil, err
	}

	if !s.cfg.Enabled {
		return nil, fmt.Errorf("dbsync is disabled in config")
	}

	// Load dbsync state from file if exists
	s.state = &DBSyncState{
		mu:       sync.Mutex{},
		filepath: s.cfg.StateFile,
		States:   map[string]*TableSyncState{},
	}
	if err := s.state.Load(); err != nil {
		return nil, err
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
	us := &usersSync{
		logger: s.logger,
		db:     s.db,
		cfg:    s.cfg,
	}
	if _, err := us.Sync(ctx); err != nil {
		s.logger.Error("error during users sync", zap.Error(err))
	}

	// TODO run one loop per source table
}

func prepareStringQuery(in string, offset int, limit int) string {
	offsetStr := strconv.Itoa(offset)
	limitStr := strconv.Itoa(limit)

	return strings.ReplaceAll(strings.ReplaceAll(in, "$offset", offsetStr), "$limit", limitStr)
}
