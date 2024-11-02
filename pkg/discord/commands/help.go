package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func init() {
	CommandsFactories["help"] = NewHelpCommand
}

var helpTopics = []string{
	"discord",
	"registration",
}

type HelpCommand struct {
	l *lang.I18n

	url string
}

func NewHelpCommand(router *cmdroute.Router, cfg *config.Config, p CommandParams) (api.CreateCommandData, error) {
	lEN := p.L.I18n("en")
	lDE := p.L.I18n("de")

	cmd := api.CreateCommandData{
		Type: discord.ChatInputCommand,
		Name: lEN.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "discord.commands.help.name",
		}),
		NameLocalizations: discord.StringLocales{
			discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.help.name",
			}),
		},

		Description: lEN.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "discord.commands.help.desc",
		}),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.help.desc",
			}),
		},

		Options:                  []discord.CommandOption{},
		DefaultMemberPermissions: discord.NewPermissions(discord.PermissionSendMessages),
	}

	choices := &discord.StringOption{
		OptionName: lEN.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "discord.commands.help.topic.name",
		}),
		OptionNameLocalizations: discord.StringLocales{
			discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.help.topic.name",
			}),
		},

		Description: lEN.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "discord.commands.help.topic.desc",
		}),
		DescriptionLocalizations: discord.StringLocales{
			discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.help.topic.desc",
			}),
		},

		Choices: []discord.StringChoice{},

		Required: true,
	}
	cmd.Options = append(cmd.Options, choices)

	for _, option := range helpTopics {
		choices.Choices = append(choices.Choices, discord.StringChoice{
			Name: lEN.MustLocalize(&i18n.LocalizeConfig{
				MessageID: fmt.Sprintf("discord.commands.help.%s.name", option),
			}),
			NameLocalizations: discord.StringLocales{
				discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
					MessageID: fmt.Sprintf("discord.commands.help.%s.name", option),
				}),
			},
			Value: option,
		})
	}

	router.Add("help", &HelpCommand{
		l:   p.L,
		url: cfg.HTTP.PublicURL,
	})

	return cmd, nil
}

func (c *HelpCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.I18n(string(cmd.Event.Locale))

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
				Title: localizer.MustLocalize(&i18n.LocalizeConfig{
					MessageID: messageId + ".title",
				}),
				Description: localizer.MustLocalize(&i18n.LocalizeConfig{
					MessageID: messageId + ".msg",
					TemplateData: map[string]string{
						"URL": c.url,
					},
				}),
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
