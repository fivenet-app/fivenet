package auth

import (
	"context"

	"github.com/galexrt/arpanet/pkg/session"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPC struct {
}

func NewGRPC() *GRPC {
	return &GRPC{}
}

func (g *GRPC) GRPCAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := session.Tokens.ParseWithClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", tokenInfo.Subject)
	grpc_ctxtags.Extract(ctx).Set("auth.accid", tokenInfo.AccountID)
	grpc_ctxtags.Extract(ctx).Set("auth.charidx", tokenInfo.CharIndex)

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "userInfo", tokenInfo)

	return newCtx, nil
}
