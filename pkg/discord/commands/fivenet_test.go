package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/mavolin/dismock/v3/pkg/dismock"
	"github.com/stretchr/testify/require"
)

func TestNewHandleFivenetCommand(t *testing.T) {
	l, err := lang.New()
	require.NoError(t, err)

	m := dismock.New(t)

	s, _ := discordgo.New("Bot abc")
	s.StateEnabled = false
	s.Client = m.Client

	cfg, err := config.LoadTestConfig()
	require.NoError(t, err)
	url := "https://example.fivenet.app/"
	cfg.HTTP.PublicURL = url

	_, handler, err := NewFivenetCommand(cfg, l)
	require.NoError(t, err)

	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:      discord.NullInteractionID.String(),
			AppID:   discord.NullAppID.String(),
			GuildID: discord.NullGuildID.String(),
			Token:   "fivenet",
			Type:    discordgo.InteractionApplicationCommand,

			Locale: "en",
		},
	}

	data := &api.InteractionResponseData{
		Flags:   discord.EphemeralMessage,
		Content: option.NewNullableString(""),
		Embeds: &[]discord.Embed{
			{
				Type:        discord.LinkEmbed,
				Title:       "FiveNet",
				Description: "FiveNet is also available in your browser! Link to the FiveNet web app.",
				URL:         cfg.HTTP.PublicURL,
				Thumbnail: &discord.EmbedThumbnail{
					URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
					Height: 128,
					Width:  128,
				},
				Footer: &discord.EmbedFooter{
					Text: embeds.EmbedFooterMadeBy.Text,
					Icon: embeds.EmbedFooterMadeBy.IconURL,
				},
				Provider: &discord.EmbedProvider{
					Name: "FiveNet",
					URL:  url,
				},
			},
		},
	}

	// English
	m.RespondInteraction(discord.NullInteractionID, "fivenet", api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: data,
	})
	handler(s, interaction)

	// German
	interaction.Locale = "de"
	(*data.Embeds)[0].Description = "FiveNet auch im Browser nutzen! Link zur FiveNet Web App."
	m.RespondInteraction(discord.NullInteractionID, "fivenet", api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: data,
	})
	handler(s, interaction)
}
