package calendarreminders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	discordstate "github.com/diamondburned/arikawa/v3/state"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	discordembeds "github.com/fivenet-app/fivenet/v2026/pkg/discord/embeds"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type sendReminderMessageFn func(channelID discord.ChannelID, data api.SendMessageData) error

type Worker struct {
	logger *zap.Logger
	db     *sql.DB
	dc     *discordstate.State
	i18n   i18n.Ii18n

	enabled   bool
	publicURL string
	nowFn     func() time.Time
	sendFn    sendReminderMessageFn
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	DB      *sql.DB
	Config  *config.Config
	Discord *discordstate.State
	I18n    i18n.Ii18n
}

func NewWorker(p Params) *Worker {
	w := &Worker{
		logger: p.Logger.Named("discord.calendar_reminders"),
		db:     p.DB,
		dc:     p.Discord,
		i18n:   p.I18n,

		enabled:   p.Config.Discord.Enabled && p.Config.Discord.CalendarReminders,
		publicURL: p.Config.HTTP.PublicURL,
		nowFn:     time.Now,
	}
	w.sendFn = func(channelID discord.ChannelID, data api.SendMessageData) error {
		if w.dc == nil {
			return errors.New("discord bot is not enabled")
		}

		_, err := w.dc.SendMessageComplex(channelID, data)
		return err
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		// Discord bot not enabled
		if !w.enabled {
			return nil
		}

		wg.Go(func() {
			for {
				select {
				case <-ctxCancel.Done():
					return

				case <-time.After(1 * time.Minute):
				}

				if err := w.run(ctxCancel); err != nil {
					w.logger.Error("error during discord calendar reminders run", zap.Error(err))
				}
			}
		})

		return nil
	}))
	p.LC.Append(fx.StopHook(func(ctxStartup context.Context) error {
		cancel()

		wg.Wait()
		return nil
	}))

	return w
}

func (w *Worker) run(ctx context.Context) error {
	calendars, err := w.loadReminderCalendars(ctx)
	if err != nil {
		return err
	}
	if len(calendars) == 0 {
		return nil
	}

	maxStepMinutes := reminderMaxStepMinutes(calendars)
	if maxStepMinutes <= 0 {
		return nil
	}

	now := w.nowFn().UTC()
	rangeEnd := now.Add(time.Duration(maxStepMinutes) * time.Minute)

	for i := range calendars {
		if err := w.processCalendar(
			ctx,
			calendars[i].Calendar,
			calendars[i].JobProps,
			now,
			rangeEnd,
		); err != nil {
			w.logger.Error(
				"failed to process discord reminder calendar",
				zap.Int64("calendar_id", calendars[i].Calendar.GetId()),
				zap.Error(err),
			)
		}
	}

	return nil
}

type reminderCalendar struct {
	Calendar *pbcalendar.Calendar `alias:"calendar"`
	JobProps *jobsprops.JobProps  `alias:"job_props"`
}

func (w *Worker) loadReminderCalendars(
	ctx context.Context,
) ([]*reminderCalendar, error) {
	tCalendar := table.FivenetCalendar.AS("calendar")
	tJobProps := table.FivenetJobProps.AS("job_props")
	tFiles := table.FivenetFiles.AS("logo_file")

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
			tCalendar.Job,
			tCalendar.SystemKind,
			tCalendar.Name,
			tCalendar.DiscordSettings,
			tJobProps.LogoFileID,
			tJobProps.LogoFileID,
			tFiles.ID,
			tFiles.FilePath,
		).
		FROM(tCalendar.
			LEFT_JOIN(tJobProps,
				tJobProps.Job.EQ(tCalendar.Job),
			).
			LEFT_JOIN(tFiles,
				tFiles.ID.EQ(tJobProps.LogoFileID),
			),
		).
		WHERE(mysql.AND(
			tCalendar.DeletedAt.IS_NULL(),
			tCalendar.Job.IS_NOT_NULL(),
			mysql.OR(
				tCalendar.SystemKind.IS_NULL(),
				tCalendar.SystemKind.EQ(
					mysql.Int32(
						int32(pbcalendar.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED),
					),
				),
			),
			tCalendar.DiscordSettings.IS_NOT_NULL(),
		)).
		ORDER_BY(tCalendar.ID.ASC())

	calendars := []*reminderCalendar{}
	if err := stmt.QueryContext(ctx, w.db, &calendars); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	filtered := make([]*reminderCalendar, 0, len(calendars))
	for i := range calendars {
		if calendars[i] == nil || calendars[i].Calendar.GetDiscordSettings() == nil {
			continue
		}
		if !calendars[i].Calendar.GetDiscordSettings().GetEnabled() {
			continue
		}
		if strings.TrimSpace(calendars[i].Calendar.GetDiscordSettings().GetChannelId()) == "" {
			continue
		}
		if len(calendars[i].Calendar.GetDiscordSettings().GetReminderSteps()) == 0 {
			continue
		}
		filtered = append(filtered, calendars[i])
	}

	return filtered, nil
}

