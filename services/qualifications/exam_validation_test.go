package qualifications

import (
	"testing"

	qualificationspb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckIfUserCanTakeExamRejectsUnavailableQualifications(t *testing.T) {
	t.Parallel()
	server := &Server{}
	user := &userinfo.UserInfo{UserId: 1}

	for _, qualification := range []*qualificationspb.QualificationShort{
		{Closed: true, ExamMode: qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_ENABLED},
		{Draft: true, ExamMode: qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_ENABLED},
		{ExamMode: qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED},
	} {
		canTake, _ := server.checkIfUserCanTakeExam(t.Context(), qualification, user)
		assert.False(t, canTake)
	}
}

func TestNormalizeExamResponsesRejectsDuplicateQuestionIDs(t *testing.T) {
	t.Parallel()
	exam := &qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{{
		Id: 1,
		Data: &qualificationsexam.ExamQuestionData{Data: &qualificationsexam.ExamQuestionData_Yesno{
			Yesno: &qualificationsexam.ExamQuestionYesNo{},
		}},
	}}}
	responses := &qualificationsexam.ExamResponses{Responses: []*qualificationsexam.ExamResponse{
		{
			QuestionId: 1,
			Response: &qualificationsexam.ExamResponseData{
				Response: &qualificationsexam.ExamResponseData_Yesno{
					Yesno: &qualificationsexam.ExamResponseYesNo{},
				},
			},
		},
		{
			QuestionId: 1,
			Response: &qualificationsexam.ExamResponseData{
				Response: &qualificationsexam.ExamResponseData_Yesno{
					Yesno: &qualificationsexam.ExamResponseYesNo{},
				},
			},
		},
	}}

	_, err := normalizeExamResponses(exam, responses, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestNormalizeExamResponsesCanonicalizesQuestionSnapshot(t *testing.T) {
	t.Parallel()
	question := &qualificationsexam.ExamQuestion{
		Id:    1,
		Title: "Authoritative question",
		Answer: &qualificationsexam.ExamQuestionAnswerData{
			AnswerKey: "authoritative answer",
			Answer: &qualificationsexam.ExamQuestionAnswerData_Yesno{
				Yesno: &qualificationsexam.ExamResponseYesNo{Value: true},
			},
		},
		Data: &qualificationsexam.ExamQuestionData{
			Data: &qualificationsexam.ExamQuestionData_Yesno{
				Yesno: &qualificationsexam.ExamQuestionYesNo{},
			},
		},
	}
	response := &qualificationsexam.ExamResponse{
		QuestionId: 1,
		Question:   &qualificationsexam.ExamQuestion{Id: 1, Title: "Client supplied question"},
		Response: &qualificationsexam.ExamResponseData{
			Response: &qualificationsexam.ExamResponseData_Yesno{
				Yesno: &qualificationsexam.ExamResponseYesNo{},
			},
		},
	}

	normalized, err := normalizeExamResponses(
		&qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{question}},
		&qualificationsexam.ExamResponses{Responses: []*qualificationsexam.ExamResponse{response}},
		false,
	)
	require.NoError(t, err)
	require.Len(t, normalized.GetResponses(), 1)
	normalizedQuestion := normalized.GetResponses()[0].GetQuestion()
	assert.NotSame(t, question, normalizedQuestion)
	assert.Equal(t, question.GetId(), normalizedQuestion.GetId())
	assert.Equal(t, question.GetTitle(), normalizedQuestion.GetTitle())
	assert.Nil(t, normalizedQuestion.GetAnswer())
	assert.Equal(t, "authoritative answer", question.GetAnswer().GetAnswerKey())
}

func TestNormalizeExamResponsesValidatesChoiceAndTextConstraints(t *testing.T) {
	t.Parallel()
	exam := &qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{
		{
			Id: 1,
			Data: &qualificationsexam.ExamQuestionData{
				Data: &qualificationsexam.ExamQuestionData_SingleChoice{
					SingleChoice: &qualificationsexam.ExamQuestionSingleChoice{
						Choices: []string{"A"},
					},
				},
			},
		},
		{
			Id: 2,
			Data: &qualificationsexam.ExamQuestionData{
				Data: &qualificationsexam.ExamQuestionData_FreeText{
					FreeText: &qualificationsexam.ExamQuestionText{MinLength: 3, MaxLength: 5},
				},
			},
		},
	}}

	_, err := normalizeExamResponses(
		exam,
		&qualificationsexam.ExamResponses{Responses: []*qualificationsexam.ExamResponse{
			{
				QuestionId: 1,
				Response: &qualificationsexam.ExamResponseData{
					Response: &qualificationsexam.ExamResponseData_SingleChoice{
						SingleChoice: &qualificationsexam.ExamResponseSingleChoice{Choice: "B"},
					},
				},
			},
			{
				QuestionId: 2,
				Response: &qualificationsexam.ExamResponseData{
					Response: &qualificationsexam.ExamResponseData_FreeText{
						FreeText: &qualificationsexam.ExamResponseText{Text: "ab"},
					},
				},
			},
		}},
		false,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid choice")
}

func TestNormalizeExamResponsesRejectsOverlongFreeText(t *testing.T) {
	t.Parallel()
	exam := &qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{{
		Id: 1,
		Data: &qualificationsexam.ExamQuestionData{
			Data: &qualificationsexam.ExamQuestionData_FreeText{
				FreeText: &qualificationsexam.ExamQuestionText{MaxLength: 3},
			},
		},
	}}}
	response := &qualificationsexam.ExamResponseData{
		Response: &qualificationsexam.ExamResponseData_FreeText{
			FreeText: &qualificationsexam.ExamResponseText{Text: "longer"},
		},
	}

	_, err := normalizeExamResponses(
		exam,
		&qualificationsexam.ExamResponses{Responses: []*qualificationsexam.ExamResponse{{
			QuestionId: 1,
			Response:   response,
		}}},
		true,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "too long")
	assert.Equal(t, "longer", response.GetFreeText().GetText())
}
