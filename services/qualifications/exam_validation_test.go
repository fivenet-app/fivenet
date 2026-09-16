package qualifications

import (
	"context"
	"testing"

	qualificationspb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type effectiveExamTestStore struct {
	qualificationsstore.IStore

	exam            *qualificationsexam.ExamQuestions
	qualificationID int64
	withAnswers     bool
}

func (s *effectiveExamTestStore) GetExamQuestions(
	_ context.Context,
	_ qrm.DB,
	qualificationID int64,
	withAnswers bool,
) (*qualificationsexam.ExamQuestions, error) {
	s.qualificationID = qualificationID
	s.withAnswers = withAnswers
	return s.exam, nil
}

func TestEffectiveExamForAutoGradingLoadsPersistedQuestionsWhenOmitted(t *testing.T) {
	t.Parallel()

	persistedExam := &qualificationsexam.ExamQuestions{
		Questions: []*qualificationsexam.ExamQuestion{
			{
				Data: &qualificationsexam.ExamQuestionData{
					Data: &qualificationsexam.ExamQuestionData_FreeText{
						FreeText: &qualificationsexam.ExamQuestionText{},
					},
				},
			},
		},
	}
	store := &effectiveExamTestStore{exam: persistedExam}
	server := &Server{store: store}

	effectiveExam, err := server.effectiveExamForAutoGrading(t.Context(), 42, nil)
	require.NoError(t, err)
	assert.Same(t, persistedExam, effectiveExam)
	assert.Equal(t, int64(42), store.qualificationID)
	assert.True(t, store.withAnswers)
	assert.ErrorIs(
		t,
		validateExamAutoGrading(effectiveExam, &qualificationsexam.QualificationExamSettings{
			AutoGrade: true,
		}),
		errorsqualifications.ErrExamAutoGradingFreeText,
	)
}

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

func TestValidateExamAutoGradingRejectsInvalidConfigurations(t *testing.T) {
	t.Parallel()
	settings := &qualificationsexam.QualificationExamSettings{AutoGrade: true}

	err := validateExamAutoGrading(
		&qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{
			{
				Data: &qualificationsexam.ExamQuestionData{
					Data: &qualificationsexam.ExamQuestionData_FreeText{
						FreeText: &qualificationsexam.ExamQuestionText{},
					},
				},
			},
		}},
		settings,
	)
	require.ErrorIs(t, err, errorsqualifications.ErrExamAutoGradingFreeText)

	err = validateExamAutoGrading(
		&qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{
			{
				Data: &qualificationsexam.ExamQuestionData{
					Data: &qualificationsexam.ExamQuestionData_Yesno{
						Yesno: &qualificationsexam.ExamQuestionYesNo{},
					},
				},
			},
		}},
		settings,
	)
	require.ErrorIs(t, err, errorsqualifications.ErrExamAutoGradingInvalid)

	settings.MinimumPoints = 11
	err = validateExamAutoGrading(
		&qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{
			{
				Points: new(int32(10)),
				Data: &qualificationsexam.ExamQuestionData{
					Data: &qualificationsexam.ExamQuestionData_Yesno{
						Yesno: &qualificationsexam.ExamQuestionYesNo{},
					},
				},
				Answer: &qualificationsexam.ExamQuestionAnswerData{
					Answer: &qualificationsexam.ExamQuestionAnswerData_Yesno{
						Yesno: &qualificationsexam.ExamResponseYesNo{},
					},
				},
			},
		}},
		settings,
	)
	require.ErrorIs(t, err, errorsqualifications.ErrExamAutoGradingInvalid)
}

func TestValidateExamAutoGradingRejectsMismatchedAnswers(t *testing.T) {
	t.Parallel()
	settings := &qualificationsexam.QualificationExamSettings{AutoGrade: true}

	for _, question := range []*qualificationsexam.ExamQuestion{
		{
			Data: &qualificationsexam.ExamQuestionData{Data: &qualificationsexam.ExamQuestionData_SingleChoice{
				SingleChoice: &qualificationsexam.ExamQuestionSingleChoice{Choices: []string{"A", "B"}},
			}},
			Answer: &qualificationsexam.ExamQuestionAnswerData{Answer: &qualificationsexam.ExamQuestionAnswerData_SingleChoice{
				SingleChoice: &qualificationsexam.ExamResponseSingleChoice{Choice: "C"},
			}},
		},
		{
			Data: &qualificationsexam.ExamQuestionData{Data: &qualificationsexam.ExamQuestionData_MultipleChoice{
				MultipleChoice: &qualificationsexam.ExamQuestionMultipleChoice{Choices: []string{"A", "B"}},
			}},
			Answer: &qualificationsexam.ExamQuestionAnswerData{Answer: &qualificationsexam.ExamQuestionAnswerData_MultipleChoice{
				MultipleChoice: &qualificationsexam.ExamResponseMultipleChoice{Choices: []string{"A", "C"}},
			}},
		},
	} {
		err := validateExamAutoGrading(
			&qualificationsexam.ExamQuestions{
				Questions: []*qualificationsexam.ExamQuestion{question},
			},
			settings,
		)
		require.ErrorIs(t, err, errorsqualifications.ErrExamAutoGradingInvalid)
	}
}

