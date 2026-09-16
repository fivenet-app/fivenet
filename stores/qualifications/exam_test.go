package qualificationsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreClaimActiveExamUserAllowsPartialSaveDuringGracePeriod(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)

	expectedQuery := regexp.QuoteMeta(`SELECT`) +
		`(?s).*` + regexp.QuoteMeta(`FROM fivenet_qualifications_exam_users`) +
		`(?s).*` + regexp.QuoteMeta(`qualification_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`attempt_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`FOR UPDATE`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int32(7), "attempt-1", sqlmock.AnyArg(), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"ended_at", "ends_at"}).AddRow(nil, time.Now().Add(-time.Second)))

	active, err := store.ClaimActiveExamUser(
		t.Context(),
		tx,
		42,
		7,
		"attempt-1",
		false,
		30*time.Second,
	)
	require.NoError(t, err)
	assert.True(t, active)
	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreClaimActiveExamUserRejectsExpiredCancellation(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)

	expectedQuery := regexp.QuoteMeta(`SELECT`) +
		`(?s).*` + regexp.QuoteMeta(`FROM fivenet_qualifications_exam_users`) +
		`(?s).*` + regexp.QuoteMeta(`attempt_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`FOR UPDATE`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int32(7), "attempt-1", sqlmock.AnyArg(), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"ended_at", "ends_at"}))

	active, err := store.ClaimActiveExamUser(t.Context(), tx, 42, 7, "attempt-1", false, 0)
	require.NoError(t, err)
	assert.False(t, active)
	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetExamUser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_qualifications_exam_users AS exam_user`) +
		`(?s).*` + regexp.QuoteMeta(`exam_user.qualification_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`exam_user.user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	now := time.Unix(0, 0).UTC()
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int32(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"exam_user.qualification_id",
			"exam_user.user_id",
			"exam_user.created_at",
			"exam_user.started_at",
			"exam_user.ends_at",
			"exam_user.ended_at",
		}).AddRow(int64(42), int32(7), now, now, now, nil))

	examUser, err := store.GetExamUser(t.Context(), db, 42, 7)
	require.NoError(t, err)
	require.NotNil(t, examUser)
	assert.Equal(t, int64(42), examUser.GetQualificationId())
	assert.Equal(t, int32(7), examUser.GetUserId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCountExamQuestions(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_qualifications_exam_questions AS exam_question`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`exam_question.qualification_id = ?`,
	)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(3)))

	total, err := store.CountExamQuestions(t.Context(), 42)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreExpireExamUserClaimsOnlyActiveExpiredAttempt(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectBegin()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)

	expectedQuery := regexp.QuoteMeta(`UPDATE fivenet_qualifications_exam_users`) +
		`(?s).*` + regexp.QuoteMeta(`SET ended_at = fivenet_qualifications_exam_users.ends_at`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_qualifications_exam_users.qualification_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_qualifications_exam_users.user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_qualifications_exam_users.attempt_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_qualifications_exam_users.ended_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_qualifications_exam_users.ends_at < TIMESTAMP(?)`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(42), int32(7), "attempt-1", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	expired, err := store.ExpireExamUser(t.Context(), tx, 42, 7, "attempt-1")
	require.NoError(t, err)
	assert.True(t, expired)
	mock.ExpectCommit()
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListExamUsersPastRetentionExcludesPendingGrading(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	expectedQuery := regexp.QuoteMeta(`FROM fivenet_qualifications_exam_users`) +
		`(?s).*` + regexp.QuoteMeta(`INNER JOIN fivenet_qualifications_requests AS qualification_request`) +
		`(?s).*` + regexp.QuoteMeta(`qualification_request.status = ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(sqlmock.AnyArg(), int32(resqualifications.RequestStatus_REQUEST_STATUS_COMPLETED), int64(1000)).
		WillReturnRows(sqlmock.NewRows([]string{
			"qualification_id",
			"user_id",
			"attempt_id",
			"created_at",
			"started_at",
			"ends_at",
			"ended_at",
		}))

	attempts, err := store.ListExamUsersPastRetention(t.Context(), time.Now(), 1000)
	require.NoError(t, err)
	assert.Empty(t, attempts)
	require.NoError(t, mock.ExpectationsWereMet())
}
