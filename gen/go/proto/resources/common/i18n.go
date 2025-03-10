package common

func NewTranslateItem(key string) *TranslateItem {
	return &TranslateItem{
		Key: key,
	}
}

func NewTranslateItemWithParams(key string, params map[string]string) *TranslateItem {
	return &TranslateItem{
		Key:        key,
		Parameters: params,
	}
}
