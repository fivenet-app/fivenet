package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) setTokenCookie(ctx context.Context, token string) {
	cookie := http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(auth.TokenExpireTime),
		Path:     "/",
		Domain:   s.domain,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
	}
	header := metadata.Pairs("set-cookie", cookie.String())
	// send the header back to the gateway
	grpc.SendHeader(ctx, header)
}

func (s *Server) destroyTokenCookie(ctx context.Context) {
	cookie := http.Cookie{
		Name:     auth.CookieName,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Time{},
		Path:     "/",
		Domain:   s.domain,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
	}
	header := metadata.Pairs("set-cookie", cookie.String())
	// send the header back to the gateway
	grpc.SendHeader(ctx, header)
}
