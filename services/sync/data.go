package sync

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SendData(
	ctx context.Context,
	req *pbsync.SendDataRequest,
) (*pbsync.SendDataResponse, error) {
	resp := &pbsync.SendDataResponse{
		AffectedRows: 0,
	}

	s.lastSyncedData.Store(time.Now().Unix())

	var err error
	switch d := req.GetData().(type) {
	case *pbsync.SendDataRequest_Jobs:
		if resp.AffectedRows, err = s.handleJobsData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle jobs data. %w", err)
		}

	case *pbsync.SendDataRequest_Licenses:
		if resp.AffectedRows, err = s.handleLicensesData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle licenses data. %w", err)
		}

	case *pbsync.SendDataRequest_Users:
		if resp.AffectedRows, err = s.handleUsersData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle users data. %w", err)
		}

	case *pbsync.SendDataRequest_Vehicles:
		if resp.AffectedRows, err = s.handleVehiclesData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle vehicles data. %w", err)
		}

	case *pbsync.SendDataRequest_Accounts:
		if resp.AffectedRows, err = s.handleAccountsData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle accounts data. %w", err)
		}

	case *pbsync.SendDataRequest_UserLocations:
		if resp.AffectedRows, err = s.handleUserLocations(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle user locations data. %w", err)
		}

	case *pbsync.SendDataRequest_LastCharId:
		if resp.AffectedRows, err = s.handleLastCharId(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle user locations data. %w", err)
		}
	}

	return resp, nil
}

func (s *Server) handleLicensesData(
	ctx context.Context,
	data *pbsync.SendDataRequest_Licenses,
) (int64, error) {
	if len(data.Licenses.GetLicenses()) == 0 {
		return 0, nil
	}

	tLicenses := table.FivenetLicenses

	stmt := tLicenses.
		INSERT(
			tLicenses.Type,
			tLicenses.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tLicenses.Label.SET(mysql.StringExp(mysql.Raw("VALUES(`label`)"))),
		)

	for _, license := range data.Licenses.GetLicenses() {
		stmt = stmt.VALUES(
			license.GetType(),
			license.GetLabel(),
		)
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

func (s *Server) handleAccountsData(
	ctx context.Context,
	data *pbsync.SendDataRequest_Accounts,
) (int64, error) {
	if len(data.Accounts.GetAccountUpdates()) == 0 {
		return 0, nil
	}

	tAccounts := table.FivenetAccounts

	stmt := tAccounts.
		INSERT(
			tAccounts.License,
			tAccounts.Groups,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tAccounts.Groups.SET(mysql.StringExp(mysql.Raw("VALUES(`groups`)"))),
		)

	for _, account := range data.Accounts.GetAccountUpdates() {
		var groups *accounts.AccountGroups
		gs := account.GetGroups()
		if len(gs) > 0 {
			groups = &accounts.AccountGroups{
				Groups: gs,
			}
		} else if account.GetGroup() != "" {
			groups = &accounts.AccountGroups{
				Groups: []string{account.GetGroup()},
			}
		}

		stmt = stmt.VALUES(
			account.GetLicense(),
			groups,
		)
	}

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to execute accounts insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for accounts insert. %w", err)
	}

	return rowsAffected, nil
}

func (s *Server) handleLastCharId(
	ctx context.Context,
	data *pbsync.SendDataRequest_LastCharId,
) (int64, error) {
	if data.LastCharId == nil || data.LastCharId.GetIdentifier() == "" ||
		data.LastCharId.LastCharId == nil ||
		data.LastCharId.GetLastCharId() == 0 {
		return 0, status.Error(
			codes.InvalidArgument,
			"LastCharId must contain UserId and CharacterId",
		)
	}

	tAccounts := table.FivenetAccounts

	stmt := tAccounts.
		UPDATE(
			tAccounts.LastChar,
		).
		SET(
			tAccounts.LastChar.SET(mysql.Int32(data.LastCharId.GetLastCharId())),
		).
		WHERE(
			tAccounts.License.EQ(mysql.String(data.LastCharId.GetIdentifier())),
		).
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

func (s *Server) DeleteData(
	ctx context.Context,
	req *pbsync.DeleteDataRequest,
) (*pbsync.DeleteDataResponse, error) {
	rowsAffected := int64(0)

	switch d := req.GetData().(type) {
	case *pbsync.DeleteDataRequest_Users:
		userIds := []mysql.Expression{}
		for _, identifier := range d.Users.GetUserIds() {
			userIds = append(userIds, mysql.Int32(identifier))
		}

		tUsers := table.FivenetUser

		delStmt := tUsers.
			DELETE().
			WHERE(tUsers.ID.IN(userIds...)).
			LIMIT(int64(len(d.Users.GetUserIds())))

		res, err := delStmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, fmt.Errorf("failed to execute users delete statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve rows affected for users delete. %w", err)
		}

		rowsAffected += rows

	case *pbsync.DeleteDataRequest_Vehicles:
		plates := []mysql.Expression{}
		for _, plate := range d.Vehicles.GetPlates() {
			plates = append(plates, mysql.String(plate))
		}

		tVehicles := table.FivenetOwnedVehicles

		delStmt := tVehicles.
			DELETE().
			WHERE(tVehicles.Plate.IN(plates...)).
			LIMIT(int64(len(d.Vehicles.GetPlates())))

		res, err := delStmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, fmt.Errorf("failed to execute vehicles delete statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve rows affected for vehicles delete. %w", err)
		}

		rowsAffected += rows
	}

	return &pbsync.DeleteDataResponse{
		AffectedRows: rowsAffected,
	}, nil
}

func getLicenseFromIdentifier(identifier string) string {
	parts := strings.SplitN(identifier, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return identifier
}
