package tiptapsanitizer

import (
	"encoding/json"
	"testing"
)

func TestBuildAllowedContainsAllNodePolicies(t *testing.T) {
	t.Parallel()
	s := buildAllowed()

	nodeTypes := []string{
		NodeTypeDoc,
		NodeTypeText,
		NodeTypeParagraph,
		NodeTypeBlockquote,
		NodeTypeBulletList,
		NodeTypeOrderedList,
		NodeTypeListItem,
		NodeTypeHardBreak,
		NodeTypeHorizontalRule,
		NodeTypeHeading,
		NodeTypeCodeBlock,
		NodeTypeTable,
		NodeTypeTableRow,
		NodeTypeTableHeader,
		NodeTypeTableCell,
		NodeTypeTaskList,
		NodeTypeTaskItem,
		NodeTypeCheckboxStandalone,
		NodeTypeImage,
		NodeTypeMention,
		NodeTypeTemplateVar,
		NodeTypeTemplateBlock,
	}

	for _, typ := range nodeTypes {
		if _, ok := s.Nodes[typ]; !ok {
			t.Fatalf("missing node policy for %q", typ)
		}
	}
}

func TestBuildAllowedContainsAllMarkPolicies(t *testing.T) {
	t.Parallel()
	s := buildAllowed()

	markTypes := []string{
		MarkTypeBold,
		MarkTypeItalic,
		MarkTypeUnderline,
		MarkTypeStrike,
		MarkTypeCode,
		MarkTypeSubscript,
		MarkTypeSuperscript,
		MarkTypeHighlight,
		MarkTypeTextStyle,
		MarkTypeLink,
	}

	for _, typ := range markTypes {
		if _, ok := s.Marks[typ]; !ok {
			t.Fatalf("missing mark policy for %q", typ)
		}
	}
}

func TestNodePoliciesBasicValidation(t *testing.T) {
	t.Parallel()
	s := buildAllowed()

	tests := []struct {
		name  string
		typ   string
		attrs map[string]any
		ok    bool
	}{
		{name: "doc", typ: NodeTypeDoc, attrs: map[string]any{}, ok: true},
		{name: "text", typ: NodeTypeText, attrs: map[string]any{}, ok: true},
		{
			name:  "paragraph",
			typ:   NodeTypeParagraph,
			attrs: map[string]any{"textAlign": "center"},
			ok:    true,
		},
		{name: "blockquote", typ: NodeTypeBlockquote, attrs: map[string]any{}, ok: true},
		{name: "bulletList", typ: NodeTypeBulletList, attrs: map[string]any{}, ok: true},
		{name: "orderedList", typ: NodeTypeOrderedList, attrs: map[string]any{}, ok: true},
		{name: "listItem", typ: NodeTypeListItem, attrs: map[string]any{}, ok: true},
		{name: "hardBreak", typ: NodeTypeHardBreak, attrs: map[string]any{}, ok: true},
		{name: "horizontalRule", typ: NodeTypeHorizontalRule, attrs: map[string]any{}, ok: true},
		{
			name:  "heading",
			typ:   NodeTypeHeading,
			attrs: map[string]any{"level": float64(2), "textAlign": "left"},
			ok:    true,
		},
		{
			name:  "codeBlock",
			typ:   NodeTypeCodeBlock,
			attrs: map[string]any{"language": "go"},
			ok:    true,
		},
		{name: "table", typ: NodeTypeTable, attrs: map[string]any{}, ok: true},
		{name: "tableRow", typ: NodeTypeTableRow, attrs: map[string]any{}, ok: true},
		{name: "tableHeader", typ: NodeTypeTableHeader, attrs: map[string]any{}, ok: true},
		{
			name:  "tableCell",
			typ:   NodeTypeTableCell,
			attrs: map[string]any{"colspan": float64(2), "rowspan": float64(3)},
			ok:    true,
		},
		{name: "taskList", typ: NodeTypeTaskList, attrs: map[string]any{}, ok: true},
		{name: "taskItem", typ: NodeTypeTaskItem, attrs: map[string]any{"checked": true}, ok: true},
		{
			name:  "checkboxStandalone",
			typ:   NodeTypeCheckboxStandalone,
			attrs: map[string]any{"checked": true, "label": "x"},
			ok:    true,
		},
		{
			name: "image-valid",
			typ:  NodeTypeImage,
			attrs: map[string]any{
				"src":    "https://example.com/a.png",
				"width":  float64(100),
				"height": float64(50),
			},
			ok: true,
		},
		{name: "image-missing-src", typ: NodeTypeImage, attrs: map[string]any{}, ok: false},
		{
			name:  "mention-valid-id",
			typ:   NodeTypeMention,
			attrs: map[string]any{"id": "user-42"},
			ok:    true,
		},
		{
			name:  "mention-valid-label",
			typ:   NodeTypeMention,
			attrs: map[string]any{"label": "Ada"},
			ok:    true,
		},
		{
			name:  "mention-invalid-empty",
			typ:   NodeTypeMention,
			attrs: map[string]any{"id": " "},
			ok:    false,
		},
		{
			name:  "templateVar-valid",
			typ:   NodeTypeTemplateVar,
			attrs: map[string]any{"data-template-var": "person.name"},
			ok:    true,
		},
		{
			name:  "templateVar-invalid-empty",
			typ:   NodeTypeTemplateVar,
			attrs: map[string]any{"data-template-var": " "},
			ok:    false,
		},
		{
			name:  "templateBlock-valid",
			typ:   NodeTypeTemplateBlock,
			attrs: map[string]any{"data-template-block": "if person.active"},
			ok:    true,
		},
		{
			name:  "templateBlock-invalid-empty",
			typ:   NodeTypeTemplateBlock,
			attrs: map[string]any{"data-template-block": " "},
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			policy := s.Nodes[tt.typ]
			ok, _ := policy.Validate(tt.attrs)
			if ok != tt.ok {
				t.Fatalf("policy %q validate ok=%v, want %v", tt.typ, ok, tt.ok)
			}
		})
	}
}

