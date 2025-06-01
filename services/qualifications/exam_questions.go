package qualifications

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getExamQuestions(ctx context.Context, tx qrm.DB, qualificationId uint64, withAnswers bool) (*qualifications.ExamQuestions, error) {
	columns := []jet.Projection{
		tExamQuestions.QualificationID,
		tExamQuestions.CreatedAt,
		tExamQuestions.UpdatedAt,
		tExamQuestions.Title,
		tExamQuestions.Description,
		tExamQuestions.Data,
		tExamQuestions.Points,
	}

	if withAnswers {
		columns = append(columns, tExamQuestions.Answer)
	}

	stmt := tExamQuestions.
		SELECT(
			tExamQuestions.ID,
			columns...,
		).
		FROM(tExamQuestions).
		WHERE(jet.AND(
			tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId)),
		))

	var dest qualifications.ExamQuestions
	if err := stmt.QueryContext(ctx, tx, &dest.Questions); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) countExamQuestions(ctx context.Context, qualificationid uint64) (int32, error) {
	stmt := tExamQuestions.
		SELECT(
			jet.COUNT(jet.DISTINCT(tExamQuestions.ID)).AS("data_count.total"),
		).
		FROM(tExamQuestions).
		WHERE(
			tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationid)),
		)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		return 0, err
	}

	return int32(count.Total), nil
}

func (s *Server) handleExamQuestionsChanges(ctx context.Context, tx *sql.Tx, qualificationId uint64, questions *qualifications.ExamQuestions) ([]*file.File, error) {
	files := []*file.File{}

	tExamQuestions := table.FivenetQualificationsExamQuestions

	if len(questions.Questions) == 0 {
		stmt := tExamQuestions.
			DELETE().
			WHERE(tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId))).
			LIMIT(100)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		return nil, nil
	}

	current, err := s.getExamQuestions(ctx, tx, qualificationId, false)
	if err != nil {
		return nil, err
	}

	toCreate, toUpdate, toDelete := s.compareExamQuestions(current.Questions, questions.Questions)

	for _, question := range toCreate {
		if question.Data == nil {
			continue
		}

		switch data := question.Data.Data.(type) {
		case *qualifications.ExamQuestionData_Image:
			if data.Image.Image == nil {
				continue
			}

			files = append(files, data.Image.Image)
		}

		stmt := tExamQuestions.
			INSERT(
				tExamQuestions.QualificationID,
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
				tExamQuestions.Points,
			).
			VALUES(
				qualificationId,
				question.Title,
				question.Description,
				question.Data,
				question.Answer,
				question.Points,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	for _, question := range toUpdate {
		if question.Data != nil {
			switch data := question.Data.Data.(type) {
			case *qualifications.ExamQuestionData_Image:
				if data.Image.Image == nil {
					continue
				}

				files = append(files, data.Image.Image)
			}
		}

		stmt := tExamQuestions.
			UPDATE(
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
				tExamQuestions.Points,
			).
			SET(
				question.Title,
				question.Description,
				question.Data,
				question.Answer,
				question.Points,
			).
			WHERE(jet.AND(
				tExamQuestions.ID.EQ(jet.Uint64(question.Id)),
				tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	if len(toDelete) > 0 {
		questionIds := []jet.Expression{}
		for _, question := range toDelete {
			questionIds = append(questionIds, jet.Uint64(question.Id))
		}

		stmt := tExamQuestions.
			DELETE().
			WHERE(jet.AND(
				tExamQuestions.ID.IN(questionIds...),
				tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		// Don't include deleted questions files in the files list
	}

	return files, nil
}

func (s *Server) compareExamQuestions(current, in []*qualifications.ExamQuestion) (toCreate []*qualifications.ExamQuestion, toUpdate []*qualifications.ExamQuestion, toDelete []*qualifications.ExamQuestion) {
	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b *qualifications.ExamQuestion) int {
		return int(a.Id - b.Id)
	})

	if len(current) == 0 {
		toCreate = in
	} else {
		foundTracker := []int{}
		for _, cj := range current {
			idx := slices.IndexFunc(in, func(a *qualifications.ExamQuestion) bool {
				return cj.Id == a.Id
			})
			// No match in incoming questions, needs to be deleted
			if idx == -1 {
				toDelete = append(toDelete, cj)
				continue
			}

			foundTracker = append(foundTracker, idx)
			toUpdate = append(toUpdate, in[idx])
		}

		for i, eq := range in {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, eq)
			}
		}
	}

	return
}

type examResponses struct {
	ExamResponses *qualifications.ExamResponses `alias:"responses"`
	ExamGrading   *qualifications.ExamGrading   `alias:"grading"`
}

func (s *Server) getExamResponses(ctx context.Context, qualificationId uint64, userId int32) (*qualifications.ExamResponses, *qualifications.ExamGrading, error) {
	tExamResponses := tExamResponses.AS("examresponses")
	stmt := tExamResponses.
		SELECT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.Responses,
			tExamResponses.Grading,
		).
		FROM(tExamResponses).
		WHERE(jet.AND(
			tExamResponses.QualificationID.EQ(jet.Uint64(qualificationId)),
			tExamResponses.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	dest := &examResponses{
		ExamResponses: &qualifications.ExamResponses{},
		ExamGrading:   &qualifications.ExamGrading{},
	}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	dest.ExamResponses.QualificationId = qualificationId
	dest.ExamResponses.UserId = userId

	return dest.ExamResponses, dest.ExamGrading, nil
}
