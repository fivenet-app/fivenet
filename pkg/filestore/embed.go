package filestore

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbfilestore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tFiles = table.FivenetFiles

// ParentID is a type constraint that matches any type that can be a parent ID.
// It can be an unsigned 64-bit integer, a 32-bit integer, or a string.
type ParentID interface {
	~uint64 | ~int32 | ~string
	comparable
}

// ParentColBoolExpFn is a function type that returns a Jet boolean expression for a given parent ID.
type ParentColBoolExpFn[P ParentID] func(parentId P) jet.BoolExpression

// JoinRowInserterFn is a function type for inserting a join row into a table linking parent and file.
type JoinRowInserterFn[P ParentID] func(ctx context.Context, tx *sql.Tx, join jet.Table, parentCol jet.Column, fileCol jet.ColumnInteger, parentId P, _ jet.BoolExpression, fileID uint64) error

// Handler is a generic, embeddable file-upload helper.
type Handler[P ParentID] struct {
	// store is the storage adapter (e.g. S3, filesystem).
	store storage.IStorage
	// db is the SQL database connection.
	db *sql.DB
	// joinTable is the table that joins parent and file.
	joinTable jet.Table
	// parentCol is the column in joinTable for the parent (may be integer or string).
	parentCol jet.Column
	// fileCol is the column in joinTable for the file ID.
	fileCol jet.ColumnInteger
	// sizeLimit is the maximum allowed file size in bytes.
	sizeLimit int64

	// parentColBoolExp is a function that converts the parent ID to a Jet boolean expression.
	parentColBoolExp ParentColBoolExpFn[P]
	// joinRowInserter is a function for inserting a join row.
	joinRowInserter JoinRowInserterFn[P]
	// nullOnlyParentRow indicates if only the parent row should be nulled instead of deleted.
	nullOnlyParentRow bool
}

// NewHandler creates a new Handler for file uploads and associations.
func NewHandler[P ParentID](
	st storage.IStorage,
	db *sql.DB,
	join jet.Table,
	parentCol jet.Column,
	fileCol jet.ColumnInteger,
	sizeLimit int64,
	parentColBoolExp ParentColBoolExpFn[P],
	joinRowInserter JoinRowInserterFn[P],
	nullOnlyParentRow bool,
) *Handler[P] {
	AddTable(joinInfo{
		Table:   join,
		FileCol: fileCol,
	})

	return &Handler[P]{
		store:     st,
		db:        db,
		joinTable: join,
		parentCol: parentCol,
		fileCol:   fileCol,
		sizeLimit: sizeLimit,
		// parentColBoolExp is a function that converts the parent ID to a jet.BoolExpression
		parentColBoolExp:  parentColBoolExp,
		joinRowInserter:   joinRowInserter,
		nullOnlyParentRow: nullOnlyParentRow,
	}
}

// AwaitHandshake reads the first upload packet to extract metadata. No handshake logic is required for this handler.
func (h *Handler[P]) AwaitHandshake(
	srv pbfilestore.FilestoreService_UploadServer,
) (*file.UploadMeta, error) {
	// First packet must be metadata
	first, err := srv.Recv()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata: %v", err)
	}
	meta := first.GetMeta()
	if meta == nil {
		return nil, status.Error(codes.InvalidArgument, "first packet must contain meta")
	}

	if h.sizeLimit > 0 && meta.GetSize() > h.sizeLimit {
		return nil, status.Errorf(
			codes.ResourceExhausted,
			"file size exceeds limit: %d > %d",
			meta.GetSize(),
			h.sizeLimit,
		)
	}

	return meta, nil
}

