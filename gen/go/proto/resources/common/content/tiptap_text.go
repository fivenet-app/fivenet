package content

import (
	"fmt"
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
	walkNodeWithIndent(node, b, firstHeading, inHeading, 0, listNone, 0)
}

func walkNodeWithIndent(
	node any,
	b *strings.Builder,
	firstHeading *string,
	inHeading bool,
	indent int,
	lk listKind,
	num int,
) {
	m, ok := node.(map[string]any)
	if !ok {
		return
	}

	typ, _ := m["type"].(string)

	// Text nodes are only meaningful in inline contexts
	if typ == "text" {
		if s, _ := m["text"].(string); s != "" {
			b.WriteString(s)
		}
		return
	}

	// First heading detection (keep your logic; using inline collector)
	isHeading := typ == "heading"
	if isHeading && *firstHeading == "" {
		var hb strings.Builder
		if content, _ := m["content"].([]any); len(content) > 0 {
			for _, c := range content {
				collectInlineText(c, &hb)
			}
		}
		*firstHeading = hb.String()
	}

	switch typ {
	case "table":
		writeTable(m, b, indent)
		return

	case "image":
		// Block-level image
		var ib strings.Builder
		collectInlineText(m, &ib)
		b.WriteString(ib.String())
		b.WriteString("\n")
		return

	case "bulletList":
		writeBulletList(m, b, indent, firstHeading)
		return

	case "orderedList":
		writeOrderedList(m, b, indent, firstHeading)
		return

	case "listItem":
		writeListItem(m, b, firstHeading, indent, lk, num)
		return

	case "taskList":
		writeTaskList(m, b, indent, firstHeading)
		return

	case "taskItem":
		writeTaskItem(m, b, firstHeading, indent)
		return
	}

	// Default: your existing block separators + recursion
	beforeBlock(typ, b)

	if content, _ := m["content"].([]any); len(content) > 0 {
		for _, c := range content {
			walkNodeWithIndent(c, b, firstHeading, inHeading || isHeading, indent, lk, num)
		}
	}

	afterBlock(typ, b)
}

func collectInlineText(node any, b *strings.Builder) {
	m, ok := node.(map[string]any)
	if !ok {
		return
	}

	typ, _ := m["type"].(string)

	switch typ {
	case "text":
		if s, _ := m["text"].(string); s != "" {
			b.WriteString(s)
		}
		return

	case "hardBreak":
		b.WriteString("\n")
		return

	case "mention":
		attrs, _ := m["attrs"].(map[string]any)
		if label, _ := attrs["label"].(string); label != "" {
			b.WriteString(label)
		} else if id, _ := attrs["id"].(string); id != "" {
			b.WriteString("@")
			b.WriteString(id)
		} else {
			b.WriteString("@mention")
		}
		return

	case "image":
		// Inline image token (stable & readable)
		attrs, _ := m["attrs"].(map[string]any)
		alt, _ := attrs["alt"].(string)

		// Avoid dumping huge URLs into diffs unless you really want that:
		if alt != "" {
			b.WriteString("[Image alt=\"")
			b.WriteString(alt)
			b.WriteString("\"]")
		} else {
			b.WriteString("[Image]")
		}
		return
	}

	// default recurse
	if content, _ := m["content"].([]any); len(content) > 0 {
		for _, c := range content {
			collectInlineText(c, b)
		}
	}
}

func writeTable(node map[string]any, b *strings.Builder, indent int) {
	pad := strings.Repeat("  ", indent)

	b.WriteString(pad)
	b.WriteString("[Table]\n")

	content, _ := node["content"].([]any) // tableRow[]
	if len(content) == 0 {
		b.WriteString(pad)
		b.WriteString("[/Table]\n")
		return
	}

	rows := make([][]string, 0, len(content))
	maxCols := 0
	isHeader := false

	for rowIdx, r := range content {
		rm, ok := r.(map[string]any)
		if !ok {
			continue
		}
		if t, _ := rm["type"].(string); t != "tableRow" {
			continue
		}

		cellsAny, _ := rm["content"].([]any) // tableCell/tableHeader[]
		row := make([]string, 0, len(cellsAny))

		for _, c := range cellsAny {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}

			// detect header if first row contains any tableHeader
			if rowIdx == 0 {
				if ct, _ := cm["type"].(string); ct == "tableHeader" {
					isHeader = true
				}
			}

			var cellText strings.Builder
			if cc, _ := cm["content"].([]any); len(cc) > 0 {
				for _, inner := range cc {
					collectInlineText(inner, &cellText)
				}
			}
			row = append(row, strings.TrimSpace(cellText.String()))
		}

		if len(row) > maxCols {
			maxCols = len(row)
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		b.WriteString(pad)
		b.WriteString("[/Table]\n")
		return
	}

	// pad to max cols for stable output
	for i := range rows {
		for len(rows[i]) < maxCols {
			rows[i] = append(rows[i], "")
		}
	}

	renderRow := func(row []string) {
		b.WriteString(pad)
		b.WriteString("| ")
		for i, cell := range row {
			if i > 0 {
				b.WriteString(" | ")
			}
			if cell == "" {
				b.WriteString(" ")
			} else {
				b.WriteString(cell)
			}
		}
		b.WriteString(" |\n")
	}

	renderRow(rows[0])

	if isHeader {
		b.WriteString(pad)
		b.WriteString("|")
		for i := 0; i < maxCols; i++ {
			b.WriteString(" ---- |")
		}
		b.WriteString("\n")

		for _, row := range rows[1:] {
			renderRow(row)
		}
	} else {
		for _, row := range rows[1:] {
			renderRow(row)
		}
	}

	b.WriteString(pad)
	b.WriteString("[/Table]\n")
}

