package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/helpers"
	"github.com/galexrt/arpanet/query"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

var (
	Auth = &authAPI{}
)

type authAPI struct {
}

func (a *authAPI) GetUserFromContext(ctx context.Context) (*model.User, error) {
	values := grpc_ctxtags.Extract(ctx).Values()

	activeCharIdentifier := values[auth.AuthActiveCharCtxTag].(string)
	license := values[auth.AuthSubCtxTag].(string)
	if !strings.Contains(activeCharIdentifier, license) {
		return nil, fmt.Errorf("wrong license for char identifier")
	}

	return a.GetCharByIdentifier(activeCharIdentifier)
}

func (a *authAPI) GetCharByIdentifier(identifier string) (*model.User, error) {
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

func (a *authAPI) GetCharsByLicense(license string) ([]*model.User, error) {
	licenseSearch := helpers.BuildCharSearchIdentifier(license)

	users := query.User
	return users.Preload(users.UserLicenses.RelationField).
		Where(users.Identifier.Like(licenseSearch)).
		Limit(5).
		Find()
}
