package rector

import (
	"context"
	"errors"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
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
	OwnRoleDeletionErr   = status.Error(codes.InvalidArgument, "Can't delete your own role!")
)

var (
	ignoredGuardPermissions = []string{
		"authservice-setjob",
	}
)

func (s *Server) ensureUserCanAccessRole(ctx context.Context, roleId uint64) (*model.FivenetRoles, bool, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, err := s.p.GetRole(ctx, roleId)
	if err != nil {
		return nil, false, err
	}

	// Make sure the user is from the job
	if role.Job != userInfo.Job {
		return nil, false, InvalidRequestErr
	}

	if role.Grade > userInfo.JobGrade {
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

	perms, err := s.p.GetPermissionsByIDs(ctx, ids...)
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

func (s *Server) filterAttributes(ctx context.Context, attrs []*permissions.RoleAttribute) error {
	if len(attrs) == 0 {
		return nil
	}

	for i := 0; i < len(attrs); i++ {
		dbAttrs, err := s.p.GetAttributeByIDs(ctx, attrs[i].AttrId)
		if err != nil {
			return err
		}
		attr := dbAttrs[0]

		// If the attribute valid values is null, nothing to validate
		if attrs[i].Value != nil && attrs[i].Value.ValidValues != nil && attr.ValidValues != nil {
			switch perms.AttributeTypes(attr.Type) {
			case perms.StringListAttributeType:
				if attr.ValidValues.GetStringList() != nil && attr.ValidValues.GetStringList().Strings != nil {
					if !perms.ValidateStringList(attrs[i].Value.GetStringList().Strings, attr.ValidValues.GetStringList().Strings) {
						return fmt.Errorf("invalid option in list")
					}
				}
			case perms.JobListAttributeType:
				if attr.ValidValues.GetJobList() != nil && attr.ValidValues.GetJobList().Strings != nil {
					if !perms.ValidateJobList(attrs[i].Value.GetJobList().Strings, attr.ValidValues.GetJobList().Strings) {
						return fmt.Errorf("invalid job in job list")
					}
				}
			case perms.JobGradeListAttributeType:
				if attr.ValidValues.GetJobGradeList() != nil && attr.ValidValues.GetJobGradeList().Jobs != nil {
					if !perms.ValidateJobGradeList(attrs[i].Value.GetJobGradeList().Jobs) {
						return fmt.Errorf("invalid job/ grade in job grade list")
					}
				}
			}
		}

		attrs[i].AttrId = attr.AttrId
		attrs[i].Type = attr.Type
	}

	return nil
}

func (s *Server) GetRoles(ctx context.Context, req *GetRolesRequest) (*GetRolesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	roles, err := s.p.GetJobRolesUpTo(ctx, userInfo.Job, userInfo.JobGrade)
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

	perms, err := s.p.GetRolePermissions(ctx, role.ID)
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

	resp.Role.Attributes, err = s.p.GetRoleAttributes(role.Job, role.Grade)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "CreateRole",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	role, err := s.p.GetRoleByJobAndGrade(ctx, userInfo.Job, req.Grade)
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}
	if role != nil {
		return nil, RoleAlreadyExistsErr
	}

	cr, err := s.p.CreateRole(ctx, userInfo.Job, req.Grade)
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
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteRole",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
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

	roleCount, err := s.p.CountRolesForJob(ctx, userInfo.Job)
	if err != nil {
		return nil, InvalidRequestErr
	}

	// Don't allow deleting the "last" role, one role should always remain
	if roleCount <= 1 {
		return nil, InvalidRequestErr
	}

	// Don't allow deleting the own or higher role
	if role.Grade >= userInfo.JobGrade {
		return nil, OwnRoleDeletionErr
	}

	if err := s.p.DeleteRole(ctx, role.ID); err != nil {
		return nil, InvalidRequestErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteRoleResponse{}, nil
}

func (s *Server) UpdateRolePerms(ctx context.Context, req *UpdateRolePermsRequest) (*UpdateRolePermsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "UpdateRolePerms",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
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
		if err := s.p.UpdateRolePermissions(ctx, role.ID, toUpdatePerms...); err != nil {
			return err
		}
	}
	if len(toDelete) > 0 {
		if err := s.p.RemovePermissionsFromRole(ctx, role.ID, toDelete...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleAttributeUpdate(ctx context.Context, role *model.FivenetRoles, attrUpdates *AttrsUpdate) error {
	if err := s.filterAttributes(ctx, attrUpdates.ToUpdate); err != nil {
		return err
	}

	if err := s.filterAttributes(ctx, attrUpdates.ToRemove); err != nil {
		return err
	}

	if len(attrUpdates.ToUpdate) > 0 {
		if err := s.p.AddOrUpdateAttributesToRole(ctx, role.ID, attrUpdates.ToUpdate...); err != nil {
			return err
		}
	}
	if len(attrUpdates.ToRemove) > 0 {
		if err := s.p.RemoveAttributesFromRole(ctx, role.ID, attrUpdates.ToRemove...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetPermissions(ctx context.Context, req *GetPermissionsRequest) (*GetPermissionsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	perms, err := s.p.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, perms)
	if err != nil {
		return nil, InvalidRequestErr
	}

	resp := &GetPermissionsResponse{}
	resp.Permissions = filtered

	attrs, err := s.p.GetAllAttributes(ctx, userInfo.Job)
	if err != nil {
		return nil, InvalidRequestErr
	}
	resp.Attributes = attrs

	return resp, nil
}
