package embeds

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/pkg/version"
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
	Icon: "https://raw.githubusercontent.com/fivenet-app/fivenet/main/public/images/logo-200x200.png",
}

var EmbedFooterVersion = &discord.EmbedFooter{
	Text: "Version: " + version.Version,
}

var EmbedFooterMadeBy = &discord.EmbedFooter{
	Text: "Made by Galexrt",
	Icon: "https://cdn.discordapp.com/avatars/290472392084422658/58e38b558fb3a54b4864584e7b7297f6.png",
}
