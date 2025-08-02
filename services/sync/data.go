package sync

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultUserGroupFallback = "user"

var ErrSendDataDisabled = status.Error(codes.FailedPrecondition, "Sync API: SendData is disabled due to ESXCompat being enabled")

func (s *Server) SendData(ctx context.Context, req *pbsync.SendDataRequest) (*pbsync.SendDataResponse, error) {
	resp := &pbsync.SendDataResponse{
		AffectedRows: 0,
	}

	s.lastSyncedData.Store(time.Now().Unix())

	var err error
	switch d := req.Data.(type) {
	case *pbsync.SendDataRequest_Jobs:
		if s.esxCompat {
			return nil, ErrSendDataDisabled
		}

		if resp.AffectedRows, err = s.handleJobsData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle jobs data. %w", err)
		}

	case *pbsync.SendDataRequest_Licenses:
		if s.esxCompat {
			return nil, ErrSendDataDisabled
		}

		if resp.AffectedRows, err = s.handleLicensesData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle licenses data. %w", err)
		}

	case *pbsync.SendDataRequest_Users:
		if s.esxCompat {
			return nil, ErrSendDataDisabled
		}

		if resp.AffectedRows, err = s.handleUsersData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle users data. %w", err)
		}

	case *pbsync.SendDataRequest_Vehicles:
		if s.esxCompat {
			return nil, ErrSendDataDisabled
		}

		if resp.AffectedRows, err = s.handleVehiclesData(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle vehicles data. %w", err)
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

func (s *Server) handleJobsData(ctx context.Context, data *pbsync.SendDataRequest_Jobs) (int64, error) {
	if len(data.Jobs.Jobs) == 0 {
		return 0, nil
	}

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
		return 0, fmt.Errorf("failed to execute job insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for job insert. %w", err)
	}

	for _, job := range data.Jobs.Jobs {
		rowCounts, err := s.handleJobGrades(ctx, job)
		if err != nil {
			return 0, fmt.Errorf("failed to handle job grades for job %s. %w", job.Name, err)
		}

		rowsAffected += rowCounts
	}

	return rowsAffected, nil
}

func (s *Server) handleJobGrades(ctx context.Context, job *jobs.Job) (int64, error) {
	if len(job.Grades) == 0 {
		return 0, nil
	}

	rowsAffectedCount := int64(0)

	tJobsGrades := tables.JobsGrades().AS("job_grade")

	selectStmt := tJobsGrades.
		SELECT(
			tJobsGrades.JobName.AS("job_grade.job_name"),
			tJobsGrades.Grade,
			tJobsGrades.Label,
		).
		FROM(tJobsGrades).
		ORDER_BY(
			tJobsGrades.Grade.ASC(),
		).
		WHERE(tJobsGrades.JobName.EQ(jet.String(job.Name)))

	currentGrades := []*jobs.JobGrade{}
	if err := selectStmt.QueryContext(ctx, s.db, &currentGrades); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, fmt.Errorf("failed to query current job grades for job %s. %w", job.Name, err)
		}
	}

	toCreate, toUpdate, toDelete := []*jobs.JobGrade{}, []*jobs.JobGrade{}, []*jobs.JobGrade{}
	if len(currentGrades) == 0 {
		toCreate = job.Grades
	} else {
		// Update cache
		foundTracker := []int{}
		for _, cg := range currentGrades {
			var found *jobs.JobGrade
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

	tJobsGrades = tables.JobsGrades()

	if len(toCreate) > 0 {
		stmt := tJobsGrades.
			INSERT(
				tJobsGrades.JobName,
				tJobsGrades.Grade,
				tJobsGrades.Label,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tJobsGrades.JobName.SET(jet.StringExp(jet.Raw("VALUES(`job_name`)"))),
				tJobsGrades.Grade.SET(jet.IntExp(jet.Raw("VALUES(`grade`)"))),
				tJobsGrades.Label.SET(jet.StringExp(jet.Raw("VALUES(`label`)"))),
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
			return 0, fmt.Errorf("failed to execute job grades insert statement. %w", err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for job grades insert. %w", err)
		}

		rowsAffectedCount += rowsAffected
	}

	if len(toUpdate) > 0 {
		for _, grade := range toUpdate {
			stmt := tJobsGrades.
				UPDATE(
					tJobsGrades.JobName,
					tJobsGrades.Grade,
					tJobsGrades.Label,
				).
				SET(
					grade.JobName,
					grade.Grade,
					grade.Label,
				).
				WHERE(jet.AND(
					tJobsGrades.JobName.EQ(jet.String(job.Name)),
					tJobsGrades.Grade.EQ(jet.Int32(grade.Grade)),
				))

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, fmt.Errorf("failed to execute job grades update statement for grade %d. %w", grade.Grade, err)
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve rows affected for job grades update. %w", err)
			}

			rowsAffectedCount += rowsAffected
		}
	}

	if len(toDelete) > 0 {
		for _, grade := range toDelete {
			stmt := tJobsGrades.
				DELETE().
				WHERE(jet.AND(
					tJobsGrades.JobName.EQ(jet.String(job.Name)),
					tJobsGrades.Grade.EQ(jet.Int32(grade.Grade)),
				)).
				LIMIT(1)

			res, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, fmt.Errorf("failed to execute job grades delete statement for grade %d. %w", grade.Grade, err)
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve rows affected for job grades delete. %w", err)
			}

			rowsAffectedCount += rowsAffected
		}
	}

	return rowsAffectedCount, nil
}

