package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

var helpTopics = []string{
	"registration",
	"browser",
	"discord",
}

type HelpCommand struct {
	l i18n.Ii18n

	url string
}

func NewHelpCommand(p CommandParams) (Command, error) {
	return &HelpCommand{
		l:   p.I18n,
		url: p.Cfg.HTTP.PublicURL,
	}, nil
}

func (c *HelpCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	tr := newCommandLocalizer(c.l, "discord.commands.help")

	cmdData := api.CreateCommandData{
		Type:              discord.ChatInputCommand,
		Name:              tr.text("name"),
		NameLocalizations: tr.localizations("name"),

		Description:              tr.text("desc"),
		DescriptionLocalizations: tr.localizations("desc"),

		Options:                  []discord.CommandOption{},
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
	}

	choices := &discord.StringOption{
		OptionName:              tr.text("topic.name"),
		OptionNameLocalizations: tr.localizations("topic.name"),

		Description:              tr.text("topic.desc"),
		DescriptionLocalizations: tr.localizations("topic.desc"),

		Choices: []discord.StringChoice{},

		Required: false,
	}
	cmdData.Options = append(cmdData.Options, choices)

	for _, option := range helpTopics {
		choices.Choices = append(choices.Choices, discord.StringChoice{
			Name:              tr.text(option + ".name"),
			NameLocalizations: tr.localizations(option + ".name"),
			Value:             option,
		})
	}

	router.Add("help", c)

	for _, opt := range helpTopics {
		router.AddComponent(fmt.Sprintf("help_%s", opt), c)
	}

	return cmdData
}

func (c *HelpCommand) HandleCommand(
	ctx context.Context,
	cmd cmdroute.CommandData,
) *api.InteractionResponseData {
	options := cmd.Options
	var topic string
	if len(options) > 0 {
		topic = strings.ReplaceAll(strings.ToLower(options[0].Value.String()), "\"", "")
	}

	return c.getHelpTopicResponse(cmd.Event.Locale, topic)
}

func (c *HelpCommand) HandleComponent(
	ctx context.Context,
	data cmdroute.ComponentData,
) *api.InteractionResponse {
	topic, found := strings.CutPrefix(string(data.ComponentInteraction.ID()), "help_")
	if !found {
		topic = ""
	}

	return &api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: c.getHelpTopicResponse(data.Event.Locale, topic),
	}
}

func (c *HelpCommand) getHelpTopicResponse(
	locale discord.Language,
	topic string,
) *api.InteractionResponseData {
	t := c.l.Translator(string(locale))

	var title string
	var desc string

	if len(topic) > 0 && slices.Contains(helpTopics, topic) {
		messageId := "discord.commands.help." + topic
		title = t(messageId+".title", nil)
		desc = t(messageId+".msg", map[string]any{"url": c.url})
	} else {
		messageId := "discord.commands.help.empty"
		title = t(messageId+".title", nil)
		desc = t(messageId+".msg", nil)
	}

	helpOpts := []discord.InteractiveComponent{}
	for _, opt := range helpTopics {
		btnStyle := discord.SecondaryButtonStyle()
		if opt == topic {
			btnStyle = discord.PrimaryButtonStyle()
		}

		helpOpts = append(helpOpts, &discord.ButtonComponent{
			Label:    t(fmt.Sprintf("discord.commands.help.%s.title", opt), nil),
			CustomID: discord.ComponentID(fmt.Sprintf("help_%s", opt)),
			Style:    btnStyle,
		})
	}
	actionRow := discord.ActionRowComponent(helpOpts)

	return &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{
			{
				Title:       title,
				Description: desc,
				Provider: &discord.EmbedProvider{
					Name: version.Project,
					URL:  c.url,
				},
				Color:  embeds.ColorInfo,
				Footer: embeds.EmbedFooterFiveNet,
			},
		},
		Components: discord.ComponentsPtr(&actionRow),
	}
}
