package sync

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

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
	us := data.Users.GetUsers()
	for i := range us {
		userIds = append(userIds, mysql.Int32(us[i].GetUserId()))

		// Ensure that user has valid jobs and phone_numbers data
		if len(us[i].Jobs) == 0 {
			// If no jobs are set, create one from the user job field
			us[i].Jobs = []*users.UserJob{
				{
					Job:       us[i].GetJob(),
					Grade:     us[i].GetJobGrade(),
					IsPrimary: true,
				},
			}
		} else {
			// Sort the user's jobs by is primary and then alphabetically to ensure consistent order
			slices.SortFunc(us[i].GetJobs(), func(a *users.UserJob, b *users.UserJob) int {
				if a.GetIsPrimary() && !b.GetIsPrimary() {
					return -1
				}
				if !a.GetIsPrimary() && b.GetIsPrimary() {
					return 1
				}

				return strings.Compare(a.GetJob(), b.GetJob())
			})

			foundPrimary := false
			primaryJob := us[i].GetJob()
			for _, job := range us[i].GetJobs() {
				if job.GetJob() == primaryJob {
					// Make sure the "primary" job (user's job field if set) is marked as primary
					foundPrimary = true
					job.IsPrimary = true
				} else {
					job.IsPrimary = false
				}
			}

			// If not ensure user has at least one primary job set
			if !foundPrimary {
				us[i].Jobs[0].IsPrimary = true
			}
		}

		if len(us[i].PhoneNumbers) == 0 {
			// If no phone numbers are set, create one from the user phone number field if it exists
			if us[i].GetPhoneNumber() != "" {
				us[i].PhoneNumbers = []*users.PhoneNumber{
					{
						Number:    us[i].GetPhoneNumber(),
						IsPrimary: true,
					},
				}
			} else {
				foundPrimary := false

				for _, phoneNumber := range us[i].GetPhoneNumbers() {
					if phoneNumber.GetNumber() == us[i].GetPhoneNumber() {
						foundPrimary = true
						phoneNumber.IsPrimary = true
					} else {
						phoneNumber.IsPrimary = false
					}
				}

				// If not ensure user has at least one primary phone number set
				if !foundPrimary && len(us[i].GetPhoneNumbers()) > 0 {
					us[i].PhoneNumbers[0].IsPrimary = true
				}
			}
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

	rowsAffected := int64(0)
	if len(toCreate) > 0 {
		for _, user := range toCreate {
			affected, err := s.createUser(ctx, user)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to create user %d (%s). %w",
					user.GetUserId(),
					user.GetIdentifier(),
					err,
				)
			}
			rowsAffected += affected
		}
	}

	if len(toUpdate) > 0 {
		for _, user := range toCreate {
			affected, err := s.updateUser(ctx, user)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to update user %d (%s). %w",
					user.GetUserId(),
					user.GetIdentifier(),
					err,
				)
			}
			rowsAffected += affected
		}
	}

	return rowsAffected, nil
}

