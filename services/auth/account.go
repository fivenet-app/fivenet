package auth

import (
	"context"
	"slices"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	accountsoauth2 "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts/oauth2"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
)

func (s *Server) getAccountFromAccToken(
	ctx context.Context,
) (*accounts.Account, *authclaims.AccountInfoClaims, error) {
	token, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseAccToken(token)
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}

	acc, err := s.store.GetAccountByID(ctx, claims.AccID, false)
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
	}
	if acc == nil || acc.ID == 0 {
		return nil, nil, errorsauth.ErrGenericAccount
	}

	account := accounts.ConvertFromModelAcc(acc)
	return account, claims, nil
}

func (s *Server) GetAccountInfo(
	ctx context.Context,
	req *pbauth.GetAccountInfoRequest,
) (*pbauth.GetAccountInfoResponse, error) {
	account, _, err := s.getAccountFromAccToken(ctx)
	if err != nil {
		return nil, err
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

	oauth2Conns, err := s.store.ListOAuth2Connections(ctx, account.GetId())
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
		Account:           account,
		Oauth2Providers:   oauth2Providers,
		Oauth2Connections: oauth2Conns,
	}, nil
}

func (s *Server) RefreshAccountSession(
	ctx context.Context,
	req *pbauth.RefreshAccountSessionRequest,
) (*pbauth.RefreshAccountSessionResponse, error) {
	account, claims, err := s.getAccountFromAccToken(ctx)
	if err != nil {
		return nil, err
	}

	canBeConfigAdmin := s.canAccountBeConfigAdmin(account.GetGroups(), account.GetLicense())

	responseClaims := claims
	if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) <= auth.TokenRenewalTime {
		responseClaims = auth.MapAccountToClaims(account, s.canAccountBeSuperuser(account.GetGroups(), account.GetLicense()))
		if err := s.setCookies(ctx, responseClaims); err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericAccount)
		}
	}

	return &pbauth.RefreshAccountSessionResponse{
		Expires:          timestamp.New(responseClaims.ExpiresAt.Time),
		AccountId:        account.GetId(),
		Username:         account.GetUsername(),
		CanBeConfigAdmin: canBeConfigAdmin,
	}, nil
}
