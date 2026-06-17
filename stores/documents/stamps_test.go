package documentsstore

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListUsableStampsUsesVisibilityCteForNonSuperuser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	countQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*COUNT\(doc_ids\.id\) AS "data_count\.total".*`
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	listQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*INNER JOIN fivenet_documents_stamps AS stamp ON.*ORDER BY stamp\.sort_key ASC, stamp\.created_at DESC LIMIT \? OFFSET \?;`
	mock.ExpectQuery(listQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"stamp.id",
			"stamp.name",
			"stamp.svg_template",
			"stamp.variants_json",
			"stamp.created_at",
		}).AddRow(int64(11), "fire", "<svg/>", nil, time.Unix(0, 0).UTC()))

	pageSize := int64(20)
	pag, stamps, err := store.ListUsableStamps(t.Context(), ListUsableStampsQuery{
		UserInfo:   &userinfo.UserInfo{UserId: 3, Job: "doj", JobGrade: 16},
		Pagination: &resourcesdatabase.PaginationRequest{PageSize: &pageSize},
	})
	require.NoError(t, err)
	assert.NotNil(t, pag)
	require.Len(t, stamps, 1)
	assert.Equal(t, int64(11), stamps[0].GetId())
	assert.Equal(t, "fire", stamps[0].GetName())
	require.NoError(t, mock.ExpectationsWereMet())
}