func (s *Server) createUser(
	ctx context.Context,
	user *syncdata.DataUser,
) (int64, error) {
	tAccounts := table.FivenetAccounts
	tUsers := table.FivenetUser

	stmt := tUsers.
		INSERT(
			tUsers.ID,
			tUsers.AccountID,
			tUsers.License,
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

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	var accountIdStmt mysql.SelectStatement = nil
	if user.GetUserId() <= 0 || user.GetIdentifier() != "" {
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
			utils.GetLicenseFromIdentifier(user.Identifier),
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
			tUsers.Height.SET(mysql.FloatExp(mysql.Raw("VALUES(`height`)"))),
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

	if err := s.handleUserJobs(ctx, tx, user.GetUserId(), user.GetJobs()); err != nil {
		return 0, fmt.Errorf(
			"failed to handle user jobs for user %d (%s). %w",
			user.GetUserId(),
			user.GetIdentifier(),
			err,
		)
	}

	if err := s.handleUserLicenses(ctx, tx, user.GetUserId(), user.GetLicenses()); err != nil {
		return 0, fmt.Errorf(
			"failed to handle user licenses for user %d (%s). %w",
			user.GetUserId(),
			user.GetIdentifier(),
			err,
		)
	}

	if err := s.handleUserPhoneNumbers(ctx, tx, user.GetUserId(), user.GetPhoneNumbers()); err != nil {
		return 0, fmt.Errorf(
			"failed to handle user phone numbers for user %d (%s). %w",
			user.GetUserId(),
			user.GetIdentifier(),
			err,
		)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *Server) updateUser(
	ctx context.Context,
	user *syncdata.DataUser,
) (int64, error) {
	tAccounts := table.FivenetAccounts
	tUsers := table.FivenetUser

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	accountIdStmt := tAccounts.
		SELECT(
			mysql.COALESCE(tAccounts.ID, mysql.NULL),
		).
		FROM(tAccounts).
		WHERE(tAccounts.License.EQ(mysql.String(getLicenseFromIdentifier(user.GetIdentifier())))).
		LIMIT(1)

	stmt := tUsers.
		UPDATE(
			tUsers.ID,
			tUsers.AccountID,
			tUsers.License,
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
			user.GetUserId(),
			accountIdStmt,
			utils.GetLicenseFromIdentifier(user.Identifier),
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
		WHERE(mysql.OR(
			tUsers.ID.EQ(mysql.Int32(user.GetUserId())),
			tUsers.Identifier.EQ(mysql.String(user.GetIdentifier())),
		)).
		LIMIT(1)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("failed to execute user update statement. %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected for user update. %w", err)
	}

	if err := s.handleUserJobs(ctx, tx, user.GetUserId(), user.GetJobs()); err != nil {
		return 0, fmt.Errorf(
			"failed to handle user jobs for user %d (%s). %w",
			user.GetUserId(),
			user.GetIdentifier(),
			err,
		)
	}

	if err := s.handleUserLicenses(ctx, tx, user.GetUserId(), user.GetLicenses()); err != nil {
		return 0, fmt.Errorf(
			"failed to handle user licenses for user %d (%s). %w",
			user.GetUserId(),
			user.GetIdentifier(),
			err,
		)
	}

	if err := s.handleUserPhoneNumbers(ctx, tx, user.GetUserId(), user.GetPhoneNumbers()); err != nil {
		return 0, fmt.Errorf(
			"failed to handle user phone numbers for user %d (%s). %w",
			user.GetUserId(),
			user.GetIdentifier(),
			err,
		)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rows, nil
}

func (s *Server) handleUserLicenses(
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
		WHERE(tCitizensLicenses.UserID.EQ(mysql.Int32(userId))).
		ORDER_BY(tCitizensLicenses.Type)

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

func (s *Server) handleUserJobs(
	ctx context.Context,
	tx *sql.Tx,
	userId int32,
	jobs []*users.UserJob,
) error {
	// Make sure to sort array by is primary and then alphabetically to ensure consistent order
	slices.SortFunc(jobs, func(a *users.UserJob, b *users.UserJob) int {
		if a.GetIsPrimary() && !b.GetIsPrimary() {
			return -1
		}
		if !a.GetIsPrimary() && b.GetIsPrimary() {
			return 1
		}

		return strings.Compare(a.GetJob(), b.GetJob())
	})

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

		return nil
	}

	var selectStmt mysql.SelectStatement
	{
		tCitizensJobs := tCitizensJobs.AS("user_job")
		selectStmt = tCitizensJobs.
			SELECT(
				tCitizensJobs.Job,
				tCitizensJobs.Grade,
				tCitizensJobs.IsPrimary,
			).
			FROM(tCitizensJobs).
			WHERE(tCitizensJobs.UserID.EQ(mysql.Int32(userId))).
			ORDER_BY(tCitizensJobs.IsPrimary, tCitizensJobs.Job, tCitizensJobs.Grade)
	}

	currentJobs := []*users.UserJob{}
	if err := selectStmt.QueryContext(ctx, tx, &currentJobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf(
				"failed to query current user jobs for user ID %d. %w",
				userId,
				err,
			)
		}
	}

	toAdd, toUpdate, toRemove := compareJobs(currentJobs, jobs)

	if len(toAdd) > 0 || len(toUpdate) > 0 {
		stmt := tCitizensJobs.
			INSERT(
				tCitizensJobs.UserID,
				tCitizensJobs.Job,
				tCitizensJobs.Grade,
				tCitizensJobs.IsPrimary,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCitizensJobs.Job.SET(mysql.StringExp(mysql.Raw("VALUES(`job`)"))),
				tCitizensJobs.Grade.SET(mysql.IntExp(mysql.Raw("VALUES(`grade`)"))),
				tCitizensJobs.IsPrimary.SET(mysql.BoolExp(mysql.Raw("VALUES(`is_primary`)"))),
			)

		for _, t := range append(toAdd, toUpdate...) {
			stmt = stmt.
				VALUES(
					userId,
					t.GetJob(),
					t.GetGrade(),
					t.GetIsPrimary(),
				)
		}

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user jobs insert statement. %w", err)
		}
	}

	if len(toRemove) > 0 {
		jobs := []mysql.Expression{}
		for _, t := range toRemove {
			jobs = append(jobs, mysql.String(t.GetJob()))
		}

		stmt := tCitizensJobs.
			DELETE().
			WHERE(mysql.AND(
				tCitizensJobs.UserID.EQ(mysql.Int32(userId)),
				tCitizensJobs.Job.IN(jobs...),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user jobs delete statement. %w", err)
		}
	}

	return nil
}

// compareJobs compares currentJobs and jobs, returning toAdd, toUpdate, and toRemove.
func compareJobs(currentJobs, jobs []*users.UserJob) (toAdd, toUpdate, toRemove []*users.UserJob) {
	// Create a map for current jobs by job name
	currentJobsMap := make(map[string]*users.UserJob)
	for _, job := range currentJobs {
		currentJobsMap[job.GetJob()] = job
	}

	// Create a map for incoming jobs by job name
	incomingJobsMap := make(map[string]*users.UserJob)
	for _, job := range jobs {
		incomingJobsMap[job.GetJob()] = job
	}

	// Determine toAdd and toUpdate
	for jobName, incomingJob := range incomingJobsMap {
		if currentJob, exists := currentJobsMap[jobName]; exists {
			// Check if the job needs an update (ignoring grade)
			if currentJob.GetIsPrimary() != incomingJob.GetIsPrimary() {
				toUpdate = append(toUpdate, incomingJob)
			}
		} else {
			// Job does not exist in current jobs, add it
			toAdd = append(toAdd, incomingJob)
		}
	}

	// Determine toRemove
	for jobName, currentJob := range currentJobsMap {
		if _, exists := incomingJobsMap[jobName]; !exists {
			toRemove = append(toRemove, currentJob)
		}
	}

	return toAdd, toUpdate, toRemove
}

func (s *Server) handleUserPhoneNumbers(
	ctx context.Context,
	tx *sql.Tx,
	userId int32,
	phoneNumbers []*users.PhoneNumber,
) error {
	// Make sure to sort array by is primary and then alphabetically to ensure consistent order
	slices.SortFunc(phoneNumbers, func(a *users.PhoneNumber, b *users.PhoneNumber) int {
		if a.GetIsPrimary() && !b.GetIsPrimary() {
			return -1
		}
		if !a.GetIsPrimary() && b.GetIsPrimary() {
			return 1
		}

		return strings.Compare(a.GetNumber(), b.GetNumber())
	})

	tCitizensPhoneNumbers := table.FivenetUserPhoneNumbers

	if len(phoneNumbers) == 0 {
		// User has no phone numbers? Delete all user phone numbers from the database.
		stmt := tCitizensPhoneNumbers.
			DELETE().
			WHERE(tCitizensPhoneNumbers.UserID.EQ(mysql.Int32(userId))).
			LIMIT(25)
		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user phone numbers delete statement. %w", err)
		}

		return nil
	}

	var selectStmt mysql.SelectStatement
	{
		tCitizensPhoneNumbers := tCitizensPhoneNumbers.AS("phone_number")
		selectStmt = tCitizensPhoneNumbers.
			SELECT(
				tCitizensPhoneNumbers.UserID,
				tCitizensPhoneNumbers.PhoneNumber,
				tCitizensPhoneNumbers.IsPrimary,
			).
			FROM(tCitizensPhoneNumbers).
			WHERE(tCitizensPhoneNumbers.UserID.EQ(mysql.Int32(userId))).
			ORDER_BY(tCitizensPhoneNumbers.IsPrimary, tCitizensPhoneNumbers.PhoneNumber)
	}

	currentPhoneNumbers := []*users.PhoneNumber{}
	if err := selectStmt.QueryContext(ctx, tx, &currentPhoneNumbers); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf(
				"failed to query current user phone numbers for user ID %d. %w",
				userId,
				err,
			)
		}
	}

	toAdd, toUpdate, toRemove := comparePhoneNumbers(currentPhoneNumbers, phoneNumbers)

	if len(toAdd) > 0 || len(toUpdate) > 0 {
		stmt := tCitizensPhoneNumbers.
			INSERT(
				tCitizensPhoneNumbers.UserID,
				tCitizensPhoneNumbers.PhoneNumber,
				tCitizensPhoneNumbers.IsPrimary,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tCitizensPhoneNumbers.PhoneNumber.SET(mysql.StringExp(mysql.Raw("VALUES(`phone_number`)"))),
				tCitizensPhoneNumbers.IsPrimary.SET(mysql.BoolExp(mysql.Raw("VALUES(`is_primary`)"))),
			)

		phoneNumbersToAdd := []*users.PhoneNumber{}
		for _, phoneNumber := range phoneNumbers {
			if slices.ContainsFunc(toAdd, func(t *users.PhoneNumber) bool {
				return t.GetNumber() == phoneNumber.GetNumber()
			}) {
				phoneNumbersToAdd = append(phoneNumbersToAdd, phoneNumber)
			}
		}

		for _, t := range phoneNumbersToAdd {
			stmt = stmt.
				VALUES(
					userId,
					t.GetNumber(),
					t.GetIsPrimary(),
				)
		}

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user phone numbers insert statement. %w", err)
		}
	}

	if len(toRemove) > 0 {
		phoneNumbers := []mysql.Expression{}
		for _, t := range toRemove {
			phoneNumbers = append(phoneNumbers, mysql.String(t.GetNumber()))
		}

		stmt := tCitizensPhoneNumbers.
			DELETE().
			WHERE(mysql.AND(
				tCitizensPhoneNumbers.UserID.EQ(mysql.Int32(userId)),
				tCitizensPhoneNumbers.PhoneNumber.IN(phoneNumbers...),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to execute user phone numbers delete statement. %w", err)
		}
	}

	return nil
}

