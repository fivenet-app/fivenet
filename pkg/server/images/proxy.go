package images

import (
	"fmt"
	"net/http"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"willnorris.com/go/imageproxy"
)

const Path = "/api/image_proxy"

// ImageProxy provides an HTTP image proxy handler using the imageproxy package.
type ImageProxy struct {
	// logger is used for logging proxy activity and errors.
	logger *zap.Logger

	// config holds the image proxy configuration options.
	config config.ImageProxy
}

// New creates a new ImageProxy instance with the provided logger and configuration.
func New(logger *zap.Logger, cfg *config.Config) *ImageProxy {
	ip := &ImageProxy{
		logger: logger,
		config: cfg.ImageProxy,
	}

	return ip
}

// RegisterHTTP registers the image proxy HTTP handler on the provided Gin engine.
// If the proxy is not enabled in the config, this function does nothing.
// The handler proxies image requests and applies restrictions from the config.
func (p *ImageProxy) RegisterHTTP(e *gin.Engine) {
	if !p.config.Enabled {
		return
	}

	// Image Proxy setup
	proxy := imageproxy.NewProxy(http.DefaultTransport, imageproxy.NopCache)
	proxy.Logger = zap.NewStdLog(p.logger.Named("image_proxy"))

	proxy.AllowHosts = p.config.Options.AllowHosts
	proxy.DenyHosts = p.config.Options.DenyHosts
	proxy.UserAgent = fmt.Sprintf("FiveNet Image Proxy %s", version.Version)
	proxy.ContentTypes = []string{"image/*"}
	proxy.ScaleUp = false

	// Example URLs:
	// - http://localhost:3000/api/image_proxy/500/https://octodex.github.com/images/codercat.jpg
	// - http://localhost:3000/api/image_proxy/x/aHR0cHM6Ly9vY3RvZGV4LmdpdGh1Yi5jb20vaW1hZ2VzL2NvZGVyY2F0LmpwZw
	e.GET(Path+"/*url", gin.WrapH(http.StripPrefix(Path, proxy)))
}
