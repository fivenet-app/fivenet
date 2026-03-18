package demo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/services/sync"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	demoAccountUsername = "demo"

	targetSeedPrefix = "demotarget"
)

//nolint:gosec // This password is only used when the demo mode is enabled. It is inherently non-sensitive.
var demoAccountPassword = getDemoAccountPassword()

func getDemoAccountPassword() string {
	if v := os.Getenv("FIVENET_DEMO_PASSWORD"); v != "" {
		return v
	}

	return "fivenet-demo"
}

type fakeUserJob struct {
	Job       string
	Grade     int32
	IsPrimary bool
}

type fakeUserProfile struct {
	Identifier      string
	AccountID       *int64
	License         string
	Firstname       string
	Lastname        string
	DateOfBirth     string
	Sex             string
	Height          float64
	PhoneNumber     string
	Visum           int32
	Playtime        int32
	PrimaryJob      string
	PrimaryJobGrade int32
	Jobs            []fakeUserJob
	Licenses        []string
	BloodType       string
}

type fakeUsersGenerator struct{}

func (g fakeUsersGenerator) Name() string {
	return "fake_users"
}

func (g fakeUsersGenerator) Enabled(d *Demo) bool {
	return d.cfg.Demo.Features.Users
}

func (g fakeUsersGenerator) Run(ctx context.Context, d *Demo) error {
	return d.seedFakeUsers(ctx)
}

func (d *Demo) initDemoJobCatalog() {
	jobGrades := map[string][]int32{}
	for _, grade := range demoSeedJobGrades {
		jobGrades[grade.JobName] = append(jobGrades[grade.JobName], grade.Grade)
	}

	demoJobNames := make([]string, 0, len(demoSeedJobs))
	for _, job := range demoSeedJobs {
		grades := uniqueInt32Sorted(jobGrades[job.Name])
		if len(grades) == 0 {
			continue
		}

		demoJobNames = append(demoJobNames, job.Name)
		jobGrades[job.Name] = grades
	}

	d.demoJobNames = demoJobNames
	d.demoJobGrades = jobGrades
}

func uniqueInt32Sorted(in []int32) []int32 {
	if len(in) == 0 {
		return nil
	}

	seen := map[int32]struct{}{}
	out := make([]int32, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	slices.Sort(out)

	return out
}

func (d *Demo) seedFakeUsers(ctx context.Context) error {
	demoCharsCreated, err := d.seedDemoAccountWithChars(ctx)
	if err != nil {
		return err
	}

	count := max(0, d.cfg.Demo.FakeUsers.Count)
	if count == 0 {
		d.logger.Info(
			"no additional demo fake users to seed (count is 0); creating demo account and runtime target-job users only",
		)
	} else {
		availableLicenses, err := d.lookupAvailableLicenseTypes(ctx)
		if err != nil {
			return err
		}

		for i := range count {
			profile := d.buildFakeUserProfile(i+1, "char", availableLicenses)
			if err := d.upsertFakeUser(ctx, profile); err != nil {
				return fmt.Errorf("failed to upsert fake user %s. %w", profile.Identifier, err)
			}
		}
	}

	if err := d.ensureRuntimeTargetJobUsers(ctx, minRuntimeTargetJobUsers); err != nil {
		return err
	}

	d.logger.Info(
		"completed demo fake user seeding",
		zap.Int("count", count),
		zap.Int("demo_account_chars", demoCharsCreated),
	)

	return nil
}

func (d *Demo) getMainCharacterIdentifier() string {
	return d.charIdentifier(1, stableLicenseToken("demo-account", 1))
}

func (d *Demo) ensureRuntimeTargetJobUsers(ctx context.Context, minimum int) error {
	if minimum <= 0 {
		return nil
	}

	targetJobUsers, err := d.lookupUsers(ctx, nil, int64(max(defaultUserLookupLimit, minimum*2)))
	if err != nil {
		return err
	}
	if len(targetJobUsers) >= minimum {
		return nil
	}

	if !d.cfg.Demo.Features.Users {
		d.logger.Warn(
			"not enough target-job users available for demo runtime list",
			zap.Int("required", minimum),
			zap.Int("found", len(targetJobUsers)),
			zap.String("job", d.cfg.Demo.TargetJob),
		)
		return nil
	}

	availableLicenses, err := d.lookupAvailableLicenseTypes(ctx)
	if err != nil {
		return err
	}

	missing := minimum - len(targetJobUsers)
	for i := range missing {
		profile := d.buildTargetJobUserProfile(i+1, availableLicenses)
		if err := d.upsertFakeUser(ctx, profile); err != nil {
			return fmt.Errorf(
				"failed to upsert required target-job demo user %s. %w",
				profile.Identifier,
				err,
			)
		}
	}

	targetJobUsers, err = d.lookupUsers(ctx, nil, int64(max(defaultUserLookupLimit, minimum*2)))
	if err != nil {
		return err
	}
	if len(targetJobUsers) < minimum {
		d.logger.Warn(
			"insufficient target-job users after forced demo seeding",
			zap.Int("required", minimum),
			zap.Int("found", len(targetJobUsers)),
			zap.String("job", d.cfg.Demo.TargetJob),
		)
	}

	return nil
}

func (d *Demo) lookupAvailableLicenseTypes(ctx context.Context) ([]string, error) {
	stmt := tLicenses.
		SELECT(tLicenses.Type).
		FROM(tLicenses).
		ORDER_BY(tLicenses.Type)

	licenses := []string{}
	if err := stmt.QueryContext(ctx, d.db, &licenses); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to lookup available licenses. %w", err)
	}

	return licenses, nil
}