// comparePhoneNumbers compares currentPhoneNumbers and phoneNumbers, ensuring only one primary number, and returns toAdd, toUpdate, and toRemove.
func comparePhoneNumbers(
	currentPhoneNumbers, phoneNumbers []*users.PhoneNumber,
) (toAdd, toUpdate, toRemove []*users.PhoneNumber) {
	// Create a map for current phone numbers by number
	currentPhoneNumbersMap := make(map[string]*users.PhoneNumber)
	for _, phoneNumber := range currentPhoneNumbers {
		currentPhoneNumbersMap[phoneNumber.GetNumber()] = phoneNumber
	}

	// Create a map for incoming phone numbers by number
	incomingPhoneNumbersMap := make(map[string]*users.PhoneNumber)
	for _, phoneNumber := range phoneNumbers {
		incomingPhoneNumbersMap[phoneNumber.GetNumber()] = phoneNumber
	}

	// Track the new primary number
	var newPrimary *users.PhoneNumber
	for _, phoneNumber := range phoneNumbers {
		if phoneNumber.GetIsPrimary() {
			newPrimary = phoneNumber
			break
		}
	}

	// Determine toAdd and toUpdate
	for number, incomingPhoneNumber := range incomingPhoneNumbersMap {
		if currentPhoneNumber, exists := currentPhoneNumbersMap[number]; exists {
			// Check if the phone number needs an update
			if currentPhoneNumber.GetIsPrimary() != incomingPhoneNumber.GetIsPrimary() {
				toUpdate = append(toUpdate, incomingPhoneNumber)
			}
		} else {
			// Phone number does not exist in current phone numbers, add it
			toAdd = append(toAdd, incomingPhoneNumber)
		}
	}

	// Ensure only the new primary number is marked as primary
	if newPrimary != nil {
		for _, currentPhoneNumber := range currentPhoneNumbers {
			if currentPhoneNumber.GetIsPrimary() &&
				currentPhoneNumber.GetNumber() != newPrimary.GetNumber() {
				// Mark the old primary number as non-primary
				currentPhoneNumber.IsPrimary = false
				toUpdate = append(toUpdate, currentPhoneNumber)
			}
		}
	}

	// Determine toRemove
	for number, currentPhoneNumber := range currentPhoneNumbersMap {
		if _, exists := incomingPhoneNumbersMap[number]; !exists {
			toRemove = append(toRemove, currentPhoneNumber)
		}
	}

	return toAdd, toUpdate, toRemove
}
