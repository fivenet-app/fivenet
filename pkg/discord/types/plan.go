package types

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type UserProcessorHandler func(ctx context.Context, guildId string, member *discordgo.Member, u *User) (*User, []*discordgo.MessageEmbed, error)

type Plan struct {
	GuildID string `yaml:"guildId"`
	DryRun  bool   `yaml:"dryRun"`

	UserProcessors []UserProcessorHandler `yaml:"-"`

	Roles *PlanRoles `yaml:"roleChanges"`
	Users []*User    `yaml:"userChanges"`
}

type PlanRoles struct {
	ToCreate Roles `yaml:"toCreate,omitempty"`
	ToUpdate Roles `yaml:"toUpdate,omitempty"`
}

func NewPlan(guildId string, dryRun bool) *Plan {
	return &Plan{
		GuildID: guildId,
		DryRun:  dryRun,

		Roles: &PlanRoles{},
		Users: []*User{},
	}
}
