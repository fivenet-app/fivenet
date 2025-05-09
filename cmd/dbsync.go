package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/dbsync"
	"go.uber.org/fx"
)

type DBSyncCmd struct{}

func (c *DBSyncCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)
	fxOpts = append(fxOpts, fx.Invoke(func(*dbsync.Sync) {}))

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
