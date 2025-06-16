package wk

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// securityTxtTpl is the template for the security.txt file served under /.well-known/security.txt.
const securityTxtTpl = `Contact: https://github.com/fivenet-app/fivenet/v2025/blob/main/SECURITY.md
Expires: %s
Preferred-Languages: en, de
Policy: https://github.com/fivenet-app/fivenet/v2025/blob/main/SECURITY.md
`

// WK provides handlers for /.well-known endpoints such as security.txt and change-password.
type WK struct {
	// securityTxt contains the generated security.txt content with an expiration date.
	securityTxt string
}

// New creates a new WK instance with a generated security.txt content.
func New() *WK {
	return &WK{
		// ISO-8601 format see https://www.rfc-editor.org/rfc/rfc9116#section-2.5.5
		securityTxt: fmt.Sprintf(securityTxtTpl, time.Now().AddDate(0, 4, 0).Format("2006-01-02T15:04:05:07.000Z")),
	}
}

// RegisterHTTP registers the /.well-known endpoints on the provided Gin engine.
// - /.well-known/change-password redirects to the forgot password tab on the login page.
// - /.well-known/security.txt serves the security.txt file.
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
