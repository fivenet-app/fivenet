package qualificationsexam

import (
	"slices"
)

func (e *ExamQuestions) Grade(
	mode AutoGradeMode,
	questions *ExamResponses,
) (float32, *ExamGrading) {
	var totalPoints float32
	var earnedPoints float32
	grading := &ExamGrading{Responses: []*ExamGradingResponse{}}

	for _, question := range e.GetQuestions() {
		if question.GetData().GetFreeText() != nil || question.GetData().GetImage() != nil ||
			question.GetData().GetSeparator() != nil {
			// Skip free text questions
			continue
		}

		response := findResponse(questions.GetResponses(), question.GetId())
		if response == nil {
			continue
		}

		if question.GetAnswer() == nil || response.GetResponse() == nil {
			continue
		}

		switch {
		case response.GetResponse().GetYesno() != nil && question.GetAnswer().GetYesno() != nil:
			if response.GetResponse().
				GetYesno().
				GetValue() ==
				question.GetAnswer().
					GetYesno().
					GetValue() {
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     float32(question.GetPoints()),
					Checked:    new(true),
				})
				earnedPoints += float32(question.GetPoints())
			} else {
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     0,
					Checked:    new(true),
				})
			}
		case response.GetResponse().GetSingleChoice() != nil && question.GetAnswer().GetSingleChoice() != nil:
			if response.GetResponse().
				GetSingleChoice().
				GetChoice() ==
				question.GetAnswer().
					GetSingleChoice().
					GetChoice() {
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     float32(question.GetPoints()),
					Checked:    new(true),
				})
				earnedPoints += float32(question.GetPoints())
			} else {
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     0,
					Checked:    new(true),
				})
			}
		case response.GetResponse().GetMultipleChoice() != nil && question.GetAnswer().GetMultipleChoice() != nil:
			answerChoices := question.GetAnswer().GetMultipleChoice().GetChoices()
			if len(answerChoices) == 0 {
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(), Points: 0, Checked: new(true),
				})
				break
			}
			correctChoices := 0
			for _, choice := range response.GetResponse().GetMultipleChoice().GetChoices() {
				if slices.Contains(
					answerChoices,
					choice,
				) {
					correctChoices++
				}
			}

			switch {
			case mode == AutoGradeMode_AUTO_GRADE_MODE_PARTIAL_CREDIT:
				points := float32(
					question.GetPoints(),
				) * (float32(correctChoices) / float32(len(answerChoices)))
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     points,
					Checked:    new(true),
				})
				earnedPoints += points

			case correctChoices == len(answerChoices) && len(response.GetResponse().GetMultipleChoice().GetChoices()) == len(answerChoices):
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     float32(question.GetPoints()),
					Checked:    new(true),
				})
				earnedPoints += float32(question.GetPoints())

			default:
				grading.Responses = append(grading.Responses, &ExamGradingResponse{
					QuestionId: question.GetId(),
					Points:     0,
					Checked:    new(true),
				})
			}
		}

		totalPoints += float32(question.GetPoints())
	}

	if totalPoints == 0 {
		return 0, grading
	}

	return earnedPoints, grading
}

func findResponse(responses []*ExamResponse, questionID int64) *ExamResponse {
	for _, response := range responses {
		if response.GetQuestionId() == questionID {
			return response
		}
	}
	return nil
}
