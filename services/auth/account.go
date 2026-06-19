package auth

import (
	"context"
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	accountsoauth2 "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts/oauth2"
	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
)

func (s *Server) GetAccountInfo(
	ctx context.Context,
	req *pbauth.GetAccountInfoRequest,
) (*pbauth.GetAccountInfoResponse, error) {
	token, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseAccToken(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	// Load account
	acc, err := s.store.GetAccountByID(ctx, claims.AccID, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}
	if acc == nil || acc.ID == 0 {
		return nil, errorsauth.ErrGenericAccount
	}

	oauth2Providers := make([]*accountsoauth2.OAuth2Provider, len(s.oauth2Providers))
	for i := range oauth2Providers {
		p := s.oauth2Providers[i]
		oauth2Providers[i] = &accountsoauth2.OAuth2Provider{
			Name:     p.Name,
			Label:    p.Label,
			Homepage: p.Homepage,
			Icon:     p.Icon,
		}
	}

	oauth2Conns, err := s.store.ListOAuth2Connections(ctx, acc.ID)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	// Set provider in the connections
	for i := range oauth2Conns {
		idx := slices.IndexFunc(oauth2Providers, func(p *accountsoauth2.OAuth2Provider) bool {
			return p.GetName() == oauth2Conns[i].GetProviderName()
		})
		if idx > -1 {
			oauth2Conns[i].Provider = oauth2Providers[idx]
		}
	}

	return &pbauth.GetAccountInfoResponse{
		Account:           accounts.ConvertFromModelAcc(acc),
		Oauth2Providers:   oauth2Providers,
		Oauth2Connections: oauth2Conns,
	}, nil
}
