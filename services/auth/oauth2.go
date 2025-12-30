package auth

import (
	"context"
	"slices"

	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsauth "github.com/fivenet-app/fivenet/v2025/services/auth/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tOAuth2Accounts = table.FivenetAccountsOauth2

func (s *Server) DeleteSocialLogin(
	ctx context.Context,
	req *pbauth.DeleteSocialLoginRequest,
) (*pbauth.DeleteSocialLoginResponse, error) {
	if ok := s.oauth2ProviderExists(req.GetProvider()); !ok {
		return nil, errorsauth.ErrGenericAccount
	}

	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.oauth2_provider", req.GetProvider()})

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	stmt := tOAuth2Accounts.
		DELETE().
		WHERE(mysql.AND(
			tOAuth2Accounts.AccountID.EQ(mysql.Int64(claims.AccID)),
			tOAuth2Accounts.Provider.EQ(mysql.String(req.GetProvider())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	return &pbauth.DeleteSocialLoginResponse{
		Success: true,
	}, nil
}

func (s *Server) oauth2ProviderExists(name string) bool {
	return slices.ContainsFunc(s.oauth2Providers, func(p *config.OAuth2Provider) bool {
		return p.Name == name
	})
}
