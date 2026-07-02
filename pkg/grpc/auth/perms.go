package auth

import (
	"context"
	"strings"

	goproto "github.com/fivenet-app/fivenet/v2026/gen/go/proto"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	pkgperms "github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var PermsModule = fx.Module("grpc.perms",
	fx.Provide(
		NewGRPCPerms,
	),
)

type GRPCPerm struct {
	ps pkgperms.Permissions
}

func NewGRPCPerms(p pkgperms.Permissions) *GRPCPerm {
	return &GRPCPerm{
		ps: p,
	}
}

func (g *GRPCPerm) GRPCPermissionUnaryFunc(
	ctx context.Context,
	info *grpc.UnaryServerInfo,
) (context.Context, error) {
	return g.checkPermission(ctx, info.FullMethod)
}

func (g *GRPCPerm) GRPCPermissionStreamFunc(
	ctx context.Context,
	srv any,
	info *grpc.StreamServerInfo,
) (context.Context, error) {
	return g.checkPermission(ctx, info.FullMethod)
}

func (g *GRPCPerm) checkPermission(
	ctx context.Context,
	svcAndMethod string,
) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if !strings.HasPrefix(svcAndMethod, "/services.") {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	userInfo, ok := FromContext(ctx)
	if !ok {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	perm, found := strings.CutPrefix(svcAndMethod, "/services.")
	if !found {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	if ps, ok := goproto.PermsRemap[perm]; ok {
		for _, p := range ps {
			if p == pkgperms.PermAnyRef {
				return ctx, nil
			}
			if g.ps.Can(userInfo, p) {
				return ctx, nil
			}
		}

		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	// Keep the fast path for non-remapped RPCs, but let remapped config-admin gates
	// evaluate through the permission system instead of bypassing them here.
	if userInfo.GetJobAdmin() {
		return ctx, nil
	}

	if g.ps.CanServiceMethod(userInfo, perm) {
		return ctx, nil
	}

	return nil, errorsgrpcauth.ErrPermissionDenied
}
