package tiptapsanitizer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/structpb"
)

// Module provides the Fx module for the HTML sanitizer, wiring up dependency injection.
var Module = fx.Module("tiptapsanitizer",
	fx.Provide(
		New,
	),
)

var (
	// sanitizerOnce ensures the sanitizer is initialized only once.
	sanitizerOnce sync.Once
	// sanitizer is the main for TipTap sanitization.
	sanitizer *Sanitizer
)

type Sanitizer struct {
	Nodes map[string]AttrPolicy
	Marks map[string]AttrPolicy
}

type AttrPolicy struct {
	// Validate returns (ok, normalizedAttrs)
	Validate func(attrs map[string]any) (bool, map[string]any)
}

var hexColor = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

// normalizeLinkHref ensures the href is safe and normalized.
// It allows relative URLs and absolute HTTPS URLs only.
func normalizeLinkHref(raw string) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", false
	}

	// Allow relative
	if !u.IsAbs() {
		// Disallow javascript:, data:, etc. (relative is fine)
		return u.String(), true
	}

	// Absolute: https only
	if strings.EqualFold(u.Scheme, "https") {
		return (&url.URL{
			Scheme:   "https",
			Host:     u.Host,
			Path:     u.Path,
			RawQuery: u.RawQuery,
			Fragment: u.Fragment,
		}).String(), true
	}
	return "", false
}

func normalizeColor(v any) (string, bool) {
	s, _ := v.(string)
	if s == "" {
		return "", true
	}
	// keep named colors you support, else hex
	if hexColor.MatchString(s) {
		return strings.ToLower(s), true
	}
	// allow a small whitelist of named colors, e.g., "red", "yellow"
	switch strings.ToLower(s) {
	case "yellow", "red", "green", "blue", "gray", "black", "white":
		return s, true
	}
	return "", false
}

func New(cfg *config.Config) *Sanitizer {
	sanitizerOnce.Do(func() {
		sanitizer = buildAllowed(cfg)
	})

	return buildAllowed(cfg)
}

// Build allowlist tailored to your manifest

