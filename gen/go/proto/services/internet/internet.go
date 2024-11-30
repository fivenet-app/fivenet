package internet

import (
	"database/sql"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	InternetServiceServer
	AdsServiceServer

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
	RegisterInternetServiceServer(srv, s)
	RegisterAdsServiceServer(srv, s)
}
