package auth

import "go.uber.org/zap"

type Server struct {
	AccountServiceServer

	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

// TODO
