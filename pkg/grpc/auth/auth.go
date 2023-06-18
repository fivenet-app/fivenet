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

const (
	PermSuperUser = "SuperUser"
	PermAny       = "Any"
)

var UserInfoKey struct{}

var (
	ErrNoToken          = status.Errorf(codes.Unauthenticated, "errors.pkg-auth.ErrNoToken")
	ErrInvalidToken     = status.Error(codes.Unauthenticated, "errors.pkg-auth.ErrInvalidToken")
	ErrCheckToken       = status.Error(codes.Unauthenticated, "errors.pkg-auth.ErrCheckToken")
	ErrUserNoPerms      = status.Error(codes.PermissionDenied, "errors.pkg-auth.ErrUserNoPerms")
	ErrNoUserInfo       = status.Error(codes.Unauthenticated, "errors.pkg-auth.ErrNoUserInfo")
	ErrPermissionDenied = status.Errorf(codes.PermissionDenied, "errors.pkg-auth.ErrPermissionDenied")
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
		return nil, ErrNoToken
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userInfo, err := g.ui.GetUserInfo(ctx, tInfo.CharID, tInfo.AccID)
	if err != nil {
		return nil, err
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, tInfo.Subject,
		AuthAccIDCtxTag, tInfo.CharID,
	})

	return context.WithValue(ctx, UserInfoKey, userInfo), nil
}

func (g *GRPCAuth) GRPCAuthFuncWithoutUserInfo(ctx context.Context, fullMethod string) (context.Context, error) {
	t, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, ErrNoToken
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, ErrInvalidToken
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
		panic(ErrNoUserInfo)
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

			if perm == PermSuperUser && userInfo.SuperUser {
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

			if perm == PermSuperUser && userInfo.SuperUser {
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
