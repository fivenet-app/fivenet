package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func init() {
	CommandsFactories["sync"] = NewSyncCommand
}

type SyncCommand struct {
	l *lang.I18n

	url string
}

func NewSyncCommand(router *cmdroute.Router, cfg *config.Config, p CommandParams) (api.CreateCommandData, error) {
	lEN := p.L.I18n("en")
	lDE := p.L.I18n("de")

	router.Add("sync", &SyncCommand{
		l:   p.L,
		url: cfg.HTTP.PublicURL,
	})

	return api.CreateCommandData{
			Type: discord.ChatInputCommand,
			Name: "sync",
			Description: lEN.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.fivenet.desc",
			}),
			DescriptionLocalizations: discord.StringLocales{
				discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "discord.commands.fivenet.desc",
				}),
			},
			DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
		},
		nil
}

func (c *SyncCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.I18n(string(cmd.Event.Locale))
	_ = localizer

	return &api.InteractionResponseData{}
}
