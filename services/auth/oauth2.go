package auth

import (
	"context"
	"slices"

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
		return nil, errorsauth.ErrGenericAccount
	}

	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.oauth2_provider", req.GetProvider()})

	token, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	claims, err := s.tm.ParseAccToken(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	if err := s.store.DeleteSocialLogin(ctx, claims.AccID, req.GetProvider()); err != nil {
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
