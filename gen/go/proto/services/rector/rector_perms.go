package rector

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms/collections"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidRequestErr    = status.Error(codes.InvalidArgument, "Invalid role action requested!")
	NoPermissionErr      = status.Error(codes.PermissionDenied, "No permission to create/change/delete role!")
	RoleAlreadyExistsErr = status.Error(codes.InvalidArgument, "Role already exists!")
)

var (
	ignoredGuardPermissions = []string{
		"authservice-setjob",
		common.SuperuserAnyAccessGuard,
	}
)

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

func (s *Server) filterPermissions(ctx context.Context, perms collections.Permissions, jobFilter bool) (collections.Permissions, error) {
	userId := auth.GetUserIDFromContext(ctx)
	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, "RectorService.GetPermissions")
	if err != nil {
		return nil, err
	}

	// Disable job filter when superuser
	if s.p.Can(userId, common.SuperuserAnyAccess) {
		jobFilter = false
	}

	filtered := collections.Permissions{}

outer:
	for _, p := range perms {
		for i := 0; i < len(ignoredGuardPermissions); i++ {
			if p.GuardName == ignoredGuardPermissions[i] {
				continue outer
			}
			if jobFilter {
				for _, jc := range config.C.Game.PermissionRoleJobs {
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
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}

func (s *Server) filterPermissionIDs(ctx context.Context, ids []uint64, jobFilter bool) ([]uint64, error) {
	perms, err := s.p.GetPermissionsByIDs(ids...)
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, perms, jobFilter)
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
		return nil, NoPermissionErr
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

	fPerms, err := s.filterPermissions(ctx, perms, false)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp.Role.Permissions = make([]*permissions.Permission, len(fPerms))
	for k := 0; k < len(fPerms); k++ {
		resp.Role.Permissions[k] = permissions.ConvertFromPerm(fPerms[k])
	}

	return resp, nil
}

func (s *Server) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "CreateRole",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	name := fmt.Sprintf("%s - Rank: %d", strings.ToTitle(job), req.Grade)

	role, err := s.p.GetRoleByGuardName(name)
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}
	if role != nil {
		return nil, RoleAlreadyExistsErr
	}

	guard := fmt.Sprintf("job-%s-%d", job, req.Grade)
	description := fmt.Sprintf("Role for %s (Rank: %d)", job, req.Grade)

	cr, err := s.p.CreateRoleWithGuard(name, guard, description)
	if err != nil {
		return nil, err
	}

	if cr == nil {
		return nil, InvalidRequestErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)
	return &CreateRoleResponse{
		Role: permissions.ConvertFromRole(cr),
	}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteRole",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, NoPermissionErr
	}

	jobRoleKey := fmt.Sprintf("job-%s-", job)
	roleCount, err := s.p.CountRoles(jobRoleKey)
	if err != nil {
		return nil, InvalidRequestErr
	}

	// Don't allow deleting the "last" role, one role should always remain
	if roleCount <= 1 {
		return nil, InvalidRequestErr
	}

	if err := s.p.DeleteRole(role.ID); err != nil {
		return nil, InvalidRequestErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteRoleResponse{}, nil
}

func (s *Server) AddPermToRole(ctx context.Context, req *AddPermToRoleRequest) (*AddPermToRoleResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "AddPermToRole",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, NoPermissionErr
	}

	perms, err := s.filterPermissionIDs(ctx, req.Permissions, true)
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

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return resp, nil
}

func (s *Server) RemovePermFromRole(ctx context.Context, req *RemovePermFromRoleRequest) (*RemovePermFromRoleResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "RemovePermFromRole",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, InvalidRequestErr
	}
	if !check {
		return nil, NoPermissionErr
	}

	perms, err := s.filterPermissionIDs(ctx, req.Permissions, true)
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

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return resp, nil
}

func (s *Server) GetPermissions(ctx context.Context, req *GetPermissionsRequest) (*GetPermissionsResponse, error) {
	perms, err := s.p.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, perms, true)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &GetPermissionsResponse{}
	resp.Permissions = make([]*permissions.Permission, len(filtered))
	for k := 0; k < len(filtered); k++ {
		resp.Permissions[k] = permissions.ConvertFromPerm(filtered[k])
	}

	return resp, nil
}
