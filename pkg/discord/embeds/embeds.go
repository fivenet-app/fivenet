package embeds

import "github.com/bwmarrin/discordgo"

var EmbedAuthor = &discordgo.MessageEmbedAuthor{
	Name:    "FiveNet Discord Bot",
	URL:     "https://fivenet.app/",
	IconURL: "https://raw.githubusercontent.com/galexrt/fivenet/main/src/public/images/logo-200x200.png",
}

const (
	ColorInfo    = 5627360
	ColorWarn    = 16166445
	ColorFailure = 14100024
	ColorSuccess = 2346603
)
