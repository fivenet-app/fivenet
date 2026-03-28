package commands

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	lang "github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var metricCommandCalls = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "discord_commands",
	Name:      "call_count",
	Help:      "Number of per command call count.",
}, []string{"command"})

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("discord.bot").Named("commands")
}

var Module = fx.Module("discord.commands",
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
	BotState discordtypes.BotState
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
	}

	c.router.Use(newMiddlewareLogger(c.logger))
	// Automatically defer handles if they're slow.
	c.router.Use(cmdroute.Deferrable(p.DC, cmdroute.DeferOpts{}))

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		c.logger.Info("registering commands", zap.Int("count", len(p.Commands)))

		commands := []api.CreateCommandData{}
		for _, cmd := range p.Commands {
			if cmd == nil {
				continue
			}

			cmdData := cmd.RegisterCommand(c.router)
			commands = append(commands, cmdData)

			metricCommandCalls.WithLabelValues(cmdData.Name).Set(0)
		}

		if err := cmdroute.OverwriteCommands(c.dc, commands); err != nil {
			return fmt.Errorf("cannot update discord bot commands. %w", err)
		}

		c.dc.AddInteractionHandler(c.router)

		return nil
	}))

	return c
}

func newMiddlewareLogger(logger *zap.Logger) cmdroute.Middleware {
	return func(next cmdroute.InteractionHandler) cmdroute.InteractionHandler {
		return cmdroute.InteractionHandlerFunc(
			func(ctx context.Context, ev *discord.InteractionEvent) *api.InteractionResponse {
				switch data := ev.Data.(type) {
				case *discord.CommandInteraction:
					logger.Info(
						"received interaction event",
						zap.Uint64("sender_id", uint64(ev.SenderID())),
						zap.String("command", data.Name),
					)
				}

				resp := next.HandleInteraction(ctx, ev)
				if resp == nil {
					// Most likely command not found error
					return nil
				}

				switch data := ev.Data.(type) {
				case *discord.CommandInteraction:
					metricCommandCalls.WithLabelValues(data.Name).Inc()
				}
				return resp
			},
		)
	}
}
