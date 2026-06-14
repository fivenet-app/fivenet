package documentsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListUserDocuments(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	countRows := sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(DISTINCT document_relation.document_id) AS "data_count.total" FROM fivenet_documents_relations AS document_relation`) + `(?s).*`).
		WillReturnRows(countRows)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_relations AS document_relation`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY document.created_at DESC`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?`)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	pag, relations, err := store.ListUserDocuments(t.Context(), ListUserDocumentsQuery{
		UserID:         3,
		IncludeCreated: true,
		Relations: []documentsrelations.DocRelation{
			documentsrelations.DocRelation_DOC_RELATION_MENTIONED,
		},
		Pagination: &resourcesdatabase.PaginationRequest{},
		UserInfo:   &userinfo.UserInfo{UserId: 3, Job: "doj"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), pag.Total)
	assert.Empty(t, relations)
	require.NoError(t, mock.ExpectationsWereMet())
}
