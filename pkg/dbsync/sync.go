package dbsync

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/XSAM/otelsql"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	cli pbsync.SyncServiceClient
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
	s.state = NewDBSyncState(s.logger, s.cfg.StateFile)
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

	// Create GRPC client for sync if destination is given
	if s.cfg.Destination.URL != "" {
		cli, err := grpc.NewClient(s.cfg.Destination.URL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithPerRPCCredentials(auth.NewClientTokenAuth(s.cfg.Destination.Token)),
		)
		if err != nil {
			return nil, err
		}
		s.cli = pbsync.NewSyncServiceClient(cli)
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go s.Run(ctx)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		if err := s.state.Save(); err != nil {
			return err
		}

		if err := db.Close(); err != nil {
			return err
		}

		return nil
	}))

	return s, nil
}

func (s *Sync) Run(ctx context.Context) {
	for {
		if err := s.run(ctx); err != nil {
			s.logger.Error("error during sync run", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func (s *Sync) run(ctx context.Context) error {
	syncer := &syncer{
		logger: s.logger,
		db:     s.db,
		cfg:    s.cfg,
		cli:    s.cli,
	}

	// On startup sync jobs, job grades, licenses before the "main" sync loop starts
	jobs, err := NewJobsSync(syncer, s.state.Jobs)
	if err != nil {
		return err
	}
	licenses, err := NewLicensesSync(syncer, s.state.Licenses)
	if err != nil {
		return err
	}
	users, err := NewUsersSync(syncer, s.state.Users)
	if err != nil {
		return err
	}
	userLicenses, err := NewUserLicensesSync(syncer, s.state.UserLicenses)
	if err != nil {
		return err
	}

	if err := jobs.Sync(ctx); err != nil {
		s.logger.Error("error during jobs sync", zap.Error(err))
	}
	if err := licenses.Sync(ctx); err != nil {
		s.logger.Error("error during licenses sync", zap.Error(err))
	}

	// User data loop
	_ = userLicenses

	for {
		if err := users.Sync(ctx); err != nil {
			s.logger.Error("error during users sync", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return nil

		case <-time.After(1 * time.Second):
		}
	}
}

func prepareStringQuery(query config.DBSyncTable, state *TableSyncState, offset uint64, limit int) string {
	offsetStr := strconv.Itoa(int(offset))
	limitStr := strconv.Itoa(limit)

	q := strings.ReplaceAll(query.Query, "$offset", offsetStr)
	q = strings.ReplaceAll(q, "$limit", limitStr)

	where := ""
	if state == nil || (query.UpdatedField == nil || state.LastCheck.IsZero()) {
		q = strings.ReplaceAll(q, "$whereCondition", "")
	} else {
		where = fmt.Sprintf("WHERE `%s` >= '%s'\n",
			*query.UpdatedField,
			state.LastCheck.Format("2006-01-02 15:04:05"),
		)
	}

	q = strings.ReplaceAll(q, "$whereCondition", where)

	return q
}
