package commands

import (
	"context"

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

type HelpCommand struct {
	l *lang.I18n

	url string
}

func NewHelpCommand(cfg *config.Config, l *lang.I18n) (api.CreateCommandData, cmdroute.CommandHandler, error) {
	lEN := l.I18n("en")
	lDE := l.I18n("de")

	return api.CreateCommandData{
			Name: "help",
			Description: lEN.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.help.desc",
			}),
			DescriptionLocalizations: discord.StringLocales{
				discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "discord.commands.help.desc",
				}),
			},
			Type: discord.ChatInputCommand,
			Options: []discord.CommandOption{
				&discord.SubcommandOption{
					OptionName: "discord",
					Description: lEN.MustLocalize(&i18n.LocalizeConfig{
						MessageID: "discord.commands.help.discord.desc",
					}),
					DescriptionLocalizations: discord.StringLocales{
						discord.German: lDE.MustLocalize(&i18n.LocalizeConfig{
							MessageID: "discord.commands.help.discord.desc",
						}),
					},
				},
			},
		},
		&HelpCommand{
			url: cfg.HTTP.PublicURL,
		}, nil
}

func (c *HelpCommand) HandleCommand(ctx context.Context, cmd cmdroute.CommandData) *api.InteractionResponseData {
	localizer := c.l.I18n(string(cmd.Event.Locale))

	options := cmd.CommandInteractionOption.Options

	messageId := "discord.commands.help.empty"
	if len(options) > 0 {
		switch options[0].Name {
		case "discord":
			messageId = "discord.commands.help.discord"
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
						"URL": c.url + "/auth/account-info?tab=oauth2Connections#",
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
