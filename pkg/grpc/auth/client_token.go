package auth

import "context"

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

// Return value is mapped to request headers.
//
//nolint:unparam // nil error return is required to fullfil the interface.
func (t *ClientTokenAuth) GetRequestMetadata(
	ctx context.Context,
	in ...string,
) (map[string]string,
	error,
) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t *ClientTokenAuth) RequireTransportSecurity() bool {
	return t.requireTs
}
