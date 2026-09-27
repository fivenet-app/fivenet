package filestore

import (
	"errors"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCheckFileLimit(t *testing.T) {
	t.Parallel()

	t.Run("skips enforcement when disabled", func(t *testing.T) {
		t.Parallel()

		db, _ := newFilestoreDB(t)
		h := &Handler[int64]{db: db, fileLimit: 0}
		require.NoError(t, h.CheckFileLimit(t.Context(), 7))
	})

	t.Run("allows parent below limit", func(t *testing.T) {
		t.Parallel()

		db, mock := newFilestoreDB(t)
		h := newDeleteHandler(db, &uploadStorage{}, false)
		h.fileLimit = 2
		mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*document_id = \?;`).
			WithArgs(int64(7)).
			WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

		require.NoError(t, h.CheckFileLimit(t.Context(), 7))
	})

	t.Run("rejects parent at limit", func(t *testing.T) {
		t.Parallel()

		db, mock := newFilestoreDB(t)
		h := newDeleteHandler(db, &uploadStorage{}, false)
		h.fileLimit = 2
		mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*document_id = \?;`).
			WithArgs(int64(7)).
			WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(2)))

		err := h.CheckFileLimit(t.Context(), 7)
		require.Error(t, err)
		assert.Equal(t, codes.ResourceExhausted, status.Code(err))
	})

	t.Run("returns count query error", func(t *testing.T) {
		t.Parallel()

		db, mock := newFilestoreDB(t)
		h := newDeleteHandler(db, &uploadStorage{}, false)
		h.fileLimit = 2
		queryErr := errors.New("count failed")
		mock.ExpectQuery(`SELECT COUNT\(fivenet_documents_files\.file_id\).*FROM fivenet_documents_files.*WHERE .*document_id = \?;`).
			WithArgs(int64(7)).
			WillReturnError(queryErr)

		require.ErrorIs(t, h.CheckFileLimit(t.Context(), 7), queryErr)
	})
}

func TestUploadFromMetaSanitizesNameAndDelegatesUpload(t *testing.T) {
	t.Parallel()
	db, mock := newFilestoreDB(t)
	st := &uploadStorage{}
	stream := &uploadStream{
		packets: []*file.UploadFileRequest{dataRequest([]byte("hello"))},
		recvErr: io.EOF,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT .*fivenet_files.*WHERE .*file_path = \?.*FOR UPDATE;`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(`INSERT INTO fivenet_files .*file_path.*byte_size.*content_type.*meta`).
		WithArgs(sqlmock.AnyArg(), int64(5), "text/plain").
		WillReturnResult(sqlmock.NewResult(23, 1))
	mock.ExpectCommit()

	h := &Handler[int64]{db: db, store: st}
	resp, err := h.UploadFromMeta(t.Context(), &file.UploadMeta{
		Namespace:    "documents",
		OriginalName: "my file.txt",
		Size:         5,
	}, 7, stream)
	require.NoError(t, err)
	assert.Regexp(
		t,
		`^/api/filestore/documents/[0-9]{8}/[0-9a-f-]{36}-my_file\.txt$`,
		resp.GetUrl(),
	)
	assert.Equal(t, "hello", st.putData)
}
