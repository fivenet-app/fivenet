package documentsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetDocumentPinMergesRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_pins AS document_pin.*document_pin\.document_id = \?.*document_pin\.job = \?.*document_pin\.user_id = \?.*LIMIT \?.*`).
		WithArgs(int64(42), "doj", int32(3), int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"document_pin.document_id", "document_pin.job", "document_pin.user_id", "document_pin.created_at", "document_pin.state", "document_pin.creator_id"}).
			AddRow(int64(42), nil, nil, nil, true, int32(7)).
			AddRow(int64(42), "doj", int32(3), nil, false, int32(8)))

	pin, err := store.GetDocumentPin(t.Context(), 42, &userinfo.UserInfo{UserId: 3, Job: "doj"})
	require.NoError(t, err)
	require.NotNil(t, pin)
	assert.Equal(t, int64(42), pin.GetDocumentId())
	assert.Equal(t, "doj", pin.GetJob())
	assert.Equal(t, int32(3), pin.GetUserId())
	assert.False(t, pin.GetState())
	assert.Equal(t, int32(8), pin.GetCreatorId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListDocumentPinsEmpty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_pins AS document_pin.*`).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))

	pag, docs, err := store.ListDocumentPins(t.Context(), ListDocumentPinsQuery{
		Pagination: &resourcesdatabase.PaginationRequest{},
		UserInfo:   &userinfo.UserInfo{UserId: 3, Job: "doj"},
	})
	require.NoError(t, err)
	require.NotNil(t, pag)
	assert.Empty(t, docs)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreToggleDocumentPinWrites(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	userInfo := &userinfo.UserInfo{UserId: 3, Job: "doj"}

	mock.ExpectExec(`(?s).*INSERT INTO fivenet_documents_pins.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.CreateDocumentPin(t.Context(), db, 42, userInfo, false))

	mock.ExpectExec(`(?s).*DELETE FROM fivenet_documents_pins.*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.DeleteDocumentPin(t.Context(), db, 42, userInfo, false))

	require.NoError(t, mock.ExpectationsWereMet())
}
