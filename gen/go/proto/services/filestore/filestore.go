package filestore

import (
	"context"

	"github.com/galexrt/fivenet/pkg/storage"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	FileStoreServiceServer

	st *storage.Storage
}

type Params struct {
	fx.In

	Storage *storage.Storage
}

func NewServer(p Params) *Server {
	return &Server{
		st: p.Storage,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterFileStoreServiceServer(srv, s)
}

func (s *Server) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {

	// TODO implement access checks and upload logic

	return nil, nil
}
