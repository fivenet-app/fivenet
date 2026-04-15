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
	"github.com/fivenet-app/fivenet/v2026/pkg/dbsync/syncers"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	ctx    context.Context //nolint:containedctx // Sync process retains service lifecycle context for long-running loops.
	cancel context.CancelFunc

	cfg     *dbsyncconfig.Config
	state   *dbsyncconfig.State
	cli     *grpc.ClientConn
	syncCli pbsync.SyncServiceClient

	streamCh chan *pbsync.StreamResponse

	// Base data syncers
	jobs     *syncers.JobsSync
	licenses *syncers.LicensesSync

	// Main data syncers
	accounts *syncers.AccountsSync

	users       *syncers.UsersSync
	usersResync *syncers.UsersSync

	vehicles       *syncers.VehiclesSync
	vehiclesResync *syncers.VehiclesSync
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

	if err := s.createGRPCClient(); err != nil {
		return nil, fmt.Errorf("failed to create sync API gRPC client. %w", err)
	}

	p.LC.Append(fx.StartHook(s.start))
	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		return s.stop(true)
	}))

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

	ctx, cancel := context.WithTimeout(ctx, timeout)
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

	// Setup data syncers
	syncer := syncers.New(s.logger, s.db, s.cfg.Load(), s.syncCli)
	s.jobs = syncers.NewJobsSync(syncer, s.state.Jobs)
	s.licenses = syncers.NewLicensesSync(syncer, s.state.Licenses)

	s.accounts = syncers.NewAccountsSync(syncer, s.state.Accounts)
	s.users = syncers.NewUsersSync(syncer, s.state.Users, true)
	s.usersResync = syncers.NewUsersSync(syncer, s.state.UsersResync, false)
	s.vehicles = syncers.NewVehiclesSync(syncer, s.state.Vehicles, true)
	s.vehiclesResync = syncers.NewVehiclesSync(syncer, s.state.VehiclesResync, false)

	s.wg.Go(func() {
		s.runLoop(s.ctx)
	})

	// Run stream only when not in dry run mode
	if !s.cfg.Load().Destination.DryRun {
		s.wg.Go(func() {
			s.RunStream(s.ctx)
		})
	} else {
		s.logger.Warn("dry run enabled, not sending data to server")
	}

	return nil
}

func (s *Sync) stop(closeCli bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
	}

	s.wg.Wait()

	if s.cli != nil && closeCli {
		if err := s.cli.Close(); err != nil {
			st := status.Convert(err)
			if st.Code() == codes.Unavailable || st.Code() == codes.Canceled {
				return nil
			}

			return fmt.Errorf("failed to close gRPC client. %w", err)
		}
	}

	return nil
}

