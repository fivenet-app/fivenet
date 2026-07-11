package settings

import (
	"context"
	"testing"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	pkgperms "github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetAppConfigRequiresConfigAdmin(t *testing.T) {
	t.Parallel()

	server := &Server{
		appCfg: appconfig.NewTest(appconfig.TestParams{
			LC: fxtest.NewLifecycle(t),
		}),
	}

	grpcPerm := grpcauth.NewGRPCPerms(&permsstub.Permissions{
		CanFunc: func(userInfo *pbuserinfo.UserInfo, perm pkgperms.PermissionRef) bool {
			switch perm {
			case pkgperms.PermConfigAdminRef:
				return userInfo.GetCanBeConfigAdmin()
			case pkgperms.PermJobAdminRef:
				return userInfo.GetJobAdmin()
			default:
				return false
			}
		},
	})

	jobAdminCtx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		Superuser: true,
	})
	outCtx, err := grpcPerm.GRPCPermissionUnaryFunc(
		jobAdminCtx,
		&grpc.UnaryServerInfo{FullMethod: "/services.settings.ConfigService/GetAppConfig"},
	)
	require.ErrorIs(t, err, errorsgrpcauth.ErrPermissionDenied)
	assert.Nil(t, outCtx)

	configAdminCtx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		CanBeConfigAdmin: true,
	})
	outCtx, err = grpcPerm.GRPCPermissionUnaryFunc(
		configAdminCtx,
		&grpc.UnaryServerInfo{FullMethod: "/services.settings.ConfigService/GetAppConfig"},
	)
	require.NoError(t, err)
	require.NotNil(t, outCtx)

	resp, err := server.GetAppConfig(outCtx, &pbsettings.GetAppConfigRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.GetConfig())
}

func TestConfigAuthFuncOverrideRespectsAdminLevels(t *testing.T) {
	t.Parallel()

	tm := grpcauth.NewTokenMgr("test-secret")
	appCfg := appconfig.NewTest(appconfig.TestParams{
		LC: fxtest.NewLifecycle(t),
	})

	server := &Server{
		cfg: &config.Config{
			Auth: config.Auth{
				ConfigAdminGroups: []string{"config-admin"},
			},
		},
		appCfg: appCfg,
		auth: grpcauth.NewGRPCAuth(
			userinfo.NewMockUserInfoRetriever(map[int32]*pbuserinfo.UserInfo{}),
			tm,
			appCfg,
		),
		tm: tm,
	}

	grpcPerm := grpcauth.NewGRPCPerms(&permsstub.Permissions{
		CanFunc: func(userInfo *pbuserinfo.UserInfo, perm pkgperms.PermissionRef) bool {
			switch perm {
			case pkgperms.PermConfigAdminRef:
				return userInfo.GetCanBeConfigAdmin()
			case pkgperms.PermJobAdminRef:
				return userInfo.GetJobAdmin()
			default:
				return false
			}
		},
	})

	makeCtx := func(t *testing.T, accID int64, username string, groups []string) context.Context {
		t.Helper()

		token, err := tm.FromAccClaims(&authclaims.AccountInfoClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "license-1",
			},
			AccID:    accID,
			Username: username,
			Groups:   groups,
		})
		require.NoError(t, err)

		return metadata.NewIncomingContext(
			t.Context(),
			metadata.Pairs("cookie", "fivenet_acc="+token),
		)
	}

	t.Run("config-admin account can access both config rpc methods", func(t *testing.T) {
		t.Parallel()

		ctx := makeCtx(t, 42, "config-admin", []string{"config-admin"})

		for _, method := range []string{
			pbsettings.ConfigService_GetAppConfig_FullMethodName,
			pbsettings.ConfigService_UpdateAppConfig_FullMethodName,
		} {
			outCtx, err := server.AuthFuncOverride(ctx, method)
			require.NoError(t, err)
			require.NotNil(t, outCtx)

			userInfo, ok := grpcauth.GetUserInfoFromContext(outCtx)
			require.True(t, ok)
			assert.Equal(t, int64(42), userInfo.GetAccountId())
			assert.Equal(t, int32(0), userInfo.GetUserId())
			assert.False(t, userInfo.GetJobAdmin())
			assert.False(t, userInfo.GetCanBeSuperuser())
			assert.True(t, userInfo.GetCanBeConfigAdmin())

			permCtx, err := grpcPerm.GRPCPermissionUnaryFunc(
				outCtx,
				&grpc.UnaryServerInfo{FullMethod: method},
			)
			require.NoError(t, err)
			require.NotNil(t, permCtx)
		}
	})

	t.Run("job-admin-only account cannot access config rpc methods", func(t *testing.T) {
		t.Parallel()

		ctx := makeCtx(t, 43, "job-admin", []string{"job-admin"})

		for _, method := range []string{
			pbsettings.ConfigService_GetAppConfig_FullMethodName,
			pbsettings.ConfigService_UpdateAppConfig_FullMethodName,
		} {
			outCtx, err := server.AuthFuncOverride(ctx, method)
			require.NoError(t, err)
			require.NotNil(t, outCtx)

			userInfo, ok := grpcauth.GetUserInfoFromContext(outCtx)
			require.True(t, ok)
			assert.Equal(t, int64(43), userInfo.GetAccountId())
			assert.Equal(t, int32(0), userInfo.GetUserId())
			assert.False(t, userInfo.GetJobAdmin())
			assert.False(t, userInfo.GetCanBeSuperuser())
			assert.False(t, userInfo.GetCanBeConfigAdmin())

			permCtx, err := grpcPerm.GRPCPermissionUnaryFunc(
				outCtx,
				&grpc.UnaryServerInfo{FullMethod: method},
			)
			require.ErrorIs(t, err, errorsgrpcauth.ErrPermissionDenied)
			assert.Nil(t, permCtx)
		}
	})
}
