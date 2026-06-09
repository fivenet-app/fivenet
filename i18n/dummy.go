package i18n

import "fmt"

type Dummy struct{}

func NewDummy() Ii18n {
	return &Dummy{}
}

func (d *Dummy) Langs() []string {
	return []string{"en", "de"}
}

func (d *Dummy) GetFallbackLanguage() string {
	return "en"
}
func (d *Dummy) SetFallbackLanguage(lang string) {}

func (d *Dummy) Translator(lang string) TFunc {
	return DummyTranslator()
}

func DummyTranslator() TFunc {
	return func(key string, vars map[string]any) string {
		return fmt.Sprintf("%s(%+v)", key, vars)
	}
}
