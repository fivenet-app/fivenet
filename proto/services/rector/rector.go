package rector

import (
	"context"
	"database/sql"
	"strings"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/proto/resources/permissions"
	"github.com/galexrt/fivenet/proto/resources/timestamp"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidRequestErr = status.Error(codes.InvalidArgument, "Invalid role action requested!")
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

func (s *Server) ensureUserCanAccessRole(ctx context.Context, roleId uint64) (*model.FivenetRoles, bool, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	role, err := s.p.GetRole(roleId)
	if err != nil {
		return nil, false, err
	}

	// Make sure the user is from the job
	if !strings.HasPrefix(role.GuardName, "job-"+job+"-") {
		return nil, false, InvalidRequestErr
	}

	return role, true, nil
}

func (s *Server) GetRoles(ctx context.Context, req *GetRolesRequest) (*GetRolesResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	rolePrefix := "job-" + job + "-"

	roles, err := s.p.GetRoles(rolePrefix)
	if err != nil {
		return nil, err
	}

	resp := &GetRolesResponse{}
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
			Description: r.Description,
			Permissions: []*permissions.Permission{},
		})
	}

	return resp, nil
}

func (s *Server) GetRole(ctx context.Context, req *GetRoleRequest) (*GetRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	perms, err := s.p.GetRolePermissions(role.ID)
	if err != nil {
		return nil, InvalidRequestErr
	}

	rPerms := make([]*permissions.Permission, len(perms))
	for k := 0; k < len(perms); k++ {
		var createdAt *timestamp.Timestamp
		if perms[k].CreatedAt != nil {
			createdAt = timestamp.New(*perms[k].CreatedAt)
		}
		var updatedAt *timestamp.Timestamp
		if perms[k].UpdatedAt != nil {
			updatedAt = timestamp.New(*perms[k].UpdatedAt)
		}

		rPerms[k] = &permissions.Permission{
			Id:          perms[k].ID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Name:        perms[k].Name,
			GuardName:   perms[k].GuardName,
			Description: perms[k].Description,
		}
	}

	resp := &GetRoleResponse{}
	var updatedAt *timestamp.Timestamp
	if role.UpdatedAt != nil {
		updatedAt = timestamp.New(*role.UpdatedAt)
	}

	resp.Role = &permissions.Role{
		Id:          role.ID,
		CreatedAt:   timestamp.New(*role.CreatedAt),
		UpdatedAt:   updatedAt,
		Name:        role.Name,
		GuardName:   role.GuardName,
		Description: role.Description,
		Permissions: rPerms,
	}

	return resp, nil
}

func (s *Server) AddPermToRole(ctx context.Context, req *AddPermToRoleRequest) (*AddPermToRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	// TODO
	_ = role

	return &AddPermToRoleResponse{}, nil
}

func (s *Server) RemovePermFromRole(ctx context.Context, req *RemovePermFromRoleRequest) (*RemovePermFromRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	// TODO
	_ = role

	return &RemovePermFromRoleResponse{}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	// TODO
	_ = role

	return &DeleteRoleResponse{}, nil
}