func buildAllowed(cfg *config.Config) *Sanitizer {
	textAlignOK := func(v any) (string, bool) {
		s, _ := v.(string)
		switch s {
		case "left", "center", "right", "justify", "":
			return s, true

		default:
			return "", false
		}
	}

	node := map[string]AttrPolicy{
		"doc": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"text": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"paragraph": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if v, ok := a["textAlign"]; ok {
				if vv, ok := textAlignOK(v); ok && vv != "" {
					out["textAlign"] = vv
				}
			}
			return true, out
		}},
		"blockquote": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"bulletList": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"orderedList": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"listItem": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"hardBreak": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"horizontalRule": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"heading": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			lv, ok := a["level"].(float64)
			if ok && lv >= 1 && lv <= 6 {
				out["level"] = int(lv)
			} else {
				out["level"] = 1
				lv = 1
			}

			if v, ok := a["textAlign"]; ok {
				if vv, ok := textAlignOK(v); ok && vv != "" {
					out["textAlign"] = vv
				}
			}

			// Allow custom ID or generate one
			if id, ok := a["id"].(string); ok && id != "" {
				// Optionally whitelist pattern: ^[a-z0-9\-]{1,64}$
				out["id"] = id
			} else {
				out["id"] = fmt.Sprintf("h%d-%s", int(lv), uuid.New().String())
			}

			return true, out
		}},
		"codeBlock": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if lang, ok := a["language"].(string); ok && lang != "" {
				out["language"] = lang
			}
			return true, out
		}},
		"table": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"tableRow": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"tableHeader": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"tableCell": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if cs, ok := a["colspan"].(float64); ok && cs >= 1 && cs <= 12 {
				out["colspan"] = int(cs)
			}
			if rs, ok := a["rowspan"].(float64); ok && rs >= 1 && rs <= 100 {
				out["rowspan"] = int(rs)
			}
			return true, out
		}},
		"taskList": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"taskItem": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if ck, ok := a["checked"].(bool); ok {
				out["checked"] = ck
			}
			return true, out
		}},
		// Custom
		"checkboxStandalone": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if ck, ok := a["checked"].(bool); ok {
				out["checked"] = ck
			}
			if lbl, ok := a["label"].(string); ok && len(lbl) <= 200 {
				out["label"] = lbl
			}
			return true, out
		}},
		"image": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if src, ok := a["src"].(string); ok && src != "" {
				out["src"] = src
			} else {
				return false, nil
			}
			if alt, ok := a["alt"].(string); ok {
				out["alt"] = alt
			}
			if title, ok := a["title"].(string); ok {
				out["title"] = title
			}
			if fileId, ok := a["fileId"].(uint64); ok {
				out["fileId"] = fileId
			}
			if style, ok := a["style"].(string); ok {
				split := strings.Split(style, ";")
				if len(split) > 6 {
					return false, nil
				}

				for _, part := range split {
					keyValue := strings.SplitN(part, ":", 2)
					if len(keyValue) != 2 {
						return false, nil
					}
					if keyValue[0] != "width" {
						continue
					}

					// Set width from style attribute
					a["width"] = strings.TrimSpace(keyValue[1])
				}
			}
			if width, ok := a["width"].(float64); ok && width > 0 && width <= 5000 {
				out["width"] = int(width)
			}
			if height, ok := a["height"].(float64); ok && height > 0 && height <= 5000 {
				out["height"] = int(height)
			}
			return true, out
		}},
		"templateVar": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}

			val, _ := a["data-template-var"].(string)
			val = strings.TrimSpace(val)
			if val == "" || len(val) > 512 {
				return false, nil
			}
			out["data-template-var"] = val

			if lt, ok := a["data-left-trim"].(bool); ok {
				out["data-left-trim"] = lt
			}
			if rt, ok := a["data-right-trim"].(bool); ok {
				out["data-right-trim"] = rt
			}
			return true, out
		}},
		"templateBlock": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}

			val, _ := a["data-template-block"].(string)
			val = strings.TrimSpace(val)
			if val == "" || len(val) > 512 {
				return false, nil
			}
			out["data-template-block"] = val

			if lt, ok := a["data-left-trim"].(bool); ok {
				out["data-left-trim"] = lt
			}
			if rt, ok := a["data-right-trim"].(bool); ok {
				out["data-right-trim"] = rt
			}

			return true, out
		}},
	}

	mark := map[string]AttrPolicy{
		"bold": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"italic": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"underline": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"strike": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"code": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"subscript": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"superscript": {
			Validate: func(a map[string]any) (bool, map[string]any) { return true, map[string]any{} },
		},
		"highlight": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if c, ok := normalizeColor(a["color"]); ok && c != "" {
				out["color"] = c
			}
			return true, out
		}},
		"textStyle": {Validate: func(a map[string]any) (bool, map[string]any) {
			out := map[string]any{}
			if c, ok := normalizeColor(a["color"]); ok && c != "" {
				out["color"] = c
			}
			if c, ok := normalizeColor(a["backgroundColor"]); ok && c != "" {
				out["backgroundColor"] = c
			}
			// fontFamily/fontSize/lineHeight: optionally whitelist before keeping
			return true, out
		}},
		"link": {Validate: func(a map[string]any) (bool, map[string]any) {
			href, _ := a["href"].(string)
			if href == "" {
				return false, nil
			}
			norm, ok := normalizeLinkHref(href)
			if !ok {
				return false, nil
			}
			out := map[string]any{"href": norm}
			// enforce rel + optional target
			out["rel"] = "noopener noreferrer nofollow"
			if tgt, _ := a["target"].(string); tgt == "_blank" {
				out["target"] = "_blank"
			}
			return true, out
		}},
	}

	return &Sanitizer{
		Nodes: node,
		Marks: mark,
	}
}

