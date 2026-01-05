package content

import (
	"strings"
	"unicode"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

// ExtractFromTiptap converts a Tiptap JSON doc (Struct) into text + derived fields.
// Assumes you've already sanitized/validated the doc (allowed nodes/marks, URL policy, etc.).
func ExtractFromTiptap(doc *structpb.Struct) *ExtractedContent {
	if doc == nil {
		return &ExtractedContent{}
	}
	m := doc.AsMap()

	var b strings.Builder
	var firstHeading string

	walkNode(m, &b, &firstHeading, false)

	text := normalizeWhitespace(b.String())
	return &ExtractedContent{
		Text:         text,
		WordCount:    uint32(len(strings.Fields(text))),
		FirstHeading: normalizeWhitespace(firstHeading),
	}
}

func walkNode(node any, b *strings.Builder, firstHeading *string, inHeading bool) {
	m, ok := node.(map[string]any)
	if !ok {
		return
	}

	typ, _ := m["type"].(string)

	// Handle text nodes
	if typ == "text" {
		if s, _ := m["text"].(string); s != "" {
			b.WriteString(s)
		}
		return
	}

	// Block-ish separators
	beforeBlock(typ, b)

	// Detect heading for first_heading
	isHeading := typ == "heading"
	if isHeading && *firstHeading == "" {
		// Extract heading text into a temporary builder while still writing to main output
		var hb strings.Builder
		content, _ := m["content"].([]any)
		for _, c := range content {
			walkNodeHeadingOnly(c, &hb)
		}
		*firstHeading = hb.String()
	}

	// Recurse into children
	if content, _ := m["content"].([]any); len(content) > 0 {
		for _, c := range content {
			walkNode(c, b, firstHeading, inHeading || isHeading)
		}
	}

	afterBlock(typ, b)
}

func walkNodeHeadingOnly(node any, b *strings.Builder) {
	m, ok := node.(map[string]any)
	if !ok {
		return
	}
	typ, _ := m["type"].(string)
	if typ == "text" {
		if s, _ := m["text"].(string); s != "" {
			b.WriteString(s)
			b.WriteString(" ")
		}
		return
	}
	if content, _ := m["content"].([]any); len(content) > 0 {
		for _, c := range content {
			walkNodeHeadingOnly(c, b)
		}
	}
}

func beforeBlock(typ string, b *strings.Builder) {
	switch typ {
	case "paragraph", "heading", "blockquote", "listItem", "taskItem",
		"tableRow", "tableCell", "tableHeader", "codeBlock":
		ensureNewline(b)
	}
}

func afterBlock(typ string, b *strings.Builder) {
	switch typ {
	case "hardBreak":
		b.WriteString("\n")
	case "paragraph", "heading", "blockquote", "listItem", "taskItem",
		"horizontalRule", "tableRow", "codeBlock":
		ensureNewline(b)
	}
}

func ensureNewline(b *strings.Builder) {
	s := b.String()
	if s == "" {
		return
	}
	// avoid piling up many newlines
	if !strings.HasSuffix(s, "\n") {
		b.WriteString("\n")
	}
}

func normalizeWhitespace(s string) string {
	// Collapse runs of whitespace but preserve newlines as separators.
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.Join(strings.Fields(lines[i]), " ")
	}
	out := strings.Join(lines, "\n")
	out = strings.TrimSpace(out)
	// Collapse multiple blank lines
	out = collapseBlankLines(out)
	return out
}

func collapseBlankLines(s string) string {
	var out strings.Builder
	prevNL := 0
	for _, r := range s {
		if r == '\n' {
			prevNL++
			if prevNL <= 2 {
				out.WriteRune(r)
			}
			continue
		}
		if !unicode.IsSpace(r) {
			prevNL = 0
		}
		out.WriteRune(r)
	}
	return strings.TrimSpace(out.String())
}
