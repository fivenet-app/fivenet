package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	lang "github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
)

var helpTopics = []string{
	"registration",
	"browser",
	"discord",
}

type HelpCommand struct {
	l *lang.I18n

	url string
}

func NewHelpCommand(p CommandParams) (Command, error) {
	return &HelpCommand{
		l:   p.L,
		url: p.Cfg.HTTP.PublicURL,
	}, nil
}

func (c *HelpCommand) RegisterCommand(router *cmdroute.Router) api.CreateCommandData {
	lEN := c.l.Translator("en")
	lDE := c.l.Translator("de")

	cmdData := api.CreateCommandData{
		Type: discord.ChatInputCommand,
		Name: lEN("discord.commands.help.name", nil),
		NameLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.help.name", nil),
		},

		Description: lEN("discord.commands.help.desc", nil),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.help.desc", nil),
		},

		Options:                  []discord.CommandOption{},
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
	}

	choices := &discord.StringOption{
		OptionName: lEN("discord.commands.help.topic.name", nil),
		OptionNameLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.help.topic.name", nil),
		},

		Description: lEN("discord.commands.help.topic.desc", nil),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE("discord.commands.help.topic.desc", nil),
		},

		Choices: []discord.StringChoice{},

		Required: false,
	}
	cmdData.Options = append(cmdData.Options, choices)

	for _, option := range helpTopics {
		choices.Choices = append(choices.Choices, discord.StringChoice{
			Name: lEN(fmt.Sprintf("discord.commands.help.%s.name", option), nil),
			NameLocalizations: discord.StringLocales{
				discord.German: lDE(fmt.Sprintf("discord.commands.help.%s.name", option), nil),
			},
			Value: option,
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
	parts := strings.Split(string(data.ComponentInteraction.ID()), "_")
	return &api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: c.getHelpTopicResponse(data.Event.Locale, parts[1]),
	}
}

func (c *HelpCommand) getHelpTopicResponse(
	locale discord.Language,
	topic string,
) *api.InteractionResponseData {
	localizer := c.l.Translator(string(locale))

	var title string
	var desc string

	if len(topic) > 0 && slices.Contains(helpTopics, topic) {
		messageId := "discord.commands.help." + topic
		title = localizer(messageId+".title", nil)
		desc = localizer(messageId+".msg", map[string]any{"url": c.url})
	} else {
		messageId := "discord.commands.help.empty"
		title = localizer(messageId+".title", nil)
		desc = localizer(messageId+".msg", nil)
	}

	helpOpts := []discord.InteractiveComponent{}
	for _, opt := range helpTopics {
		btnStyle := discord.SecondaryButtonStyle()
		if opt == topic {
			btnStyle = discord.PrimaryButtonStyle()
		}

		helpOpts = append(helpOpts, &discord.ButtonComponent{
			Label:    localizer(fmt.Sprintf("discord.commands.help.%s.title", opt), nil),
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
					Name: "FiveNet",
					URL:  c.url,
				},
				Color:  embeds.ColorInfo,
				Footer: embeds.EmbedFooterFiveNet,
			},
		},
		Components: discord.ComponentsPtr(&actionRow),
	}
}
