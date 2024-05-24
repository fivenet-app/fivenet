package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getExam(ctx context.Context, qualificationId uint64, userInfo *userinfo.UserInfo, withQuestions bool, withQualification bool) (*qualifications.Exam, error) {
	columns := []jet.Projection{
		tExam.QualificationID,
		tExam.CreatedAt,
		tExam.DeletedAt,
		tExam.Settings,
	}

	if withQuestions {
		columns = append(columns, tExam.Questions)
	}

	stmt := tExam.
		SELECT(
			tExam.ID,
			columns...,
		).
		FROM(tExam).
		WHERE(jet.AND(
			tExam.QualificationID.EQ(jet.Uint64(qualificationId)),
		)).
		LIMIT(1)

	var dest qualifications.Exam
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	if dest.Id == 0 || dest.QualificationId == 0 {
		return nil, nil
	}

	if withQualification {
		quali, err := s.getQualificationShort(ctx, qualificationId, nil, userInfo)
		if err != nil {
			return nil, err
		}
		dest.Qualification = quali
	}

	return &dest, nil
}
