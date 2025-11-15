package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tAccounts = table.FivenetAccounts.AS("account")
	tOauth2   = table.FivenetAccountsOauth2.AS("oauth2account")
)

func (s *Server) ListAccounts(
	ctx context.Context,
	req *pbsettings.ListAccountsRequest,
) (*pbsettings.ListAccountsResponse, error) {
	var t mysql.ReadableTable = tAccounts
	condition := mysql.Bool(true)
	if req.License != nil && req.GetLicense() != "" {
		condition = condition.AND(
			tAccounts.License.LIKE(mysql.String(fmt.Sprintf("%%%s%%", req.GetLicense()))),
		)
	}
	if req.Enabled != nil {
		condition = condition.AND(tAccounts.Enabled.EQ(mysql.Bool(req.GetEnabled())))
	}
	if req.Username != nil && req.GetUsername() != "" {
		condition = condition.AND(
			tAccounts.Username.LIKE(mysql.String(fmt.Sprintf("%%%s%%", req.GetUsername()))),
		)
	}
	if req.ExternalId != nil && req.GetExternalId() != "" {
		condition = condition.AND(
			tOauth2.ExternalID.LIKE(mysql.String(fmt.Sprintf("%%%s%%", req.GetExternalId()))),
		)
		t = t.INNER_JOIN(tOauth2,
			tOauth2.AccountID.EQ(tAccounts.ID),
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
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 30)
	resp := &pbsettings.ListAccountsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "license":
			column = tAccounts.License
		case "username":
			fallthrough
		case "id":
			fallthrough
		default:
			column = tAccounts.ID
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tAccounts.CreatedAt.DESC())
	}

	// First, fetch the distinct account IDs for the current page
	var accountIDs []int64
	idStmt := tAccounts.
		SELECT(
			tAccounts.ID,
		).
		FROM(t).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := idStmt.QueryContext(ctx, s.db, &accountIDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}
	if len(accountIDs) == 0 {
		return resp, nil
	}

	ids := make([]mysql.Expression, len(accountIDs))
	for i, id := range accountIDs {
		ids[i] = mysql.Int64(id)
	}

	// Now, fetch all accounts and their oauth2 connections for these IDs
	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.CreatedAt,
			tAccounts.UpdatedAt,
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.License,
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
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) getAccount(ctx context.Context, id int64) (*accounts.Account, error) {
	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.CreatedAt,
			tAccounts.UpdatedAt,
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.License,
			tAccounts.LastChar,
		).
		FROM(tAccounts).
		WHERE(
			tAccounts.ID.EQ(mysql.Int64(id)),
		)

	var account accounts.Account
	if err := stmt.QueryContext(ctx, s.db, &account); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	if account.GetId() == 0 {
		return nil, nil
	}

	return &account, nil
}

func (s *Server) UpdateAccount(
	ctx context.Context,
	req *pbsettings.UpdateAccountRequest,
) (*pbsettings.UpdateAccountResponse, error) {
	tAccounts := table.FivenetAccounts

	updateSets := []interface{}{}

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
			)
		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	acc, err := s.getAccount(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbsettings.UpdateAccountResponse{
		Account: acc,
	}, nil
}

func (s *Server) DisconnectSocialLogin(
	ctx context.Context,
	req *pbsettings.DisconnectSocialLoginRequest,
) (*pbsettings.DisconnectSocialLoginResponse, error) {
	tOauth2 := table.FivenetAccountsOauth2

	stmt := tOauth2.
		DELETE().
		WHERE(mysql.AND(
			tOauth2.AccountID.EQ(mysql.Int64(req.GetId())),
			tOauth2.Provider.EQ(mysql.String(req.GetProviderName())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return nil, nil
}

func (s *Server) DeleteAccount(
	ctx context.Context,
	req *pbsettings.DeleteAccountRequest,
) (*pbsettings.DeleteAccountResponse, error) {
	tAccounts := table.FivenetAccounts

	stmt := tAccounts.
		DELETE().
		WHERE(tAccounts.ID.EQ(mysql.Int64(req.GetId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbsettings.DeleteAccountResponse{}, nil
}
