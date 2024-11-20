package utils

import (
	"strings"

	"github.com/gosimple/slug"
)

func Slug(in string) string {
	return strings.ReplaceAll(slug.Make(in), "-", ".")
}
