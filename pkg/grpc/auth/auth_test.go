package auth

import (
	"testing"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	grpc_metadata "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Inspired by https://github.com/grpc-ecosystem/go-grpc-middleware/blob/da1b13ec28bbdd492bdc876045791b69c4be5b81/auth/metadata_test.go
func TestGRPCAuthFunc(t *testing.T) {
	t.Parallel()

	// Valid JWT token with claims matching testUserCombinedClaim
	tm := NewTokenMgr(jwtTokenTestSecret)
	assert.NotNil(t, tm)
	token, err := tm.FromCombinedClaims(testUserCombinedClaim)
	require.NoError(t, err)
	ui := userinfo.NewMockUserInfoRetriever(map[int32]*pbuserinfo.UserInfo{
		testUserCombinedClaim.UserID: {
			AccountId: testUserCombinedClaim.AccID,
		},
	})
	appCfg := appconfig.NewTest(appconfig.TestParams{
		LC: fxtest.NewLifecycle(t),
	})

	grpcAuth := NewGRPCAuth(ui, tm, appCfg)

	for _, run := range []struct {
		md        metadata.MD
		outputNil bool
		errCode   codes.Code
		msg       string
	}{
		{
			md:        metadata.Pairs("Authorization", ""),
			outputNil: true,
			errCode:   codes.Unauthenticated,
			msg:       "authorization string must not be empty",
		},
		{
			md:        metadata.Pairs("Authorization", "invalid-jwt-token"),
			outputNil: true,
			errCode:   codes.Unauthenticated,
			msg:       "invalid auth token: ",
		},
		{
			md:        metadata.Pairs("Authorization", "Bearer "+token),
			outputNil: false,
			errCode:   codes.OK,
			msg:       "valid token",
		},
	} {
		ctx := grpc_metadata.MD(run.md).ToIncoming(t.Context())
		out, err := grpcAuth.GRPCAuthFunc(ctx, "/services.Example/GetExample")
		if run.errCode != codes.OK {
			assert.Equal(t, run.errCode, status.Code(err), run.msg)
		} else {
			require.NoError(t, err, run.msg)
		}
		if run.outputNil {
			assert.Nil(t, out, run.msg)
		} else {
			assert.NotNil(t, out, run.msg)
		}
	}
}

func TestGRPCAuthFuncWithoutUserInfoUsesLiveAccountState(t *testing.T) {
	t.Parallel()

	tm := NewTokenMgr(jwtTokenTestSecret)
	assert.NotNil(t, tm)

	token, err := tm.FromAccClaims(&authclaims.AccountInfoClaims{
		RegisteredClaims: testUserCombinedClaim.RegisteredClaims,
		AccID:            testUserCombinedClaim.AccID,
		Username:         testUserCombinedClaim.Username,
		Groups:           []string{"config-admin"},
	})
	require.NoError(t, err)

	ui := &userinfo.MockUserInfoRetriever{
		AccountInfo: map[int64]*pbuserinfo.UserInfo{
			testUserCombinedClaim.AccID: {
				AccountId: testUserCombinedClaim.AccID,
				Enabled:   true,
				License:   testUserCombinedClaim.Subject,
				Groups: &accounts.AccountGroups{
					Groups: []string{"job-admin"},
				},
				CanBeConfigAdmin: false,
			},
		},
	}

	appCfg := appconfig.NewTest(appconfig.TestParams{
		LC: fxtest.NewLifecycle(t),
	})
	grpcAuth := NewGRPCAuth(ui, tm, appCfg)

	ctx := grpc_metadata.MD(metadata.Pairs("cookie", "fivenet_acc="+token)).ToIncoming(t.Context())
	out, err := grpcAuth.GRPCAuthFuncWithoutUserInfo(
		ctx,
		"/services.settings.ConfigService/GetAppConfig",
	)
	require.NoError(t, err)
	require.NotNil(t, out)

	userInfo, ok := GetUserInfoFromContext(out)
	require.True(t, ok)
	assert.Equal(t, testUserCombinedClaim.AccID, userInfo.GetAccountId())
	assert.False(t, userInfo.GetCanBeConfigAdmin())
}

func TestGRPCAuthFuncWithoutUserInfoRejectsDisabledAccount(t *testing.T) {
	t.Parallel()

	tm := NewTokenMgr(jwtTokenTestSecret)
	assert.NotNil(t, tm)

	token, err := tm.FromAccClaims(testUserCombinedClaim.GetAccountInfoClaims())
	require.NoError(t, err)

	ui := &userinfo.MockUserInfoRetriever{
		AccountInfo: map[int64]*pbuserinfo.UserInfo{
			testUserCombinedClaim.AccID: {
				AccountId: testUserCombinedClaim.AccID,
				Enabled:   false,
				License:   testUserCombinedClaim.Subject,
				Groups: &accounts.AccountGroups{
					Groups: []string{"config-admin"},
				},
			},
		},
	}

	appCfg := appconfig.NewTest(appconfig.TestParams{
		LC: fxtest.NewLifecycle(t),
	})
	grpcAuth := NewGRPCAuth(ui, tm, appCfg)

	ctx := grpc_metadata.MD(metadata.Pairs("cookie", "fivenet_acc="+token)).ToIncoming(t.Context())
	out, err := grpcAuth.GRPCAuthFuncWithoutUserInfo(
		ctx,
		"/services.settings.ConfigService/GetAppConfig",
	)
	require.Error(t, err)
	assert.Nil(t, out)
}
