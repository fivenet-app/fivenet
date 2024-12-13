package content

import (
	"html"
	"regexp"
	"strings"

	"github.com/aymanbagabas/go-udiff"
)

var removeImgData = regexp.MustCompile(`(?m)"data:image/[^"]+"`)

func DiffHTML(old string, new string) string {
	old = removeImgData.ReplaceAllString(old, "\"IMAGE_DATA_OMITTED\"")
	old = strings.ReplaceAll(old, "<br/>", "<br>")
	old = html.UnescapeString(old)

	new = removeImgData.ReplaceAllString(new, "\"IMAGE_DATA_OMITTED\"")
	new = strings.ReplaceAll(new, "<br/>", "<br>")
	new = html.UnescapeString(new)

	if strings.EqualFold(old, new) {
		return ""
	}

	d := udiff.Unified("a", "b", old, new)
	return d
}
