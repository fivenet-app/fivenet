package api

import (
	"context"
	"fmt"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/query"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

func BuildIdentifierFromLicense(activeCharIdentifier int, license string) string {
	return fmt.Sprintf("char%d:%s", activeCharIdentifier, license)
}

func BuildCharSearchIdentifier(license string) string {
	return fmt.Sprintf("char%%:%s", license)
}

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	values := grpc_ctxtags.Extract(ctx).Values()

	activeCharIdentifier := values[auth.AuthCharIdxCtxTag].(int)
	license := values[auth.AuthSubCtxTag].(string)

	return GetCharByIdentifier(BuildIdentifierFromLicense(activeCharIdentifier, license))
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
		Limit(5).
		Find()
}
