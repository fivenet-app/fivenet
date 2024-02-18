package images

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"willnorris.com/go/imageproxy"
)

var Module = fx.Module("image_proxy",
	fx.Provide(
		New,
	),
)

type ImageProxy struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) (*ImageProxy, error) {
	ip := &ImageProxy{
		logger: logger,
	}

	return ip, nil
}

func (p *ImageProxy) Register(e *gin.Engine) {
	// Image Proxy
	proxy := imageproxy.NewProxy(http.DefaultTransport, imageproxy.NopCache)

	// Example URL: http://localhost:3000/api/image_proxy/500/https://octodex.github.com/images/codercat.jpg
	e.GET("/api/image_proxy/*url", gin.WrapH(http.StripPrefix("/api/image_proxy", proxy)))
}
