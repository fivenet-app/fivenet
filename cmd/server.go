package cmd

import (
	"github.com/fivenet-app/fivenet/v2026/cmd/fxopts"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"go.uber.org/fx"
)

type ServerCmd struct {
	ModuleCronAgent bool `default:"false" help:"Run the cron agent, should only be used for single container/binary deployments."`
}

func (c *ServerCmd) Run(cli *CLI) error {
	instance.SetComponent("server")

	fxOpts := fxopts.GetFxBaseOpts(cli.StartTimeout, true, true)
	fxOpts = append(fxOpts, fxopts.FxServerOpts()...)

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, fxopts.FxCronerOpts()...)
		fxOpts = append(fxOpts, fxopts.FxServiceHousekeeperOpts()...)
		fxOpts = append(fxOpts, fxopts.FxHousekeeperOpts()...)
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
