package types

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
)

type BotState interface {
	GetJobFromGuildID(guildId discord.GuildID) (string, bool)

	RunSync(guildID discord.GuildID) (bool, error)

	IsUserGuildAdmin(ctx context.Context, channelId discord.ChannelID, userId discord.UserID) (bool, error)
}
