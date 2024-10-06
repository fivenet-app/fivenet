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

	i := &I18n{
		bundle: bundle,
	}

	for _, filePath := range map[language.Tag]string{
		language.English: "en.json",
		language.German:  "de.json",
	} {
		f, err := bundle.LoadMessageFileFS(langFS, filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load message file. %w", err)
		}

		i.availableLangs = append(i.availableLangs, f.Tag.String())
	}

	return i, nil
}

type I18n struct {
	bundle *i18n.Bundle

	availableLangs []string
}

func (i *I18n) I18n(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(i.bundle, lang)
}

func (i *I18n) Langs() []string {
	return i.availableLangs
}
