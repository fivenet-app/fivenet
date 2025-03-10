package calendar

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	pbcalendar "github.com/fivenet-app/fivenet/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
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
		TimestampColumn: table.FivenetCalendar.DeletedAt,
		MinDays:         60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCalendarEntries,
		TimestampColumn: table.FivenetCalendarEntries.DeletedAt,
		MinDays:         60,
	})
}

type Server struct {
	pbcalendar.CalendarServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
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
	Aud       audit.IAuditer
	AppConfig appconfig.IConfig
	Notif     notifi.INotifi
	JS        *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		p:        p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
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
					Job:          table.FivenetCalendarAccess.AS("calendar_job_access").Job,
					MinimumGrade: table.FivenetCalendarAccess.AS("calendar_job_access").MinimumGrade,
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
					UserId: table.FivenetCalendarAccess.UserID,
				},
				table.FivenetCalendarAccess.AS("calendar_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.AS("calendar_user_access").ID,
						TargetID: table.FivenetCalendarAccess.AS("calendar_user_access").TargetID,
						Access:   table.FivenetCalendarAccess.AS("calendar_user_access").Access,
					},
					UserId: table.FivenetCalendarAccess.AS("calendar_user_access").UserID,
				},
			),
			nil,
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcalendar.RegisterCalendarServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbcalendar.PermsRemap
}