func (d *Demo) seedDemoAccountWithChars(ctx context.Context) (int, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	demoLicense := stableLicenseToken("demo-account", 1)
	targetIdentifier := d.charIdentifier(1, demoLicense)

	accountID, err := d.upsertDemoAccount(ctx, tx, demoLicense)
	if err != nil {
		return 0, err
	}

	profiles := d.buildDemoAccountCharProfiles(accountID, demoLicense)
	var targetUserID int32
	for _, profile := range profiles {
		userID, err := d.upsertFakeUserTx(ctx, tx, profile)
		if err != nil {
			return 0, fmt.Errorf(
				"failed to upsert demo account char %s. %w",
				profile.Identifier,
				err,
			)
		}
		if profile.Identifier == targetIdentifier {
			targetUserID = userID
		}
	}

	if targetUserID == 0 {
		return 0, fmt.Errorf(
			"failed to determine target character user id for %s",
			targetIdentifier,
		)
	}
	if err := d.updateDemoAccountLastChar(ctx, tx, accountID, targetUserID); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(profiles), nil
}

func (d *Demo) upsertDemoAccount(
	ctx context.Context,
	tx *sql.Tx,
	license string,
) (int64, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(demoAccountPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to hash demo account password. %w", err)
	}
	passwordHash := string(passwordHashBytes)

	var existing struct {
		ID int64 `alias:"id"`
	}
	selectStmt := tAccounts.
		SELECT(tAccounts.ID.AS("id")).
		FROM(tAccounts).
		WHERE(tAccounts.Username.EQ(mysql.String(demoAccountUsername))).
		LIMIT(1)

	err = selectStmt.QueryContext(ctx, tx, &existing)
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return 0, fmt.Errorf("failed to lookup demo account. %w", err)
	}

	if errors.Is(err, qrm.ErrNoRows) {
		insertStmt := tAccounts.
			INSERT(
				tAccounts.Enabled,
				tAccounts.Username,
				tAccounts.Password,
				tAccounts.License,
				tAccounts.Groups,
				tAccounts.LastChar,
			).
			VALUES(
				true,
				demoAccountUsername,
				passwordHash,
				license,
				"[\"demo\"]",
				0,
			)

		res, err := insertStmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to insert demo account. %w", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve demo account id. %w", err)
		}

		return id, nil
	}

	updateStmt := tAccounts.
		UPDATE(
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.Password,
			tAccounts.License,
			tAccounts.Groups,
			tAccounts.LastChar,
			tAccounts.RegToken,
		).
		SET(
			true,
			demoAccountUsername,
			passwordHash,
			license,
			"[\"demo\"]",
			0,
			nil,
		).
		WHERE(tAccounts.ID.EQ(mysql.Int64(existing.ID))).
		LIMIT(1)

	if _, err := updateStmt.ExecContext(ctx, tx); err != nil {
		return 0, fmt.Errorf("failed to update demo account. %w", err)
	}

	return existing.ID, nil
}

