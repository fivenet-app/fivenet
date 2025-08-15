// Package i18n provides internationalization support for translations in multiple languages.
// It allows loading translations from JSON, setting fallback languages, and translating keys with variable replacements.
// The translations can include plural forms and variable replacements.
package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// I18n is a struct that holds translations for multiple languages and handles the translation.
type I18n struct {
	fallbackLang   string
	availableLangs []string
	// Lang → Keys
	translations map[string]map[string]any
}

// Langs returns the available languages as a string slice.
func (i *I18n) Langs() []string {
	return i.availableLangs
}

// LoadFromJSON initializes I18n with a JSON string for a specific language.
func (i *I18n) LoadFromJSON(lang string, data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal JSON for language %s. %w", lang, err)
	}

	if i.translations == nil {
		i.translations = make(map[string]map[string]any)
	}

	if _, ok := i.translations[lang]; ok {
		return fmt.Errorf("language %s already loaded", lang)
	}
	i.translations[lang] = raw

	i.availableLangs = append(i.availableLangs, lang)

	return nil
}

// SetFallbackLanguage sets the fallback language.
func (i *I18n) SetFallbackLanguage(lang string) {
	i.fallbackLang = lang
}

// Translator returns a function that translates keys for a specific language.
// The returned function takes a key and a map of variables to replace in the translation.
// It returns the translated string, if the translation is not found at all returns nothing.
//
// The map of variables can contain:
// - "n": a number for pluralization
// - Any other key-value pairs for variable replacement.
func (i *I18n) Translator(lang string) func(string, map[string]any) string {
	return func(key string, vars map[string]any) string {
		t, _ := i.translateWithFallback(lang, key, vars)
		return t
	}
}

// translateWithFallback attempts to translate with fallback.
func (i *I18n) translateWithFallback(lang, key string, vars map[string]any) (string, error) {
	langs := expandLangFallbacks(lang, i.fallbackLang)

	for _, l := range langs {
		data, ok := i.translations[l]
		if !ok {
			continue
		}
		val, err := getNestedValue(data, strings.Split(key, "."))
		if err != nil {
			continue
		}

		switch v := val.(type) {
		case string:
			return replaceVars(selectPluralForm(v, vars), vars), nil

		default:
			return "", fmt.Errorf("invalid translation format for key '%s' in language '%s'", key, l)
		}
	}

	return "", fmt.Errorf("translation for key '%s' not found in any language", key)
}

// getNestedValue retrieves a value from nested maps.
func getNestedValue(data map[string]any, keys []string) (any, error) {
	var val any = data
	for _, key := range keys {
		m, ok := val.(map[string]any)
		if !ok {
			return nil, errors.New("invalid locale json structure")
		}
		val, ok = m[key]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found", key)
		}
	}
	return val, nil
}

// replaceVars replace vars of `{KEY}` format in s and handles if variables are missing.
// It returns the original string if no replacements are made.
func replaceVars(s string, vars map[string]any) string {
	re := regexp.MustCompile(`\{(\w+)}`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		if val, ok := vars[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}

// selectPluralForm splits and selects the appropriate plural form from a single string.
func selectPluralForm(text string, vars map[string]any) string {
	forms := strings.Split(text, "|")
	for i := range forms {
		forms[i] = strings.TrimSpace(forms[i])
	}

	n := 0
	if val, ok := vars["n"]; ok {
		switch v := val.(type) {
		case int:
			n = v
		case float64:
			n = int(v)
		case string:
			if num, err := strconv.Atoi(v); err == nil {
				n = num
			}
		}
	}

	switch {
	case len(forms) == 0:
		return ""
	case len(forms) == 1:
		return forms[0]
	case len(forms) == 2:
		if n == 0 {
			return forms[0]
		}
		return forms[1]
	case n == 0:
		return forms[0]
	case n == 1:
		return forms[1]
	case len(forms) > 2:
		return forms[2]
	}

	return forms[len(forms)-1]
}

// Optimized expandLangFallbacks for readability and performance.
func expandLangFallbacks(lang, fallbackLang string) []string {
	seen := make(map[string]struct{})
	langs := []string{}

	parts := strings.FieldsFunc(lang, func(r rune) bool { return r == '_' || r == '-' })
	for i := len(parts); i >= 1; i-- {
		parent := strings.Join(parts[:i], "_")
		if _, exists := seen[parent]; !exists {
			langs = append(langs, parent)
			seen[parent] = struct{}{}
		}
	}

	if fallbackLang != "" {
		if _, exists := seen[fallbackLang]; !exists {
			langs = append(langs, fallbackLang)
		}
	}

	return langs
}
