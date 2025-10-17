package auth

import (
	"context"
	"errors"
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
)

var ErrGenericAccount = common.NewI18nErr(
	codes.Internal,
	&common.I18NItem{Key: "errors.AuthService.ErrGenericAccount"},
	nil,
)

func (s *Server) GetAccountInfo(
	ctx context.Context,
	req *pbauth.GetAccountInfoRequest,
) (*pbauth.GetAccountInfoResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}

	// Load account
	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(mysql.Int64(claims.AccID)), false)
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

	tOAuth2Accounts := table.FivenetAccountsOauth2.AS("oauth2_account")

	stmt := tOAuth2Accounts.
		SELECT(
			tOAuth2Accounts.AccountID,
			tOAuth2Accounts.CreatedAt,
			tOAuth2Accounts.Provider.AS("oauth2_account.providername"),
			tOAuth2Accounts.ExternalID,
			tOAuth2Accounts.Username,
			tOAuth2Accounts.Avatar,
		).
		FROM(
			tOAuth2Accounts,
		).
		WHERE(
			tOAuth2Accounts.AccountID.EQ(mysql.Int64(acc.ID)),
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
