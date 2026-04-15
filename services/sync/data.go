package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	citizenslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/licenses"
	syncactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
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
		RowsAffected: 0,
	}

	s.lastSyncedData.Store(time.Now().Unix())

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

func (s *Server) SendLicenses(
	ctx context.Context,
	req *pbsync.SendLicensesRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())

	rowsAffected, err := s.handleLicensesData(ctx, req.GetLicenses())
	if err != nil {
		return nil, fmt.Errorf("failed to handle licenses data. %w", err)
	}

	return &pbsync.SendDataResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (s *Server) handleLicensesData(
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
	data []*syncactivity.AccountUpdate,
) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	tAccounts := table.FivenetAccounts

	rowsAffected := int64(0)
	for _, account := range data {
		var groups *accounts.AccountGroups
		if account.GetGroups() != nil && len(account.GetGroups().GetGroups()) > 0 {
			groups = account.GetGroups()
		} else if account.GetGroup() != "" {
			groups = &accounts.AccountGroups{
				Groups: []string{account.GetGroup()},
			}
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
			WHERE(
				tAccounts.License.EQ(mysql.String(account.License)),
			).
			LIMIT(1)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return 0, fmt.Errorf("failed to execute accounts insert statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for accounts insert. %w", err)
		}
		rowsAffected += rows
	}

	return rowsAffected, nil
}

func (s *Server) SendAccounts(
	ctx context.Context,
	req *pbsync.SendAccountsRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())

	rowsAffected, err := s.handleAccountsData(ctx, req.GetAccountUpdates())
	if err != nil {
		return nil, fmt.Errorf("failed to handle accounts data. %w", err)
	}

	return &pbsync.SendDataResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (s *Server) handleLastCharId(
	ctx context.Context,
	data *syncdata.LastCharID,
) (int64, error) {
	if data.LastCharId == nil || data.GetIdentifier() == "" ||
		data.LastCharId == nil ||
		data.GetLastCharId() == 0 {
		return 0, status.Error(
			codes.InvalidArgument,
			"LastCharId must contain CharacterId and an UserId or Identifier",
		)
	}

	tAccounts := table.FivenetAccounts
	stmt := tAccounts.
		UPDATE(
			tAccounts.LastChar,
		).
		SET(
			tAccounts.LastChar.SET(mysql.Int32(data.GetLastCharId())),
		).
		WHERE(
			tAccounts.License.EQ(mysql.String(data.GetIdentifier())),
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

func (s *Server) SetLastCharID(
	ctx context.Context,
	req *pbsync.SetLastCharIDRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())

	rowsAffected, err := s.handleLastCharId(ctx, req.GetLastCharId())
	if err != nil {
		return nil, fmt.Errorf("failed to handle last character data. %w", err)
	}

	return &pbsync.SendDataResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (s *Server) DeleteData(
	ctx context.Context,
	req *pbsync.DeleteDataRequest,
) (*pbsync.DeleteDataResponse, error) {
	switch d := req.GetData().(type) {
	case *pbsync.DeleteDataRequest_Users:
		return s.DeleteUsers(ctx, &pbsync.DeleteUsersRequest{
			UserIds: d.Users.GetUserIds(),
		})

	case *pbsync.DeleteDataRequest_Vehicles:
		return s.DeleteVehicles(ctx, &pbsync.DeleteVehiclesRequest{
			Plates: d.Vehicles.GetPlates(),
		})
	}

	return &pbsync.DeleteDataResponse{}, nil
}
