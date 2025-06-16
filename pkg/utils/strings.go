package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// StringFirstN returns the first n runes of a string, handling multi-byte characters safely.
// Taken from "KAdot" here: https://stackoverflow.com/a/41604514
func StringFirstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

// StringFirstToLower lowercases the first rune of a string, if possible.
// Taken from "rocka2q" here: https://stackoverflow.com/a/75989905
func StringFirstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

// commonTitlePrefixes is a list of common title prefixes to remove from names.
var commonTitlePrefixes = []string{
	"prof.", "prof ",
	"dr.", "dr ",
	"sr.", "sr ",
}

// RemoveTitlePrefixes removes known title prefixes from the start of a string.
func RemoveTitlePrefixes(s string) string {
	s = strings.TrimSpace(s)
	prefixes := commonTitlePrefixes

	for {
		lower := strings.ToLower(s)
		matched := false
		for _, p := range prefixes {
			if strings.HasPrefix(lower, p) {
				// Remove using the original string slice
				s = strings.TrimSpace(s[len(p):])
				matched = true
				break
			}
		}
		if !matched {
			break
		}
	}

	return s
}
