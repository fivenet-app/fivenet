package filestore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/require"
)

func TestListFileReferences(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	createdAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

	expectVehicleFileReferencesQuery(mock, sqlmock.NewRows([]string{
		"file.id",
		"file.file_path",
		"file.byte_size",
		"file.content_type",
		"file.created_at",
		"file.kind",
	}).AddRow(23, "uploads/file.txt", 5, "text/plain", createdAt, "attachment"))

	join := table.FivenetVehicleFiles
	references, err := ListFileReferences(
		t.Context(),
		db,
		join,
		join.Plate.EQ(mysql.String("ABC123")),
		join.FileID,
		join.Kind,
		join.CreatedAt,
	)
	require.NoError(t, err)
	require.Len(t, references, 1)
	require.Equal(t, "attachment", references[0].Kind)
	require.NotNil(t, references[0].File)
	require.Equal(t, int64(23), references[0].File.GetId())
	require.Equal(t, "uploads/file.txt", references[0].File.GetFilePath())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindFileReference(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	createdAt := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	expectVehicleFileReferenceQuery(mock, 23, sqlmock.NewRows([]string{
		"file.id",
		"file.file_path",
		"file.byte_size",
		"file.content_type",
		"file.created_at",
		"file.kind",
	}).AddRow(23, "uploads/file.txt", 5, "text/plain", createdAt, "attachment"))

	join := table.FivenetVehicleFiles
	reference, err := FindFileReference(
		t.Context(),
		db,
		join,
		join.Plate.EQ(mysql.String("ABC123")),
		join.FileID,
		join.Kind,
		23,
	)
	require.NoError(t, err)
	require.NotNil(t, reference)
	require.Equal(t, int64(23), reference.GetFileID())
	require.Equal(t, "attachment", reference.Kind)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindFileReferenceReturnsNilWhenNotRelated(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	expectVehicleFileReferenceQuery(mock, 23, sqlmock.NewRows([]string{
		"file.id",
		"file.file_path",
		"file.byte_size",
		"file.content_type",
		"file.created_at",
		"file.kind",
	}))

	join := table.FivenetVehicleFiles
	reference, err := FindFileReference(
		t.Context(),
		db,
		join,
		join.Plate.EQ(mysql.String("ABC123")),
		join.FileID,
		join.Kind,
		23,
	)
	require.NoError(t, err)
	require.Nil(t, reference)
	require.NoError(t, mock.ExpectationsWereMet())
}

func expectVehicleFileReferencesQuery(mock sqlmock.Sqlmock, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT fivenet_files.id AS \"file.id\", fivenet_files.file_path AS \"file.file_path\", fivenet_files.byte_size AS \"file.byte_size\", fivenet_files.content_type AS \"file.content_type\", fivenet_files.created_at AS \"file.created_at\", fivenet_vehicle_files.kind AS \"file.kind\" FROM fivenet_vehicle_files INNER JOIN fivenet_files ON (fivenet_vehicle_files.file_id = fivenet_files.id) WHERE fivenet_vehicle_files.plate = ? ORDER BY fivenet_vehicle_files.created_at ASC;",
	)).
		WithArgs("ABC123").
		WillReturnRows(rows)
}

func expectVehicleFileReferenceQuery(mock sqlmock.Sqlmock, fileID int64, rows *sqlmock.Rows) {
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT fivenet_files.id AS \"file.id\", fivenet_files.file_path AS \"file.file_path\", fivenet_files.byte_size AS \"file.byte_size\", fivenet_files.content_type AS \"file.content_type\", fivenet_files.created_at AS \"file.created_at\", fivenet_vehicle_files.kind AS \"file.kind\" FROM fivenet_vehicle_files INNER JOIN fivenet_files ON (fivenet_vehicle_files.file_id = fivenet_files.id) WHERE ( (fivenet_vehicle_files.plate = ?) AND (fivenet_vehicle_files.file_id = ?) ) LIMIT ?;",
	)).
		WithArgs("ABC123", fileID, int64(1)).
		WillReturnRows(rows)
}
