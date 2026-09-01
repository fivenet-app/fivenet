package discordtypes

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type BotState interface {
	GetJobFromGuildID(guildId discord.GuildID) (string, bool)
	GetStatus() BotStatus
	IsUserConfigAdmin(ctx context.Context, userID discord.UserID) (bool, error)
	DebugUser(ctx context.Context, guildID discord.GuildID, userID discord.UserID) (*UserDebug, error)

	RunSync(guildID discord.GuildID) (bool, error)

	IsUserGuildAdmin(
		ctx context.Context,
		channelId discord.ChannelID,
		userId discord.UserID,
	) (bool, error)
}

type UserDebug struct {
	DiscordUserID discord.UserID
	MemberFound   bool
	DiscordLinked bool
	AccountFound  bool
	AccountActive bool

	UserID    int32
	Job       string
	JobGrade  int32
	PartOfJob bool

	Groups        []string
	MappedGroups  []string
	MissingRoles  []string
	ActualRoles   []string
	ExpectedRoles []string

	GroupSyncEnabled bool
	UserInfoEnabled  bool
	SyncRunning      bool
	LastSync         time.Time
}

// BotStatus contains the read-only state exposed by the Discord bot status command.
type BotStatus struct {
	Guilds       []GuildStatus
	SyncInterval time.Duration
}

type GuildStatus struct {
	GuildID  discord.GuildID
	Job      string
	Running  bool
	LastSync time.Time
}
