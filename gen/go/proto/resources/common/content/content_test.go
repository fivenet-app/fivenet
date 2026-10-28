package content

import (
	"testing"

	"github.com/stretchr/testify/require"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestContentExtract(t *testing.T) {
	t.Run("nil content", func(t *testing.T) {
		var content *Content

		require.Empty(t, content.Extract().GetText())
	})

	t.Run("tiptap content", func(t *testing.T) {
		content := &Content{
			TiptapJson: mustStructForContentTest(t, map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{"type": "text", "text": "Hello"},
						},
					},
				},
			}),
		}

		require.Equal(t, "Hello", content.Extract().GetText())
	})
}

func mustStructForContentTest(t *testing.T, value map[string]any) *structpb.Struct {
	t.Helper()

	result, err := structpb.NewStruct(value)
	require.NoError(t, err)

	return result
}
