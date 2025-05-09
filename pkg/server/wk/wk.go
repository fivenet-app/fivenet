package wk

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const securityTxtTpl = `Contact: https://github.com/fivenet-app/fivenet/v2025/blob/main/SECURITY.md
Expires: %s
Preferred-Languages: en, de
Policy: https://github.com/fivenet-app/fivenet/v2025/blob/main/SECURITY.md
`

type WK struct {
	securityTxt string
}

func New() *WK {
	return &WK{
		// ISO-8601 format see https://www.rfc-editor.org/rfc/rfc9116#section-2.5.5
		securityTxt: fmt.Sprintf(securityTxtTpl, time.Now().AddDate(0, 4, 0).Format("2006-01-02T15:04:05:07.000Z")),
	}
}

func (w *WK) RegisterHTTP(e *gin.Engine) {
	g := e.Group("/.well-known")
	{
		g.GET("change-password", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/login?tab=forgotPassword#")
		})
		g.GET("security.txt", func(c *gin.Context) {
			c.String(http.StatusOK, w.securityTxt)
		})
	}
}
