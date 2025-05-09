package discord

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var StateModule = fx.Module("discord",
	fx.Provide(
		NewDCState,
	),
	fx.Decorate(wrapLogger),
)

type StateParams struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner

	Logger *zap.Logger
	Config *config.Config
}

func NewDCState(p StateParams) *state.State {
	// Discord bot not enabled
	if !p.Config.Discord.Enabled {
		return nil
	}

	cancelCtx, cancel := context.WithCancel(context.Background())

	// Create a new Discord session using the provided login information.
	state := state.New("Bot " + p.Config.Discord.Token)
	state.AddIntents(gateway.IntentGuildMembers)
	state.AddIntents(gateway.IntentGuildPresences)
	state.AddIntents(gateway.IntentGuildIntegrations)

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go func() {
			if err := state.Connect(cancelCtx); err != nil {
				p.Logger.Error("failed to connect to discord gateway", zap.Error(err))

				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					p.Logger.Fatal("failed to shutdown app via shutdowner", zap.Error(err))
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		if err := state.Close(); err != nil {
			p.Logger.Warn("error during discord client close", zap.Error(err))
		}

		cancel()

		return nil
	}))

	return state
}
