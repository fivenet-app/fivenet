package htmlsanitizer

import (
	"html"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
)

var (
	stripTagsOnce sync.Once
	stripTags     *bluemonday.Policy

	sanitizerOnce sync.Once
	sanitizer     *bluemonday.Policy
)

var colorRegex = regexp.MustCompile(`(?m)(?i)^(#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})|rgb\(\d{1,3},[ ]*\d{1,3},[ ]*\d{1,3}\))$`)

var Module = fx.Module("htmlsanitizer",
	fx.Provide(
		New,
	),
)

func setupSanitizer() {
	// Custom UGC Policy
	sanitizer = bluemonday.UGCPolicy()

	// "img" is permitted
	sanitizer.AllowAttrs("align").Matching(bluemonday.ImageAlign).OnElements("img")
	sanitizer.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	sanitizer.AllowAttrs("height", "width").Matching(bluemonday.NumberOrPercent).OnElements("img")

	// Standard URLs enabled
	sanitizer.AllowAttrs("src").OnElements("img")

	// Allow in-line images (for now)
	sanitizer.AllowDataURIImages()

	// Style
	sanitizer.AllowAttrs("style").OnElements("span", "p", "img")
	// Image centering
	sanitizer.AllowStyles("display").OnElements("span", "p", "img")
	sanitizer.AllowStyles("margin-left").OnElements("span", "p", "img")
	sanitizer.AllowStyles("margin-right").OnElements("span", "p", "img")
	// Allow the 'color' property with valid RGB(A) hex values only (on any element allowed a 'style' attribute)
	sanitizer.AllowStyles("color").Matching(colorRegex).Globally()
	sanitizer.AllowStyles("text-align").Globally()
	sanitizer.AllowStyles("font-weight").Globally()
	sanitizer.AllowStyles("font-size").Globally()
	sanitizer.AllowStyles("line-height").Globally()
	// Allow the 'text-decoration' property to be set to 'underline', 'line-through' or 'none'
	// on 'span' and 'p' elements only
	sanitizer.AllowStyles("text-decoration").MatchingEnum("underline", "line-through", "none").OnElements("span", "p")

	// Links
	// Custom policy based on the origional "AllowStandardURLs" helper func
	// URLs must be parseable by net/url.Parse()
	sanitizer.RequireParseableURLs(true)

	// Allow relative URLs (!url.IsAbs() is permitted)
	sanitizer.AllowRelativeURLs(true)

	// Most common URL schemes only
	sanitizer.AllowURLSchemes("https")

	// For linking elements we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	sanitizer.RequireNoFollowOnLinks(true)

	sanitizer.AllowAttrs("cite").OnElements("blockquote", "q")
	sanitizer.AllowAttrs("href").OnElements("a", "area")
	sanitizer.AllowAttrs("src").OnElements("img")
	sanitizer.AllowElements("hr", "sup", "sub", "h1", "h2", "h3", "h4", "h5", "code", "em", "pre")
	sanitizer.AllowTables()
	sanitizer.AllowLists()

	// Checkboxes
	sanitizer.AllowAttrs("contenteditable").Matching(regexp.MustCompile(`(?i)false`)).OnElements("label")
	sanitizer.AllowAttrs("type").Matching(regexp.MustCompile(`(?i)checkbox`)).OnElements("input")
	sanitizer.AllowAttrs("checked").Matching(regexp.MustCompile(`(?i)true`)).OnElements("input")
}

func New(cfg *config.Config) (*bluemonday.Policy, error) {
	sanitizerOnce.Do(setupSanitizer)

	// Use Image Proxy if enabled
	if cfg.ImageProxy.Enabled {
		proxyUrl, err := url.Parse(cfg.ImageProxy.URL)
		if err != nil {
			return nil, err
		}

		sanitizer.RewriteSrc(func(u *url.URL) {
			// Rewrite URLs to image proxy to proxy all requests through a single URL.
			imgUrl, _ := url.PathUnescape(u.String())

			if u.Scheme == "data" {
				return
			}

			if u.Scheme != proxyUrl.Scheme {
				u.Scheme = proxyUrl.Scheme
			}
			if u.Host != proxyUrl.Host {
				u.Host = proxyUrl.Host
			}
			if !strings.HasPrefix(u.Path, proxyUrl.Path) {
				u.Path = proxyUrl.Path + imgUrl
			}
			u.RawQuery = ""
		})
	}

	return sanitizer, nil
}

func Sanitize(in string) string {
	sanitizerOnce.Do(setupSanitizer)

	out := sanitizer.Sanitize(in)
	return strings.TrimSuffix(out, "<p><br></p>")
}

func StripTags(in string) string {
	stripTagsOnce.Do(func() {
		stripTags = bluemonday.StrictPolicy()
	})

	return html.UnescapeString(stripTags.Sanitize(in))
}
