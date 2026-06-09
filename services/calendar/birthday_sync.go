package calendar

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	birthdaySyncCronName      = "calendar.birthday_sync"
	birthdaySyncBatchSize     = 5
	birthdaySyncOffsetAttrKey = "job_offset"
)

type BirthdaySyncer struct {
	logger   *zap.Logger
	db       *sql.DB
	i18n     *i18n.I18n
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.UserAwareEnricher
}

type BirthdaySyncerParams struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	I18n      *i18n.I18n
	AppConfig appconfig.IConfig
	Enricher  *mstlystcdata.UserAwareEnricher
}

type BirthdaySyncerResult struct {
	fx.Out

	Syncer       *BirthdaySyncer
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewBirthdaySyncer(p BirthdaySyncerParams) BirthdaySyncerResult {
	s := &BirthdaySyncer{
		logger:   p.Logger.Named("calendar.birthday_sync"),
		db:       p.DB,
		i18n:     p.I18n,
		appCfg:   p.AppConfig,
		enricher: p.Enricher,
	}

	return BirthdaySyncerResult{
		Syncer:       s,
		CronRegister: s,
	}
}

func (s *BirthdaySyncer) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	return registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     birthdaySyncCronName,
		Schedule: "*/4 * * * *",
		Timeout:  durationpb.New(45 * time.Second),
	})
}

func (s *BirthdaySyncer) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add(birthdaySyncCronName, func(ctx context.Context, data *cron.CronjobData) error {
		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if err := data.Unmarshal(dest); err != nil {
			s.logger.Warn("failed to unmarshal birthday sync cron data", zap.Error(err))
		}

		offset, err := strconv.Atoi(strings.TrimSpace(dest.GetAttribute(birthdaySyncOffsetAttrKey)))
		if err != nil || offset < 0 {
			offset = 0
		}

		nextOffset, finished, err := s.syncBirthdayJobsBatch(ctx, offset)
		if err != nil {
			s.logger.Error("error during birthday sync", zap.Error(err))
			return err
		}

		if finished {
			dest.SetAttribute(birthdaySyncOffsetAttrKey, "0")
		} else {
			dest.SetAttribute(birthdaySyncOffsetAttrKey, strconv.Itoa(nextOffset))
		}

		if err := data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated birthday sync cron data. %w", err)
		}

		return nil
	})

	return nil
}

func (s *BirthdaySyncer) syncBirthdayJobsBatch(ctx context.Context, offset int) (int, bool, error) {
	jobs, err := s.listBirthdayJobs(ctx, offset, birthdaySyncBatchSize+1)
	if err != nil {
		return 0, false, err
	}

	if len(jobs) == 0 {
		return 0, true, nil
	}

	finished := len(jobs) <= birthdaySyncBatchSize
	if !finished {
		jobs = jobs[:birthdaySyncBatchSize]
	}

	for i := range jobs {
		if err := s.syncBirthdayJob(ctx, jobs[i]); err != nil {
			return 0, false, err
		}
	}

	return offset + len(jobs), finished, nil
}

func (s *BirthdaySyncer) listBirthdayJobs(
	ctx context.Context,
	offset, limit int,
) ([]string, error) {
	unemployedJob := s.appCfg.Get().GetJobInfo().GetUnemployedJob()

	tJobs := table.FivenetJobs.AS("job")
	stmt := tJobs.
		SELECT(
			tJobs.Name.AS("job"),
		).
		FROM(tJobs).
		WHERE(mysql.AND(
			tJobs.DeletedAt.IS_NULL(),
			tJobs.Name.NOT_IN(
				mysql.String(""),
				mysql.String(unemployedJob.GetName()),
			),
		)).
		ORDER_BY(tJobs.Name.ASC())

	jobs := []string{}
	if err := stmt.
		OFFSET(int64(offset)).
		LIMIT(int64(limit)).
		QueryContext(
			ctx,
			s.db,
			&jobs,
		); err != nil &&
		!errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	for i := range jobs {
		jobs[i] = strings.TrimSpace(jobs[i])
	}

	return jobs, nil
}

func (s *BirthdaySyncer) syncBirthdayJob(ctx context.Context, job string) error {
	job = strings.TrimSpace(job)
	if job == "" {
		return nil
	}

	jobInfo := s.enricher.GetJobByName(job)
	title := birthdayCalendarTitle(s.i18n, s.appCfg, job, jobInfo)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tCalendarEntry := table.FivenetCalendarEntries
	calendarID, err := s.upsertBirthdayCalendar(ctx, tx, job, title)
	if err != nil {
		return err
	}

	if err := ensureBirthdayCalendarAccess(ctx, tx, calendarID, job, jobInfo); err != nil {
		return err
	}

	if _, err := tCalendarEntry.
		DELETE().
		WHERE(tCalendarEntry.CalendarID.EQ(mysql.Int64(calendarID))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	colleagues, err := s.loadBirthdayColleagues(ctx, tx, job)
	if err != nil {
		return err
	}

	for i := range colleagues {
		if err := s.insertBirthdayEntry(ctx, tx, calendarID, job, colleagues[i]); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *BirthdaySyncer) upsertBirthdayCalendar(
	ctx context.Context,
	tx *sql.Tx,
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

func (s *BirthdaySyncer) loadBirthdayColleagues(
	ctx context.Context,
	tx *sql.Tx,
	job string,
) ([]*birthdayColleague, error) {
	tUsers := table.FivenetUser.AS("birthday_colleague")

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
		// Sane limit
		LIMIT(500)

	colleagues := []*birthdayColleague{}
	if err := stmt.QueryContext(ctx, tx, &colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return colleagues, nil
}

func (s *BirthdaySyncer) insertBirthdayEntry(
	ctx context.Context,
	tx *sql.Tx,
	calendarID int64,
	job string,
	colleague *birthdayColleague,
) error {
	birthDate, err := time.Parse("02.01.2006", colleague.Dateofbirth)
	if err != nil {
		return nil
	}

	startTime := birthdayForYear(birthDate.Year(), birthDate.Month(), birthDate.Day())
	recurring := &calendarentries.CalendarEntryRecurring{
		Every: calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_YEAR,
		Count: 1,
	}

	tCalendarEntry := table.FivenetCalendarEntries

	stmt := tCalendarEntry.
		INSERT(
			tCalendarEntry.CalendarID,
			tCalendarEntry.Job,
			tCalendarEntry.StartTime,
			tCalendarEntry.EndTime,
			tCalendarEntry.Title,
			tCalendarEntry.Content,
			tCalendarEntry.Closed,
			tCalendarEntry.RsvpOpen,
			tCalendarEntry.Recurring,
			tCalendarEntry.CreatorID,
			tCalendarEntry.CreatorJob,
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
