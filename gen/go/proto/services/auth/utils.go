package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) setTokenCookie(ctx context.Context, token string) error {
	cookie := http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(auth.TokenExpireTime),
		Path:     "/",
		Domain:   s.domain,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	header := metadata.Pairs("set-cookie", cookie.String())
	// Send the cookie back to the client
	return grpc.SendHeader(ctx, header)
}

func (s *Server) destroyTokenCookie(ctx context.Context) error {
	cookie := http.Cookie{
		Name:     auth.CookieName,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Time{},
		MaxAge:   -1,
		Path:     "/",
		Domain:   s.domain,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	header := metadata.Pairs("set-cookie", cookie.String())
	// Send the cookie back to the client
	return grpc.SendHeader(ctx, header)
}
