package cmd

import (
	"github.com/fivenet-app/fivenet/v2026/cmd/fxopts"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord"
	discordcalendarreminders "github.com/fivenet-app/fivenet/v2026/pkg/discord/calendarreminders"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/commands"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"go.uber.org/fx"
)

type DiscordCmd struct {
	ModuleCronAgent bool `default:"false" help:"Run the cron agent."`
}

func (c *DiscordCmd) Run(cli *CLI) error {
	instance.SetComponent("discord")

	fxOpts := fxopts.GetFxBaseOpts(cli.StartTimeout, true, true)
	fxOpts = append(fxOpts,
		fx.Invoke(func(*discord.Bot) {}),
		fx.Invoke(func(*commands.Cmds) {}),
		fx.Invoke(func(*discordcalendarreminders.Worker) {}),
	)

	if c.ModuleCronAgent {
		fxOpts = append(fxOpts,
			fx.Invoke(func(*croner.Executor) {}),
		)
	}

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}
