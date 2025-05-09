package auth

import (
	"context"
	"errors"
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
)

var ErrGenericAccount = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.AuthService.ErrGenericAccount"}, nil)

func (s *Server) GetAccountInfo(ctx context.Context, req *pbauth.GetAccountInfoRequest) (*pbauth.GetAccountInfoResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}

	// Load account
	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)))
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}
	if acc == nil || acc.ID == 0 {
		return nil, ErrGenericAccount
	}

	oauth2Providers := make([]*accounts.OAuth2Provider, len(s.oauth2Providers))
	for i := range oauth2Providers {
		p := s.oauth2Providers[i]
		oauth2Providers[i] = &accounts.OAuth2Provider{
			Name:     p.Name,
			Label:    p.Label,
			Homepage: p.Homepage,
			Icon:     p.Icon,
		}
	}

	oAuth2Accounts := table.FivenetOauth2Accounts.AS("oauth2account")
	stmt := oAuth2Accounts.
		SELECT(
			oAuth2Accounts.AccountID,
			oAuth2Accounts.CreatedAt,
			oAuth2Accounts.Provider.AS("oauth2account.providername"),
			oAuth2Accounts.ExternalID,
			oAuth2Accounts.Username,
			oAuth2Accounts.Avatar,
		).
		FROM(
			oAuth2Accounts,
		).
		WHERE(
			oAuth2Accounts.AccountID.EQ(jet.Uint64(acc.ID)),
		).
		LIMIT(5)

	oauth2Conns := []*accounts.OAuth2Account{}
	if err := stmt.QueryContext(ctx, s.db, &oauth2Conns); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrGenericAccount)
		}
	}
	for

	// Set provider in the connections
	i := range oauth2Conns {
		idx := slices.IndexFunc(oauth2Providers, func(p *accounts.OAuth2Provider) bool {
			return p.Name == oauth2Conns[i].GetProviderName()
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
