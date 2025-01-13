package sync

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) RegisterAccount(ctx context.Context, req *RegisterAccountRequest) (*RegisterAccountResponse, error) {
	tAccounts := table.FivenetAccounts

	// Check if an account already exists
	selectStmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.Username,
		).
		FROM(tAccounts).
		WHERE(jet.AND(
			tAccounts.License.EQ(jet.String(req.Identifier)),
		)).
		LIMIT(1)

	acc := &accounts.Account{}
	if err := selectStmt.QueryContext(ctx, s.db, acc); err != nil {
		return nil, err
	}

	if acc.Id != 0 {
		// Account exists and no token reset has been requested
		if !req.ResetToken {
			return &RegisterAccountResponse{
				Username: &acc.Username,
			}, nil
		}
	}

	// Generate new token
	regToken, err := utils.GenerateRandomNumberString(6)
	if err != nil {
		return nil, err
	}

	var stmt jet.Statement

	// No account found, insert new account
	if acc.Id == 0 {
		stmt = tAccounts.
			INSERT(
				tAccounts.Enabled,
				tAccounts.License,
				tAccounts.RegToken,
				tAccounts.LastChar,
			).
			VALUES(
				true,
				req.Identifier,
				regToken,
				req.CharId,
			)
	} else {
		// Account exists, and token reset requested
		stmt = tAccounts.
			UPDATE().
			SET(
				tAccounts.Password.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(jet.AND(
				tAccounts.ID.EQ(jet.Uint64(acc.Id)),
				// Make sure the license is (still) the same
				tAccounts.License.EQ(jet.String(req.Identifier)),
			))
	}

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	resp := &RegisterAccountResponse{
		RegToken: &regToken,
	}
	// Set account info if it is set
	if acc.Id != 0 {
		resp.AccountId = &acc.Id
	}
	if acc.Username != "" {
		resp.Username = &acc.Username
	}

	return resp, nil
}
