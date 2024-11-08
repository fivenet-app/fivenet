package htmldiffer

import (
	"regexp"
	"strings"

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
			Granularity:  5,
			InsertedSpan: []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-success-600"}},
			DeletedSpan:  []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-error-600"}},
			ReplacedSpan: []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-info-600"}},
			CleanTags:    []string{""},
		},
	}
}

func (d *Differ) Diff(oldContent string, newContent string) (string, error) {
	oldContent = brFixer.ReplaceAllString(oldContent, "<br>")
	newContent = brFixer.ReplaceAllString(newContent, "<br>")
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

// TODO the diff needs to be reduced to the changes 1-2 elements around them to reduce data needs
