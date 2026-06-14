package documentsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListDocumentReqsEmpty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_requests AS doc_request.*`).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	count, reqs, err := store.ListDocumentReqs(t.Context(), ListDocumentReqsQuery{
		DocumentID: 42,
		UserInfo:   &userinfo.UserInfo{UserId: 3, Job: "doj"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count.Total)
	assert.Empty(t, reqs)
	require.NoError(t, mock.ExpectationsWereMet())
}
