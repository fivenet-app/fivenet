package settingsstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tAccounts = table.FivenetAccounts.AS("account")
	tOauth2   = table.FivenetAccountsOauth2.AS("oauth2account")
)

func (s *Store) ListAccounts(
	ctx context.Context,
	opts ListAccountsOptions,
) (*pbsettings.ListAccountsResponse, error) {
	var t mysql.ReadableTable = tAccounts

	condition := mysql.Bool(true)
	if opts.License != "" {
		condition = condition.AND(
			tAccounts.License.LIKE(mysql.String(fmt.Sprintf("%%%s%%", opts.License))),
		)
	}
	if opts.OnlyDisabled {
		condition = condition.AND(tAccounts.Enabled.EQ(mysql.Bool(false)))
	}
	if opts.Username != "" {
		condition = condition.AND(
			tAccounts.Username.LIKE(mysql.String(fmt.Sprintf("%%%s%%", opts.Username))),
		)
	}
	if opts.ExternalID != "" {
		condition = condition.AND(
			tOauth2.ExternalID.LIKE(mysql.String(fmt.Sprintf("%%%s%%", opts.ExternalID))),
		)
		t = t.
			INNER_JOIN(tOauth2,
				tOauth2.AccountID.EQ(tAccounts.ID),
			)
	}
	if opts.Group != "" {
		condition = condition.AND(
			dbutils.JSON_CONTAINS(tAccounts.Groups, mysql.String(opts.Group)),
		)
	}

	countStmt := tAccounts.
		SELECT(
			mysql.COUNT(tAccounts.ID).AS("data_count.total"),
		).
		FROM(t).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag, limit := opts.Pagination.GetResponseWithPageSize(count.Total, 30)
	resp := &pbsettings.ListAccountsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	orderBys := s.accountSorter.Build(opts.Sort)

	var accountIDs []int64
	idStmt := tAccounts.
		SELECT(
			tAccounts.ID,
		).
		FROM(t).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(opts.Pagination.GetOffset()).
		LIMIT(limit)

	if err := idStmt.QueryContext(ctx, s.db, &accountIDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if len(accountIDs) == 0 {
		return resp, nil
	}

	ids := make([]mysql.Expression, len(accountIDs))
	for i, id := range accountIDs {
		ids[i] = mysql.Int64(id)
	}

	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.CreatedAt,
			tAccounts.UpdatedAt,
			tAccounts.DeletedAt,
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.License,
			tAccounts.Groups,
			tAccounts.LastChar,
			tOauth2.AccountID,
			tOauth2.CreatedAt,
			tOauth2.Provider.AS("oauth2account.providername"),
			tOauth2.ExternalID,
			tOauth2.Username,
			tOauth2.Avatar,
		).
		FROM(
			tAccounts.
				LEFT_JOIN(tOauth2,
					tOauth2.AccountID.EQ(tAccounts.ID),
				),
		).
		WHERE(tAccounts.ID.IN(ids...))

	if err := stmt.QueryContext(ctx, s.db, &resp.Accounts); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Store) getAccount(
	ctx context.Context,
	condition mysql.BoolExpression,
	withPass bool,
	includeDeleted bool,
) (*accounts.Account, error) {
	columns := mysql.ProjectionList{
		tAccounts.ID,
		tAccounts.CreatedAt,
		tAccounts.UpdatedAt,
		tAccounts.DeletedAt,
		tAccounts.Enabled,
		tAccounts.Username,
		tAccounts.License,
		tAccounts.RegToken,
		tAccounts.Groups,
		tAccounts.LastChar,
	}
	if withPass {
		columns = append(columns, tAccounts.Password)
	}

	stmt := tAccounts.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(tAccounts).
		WHERE(mysql.AND(
			tAccounts.Enabled.IS_TRUE(),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tAccounts.DeletedAt.IS_NULL(),
			),
			condition,
		)).
		LIMIT(1)

	acc := &accounts.Account{}
	if err := stmt.QueryContext(ctx, s.db, acc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if acc.GetId() == 0 {
		return nil, nil
	}

	return acc, nil
}

func (s *Store) GetAccountByID(
	ctx context.Context,
	accountID int64,
	includeDeleted bool,
) (*accounts.Account, error) {
	tAccounts := table.FivenetAccounts.AS("account")
	return s.getAccount(ctx, tAccounts.ID.EQ(mysql.Int64(accountID)), false, includeDeleted)
}

func (s *Store) UpdateAccount(
	ctx context.Context,
	req *pbsettings.UpdateAccountRequest,
) (*pbsettings.UpdateAccountResponse, error) {
	updateSets := []any{}

	tAccounts := table.FivenetAccounts
	if req.Enabled != nil {
		updateSets = append(updateSets, tAccounts.Enabled.SET(mysql.Bool(req.GetEnabled())))
	}

	if req.LastChar != nil && req.GetLastChar() > 0 {
		updateSets = append(updateSets, tAccounts.LastChar.SET(mysql.Int32(req.GetLastChar())))
	}

	if len(updateSets) > 0 {
		stmt := tAccounts.
			UPDATE()
		if len(updateSets) == 1 {
			stmt = stmt.SET(updateSets[0])
		} else {
			stmt = stmt.SET(updateSets[0], updateSets[1:]...)
		}

		stmt = stmt.
			WHERE(
				tAccounts.ID.EQ(mysql.Int64(req.GetId())),
			).
			LIMIT(1)
		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}
	}

	acc, err := s.GetAccountByID(ctx, req.GetId(), false)
	if err != nil {
		return nil, err
	}

	return &pbsettings.UpdateAccountResponse{Account: acc}, nil
}

func (s *Store) DisconnectSocialLogin(ctx context.Context, accountID int64, provider string) error {
	tOauth2 := table.FivenetAccountsOauth2
	stmt := tOauth2.
		DELETE().
		WHERE(mysql.AND(
			tOauth2.AccountID.EQ(mysql.Int64(accountID)),
			tOauth2.Provider.EQ(mysql.String(provider)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) DeleteAccount(
	ctx context.Context,
	accountID int64,
	deletedAtTime *timestamp.Timestamp,
) (*pbsettings.DeleteAccountResponse, error) {
	tAccounts := table.FivenetAccounts
	stmt := tAccounts.
		UPDATE().
		SET(
			tAccounts.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(tAccounts.ID.EQ(mysql.Int64(accountID))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &pbsettings.DeleteAccountResponse{DeletedAt: deletedAtTime}, nil
}
