package auth

import (
	"context"

	"github.com/galexrt/arpanet/pkg/session"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthAccIDCtxTag      = "auth.accid"
	AuthActiveCharCtxTag = "auth.act_char"
	AuthSubCtxTag        = "auth.sub"
)

func GRPCAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := session.Tokens.ParseWithClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set(AuthAccIDCtxTag, tokenInfo.AccountID)
	grpc_ctxtags.Extract(ctx).Set(AuthActiveCharCtxTag, tokenInfo.ActiveChar)
	grpc_ctxtags.Extract(ctx).Set(AuthSubCtxTag, tokenInfo.Subject)

	// WARNING: in production define your own type to avoid context collisions
	//newCtx := context.WithValue(ctx, "userInfo", tokenInfo)

	return ctx, nil
}

func GetTokenFromGRPCContext(ctx context.Context) (string, error) {
	return grpc_auth.AuthFromMD(ctx, "bearer")
}

func MustGetTokenFromGRPCContext(ctx context.Context) (token string) {
	token, _ = GetTokenFromGRPCContext(ctx)
	return
}
