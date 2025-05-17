package i18n

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslateWithFallback(t *testing.T) {
	i18n := &I18n{}
	err := loadTestTranslations(i18n)
	require.NoError(t, err)

	i18n.SetFallbackLanguage("en")

	tests := []struct {
		name     string
		lang     string
		key      string
		vars     map[string]any
		expected string
		err      string
	}{
		{
			name:     "Basic translation",
			lang:     "en",
			key:      "greeting.hello",
			vars:     nil,
			expected: "Hello",
		},
		{
			name:     "\"Deeply\" nested translation",
			lang:     "en",
			key:      "greeting.nested.hello",
			vars:     nil,
			expected: "Hello from nested",
		},
		{
			name:     "Translation with variables",
			lang:     "en",
			key:      "greeting.welcome",
			vars:     map[string]any{"name": "John"},
			expected: "Welcome, John!",
		},
		{
			name:     "Fallback language translation",
			lang:     "de",
			key:      "greeting.non_existing_in_other_lang",
			vars:     nil,
			expected: "Not existing, fallback :-(",
		},
		{
			name:     "Missing key in fallback",
			lang:     "de",
			key:      "nonexistent.key",
			vars:     nil,
			expected: "",
			err:      "translation for key 'nonexistent.key' not found in any language",
		},
		{
			name:     "Plural form translation",
			lang:     "en",
			key:      "items",
			vars:     map[string]any{"n": 1},
			expected: "1 item",
		},
		{
			name:     "Plural form translation (multiple)",
			lang:     "en",
			key:      "items",
			vars:     map[string]any{"n": 5},
			expected: "5 items",
		},
		{
			name:     "Plural form translation (multiple) with only 2 variants",
			lang:     "en",
			key:      "books",
			vars:     map[string]any{"n": 5},
			expected: "5 books",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := i18n.translateWithFallback(tt.lang, tt.key, tt.vars)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func loadTestTranslations(i18n *I18n) error {
	if err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}

		lang := strings.TrimSuffix(d.Name(), ".json")

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s. %w", path, err)
		}

		if err := i18n.LoadFromJSON(lang, data); err != nil {
			return fmt.Errorf("failed to load translations from file %s. %w", path, err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk testdata locales dir. %w", err)
	}

	return nil
}

func TestTranslateViaTranslatorFunc(t *testing.T) {
	i18n := &I18n{}
	err := loadTestTranslations(i18n)
	require.NoError(t, err)

	i18n.SetFallbackLanguage("en")

	translatorFunc := i18n.Translator("en")
	result := translatorFunc("greeting.hello", nil)
	assert.Equal(t, "Hello", result)

	result = translatorFunc("items", map[string]any{"n": 3})
	assert.Equal(t, "3 items", result)

	translatorFunc = i18n.Translator("de")
	result = translatorFunc("greeting.hello", nil)
	assert.Equal(t, "Hallo", result)

	result = translatorFunc("computers", map[string]any{"n": 3})
	assert.Equal(t, "3 computers", result)
}

// Unit test for missing variables
func TestReplaceVars_MissingVariables(t *testing.T) {
	s := "Hello, {name}!"
	vars := map[string]any{}
	result := replaceVars(s, vars)
	assert.Equal(t, "Hello, {name}!", result)
}
