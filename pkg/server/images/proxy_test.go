package images

import (
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestTargetURL(t *testing.T) {
	t.Parallel()

	request := &http.Request{URL: &url.URL{
		Path:    Path + "/https://example.com/image.png",
		RawPath: Path + "/https%3A%2F%2Fexample.com%2Fimage.png",
	}}
	target, err := targetURL(request)

	require.NoError(t, err)
	require.Equal(t, "https://example.com/image.png", target.String())
}

func TestTargetURLRejectsCredentialsAndUnsupportedSchemes(t *testing.T) {
	t.Parallel()

	for _, target := range []struct {
		path    string
		rawPath string
	}{
		{path: "file:///etc/passwd", rawPath: "file%3A%2F%2F%2Fetc%2Fpasswd"},
		{path: "https://user:pass@example.com/image.png", rawPath: "https%3A%2F%2Fuser%3Apass%40example.com%2Fimage.png"},
	} {
		request := &http.Request{
			URL: &url.URL{Path: Path + "/" + target.path, RawPath: Path + "/" + target.rawPath},
		}
		_, err := targetURL(request)
		require.Error(t, err)
	}
}

func TestIsAllowedURL(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name    string
		options config.ImageProxy
		target  string
		want    bool
	}{
		{name: "unrestricted hostname", target: "https://example.com/photo.png", want: true},
		{name: "matching wildcard", options: config.ImageProxy{AllowHosts: []string{"*.example.com"}}, target: "https://images.example.com/photo.png", want: true},
		{name: "nested wildcard", options: config.ImageProxy{AllowHosts: []string{"*.example.com"}}, target: "https://cdn.images.example.com/photo.png", want: true},
		{name: "wildcard excludes apex", options: config.ImageProxy{AllowHosts: []string{"*.example.com"}}, target: "https://example.com/photo.png"},
		{name: "wildcard excludes suffix bypass", options: config.ImageProxy{AllowHosts: []string{"*.example.com"}}, target: "https://notexample.com/photo.png"},
		{name: "normalizes host and rule", options: config.ImageProxy{AllowHosts: []string{" *.example.com. "}}, target: "https://IMAGES.EXAMPLE.COM./photo.png", want: true},
		{name: "deny overrides allow", options: config.ImageProxy{AllowHosts: []string{"*.example.com"}, DenyHosts: []string{"blocked.example.com"}}, target: "https://blocked.example.com/photo.png"},
		{name: "wildcard deny", options: config.ImageProxy{DenyHosts: []string{"*.internal.example.com"}}, target: "https://api.internal.example.com/photo.png"},
		{name: "public IPv4 disabled by default", target: "https://8.8.8.8/photo.png"},
		{name: "public IPv4 enabled", options: config.ImageProxy{AllowIPs: true}, target: "https://8.8.8.8/photo.png", want: true},
		{name: "public IPv6 enabled", options: config.ImageProxy{AllowIPs: true}, target: "https://[2606:4700:4700::1111]/photo.png", want: true},
		{name: "private IPv4 remains blocked", options: config.ImageProxy{AllowIPs: true}, target: "https://10.0.0.1/photo.png"},
		{name: "mapped loopback remains blocked", options: config.ImageProxy{AllowIPs: true}, target: "https://[::ffff:127.0.0.1]/photo.png"},
		{name: "IPv6 zone is blocked", options: config.ImageProxy{AllowIPs: true}, target: "https://[fe80::1%25eth0]/photo.png"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			target, err := url.Parse(test.target)
			require.NoError(t, err)
			require.Equal(t, test.want, isAllowedURL(test.options, target))
		})
	}
}

func TestHostMatches(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		host      string
		candidate string
		want      bool
	}{
		{host: "example.com", candidate: "example.com", want: true},
		{host: "images.example.com", candidate: "*.example.com", want: true},
		{host: "example.com", candidate: "*.example.com"},
		{host: "notexample.com", candidate: "*.example.com"},
		{host: "images.example.com", candidate: " *.EXAMPLE.COM. ", want: true},
	} {
		t.Run(test.host+"/"+test.candidate, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.want, hostMatches(test.host, test.candidate))
		})
	}
}

func TestIsPublicIP(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string
		ip   string
		want bool
	}{
		{name: "public IPv4", ip: "8.8.8.8", want: true},
		{name: "public IPv6", ip: "2606:4700:4700::1111", want: true},
		{name: "private IPv4", ip: "10.0.0.1"},
		{name: "IPv4 loopback", ip: "127.0.0.1"},
		{name: "IPv4 mapped loopback", ip: "::ffff:127.0.0.1"},
		{name: "IPv4 link local", ip: "169.254.169.254"},
		{name: "carrier-grade NAT", ip: "100.64.0.1"},
		{name: "benchmarking range", ip: "198.18.0.1"},
		{name: "IPv4 documentation range", ip: "203.0.113.1"},
		{name: "IPv6 unique local", ip: "fc00::1"},
		{name: "IPv6 loopback", ip: "::1"},
		{name: "IPv6 protocol assignment", ip: "2001::1"},
		{name: "IPv6 documentation range", ip: "2001:db8::1"},
		{name: "IPv6 6to4", ip: "2002::1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ip := net.ParseIP(test.ip)
			require.NotNil(t, ip)
			require.Equal(t, test.want, isPublicIP(ip))
		})
	}
}

func TestSetMinimumCacheDuration(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name    string
		initial string
		want    string
	}{
		{name: "empty", want: "public, max-age=1800"},
		{name: "longer max age", initial: "public, max-age=3600", want: "public, max-age=3600"},
		{name: "shorter max age", initial: "public, max-age=60", want: "public, max-age=1800"},
		{name: "shared max age", initial: "public, s-maxage=60", want: "public, s-maxage=1800, max-age=1800"},
		{name: "no cache", initial: "no-cache, max-age=60", want: "no-cache, max-age=60"},
		{name: "private", initial: "private, max-age=60", want: "private, max-age=60"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			header := make(http.Header)
			header.Set("Cache-Control", test.initial)
			setMinimumCacheDuration(header, 30*time.Minute)
			require.Equal(t, test.want, header.Get("Cache-Control"))
		})
	}
}
