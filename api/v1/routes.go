package v1

import (
	"net/http"

	"github.com/galexrt/rphub/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(r *gin.Engine) {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "RPHub API v1"
	docs.SwaggerInfo.Description = "RPHub Server"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", Ping)
		AuthRoutes(v1.Group("/auth"))
		JobRoutes(v1.Group("/jobs"))
		DocumentsRoutes(v1.Group("/documents"))
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

//	@BasePath	/api/v1

// Ping godoc
//
//	@Summary	ping
//	@Schemes
//	@Description	do ping
//	@Tags			ping
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Pong
//	@Router			/ping [get]
func Ping(g *gin.Context) {
	g.JSON(http.StatusOK, "Pong")
}
