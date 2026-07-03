package settings

import (
	"testing"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	pkgperms "github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
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
