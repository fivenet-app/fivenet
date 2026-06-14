package calendarstore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	calendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const birthdayCalendarColor = "neutral"

var ErrInvalidDateofbirth = errors.New("failed to parse dateofbirth")

type BirthdayColleague struct {
	UserID      int32  `alias:"user_id"`
	Firstname   string `alias:"firstname"`
	Lastname    string `alias:"lastname"`
	Dateofbirth string `alias:"dateofbirth"`
}

var birthdaySyncSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(calendaraccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_SHARE),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	},
}

func birthdayCalendarAccessEntries(
	calendarID int64,
	job string,
	jobInfo *jobs.Job,
) []*calendaraccess.CalendarJobAccess {
	minimumGrade := int32(0)
	highestGrade := int32(0)

	if jobInfo != nil && len(jobInfo.GetGrades()) > 0 {
		minimumGrade = jobInfo.GetGrades()[0].GetGrade()
		highestGrade = jobInfo.GetGrades()[len(jobInfo.GetGrades())-1].GetGrade()
	}

	entries := []*calendaraccess.CalendarJobAccess{{
		TargetId:     calendarID,
		Job:          job,
		MinimumGrade: minimumGrade,
		Access:       int32(calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW),
	}}

	if highestGrade > minimumGrade {
		entries = append(entries, &calendaraccess.CalendarJobAccess{
			TargetId:     calendarID,
			Job:          job,
			MinimumGrade: highestGrade,
			Access:       int32(calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT),
		})
		return entries
	}

	entries[0].Access = int32(calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT)
	return entries
}

func birthdayForYear(year int, month time.Month, day int) time.Time {
	if day <= 0 {
		day = 1
	}

	lastDay := time.Date(year, month+1, 0, 12, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		day = lastDay
	}

	return time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
}

func (s *Store) ListBirthdayJobs(
	ctx context.Context,
	offset, limit int,
	unemployedJob string,
) ([]string, error) {
	tJobs := table.FivenetJobs.AS("job")
	unemployedJob = strings.TrimSpace(unemployedJob)

	stmt := tJobs.
		SELECT(
			tJobs.Name.AS("job"),
		).
		FROM(tJobs).
		WHERE(mysql.AND(
			tJobs.DeletedAt.IS_NULL(),
			tJobs.Name.NOT_IN(
				mysql.String(""),
				mysql.String(unemployedJob),
			),
		)).
		ORDER_BY(tJobs.Name.ASC())

	jobs := []string{}
	if err := stmt.
		OFFSET(int64(offset)).
		LIMIT(int64(limit)).
		QueryContext(ctx, s.db, &jobs); err != nil &&
		!errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}

	for i := range jobs {
		jobs[i] = strings.TrimSpace(jobs[i])
	}

	return jobs, nil
}

func (s *Store) EnsureBirthdayCalendarAccess(
	ctx context.Context,
	q qrm.DB,
	calendarID int64,
	job string,
	jobInfo *jobs.Job,
) error {
	jobAccess := birthdayCalendarAccessEntries(calendarID, job, jobInfo)

	_, err := access.NewCalendarSubjectObjectAccess(nil).ReplaceTargetAccess(
		ctx,
		q,
		access.NewSubjectResolver(nil),
		calendarID,
		&calendaraccess.CalendarAccess{Jobs: jobAccess},
		birthdaySyncSubjectAccessOptions,
	)
	return err
}

