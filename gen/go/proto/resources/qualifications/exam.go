package qualifications

import (
	"slices"

	"google.golang.org/protobuf/proto"
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

		for _, response := range questions.GetResponses() {
			if response.GetQuestionId() != question.GetId() {
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
						Checked:    proto.Bool(true),
					})
					earnedPoints += float32(question.GetPoints())
				} else {
					grading.Responses = append(grading.Responses, &ExamGradingResponse{
						QuestionId: question.GetId(),
						Points:     0,
						Checked:    proto.Bool(true),
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
						Checked:    proto.Bool(true),
					})
					earnedPoints += float32(question.GetPoints())
				} else {
					grading.Responses = append(grading.Responses, &ExamGradingResponse{
						QuestionId: question.GetId(),
						Points:     0,
						Checked:    proto.Bool(true),
					})
				}
			case response.GetResponse().GetMultipleChoice() != nil && question.GetAnswer().GetMultipleChoice() != nil:
				correctChoices := 0
				for _, choice := range response.GetResponse().GetMultipleChoice().GetChoices() {
					if slices.Contains(
						question.GetAnswer().GetMultipleChoice().GetChoices(),
						choice,
					) {
						correctChoices++
					}
				}

				switch {
				case mode == AutoGradeMode_AUTO_GRADE_MODE_PARTIAL_CREDIT:
					points := float32(
						question.GetPoints(),
					) * (float32(correctChoices) / float32(len(question.GetAnswer().GetMultipleChoice().GetChoices())))
					grading.Responses = append(grading.Responses, &ExamGradingResponse{
						QuestionId: question.GetId(),
						Points:     points,
						Checked:    proto.Bool(true),
					})
					earnedPoints += points

				case correctChoices == len(question.GetAnswer().GetMultipleChoice().GetChoices()):
					grading.Responses = append(grading.Responses, &ExamGradingResponse{
						QuestionId: question.GetId(),
						Points:     float32(question.GetPoints()),
						Checked:    proto.Bool(true),
					})
					earnedPoints += float32(question.GetPoints())

				default:
					grading.Responses = append(grading.Responses, &ExamGradingResponse{
						QuestionId: question.GetId(),
						Points:     0,
						Checked:    proto.Bool(true),
					})
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
