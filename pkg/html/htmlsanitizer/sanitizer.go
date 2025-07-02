package htmlsanitizer

import (
	"html"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/fx"
)

var (
	// stripTagsOnce ensures the stripTags policy is initialized only once.
	stripTagsOnce sync.Once
	// stripTags is a bluemonday policy for strict tag stripping.
	stripTags *bluemonday.Policy

	// sanitizerOnce ensures the sanitizer policy is initialized only once.
	sanitizerOnce sync.Once
	// sanitizer is the main bluemonday policy for HTML sanitization.
	sanitizer *bluemonday.Policy
)

var (
	// colorRegex matches valid color values for style attributes.
	colorRegex = regexp.MustCompile(`(?mi)^(#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})|rgb\(\d{1,3},[ ]*\d{1,3},[ ]*\d{1,3}\))$`)
	// fontFamilyRegex matches valid font family names for style attributes.
	fontFamilyRegex = regexp.MustCompile(`(?mi)^(arial,\shelvetica,\ssans-serif|times new roman,\stimes,\sserif|Comic Sans MS,\sComic Sans|serif|monospace|DM Sans)$`)

	// prosemirrorClassRegex matches ProseMirror class names for editor compatibility.
	prosemirrorClassRegex = regexp.MustCompile(`(?m)^ProseMirror-[A-Za-z]+$`)

	// boolFalseRegex matches the string "false" (case-insensitive).
	boolFalseRegex = regexp.MustCompile(`(?i)^false$`)
	// boolTrueRegex matches the string "true" or "checked" (case-insensitive).
	boolTrueRegex = regexp.MustCompile(`(?i)^(true|checked)$`)
	// inputTypeCheckbox matches the string "checkbox" (case-insensitive).
	inputTypeCheckbox = regexp.MustCompile(`(?i)checkbox`)
)

// Module provides the Fx module for the HTML sanitizer, wiring up dependency injection.
var Module = fx.Module("htmlsanitizer",
	fx.Provide(
		New,
	),
)

// setupSanitizer initializes the main bluemonday sanitizer policy with custom rules for UGC, images, styles, links, and editor compatibility.
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
	// Image centering + positioning
	sanitizer.AllowStyles("display").OnElements("span", "p", "img")
	sanitizer.AllowStyles("margin-left").OnElements("span", "p", "img")
	sanitizer.AllowStyles("margin-right").OnElements("span", "p", "img")
	sanitizer.AllowStyles("height").OnElements("img")
	sanitizer.AllowStyles("width").OnElements("img")
	sanitizer.AllowStyles("margin").OnElements("img")

	// Allow the 'color' property with valid RGB(A) hex values only (on any element allowed a 'style' attribute)
	sanitizer.AllowStyles("color").Matching(colorRegex).Globally()
	sanitizer.AllowStyles("font-family").Matching(fontFamilyRegex).Globally()
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
	sanitizer.AllowNoAttrs().OnElements("label")
	sanitizer.AllowAttrs("contenteditable").Matching(boolFalseRegex).OnElements("label")
	sanitizer.AllowAttrs("type").Matching(inputTypeCheckbox).OnElements("input")
	sanitizer.AllowAttrs("checked").Matching(boolTrueRegex).OnElements("input")

	// # ProseMirror / Tiptap Editor
	sanitizer.AllowAttrs("class").Matching(prosemirrorClassRegex).OnElements("br")
	// ## Checkboxes
	sanitizer.AllowAttrs("data-checked").OnElements("li", "span")
	sanitizer.AllowAttrs("data-type").OnElements("ul", "ol", "li", "span")
}

// New creates and returns a new bluemonday.Policy for HTML sanitization, optionally enabling image proxy rewriting if configured.
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
			if !strings.HasPrefix(u.Path, proxyUrl.Path) && !strings.HasPrefix(u.Path, "/api/filestore/") {
				u.Path = proxyUrl.Path + imgUrl
			}
			u.RawQuery = ""
		})
	}

	return sanitizer, nil
}

// Sanitize applies the main HTML sanitizer policy to the input string and trims trailing empty paragraphs.
func Sanitize(in string) string {
	sanitizerOnce.Do(setupSanitizer)

	out := sanitizer.Sanitize(in)
	return strings.TrimSuffix(out, "<p><br></p>")
}

// StripTags removes all HTML tags from the input string using a strict bluemonday policy and returns the unescaped result.
func StripTags(in string) string {
	stripTagsOnce.Do(func() {
		stripTags = bluemonday.StrictPolicy()
	})

	return html.UnescapeString(stripTags.Sanitize(in))
}
