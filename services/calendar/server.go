package calendar

import (
	"database/sql"

	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
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

	tCalendarEntry = table.FivenetCalendarEntries.AS("calendar_entry")
	tCalendarRSVP  = table.FivenetCalendarRsvp.AS("calendar_entry_rsvp")

	tCAccess = table.FivenetCalendarAccess.AS("calendar_access")

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
				},
			},
		},
	})
}

type Server struct {
	pbcalendar.CalendarServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	appCfg   appconfig.IConfig
	notif    notifi.INotifi
	js       *events.JSWrapper

	access *access.Grouped[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, access.DummyQualificationAccess[calendaraccess.AccessLevel], *access.DummyQualificationAccess[calendaraccess.AccessLevel], calendaraccess.AccessLevel]
}

type Params struct {
	fx.In

	DB        *sql.DB
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	AppConfig appconfig.IConfig
	Notif     notifi.INotifi
	JS        *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		ps:       p.P,
		enricher: p.Enricher,
		appCfg:   p.AppConfig,
		notif:    p.Notif,
		js:       p.JS,
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
}
