package rector

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/proto/resources/permissions"
	"github.com/galexrt/fivenet/proto/resources/timestamp"
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
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	rolePrefix := "job-" + job + "-"

	roles, err := s.p.GetRoles(rolePrefix)
	if err != nil {
		return nil, err
	}

	resp := &GetRolesResponse{}
	_ = roles

	for _, r := range roles {
		var updatedAt *timestamp.Timestamp
		if r.UpdatedAt != nil {
			updatedAt = timestamp.New(*r.UpdatedAt)
		}

		resp.Roles = append(resp.Roles, &permissions.Role{
			Id:          r.ID,
			CreatedAt:   timestamp.New(*r.CreatedAt),
			UpdatedAt:   updatedAt,
			Name:        r.Name,
			GuardName:   r.GuardName,
			Description: *r.Description,
			Permissions: []*permissions.Permission{},
		})
	}

	return resp, nil
}

func (s *Server) UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*UpdateRoleResponse, error) {

	// TODO

	return &UpdateRoleResponse{}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error) {

	// TODO

	return &DeleteRoleResponse{}, nil
}
