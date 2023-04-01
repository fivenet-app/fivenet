package rector

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/perms/collections"
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

var (
	ignoredGuardPermissions = []string{
		"authservice-setjob",
		"superuser-override",
	}
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

func (s *Server) ensureUserCanAccessRoleByGuardName(ctx context.Context, name string) (*model.FivenetRoles, bool, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	role, err := s.p.GetRoleByGuardName(name)
	if err != nil {
		return nil, false, err
	}

	// Make sure the user is from the job
	if !strings.HasPrefix(role.GuardName, "job-"+job+"-") {
		return nil, false, InvalidRequestErr
	}

	return role, true, nil
}

func (s *Server) filterPermissions(ctx context.Context, perms collections.Permissions) (collections.Permissions, error) {
	userId := auth.GetUserIDFromContext(ctx)
	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, "RectorService.GetPermissions")
	if err != nil {
		return nil, err
	}

	filtered := collections.Permissions{}

outer:
	for _, p := range perms {
		for i := 0; i < len(ignoredGuardPermissions); i++ {
			if p.GuardName == ignoredGuardPermissions[i] {
				continue outer
			}
			for _, jc := range config.C.FiveM.PermissionRoleJobs {
				if strings.HasSuffix(p.GuardName, "-"+jc) {
					if len(jobs) == 0 {
						continue outer
					}
					for _, j := range jobs {
						if !strings.HasSuffix(p.GuardName, "-"+j) {
							continue outer
						}
					}
				}
			}
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}

func (s *Server) filterPermissionIDs(ctx context.Context, ids []uint64) ([]uint64, error) {
	perms, err := s.p.GetPermissionsByIDs(ids...)
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, perms)
	if err != nil {
		return nil, err
	}

	return filtered.IDs(), nil
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
	}

	fPerms, err := s.filterPermissions(ctx, perms)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp.Role.Permissions = make([]*permissions.Permission, len(fPerms))
	for k := 0; k < len(fPerms); k++ {
		var createdAt *timestamp.Timestamp
		if fPerms[k].CreatedAt != nil {
			createdAt = timestamp.New(*fPerms[k].CreatedAt)
		}
		var updatedAt *timestamp.Timestamp
		if fPerms[k].UpdatedAt != nil {
			updatedAt = timestamp.New(*fPerms[k].UpdatedAt)
		}

		resp.Role.Permissions[k] = &permissions.Permission{
			Id:          fPerms[k].ID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Name:        fPerms[k].Name,
			GuardName:   fPerms[k].GuardName,
			Description: fPerms[k].Description,
		}
	}

	return resp, nil
}

func (s *Server) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	name := fmt.Sprintf("job-%s-%d", req.Job, req.Grade)

	role, check, err := s.ensureUserCanAccessRoleByGuardName(ctx, name)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	if role != nil {
		return nil, InvalidRequestErr
	}

	description := fmt.Sprintf("Role for ambulance %s (Rank: %d)", req.Job, req.Grade)

	roleId, err := s.p.CreateRole(name, description)
	if err != nil {
		return nil, err
	}

	return &CreateRoleResponse{
		Id: roleId,
	}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	if err := s.p.DeleteRole(role.ID); err != nil {
		return nil, err
	}

	return &DeleteRoleResponse{}, nil
}

func (s *Server) AddPermToRole(ctx context.Context, req *AddPermToRoleRequest) (*AddPermToRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	perms, err := s.filterPermissionIDs(ctx, req.Permissions)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &AddPermToRoleResponse{}
	if len(perms) == 0 {
		return resp, nil
	}

	if err := s.p.AddPermissionsToRole(role.ID, perms); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) RemovePermFromRole(ctx context.Context, req *RemovePermFromRoleRequest) (*RemovePermFromRoleResponse, error) {
	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, InvalidRequestErr
	}

	perms, err := s.filterPermissionIDs(ctx, req.Permissions)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &RemovePermFromRoleResponse{}
	if len(perms) == 0 {
		return resp, nil
	}

	if err := s.p.RemovePermissionsFromRole(role.ID, perms); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetPermissions(ctx context.Context, req *GetPermissionsRequest) (*GetPermissionsResponse, error) {
	perms, err := s.p.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissionIDs(ctx, perms.IDs())
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &GetPermissionsResponse{}
	resp.Permissions = make([]*permissions.Permission, len(filtered))
	for k := 0; k < len(filtered); k++ {
		var createdAt *timestamp.Timestamp
		if perms[k].CreatedAt != nil {
			createdAt = timestamp.New(*perms[k].CreatedAt)
		}
		var updatedAt *timestamp.Timestamp
		if perms[k].UpdatedAt != nil {
			updatedAt = timestamp.New(*perms[k].UpdatedAt)
		}

		resp.Permissions[k] = &permissions.Permission{
			Id:          perms[k].ID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Name:        perms[k].Name,
			GuardName:   perms[k].GuardName,
			Description: perms[k].Description,
		}
	}

	return resp, nil
}
