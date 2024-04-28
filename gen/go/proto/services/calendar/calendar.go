package calendar

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/config/appconfig"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/storage"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tCalendar      = table.FivenetCalendar.AS("calendar")
	tCalendarEntry = table.FivenetCalendarEntries.AS("calendar_entry")

	tCJobAccess  = table.FivenetCalendarJobAccess.AS("job_access")
	tCUserAccess = table.FivenetCalendarUserAccess.AS("user_access")

	tUsers   = table.Users
	tCreator = tUsers.AS("creator")
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
