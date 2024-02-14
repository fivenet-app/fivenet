package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/galexrt/fivenet/pkg/config"
)

func init() {
	CommandsFactories["fivenet"] = NewFivenetCommand
}

func NewFivenetCommand(cfg *config.Config) (*discordgo.ApplicationCommand, CommandHandler, error) {
	url := cfg.HTTP.PublicURL
	return &discordgo.ApplicationCommand{
			Name:        "fivenet",
			Description: "FiveNet Website URL",
			Type:        discordgo.ChatApplicationCommand,
			GuildID:     GlobalCommandGuildID,
		},
		NewHandleFivenetCommand(url),
		nil
}

func NewHandleFivenetCommand(url string) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeLink,
						Title:       "FiveNet",
						Description: "Link zur FiveNet Web App.",
						URL:         url,
						Provider: &discordgo.MessageEmbedProvider{
							Name: "FiveNet",
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
							Width:  128,
							Height: 128,
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Made by Galexrt",
							IconURL: "https://cdn.discordapp.com/avatars/290472392084422658/58e38b558fb3a54b4864584e7b7297f6.png",
						},
					},
				},
			},
		})
	}
}
