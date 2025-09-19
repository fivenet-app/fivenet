package sync

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/sync"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) GetStatus(
	ctx context.Context,
	req *pbsync.GetStatusRequest,
) (*pbsync.GetStatusResponse, error) {
	resp := &pbsync.GetStatusResponse{}

	tJobs := tables.Jobs()
	tUsers := tables.User()
	tVehicles := tables.OwnedVehicles()
	tLicenses := tables.Licenses()

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
	resp.Jobs = &sync.DataStatus{
		Count: jobsCount.Total,
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
	resp.Users = &sync.DataStatus{
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
	resp.Vehicles = &sync.DataStatus{
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
	resp.Licenses = &sync.DataStatus{
		Count: licensesCount.Total,
	}

	return resp, nil
}
