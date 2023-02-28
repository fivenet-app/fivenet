package auth

import (
	"context"
	"fmt"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/query"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthAccIDCtxTag   = "auth.accid"
	AuthCharIdxCtxTag = "auth.charidx"
	AuthSubCtxTag     = "auth.sub"
)

func GRPCAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := session.Tokens.ParseWithClaims(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set(AuthAccIDCtxTag, tokenInfo.AccountID)
	grpc_ctxtags.Extract(ctx).Set(AuthCharIdxCtxTag, tokenInfo.ActiveChar)
	grpc_ctxtags.Extract(ctx).Set(AuthSubCtxTag, tokenInfo.Subject)

	// WARNING: in production define your own type to avoid context collisions
	//newCtx := context.WithValue(ctx, "userInfo", tokenInfo)

	return ctx, nil
}

func BuildIdentifierFromLicense(activeChar int, license string) string {
	return fmt.Sprintf("char%d:%s", activeChar, license)
}

func BuildCharSearchIdentifier(license string) string {
	return fmt.Sprintf("char%%:%s", license)
}

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	values := grpc_ctxtags.Extract(ctx).Values()

	activeChar := values[AuthCharIdxCtxTag].(int)
	license := values[AuthSubCtxTag].(string)

	return GetCharByIdentifier(BuildIdentifierFromLicense(activeChar, license))
}

func GetCharByIdentifier(identifier string) (*model.User, error) {
	// Find user info for the new/old char index in the claim
	users := query.User
	user, err := users.Where(users.Identifier.Like(identifier)).
		Preload(users.UserLicenses.Where(query.UserLicense.Owner.Eq(identifier))).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetCharsByLicense(license string) ([]*model.User, error) {
	licenseSearch := BuildCharSearchIdentifier(license)

	users := query.User
	return users.Preload(users.UserLicenses.RelationField).
		Where(users.Identifier.Like(licenseSearch)).
		Limit(10).
		Find()
}
