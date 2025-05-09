package commands

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/lang"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord_bot").Named("commands")
}

var Module = fx.Module("discord_commands",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

type CommandParams struct {
	fx.In

	Cfg      *config.Config
	JS       *events.JSWrapper
	DB       *sql.DB
	L        *lang.I18n
	BotState types.BotState
	Perms    perms.Permissions
}

type CommandFactory = func(p CommandParams) (Command, error)

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	Cfg    *config.Config
	DC     *state.State

	Commands []Command `group:"discordcommands"`
}

type Cmds struct {
	logger *zap.Logger
	cfg    *config.Config
	dc     *state.State

	router *cmdroute.Router

	cmds []Command
}

func New(p Params) *Cmds {
	if !p.Cfg.Discord.Commands.Enabled {
		return nil
	}

	c := &Cmds{
		logger: p.Logger,
		cfg:    p.Cfg,
		dc:     p.DC,

		router: cmdroute.NewRouter(),

		cmds: p.Commands,
	}

	c.router.Use(newMiddlewareLogger(c.logger))
	// Automatically defer handles if they're slow.
	c.router.Use(cmdroute.Deferrable(p.DC, cmdroute.DeferOpts{}))

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := c.registerCommands(); err != nil {
			return fmt.Errorf("failed to register commands. %w", err)
		}

		return nil
	}))

	return c
}

func (c *Cmds) registerCommands() error {
	c.logger.Info("registering commands", zap.Int("count", len(c.cmds)))

	commands := []api.CreateCommandData{}
	for _, cmd := range c.cmds {
		cmdData := cmd.RegisterCommand(c.router)

		commands = append(commands, cmdData)
	}

	if err := cmdroute.OverwriteCommands(c.dc, commands); err != nil {
		return fmt.Errorf("cannot update discord bot commands. %w", err)
	}

	c.dc.AddInteractionHandler(c.router)

	return nil
}

func newMiddlewareLogger(logger *zap.Logger) cmdroute.Middleware {
	return func(next cmdroute.InteractionHandler) cmdroute.InteractionHandler {
		return cmdroute.InteractionHandlerFunc(func(ctx context.Context, ev *discord.InteractionEvent) *api.InteractionResponse {
			switch data := ev.Data.(type) {
			case *discord.CommandInteraction:
				logger.Info("received interaction event", zap.Uint64("sender_id", uint64(ev.SenderID())), zap.String("command", data.Name))
			}

			return next.HandleInteraction(ctx, ev)
		})
	}
}
