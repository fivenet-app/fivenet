package qualificationsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateQualificationResultUpsertsSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectExec("(?s)INSERT INTO .*fivenet_qualifications_results.*").
		WithArgs(int64(42), int32(7), int32(resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL), nil, "ok", int32(1), "police").
		WillReturnResult(sqlmock.NewResult(99, 1))
	mock.ExpectExec("(?s)DELETE FROM .*fivenet_qualifications_result_success_map.*LIMIT \\?.*").
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)INSERT INTO .*fivenet_qualifications_result_success_map.*ON DUPLICATE KEY UPDATE.*").
		WithArgs(int64(42), int32(7), int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	creator := &userinfo.UserInfo{UserId: 1, Job: "police"}
	resultID, err := store.CreateQualificationResult(
		t.Context(),
		db,
		42,
		7,
		resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
		nil,
		"ok",
		creator,
	)
	require.NoError(t, err)
	require.Equal(t, int64(99), resultID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateQualificationResultSuccessfulSwapsSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_results SET .*WHERE .*deleted_at IS NULL.*LIMIT \\?;").
		WithArgs(int64(42), int32(7), int32(resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL), nil, "ok", int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)DELETE FROM .*fivenet_qualifications_result_success_map.*LIMIT \\?.*").
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)INSERT INTO .*fivenet_qualifications_result_success_map.*ON DUPLICATE KEY UPDATE.*").
		WithArgs(int64(42), int32(7), int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(
		t,
		store.UpdateQualificationResult(
			t.Context(),
			db,
			42,
			99,
			7,
			resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
			nil,
			"ok",
		),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateQualificationResultNonSuccessfulDeletesSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_results SET .*WHERE .*deleted_at IS NULL.*LIMIT \\?;").
		WithArgs(int64(42), int32(7), int32(resqualifications.ResultStatus_RESULT_STATUS_FAILED), nil, "nope", int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)DELETE FROM .*fivenet_qualifications_result_success_map.*LIMIT \\?.*").
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(
		t,
		store.UpdateQualificationResult(
			t.Context(),
			db,
			42,
			99,
			7,
			resqualifications.ResultStatus_RESULT_STATUS_FAILED,
			nil,
			"nope",
		),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteQualificationResultDeletesSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectExec("(?s)UPDATE fivenet_qualifications_results SET deleted_at = CURRENT_TIMESTAMP .*LIMIT \\?;").
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)DELETE FROM .*fivenet_qualifications_result_success_map.*LIMIT \\?.*").
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteQualificationResult(t.Context(), db, 99))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteQualificationClearsAndRebuildsSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	now := timestamp.Now()

	mock.ExpectExec("(?s)UPDATE fivenet_qualifications SET .*deleted_at.*WHERE fivenet_qualifications\\.id = \\?.*LIMIT \\?;").
		WithArgs(now.AsTime(), int64(42), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)DELETE FROM .*fivenet_qualifications_result_success_map.*WHERE .*qualification_id = \\?.*").
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteQualification(t.Context(), db, 42, now))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteQualificationRestoreRebuildsSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	mock.ExpectExec("(?s)UPDATE fivenet_qualifications SET .*deleted_at = NULL.*WHERE fivenet_qualifications\\.id = \\?.*LIMIT \\?;").
		WithArgs(int64(42), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("(?s)DELETE FROM .*fivenet_qualifications_result_success_map.*WHERE .*qualification_id = \\?.*").
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("(?s)SELECT qualification_result\\.qualification_id AS \"qualification_id\".*MAX\\(qualification_result\\.id\\) AS \"result_id\".*GROUP BY qualification_result\\.qualification_id, qualification_result\\.user_id").
		WithArgs(int64(42), int32(resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)).
		WillReturnRows(sqlmock.NewRows([]string{"qualification_id", "user_id", "result_id"}).
			AddRow(int64(42), int32(7), int64(99)))
	mock.ExpectExec("(?s)INSERT INTO .*fivenet_qualifications_result_success_map.*ON DUPLICATE KEY UPDATE.*").
		WithArgs(int64(42), int32(7), int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteQualification(t.Context(), db, 42, nil))
	require.NoError(t, mock.ExpectationsWereMet())
}
