package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	jet "github.com/go-jet/jet/v2/mysql"
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
		)).
		LIMIT(1)

	var dest qualifications.ExamQuestions
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
