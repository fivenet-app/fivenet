package v1

import (
	"net/http"

	"github.com/galexrt/arpanet/query"
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
	jobs, err := query.Job.Find()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, jobs)
}
