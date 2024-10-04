package lang

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"golang.org/x/text/language"
)

var Module = fx.Module("lang",
	fx.Provide(New),
)

//go:embed *.json
var langFS embed.FS

func New() (*I18n, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, filePath := range map[language.Tag]string{
		language.English: "en.json",
		language.German:  "de.json",
	} {
		if _, err := bundle.LoadMessageFileFS(langFS, filePath); err != nil {
			return nil, fmt.Errorf("failed to load message file. %w", err)
		}
	}

	return &I18n{
		bundle: bundle,
	}, nil
}

type I18n struct {
	bundle *i18n.Bundle
}

func (i *I18n) I18n(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(i.bundle, lang)
}
