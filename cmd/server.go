package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type ServerCmd struct {
	ModuleCronAgent bool `default:"false" help:"Run the cron agent, should only be used for single container/binary deployments."`
}

func (c *ServerCmd) Run(_ *Context) error {
	instance.SetComponent("server")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, true, true)
	fxOpts = append(fxOpts, FxServerOpts()...)

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, FxCronerOpts()...)
		fxOpts = append(fxOpts, FxJobsHousekeeperOpts()...)
		fxOpts = append(fxOpts, FxHousekeeperOpts()...)
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
