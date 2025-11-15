package dbsync

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/XSAM/otelsql"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/dsn"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	ctx    context.Context
	cancel context.CancelFunc

	cfg     *Config
	state   *DBSyncState
	cli     *grpc.ClientConn
	syncCli pbsync.SyncServiceClient

	jobs     *jobsSync
	licenses *licensesSync
	users    *usersSync
	vehicles *vehiclesSync

	streamCh chan *pbsync.StreamResponse
}

type Params struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner

	Logger *zap.Logger
	Config *Config
}

func New(p Params) (*Sync, error) {
	logger := p.Logger.Named("dbsync")
	s := &Sync{
		wg: &sync.WaitGroup{},

		logger: logger,
		cfg:    p.Config,

		streamCh: make(chan *pbsync.StreamResponse, 12),
	}

	p.Config.setupWatch(logger.Named("config"), s.restart)

	// Load dbsync state from file if exists
	s.state = NewDBSyncState(s.logger, s.cfg.Load().StateFile)
	if err := s.state.Load(); err != nil {
		return nil, fmt.Errorf("failed to load dbsync state. %w", err)
	}

	dsn, err := dsn.PrepareDSN(s.cfg.Load().Source.DSN, false)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare source database connection (bad DSN). %w", err)
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
	s.db = db

	if err := otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	)); err != nil {
		return nil, fmt.Errorf("failed to register db stats metrics. %w", err)
	}

	// Setup SQL Prometheus metrics collector
	prometheus.MustRegister(collectors.NewDBStatsCollector(db, "fivenet"))

	s.ctx, s.cancel = context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(s.start))
	p.LC.Append(fx.StopHook(func() error {
		if err := s.stop(); err != nil {
			return err
		}

		if err := s.db.Close(); err != nil {
			return err
		}

		return nil
	}))

	return s, nil
}

func (s *Sync) createGRPCClient() error {
	// Create GRPC client for sync if destination is given
	if s.cfg.Load().Destination.URL != "" {
		transportCreds := auth.NewCredentials()
		if !s.cfg.Load().Destination.Insecure {
			//nolint:gosec // G402: TLS MinVersion is set to TLS 1.1 for compatibility (gameservers may not support TLS 1.2 and higher)
			transportCreds = credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS11,
				ClientAuth: tls.NoClientCert,
			})
		}

		cli, err := grpc.NewClient(
			s.cfg.Load().Destination.URL,
			grpc.WithTransportCredentials(transportCreds),
			// Require transport security for release mode
			grpc.WithPerRPCCredentials(
				auth.NewClientTokenAuth(
					s.cfg.Load().Destination.Token,
					s.cfg.Load().Mode == "release",
				),
			),
		)
		if err != nil {
			return err
		}
		s.cli = cli
		s.syncCli = pbsync.NewSyncServiceClient(cli)
	}

	return nil
}

func (s *Sync) start() error {
	s.logger.Info("starting dbsync process")

	if err := s.createGRPCClient(); err != nil {
		return err
	}

	// Setup table syncers
	syncer := &syncer{
		logger: s.logger,
		db:     s.db,
		cfg:    s.cfg.Load(),
		cli:    s.syncCli,
	}
	s.jobs = newJobsSync(syncer, s.state.Jobs)
	s.licenses = newLicensesSync(syncer, s.state.Licenses)
	s.users = newUsersSync(syncer, s.state.Users)
	s.vehicles = newVehiclesSync(syncer, s.state.OwnedVehicles)

	s.wg.Add(1)
	go s.Run(s.ctx)

	if s.syncCli != nil {
		s.wg.Add(1)
		go s.RunStream(s.ctx)
	}

	return nil
}

func (s *Sync) stop() error {
	s.cancel()

	s.wg.Wait()

	var errs error
	if err := s.cli.Close(); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to close gRPC client. %w", err))
	}

	if err := s.state.Save(); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to save state. %w", err))
	}

	return errs
}

func (s *Sync) restart() error {
	s.logger.Info("stopping sync process")
	if err := s.stop(); err != nil {
		return err
	}
	s.logger.Info("stopped sync process")

	s.logger.Info("starting sync process")
	if err := s.start(); err != nil {
		return err
	}
	s.logger.Info("sync process started")

	return nil
}

func (s *Sync) Run(ctx context.Context) {
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
	var wg sync.WaitGroup
	// On startup sync base data (jobs, job grades and license types) before the "main" sync loop starts,
	// then sync in 5 minute interval to keep the data fresh
	if err := s.syncBaseData(ctx); err != nil {
		s.logger.Error("error during jobs sync", zap.Error(err))
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

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
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := s.users.Sync(ctx); err != nil {
				s.logger.Error("error during users sync", zap.Error(err))
			}

			select {
			case <-ctx.Done():
				return

			case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Users)):
			}
		}
	}()

	// Vehicles data sync loop
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := s.vehicles.Sync(ctx); err != nil {
				s.logger.Error("error during vehicles sync", zap.Error(err))
			}

			select {
			case <-ctx.Done():
				return

			case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Vehicles)):
			}
		}
	}()

	wg.Wait()

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
