package calendar

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/storage"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tCalendar      = table.FivenetCalendar.AS("calendar")
	tCalendarShort = tCalendar.AS("calendar_short")
	tCalendarEntry = table.FivenetCalendarEntries.AS("calendar_entry")
	tCalendarRSVP  = table.FivenetCalendarRsvp.AS("calendar_entry_rsvp")
	tCalendarSubs  = table.FivenetCalendarSubs.AS("calendar_subs")

	tCJobAccess  = table.FivenetCalendarJobAccess.AS("job_access")
	tCUserAccess = table.FivenetCalendarUserAccess.AS("user_access")

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
}

type Params struct {
	fx.In

	DB        *sql.DB
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	Aud       audit.IAuditer
	Storage   storage.IStorage
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		p:        p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		st:       p.Storage,
		appCfg:   p.AppConfig,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCalendarServiceServer(srv, s)
}
