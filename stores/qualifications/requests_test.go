package qualificationsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/require"
)

func TestStoreUpdateRequestStatus(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_qualifications_requests`) +
		`(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)

	mock.ExpectExec(expectedQuery).
		WithArgs(int64(42), int32(7), int32(resqualifications.RequestStatus_REQUEST_STATUS_ACCEPTED), int32(resqualifications.RequestStatus_REQUEST_STATUS_ACCEPTED)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(
		t,
		store.UpdateRequestStatus(
			t.Context(),
			db,
			42,
			7,
			resqualifications.RequestStatus_REQUEST_STATUS_ACCEPTED,
		),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteQualificationRequestByAttemptIDScopesToAttempt(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_requests SET deleted_at = CURRENT_TIMESTAMP.*WHERE fivenet_qualifications_requests\\.exam_attempt_id = \\?.*LIMIT \\?;").
		WithArgs("attempt-older", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(
		t,
		store.DeleteQualificationRequestByAttemptID(t.Context(), db, "attempt-older"),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreRestoreQualificationRequestScopesToAttempt(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_requests SET deleted_at = NULL, status = \\?.*WHERE .*qualification_id = \\?.*user_id = \\?.*exam_attempt_id = \\?.*LIMIT \\?;").
		WithArgs(int32(resqualifications.RequestStatus_REQUEST_STATUS_COMPLETED), int64(42), int32(7), "attempt-current", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RestoreQualificationRequest(
		t.Context(), db, 42, 7,
		resqualifications.RequestStatus_REQUEST_STATUS_COMPLETED,
		"attempt-current",
	))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreRestoreQualificationRequestWithoutAttemptID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_requests SET deleted_at = NULL, status = \\?.*WHERE .*qualification_id = \\?.*user_id = \\?.*\\(.*exam_attempt_id IS NULL.*OR.*exam_attempt_id = \\?.*\\).*LIMIT \\?;").
		WithArgs(int32(resqualifications.RequestStatus_REQUEST_STATUS_DENIED), int64(42), int32(7), "", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RestoreQualificationRequest(
		t.Context(), db, 42, 7,
		resqualifications.RequestStatus_REQUEST_STATUS_DENIED,
		"",
	))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreSetQualificationRequestExamAttemptIDIgnoresDeletedRequest(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_requests SET exam_attempt_id = \\?.*WHERE .*qualification_id = \\?.*user_id = \\?.*deleted_at IS NULL.*LIMIT \\?;").
		WithArgs("attempt-new", int64(42), int32(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	require.NoError(
		t,
		store.SetQualificationRequestExamAttemptID(t.Context(), db, 42, 7, "attempt-new"),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpsertQualificationRequestClearsPreviousExamAttempt(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	mock.ExpectExec("(?s)INSERT INTO .*fivenet_qualifications_requests.*ON DUPLICATE KEY UPDATE.*exam_attempt_id = NULL.*").
		WithArgs(
			int64(42), int32(7), nil,
			int32(resqualifications.RequestStatus_REQUEST_STATUS_PENDING),
			int32(resqualifications.RequestStatus_REQUEST_STATUS_PENDING),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(
		t,
		store.UpsertQualificationRequest(t.Context(), db, &resqualifications.QualificationRequest{
			QualificationId: 42,
			UserId:          7,
			UserComment:     nil,
		}),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListQualificationRequestsUsesVisibilityCte(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))

	countQuery := regexp.QuoteMeta(`WITH user_subjects AS`) +
		`(?s).*` + regexp.QuoteMeta(`visible_sources AS`) +
		`(?s).*` + regexp.QuoteMeta(`winning_visibility AS`) +
		`(?s).*` + regexp.QuoteMeta(`COUNT(DISTINCT qualification_request.user_id) AS "data_count.total"`)
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	pageSize := int64(10)
	resp, err := store.ListQualificationRequests(
		t.Context(),
		ListQualificationRequestsOptions{
			Pagination:      &database.PaginationRequest{PageSize: &pageSize},
			QualificationID: 42,
		},
		&userinfo.UserInfo{UserId: 7, Job: "doj"},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}
