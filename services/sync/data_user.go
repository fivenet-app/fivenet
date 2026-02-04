package sync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	userslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/licenses"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) handleUsersData(
	ctx context.Context,
	data *pbsync.SendDataRequest_Users,
) (int64, error) {
	tUsers := table.FivenetUser

	userIds := []mysql.Expression{}
	for _, user := range data.Users.GetUsers() {
		userIds = append(userIds, mysql.Int32(user.GetUserId()))
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

	toCreate, toUpdate := []*syncdata.DataUser{}, []*syncdata.DataUser{}
	// Check which user ids already exist in the database and create/update them accordingly
	if len(existing) == 0 {
		toCreate = data.Users.GetUsers()
	} else {
		for _, user := range data.Users.GetUsers() {
			if idx := slices.IndexFunc(existing, func(userId int32) bool {
				return userId == user.GetUserId()
			}); idx == -1 {
				toCreate = append(toCreate, user)
			} else {
				toUpdate = append(toUpdate, user)
			}
		}
	}

	tAccounts := table.FivenetAccounts

	rowsAffected := int64(0)
	if len(toCreate) > 0 {
		// Begin transaction
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return 0, err
		}
		// Defer a rollback in case anything fails
		defer tx.Rollback()

		stmt := tUsers.
			INSERT(
				tUsers.ID,
				tUsers.AccountID,
				tUsers.Identifier,
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
			var accountIdStmt mysql.SelectStatement = nil
			if user.GetIdentifier() != "" {
				accountIdStmt = tAccounts.
					SELECT(
						mysql.COALESCE(tAccounts.ID, mysql.NULL),
					).
					FROM(tAccounts).
					WHERE(tAccounts.License.EQ(mysql.String(getLicenseFromIdentifier(user.GetIdentifier())))).
					LIMIT(1)
			}

			insertStmt := stmt.
				VALUES(
					user.GetUserId(),
					accountIdStmt,
					user.Identifier,
					user.Firstname,
					user.Lastname,
					user.GetDateofbirth(),
					user.GetJob(),
					user.GetJobGrade(),
					user.Sex,
					user.PhoneNumber,
					user.Height,
					user.Visum,
					user.Playtime,
				).
				ON_DUPLICATE_KEY_UPDATE(
					tUsers.AccountID.SET(mysql.IntExp(mysql.Raw("VALUES(`account_id`)"))),
					tUsers.Firstname.SET(mysql.StringExp(mysql.Raw("VALUES(`firstname`)"))),
					tUsers.Lastname.SET(mysql.StringExp(mysql.Raw("VALUES(`lastname`)"))),
					tUsers.Dateofbirth.SET(mysql.StringExp(mysql.Raw("VALUES(`dateofbirth`)"))),
					tUsers.Job.SET(mysql.StringExp(mysql.Raw("VALUES(`job`)"))),
					tUsers.JobGrade.SET(mysql.IntExp(mysql.Raw("VALUES(`job_grade`)"))),
					tUsers.Sex.SET(mysql.StringExp(mysql.Raw("VALUES(`sex`)"))),
					tUsers.PhoneNumber.SET(mysql.StringExp(mysql.Raw("VALUES(`phone_number`)"))),
					tUsers.Height.SET(mysql.StringExp(mysql.Raw("VALUES(`height`)"))),
					tUsers.Visum.SET(mysql.IntExp(mysql.Raw("VALUES(`visum`)"))),
					tUsers.Playtime.SET(mysql.IntExp(mysql.Raw("VALUES(`playtime`)"))),
				)

			res, err := insertStmt.ExecContext(ctx, tx)
			if err != nil {
				return 0, fmt.Errorf("failed to execute user insert statement. %w", err)
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, fmt.Errorf("failed to retrieve rows affected for user insert. %w", err)
			}

			rowsAffected += rows

			if err := s.handleCitizenLicenses(ctx, tx, user.GetUserId(), user.GetLicenses()); err != nil {
				return 0, fmt.Errorf(
					"failed to handle user licenses for user %d (%s). %w",
					user.GetUserId(),
					user.GetIdentifier(),
					err,
				)
			}

			if err := s.handleCitizensJobs(ctx, tx, user.GetUserId(), user.GetJobs()); err != nil {
				return 0, fmt.Errorf(
					"failed to handle user jobs for user %d (%s). %w",
					user.GetUserId(),
					user.GetIdentifier(),
					err,
				)
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return 0, err
		}
	}

	// TODO insert user's job(s) to fivenet_user_jobs table and phone_number(s) to fivenet_user_phone_numbers table

	if len(toUpdate) > 0 {
		// Begin transaction
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return 0, err
		}
		// Defer a rollback in case anything fails
		defer tx.Rollback()

		for _, user := range toUpdate {
			accountIdStmt := tAccounts.
				SELECT(
					mysql.COALESCE(tAccounts.ID, mysql.NULL),
				).
				FROM(tAccounts).
				WHERE(tAccounts.License.EQ(mysql.String(getLicenseFromIdentifier(user.GetIdentifier())))).
				LIMIT(1)

			stmt := tUsers.
				UPDATE(
					tUsers.AccountID,
					tUsers.Identifier,
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
					accountIdStmt,
					user.Identifier,
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
					tUsers.ID.EQ(mysql.Int32(user.GetUserId())),
				).
				LIMIT(1)

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

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return 0, err
		}
	}

	return rowsAffected, nil
}

