package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Routes {
	return &Routes{
		logger: logger,
	}
}

func (r *Routes) Register(e *gin.Engine) {
	e.GET("/api/ping", r.PingGET)
}

func (r *Routes) PingGET(g *gin.Context) {
	g.JSON(http.StatusOK, "Pong")
}
