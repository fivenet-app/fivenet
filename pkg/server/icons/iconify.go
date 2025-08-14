package icons

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	iconifygo "github.com/andyburri/iconify-go"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const cacheExpirationDuration = 4 * time.Hour

type IconifyAPI struct {
	server *iconifygo.IconifyServer

	proxy  bool
	apiURL string

	cache *cache.Cache[string, []byte]

	client *http.Client
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Config *config.Config
}

// New creates a new ImageProxy instance with the provided logger and configuration.
func New(p Params) *IconifyAPI {
	ctxCancel, cancel := context.WithCancel(context.Background())

	ip := &IconifyAPI{
		server: iconifygo.NewIconifyServer("/api/icons", p.Config.Icons.Path, "json"),

		proxy:  p.Config.Icons.Proxy,
		apiURL: p.Config.Icons.APIURL,

		cache: cache.NewContext(ctxCancel, cache.AsLRU[string, []byte]()),

		client: http.DefaultClient,
	}

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return ip
}

func (i *IconifyAPI) RegisterHTTP(e *gin.Engine) {
	// Register the iconify API handler for local serving or proxy
	if !i.proxy {
		e.GET("/api/icons/*path", func(c *gin.Context) {
			recorder := &responseCapture{ResponseWriter: c.Writer, buf: &bytes.Buffer{}}

			if !validateIconRequest(c.Param("path"), c.Request.URL.Query()) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid icon request"})
				return
			}

			i.server.HandlerFunc().ServeHTTP(recorder, c.Request)

			// Cache the response body
			i.cache.Set(
				c.Request.URL.String(),
				recorder.buf.Bytes(),
				cache.WithExpiration(cacheExpirationDuration),
			)
		})
	} else {
		// Proxy requests to iconify API if enabled (make sure the request is a valid json icon request)
		e.GET("/api/icons/:path", func(c *gin.Context) {
			// Validate the request and extract the target URL
			path := c.Param("path")
			query := c.Request.URL.Query()
			if !validateIconRequest(path, query) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid icon request"})
				return
			}

			// Build the target URL for the iconify API request
			targetURL := buildTargetURL(i.apiURL, path, query)
			if len(targetURL) > 1024 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request URL too long"})
				return
			}

			if cachedResponse, ok := i.cache.Get(targetURL); ok {
				c.Writer.Header().Set("Content-Type", "application/json")
				c.Writer.Write(cachedResponse)
				return
			}

			req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, targetURL, c.Request.Body)
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

			// Need to read the body to be able to cache it..
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to read response body"})
				return
			}

			i.cache.Set(targetURL, body, cache.WithExpiration(cacheExpirationDuration))

			c.Writer.Header().Set("Content-Type", "application/json")
			c.Status(resp.StatusCode)
			c.Writer.Write(body)
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

func buildTargetURL(apiURL string, path string, query url.Values) string {
	targetURL := apiURL + path
	if q := query.Encode(); q != "" {
		targetURL += "?" + q
	}
	return targetURL
}

// responseCapture wraps gin.ResponseWriter to capture output for caching icons.
type responseCapture struct {
	gin.ResponseWriter

	buf *bytes.Buffer
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}
