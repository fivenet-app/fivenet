package accounts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/crypt"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/oauth2utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tOauth2 = table.FivenetAccountsOauth2

func RetrieveOAuth2Account(ctx context.Context, db qrm.Queryable, ct *crypt.Crypt, accountId uint64, provider string) (*model.FivenetAccountsOauth2, error) {
	stmt := tOauth2.
		SELECT(
			tOauth2.AccountID,
			tOauth2.CreatedAt,
			tOauth2.Provider,
			tOauth2.ExternalID,
			tOauth2.Username,
			tOauth2.Avatar,
			tOauth2.Provider,
			tOauth2.AccessToken,
			tOauth2.RefreshToken,
			tOauth2.TokenType,
			tOauth2.Scope,
			tOauth2.ExpiresIn,
			tOauth2.ObtainedAt,
		).
		FROM(tOauth2).
		WHERE(jet.AND(
			tOauth2.AccountID.EQ(jet.Uint64(accountId)),
			tOauth2.Provider.EQ(jet.String(provider)),
		)).
		LIMIT(1)

	acc := &model.FivenetAccountsOauth2{}
	if err := stmt.QueryContext(ctx, db, acc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if acc.AccountID == 0 {
		return nil, nil
	}

	var err error
	acc.AccessToken, err = ct.DecryptPointerString(acc.AccessToken)
	if err != nil {
		return nil, err
	}
	acc.RefreshToken, err = ct.DecryptPointerString(acc.RefreshToken)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func UpdateOAuth2Account(ctx context.Context, db qrm.Executable, ct *crypt.Crypt, accountId uint64, oauth2Acc *model.FivenetAccountsOauth2) error {
	accessToken, err := ct.EncryptPointerString(oauth2Acc.AccessToken)
	if err != nil {
		return err
	}
	refreshToken, err := ct.EncryptPointerString(oauth2Acc.RefreshToken)
	if err != nil {
		return err
	}

	stmt := tOauth2.
		UPDATE(
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
		SET(
			oauth2Acc.ExternalID,
			oauth2Acc.Username,
			oauth2Acc.Avatar,
			accessToken,
			refreshToken,
			oauth2Acc.TokenType,
			oauth2Acc.Scope,
			oauth2Acc.ExpiresIn,
			oauth2Acc.ObtainedAt,
		).
		WHERE(jet.AND(
			tOauth2.AccountID.EQ(jet.Uint64(accountId)),
			tOauth2.Provider.EQ(jet.String(oauth2Acc.Provider)),
			tOauth2.ExternalID.EQ(jet.String(oauth2Acc.ExternalID)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		return err
	}

	return err
}

func GetAccessToken(ctx context.Context, db qrm.Executable, cryp *crypt.Crypt, acc *model.FivenetAccountsOauth2, clientID string, clientSecret string, redirectURI string) (string, error) {
	if acc == nil || acc.RefreshToken == nil {
		return "", nil
	}

	// Check expiry (if you store obtained_at + expires_in)
	if time.Now().After(acc.ObtainedAt.Add(time.Duration(*acc.ExpiresIn) * time.Second)) {
		// Token expired, refresh
		newAT, newRT, expiresIn, err := oauth2utils.RefreshDiscordAccessToken(ctx, clientID, clientSecret, *acc.RefreshToken, redirectURI)
		if err != nil {
			return "", fmt.Errorf("refresh discord access token error. %w", err)
		}

		now := time.Now()
		acc.AccessToken = &newAT
		acc.RefreshToken = &newRT
		acc.ExpiresIn = &expiresIn
		acc.ObtainedAt = &now

		if err := UpdateOAuth2Account(ctx, db, cryp, acc.AccountID, acc); err != nil {
			return "", err
		}
	}

	return *acc.AccessToken, nil
}