func reminderMaxStepMinutes(calendars []*reminderCalendar) int32 {
	var maxMinutes int32
	for i := range calendars {
		settings := calendars[i].Calendar.GetDiscordSettings()
		for j := range settings.GetReminderSteps() {
			maxMinutes = max(maxMinutes, settings.GetReminderSteps()[j].GetAtMinute())
		}
	}
	return maxMinutes
}

func (w *Worker) processCalendar(
	ctx context.Context,
	cal *pbcalendar.Calendar,
	jps *jobsprops.JobProps,
	now time.Time,
	rangeEnd time.Time,
) error {
	channelIDNum, err := strconv.ParseUint(
		strings.TrimSpace(cal.GetDiscordSettings().GetChannelId()),
		10,
		64,
	)
	if err != nil {
		return fmt.Errorf("invalid discord channel id: %w", err)
	}

	entries, err := w.loadReminderEntries(ctx, cal.GetId(), rangeEnd)
	if err != nil {
		return err
	}

	for i := range entries {
		entries[i].Calendar = cal

		occurrences := expandReminderOccurrences(entries[i], now, rangeEnd)
		for j := range occurrences {
			if err := w.processOccurrence(
				ctx,
				cal,
				jps,
				occurrences[j],
				discord.ChannelID(channelIDNum),
				now,
			); err != nil {
				w.logger.Error(
					"failed to process discord reminder occurrence",
					zap.Int64("calendar_id", cal.GetId()),
					zap.Int64("entry_id", occurrences[j].GetId()),
					zap.String("occurrence_key", occurrences[j].GetOccurrence().GetKey()),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

func (w *Worker) loadReminderEntries(
	ctx context.Context,
	calendarID int64,
	rangeEnd time.Time,
) ([]*calendarentries.CalendarEntry, error) {
	tCalendarEntries := table.FivenetCalendarEntries.AS("calendar_entry")

	stmt := tCalendarEntries.
		SELECT(
			tCalendarEntries.ID,
			tCalendarEntries.CalendarID,
			tCalendarEntries.Job,
			tCalendarEntries.StartTime,
			tCalendarEntries.EndTime,
			tCalendarEntries.Title,
			tCalendarEntries.Content,
			tCalendarEntries.Closed,
			tCalendarEntries.RsvpOpen,
			tCalendarEntries.CreatorID,
			tCalendarEntries.CreatorJob,
			tCalendarEntries.Recurring,
		).
		FROM(tCalendarEntries).
		WHERE(mysql.AND(
			tCalendarEntries.DeletedAt.IS_NULL(),
			tCalendarEntries.CalendarID.EQ(mysql.Int64(calendarID)),
			tCalendarEntries.Closed.IS_FALSE(),
			tCalendarEntries.StartTime.LT_EQ(mysql.TimestampT(rangeEnd)),
		)).
		ORDER_BY(
			tCalendarEntries.StartTime.ASC(),
			tCalendarEntries.ID.ASC(),
		)

	entries := []*calendarentries.CalendarEntry{}
	if err := stmt.QueryContext(ctx, w.db, &entries); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return entries, nil
}

func (w *Worker) processOccurrence(
	ctx context.Context,
	cal *pbcalendar.Calendar,
	jps *jobsprops.JobProps,
	occurrence *calendarentries.CalendarEntry,
	channelID discord.ChannelID,
	now time.Time,
) error {
	if occurrence == nil || occurrence.GetOccurrence() == nil {
		return nil
	}

	for i := range cal.GetDiscordSettings().GetReminderSteps() {
		step := cal.GetDiscordSettings().GetReminderSteps()[i]
		if step == nil || !reminderStepDue(now, occurrence.GetStartTime().AsTime(), step) {
			continue
		}

		if sent, err := w.hasReminderSendLog(
			ctx,
			cal.GetId(),
			occurrence.GetId(),
			occurrence.GetOccurrence().GetKey(),
			step.GetAtMinute(),
		); err != nil {
			return err
		} else if sent {
			continue
		}

		data := buildReminderSendMessage(
			reminderTranslator(w.i18n),
			jps,
			step,
			occurrence,
			w.publicURL,
		)
		if err := w.sendFn(channelID, data); err != nil {
			return err
		}

		if err := w.recordReminderSendLog(
			ctx,
			cal.GetId(),
			occurrence.GetId(),
			occurrence.GetOccurrence().GetKey(),
			step.GetAtMinute(),
		); err != nil {
			return err
		}
	}

	return nil
}

func reminderTranslator(localizer i18n.Ii18n) i18n.TFunc {
	if localizer == nil {
		return i18n.DummyTranslator()
	}

	locale := strings.TrimSpace(localizer.GetFallbackLanguage())
	if locale == "" {
		locale = "en"
	}

	return localizer.Translator(locale)
}

func reminderStepDue(
	now time.Time,
	start time.Time,
	step *pbcalendar.CalendarDiscordReminderStep,
) bool {
	scheduled := start.Add(-time.Duration(step.GetAtMinute()) * time.Minute)
	return !scheduled.After(now)
}

func buildReminderSendMessage(
	i18n i18n.TFunc,
	jps *jobsprops.JobProps,
	step *pbcalendar.CalendarDiscordReminderStep,
	entry *calendarentries.CalendarEntry,
	publicURL string,
) api.SendMessageData {
	embeds := make([]discord.Embed, 0, 2)

	embed := step.GetEmbed()

	if custom := buildCustomReminderEmbed(i18n, embed, jps, publicURL); custom != nil {
		embeds = append(embeds, *custom)
	}
	if summary := buildReminderSummaryEmbed(i18n, entry, embed, publicURL); summary != nil {
		embeds = append(embeds, *summary)
	}

	return api.SendMessageData{
		Content: step.GetMessage(),
		Embeds:  embeds,
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{
				api.AllowRoleMention,
			},
		},
	}
}

func buildCustomReminderEmbed(
	i18n i18n.TFunc,
	embed *pbcalendar.CalendarDiscordReminderEmbed,
	jps *jobsprops.JobProps,
	publicURL string,
) *discord.Embed {
	if !calendarDiscordEmbedHasContent(embed) {
		return nil
	}
	footer := &discord.EmbedFooter{
		Text: i18n("discord.calendar.custom_embed_author", nil),
		Icon: discordembeds.EmbedFooterFiveNet.Icon,
	}
	if jps != nil {
		iconURL, err := url.JoinPath(publicURL, jps.GetLogoFile().GetFilePath())
		if err == nil {
			footer.Icon = iconURL
		}
	}

	out := &discord.Embed{
		Type:        discord.NormalEmbed,
		Title:       embed.GetTitle(),
		Description: embed.GetDescription(),
		Footer:      footer,
	}
	if color := parseDiscordHexColor(embed.GetColor()); color != 0 {
		out.Color = color
	}

	return out
}

func buildReminderSummaryEmbed(
	i18n i18n.TFunc,
	entry *calendarentries.CalendarEntry,
	embed *pbcalendar.CalendarDiscordReminderEmbed,
	publicURL string,
) *discord.Embed {
	if entry == nil || entry.GetStartTime() == nil || entry.GetCalendar() == nil {
		return nil
	}

	link := buildCalendarEntryLink(entry, publicURL)
	lines := []string{
		i18n(
			"discord.calendar.calendar_name",
			map[string]any{"name": entry.GetCalendar().GetName()},
		),
		i18n(
			"discord.calendar.start_time",
			map[string]any{"ts": entry.GetStartTime().AsTime().Unix()},
		),
	}
	if entry.GetEndTime() != nil {
		lines = append(
			lines,
			i18n(
				"discord.calendar.end_time",
				map[string]any{"ts": entry.GetEndTime().AsTime().Unix()},
			),
		)
	}
	if link != "" {
		lines = append(
			lines,
			i18n(
				"discord.calendar.link",
				map[string]any{"link": link},
			),
		)
	}

	color := discordembeds.ColorInfo
	if c := parseDiscordHexColor(embed.GetColor()); c != 0 {
		color = c
	}

	return &discord.Embed{
		Type:        discord.NormalEmbed,
		Title:       entry.GetTitle(),
		Description: strings.Join(lines, "\n"),
		URL:         link,
		Color:       color,
		Footer: &discord.EmbedFooter{
			Text: i18n("discord.calendar.custom_embed_author", nil),
			Icon: discordembeds.EmbedAuthor.Icon,
		},
	}
}

func buildCalendarEntryLink(entry *calendarentries.CalendarEntry, publicURL string) string {
	base := strings.TrimSpace(publicURL)
	if base == "" || entry == nil {
		return ""
	}

	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/calendar"

	query := u.Query()
	if entry.GetOccurrence() != nil &&
		entry.GetOccurrence().
			GetKind() ==
			calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING &&
		entry.GetOccurrence().GetKey() != "" {
		query.Set("entryKey", entry.GetOccurrence().GetKey())
	} else {
		query.Set("entryId", strconv.FormatInt(entry.GetId(), 10))
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func parseDiscordHexColor(color string) discord.Color {
	color = strings.TrimSpace(strings.TrimPrefix(color, "#"))
	if len(color) != 6 {
		return 0
	}

	value, err := strconv.ParseUint(color, 16, 32)
	if err != nil {
		return 0
	}

	// Check if parsed value is in int32 bounds
	if value > math.MaxInt32 {
		return 0
	}

	return discord.Color(value)
}

func (w *Worker) hasReminderSendLog(
	ctx context.Context,
	calendarID int64,
	entryID int64,
	occurrenceKey string,
	stepAtMinute int32,
) (bool, error) {
	stmt := table.FivenetCalendarDiscordReminderSends.
		SELECT(table.FivenetCalendarDiscordReminderSends.ID.AS("id")).
		FROM(table.FivenetCalendarDiscordReminderSends).
		WHERE(mysql.AND(
			table.FivenetCalendarDiscordReminderSends.CalendarID.EQ(mysql.Int64(calendarID)),
			table.FivenetCalendarDiscordReminderSends.EntryID.EQ(mysql.Int64(entryID)),
			table.FivenetCalendarDiscordReminderSends.OccurrenceKey.EQ(mysql.String(occurrenceKey)),
			table.FivenetCalendarDiscordReminderSends.StepAtMinute.EQ(mysql.Int32(stepAtMinute)),
		)).
		LIMIT(1)

	dest := struct {
		ID int64 `alias:"id"`
	}{}
	if err := stmt.QueryContext(ctx, w.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return dest.ID > 0, nil
}

func (w *Worker) recordReminderSendLog(
	ctx context.Context,
	calendarID int64,
	entryID int64,
	occurrenceKey string,
	stepAtMinute int32,
) error {
	_, err := table.FivenetCalendarDiscordReminderSends.
		INSERT(
			table.FivenetCalendarDiscordReminderSends.CalendarID,
			table.FivenetCalendarDiscordReminderSends.EntryID,
			table.FivenetCalendarDiscordReminderSends.OccurrenceKey,
			table.FivenetCalendarDiscordReminderSends.StepAtMinute,
		).
		VALUES(
			calendarID,
			entryID,
			occurrenceKey,
			stepAtMinute,
		).
		ExecContext(ctx, w.db)

	return err
}

func expandReminderOccurrences(
	entry *calendarentries.CalendarEntry,
	rangeStart time.Time,
	rangeEnd time.Time,
) []*calendarentries.CalendarEntry {
	if entry == nil || entry.GetStartTime() == nil {
		return nil
	}

	if entry.GetRecurring() == nil {
		if !entryOverlapsRange(
			entry.GetStartTime().AsTime(),
			entry.GetEndTime(),
			rangeStart,
			rangeEnd,
		) {
			return nil
		}

		clone := proto.Clone(entry).(*calendarentries.CalendarEntry)
		clone.Occurrence = &calendarentries.CalendarEntryOccurrence{
			Key: fmt.Sprintf(
				"manual:%d:%d",
				clone.GetId(),
				clone.GetStartTime().AsTime().Unix(),
			),
			Kind:          calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_MANUAL,
			SourceEntryId: &clone.Id,
			AllDay:        clone.GetEndTime() == nil,
		}
		return []*calendarentries.CalendarEntry{clone}
	}

	interval := entry.GetRecurring().GetCount()
	if interval <= 0 {
		interval = 1
	}

	duration := time.Duration(0)
	if entry.GetEndTime() != nil {
		duration = entry.GetEndTime().AsTime().Sub(entry.GetStartTime().AsTime())
	}

	occurrenceStart := fastForwardRecurringStart(
		entry.GetStartTime().AsTime(),
		rangeStart,
		interval,
		entry.GetRecurring().GetEvery(),
	)

	out := []*calendarentries.CalendarEntry{}
	for !occurrenceStart.After(rangeEnd) {
		if until := entry.GetRecurring().
			GetUntil(); until != nil &&
			occurrenceStart.After(until.AsTime()) {
			break
		}

		if entryOverlapsRange(occurrenceStart, entry.GetEndTime(), rangeStart, rangeEnd) {
			clone := proto.Clone(entry).(*calendarentries.CalendarEntry)
			clone.StartTime = timestamp.New(occurrenceStart)
			if entry.GetEndTime() != nil {
				clone.EndTime = timestamp.New(occurrenceStart.Add(duration))
			}
			clone.Occurrence = &calendarentries.CalendarEntryOccurrence{
				Key: fmt.Sprintf(
					"recurring:%d:%d",
					entry.GetId(),
					occurrenceStart.Unix(),
				),
				Kind:          calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING,
				SourceEntryId: &clone.Id,
				AllDay:        clone.GetEndTime() == nil,
			}
			out = append(out, clone)
		}

		occurrenceStart = nextRecurringOccurrence(
			occurrenceStart,
			interval,
			entry.GetRecurring().GetEvery(),
		)
	}

	return out
}

func fastForwardRecurringStart(
	start time.Time,
	rangeStart time.Time,
	interval int32,
	every calendarentries.CalendarEntryRecurringEvery,
) time.Time {
	if !start.Before(rangeStart) {
		return start
	}

	switch every {
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_DAY:
		diffDays := int(math.Max(0, rangeStart.Sub(start).Hours()/24))
		steps := diffDays / int(interval)
		start = start.AddDate(0, 0, steps*int(interval))
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_WEEK:
		diffDays := int(math.Max(0, rangeStart.Sub(start).Hours()/24))
		steps := diffDays / (7 * int(interval))
		start = start.AddDate(0, 0, steps*7*int(interval))
	}

	for start.Before(rangeStart) {
		start = nextRecurringOccurrence(start, interval, every)
	}

	return start
}

func nextRecurringOccurrence(
	start time.Time,
	interval int32,
	every calendarentries.CalendarEntryRecurringEvery,
) time.Time {
	switch every {
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_DAY:
		return start.AddDate(0, 0, int(interval))
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_WEEK:
		return start.AddDate(0, 0, 7*int(interval))
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_MONTH:
		return start.AddDate(0, int(interval), 0)
	case calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_YEAR:
		return start.AddDate(int(interval), 0, 0)
	default:
		return start.AddDate(0, 0, int(interval))
	}
}

func entryOverlapsRange(
	start time.Time,
	end *timestamp.Timestamp,
	rangeStart time.Time,
	rangeEnd time.Time,
) bool {
	if start.After(rangeEnd) {
		return false
	}

	if end == nil {
		return !start.Before(rangeStart) && !start.After(rangeEnd)
	}

	endTime := end.AsTime()
	return !endTime.Before(rangeStart) && !start.After(rangeEnd)
}

func calendarDiscordEmbedHasContent(embed *pbcalendar.CalendarDiscordReminderEmbed) bool {
	if embed == nil {
		return false
	}

	return embed.Title != nil || embed.Description != nil
}
