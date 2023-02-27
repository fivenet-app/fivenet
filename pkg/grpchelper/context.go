package grpchelper

import (
	"context"
	"fmt"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/query"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	values := grpc_ctxtags.Extract(ctx).Values()

	charIndex := values[auth.AuthCharIdxCtxTag]
	license := values[auth.AuthSubCtxTag]

	charIdentifier := fmt.Sprintf("char%d:%s", charIndex, license)

	// Find user info for the new/old char index in the claim
	users := query.User
	user, err := users.Where(users.Identifier.Like(charIdentifier)).
		Preload(users.UserLicenses.Where(query.UserLicense.Owner.Eq(charIdentifier))).
		Limit(1).
		First()
	if err != nil {
		return nil, err
	}

	return user, nil
}
