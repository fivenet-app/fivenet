package calendar

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	discordstate "github.com/diamondburned/arikawa/v3/state"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendarstore "github.com/fivenet-app/fivenet/v2026/stores/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarDiscordSettingsValueScanRoundTrip(t *testing.T) {
	t.Parallel()

	msg := "Reminder text"
	title := "Embed title"
	description := "Embed description"
	color := "#22C55E"

	original := &calendarresource.CalendarDiscordSettings{
		Enabled:   true,
		ChannelId: "1234567890",
		ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
			{
				AtMinute: 30,
				Message:  &msg,
				Embed: &calendarresource.CalendarDiscordReminderEmbed{
					Title:       &title,
					Description: &description,
					Color:       &color,
				},
			},
		},
	}

	value, err := original.Value()
	require.NoError(t, err)

	var scanned calendarresource.CalendarDiscordSettings
	require.NoError(t, scanned.Scan(value))

	require.Equal(t, original.GetEnabled(), scanned.GetEnabled())
	require.Equal(t, original.GetChannelId(), scanned.GetChannelId())
	require.Len(t, scanned.GetReminderSteps(), 1)
	assert.Equal(t, int32(30), scanned.GetReminderSteps()[0].GetAtMinute())
	assert.Equal(t, msg, scanned.GetReminderSteps()[0].GetMessage())
	require.NotNil(t, scanned.GetReminderSteps()[0].GetEmbed())
	assert.Equal(t, title, scanned.GetReminderSteps()[0].GetEmbed().GetTitle())
	assert.Equal(t, description, scanned.GetReminderSteps()[0].GetEmbed().GetDescription())
	assert.Equal(t, color, scanned.GetReminderSteps()[0].GetEmbed().GetColor())
}

func TestValidateCalendarDiscordSettingsRejectsInvalidConfigurations(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		job      string
		kind     calendarresource.CalendarSystemKind
		settings *calendarresource.CalendarDiscordSettings
		server   *Server
	}{
		{
			name: "duplicate step times",
			job:  "police",
			settings: &calendarresource.CalendarDiscordSettings{
				ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
					{AtMinute: 10, Message: new("a")},
					{AtMinute: 10, Message: new("b")},
				},
			},
			server: &Server{},
		},
		{
			name: "empty step payload",
			job:  "police",
			settings: &calendarresource.CalendarDiscordSettings{
				ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
					{AtMinute: 10},
				},
			},
			server: &Server{},
		},
		{
			name: "private calendar enabled",
			settings: &calendarresource.CalendarDiscordSettings{
				Enabled:   true,
				ChannelId: "123",
				ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
					{AtMinute: 5, Message: new("test")},
				},
			},
			server: &Server{},
		},
		{
			name: "system calendar enabled",
			job:  "police",
			kind: calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS,
			settings: &calendarresource.CalendarDiscordSettings{
				Enabled:   true,
				ChannelId: "123",
				ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
					{AtMinute: 5, Message: new("test")},
				},
			},
			server: &Server{},
		},
		{
			name: "bot disabled",
			job:  "police",
			settings: &calendarresource.CalendarDiscordSettings{
				Enabled:   true,
				ChannelId: "123",
				ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
					{AtMinute: 5, Message: new("test")},
				},
			},
			server: &Server{},
		},
		{
			name: "missing channel",
			job:  "police",
			settings: &calendarresource.CalendarDiscordSettings{
				Enabled: true,
				ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
					{AtMinute: 5, Message: new("test")},
				},
			},
			server: &Server{dc: discordstate.New("Bot test")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.server.validateCalendarDiscordSettings(
				t.Context(),
				tc.job,
				tc.kind,
				tc.settings,
			)
			require.Error(t, err)
		})
	}
}

func TestValidateCalendarDiscordSettingsRejectsInvalidChannelID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .*discord_guild_id.*fivenet_job_props").
		WillReturnRows(sqlmock.NewRows([]string{"discord_guild_id"}).AddRow("999"))

	srv := &Server{
		db:    db,
		dc:    discordstate.New("Bot test"),
		store: calendarstore.New(db),
	}

	err = srv.validateCalendarDiscordSettings(
		t.Context(),
		"police",
		calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED,
		&calendarresource.CalendarDiscordSettings{
			Enabled:   true,
			ChannelId: "not-a-channel",
			ReminderSteps: []*calendarresource.CalendarDiscordReminderStep{
				{AtMinute: 5, Message: new("test")},
			},
		},
	)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

var _ driver.Valuer = (*calendarresource.CalendarDiscordSettings)(nil)
