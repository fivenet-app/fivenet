package rector

import (
	"context"
	"errors"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms"
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
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	role, err := s.p.GetRole(roleId)
	if err != nil {
		return nil, false, err
	}

	// Make sure the user is from the job
	if role.Job != job {
		return nil, false, InvalidRequestErr
	}

	if role.Grade > jobGrade {
		return nil, false, InvalidRequestErr
	}

	return role, true, nil
}

func (s *Server) filterPermissions(ctx context.Context, ps collections.Permissions, jobFilter bool) (collections.Permissions, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	jobsAttr, err := s.p.Attr(userId, job, jobGrade, RectorServicePerm, RectorServiceGetPermissionsPerm, perms.Key("Jobs"))
	if err != nil {
		return nil, err
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.(perms.StringList)
	}

	// Disable job filter when superuser
	if s.p.Can(userId, job, jobGrade, common.SuperuserCategoryPerm, common.SuperuserAnyAccessName) {
		jobFilter = false
	}

	filtered := collections.Permissions{}

outer:
	for _, p := range ps {
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
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	roles, err := s.p.GetJobRolesUpTo(job, jobGrade)
	if err != nil {
		return nil, err
	}

	resp := &GetRolesResponse{}
	for _, r := range roles {
		role := &permissions.Role{
			Id:          r.ID,
			CreatedAt:   timestamp.New(*r.CreatedAt),
			Job:         r.Job,
			Grade:       r.Grade,
			Permissions: []*permissions.Permission{},
		}

		s.c.EnrichJobInfo(role)

		resp.Roles = append(resp.Roles, role)
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

	resp.Role = &permissions.Role{
		Id:        role.ID,
		CreatedAt: timestamp.New(*role.CreatedAt),
		Job:       role.Job,
		Grade:     role.Grade,
	}
	s.c.EnrichJobInfo(resp.Role)

	fPerms, err := s.filterPermissions(ctx, perms, false)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp.Role.Permissions = make([]*permissions.Permission, len(fPerms))
	for k := 0; k < len(fPerms); k++ {
		resp.Role.Permissions[k] = permissions.ConvertFromPerm(fPerms[k])
	}

	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	resp.Role.Attributes, err = s.p.GetRoleAttributes(job, jobGrade)
	if err != nil {
		return nil, InvalidRequestErr
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

	role, err := s.p.GetRoleByJobAndGrade(job, req.Grade)
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}
	if role != nil {
		return nil, RoleAlreadyExistsErr
	}

	cr, err := s.p.CreateRole(job, req.Grade)
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
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

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

	roleCount, err := s.p.CountRolesForJob(job)
	if err != nil {
		return nil, InvalidRequestErr
	}

	// Don't allow deleting the "last" role, one role should always remain
	if roleCount <= 1 {
		return nil, InvalidRequestErr
	}

	// Don't allow deleting the own or higher role
	if role.Grade >= jobGrade {
		return nil, InvalidRequestErr
	}

	if err := s.p.DeleteRole(role.ID); err != nil {
		return nil, InvalidRequestErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteRoleResponse{}, nil
}

func (s *Server) UpdateRolePerms(ctx context.Context, req *UpdateRolePermsRequest) (*UpdateRolePermsResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "UpdateRolePerms",
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

	toAdd, err := s.filterPermissionIDs(ctx, req.ToAdd, true)
	if err != nil {
		return nil, InvalidRequestErr
	}

	toDelete, err := s.filterPermissionIDs(ctx, req.ToRemove, true)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &UpdateRolePermsResponse{}
	if len(toAdd) > 0 {
		if err := s.p.AddPermissionsToRole(role.ID, toAdd...); err != nil {
			return nil, err
		}
	}
	if len(toDelete) > 0 {
		if err := s.p.RemovePermissionsFromRole(role.ID, toDelete...); err != nil {
			return nil, err
		}
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

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
