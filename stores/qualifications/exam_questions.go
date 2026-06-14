package qualifications

import (
	"context"
	"database/sql"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Store) HandleExamQuestionsChanges(
	ctx context.Context,
	tx *sql.Tx,
	qualificationId int64,
	questions *qualificationsexam.ExamQuestions,
) ([]*file.File, error) {
	files := []*file.File{}

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

	toCreate, toUpdate, toDelete := compareExamQuestions(
		current.GetQuestions(),
		questions.GetQuestions(),
	)

	for _, question := range toCreate {
		if question.GetData() == nil {
			continue
		}

		switch data := question.GetData().GetData().(type) {
		case *qualificationsexam.ExamQuestionData_Image:
			if data.Image.GetImage() == nil {
				continue
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
					continue
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
) ([]*qualificationsexam.ExamQuestion, []*qualificationsexam.ExamQuestion, []*qualificationsexam.ExamQuestion) {
	toCreate := []*qualificationsexam.ExamQuestion{}
	toUpdate := []*qualificationsexam.ExamQuestion{}
	toDelete := []*qualificationsexam.ExamQuestion{}
	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b *qualificationsexam.ExamQuestion) int {
		return int(a.GetId() - b.GetId())
	})

	foundTracker := []int{}
	for _, cj := range current {
		idx := slices.IndexFunc(in, func(a *qualificationsexam.ExamQuestion) bool {
			return cj.GetId() == a.GetId()
		})
		if idx == -1 {
			toDelete = append(toDelete, cj)
			continue
		}

		foundTracker = append(foundTracker, idx)
		toUpdate = append(toUpdate, in[idx])
	}

	for i, eq := range in {
		if idx := slices.Index(foundTracker, i); idx == -1 {
			toCreate = append(toCreate, eq)
		}
	}

	return toCreate, toUpdate, toDelete
}
