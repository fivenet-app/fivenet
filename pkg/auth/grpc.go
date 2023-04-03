package auth

import (
	"context"
	"strings"

	grpc_permission "github.com/galexrt/fivenet/pkg/grpc/permission"
	"github.com/galexrt/fivenet/pkg/perms"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthAccIDCtxTag              = "auth.accid"
	AuthActiveCharIDCtxTag       = "auth.chrid"
	AuthActiveCharJobCtxTag      = "auth.chrjob"
	AuthActiveCharJobGradeCtxTag = "auth.chrjobg"
	AuthSubCtxTag                = "auth.sub"
)

var (
	InvalidTokenErr = status.Error(codes.Unauthenticated, "Token invalid/ expired!")
)

type GRPCAuth struct {
	tm *TokenManager
}

func NewGRPCAuth(tm *TokenManager) *GRPCAuth {
	return &GRPCAuth{
		tm: tm,
	}
}

func (g *GRPCAuth) GRPCAuthFunc(ctx context.Context, fullMethod string) (context.Context, error) {
	token, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if token == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authorization string must not be empty")
	}

	// Parse token only returns the token info when the token is still valid
	tokenInfo, err := g.tm.ParseWithClaims(token)
	if err != nil {
		return nil, InvalidTokenErr
	}

	grpc_ctxtags.Extract(ctx).Set(AuthAccIDCtxTag, tokenInfo.AccountID)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharIDCtxTag, tokenInfo.ActiveCharID)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharJobCtxTag, tokenInfo.ActiveCharJob)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharJobGradeCtxTag, tokenInfo.ActiveCharJobGrade)
	grpc_ctxtags.Extract(ctx).Set(AuthSubCtxTag, tokenInfo.Subject)

	return ctx, nil
}

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
	if strings.HasPrefix(info.FullMethod, "/services") {
		values := grpc_ctxtags.Extract(ctx).Values()
		if activeCharID, ok := values[AuthActiveCharIDCtxTag].(int32); ok {
			split := strings.Split(info.FullMethod[10:], ".")
			perm := strings.Join(split[1:], "-")

			if overrideSrv, ok := info.Server.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			if g.p.Can(activeCharID, perm) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func (g *GRPCPerm) GRPCPermissionStreamFunc(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		values := grpc_ctxtags.Extract(ctx).Values()
		if activeCharID, ok := values[AuthActiveCharIDCtxTag].(int32); ok {
			split := strings.Split(info.FullMethod[10:], ".")
			perm := strings.Join(split[1:], "-")

			if overrideSrv, ok := srv.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			if g.p.Can(activeCharID, perm) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	return grpc_auth.AuthFromMD(ctx, "bearer")
}

func MustGetTokenFromGRPCContext(ctx context.Context) string {
	token, _ := GetTokenFromGRPCContext(ctx)
	return token
}

func GetUserIDFromContext(ctx context.Context) int32 {
	values := grpc_ctxtags.Extract(ctx).Values()

	return values[AuthActiveCharIDCtxTag].(int32)
}

func GetUserInfoFromContext(ctx context.Context) (int32, string, int32) {
	values := grpc_ctxtags.Extract(ctx).Values()

	return values[AuthActiveCharIDCtxTag].(int32), values[AuthActiveCharJobCtxTag].(string), values[AuthActiveCharJobGradeCtxTag].(int32)
}
