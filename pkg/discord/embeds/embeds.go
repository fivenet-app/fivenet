package embeds

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

const (
	ColorInfo    = 5627360
	ColorWarn    = 16166445
	ColorError   = 14100024
	ColorSuccess = 2346603
)

var EmbedAuthor = &discord.EmbedAuthor{
	Name: "FiveNet Discord Bot",
	URL:  "https://fivenet.app/",
	Icon: "https://raw.githubusercontent.com/fivenet-app/fivenet/main/public/images/logo-200x200.webp",
}

var EmbedFooterVersion = &discord.EmbedFooter{
	Text: "Version: " + version.Version,
}

var EmbedFooterMadeBy = &discord.EmbedFooter{
	Text: "Made by Galexrt",
	Icon: "https://galexrt.moe/favicon.png",
}

var EmbedFooterFiveNet = &discord.EmbedFooter{
	Text: "FiveNet",
	Icon: "https://raw.githubusercontent.com/fivenet-app/fivenet/main/public/images/logo-200x200.webp",
}

var EmbedThumbnailLogo = &discord.EmbedThumbnail{
	URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
	Width:  128,
	Height: 128,
}
