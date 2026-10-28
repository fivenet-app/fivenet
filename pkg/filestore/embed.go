package filestore

import (
	bytespkg "bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbfilestore "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/filestore"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ApplicationOctetStreamMIMEType = "application/octet-stream"

// MaxUploadChunks bounds the number of data packets accepted for one upload.
// The web client uses 128 KiB packets, so this leaves substantial room for
// small-chunk clients while bounding per-message overhead.
const MaxUploadChunks = 4096

var tFiles = table.FivenetFiles

// ParentID is a type constraint that matches any type that can be a parent ID.
// It can be an unsigned 64-bit integer, a 32-bit integer, or a string.
type ParentID interface {
	~int64 | ~int32 | ~string
	comparable
}

// ParentColBoolExpFn is a function type that returns a Jet boolean expression for a given parent ID.
type ParentColBoolExpFn[P ParentID] func(parentId P) mysql.BoolExpression

// JoinRowInserterFn is a function type for inserting a join row into a table linking parent and file.
type JoinRowInserterFn[P ParentID] func(ctx context.Context, tx *sql.Tx, join mysql.Table, parentCol mysql.Column, fileCol mysql.ColumnInteger, parentId P, _ mysql.BoolExpression, fileID int64) error

// JoinRowReplacerFn replaces a singleton relation and returns the old file ID.
type JoinRowReplacerFn[P ParentID] func(ctx context.Context, tx *sql.Tx, join mysql.Table, parentCol mysql.Column, fileCol mysql.ColumnInteger, parentID P, parentIDBoolExp mysql.BoolExpression, fileID int64) (int64, error)

// FileLimitLockFn locks the parent row used to serialize a file-count check
// with the subsequent join-row insert.
type FileLimitLockFn[P ParentID] func(ctx context.Context, tx *sql.Tx, parentID P) error

// RelationFileIDFn finds the relation that should be deleted while the
// handler's parent lock is held.
type RelationFileIDFn[P ParentID] func(ctx context.Context, tx *sql.Tx, join mysql.Table, parentCol mysql.Column, fileCol mysql.ColumnInteger, parentID P) (int64, error)

// UploadCallbackFn runs after the file and its relation have been written, but
// before their transaction commits. It allows callers to atomically persist
// related records, such as an activity entry.
type UploadCallbackFn func(ctx context.Context, tx *sql.Tx, response *file.UploadFileResponse) error

// DeleteCallbackFn runs before the file relation is removed, while the
// deletion transaction is still open. It allows callers to atomically persist
// related records, such as an activity entry.
//
// The callback intentionally runs before the relation deletion because an
// orphaned file may be removed from storage during deletion. A callback error
// must therefore abort before storage is touched.
type DeleteCallbackFn func(ctx context.Context, tx *sql.Tx) error

// Handler is a generic, embeddable file-upload helper.
type Handler[P ParentID] struct {
	// store is the storage adapter (e.g. S3, filesystem).
	store storage.IStorage
	// db is the SQL database connection.
	db *sql.DB
	// joinTable is the table that joins parent and file.
	joinTable mysql.Table
	// parentCol is the column in joinTable for the parent (may be integer or string).
	parentCol mysql.Column
	// fileCol is the column in joinTable for the file ID.
	fileCol mysql.ColumnInteger
	// sizeLimit is the maximum allowed file size in bytes.
	sizeLimit int64
	// fileLimit is the maximum allowed number of files per parent.
	fileLimit int64
	// uploadFilter optionally restricts allowed upload file types.
	uploadFilter  *UploadFilter
	fileLimitLock FileLimitLockFn[P]
	// fileLimitMu serializes additive uploads so the count check and insert are
	// atomic from the handler's perspective.
	fileLimitMu sync.Mutex

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
	join mysql.Table,
	parentCol mysql.Column,
	fileCol mysql.ColumnInteger,
	sizeLimit int64,
	fileLimit int64,
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
		fileLimit: fileLimit,
		// parentColBoolExp is a function that converts the parent ID to a mysql.BoolExpression
		parentColBoolExp:  parentColBoolExp,
		joinRowInserter:   joinRowInserter,
		nullOnlyParentRow: nullOnlyParentRow,
	}
}