func (s *Store) UpsertBirthdayCalendar(
	ctx context.Context,
	tx qrm.DB,
	job string,
	title string,
) (int64, error) {
	tCalendar := table.FivenetCalendar

	stmt := tCalendar.
		INSERT(
			tCalendar.Job,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
			tCalendar.SystemKind,
		).
		VALUES(
			mysql.String(job),
			title,
			mysql.String("System-managed birthday calendar"),
			mysql.Bool(false),
			mysql.Bool(true),
			mysql.String(birthdayCalendarColor),
			mysql.NULL,
			mysql.String(job),
			mysql.Int32(int32(calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS)),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCalendar.Name.SET(mysql.String(title)),
			tCalendar.Description.SET(mysql.String("System-managed birthday calendar")),
			tCalendar.Public.SET(mysql.Bool(false)),
			tCalendar.Closed.SET(mysql.Bool(true)),
			tCalendar.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tCalendar.CreatorID.SET(mysql.IntExp(mysql.NULL)),
			tCalendar.CreatorJob.SET(mysql.String(job)),
			tCalendar.SystemKind.SET(mysql.Int32(int32(calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS))),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if lastID > 0 {
		return lastID, nil
	}

	var calendarID struct {
		ID int64
	}
	calendarTable := table.FivenetCalendar

	selectStm := calendarTable.
		SELECT(
			calendarTable.ID.AS("id"),
		).
		FROM(calendarTable).
		WHERE(mysql.AND(
			calendarTable.DeletedAt.IS_NULL(),
			calendarTable.Job.EQ(mysql.String(job)),
			calendarTable.SystemKind.EQ(
				mysql.Int32(int32(calendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS)),
			),
		)).
		LIMIT(1)

	if err := selectStm.QueryContext(ctx, tx, &calendarID); err != nil {
		return 0, err
	}

	return calendarID.ID, nil
}

func (s *Store) DeleteBirthdayEntries(
	ctx context.Context,
	tx qrm.DB,
	calendarID int64,
) error {
	stmt := table.FivenetCalendarEntries.
		DELETE().
		WHERE(table.FivenetCalendarEntries.CalendarID.EQ(mysql.Int64(calendarID)))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadBirthdayColleagues(
	ctx context.Context,
	tx qrm.DB,
	job string,
) ([]*BirthdayColleague, error) {
	tUsers := table.FivenetUser.AS("birthday_colleague")
	tUserJobs := table.FivenetUserJobs.AS("user_jobs")

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("birthday_colleague.user_id"),
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
		).
		FROM(
			tUsers.
				INNER_JOIN(tUserJobs,
					tUserJobs.UserID.EQ(tUsers.ID),
				),
		).
		WHERE(mysql.AND(
			tUserJobs.Job.EQ(mysql.String(job)),
			tUsers.DeletedAt.IS_NULL(),
			tUsers.Dateofbirth.IS_NOT_NULL(),
			tUsers.Firstname.NOT_EQ(mysql.String("")),
			tUsers.Lastname.NOT_EQ(mysql.String("")),
		)).
		ORDER_BY(tUsers.ID.ASC()).
		LIMIT(500)

	colleagues := []*BirthdayColleague{}
	if err := stmt.QueryContext(ctx, tx, &colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return colleagues, nil
}

func (s *Store) InsertBirthdayEntry(
	ctx context.Context,
	tx qrm.DB,
	calendarID int64,
	job string,
	colleague *BirthdayColleague,
) error {
	birthDate, err := time.Parse("02.01.2006", colleague.Dateofbirth)
	if err != nil {
		return ErrInvalidDateofbirth
	}

	startTime := birthdayForYear(birthDate.Year(), birthDate.Month(), birthDate.Day())
	recurring := &calendarentries.CalendarEntryRecurring{
		Every: calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_YEAR,
		Count: 1,
	}

	stmt := table.FivenetCalendarEntries.
		INSERT(
			table.FivenetCalendarEntries.CalendarID,
			table.FivenetCalendarEntries.Job,
			table.FivenetCalendarEntries.StartTime,
			table.FivenetCalendarEntries.EndTime,
			table.FivenetCalendarEntries.Title,
			table.FivenetCalendarEntries.Content,
			table.FivenetCalendarEntries.Closed,
			table.FivenetCalendarEntries.RsvpOpen,
			table.FivenetCalendarEntries.Recurring,
			table.FivenetCalendarEntries.CreatorID,
			table.FivenetCalendarEntries.CreatorJob,
		).
		VALUES(
			calendarID,
			job,
			startTime,
			mysql.NULL,
			fmt.Sprintf("%s %s", colleague.Firstname, colleague.Lastname),
			mysql.NULL,
			true,
			false,
			recurring,
			colleague.UserID,
			job,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
