package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	discordtypes "github.com/fivenet-app/fivenet/v2026/pkg/discord/types"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

type DebugUserCommand struct {
	l i18n.Ii18n
	b discordtypes.BotState
}

func NewDebugUserCommand(p CommandParams) (Command, error) {
	if p.BotState == nil {
		return nil, nil
	}
	return &DebugUserCommand{l: p.I18n, b: p.BotState}, nil
}

func (c *DebugUserCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	tr := newCommandLocalizer(c.l, "discord.commands.debug-user")
	router.Add("debug-user", c)
	return api.CreateCommandData{
		Type:                     discord.ChatInputCommand,
		Name:                     "debug-user",
		Description:              tr.text("desc"),
		DescriptionLocalizations: tr.localizations("desc"),
		Options: discord.CommandOptions{&discord.UserOption{
			OptionName:               "user",
			OptionNameLocalizations:  tr.localizations("options.user.name"),
			Description:              tr.text("options.user.desc"),
			DescriptionLocalizations: tr.localizations("options.user.desc"),
			Required:                 true,
		}},
		// Config admins may not have Discord's Administrator permission, but it is
		// required to use this command, so we set it here to avoid confusion.
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionAdministrator),
	}
}

func (c *DebugUserCommand) HandleCommand(
	ctx context.Context,
	cmd cmdroute.CommandData,
) *api.InteractionResponseData {
	t := c.l.Translator(string(cmd.Event.Locale))
	resp := &api.InteractionResponseData{Flags: discord.EphemeralMessage, Embeds: &[]discord.Embed{
		{
			Type:      discord.NormalEmbed,
			Provider:  &discord.EmbedProvider{Name: version.ProjectName},
			Thumbnail: embeds.EmbedThumbnailLogo,
			Footer:    embeds.EmbedFooterMadeBy,
			Color:     embeds.ColorInfo,
		},
	}}
	embed := &(*resp.Embeds)[0]
	if cmd.Event.GuildID == discord.NullGuildID || cmd.Event.Member == nil ||
		cmd.Event.Channel == nil {
		embed.Title = t("discord.commands.debug-user.results.wrong_discord.title", nil)
		embed.Description = t("discord.commands.debug-user.results.wrong_discord.desc", nil)
		return resp
	}
	dcAdmin, err := c.b.IsUserGuildAdmin(ctx, cmd.Event.ChannelID, cmd.Event.Member.User.ID)
	if err != nil {
		dcAdmin = false
	}
	configAdmin, err := c.b.IsUserConfigAdmin(ctx, cmd.Event.Member.User.ID)
	if err != nil {
		configAdmin = false
	}
	if !dcAdmin && !configAdmin {
		embed.Title = t("discord.commands.debug-user.results.permission_denied.title", nil)
		embed.Description = t("discord.commands.debug-user.results.permission_denied.desc", nil)
		embed.Color = embeds.ColorError
		return resp
	}

	opt := cmd.Options.Find("user")
	if opt.Name == "" {
		embed.Title = t("discord.commands.debug-user.results.invalid_user.title", nil)
		embed.Description = t("discord.commands.debug-user.results.invalid_user.desc", nil)
		return resp
	}
	userID, err := strconv.ParseUint(strings.Trim(opt.Value.String(), `"`), 10, 64)
	if err != nil {
		embed.Title = t("discord.commands.debug-user.results.invalid_user.title", nil)
		embed.Description = t("discord.commands.debug-user.results.invalid_user.desc", nil)
		return resp
	}

	debug, err := c.b.DebugUser(ctx, cmd.Event.GuildID, discord.UserID(userID))
	if err != nil {
		embed.Title = t("discord.commands.debug-user.results.failed.title", nil)
		embed.Description = t(
			"discord.commands.debug-user.results.failed.desc",
			map[string]any{"error": err.Error()},
		)
		embed.Color = embeds.ColorError
		return resp
	}
	username := strconv.FormatUint(userID, 10)
	if cmd.Data != nil {
		if user, ok := cmd.Data.Resolved.Users[discord.UserID(userID)]; ok {
			if name := user.DisplayOrUsername(); name != "" {
				username = name
			}
		}
	}
	resp.Content = option.NewNullableString(fmt.Sprintf("<@%d>", userID))
	embed.Title = t(
		"discord.commands.debug-user.results.title",
		map[string]any{"user": username},
	)
	lines := []string{
		fmt.Sprintf(
			"**%s:** %s",
			t("discord.commands.debug-user.results.member", nil),
			yesNo(t, debug.MemberFound),
		),
		fmt.Sprintf(
			"**%s:** %s",
			t("discord.commands.debug-user.results.part_of_job", nil),
			yesNo(t, debug.PartOfJob),
		),
		fmt.Sprintf(
			"**%s:** %s",
			t("discord.commands.debug-user.results.group_sync", nil),
			yesNo(t, len(debug.MappedGroups) > 0),
		),
		fmt.Sprintf(
			"**%s:** %s",
			t("discord.commands.debug-user.results.sync_status", nil),
			syncState(t, debug.SyncRunning),
		),
	}
	causes := []string{}
	if !debug.MemberFound {
		causes = append(causes, t("discord.commands.debug-user.results.causes.not_member", nil))
	}
	if !debug.DiscordLinked {
		causes = append(causes, t("discord.commands.debug-user.results.causes.not_linked", nil))
	}
	if !debug.AccountActive {
		causes = append(
			causes,
			t("discord.commands.debug-user.results.causes.account_inactive", nil),
		)
	}
	if !debug.PartOfJob {
		causes = append(causes, t("discord.commands.debug-user.results.causes.not_job", nil))
	}
	if debug.GroupSyncEnabled && len(debug.MappedGroups) == 0 {
		causes = append(causes, t("discord.commands.debug-user.results.causes.no_group", nil))
	}
	if len(debug.MissingRoles) > 0 {
		causes = append(causes, t("discord.commands.debug-user.results.causes.missing_role", nil))
	}
	if len(causes) == 0 {
		causes = append(causes, t("discord.commands.debug-user.results.causes.none", nil))
	}
	lines = append(
		lines,
		"",
		"**"+t("discord.commands.debug-user.results.likely_causes", nil)+":**",
		"- "+strings.Join(causes, "\n- "),
	)
	if configAdmin {
		lines = append(lines, "", "**"+t("discord.commands.debug-user.results.details", nil)+":**")
		lines = append(
			lines,
			fmt.Sprintf(
				"**%s:** %s",
				t("discord.commands.debug-user.results.linked", nil),
				yesNo(t, debug.DiscordLinked),
			),
			fmt.Sprintf(
				"**%s:** %s",
				t("discord.commands.debug-user.results.account", nil),
				yesNo(t, debug.AccountFound),
			),
			fmt.Sprintf("**%s:** %s", t("discord.commands.debug-user.results.job", nil), debug.Job),
			fmt.Sprintf(
				"**%s:** %d",
				t("discord.commands.debug-user.results.grade", nil),
				debug.JobGrade,
			),
			fmt.Sprintf(
				"**%s:** %s",
				t("discord.commands.debug-user.results.groups", nil),
				listOrNone(t, debug.Groups),
			),
			fmt.Sprintf(
				"**%s:** %s",
				t("discord.commands.debug-user.results.mapped_groups", nil),
				listOrNone(t, debug.MappedGroups),
			),
			fmt.Sprintf(
				"**%s:** %s",
				t("discord.commands.debug-user.results.missing_roles", nil),
				listOrNone(t, debug.MissingRoles),
			),
		)
	}
	embed.Description = strings.Join(lines, "\n")
	return resp
}

func yesNo(t i18n.TFunc, value bool) string {
	if value {
		return t("discord.commands.debug-user.results.yes", nil)
	}
	return t("discord.commands.debug-user.results.no", nil)
}

func syncState(t i18n.TFunc, running bool) string {
	if running {
		return t("discord.commands.debug-user.results.running", nil)
	}
	return t("discord.commands.debug-user.results.idle", nil)
}

func listOrNone(t i18n.TFunc, values []string) string {
	if len(values) == 0 {
		return t("discord.commands.debug-user.results.none", nil)
	}
	return "`" + strings.Join(values, "`, `") + "`"
}
