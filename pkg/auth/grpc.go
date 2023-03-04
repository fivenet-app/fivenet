package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/query"
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

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	values := grpc_ctxtags.Extract(ctx).Values()

	activeCharIdentifier := values[AuthActiveCharCtxTag].(string)
	license := values[AuthSubCtxTag].(string)
	if !strings.Contains(activeCharIdentifier, license) {
		return nil, fmt.Errorf("wrong license for char identifier")
	}

	return getCharByIdentifier(activeCharIdentifier)
}

func getCharByIdentifier(identifier string) (*model.User, error) {
	// Find user info for the new/old char index in the claim
	u := query.User
	user, err := u.Where(u.Identifier.Like(identifier)).
		Preload(u.UserLicenses.Where(query.UserLicense.Owner.Eq(identifier))).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CanUser(user *model.User, perm string, field ...string) bool {
	return CanUserID(uint(user.ID), perm)
}

func CanUserAccessField(user *model.User, perm string, field string) bool {
	return CanUserID(uint(user.ID), perm+"."+field)
}

func CanUserID(userID uint, perm string) bool {
	can, err := query.Perms.UserHasPermission(userID, perm)
	if err != nil {
		return false
	}

	return can
}