// UploadFile streams file data from the gRPC server to storage, then records metadata and associations in the database.
func (h *Handler[P]) UploadFile(
	ctx context.Context,
	parentID P,
	key string,
	size int64,
	ctype string,
	srv pbfilestore.FilestoreService_UploadServer,
) (*file.UploadFileResponse, error) {
	if h.sizeLimit > 0 && size > h.sizeLimit {
		return nil, ErrUploadFileTooLarge(
			map[string]any{"maxSize": h.sizeLimit / 8},
		) // Convert bytes to megabytes
	}

	// pipe chunks to the storage backend
	pr, pw := io.Pipe()
	go func() {
		for {
			pkt, err := srv.Recv()
			if errors.Is(err, io.EOF) {
				_ = pw.Close()
				return
			}
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			if data := pkt.GetData(); data != nil {
				if _, w := pw.Write(data); w != nil {
					pw.CloseWithError(w)
					return
				}
			}
		}
	}()

	// Size - if client sent one, honour it; otherwise -1 (unknown)
	szHint := int64(-1)
	if size > 0 {
		szHint = size
	}

	urlKey, bytes, err := putToStorage(ctx, h.store, key, pr, ctype, szHint)
	if err != nil {
		return nil, err
	}

	// DB transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	fileID, err := upsertFileRow(ctx, tx, key, ctype, bytes)
	if err != nil {
		return nil, err
	}

	if h.joinTable != nil {
		parentIdBoolExp := h.parentColBoolExp(parentID)
		if err := h.joinRowInserter(ctx, tx, h.joinTable, h.parentCol, h.fileCol, parentID, parentIdBoolExp, fileID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	full, err := url.JoinPath(FilestoreURLPrefix, urlKey)
	if err != nil {
		return nil, err
	}

	resp := &file.UploadFileResponse{
		Id:  fileID,
		Url: full,
		File: &file.File{
			Id:          fileID,
			FilePath:    key,
			ContentType: ctype,
			ByteSize:    bytes,
			CreatedAt:   timestamp.Now(),
		},
	}
	return resp, srv.SendAndClose(resp)
}

// UploadFromMeta streams the remainder of the gRPC Upload after the caller has read & validated the first UploadMeta packet.
func (h *Handler[P]) UploadFromMeta(
	ctx context.Context,
	meta *file.UploadMeta,
	parentID P,
	srv pbfilestore.FilestoreService_UploadServer,
) (*file.UploadFileResponse, error) {
	key := buildKey(meta.GetNamespace(), SanitizeFileName(meta.GetOriginalName()))
	ctype := sniff(meta.GetContentType(), key)

	return h.UploadFile(ctx, parentID, key, meta.GetSize(), ctype, srv)
}

// GetFileByPath retrieves the file ID and path for a given file path string.
func (h *Handler[P]) GetFileByPath(ctx context.Context, path string) (uint64, string, error) {
	path = strings.TrimPrefix(path, FilestoreURLPrefix)

	tFiles := table.FivenetFiles.AS("file")

	stmt := tFiles.
		SELECT(
			tFiles.ID,
		).
		FROM(tFiles).
		WHERE(
			tFiles.FilePath.EQ(jet.String(path)),
		).
		LIMIT(1)

	f := &struct {
		ID       uint64 `jet:"id"`
		ParentID P      `jet:"parent_id"`
		FilePath string `jet:"file_path"`
	}{}
	if err := stmt.QueryContext(ctx, h.db, f); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return f.ID, "", err
		}
	}

	if f.ID == 0 {
		return 0, "", nil
	}

	return f.ID, f.FilePath, nil
}

// deleteJoinRow removes or nulls the join row for a parent/file association, depending on nullOnlyParentRow.
func (h *Handler[P]) deleteJoinRow(
	ctx context.Context,
	tx *sql.Tx,
	parentID P,
	fileID uint64,
) error {
	if h.nullOnlyParentRow {
		_, err := h.joinTable.
			UPDATE(
				h.parentCol,
				h.fileCol,
			).
			SET(
				parentID,
				jet.NULL,
			).
			WHERE(jet.AND(
				h.parentColBoolExp(parentID),
				h.fileCol.EQ(jet.Uint64(fileID)),
			)).
			ExecContext(ctx, tx)
		return err
	}

	_, err := h.joinTable.
		DELETE().
		WHERE(jet.AND(
			h.parentColBoolExp(parentID),
			h.fileCol.EQ(jet.Uint64(fileID)),
		)).
		ExecContext(ctx, tx)
	return err
}

