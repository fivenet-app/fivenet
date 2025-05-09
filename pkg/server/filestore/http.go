package filestore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	cachecontrol "go.eigsys.de/gin-cachecontrol/v2"
)

type FilestoreHTTP struct {
	st storage.IStorage
	tm *auth.TokenMgr
}

func New(st storage.IStorage, tm *auth.TokenMgr) *FilestoreHTTP {
	return &FilestoreHTTP{
		st: st,
		tm: tm,
	}
}

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

func (s *FilestoreHTTP) HEAD(c *gin.Context) {
	prefix := c.Param("prefix")
	prefix = filepath.Clean(prefix)
	fileName := c.Param("fileName")
	fileName = filepath.Clean(fileName)
	if prefix == "" || fileName == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested"))
		return
	}

	filePath := path.Join(prefix, fileName)

	if _, err := s.st.Stat(c, filePath); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (s *FilestoreHTTP) GET(c *gin.Context) {
	prefix := c.Param("prefix")
	prefix = filepath.Clean(prefix)
	fileName := c.Param("fileName")
	fileName = filepath.Clean(fileName)
	if prefix == "" || fileName == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested"))
		return
	}

	filePath := path.Join(prefix, fileName)

	ctx, cancel := context.WithTimeout(c, 6*time.Second)
	defer cancel()

	url, err := s.st.GetURL(ctx, filePath, 1*time.Hour, url.Values{})
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to retrieve file url from store. %w", err))
		return
	}

	if url != nil {
		c.Redirect(http.StatusTemporaryRedirect, *url)
		return
	}

	object, objInfo, err := s.st.Get(ctx, path.Join(prefix, fileName))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to retrieve file from store. %w", err))
		return
	}
	defer object.Close()

	mimeType := filetype.GetType(objInfo.GetExtension())

	c.DataFromReader(200, objInfo.GetSize(), mimeType.MIME.Value, object, nil)
}
