package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"go.uber.org/fx"
)

type Command interface {
	RegisterCommand(router *cmdroute.Router) api.CreateCommandData
	HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData
}

// AsCommand annotates the given constructor to state that
// it provides a Discord command to the "discordcommands" group.
func AsCommand(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Command)),
		fx.ResultTags(`group:"discordcommands"`),
	)
}
