package auth

import (
	"context"
	"strings"

	goproto "github.com/fivenet-app/fivenet/v2025/gen/go/proto"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	pkgperms "github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"google.golang.org/grpc"
)

type GRPCPerm struct {
	p pkgperms.Permissions
}

func NewGRPCPerms(p pkgperms.Permissions) *GRPCPerm {
	return &GRPCPerm{
		p: p,
	}
}

func (g *GRPCPerm) GRPCPermissionUnaryFunc(
	ctx context.Context,
	info *grpc.UnaryServerInfo,
) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if !strings.HasPrefix(info.FullMethod, "/services.") {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	userInfo, ok := FromContext(ctx)
	if !ok {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	perm, found := strings.CutPrefix(info.FullMethod, "/services.")
	if !found {
		return nil, errorsgrpcauth.ErrPermissionDenied
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

			if g.p.Can(userInfo, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, errorsgrpcauth.ErrPermissionDenied
}

func (g *GRPCPerm) GRPCPermissionStreamFunc(
	ctx context.Context,
	srv any,
	info *grpc.StreamServerInfo,
) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if !strings.HasPrefix(info.FullMethod, "/services.") {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	userInfo, ok := FromContext(ctx)
	if !ok {
		return nil, errorsgrpcauth.ErrPermissionDenied
	}

	perm, found := strings.CutPrefix(info.FullMethod, "/services.")
	if !found {
		return nil, errorsgrpcauth.ErrPermissionDenied
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

			if g.p.Can(userInfo, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, errorsgrpcauth.ErrPermissionDenied
}
