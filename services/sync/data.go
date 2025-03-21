package sync

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultUserGroupFallback = "user"

var ErrSendDataDisabled = status.Error(codes.FailedPrecondition, "Sync API: SendData is disabled due to ESXCompat being enabled")

func (s *Server) SendData(ctx context.Context, req *pbsync.SendDataRequest) (*pbsync.SendDataResponse, error) {
	resp := &pbsync.SendDataResponse{
		AffectedRows: 0,
	}

	if s.esxCompat {
		return nil, ErrSendDataDisabled
	}

	var err error
	switch d := req.Data.(type) {
	case *pbsync.SendDataRequest_Jobs:
		if resp.AffectedRows, err = s.handleJobsData(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.SendDataRequest_Licenses:
		if resp.AffectedRows, err = s.handleLicensesData(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.SendDataRequest_Users:
		if resp.AffectedRows, err = s.handleUsersData(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.SendDataRequest_Vehicles:
		if resp.AffectedRows, err = s.handleVehiclesData(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.SendDataRequest_UserLocations:
		if resp.AffectedRows, err = s.handleUserLocations(ctx, d); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) handleJobsData(ctx context.Context, data *pbsync.SendDataRequest_Jobs) (int64, error) {
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

	// ??? Shoud we delete jobs, that are not part of the list, from the database?

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

	rowsAffectedCount := int64(0)

	tJobGrades := tables.JobGrades().AS("jobgrade")

	selectStmt := tJobGrades.
		SELECT(
			tJobGrades.JobName.AS("job_grade.job_name"),
			tJobGrades.Grade,
			tJobGrades.Label,
		).
		FROM(tJobGrades).
		ORDER_BY(
			tJobGrades.Grade.ASC(),
		).
		WHERE(tJobGrades.JobName.EQ(jet.String(job.Name)))

	currentGrades := []*users.JobGrade{}
	if err := selectStmt.QueryContext(ctx, s.db, &currentGrades); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	toCreate, toUpdate, toDelete := []*users.JobGrade{}, []*users.JobGrade{}, []*users.JobGrade{}
	if len(currentGrades) == 0 {
		toCreate = job.Grades
	} else {
		// Update cache
		foundTracker := []int{}
		for _, cg := range currentGrades {
			var found *users.JobGrade
			var foundIdx int

			for i, ug := range job.Grades {
				if cg.Grade != ug.Grade {
					continue
				}

				found = ug
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete = append(toDelete, cg)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cg.Label != found.Label {
				cg.Label = found.Label
				changed = true
			}

			if changed {
				toUpdate = append(toUpdate, cg)
			}
		}

		for i, uj := range job.Grades {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	tJobGrades = tables.JobGrades()

	if len(toCreate) > 0 {
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

		for _, grade := range toCreate {
			stmt = stmt.VALUES(
				grade.JobName,
				grade.Grade,
				grade.Label,
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

		rowsAffectedCount += rowsAffected
	}

	if len(toUpdate) > 0 {
		for _, grade := range toUpdate {
			stmt := tJobGrades.
				UPDATE(
					tJobGrades.JobName,
					tJobGrades.Grade,
					tJobGrades.Label,
				).
				SET(
					grade.JobName,
					grade.Grade,
					grade.Label,
				).
				WHERE(jet.AND(
					tJobGrades.JobName.EQ(jet.String(job.Name)),
					tJobGrades.Grade.EQ(jet.Int32(grade.Grade)),
				))

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, err
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}

			rowsAffectedCount += rowsAffected
		}
	}

	if len(toDelete) > 0 {
		for _, grade := range toDelete {
			stmt := tJobGrades.
				DELETE().
				WHERE(jet.AND(
					tJobGrades.JobName.EQ(jet.String(job.Name)),
					tJobGrades.Grade.EQ(jet.Int32(grade.Grade)),
				)).
				LIMIT(1)

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, err
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}

			rowsAffectedCount += rowsAffected
		}
	}

	return rowsAffectedCount, nil
}

func (s *Server) handleLicensesData(ctx context.Context, data *pbsync.SendDataRequest_Licenses) (int64, error) {
	tLicenses := tables.Licenses()

	stmt := tLicenses.
		INSERT(
			tLicenses.Type,
			tLicenses.Label,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tLicenses.Label.SET(jet.StringExp(jet.Raw("VALUES(`label`)"))),
		)

	for _, license := range data.Licenses.Licenses {
		stmt = stmt.VALUES(
			license.Type,
			license.Label,
		)
	}

	// ??? Shoud we delete jobs, that are not part of the list, from the database?

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

func (s *Server) handleUsersData(ctx context.Context, data *pbsync.SendDataRequest_Users) (int64, error) {
	tUsers := tables.Users()

	defaultUserGroup := defaultUserGroupFallback

	userIds := []jet.Expression{}
	for _, user := range data.Users.Users {
		userIds = append(userIds, jet.Int32(user.UserId))

		if user.Group == nil {
			user.Group = &defaultUserGroup
		}
	}

	checkStmt := tUsers.
		SELECT(
			tUsers.ID,
		).
		FROM(tUsers).
		WHERE(
			tUsers.ID.IN(userIds...),
		)

	var existing []int32
	if err := checkStmt.QueryContext(ctx, s.db, &existing); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	toCreate, toUpdate := []*users.User{}, []*users.User{}
	// Check which user ids already exist in the database and create/update them accordingly
	if len(existing) == 0 {
		toCreate = data.Users.Users
	} else {
		for _, user := range data.Users.Users {
			if idx := slices.IndexFunc(existing, func(userId int32) bool {
				return userId == user.UserId
			}); idx == -1 {
				toCreate = append(toCreate, user)
			} else {
				toUpdate = append(toUpdate, user)
			}
		}
	}

	rowsAffected := int64(0)
	if len(toCreate) > 0 {
		stmt := tUsers.
			INSERT(
				tUsers.ID,
				tUsers.Identifier,
				tUsers.Group,
				tUsers.Firstname,
				tUsers.Lastname,
				tUsers.Dateofbirth,
				tUsers.Job,
				tUsers.JobGrade,
				tUsers.Sex,
				tUsers.PhoneNumber,
				tUsers.Height,
				tUsers.Visum,
				tUsers.Playtime,
			)

		for _, user := range toCreate {
			insertStmt := stmt.
				VALUES(
					user.UserId,
					user.Identifier,
					user.Group,
					user.Firstname,
					user.Lastname,
					user.Dateofbirth,
					user.Job,
					user.JobGrade,
					user.Sex,
					user.PhoneNumber,
					user.Height,
					user.Visum,
					user.Playtime,
				)

			res, err := insertStmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}

			rowsAffected += rows

			if err := s.handleUserLicenses(ctx, *user.Identifier, user.Licenses); err != nil {
				return 0, err
			}
		}
	}

	if len(toUpdate) > 0 {
		for _, user := range toUpdate {
			stmt := tUsers.
				UPDATE(
					tUsers.Identifier,
					tUsers.Group,
					tUsers.Firstname,
					tUsers.Lastname,
					tUsers.Dateofbirth,
					tUsers.Job,
					tUsers.JobGrade,
					tUsers.Sex,
					tUsers.PhoneNumber,
					tUsers.Height,
					tUsers.Visum,
					tUsers.Playtime,
				).
				SET(
					user.Identifier,
					user.Group,
					user.Firstname,
					user.Lastname,
					user.Dateofbirth,
					user.Job,
					user.JobGrade,
					user.Sex,
					user.PhoneNumber,
					user.Height,
					user.Visum,
					user.Playtime,
				).
				WHERE(
					tUsers.ID.EQ(jet.Int32(user.UserId)),
				)

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}

			rowsAffected += rows
		}
	}

	return rowsAffected, nil
}

func (s *Server) handleUserLicenses(ctx context.Context, identifier string, licenses []*users.License) error {
	tUserLicenses := tables.UserLicenses()

	if len(licenses) == 0 {
		// User has no licenses? Delete all from the database.
		stmt := tUserLicenses.
			DELETE().
			WHERE(tUserLicenses.Owner.EQ(jet.String(identifier))).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}

		return nil
	}

	selectStmt := tUserLicenses.
		SELECT(
			tUserLicenses.Type,
		).
		FROM(tUserLicenses).
		WHERE(tUserLicenses.Owner.EQ(jet.String(identifier)))

	currentLicenses := []string{}
	if err := selectStmt.QueryContext(ctx, s.db, &currentLicenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	licensesList := []string{}
	for _, license := range licenses {
		licensesList = append(licensesList, license.Type)
	}

	toAdd, toRemove := utils.SlicesDifference(currentLicenses, licensesList)

	if len(toAdd) > 0 {
		stmt := tUserLicenses.
			INSERT(
				tUserLicenses.Owner,
				tUserLicenses.Type,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tUserLicenses.Type.SET(jet.StringExp(jet.Raw("VALUES(`type`)"))),
			)

		for _, t := range toAdd {
			stmt = stmt.
				VALUES(
					identifier,
					t,
				)
		}

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	if len(toRemove) > 0 {
		types := []jet.Expression{}
		for _, t := range toRemove {
			types = append(types, jet.String(t))
		}

		stmt := tUserLicenses.
			DELETE().
			WHERE(jet.AND(
				tUserLicenses.Owner.EQ(jet.String(identifier)),
				tUserLicenses.Type.IN(types...),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleVehiclesData(ctx context.Context, data *pbsync.SendDataRequest_Vehicles) (int64, error) {
	tVehicles := tables.OwnedVehicles()

	stmt := tVehicles.
		INSERT(
			tVehicles.Owner,
			tVehicles.Plate,
			tVehicles.Model,
			tVehicles.Type,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tVehicles.Plate.SET(jet.StringExp(jet.Raw("VALUES(`plate`)"))),
			tVehicles.Model.SET(jet.StringExp(jet.Raw("VALUES(`model`)"))),
			tVehicles.Type.SET(jet.StringExp(jet.Raw("VALUES(`type`)"))),
		)

	for _, vehicle := range data.Vehicles.Vehicles {
		var ownerId jet.Expression
		if vehicle.Owner != nil && vehicle.Owner.Identifier != nil && *vehicle.Owner.Identifier != "" {
			ownerId = jet.String(*vehicle.Owner.Identifier)
		} else if vehicle.OwnerId != nil {
			ownerId = jet.Int32(*vehicle.OwnerId)
		}

		stmt = stmt.VALUES(
			ownerId,
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

func (s *Server) handleUserLocations(ctx context.Context, data *pbsync.SendDataRequest_UserLocations) (int64, error) {
	tLocations := table.FivenetUserLocations

	// Handle clear all
	if data.UserLocations.ClearAll != nil && *data.UserLocations.ClearAll {
		stmt := tLocations.
			DELETE()

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return 0, err
		}
	}

	stmt := tLocations.
		INSERT(
			tLocations.Identifier,
			tLocations.Job,
			tLocations.X,
			tLocations.Y,
			tLocations.Hidden,
		)

	atLeastOne := false
	toDelete := []string{}
	for _, location := range data.UserLocations.Users {
		// Collect user locations are marked for removal
		if location.Remove {
			toDelete = append(toDelete, location.Identifier)
			continue
		}

		stmt = stmt.
			VALUES(
				location.Identifier,
				location.Job,
				location.Coords.X,
				location.Coords.Y,
				location.Hidden,
			)
		atLeastOne = true
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(
			tLocations.Job.SET(jet.StringExp(jet.Raw("VALUES(`job`)"))),
			tLocations.X.SET(jet.FloatExp(jet.Raw("VALUES(`x`)"))),
			tLocations.Y.SET(jet.FloatExp(jet.Raw("VALUES(`y`)"))),
			tLocations.Hidden.SET(jet.BoolExp(jet.Raw("VALUES(`hidden`)"))),
		)

	rowsAffected := int64(0)
	if atLeastOne {
		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return 0, err
		}

		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return 0, err
		}
	}

	// Delete any user locations that have been marked for removal
	if len(toDelete) > 0 {
		identifiers := []jet.Expression{}
		for _, identifier := range toDelete {
			identifiers = append(identifiers, jet.String(identifier))
		}

		delStmt := tLocations.
			DELETE().
			WHERE(tLocations.Identifier.IN(identifiers...)).
			LIMIT(int64(len(toDelete)))

		res, err := delStmt.ExecContext(ctx, s.db)
		if err != nil {
			return 0, err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		rowsAffected += rows
	}

	return rowsAffected, nil
}

func (s *Server) DeleteData(ctx context.Context, req *pbsync.DeleteDataRequest) (*pbsync.DeleteDataResponse, error) {
	if s.esxCompat {
		return nil, ErrSendDataDisabled
	}

	rowsAffected := int64(0)

	switch d := req.Data.(type) {
	case *pbsync.DeleteDataRequest_Users:
		userIds := []jet.Expression{}
		for _, identifier := range d.Users.UserIds {
			userIds = append(userIds, jet.Int32(identifier))
		}

		tUsers := tables.Users()

		delStmt := tUsers.
			DELETE().
			WHERE(tUsers.ID.IN(userIds...)).
			LIMIT(int64(len(d.Users.UserIds)))

		res, err := delStmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		rowsAffected += rows

	case *pbsync.DeleteDataRequest_Vehicles:
		plates := []jet.Expression{}
		for _, plate := range d.Vehicles.Plates {
			plates = append(plates, jet.String(plate))
		}

		tVehicles := tables.OwnedVehicles()

		delStmt := tVehicles.
			DELETE().
			WHERE(tVehicles.Plate.IN(plates...)).
			LIMIT(int64(len(d.Vehicles.Plates)))

		res, err := delStmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		rowsAffected += rows
	}

	return &pbsync.DeleteDataResponse{
		AffectedRows: rowsAffected,
	}, nil
}