// Delete removes the association between a parent and a file, and deletes the file if it is orphaned.
func (h *Handler[P]) Delete(ctx context.Context, parentID P, fileID uint64) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var zeroVal P
	if parentID != zeroVal {
		// 1. Remove or empty join row
		if err := h.deleteJoinRow(ctx, tx, parentID, fileID); err != nil {
			return err
		}
	}

	// 2. Is the file now orphaned?
	var refs struct {
		Count int64
	}
	err = h.joinTable.
		SELECT(jet.COUNT(h.fileCol).AS("count")).
		FROM(h.joinTable).
		WHERE(jet.AND(h.fileCol.EQ(jet.Uint64(fileID)))).
		QueryContext(ctx, tx, &refs)
	if err != nil {
		return err
	}

	if refs.Count == 0 {
		var key struct {
			FilePath string
		}
		tFiles := table.FivenetFiles
		err := tFiles.
			SELECT(tFiles.FilePath.AS("file_path")).
			WHERE(tFiles.ID.EQ(jet.Uint64(fileID))).
			QueryContext(ctx, tx, &key)
		if err != nil {
			return err
		}

		if err := h.store.Delete(ctx, key.FilePath); err != nil {
			return err
		}

		_, err = tFiles.
			UPDATE().
			SET(tFiles.DeletedAt.SET(jet.CURRENT_TIMESTAMP())).
			WHERE(tFiles.ID.EQ(jet.Uint64(fileID))).
			LIMIT(1).
			ExecContext(ctx, tx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteFileByPath deletes a file by its path and parent ID, if it exists.
func (h *Handler[P]) DeleteFileByPath(ctx context.Context, parentId P, path string) error {
	fileId, _, err := h.GetFileByPath(ctx, path)
	if err != nil {
		return err
	}

	if fileId == 0 {
		return nil
	}

	if err := h.Delete(ctx, parentId, fileId); err != nil {
		return err
	}

	return nil
}

// upsertFileRow inserts or updates a file row in the database, returning the file ID.
func upsertFileRow(ctx context.Context, tx *sql.Tx, key, ctype string, size int64) (uint64, error) {
	// 1. Try to lock an existing row via the UNIQUE(file_path) index
	var fileId struct {
		ID uint64 `jet:"id"`
	}
	err := tFiles.
		SELECT(tFiles.ID.AS("id")).
		WHERE(tFiles.FilePath.EQ(jet.String(key))).
		FOR(jet.UPDATE()). // ← row-lock
		QueryContext(ctx, tx, &fileId)

	switch {
	case errors.Is(err, qrm.ErrNoRows):
		// 2a) insert fresh row
		res, err := tFiles.
			INSERT(
				tFiles.FilePath,
				tFiles.ByteSize,
				tFiles.ContentType,
				tFiles.Meta,
			).
			VALUES(key, size, ctype, jet.NULL).
			ExecContext(ctx, tx)
		if err != nil {
			return 0, err
		}
		newID, _ := res.LastInsertId()
		return uint64(newID), nil

	case err == nil:
		// 2b) Row exists - overwrite metadata
		if _, err := tFiles.
			UPDATE().
			SET(
				tFiles.ByteSize.SET(jet.Int64(size)),
				tFiles.ContentType.SET(jet.String(ctype)),
				tFiles.DeletedAt.SET(jet.TimestampExp(jet.NULL)), // Revive if file was soft-deleted
			).
			WHERE(tFiles.ID.EQ(jet.Uint64(fileId.ID))).
			ExecContext(ctx, tx); err != nil {
			return 0, err
		}
		return fileId.ID, nil

	default:
		return 0, err // Unexpected query error
	}
}

// InsertJoinRow inserts a new join row for a parent and file.
func InsertJoinRow[P ParentID](
	ctx context.Context,
	tx *sql.Tx,
	join jet.Table,
	parentCol jet.Column,
	fileCol jet.ColumnInteger,
	parentId P,
	_ jet.BoolExpression,
	fileID uint64,
) error {
	_, err := join.
		INSERT(
			parentCol,
			fileCol,
		).
		VALUES(
			parentId,
			fileID,
		).
		ExecContext(ctx, tx)
	return err
}

// UpdateJoinRow updates the join row for a parent and file.
func UpdateJoinRow[P ParentID](
	ctx context.Context,
	tx *sql.Tx,
	join jet.Table,
	_ jet.Column,
	fileCol jet.ColumnInteger,
	_ P,
	parentIdBoolExp jet.BoolExpression,
	fileID uint64,
) error {
	_, err := join.
		UPDATE(
			fileCol,
		).
		SET(
			fileID,
		).
		WHERE(parentIdBoolExp).
		LIMIT(1).
		ExecContext(ctx, tx)
	return err
}

// putToStorage writes file data to the storage backend and returns the key and size.
func putToStorage(
	ctx context.Context,
	st storage.IStorage,
	key string,
	r io.Reader,
	ctype string,
	userSize int64,
) (string, int64, error) {
	cr := &countingReader{Reader: r}

	// Unknown total size: pass -1, backend switches to multipart automatically
	s := int64(-1)
	// If the user provided a size, use it instead
	if userSize > 0 {
		s = userSize
	}
	if _, err := st.Put(ctx, key, cr, s, ctype); err != nil {
		return "", 0, err
	}

	return key, cr.n, nil
}

// CountFilesForParentID returns the number of files associated with a given parent ID.
func (h *Handler[P]) CountFilesForParentID(ctx context.Context, parentID P) (int64, error) {
	stmt := h.joinTable.
		SELECT(jet.COUNT(h.fileCol).AS("count")).
		FROM(h.joinTable).
		WHERE(h.parentColBoolExp(parentID))

	var count struct {
		Count int64 `jet:"count"`
	}
	if err := stmt.QueryContext(ctx, h.db, &count); err != nil {
		return 0, err
	}

	return count.Count, nil
}

// ListFilesForParentID returns a list of files associated with a given parent ID.
func (h *Handler[P]) ListFilesForParentID(ctx context.Context, parentID P) ([]*file.File, error) {
	stmt := h.joinTable.
		SELECT(
			h.fileCol.AS("file.id"),
			h.parentCol.AS("file.parent_id"),
			tFiles.FilePath.AS("file.file_path"),
			tFiles.ByteSize.AS("file.byte_size"),
			tFiles.ContentType.AS("file.content_type"),
			tFiles.CreatedAt.AS("file.created_at"),
		).
		FROM(
			h.joinTable.
				INNER_JOIN(tFiles,
					h.fileCol.EQ(tFiles.ID),
				),
		).
		WHERE(h.parentColBoolExp(parentID)).
		LIMIT(15)

	var files []*file.File
	if err := stmt.QueryContext(ctx, h.db, &files); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return files, nil
}

// HandleFileChangesForParent synchronizes the join table for a parent with the provided list of files.
// It deletes associations for files no longer present and returns the number of added and deleted files.
func (h *Handler[P]) HandleFileChangesForParent(
	ctx context.Context,
	tx *sql.Tx,
	parentID P,
	updatedFiles []*file.File,
) (int64, int64, error) {
	current, err := h.ListFilesForParentID(ctx, parentID)
	if err != nil {
		return 0, 0, err
	}

	currentMap := make(map[uint64]*file.File)
	for _, f := range current {
		currentMap[f.GetId()] = f
	}

	updatedMap := make(map[uint64]*file.File)
	for _, f := range updatedFiles {
		updatedMap[f.GetId()] = f
	}

	var added []*file.File
	var deleted []*file.File

	for id, f := range updatedMap {
		if _, exists := currentMap[id]; !exists {
			added = append(added, f)
		}
	}

	for id, f := range currentMap {
		if _, exists := updatedMap[id]; !exists {
			deleted = append(deleted, f)
		}
	}

	if len(added) == 0 && len(deleted) == 0 {
		return 0, 0, nil // No changes
	}

	for _, f := range deleted {
		if err := h.deleteJoinRow(ctx, tx, parentID, f.GetId()); err != nil {
			return 0, 0, err
		}
	}

	// No need to add files as they are already present in the filestore and join table (mapping)

	return int64(len(added)), int64(len(deleted)), nil
}
