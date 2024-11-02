package types

import "github.com/diamondburned/arikawa/v3/discord"

type BotState interface {
	GetJobFromGuildID(guildId discord.GuildID) (string, bool)
}