func (d *Demo) updateDemoAccountLastChar(
	ctx context.Context,
	tx *sql.Tx,
	accountID int64,
	targetUserID int32,
) error {
	updateStmt := tAccounts.
		UPDATE(tAccounts.LastChar).
		SET(targetUserID).
		WHERE(tAccounts.ID.EQ(mysql.Int64(accountID))).
		LIMIT(1)

	if _, err := updateStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to update demo account last char user id. %w", err)
	}

	return nil
}

func (d *Demo) buildDemoAccountCharProfiles(accountID int64, license string) []*fakeUserProfile {
	targetJob := d.targetJobName()

	first := d.newAccountCharProfile(accountID, license, 1, targetJob)
	first.PrimaryJobGrade = d.highestJobGrade(targetJob)
	first.Jobs = []fakeUserJob{
		{Job: targetJob, Grade: first.PrimaryJobGrade, IsPrimary: true},
	}

	secondJob := d.pickNonTargetJob(targetJob, nil, "ambulance", "doj")
	second := d.newAccountCharProfile(accountID, license, 2, secondJob)

	exclude := map[string]struct{}{
		targetJob: {},
		secondJob: {},
	}
	thirdJob := d.pickNonTargetJob(targetJob, exclude, "doj", "mechanic", "ambulance", "cafe")
	third := d.newAccountCharProfile(accountID, license, 3, thirdJob)

	return []*fakeUserProfile{first, second, third}
}

func (d *Demo) newAccountCharProfile(
	accountID int64,
	license string,
	charIndex int,
	job string,
) *fakeUserProfile {
	jobs := []fakeUserJob{
		{
			Job:       job,
			Grade:     d.highestJobGrade(job),
			IsPrimary: true,
		},
	}

	first := strings.TrimSpace(d.fake.FirstName())
	if first == "" {
		first = fmt.Sprintf("Demo%d", charIndex)
	}
	last := strings.TrimSpace(d.fake.LastName())
	if last == "" {
		last = "Account"
	}

	birthDate := d.fake.DateRange(
		time.Date(1960, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2005, time.December, 31, 0, 0, 0, 0, time.UTC),
	)
	sex := "m"
	if d.randIntN(2) == 0 {
		sex = "f"
	}

	identifier := d.charIdentifier(charIndex, license)
	primary := jobs[0]

	return &fakeUserProfile{
		Identifier:      identifier,
		AccountID:       &accountID,
		License:         license,
		Firstname:       first,
		Lastname:        last,
		DateOfBirth:     birthDate.Format("02.01.2006"),
		Sex:             sex,
		Height:          float64(d.randIntN(46) + 160),
		PhoneNumber:     d.fakePhoneNumber(charIndex),
		Visum:           int32(d.randIntN(350) + 1),
		Playtime:        int32(d.randIntN(7_500_000) + 1_000),
		PrimaryJob:      primary.Job,
		PrimaryJobGrade: primary.Grade,
		Jobs:            jobs,
		Licenses:        nil,
		BloodType:       sync.BloodTypes[d.randIntN(len(sync.BloodTypes))],
	}
}

func (d *Demo) highestJobGrade(jobName string) int32 {
	grades := d.demoJobGrades[jobName]
	if len(grades) == 0 {
		return 1
	}
	return grades[len(grades)-1]
}

func (d *Demo) pickNonTargetJob(
	targetJob string,
	exclude map[string]struct{},
	preferred ...string,
) string {
	for _, name := range preferred {
		if name == "" || name == targetJob {
			continue
		}
		if exclude != nil {
			if _, ok := exclude[name]; ok {
				continue
			}
		}
		if len(d.demoJobGrades[name]) > 0 {
			return name
		}
	}

	for _, name := range d.demoJobNames {
		if name == targetJob {
			continue
		}
		if exclude != nil {
			if _, ok := exclude[name]; ok {
				continue
			}
		}
		if len(d.demoJobGrades[name]) > 0 {
			return name
		}
	}

	return "unemployed"
}

