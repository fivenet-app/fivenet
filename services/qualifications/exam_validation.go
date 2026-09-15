package qualifications

import (
	"fmt"

	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
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

	normalized := &qualificationsexam.ExamResponses{Responses: make([]*qualificationsexam.ExamResponse, 0, len(responses.GetResponses()))}
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
		normalized.Responses = append(normalized.Responses, &qualificationsexam.ExamResponse{
			QuestionId: question.GetId(),
			Question:   question,
			Response:   response.GetResponse(),
		})
	}
	if !partial && len(seen) != len(questions) {
		return nil, fmt.Errorf("missing exam responses")
	}

	return normalized, nil
}

func validateExamResponse(question *qualificationsexam.ExamQuestion, response *qualificationsexam.ExamResponseData, partial bool) error {
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
		maxLength := question.GetData().GetFreeText().GetMaxLength()
		if maxLength > 0 && int32(len([]rune(text.GetText()))) > maxLength {
			text.Text = string([]rune(text.GetText())[:maxLength])
		}
		if !partial && question.GetData().GetFreeText().GetMinLength() > 0 && int32(len([]rune(text.GetText()))) < question.GetData().GetFreeText().GetMinLength() {
			return fmt.Errorf("response for exam question %d is too short", question.GetId())
		}
	case question.GetData().GetSingleChoice() != nil:
		choice := response.GetSingleChoice()
		if choice == nil || !contains(question.GetData().GetSingleChoice().GetChoices(), choice.GetChoice()) {
			return fmt.Errorf("invalid choice for exam question %d", question.GetId())
		}
	case question.GetData().GetMultipleChoice() != nil:
		choices := response.GetMultipleChoice()
		if choices == nil || !uniqueSubset(choices.GetChoices(), question.GetData().GetMultipleChoice().GetChoices()) {
			return fmt.Errorf("invalid choices for exam question %d", question.GetId())
		}
		if limit := question.GetData().GetMultipleChoice().GetLimit(); limit > 0 && int32(len(choices.GetChoices())) > limit {
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
