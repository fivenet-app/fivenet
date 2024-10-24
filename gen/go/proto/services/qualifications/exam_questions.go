package qualifications

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getExamQuestions(ctx context.Context, qualificationId uint64, withAnswers bool) (*qualifications.ExamQuestions, error) {
	columns := []jet.Projection{
		tExamQuestions.QualificationID,
		tExamQuestions.CreatedAt,
		tExamQuestions.UpdatedAt,
		tExamQuestions.Title,
		tExamQuestions.Description,
		tExamQuestions.Data,
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
	if err := stmt.QueryContext(ctx, s.db, &dest.Questions); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) countExamQuestions(ctx context.Context, qualificationid uint64) (int32, error) {
	stmt := tExamQuestions.
		SELECT(
			jet.COUNT(jet.DISTINCT(tExamQuestions.ID)).AS("datacount.totalcount"),
		).
		FROM(tExamQuestions).
		WHERE(
			tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationid)),
		)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		return 0, err
	}

	return int32(count.TotalCount), nil
}

func (s *Server) handleExamQuestionsChanges(ctx context.Context, tx qrm.DB, qualificiationId uint64, questions *qualifications.ExamQuestions) error {
	tExamQuestions := table.FivenetQualificationsExamQuestions
	if len(questions.Questions) == 0 {
		stmt := tExamQuestions.
			DELETE().
			WHERE(tExamQuestions.QualificationID.EQ(jet.Uint64(qualificiationId))).
			LIMIT(100)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}

		return nil
	}

	toCreate := []*qualifications.ExamQuestion{}
	toUpdate := []*qualifications.ExamQuestion{}

	for _, question := range questions.Questions {
		question.QualificationId = qualificiationId
		if question.Id > 0 {
			toUpdate = append(toUpdate, question)
		} else {
			toCreate = append(toCreate, question)
		}
	}

	if len(toCreate) > 0 {
		stmt := tExamQuestions.
			INSERT(
				tExamQuestions.QualificationID,
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
			)

		for _, question := range toCreate {
			stmt = stmt.VALUES(
				question.QualificationId,
				question.Title,
				question.Description,
				question.Data,
				question.Answer,
			)
		}

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	if len(toUpdate) > 0 {
		for _, question := range toUpdate {
			stmt := tExamQuestions.
				UPDATE(
					tExamQuestions.Title,
					tExamQuestions.Description,
					tExamQuestions.Data,
					tExamQuestions.Answer,
				).
				SET(
					question.Title,
					question.Description,
					question.Data,
					question.Answer,
				).
				WHERE(jet.AND(
					tExamQuestions.ID.EQ(jet.Uint64(question.Id)),
					tExamQuestions.QualificationID.EQ(jet.Uint64(question.QualificationId)),
				))

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

type examResponses struct {
	ExamResponses *qualifications.ExamResponses `alias:"responses"`
}

func (s *Server) getExamResponses(ctx context.Context, qualificationId uint64, userId int32) (*qualifications.ExamResponses, error) {
	tExamResponses := tExamResponses.AS("examresponses")
	stmt := tExamResponses.
		SELECT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.Responses,
		).
		FROM(tExamResponses).
		WHERE(jet.AND(
			tExamResponses.QualificationID.EQ(jet.Uint64(qualificationId)),
			tExamResponses.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	dest := examResponses{
		ExamResponses: &qualifications.ExamResponses{},
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	dest.ExamResponses.QualificationId = qualificationId
	dest.ExamResponses.UserId = userId

	return dest.ExamResponses, nil
}
