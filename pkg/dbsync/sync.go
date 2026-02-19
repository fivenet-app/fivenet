package dbsync

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"sync"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var Module = fx.Module("dbsync",
	fx.Provide(
		New,
	),
)

type Sync struct {
	mu sync.Mutex
	wg *sync.WaitGroup

	logger *zap.Logger
	db     *sql.DB

	ctx    context.Context
	cancel context.CancelFunc

	cfg     *dbsyncconfig.Config
	state   *dbsyncconfig.State
	cli     *grpc.ClientConn
	syncCli pbsync.SyncServiceClient

	jobs     *jobsSync
	licenses *licensesSync
	users    *usersSync
	vehicles *vehiclesSync
	accounts *accountsSync

	streamCh chan *pbsync.StreamResponse
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	Config *dbsyncconfig.Config
	DB     *sql.DB
	State  *dbsyncconfig.State
}

type Result struct {
	fx.Out

	Sync *Sync
}

func New(p Params) (*Sync, error) {
	logger := p.Logger.Named("dbsync")
	s := &Sync{
		wg: &sync.WaitGroup{},

		logger: logger,
		cfg:    p.Config,
		db:     p.DB,
		state:  p.State,

		streamCh: make(chan *pbsync.StreamResponse, 12),
	}

	p.Config.SetupWatch(logger.Named("config"), s.restart)

	// Load dbsync state from file if exists
	if err := s.state.Load(); err != nil {
		return nil, fmt.Errorf("failed to load dbsync state. %w", err)
	}

	if err := s.createGRPCClient(); err != nil {
		return nil, fmt.Errorf("failed to create sync API gRPC client. %w", err)
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(s.start))
	p.LC.Append(fx.StopHook(s.stop))

	return s, nil
}

func (s *Sync) createGRPCClient() error {
	// Create GRPC client for sync if destination is given
	cfg := s.cfg.Load()
	api := cfg.Destination.API
	if api.URL != "" {
		var transportCreds credentials.TransportCredentials
		if api.Insecure {
			transportCreds = insecure.NewCredentials()
		} else {
			transportCreds = credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS13,
				ClientAuth: tls.NoClientCert,
			})
		}

		cli, err := grpc.NewClient(api.URL,
			grpc.WithTransportCredentials(transportCreds),
			grpc.WithPerRPCCredentials(
				auth.NewClientTokenAuth(
					api.Token,
					api.TransportSecurity,
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

func (s *Sync) getSyncStatus(ctx context.Context) error {
	timeout := s.cfg.Load().Destination.API.Timeout

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := s.syncCli.GetStatus(ctx, &pbsync.GetStatusRequest{})
	if err != nil {
		return fmt.Errorf("failed to connect to sync service. %w", err)
	}

	out, err := protoutils.MarshalToPrettyJSON(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal sync status response. %w", err)
	}
	s.logger.Info("sync API status", zap.String("status", string(out)))

	return nil
}

func (s *Sync) start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("starting dbsync process")

	s.ctx, s.cancel = context.WithCancel(context.Background())

	// Setup table syncers
	syncer := &syncer{
		logger: s.logger,
		db:     s.db,
		cfg:    s.cfg.Load(),
		cli:    s.syncCli,
	}
	s.jobs = newJobsSync(syncer, s.state.Jobs)
	s.licenses = newLicensesSync(syncer, s.state.Licenses)

	s.accounts = newAccountsSync(syncer, s.state.Accounts)
	s.users = newUsersSync(syncer, s.state.Users)
	s.vehicles = newVehiclesSync(syncer, s.state.OwnedVehicles)

	s.wg.Go(func() {
		s.Run(s.ctx)
	})
	s.wg.Go(func() {
		s.RunStream(s.ctx)
	})

	return nil
}

func (s *Sync) stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
	}

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
	for {
		s.logger.Info("started dbsync loop")

		if err := s.run(ctx); err != nil {
			s.logger.Error("error during sync run", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(5 * time.Second):
		}
	}
}

func (s *Sync) run(ctx context.Context) error {
	if err := s.getSyncStatus(s.ctx); err != nil {
		s.logger.Error("failed to get sync status", zap.Error(err))
	}

	var wg sync.WaitGroup
	// On startup sync base data (jobs, job grades and license types) before the "main" sync loop starts,
	// then sync in 5 minute interval to keep the data fresh
	if err := s.syncBaseData(ctx); err != nil {
		s.logger.Error("error during initial base data sync", zap.Error(err))
		return err
	}

	// Base data sync loop
	wg.Go(func() {
		for {
			s.syncBaseData(ctx)

			select {
			case <-ctx.Done():
				return

			case <-time.After(5 * time.Minute):
			}
		}
	})

	cfg := s.cfg.Load().Tables

	// User data sync loop
	wg.Go(func() {
		if !cfg.Users.Enabled {
			return
		}

		for {
			if count, offset, err := s.users.Sync(ctx); err != nil {
				s.logger.Error("error during users sync", zap.Error(err))
			} else {
				s.logger.Info("users synced", zap.Int64("count", count), zap.Int64("offset", offset))
			}

			select {
			case <-ctx.Done():
				return

			case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Users)):
			}
		}
	})

	// Vehicles data sync loop
	wg.Go(func() {
		if !cfg.Vehicles.Enabled {
			return
		}

		for {
			if count, offset, err := s.vehicles.Sync(ctx); err != nil {
				s.logger.Error("error during vehicles sync", zap.Error(err))
			} else {
				s.logger.Info("vehicles synced", zap.Int64("count", count), zap.Int64("offset", offset))
			}

			select {
			case <-ctx.Done():
				return

			case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Vehicles)):
			}
		}
	})

	// Accounts data sync loop
	wg.Go(func() {
		if !cfg.Accounts.Enabled {
			return
		}

		for {
			if count, err := s.accounts.Sync(ctx); err != nil {
				s.logger.Error("error during accounts sync", zap.Error(err))
			} else {
				s.logger.Info("accounts synced", zap.Int64("count", count))
			}

			select {
			case <-ctx.Done():
				return

			case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Users)):
			}
		}
	})

	wg.Wait()

	return nil
}

func (s *Sync) syncBaseData(ctx context.Context) error {
	errs := multierr.Combine()

	if s.cfg.Load().Tables.Jobs.Enabled {
		if count, err := s.jobs.Sync(ctx); err != nil {
			errs = multierr.Append(errs, err)
			s.logger.Error("error during jobs sync", zap.Error(err))
		} else {
			s.logger.Info("jobs synced", zap.Int64("count", count))
		}
	}

	if s.cfg.Load().Tables.Licenses.Enabled {
		if count, err := s.licenses.Sync(ctx); err != nil {
			errs = multierr.Append(errs, err)
			s.logger.Error("error during licenses sync", zap.Error(err))
		} else {
			s.logger.Info("licenses synced", zap.Int64("count", count))
		}
	}

	return errs
}
