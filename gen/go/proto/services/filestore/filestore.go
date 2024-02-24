package filestore

import (
	"context"

	"github.com/galexrt/fivenet/pkg/storage"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	FileStoreServiceServer

	st storage.IStorage
}

type Params struct {
	fx.In

	Storage storage.IStorage
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

	// TODO sanitize image file (https://www.libvips.org/), implement access checks and upload file process

	return nil, nil
}

func (s *Server) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {

	return nil, nil
}
