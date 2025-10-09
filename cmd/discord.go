package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type DiscordCmd struct {
	ModuleCronAgent bool `default:"false" help:"Run the cron agent."`
}

func (c *DiscordCmd) Run(_ *Context) error {
	instance.SetComponent("discord")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, true, true)
	fxOpts = append(fxOpts, FxDiscordOpts()...)

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, fx.Invoke(func(*croner.Executor) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
