package calendar

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/calendar"
	pbcalendar "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
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

	access *access.Grouped[calendar.CalendarJobAccess, *calendar.CalendarJobAccess, calendar.CalendarUserAccess, *calendar.CalendarUserAccess, access.DummyQualificationAccess[calendar.AccessLevel], *access.DummyQualificationAccess[calendar.AccessLevel], calendar.AccessLevel]
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
		access: access.NewGrouped[calendar.CalendarJobAccess, *calendar.CalendarJobAccess, calendar.CalendarUserAccess, *calendar.CalendarUserAccess, access.DummyQualificationAccess[calendar.AccessLevel], *access.DummyQualificationAccess[calendar.AccessLevel], calendar.AccessLevel](
			p.DB,
			table.FivenetDocuments,
			&access.TargetTableColumns{
				ID:         table.FivenetDocuments.ID,
				DeletedAt:  table.FivenetDocuments.DeletedAt,
				CreatorJob: table.FivenetDocuments.CreatorJob,
				CreatorID:  table.FivenetDocuments.CreatorID,
			},
			access.NewJobs[calendar.CalendarJobAccess, *calendar.CalendarJobAccess, calendar.AccessLevel](
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
			access.NewUsers[calendar.CalendarUserAccess, *calendar.CalendarUserAccess, calendar.AccessLevel](
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

// GetPermsRemap returns the permissions re-mapping for the services.
func (s *Server) GetPermsRemap() map[string]string {
	return pbcalendar.PermsRemap
}