func TestMarkPoliciesBasicValidation(t *testing.T) {
	t.Parallel()
	s := buildAllowed()

	tests := []struct {
		name  string
		typ   string
		attrs map[string]any
		ok    bool
	}{
		{name: "bold", typ: MarkTypeBold, attrs: map[string]any{}, ok: true},
		{name: "italic", typ: MarkTypeItalic, attrs: map[string]any{}, ok: true},
		{name: "underline", typ: MarkTypeUnderline, attrs: map[string]any{}, ok: true},
		{name: "strike", typ: MarkTypeStrike, attrs: map[string]any{}, ok: true},
		{name: "code", typ: MarkTypeCode, attrs: map[string]any{}, ok: true},
		{name: "subscript", typ: MarkTypeSubscript, attrs: map[string]any{}, ok: true},
		{name: "superscript", typ: MarkTypeSuperscript, attrs: map[string]any{}, ok: true},
		{
			name:  "highlight",
			typ:   MarkTypeHighlight,
			attrs: map[string]any{"color": "#FF00AA"},
			ok:    true,
		},
		{
			name:  "textStyle",
			typ:   MarkTypeTextStyle,
			attrs: map[string]any{"color": "red", "backgroundColor": "#00ff00"},
			ok:    true,
		},
		{
			name:  "link-valid-relative",
			typ:   MarkTypeLink,
			attrs: map[string]any{"href": "/local/path"},
			ok:    true,
		},
		{
			name:  "link-valid-https",
			typ:   MarkTypeLink,
			attrs: map[string]any{"href": "https://example.com/path"},
			ok:    true,
		},
		{
			name:  "link-invalid-http",
			typ:   MarkTypeLink,
			attrs: map[string]any{"href": "http://example.com/path"},
			ok:    false,
		},
		{name: "link-invalid-missing", typ: MarkTypeLink, attrs: map[string]any{}, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			policy := s.Marks[tt.typ]
			ok, _ := policy.Validate(tt.attrs)
			if ok != tt.ok {
				t.Fatalf("policy %q validate ok=%v, want %v", tt.typ, ok, tt.ok)
			}
		})
	}
}

func TestSanitizeTextNodeWithMarks(t *testing.T) {
	t.Parallel()
	s := buildAllowed()
	stats := &Stats{}
	in := map[string]any{
		"type": NodeTypeText,
		"text": "hello world",
		"marks": []any{
			map[string]any{"type": MarkTypeBold},
			map[string]any{
				"type":  MarkTypeLink,
				"attrs": map[string]any{"href": "https://example.com"},
			},
			map[string]any{"type": "unknown"},
		},
	}

	out, ok := sanitizeNode(in, s, 0, 8, stats)
	if !ok {
		t.Fatal("sanitizeNode returned not ok")
	}

	if got, _ := out["type"].(string); got != NodeTypeText {
		t.Fatalf("sanitized type = %q, want %q", got, NodeTypeText)
	}

	if got, _ := out["text"].(string); got != "hello world" {
		t.Fatalf("sanitized text = %q, want %q", got, "hello world")
	}

	marks, _ := out["marks"].([]any)
	if len(marks) != 2 {
		t.Fatalf("sanitized marks count = %d, want 2", len(marks))
	}

	if stats.Words != 2 {
		t.Fatalf("stats words = %d, want 2", stats.Words)
	}
}

