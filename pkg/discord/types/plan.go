package types

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
)

type UserProcessorHandler func(ctx context.Context, guildId discord.GuildID, member discord.Member, u *User) ([]discord.Embed, error)

type Plan struct {
	GuildID discord.GuildID `yaml:"guildId"`
	DryRun  bool            `yaml:"dryRun"`

	UserProcessors []UserProcessorHandler `yaml:"-"`

	Roles *PlanRoles `yaml:"roleChanges"`
	Users []*User    `yaml:"userChanges"`
}

type PlanRoles struct {
	ToCreate Roles `yaml:"toCreate,omitempty"`
	ToUpdate Roles `yaml:"toUpdate,omitempty"`
}

func NewPlan(guildId discord.GuildID, dryRun bool) *Plan {
	return &Plan{
		GuildID: guildId,
		DryRun:  dryRun,

		Roles: &PlanRoles{},
		Users: []*User{},
	}
}