// WithUploadFilter sets an optional upload filter for file type validation.
func (h *Handler[P]) WithUploadFilter(filter *UploadFilter) *Handler[P] {
	h.uploadFilter = filter.clone()
	return h
}

// WithFileLimitLock configures a database row lock for the file-limit check.
func (h *Handler[P]) WithFileLimitLock(lock FileLimitLockFn[P]) *Handler[P] {
	h.fileLimitLock = lock
	return h
}

// CheckFileLimit enforces the maximum number of files per parent for additive relations.
func (h *Handler[P]) CheckFileLimit(ctx context.Context, parentID P) error {
	if h.joinTable == nil || h.fileLimit <= 0 || h.nullOnlyParentRow {
		return nil
	}

	count, err := h.CountFilesForParentID(ctx, parentID)
	if err != nil {
		return err
	}

	if count >= h.fileLimit {
		return ErrMaxFilesPerParentReached(map[string]any{
			"maxFiles": h.fileLimit,
		})
	}

	return nil
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
	resp, _, err := h.uploadFile(ctx, parentID, key, size, ctype, srv, nil, nil)
	return resp, err
}

// UploadFileWithCallback is UploadFile with a callback that participates in
// the file metadata and relation transaction.
func (h *Handler[P]) UploadFileWithCallback(
	ctx context.Context,
	parentID P,
	key string,
	size int64,
	ctype string,
	srv pbfilestore.FilestoreService_UploadServer,
	callback UploadCallbackFn,
) (*file.UploadFileResponse, error) {
	resp, _, err := h.uploadFile(ctx, parentID, key, size, ctype, srv, nil, callback)
	return resp, err
}

