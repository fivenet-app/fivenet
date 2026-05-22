package commands

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/fivenet-app/fivenet/v2026/i18n"
)

const commandDefaultLang = "en"

var commandLocaleLangs = map[discord.Language]string{
	discord.German: "de",
}

type commandLocalizer struct {
	i18n *i18n.I18n
	base string
}

func newCommandLocalizer(i18n *i18n.I18n, base string) *commandLocalizer {
	return &commandLocalizer{
		i18n: i18n,
		base: base,
	}
}

func (l *commandLocalizer) key(suffix string) string {
	return l.base + "." + suffix
}

func (l *commandLocalizer) text(suffix string) string {
	return l.i18n.Translator(commandDefaultLang)(l.key(suffix), nil)
}

func (l *commandLocalizer) localizations(suffix string) discord.StringLocales {
	localizations := discord.StringLocales{}
	for locale, langCode := range commandLocaleLangs {
		localizations[locale] = l.i18n.Translator(langCode)(l.key(suffix), nil)
	}
	return localizations
}