func (s *Server) handleLicensesData(ctx context.Context, data *pbsync.SendDataRequest_Licenses) (int64, error) {
	if len(data.Licenses.Licenses) == 0 {
		return 0, nil
	}

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

func (s *Server) handleUsersData(ctx context.Context, data *pbsync.SendDataRequest_Users) (int64, error) {
	tUsers := tables.User()

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
			return 0, fmt.Errorf("failed to query existing users. %w", err)
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
				).
				ON_DUPLICATE_KEY_UPDATE(
					tUsers.Group.SET(jet.StringExp(jet.Raw("VALUES(`group`)"))),
					tUsers.Firstname.SET(jet.StringExp(jet.Raw("VALUES(`firstname`)"))),
					tUsers.Lastname.SET(jet.StringExp(jet.Raw("VALUES(`lastname`)"))),
					tUsers.Dateofbirth.SET(jet.StringExp(jet.Raw("VALUES(`dateofbirth`)"))),
					tUsers.Job.SET(jet.StringExp(jet.Raw("VALUES(`job`)"))),
					tUsers.JobGrade.SET(jet.IntExp(jet.Raw("VALUES(`job_grade`)"))),
					tUsers.Sex.SET(jet.StringExp(jet.Raw("VALUES(`sex`)"))),
					tUsers.PhoneNumber.SET(jet.StringExp(jet.Raw("VALUES(`phone_number`)"))),
					tUsers.Height.SET(jet.StringExp(jet.Raw("VALUES(`height`)"))),
					tUsers.Visum.SET(jet.IntExp(jet.Raw("VALUES(`visum`)"))),
					tUsers.Playtime.SET(jet.IntExp(jet.Raw("VALUES(`playtime`)"))),
				)

			res, err := insertStmt.ExecContext(ctx, s.db)
			if err != nil {
				return 0, fmt.Errorf("failed to execute user insert statement. %w", err)
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve rows affected for user insert. %w", err)
			}

			rowsAffected += rows

			if err := s.handleCitizensLicenses(ctx, *user.Identifier, user.Licenses); err != nil {
				return 0, fmt.Errorf("failed to handle user licenses for user %s. %w", *user.Identifier, err)
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
				return 0, fmt.Errorf("failed to execute user update statement. %w", err)
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve rows affected for user update. %w", err)
			}

			rowsAffected += rows
		}
	}

	return rowsAffected, nil
}

func (s *Server) handleCitizensLicenses(ctx context.Context, identifier string, licenses []*users.License) error {
	tCitizensLicenses := tables.UserLicenses()

	if len(licenses) == 0 {
		// User has no licenses? Delete all user licenses from the database.
		stmt := tCitizensLicenses.
			DELETE().
			WHERE(tCitizensLicenses.Owner.EQ(jet.String(identifier))).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to execute user licenses delete statement. %w", err)
		}

		return nil
	}

	selectStmt := tCitizensLicenses.
		SELECT(
			tCitizensLicenses.Type,
		).
		FROM(tCitizensLicenses).
		WHERE(tCitizensLicenses.Owner.EQ(jet.String(identifier)))

	currentLicenses := []string{}
	if err := selectStmt.QueryContext(ctx, s.db, &currentLicenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query current user licenses for identifier %s. %w", identifier, err)
		}
	}

	licensesList := []string{}
	for _, license := range licenses {
		licensesList = append(licensesList, license.Type)
	}

	toAdd, toRemove := utils.SlicesDifference(currentLicenses, licensesList)

	if len(toAdd) > 0 {
		stmt := tCitizensLicenses.
			INSERT(
				tCitizensLicenses.Owner,
				tCitizensLicenses.Type,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCitizensLicenses.Type.SET(jet.StringExp(jet.Raw("VALUES(`type`)"))),
			)

		for _, t := range toAdd {
			stmt = stmt.
				VALUES(
					identifier,
					t,
				)
		}

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to execute user licenses insert statement. %w", err)
		}
	}

	if len(toRemove) > 0 {
		types := []jet.Expression{}
		for _, t := range toRemove {
			types = append(types, jet.String(t))
		}

		stmt := tCitizensLicenses.
			DELETE().
			WHERE(jet.AND(
				tCitizensLicenses.Owner.EQ(jet.String(identifier)),
				tCitizensLicenses.Type.IN(types...),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to execute user licenses delete statement. %w", err)
		}
	}

	return nil
}

