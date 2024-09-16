package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) getCookieBase(name string) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Now().Add(auth.TokenExpireTime),
		Domain:   s.domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *Server) setTokenCookie(ctx context.Context, token string) error {
	cookie := s.getCookieBase(auth.TokenCookieName)
	cookie.Value = token

	authedCookie := s.getCookieBase(auth.AuthedCookieName)
	authedCookie.Value = "true"
	authedCookie.HttpOnly = false

	header := metadata.Pairs("set-cookie", cookie.String(), "set-cookie", authedCookie.String())
	// Send the cookie back to the client
	return grpc.SendHeader(ctx, header)
}

func (s *Server) destroyTokenCookie(ctx context.Context) error {
	cookie := s.getCookieBase(auth.TokenCookieName)
	cookie.Expires = time.Time{}
	cookie.MaxAge = -1

	authedCookie := s.getCookieBase(auth.AuthedCookieName)
	authedCookie.Value = "false"
	authedCookie.HttpOnly = false

	header := metadata.Pairs("set-cookie", cookie.String(), "set-cookie", authedCookie.String())
	// Send the cookie back to the client
	return grpc.SendHeader(ctx, header)
}
