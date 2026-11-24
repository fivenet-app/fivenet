package textdiff

import (
	"strings"
	"unicode/utf8"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const unchangedContextLines = 3

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

	collapsedOps := collapseEqualLines(ops)

	return &content.ContentDiff{
		Stats: &content.ContentDiffStats{
			InsertedRunes: ins,
			DeletedRunes:  del,
			OpCount:       utils.ToUint32Saturated(len(collapsedOps)),
		},
		Ops: collapsedOps,
	}
}

// CollapseEqualLines normalizes an activity diff for transport. It is also
// applied when reading stored activity so entries created before collapsing was
// introduced render with the same concise context.
func CollapseEqualLines(diff *content.ContentDiff) {
	if diff == nil {
		return
	}

	diff.Ops = collapseEqualLines(diff.GetOps())
	if diff.Stats != nil {
		diff.Stats.OpCount = utils.ToUint32Saturated(len(diff.GetOps()))
	}
}

// collapseEqualLines keeps git-style context around changes while omitting long
// runs of unchanged lines. The client receives an explicit marker so it does
// not have to reimplement the collapsing logic for every activity entry.
func collapseEqualLines(ops []*content.ContentDiffOp) []*content.ContentDiffOp {
	if len(ops) < 2 {
		return ops
	}

	collapsed := make([]*content.ContentDiffOp, 0, len(ops))
	for i, op := range ops {
		if op.GetKind() != content.Kind_KIND_EQUAL {
			collapsed = append(collapsed, op)
			continue
		}

		lines := splitLines(op.GetText())
		var before, after []string
		var omitted int

		switch {
		case i == 0 && len(lines) > unchangedContextLines:
			omitted = len(lines) - unchangedContextLines
			after = lines[omitted:]

		case i == len(ops)-1 && len(lines) > unchangedContextLines:
			before = lines[:unchangedContextLines]
			omitted = len(lines) - unchangedContextLines

		case i > 0 && i < len(ops)-1 && len(lines) > unchangedContextLines*2:
			before = lines[:unchangedContextLines]
			omitted = len(lines) - unchangedContextLines*2
			after = lines[len(lines)-unchangedContextLines:]

		default:
			collapsed = append(collapsed, op)
			continue
		}

		if len(before) > 0 {
			collapsed = append(collapsed, &content.ContentDiffOp{
				Kind: content.Kind_KIND_EQUAL,
				Text: strings.Join(before, ""),
			})
		}

		collapsed = append(collapsed, &content.ContentDiffOp{
			Kind:           content.Kind_KIND_COLLAPSED,
			CollapsedLines: utils.ToUint32Saturated(omitted),
		})

		if len(after) > 0 {
			collapsed = append(collapsed, &content.ContentDiffOp{
				Kind: content.Kind_KIND_EQUAL,
				Text: strings.Join(after, ""),
			})
		}
	}

	return collapsed
}

func splitLines(text string) []string {
	lines := strings.SplitAfter(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}
	return lines
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimSpace(s)
}
