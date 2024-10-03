package commands

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/pkg/lang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandleFivenetCommand(t *testing.T) {
	l, err := lang.New()
	require.NoError(t, err)

	cfg, err := config.LoadTestConfig()
	require.NoError(t, err)
	url := "https://example.fivenet.app/"
	cfg.HTTP.PublicURL = url

	_, handler, err := NewFivenetCommand(cfg, l)
	require.NoError(t, err)

	router := cmdroute.NewRouter()
	router.Add("fivenet", handler)

	interactionEvent := &discord.InteractionEvent{
		ID: discord.NullInteractionID,
		Data: &discord.CommandInteraction{
			ID:   discord.NullCommandID,
			Name: "fivenet",
		},
		Locale: discord.EnglishUS,
	}

	data := &api.InteractionResponseData{
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
					Icon: embeds.EmbedFooterMadeBy.Icon,
				},
				Provider: &discord.EmbedProvider{
					Name: "FiveNet",
					URL:  url,
				},
			},
		},
	}

	// English
	resp := router.HandleInteraction(interactionEvent)
	require.NotNil(t, resp.Data)
	assert.Len(t, *resp.Data.Embeds, 1)
	assert.Equal(t, (*resp.Data.Embeds)[0].Description, (*data.Embeds)[0].Description)

	// German
	interactionEvent.Locale = "de"
	(*data.Embeds)[0].Description = "FiveNet auch im Browser nutzen! Link zur FiveNet Web App."
	resp = router.HandleInteraction(interactionEvent)
	require.NotNil(t, resp.Data)
	assert.Len(t, *resp.Data.Embeds, 1)
	assert.Equal(t, (*resp.Data.Embeds)[0].Description, (*data.Embeds)[0].Description)
}
