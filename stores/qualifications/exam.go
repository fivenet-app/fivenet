package qualificationsstore

import (
	"context"
	"errors"
	"time"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

const ExamSubmissionGracePeriod = 30 * time.Second

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
			tExamUser.AttemptID,
			tExamUser.CreatedAt,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
			tExamUser.Snapshot,
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
		SELECT(
			tExamQuestion.ID,
			columns...,
		).
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
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tExamQuestion.ID)).AS("data_count.total"),
		).
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
	q qrm.DB,
	attemptId string,
) (*qualificationsexam.ExamResponses, *qualificationsexam.ExamGrading, error) {
	tExamResponses := tExamResponses.AS("examresponses")
	stmt := tExamResponses.
		SELECT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.AttemptID,
			tExamResponses.Responses,
			tExamResponses.Grading,
		).
		FROM(tExamResponses).
		WHERE(mysql.AND(
			tExamResponses.AttemptID.EQ(mysql.String(attemptId)),
		)).
		LIMIT(1)

	dest := &examResponses{
		ExamResponses: &qualificationsexam.ExamResponses{},
		ExamGrading:   &qualificationsexam.ExamGrading{},
	}
	if err := stmt.QueryContext(ctx, q, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	dest.ExamResponses.SetAttemptId(attemptId)

	return dest.ExamResponses, dest.ExamGrading, nil
}

