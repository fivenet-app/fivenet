package qualifications

import (
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"google.golang.org/protobuf/proto"
)

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

func examForCandidate(exam *qualificationsexam.ExamQuestions) *qualificationsexam.ExamQuestions {
	if exam == nil {
		return nil
	}
	examCopy := proto.Clone(exam).(*qualificationsexam.ExamQuestions)
	for _, question := range examCopy.GetQuestions() {
		question.ClearAnswer()
	}
	return examCopy
}

func publicExamUser(examUser *qualificationsexam.ExamUser) *qualificationsexam.ExamUser {
	if examUser == nil {
		return nil
	}
	examUserCopy := proto.Clone(examUser).(*qualificationsexam.ExamUser)
	examUserCopy.ClearSnapshot()
	examUserCopy.SetAttemptId("")
	return examUserCopy
}

func publicExamResponses(
	exam *qualificationsexam.ExamQuestions,
	responses *qualificationsexam.ExamResponses,
) *qualificationsexam.ExamResponses {
	if responses == nil {
		return nil
	}
	copyExamResponse := sanitizeExamResponses(exam, responses)
	copyExamResponse.SetAttemptId("")
	return copyExamResponse
}
