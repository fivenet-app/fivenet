package auth

import (
	"context"
	"errors"
	"testing"

	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/proto"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func newAuthErrorTestServer(t *testing.T, store *refreshAccountSessionStore) *Server {
	t.Helper()

	return &Server{
		logger: zap.NewNop(),
		tm:     auth.NewTokenMgr("test-secret"),
		appCfg: appconfig.NewTest(appconfig.TestParams{LC: fxtest.NewLifecycle(t)}),
		store:  store,
	}
}

func newIncomingAuthCtx(token string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"Authorization", "Bearer "+token,
	))
}

func newAccountToken(t *testing.T, tm *auth.TokenMgr, accountID int64, username string) string {
	t.Helper()

	token, err := tm.FromAccClaims(&authclaims.AccountInfoClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "fivenet",
			Subject:  "license",
			Audience: []string{"fivenet"},
		},
		AccID:    accountID,
		Username: username,
	})
	require.NoError(t, err)

	return token
}

func TestLoginReturnsInvalidLoginForMissingPassword(t *testing.T) {
	t.Parallel()

	store := &refreshAccountSessionStore{
		getLoginAccountByUsernameFn: func(_ context.Context, _ string) (*model.FivenetAccounts, error) {
			username := "user"
			return &model.FivenetAccounts{
				ID:       1,
				Username: &username,
			}, nil
		},
	}
	srv := newAuthErrorTestServer(t, store)

	resp, err := srv.Login(t.Context(), &pbauth.LoginRequest{
		Username: "user",
		Password: "password",
	})
	require.Nil(t, resp)
	proto.CompareGRPCError(t, errorsauth.ErrInvalidLogin, err)
}

func TestCreateAccountReturnsGenericFailureForUsedToken(t *testing.T) {
	t.Parallel()

	store := &refreshAccountSessionStore{
		getNewAccountByRegTokenFn: func(_ context.Context, _ string) (*model.FivenetAccounts, error) {
			username := "existing-user"
			password := "hashed"
			return &model.FivenetAccounts{
				ID:       2,
				Username: &username,
				Password: &password,
			}, nil
		},
	}
	srv := newAuthErrorTestServer(t, store)

	resp, err := srv.CreateAccount(t.Context(), &pbauth.CreateAccountRequest{
		RegToken: "reg-token",
		Username: "new-user",
		Password: "password",
	})
	require.Nil(t, resp)
	proto.CompareGRPCError(t, errorsauth.ErrGenericAccount, err)
}

func TestChangePasswordReturnsCollapsedStateError(t *testing.T) {
	t.Parallel()

	hashedPassword, err := hashPassword("old")
	require.NoError(t, err)

	store := &refreshAccountSessionStore{
		getAccountByIDAndUsernameFn: func(_ context.Context, _ int64, _ string, _ bool) (*model.FivenetAccounts, error) {
			username := "user"
			return &model.FivenetAccounts{
				ID:       3,
				Username: &username,
				Password: &hashedPassword,
			}, nil
		},
		updatePasswordFn: func(_ context.Context, _ int64, _ string) error {
			return errors.New("boom")
		},
	}
	srv := newAuthErrorTestServer(t, store)
	token := newAccountToken(t, srv.tm, 3, "user")
	ctx := newIncomingAuthCtx(token)

	resp, err := srv.ChangePassword(ctx, &pbauth.ChangePasswordRequest{
		CurrentPassword: "old",
		NewPassword:     "new",
	})
	require.Nil(t, resp)
	proto.CompareGRPCError(t, errorsauth.ErrGenericAccount, err)
}

func TestChangeUsernameReturnsCollapsedStateError(t *testing.T) {
	t.Parallel()

	hashedPassword, err := hashPassword("password")
	require.NoError(t, err)

	store := &refreshAccountSessionStore{
		getAccountByIDAndUsernameFn: func(_ context.Context, _ int64, _ string, _ bool) (*model.FivenetAccounts, error) {
			username := "user"
			return &model.FivenetAccounts{
				ID:       4,
				Username: &username,
				Password: &hashedPassword,
			}, nil
		},
		getAccountByUsernameFn: func(_ context.Context, _ string, _ bool) (*model.FivenetAccounts, error) {
			return nil, qrm.ErrNoRows
		},
		updateUsernameFn: func(_ context.Context, _ int64, _ string) error {
			return errors.New("boom")
		},
	}
	srv := newAuthErrorTestServer(t, store)
	token := newAccountToken(t, srv.tm, 4, "user")
	ctx := newIncomingAuthCtx(token)

	resp, err := srv.ChangeUsername(ctx, &pbauth.ChangeUsernameRequest{
		CurrentUsername: "user",
		NewUsername:     "new-user",
	})
	require.Nil(t, resp)
	proto.CompareGRPCError(t, errorsauth.ErrGenericAccount, err)
}

func TestForgotPasswordRequiresUsername(t *testing.T) {
	t.Parallel()

	store := &refreshAccountSessionStore{
		getAccountByRegTokenFn: func(_ context.Context, _ string, _ bool) (*model.FivenetAccounts, error) {
			password := "hashed"
			return &model.FivenetAccounts{
				ID:       5,
				Password: &password,
			}, nil
		},
	}
	srv := newAuthErrorTestServer(t, store)

	resp, err := srv.ForgotPassword(t.Context(), &pbauth.ForgotPasswordRequest{
		RegToken: "reg-token",
		New:      "new-password",
	})
	require.Nil(t, resp)
	proto.CompareGRPCError(t, errorsauth.ErrForgotPassword, err)
}
