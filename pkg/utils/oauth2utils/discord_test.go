package oauth2utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/stretchr/testify/assert"
)

func TestIsDiscordTokenExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "http error invalid grant body",
			err: &httputil.HTTPError{
				Status: 400,
				Body: []byte(
					`{"error":"invalid_grant","error_description":"refresh token is invalid"}`,
				),
			},
			want: true,
		},
		{
			name: "wrapped http error invalid grant body",
			err: fmt.Errorf("discord request failed: %w", &httputil.HTTPError{
				Status: 400,
				Body:   []byte(`{"error":"invalid_grant"}`),
			}),
			want: true,
		},
		{
			name: "generic invalid grant string",
			err:  errors.New("discord oauth failed: invalid_grant"),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("discord oauth failed: rate limited"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsDiscordTokenExpired(tt.err); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
