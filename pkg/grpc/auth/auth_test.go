package auth

import (
	"testing"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
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
