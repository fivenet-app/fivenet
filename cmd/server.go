package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/server"
	pbjobs "github.com/fivenet-app/fivenet/v2025/services/jobs"
	"go.uber.org/fx"
)

type ServerCmd struct {
	ModuleCronAgent bool `help:"Run the cron agent, should only be used for single container/binary deployments." default:"false"`
}

func (c *ServerCmd) Run(ctx *Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, true)
	fxOpts = append(fxOpts, fx.Invoke(func(server.HTTPServer) {}))

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, fx.Invoke(func(*croner.Executor) {}))
		fxOpts = append(fxOpts, fx.Invoke(func(*pbjobs.Housekeeper) {}))
		fxOpts = append(fxOpts, fx.Invoke(func(*housekeeper.Housekeeper) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
