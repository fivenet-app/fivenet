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

	options config.ImageProxyOptions
}

func New(logger *zap.Logger, cfg *config.Config) *ImageProxy {
	if !cfg.ImageProxy.Enabled {
		return nil
	}

	ip := &ImageProxy{
		logger:  logger,
		options: cfg.ImageProxy.Options,
	}

	return ip
}

func (p *ImageProxy) RegisterHTTP(e *gin.Engine) {
	// Image Proxy
	proxy := imageproxy.NewProxy(http.DefaultTransport, imageproxy.NopCache)
	proxy.Logger = zap.NewStdLog(p.logger.Named("image_proxy"))

	proxy.AllowHosts = p.options.AllowHosts
	proxy.DenyHosts = p.options.DenyHosts

	// Example URL: http://localhost:3000/api/image_proxy/500/https://octodex.github.com/images/codercat.jpg
	e.GET("/api/image_proxy/*url", gin.WrapH(http.StripPrefix("/api/image_proxy", proxy)))
}
