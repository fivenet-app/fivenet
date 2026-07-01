package syncstore

import (
	"context"
	"fmt"

	citizenslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	syncactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) SendData(
	ctx context.Context,
	req *pbsync.SendDataRequest,
) (*pbsync.SendDataResponse, error) {
	resp := &pbsync.SendDataResponse{RowsAffected: 0}

	var err error
	switch d := req.GetData().(type) {
	case *pbsync.SendDataRequest_Jobs:
		if resp.RowsAffected, err = s.handleJobsData(ctx, d.Jobs.GetJobs()); err != nil {
			return nil, fmt.Errorf("failed to handle jobs data. %w", err)
		}

	case *pbsync.SendDataRequest_Licenses:
		if resp.RowsAffected, err = s.handleLicensesData(
			ctx,
			d.Licenses.GetLicenses(),
		); err != nil {
			return nil, fmt.Errorf("failed to handle licenses data. %w", err)
		}

	case *pbsync.SendDataRequest_Users:
		if resp.RowsAffected, err = s.handleUsersData(ctx, d.Users.GetUsers()); err != nil {
			return nil, fmt.Errorf("failed to handle users data. %w", err)
		}

	case *pbsync.SendDataRequest_Vehicles:
		if resp.RowsAffected, err = s.handleVehiclesData(
			ctx,
			d.Vehicles.GetVehicles(),
		); err != nil {
			return nil, fmt.Errorf("failed to handle vehicles data. %w", err)
		}

	case *pbsync.SendDataRequest_Accounts:
		if resp.RowsAffected, err = s.handleAccountsData(
			ctx,
			d.Accounts.GetAccountUpdates(),
			d.Accounts.GetClear(),
		); err != nil {
			return nil, fmt.Errorf("failed to handle accounts data. %w", err)
		}

	case *pbsync.SendDataRequest_UserLocations:
		if resp.RowsAffected, err = s.handleUserLocations(
			ctx,
			d.UserLocations.GetUsers(),
			d.UserLocations.GetClearAll(),
		); err != nil {
			return nil, fmt.Errorf("failed to handle user locations data. %w", err)
		}

	case *pbsync.SendDataRequest_LastCharId:
		if resp.RowsAffected, err = s.handleLastCharId(ctx, d.LastCharId); err != nil {
			return nil, fmt.Errorf("failed to handle last char ID data. %w", err)
		}
	}

	return resp, nil
}

func (s *Store) SendLicenses(
	ctx context.Context,
	req *pbsync.SendLicensesRequest,
) (*pbsync.SendDataResponse, error) {
	rowsAffected, err := s.handleLicensesData(ctx, req.GetLicenses())
	if err != nil {
		return nil, fmt.Errorf("failed to handle licenses data. %w", err)
	}

	return &pbsync.SendDataResponse{RowsAffected: rowsAffected}, nil
}

func (s *Store) SendAccounts(
	ctx context.Context,
	req *pbsync.SendAccountsRequest,
) (*pbsync.SendDataResponse, error) {
	rowsAffected, err := s.handleAccountsData(ctx, req.GetAccountUpdates(), req.GetClear())
	if err != nil {
		return nil, fmt.Errorf("failed to handle accounts data. %w", err)
	}

	return &pbsync.SendDataResponse{RowsAffected: rowsAffected}, nil
}

func (s *Store) SetLastCharID(
	ctx context.Context,
	req *pbsync.SetLastCharIDRequest,
) (*pbsync.SendDataResponse, error) {
	rowsAffected, err := s.handleLastCharId(ctx, req.GetLastCharId())
	if err != nil {
		return nil, fmt.Errorf("failed to handle last character data. %w", err)
	}

	return &pbsync.SendDataResponse{RowsAffected: rowsAffected}, nil
}

func (s *Store) DeleteData(
	ctx context.Context,
	req *pbsync.DeleteDataRequest,
) (*pbsync.DeleteDataResponse, error) {
	switch d := req.GetData().(type) {
	case *pbsync.DeleteDataRequest_Users:
		return s.DeleteUsers(ctx, d.Users.GetUserIds())
	case *pbsync.DeleteDataRequest_Vehicles:
		return s.DeleteVehicles(ctx, d.Vehicles.GetPlates())
	}

	return &pbsync.DeleteDataResponse{}, nil
}