func TestValidateExamAutoGradingRejectsAnswerKeyExceedingLimit(t *testing.T) {
	t.Parallel()

	err := validateExamAutoGrading(
		&qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{
			{
				Data: &qualificationsexam.ExamQuestionData{
					Data: &qualificationsexam.ExamQuestionData_MultipleChoice{
						MultipleChoice: &qualificationsexam.ExamQuestionMultipleChoice{
							Choices: []string{"A", "B"},
							Limit:   new(int32(1)),
						},
					},
				},
				Answer: &qualificationsexam.ExamQuestionAnswerData{
					Answer: &qualificationsexam.ExamQuestionAnswerData_MultipleChoice{
						MultipleChoice: &qualificationsexam.ExamResponseMultipleChoice{
							Choices: []string{"A", "B"},
						},
					},
				},
			},
		}},
		&qualificationsexam.QualificationExamSettings{AutoGrade: true},
	)
	require.ErrorIs(t, err, errorsqualifications.ErrExamAutoGradingAnswerLimit)
}

func TestExamForCandidateRedactsAnswer(t *testing.T) {
	t.Parallel()
	exam := &qualificationsexam.ExamQuestions{Questions: []*qualificationsexam.ExamQuestion{{
		Id: 1,
		Answer: &qualificationsexam.ExamQuestionAnswerData{
			AnswerKey: "yes",
			Answer: &qualificationsexam.ExamQuestionAnswerData_Yesno{
				Yesno: &qualificationsexam.ExamResponseYesNo{Value: true},
			},
		},
	}}}

	candidateExam := examForCandidate(exam)
	require.Nil(t, candidateExam.GetQuestions()[0].GetAnswer())
	require.NotNil(t, exam.GetQuestions()[0].GetAnswer())
}

func TestMergeExamResponsesPreservesUnchangedResponses(t *testing.T) {
	t.Parallel()
	existingResponse := &qualificationsexam.ExamResponse{
		QuestionId: 1,
		Question: &qualificationsexam.ExamQuestion{
			Id: 999,
			Answer: &qualificationsexam.ExamQuestionAnswerData{
				AnswerKey: "legacy-answer",
			},
		},
	}
	updatedResponse := &qualificationsexam.ExamResponse{QuestionId: 2}
	existing := &qualificationsexam.ExamResponses{
		QualificationId: 42,
		UserId:          7,
		Responses:       []*qualificationsexam.ExamResponse{existingResponse, {QuestionId: 2}},
	}
	incoming := &qualificationsexam.ExamResponses{
		Responses: []*qualificationsexam.ExamResponse{updatedResponse},
	}

	merged := mergeExamResponses(&qualificationsexam.ExamQuestions{
		Questions: []*qualificationsexam.ExamQuestion{{Id: 1}, {Id: 2}},
	}, existing, incoming)
	require.Len(t, merged.GetResponses(), 2)
	assert.NotSame(t, existingResponse, merged.GetResponses()[0])
	assert.Equal(t, int64(1), merged.GetResponses()[0].GetQuestion().GetId())
	assert.Nil(t, merged.GetResponses()[0].GetQuestion().GetAnswer())
	assert.Same(t, updatedResponse, merged.GetResponses()[1])
	assert.Equal(t, int64(42), merged.GetQualificationId())
	assert.Equal(t, int32(7), merged.GetUserId())
}

func TestPublicExamResponsesRedactsLegacyAnswer(t *testing.T) {
	t.Parallel()
	exam := &qualificationsexam.ExamQuestions{
		Questions: []*qualificationsexam.ExamQuestion{{Id: 1}},
	}
	responses := &qualificationsexam.ExamResponses{
		AttemptId: "attempt",
		Responses: []*qualificationsexam.ExamResponse{
			{
				QuestionId: 1,
				Question: &qualificationsexam.ExamQuestion{
					Id:     1,
					Answer: &qualificationsexam.ExamQuestionAnswerData{AnswerKey: "secret"},
				},
			},
		},
	}

	public := publicExamResponses(exam, responses)
	require.Nil(t, public.GetResponses()[0].GetQuestion().GetAnswer())
	assert.Empty(t, public.GetAttemptId())
	assert.NotNil(t, responses.GetResponses()[0].GetQuestion().GetAnswer())
}
