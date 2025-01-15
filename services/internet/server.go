package internet

import (
	"database/sql"

	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	pbinternet.InternetServiceServer
	pbinternet.AdsServiceServer

	db *sql.DB
}

type Params struct {
	fx.In

	DB *sql.DB
}

func NewServer(p Params) *Server {
	s := &Server{
		db: p.DB,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbinternet.RegisterInternetServiceServer(srv, s)
	pbinternet.RegisterAdsServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbinternet.PermsRemap
}