func writeBulletList(node map[string]any, b *strings.Builder, indent int, firstHeading *string) {
	// Each listItem becomes "- <text>"
	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			walkNodeWithIndent(c, b, firstHeading, false, indent, listBullet, 1)
		}
	}
}

func writeOrderedList(node map[string]any, b *strings.Builder, indent int, firstHeading *string) {
	start := 1
	if attrs, _ := node["attrs"].(map[string]any); attrs != nil {
		if s, ok := attrs["start"].(float64); ok { // JSON numbers decode as float64
			start = int(s)
		}
	}
	i := 0
	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			walkNodeWithIndent(c, b, firstHeading, false, indent, listNumber, start+i)
			i++
		}
	}
}

func writeListItem(
	node map[string]any,
	b *strings.Builder,
	firstHeading *string,
	indent int,
	lk listKind,
	num int,
) {
	// Collect text for this item (paragraphs etc.)
	var line strings.Builder

	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			ct, _ := cm["type"].(string)

			// Nested lists handled later; for line text, collect inline from other children
			switch ct {
			case "bulletList", "orderedList", "taskList":
				// skip in line text
			default:
				collectInlineText(c, &line)
			}
		}
	}

	prefix := strings.Repeat("  ", indent)
	switch lk {
	case listNumber:
		prefix += fmt.Sprintf("%d. ", num)
	default:
		prefix += "- "
	}

	b.WriteString(prefix)
	b.WriteString(strings.TrimSpace(line.String()))
	b.WriteString("\n")

	// Now render nested lists under this item (one indent deeper)
	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			ct, _ := cm["type"].(string)
			switch ct {
			case "bulletList":
				writeBulletList(cm, b, indent+1, firstHeading)
			case "orderedList":
				writeOrderedList(cm, b, indent+1, firstHeading)
			case "taskList":
				writeTaskList(cm, b, indent+1, firstHeading)
			}
		}
	}
}

func writeTaskList(node map[string]any, b *strings.Builder, indent int, firstHeading *string) {
	// taskList contains taskItem children
	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			walkNodeWithIndent(c, b, firstHeading, false, indent, listTask, 0)
		}
	}
}

func writeTaskItem(node map[string]any, b *strings.Builder, firstHeading *string, indent int) {
	checked := false
	if attrs, _ := node["attrs"].(map[string]any); attrs != nil {
		// JSON bool stays bool
		if v, ok := attrs["checked"].(bool); ok {
			checked = v
		}
	}

	var line strings.Builder
	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			ct, _ := cm["type"].(string)
			switch ct {
			case "bulletList", "orderedList", "taskList":
				// nested lists handled later
			default:
				collectInlineText(c, &line)
			}
		}
	}

	box := "[ ]"
	if checked {
		box = "[x]"
	}

	prefix := strings.Repeat("  ", indent) + "- " + box + " "
	b.WriteString(prefix)
	b.WriteString(strings.TrimSpace(line.String()))
	b.WriteString("\n")

	// nested lists under task item
	if content, _ := node["content"].([]any); len(content) > 0 {
		for _, c := range content {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			ct, _ := cm["type"].(string)
			switch ct {
			case "bulletList":
				writeBulletList(cm, b, indent+1, firstHeading)
			case "orderedList":
				writeOrderedList(cm, b, indent+1, firstHeading)
			case "taskList":
				writeTaskList(cm, b, indent+1, firstHeading)
			}
		}
	}
}

type listKind int

const (
	listNone listKind = iota
	listBullet
	listNumber
	listTask
)

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

	// Avoid piling up many newlines
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
