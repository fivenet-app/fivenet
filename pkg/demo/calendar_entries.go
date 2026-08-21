package demo

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	commoncontent "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

const (
	demoCalendarEntryPrefix = "[DEMO]"

	demoCalendarFallbackColor = "#2563eb"
	demoCalendarFallbackName  = "Demo %s Calendar"
)

var (
	demoShortMeetingTitles = []string{
		"Shift briefing",
		"Operations sync",
		"Radio check",
		"Case review",
		"Unit debrief",
		"Crew handoff",
		"Planning slot",
		"Stand-up meeting",
		"Follow-up call",
		"Supervisor update",
	}
	demoAllDayTitles = []string{
		"Training day",
		"Command stand-by",
		"Admin catch-up",
		"Public outreach",
		"Audit prep",
		"Regional conference",
	}
	demoMultiDayTitles = []string{
		"Strategy offsite",
		"Joint training exercise",
		"Weekend operations summit",
	}
	demoMeetingTopics = []string{
		"staffing",
		"incident review",
		"deployment planning",
		"resource allocation",
		"shift transitions",
		"communications",
	}
)

type demoCalendarEntriesGenerator struct{}

func (g demoCalendarEntriesGenerator) Name() string {
	return "calendar_entries"
}

func (g demoCalendarEntriesGenerator) Enabled(d *Demo) bool {
	return d.cfg.Demo.Features.CalendarEntries
}

func (g demoCalendarEntriesGenerator) Run(ctx context.Context, d *Demo) error {
	return d.seedDemoCalendarEntries(ctx)
}

func (d *Demo) seedDemoCalendarEntries(ctx context.Context) error {
	if d.calendars == nil {
		return errors.New("failed to seed demo calendar entries: calendar store is not available")
	}

	job := d.targetJobName()
	calendar, err := d.lookupDemoCalendar(ctx, job)
	if err != nil {
		return err
	}
	if calendar == nil {
		seedUser, err := d.lookupCalendarSeedUser(ctx, job)
		if err != nil {
			return err
		}
		if seedUser == nil {
			d.logger.Warn(
				"skipping demo calendar entry seeding, no suitable creator user was found",
				zap.String("job", job),
			)
			return nil
		}

		calendar, err = d.createFallbackDemoCalendar(ctx, job, seedUser)
		if err != nil {
			return err
		}
	}

	seedUser, err := d.lookupCalendarSeedUser(ctx, job)
	if err != nil {
		return err
	}
	if seedUser == nil {
		if calendar.GetCreatorId() > 0 {
			seedUser = &userinfo.UserInfo{
				UserId:    calendar.GetCreatorId(),
				Job:       job,
				JobGrade:  d.highestJobGrade(job),
				Superuser: false,
			}
		} else {
			d.logger.Warn(
				"skipping demo calendar entry seeding, no suitable creator user was found",
				zap.String("job", job),
				zap.Int64("calendar_id", calendar.GetId()),
			)
			return nil
		}
	}

	return d.seedDemoCalendarEntriesForCalendar(ctx, calendar.GetId(), seedUser, time.Now().In(time.Local))
}

