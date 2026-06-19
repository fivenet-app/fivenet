package discord

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var SessionModule = fx.Module("discord.session",
	fx.Provide(
		NewDCSession,
	),
	fx.Decorate(wrapLogger),
)

type SessionParams struct {
	fx.In

	LC         fx.Lifecycle
	Shutdowner fx.Shutdowner

	Logger *zap.Logger
	Config *config.Config
}

func NewDCSession(p SessionParams) *session.Session {
	// Discord token not provided
	if p.Config.Discord.Token == "" {
		return nil
	}

	// Create a new Discord client/session using the provided login information.
	sess := session.NewWithIntents("Bot "+p.Config.Discord.Token,
		gateway.IntentDirectMessages,
		gateway.IntentGuildMembers,
		gateway.IntentGuildPresences,
		gateway.IntentGuildIntegrations,
	)

	return sess
}