// Tree sanitizer

type Stats struct {
	Words        int
	FirstHeading string
}

func SanitizeStruct(
	doc *structpb.Struct,
	maxBytes int,
	maxDepth int,
) error {
	// Convert proto struct to map[string]any
	docMap := doc.AsMap()

	// Run through the Sanitize function
	_, _, err := Sanitize(docMap, maxBytes, maxDepth)
	if err != nil {
		return err
	}

	return nil
}

func Sanitize(
	doc map[string]any,
	maxBytes int,
	maxDepth int,
) (map[string]any, Stats, error) {
	sanitizerOnce.Do(func() {
		panic("tiptap sanitizer not initialized before first use")
	})

	b, _ := json.Marshal(doc)
	if maxBytes > 0 && len(b) > maxBytes {
		return nil, Stats{}, errors.New("document too large")
	}

	var stats Stats
	out, ok := sanitizeNode(doc, sanitizer, 0, maxDepth, &stats)
	if !ok {
		return nil, Stats{}, errors.New("invalid root")
	}

	return out, stats, nil
}

func sanitizeNode(
	n map[string]any,
	allow *Sanitizer,
	depth, maxDepth int,
	stats *Stats,
) (map[string]any, bool) {
	if depth > maxDepth {
		return nil, false
	}

	typ, _ := n["type"].(string)
	if typ == "" {
		return nil, false
	}

	attrs, _ := n["attrs"].(map[string]any)
	content, _ := n["content"].([]any)
	text, _ := n["text"].(string)
	marks, _ := n["marks"].([]any)

	// text node passthrough
	if typ == "text" {
		marksOut := sanitizeMarks(marks, allow)
		if text != "" {
			stats.Words += countWords(text)
		}
		return map[string]any{
			"type":  "text",
			"text":  text,
			"marks": marksOut,
		}, true
	}

	// node allowlist
	policy, ok := allow.Nodes[typ]
	if !ok {
		return nil, false
	}

	okv, attrsOut := policy.Validate(attrs)
	if !okv {
		return nil, false
	}

	// children
	var children []any
	for _, c := range content {
		cm, _ := c.(map[string]any)
		if cm == nil {
			continue
		}
		if sn, ok := sanitizeNode(cm, allow, depth+1, maxDepth, stats); ok {
			children = append(children, sn)
		}
	}

	// derive first heading
	if stats.FirstHeading == "" && typ == "heading" && len(children) > 0 {
		if s := plainText(children); s != "" {
			stats.FirstHeading = s
		}
	}

	out := map[string]any{"type": typ}
	if len(attrsOut) > 0 {
		out["attrs"] = attrsOut
	}
	if len(children) > 0 {
		out["content"] = children
	}
	return out, true
}

func sanitizeMarks(in []any, allow *Sanitizer) []any {
	var out []any
	for _, m := range in {
		mm, _ := m.(map[string]any)
		if mm == nil {
			continue
		}

		typ, _ := mm["type"].(string)
		attrs, _ := mm["attrs"].(map[string]any)
		pol, ok := allow.Marks[typ]
		if !ok {
			continue
		}

		if okv, aout := pol.Validate(attrs); okv {
			mOut := map[string]any{"type": typ}
			if len(aout) > 0 {
				mOut["attrs"] = aout
			}
			out = append(out, mOut)
		}
	}
	return out
}

func plainText(children []any) string {
	var sb strings.Builder
	for _, c := range children {
		m, _ := c.(map[string]any)
		if m == nil {
			continue
		}

		if t, _ := m["text"].(string); t != "" {
			sb.WriteString(t)
			sb.WriteString(" ")
		}

		if arr, _ := m["content"].([]any); len(arr) > 0 {
			sb.WriteString(plainText(arr))
		}
	}
	return strings.TrimSpace(sb.String())
}

func countWords(s string) int {
	fields := strings.Fields(s)
	return len(fields)
}
