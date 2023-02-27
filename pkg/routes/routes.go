package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	// Routes
	auth    *Auth
	testing *Testing
}

func New(logger *zap.Logger) *Routes {
	return &Routes{
		logger: logger,

		// Routes
		auth:    &Auth{},
		testing: &Testing{},
	}
}

func (rs *Routes) Register(e *gin.Engine) {
	e.GET("/ping", PingGET)

	rs.auth.Register(e)
	rs.testing.Register(e)
}

func PingGET(g *gin.Context) {
	g.JSON(http.StatusOK, "Pong")
}
