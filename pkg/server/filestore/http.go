package filestore

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/storage"
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
		MaxAge:               cachecontrol.Duration(48 * time.Hour),
		SMaxAge:              nil,
		Immutable:            false,
		StaleWhileRevalidate: cachecontrol.Duration(1 * time.Hour),
		StaleIfError:         cachecontrol.Duration(1 * time.Hour),
	}))

	{
		g.GET("/:prefix/*fileName", func(c *gin.Context) {
			prefix := c.Param("prefix")
			prefix = filepath.Clean(prefix)
			fileName := c.Param("fileName")
			fileName = filepath.Clean(fileName)
			if prefix == "" || fileName == "" {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested"))
				return
			}

			object, objInfo, err := s.st.Get(c, path.Join(prefix, fileName))
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
		})
	}
}
