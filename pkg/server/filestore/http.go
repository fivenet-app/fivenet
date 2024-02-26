package filestore

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/storage"
	"github.com/gin-gonic/gin"
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
	g := e.Group("/api/filestore")
	{
		g.GET("/:prefix/*fileName", func(c *gin.Context) {
			prefix := c.Param("prefix")
			prefix = filepath.Clean(prefix)
			fileName := c.Param("fileName")
			fileName = filepath.Clean(fileName)
			if fileName == "" {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested"))
				return
			}

			object, info, err := s.st.Get(c, path.Join(prefix, fileName))
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to retrieve file from store. %w", err))
				return
			}
			defer object.Close()

			c.DataFromReader(200, info.GetSize(), info.GetContentType(), object, nil)
		})
	}
}