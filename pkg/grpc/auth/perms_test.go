package auth

import (
	"context"
	"testing"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	pkgperms "github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGRPCPermissionUnaryFuncSplitsConfigAdminAndJobAdmin(t *testing.T) {
	t.Parallel()

	grpcPerm := NewGRPCPerms(&permsstub.Permissions{
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

	t.Run("job-admin cannot bypass config-admin rpc", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(t.Context(), userInfoCtxMarkerKey, &pbuserinfo.UserInfo{
			Superuser: true,
		})

		out, err := grpcPerm.GRPCPermissionUnaryFunc(
			ctx,
			&grpc.UnaryServerInfo{FullMethod: "/services.settings.ConfigService/UpdateAppConfig"},
		)
		require.ErrorIs(t, err, errorsgrpcauth.ErrPermissionDenied)
		assert.Nil(t, out)
	})

	t.Run("config-admin can access config-admin rpc", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(t.Context(), userInfoCtxMarkerKey, &pbuserinfo.UserInfo{
			CanBeConfigAdmin: true,
		})

		out, err := grpcPerm.GRPCPermissionUnaryFunc(
			ctx,
			&grpc.UnaryServerInfo{FullMethod: "/services.settings.ConfigService/UpdateAppConfig"},
		)
		require.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("job-admin still accesses job-admin rpc", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(t.Context(), userInfoCtxMarkerKey, &pbuserinfo.UserInfo{
			Superuser: true,
		})

		out, err := grpcPerm.GRPCPermissionUnaryFunc(
			ctx,
			&grpc.UnaryServerInfo{FullMethod: "/services.settings.SystemService/GetAllPermissions"},
		)
		require.NoError(t, err)
		assert.NotNil(t, out)
	})
}

func TestGRPCPermissionUnaryFuncAllowsDocumentReferenceAndRelationWritesWithListDocuments(t *testing.T) {
	t.Parallel()

	grpcPerm := NewGRPCPerms(&permsstub.Permissions{
		CanFunc: func(_ *pbuserinfo.UserInfo, perm pkgperms.PermissionRef) bool {
			return perm == permsdocuments.DocumentsService.ListDocuments.Perm
		},
	})

	ctx := ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 1,
		Job:    "police",
	})

	methods := []string{
		"/services.documents.DocumentsService/AddDocumentReference",
		"/services.documents.DocumentsService/RemoveDocumentReference",
		"/services.documents.DocumentsService/AddDocumentRelation",
		"/services.documents.DocumentsService/RemoveDocumentRelation",
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			out, err := grpcPerm.GRPCPermissionUnaryFunc(ctx, &grpc.UnaryServerInfo{FullMethod: method})
			require.NoError(t, err)
			assert.NotNil(t, out)
		})
	}
}
