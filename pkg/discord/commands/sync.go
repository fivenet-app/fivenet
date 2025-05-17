package commands

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	lang "github.com/fivenet-app/fivenet/v2025/i18n"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/types"
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
	lEN := c.l.Translator("en")
	lDE := c.l.Translator("de")

	router.Add("sync", c)

	return api.CreateCommandData{
		Type:        discord.ChatInputCommand,
		Name:        "sync",
		Description: lEN("discord.commands.sync.desc", nil),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.sync.desc", nil),
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
	localizer := c.l.Translator(string(cmd.Event.Locale))
	resp := c.getBaseResponse()

	// Make sure command is used on a guild's channel
	if cmd.Event.GuildID == discord.NullGuildID || cmd.Event.Member == nil || cmd.Event.Channel == nil {
		(*resp.Embeds)[0].Title = localizer("discord.commands.sync.results.wrong_discord.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.sync.results.wrong_discord.desc", nil)
		(*resp.Embeds)[0].Color = embeds.ColorInfo
		return resp
	}

	// Check if user has admin perms to guild server
	channelAdmin, err := c.b.IsUserGuildAdmin(ctx, cmd.Event.ChannelID, cmd.Event.Member.User.ID)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer("discord.commands.sync.results.permission_denied.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.sync.results.permission_denied.desc", nil)
		return resp
	}
	if !channelAdmin {
		return resp
	}

	// Try to run sync
	running, err := c.b.RunSync(cmd.Event.GuildID)
	if err != nil {
		(*resp.Embeds)[0].Title = localizer("discord.commands.sync.results.start_error.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.sync.results.start_error.desc", nil)
		return resp
	}

	(*resp.Embeds)[0].Color = embeds.ColorInfo
	if running {
		(*resp.Embeds)[0].Title = localizer("discord.commands.sync.results.already_running.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.sync.results.already_running.desc", nil)
		(*resp.Embeds)[0].Color = embeds.ColorWarn
	} else {
		(*resp.Embeds)[0].Title = localizer("discord.commands.sync.results.started.title", nil)
		(*resp.Embeds)[0].Description = localizer("discord.commands.sync.results.started.desc", nil)
		(*resp.Embeds)[0].Color = embeds.ColorSuccess
	}

	return resp
}
