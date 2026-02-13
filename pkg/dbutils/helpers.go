package dbutils

import (
	"strings"
)

func PrepareForLikeSearch(input string) string {
	// 1. Trim leading and trailing spaces
	input = strings.TrimSpace(input)

	// 2. Normalize multiple spaces to a single space
	input = strings.Join(strings.Fields(input), " ")

	// 3. Escape special characters
	input = strings.ReplaceAll(input, "%", "\\%")
	input = strings.ReplaceAll(input, "_", "\\_")
	input = strings.ReplaceAll(input, "\t", " ")

	// 4. Replace spaces with `%` for LIKE condition
	input = strings.ReplaceAll(input, " ", "%")

	// 5. Wrap with `%` if not empty
	if input != "" {
		input = "%" + input + "%"
	}

	return input
}