func (d *Demo) seedDemoCalendarEntriesForCalendar(
	ctx context.Context,
	calendarID int64,
	seedUser *userinfo.UserInfo,
	now time.Time,
) error {
	entries := d.buildDemoCalendarEntriesAt(calendarID, now)
	if len(entries) == 0 {
		return nil
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := d.clearDemoCalendarEntries(ctx, tx, calendarID); err != nil {
		return err
	}

	for _, entry := range entries {
		if _, err := d.calendars.UpsertCalendarEntry(ctx, tx, entry, nil, seedUser); err != nil {
			return fmt.Errorf("failed to upsert demo calendar entry. %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	d.logger.Info(
		"completed demo calendar entry seeding",
		zap.Int64("calendar_id", calendarID),
		zap.String("job", seedUser.GetJob()),
		zap.Int("count", len(entries)),
	)

	return nil
}

func (d *Demo) lookupDemoCalendar(
	ctx context.Context,
	job string,
) (*calendarresource.Calendar, error) {
	name := fmt.Sprintf(demoCalendarFallbackName, titleizeJob(job))

	if cal, err := d.lookupDemoCalendarByName(ctx, job, name); err != nil {
		return nil, err
	} else if cal != nil {
		return cal, nil
	}

	return d.lookupAnyDemoCalendar(ctx, job)
}

func (d *Demo) lookupDemoCalendarByName(
	ctx context.Context,
	job string,
	name string,
) (*calendarresource.Calendar, error) {
	stmt := d.demoCalendarLookupByNameStmt(job, name)

	dest := &calendarresource.Calendar{}
	if err := stmt.QueryContext(ctx, d.db, dest); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, fmt.Errorf("failed to lookup demo calendar by name. %w", err)
	}
	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (d *Demo) lookupAnyDemoCalendar(
	ctx context.Context,
	job string,
) (*calendarresource.Calendar, error) {
	stmt := d.demoAnyCalendarLookupStmt(job)

	dest := &calendarresource.Calendar{}
	if err := stmt.QueryContext(ctx, d.db, dest); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, fmt.Errorf("failed to lookup demo calendar. %w", err)
	}
	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (d *Demo) demoCalendarLookupByNameStmt(job string, name string) mysql.SelectStatement {
	return table.FivenetCalendar.
		SELECT(
			table.FivenetCalendar.ID,
			table.FivenetCalendar.CreatedAt,
			table.FivenetCalendar.UpdatedAt,
			table.FivenetCalendar.DeletedAt,
			table.FivenetCalendar.Job,
			table.FivenetCalendar.Name,
			table.FivenetCalendar.Public,
			table.FivenetCalendar.Closed,
			table.FivenetCalendar.Color,
			table.FivenetCalendar.CreatorID,
			table.FivenetCalendar.CreatorJob,
			table.FivenetCalendar.SystemKind,
		).
		FROM(table.FivenetCalendar).
		WHERE(mysql.AND(
			table.FivenetCalendar.DeletedAt.IS_NULL(),
			table.FivenetCalendar.Job.EQ(mysql.String(job)),
			table.FivenetCalendar.Name.EQ(mysql.String(name)),
		)).
		ORDER_BY(
			table.FivenetCalendar.Name.ASC(),
			table.FivenetCalendar.ID.ASC(),
		).
		LIMIT(1)
}

func (d *Demo) demoAnyCalendarLookupStmt(job string) mysql.SelectStatement {
	return table.FivenetCalendar.
		SELECT(
			table.FivenetCalendar.ID,
			table.FivenetCalendar.CreatedAt,
			table.FivenetCalendar.UpdatedAt,
			table.FivenetCalendar.DeletedAt,
			table.FivenetCalendar.Job,
			table.FivenetCalendar.Name,
			table.FivenetCalendar.Public,
			table.FivenetCalendar.Closed,
			table.FivenetCalendar.Color,
			table.FivenetCalendar.CreatorID,
			table.FivenetCalendar.CreatorJob,
			table.FivenetCalendar.SystemKind,
		).
		FROM(table.FivenetCalendar).
		WHERE(mysql.AND(
			table.FivenetCalendar.DeletedAt.IS_NULL(),
			table.FivenetCalendar.Job.EQ(mysql.String(job)),
		)).
		ORDER_BY(
			table.FivenetCalendar.Name.ASC(),
			table.FivenetCalendar.ID.ASC(),
		).
		LIMIT(1)
}

func (d *Demo) createFallbackDemoCalendar(
	ctx context.Context,
	job string,
	seedUser *userinfo.UserInfo,
) (*calendarresource.Calendar, error) {
	if seedUser == nil {
		return nil, nil
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cal := &calendarresource.Calendar{
		Job:    &job,
		Name:   fmt.Sprintf(demoCalendarFallbackName, titleizeJob(job)),
		Public: true,
		Closed: false,
		Color:  demoCalendarFallbackColor,
	}

	lastID, err := d.calendars.CreateCalendar(ctx, tx, cal, seedUser, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fallback demo calendar. %w", err)
	}
	cal.SetId(lastID)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	d.logger.Info(
		"created fallback demo calendar",
		zap.Int64("calendar_id", cal.GetId()),
		zap.String("job", job),
	)

	return cal, nil
}

func (d *Demo) lookupCalendarSeedUser(
	ctx context.Context,
	job string,
) (*userinfo.UserInfo, error) {
	queries := []mysql.SelectStatement{
		table.FivenetUser.
			SELECT(
				table.FivenetUser.ID.AS("id"),
				table.FivenetUser.Job.AS("job"),
				table.FivenetUser.JobGrade.AS("grade"),
			).
			FROM(table.FivenetUser).
			WHERE(table.FivenetUser.Identifier.EQ(mysql.String(d.getMainCharacterIdentifier()))).
			LIMIT(1),
		table.FivenetUser.
			SELECT(
				table.FivenetUser.ID.AS("id"),
				table.FivenetUser.Job.AS("job"),
				table.FivenetUser.JobGrade.AS("grade"),
			).
			FROM(table.FivenetUser).
			WHERE(table.FivenetUser.Job.EQ(mysql.String(job))).
			ORDER_BY(
				table.FivenetUser.JobGrade.DESC(),
				table.FivenetUser.ID.ASC(),
			).
			LIMIT(1),
		table.FivenetUser.
			SELECT(
				table.FivenetUser.ID.AS("id"),
				table.FivenetUser.Job.AS("job"),
				table.FivenetUser.JobGrade.AS("grade"),
			).
			FROM(table.FivenetUser).
			ORDER_BY(table.FivenetUser.ID.ASC()).
			LIMIT(1),
	}

	for _, stmt := range queries {
		var row struct {
			ID    int32  `alias:"id"`
			Job   string `alias:"job"`
			Grade int32  `alias:"grade"`
		}
		if err := stmt.QueryContext(ctx, d.db, &row); err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to lookup demo calendar seed user. %w", err)
		}
		if row.ID <= 0 {
			continue
		}

		return &userinfo.UserInfo{
			UserId:    row.ID,
			Job:       job,
			JobGrade:  d.highestJobGrade(job),
			Superuser: false,
		}, nil
	}

	return nil, nil
}

func (d *Demo) clearDemoCalendarEntries(
	ctx context.Context,
	tx qrm.DB,
	calendarID int64,
) error {
	stmt := d.clearDemoCalendarEntriesStmt(calendarID)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return fmt.Errorf("failed to clear demo calendar entries. %w", err)
	}

	return nil
}

func (d *Demo) clearDemoCalendarEntriesStmt(calendarID int64) mysql.DeleteStatement {
	return table.FivenetCalendarEntries.
		DELETE().
		WHERE(mysql.AND(
			table.FivenetCalendarEntries.CalendarID.EQ(mysql.Int64(calendarID)),
			table.FivenetCalendarEntries.Content.LIKE(mysql.String("%"+demoCalendarEntryPrefix+"%")),
		))
}

func (d *Demo) buildDemoCalendarEntriesAt(calendarID int64, now time.Time) []*calendarentries.CalendarEntry {
	total := 20 + d.randIntN(21)
	multiDayCount := min(2, max(1, total/20))
	allDayCount := min(7, max(4, total/5))
	shortCount := total - allDayCount - multiDayCount
	if shortCount < 1 {
		shortCount = 1
		allDayCount = max(1, total-shortCount-multiDayCount)
	}

	entries := make([]*calendarentries.CalendarEntry, 0, total)
	for i := range shortCount {
		entries = append(entries, d.newDemoShortMeetingEntry(calendarID, now, i))
	}
	for i := range allDayCount {
		entries = append(entries, d.newDemoAllDayEntry(calendarID, now, i))
	}
	for i := range multiDayCount {
		entries = append(entries, d.newDemoMultiDayEntry(calendarID, now, i))
	}

	if len(entries) <= 1 {
		return entries
	}

	perm := d.randInts(len(entries))
	out := make([]*calendarentries.CalendarEntry, 0, len(entries))
	for _, idx := range perm {
		out = append(out, entries[idx])
	}

	return out
}

func (d *Demo) newDemoShortMeetingEntry(
	calendarID int64,
	now time.Time,
	index int,
) *calendarentries.CalendarEntry {
	title := demoShortMeetingTitles[(index+d.randIntN(len(demoShortMeetingTitles)))%len(demoShortMeetingTitles)]
	topic := demoMeetingTopics[d.randIntN(len(demoMeetingTopics))]
	start := d.randomMeetingStart(d.randIntN(2) == 0, now, 30)
	duration := time.Duration(30+d.randIntN(4)*15) * time.Minute
	end := start.Add(duration)

	return d.newDemoCalendarEntry(
		calendarID,
		title,
		fmt.Sprintf("%s about %s.", title, topic),
		start,
		end,
		false,
	)
}

func (d *Demo) newDemoAllDayEntry(
	calendarID int64,
	now time.Time,
	index int,
) *calendarentries.CalendarEntry {
	title := demoAllDayTitles[(index+d.randIntN(len(demoAllDayTitles)))%len(demoAllDayTitles)]
	start := d.randomDayBoundary(d.randIntN(2) == 0, now, 30)
	end := start.AddDate(0, 0, 1)

	return d.newDemoCalendarEntry(
		calendarID,
		title,
		fmt.Sprintf("%s blocks the full day.", title),
		start,
		end,
		true,
	)
}

func (d *Demo) newDemoMultiDayEntry(
	calendarID int64,
	now time.Time,
	index int,
) *calendarentries.CalendarEntry {
	title := demoMultiDayTitles[(index+d.randIntN(len(demoMultiDayTitles)))%len(demoMultiDayTitles)]
	days := 2 + d.randIntN(3)
	start := d.randomDayBoundary(d.randIntN(2) == 0, now, 60-int64(days))
	end := start.AddDate(0, 0, days)

	return d.newDemoCalendarEntry(
		calendarID,
		title,
		fmt.Sprintf("%s runs for %d days.", title, days),
		start,
		end,
		true,
	)
}

func (d *Demo) newDemoCalendarEntry(
	calendarID int64,
	title string,
	summary string,
	start time.Time,
	end time.Time,
	allDay bool,
) *calendarentries.CalendarEntry {
	entry := &calendarentries.CalendarEntry{
		CalendarId:        calendarID,
		StartTime:         timestamp.New(start),
		Title:             title,
		Content:           newDemoCalendarEntryContent(summary),
		Closed:            false,
		AllDay:            allDay,
		RecurrenceVersion: 1,
	}

	if !end.IsZero() {
		entry.EndTime = timestamp.New(end)
	}

	return entry
}

func newDemoCalendarEntryContent(summary string) *commoncontent.Content {
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return nil
	}

	rawHTML := fmt.Sprintf(
		"<p>%s</p><p><em>%s</em></p>",
		html.EscapeString(summary),
		html.EscapeString(demoCalendarEntryPrefix),
	)
	return &commoncontent.Content{
		Version:     commoncontent.ContentVersionLegacyJSONV1,
		ContentType: commoncontent.ContentType_CONTENT_TYPE_HTML,
		RawHtml:     &rawHTML,
	}
}

func (d *Demo) randomWindowedTime(past bool, now time.Time, maxDays int64) time.Time {
	if maxDays < 1 {
		maxDays = 1
	}

	offsetDays := int64(1 + d.randIntN(int(maxDays)))
	if past {
		return windowAnchor(now.AddDate(0, 0, -int(offsetDays)))
	}

	return windowAnchor(now.AddDate(0, 0, int(offsetDays)))
}

func (d *Demo) randomMeetingStart(past bool, now time.Time, maxDays int64) time.Time {
	day := d.randomWindowedTime(past, now, maxDays)
	loc := day.Location()
	hour := 8 + d.randIntN(10)
	minuteOptions := []int{0, 15, 30, 45}
	minute := minuteOptions[d.randIntN(len(minuteOptions))]
	return time.Date(
		day.Year(),
		day.Month(),
		day.Day(),
		hour,
		minute,
		0,
		0,
		loc,
	)
}

func (d *Demo) randomDayBoundary(past bool, now time.Time, maxDays int64) time.Time {
	return d.randomWindowedTime(past, now, maxDays)
}

func windowAnchor(t time.Time) time.Time {
	loc := t.Location()
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

func titleizeJob(job string) string {
	job = strings.TrimSpace(job)
	if job == "" {
		return "Target"
	}

	return strings.ToUpper(job[:1]) + job[1:]
}
