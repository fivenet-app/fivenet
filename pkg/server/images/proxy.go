package images

import (
	"net/http"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"willnorris.com/go/imageproxy"
)

type ImageProxy struct {
	logger *zap.Logger

	config config.ImageProxy
}

func New(logger *zap.Logger, cfg *config.BaseConfig) *ImageProxy {
	ip := &ImageProxy{
		logger: logger,
		config: cfg.ImageProxy,
	}

	return ip
}

func (p *ImageProxy) RegisterHTTP(e *gin.Engine) {
	if !p.config.Enabled {
		return
	}

	// Image Proxy
	proxy := imageproxy.NewProxy(http.DefaultTransport, imageproxy.NopCache)
	proxy.Logger = zap.NewStdLog(p.logger.Named("image_proxy"))

	proxy.AllowHosts = p.config.Options.AllowHosts
	proxy.DenyHosts = p.config.Options.DenyHosts

	// Example URL: http://localhost:3000/api/image_proxy/500/https://octodex.github.com/images/codercat.jpg
	e.GET("/api/image_proxy/*url", gin.WrapH(http.StripPrefix("/api/image_proxy", proxy)))
}
