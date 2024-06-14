package modules

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
)

type QualificationsSync struct {
	*BaseModule
}

func init() {
	Modules["qualifications"] = NewQualifications
}

func NewQualifications(base *BaseModule) (Module, error) {
	return &QualificationsSync{
		BaseModule: base,
	}, nil
}

func (g *QualificationsSync) Plan(ctx context.Context) (*types.State, []*discordgo.MessageEmbed, error) {

	// TODO

	return &types.State{
		Roles: types.Roles{},
		Users: types.Users{},
	}, nil, nil
}
