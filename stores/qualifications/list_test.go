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
