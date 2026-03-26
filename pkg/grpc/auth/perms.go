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

	// Short circuit for superusers, they have access to everything
	if userInfo.Superuser {
		return ctx, nil
	}

	perms := []string{perm}
	if _, ok := goproto.PermsRemap[perm]; ok {
		perms = goproto.PermsRemap[perm]
	}

	for _, p := range perms {
		if p == PermSuperuser.GetName() && userInfo.GetSuperuser() {
			return ctx, nil
		} else if p == PermAny {
			return ctx, nil
		}

		permSplit := strings.Split(p, "/")
		if len(permSplit) > 1 {
			category := pkgperms.Category(permSplit[0])
			name := pkgperms.Name(permSplit[1])

			if g.ps.Can(userInfo, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, errorsgrpcauth.ErrPermissionDenied
}