func (s *Sync) restart() error {
	s.logger.Info("stopping sync process")
	if err := s.stop(false); err != nil {
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

func (s *Sync) runLoop(ctx context.Context) {
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
	if err := s.getSyncStatus(ctx); err != nil {
		s.logger.Error("failed to get sync status", zap.Error(err))
	}

	cfg := s.cfg.Load().Tables

	var wg sync.WaitGroup
	if cfg.Jobs.Enabled || cfg.Licenses.Enabled {
		s.logger.Info("starting base data sync")
		// On startup sync base data (jobs, job grades and license types) before the "main" sync loop starts,
		// then sync in 5 minute interval to keep the data fresh
		if err := s.syncBaseData(ctx); err != nil {
			s.logger.Error("error during initial base data sync", zap.Error(err))
			return err
		}

		// Base data sync loop
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return

				case <-time.After(5 * time.Minute):
				}

				s.syncBaseData(ctx)
			}
		})
	}

	// User data sync loop
	if cfg.Users.Enabled {
		wg.Go(func() {
			s.logger.Info("starting users sync")
			for {
				startedAt := time.Now()
				fetched, sent, lastID, lastUpdatedAt, err := s.users.Sync(ctx)
				recordSyncMetrics("users", startedAt, fetched, sent, err)
				if err != nil {
					s.logger.Error("error during users sync", zap.Error(err))
				} else {
					s.logger.Info("users synced",
						zap.Int64("fetched", fetched),
						zap.Int64("sent", sent),
						zap.String("last_id", lastID),
						zap.Timep("last_updated_at", lastUpdatedAt),
					)
				}

				select {
				case <-ctx.Done():
					return

				case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Users)):
				}
			}
		})

		if cfg.Users.ResyncInterval != nil && *cfg.Users.ResyncInterval > 0 {
			resyncInterval := *cfg.Users.ResyncInterval
			wg.Go(func() {
				s.logger.Info("starting users resync")
				for {
					select {
					case <-ctx.Done():
						return

					case <-time.After(resyncInterval):
						startedAt := time.Now()
						fetched, sent, lastID, _, err := s.usersResync.Resync(ctx)
						recordSyncMetrics("users_resync", startedAt, fetched, sent, err)
						if err != nil {
							s.logger.Error("error during users resync", zap.Error(err))
						} else {
							fields := []zap.Field{
								zap.Int64("fetched", fetched),
								zap.Int64("sent", sent),
								zap.String("last_id", lastID),
							}
							s.logger.Info("users resynced", fields...)
						}
					}
				}
			})
		}
	}

	// Vehicles data sync loop
	if cfg.Vehicles.Enabled {
		wg.Go(func() {
			s.logger.Info("starting vehicles sync")
			for {
				startedAt := time.Now()
				fetched, sent, plate, lastUpdatedAt, err := s.vehicles.Sync(ctx)
				recordSyncMetrics("vehicles", startedAt, fetched, sent, err)
				if err != nil {
					s.logger.Error("error during vehicles sync", zap.Error(err))
				} else {
					s.logger.Info("vehicles synced",
						zap.Int64("fetched", fetched),
						zap.Int64("sent", sent),
						zap.String("last_plate", plate),
						zap.Timep("last_updated_at", lastUpdatedAt),
					)
				}

				select {
				case <-ctx.Done():
					return

				case <-time.After(s.cfg.Load().GetSyncInterval(&s.cfg.Load().Tables.Vehicles)):
				}
			}
		})

		if cfg.Vehicles.ResyncInterval != nil && *cfg.Vehicles.ResyncInterval > 0 {
			resyncInterval := *cfg.Vehicles.ResyncInterval
			wg.Go(func() {
				s.logger.Info("starting vehicles resync")
				for {
					select {
					case <-ctx.Done():
						return

					case <-time.After(resyncInterval):
						startedAt := time.Now()
						fetched, sent, lastID, _, err := s.vehiclesResync.Resync(ctx)
						recordSyncMetrics("vehicles_resync", startedAt, fetched, sent, err)
						if err != nil {
							s.logger.Error("error during vehicles resync", zap.Error(err))
						} else {
							fields := []zap.Field{
								zap.Int64("fetched", fetched),
								zap.Int64("sent", sent),
								zap.String("last_id", lastID),
							}
							s.logger.Info("vehicles resynced", fields...)
						}
					}
				}
			})
		}
	}

	// Accounts data sync loop
	if cfg.Accounts.Enabled {
		wg.Go(func() {
			s.logger.Info("starting accounts sync")
			for {
				startedAt := time.Now()
				count, err := s.accounts.Sync(ctx)
				recordSyncMetrics("accounts", startedAt, count, count, err)
				if err != nil {
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
	}

	wg.Wait()

	return nil
}

func (s *Sync) syncBaseData(ctx context.Context) error {
	errs := multierr.Combine()

	if s.cfg.Load().Tables.Jobs.Enabled {
		startedAt := time.Now()
		count, err := s.jobs.Sync(ctx)
		recordSyncMetrics("jobs", startedAt, count, count, err)
		if err != nil {
			errs = multierr.Append(errs, err)
			s.logger.Error("error during jobs sync", zap.Error(err))
		} else {
			s.logger.Info("jobs synced", zap.Int64("count", count))
		}
	}

	if s.cfg.Load().Tables.Licenses.Enabled {
		startedAt := time.Now()
		count, err := s.licenses.Sync(ctx)
		recordSyncMetrics("licenses", startedAt, count, count, err)
		if err != nil {
			errs = multierr.Append(errs, err)
			s.logger.Error("error during licenses sync", zap.Error(err))
		} else {
			s.logger.Info("licenses synced", zap.Int64("count", count))
		}
	}

	return errs
}
