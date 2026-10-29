package filestore

import (
	"context"
	"errors"
	"time"

	filepb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// FileReference is a file joined to an entity through a relation table that
// also carries a kind (for example, a mugshot or attachment).
type FileReference struct {
	File *filepb.File
	Kind string
}

// ListFileReferences lists typed file relations for one parent. The relation
// table and columns are supplied by the caller because the parent key may be
// an integer or a string, while the file metadata and mapping are shared.
func ListFileReferences(
	ctx context.Context,
	db qrm.Queryable,
	join mysql.Table,
	parentCondition mysql.BoolExpression,
	fileID mysql.ColumnInteger,
	kind mysql.ColumnString,
	createdAt mysql.ColumnTimestamp,
) ([]*FileReference, error) {
	rows := []struct {
		ID          int64     `alias:"file.id"`
		FilePath    string    `alias:"file.file_path"`
		ByteSize    int64     `alias:"file.byte_size"`
		ContentType string    `alias:"file.content_type"`
		CreatedAt   time.Time `alias:"file.created_at"`
		Kind        string    `alias:"file.kind"`
	}{}

	stmt := join.
		SELECT(
			table.FivenetFiles.ID.AS("file.id"),
			table.FivenetFiles.FilePath.AS("file.file_path"),
			table.FivenetFiles.ByteSize.AS("file.byte_size"),
			table.FivenetFiles.ContentType.AS("file.content_type"),
			table.FivenetFiles.CreatedAt.AS("file.created_at"),
			kind.AS("file.kind"),
		).
		FROM(join.INNER_JOIN(table.FivenetFiles, fileID.EQ(table.FivenetFiles.ID))).
		WHERE(parentCondition).
		ORDER_BY(createdAt.ASC())

	if err := stmt.QueryContext(ctx, db, &rows); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}

	result := make([]*FileReference, 0, len(rows))
	for _, row := range rows {
		result = append(result, &FileReference{
			Kind: row.Kind,
			File: &filepb.File{
				Id:          row.ID,
				FilePath:    row.FilePath,
				ByteSize:    row.ByteSize,
				ContentType: row.ContentType,
				CreatedAt:   timestamp.New(row.CreatedAt),
			},
		})
	}
	return result, nil
}

// FindFileReference returns a related file by ID, or nil when it is not
// related to the requested parent.
func FindFileReference(
	ctx context.Context,
	db qrm.Queryable,
	join mysql.Table,
	parentCondition mysql.BoolExpression,
	fileID mysql.ColumnInteger,
	kind mysql.ColumnString,
	targetID int64,
) (*FileReference, error) {
	row := struct {
		ID          int64     `alias:"file.id"`
		FilePath    string    `alias:"file.file_path"`
		ByteSize    int64     `alias:"file.byte_size"`
		ContentType string    `alias:"file.content_type"`
		CreatedAt   time.Time `alias:"file.created_at"`
		Kind        string    `alias:"file.kind"`
	}{}

	stmt := join.
		SELECT(
			table.FivenetFiles.ID.AS("file.id"),
			table.FivenetFiles.FilePath.AS("file.file_path"),
			table.FivenetFiles.ByteSize.AS("file.byte_size"),
			table.FivenetFiles.ContentType.AS("file.content_type"),
			table.FivenetFiles.CreatedAt.AS("file.created_at"),
			kind.AS("file.kind"),
		).
		FROM(join.INNER_JOIN(table.FivenetFiles, fileID.EQ(table.FivenetFiles.ID))).
		WHERE(mysql.AND(parentCondition, fileID.EQ(mysql.Int64(targetID)))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, db, &row); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &FileReference{
		Kind: row.Kind,
		File: &filepb.File{
			Id:          row.ID,
			FilePath:    row.FilePath,
			ByteSize:    row.ByteSize,
			ContentType: row.ContentType,
			CreatedAt:   timestamp.New(row.CreatedAt),
		},
	}, nil
}

func (r *FileReference) GetFileID() int64 {
	if r == nil || r.File == nil {
		return 0
	}
	return r.File.GetId()
}
