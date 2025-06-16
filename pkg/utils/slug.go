package utils

import (
	"strings"

	"github.com/gosimple/slug"
)

// SlugNoDots returns a slugified version of the input string using the gosimple/slug package.
func SlugNoDots(in string) string {
	return slug.Make(in)
}

// Slug returns a slugified version of the input string, replacing dashes with dots.
func Slug(in string) string {
	return strings.ReplaceAll(slug.Make(in), "-", ".")
}
