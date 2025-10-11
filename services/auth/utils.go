package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	errorsauth "github.com/fivenet-app/fivenet/v2025/services/auth/errors"
	"github.com/go-jet/jet/v2/mysql"
	"golang.org/x/crypto/bcrypt"
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
		SameSite: http.SameSiteNoneMode,
	}
}

func (s *Server) setTokenCookie(ctx context.Context, token string) error {
	authedCookie := s.getCookieBase(auth.AuthedCookieName)
	authedCookie.Value = "true"
	authedCookie.HttpOnly = false

	cookie := s.getCookieBase(auth.TokenCookieName)
	cookie.Value = token

	header := metadata.Pairs("set-cookie", authedCookie.String(), "set-cookie", cookie.String())
	// Send the cookie back to the client
	return grpc.SendHeader(ctx, header)
}

func (s *Server) destroyTokenCookie(ctx context.Context) error {
	authedCookie := s.getCookieBase(auth.AuthedCookieName)
	authedCookie.Value = "false"
	authedCookie.HttpOnly = false

	cookie := s.getCookieBase(auth.TokenCookieName)
	cookie.Expires = time.Time{}
	cookie.MaxAge = -1

	header := metadata.Pairs("set-cookie", authedCookie.String(), "set-cookie", cookie.String())
	// Send the cookie back to the client
	return grpc.SendHeader(ctx, header)
}

// Helper to fetch account from claims.
func (s *Server) getAccountFromClaims(
	ctx context.Context,
	claims *auth.CitizenInfoClaims,
	withPassword bool,
) (*model.FivenetAccounts, error) {
	return s.getAccountFromDB(ctx, mysql.AND(
		tAccounts.ID.EQ(mysql.Int64(claims.AccID)),
		tAccounts.Username.EQ(mysql.String(claims.Username)),
	), withPassword)
}

// Helper for password hashing.
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
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

// Centralized superuser/override logic for character/job selection
// Returns updated jobProps (if any) and error.
func (s *Server) handleSuperuserOverride(
	ctx context.Context,
	account *model.FivenetAccounts,
	char *users.User,
	claims *auth.CitizenInfoClaims,
	isSuperuser bool,
) (*jobs.JobProps, error) {
	var jProps *jobs.JobProps

	if !isSuperuser &&
		((account.Superuser != nil && *account.Superuser) || account.OverrideJob != nil) {
		account.OverrideJob = nil
		account.OverrideJobGrade = nil

		if err := s.ui.SetUserInfo(ctx, claims.AccID, claims.CharID, false, account.OverrideJob, account.OverrideJobGrade); err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}

		not := false
		account.Superuser = &not
	} else if isSuperuser &&
		(account.Superuser != nil && *account.Superuser) && account.OverrideJob != nil && account.OverrideJobGrade != nil {
		char.Job = *account.OverrideJob
		char.JobGrade = *account.OverrideJobGrade

		s.enricher.EnrichJobInfo(char)

		_, _, jp, err := s.getJobWithProps(ctx, char.GetJob())
		if err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}
		jProps = jp
	}

	return jProps, nil
}
