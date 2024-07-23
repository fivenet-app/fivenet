package embeds

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fivenet-app/fivenet/pkg/version"
)

const (
	ColorInfo    = 5627360
	ColorWarn    = 16166445
	ColorError   = 14100024
	ColorSuccess = 2346603
)

var EmbedAuthor = &discordgo.MessageEmbedAuthor{
	Name:    "FiveNet Discord Bot",
	URL:     "https://fivenet.app/",
	IconURL: "https://raw.githubusercontent.com/fivenet-app/fivenet/main/public/images/logo-200x200.png",
}

var EmbedFooter = &discordgo.MessageEmbedFooter{
	Text: "Version: " + version.Version,
}
