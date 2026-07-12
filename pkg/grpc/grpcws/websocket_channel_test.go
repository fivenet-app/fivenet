package grpcws

import (
	"net/http"
	"net/http/httptest"
	"testing"

	grpcwsproto "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/grpcws"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/stretchr/testify/require"
)

func TestApplyControlAuthAllowsTokenlessDowngrade(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Cookie", grpcauth.AccCookieName+"=acc-token")

	ws := &WebsocketChannel{
		req: req,
		validateTokenFunc: func(token string) (bool, error) {
			return token == "char-token", nil
		},
	}

	err := ws.applyControlAuth(&grpcwsproto.Header{
		Operation: "auth",
		Headers: map[string]*grpcwsproto.HeaderValue{
			"Authorization": {
				Value: []string{"Bearer char-token"},
			},
		},
	})
	require.NoError(t, err)
	require.True(t, ws.authOk)
	require.Equal(t, "char-token", ws.getAuthToken())

	err = ws.applyControlAuth(&grpcwsproto.Header{
		Operation: "reauth",
		Headers:   map[string]*grpcwsproto.HeaderValue{},
	})
	require.NoError(t, err)
	require.True(t, ws.authOk)
	require.Empty(t, ws.getAuthToken())
}

func TestApplyControlAuthRejectsMissingAccountCookie(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ws := &WebsocketChannel{
		req: req,
		validateTokenFunc: func(token string) (bool, error) {
			return true, nil
		},
	}

	err := ws.applyControlAuth(&grpcwsproto.Header{
		Operation: "reauth",
		Headers:   map[string]*grpcwsproto.HeaderValue{},
	})
	require.ErrorContains(t, err, "missing authorization")
	require.False(t, ws.authOk)
	require.Empty(t, ws.getAuthToken())
}
