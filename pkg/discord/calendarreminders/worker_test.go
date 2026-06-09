package calendarreminders

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildReminderSendMessageIncludesCustomAndSummaryEmbeds(t *testing.T) {
	t.Parallel()

	message := "Reminder message"
	title := "Custom title"
	description := "Custom description"
	color := "#3B82F6"

	entry := &calendarentries.CalendarEntry{
		Id:         42,
		Title:      "Shift meeting",
		CalendarId: 7,
		StartTime:  timestamp.New(time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)),
		Calendar: &calendarresource.Calendar{
			Id:   7,
			Name: "Command",
		},
		Occurrence: &calendarentries.CalendarEntryOccurrence{
			Key:  "recurring:42:1736935200",
			Kind: calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_RECURRING,
		},
	}

	data := buildReminderSendMessage(
		i18n.DummyTranslator(),
		&calendarresource.CalendarDiscordReminderStep{
			AtMinute: 15,
			Message:  &message,
			Embed: &calendarresource.CalendarDiscordReminderEmbed{
				Title:       &title,
				Description: &description,
				Color:       &color,
			},
		},
		entry,
		"https://fivenet.example",
	)

	assert.Equal(t, message, data.Content)
	require.Len(t, data.Embeds, 2)
	assert.Equal(t, title, data.Embeds[0].Title)
	assert.Equal(t, "Shift meeting", data.Embeds[1].Title)
	assert.Contains(
		t,
		data.Embeds[1].Description,
		"discord.calendar.calendar_name(map[name:Command])",
	)
	assert.Contains(t, data.Embeds[1].Description, "entryKey=recurring%3A42%3A1736935200")
}

func TestExpandReminderOccurrencesFastForwardsRecurringEntries(t *testing.T) {
	t.Parallel()

	entry := &calendarentries.CalendarEntry{
		Id:        9,
		Title:     "Daily briefing",
		StartTime: timestamp.New(time.Date(2026, time.January, 1, 8, 0, 0, 0, time.UTC)),
		Recurring: &calendarentries.CalendarEntryRecurring{
			Every: calendarentries.CalendarEntryRecurringEvery_CALENDAR_ENTRY_RECURRING_EVERY_DAY,
			Count: 1,
		},
	}

	rangeStart := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)
	rangeEnd := time.Date(2026, time.January, 16, 23, 59, 0, 0, time.UTC)

	occurrences := expandReminderOccurrences(entry, rangeStart, rangeEnd)
	require.Len(t, occurrences, 2)
	assert.Equal(
		t,
		time.Date(2026, time.January, 15, 8, 0, 0, 0, time.UTC),
		occurrences[0].GetStartTime().AsTime(),
	)
	assert.Equal(
		t,
		time.Date(2026, time.January, 16, 8, 0, 0, 0, time.UTC),
		occurrences[1].GetStartTime().AsTime(),
	)
}

func TestProcessOccurrenceSkipsAlreadySentReminder(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .*fivenet_calendar_discord_reminder_sends").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	sendCalls := 0
	worker := &Worker{
		db: db,
		sendFn: func(channelID discord.ChannelID, data api.SendMessageData) error {
			sendCalls++
			return nil
		},
	}

	now := time.Date(2026, time.January, 15, 9, 55, 0, 0, time.UTC)
	cal := &calendarresource.Calendar{
		Id: 7,
		DiscordSettings: &calendarresource.CalendarDiscordSettings{
			Enabled:   true,
			ChannelId: "123",
			ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
				{AtMinute: 10, Message: strPtr("test")},
			},
		},
	}
	entry := &calendarentries.CalendarEntry{
		Id:        99,
		Title:     "Meeting",
		StartTime: timestamp.New(time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)),
		Calendar:  cal,
		Occurrence: &calendarentries.CalendarEntryOccurrence{
			Key:  "manual:99:1736935200",
			Kind: calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_MANUAL,
		},
	}

	require.NoError(
		t,
		worker.processOccurrence(t.Context(), cal, entry, discord.ChannelID(123), now),
	)
	assert.Equal(t, 0, sendCalls)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProcessOccurrenceSendsOnlyDueStepsAndRecordsLog(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .*fivenet_calendar_discord_reminder_sends").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec("INSERT INTO .*fivenet_calendar_discord_reminder_sends").
		WillReturnResult(sqlmock.NewResult(1, 1))

	sendCalls := 0
	worker := &Worker{
		db: db,

		i18n: i18n.NewDummy(),

		publicURL: "https://fivenet.example",
		sendFn: func(channelID discord.ChannelID, data api.SendMessageData) error {
			sendCalls++
			require.Len(t, data.Embeds, 1)
			assert.Contains(t, data.Embeds[0].Description, "entryId=100")
			return nil
		},
	}

	now := time.Date(2026, time.January, 15, 9, 55, 0, 0, time.UTC)
	cal := &calendarresource.Calendar{
		Id:   7,
		Name: "Operations",
		DiscordSettings: &calendarresource.CalendarDiscordSettings{
			Enabled:   true,
			ChannelId: "123",
			ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
				{AtMinute: 10, Message: strPtr("due")},
				{AtMinute: 1, Message: strPtr("future")},
			},
		},
	}
	entry := &calendarentries.CalendarEntry{
		Id:        100,
		Title:     "Briefing",
		StartTime: timestamp.New(time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC)),
		Calendar:  cal,
		Occurrence: &calendarentries.CalendarEntryOccurrence{
			Key:  "manual:100:1736935200",
			Kind: calendarentries.CalendarEntryOccurrenceKind_CALENDAR_ENTRY_OCCURRENCE_KIND_MANUAL,
		},
	}

	require.NoError(
		t,
		worker.processOccurrence(t.Context(), cal, entry, discord.ChannelID(123), now),
	)
	assert.Equal(t, 1, sendCalls)
	require.NoError(t, mock.ExpectationsWereMet())
}

func strPtr(v string) *string {
	return &v
}
