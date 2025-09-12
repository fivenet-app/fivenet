package auth

import (
	"context"
)

type ClientTokenAuth struct {
	token string

	requireTs bool
}

func NewClientTokenAuth(token string, requireTs bool) *ClientTokenAuth {
	return &ClientTokenAuth{
		token: token,

		requireTs: requireTs,
	}
}

// GetRequestMetadata return value is the token passed in during creation mapped to Authorization request headers.
//
//nolint:unparam // nil error return is required to fullfil the interface.
func (t *ClientTokenAuth) GetRequestMetadata(
	ctx context.Context,
	in ...string,
) (map[string]string,
	error,
) {
	return map[string]string{
		"Authorization": "Bearer " + t.token,
	}, nil
}

func (t *ClientTokenAuth) RequireTransportSecurity() bool {
	return t.requireTs
}
