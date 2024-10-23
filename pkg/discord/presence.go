package discord

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"go.uber.org/zap"
)

func (b *Bot) setBotPresence(cfg *rector.DiscordBotPresence) error {
	var activity *discord.Activity
	if cfg.Type == rector.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_GAME {
		activity = &discord.Activity{
			Type: discord.GameActivity,
			Name: *cfg.Status,
		}
	} else if cfg.Type == rector.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_LISTENING {
		activity = &discord.Activity{
			Type: discord.ListeningActivity,
			Name: *cfg.Status,
		}
	} else if cfg.Type == rector.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_STREAMING {
		activity = &discord.Activity{
			Type: discord.StreamingActivity,
			Name: *cfg.Status,
		}
	} else if cfg.Type == rector.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_WATCH {
		activity = &discord.Activity{
			Type:  discord.WatchingActivity,
			Name:  *cfg.Status,
			Flags: discord.JoinActivity,
		}
	}

	if activity != nil {
		if cfg.Url != nil {
			activity.URL = *cfg.Url
		}

		if err := b.dc.Gateway().Send(b.ctx, &gateway.UpdatePresenceCommand{
			Activities: []discord.Activity{*activity},
			Status:     discord.OnlineStatus,
		}); err != nil {
			return fmt.Errorf("failed to set bot presence. %w", err)
		}
	}

	b.logger.Info("bot presence has been set", zap.Int32("presence_type", int32(cfg.Type)))

	return nil
}