func (s *Store) handleLicensesData(
	ctx context.Context,
	data []*citizenslicenses.License,
) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	tLicenses := table.FivenetLicenses
	stmt := tLicenses.
		INSERT(
			tLicenses.Type,
			tLicenses.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tLicenses.Label.SET(mysql.RawString("VALUES(`label`)")),
		)

	for _, license := range data {
		stmt = stmt.VALUES(license.GetType(), license.GetLabel())
	}

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to execute licenses insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for licenses insert. %w", err)
	}

	return rowsAffected, nil
}

func (s *Store) handleAccountsData(
	ctx context.Context,
	data []*syncactivity.AccountUpdate,
	clearGroups bool,
) (int64, error) {
	if len(data) == 0 && !clearGroups {
		return 0, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	tAccounts := table.FivenetAccounts
	rowsAffected := int64(0)
	accountLicenses := make([]mysql.Expression, 0, len(data))
	pendingGroupChanges := make([]accountGroupChange, 0, len(data))
	for _, account := range data {
		license := account.GetLicense()
		if license == "" {
			continue
		}
		accountLicenses = append(accountLicenses, mysql.String(license))

		groups := accountGroupsFromSyncUpdate(account)

		current, err := s.loadAccountGroupState(ctx, tx, license)
		if err != nil {
			return 0, err
		}
		currentGroups := accountGroupsFromState(current)
		if current != nil && !accountGroupsEqual(currentGroups, groups) {
			pendingGroupChanges = append(pendingGroupChanges, accountGroupChange{
				accountID: current.ID,
				license:   current.License,
				groups:    groups,
			})
		}

		stmt := tAccounts.
			UPDATE(
				tAccounts.License,
				tAccounts.Groups,
			).
			SET(
				account.GetLicense(),
				groups,
			).
			WHERE(tAccounts.License.EQ(mysql.String(license))).
			LIMIT(1)

		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to execute accounts insert statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for accounts insert. %w", err)
		}
		rowsAffected += rows
	}

	// Skip clearing when the request has entries, but none of them carry a license.
	skipClearGroups := clearGroups && len(data) > 0 && len(accountLicenses) == 0
	if clearGroups && !skipClearGroups {
		clearCandidates := []*accountGroupState{}
		clearQueryStmt := tAccounts.
			SELECT(
				tAccounts.ID.AS("account_group_state.id"),
				tAccounts.License.AS("account_group_state.license"),
				tAccounts.Groups.AS("account_group_state.groups"),
			).
			FROM(tAccounts)
		clearConditions := []mysql.BoolExpression{tAccounts.Groups.IS_NOT_NULL()}
		if len(accountLicenses) > 0 {
			clearConditions = append(clearConditions, tAccounts.License.NOT_IN(accountLicenses...))
		}
		clearQueryStmt = clearQueryStmt.WHERE(mysql.AND(clearConditions...))
		if err := clearQueryStmt.QueryContext(ctx, tx, &clearCandidates); err != nil {
			return 0, fmt.Errorf("failed to query account groups to clear. %w", err)
		}
		for _, acc := range clearCandidates {
			if len(acc.Groups.GetGroups()) == 0 {
				continue
			}
			pendingGroupChanges = append(pendingGroupChanges, accountGroupChange{
				accountID: acc.ID,
				license:   acc.License,
				groups:    nil,
			})
		}

		clearUpdateStmt := tAccounts.
			UPDATE(tAccounts.Groups).
			SET(mysql.StringExp(mysql.NULL))

		if len(accountLicenses) > 0 {
			clearUpdateStmt = clearUpdateStmt.
				WHERE(tAccounts.License.NOT_IN(accountLicenses...))
		}

		res, err := clearUpdateStmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to execute accounts clear statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for accounts clear. %w", err)
		}
		rowsAffected += rows
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	for _, change := range pendingGroupChanges {
		s.publishAccountGroupsChanged(ctx, change.accountID, change.license, change.groups)
	}

	return rowsAffected, nil
}

func (s *Store) handleLastCharId(ctx context.Context, data *syncdata.LastCharID) (int64, error) {
	if data.GetLicense() == "" || data.GetLastCharId() == 0 {
		return 0, status.Error(
			codes.InvalidArgument,
			"LastCharId must contain char's identifier and lastCharId",
		)
	}

	tAccounts := table.FivenetAccounts
	stmt := tAccounts.
		UPDATE(tAccounts.LastChar).
		SET(tAccounts.LastChar.SET(mysql.Int32(data.GetLastCharId()))).
		WHERE(tAccounts.License.EQ(mysql.String(data.GetLicense()))).
		LIMIT(1)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to execute last character insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for last character insert. %w", err)
	}

	return rowsAffected, nil
}
