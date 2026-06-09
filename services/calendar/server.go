package calendar

import (
	"database/sql"

	discordstate "github.com/diamondburned/arikawa/v3/state"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tCalendar     = table.FivenetCalendar.AS("calendar")
	tCalendarSubs = table.FivenetCalendarSubs.AS("calendar_sub")

	tCalendarEntry          = table.FivenetCalendarEntries.AS("calendar_entry")
	tCalendarRSVP           = table.FivenetCalendarRsvp.AS("calendar_entry_rsvp")
	tCalendarRSVPOccurrence = table.FivenetCalendarRsvpOccurrence.AS(
		"calendar_entry_rsvp_occurrence",
	)

	tCAccess = table.FivenetCalendarAccess.AS("calendar_access")

	tUserJobs  = table.FivenetUserJobs.AS("user_jobs")
	tUserProps = table.FivenetUserProps
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCalendar,
		JobColumn:       table.FivenetCalendar.Job,
		DeletedAtColumn: table.FivenetCalendar.DeletedAt,
		IDColumn:        table.FivenetCalendar.ID,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:      table.FivenetCalendarSubs,
				ForeignKey: table.FivenetCalendarSubs.CalendarID,
			},
			{
				Table:           table.FivenetCalendarEntries,
				IDColumn:        table.FivenetCalendarEntries.ID,
				JobColumn:       table.FivenetCalendarEntries.Job,
				ForeignKey:      table.FivenetCalendarEntries.CalendarID,
				DeletedAtColumn: table.FivenetCalendarEntries.DeletedAt,

				MinDays: 60,

				DependantTables: []*housekeeper.Table{
					{
						Table:      table.FivenetCalendarRsvp,
						ForeignKey: table.FivenetCalendarRsvp.EntryID,
					},
					{
						Table:      table.FivenetCalendarRsvpOccurrence,
						ForeignKey: table.FivenetCalendarRsvpOccurrence.EntryID,
					},
					{
						Table:      table.FivenetCalendarDiscordReminderSends,
						ForeignKey: table.FivenetCalendarDiscordReminderSends.EntryID,
					},
				},
			},
		},
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCalendarDiscordReminderSends,
		IDColumn:        table.FivenetCalendarDiscordReminderSends.ID,
		TimestampColumn: table.FivenetCalendarDiscordReminderSends.CreatedAt,
		MinDays:         30,
	})
}

type Server struct {
	pbcalendar.CalendarServiceServer
	pbcalendar.EntriesServiceServer

	db       *sql.DB
	cfg      *config.Config
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	appCfg   appconfig.IConfig
	i18n     i18n.Ii18n
	notif    notifi.INotifi
	js       *events.JSWrapper
	dc       *discordstate.State

	access *access.Grouped[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, access.DummyQualificationAccess[calendaraccess.AccessLevel], *access.DummyQualificationAccess[calendaraccess.AccessLevel], calendaraccess.AccessLevel]
}

type Params struct {
	fx.In

	DB        *sql.DB
	Config    *config.Config
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	AppConfig appconfig.IConfig
	I18n      i18n.Ii18n
	Notif     notifi.INotifi
	JS        *events.JSWrapper
	Discord   *discordstate.State
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		cfg:      p.Config,
		ps:       p.P,
		enricher: p.Enricher,
		appCfg:   p.AppConfig,
		i18n:     p.I18n,
		notif:    p.Notif,
		js:       p.JS,
		dc:       p.Discord,
		access: access.NewGrouped[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, access.DummyQualificationAccess[calendaraccess.AccessLevel], *access.DummyQualificationAccess[calendaraccess.AccessLevel], calendaraccess.AccessLevel](
			p.DB,
			table.FivenetDocuments,
			&access.TargetTableColumns{
				ID:         table.FivenetDocuments.ID,
				DeletedAt:  table.FivenetDocuments.DeletedAt,
				CreatorJob: table.FivenetDocuments.CreatorJob,
				CreatorID:  table.FivenetDocuments.CreatorID,
			},
			access.NewJobs[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.AccessLevel](
				table.FivenetCalendarAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.ID,
						TargetID: table.FivenetCalendarAccess.TargetID,
						Access:   table.FivenetCalendarAccess.Access,
					},
					Job:          table.FivenetCalendarAccess.Job,
					MinimumGrade: table.FivenetCalendarAccess.MinimumGrade,
				},
				table.FivenetCalendarAccess.AS("calendar_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.AS("calendar_job_access").ID,
						TargetID: table.FivenetCalendarAccess.AS("calendar_job_access").TargetID,
						Access:   table.FivenetCalendarAccess.AS("calendar_job_access").Access,
					},
					Job: table.FivenetCalendarAccess.AS("calendar_job_access").Job,
					MinimumGrade: table.FivenetCalendarAccess.AS(
						"calendar_job_access",
					).MinimumGrade,
				},
			),
			access.NewUsers[calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, calendaraccess.AccessLevel](
				table.FivenetCalendarAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.ID,
						TargetID: table.FivenetCalendarAccess.TargetID,
						Access:   table.FivenetCalendarAccess.Access,
					},
					UserID: table.FivenetCalendarAccess.UserID,
				},
				table.FivenetCalendarAccess.AS("calendar_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.AS("calendar_user_access").ID,
						TargetID: table.FivenetCalendarAccess.AS("calendar_user_access").TargetID,
						Access:   table.FivenetCalendarAccess.AS("calendar_user_access").Access,
					},
					UserID: table.FivenetCalendarAccess.AS("calendar_user_access").UserID,
				},
			),
			nil,
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcalendar.RegisterCalendarServiceServer(srv, s)
	pbcalendar.RegisterEntriesServiceServer(srv, s)
}
