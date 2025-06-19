package cmd

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/commands"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"go.uber.org/fx"
)

type DiscordCmd struct {
	ModuleCronAgent bool `help:"Run the cron agent." default:"false"`
}

func (c *DiscordCmd) Run(ctx *Context) error {
	instance.SetComponent("discord")

	fxOpts := getFxBaseOpts(Cli.StartTimeout, true)
	fxOpts = append(fxOpts,
		fx.Invoke(func(*discord.Bot) {}),
		fx.Invoke(func(*commands.Cmds) {}),
	)

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts, fx.Invoke(func(*croner.Executor) {}))
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
