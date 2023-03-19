package htmlsanitizer

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var (
	p               *bluemonday.Policy
	policyStripTags *bluemonday.Policy
)

func init() {
	p = bluemonday.UGCPolicy()

	// Style
	p.AllowAttrs("style").OnElements("span", "p")
	// Allow the 'color' property with valid RGB(A) hex values only (on any element allowed a 'style' attribute)
	p.AllowStyles("color").Matching(regexp.MustCompile("(?i)^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$")).Globally()
	// Allow the 'text-decoration' property to be set to 'underline', 'line-through' or 'none'
	// on 'span' elements only
	p.AllowStyles("text-decoration").MatchingEnum("underline", "line-through", "none").OnElements("span", "p")

	// Links
	p.AllowStandardURLs()
	p.AllowAttrs("cite").OnElements("blockquote", "q")
	p.AllowAttrs("href").OnElements("a", "area")
	p.AllowAttrs("src").OnElements("img")

	policyStripTags = bluemonday.StripTagsPolicy()
}

func Sanitize(in string) string {
	return p.Sanitize(in)
}

func StripTags(in string) string {
	return policyStripTags.Sanitize(in)
}
