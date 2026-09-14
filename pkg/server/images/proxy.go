package images

import (
	"context"
	"errors"
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

var nonPublicIPPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("2001:db8::/32"),
}

var publicIPv6Prefix = netip.MustParsePrefix("2000::/3")

// ImageProxy provides a small HTTP image proxy for remote images.
type ImageProxy struct {
	// logger is used for logging proxy activity and errors.
	logger *zap.Logger

	// config holds the image proxy configuration options.
	config config.ImageProxy

	client *http.Client
}

// New creates a new ImageProxy instance with the provided logger and configuration.
func New(logger *zap.Logger, cfg *config.Config) *ImageProxy {
	return &ImageProxy{
		logger: logger,
		config: cfg.ImageProxy,
		client: newHTTPClient(cfg.ImageProxy),
	}
}

// RegisterHTTP registers the image proxy HTTP handler on the provided Gin engine.
// The handler proxies image requests and applies restrictions from the config.
func (p *ImageProxy) RegisterHTTP(e *gin.Engine) {
	// Example URLs for the image proxy:
	// - Plain URL: http://localhost:3000/api/image_proxy/https://octodex.github.com/images/codercat.jpg
	e.GET(Path+"/*url", p.handle)
}

func newHTTPClient(options config.ImageProxy) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = safeDialContext

	return &http.Client{
		Timeout:   proxyTimeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return errors.New("too many redirects")
			}
			if !isAllowedURL(options, req.URL) {
				return errors.New("redirect host is not allowed")
			}
			return nil
		},
	}
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

	response, err := p.client.Do(upstreamRequest)
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

	contentType := response.Header.Get("Content-Type")
	if parsed, _, parseErr := mime.ParseMediaType(contentType); parseErr == nil {
		contentType = parsed
	}
	if contentType != "" && contentType != "application/octet-stream" &&
		contentType != "binary/octet-stream" && !strings.HasPrefix(contentType, "image/") {
		c.String(http.StatusForbidden, "remote resource is not an image")
		return
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxImageSize+1))
	if err != nil {
		c.String(http.StatusBadGateway, "failed to read image")
		return
	}
	if len(body) > maxImageSize {
		c.String(http.StatusRequestEntityTooLarge, "image is too large")
		return
	}
	if contentType == "" || contentType == "application/octet-stream" ||
		contentType == "binary/octet-stream" {
		contentType = http.DetectContentType(body)
	}
	if !strings.HasPrefix(contentType, "image/") {
		c.String(http.StatusForbidden, "remote resource is not an image")
		return
	}

	copyResponseHeaders(c.Writer.Header(), response.Header)
	setMinimumCacheDuration(c.Writer.Header(), p.config.MinimumCacheDuration)
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
		return nil, errors.New("missing target URL")
	}
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}
	target, err := url.Parse(decoded)
	if err != nil || target.Scheme == "" || target.Hostname() == "" {
		return nil, errors.New("target is not an absolute URL")
	}
	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, errors.New("unsupported target scheme")
	}
	if target.User != nil {
		return nil, errors.New("target credentials are not allowed")
	}
	return target, nil
}

func (p *ImageProxy) allowedURL(target *url.URL) bool {
	return isAllowedURL(p.config, target)
}

func isAllowedURL(options config.ImageProxy, target *url.URL) bool {
	host := strings.ToLower(strings.TrimSuffix(target.Hostname(), "."))
	if address, err := netip.ParseAddr(host); err == nil {
		return options.AllowIPs && isPublicIP(address.AsSlice())
	}
	if strings.Contains(host, ":") {
		return false
	}
	for _, denied := range options.DenyHosts {
		if hostMatches(host, denied) {
			return false
		}
	}
	if len(options.AllowHosts) > 0 {
		allowed := false
		for _, candidate := range options.AllowHosts {
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
	if suffix, ok := strings.CutPrefix(candidate, "*."); ok {
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
		err = errors.New("target resolves to a non-public IP")
	}
	return nil, err
}

func isPublicIP(ip net.IP) bool {
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	address = address.Unmap()
	if !address.IsGlobalUnicast() || address.IsPrivate() || address.IsLoopback() ||
		address.IsLinkLocalUnicast() || address.IsLinkLocalMulticast() || address.IsMulticast() {
		return false
	}
	if address.Is6() && !publicIPv6Prefix.Contains(address) {
		return false
	}
	for _, prefix := range nonPublicIPPrefixes {
		if prefix.Contains(address) {
			return false
		}
	}
	return true
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
	maxAge := int64(minimum.Seconds())
	directives := strings.Split(cacheControl, ",")
	foundMaxAge := false
	for index, directive := range directives {
		directives[index] = strings.TrimSpace(directive)
		parts := strings.SplitN(directives[index], "=", 2)
		name := strings.ToLower(parts[0])
		if name == "no-store" || name == "private" || name == "no-cache" {
			return
		}
		if len(parts) == 2 && (name == "max-age" || name == "s-maxage") {
			if value, err := strconv.ParseInt(strings.Trim(parts[1], "\" "), 10, 64); err == nil {
				if name == "max-age" {
					foundMaxAge = true
				}
				if value > maxAge {
					maxAge = value
				} else {
					directives[index] = fmt.Sprintf("%s=%d", name, maxAge)
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
	header.Set("Cache-Control", strings.Join(directives, ", "))
}
