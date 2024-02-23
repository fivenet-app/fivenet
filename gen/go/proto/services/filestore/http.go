package filestore

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterHTTP(e *gin.Engine) {
	g := e.Group("/api/filestore")
	{
		g.GET("/:category/:name", func(c *gin.Context) {
			c.JSON(http.StatusOK, nil)
		})
		// TODO Figure out the http based logic to easily and safely retrieve the files
	}
}
