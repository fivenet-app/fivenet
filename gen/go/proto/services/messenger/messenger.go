package messenger

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tThreads           = table.FivenetMsgsThreads
	tThreadsUserState  = table.FivenetMsgsThreadsUserState
	tThreadsJobAccess  = table.FivenetMsgsThreadsJobAccess
	tThreadsUserAccess = table.FivenetMsgsThreadsUserAccess

	tMessages = table.FivenetMsgsMessages
)

type Server struct {
	MessengerServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
}

type Params struct {
	fx.In

	DB       *sql.DB
	P        perms.Permissions
	Enricher *mstlystcdata.UserAwareEnricher
	Aud      audit.IAuditer
}

func NewServer(p Params) *Server {
	return &Server{}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterMessengerServiceServer(srv, s)
}
