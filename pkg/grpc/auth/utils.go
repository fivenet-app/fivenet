package auth

import (
	"context"
	"net/http"
	"strings"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	grpc_auth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
)

func FromContext(ctx context.Context) (*pbuserinfo.UserInfo, bool) {
	c, ok := ctx.Value(userInfoCtxMarkerKey).(*pbuserinfo.UserInfo)
	return c, ok
}

func getTokenFromGRPCContext(ctx context.Context, cookieName string) (string, error) {
	// Try to get auth token from cookies and fallback to `authorization header`
	val := metadata.ExtractIncoming(ctx)
	if cookie := val.Get("cookie"); cookie != "" {
		for line := range strings.SplitSeq(cookie, "; ") {
			cs, err := http.ParseCookie(line)
			if err != nil {
				continue
			}
			if len(cs) == 0 {
				continue
			}

			if cs[0].Name == cookieName {
				return cs[0].Value, nil
			}
		}
	}

	return "", nil
}

// GetAccTokenFromGRPCContext extracts the account token from gRPC context.
func GetAccTokenFromGRPCContext(ctx context.Context) (string, error) {
	token, err := getTokenFromGRPCContext(ctx, AccCookieName)
	if err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}

	// Fallback to `authorization` header for account token
	token, err = grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return "", errorsgrpcauth.ErrInvalidToken
	}

	return token, nil
}

// GetUserTokenFromGRPCContext extracts the user token from gRPC context.
func GetUserTokenFromGRPCContext(ctx context.Context) (string, error) {
	// Use `authorization` header for user token
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return "", errorsgrpcauth.ErrInvalidToken
	}

	return token, nil
}

// GetTokensFromGRPCContext extracts both account and user tokens from gRPC context.
func GetTokensFromGRPCContext(ctx context.Context) (string, string, error) {
	accToken, err := GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return "", "", err
	}

	userToken, err := GetUserTokenFromGRPCContext(ctx)
	if err != nil {
		return "", "", err
	}

	return accToken, userToken, nil
}

func SetTokenInGRPCContext(ctx context.Context, token string) context.Context {
	md := metadata.ExtractIncoming(ctx).Clone("Authorization")
	return md.Set("Authorization", "Bearer "+token).ToIncoming(ctx)
}

func GetUserInfoFromContext(ctx context.Context) (*pbuserinfo.UserInfo, bool) {
	return FromContext(ctx)
}

func MustGetUserInfoFromContext(ctx context.Context) *pbuserinfo.UserInfo {
	userInfo, ok := FromContext(ctx)
	if !ok {
		panic(errorsgrpcauth.ErrNoUserInfo)
	}
	return userInfo
}
