package documentsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentsrequests "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/requests"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreRequestCRUD(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	request := &documentsrequests.DocRequest{
		DocumentId:  42,
		RequestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE,
		CreatorId:   func() *int32 { v := int32(3); return &v }(),
		CreatorJob:  "doj",
	}

	mock.ExpectExec(`(?s).*INSERT INTO fivenet_documents_requests.*`).
		WithArgs(int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE, int32(3), "doj", nil, nil, nil).
		WillReturnResult(sqlmock.NewResult(7, 1))
	lastID, err := store.AddDocumentReq(t.Context(), db, request)
	require.NoError(t, err)
	assert.Equal(t, int64(7), lastID)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_requests AS doc_request.*`).
		WithArgs(int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"doc_request.id", "doc_request.created_at", "doc_request.updated_at", "doc_request.document_id", "doc_request.request_type", "doc_request.creator_id", "doc_request.creator_job", "doc_request.reason", "doc_request.data", "creator.id", "creator.firstname", "creator.lastname", "creator.job", "creator.dateofbirth", "creator.phone_number"}).AddRow(int64(7), nil, nil, int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE, int32(3), "doj", nil, nil, int32(3), "A", "B", "doj", nil, nil))
	got, err := store.GetDocumentReq(
		t.Context(),
		db,
		table.FivenetDocumentsRequests.ID.EQ(mysql.Int64(7)),
	)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(7), got.GetId())

	mock.ExpectExec(`(?s).*UPDATE fivenet_documents_requests.*`).
		WithArgs(int64(42), documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE, int32(3), "doj", nil, nil, nil, int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.UpdateDocumentReq(t.Context(), db, 7, request))

	mock.ExpectExec(`(?s).*DELETE FROM fivenet_documents_requests.*`).
		WithArgs(int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.DeleteDocumentReq(t.Context(), db, 7))

	require.NoError(t, mock.ExpectationsWereMet())
}
