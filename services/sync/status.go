package sync

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) GetStatus(
	ctx context.Context,
	req *pbsync.GetStatusRequest,
) (*pbsync.GetStatusResponse, error) {
	resp := &pbsync.GetStatusResponse{}

	tJobs := table.FivenetJobs
	tUsers := table.FivenetUser
	tVehicles := table.FivenetOwnedVehicles
	tLicenses := table.FivenetLicenses

	// Jobs
	jobsStmt := tJobs.
		SELECT(
			mysql.COUNT(tJobs.Name),
		).
		FROM(tJobs)

	var jobsCount database.DataCount
	if err := jobsStmt.QueryContext(ctx, s.db, &jobsCount); err != nil {
		return nil, err
	}
	resp.Jobs = &syncdata.DataStatus{
		Count: jobsCount.Total,
	}

	// Accounts
	// Users
	accountsStmt := tUsers.
		SELECT(
			mysql.COUNT(tUsers.ID),
		).
		FROM(tUsers)

	var accountsCount database.DataCount
	if err := accountsStmt.QueryContext(ctx, s.db, &accountsCount); err != nil {
		return nil, err
	}
	resp.Accounts = &syncdata.DataStatus{
		Count: accountsCount.Total,
	}

	// Users
	usersStmt := tUsers.
		SELECT(
			mysql.COUNT(tUsers.ID),
		).
		FROM(tUsers)

	var usersCount database.DataCount
	if err := usersStmt.QueryContext(ctx, s.db, &usersCount); err != nil {
		return nil, err
	}
	resp.Users = &syncdata.DataStatus{
		Count: usersCount.Total,
	}

	// Vehicles
	vehiclesStmt := tVehicles.
		SELECT(
			mysql.COUNT(tVehicles.Plate),
		).
		FROM(tVehicles)

	var vehiclesCount database.DataCount
	if err := vehiclesStmt.QueryContext(ctx, s.db, &vehiclesCount); err != nil {
		return nil, err
	}
	resp.Vehicles = &syncdata.DataStatus{
		Count: vehiclesCount.Total,
	}

	// Licenses
	licensesStmt := tLicenses.
		SELECT(
			mysql.COUNT(tLicenses.Type),
		).
		FROM(tLicenses)

	var licensesCount database.DataCount
	if err := licensesStmt.QueryContext(ctx, s.db, &licensesCount); err != nil {
		return nil, err
	}
	resp.Licenses = &syncdata.DataStatus{
		Count: licensesCount.Total,
	}

	return resp, nil
}
