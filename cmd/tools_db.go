package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/v2025/cmd/envs"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/dsn"
	"github.com/fivenet-app/fivenet/v2025/query"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DBCmd struct {
	Version VersionCmd `cmd:"" help:"Display db migration version info"`
	Up      UpCmd      `cmd:"" help:"Run any outstanding migrations"`
}

type VersionCmd struct{}

func (c *VersionCmd) Run(ctx *kong.Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(fxOpts,
		fx.Invoke(func(lifecycle fx.Lifecycle, cfg *config.Config, shutdowner fx.Shutdowner) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						exitCode := 0
						if err := c.run(ctx, cfg); err != nil {
							// handle error, set non-zero exit code so caller knows the job failed
							exitCode = 1
						}
						_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
					}()
					return nil
				},
			})
		}),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

func (c *VersionCmd) run(_ context.Context, cfg *config.Config) error {
	dsn, err := dsn.PrepareDSN(cfg.Database.DSN, cfg.Database.DisableLocking, dsn.WithMultiStatements())
	if err != nil {
		return err
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	m, err := query.NewMigrate(db, cfg.Database.ESXCompat, cfg.Database.DisableLocking)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}

	fmt.Printf("Version: %d (Dirty: %t)\n", version, dirty)

	return nil
}

type UpCmd struct{}

func (c *UpCmd) Run(ctx *kong.Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(fxOpts,
		fx.Invoke(func(logger *zap.Logger, lifecycle fx.Lifecycle, cfg *config.Config, shutdowner fx.Shutdowner) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						exitCode := 0
						if err := c.run(ctx, logger, cfg); err != nil {
							logger.Error("Failed to run migrations", zap.Error(err))
							// handle error, set non-zero exit code so caller knows the job failed

							exitCode = 1
						}
						_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
					}()
					return nil
				},
			})
		}),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

func (c *UpCmd) run(_ context.Context, logger *zap.Logger, cfg *config.Config) error {
	dsn, err := dsn.PrepareDSN(cfg.Database.DSN, cfg.Database.DisableLocking, dsn.WithMultiStatements())
	if err != nil {
		return fmt.Errorf("failed to prepare dsn. %w", err)
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open mysqsl db connection. %w", err)
	}

	m, err := query.NewMigrate(db, cfg.Database.ESXCompat, cfg.Database.DisableLocking)
	if err != nil {
		return fmt.Errorf("failed to create migrationg client. %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate. %w", err)
		}

		logger.Info("No migrations to apply, database is up to date")
	}

	return nil
}
