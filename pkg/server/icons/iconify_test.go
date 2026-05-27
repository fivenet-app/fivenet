package icons

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils/httperrors"
	"github.com/stretchr/testify/assert"
)

func TestValidateIconRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		path  string
		query url.Values
		valid bool
	}{
		{
			name:  "valid request",
			path:  "/mdi.json",
			query: url.Values{"icons": []string{"home"}},
			valid: true,
		},
		{
			name:  "valid request multiple icons",
			path:  "/mdi.json",
			query: url.Values{"icons": []string{"home,account,alert-circle"}},
			valid: true,
		},
		{
			name:  "missing icons param",
			path:  "/mdi.json",
			query: url.Values{},
			valid: false,
		},
		{
			name:  "empty icons value",
			path:  "/mdi.json",
			query: url.Values{"icons": []string{""}},
			valid: false,
		},
		{
			name:  "empty path",
			path:  "",
			query: url.Values{"icons": []string{"home"}},
			valid: false,
		},
		{
			name:  "path is slash",
			path:  "/",
			query: url.Values{"icons": []string{"home"}},
			valid: false,
		},
		{
			name:  "path does not end with .json",
			path:  "/mdi",
			query: url.Values{"icons": []string{"home"}},
			valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := validateIconRequest(tc.path, tc.query)
			assert.Equal(t, tc.valid, got)
		})
	}
}

func TestBuildTargetURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		apiURL     string
		path       string
		query      url.Values
		wantURL    string
		wantErr    bool
		wantStatus int
		assertSafe bool
	}{
		{
			name:    "valid request",
			apiURL:  "https://api.iconify.design",
			path:    "mdi.json",
			query:   url.Values{"icons": []string{"home,account"}},
			wantURL: "https://api.iconify.design/mdi.json?icons=home%2Caccount",
		},
		{
			name:       "invalid path suffix",
			apiURL:     "https://api.iconify.design",
			path:       "mdi",
			query:      url.Values{"icons": []string{"home"}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown icon set",
			apiURL:     "https://api.iconify.design",
			path:       "unknown.json",
			query:      url.Values{"icons": []string{"home"}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing icons parameter",
			apiURL:     "https://api.iconify.design",
			path:       "mdi.json",
			query:      url.Values{},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "too long URL",
			apiURL:     "https://api.iconify.design",
			path:       "mdi.json",
			query:      url.Values{"icons": []string{strings.Repeat("a", 300)}},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "encodes untrusted query safely",
			apiURL:     "https://api.iconify.design",
			path:       "mdi.json",
			query:      url.Values{"icons": []string{"home&foo=bar#frag"}},
			wantErr:    false,
			assertSafe: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotURL, err := buildTargetURL(tc.apiURL, tc.path, tc.query)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, gotURL)
				if tc.wantStatus != 0 {
					var statusErr httperrors.HTTPStatusError
					assert.ErrorAs(t, err, &statusErr)
					assert.Equal(t, tc.wantStatus, statusErr.StatusCode())
				}
				return
			}

			assert.NoError(t, err)
			if tc.wantURL != "" {
				assert.Equal(t, tc.wantURL, gotURL)
			}

			if tc.assertSafe {
				parsed, parseErr := url.Parse(gotURL)
				assert.NoError(t, parseErr)
				assert.Equal(t, "home&foo=bar#frag", parsed.Query().Get("icons"))
				assert.Empty(t, parsed.Query().Get("foo"))
				assert.Empty(t, parsed.Fragment)
			}
		})
	}
}
