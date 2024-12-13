package utils

import (
	"strings"

	"github.com/gosimple/slug"
)

func SlugNoDots(in string) string {
	return slug.Make(in)
}

func Slug(in string) string {
	return strings.ReplaceAll(slug.Make(in), "-", ".")
}
