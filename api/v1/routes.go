package v1

import (
	"net/http"

	"github.com/galexrt/arpanet/docs"
	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/query"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(r *gin.Engine) {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "arpanet API v1"
	docs.SwaggerInfo.Description = "arpanet Server"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	v1 := r.Group("/api/v1")
	{
		v1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, "Welcome to the aRPaNet API!")
		})
		v1.GET("/ping", Ping)

		// Register Type Routes
		AuthRoutes(v1.Group("/auth"))
		JobRoutes(v1.Group("/jobs"))
		DocumentsRoutes(v1.Group("/documents"))
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/documents", func(c *gin.Context) {
		info, _ := auth.GetSessionInfo(c)
		_ = info

		d := query.Document
		dja := query.DocumentJobAccess
		dua := query.DocumentUserAccess
		documents, err := d.LeftJoin(dja, dja.DocumentID.EqCol(d.ID)).LeftJoin(dua, dua.DocumentID.EqCol(d.ID)).
			Where(
				d.Where(d.Public.Is(true)).
					Or(d.Creator.Eq(info.Identifier)).
					Or(
						dua.Where(
							dua.Access.Neq(model.BlockedAccessRole),
							dua.Identifier.Eq(info.Identifier),
						),
					).
					Or(
						dja.Where(
							dja.Access.Neq(model.BlockedAccessRole),
							dja.Name.Eq(info.Job),
							dja.MinimumGrade.Lte(info.JobGrade),
						),
					),
			).
			Order(d.CreatedAt.Desc()).
			Preload(
				d.Jobs.On(dja.Name.Eq(info.Job)),
				d.Users.On(dua.Identifier.Eq(info.Identifier)),
			).
			Find()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSONP(http.StatusOK, documents)
	})
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
