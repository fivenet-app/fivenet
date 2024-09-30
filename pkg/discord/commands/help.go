package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func init() {
	CommandsFactories["help"] = NewFivenetHelpCommand
}

func NewFivenetHelpCommand(cfg *config.Config, l *lang.I18n) (*discordgo.ApplicationCommand, CommandHandler, error) {
	lEN := l.I18n("en")
	lDE := l.I18n("de")

	url := cfg.HTTP.PublicURL
	return &discordgo.ApplicationCommand{
			Name: "help",
			Description: lEN.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.help.desc",
			}),
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.German: lDE.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "discord.commands.help.desc",
				}),
			},
			Type:    discordgo.ChatApplicationCommand,
			GuildID: GlobalCommandGuildID,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "discord",
					Type: discordgo.ApplicationCommandOptionSubCommand,
					Description: lEN.MustLocalize(&i18n.LocalizeConfig{
						MessageID: "discord.commands.help.discord.desc",
					}),
					DescriptionLocalizations: map[discordgo.Locale]string{
						discordgo.German: lDE.MustLocalize(&i18n.LocalizeConfig{
							MessageID: "discord.commands.help.discord.desc",
						}),
					},
				},
			},
		},
		NewHandleHelpCommand(url, l),
		nil
}

func NewHandleHelpCommand(url string, lang *lang.I18n) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		localizer := lang.I18n(string(i.Locale))

		options := i.ApplicationCommandData().Options

		if len(options) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:  discordgo.EmbedTypeLink,
							Title: "FiveNet",
							Description: localizer.MustLocalize(&i18n.LocalizeConfig{
								MessageID: "discord.commands.help.empty",
							}),
							URL: url,
							Provider: &discordgo.MessageEmbedProvider{
								Name: "FiveNet",
								URL:  url,
							},
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
								Width:  128,
								Height: 128,
							},
							Footer: embeds.EmbedFooterMadeBy,
						},
					},
				},
			})
			return
		}

		messageId := "discord.commands.help.empty"
		switch options[0].Name {
		case "discord":
			messageId = "discord.commands.help.discord"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Type: discordgo.EmbedTypeLink,
						Title: localizer.MustLocalize(&i18n.LocalizeConfig{
							MessageID: messageId + ".title",
						}),
						Description: localizer.MustLocalize(&i18n.LocalizeConfig{
							MessageID: messageId + ".msg",
							TemplateData: map[string]string{
								"URL": url + "/auth/account-info?tab=oauth2Connections#",
							},
						}),
						URL: url,
						Provider: &discordgo.MessageEmbedProvider{
							Name: "FiveNet",
							URL:  url,
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
							Width:  128,
							Height: 128,
						},
						Footer: embeds.EmbedFooterMadeBy,
					},
				},
			},
		})
	}
}
