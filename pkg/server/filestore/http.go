package filestore

import (
	"fmt"
	"net/http"
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
		g.GET("/:category/:name", func(c *gin.Context) {
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

			c.JSON(http.StatusOK, nil)
		})
	}
}
