package icons

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	iconifygo "github.com/galexrt/iconify-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type IconifyAPI struct {
	server *iconifygo.IconifyServer

	proxy  bool
	apiURL string

	client *http.Client
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Config *config.Config
}

// New creates a new ImageProxy instance with the provided logger and configuration.
func New(p Params) (*IconifyAPI, error) {
	var icServer *iconifygo.IconifyServer
	if !p.Config.Icons.Proxy {
		var err error
		icServer, err = iconifygo.NewIconifyServer(
			"/api/icons",
			p.Config.Icons.Path,
			iconifygo.WithHandlers("json"),
			iconifygo.WithPreloadIconsets([]string{"mdi", "simple-icons", "flagpack"}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create iconify server. %w", err)
		}
	}

	ip := &IconifyAPI{
		server: icServer,

		proxy:  p.Config.Icons.Proxy,
		apiURL: p.Config.Icons.APIURL,
	}

	// Set up default HTTP client for proxying requests (if enabled)
	if ip.proxy {
		ip.client = http.DefaultClient
	}

	return ip, nil
}

func (i *IconifyAPI) RegisterHTTP(e *gin.Engine) {
	// Register the iconify API handler for local serving or proxy
	if !i.proxy {
		e.GET("/api/icons/*path", func(c *gin.Context) {
			if !validateIconRequest(c.Param("path"), c.Request.URL.Query()) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid icon request"})
				return
			}

			i.server.HandlerFunc().ServeHTTP(c.Writer, c.Request)
		})
	} else {
		// Proxy requests to iconify API if enabled (make sure the request is a valid json icon request)
		e.GET("/api/icons/:path", func(c *gin.Context) {
			// Validate the request and build the target URL
			path := c.Param("path")
			query := c.Request.URL.Query()
			if !validateIconRequest(path, query) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid icon request"})
				return
			}

			// Build the target URL for the iconify API request
			targetURL, err := buildTargetURL(i.apiURL, path, query)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to build proxy request"})
				return
			}
			if len(targetURL) > 1024 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request URL too long"})
				return
			}

			req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, targetURL, nil)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create proxy request"})
				return
			}
			req.Header.Set("User-Agent", c.Request.UserAgent())
			req.Header.Set("Accept", "application/json")

			resp, err := i.client.Do(req)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "failed to proxy request"})
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				c.Writer.Header().Set("Content-Type", "application/json")
				c.AbortWithStatusJSON(resp.StatusCode, gin.H{"error": "failed to fetch icon data"})
				return
			}

			c.DataFromReader(http.StatusOK, resp.ContentLength, "application/json", resp.Body, nil)
		})
	}
}

func validateIconRequest(path string, query url.Values) bool {
	if path == "" || path == "/" || !strings.HasSuffix(path, ".json") {
		return false
	}

	if query.Get("icons") == "" {
		return false
	}

	return true
}

func buildTargetURL(apiURL string, path string, query url.Values) (string, error) {
	targetURL, err := url.JoinPath(apiURL, path)
	if err != nil {
		return "", err
	}

	q := query.Get("icons")
	if q == "" {
		return "", fmt.Errorf("missing icons query parameter")
	}

	targetURL += "?icons=" + q

	return targetURL, nil
}
