package events

import "strings"

func SanitizeKey(in string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(in, "#", "."),
				":", "."),
			"/", "."),
		"=", "_")
}
