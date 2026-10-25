package filestore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestFilestoreHTTPResponseHeaders(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		method       string
		contentType  string
		extension    string
		expectType   string
		expectInline bool
	}{
		{
			name:         "GET image is inline",
			method:       http.MethodGet,
			contentType:  "image/webp",
			extension:    "webp",
			expectType:   "image/webp",
			expectInline: true,
		},
		{
			name:         "HEAD image is inline",
			method:       http.MethodHead,
			contentType:  "image/png",
			extension:    "png",
			expectType:   "image/png",
			expectInline: true,
		},
		{
			name:        "GET non-image is attachment",
			method:      http.MethodGet,
			contentType: "application/pdf",
			extension:   "pdf",
			expectType:  "application/pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := &testObjectInfo{
				contentType: tt.contentType,
				extension:   tt.extension,
				size:        int64(len("content")),
			}
			handler := &FilestoreHTTP{st: &testStorage{info: info}}
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = httptest.NewRequest(
				tt.method,
				"/api/filestore/files/asset."+tt.extension,
				nil,
			)
			ctx.Params = gin.Params{
				{Key: "prefix", Value: "files"},
				{Key: "fileName", Value: "/asset." + tt.extension},
			}

			if tt.method == http.MethodHead {
				handler.HEAD(ctx)
			} else {
				handler.GET(ctx)
			}

			require.Equal(t, http.StatusOK, response.Code)
			require.Equal(t, tt.expectType, response.Header().Get("Content-Type"))
			require.Equal(t, "nosniff", response.Header().Get("X-Content-Type-Options"))
			if tt.expectInline {
				require.Equal(
					t,
					contentDispositionInline,
					response.Header().Get("Content-Disposition"),
				)
			} else {
				require.Equal(
					t,
					contentDispositionAttachment,
					response.Header().Get("Content-Disposition"),
				)
			}
		})
	}
}

type testStorage struct {
	info storage.IObjectInfo
}

func (s *testStorage) Get(context.Context, string) (storage.IObject, storage.IObjectInfo, error) {
	return &testObject{Reader: bytes.NewReader([]byte("content"))}, s.info, nil
}

func (s *testStorage) Stat(context.Context, string) (storage.IObjectInfo, error) {
	return s.info, nil
}

func (*testStorage) Put(context.Context, string, io.Reader, int64, string) (string, error) {
	return "", errors.New("not implemented")
}

func (*testStorage) Delete(context.Context, string) error { return errors.New("not implemented") }

func (*testStorage) List(context.Context, string, int, int) ([]*storage.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (*testStorage) GetSpaceUsage(context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

type testObject struct {
	*bytes.Reader
}

func (*testObject) Close() error { return nil }

type testObjectInfo struct {
	contentType string
	extension   string
	size        int64
}

func (*testObjectInfo) GetName() string            { return "asset" }
func (i *testObjectInfo) GetExtension() string     { return i.extension }
func (i *testObjectInfo) GetContentType() string   { return i.contentType }
func (i *testObjectInfo) GetSize() int64           { return i.size }
func (*testObjectInfo) GetLastModified() time.Time { return time.Time{} }
func (*testObjectInfo) GetExpiration() time.Time   { return time.Time{} }
