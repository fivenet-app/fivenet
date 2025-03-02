package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/discord/types"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type SyncCommand struct {
	l *lang.I18n
	b types.BotState

	url string
}

func NewSyncCommand(p CommandParams) (Command, error) {
	return &SyncCommand{
			l:   p.L,
			b:   p.BotState,
			url: p.Cfg.HTTP.PublicURL,
		},
		nil
}

func (c *SyncCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	lEN := c.l.I18n("en")
	lDE := c.l.I18n("de")

	router.Add("sync", c)

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
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionAdministrator),
	}
}

func (c *SyncCommand) getBaseResponse() *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{
			{
				Type:  discord.LinkEmbed,
				Color: embeds.ColorError,
				Provider: &discord.EmbedProvider{
					Name: "FiveNet",
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

func (c *SyncCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.I18n(string(cmd.Event.Locale))
	resp := c.getBaseResponse()

	// Make sure command is used on a guild's channel
	if cmd.Event.GuildID == discord.NullGuildID || cmd.Event.Member == nil || cmd.Event.Channel == nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.wrong_discord.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.wrong_discord.desc"})
		(*resp.Embeds)[0].Color = embeds.ColorInfo
		return resp
	}

	// Check if user has admin perms to guild server
	channelAdmin, err := c.b.IsUserGuildAdmin(ctx, cmd.Event.ChannelID, cmd.Event.Member.User.ID)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.permission_denied.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.permission_denied.desc"})
		return resp
	}
	if !channelAdmin {
		return resp
	}

	// Try to run sync
	running, err := c.b.RunSync(cmd.Event.GuildID)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.start_error.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.start_error.desc"})
		return resp
	}

	(*resp.Embeds)[0].Color = embeds.ColorInfo
	if running {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.already_running.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.already_running.desc"})
		(*resp.Embeds)[0].Color = embeds.ColorWarn
	} else {
		(*resp.Embeds)[0].Title = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.started.title"})
		(*resp.Embeds)[0].Description = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "discord.commands.sync.results.started.desc"})
		(*resp.Embeds)[0].Color = embeds.ColorSuccess
	}

	return resp
}