func (h *Handler[P]) uploadFile(
	ctx context.Context,
	parentID P,
	key string,
	size int64,
	ctype string,
	srv pbfilestore.FilestoreService_UploadServer,
	replacer JoinRowReplacerFn[P],
	callback UploadCallbackFn,
) (*file.UploadFileResponse, int64, error) {
	if h.sizeLimit > 0 && size > h.sizeLimit {
		return nil, 0, ErrUploadFileTooLarge(
			map[string]any{"maxSize": h.sizeLimit / (1 << 20)},
		) // Convert bytes to megabytes
	}

	ctype = sniff(ctype, key)
	chunkCounter := &uploadChunkCounter{max: MaxUploadChunks}
	leadingData, err := readLeadingDataWithCounter(srv, 512, chunkCounter)
	if err != nil {
		return nil, 0, err
	}
	if detected := detectContentType(leadingData); detected != "" {
		ctype = detected
	}

	if err := h.validateUploadType(key, ctype); err != nil {
		return nil, 0, err
	}

	if h.joinTable != nil && h.fileLimit > 0 && !h.nullOnlyParentRow && h.fileLimitLock == nil {
		h.fileLimitMu.Lock()
		defer h.fileLimitMu.Unlock()
	}
	if replacer == nil {
		if err := h.CheckFileLimit(ctx, parentID); err != nil {
			return nil, 0, err
		}
	}

	// Pipe chunks to the storage backend
	pr := streamUploadData(srv, leadingData, chunkCounter)
	defer pr.Close()

	// Size - if client sent one, honor it; otherwise -1 (unknown)
	szHint := int64(-1)
	if size > 0 {
		szHint = size
	}

	cleanupStorage := true
	defer func() {
		if cleanupStorage {
			_ = h.store.Delete(ctx, key)
		}
	}()

	urlKey, bytes, err := putToStorage(ctx, h.store, key, pr, ctype, szHint, h.sizeLimit)
	if err != nil {
		return nil, 0, err
	}

	// DB transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	var replacedFileID int64
	if h.joinTable != nil && h.fileLimit > 0 && !h.nullOnlyParentRow {
		if h.fileLimitLock != nil {
			if err := h.fileLimitLock(ctx, tx, parentID); err != nil {
				return nil, 0, err
			}
		}
		if replacer == nil {
			count, err := h.countFilesForParentID(ctx, tx, parentID)
			if err != nil {
				return nil, 0, err
			}
			if count >= h.fileLimit {
				return nil, 0, ErrMaxFilesPerParentReached(map[string]any{"maxFiles": h.fileLimit})
			}
		}
	}

	fileID, err := upsertFileRow(ctx, tx, key, ctype, bytes)
	if err != nil {
		return nil, 0, err
	}

	if h.joinTable != nil {
		parentIdBoolExp := h.parentColBoolExp(parentID)
		var err error
		if replacer != nil {
			replacedFileID, err = replacer(
				ctx, tx, h.joinTable, h.parentCol, h.fileCol,
				parentID, parentIdBoolExp, fileID,
			)
		} else {
			err = h.joinRowInserter(
				ctx,
				tx,
				h.joinTable,
				h.parentCol,
				h.fileCol,
				parentID,
				parentIdBoolExp,
				fileID,
			)
		}
		if err != nil {
			return nil, 0, err
		}
	}
	if replacer != nil && h.joinTable != nil && h.fileLimit > 0 && !h.nullOnlyParentRow {
		count, err := h.countFilesForParentID(ctx, tx, parentID)
		if err != nil {
			return nil, 0, err
		}
		if count > h.fileLimit {
			return nil, 0, ErrMaxFilesPerParentReached(map[string]any{"maxFiles": h.fileLimit})
		}
	}

	full, err := url.JoinPath(FilestoreURLPrefix, urlKey)
	if err != nil {
		return nil, 0, err
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
	if callback != nil {
		if err := callback(ctx, tx, resp); err != nil {
			return nil, 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	cleanupStorage = false
	if err := srv.SendAndClose(resp); err != nil {
		return resp, replacedFileID, err
	}
	grpc_audit.AddFileMeta(ctx, resp.GetFile())
	return resp, replacedFileID, nil
}

// ReplaceFile replaces the bytes and metadata of an existing file while keeping
// its storage key and database ID stable. This is useful for files referenced by
// documents through their URL.
func (h *Handler[P]) ReplaceFile(
	ctx context.Context,
	fileID int64,
	key string,
	size int64,
	ctype string,
	srv pbfilestore.FilestoreService_UploadServer,
) (*file.UploadFileResponse, error) {
	if h.sizeLimit > 0 && size > h.sizeLimit {
		return nil, ErrUploadFileTooLarge(map[string]any{"maxSize": h.sizeLimit / (1 << 20)})
	}

	ctype = sniff(ctype, key)
	chunkCounter := &uploadChunkCounter{max: MaxUploadChunks}
	leadingData, err := readLeadingDataWithCounter(srv, 512, chunkCounter)
	if err != nil {
		return nil, err
	}
	if detected := detectContentType(leadingData); detected != "" {
		ctype = detected
	}
	if err := h.validateUploadType(key, ctype); err != nil {
		return nil, err
	}

	// Keep a rollback copy because storage replacement and the metadata update
	// cannot share a transaction. This preserves the old URL contents if the
	// metadata update fails.
	oldObject, oldInfo, err := h.store.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	oldReader := io.Reader(oldObject)
	if h.sizeLimit > 0 {
		oldReader = io.LimitReader(oldObject, h.sizeLimit+1)
	}
	oldBytes, err := io.ReadAll(oldReader)
	_ = oldObject.Close()
	if err != nil {
		return nil, err
	}
	if h.sizeLimit > 0 && int64(len(oldBytes)) > h.sizeLimit {
		return nil, ErrUploadFileTooLarge(map[string]any{"maxSize": h.sizeLimit / (1 << 20)})
	}

	restoreOldFile := func(cause error) (*file.UploadFileResponse, error) {
		_, restoreErr := h.store.Put(
			ctx,
			key,
			bytespkg.NewReader(oldBytes),
			int64(len(oldBytes)),
			oldInfo.GetContentType(),
		)
		if restoreErr != nil {
			return nil, errors.Join(cause, restoreErr)
		}
		return nil, cause
	}

	pr := streamUploadData(srv, leadingData, chunkCounter)
	defer pr.Close()

	_, bytes, err := putToStorage(ctx, h.store, key, pr, ctype, size, h.sizeLimit)
	if err != nil {
		return restoreOldFile(err)
	}
	result, err := tFiles.UPDATE().SET(
		tFiles.ByteSize.SET(mysql.Int64(bytes)),
		tFiles.ContentType.SET(mysql.String(ctype)),
		tFiles.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
	).WHERE(tFiles.ID.EQ(mysql.Int64(fileID))).LIMIT(1).ExecContext(ctx, h.db)
	if err != nil {
		return restoreOldFile(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return restoreOldFile(err)
	}
	if rows == 0 {
		var existing struct {
			ID int64 `alias:"id"`
		}
		err := tFiles.
			SELECT(tFiles.ID.AS("id")).
			WHERE(tFiles.ID.EQ(mysql.Int64(fileID))).
			LIMIT(1).
			QueryContext(ctx, h.db, &existing)
		if errors.Is(err, qrm.ErrNoRows) {
			return restoreOldFile(sql.ErrNoRows)
		}
		if err != nil {
			return restoreOldFile(err)
		}
	}

	full, err := url.JoinPath(FilestoreURLPrefix, key)
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
	if err := srv.SendAndClose(resp); err != nil {
		return resp, err
	}
	grpc_audit.AddFileMeta(ctx, resp.GetFile())
	grpc_audit.AddPreviousFileMeta(ctx, fileID, key)
	return resp, nil
}

func (h *Handler[P]) validateUploadType(key, contentType string) error {
	if h.uploadFilter == nil {
		return nil
	}
	return h.uploadFilter.Validate(key, contentType)
}

func readLeadingData(
	srv pbfilestore.FilestoreService_UploadServer,
	minBytes int,
) ([]byte, error) {
	return readLeadingDataWithCounter(srv, minBytes, nil)
}

type uploadChunkCounter struct {
	count int
	max   int
}

func streamUploadData(
	srv pbfilestore.FilestoreService_UploadServer,
	leadingData []byte,
	counter *uploadChunkCounter,
) *io.PipeReader {
	pr, pw := io.Pipe()
	go func() {
		if len(leadingData) > 0 {
			if _, err := pw.Write(leadingData); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
		for {
			pkt, err := recvUploadChunk(srv, counter)
			if errors.Is(err, io.EOF) {
				_ = pw.Close()
				return
			}
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			if data := pkt.GetData(); data != nil {
				if _, err := pw.Write(data); err != nil {
					pw.CloseWithError(err)
					return
				}
			}
		}
	}()
	return pr
}

func recvUploadChunk(
	srv pbfilestore.FilestoreService_UploadServer,
	counter *uploadChunkCounter,
) (*file.UploadFileRequest, error) {
	if counter != nil && counter.count >= counter.max {
		// Permit the terminating EOF after the final allowed packet, but reject
		// any additional packet.
		_, err := srv.Recv()
		if errors.Is(err, io.EOF) {
			return nil, io.EOF
		}
		if err != nil {
			return nil, err
		}
		return nil, status.Errorf(
			codes.ResourceExhausted,
			"upload contains too many chunks: maximum is %d",
			counter.max,
		)
	}
	pkt, err := srv.Recv()
	if err != nil {
		return nil, err
	}
	if counter != nil {
		counter.count++
	}
	return pkt, nil
}

func readLeadingDataWithCounter(
	srv pbfilestore.FilestoreService_UploadServer,
	minBytes int,
	counter *uploadChunkCounter,
) ([]byte, error) {
	if minBytes <= 0 {
		return nil, nil
	}

	var buf bytespkg.Buffer
	for buf.Len() < minBytes {
		pkt, err := recvUploadChunk(srv, counter)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		data := pkt.GetData()
		if len(data) == 0 {
			continue
		}

		if _, err := buf.Write(data); err != nil {
			return nil, err
		}
	}

	if buf.Len() == 0 {
		return nil, status.Error(codes.InvalidArgument, "upload stream has no file data")
	}

	return buf.Bytes(), nil
}

func detectContentType(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	detected := normalizeContentType(http.DetectContentType(data))
	if detected == ApplicationOctetStreamMIMEType {
		return ""
	}

	return detected
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

// UploadFromMetaWithCallback is UploadFromMeta with a callback that
// participates in the file metadata and relation transaction.
func (h *Handler[P]) UploadFromMetaWithCallback(
	ctx context.Context,
	meta *file.UploadMeta,
	parentID P,
	srv pbfilestore.FilestoreService_UploadServer,
	callback UploadCallbackFn,
) (*file.UploadFileResponse, error) {
	key := buildKey(meta.GetNamespace(), SanitizeFileName(meta.GetOriginalName()))
	ctype := sniff(meta.GetContentType(), key)

	return h.UploadFileWithCallback(ctx, parentID, key, meta.GetSize(), ctype, srv, callback)
}

// UploadFromMetaWithReplacement is the replacement variant for singleton
// relations. The relation callback runs while the handler holds its parent
// lock and the upload transaction is open.
func (h *Handler[P]) UploadFromMetaWithReplacement(
	ctx context.Context,
	meta *file.UploadMeta,
	parentID P,
	srv pbfilestore.FilestoreService_UploadServer,
	replacer JoinRowReplacerFn[P],
) (*file.UploadFileResponse, int64, error) {
	key := buildKey(meta.GetNamespace(), SanitizeFileName(meta.GetOriginalName()))
	ctype := sniff(meta.GetContentType(), key)
	return h.uploadFile(ctx, parentID, key, meta.GetSize(), ctype, srv, replacer, nil)
}

// GetFileByPath retrieves the file ID and path for a given file path string.
func (h *Handler[P]) GetFileByPath(ctx context.Context, path string) (int64, string, error) {
	path = strings.TrimPrefix(path, FilestoreURLPrefix)

	tFiles := table.FivenetFiles.AS("file")

	stmt := tFiles.
		SELECT(
			tFiles.ID,
			tFiles.FilePath,
		).
		FROM(tFiles).
		WHERE(
			tFiles.FilePath.EQ(mysql.String(path)),
		).
		LIMIT(1)

	f := &struct {
		ID       int64  `alias:"id"`
		ParentID P      `alias:"parent_id"`
		FilePath string `alias:"file_path"`
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
	fileID int64,
) error {
	if h.nullOnlyParentRow {
		_, err := h.joinTable.
			UPDATE(
				h.parentCol,
				h.fileCol,
			).
			SET(
				parentID,
				mysql.NULL,
			).
			WHERE(mysql.AND(
				h.parentColBoolExp(parentID),
				h.fileCol.EQ(mysql.Int64(fileID)),
			)).
			LIMIT(1).
			ExecContext(ctx, tx)
		return err
	}

	_, err := h.joinTable.
		DELETE().
		WHERE(mysql.AND(
			h.parentColBoolExp(parentID),
			h.fileCol.EQ(mysql.Int64(fileID)),
		)).
		LIMIT(1).
		ExecContext(ctx, tx)
	return err
}

// Delete removes the association between a parent and a file, and deletes the file if it is orphaned.
func (h *Handler[P]) Delete(ctx context.Context, parentID P, fileID int64) error {
	return h.DeleteWithCallback(ctx, parentID, fileID, nil)
}

// DeleteWithCallback removes the association between a parent and a file, and
// deletes the file if it is orphaned. callback participates in the deletion
// transaction and runs before the relation or storage object is removed.
func (h *Handler[P]) DeleteWithCallback(
	ctx context.Context,
	parentID P,
	fileID int64,
	callback DeleteCallbackFn,
) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if callback != nil {
		if err := callback(ctx, tx); err != nil {
			return err
		}
	}
	if err := h.deleteInTx(ctx, tx, parentID, fileID); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteWithParentLock finds and deletes a relation while holding the
// caller-provided parent lock. This prevents replacement and deletion from
// racing on singleton relations.
func (h *Handler[P]) DeleteWithParentLock(
	ctx context.Context,
	parentID P,
	lock FileLimitLockFn[P],
	find RelationFileIDFn[P],
) (int64, error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if err := lock(ctx, tx, parentID); err != nil {
		return 0, err
	}
	fileID, err := find(ctx, tx, h.joinTable, h.parentCol, h.fileCol, parentID)
	if err != nil {
		return 0, err
	}
	if fileID == 0 {
		return 0, tx.Commit()
	}
	if err := h.deleteInTx(ctx, tx, parentID, fileID); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return fileID, nil
}

func (h *Handler[P]) deleteInTx(ctx context.Context, tx *sql.Tx, parentID P, fileID int64) error {
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
	err := h.joinTable.
		SELECT(mysql.COUNT(h.fileCol).AS("count")).
		FROM(h.joinTable).
		WHERE(h.fileCol.EQ(mysql.Int64(fileID))).
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
			WHERE(tFiles.ID.EQ(mysql.Int64(fileID))).
			QueryContext(ctx, tx, &key)
		if err != nil {
			return err
		}

		if err := h.store.Delete(ctx, key.FilePath); err != nil {
			return err
		}

		_, err = tFiles.
			UPDATE().
			SET(tFiles.DeletedAt.SET(mysql.CURRENT_TIMESTAMP())).
			WHERE(tFiles.ID.EQ(mysql.Int64(fileID))).
			LIMIT(1).
			ExecContext(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
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
func upsertFileRow(ctx context.Context, tx *sql.Tx, key, ctype string, size int64) (int64, error) {
	// 1. Try to lock an existing row via the UNIQUE(file_path) index
	var fileId struct {
		ID int64 `alias:"id"`
	}
	err := tFiles.
		SELECT(tFiles.ID.AS("id")).
		WHERE(tFiles.FilePath.EQ(mysql.String(key))).
		FOR(mysql.UPDATE()). // Row-lock
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
			VALUES(key, size, ctype, mysql.NULL).
			ExecContext(ctx, tx)
		if err != nil {
			return 0, err
		}
		newID, _ := res.LastInsertId()
		return newID, nil

	case err == nil:
		// 2b) Row exists - overwrite metadata
		if _, err := tFiles.
			UPDATE().
			SET(
				tFiles.ByteSize.SET(mysql.Int64(size)),
				tFiles.ContentType.SET(mysql.String(ctype)),
				tFiles.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)), // Revive if file was soft-deleted
			).
			WHERE(tFiles.ID.EQ(mysql.Int64(fileId.ID))).
			LIMIT(1).
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
	join mysql.Table,
	parentCol mysql.Column,
	fileCol mysql.ColumnInteger,
	parentId P,
	_ mysql.BoolExpression,
	fileID int64,
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
	join mysql.Table,
	_ mysql.Column,
	fileCol mysql.ColumnInteger,
	_ P,
	parentIdBoolExp mysql.BoolExpression,
	fileID int64,
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
	maxSize int64,
) (string, int64, error) {
	cr := &countingReader{Reader: r, max: maxSize}

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
	return h.countFilesForParentID(ctx, h.db, parentID)
}

func (h *Handler[P]) countFilesForParentID(
	ctx context.Context,
	db qrm.Queryable,
	parentID P,
) (int64, error) {
	stmt := h.joinTable.
		SELECT(mysql.COUNT(h.fileCol).AS("data_count.total")).
		FROM(h.joinTable).
		WHERE(h.parentColBoolExp(parentID))

	var count database.DataCount
	if err := stmt.QueryContext(ctx, db, &count); err != nil {
		return 0, err
	}

	return count.Total, nil
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
		LIMIT(20)

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

	currentMap := make(map[int64]*file.File)
	for _, f := range current {
		currentMap[f.GetId()] = f
	}

	updatedMap := make(map[int64]*file.File)
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
