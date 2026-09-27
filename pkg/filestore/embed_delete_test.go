package filestore

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteOrphanedFileRemovesAssociationStorageAndRow(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_documents_files .*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*file_id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectQuery(`SELECT .*fivenet_files\.file_path.*WHERE .*id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"file_path"}).AddRow("uploads/file.txt"))
	mock.ExpectExec(`UPDATE fivenet_files SET .*deleted_at = CURRENT_TIMESTAMP.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	h := newDeleteHandler(db, st, false)
	require.NoError(t, h.Delete(t.Context(), int64(7), 23))
	assert.Equal(t, "uploads/file.txt", st.deleted)
}

func TestDeleteReferencedFileKeepsFileRowAndStorageObject(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_documents_files .*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*file_id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))
	mock.ExpectCommit()

	h := newDeleteHandler(db, st, false)
	require.NoError(t, h.Delete(t.Context(), int64(7), 23))
	assert.Empty(t, st.deleted)
}

func TestDeleteWithZeroParentSkipsAssociationDeletion(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*file_id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))
	mock.ExpectCommit()

	h := newDeleteHandler(db, st, false)
	require.NoError(t, h.Delete(t.Context(), int64(0), 23))
}

func TestDeleteNullsFileReferenceInsteadOfDeletingJoinRow(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE fivenet_documents_files SET .*document_id = \?.*file_id = NULL.*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(7), int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*file_id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))
	mock.ExpectCommit()

	h := newDeleteHandler(db, st, true)
	require.NoError(t, h.Delete(t.Context(), int64(7), 23))
	assert.Empty(t, st.deleted)
}

func TestDeleteStorageFailureLeavesDatabaseTransactionUncommitted(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{deleteErr: errors.New("storage unavailable")}

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_documents_files .*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*file_id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectQuery(`SELECT .*fivenet_files\.file_path.*WHERE .*id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"file_path"}).AddRow("uploads/file.txt"))
	mock.ExpectRollback()

	h := newDeleteHandler(db, st, false)
	require.ErrorIs(t, h.Delete(t.Context(), int64(7), 23), st.deleteErr)
}

func TestDeleteFileByPathReturnsWithoutDeletingMissingFile(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectQuery(`SELECT .*file\.id.*file_path.*FROM fivenet_files.*WHERE .*file_path = \?.*LIMIT \?;`).
		WithArgs("missing", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "file_path"}))

	h := newDeleteHandler(db, &uploadStorage{}, false)
	require.NoError(t, h.DeleteFileByPath(t.Context(), int64(7), "missing"))
}

func TestDeleteFileByPathDeletesExistingFile(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}

	mock.ExpectQuery(`SELECT .*file\.id.*file_path.*FROM fivenet_files.*WHERE .*file_path = \?.*LIMIT \?;`).
		WithArgs("uploads/file.txt", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "file_path"}).AddRow(int64(23), "uploads/file.txt"))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_documents_files .*WHERE .*document_id = \?.*file_id = \?.*LIMIT \?;`).
		WithArgs(int64(7), int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*file_id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))
	mock.ExpectQuery(`SELECT .*fivenet_files\.file_path.*WHERE .*id = \?;`).
		WithArgs(int64(23)).
		WillReturnRows(sqlmock.NewRows([]string{"file_path"}).AddRow("uploads/file.txt"))
	mock.ExpectExec(`UPDATE fivenet_files SET .*deleted_at = CURRENT_TIMESTAMP.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	h := newDeleteHandler(db, st, false)
	require.NoError(
		t,
		h.DeleteFileByPath(t.Context(), int64(7), FilestoreURLPrefix+"uploads/file.txt"),
	)
	assert.Equal(t, "uploads/file.txt", st.deleted)
}

func TestUpdateJoinRow(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE fivenet_documents_files SET .*file_id = \?.*WHERE .*document_id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)
	err = UpdateJoinRow(
		t.Context(),
		tx,
		table.FivenetDocumentsFiles,
		table.FivenetDocumentsFiles.DocumentID,
		table.FivenetDocumentsFiles.FileID,
		int64(7),
		table.FivenetDocumentsFiles.DocumentID.EQ(mysql.Int64(7)),
		23,
	)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
}

func newDeleteHandler(db *sql.DB, st *uploadStorage, nullOnly bool) *Handler[int64] {
	return &Handler[int64]{
		db:                db,
		store:             st,
		joinTable:         table.FivenetDocumentsFiles,
		parentCol:         table.FivenetDocumentsFiles.DocumentID,
		fileCol:           table.FivenetDocumentsFiles.FileID,
		nullOnlyParentRow: nullOnly,
		parentColBoolExp:  func(id int64) mysql.BoolExpression { return table.FivenetDocumentsFiles.DocumentID.EQ(mysql.Int64(id)) },
	}
}
