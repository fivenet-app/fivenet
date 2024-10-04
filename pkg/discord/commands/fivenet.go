package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func init() {
	CommandsFactories["fivenet"] = NewFivenetCommand
}

func NewFivenetCommand(cfg *config.Config, l *lang.I18n) (*discordgo.ApplicationCommand, CommandHandler, error) {
	lEN := l.I18n("en")
	lDE := l.I18n("de")

	url := cfg.HTTP.PublicURL
	return &discordgo.ApplicationCommand{
			Name: "fivenet",
			Description: lEN.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "discord.commands.fivenet.desc",
			}),
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.German: lDE.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "discord.commands.fivenet.desc",
				}),
			},
			Type:    discordgo.ChatApplicationCommand,
			GuildID: GlobalCommandGuildID,
		},
		NewHandleFivenetCommand(url, l),
		nil
}

func NewHandleFivenetCommand(url string, lang *lang.I18n) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		localizer := lang.I18n(string(i.Locale))
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:  discordgo.EmbedTypeLink,
						Title: "FiveNet",
						Description: localizer.MustLocalize(&i18n.LocalizeConfig{
							MessageID: "discord.commands.fivenet.summary",
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
