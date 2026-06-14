package documentsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListDocumentActivityEmpty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_activity AS doc_activity.*`).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	count, items, err := store.ListDocumentActivity(t.Context(), ListDocumentActivityQuery{
		DocumentID: 42,
		Pagination: &resourcesdatabase.PaginationRequest{},
		UserInfo:   &userinfo.UserInfo{UserId: 3, Job: "doj"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count.Total)
	assert.Empty(t, items)
	require.NoError(t, mock.ExpectationsWereMet())
}
