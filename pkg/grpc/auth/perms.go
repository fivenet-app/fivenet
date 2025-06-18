package auth

import (
	"context"
	"strings"

	grpc_permission "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/permission"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"google.golang.org/grpc"
)

type GRPCPerm struct {
	p perms.Permissions
}

func NewGRPCPerms(p perms.Permissions) *GRPCPerm {
	return &GRPCPerm{
		p: p,
	}
}

func (g *GRPCPerm) GRPCPermissionUnaryFunc(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services.") {
		userInfo, ok := FromContext(ctx)
		if ok {
			perm, found := strings.CutPrefix(info.FullMethod, "/services.")
			if !found {
				return nil, ErrPermissionDenied
			}

			if overrideSrv, ok := info.Server.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			if perm == PermSuperuser.Name && userInfo.Superuser {
				return ctx, nil
			} else if perm == PermAny {
				return ctx, nil
			}

			permSplit := strings.Split(perm, "/")
			if len(permSplit) > 1 {
				category := perms.Category(permSplit[0])
				name := perms.Name(permSplit[1])

				if g.p.Can(userInfo, category, name) {
					return ctx, nil
				}
			}
		}
	}

	return nil, ErrPermissionDenied
}

func (g *GRPCPerm) GRPCPermissionStreamFunc(ctx context.Context, srv any, info *grpc.StreamServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services.") {
		userInfo, ok := FromContext(ctx)
		if ok {
			perm, found := strings.CutPrefix(info.FullMethod, "/services.")
			if !found {
				return nil, ErrPermissionDenied
			}

			if overrideSrv, ok := srv.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			if perm == PermSuperuser.Name && userInfo.Superuser {
				return ctx, nil
			} else if perm == PermAny {
				return ctx, nil
			}

			permSplit := strings.Split(perm, "/")
			if len(permSplit) > 1 {
				category := perms.Category(permSplit[0])
				name := perms.Name(permSplit[1])

				if g.p.Can(userInfo, category, name) {
					return ctx, nil
				}
			}
		}
	}

	return nil, ErrPermissionDenied
}
