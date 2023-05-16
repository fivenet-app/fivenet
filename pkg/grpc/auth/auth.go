package auth

import (
	"context"
	"strings"

	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
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
	UserInfoKey     struct{}
	NoTokenErr      = status.Errorf(codes.Unauthenticated, "authorization string must not be empty")
	InvalidTokenErr = status.Error(codes.Unauthenticated, "Token invalid/ expired!")
	CheckTokenErr   = status.Error(codes.Unauthenticated, "Token check failed!")
	NoPermsErr      = status.Error(codes.PermissionDenied, "No permissions associated with your user!")
	NoUserInfoErr   = status.Error(codes.Unauthenticated, "Something went wrong, please logout and login again!")
)

type GRPCAuth struct {
	ui userinfo.UserInfoRetriever
	tm *TokenMgr
}

func NewGRPCAuth(ui userinfo.UserInfoRetriever, tm *TokenMgr) *GRPCAuth {
	return &GRPCAuth{
		ui: ui,
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

	userInfo, err := g.ui.GetUserInfo(ctx, tInfo.CharID, tInfo.AccID)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, UserInfoKey, userInfo), nil
}

func (g *GRPCAuth) GRPCAuthFuncWithoutUserInfo(ctx context.Context, fullMethod string) (context.Context, error) {
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

func FromContext(ctx context.Context) (*userinfo.UserInfo, bool) {
	c, ok := ctx.Value(UserInfoKey).(*userinfo.UserInfo)
	return c, ok
}

func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	return grpc_auth.AuthFromMD(ctx, "bearer")
}

func GetUserInfoFromContext(ctx context.Context) (*userinfo.UserInfo, bool) {
	return FromContext(ctx)
}

func MustGetUserInfoFromContext(ctx context.Context) *userinfo.UserInfo {
	userInfo, ok := FromContext(ctx)
	if !ok {
		panic(NoUserInfoErr)
	}
	return userInfo
}

func (g *GRPCPerm) GRPCPermissionUnaryFunc(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		userInfo, ok := FromContext(ctx)
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

			if g.p.Can(userInfo, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}

func (g *GRPCPerm) GRPCPermissionStreamFunc(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Check if the method is from a service otherwise the request must be invalid
	if strings.HasPrefix(info.FullMethod, "/services") {
		userInfo, ok := FromContext(ctx)
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

			if g.p.Can(userInfo, category, name) {
				return ctx, nil
			}
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to do that! Permission: "+info.FullMethod)
}
