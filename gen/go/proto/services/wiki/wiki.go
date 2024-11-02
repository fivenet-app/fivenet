package wiki

import (
	"context"
	"database/sql"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	WikiServiceServer

	logger *zap.Logger
	db     *sql.DB
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("wiki"),
		db:     p.DB,
	}

	// TODO

	return s
}

func (s *Server) ListPages(ctx context.Context, req *ListPagesRequest) (*ListPagesResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) CreateOrUpdatePage(ctx context.Context, req *CreateOrUpdatePageRequest) (*CreateOrUpdatePageResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) DeletePage(ctx context.Context, req *DeletePageRequest) (*DeletePageResponse, error) {
	// TODO

	return nil, nil
}
