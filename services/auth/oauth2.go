package auth

import (
	"context"
	"slices"
	"strconv"

	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) DeleteSocialLogin(
	ctx context.Context,
	req *pbauth.DeleteSocialLoginRequest,
) (*pbauth.DeleteSocialLoginResponse, error) {
	if ok := s.oauth2ProviderExists(req.GetProvider()); !ok {
		auditAuthFailure(ctx, "delete_social_login", "provider_not_supported", map[string]string{
			"provider": req.GetProvider(),
		})
		return nil, errorsauth.ErrGenericAccount
	}

	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.oauth2_provider", req.GetProvider()})

	token, err := auth.GetTokenFromAuthHeaderGRPCContext(ctx)
	if err != nil {
		auditAuthFailure(ctx, "delete_social_login", "token_missing_or_invalid", map[string]string{
			"provider": req.GetProvider(),
		})
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	claims, err := s.tm.ParseAccToken(token)
	if err != nil {
		auditAuthFailure(
			ctx,
			"delete_social_login",
			"account_token_parse_failed",
			map[string]string{
				"provider": req.GetProvider(),
			},
		)
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	if err := s.store.DeleteSocialLogin(ctx, claims.AccID, req.GetProvider()); err != nil {
		auditAuthFailure(ctx, "delete_social_login", "delete_failed", map[string]string{
			"account_id": strconv.FormatInt(claims.AccID, 10),
			"provider":   req.GetProvider(),
		})
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
