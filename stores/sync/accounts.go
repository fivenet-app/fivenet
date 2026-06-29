package syncstore

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (s *Store) getAccount(
	ctx context.Context,
	identifier string,
) (*accounts.Account, *string, error) {
	tAccounts := table.FivenetAccounts.AS("account")
	selectStmt := tAccounts.
		SELECT(tAccounts.ID, tAccounts.Username, tAccounts.RegToken.AS("reg_token")).
		FROM(tAccounts).
		WHERE(mysql.AND(tAccounts.License.EQ(mysql.String(identifier)))).
		LIMIT(1)

	acc := &struct {
		*accounts.Account

		RegToken *string
	}{Account: &accounts.Account{}}
	if err := selectStmt.QueryContext(ctx, s.db, acc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	return acc.Account, acc.RegToken, nil
}

func (s *Store) RegisterAccount(
	ctx context.Context,
	req *pbsync.RegisterAccountRequest,
) (*pbsync.RegisterAccountResponse, error) {
	acc, accToken, err := s.getAccount(ctx, req.GetIdentifier())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account. %w", err)
	}

	if acc == nil || acc.GetId() > 0 {
		if !req.GetResetToken() {
			return &pbsync.RegisterAccountResponse{Username: &acc.Username, RegToken: accToken}, nil
		}
	}

	regToken, err := utils.GenerateRandomNumberString(6)
	if err != nil {
		return nil, fmt.Errorf("failed to generate registration token. %w", err)
	}

	tAccounts := table.FivenetAccounts
	var stmt mysql.Statement
	if acc.GetId() == 0 {
		stmt = tAccounts.
			INSERT(
				tAccounts.Enabled,
				tAccounts.License,
				tAccounts.RegToken,
				tAccounts.LastChar,
			).
			VALUES(
				true,
				req.GetIdentifier(),
				regToken,
				req.LastCharId,
			)
	} else {
		stmt = tAccounts.
			UPDATE().
			SET(
				tAccounts.Password.SET(mysql.StringExp(mysql.NULL)),
				tAccounts.RegToken.SET(mysql.String(regToken)),
			).
			WHERE(mysql.AND(
				tAccounts.ID.EQ(mysql.Int64(acc.GetId())),
				tAccounts.License.EQ(mysql.String(req.GetIdentifier())),
			)).
			LIMIT(1)
	}

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to execute statement. %w", err)
	}

	resp := &pbsync.RegisterAccountResponse{RegToken: &regToken}
	if acc.GetId() != 0 {
		resp.AccountId = &acc.Id
	}
	if acc.GetUsername() != "" {
		resp.Username = &acc.Username
	}

	return resp, nil
}

func (s *Store) TransferAccount(
	ctx context.Context,
	req *pbsync.TransferAccountRequest,
) (*pbsync.TransferAccountResponse, error) {
	acc, _, err := s.getAccount(ctx, req.GetOldLicense())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account with old license. %w", err)
	}

	resp := &pbsync.TransferAccountResponse{}
	if acc.GetId() == 0 {
		return resp, nil
	}

	tAccounts := table.FivenetAccounts
	delStmt := tAccounts.
		DELETE().
		WHERE(tAccounts.ID.EQ(mysql.Int64(acc.GetId()))).
		LIMIT(1)
	if _, err := delStmt.ExecContext(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to delete new account. %w", err)
	}

	stmt := tAccounts.
		UPDATE(tAccounts.License).
		SET(tAccounts.License.SET(mysql.String(req.GetNewLicense()))).
		WHERE(tAccounts.ID.EQ(mysql.Int64(acc.GetId()))).
		LIMIT(1)
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to update old account's license. %w", err)
	}

	return resp, nil
}

func (s *Store) AddAccountUpdate(
	ctx context.Context,
	req *pbsync.AddAccountUpdateRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleAccountUpdate(ctx, req.GetAccountUpdate()); err != nil {
		return nil, fmt.Errorf("failed to handle account update. %w", err)
	}

	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddUserOAuth2Conn(
	ctx context.Context,
	req *pbsync.AddUserOAuth2ConnRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleUserOauth2(ctx, req.GetUserOauth2()); err != nil {
		return nil, fmt.Errorf("failed to handle UserOauth2 activity. %w", err)
	}

	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) handleAccountUpdate(ctx context.Context, data *activity.AccountUpdate) error {
	tAccounts := table.FivenetAccounts

	var groups *accounts.AccountGroups
	if data.GetGroups() != nil {
		if data.GetGroups() != nil && len(data.GetGroups().GetGroups()) > 0 {
			groups = data.GetGroups()
		} else if data.GetGroup() != "" {
			groups = &accounts.AccountGroups{Groups: []string{data.GetGroup()}}
		}
	}

	stmt := tAccounts.
		UPDATE(tAccounts.Groups).
		SET(groups).
		WHERE(tAccounts.License.EQ(mysql.String(data.GetLicense()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return fmt.Errorf("failed to update account. %w", err)
	}

	return nil
}

func (s *Store) handleUserOauth2(ctx context.Context, data *activity.UserOAuth2Conn) error {
	idx := slices.IndexFunc(s.cfg.OAuth2.Providers, func(in *config.OAuth2Provider) bool {
		return in.Name == data.GetProviderName()
	})
	if idx == -1 {
		return fmt.Errorf("invalid provider name. %s", data.GetProviderName())
	}

	provider := s.cfg.OAuth2.Providers[idx]
	tAccounts := table.FivenetAccounts

	type Account struct{ ID int64 }
	var account Account
	stmt := tAccounts.
		SELECT(tAccounts.ID).
		FROM(tAccounts).
		WHERE(tAccounts.License.EQ(mysql.String(data.GetIdentifier()))).
		LIMIT(1)
	if err := stmt.QueryContext(ctx, s.db, &account); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query account by identifier. %w", err)
		}
	}

	if account.ID == 0 {
		s.logger.Warn("no fivenet account found for identifier in user oauth2 sync connect",
			zap.String("provider", data.GetProviderName()),
			zap.String("identifier", data.GetIdentifier()),
			zap.String("external_id", data.GetExternalId()),
		)
		return nil
	}

	tOAuth2Accs := table.FivenetAccountsOauth2
	insertStmt := tOAuth2Accs.
		INSERT(
			tOAuth2Accs.AccountID,
			tOAuth2Accs.Provider,
			tOAuth2Accs.ExternalID,
			tOAuth2Accs.Username,
			tOAuth2Accs.Avatar,
		).
		VALUES(
			account.ID,
			provider.Name,
			data.GetExternalId(),
			data.GetUsername(),
			provider.DefaultAvatar,
		)
	if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to insert OAuth2 account. %w", err)
		}
	}

	return nil
}
