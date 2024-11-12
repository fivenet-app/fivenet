package mailer

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tUsers   = table.Users.AS("usershort")
	tCreator = table.Users.AS("creator")

	tUserProps = table.FivenetUserProps

	tThreads           = table.FivenetMsgsThreads.AS("thread")
	tThreadsUserState  = table.FivenetMsgsThreadsUserState.AS("threaduserstate")
	tThreadsJobAccess  = table.FivenetMsgsThreadsJobAccess
	tThreadsUserAccess = table.FivenetMsgsThreadsUserAccess

	tMessages = table.FivenetMsgsMessages.AS("message")
)

type Server struct {
	MailerServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	js       *events.JSWrapper
}

type Params struct {
	fx.In

	DB       *sql.DB
	P        perms.Permissions
	Enricher *mstlystcdata.UserAwareEnricher
	Aud      audit.IAuditer
	JS       *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		p:        p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		js:       p.JS,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterMailerServiceServer(srv, s)
}
