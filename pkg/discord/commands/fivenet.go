package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
)

type FivenetCommand struct {
	l i18n.Ii18n

	url string
}

func NewFivenetCommand(p CommandParams) (Command, error) {
	return &FivenetCommand{
		l:   p.I18n,
		url: p.Cfg.HTTP.PublicURL,
	}, nil
}

func (c *FivenetCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	tr := newCommandLocalizer(c.l, "discord.commands.fivenet")

	router.Add("fivenet", c)

	return api.CreateCommandData{
		Type:                     discord.ChatInputCommand,
		Name:                     "fivenet",
		Description:              tr.text("desc"),
		DescriptionLocalizations: tr.localizations("desc"),
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
	}
}

func (c *FivenetCommand) HandleCommand(
	ctx context.Context,
	cmd cmdroute.CommandData,
) *api.InteractionResponseData {
	t := c.l.Translator(string(cmd.Event.Locale))

	return &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{
			{
				Type:        discord.LinkEmbed,
				Title:       "FiveNet",
				Description: t("discord.commands.fivenet.summary", nil),
				URL:         c.url,
				Thumbnail:   embeds.EmbedThumbnailLogo,
				Provider: &discord.EmbedProvider{
					Name: "FiveNet",
					URL:  c.url,
				},
				Color:  embeds.ColorInfo,
				Footer: embeds.EmbedFooterMadeBy,
			},
		},
		Components: discord.ComponentsPtr(
			&discord.ActionRowComponent{
				&discord.ButtonComponent{
					Label: t("discord.commands.fivenet.open_link", nil),
					Style: discord.LinkButtonStyle(c.url),
				},
			},
		),
	}
}
