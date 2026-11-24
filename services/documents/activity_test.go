package documents

import (
	"testing"
	"unicode/utf8"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	resdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGenerateDocumentDiffCreatesActivityPayloadForTiptapContentChanges(t *testing.T) {
	t.Parallel()

	oldContent := testTiptapContent(t, "Shared introduction", "The original paragraph")
	newContent := testTiptapContent(t, "Shared introduction", "The updated paragraph")

	updated, err := (&Server{}).generateDocumentDiff(
		&resdocuments.Document{Content: oldContent},
		&resdocuments.Document{Content: newContent},
	)
	require.NoError(t, err)
	require.NotNil(t, updated.GetContentCdiff())
	assert.True(t, documentUpdatedHasChanges(updated))

	stats := updated.GetContentCdiff().GetStats()
	assert.Positive(t, stats.GetInsertedRunes())
	assert.Positive(t, stats.GetDeletedRunes())
	assert.Less(t, stats.GetInsertedRunes(), uint32(utf8.RuneCountInString(newContent.Extract().GetText())))
	assert.Less(t, stats.GetDeletedRunes(), uint32(utf8.RuneCountInString(oldContent.Extract().GetText())))
	assert.Contains(t, updated.GetContentCdiff().GetOps(), &content.ContentDiffOp{
		Kind: content.Kind_KIND_EQUAL,
		Text: "Shared introduction\nThe ",
	})

	activity := &documentsactivity.DocActivity{
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED,
		Data: &documentsactivity.DocActivityData{
			Data: &documentsactivity.DocActivityData_Updated{Updated: updated},
		},
	}
	assert.True(t, documentUpdatedHasChanges(activity.GetData().GetUpdated()))

	stateChanged, err := (&Server{}).generateDocumentDiff(
		&resdocuments.Document{Meta: &resdocuments.DocumentMeta{State: "draft"}},
		&resdocuments.Document{Meta: &resdocuments.DocumentMeta{State: "published"}},
	)
	require.NoError(t, err)
	assert.NotNil(t, stateChanged.GetStateCdiff())
	assert.True(t, documentUpdatedHasChanges(stateChanged))
}

func TestGenerateDocumentDiffForStructuredTiptapContent(t *testing.T) {
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

			updated, err := (&Server{}).generateDocumentDiff(
				&resdocuments.Document{Content: tt.oldContent},
				&resdocuments.Document{Content: tt.newContent},
			)
			require.NoError(t, err)
			assert.Equal(t, tt.wantActivity, documentUpdatedHasChanges(updated))

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

func TestGenerateDocumentDiffRecordsCaseOnlyMetadataChanges(t *testing.T) {
	t.Parallel()

	updated, err := (&Server{}).generateDocumentDiff(
		&resdocuments.Document{Title: "report", Meta: &resdocuments.DocumentMeta{State: "draft"}},
		&resdocuments.Document{Title: "Report", Meta: &resdocuments.DocumentMeta{State: "Draft"}},
	)
	require.NoError(t, err)
	assert.NotNil(t, updated.GetTitleCdiff())
	assert.NotNil(t, updated.GetStateCdiff())
	assert.True(t, documentUpdatedHasChanges(updated))
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
