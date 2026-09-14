package images

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// Path is the base path for the image proxy API endpoint.
	Path = "/api/image_proxy"

	// UserAgentPrefix is the prefix for the User-Agent header sent by the image proxy.
	UserAgentPrefix = "FiveNet Image Proxy "

	proxyTimeout = 15 * time.Second
	maxImageSize = 10 << 20
	maxRedirects = 10
)

// ImageProxy provides a small HTTP image proxy for remote images.
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
	// Example URLs for the image proxy:
	// - Plain URL: http://localhost:3000/api/image_proxy/https://octodex.github.com/images/codercat.jpg
	// - Base64 encoded URL: http://localhost:3000/api/image_proxy/aHR0cHM6Ly9vY3RvZGV4LmdpdGh1Yi5jb20vaW1hZ2VzL2NvZGVyY2F0LmpwZw
	e.GET(Path+"/*url", p.handle)
}

func (p *ImageProxy) handle(c *gin.Context) {
	target, err := targetURL(c.Request)
	if err != nil {
		p.logger.Warn("invalid image proxy URL", zap.Error(err))
		c.String(http.StatusBadRequest, "invalid image URL")
		return
	}
	if !p.allowedURL(target) {
		p.logger.Warn("image proxy host is not allowed", zap.String("host", target.Hostname()))
		c.String(http.StatusForbidden, "image host is not allowed")
		return
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = safeDialContext
	client := &http.Client{
		Timeout:   proxyTimeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("too many redirects")
			}
			if !p.allowedURL(req.URL) {
				return fmt.Errorf("redirect host is not allowed")
			}
			return nil
		},
	}

	upstreamRequest, err := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodGet,
		target.String(),
		nil,
	)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid image URL")
		return
	}
	upstreamRequest.Header.Set("Accept", "image/*")
	upstreamRequest.Header.Set("User-Agent", UserAgentPrefix+version.Version)

	response, err := client.Do(upstreamRequest)
	if err != nil {
		p.logger.Warn(
			"failed to fetch proxied image",
			zap.Error(err),
			zap.String("url", target.Redacted()),
		)
		c.String(http.StatusBadGateway, "failed to fetch image")
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		c.Status(http.StatusNotFound)
		return
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		c.Status(http.StatusBadGateway)
		return
	}
	if response.ContentLength > maxImageSize {
		c.String(http.StatusRequestEntityTooLarge, "image is too large")
		return
	}

	reader := bufio.NewReader(io.LimitReader(response.Body, maxImageSize+1))
	contentType := response.Header.Get("Content-Type")
	if parsed, _, parseErr := mime.ParseMediaType(contentType); parseErr == nil {
		contentType = parsed
	}
	if contentType == "" || contentType == "application/octet-stream" ||
		contentType == "binary/octet-stream" {
		peek, peekErr := reader.Peek(512)
		if peekErr != nil && peekErr != io.EOF && peekErr != bufio.ErrBufferFull {
			c.String(http.StatusBadGateway, "could not inspect image")
			return
		}
		contentType = http.DetectContentType(peek)
	}
	if !strings.HasPrefix(contentType, "image/") {
		c.String(http.StatusForbidden, "remote resource is not an image")
		return
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		c.String(http.StatusBadGateway, "failed to read image")
		return
	}
	if len(body) > maxImageSize {
		c.String(http.StatusRequestEntityTooLarge, "image is too large")
		return
	}

	copyResponseHeaders(c.Writer.Header(), response.Header)
	setMinimumCacheDuration(c.Writer.Header(), p.config.Options.MinimumCacheDuration)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.Itoa(len(body)))
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Security-Policy", "script-src 'none'")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(body)
}

func targetURL(request *http.Request) (*url.URL, error) {
	path := strings.TrimPrefix(request.URL.EscapedPath(), Path)
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return nil, fmt.Errorf("missing target URL")
	}
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}
	target, err := url.Parse(decoded)
	if err != nil || target.Scheme == "" || target.Hostname() == "" {
		return nil, fmt.Errorf("target is not an absolute URL")
	}
	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, fmt.Errorf("unsupported target scheme")
	}
	if target.User != nil {
		return nil, fmt.Errorf("target credentials are not allowed")
	}
	return target, nil
}

func (p *ImageProxy) allowedURL(target *url.URL) bool {
	host := strings.ToLower(strings.TrimSuffix(target.Hostname(), "."))
	for _, denied := range p.config.Options.DenyHosts {
		if hostMatches(host, denied) {
			return false
		}
	}
	if len(p.config.Options.AllowHosts) > 0 {
		allowed := false
		for _, candidate := range p.config.Options.AllowHosts {
			if hostMatches(host, candidate) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

func hostMatches(host, candidate string) bool {
	candidate = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(candidate), "."))
	if strings.HasPrefix(candidate, "*.") {
		suffix := strings.TrimPrefix(candidate, "*.")
		return strings.HasSuffix(host, "."+suffix)
	}
	return host == candidate
}

func safeDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	dialer := net.Dialer{}
	for _, resolved := range addresses {
		if !isPublicIP(resolved.IP) {
			continue
		}
		connection, dialErr := dialer.DialContext(
			ctx,
			network,
			net.JoinHostPort(resolved.IP.String(), port),
		)
		if dialErr == nil {
			return connection, nil
		}
		err = dialErr
	}
	if err == nil {
		err = fmt.Errorf("target resolves to a non-public IP")
	}
	return nil, err
}

func isPublicIP(ip net.IP) bool {
	address, err := netip.ParseAddr(ip.String())
	return err == nil && address.IsGlobalUnicast() && !address.IsPrivate()
}

func copyResponseHeaders(dst, src http.Header) {
	for _, name := range []string{"Cache-Control", "ETag", "Last-Modified"} {
		for _, value := range src.Values(name) {
			dst.Add(name, value)
		}
	}
}

func setMinimumCacheDuration(header http.Header, minimum time.Duration) {
	if minimum <= 0 {
		return
	}
	cacheControl := header.Get("Cache-Control")
	lowerCacheControl := strings.ToLower(cacheControl)
	if strings.Contains(lowerCacheControl, "no-store") ||
		strings.Contains(lowerCacheControl, "private") {
		return
	}

	maxAge := int64(minimum.Seconds())
	directives := strings.Split(cacheControl, ",")
	foundMaxAge := false
	for index, directive := range directives {
		parts := strings.SplitN(strings.TrimSpace(directive), "=", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "max-age") {
			if value, err := strconv.ParseInt(strings.Trim(parts[1], "\" "), 10, 64); err == nil {
				foundMaxAge = true
				if value > maxAge {
					maxAge = value
				} else {
					directives[index] = fmt.Sprintf("max-age=%d", maxAge)
				}
			}
		}
	}

	if cacheControl == "" {
		header.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		return
	}
	if !foundMaxAge {
		directives = append(directives, fmt.Sprintf("max-age=%d", maxAge))
	}
	header.Set("Cache-Control", strings.Join(directives, ","))
}
