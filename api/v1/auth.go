package v1

import (
	"fmt"
	"net/http"

	"github.com/galexrt/rphub/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(g *gin.RouterGroup) {
	g.GET("/", GetAuth)
	g.POST("/", PostAuth)
}

// Auth godoc
//
//	@Summary	Show authentication status
//	@Schemes
//	@Description	Shows your authentication status
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Pong
//	@Router			/auth/ [get]
func GetAuth(c *gin.Context) {
	info, _ := auth.GetSessionInfo(c)

	c.JSON(http.StatusOK, info)
}

// Auth godoc
//
//	@Summary	Authenticate yourself
//	@Schemes
//	@Description	Authenticate yourself against the API
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Pong
//	@Router			/auth/ [post]
func PostAuth(c *gin.Context) {
	err := auth.SaveSessionInfo(c, &auth.SessionInfo{
		ID:         1,
		Identifier: "fcee377a1fda007a8d2cc764a0a272e04d8c5d57",
		Job:        "ambulance",
		JobGrade:   20,
		Admin:      false,
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	info, _ := auth.GetSessionInfo(c)
	c.JSON(http.StatusOK, info)
}
