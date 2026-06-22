package embeds

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

const (
	ColorInfo    discord.Color = 5627360
	ColorWarn    discord.Color = 16166445
	ColorError   discord.Color = 14100024
	ColorSuccess discord.Color = 2346603
)

var (
	EmbedAuthor = &discord.EmbedAuthor{
		Name: version.ProjectName + " Discord Bot",
		URL:  "https://fivenet.app/",
		Icon: "https://raw.githubusercontent.com/fivenet-app/fivenet/main/public/images/logo-200x200.webp",
	}

	EmbedFooterVersion = &discord.EmbedFooter{
		Text: "Version: " + version.Version,
	}

	EmbedFooterMadeBy = &discord.EmbedFooter{
		Text: "Made by Galexrt",
		Icon: "https://galexrt.moe/favicon.png",
	}

	EmbedFooterFiveNet = &discord.EmbedFooter{
		Text: version.ProjectName,
		Icon: "https://raw.githubusercontent.com/fivenet-app/fivenet/main/public/images/logo-200x200.webp",
	}

	EmbedThumbnailLogo = &discord.EmbedThumbnail{
		URL:    "https://cdn.discordapp.com/app-icons/1101207666652618865/94429951df15108c737949ff2770cd8f.png",
		Width:  128,
		Height: 128,
	}
)
