package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/cmd/envs"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/dbutils/dsn"
	"github.com/fivenet-app/fivenet/query"
	"go.uber.org/fx"
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
	dsn, err := dsn.PrepareDSN(cfg.Database.DSN, dsn.WithMultiStatements())
	if err != nil {
		return err
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	m, err := query.NewMigrate(db, cfg.Database.ESXCompat)
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

func (c *UpCmd) run(_ context.Context, cfg *config.Config) error {
	dsn, err := dsn.PrepareDSN(cfg.Database.DSN, dsn.WithMultiStatements())
	if err != nil {
		return err
	}

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	m, err := query.NewMigrate(db, cfg.Database.ESXCompat)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
