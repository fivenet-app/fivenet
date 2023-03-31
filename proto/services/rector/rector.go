package rector

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/pkg/perms"
	"go.uber.org/zap"
)

type Server struct {
	RectorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
	}
}

func (s *Server) GetRoles(ctx context.Context, req *GetRolesRequest) (*GetRolesResponse, error) {
	resp := &GetRolesResponse{}

	return resp, nil
}

func (s *Server) UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*UpdateRoleResponse, error) {
	resp := &UpdateRoleResponse{}

	return resp, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	resp := &DeleteRoleResponse{}

	return resp, nil
}
