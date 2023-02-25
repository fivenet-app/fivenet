package echo

import (
	context "context"
)

type Server struct {
	EchoServiceServer
}

func (s *Server) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{
		Message: req.Message,
		Value:   "TEST",
	}, nil
}
