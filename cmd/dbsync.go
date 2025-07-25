package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type DBSyncCmd struct{}

func (c *DBSyncCmd) Run(ctx *Context) error {
	instance.SetComponent("dbsync")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)
	fxOpts = append(fxOpts, FxDBSyncOpts()...)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
