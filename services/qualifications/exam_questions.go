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

func (s *Server) getExamQuestions(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	withAnswers bool,
) (*qualifications.ExamQuestions, error) {
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
			tExamQuestions.QualificationID.EQ(jet.Int64(qualificationId)),
		)).
		ORDER_BY(tExamQuestions.Order.ASC()).
		LIMIT(100)

	var dest qualifications.ExamQuestions
	if err := stmt.QueryContext(ctx, tx, &dest.Questions); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) countExamQuestions(ctx context.Context, qualificationid int64) (int64, error) {
	stmt := tExamQuestions.
		SELECT(
			jet.COUNT(jet.DISTINCT(tExamQuestions.ID)).AS("data_count.total"),
		).
		FROM(tExamQuestions).
		WHERE(
			tExamQuestions.QualificationID.EQ(jet.Int64(qualificationid)),
		)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		return 0, err
	}

	return count.Total, nil
}

func (s *Server) handleExamQuestionsChanges(
	ctx context.Context,
	tx *sql.Tx,
	qualificationId int64,
	questions *qualifications.ExamQuestions,
) ([]*file.File, error) {
	files := []*file.File{}

	tExamQuestions := table.FivenetQualificationsExamQuestions

	if len(questions.GetQuestions()) == 0 {
		stmt := tExamQuestions.
			DELETE().
			WHERE(tExamQuestions.QualificationID.EQ(jet.Int64(qualificationId))).
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

	toCreate, toUpdate, toDelete := s.compareExamQuestions(
		current.GetQuestions(),
		questions.GetQuestions(),
	)

	for _, question := range toCreate {
		if question.GetData() == nil {
			continue
		}

		switch data := question.GetData().GetData().(type) {
		case *qualifications.ExamQuestionData_Image:
			if data.Image.GetImage() == nil {
				continue
			}

			files = append(files, data.Image.GetImage())
		}

		stmt := tExamQuestions.
			INSERT(
				tExamQuestions.QualificationID,
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
				tExamQuestions.Points,
				tExamQuestions.Order,
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
			case *qualifications.ExamQuestionData_Image:
				if data.Image.GetImage() == nil {
					continue
				}

				files = append(files, data.Image.GetImage())
			}
		}

		stmt := tExamQuestions.
			UPDATE(
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
				tExamQuestions.Points,
				tExamQuestions.Order,
			).
			SET(
				question.GetTitle(),
				question.Description,
				question.GetData(),
				question.GetAnswer(),
				question.Points,
				question.GetOrder(),
			).
			WHERE(jet.AND(
				tExamQuestions.ID.EQ(jet.Int64(question.GetId())),
				tExamQuestions.QualificationID.EQ(jet.Int64(qualificationId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	if len(toDelete) > 0 {
		questionIds := []jet.Expression{}
		for _, question := range toDelete {
			questionIds = append(questionIds, jet.Int64(question.GetId()))
		}

		stmt := tExamQuestions.
			DELETE().
			WHERE(jet.AND(
				tExamQuestions.ID.IN(questionIds...),
				tExamQuestions.QualificationID.EQ(jet.Int64(qualificationId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		// Don't include deleted questions files in the files list
	}

	return files, nil
}

func (s *Server) compareExamQuestions(
	current, in []*qualifications.ExamQuestion,
) ([]*qualifications.ExamQuestion, []*qualifications.ExamQuestion, []*qualifications.ExamQuestion) {
	toCreate := []*qualifications.ExamQuestion{}
	toUpdate := []*qualifications.ExamQuestion{}
	toDelete := []*qualifications.ExamQuestion{}
	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b *qualifications.ExamQuestion) int {
		return int(a.GetId() - b.GetId())
	})

	foundTracker := []int{}
	for _, cj := range current {
		idx := slices.IndexFunc(in, func(a *qualifications.ExamQuestion) bool {
			return cj.GetId() == a.GetId()
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

	return toCreate, toUpdate, toDelete
}

type examResponses struct {
	ExamResponses *qualifications.ExamResponses `alias:"responses"`
	ExamGrading   *qualifications.ExamGrading   `alias:"grading"`
}

func (s *Server) getExamResponses(
	ctx context.Context,
	qualificationId int64,
	userId int32,
) (*qualifications.ExamResponses, *qualifications.ExamGrading, error) {
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
			tExamResponses.QualificationID.EQ(jet.Int64(qualificationId)),
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
