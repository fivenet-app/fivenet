package oauth2

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2025/pkg/crypt"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/oauth2/providers"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type userInfoStore interface {
	storeUserInfo(ctx context.Context, accountId uint64, provider string, userInfo *providers.UserInfo) error
	updateUserInfo(ctx context.Context, accountId uint64, provider string, userInfo *providers.UserInfo) error
	getAccountInfo(ctx context.Context, provider string, userInfo *providers.UserInfo) (*model.FivenetAccounts, error)
}

type oauth2UserInfo struct {
	db    *sql.DB
	crypt *crypt.Crypt
}

func (o *oauth2UserInfo) getAccountInfo(ctx context.Context, provider string, userInfo *providers.UserInfo) (*model.FivenetAccounts, error) {
	stmt := tOauth2.
		SELECT(
			tAccs.ID,
			tAccs.Enabled,
			tAccs.Username,
			tAccs.License,
		).
		FROM(tOauth2.
			INNER_JOIN(tAccs,
				tAccs.ID.EQ(tOauth2.AccountID),
			),
		).
		WHERE(jet.AND(
			tOauth2.Provider.EQ(jet.String(provider)),
			tOauth2.ExternalID.EQ(jet.String(userInfo.ID)),
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
	accessToken, err := o.crypt.EncryptPointerString(userInfo.AccessToken)
	if err != nil {
		return err
	}
	refreshToken, err := o.crypt.EncryptPointerString(userInfo.RefreshToken)
	if err != nil {
		return err
	}

	stmt := tOauth2.
		INSERT(
			tOauth2.AccountID,
			tOauth2.Provider,
			tOauth2.ExternalID,
			tOauth2.Username,
			tOauth2.Avatar,
			tOauth2.AccessToken,
			tOauth2.RefreshToken,
			tOauth2.TokenType,
			tOauth2.Scope,
			tOauth2.ExpiresIn,
			tOauth2.ObtainedAt,
		).
		VALUES(
			accountId,
			provider,
			userInfo.ID,
			userInfo.Username,
			userInfo.Avatar,
			accessToken,
			refreshToken,
			userInfo.TokenType,
			userInfo.Scope,
			userInfo.ExpiresIn,
			userInfo.ObtainedAt,
		)

	if _, dberr := stmt.ExecContext(ctx, o.db); dberr != nil {
		if !dbutils.IsDuplicateError(dberr) {
			return dberr
		}

		// Retrieve oauth2 connection to make sure the external ID matches before updating the user info
		acc, err := accounts.RetrieveOAuth2Account(ctx, o.db, o.crypt, accountId, provider)
		if err != nil {
			return err
		}
		if acc == nil || acc.ExternalID != userInfo.ID {
			// Either no valid oauth2 connection found or the external ID doesn't match, return the db error
			// as the logic uses the "duplicate" error to detect that case
			return dberr
		}

		if err := o.updateUserInfo(ctx, accountId, provider, userInfo); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (o *oauth2UserInfo) updateUserInfo(ctx context.Context, accountId uint64, provider string, userInfo *providers.UserInfo) error {
	expiresIn := int32(0)
	if userInfo.ExpiresIn != nil {
		expiresIn = int32(*userInfo.ExpiresIn)
	}

	if err := accounts.UpdateOAuth2Account(ctx, o.db, o.crypt, accountId, &model.FivenetAccountsOauth2{
		AccountID:    accountId,
		ExternalID:   userInfo.ID,
		Username:     userInfo.Username,
		Avatar:       userInfo.Avatar,
		AccessToken:  userInfo.AccessToken,
		RefreshToken: userInfo.RefreshToken,
		TokenType:    userInfo.TokenType,
		Scope:        userInfo.Scope,
		ExpiresIn:    &expiresIn,
		ObtainedAt:   userInfo.ObtainedAt,
	}); err != nil {
		return err
	}

	return nil
}
