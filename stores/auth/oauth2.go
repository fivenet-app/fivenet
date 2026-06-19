package authstore

import (
	"context"
	"errors"

	accountsoauth2 "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts/oauth2"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListOAuth2Connections(
	ctx context.Context,
	accountID int64,
) ([]*accountsoauth2.OAuth2Account, error) {
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
			tOAuth2Accounts.AccountID.EQ(mysql.Int64(accountID)),
		).
		LIMIT(5)

	oauth2Conns := []*accountsoauth2.OAuth2Account{}
	if err := stmt.QueryContext(ctx, s.db, &oauth2Conns); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return oauth2Conns, nil
}

func (s *Store) DeleteSocialLogin(ctx context.Context, accountID int64, provider string) error {
	tOAuth2Accounts := table.FivenetAccountsOauth2
	stmt := tOAuth2Accounts.
		DELETE().
		WHERE(mysql.AND(
			tOAuth2Accounts.AccountID.EQ(mysql.Int64(accountID)),
			tOAuth2Accounts.Provider.EQ(mysql.String(provider)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}