func (d *Demo) buildFakeUserProfile(
	index int,
	identifierSeed string,
	availableLicenses []string,
) *fakeUserProfile {
	license := stableLicenseToken(identifierSeed, index)
	identifier := d.charIdentifier(1, license)
	jobs := d.pickUserJobs()
	primaryJob := jobs[0]

	first := strings.TrimSpace(d.fake.FirstName())
	if first == "" {
		first = fmt.Sprintf("Demo%d", index)
	}
	last := strings.TrimSpace(d.fake.LastName())
	if last == "" {
		last = fmt.Sprintf("User%d", index)
	}

	birthDate := d.fake.DateRange(
		time.Date(1960, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2005, time.December, 31, 0, 0, 0, 0, time.UTC),
	)
	sex := "m"
	if d.randIntN(2) == 0 {
		sex = "f"
	}

	return &fakeUserProfile{
		Identifier:      identifier,
		License:         license,
		Firstname:       first,
		Lastname:        last,
		DateOfBirth:     birthDate.Format("02.01.2006"),
		Sex:             sex,
		Height:          float64(d.randIntN(46) + 160),
		PhoneNumber:     d.fakePhoneNumber(index),
		Visum:           int32(d.randIntN(350) + 1),
		Playtime:        int32(d.randIntN(7_500_000) + 1_000),
		PrimaryJob:      primaryJob.Job,
		PrimaryJobGrade: primaryJob.Grade,
		Jobs:            jobs,
		Licenses:        d.pickUserLicenses(availableLicenses),
		BloodType:       sync.BloodTypes[d.randIntN(len(sync.BloodTypes))],
	}
}

func (d *Demo) buildTargetJobUserProfile(index int, availableLicenses []string) *fakeUserProfile {
	targetJob := d.targetJobName()
	seedPrefix := fmt.Sprintf("%s-%s", targetSeedPrefix, targetJob)
	license := stableLicenseToken(seedPrefix, index)
	profile := d.buildFakeUserProfile(index, seedPrefix, availableLicenses)
	profile.Identifier = d.charIdentifier(1, license)
	profile.License = license
	profile.PrimaryJob = targetJob
	profile.PrimaryJobGrade = d.highestJobGrade(targetJob)
	profile.Jobs = []fakeUserJob{
		{
			Job:       targetJob,
			Grade:     profile.PrimaryJobGrade,
			IsPrimary: true,
		},
	}

	return profile
}

func (d *Demo) targetJobName() string {
	targetJob := strings.TrimSpace(d.cfg.Demo.TargetJob)
	if targetJob == "" {
		return "police"
	}
	return targetJob
}

func (d *Demo) charIdentifier(charNumber int, license string) string {
	if charNumber < 1 {
		charNumber = 1
	}

	identifier := fmt.Sprintf("char%d:%s", charNumber, license)
	if len(identifier) > 64 {
		identifier = identifier[:64]
	}
	return identifier
}

func stableLicenseToken(prefix string, index int) string {
	raw := fmt.Sprintf("%s-%d", prefix, index)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (d *Demo) fakePhoneNumber(index int) string {
	prefix := d.fake.Number(200, 999)
	number := fmt.Sprintf("%03d%03d", prefix, index)
	if len(number) > 10 {
		number = number[:10]
	}
	return number
}

func (d *Demo) pickUserJobs() []fakeUserJob {
	if len(d.demoJobNames) == 0 {
		return []fakeUserJob{{Job: "unemployed", Grade: 1, IsPrimary: true}}
	}

	primaryJob := d.demoJobNames[d.randIntN(len(d.demoJobNames))]
	jobs := []fakeUserJob{{Job: primaryJob, Grade: d.randJobGrade(primaryJob), IsPrimary: true}}

	if len(d.demoJobNames) <= 1 || d.randIntN(4) != 0 {
		return jobs
	}

	var secondaryJob string
	for {
		secondaryJob = d.demoJobNames[d.randIntN(len(d.demoJobNames))]
		if secondaryJob != primaryJob {
			break
		}
	}

	jobs = append(
		jobs,
		fakeUserJob{Job: secondaryJob, Grade: d.randJobGrade(secondaryJob), IsPrimary: false},
	)
	return jobs
}

func (d *Demo) randJobGrade(jobName string) int32 {
	grades := d.demoJobGrades[jobName]
	if len(grades) == 0 {
		return 1
	}
	return grades[d.randIntN(len(grades))]
}

func (d *Demo) pickUserLicenses(available []string) []string {
	if len(available) == 0 {
		return nil
	}

	maxLicenses := min(3, len(available))
	count := d.randIntN(maxLicenses) + 1
	perm := d.randPerm(len(available))

	out := make([]string, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, available[perm[i]])
	}

	slices.Sort(out)

	return out
}

