package filestore

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFilestoreDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
		require.NoError(t, mock.ExpectationsWereMet())
	})
	return db, mock
}

func TestGetFileByPathReturnsIDAndPath(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectQuery(`SELECT .*file\.id.*file_path.*FROM fivenet_files.*WHERE .*file_path = \?.*LIMIT \?;`).
		WithArgs("uploads/file.pdf", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "file_path"}).AddRow(int64(42), "uploads/file.pdf"))

	h := &Handler[string]{db: db}
	id, path, err := h.GetFileByPath(t.Context(), FilestoreURLPrefix+"uploads/file.pdf")
	require.NoError(t, err)
	assert.Equal(t, int64(42), id)
	assert.Equal(t, "uploads/file.pdf", path)
}

func TestGetFileByPathNotFound(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectQuery(`SELECT .*file\.id.*file_path.*FROM fivenet_files.*WHERE .*file_path = \?.*LIMIT \?;`).
		WithArgs("missing", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "file_path"}))

	h := &Handler[string]{db: db}
	id, path, err := h.GetFileByPath(t.Context(), "missing")
	require.NoError(t, err)
	assert.Zero(t, id)
	assert.Empty(t, path)
}

func TestUpsertFileRowInsertsNewFile(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*file_path = \?.*FOR UPDATE;`).
		WithArgs("uploads/new.txt").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(`INSERT INTO fivenet_files .*file_path.*byte_size.*content_type.*meta`).
		WithArgs("uploads/new.txt", int64(12), "text/plain").
		WillReturnResult(sqlmock.NewResult(17, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)
	id, err := upsertFileRow(t.Context(), tx, "uploads/new.txt", "text/plain", 12)
	require.NoError(t, err)
	assert.Equal(t, int64(17), id)
	require.NoError(t, tx.Commit())
}

func TestUpsertFileRowUpdatesExistingFile(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*file_path = \?.*FOR UPDATE;`).
		WithArgs("uploads/existing.txt").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(9)))
	mock.ExpectExec(`UPDATE fivenet_files SET .*byte_size.*content_type.*deleted_at.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(20), "text/plain", int64(9), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)
	id, err := upsertFileRow(t.Context(), tx, "uploads/existing.txt", "text/plain", 20)
	require.NoError(t, err)
	assert.Equal(t, int64(9), id)
	require.NoError(t, tx.Commit())
}

func TestInsertJoinRow(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO fivenet_documents_files .*document_id.*file_id`).
		WithArgs(int64(3), int64(8)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)
	err = InsertJoinRow(
		t.Context(),
		tx,
		table.FivenetDocumentsFiles,
		table.FivenetDocumentsFiles.DocumentID,
		table.FivenetDocumentsFiles.FileID,
		int64(3),
		mysql.Bool(true),
		8,
	)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
}
