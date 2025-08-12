package i18n

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/fx"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var Module = fx.Module("lang",
	fx.Provide(New),
)

//go:embed locales/*.json
var langFS embed.FS

const defaultLanguage = "en"

func New() (*I18n, error) {
	i := &I18n{}
	i.SetFallbackLanguage(defaultLanguage)

	if err := fs.WalkDir(langFS, "locales", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}

		lang := strings.TrimSuffix(d.Name(), ".json")

		data, err := langFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s. %w", path, err)
		}

		if err := i.LoadFromJSON(lang, data); err != nil {
			return fmt.Errorf("failed to load translations from file %s. %w", path, err)
		}

		i.availableLangs = append(i.availableLangs, lang)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk locales dir. %w", err)
	}

	return i, nil
}
