package filestore

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAwaitHandshake(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		packets  []*file.UploadFileRequest
		recvErr  error
		limit    int64
		wantCode int
		wantMeta bool
	}{
		{
			name:     "valid metadata",
			packets:  []*file.UploadFileRequest{uploadMetaRequest(5)},
			wantMeta: true,
		},
		{
			name:     "missing metadata",
			recvErr:  io.EOF,
			wantCode: 3, // InvalidArgument
		},
		{
			name:     "first packet is data",
			packets:  []*file.UploadFileRequest{{}},
			wantCode: 3,
		},
		{
			name:     "metadata exceeds limit",
			packets:  []*file.UploadFileRequest{uploadMetaRequest(11)},
			limit:    10,
			wantCode: 8, // ResourceExhausted
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler[int64]{sizeLimit: tt.limit}
			stream := &uploadStream{packets: tt.packets, recvErr: tt.recvErr}
			meta, err := h.AwaitHandshake(stream)
			if tt.wantMeta {
				require.NoError(t, err)
				assert.Equal(t, int64(5), meta.GetSize())
				return
			}
			require.Error(t, err)
			assert.Equal(t, tt.wantCode, int(status.Code(err)))
		})
	}
}

func TestReadLeadingData(t *testing.T) {
	t.Parallel()

	stream := &uploadStream{
		packets: []*file.UploadFileRequest{
			dataRequest(nil),
			dataRequest([]byte("abc")),
			dataRequest([]byte("def")),
		},
		recvErr: io.EOF,
	}
	data, err := readLeadingData(stream, 5)
	require.NoError(t, err)
	assert.Equal(t, []byte("abcdef"), data)

	empty := &uploadStream{recvErr: io.EOF}
	_, err = readLeadingData(empty, 5)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestUploadFileSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("hello"))},
		recvErr: io.EOF,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*file_path = \?.*FOR UPDATE;`).
		WithArgs("uploads/file.txt").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(`INSERT INTO fivenet_files .*file_path.*byte_size.*content_type.*meta`).
		WithArgs("uploads/file.txt", int64(5), "text/plain").
		WillReturnResult(sqlmock.NewResult(23, 1))
	mock.ExpectCommit()

	h := &Handler[int64]{store: st, db: db}
	resp, err := h.UploadFile(t.Context(), 7, "uploads/file.txt", 5, "", stream)
	require.NoError(t, err)
	assert.Equal(t, int64(23), resp.GetId())
	assert.Equal(t, "/api/filestore/uploads/file.txt", resp.GetUrl())
	assert.Equal(t, "text/plain", resp.GetFile().GetContentType())
	assert.Equal(t, int64(5), resp.GetFile().GetByteSize())
	assert.Equal(t, "hello", st.putData)
	assert.Empty(t, st.deleted)
}

func TestUploadFileRejectsDisallowedTypeBeforeStorage(t *testing.T) {
	t.Parallel()
	st := &uploadStorage{}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("plain text"))},
		recvErr: io.EOF,
	}
	h := (&Handler[int64]{store: st}).WithUploadFilter(NewImageUploadFilter())

	_, err := h.UploadFile(t.Context(), 1, "uploads/file.txt", 10, "", stream)
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Empty(t, st.putData)
}

func TestUploadFileCleansStorageWhenDatabaseBeginFails(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("hello"))},
		recvErr: io.EOF,
	}
	mock.ExpectBegin().WillReturnError(errors.New("database unavailable"))

	h := &Handler[int64]{store: st, db: db}
	_, err := h.UploadFile(t.Context(), 1, "uploads/file.txt", 5, "text/plain", stream)
	require.Error(t, err)
	assert.Equal(t, "uploads/file.txt", st.deleted)
}

func TestUploadFileCleansStorageWhenStreamReceiveFails(t *testing.T) {
	t.Parallel()
	st := &uploadStorage{}
	streamErr := errors.New("stream receive failed")
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest(bytes.Repeat([]byte("x"), 512))},
		recvErr: streamErr,
	}

	h := &Handler[int64]{store: st}
	_, err := h.UploadFile(t.Context(), 1, "uploads/file.txt", 512, "text/plain", stream)
	require.ErrorIs(t, err, streamErr)
	assert.Equal(t, "uploads/file.txt", st.deleted)
}

func TestUploadFileCleansStorageWhenCommitFails(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("hello"))},
		recvErr: io.EOF,
	}
	commitErr := errors.New("commit failed")

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*file_path = \?.*FOR UPDATE;`).
		WithArgs("uploads/file.txt").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(`INSERT INTO fivenet_files .*file_path.*byte_size.*content_type.*meta`).
		WithArgs("uploads/file.txt", int64(5), "text/plain").
		WillReturnResult(sqlmock.NewResult(23, 1))
	mock.ExpectCommit().WillReturnError(commitErr)

	h := &Handler[int64]{store: st, db: db}
	_, err := h.UploadFile(t.Context(), 1, "uploads/file.txt", 5, "text/plain", stream)
	require.ErrorIs(t, err, commitErr)
	assert.Equal(t, "uploads/file.txt", st.deleted)
}

