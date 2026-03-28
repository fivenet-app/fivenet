package textdiff

import (
	"strings"
	"unicode/utf8"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func DiffText(oldText string, newText string) *content.ContentDiff {
	oldText = normalize(oldText)
	newText = normalize(newText)

	dmp := diffmatchpatch.New()

	// Main diff
	diffs := dmp.DiffMain(oldText, newText, false)

	// Optional but recommended:
	// improves human readability (word/line alignment)
	dmp.DiffCleanupSemantic(diffs)

	ops := make([]*content.ContentDiffOp, 0, len(diffs))

	var ins, del uint32
	for _, d := range diffs {
		if d.Text == "" {
			continue
		}

		kind := content.Kind_KIND_UNSPECIFIED
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			kind = content.Kind_KIND_EQUAL
		case diffmatchpatch.DiffInsert:
			kind = content.Kind_KIND_INSERT
			ins = utils.SaturatingAddUint32(ins, utf8.RuneCountInString(d.Text))
		case diffmatchpatch.DiffDelete:
			kind = content.Kind_KIND_DELETE
			del = utils.SaturatingAddUint32(del, utf8.RuneCountInString(d.Text))
		}

		ops = append(ops, &content.ContentDiffOp{
			Kind: kind,
			Text: d.Text,
		})
	}

	return &content.ContentDiff{
		Stats: &content.ContentDiffStats{
			InsertedRunes: ins,
			DeletedRunes:  del,
			OpCount:       utils.ToUint32Saturated(len(ops)),
		},
		Ops: ops,
	}
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimSpace(s)
}
