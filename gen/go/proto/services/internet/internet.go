package internet

import (
	"context"
	"database/sql"

	errorsinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet/errors"
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

func (s *Server) GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error) {
	// TODO

	return nil, errorsinternet.ErrDomainNotFound
}
