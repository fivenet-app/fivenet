package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

type StatusCommand struct {
	l i18n.Ii18n
	b discordtypes.BotState
}

func NewStatusCommand(p CommandParams) (Command, error) {
	if p.BotState == nil {
		return nil, nil
	}
	return &StatusCommand{l: p.I18n, b: p.BotState}, nil
}

func (c *StatusCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	tr := newCommandLocalizer(c.l, "discord.commands.status")
	router.Add("status", c)

	return api.CreateCommandData{
		Type:                     discord.ChatInputCommand,
		Name:                     "status",
		Description:              tr.text("desc"),
		DescriptionLocalizations: tr.localizations("desc"),
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionAdministrator),
	}
}

func (c *StatusCommand) HandleCommand(
	ctx context.Context,
	cmd cmdroute.CommandData,
) *api.InteractionResponseData {
	t := c.l.Translator(string(cmd.Event.Locale))
	resp := &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{{
			Type:      discord.NormalEmbed,
			Provider:  &discord.EmbedProvider{Name: version.ProjectName},
			Thumbnail: embeds.EmbedThumbnailLogo,
			Footer:    embeds.EmbedFooterMadeBy,
			Color:     embeds.ColorInfo,
		}},
	}
	embed := &(*resp.Embeds)[0]

	if cmd.Event.GuildID == discord.NullGuildID || cmd.Event.Member == nil ||
		cmd.Event.Channel == nil {
		embed.Title = t("discord.commands.status.results.wrong_discord.title", nil)
		embed.Description = t("discord.commands.status.results.wrong_discord.desc", nil)
		return resp
	}

	admin, err := c.b.IsUserGuildAdmin(ctx, cmd.Event.ChannelID, cmd.Event.Member.User.ID)
	if err != nil || !admin {
		embed.Title = t("discord.commands.status.results.permission_denied.title", nil)
		embed.Description = t("discord.commands.status.results.permission_denied.desc", nil)
		embed.Color = embeds.ColorError
		return resp
	}

	status, ok := c.b.GetStatusForGuild(cmd.Event.GuildID)
	if !ok {
		embed.Title = t("discord.commands.status.results.wrong_discord.title", nil)
		embed.Description = t("discord.commands.status.results.wrong_discord.desc", nil)
		return resp
	}
	lines := []string{
		t(
			"discord.commands.status.results.interval",
			map[string]any{"interval": status.SyncInterval.String()},
		),
	}
	for _, guild := range status.Guilds {
		state := t("discord.commands.status.results.idle", nil)
		if guild.Running {
			state = t("discord.commands.status.results.running", nil)
		}
		lastSync := t("discord.commands.status.results.never", nil)
		if !guild.LastSync.IsZero() {
			lastSync = fmt.Sprintf("<t:%d:R>", guild.LastSync.Unix())
		}
		lines = append(
			lines,
			fmt.Sprintf(
				"**%s** (`%d`) — **%s:** %s, **%s:** %s",
				guild.Job,
				guild.GuildID,
				t("discord.commands.status.results.state", nil),
				state,
				t("discord.commands.status.results.last_sync", nil),
				lastSync,
			),
		)
	}
	if len(status.Guilds) == 0 {
		lines = append(lines, t("discord.commands.status.results.no_guilds", nil))
	}
	embed.Title = t("discord.commands.status.results.title", nil)
	embed.Description = strings.Join(lines, "\n")
	return resp
}
