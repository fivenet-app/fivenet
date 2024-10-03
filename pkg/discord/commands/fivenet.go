package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func init() {
	CommandsFactories["fivenet"] = NewFivenetCommand
}

type FiveNetCommand struct {
	l *lang.I18n

	url string
}

func NewFivenetCommand(cfg *config.Config, l *lang.I18n) (api.CreateCommandData, cmdroute.CommandHandler, error) {
	lEN := l.I18n("en")
	lDE := l.I18n("de")

	return api.CreateCommandData{
			Name: "fivenet",
			Description: lEN.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.fivenet.desc",
			}),
			DescriptionLocalizations: discord.StringLocales{
				discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "discord.commands.fivenet.desc",
				}),
			},
			Type: discord.ChatInputCommand,
		},
		&FiveNetCommand{
			l:   l,
			url: cfg.HTTP.PublicURL,
		},
		nil
}

func (c *FiveNetCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.I18n(string(cmd.Event.Locale))

	return &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{
			{
				Title: "FiveNet",
				Description: localizer.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "discord.commands.fivenet.summary",
				}),
				URL:  c.url,
				Type: discord.LinkEmbed,
				Provider: &discord.EmbedProvider{
					Name: "FiveNet",
					URL:  c.url,
				},
				Thumbnail: &discord.EmbedThumbnail{
					URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
					Width:  128,
					Height: 128,
				},
				Footer: embeds.EmbedFooterMadeBy,
			},
		},
	}
}