func (s *Server) handleCitizenLicenses(
	ctx context.Context,
	tx *sql.Tx,
	userId int32,
	licenses []*userslicenses.License,
) error {
	tCitizensLicenses := table.FivenetUserLicenses

	if len(licenses) == 0 {
		// User has no licenses? Delete all user licenses from the database.
		stmt := tCitizensLicenses.
			DELETE().
			WHERE(tCitizensLicenses.UserID.EQ(mysql.Int32(userId))).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user licenses delete statement. %w", err)
		}

		return nil
	}

	selectStmt := tCitizensLicenses.
		SELECT(
			tCitizensLicenses.Type,
		).
		FROM(tCitizensLicenses).
		WHERE(tCitizensLicenses.UserID.EQ(mysql.Int32(userId)))

	currentLicenses := []string{}
	if err := selectStmt.QueryContext(ctx, tx, &currentLicenses); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf(
				"failed to query current user licenses for user ID %d. %w",
				userId,
				err,
			)
		}
	}

	licensesList := []string{}
	for _, license := range licenses {
		licensesList = append(licensesList, license.GetType())
	}

	toAdd, toRemove := utils.SlicesDifference(currentLicenses, licensesList)

	if len(toAdd) > 0 {
		stmt := tCitizensLicenses.
			INSERT(
				tCitizensLicenses.UserID,
				tCitizensLicenses.Type,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCitizensLicenses.Type.SET(mysql.StringExp(mysql.Raw("VALUES(`type`)"))),
			)

		for _, t := range toAdd {
			stmt = stmt.
				VALUES(
					userId,
					t,
				)
		}

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user licenses insert statement. %w", err)
		}
	}

	if len(toRemove) > 0 {
		types := []mysql.Expression{}
		for _, t := range toRemove {
			types = append(types, mysql.String(t))
		}

		stmt := tCitizensLicenses.
			DELETE().
			WHERE(mysql.AND(
				tCitizensLicenses.UserID.EQ(mysql.Int32(userId)),
				tCitizensLicenses.Type.IN(types...),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user licenses delete statement. %w", err)
		}
	}

	return nil
}

func (s *Server) handleCitizensJobs(
	ctx context.Context,
	tx *sql.Tx,
	userId int32,
	jobs []*users.UserJob,
) error {
	tCitizensJobs := table.FivenetUserJobs

	if len(jobs) == 0 {
		// User has no jobs? Delete all user jobs from the database.
		stmt := tCitizensJobs.
			DELETE().
			WHERE(tCitizensJobs.UserID.EQ(mysql.Int32(userId))).
			LIMIT(25)
		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user jobs delete statement. %w", err)
		}
	} else {
		// TODO
	}

	return nil
}
