package sync

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tJobs         = table.Jobs
	tUsers        = table.Users
	tVehicles     = table.OwnedVehicles
	tLicenses     = table.Licenses
	tUserLicenses = table.UserLicenses
)

func (s *Server) GetStatus(ctx context.Context, req *GetStatusRequest) (*GetStatusResponse, error) {
	resp := &GetStatusResponse{}

	// Jobs
	jobsStmt := tJobs.
		SELECT(
			jet.COUNT(tJobs.Name),
		).
		FROM(tJobs)

	var jobsCount database.DataCount
	if err := jobsStmt.QueryContext(ctx, s.db, &jobsCount); err != nil {
		return nil, err
	}
	resp.Jobs = &sync.DataStatus{
		Count: jobsCount.TotalCount,
	}

	// Users
	usersStmt := tUsers.
		SELECT(
			jet.COUNT(tUsers.ID),
		).
		FROM(tUsers)

	var usersCount database.DataCount
	if err := usersStmt.QueryContext(ctx, s.db, &usersCount); err != nil {
		return nil, err
	}
	resp.Users = &sync.DataStatus{
		Count: usersCount.TotalCount,
	}

	// Vehicles
	vehiclesStmt := tVehicles.
		SELECT(
			jet.COUNT(tVehicles.Plate),
		).
		FROM(tVehicles)

	var vehiclesCount database.DataCount
	if err := vehiclesStmt.QueryContext(ctx, s.db, &vehiclesCount); err != nil {
		return nil, err
	}
	resp.Vehicles = &sync.DataStatus{
		Count: vehiclesCount.TotalCount,
	}

	// Licenses
	licensesStmt := tLicenses.
		SELECT(
			jet.COUNT(tLicenses.Type),
		).
		FROM(tLicenses)

	var licensesCount database.DataCount
	if err := licensesStmt.QueryContext(ctx, s.db, &licensesCount); err != nil {
		return nil, err
	}
	resp.Licenses = &sync.DataStatus{
		Count: licensesCount.TotalCount,
	}

	// User Licenses
	userLicensesStmt := tUserLicenses.
		SELECT(
			jet.COUNT(tUserLicenses.Type),
		).
		FROM(tUserLicenses)

	var userLicensesCount database.DataCount
	if err := userLicensesStmt.QueryContext(ctx, s.db, &userLicensesCount); err != nil {
		return nil, err
	}
	resp.UserLicenses = &sync.DataStatus{
		Count: userLicensesCount.TotalCount,
	}

	return resp, nil
}
