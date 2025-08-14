package sync

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getAccount(
	ctx context.Context,
	identifier string,
) (*accounts.Account, *string, error) {
	tAccounts := table.FivenetAccounts.AS("account")
	// Check if an account already exists
	selectStmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.Username,
			tAccounts.RegToken.AS("reg_token"),
		).
		FROM(tAccounts).
		WHERE(jet.AND(
			tAccounts.License.EQ(jet.String(identifier)),
		)).
		LIMIT(1)

	acc := &struct {
		*accounts.Account

		RegToken *string
	}{
		Account: &accounts.Account{},
	}
	if err := selectStmt.QueryContext(ctx, s.db, acc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	return acc.Account, acc.RegToken, nil
}

func (s *Server) RegisterAccount(
	ctx context.Context,
	req *pbsync.RegisterAccountRequest,
) (*pbsync.RegisterAccountResponse, error) {
	acc, accToken, err := s.getAccount(ctx, req.GetIdentifier())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account. %w", err)
	}

	if acc == nil || acc.GetId() > 0 {
		// Account exists and no token reset has been requested
		if !req.GetResetToken() {
			return &pbsync.RegisterAccountResponse{
				Username: &acc.Username,
				RegToken: accToken,
			}, nil
		}
	}

	// Generate new token
	regToken, err := utils.GenerateRandomNumberString(6)
	if err != nil {
		return nil, fmt.Errorf("failed to generate registration token. %w", err)
	}

	var stmt jet.Statement

	tAccounts := table.FivenetAccounts

	// No account found, insert new account
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
		// Account exists, and token reset requested
		stmt = tAccounts.
			UPDATE().
			SET(
				tAccounts.Password.SET(jet.StringExp(jet.NULL)),
				tAccounts.RegToken.SET(jet.String(regToken)),
			).
			WHERE(jet.AND(
				tAccounts.ID.EQ(jet.Uint64(acc.GetId())),
				// Make sure the license is (still) the same
				tAccounts.License.EQ(jet.String(req.GetIdentifier())),
			))
	}

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to execute statement. %w", err)
	}

	resp := &pbsync.RegisterAccountResponse{
		RegToken: &regToken,
	}
	// Set account info if it is set
	if acc.GetId() != 0 {
		resp.AccountId = &acc.Id
	}
	if acc.GetUsername() != "" {
		resp.Username = &acc.Username
	}

	return resp, nil
}

func (s *Server) TransferAccount(
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

	// Delete new account
	delStmt := tAccounts.
		DELETE().
		WHERE(tAccounts.ID.EQ(jet.Uint64(acc.GetId()))).
		LIMIT(1)

	if _, err := delStmt.ExecContext(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to delete new account. %w", err)
	}

	// Update old account's license
	stmt := tAccounts.
		UPDATE(
			tAccounts.License,
		).
		SET(
			tAccounts.License.SET(jet.String(req.GetNewLicense())),
		).
		WHERE(
			tAccounts.ID.EQ(jet.Uint64(acc.GetId())),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to update old account's license. %w", err)
	}

	return resp, nil
}
