package commands

import (
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/stretchr/testify/require"
)

func TestHelpTopicResponse(t *testing.T) {
	t.Parallel()
	l, err := i18n.New()
	require.NoError(t, err)
	c := &HelpCommand{l: l, url: "https://example.test/"}

	tests := []struct {
		name       string
		topic      string
		wantTitle  string
		wantButton int
	}{
		{name: "empty topic", topic: "", wantTitle: "No help topic chosen!", wantButton: -1},
		{name: "valid topic", topic: "browser", wantTitle: "Accessing FiveNet in your browser", wantButton: 1},
		{name: "invalid topic", topic: "unknown", wantTitle: "No help topic chosen!", wantButton: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := c.getHelpTopicResponse(discord.EnglishUS, tt.topic)
			require.NotNil(t, resp)
			require.Len(t, *resp.Embeds, 1)
			require.Equal(t, tt.wantTitle, (*resp.Embeds)[0].Title)

			require.Len(t, *resp.Components, 1)
			row, ok := (*resp.Components)[0].(*discord.ActionRowComponent)
			require.True(t, ok)
			require.Len(t, *row, len(helpTopics))
			for i, component := range *row {
				button := component.(*discord.ButtonComponent)
				if i == tt.wantButton {
					require.Equal(t, discord.PrimaryButtonStyle(), button.Style)
				} else {
					require.Equal(t, discord.SecondaryButtonStyle(), button.Style)
				}
			}
		})
	}
}
