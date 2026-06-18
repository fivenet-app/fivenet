package qualificationsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListQualificationsUsesVisibilityCte(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	countQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*COUNT\(DISTINCT qualification\.id\) AS "data_count\.total".*qualification_result\.deleted_at IS NULL.*`
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	listQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*ORDER BY qualification\.draft ASC, qualification_result\.id DESC LIMIT \? OFFSET \?;`
	mock.ExpectQuery(listQuery).
		WillReturnRows(sqlmock.NewRows([]string{}))

	pageSize := int64(10)
	resp, err := store.ListQualifications(
		t.Context(),
		ListQualificationsOptions{
			Pagination: &database.PaginationRequest{PageSize: &pageSize},
		},
		&userinfo.UserInfo{UserId: 7, Job: "doj", JobGrade: 16},
		false,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.GetPagination().GetTotalCount())
	assert.Empty(t, resp.GetQualifications())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCheckRequirementsMetForQualificationUsesSuccessMap(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	query := `(?s).*FROM fivenet_qualifications_requirements AS qualification_requirement LEFT JOIN fivenet_qualifications_result_success_map AS qualification_result_success_map ON .*qualification_result_success_map\.qualification_id = qualification_requirement\.target_qualification_id.*qualification_result_success_map\.user_id = \?.*`
	mock.ExpectQuery(query).
		WithArgs(int32(7), int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"qualification_id", "userid"}))

	ok, err := store.CheckRequirementsMetForQualification(t.Context(), 42, 7)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}
