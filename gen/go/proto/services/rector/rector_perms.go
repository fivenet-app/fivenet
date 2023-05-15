package rector

import (
	"context"
	"errors"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms"
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

func (s *Server) filterPermissions(ctx context.Context, ps []*permissions.Permission) ([]*permissions.Permission, error) {
	filtered := []*permissions.Permission{}

outer:
	for _, p := range ps {
		for i := 0; i < len(ignoredGuardPermissions); i++ {
			if p.GuardName == ignoredGuardPermissions[i] {
				continue outer
			}
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}

func (s *Server) filterPermissionIDs(ctx context.Context, ids []uint64) ([]uint64, error) {
	if len(ids) == 0 {
		return ids, nil
	}

	perms, err := s.p.GetPermissionsByIDs(ids...)
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, perms)
	if err != nil {
		return nil, err
	}

	permIds := make([]uint64, len(filtered))
	for i := 0; i < len(filtered); i++ {
		permIds[i] = filtered[i].Id
	}
	return permIds, nil
}

func (s *Server) filterAttributes(ctx context.Context, attrs []*permissions.RoleAttribute) ([]*permissions.RoleAttribute, error) {
	if len(attrs) == 0 {
		return attrs, nil
	}

	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if s.p.Can(userId, job, jobGrade, common.SuperuserCategoryPerm, common.SuperuserAnyAccessName) {
		return attrs, nil
	}

	for _, a := range attrs {
		attr, err := s.p.GetAttribute(perms.Category(a.Category), perms.Name(a.Name), perms.Key(a.Key))
		if err != nil {
			return nil, err
		}

		switch perms.AttributeTypes(a.Type) {
		case perms.StringListAttributeType:
			if !perms.ValidateStringList(a.Value.GetStringList().Strings, attr.ValidValues.GetStringList().Strings) {
				return nil, fmt.Errorf("invalid option in list")
			}
		case perms.JobListAttributeType:
			if !perms.ValidateJobList(a.Value.GetJobList().Strings, config.C.Game.PermissionRoleJobs) {
				return nil, fmt.Errorf("invalid job in job list")
			}
		case perms.JobGradeListAttributeType:
			if !perms.ValidateJobGradeList(a.Value.GetJobGradeList().Jobs) {
				return nil, fmt.Errorf("invalid job/ grade in job grade list")
			}
		}
	}

	return attrs, nil
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

	fPerms, err := s.filterPermissions(ctx, perms)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp.Role.Permissions = make([]*permissions.Permission, len(fPerms))
	copy(resp.Role.Permissions, fPerms)

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

	if req.Perms != nil {
		if err := s.handlPermissionsUpdate(ctx, role, req.Perms); err != nil {
			return nil, err
		}
	}
	if req.Attrs != nil {
		if err := s.handleAttributeUpdate(ctx, role, req.Attrs); err != nil {
			return nil, err
		}
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &UpdateRolePermsResponse{}, nil
}

func (s *Server) handlPermissionsUpdate(ctx context.Context, role *model.FivenetRoles, permsUpdate *PermsUpdate) error {
	updatePermIds := make([]uint64, len(permsUpdate.ToUpdate))
	for i := 0; i < len(permsUpdate.ToUpdate); i++ {
		updatePermIds[i] = permsUpdate.ToUpdate[i].Id
	}
	toUpdate, err := s.filterPermissionIDs(ctx, updatePermIds)
	if err != nil {
		return InvalidRequestErr
	}

	toDelete, err := s.filterPermissionIDs(ctx, permsUpdate.ToRemove)
	if err != nil {
		return InvalidRequestErr
	}

	if len(toUpdate) > 0 {
		toUpdatePerms := make([]perms.AddPerm, len(permsUpdate.ToUpdate))
		for _, v := range toUpdate {
			for i := 0; i < len(permsUpdate.ToUpdate); i++ {
				if v == permsUpdate.ToUpdate[i].Id {
					toUpdatePerms[i] = perms.AddPerm{
						Id:  permsUpdate.ToUpdate[i].Id,
						Val: permsUpdate.ToUpdate[i].Val,
					}
					break
				}
			}
		}
		if err := s.p.UpdateRolePermissions(role.ID, toUpdatePerms...); err != nil {
			return err
		}
	}
	if len(toDelete) > 0 {
		if err := s.p.RemovePermissionsFromRole(role.ID, toDelete...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleAttributeUpdate(ctx context.Context, role *model.FivenetRoles, attrUpdates *AttrsUpdate) error {
	toUpdate, err := s.filterAttributes(ctx, attrUpdates.ToUpdate)
	if err != nil {
		return InvalidRequestErr
	}

	toDelete, err := s.filterAttributes(ctx, attrUpdates.ToRemove)
	if err != nil {
		return InvalidRequestErr
	}

	// TODO validate each attribute by type

	if len(toUpdate) > 0 {
		if err := s.p.AddOrUpdateAttributesToRole(toUpdate...); err != nil {
			return err
		}
	}
	if len(toDelete) > 0 {
		if err := s.p.RemoveAttributesFromRole(role.ID, toDelete...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetPermissions(ctx context.Context, req *GetPermissionsRequest) (*GetPermissionsResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	perms, err := s.p.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, perms)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &GetPermissionsResponse{}
	resp.Permissions = filtered

	attrs, err := s.p.GetAllAttributes(job)
	if err != nil {
		return nil, InvalidRequestErr
	}
	resp.Attributes = attrs

	return resp, nil
}
