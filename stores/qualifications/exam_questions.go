package qualificationsstore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Store) HandleExamQuestionsChanges(
	ctx context.Context,
	tx *sql.Tx,
	qualificationId int64,
	questions *qualificationsexam.ExamQuestions,
) ([]*file.File, error) {
	files := []*file.File{}

	tExamQuestion := table.FivenetQualificationsExamQuestions
	if len(questions.GetQuestions()) == 0 {
		stmt := tExamQuestion.
			DELETE().
			WHERE(tExamQuestion.QualificationID.EQ(mysql.Int64(qualificationId))).
			LIMIT(100)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		return nil, nil
	}

	current, err := s.GetExamQuestions(ctx, tx, qualificationId, false)
	if err != nil {
		return nil, err
	}

	toCreate, toUpdate, toDelete, err := compareExamQuestions(
		current.GetQuestions(),
		questions.GetQuestions(),
	)
	if err != nil {
		return nil, err
	}

	for _, question := range toCreate {
		if question.GetData() == nil {
			continue
		}

		switch data := question.GetData().GetData().(type) {
		case *qualificationsexam.ExamQuestionData_Image:
			if data.Image.GetImage() == nil {
				return nil, fmt.Errorf("image question requires an image")
			}
			files = append(files, data.Image.GetImage())
		}

		stmt := tExamQuestion.
			INSERT(
				tExamQuestion.QualificationID,
				tExamQuestion.Title,
				tExamQuestion.Description,
				tExamQuestion.Data,
				tExamQuestion.Answer,
				tExamQuestion.Points,
				tExamQuestion.Order,
			).
			VALUES(
				qualificationId,
				question.GetTitle(),
				question.Description,
				question.GetData(),
				question.GetAnswer(),
				question.GetPoints(),
				question.GetOrder(),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	for _, question := range toUpdate {
		if question.GetData() != nil {
			switch data := question.GetData().GetData().(type) {
			case *qualificationsexam.ExamQuestionData_Image:
				if data.Image.GetImage() == nil {
					return nil, fmt.Errorf("image question requires an image")
				}
				files = append(files, data.Image.GetImage())
			}
		}

		stmt := tExamQuestion.
			UPDATE(
				tExamQuestion.Title,
				tExamQuestion.Description,
				tExamQuestion.Data,
				tExamQuestion.Answer,
				tExamQuestion.Points,
				tExamQuestion.Order,
			).
			SET(
				question.GetTitle(),
				question.Description,
				question.GetData(),
				question.GetAnswer(),
				question.Points,
				question.GetOrder(),
			).
			WHERE(mysql.AND(
				tExamQuestion.ID.EQ(mysql.Int64(question.GetId())),
				tExamQuestion.QualificationID.EQ(mysql.Int64(qualificationId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	if len(toDelete) > 0 {
		questionIds := []mysql.Expression{}
		for _, question := range toDelete {
			questionIds = append(questionIds, mysql.Int64(question.GetId()))
		}

		stmt := tExamQuestion.
			DELETE().
			WHERE(mysql.AND(
				tExamQuestion.ID.IN(questionIds...),
				tExamQuestion.QualificationID.EQ(mysql.Int64(qualificationId)),
			)).
			LIMIT(int64(len(questionIds)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	return files, nil
}

func compareExamQuestions(
	current, in []*qualificationsexam.ExamQuestion,
) ([]*qualificationsexam.ExamQuestion, []*qualificationsexam.ExamQuestion, []*qualificationsexam.ExamQuestion, error) {
	toCreate := []*qualificationsexam.ExamQuestion{}
	toUpdate := []*qualificationsexam.ExamQuestion{}
	toDelete := []*qualificationsexam.ExamQuestion{}
	persisted := make(map[int64]*qualificationsexam.ExamQuestion, len(current))
	for _, question := range current {
		if question == nil || question.GetId() <= 0 {
			return nil, nil, nil, fmt.Errorf("invalid persisted exam question")
		}
		persisted[question.GetId()] = question
	}

	incoming := make(map[int64]*qualificationsexam.ExamQuestion, len(in))
	for _, question := range in {
		if question == nil {
			return nil, nil, nil, fmt.Errorf("invalid exam question")
		}
		if question.GetId() <= 0 {
			toCreate = append(toCreate, question)
			continue
		}
		if _, duplicate := incoming[question.GetId()]; duplicate {
			return nil, nil, nil, fmt.Errorf("duplicate exam question id %d", question.GetId())
		}
		if _, exists := persisted[question.GetId()]; !exists {
			return nil, nil, nil, fmt.Errorf("unknown exam question id %d", question.GetId())
		}
		incoming[question.GetId()] = question
		toUpdate = append(toUpdate, question)
	}

	for _, question := range current {
		if _, exists := incoming[question.GetId()]; !exists {
			toDelete = append(toDelete, question)
		}
	}

	return toCreate, toUpdate, toDelete, nil
}
