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

func TestCollapseEqualLines(t *testing.T) {
	t.Parallel()

	ops := collapseEqualLines([]*content.ContentDiffOp{
		{Kind: content.Kind_KIND_EQUAL, Text: "a\nb\nc\nd\ne\nf\n"},
		{Kind: content.Kind_KIND_INSERT, Text: "changed"},
		{Kind: content.Kind_KIND_EQUAL, Text: "g\nh\ni\nj\nk\nl\n"},
	})

	assert.Equal(t, []*content.ContentDiffOp{
		{Kind: content.Kind_KIND_COLLAPSED, CollapsedLines: 3},
		{Kind: content.Kind_KIND_EQUAL, Text: "d\ne\nf\n"},
		{Kind: content.Kind_KIND_INSERT, Text: "changed"},
		{Kind: content.Kind_KIND_EQUAL, Text: "g\nh\ni\n"},
		{Kind: content.Kind_KIND_COLLAPSED, CollapsedLines: 3},
	}, ops)
}

func TestCollapseEqualLinesKeepsShortContext(t *testing.T) {
	t.Parallel()

	ops := []*content.ContentDiffOp{
		{Kind: content.Kind_KIND_EQUAL, Text: "a\nb\nc\n"},
		{Kind: content.Kind_KIND_INSERT, Text: "changed"},
		{Kind: content.Kind_KIND_EQUAL, Text: "d\ne\nf"},
	}

	assert.Equal(t, ops, collapseEqualLines(ops))
}

func TestCollapseEqualLinesHandlesDeletionAndMultipleHunks(t *testing.T) {
	t.Parallel()

	ops := collapseEqualLines([]*content.ContentDiffOp{
		{Kind: content.Kind_KIND_EQUAL, Text: "a\nb\nc\nd\ne\nf\n"},
		{Kind: content.Kind_KIND_DELETE, Text: "old"},
		{Kind: content.Kind_KIND_EQUAL, Text: "g\nh\ni\nj\nk\nl\nm\nn\n"},
		{Kind: content.Kind_KIND_INSERT, Text: "new"},
		{Kind: content.Kind_KIND_EQUAL, Text: "o\np\nq\nr\ns\nt\n"},
	})

	assert.Equal(t, []*content.ContentDiffOp{
		{Kind: content.Kind_KIND_COLLAPSED, CollapsedLines: 3},
		{Kind: content.Kind_KIND_EQUAL, Text: "d\ne\nf\n"},
		{Kind: content.Kind_KIND_DELETE, Text: "old"},
		{Kind: content.Kind_KIND_EQUAL, Text: "g\nh\ni\n"},
		{Kind: content.Kind_KIND_COLLAPSED, CollapsedLines: 2},
		{Kind: content.Kind_KIND_EQUAL, Text: "l\nm\nn\n"},
		{Kind: content.Kind_KIND_INSERT, Text: "new"},
		{Kind: content.Kind_KIND_EQUAL, Text: "o\np\nq\n"},
		{Kind: content.Kind_KIND_COLLAPSED, CollapsedLines: 3},
	}, ops)
}

func TestCollapseEqualLinesIsIdempotent(t *testing.T) {
	t.Parallel()

	first := collapseEqualLines([]*content.ContentDiffOp{
		{Kind: content.Kind_KIND_EQUAL, Text: "a\nb\nc\nd\ne\nf\n"},
		{Kind: content.Kind_KIND_INSERT, Text: "changed"},
		{Kind: content.Kind_KIND_EQUAL, Text: "g\nh\ni\nj\nk\nl\n"},
	})

	assert.Equal(t, first, collapseEqualLines(first))
}

func TestDiffTextCollapsesUnchangedLineRuns(t *testing.T) {
	t.Parallel()

	diff := DiffText(
		"first\none\ntwo\nthree\nfour\nfive\nsix\nseven\nbefore\nlast",
		"first\none\ntwo\nthree\nfour\nfive\nsix\nseven\nafter\nlast",
	)

	assert.Equal(t, uint32(len(diff.GetOps())), diff.GetStats().GetOpCount())
	assert.Contains(t, diff.GetOps(), &content.ContentDiffOp{
		Kind:           content.Kind_KIND_COLLAPSED,
		CollapsedLines: 5,
	})
}

func TestCollapseEqualLinesUpdatesOperationCount(t *testing.T) {
	t.Parallel()

	diff := &content.ContentDiff{
		Stats: &content.ContentDiffStats{},
		Ops: []*content.ContentDiffOp{
			{Kind: content.Kind_KIND_EQUAL, Text: "a\nb\nc\nd\ne\nf\n"},
			{Kind: content.Kind_KIND_INSERT, Text: "changed"},
		},
	}

	CollapseEqualLines(diff)

	assert.Equal(t, uint32(len(diff.GetOps())), diff.GetStats().GetOpCount())
	assert.Contains(t, diff.GetOps(), &content.ContentDiffOp{
		Kind:           content.Kind_KIND_COLLAPSED,
		CollapsedLines: 3,
	})
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
