package filestore

import (
	"context"
)

type Server struct {
	FileStoreServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {

	// TODO implement access checks and upload logic

	return nil, nil
}
