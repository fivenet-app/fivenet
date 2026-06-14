package qualifications

import (
	"context"
	"errors"
	"time"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type examResponses struct {
	ExamResponses *qualificationsexam.ExamResponses `alias:"responses"`
	ExamGrading   *qualificationsexam.ExamGrading   `alias:"grading"`
}

func (s *Store) GetExamUser(
	ctx context.Context,
	qualificationId int64,
	userId int32,
) (*qualificationsexam.ExamUser, error) {
	tExamUser := tExamUser.AS("exam_user")
	stmt := tExamUser.
		SELECT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.CreatedAt,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
		).
		FROM(tExamUser).
		WHERE(mysql.AND(
			tExamUser.QualificationID.EQ(mysql.Int64(qualificationId)),
			tExamUser.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	var dest qualificationsexam.ExamUser
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if dest.GetQualificationId() == 0 || dest.GetUserId() == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Store) GetExamQuestions(
	ctx context.Context,
	q qrm.DB,
	qualificationId int64,
	withAnswers bool,
) (*qualificationsexam.ExamQuestions, error) {
	columns := mysql.ProjectionList{
		tExamQuestion.QualificationID,
		tExamQuestion.CreatedAt,
		tExamQuestion.UpdatedAt,
		tExamQuestion.Title,
		tExamQuestion.Description,
		tExamQuestion.Data,
		tExamQuestion.Points,
	}
	if withAnswers {
		columns = append(columns, tExamQuestion.Answer)
	}

	stmt := tExamQuestion.
		SELECT(tExamQuestion.ID, columns...).
		FROM(tExamQuestion).
		WHERE(tExamQuestion.QualificationID.EQ(mysql.Int64(qualificationId))).
		ORDER_BY(tExamQuestion.Order.ASC()).
		LIMIT(100)

	var dest qualificationsexam.ExamQuestions
	if err := stmt.QueryContext(ctx, q, &dest.Questions); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Store) CountExamQuestions(ctx context.Context, qualificationId int64) (int64, error) {
	tExamQuestion := tExamQuestion.AS("exam_question")
	stmt := tExamQuestion.
		SELECT(mysql.COUNT(mysql.DISTINCT(tExamQuestion.ID)).AS("data_count.total")).
		FROM(tExamQuestion).
		WHERE(tExamQuestion.QualificationID.EQ(mysql.Int64(qualificationId)))

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		return 0, err
	}

	return count.Total, nil
}

func (s *Store) GetExamResponses(
	ctx context.Context,
	qualificationId int64,
	userId int32,
) (*qualificationsexam.ExamResponses, *qualificationsexam.ExamGrading, error) {
	tExamResponses := tExamResponse.AS("examresponses")
	stmt := tExamResponses.
		SELECT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.Responses,
			tExamResponses.Grading,
		).
		FROM(tExamResponses).
		WHERE(mysql.AND(
			tExamResponses.QualificationID.EQ(mysql.Int64(qualificationId)),
			tExamResponses.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	dest := &examResponses{
		ExamResponses: &qualificationsexam.ExamResponses{},
		ExamGrading:   &qualificationsexam.ExamGrading{},
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

func (s *Store) DeleteExamUser(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
) error {
	stmt := tExamUser.
		DELETE().
		WHERE(mysql.AND(
			tExamUser.QualificationID.EQ(mysql.Int64(qualificationId)),
			tExamUser.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) CreateExamUser(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	endsAt time.Time,
) error {
	stmt := tExamUser.
		INSERT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
		).
		VALUES(
			qualificationId,
			userId,
			mysql.CURRENT_TIMESTAMP(),
			mysql.TimestampT(endsAt),
			mysql.NULL,
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) UpsertExamResponses(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	responses *qualificationsexam.ExamResponses,
) error {
	stmt := tExamResponse.
		INSERT(
			tExamResponse.QualificationID,
			tExamResponse.UserID,
			tExamResponse.Responses,
			tExamResponse.Grading,
		).
		VALUES(
			qualificationId,
			userId,
			responses,
			mysql.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tExamResponse.Responses.SET(mysql.RawString("VALUES(`responses`)")),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) UpsertExamUserEndedAt(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	endedAt time.Time,
) error {
	stmt := tExamUser.
		INSERT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.EndedAt,
		).
		VALUES(
			qualificationId,
			userId,
			mysql.DateTimeT(endedAt),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tExamUser.EndedAt.SET(mysql.TimestampT(endedAt)),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
