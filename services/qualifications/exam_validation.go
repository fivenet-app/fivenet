package qualifications

import (
	"errors"
	"fmt"

	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	"google.golang.org/protobuf/proto"
)

// normalizeExamResponses removes client-controlled question snapshots and accepts
// exactly one type-compatible response for every submitted question ID.
func normalizeExamResponses(
	exam *qualificationsexam.ExamQuestions,
	responses *qualificationsexam.ExamResponses,
	partial bool,
) (*qualificationsexam.ExamResponses, error) {
	questions := make(map[int64]*qualificationsexam.ExamQuestion, len(exam.GetQuestions()))
	for _, question := range exam.GetQuestions() {
		questions[question.GetId()] = question
	}

	normalized := &qualificationsexam.ExamResponses{
		Responses: make([]*qualificationsexam.ExamResponse, 0, len(responses.GetResponses())),
	}
	seen := make(map[int64]struct{}, len(responses.GetResponses()))
	for _, response := range responses.GetResponses() {
		question, ok := questions[response.GetQuestionId()]
		if !ok {
			return nil, fmt.Errorf("unknown exam question %d", response.GetQuestionId())
		}
		if _, duplicate := seen[question.GetId()]; duplicate {
			return nil, fmt.Errorf("duplicate response for exam question %d", question.GetId())
		}
		seen[question.GetId()] = struct{}{}

		if err := validateExamResponse(question, response.GetResponse(), partial); err != nil {
			return nil, err
		}
		questionCopy := proto.Clone(question).(*qualificationsexam.ExamQuestion)
		questionCopy.ClearAnswer()
		normalized.Responses = append(normalized.Responses, &qualificationsexam.ExamResponse{
			QuestionId: question.GetId(),
			Question:   questionCopy,
			Response:   response.GetResponse(),
		})
	}
	if !partial && len(seen) != len(questions) {
		return nil, errors.New("missing exam responses")
	}

	return normalized, nil
}

func examHasFreeText(exam *qualificationsexam.ExamQuestions) bool {
	for _, question := range exam.GetQuestions() {
		if question.GetData().GetFreeText() != nil {
			return true
		}
	}
	return false
}

func validateExamAutoGrading(
	exam *qualificationsexam.ExamQuestions,
	settings *qualificationsexam.QualificationExamSettings,
) error {
	if !settings.GetAutoGrade() {
		return nil
	}
	if examHasFreeText(exam) {
		return errorsqualifications.ErrExamAutoGradingFreeText
	}

	var autoGradePoints int32
	for _, question := range exam.GetQuestions() {
		switch {
		case question.GetData().GetYesno() != nil:
			if question.GetAnswer().GetYesno() == nil {
				return errorsqualifications.ErrExamAutoGradingInvalid
			}
			autoGradePoints += question.GetPoints()
		case question.GetData().GetSingleChoice() != nil:
			answer := question.GetAnswer().GetSingleChoice()
			if answer == nil ||
				!contains(question.GetData().GetSingleChoice().GetChoices(), answer.GetChoice()) {
				return errorsqualifications.ErrExamAutoGradingInvalid
			}
			autoGradePoints += question.GetPoints()
		case question.GetData().GetMultipleChoice() != nil:
			answer := question.GetAnswer().GetMultipleChoice()
			if answer == nil || len(answer.GetChoices()) == 0 ||
				!uniqueSubset(
					answer.GetChoices(),
					question.GetData().GetMultipleChoice().GetChoices(),
				) {
				return errorsqualifications.ErrExamAutoGradingInvalid
			}
			limit := question.GetData().GetMultipleChoice().GetLimit()
			if limit > 0 && len(answer.GetChoices()) > int(limit) {
				return errorsqualifications.ErrExamAutoGradingAnswerLimit
			}
			autoGradePoints += question.GetPoints()
		}
	}
	if settings.GetMinimumPoints() > autoGradePoints {
		return errorsqualifications.ErrExamAutoGradingInvalid
	}
	return nil
}

