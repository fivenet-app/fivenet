package routes

import (
	"fmt"
	"net/http"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(g *gin.RouterGroup) {
	g.GET("/", GetAuth)
	g.POST("/", PostAuth)
}

func GetAuth(c *gin.Context) {
	info, _ := auth.GetSessionInfo(c)

	c.JSON(http.StatusOK, info)
}

func PostAuth(c *gin.Context) {
	err := auth.SaveSessionInfo(c, &auth.SessionInfo{
		ID:         26061,
		Identifier: "char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57",
		CharIndex:  1,
		Job:        "ambulance",
		JobGrade:   20,
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	info, _ := auth.GetSessionInfo(c)
	c.JSON(http.StatusOK, info)
}
