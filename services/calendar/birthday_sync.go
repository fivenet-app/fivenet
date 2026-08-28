package calendar

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	calendarstore "github.com/fivenet-app/fivenet/v2026/stores/calendar"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	birthdaySyncCronName                 = "calendar.birthday_sync"
	birthdaySyncBatchSize                = 5
	birthdaySyncOffsetAttrKey            = "job_offset"
	birthdaySyncJobsFetchedAttrKey       = "jobs_fetched"
	birthdaySyncJobsSyncedAttrKey        = "jobs_synced"
	birthdaySyncCalendarsUpsertedAttrKey = "calendars_upserted"
	birthdaySyncEntriesDeletedAttrKey    = "entries_deleted"
	birthdaySyncEntriesInsertedAttrKey   = "entries_inserted"
	birthdaySyncColleaguesLoadedAttrKey  = "colleagues_loaded"
)

type BirthdaySyncer struct {
	logger   *zap.Logger
	db       *sql.DB
	i18n     i18n.Ii18n
	appCfg   appconfig.IConfig
	enricher mstlystcdata.IUserAwareEnricher
	store    calendarstore.IStore
}

type BirthdaySyncerParams struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	I18n      i18n.Ii18n
	AppConfig appconfig.IConfig
	Enricher  mstlystcdata.IUserAwareEnricher
	Store     calendarstore.IStore
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
		store:    p.Store,
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

		stats, nextOffset, finished, err := s.syncBirthdayJobsBatch(ctx, offset)
		if err != nil {
			s.logger.Error("error during birthday sync", zap.Error(err))
			return err
		}

		dest.SetAttribute(birthdaySyncJobsFetchedAttrKey, strconv.Itoa(stats.jobsFetched))
		dest.SetAttribute(birthdaySyncJobsSyncedAttrKey, strconv.Itoa(stats.jobsSynced))
		dest.SetAttribute(
			birthdaySyncCalendarsUpsertedAttrKey,
			strconv.Itoa(stats.calendarsUpserted),
		)
		dest.SetAttribute(
			birthdaySyncEntriesDeletedAttrKey,
			strconv.Itoa(stats.entriesDeleted),
		)
		dest.SetAttribute(
			birthdaySyncEntriesInsertedAttrKey,
			strconv.Itoa(stats.entriesInserted),
		)
		dest.SetAttribute(
			birthdaySyncColleaguesLoadedAttrKey,
			strconv.Itoa(stats.colleaguesLoaded),
		)
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

type birthdaySyncBatchStats struct {
	jobsFetched       int
	jobsSynced        int
	calendarsUpserted int
	entriesDeleted    int
	entriesInserted   int
	colleaguesLoaded  int
}

type birthdaySyncJobStats struct {
	entriesDeleted   int
	entriesInserted  int
	colleaguesLoaded int
}

func (s *BirthdaySyncer) syncBirthdayJobsBatch(
	ctx context.Context,
	offset int,
) (birthdaySyncBatchStats, int, bool, error) {
	jobs, err := s.listBirthdayJobs(ctx, offset, birthdaySyncBatchSize+1)
	if err != nil {
		return birthdaySyncBatchStats{}, 0, false, err
	}

	if len(jobs) == 0 {
		return birthdaySyncBatchStats{}, 0, true, nil
	}

	stats := birthdaySyncBatchStats{
		jobsFetched: len(jobs),
	}
	finished := len(jobs) <= birthdaySyncBatchSize
	if !finished {
		jobs = jobs[:birthdaySyncBatchSize]
	}
	stats.jobsSynced = len(jobs)

	for i := range jobs {
		jobStats, err := s.syncBirthdayJob(ctx, jobs[i])
		if err != nil {
			return birthdaySyncBatchStats{}, 0, false, err
		}
		stats.calendarsUpserted++
		stats.entriesDeleted += jobStats.entriesDeleted
		stats.entriesInserted += jobStats.entriesInserted
		stats.colleaguesLoaded += jobStats.colleaguesLoaded
	}

	return stats, offset + len(jobs), finished, nil
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

func (s *BirthdaySyncer) syncBirthdayJob(
	ctx context.Context,
	job string,
) (birthdaySyncJobStats, error) {
	stats := birthdaySyncJobStats{}
	job = strings.TrimSpace(job)
	if job == "" {
		return stats, nil
	}

	jobInfo := s.enricher.GetJobByName(job)
	title := birthdayCalendarTitle(s.i18n, s.appCfg, job, jobInfo)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return stats, err
	}
	defer tx.Rollback()

	tCalendarEntry := table.FivenetCalendarEntries
	calendarID, err := s.store.UpsertBirthdayCalendar(ctx, tx, job, title)
	if err != nil {
		return stats, err
	}

	if err := s.store.EnsureBirthdayCalendarAccess(ctx, tx, calendarID, job, jobInfo); err != nil {
		return stats, err
	}

	res, err := tCalendarEntry.
		DELETE().
		WHERE(tCalendarEntry.CalendarID.EQ(mysql.Int64(calendarID))).
		ExecContext(ctx, tx)
	if err != nil {
		return stats, err
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		return stats, err
	}
	stats.entriesDeleted = int(deleted)

	colleagues, err := s.store.LoadBirthdayColleagues(ctx, tx, job)
	if err != nil {
		return stats, err
	}
	stats.colleaguesLoaded = len(colleagues)

	for i := range colleagues {
		if err := s.store.InsertBirthdayEntry(ctx, tx, calendarID, job, colleagues[i]); err != nil {
			return stats, err
		}
	}
	stats.entriesInserted = len(colleagues)

	if err := tx.Commit(); err != nil {
		return stats, err
	}

	return stats, nil
}
