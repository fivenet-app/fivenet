package calendar

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
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
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCalendarServiceServer(srv, s)
}
