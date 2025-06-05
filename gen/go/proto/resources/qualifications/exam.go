package qualifications

import (
	"slices"

	"google.golang.org/protobuf/proto"
)

func (e *ExamQuestions) Grade(mode AutoGradeMode, questions *ExamResponses) (float32, *ExamGrading) {
	var totalPoints float32
	var earnedPoints float32
	grading := &ExamGrading{Responses: []*ExamGradingResponse{}}

	for _, question := range e.Questions {
		if question.Data.GetFreeText() != nil || question.Data.GetImage() != nil || question.Data.GetSeparator() != nil {
			// Skip free text questions
			continue
		}

		for _, response := range questions.Responses {
			if response.QuestionId == question.Id {
				if question.Answer != nil && response.Response != nil {
					switch {
					case response.Response.GetYesno() != nil && question.Answer.GetYesno() != nil:
						if response.Response.GetYesno().Value == question.Answer.GetYesno().Value {
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     float32(question.GetPoints()),
								Checked:    proto.Bool(true),
							})
							earnedPoints += float32(question.GetPoints())
						} else {
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     0,
								Checked:    proto.Bool(true),
							})
						}
					case response.Response.GetSingleChoice() != nil && question.Answer.GetSingleChoice() != nil:
						if response.Response.GetSingleChoice().Choice == question.Answer.GetSingleChoice().Choice {
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     float32(question.GetPoints()),
								Checked:    proto.Bool(true),
							})
							earnedPoints += float32(question.GetPoints())
						} else {
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     0,
								Checked:    proto.Bool(true),
							})
						}
					case response.Response.GetMultipleChoice() != nil && question.Answer.GetMultipleChoice() != nil:
						correctChoices := 0
						for _, choice := range response.Response.GetMultipleChoice().Choices {
							if slices.Contains(question.Answer.GetMultipleChoice().Choices, choice) {
								correctChoices++
							}
						}
						if mode == AutoGradeMode_AUTO_GRADE_MODE_PARTIAL_CREDIT {
							points := float32(question.GetPoints()) * (float32(correctChoices) / float32(len(question.Answer.GetMultipleChoice().Choices)))
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     points,
								Checked:    proto.Bool(true),
							})
							earnedPoints += points
						} else if correctChoices == len(question.Answer.GetMultipleChoice().Choices) {
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     float32(question.GetPoints()),
								Checked:    proto.Bool(true),
							})
							earnedPoints += float32(question.GetPoints())
						} else {
							grading.Responses = append(grading.Responses, &ExamGradingResponse{
								QuestionId: question.Id,
								Points:     0,
								Checked:    proto.Bool(true),
							})
						}
					}
				}
			}
		}
		totalPoints += float32(question.GetPoints())
	}

	if totalPoints == 0 {
		return 0, grading
	}

	return earnedPoints, grading
}
