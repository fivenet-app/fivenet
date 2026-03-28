package commands

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	lang "github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
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

	router := cmdroute.NewRouter()
	cmd, err := NewFivenetCommand(CommandParams{
		Cfg: cfg,
		L:   l,
	})
	require.NoError(t, err)
	require.NotNil(t, cmd)
	cmd.RegisterCommand(router)

	interactionEvent := &discord.InteractionEvent{
		ID: discord.NullInteractionID,
		Data: &discord.CommandInteraction{
			ID:   discord.NullCommandID,
			Name: "fivenet",
		},
		Locale: discord.EnglishUS,
	}

	// English
	resp := router.HandleInteraction(interactionEvent)
	require.NotNil(t, resp.Data)
	// Check embeds
	assert.Len(t, *resp.Data.Embeds, 1)
	assert.Equal(
		t,
		"FiveNet is also available in your browser! Link to the FiveNet web app.",
		(*resp.Data.Embeds)[0].Description,
	)
	// Check components
	require.Equal(t, discord.ActionRowComponentType, (*resp.Data.Components)[0].Type())
	actionRow, ok := (*resp.Data.Components)[0].(*discord.ActionRowComponent)
	require.True(t, ok)

	button := (*actionRow)[0].(*discord.ButtonComponent)
	assert.Equal(t, "Open FiveNet", button.Label)
	assert.Equal(t, discord.LinkButtonStyle(url), button.Style)

	// German
	interactionEvent.Locale = "de"
	resp = router.HandleInteraction(interactionEvent)
	require.NotNil(t, resp.Data)
	// Check embeds
	assert.Len(t, *resp.Data.Embeds, 1)
	assert.Equal(
		t,
		"FiveNet auch im Browser nutzen! Link zur FiveNet Web App.",
		(*resp.Data.Embeds)[0].Description,
	)
	// Check components
	require.Equal(t, discord.ActionRowComponentType, (*resp.Data.Components)[0].Type())
	actionRow, ok = (*resp.Data.Components)[0].(*discord.ActionRowComponent)
	require.True(t, ok)

	button = (*actionRow)[0].(*discord.ButtonComponent)
	assert.Equal(t, "FiveNet Web App öffnen", button.Label)
	assert.Equal(t, discord.LinkButtonStyle(url), button.Style)
}