func (s *Server) handleVehiclesData(ctx context.Context, data *pbsync.SendDataRequest_Vehicles) (int64, error) {
	if len(data.Vehicles.Vehicles) == 0 {
		return 0, nil
	}

	tVehicles := tables.OwnedVehicles()

	var stmt jet.InsertStatement
	if !tables.ESXCompatEnabled {
		stmt = tVehicles.
			INSERT(
				tVehicles.Owner,
				tVehicles.Plate,
				tVehicles.Model,
				tVehicles.Type,
				tVehicles.Job,
				tVehicles.Data,
			)
	} else {
		stmt = tVehicles.
			INSERT(
				tVehicles.Owner,
				tVehicles.Plate,
				tVehicles.Model,
				tVehicles.Type,
			)
	}

	for _, vehicle := range data.Vehicles.Vehicles {
		var ownerId jet.Expression
		if vehicle.OwnerIdentifier != nil && *vehicle.OwnerIdentifier != "" {
			ownerId = jet.String(*vehicle.OwnerIdentifier)
		} else if vehicle.OwnerId != nil {
			ownerId = jet.Int32(*vehicle.OwnerId)
		}

		if !tables.ESXCompatEnabled {
			stmt = stmt.VALUES(
				ownerId,
				vehicle.Plate,
				vehicle.Model,
				vehicle.Type,
				vehicle.Job,
				jet.NULL,
			)
		} else {
			stmt = stmt.VALUES(
				ownerId,
				vehicle.Plate,
				vehicle.Model,
				vehicle.Type,
			)
		}
	}

	assignments := []jet.ColumnAssigment{
		tVehicles.Owner.SET(jet.StringExp(jet.Raw("VALUES(`owner`)"))),
		tVehicles.Plate.SET(jet.StringExp(jet.Raw("VALUES(`plate`)"))),
		tVehicles.Model.SET(jet.StringExp(jet.Raw("VALUES(`model`)"))),
		tVehicles.Type.SET(jet.StringExp(jet.Raw("VALUES(`type`)"))),
	}

	if !tables.ESXCompatEnabled {
		assignments = append(assignments,
			tVehicles.Job.SET(jet.StringExp(jet.Raw("VALUES(`job`)"))),
			tVehicles.Data.SET(jet.StringExp(jet.Raw("VALUES(`data`)"))),
		)
	}

	stmt = stmt.
		ON_DUPLICATE_KEY_UPDATE(assignments...)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to execute vehicles insert statement. %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for vehicles insert. %w", err)
	}

	return rowsAffected, nil
}

func (s *Server) handleUserLocations(ctx context.Context, data *pbsync.SendDataRequest_UserLocations) (int64, error) {
	tLocations := table.FivenetCentrumUserLocations

	// Handle clear all
	if data.UserLocations.ClearAll != nil && *data.UserLocations.ClearAll {
		stmt := tLocations.
			DELETE().
			WHERE(tLocations.Identifier.IS_NOT_NULL().OR(tLocations.Identifier.IS_NULL()))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return 0, fmt.Errorf("failed to execute user locations clear all statement. %w", err)
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
			s.logger.Debug("failed to execute user locations insert statement", zap.Any("data", data), zap.Error(err))
			return 0, fmt.Errorf("failed to execute user locations insert statement. %w", err)
		}

		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for user locations insert. %w", err)
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
			return 0, fmt.Errorf("failed to execute user locations delete statement. %w", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve rows affected for user locations delete. %w", err)
		}
		rowsAffected += rows
	}

	return rowsAffected, nil
}

func (s *Server) handleLastCharId(ctx context.Context, data *pbsync.SendDataRequest_LastCharId) (int64, error) {
	if data.LastCharId == nil || data.LastCharId.Identifier == "" || data.LastCharId.LastCharId == nil || *data.LastCharId.LastCharId == 0 {
		return 0, status.Error(codes.InvalidArgument, "LastCharId must contain UserId and CharacterId")
	}

	tAccounts := table.FivenetAccounts

	stmt := tAccounts.
		UPDATE(
			tAccounts.LastChar,
		).
		SET(
			tAccounts.LastChar.SET(jet.Int32(*data.LastCharId.LastCharId)),
		).
		WHERE(
			tAccounts.License.EQ(jet.String(data.LastCharId.Identifier)),
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

		tUsers := tables.User()

		delStmt := tUsers.
			DELETE().
			WHERE(tUsers.ID.IN(userIds...)).
			LIMIT(int64(len(d.Users.UserIds)))

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
