package filestore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	cachecontrol "go.eigsys.de/gin-cachecontrol/v2"
)

// FilestoreHTTP provides HTTP handlers for file storage operations (GET/HEAD) via Gin.
type FilestoreHTTP struct {
	// st is the storage backend implementing IStorage.
	st storage.IStorage
	// tm is the token manager for authentication tokens.
	tm *auth.TokenMgr
}

// New creates a new FilestoreHTTP handler with the given storage backend and token manager.
func New(st storage.IStorage, tm *auth.TokenMgr) *FilestoreHTTP {
	return &FilestoreHTTP{
		st: st,
		tm: tm,
	}
}

// RegisterHTTP registers the HTTP endpoints for the filestore on the provided Gin engine.
// It sets up GET and HEAD routes with cache control middleware.
func (s *FilestoreHTTP) RegisterHTTP(e *gin.Engine) {
	g := e.Group("/api/filestore", cachecontrol.New(cachecontrol.Config{
		MustRevalidate:       true,
		NoCache:              false,
		NoStore:              false,
		NoTransform:          false,
		Public:               true,
		Private:              false,
		ProxyRevalidate:      true,
		MaxAge:               cachecontrol.Duration(4 * 24 * time.Hour),
		SMaxAge:              nil,
		Immutable:            false,
		StaleWhileRevalidate: cachecontrol.Duration(1 * 24 * time.Hour),
		StaleIfError:         cachecontrol.Duration(1 * 24 * time.Hour),
	}))

	g.HEAD("/:prefix/*fileName", s.HEAD)
	g.GET("/:prefix/*fileName", s.GET)
}

const (
	// fileRetrievalTimeout is the timeout for file retrieval operations.
	fileRetrievalTimeout = 6 * time.Second
)

// isPathSafe checks for path traversal attempts.
func isPathSafe(p string) bool {
	if p == "" || p == "." || p == "/" {
		return false
	}
	if p == ".." || len(p) >= 3 && p[:3] == "../" || len(p) >= 3 && p[len(p)-3:] == "/.." || p == "." {
		return false
	}
	if len(p) >= 2 && (p[:2] == ".." || p[len(p)-2:] == "..") {
		return false
	}
	return true
}

// HEAD handles HTTP HEAD requests for files in the filestore.
// It checks if the requested file exists and returns 404 if not found, 400 for invalid requests.
func (s *FilestoreHTTP) HEAD(c *gin.Context) {
	prefix := filepath.Clean(c.Param("prefix"))
	fileName := filepath.Clean(c.Param("fileName"))
	if !isPathSafe(prefix) || !isPathSafe(fileName) {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested (unsafe path)"))
		return
	}
	if prefix == "" || fileName == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested (empty)"))
		return
	}

	filePath := path.Join(prefix, fileName)

	objInfo, err := s.st.Stat(c, filePath)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ext := objInfo.GetExtension()
	mimeType := filetype.GetType(ext)
	mimeVal := "application/octet-stream"
	if mimeType != filetype.Unknown {
		mimeVal = mimeType.MIME.Value
	}
	c.Header("Content-Type", mimeVal)
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.GetSize()))
	// No body for HEAD
	c.Status(http.StatusOK)
}

// GET handles HTTP GET requests for files in the filestore.
// It streams the file to the client if found, or returns 404/400 for errors.
func (s *FilestoreHTTP) GET(c *gin.Context) {
	prefix := filepath.Clean(c.Param("prefix"))
	fileName := filepath.Clean(c.Param("fileName"))
	if !isPathSafe(prefix) || !isPathSafe(fileName) {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested (unsafe path)"))
		return
	}
	if prefix == "" || fileName == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested (empty)"))
		return
	}

	filePath := path.Join(prefix, fileName)

	ctx, cancel := context.WithTimeout(c, fileRetrievalTimeout)
	defer cancel()

	object, objInfo, err := s.st.Get(ctx, filePath)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Error(err)
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to retrieve file from store. %w", err))
		return
	}
	defer object.Close()

	ext := objInfo.GetExtension()
	mimeType := filetype.GetType(ext)
	mimeVal := "application/octet-stream"
	if mimeType != filetype.Unknown {
		mimeVal = mimeType.MIME.Value
	}

	c.DataFromReader(http.StatusOK, objInfo.GetSize(), mimeVal, object, nil)
}
