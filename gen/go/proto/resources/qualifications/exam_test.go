package qualifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestGrade(t *testing.T) {
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
				Points: proto.Int32(10),
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
				Points: proto.Int32(20),
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
				Points: proto.Int32(30),
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
