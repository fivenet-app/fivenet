package htmldiffer

import (
	"html"
	"regexp"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	htmldiff "github.com/documize/html-diff"
	"go.uber.org/fx"
)

// Module provides the Fx module for the HTML differ, wiring up dependency injection.
var Module = fx.Module("htmldiffer",
	fx.Provide(
		New,
	),
)

// brFixer is a regexp that normalizes multiple <br> tags and whitespace into a single <br/>.
var brFixer = regexp.MustCompile(`(?m)(<br>)+([ \n]*)(<br>)+`)

// Differ wraps the html-diff configuration for HTML diffing operations.
type Differ struct {
	// htmldiff is the configuration for the html-diff library.
	htmldiff *htmldiff.Config
}

// New creates a new Differ instance with default configuration for HTML diffing.
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

// FancyDiff computes a highlighted HTML diff between oldContent and newContent.
// If no changes are detected, returns an empty string. Falls back to newContent on error.
func (d *Differ) FancyDiff(old string, new string) (string, error) {
	old = brFixer.ReplaceAllString(old, "<br/>")
	new = brFixer.ReplaceAllString(new, "<br/>")
	res, err := d.htmldiff.HTMLdiff([]string{old, new})
	if err != nil {
		// Fallback to the new content
		//nolint:nilerr // If diffing fails, return the new content as is so no data is lost.
		return new, nil
	}

	out := res[0]
	// If no "htmldiff" change markers are found, return empty string
	if !strings.Contains(out, "htmldiff") {
		return "", nil
	}

	return out, nil
}

// removeImgData is a regexp that strips inline image data from HTML for diffing.
var removeImgData = regexp.MustCompile(`(?i)data:image/[^"']+`)

// PatchDiff computes a unified diff (udiff) between two HTML strings, omitting image data and normalizing <br> tags.
// Returns an empty string if the normalized HTML is equal.
func (d *Differ) PatchDiff(old string, new string) string {
	old = removeImgData.ReplaceAllString(old, "IMAGE_DATA_OMITTED")
	old = strings.ReplaceAll(old, "<br/>", "<br>")
	old = html.UnescapeString(old)

	new = removeImgData.ReplaceAllString(new, "IMAGE_DATA_OMITTED")
	new = strings.ReplaceAll(new, "<br/>", "<br>")
	new = html.UnescapeString(new)

	if strings.EqualFold(old, new) {
		return ""
	}

	out := udiff.Unified("a", "b", old, new)
	return out
}
