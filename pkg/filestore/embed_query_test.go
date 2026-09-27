package filestore

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountFilesForParentID(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\) AS "data_count\.total" FROM fivenet_documents_files WHERE .*document_id = \?;`).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(3)))

	count, err := h.CountFilesForParentID(t.Context(), 7)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestListFilesForParentID(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)
	createdAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`SELECT .*file\.id.*file\.parent_id.*file\.file_path.*file\.byte_size.*file\.content_type.*file\.created_at.*FROM fivenet_documents_files INNER JOIN fivenet_files.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{
			"file.id",
			"file.parent_id",
			"file.file_path",
			"file.byte_size",
			"file.content_type",
			"file.created_at",
		}).AddRow(int64(23), int64(7), "uploads/file.txt", int64(5), "text/plain", createdAt))

	files, err := h.ListFilesForParentID(t.Context(), 7)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, int64(23), files[0].GetId())
	assert.Equal(t, int64(7), files[0].GetParentId())
	assert.Equal(t, "uploads/file.txt", files[0].GetFilePath())
	assert.Equal(t, int64(5), files[0].GetByteSize())
	assert.Equal(t, "text/plain", files[0].GetContentType())
}

func TestListFilesForParentIDEmpty(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT .*FROM fivenet_documents_files INNER JOIN fivenet_files.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	files, err := h.ListFilesForParentID(t.Context(), 7)
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestHandleFileChangesForParentNoOp(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT .*FROM fivenet_documents_files INNER JOIN fivenet_files.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id"}).AddRow(int64(23)))

	mock.ExpectBegin()
	mock.ExpectRollback()
	tx, err := db.Begin()
	require.NoError(t, err)
	added, deleted, err := h.HandleFileChangesForParent(t.Context(), tx, 7, []*file.File{{Id: 23}})
	require.NoError(t, err)
	assert.Zero(t, added)
	assert.Zero(t, deleted)
	require.NoError(t, tx.Rollback())
}

func TestHandleFileChangesForParentAddsFilesWithoutInsert(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT .*FROM fivenet_documents_files INNER JOIN fivenet_files.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectBegin()
	mock.ExpectRollback()
	tx, err := db.Begin()
	require.NoError(t, err)
	added, deleted, err := h.HandleFileChangesForParent(t.Context(), tx, 7, []*file.File{{Id: 23}})
	require.NoError(t, err)
	assert.Equal(t, int64(1), added)
	assert.Zero(t, deleted)
	require.NoError(t, tx.Rollback())
}

func TestHandleFileChangesForParentDeletesRemovedFiles(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT .*FROM fivenet_documents_files INNER JOIN fivenet_files.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id"}).AddRow(int64(23)))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_documents_files .*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)
	added, deleted, err := h.HandleFileChangesForParent(t.Context(), tx, 7, nil)
	require.NoError(t, err)
	assert.Zero(t, added)
	assert.Equal(t, int64(1), deleted)
	require.NoError(t, tx.Rollback())
}

func TestHandleFileChangesForParentReturnsDeleteError(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	h := newDeleteHandler(db, &uploadStorage{}, false)
	mock.MatchExpectationsInOrder(false)
	deleteErr := errors.New("association delete failed")

	mock.ExpectQuery(`SELECT .*FROM fivenet_documents_files INNER JOIN fivenet_files.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id"}).AddRow(int64(23)))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_documents_files .*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(23), int64(1)).
		WillReturnError(deleteErr)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)
	_, _, err = h.HandleFileChangesForParent(t.Context(), tx, 7, nil)
	require.ErrorIs(t, err, deleteErr)
	require.NoError(t, tx.Rollback())
}
