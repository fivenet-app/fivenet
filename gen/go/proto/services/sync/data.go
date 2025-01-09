package sync

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrSendDataDisabled = status.Error(codes.FailedPrecondition, "Sync API: SendData is disabled due to ESXCompat being enabled")

func (s *Server) SendData(ctx context.Context, req *SendDataRequest) (*SendDataResponse, error) {
	resp := &SendDataResponse{
		AffectedRows: 0,
	}

	if s.esxCompat {
		return nil, ErrSendDataDisabled
	}

	var err error
	switch d := req.Data.(type) {
	case *SendDataRequest_Jobs:
		if resp.AffectedRows, err = s.handleJobsData(ctx, d); err != nil {
			return nil, err
		}

	case *SendDataRequest_Licenses:
		if resp.AffectedRows, err = s.handleLicensesData(ctx, d); err != nil {
			return nil, err
		}

	case *SendDataRequest_Users:
		if resp.AffectedRows, err = s.handleUsersData(ctx, d); err != nil {
			return nil, err
		}

	case *SendDataRequest_Vehicles:
		if resp.AffectedRows, err = s.handleVehiclesData(ctx, d); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) handleJobsData(ctx context.Context, data *SendDataRequest_Jobs) (int64, error) {
	tJobs := tables.Jobs()

	stmt := tJobs.
		INSERT(
			tJobs.Name,
			tJobs.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobs.Name.SET(jet.StringExp(jet.Raw("VALUES(`name`)"))),
			tJobs.Label.SET(jet.StringExp(jet.Raw("VALUES(`label`)"))),
		)

	for _, job := range data.Jobs.Jobs {
		stmt = stmt.VALUES(
			job.Name,
			job.Label,
		)
	}

	// How to handle removed/missing jobs?

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	for _, job := range data.Jobs.Jobs {
		rowCounts, err := s.handleJobGrades(ctx, job)
		if err != nil {
			return 0, err
		}

		rowsAffected += rowCounts
	}

	return rowsAffected, nil
}

func (s *Server) handleJobGrades(ctx context.Context, job *users.Job) (int64, error) {
	if len(job.Grades) == 0 {
		return 0, nil
	}

	tJobGrades := tables.JobGrades()

	stmt := tJobGrades.
		INSERT(
			tJobGrades.JobName,
			tJobGrades.Grade,
			tJobGrades.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobGrades.JobName.SET(jet.StringExp(jet.Raw("VALUES(`job_name`)"))),
			tJobGrades.Grade.SET(jet.IntExp(jet.Raw("VALUES(`grade`)"))),
			tJobGrades.Label.SET(jet.StringExp(jet.Raw("VALUES(`label`)"))),
		)

	for _, grade := range job.Grades {
		stmt = stmt.VALUES(
			grade.JobName,
			grade.Grade,
			grade.Label,
		)
	}

	// TODO delete missing job grades

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *Server) handleLicensesData(ctx context.Context, data *SendDataRequest_Licenses) (int64, error) {
	tLicenses := tables.Licenses()

	stmt := tLicenses.
		INSERT(
			tLicenses.Type,
			tLicenses.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tLicenses.Type.SET(jet.StringExp(jet.Raw("VALUES(`type`)"))),
			tLicenses.Label.SET(jet.StringExp(jet.Raw("VALUES(`label`)"))),
		)

	for _, license := range data.Licenses.Licenses {
		stmt = stmt.VALUES(
			license.Type,
			license.Label,
		)
	}

	// TODO delete missing licenses

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *Server) handleUsersData(ctx context.Context, data *SendDataRequest_Users) (int64, error) {
	tUsers := tables.Users()
	_ = tUsers

	toCreate, toUpdate := []*users.User{}, []*users.User{}
	_, _ = toCreate, toUpdate

	// TODO check which user ids already exist in the database and create/update them accordingly

	return 0, nil
}

func (s *Server) handleVehiclesData(ctx context.Context, data *SendDataRequest_Vehicles) (int64, error) {
	tVehicles := tables.OwnedVehicles()

	stmt := tVehicles.
		INSERT(
			tVehicles.Owner,
			tVehicles.Plate,
			tVehicles.Model,
			tVehicles.Type,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tVehicles.Owner.SET(jet.StringExp(jet.Raw("VALUES(`owner`)"))),
			tVehicles.Plate.SET(jet.StringExp(jet.Raw("VALUES(`plate`)"))),
			tVehicles.Model.SET(jet.StringExp(jet.Raw("VALUES(`model`)"))),
			tVehicles.Type.SET(jet.StringExp(jet.Raw("VALUES(`type`)"))),
		)

	for _, vehicle := range data.Vehicles.Vehicles {
		stmt = stmt.VALUES(
			vehicle.Owner,
			vehicle.Plate,
			vehicle.Model,
			vehicle.Type,
		)
	}

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
