package htmldiffer

import (
	"html"
	"regexp"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	htmldiff "github.com/documize/html-diff"
	"go.uber.org/fx"
)

var Module = fx.Module("htmldiffer",
	fx.Provide(
		New,
	),
)

var brFixer = regexp.MustCompile(`(?m)(<br>)+([ \n]*)(<br>)+`)

type Differ struct {
	htmldiff *htmldiff.Config
}

func New() *Differ {
	return &Differ{
		htmldiff: &htmldiff.Config{
			Granularity:  6,
			InsertedSpan: []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-success-600"}},
			DeletedSpan:  []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-error-600"}},
			ReplacedSpan: []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-info-600"}},
			CleanTags:    []string{},
		},
	}
}

func (d *Differ) FancyDiff(oldContent string, newContent string) (string, error) {
	oldContent = brFixer.ReplaceAllString(oldContent, "<br/>")
	newContent = brFixer.ReplaceAllString(newContent, "<br/>")
	res, err := d.htmldiff.HTMLdiff([]string{oldContent, newContent})
	if err != nil {
		// Fallback to the new content
		return newContent, nil
	}

	out := res[0]
	// If no "htmldiff" change markers are found, return empty string
	if !strings.Contains(out, "htmldiff") {
		return "", nil
	}

	return out, nil
}

var removeImgData = regexp.MustCompile(`(?i)data:image/[^"']+`)

func (d *Differ) PatchDiff(old string, new string) string {
	old = removeImgData.ReplaceAllString(old, "\"IMAGE_DATA_OMITTED\"")
	old = strings.ReplaceAll(old, "<br/>", "<br>")
	old = html.UnescapeString(old)

	new = removeImgData.ReplaceAllString(new, "\"IMAGE_DATA_OMITTED\"")
	new = strings.ReplaceAll(new, "<br/>", "<br>")
	new = html.UnescapeString(new)

	if strings.EqualFold(old, new) {
		return ""
	}

	out := udiff.Unified("a", "b", old, new)
	return out
}