func TestReplaceFileSuccess(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{getData: "old content", getType: "text/plain"}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("new content"))},
		recvErr: io.EOF,
	}

	mock.ExpectExec(`UPDATE fivenet_files SET .*byte_size.*content_type.*deleted_at.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(11), "text/plain", int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	h := &Handler[int64]{store: st, db: db}
	resp, err := h.ReplaceFile(t.Context(), 23, "uploads/file.txt", 11, "", stream)
	require.NoError(t, err)
	assert.Equal(t, int64(23), resp.GetId())
	assert.Equal(t, "/api/filestore/uploads/file.txt", resp.GetUrl())
	assert.Equal(t, "new content", st.putData)
	assert.Equal(t, []string{"new content"}, st.putHistory)
	assert.True(t, stream.closed)
}

func TestReplaceFileRestoresOldObjectWhenMetadataUpdateFails(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{getData: "old content", getType: "image/png"}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("new content"))},
		recvErr: io.EOF,
	}
	databaseErr := errors.New("metadata update failed")
	mock.ExpectExec(`UPDATE fivenet_files SET .*byte_size.*content_type.*deleted_at.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(11), "text/plain", int64(23), int64(1)).
		WillReturnError(databaseErr)

	h := &Handler[int64]{store: st, db: db}
	_, err := h.ReplaceFile(t.Context(), 23, "uploads/file.txt", 11, "", stream)
	require.ErrorIs(t, err, databaseErr)
	assert.Equal(t, []string{"new content", "old content"}, st.putHistory)
	assert.Equal(t, "image/png", st.putTypes[1])
}