func (d *Demo) upsertFakeUser(ctx context.Context, profile *fakeUserProfile) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := d.upsertFakeUserTx(ctx, tx, profile); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Demo) upsertFakeUserTx(
	ctx context.Context,
	tx *sql.Tx,
	profile *fakeUserProfile,
) (int32, error) {
	userID, err := d.upsertFakeUserCore(ctx, tx, profile)
	if err != nil {
		return 0, err
	}

	if err := d.upsertFakeUserJobs(ctx, tx, userID, profile.Jobs); err != nil {
		return 0, err
	}
	if err := d.upsertFakeUserPhoneNumbers(ctx, tx, userID, profile.PhoneNumber); err != nil {
		return 0, err
	}
	if err := d.upsertFakeUserLicenses(ctx, tx, userID, profile.Licenses); err != nil {
		return 0, err
	}
	if err := d.upsertFakeUserProps(ctx, tx, userID, profile); err != nil {
		return 0, err
	}

	return userID, nil
}

func (d *Demo) upsertFakeUserCore(
	ctx context.Context,
	tx *sql.Tx,
	profile *fakeUserProfile,
) (int32, error) {
	var existing struct {
		ID int32 `alias:"id"`
	}

	selectStmt := tUsers.
		SELECT(tUsers.ID.AS("id")).
		FROM(tUsers).
		WHERE(tUsers.Identifier.EQ(mysql.String(profile.Identifier))).
		LIMIT(1)

	err := selectStmt.QueryContext(ctx, tx, &existing)
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return 0, fmt.Errorf("failed to lookup user by identifier. %w", err)
	}

	if errors.Is(err, qrm.ErrNoRows) {
		insertStmt := tUsers.
			INSERT(
				tUsers.AccountID,
				tUsers.License,
				tUsers.Identifier,
				tUsers.Group,
				tUsers.Job,
				tUsers.JobGrade,
				tUsers.Firstname,
				tUsers.Lastname,
				tUsers.Dateofbirth,
				tUsers.Sex,
				tUsers.Height,
				tUsers.PhoneNumber,
				tUsers.Visum,
				tUsers.Playtime,
			).
			VALUES(
				profile.AccountID,
				profile.License,
				profile.Identifier,
				"user",
				profile.PrimaryJob,
				profile.PrimaryJobGrade,
				profile.Firstname,
				profile.Lastname,
				profile.DateOfBirth,
				profile.Sex,
				profile.Height,
				profile.PhoneNumber,
				profile.Visum,
				profile.Playtime,
			)

		res, err := insertStmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, fmt.Errorf("failed to insert fake user. %w", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve fake user id. %w", err)
		}

		return int32(id), nil
	}

	updateStmt := tUsers.
		UPDATE(
			tUsers.AccountID,
			tUsers.License,
			tUsers.Identifier,
			tUsers.Group,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.Sex,
			tUsers.Height,
			tUsers.PhoneNumber,
			tUsers.Visum,
			tUsers.Playtime,
			tUsers.DeletedAt,
			tUsers.DeletedReason,
		).
		SET(
			profile.AccountID,
			profile.License,
			profile.Identifier,
			"user",
			profile.PrimaryJob,
			profile.PrimaryJobGrade,
			profile.Firstname,
			profile.Lastname,
			profile.DateOfBirth,
			profile.Sex,
			profile.Height,
			profile.PhoneNumber,
			profile.Visum,
			profile.Playtime,
			nil,
			nil,
		).
		WHERE(tUsers.ID.EQ(mysql.Int32(existing.ID))).
		LIMIT(1)

	if _, err := updateStmt.ExecContext(ctx, tx); err != nil {
		return 0, fmt.Errorf("failed to update fake user. %w", err)
	}

	return existing.ID, nil
}

