package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func mustStruct(t *testing.T, m map[string]any) *structpb.Struct {
	t.Helper()

	s, err := structpb.NewStruct(m)
	require.NoError(t, err)
	return s
}

func TestExtractFromTiptapNil(t *testing.T) {
	t.Parallel()
	got := ExtractFromTiptap(nil)
	require.NotNil(t, got)
	assert.Empty(t, got.Text)
	assert.Equal(t, uint32(0), got.WordCount)
	assert.Empty(t, got.FirstHeading)
}

func TestExtractFromTiptapComplexGolden(t *testing.T) {
	t.Parallel()
	doc := mustStruct(t, map[string]any{
		"type": "doc",
		"content": []any{
			map[string]any{
				"type": "heading",
				"content": []any{
					map[string]any{"type": "text", "text": "Main   Title"},
				},
			},
			map[string]any{
				"type": "paragraph",
				"content": []any{
					map[string]any{"type": "text", "text": "Hello     world"},
				},
			},
			map[string]any{
				"type": "bulletList",
				"content": []any{
					map[string]any{
						"type": "listItem",
						"content": []any{
							map[string]any{
								"type": "paragraph",
								"content": []any{
									map[string]any{"type": "text", "text": "One"},
								},
							},
						},
					},
					map[string]any{
						"type": "listItem",
						"content": []any{
							map[string]any{
								"type": "paragraph",
								"content": []any{
									map[string]any{"type": "text", "text": "Two"},
								},
							},
							map[string]any{
								"type":  "orderedList",
								"attrs": map[string]any{"start": float64(3)},
								"content": []any{
									map[string]any{
										"type": "listItem",
										"content": []any{
											map[string]any{
												"type": "paragraph",
												"content": []any{
													map[string]any{
														"type": "text",
														"text": "Two.Point",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			map[string]any{
				"type": "taskList",
				"content": []any{
					map[string]any{
						"type":  "taskItem",
						"attrs": map[string]any{"checked": true},
						"content": []any{
							map[string]any{
								"type": "paragraph",
								"content": []any{
									map[string]any{"type": "text", "text": "Done"},
								},
							},
						},
					},
					map[string]any{
						"type":  "taskItem",
						"attrs": map[string]any{"checked": false},
						"content": []any{
							map[string]any{
								"type": "paragraph",
								"content": []any{
									map[string]any{"type": "text", "text": "Todo"},
								},
							},
						},
					},
				},
			},
			map[string]any{
				"type": "table",
				"content": []any{
					map[string]any{
						"type": "tableRow",
						"content": []any{
							map[string]any{
								"type": "tableHeader",
								"content": []any{
									map[string]any{"type": "text", "text": "A"},
								},
							},
							map[string]any{
								"type": "tableHeader",
								"content": []any{
									map[string]any{"type": "text", "text": "B"},
								},
							},
						},
					},
					map[string]any{
						"type": "tableRow",
						"content": []any{
							map[string]any{
								"type": "tableCell",
								"content": []any{
									map[string]any{"type": "text", "text": "1"},
								},
							},
							map[string]any{
								"type": "tableCell",
								"content": []any{
									map[string]any{"type": "text", "text": "2"},
								},
							},
						},
					},
				},
			},
		},
	})

	got := ExtractFromTiptap(doc)
	require.NotNil(t, got)

	expected := "Main Title\n" +
		"Hello world\n" +
		"- One\n" +
		"- Two\n" +
		"3. Two.Point\n" +
		"- [x] Done\n" +
		"- [ ] Todo\n" +
		"[Table]\n" +
		"| A | B |\n" +
		"| ---- | ---- |\n" +
		"| 1 | 2 |\n" +
		"[/Table]"

	assert.Equal(t, expected, got.Text)
	assert.Equal(t, "Main Title", got.FirstHeading)
	assert.NotZero(t, got.WordCount)
}

func TestExtractFromTiptapInlineNodesInParagraph(t *testing.T) {
	t.Parallel()
	doc := mustStruct(t, map[string]any{
		"type": "doc",
		"content": []any{
			map[string]any{
				"type": "paragraph",
				"content": []any{
					map[string]any{
						"type":  "mention",
						"attrs": map[string]any{"label": "Ada"},
					},
					map[string]any{"type": "text", "text": " "},
					map[string]any{
						"type":  "mention",
						"attrs": map[string]any{"id": "user-42"},
					},
					map[string]any{"type": "text", "text": " "},
					map[string]any{"type": "mention"},
					map[string]any{"type": "hardBreak"},
					map[string]any{
						"type":  "image",
						"attrs": map[string]any{"alt": "diagram"},
					},
					map[string]any{"type": "text", "text": " "},
					map[string]any{"type": "image"},
				},
			},
		},
	})

	got := ExtractFromTiptap(doc)
	require.NotNil(t, got)

	assert.Equal(t, "Ada @user-42 @mention\n[Image alt=\"diagram\"]\n[Image]", got.Text)
	assert.Empty(t, got.FirstHeading)
}

func TestExtractFromTiptapUnknownAndMalformedNodes(t *testing.T) {
	t.Parallel()
	doc := mustStruct(t, map[string]any{
		"type": "doc",
		"content": []any{
			"not-a-node",
			map[string]any{
				"type": "unknownContainer",
				"content": []any{
					123,
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{"type": "text", "text": "reachable"},
						},
					},
				},
			},
		},
	})

	got := ExtractFromTiptap(doc)
	require.NotNil(t, got)
	assert.Equal(t, "reachable", got.Text)
	assert.Equal(t, uint32(1), got.WordCount)
}

func TestExtractFromTiptapWhitespaceNormalization(t *testing.T) {
	t.Parallel()
	doc := mustStruct(t, map[string]any{
		"type": "doc",
		"content": []any{
			map[string]any{
				"type": "paragraph",
				"content": []any{
					map[string]any{"type": "text", "text": "alpha    beta"},
					map[string]any{"type": "hardBreak"},
					map[string]any{"type": "text", "text": "gamma"},
					map[string]any{"type": "hardBreak"},
					map[string]any{"type": "hardBreak"},
					map[string]any{"type": "hardBreak"},
					map[string]any{"type": "text", "text": "delta"},
				},
			},
		},
	})

	got := ExtractFromTiptap(doc)
	require.NotNil(t, got)

	assert.Equal(t, "alpha beta\ngamma\n\ndelta", got.Text)
	assert.Equal(t, uint32(4), got.WordCount)
}