func (s *Store) UpsertExamResponses(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	attemptId string,
	responses *qualificationsexam.ExamResponses,
) error {
	tExamResponses := table.FivenetQualificationsExamResponses
	stmt := tExamResponses.
		INSERT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.AttemptID,
			tExamResponses.Responses,
			tExamResponses.Grading,
		).
		VALUES(
			qualificationId,
			userId,
			attemptId,
			responses,
			mysql.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tExamResponses.Responses.SET(mysql.RawString("VALUES(`responses`)")),
			tExamResponses.AttemptID.SET(mysql.RawString("VALUES(`attempt_id`)")),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteExamUser(
	ctx context.Context,
	tx qrm.DB,
	attemptId string,
) error {
	tExamUser := table.FivenetQualificationsExamUsers
	stmt := tExamUser.
		DELETE().
		WHERE(mysql.AND(
			tExamUser.AttemptID.EQ(mysql.String(attemptId)),
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
	snapshot *qualificationsexam.ExamSnapshot,
) (string, error) {
	tExamUser := table.FivenetQualificationsExamUsers
	attemptId := uuid.NewString()
	stmt := tExamUser.
		INSERT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.AttemptID,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
			tExamUser.Snapshot,
		).
		VALUES(
			qualificationId,
			userId,
			attemptId,
			mysql.CURRENT_TIMESTAMP(),
			mysql.TimestampT(endsAt),
			mysql.NULL,
			snapshot,
		)

	_, err := stmt.ExecContext(ctx, tx)
	return attemptId, err
}

// ClaimActiveExamUser serializes a response write with expiry. When complete
// is true it also marks the attempt ended, but only while it is still active.
func (s *Store) ClaimActiveExamUser(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	attemptId string,
	complete bool,
	gracePeriod time.Duration,
) (bool, error) {
	tExamUser := table.FivenetQualificationsExamUsers
	conditions := mysql.AND(
		tExamUser.QualificationID.EQ(mysql.Int64(qualificationId)),
		tExamUser.UserID.EQ(mysql.Int32(userId)),
		tExamUser.AttemptID.EQ(mysql.String(attemptId)),
		tExamUser.EndedAt.IS_NULL(),
		tExamUser.EndsAt.IS_NOT_NULL(),
		tExamUser.EndsAt.GT(mysql.TimestampT(time.Now().Add(-gracePeriod))),
	)
	// Lock the attempt before checking or changing its state. In particular,
	// partial submissions must not rely on RowsAffected from a no-op update:
	// MySQL reports zero affected rows when the value remains unchanged.
	var examUser qualificationsexam.ExamUser
	selectStmt := tExamUser.
		SELECT(tExamUser.EndedAt, tExamUser.EndsAt).
		FROM(tExamUser).
		WHERE(conditions).
		LIMIT(1).
		FOR(mysql.UPDATE())
	if err := selectStmt.QueryContext(ctx, tx, &examUser); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if !complete {
		return true, nil
	}

	updateStmt := tExamUser.UPDATE(tExamUser.EndedAt).
		SET(mysql.CURRENT_TIMESTAMP()).
		WHERE(conditions).
		LIMIT(1)
	result, err := updateStmt.ExecContext(ctx, tx)
	if err != nil {
		return false, err
	}
	updated, err := result.RowsAffected()
	return updated > 0, err
}

func (s *Store) ListExpiredExamUsers(
	ctx context.Context,
	limit int64,
) ([]*qualificationsexam.ExamUser, error) {
	tExamUser := table.FivenetQualificationsExamUsers
	stmt := tExamUser.
		SELECT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.AttemptID,
			tExamUser.CreatedAt,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
			tExamUser.Snapshot,
		).
		FROM(tExamUser).
		WHERE(mysql.AND(
			tExamUser.EndedAt.IS_NULL(),
			tExamUser.EndsAt.IS_NOT_NULL(),
			tExamUser.EndsAt.LT(mysql.TimestampT(time.Now().Add(-ExamSubmissionGracePeriod))),
		)).
		ORDER_BY(tExamUser.EndsAt.ASC()).
		LIMIT(limit)

	var attempts []*qualificationsexam.ExamUser
	if err := stmt.QueryContext(
		ctx,
		s.db,
		&attempts,
	); err != nil &&
		!errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	return attempts, nil
}

// ListExamUsersPastRetention returns only attempts whose requests completed.
// In particular, manually graded exams remain available while grading is pending.
func (s *Store) ListExamUsersPastRetention(
	ctx context.Context,
	olderThan time.Time,
	limit int64,
) ([]*qualificationsexam.ExamUser, error) {
	tExamUser := table.FivenetQualificationsExamUsers
	stmt := tExamUser.
		SELECT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.AttemptID,
			tExamUser.CreatedAt,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
		).
		FROM(tExamUser.INNER_JOIN(tQualiReq,
			tQualiReq.QualificationID.EQ(tExamUser.QualificationID).
				AND(tQualiReq.UserID.EQ(tExamUser.UserID)),
		)).
		WHERE(mysql.AND(
			tExamUser.EndedAt.IS_NOT_NULL(),
			tExamUser.EndedAt.LT_EQ(mysql.TimestampT(olderThan)),
			tQualiReq.Status.EQ(mysql.Int32(int32(
				resqualifications.RequestStatus_REQUEST_STATUS_COMPLETED,
			))),
		)).
		ORDER_BY(tExamUser.EndedAt.ASC()).
		LIMIT(limit)

	var attempts []*qualificationsexam.ExamUser
	if err := stmt.QueryContext(
		ctx,
		s.db,
		&attempts,
	); err != nil &&
		!errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	return attempts, nil
}

func (s *Store) ExpireExamUser(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	attemptId string,
) (bool, error) {
	tExamUser := table.FivenetQualificationsExamUsers
	stmt := tExamUser.
		UPDATE(tExamUser.EndedAt).
		SET(tExamUser.EndsAt).
		WHERE(mysql.AND(
			tExamUser.QualificationID.EQ(mysql.Int64(qualificationId)),
			tExamUser.UserID.EQ(mysql.Int32(userId)),
			tExamUser.AttemptID.EQ(mysql.String(attemptId)),
			tExamUser.EndedAt.IS_NULL(),
			tExamUser.EndsAt.IS_NOT_NULL(),
			tExamUser.EndsAt.LT(mysql.TimestampT(time.Now().Add(-ExamSubmissionGracePeriod))),
		)).
		LIMIT(1)
	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return false, err
	}
	updated, err := result.RowsAffected()
	return updated > 0, err
}

func (s *Store) DeleteExamResponses(
	ctx context.Context,
	tx qrm.DB,
	attemptId string,
) error {
	tExamResponses := table.FivenetQualificationsExamResponses
	stmt := tExamResponses.
		DELETE().
		WHERE(mysql.AND(
			tExamResponses.AttemptID.EQ(mysql.String(attemptId)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
