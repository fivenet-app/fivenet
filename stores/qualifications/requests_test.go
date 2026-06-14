package qualificationsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/stretchr/testify/require"
)

func TestStoreUpdateRequestStatus(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

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
