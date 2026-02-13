package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/go-jet/jet/v2/mysql"
	"golang.org/x/crypto/bcrypt"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Cost parameter for bcrypt hashing.
const BCryptCost = 14

func (s *Server) getCookieBase(name string, value string) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(auth.TokenExpireTime),
		MaxAge:   int(auth.TokenExpireTime.Seconds()),
		Domain:   s.domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}

func (s *Server) setCookies(
	ctx context.Context,
	accClaims *authclaims.AccountInfoClaims,
) error {
	authedCookie := s.getCookieBase(auth.AuthedCookieName, "true")
	authedCookie.HttpOnly = false
	header := metadata.Pairs("set-cookie", authedCookie.String())

	accToken, err := s.tm.FromAccClaims(accClaims)
	if err != nil {
		return err
	}
	accCookie := s.getCookieBase(auth.AccCookieName, accToken)
	header.Append("set-cookie", accCookie.String())

	// Send the cookies back to the client
	return grpc.SendHeader(ctx, header)
}

func (s *Server) destroyCookies(ctx context.Context) error {
	authedCookie := s.getCookieBase(auth.AuthedCookieName, "false")
	authedCookie.HttpOnly = false

	accCookie := s.getCookieBase(auth.AccCookieName, "")
	accCookie.Expires = time.Time{}
	accCookie.MaxAge = -1

	tokenCookie := s.getCookieBase(auth.UserCookieName, "")
	tokenCookie.Expires = time.Time{}
	tokenCookie.MaxAge = -1

	// Send the cookies back to the client
	header := metadata.Pairs(
		"set-cookie", accCookie.String(),
		"set-cookie", authedCookie.String(),
		"set-cookie", tokenCookie.String(),
	)
	return grpc.SendHeader(ctx, header)
}

// Helper to fetch account from claims.
func (s *Server) getAccountFromIDAndUsername(
	ctx context.Context,
	accId int64, username string,
	withPassword bool,
) (*model.FivenetAccounts, error) {
	return s.getAccountFromDB(ctx, mysql.AND(
		tAccounts.ID.EQ(mysql.Int64(accId)),
		tAccounts.Username.EQ(mysql.String(username)),
	), withPassword)
}

// Helper for password hashing.
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), BCryptCost)
	return string(hashed), err
}

// Helper for password validation.
func checkPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// Helper for username normalization.
func normalizeUsername(username string) string {
	return strings.TrimSpace(username)
}
