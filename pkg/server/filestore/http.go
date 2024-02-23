package filestore

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/storage"
	"github.com/gin-gonic/gin"
)

type FilestoreHTTP struct {
	st *storage.Storage
	tm *auth.TokenMgr
}

func New(st *storage.Storage, tm *auth.TokenMgr) *FilestoreHTTP {
	return &FilestoreHTTP{
		st: st,
		tm: tm,
	}
}

func (s *FilestoreHTTP) RegisterHTTP(e *gin.Engine) {
	g := e.Group("/api/filestore")
	{
		g.GET("/get/:fileName", func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				c.AbortWithError(http.StatusForbidden, fmt.Errorf("invalid authorization token"))
				return
			}
			claims, err := s.tm.ParseWithClaims(token)
			if err != nil {
				c.AbortWithError(http.StatusForbidden, fmt.Errorf("invalid authorization token contents"))
				return
			}
			_ = claims

			fileName := c.Param("fileName")
			fileName = filepath.Clean(fileName)
			if fileName == "" {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid file requested"))
				return
			}

			object, info, err := s.st.Get(c, fileName)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to retrieve file from store. %w", err))
				return
			}

			c.DataFromReader(200, info.GetSize(), info.GetContentType(), object, nil)
		})
	}
}