func (d *Demo) upsertFakeUserJobs(
	ctx context.Context,
	tx *sql.Tx,
	userID int32,
	jobs []fakeUserJob,
) error {
	if len(jobs) == 0 {
		return nil
	}

	insertStmt := tUserJobs.INSERT(
		tUserJobs.UserID,
		tUserJobs.Job,
		tUserJobs.Grade,
		tUserJobs.IsPrimary,
	)

	jobNames := make([]mysql.Expression, 0, len(jobs))
	for _, job := range jobs {
		jobNames = append(jobNames, mysql.String(job.Job))
		insertStmt = insertStmt.VALUES(userID, job.Job, job.Grade, job.IsPrimary)
	}

	insertStmt = insertStmt.ON_DUPLICATE_KEY_UPDATE(
		tUserJobs.Grade.SET(mysql.RawInt("VALUES(`grade`)")),
		tUserJobs.IsPrimary.SET(mysql.RawBool("VALUES(`is_primary`)")),
	)

	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert fake user jobs. %w", err)
	}

	deleteStmt := tUserJobs.
		DELETE().
		WHERE(mysql.AND(
			tUserJobs.UserID.EQ(mysql.Int32(userID)),
			tUserJobs.Job.NOT_IN(jobNames...),
		)).
		LIMIT(25)

	if _, err := deleteStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to cleanup stale fake user jobs. %w", err)
	}

	return nil
}

func (d *Demo) upsertFakeUserPhoneNumbers(
	ctx context.Context,
	tx *sql.Tx,
	userID int32,
	primaryNumber string,
) error {
	insertStmt := tUserPhoneNumbers.
		INSERT(
			tUserPhoneNumbers.UserID,
			tUserPhoneNumbers.PhoneNumber,
			tUserPhoneNumbers.IsPrimary,
		).
		VALUES(userID, primaryNumber, true).
		ON_DUPLICATE_KEY_UPDATE(
			tUserPhoneNumbers.IsPrimary.SET(mysql.Bool(true)),
		)

	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert fake user primary phone. %w", err)
	}

	deleteStmt := tUserPhoneNumbers.
		DELETE().
		WHERE(mysql.AND(
			tUserPhoneNumbers.UserID.EQ(mysql.Int32(userID)),
			tUserPhoneNumbers.PhoneNumber.NOT_EQ(mysql.String(primaryNumber)),
		)).
		LIMIT(25)

	if _, err := deleteStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to cleanup stale fake user phone numbers. %w", err)
	}

	return nil
}

func (d *Demo) upsertFakeUserLicenses(
	ctx context.Context,
	tx *sql.Tx,
	userID int32,
	licenses []string,
) error {
	if len(licenses) == 0 {
		deleteStmt := tUserLicenses.
			DELETE().
			WHERE(tUserLicenses.UserID.EQ(mysql.Int32(userID))).
			LIMIT(25)
		if _, err := deleteStmt.ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to clear fake user licenses. %w", err)
		}
		return nil
	}

	insertStmt := tUserLicenses.
		INSERT(
			tUserLicenses.UserID,
			tUserLicenses.Type,
		)

	types := make([]mysql.Expression, 0, len(licenses))
	for _, licenseType := range licenses {
		types = append(types, mysql.String(licenseType))
		insertStmt = insertStmt.VALUES(userID, licenseType)
	}

	insertStmt = insertStmt.ON_DUPLICATE_KEY_UPDATE(
		tUserLicenses.Type.SET(mysql.RawString("VALUES(`type`)")),
	)

	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert fake user licenses. %w", err)
	}

	deleteStmt := tUserLicenses.
		DELETE().
		WHERE(mysql.AND(
			tUserLicenses.UserID.EQ(mysql.Int32(userID)),
			tUserLicenses.Type.NOT_IN(types...),
		)).
		LIMIT(25)

	if _, err := deleteStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to cleanup stale fake user licenses. %w", err)
	}

	return nil
}

func (d *Demo) upsertFakeUserProps(
	ctx context.Context,
	tx *sql.Tx,
	userID int32,
	profile *fakeUserProfile,
) error {
	var job *string
	var jobGrade *int32
	if d.randIntN(2) == 0 && profile.PrimaryJob != demoUnemployedJobName {
		job = &profile.PrimaryJob
		jobGrade = &profile.PrimaryJobGrade
	}

	insertStmt := tUserProps.
		INSERT(
			tUserProps.UserID,
			tUserProps.BloodType,
			tUserProps.Job,
			tUserProps.JobGrade,
		).
		VALUES(
			userID,
			profile.BloodType,
			job,
			jobGrade,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tUserProps.BloodType.SET(mysql.RawString("VALUES(`blood_type`)")),
			tUserProps.Job.SET(mysql.RawString("VALUES(`job`)")),
			tUserProps.JobGrade.SET(mysql.RawInt("VALUES(`job_grade`)")),
		)

	if _, err := insertStmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to upsert fake user props. %w", err)
	}

	return nil
}
