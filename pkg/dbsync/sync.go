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
	"go.uber.org/multierr"
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

	acfg  *config.Config
	cfg   *DBSync
	state *DBSyncState
	cli   pbsync.SyncServiceClient

	jobs     *jobsSync
	licenses *licensesSync
	users    *usersSync
	vehicles *vehiclesSync

	streamCh chan *pbsync.StreamResponse
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	Config *config.Config
}

func New(p Params) (*Sync, error) {
	s := &Sync{
		wg: &sync.WaitGroup{},

		logger: p.Logger.Named("dbsync"),
		acfg:   p.Config,

		streamCh: make(chan *pbsync.StreamResponse, 12),
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
			// Require transport security for release mode
			grpc.WithPerRPCCredentials(auth.NewClientTokenAuth(s.cfg.Destination.Token, s.acfg.Mode == "release")),
		)
		if err != nil {
			return nil, err
		}
		s.cli = pbsync.NewSyncServiceClient(cli)
	}

	// Setup table syncers
	syncer := &syncer{
		logger: s.logger,
		db:     s.db,
		cfg:    s.cfg,
		cli:    s.cli,
	}
	s.jobs = newJobsSync(syncer, s.state.Jobs)
	s.licenses = newLicensesSync(syncer, s.state.Licenses)
	s.users = newUsersSync(syncer, s.state.Users)
	s.vehicles = newVehiclesSync(syncer, s.state.OwnedVehicles)

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go s.Run(ctx)

		if s.cli != nil {
			go s.RunStream(ctx)
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		s.wg.Wait()

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
	s.wg.Add(1)
	defer s.wg.Done()

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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// On startup sync base data (jobs, job grades and license types) before the "main" sync loop starts,
	// then sync in 5 minute interval to keep the data fresh
	if err := s.syncBaseData(ctx); err != nil {
		s.logger.Error("error during jobs sync", zap.Error(err))
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			s.syncBaseData(ctx)

			select {
			case <-ctx.Done():
				return

			case <-time.After(5 * time.Minute):
			}
		}
	}()

	// User data sync loop
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		for {
			if err := s.users.Sync(ctx); err != nil {
				s.logger.Error("error during users sync", zap.Error(err))
			}

			if err := s.vehicles.Sync(ctx); err != nil {
				s.logger.Error("error during users sync", zap.Error(err))
			}

			select {
			case <-ctx.Done():
				return

			case <-time.After(2 * time.Second):
			}
		}
	}()

	return nil
}

func (s *Sync) syncBaseData(ctx context.Context) error {
	errs := multierr.Combine()

	if err := s.jobs.Sync(ctx); err != nil {
		errs = multierr.Append(errs, err)
		s.logger.Error("error during jobs sync", zap.Error(err))
	}

	if err := s.licenses.Sync(ctx); err != nil {
		errs = multierr.Append(errs, err)
		s.logger.Error("error during licenses sync", zap.Error(err))
	}

	return errs
}

func prepareStringQuery(query DBSyncTable, state *TableSyncState, offset uint64, limit int) string {
	offsetStr := strconv.Itoa(int(offset))
	limitStr := strconv.Itoa(limit)

	q := strings.ReplaceAll(query.Query, "$offset", offsetStr)
	q = strings.ReplaceAll(q, "$limit", limitStr)

	where := ""
	// Add "updatedAt" column condition if available
	if state == nil || (query.UpdatedTimeColumn == nil || state.LastCheck.IsZero()) {
		q = strings.ReplaceAll(q, "$whereCondition", "")
	} else {
		where = fmt.Sprintf("WHERE `%s` >= '%s'\n",
			*query.UpdatedTimeColumn,
			state.LastCheck.Format("2006-01-02 15:04:05"),
		)
	}

	q = strings.ReplaceAll(q, "$whereCondition", where)

	return q
}
