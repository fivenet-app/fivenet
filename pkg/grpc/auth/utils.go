package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
)

const (
	TokenCookieName  = "fivenet_token"
	AuthedCookieName = "fivenet_authed"
)

func FromContext(ctx context.Context) (*userinfo.UserInfo, bool) {
	c, ok := ctx.Value(userInfoCtxMarkerKey).(*userinfo.UserInfo)
	return c, ok
}

func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	// Try to get auth token from cookies and fallback to `authorization header`
	val := metadata.ExtractIncoming(ctx)
	cookie := val.Get("cookie")
	if cookie != "" {
		for _, line := range strings.Split(cookie, "; ") {
			cs, err := http.ParseCookie(line)
			if err != nil {
				continue
			}
			if len(cs) == 0 {
				continue
			}

			if cs[0].Name == TokenCookieName {
				return cs[0].Value, nil
			}
		}
	}

	return grpc_auth.AuthFromMD(ctx, "bearer")
}

func SetTokenInGRPCContext(ctx context.Context, token string) context.Context {
	md := metadata.ExtractIncoming(ctx).Clone("authorization")
	return md.Set("authorization", "bearer "+token).ToIncoming(ctx)
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
