package v1

import (
	"net/http"

	"github.com/galexrt/rphub/model"
	"github.com/gin-gonic/gin"
)

func JobRoutes(g *gin.RouterGroup) {
	g.GET("/", Jobs)
}

// Jobs godoc
//
//	@Summary	Return registered Jobs
//	@Schemes
//	@Description	Return list of currently registered jobs
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Pong
//	@Router			/jobs/ [get]
func Jobs(g *gin.Context) {
	var jobs []model.Job
	model.DB.Find(&jobs)

	g.JSON(http.StatusOK, jobs)
}
