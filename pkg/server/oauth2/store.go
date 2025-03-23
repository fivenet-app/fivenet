package oauth2

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/pkg/server/oauth2/providers"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type userInfoStore interface {
	storeUserInfo(ctx context.Context, accountId uint64, provider string, userInfo *providers.UserInfo) error
	getUserInfo(ctx context.Context, provider string, userInfo *providers.UserInfo) (*model.FivenetAccounts, error)
}

type oauth2UserInfo struct {
	db *sql.DB
}

func (o *oauth2UserInfo) getUserInfo(ctx context.Context, provider string, userInfo *providers.UserInfo) (*model.FivenetAccounts, error) {
	stmt := tOAuthAccs.
		SELECT(
			tAccs.ID,
			tAccs.Username,
			tAccs.License,
		).
		FROM(tOAuthAccs.
			INNER_JOIN(tAccs,
				tAccs.ID.EQ(tOAuthAccs.AccountID),
			),
		).
		WHERE(jet.AND(
			tOAuthAccs.Provider.EQ(jet.String(provider)),
			tOAuthAccs.ExternalID.EQ(jet.String(userInfo.ID)),
			tAccs.Enabled.IS_TRUE(),
		)).
		LIMIT(1)

	var dest model.FivenetAccounts
	if err := stmt.QueryContext(ctx, o.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
		return nil, nil
	}

	return &dest, nil
}

func (o *oauth2UserInfo) storeUserInfo(ctx context.Context, accountId uint64, provider string, userInfo *providers.UserInfo) error {
	stmt := tOAuthAccs.
		INSERT(
			tOAuthAccs.AccountID,
			tOAuthAccs.Provider,
			tOAuthAccs.ExternalID,
			tOAuthAccs.Username,
			tOAuthAccs.Avatar,
		).
		VALUES(
			accountId,
			provider,
			userInfo.ID,
			userInfo.Username,
			userInfo.Avatar,
		)

	if _, err := stmt.ExecContext(ctx, o.db); err != nil {
		return err
	}

	return nil
}
