package auth

import (
	"context"
	"strings"

	grpc_permission "github.com/galexrt/fivenet/pkg/grpc/interceptors/permission"
	"github.com/galexrt/fivenet/pkg/perms"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
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
	AuthInfoKey     struct{}
	NoTokenErr      = status.Errorf(codes.Unauthenticated, "authorization string must not be empty")
	InvalidTokenErr = status.Error(codes.Unauthenticated, "Token invalid/ expired!")
	CheckTokenErr   = status.Error(codes.Unauthenticated, "Token check failed!")
)

type GRPCAuth struct {
	tm *TokenMgr
}

func NewGRPCAuth(tm *TokenMgr) *GRPCAuth {
	return &GRPCAuth{
		tm: tm,
	}
}

func (g *GRPCAuth) GRPCAuthFunc(ctx context.Context, fullMethod string) (context.Context, error) {
	t, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, NoTokenErr
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, InvalidTokenErr
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, tInfo.Subject,
		AuthAccIDCtxTag, tInfo.CharID,
	})

	return context.WithValue(ctx, AuthInfoKey, tInfo), nil
}

type GRPCPerm struct {
	p perms.Permissions
}

func NewGRPCPerms(p perms.Permissions) *GRPCPerm {
	return &GRPCPerm{
		p: p,
	}
}

func FromContext(ctx context.Context) (*CitizenInfoClaims, bool) {
	c, ok := ctx.Value(AuthInfoKey).(*CitizenInfoClaims)
	return c, ok
}

func (g *GRPCPerm) GRPCPermissionUnaryFunc(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		claims, ok := FromContext(ctx)
		if ok {
			split := strings.Split(info.FullMethod[10:], ".")
			perm := strings.Join(split[1:], "-")

			if overrideSrv, ok := info.Server.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			permSplit := strings.Split(perm, "/")
			category := perms.Category(permSplit[0])
			name := perms.Name(permSplit[1])

			if g.p.Can(claims.CharID, claims.CharJob, claims.CharJobGrade, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func (g *GRPCPerm) GRPCPermissionStreamFunc(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		claims, ok := FromContext(ctx)
		if ok {
			split := strings.Split(info.FullMethod[10:], ".")
			perm := strings.Join(split[1:], "-")

			if overrideSrv, ok := srv.(grpc_permission.GetPermsRemapFunc); ok {
				remap := overrideSrv.GetPermsRemap()
				if _, ok := remap[perm]; ok {
					perm = remap[perm]
				}
			}

			permSplit := strings.Split(perm, "/")
			category := perms.Category(permSplit[0])
			name := perms.Name(permSplit[1])

			if g.p.Can(claims.CharID, claims.CharJob, claims.CharJobGrade, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	return grpc_auth.AuthFromMD(ctx, "bearer")
}

func GetUserIDFromContext(ctx context.Context) int32 {
	claims, ok := FromContext(ctx)
	if !ok {
		return -1
	}
	return claims.CharID
}

func GetUserInfoFromContext(ctx context.Context) (int32, string, int32) {
	claims, ok := FromContext(ctx)
	if !ok {
		return -1, "N/A", 1
	}
	return claims.CharID, claims.CharJob, claims.CharJobGrade
}
