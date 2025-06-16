package common

func NewI18NItem(key string) *I18NItem {
	return &I18NItem{
		Key: key,
	}
}

func NewI18NItemWithParams(key string, params map[string]string) *I18NItem {
	return &I18NItem{
		Key:        key,
		Parameters: params,
	}
}
