package htmlsanitizer

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	p               *bluemonday.Policy
	policyStripTags *bluemonday.Policy
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
	p.AllowStyles("color").Matching(regexp.MustCompile(`(?m)(?i)^(#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})|rgb\(\d{1,3},[ ]*\d{1,3},[ ]*\d{1,3}\))$`)).Globally()
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
	p.AllowURLSchemes("mailto", "https")

	// For linking elements we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	p.RequireNoFollowOnLinks(true)
	// Custom end

	p.AllowAttrs("cite").OnElements("blockquote", "q")
	p.AllowAttrs("href").OnElements("a", "area")
	p.AllowAttrs("src").OnElements("img")

	policyStripTags = bluemonday.StripTagsPolicy()
}

func Sanitize(in string) string {
	out := p.Sanitize(in)
	return strings.TrimSuffix(out, "<p><br></p>")
}

func StripTags(in string) string {
	return policyStripTags.Sanitize(in)
}
