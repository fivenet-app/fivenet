package qualificationsstore

import (
	"testing"

	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareExamQuestions(t *testing.T) {
	t.Parallel()
	current := []*qualificationsexam.ExamQuestion{{Id: 10}, {Id: 20}}
	incoming := []*qualificationsexam.ExamQuestion{{Id: 20}, {Id: -1}}

	toCreate, toUpdate, toDelete, err := compareExamQuestions(current, incoming)
	require.NoError(t, err)
	assert.Equal(t, []int64{-1}, questionIDs(toCreate))
	assert.Equal(t, []int64{20}, questionIDs(toUpdate))
	assert.Equal(t, []int64{10}, questionIDs(toDelete))
}

func TestCompareExamQuestionsRejectsInvalidPositiveIDs(t *testing.T) {
	t.Parallel()
	current := []*qualificationsexam.ExamQuestion{{Id: 10}}

	_, _, _, err := compareExamQuestions(current, []*qualificationsexam.ExamQuestion{{Id: 11}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown")

	_, _, _, err = compareExamQuestions(
		current,
		[]*qualificationsexam.ExamQuestion{{Id: 10}, {Id: 10}},
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func questionIDs(questions []*qualificationsexam.ExamQuestion) []int64 {
	ids := make([]int64, 0, len(questions))
	for _, question := range questions {
		ids = append(ids, question.GetId())
	}
	return ids
}
