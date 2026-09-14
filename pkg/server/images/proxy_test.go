package images

import (
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
		request := &http.Request{URL: &url.URL{Path: Path + "/" + target.path, RawPath: Path + "/" + target.rawPath}}
		_, err := targetURL(request)
		require.Error(t, err)
	}
}

func TestAllowedURL(t *testing.T) {
	t.Parallel()

	proxy := &ImageProxy{config: config.ImageProxy{Options: config.ImageProxyOptions{
		AllowHosts: []string{"*.example.com"},
		DenyHosts:  []string{"blocked.example.com"},
	}}}

	allowed, _ := url.Parse("https://images.example.com/photo.png")
	denied, _ := url.Parse("https://blocked.example.com/photo.png")
	wrong, _ := url.Parse("https://example.net/photo.png")

	require.True(t, proxy.allowedURL(allowed))
	require.False(t, proxy.allowedURL(denied))
	require.False(t, proxy.allowedURL(wrong))
}

func TestSetMinimumCacheDuration(t *testing.T) {
	t.Parallel()

	got := make(http.Header)
	setMinimumCacheDuration(got, 30*time.Minute)
	require.Equal(t, "public, max-age=1800", got.Get("Cache-Control"))

	got.Set("Cache-Control", "public, max-age=3600")
	setMinimumCacheDuration(got, 30*time.Minute)
	require.Equal(t, "public, max-age=3600", got.Get("Cache-Control"))
}
