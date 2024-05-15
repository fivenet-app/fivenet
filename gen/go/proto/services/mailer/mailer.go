package mailer

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	MailerServiceServer
}

type Params struct {
	fx.In
}

func NewServer(p Params) *Server {
	return &Server{}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterMailerServiceServer(srv, s)
}
