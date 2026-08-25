package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestClearSiteDataReturnsHtmlRedirectPage(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	_, router := gin.CreateTestContext(w)

	routes := &Routes{
		logger: zaptest.NewLogger(t),
		cfg: &config.Config{
			HTTP: config.HTTP{
				Sessions: config.Sessions{
					Domain: "localhost",
				},
			},
		},
	}
	routes.RegisterHTTP(router)

	req := httptest.NewRequest(http.MethodGet, "/api/clear-site-data", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Clear-Site-Data"), `"cache"`)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, w.Body.String(), `meta http-equiv="refresh" content="2;url=/"`)
	assert.Contains(t, w.Body.String(), "Reset complete")
}