func TestReplaceFileMissingIDRestoresOldObject(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{getData: "old content", getType: "text/plain"}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("new content"))},
		recvErr: io.EOF,
	}
	mock.ExpectExec(`UPDATE fivenet_files SET .*byte_size.*content_type.*deleted_at.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(11), "text/plain", int64(23), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(23), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	h := &Handler[int64]{store: st, db: db}
	_, err := h.ReplaceFile(t.Context(), 23, "uploads/file.txt", 11, "", stream)
	require.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, []string{"new content", "old content"}, st.putHistory)
}

func TestReplaceFileReturnsStorageReadError(t *testing.T) {
	t.Parallel()
	st := &uploadStorage{getErr: errors.New("storage read failed")}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("new content"))},
		recvErr: io.EOF,
	}

	h := &Handler[int64]{store: st}
	_, err := h.ReplaceFile(t.Context(), 23, "uploads/file.txt", 11, "text/plain", stream)
	require.ErrorIs(t, err, st.getErr)
	assert.Empty(t, st.putHistory)
}

func TestReplaceFileRejectsSizeLimitBeforeStorageRead(t *testing.T) {
	t.Parallel()
	st := &uploadStorage{}
	h := &Handler[int64]{store: st, sizeLimit: 10}

	_, err := h.ReplaceFile(t.Context(), 23, "uploads/file.txt", 11, "text/plain", &uploadStream{})
	require.Error(t, err)
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))
	assert.Empty(t, st.putHistory)
}

func TestReplaceFileReturnsRestorationError(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{
		getData: "old content",
		getType: "text/plain",
		putErrs: []error{nil, errors.New("restore failed")},
	}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("new content"))},
		recvErr: io.EOF,
	}
	databaseErr := errors.New("metadata update failed")
	mock.ExpectExec(`UPDATE fivenet_files SET .*byte_size.*content_type.*deleted_at.*WHERE .*id = \?.*LIMIT \?;`).
		WithArgs(int64(11), "text/plain", int64(23), int64(1)).
		WillReturnError(databaseErr)

	h := &Handler[int64]{store: st, db: db}
	_, err := h.ReplaceFile(t.Context(), 23, "uploads/file.txt", 11, "", stream)
	require.Error(t, err)
	require.ErrorIs(t, err, databaseErr)
	require.ErrorContains(t, err, "restore failed")
}

func uploadMetaRequest(size int64) *file.UploadFileRequest {
	request := &file.UploadFileRequest{}
	request.SetMeta(&file.UploadMeta{Size: size, OriginalName: "file.txt"})
	return request
}

func dataRequest(data []byte) *file.UploadFileRequest {
	request := &file.UploadFileRequest{}
	request.SetData(data)
	return request
}

type uploadStream struct {
	packets []*file.UploadFileRequest
	recvErr error
	sent    *file.UploadFileResponse
	closed  bool
}

func (s *uploadStream) Recv() (*file.UploadFileRequest, error) {
	if len(s.packets) == 0 {
		if s.recvErr != nil {
			return nil, s.recvErr
		}
		return nil, io.EOF
	}
	packet := s.packets[0]
	s.packets = s.packets[1:]
	return packet, nil
}

func (s *uploadStream) SendAndClose(resp *file.UploadFileResponse) error {
	s.sent = resp
	s.closed = true
	return nil
}

func (s *uploadStream) Context() context.Context     { return context.Background() }
func (s *uploadStream) SetHeader(metadata.MD) error  { return nil }
func (s *uploadStream) SendHeader(metadata.MD) error { return nil }
func (s *uploadStream) SetTrailer(metadata.MD)       {}
func (s *uploadStream) RecvMsg(any) error            { return errors.New("not implemented") }
func (s *uploadStream) SendMsg(any) error            { return errors.New("not implemented") }

type uploadStorage struct {
	putData    string
	putSize    int64
	putType    string
	putHistory []string
	putTypes   []string
	deleted    string
	getData    string
	getType    string
	getErr     error
	deleteErr  error
	putErrs    []error
}

func (s *uploadStorage) Get(context.Context, string) (storage.IObject, storage.IObjectInfo, error) {
	if s.getErr != nil {
		return nil, nil, s.getErr
	}
	return &uploadObject{
		Reader: bytes.NewReader([]byte(s.getData)),
	}, &uploadObjectInfo{
		contentType: s.getType,
	}, nil
}

func (s *uploadStorage) Stat(context.Context, string) (storage.IObjectInfo, error) {
	return nil, errors.New("not implemented")
}

func (s *uploadStorage) Put(
	_ context.Context,
	_ string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	if len(s.putErrs) > 0 {
		err := s.putErrs[0]
		s.putErrs = s.putErrs[1:]
		if err != nil {
			return "", err
		}
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	s.putData = string(data)
	s.putSize = size
	s.putType = contentType
	s.putHistory = append(s.putHistory, string(data))
	s.putTypes = append(s.putTypes, contentType)
	return "stored", nil
}

func (s *uploadStorage) Delete(_ context.Context, key string) error {
	s.deleted = key
	return s.deleteErr
}

func (s *uploadStorage) List(context.Context, string, int, int) ([]*storage.FileInfo, error) {
	return nil, nil
}

func (s *uploadStorage) GetSpaceUsage(context.Context) (int64, error) { return 0, nil }

type uploadObject struct {
	*bytes.Reader
}

func (o *uploadObject) Close() error { return nil }

type uploadObjectInfo struct {
	contentType string
}

func (i *uploadObjectInfo) GetName() string            { return "file" }
func (i *uploadObjectInfo) GetExtension() string       { return "" }
func (i *uploadObjectInfo) GetContentType() string     { return i.contentType }
func (i *uploadObjectInfo) GetSize() int64             { return 0 }
func (i *uploadObjectInfo) GetLastModified() time.Time { return time.Time{} }
func (i *uploadObjectInfo) GetExpiration() time.Time   { return time.Time{} }
