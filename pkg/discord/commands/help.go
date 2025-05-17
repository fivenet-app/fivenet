package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	lang "github.com/fivenet-app/fivenet/v2025/i18n"
	"github.com/fivenet-app/fivenet/v2025/pkg/discord/embeds"
)

var helpTopics = []string{
	"discord",
	"registration",
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

		Required: true,
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

	return cmdData
}

func (c *HelpCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.Translator(string(cmd.Event.Locale))

	messageId := "discord.commands.help.empty"

	options := cmd.CommandInteractionOption.Options
	if len(options) > 0 {
		option := strings.ReplaceAll(strings.ToLower(options[0].Value.String()), "\"", "")

		if slices.Contains(helpTopics, option) {
			messageId = "discord.commands.help." + option
		}
	}

	return &api.InteractionResponseData{
		Flags: discord.EphemeralMessage,
		Embeds: &[]discord.Embed{
			{
				Title: localizer(messageId+".title", nil),
				Description: localizer(messageId+".msg",
					map[string]any{"url": c.url}),
				URL: c.url,
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
