package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2025/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type FivenetCommand struct {
	l *lang.I18n

	url string
}

func NewFivenetCommand(p CommandParams) (Command, error) {
	return &FivenetCommand{
		l:   p.L,
		url: p.Cfg.HTTP.PublicURL,
	}, nil
}

func (c *FivenetCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	lEN := c.l.I18n("en")
	lDE := c.l.I18n("de")

	router.Add("fivenet", c)

	return api.CreateCommandData{
		Type: discord.ChatInputCommand,
		Name: "fivenet",
		Description: lEN.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "discord.commands.fivenet.desc",
		}),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.fivenet.desc",
			}),
		},
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
	}
}

func (c *FivenetCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
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
