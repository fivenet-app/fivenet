package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	accountsoauth2 "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts/oauth2"
	jobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	authstore "github.com/fivenet-app/fivenet/v2026/stores/auth"
	"github.com/golang-jwt/jwt/v5"
	grpcmd "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type refreshAccountSessionStore struct {
	account *model.FivenetAccounts
}

func (s *refreshAccountSessionStore) GetAccountByID(
	_ context.Context,
	_ int64,
	_ bool,
) (*model.FivenetAccounts, error) {
	return s.account, nil
}

func (s *refreshAccountSessionStore) GetAccountByUsername(
	_ context.Context,
	_ string,
	_ bool,
) (*model.FivenetAccounts, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) GetLoginAccountByUsername(
	_ context.Context,
	_ string,
) (*model.FivenetAccounts, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) GetAccountByIDAndUsername(
	_ context.Context,
	_ int64,
	_ string,
	_ bool,
) (*model.FivenetAccounts, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) GetAccountByRegToken(
	_ context.Context,
	_ string,
	_ bool,
) (*model.FivenetAccounts, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) GetNewAccountByRegToken(
	_ context.Context,
	_ string,
) (*model.FivenetAccounts, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) ActivateAccount(
	_ context.Context,
	_ int64,
	_ string,
	_ string,
	_ string,
	_ string,
) error {
	return errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) UpdatePassword(_ context.Context, _ int64, _ string) error {
	return errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) UpdateUsername(_ context.Context, _ int64, _ string) error {
	return errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) ForgotPassword(_ context.Context, _ int64, _ string) error {
	return errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) ListCharacters(
	_ context.Context,
	_ int64,
	_ string,
) ([]*accounts.Character, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) GetCharacter(
	_ context.Context,
	_ int32,
) (*users.User, *jobsprops.JobProps, error) {
	return nil, nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) GetJobWithProps(
	_ context.Context,
	_ string,
) (*jobs.Job, int32, *jobsprops.JobProps, error) {
	return nil, 0, nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) ListOAuth2Connections(
	_ context.Context,
	_ int64,
) ([]*accountsoauth2.OAuth2Account, error) {
	return nil, errors.New("unexpected call")
}

func (s *refreshAccountSessionStore) DeleteSocialLogin(_ context.Context, _ int64, _ string) error {
	return errors.New("unexpected call")
}

type fakeTransportStream struct {
	headers metadata.MD
}

func (s *fakeTransportStream) Method() string {
	return "/services.auth.AuthService/RefreshAccountSession"
}

func (s *fakeTransportStream) SetHeader(md metadata.MD) error {
	if s.headers == nil {
		s.headers = metadata.MD{}
	}
	for k, v := range md {
		s.headers[k] = append(s.headers[k], v...)
	}
	return nil
}

func (s *fakeTransportStream) SendHeader(md metadata.MD) error {
	return s.SetHeader(md)
}

func (s *fakeTransportStream) SetTrailer(md metadata.MD) error {
	return s.SetHeader(md)
}

func TestRefreshAccountSession(t *testing.T) {
	t.Parallel()

	const accountID = int64(42)
	const eligibleUsername = "bootstrap-admin"
	const ineligibleUsername = "plain-user"

	now := time.Now()

	tests := []struct {
		name            string
		username        string
		license         string
		configAdminUser string
		wantConfigAdmin bool
	}{
		{
			name:            "eligible",
			username:        eligibleUsername,
			license:         eligibleUsername,
			configAdminUser: eligibleUsername,
			wantConfigAdmin: true,
		},
		{
			name:            "ineligible",
			username:        ineligibleUsername,
			license:         ineligibleUsername,
			configAdminUser: eligibleUsername,
			wantConfigAdmin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := &refreshAccountSessionStore{
				account: &model.FivenetAccounts{
					ID:       accountID,
					Username: &tt.username,
					License:  tt.license,
				},
			}

			srv := &Server{
				tm:    auth.NewTokenMgr("test-secret"),
				store: store,
				appCfg: appconfig.NewTest(
					appconfig.TestParams{LC: fxtest.NewLifecycle(t)},
				),
				configAdminUsers:  []string{tt.configAdminUser},
				configAdminGroups: nil,
				jobAdminGroups:    nil,
				jobAdminUsers:     nil,
			}

			accountProto := accounts.ConvertFromModelAcc(store.account)
			claims := auth.MapAccountToClaims(accountProto, false)
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Hour))
			claims.IssuedAt = jwt.NewNumericDate(now)
			claims.NotBefore = jwt.NewNumericDate(now)

			token, err := srv.tm.FromAccClaims(claims)
			require.NoError(t, err)

			incoming := grpcmd.MD(metadata.Pairs("cookie", auth.AccCookieName+"="+token)).
				ToIncoming(t.Context())
			stream := &fakeTransportStream{}
			ctx := grpc.NewContextWithServerTransportStream(incoming, stream)

			resp, err := srv.RefreshAccountSession(ctx, &pbauth.RefreshAccountSessionRequest{})
			require.NoError(t, err)
			require.NotNil(t, resp)

			require.NotNil(t, resp.GetExpires())
			require.Equal(t, accountID, resp.GetAccountId())
			require.Equal(t, tt.username, resp.GetUsername())
			require.Equal(t, tt.wantConfigAdmin, resp.GetCanBeConfigAdmin())

			setCookies := stream.headers.Get("set-cookie")
			require.Len(t, setCookies, 2)
			require.Contains(t, setCookies[0], auth.AuthedCookieName+"=")
			require.Contains(t, setCookies[1], auth.AccCookieName+"=")

			// The renewed account cookie should remain httpOnly and cookie-shaped.
			require.Contains(t, setCookies[1], "; HttpOnly")
		})
	}
}

var _ authstore.IStore = (*refreshAccountSessionStore)(nil)
