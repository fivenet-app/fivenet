package icons

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/httperrors"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	iconifygo "github.com/galexrt/iconify-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const (
	// Path is the base path for the icon API endpoint.
	Path = "/api/icons"

	// UserAgentPrefix is the prefix for the User-Agent header sent by the icon API proxy.
	UserAgentPrefix = "FiveNet Icon Proxy "
)

var (
	loadedIconSets      = []string{"mdi", "simple-icons", "flagpack"}
	loadedIconSetLookup = map[string]struct{}{
		"mdi":          {},
		"simple-icons": {},
		"flagpack":     {},
	}

	ErrUnknownIconSet         = errors.New("unknown icon set")
	ErrInvalidCollectionName  = errors.New("invalid icon collection file name")
	ErrMissingIconsQueryParam = errors.New("missing icons query parameter")
	ErrRequestURLTooLong      = errors.New("request URL too long")
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
			Path,
			p.Config.Icons.Path,
			iconifygo.WithHandlers("json"),
			iconifygo.WithPreloadIconsets(loadedIconSets),
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
		e.GET(Path+"/*path", func(c *gin.Context) {
			if !validateIconRequest(c.Param("path"), c.Request.URL.Query()) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid icon request"})
				return
			}

			i.server.HandlerFunc().ServeHTTP(c.Writer, c.Request)
		})
	} else {
		i.registerProxyHandler(e)
	}
}

func (i *IconifyAPI) registerProxyHandler(e *gin.Engine) {
	// Proxy requests to iconify API if enabled (make sure the request is a valid json icon request)
	e.GET(Path+":path", func(c *gin.Context) {
		// Validate the request and build the target URL
		path := c.Param("path")
		query := c.Request.URL.Query()
		if !validateIconRequest(path, query) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid icon request"})
			return
		}

		req, err := validateAndBuildRequest(c, i.apiURL, path, query)
		if err != nil {
			// Error response already handled in validateAndBuildRequest
			return
		}

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

func validateAndBuildRequest(
	c *gin.Context,
	apiURL string,
	path string,
	query url.Values,
) (*http.Request, error) {
	// Build the target URL for the iconify API request
	targetURL, err := buildTargetURL(apiURL, path, query)
	if err != nil {
		if reqErr, ok := errors.AsType[httperrors.HTTPStatusError](err); ok {
			c.AbortWithStatusJSON(reqErr.StatusCode(), gin.H{"error": err.Error()})
			return nil, err
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to build proxy request"},
		)
		return nil, err
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, targetURL, nil)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to create proxy request"},
		)
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgentPrefix+version.Version)
	req.Header.Set("Accept", "application/json")
	return req, nil
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
	// Ensure the path (file name) is valid
	if !strings.HasSuffix(path, ".json") {
		return "", httperrors.NewHTTPRequestError(http.StatusBadRequest, ErrInvalidCollectionName)
	}
	// Ensure the path corresponds to a known/loaded icon set
	iconSet := strings.TrimSuffix(path, ".json")
	if _, ok := loadedIconSetLookup[iconSet]; !ok {
		return "", httperrors.NewHTTPRequestError(http.StatusBadRequest, ErrInvalidCollectionName)
	}

	// Safely construct the target URL
	targetURL, err := url.JoinPath(apiURL, iconSet+".json")
	if err != nil {
		return "", err
	}
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	q := query.Get("icons")
	if q == "" {
		return "", httperrors.NewHTTPRequestError(http.StatusBadRequest, ErrMissingIconsQueryParam)
	}
	// Encode query values so untrusted input cannot alter URL structure.
	values := parsedURL.Query()
	values.Set("icons", q)
	parsedURL.RawQuery = values.Encode()

	u := parsedURL.String()
	// Arbitrary limit to prevent "excessively" long URLs
	if len(u) > 256 {
		return "", httperrors.NewHTTPRequestError(http.StatusBadRequest, ErrRequestURLTooLong)
	}
	return u, nil
}
