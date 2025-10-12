package common

func NewI18nItem(key string) *I18NItem {
	return &I18NItem{
		Key: key,
	}
}

func NewI18nItemWithParams(key string, params map[string]string) *I18NItem {
	return &I18NItem{
		Key:        key,
		Parameters: params,
	}
}
