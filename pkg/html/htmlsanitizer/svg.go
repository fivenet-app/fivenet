package htmlsanitizer

import (
	"strings"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

var (
	sanitizerSVG *bluemonday.Policy

	// sanitizerOnce ensures the sanitizer policy is initialized only once.
	sanitizerSVGOnce sync.Once
)

func setupSanitizerSVG() {
	sanitizerSVG = bluemonday.UGCPolicy()

	sanitizerSVG.AllowAttrs("xmlns", "xmlns:xlink", "viewBox", "width", "height").
		OnElements("svg")

	sanitizerSVG.
		AllowAttrs("d", "stroke-width", "stroke", "fill", "stroke-linecap").
		OnElements("path")
}

func SanitizeSVG(in string) string {
	sanitizerSVGOnce.Do(setupSanitizerSVG)

	out := sanitizerSVG.Sanitize(in)

	return strings.Replace(out, "viewbox=\"", "viewBox=\"", 1)
}
