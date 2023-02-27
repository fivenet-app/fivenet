package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	// Routes
	testing *Testing
}

func New(logger *zap.Logger) *Routes {
	return &Routes{
		logger: logger,

		// Routes
		testing: &Testing{},
	}
}

func (r *Routes) Register(e *gin.Engine) {
	e.GET("/ping", r.PingGET)

	r.testing.Register(e)
}

func (r *Routes) PingGET(g *gin.Context) {
	g.JSON(http.StatusOK, "Pong")
}