func mergeExamResponses(
	exam *qualificationsexam.ExamQuestions,
	existing *qualificationsexam.ExamResponses,
	incoming *qualificationsexam.ExamResponses,
) *qualificationsexam.ExamResponses {
	merged := &qualificationsexam.ExamResponses{
		QualificationId: existing.GetQualificationId(),
		UserId:          existing.GetUserId(),
		Responses:       make([]*qualificationsexam.ExamResponse, 0),
	}
	byQuestionID := make(map[int64]int)
	existing = sanitizeExamResponses(exam, existing)
	for _, response := range existing.GetResponses() {
		byQuestionID[response.GetQuestionId()] = len(merged.Responses)
		merged.Responses = append(merged.Responses, response)
	}
	for _, response := range incoming.GetResponses() {
		if index, ok := byQuestionID[response.GetQuestionId()]; ok {
			merged.Responses[index] = response
			continue
		}
		byQuestionID[response.GetQuestionId()] = len(merged.Responses)
		merged.Responses = append(merged.Responses, response)
	}
	return merged
}

// sanitizeExamResponses rebuilds embedded questions from the server-side exam
// and always removes answer keys, including from legacy stored responses.
func sanitizeExamResponses(
	exam *qualificationsexam.ExamQuestions,
	responses *qualificationsexam.ExamResponses,
) *qualificationsexam.ExamResponses {
	if responses == nil {
		return nil
	}
	questions := make(map[int64]*qualificationsexam.ExamQuestion, len(exam.GetQuestions()))
	for _, question := range exam.GetQuestions() {
		questions[question.GetId()] = question
	}
	sanitized := proto.Clone(responses).(*qualificationsexam.ExamResponses)
	for _, response := range sanitized.GetResponses() {
		question, ok := questions[response.GetQuestionId()]
		if !ok {
			response.Question = nil
			continue
		}
		response.Question = proto.Clone(question).(*qualificationsexam.ExamQuestion)
		response.Question.ClearAnswer()
	}
	return sanitized
}

func validateExamResponse(
	question *qualificationsexam.ExamQuestion,
	response *qualificationsexam.ExamResponseData,
	partial bool,
) error {
	if response == nil {
		return fmt.Errorf("missing response for exam question %d", question.GetId())
	}

	switch {
	case question.GetData().GetYesno() != nil:
		if response.GetYesno() == nil {
			return fmt.Errorf("invalid response type for exam question %d", question.GetId())
		}
	case question.GetData().GetFreeText() != nil:
		text := response.GetFreeText()
		if text == nil {
			return fmt.Errorf("invalid response type for exam question %d", question.GetId())
		}
		textLength := utils.ToUint32Saturated(len([]rune(text.GetText())))
		maxLength := question.GetData().GetFreeText().GetMaxLength()
		if maxLength > 0 && textLength > utils.ToUint32Saturated(int(maxLength)) {
			return fmt.Errorf("response for exam question %d is too long", question.GetId())
		}
		minLength := question.GetData().GetFreeText().GetMinLength()
		if !partial && minLength > 0 && textLength < utils.ToUint32Saturated(int(minLength)) {
			return fmt.Errorf("response for exam question %d is too short", question.GetId())
		}
	case question.GetData().GetSingleChoice() != nil:
		choice := response.GetSingleChoice()
		if choice == nil ||
			!contains(question.GetData().GetSingleChoice().GetChoices(), choice.GetChoice()) {
			return fmt.Errorf("invalid choice for exam question %d", question.GetId())
		}
	case question.GetData().GetMultipleChoice() != nil:
		choices := response.GetMultipleChoice()
		if choices == nil ||
			!uniqueSubset(
				choices.GetChoices(),
				question.GetData().GetMultipleChoice().GetChoices(),
			) {
			return fmt.Errorf("invalid choices for exam question %d", question.GetId())
		}
		if limit := question.GetData().
			GetMultipleChoice().
			GetLimit(); limit > 0 &&
			utils.ToUint32Saturated(
				len(choices.GetChoices()),
			) > utils.ToUint32Saturated(
				int(limit),
			) {
			return fmt.Errorf("too many choices for exam question %d", question.GetId())
		}
	case question.GetData().GetSeparator() != nil || question.GetData().GetImage() != nil:
		if response.GetSeparator() == nil {
			return fmt.Errorf("invalid response type for exam question %d", question.GetId())
		}
	default:
		return fmt.Errorf("invalid exam question %d", question.GetId())
	}
	return nil
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func uniqueSubset(values, allowed []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !contains(allowed, value) {
			return false
		}
		if _, duplicate := seen[value]; duplicate {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}
