package qualificationsexam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrade(t *testing.T) {
	t.Parallel()
	questions := &ExamQuestions{
		Questions: []*ExamQuestion{
			{
				Id: 1,
				Data: &ExamQuestionData{
					Data: &ExamQuestionData_Yesno{
						Yesno: &ExamQuestionYesNo{},
					},
				},
				Answer: &ExamQuestionAnswerData{
					Answer: &ExamQuestionAnswerData_Yesno{
						Yesno: &ExamResponseYesNo{Value: true},
					},
				},
				Points: new(int32(10)),
			},
			{
				Id: 2,
				Data: &ExamQuestionData{
					Data: &ExamQuestionData_SingleChoice{
						SingleChoice: &ExamQuestionSingleChoice{
							Choices: []string{"A", "B", "C"},
						},
					},
				},
				Answer: &ExamQuestionAnswerData{
					Answer: &ExamQuestionAnswerData_SingleChoice{
						SingleChoice: &ExamResponseSingleChoice{Choice: "A"},
					},
				},
				Points: new(int32(20)),
			},
			{
				Id: 3,
				Data: &ExamQuestionData{
					Data: &ExamQuestionData_MultipleChoice{
						MultipleChoice: &ExamQuestionMultipleChoice{
							Choices: []string{"A", "B", "C"},
						},
					},
				},
				Answer: &ExamQuestionAnswerData{
					Answer: &ExamQuestionAnswerData_MultipleChoice{
						MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"A", "B"}},
					},
				},
				Points: new(int32(30)),
			},
		},
	}

	responses := &ExamResponses{
		Responses: []*ExamResponse{
			{
				QuestionId: 1,
				Response: &ExamResponseData{
					Response: &ExamResponseData_Yesno{
						Yesno: &ExamResponseYesNo{Value: true},
					},
				},
			},
			{
				QuestionId: 2,
				Response: &ExamResponseData{
					Response: &ExamResponseData_SingleChoice{
						SingleChoice: &ExamResponseSingleChoice{Choice: "A"},
					},
				},
			},
			{
				QuestionId: 3,
				Response: &ExamResponseData{
					Response: &ExamResponseData_MultipleChoice{
						// Only one of two correct choices is provided in the response
						MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"B"}},
					},
				},
			},
		},
	}

	score, grading := questions.Grade(AutoGradeMode_AUTO_GRADE_MODE_STRICT, responses)
	assert.InEpsilon(t, float32(30), score, 0.0001, "Expected score to be 30")
	assert.NotNil(t, grading, "Expected grading to be not nil")
	assert.Len(
		t,
		grading.GetResponses(),
		len(questions.GetQuestions()),
		"Expected grading responses to be equal to the number of questions",
	)

	score, grading = questions.Grade(AutoGradeMode_AUTO_GRADE_MODE_PARTIAL_CREDIT, responses)
	assert.InEpsilon(t, float32(45), score, 0.0001, "Expected score to be 45 (partial credit)")
	assert.Len(
		t,
		grading.GetResponses(),
		len(questions.GetQuestions()),
		"Expected grading responses to be equal to the number of questions",
	)
}

func TestGradeDoesNotAwardDuplicateOrExtraChoices(t *testing.T) {
	t.Parallel()
	question := &ExamQuestion{
		Id:     1,
		Points: new(int32(10)),
		Data: &ExamQuestionData{Data: &ExamQuestionData_MultipleChoice{
			MultipleChoice: &ExamQuestionMultipleChoice{Choices: []string{"A", "B", "C"}},
		}},
		Answer: &ExamQuestionAnswerData{Answer: &ExamQuestionAnswerData_MultipleChoice{
			MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"A", "B"}},
		}},
	}

	exam := &ExamQuestions{Questions: []*ExamQuestion{question}}
	responses := &ExamResponses{Responses: []*ExamResponse{{
		QuestionId: 1,
		Response: &ExamResponseData{Response: &ExamResponseData_MultipleChoice{
			MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"A", "B", "C"}},
		}},
	}}}

	score, _ := exam.Grade(
		AutoGradeMode_AUTO_GRADE_MODE_STRICT,
		responses,
	)
	assert.Zero(t, score)

	score, grading := exam.Grade(
		AutoGradeMode_AUTO_GRADE_MODE_PARTIAL_CREDIT,
		&ExamResponses{Responses: []*ExamResponse{{
			QuestionId: 1,
			Response: &ExamResponseData{Response: &ExamResponseData_MultipleChoice{
				MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"A", "B", "C", "A"}},
			}},
		}}},
	)
	assert.InEpsilon(t, float32(5), score, 0.0001)
	assert.InEpsilon(t, float32(5), grading.GetResponses()[0].GetPoints(), 0.0001)
}

func TestGradeDeduplicatesAnswerChoicesForPartialCredit(t *testing.T) {
	t.Parallel()
	question := &ExamQuestion{
		Id:     1,
		Points: new(int32(10)),
		Data: &ExamQuestionData{Data: &ExamQuestionData_MultipleChoice{
			MultipleChoice: &ExamQuestionMultipleChoice{Choices: []string{"A", "B"}},
		}},
		Answer: &ExamQuestionAnswerData{Answer: &ExamQuestionAnswerData_MultipleChoice{
			MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"A", "A"}},
		}},
	}

	score, grading := (&ExamQuestions{Questions: []*ExamQuestion{question}}).Grade(
		AutoGradeMode_AUTO_GRADE_MODE_PARTIAL_CREDIT,
		&ExamResponses{Responses: []*ExamResponse{{
			QuestionId: 1,
			Response: &ExamResponseData{Response: &ExamResponseData_MultipleChoice{
				MultipleChoice: &ExamResponseMultipleChoice{Choices: []string{"A"}},
			}},
		}}},
	)
	assert.InEpsilon(t, float32(10), score, 0.0001)
	assert.InEpsilon(t, float32(10), grading.GetResponses()[0].GetPoints(), 0.0001)
}
