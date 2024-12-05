package calendar

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/storage"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tCalendar     = table.FivenetCalendar.AS("calendar")
	tCalendarSubs = table.FivenetCalendarSubs.AS("calendar_sub")

	tCalendarEntry = table.FivenetCalendarEntries.AS("calendar_entry")
	tCalendarRSVP  = table.FivenetCalendarRsvp.AS("calendar_entry_rsvp")

	tCJobAccess  = table.FivenetCalendarJobAccess.AS("calendar_job_access")
	tCUserAccess = table.FivenetCalendarUserAccess.AS("calendar_user_access")

	tUsers   = table.Users.AS("user_short")
	tCreator = tUsers.AS("creator")

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
	CalendarServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	st       storage.IStorage
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
	Storage   storage.IStorage
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
		st:       p.Storage,
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
				table.FivenetCalendarJobAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCalendarJobAccess.ID,
						CreatedAt: table.FivenetCalendarJobAccess.CreatedAt,
						TargetID:  table.FivenetCalendarJobAccess.CalendarID,
						Access:    table.FivenetCalendarJobAccess.Access,
					},
					Job:          table.FivenetCalendarJobAccess.Job,
					MinimumGrade: table.FivenetCalendarJobAccess.MinimumGrade,
				},
				table.FivenetCalendarJobAccess.AS("calendar_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCalendarJobAccess.AS("calendar_job_access").ID,
						CreatedAt: table.FivenetCalendarJobAccess.AS("calendar_job_access").CreatedAt,
						TargetID:  table.FivenetCalendarJobAccess.AS("calendar_job_access").CalendarID,
						Access:    table.FivenetCalendarJobAccess.AS("calendar_job_access").Access,
					},
					Job:          table.FivenetCalendarJobAccess.AS("calendar_job_access").Job,
					MinimumGrade: table.FivenetCalendarJobAccess.AS("calendar_job_access").MinimumGrade,
				},
			),
			access.NewUsers[calendar.CalendarUserAccess, *calendar.CalendarUserAccess, calendar.AccessLevel](
				table.FivenetCalendarUserAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCalendarUserAccess.ID,
						CreatedAt: table.FivenetCalendarUserAccess.CreatedAt,
						TargetID:  table.FivenetCalendarUserAccess.CalendarID,
						Access:    table.FivenetCalendarUserAccess.Access,
					},
					UserId: table.FivenetCalendarUserAccess.UserID,
				},
				table.FivenetCalendarUserAccess.AS("calendar_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetCalendarUserAccess.AS("calendar_user_access").ID,
						CreatedAt: table.FivenetCalendarUserAccess.AS("calendar_user_access").CreatedAt,
						TargetID:  table.FivenetCalendarUserAccess.AS("calendar_user_access").CalendarID,
						Access:    table.FivenetCalendarUserAccess.AS("calendar_user_access").Access,
					},
					UserId: table.FivenetCalendarUserAccess.AS("calendar_user_access").UserID,
				},
			),
			nil,
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCalendarServiceServer(srv, s)
}
