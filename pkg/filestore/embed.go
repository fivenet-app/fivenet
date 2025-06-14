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

type ParentID interface {
	~uint64 | ~int32 | ~string
	comparable
}

type ParentColBoolExpFn[P ParentID] func(parentId P) jet.BoolExpression

type JoinRowInserterFn[P ParentID] func(ctx context.Context, tx *sql.Tx, join jet.Table, parentCol jet.Column, fileCol jet.ColumnInteger, parentId P, _ jet.BoolExpression, fileID uint64) error

// Generic, embeddable file‑upload helper.
type Handler[P ParentID] struct {
	store     storage.IStorage // Storage adapter (e.g. S3, FS)
	db        *sql.DB
	joinTable jet.Table
	parentCol jet.Column // May be ColumnInteger or ColumnString
	fileCol   jet.ColumnInteger
	sizeLimit int64

	parentColBoolExp  ParentColBoolExpFn[P]
	joinRowInserter   JoinRowInserterFn[P]
	nullOnlyParentRow bool
}

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

// AwaitHandshake is a no-op for the filestore handler, as it does not require any
// special handshake logic. It simply reads the first packet to extract metadata
// for the upload.
func (h *Handler[P]) AwaitHandshake(srv pbfilestore.FilestoreService_UploadServer) (*file.UploadMeta, error) {
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
		return nil, status.Errorf(codes.ResourceExhausted, "file size exceeds limit: %d > %d", meta.GetSize(), h.sizeLimit)
	}

	return meta, nil
}

func (h *Handler[P]) UploadFile(ctx context.Context, parentID P, key string, size int64, ctype string, srv pbfilestore.FilestoreService_UploadServer) (*file.UploadResponse, error) {
	if h.sizeLimit > 0 && size > h.sizeLimit {
		return nil, status.Errorf(codes.ResourceExhausted, "file too large: %d > %d", size, h.sizeLimit)
	}

	// pipe chunks to the storage backend
	pr, pw := io.Pipe()
	go func() {
		for {
			pkt, err := srv.Recv()
			if err == io.EOF {
				pw.Close()
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

	// Size – if client sent one, honour it; otherwise -1 (unknown)
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

	resp := &file.UploadResponse{
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

// UploadFromMeta streams the remainder of the gRPC Upload after the caller
// has read & validated the first UploadMeta packet.
func (h *Handler[P]) UploadFromMeta(
	ctx context.Context,
	meta *file.UploadMeta,
	parentID P,
	srv pbfilestore.FilestoreService_UploadServer,
) (*file.UploadResponse, error) {
	key := buildKey(meta.GetNamespace(), SanitizeFileName(meta.GetOriginalName()))
	ctype := sniff(meta.GetContentType(), key)

	return h.UploadFile(ctx, parentID, key, meta.GetSize(), ctype, srv)
}

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

func (h *Handler[P]) deleteJoinRow(ctx context.Context, tx *sql.Tx, parentID P, fileID uint64) error {
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

// Delete (unary)
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

// Helpers

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

	switch err {
	case qrm.ErrNoRows:
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

	case nil:
		// 2b) Row exists – overwrite metadata
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

func InsertJoinRow[P ParentID](ctx context.Context, tx *sql.Tx, join jet.Table, parentCol jet.Column, fileCol jet.ColumnInteger, parentId P, _ jet.BoolExpression, fileID uint64) error {
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

func UpdateJoinRow[P ParentID](ctx context.Context, tx *sql.Tx, join jet.Table, parentCol jet.Column, fileCol jet.ColumnInteger, parentId P, parentIdBoolExp jet.BoolExpression, fileID uint64) error {
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

func putToStorage(ctx context.Context, st storage.IStorage, key string, r io.Reader, ctype string, userSize int64) (url string, size int64, err error) {
	cr := &countingReader{Reader: r}

	// Unknown total size: pass -1, backend switches to multipart automatically
	s := int64(-1)
	// If the user provided a size, use it instead
	if userSize > 0 {
		s = userSize
	}
	if _, err = st.Put(ctx, key, cr, s, ctype); err != nil {
		return "", 0, err
	}

	return key, cr.n, nil
}

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

func (h *Handler[P]) HandleFileChangesForParent(ctx context.Context, tx *sql.Tx, parentID P, updatedFiles []*file.File) (int64, int64, error) {
	current, err := h.ListFilesForParentID(ctx, parentID)
	if err != nil {
		return 0, 0, err
	}

	currentMap := make(map[uint64]*file.File)
	for _, f := range current {
		currentMap[f.Id] = f
	}

	updatedMap := make(map[uint64]*file.File)
	for _, f := range updatedFiles {
		updatedMap[f.Id] = f
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
		if err := h.deleteJoinRow(ctx, tx, parentID, f.Id); err != nil {
			return 0, 0, err
		}
	}

	// No need to add files as they are already present in the filestore and join table (mapping)

	return int64(len(added)), int64(len(deleted)), nil
}
