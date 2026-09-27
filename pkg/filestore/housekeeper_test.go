package filestore

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

func newTestHousekeeper(db *sql.DB, st *uploadStorage) *Housekeeper {
	return &Housekeeper{
		logger:          zap.NewNop(),
		tracer:          noop.NewTracerProvider().Tracer("filestore-test"),
		db:              db,
		storage:         st,
		getTablesListFn: func() []joinInfo { return nil },
		gracePeriod:     24 * time.Hour,
		batchSize:       100,
	}
}

func TestHousekeeperRunWithNoCandidates(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*ORDER BY .*id.*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}))

	h := newTestHousekeeper(db, &uploadStorage{})
	deleted, err := h.Run(t.Context())
	require.NoError(t, err)
	assert.Zero(t, deleted)
}

func TestHousekeeperRunDeletesStorageAndDatabaseRows(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*ORDER BY .*id.*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}).AddRow(int64(23), "uploads/file.txt"))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_files .*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*ORDER BY .*id.*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}))

	h := newTestHousekeeper(db, st)
	deleted, err := h.Run(t.Context())
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted)
	assert.Equal(t, "uploads/file.txt", st.deleted)
}

func TestHousekeeperRunDoesNotCountStorageDeletionFailures(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{deleteErr: errors.New("storage unavailable")}
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*ORDER BY .*id.*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}).AddRow(int64(23), "uploads/file.txt"))
	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*ORDER BY .*id.*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}))

	h := newTestHousekeeper(db, st)
	deleted, err := h.Run(t.Context())
	require.NoError(t, err)
	assert.Zero(t, deleted)
}

func TestHousekeeperRunRollsBackOnDatabaseDeleteFailure(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	databaseErr := errors.New("database unavailable")
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*ORDER BY .*id.*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}).AddRow(int64(23), "uploads/file.txt"))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_files .*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(1)).
		WillReturnError(databaseErr)
	mock.ExpectRollback()

	h := newTestHousekeeper(db, &uploadStorage{})
	deleted, err := h.Run(t.Context())
	require.Error(t, err)
	assert.Zero(t, deleted)
	assert.ErrorContains(t, err, databaseErr.Error())
}

func TestHousekeeperRunIncludesJoinTableOrphanChecks(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)

	// Both tables must be absent for a non-soft-deleted file to be considered orphaned.
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_files.*WHERE .*NOT .*EXISTS.*SELECT 1.*FROM fivenet_documents_files.*file_id = fivenet_files\.id.*NOT .*EXISTS.*SELECT 1.*FROM fivenet_wiki_pages_files.*file_id = fivenet_files\.id.*ORDER BY .*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}))

	h := newTestHousekeeper(db, &uploadStorage{})
	h.getTablesListFn = func() []joinInfo {
		return []joinInfo{
			{Table: table.FivenetDocumentsFiles, FileCol: table.FivenetDocumentsFiles.FileID},
			{Table: table.FivenetWikiPagesFiles, FileCol: table.FivenetWikiPagesFiles.FileID},
		}
	}

	deleted, err := h.Run(t.Context())
	require.NoError(t, err)
	assert.Zero(t, deleted)
}

func TestHousekeeperRunDeletesOrphanWithJoinTables(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}
	candidateQuery := `(?s)SELECT .*FROM fivenet_files.*WHERE .*NOT .*EXISTS.*SELECT 1.*FROM fivenet_documents_files.*file_id = fivenet_files\.id.*ORDER BY .*LIMIT \?;`
	mock.ExpectQuery(candidateQuery).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}).AddRow(int64(23), "uploads/orphan.txt"))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM fivenet_files .*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(candidateQuery).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(sqlmock.NewRows([]string{"file.id", "file.file_path"}))

	h := newTestHousekeeper(db, st)
	h.getTablesListFn = func() []joinInfo {
		return []joinInfo{{
			Table:   table.FivenetDocumentsFiles,
			FileCol: table.FivenetDocumentsFiles.FileID,
		}}
	}

	deleted, err := h.Run(t.Context())
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted)
	assert.Equal(t, "uploads/orphan.txt", st.deleted)
}

func TestHousekeeperRunStopsAtMaximumBatchAttempts(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	rows := sqlmock.NewRows([]string{"file.id", "file.file_path"})
	for id := int64(1); id <= 200; id++ {
		rows.AddRow(id, "uploads/orphan.txt")
	}
	mock.ExpectQuery(`(?s)SELECT .*FROM fivenet_files.*WHERE .*NOT .*EXISTS.*SELECT 1.*FROM fivenet_documents_files.*file_id = fivenet_files\.id.*ORDER BY .*LIMIT \?;`).
		WithArgs(sqlmock.AnyArg(), int64(100)).
		WillReturnRows(rows)
	mock.ExpectBegin()
	for id := int64(1); id <= 200; id++ {
		mock.ExpectExec(`DELETE FROM fivenet_files .*WHERE .*id = \?.*LIMIT \?;`).
			WithArgs(id, int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()

	h := newTestHousekeeper(db, &uploadStorage{})
	h.getTablesListFn = func() []joinInfo {
		return []joinInfo{{
			Table:   table.FivenetDocumentsFiles,
			FileCol: table.FivenetDocumentsFiles.FileID,
		}}
	}

	deleted, err := h.Run(t.Context())
	require.NoError(t, err)
	assert.Equal(t, int64(200), deleted)
}
