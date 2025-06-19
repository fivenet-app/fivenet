package notifications

import "strings"

func (x ObjectType) ToNatsKey() string {
	return strings.ToLower(x.String()[12:])
}

func (x ObjectType) ToAccessKey() string {
	switch x {
	case ObjectType_OBJECT_TYPE_DOCUMENT:
		return "documents"

	case ObjectType_OBJECT_TYPE_WIKI_PAGE:
		return "wiki_page"

	case ObjectType_OBJECT_TYPE_CITIZEN:
		return "citizen"

	default:
		return ""
	}
}
