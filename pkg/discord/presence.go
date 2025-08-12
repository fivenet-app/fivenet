package discord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/settings"
	"go.uber.org/zap"
)

func (b *Bot) setBotPresence(ctx context.Context, cfg *settings.DiscordBotPresence) error {
	var activity *discord.Activity
	switch cfg.GetType() {
	case settings.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_GAME:
		activity = &discord.Activity{
			Type: discord.GameActivity,
			Name: cfg.GetStatus(),
		}
	case settings.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_LISTENING:
		activity = &discord.Activity{
			Type: discord.ListeningActivity,
			Name: cfg.GetStatus(),
		}
	case settings.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_STREAMING:
		activity = &discord.Activity{
			Type: discord.StreamingActivity,
			Name: cfg.GetStatus(),
		}
	case settings.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_WATCH:
		fallthrough
	case settings.DiscordBotPresenceType_DISCORD_BOT_PRESENCE_TYPE_UNSPECIFIED:
		fallthrough
	default:
		activity = &discord.Activity{
			Type:  discord.WatchingActivity,
			Name:  cfg.GetStatus(),
			Flags: discord.JoinActivity,
		}
	}

	if cfg.Url != nil {
		activity.URL = cfg.GetUrl()
	}

	if err := b.dc.Gateway().Send(ctx, &gateway.UpdatePresenceCommand{
		Activities: []discord.Activity{*activity},
		Status:     discord.OnlineStatus,
	}); err != nil {
		return fmt.Errorf("failed to set bot presence. %w", err)
	}

	b.logger.Info("bot presence has been set", zap.Int32("presence_type", int32(cfg.GetType())))

	return nil
}
