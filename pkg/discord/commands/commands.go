package commands

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"go.uber.org/zap"
)

const GlobalCommandGuildID = "-1"

type CommandFactory = func(cfg *config.Config, l *lang.I18n) (api.CreateCommandData, cmdroute.CommandHandler, error)

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
	// Automatically defer handles if they're slow.
	c.router.Use(cmdroute.Deferrable(c.discord, cmdroute.DeferOpts{}))

	commands := []api.CreateCommandData{}
	for name, fn := range CommandsFactories {
		cmdData, handler, err := fn(c.cfg, c.i18n)
		if err != nil {
			return err
		}

		commands = append(commands, cmdData)

		c.router.Add(name, handler)
	}

	if err := cmdroute.OverwriteCommands(c.discord, commands); err != nil {
		c.logger.Fatal("cannot update discord bot commands", zap.Error(err))
	}

	c.discord.AddInteractionHandler(c.router)

	return nil
}
