package auth

import "context"

type ClientTokenAuth struct {
	token string
}

func NewClientTokenAuth(token string) *ClientTokenAuth {
	return &ClientTokenAuth{
		token: token,
	}
}

// Return value is mapped to request headers.
func (t ClientTokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}
