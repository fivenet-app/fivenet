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

func TestStoreListQualificationRequestsUsesVisibilityCte(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

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
		false,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}