func TestSanitizeMentionNode(t *testing.T) {
	t.Parallel()
	s := buildAllowed()
	stats := &Stats{}
	in := map[string]any{
		"type":  NodeTypeMention,
		"attrs": map[string]any{"id": "user-42", "label": "Ada"},
	}

	out, ok := sanitizeNode(in, s, 0, 8, stats)
	if !ok {
		t.Fatal("sanitizeNode mention returned not ok")
	}

	if got, _ := out["type"].(string); got != NodeTypeMention {
		t.Fatalf("sanitized type = %q, want %q", got, NodeTypeMention)
	}
	attrs, _ := out["attrs"].(map[string]any)
	if attrs["id"] != "user-42" {
		t.Fatalf("sanitized mention id = %v, want user-42", attrs["id"])
	}
	if attrs["label"] != "Ada" {
		t.Fatalf("sanitized mention label = %v, want Ada", attrs["label"])
	}
}

func TestSanitizeNestedContentAndHeadingExtraction(t *testing.T) {
	t.Parallel()
	New()

	doc := map[string]any{
		"type": NodeTypeDoc,
		"content": []any{
			map[string]any{
				"type": NodeTypeBulletList,
				"content": []any{
					map[string]any{
						"type": NodeTypeListItem,
						"content": []any{
							map[string]any{
								"type": NodeTypeParagraph,
								"content": []any{
									map[string]any{"type": NodeTypeText, "text": "nested item"},
								},
							},
						},
					},
				},
			},
			map[string]any{
				"type":  NodeTypeHeading,
				"attrs": map[string]any{"level": float64(2)},
				"content": []any{
					map[string]any{"type": NodeTypeText, "text": "First"},
					map[string]any{"type": NodeTypeText, "text": "Heading"},
				},
			},
			map[string]any{
				"type":  NodeTypeHeading,
				"attrs": map[string]any{"level": float64(3)},
				"content": []any{
					map[string]any{"type": NodeTypeText, "text": "Second"},
				},
			},
		},
	}

	out, stats, err := Sanitize(doc, 0, 10)
	if err != nil {
		t.Fatalf("sanitize returned error: %v", err)
	}

	if stats.FirstHeading != "First Heading" {
		t.Fatalf("first heading = %q, want %q", stats.FirstHeading, "First Heading")
	}
	if stats.Words != 5 {
		t.Fatalf("word count = %d, want 5", stats.Words)
	}

	content, _ := out["content"].([]any)
	if len(content) != 3 {
		t.Fatalf("root content length = %d, want 3", len(content))
	}

	// Validate nested structure survived sanitize.
	first, _ := content[0].(map[string]any)
	firstContent, _ := first["content"].([]any)
	if len(firstContent) != 1 {
		t.Fatalf("bulletList children = %d, want 1", len(firstContent))
	}
}

func TestSanitizeMaxBytesLimit(t *testing.T) {
	t.Parallel()
	New()

	doc := map[string]any{
		"type": NodeTypeDoc,
		"content": []any{
			map[string]any{
				"type": NodeTypeParagraph,
				"content": []any{
					map[string]any{"type": NodeTypeText, "text": "hello"},
				},
			},
		},
	}

	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal doc failed: %v", err)
	}

	_, _, err = Sanitize(doc, len(b)-1, 10)
	if err == nil {
		t.Fatal("expected document too large error, got nil")
	}
	if err.Error() != "document too large" {
		t.Fatalf("error = %q, want %q", err.Error(), "document too large")
	}
}

func TestSanitizeMaxDepthPrunesTooDeepNodes(t *testing.T) {
	t.Parallel()
	New()

	doc := map[string]any{
		"type": NodeTypeDoc,
		"content": []any{
			map[string]any{
				"type": NodeTypeParagraph,
				"content": []any{
					map[string]any{
						"type": NodeTypeParagraph,
						"content": []any{
							map[string]any{"type": NodeTypeText, "text": "too deep"},
						},
					},
				},
			},
		},
	}

	// maxDepth=2 means depth 3 nodes are dropped.
	out, stats, err := Sanitize(doc, 0, 2)
	if err != nil {
		t.Fatalf("sanitize returned error: %v", err)
	}

	if stats.Words != 0 {
		t.Fatalf("word count = %d, want 0 after depth pruning", stats.Words)
	}

	content, _ := out["content"].([]any)
	if len(content) != 1 {
		t.Fatalf("root content length = %d, want 1", len(content))
	}
	first, _ := content[0].(map[string]any)
	children, _ := first["content"].([]any)
	if len(children) != 1 {
		t.Fatalf("first paragraph children = %d, want 1", len(children))
	}
	second, _ := children[0].(map[string]any)
	if deepChildren, ok := second["content"]; ok && len(deepChildren.([]any)) > 0 {
		t.Fatalf("expected deep content to be pruned, got %v", deepChildren)
	}
}
