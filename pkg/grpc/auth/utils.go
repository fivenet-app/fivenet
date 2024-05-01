package auth

import (
	"context"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

func FromContext(ctx context.Context) (*userinfo.UserInfo, bool) {
	c, ok := ctx.Value(userInfoCtxMarkerKey).(*userinfo.UserInfo)
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
