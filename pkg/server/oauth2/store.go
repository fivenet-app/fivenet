package oauth2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	accountsoauth2 "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts/oauth2"
	"github.com/fivenet-app/fivenet/v2026/pkg/crypt"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/types"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// userInfoStore defines the interface for storing, updating, and retrieving OAuth2 user info.
type userInfoStore interface {
	storeUserInfo(
		ctx context.Context,
		accountId int64,
		provider string,
		userInfo *types.UserInfo,
	) error
	updateUserInfo(
		ctx context.Context,
		accountId int64,
		provider string,
		userInfo *types.UserInfo,
	) error
	getAccountInfo(
		ctx context.Context,
		provider string,
		userInfo *types.UserInfo,
	) (*accounts.Account, error)
}

// oauth2UserInfo implements userInfoStore for handling OAuth2 user info in the database.
type oauth2UserInfo struct {
	// db is the SQL database connection.
	db *sql.DB
	// crypt is the cryptographic utility for secure operations.
	crypt *crypt.Crypt
}

// getAccountInfo retrieves the account info for a given provider and user info.
func (o *oauth2UserInfo) getAccountInfo(
	ctx context.Context,
	provider string,
	userInfo *types.UserInfo,
) (*accounts.Account, error) {
	tAccs := tAccs.AS("account")
	stmt := tAccOauth2.
		SELECT(
			tAccs.ID,
			tAccs.Enabled,
			tAccs.Username,
			tAccs.License,
		).
		FROM(tAccOauth2.
			INNER_JOIN(tAccs,
				tAccs.ID.EQ(tAccOauth2.AccountID),
			),
		).
		WHERE(mysql.AND(
			tAccOauth2.Provider.EQ(mysql.String(provider)),
			tAccOauth2.ExternalID.EQ(mysql.String(userInfo.ID)),
			tAccs.Enabled.IS_TRUE(),
		)).
		LIMIT(1)

	var dest accounts.Account
	if err := stmt.QueryContext(ctx, o.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return &dest, nil
}

// storeUserInfo stores new OAuth2 user info or updates it if a duplicate exists.
// If a duplicate is found, it checks the external ID before updating.
func (o *oauth2UserInfo) storeUserInfo(
	ctx context.Context,
	accountId int64,
	provider string,
	userInfo *types.UserInfo,
) error {
	accessToken, err := o.crypt.EncryptPointerString(userInfo.AccessToken)
	if err != nil {
		return err
	}
	refreshToken, err := o.crypt.EncryptPointerString(userInfo.RefreshToken)
	if err != nil {
		return err
	}

	stmt := tAccOauth2.
		INSERT(
			tAccOauth2.AccountID,
			tAccOauth2.Provider,
			tAccOauth2.ExternalID,
			tAccOauth2.Username,
			tAccOauth2.Avatar,
			tAccOauth2.AccessToken,
			tAccOauth2.RefreshToken,
			tAccOauth2.TokenType,
			tAccOauth2.Scope,
			tAccOauth2.ExpiresIn,
			tAccOauth2.ObtainedAt,
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
		acc, err := accountsoauth2.RetrieveOAuth2Account(ctx, o.db, o.crypt, accountId, provider)
		if err != nil {
			return err
		}
		if acc == nil || acc.ExternalID != userInfo.ID {
			// Either no valid oauth2 connection found or the external ID doesn't match, return the db error
			// as the logic uses the "duplicate" error to detect that case
			return dberr
		}

		if err := o.updateUserInfo(ctx, accountId, provider, userInfo); err != nil {
			return fmt.Errorf("failed to update oauth2 user info. %w", err)
		}

		return nil
	}

	return nil
}

// updateUserInfo updates the OAuth2 user info for the given account and provider.
func (o *oauth2UserInfo) updateUserInfo(
	ctx context.Context,
	accountId int64,
	provider string,
	userInfo *types.UserInfo,
) error {
	expiresIn := int64(0)
	if userInfo.ExpiresIn != nil {
		expiresIn = *userInfo.ExpiresIn
	}

	if err := accountsoauth2.UpdateOAuth2Account(ctx, o.db, o.crypt, accountId, &model.FivenetAccountsOauth2{
		AccountID:    accountId,
		ExternalID:   userInfo.ID,
		Provider:     provider,
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
