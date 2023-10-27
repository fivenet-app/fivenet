package htmlsanitizer

import (
	"html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	p         *bluemonday.Policy
	stripTags *bluemonday.Policy
)

var (
	colorRegex = regexp.MustCompile(`(?m)(?i)^(#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})|rgb\(\d{1,3},[ ]*\d{1,3},[ ]*\d{1,3}\))$`)
)

func init() {
	p = bluemonday.UGCPolicy()

	// "img" is permitted
	p.AllowAttrs("align").Matching(bluemonday.ImageAlign).OnElements("img")
	p.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(bluemonday.NumberOrPercent).OnElements("img")

	// Standard URLs enabled
	p.AllowAttrs("src").OnElements("img")

	// Allow in-line images (for now)
	p.AllowDataURIImages()

	// Style
	p.AllowAttrs("style").OnElements("span", "p")
	// Allow the 'color' property with valid RGB(A) hex values only (on any element allowed a 'style' attribute)
	p.AllowStyles("color").Matching(colorRegex).Globally()
	p.AllowStyles("text-align").Globally()
	p.AllowStyles("font-weight").Globally()
	p.AllowStyles("font-size").Globally()
	p.AllowStyles("line-height").Globally()
	// Allow the 'text-decoration' property to be set to 'underline', 'line-through' or 'none'
	// on 'span' elements only
	p.AllowStyles("text-decoration").MatchingEnum("underline", "line-through", "none").OnElements("span", "p")

	// Links
	// Custom policy based on the origional "AllowStandardURLs" helper func
	// URLs must be parseable by net/url.Parse()
	p.RequireParseableURLs(true)

	// !url.IsAbs() is permitted
	p.AllowRelativeURLs(true)

	// Most common URL schemes only
	p.AllowURLSchemes("https")

	// For linking elements we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	p.RequireNoFollowOnLinks(true)
	// CUSTOM END

	p.AllowAttrs("cite").OnElements("blockquote", "q")
	p.AllowAttrs("href").OnElements("a", "area")
	p.AllowAttrs("src").OnElements("img")
	p.AllowElements("hr", "sup", "sub", "h1", "h2", "h3", "h4", "h5", "code")
	p.AllowTables()
	p.AllowLists()

	stripTags = bluemonday.StrictPolicy()
}

func Sanitize(in string) string {
	out := p.Sanitize(in)
	return strings.TrimSuffix(out, "<p><br></p>")
}

func StripTags(in string) string {
	return html.UnescapeString(stripTags.Sanitize(in))
}
