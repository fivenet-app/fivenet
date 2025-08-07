package icons

import (
	"net/url"
	"testing"
)

func TestValidateIconRequest(t *testing.T) {
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
			got := validateIconRequest(tc.path, tc.query)
			if got != tc.valid {
				t.Errorf("validateIconRequest(%q, %v) = %v; want %v", tc.path, tc.query, got, tc.valid)
			}
		})
	}
}
