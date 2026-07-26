package textdiff

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/stretchr/testify/assert"
)

func TestDiffText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		oldText       string
		newText       string
		expectedOps   []diffOp
		insertedRunes uint32
		deletedRunes  uint32
		hasChanges    bool
	}{
		{
			name:    "insert text",
			oldText: "hello world",
			newText: "hello brave world",
			expectedOps: []diffOp{
				{kind: content.Kind_KIND_EQUAL, text: "hello "},
				{kind: content.Kind_KIND_INSERT, text: "brave "},
				{kind: content.Kind_KIND_EQUAL, text: "world"},
			},
			insertedRunes: 6,
			hasChanges:    true,
		},
		{
			name:    "delete text",
			oldText: "hello brave world",
			newText: "hello world",
			expectedOps: []diffOp{
				{kind: content.Kind_KIND_EQUAL, text: "hello "},
				{kind: content.Kind_KIND_DELETE, text: "brave "},
				{kind: content.Kind_KIND_EQUAL, text: "world"},
			},
			deletedRunes: 6,
			hasChanges:   true,
		},
		{
			name:    "counts utf8 runes",
			oldText: "unit",
			newText: "unit 🚓",
			expectedOps: []diffOp{
				{kind: content.Kind_KIND_EQUAL, text: "unit"},
				{kind: content.Kind_KIND_INSERT, text: " 🚓"},
			},
			insertedRunes: 2,
			hasChanges:    true,
		},
		{
			name:    "normalizes line endings and trims outer whitespace",
			oldText: "  alpha\r\nbeta  ",
			newText: "alpha\nbeta",
			expectedOps: []diffOp{
				{kind: content.Kind_KIND_EQUAL, text: "alpha\nbeta"},
			},
		},
		{
			name:        "empty normalized text",
			oldText:     " \r\n ",
			newText:     "\t",
			expectedOps: []diffOp{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diff := DiffText(tt.oldText, tt.newText)

			assert.Equal(t, tt.expectedOps, diffOps(diff.GetOps()))
			assert.Equal(t, tt.insertedRunes, diff.GetStats().GetInsertedRunes())
			assert.Equal(t, tt.deletedRunes, diff.GetStats().GetDeletedRunes())
			assert.Equal(t, uint32(len(tt.expectedOps)), diff.GetStats().GetOpCount())
			assert.Equal(t, tt.hasChanges, diff.HasChanges())
		})
	}
}

type diffOp struct {
	kind content.Kind
	text string
}

func diffOps(ops []*content.ContentDiffOp) []diffOp {
	out := make([]diffOp, 0, len(ops))
	for _, op := range ops {
		out = append(out, diffOp{
			kind: op.GetKind(),
			text: op.GetText(),
		})
	}
	return out
}
