package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"go.uber.org/zap"
)

const GlobalCommandGuildID = "-1"

type CommandFactory = func(router *cmdroute.Router, cfg *config.Config, l *lang.I18n) (api.CreateCommandData, error)

var CommandsFactories = map[string]CommandFactory{}

type Cmds struct {
	logger *zap.Logger

	router  *cmdroute.Router
	discord *state.State
	cfg     *config.Config
	i18n    *lang.I18n
}

func New(logger *zap.Logger, s *state.State, cfg *config.Config, i18n *lang.I18n) (*Cmds, error) {
	c := &Cmds{
		logger:  logger,
		discord: s,
		cfg:     cfg,
		i18n:    i18n,
	}

	return c, nil
}

func (c *Cmds) RegisterCommands() error {
	c.logger.Info("registering commands", zap.Int("count", len(CommandsFactories)))

	c.router = cmdroute.NewRouter()
	c.router.Use(Logger(c.logger.Named("discord_bot.commands")))
	// Automatically defer handles if they're slow.
	c.router.Use(cmdroute.Deferrable(c.discord, cmdroute.DeferOpts{}))

	commands := []api.CreateCommandData{}
	for _, fn := range CommandsFactories {
		cmdData, err := fn(c.router, c.cfg, c.i18n)
		if err != nil {
			return err
		}

		commands = append(commands, cmdData)
	}

	if err := cmdroute.OverwriteCommands(c.discord, commands); err != nil {
		c.logger.Fatal("cannot update discord bot commands", zap.Error(err))
	}

	c.discord.AddInteractionHandler(c.router)

	return nil
}

func Logger(logger *zap.Logger) cmdroute.Middleware {
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
