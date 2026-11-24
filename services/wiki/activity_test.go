package wiki

import (
	"testing"
	"unicode/utf8"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGeneratePageDiffUsesExtractedTextForTiptapContent(t *testing.T) {
	t.Parallel()

	oldContent := testTiptapContent(t, "Shared introduction", "The original paragraph")
	newContent := testTiptapContent(t, "Shared introduction", "The updated paragraph")

	updated, err := (&Server{}).generatePageDiff(
		&reswiki.Page{Content: oldContent},
		&reswiki.Page{Content: newContent},
	)
	require.NoError(t, err)
	require.NotNil(t, updated.GetContentCdiff())
	assert.True(t, pageUpdatedHasChanges(updated))

	stats := updated.GetContentCdiff().GetStats()
	assert.Positive(t, stats.GetInsertedRunes())
	assert.Positive(t, stats.GetDeletedRunes())
	assert.Less(
		t,
		stats.GetInsertedRunes(),
		uint32(utf8.RuneCountInString(newContent.Extract().GetText())),
	)
	assert.Less(
		t,
		stats.GetDeletedRunes(),
		uint32(utf8.RuneCountInString(oldContent.Extract().GetText())),
	)
	assert.Contains(t, updated.GetContentCdiff().GetOps(), &content.ContentDiffOp{
		Kind: content.Kind_KIND_EQUAL,
		Text: "Shared introduction\nThe ",
	})

	activity := &wikiactivity.PageActivity{
		ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_UPDATED,
		Data: &wikiactivity.PageActivityData{
			Data: &wikiactivity.PageActivityData_Updated{Updated: updated},
		},
	}
	assert.True(t, pageUpdatedHasChanges(activity.GetData().GetUpdated()))

	unchanged, err := (&Server{}).generatePageDiff(
		&reswiki.Page{Content: oldContent},
		&reswiki.Page{Content: oldContent},
	)
	require.NoError(t, err)
	assert.False(t, pageUpdatedHasChanges(unchanged))
}

func TestGeneratePageDiffForStructuredTiptapContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		oldContent   *content.Content
		newContent   *content.Content
		wantActivity bool
	}{
		{
			name:         "table cell edit",
			oldContent:   testTableContent(t, "Original", 120),
			newContent:   testTableContent(t, "Updated", 120),
			wantActivity: true,
		},
		{
			name:         "table presentation only",
			oldContent:   testTableContent(t, "Unchanged", 120),
			newContent:   testTableContent(t, "Unchanged", 240),
			wantActivity: false,
		},
		{
			name:         "map location edit",
			oldContent:   testCustomBlockContent(t, 100, 200, 2),
			newContent:   testCustomBlockContent(t, 100, 250, 2),
			wantActivity: true,
		},
		{
			name:         "map presentation only",
			oldContent:   testCustomBlockContent(t, 100, 200, 2),
			newContent:   testCustomBlockContent(t, 100, 200, 5),
			wantActivity: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			updated, err := (&Server{}).generatePageDiff(
				&reswiki.Page{Content: tt.oldContent},
				&reswiki.Page{Content: tt.newContent},
			)
			require.NoError(t, err)
			assert.Equal(t, tt.wantActivity, pageUpdatedHasChanges(updated))

			if !tt.wantActivity {
				assert.Nil(t, updated.GetContentCdiff())
				return
			}

			require.NotNil(t, updated.GetContentCdiff())
			assert.Less(t, updated.GetContentCdiff().GetStats().GetInsertedRunes(),
				uint32(utf8.RuneCountInString(tt.newContent.Extract().GetText())))
		})
	}
}

func TestGeneratePageDiffRecordsCaseOnlyMetadataChanges(t *testing.T) {
	t.Parallel()

	updated, err := (&Server{}).generatePageDiff(
		&reswiki.Page{Meta: &reswiki.PageMeta{Title: "guide", Description: "overview"}},
		&reswiki.Page{Meta: &reswiki.PageMeta{Title: "Guide", Description: "Overview"}},
	)
	require.NoError(t, err)
	assert.NotNil(t, updated.GetTitleCdiff())
	assert.NotNil(t, updated.GetDescriptionCdiff())
	assert.True(t, pageUpdatedHasChanges(updated))
}

func testTiptapContent(t *testing.T, paragraphs ...string) *content.Content {
	t.Helper()

	nodes := make([]any, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		nodes = append(nodes, map[string]any{
			"type": "paragraph",
			"content": []any{
				map[string]any{"type": "text", "text": paragraph},
			},
		})
	}

	return testTiptapContentWithNodes(t, nodes...)
}

func testTableContent(t *testing.T, value string, columnWidth int) *content.Content {
	t.Helper()

	return testTiptapContentWithNodes(t, map[string]any{
		"type": "table",
		"content": []any{
			map[string]any{
				"type": "tableRow",
				"content": []any{
					map[string]any{
						"type":    "tableHeader",
						"content": []any{map[string]any{"type": "text", "text": "Status"}},
					},
				},
			},
			map[string]any{
				"type": "tableRow",
				"content": []any{
					map[string]any{
						"type":    "tableCell",
						"attrs":   map[string]any{"colwidth": []any{columnWidth}},
						"content": []any{map[string]any{"type": "text", "text": value}},
					},
				},
			},
		},
	})
}

func testCustomBlockContent(t *testing.T, x, y, zoom int) *content.Content {
	t.Helper()

	return testTiptapContentWithNodes(t,
		map[string]any{
			"type": "paragraph",
			"content": []any{
				map[string]any{"type": "text", "text": "Location: "},
				map[string]any{
					"type": "mapBlock",
					"attrs": map[string]any{
						"postal": "Davis",
						"x":      x,
						"y":      y,
						"zoom":   zoom,
						"layer":  "satellite",
					},
				},
			},
		},
		map[string]any{"type": "penaltyCalculator"},
	)
}

func testTiptapContentWithNodes(t *testing.T, nodes ...any) *content.Content {
	t.Helper()

	doc, err := structpb.NewStruct(map[string]any{
		"type":    "doc",
		"content": nodes,
	})
	require.NoError(t, err)

	return &content.Content{
		ContentType: content.ContentType_CONTENT_TYPE_TIPTAP_JSON,
		TiptapJson:  doc,
	}
}
