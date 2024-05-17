package messenger

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
// TODO
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
