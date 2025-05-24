package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ignoredGuardPermissions = []string{}

func (s *Server) ensureUserCanAccessRole(ctx context.Context, roleId uint64) (*model.FivenetRbacRoles, bool, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, err := s.ps.GetRole(ctx, roleId)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get role %d for user access check. %w", roleId, err)
	}

	if userInfo.Superuser {
		return role, true, nil
	}

	// Make sure the user is from the job
	if role.Job != userInfo.Job {
		return nil, false, errorssettings.ErrInvalidRequest
	}

	if role.Grade > userInfo.JobGrade {
		return nil, false, errorssettings.ErrInvalidRequest
	}

	return role, true, nil
}

func (s *Server) filterPermissions(ctx context.Context, job string, fullFilter bool, ps []*permissions.Permission) ([]*permissions.Permission, error) {
	filtered := []*permissions.Permission{}

	filters, err := s.ps.GetJobPermissions(ctx, job)
	if err != nil {
		return nil, err
	}

outer:
	for _, p := range ps {
		for i := range ignoredGuardPermissions {
			if p.GuardName == ignoredGuardPermissions[i] {
				continue outer
			}
		}

		if fullFilter {
			found := false
			for _, filter := range filters {
				if p.Id == filter.Id {
					if !filter.Val {
						continue outer
					}
					found = true
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}

func (s *Server) filterPermissionIDs(ctx context.Context, job string, full bool, ids []uint64) ([]uint64, error) {
	if len(ids) == 0 {
		return ids, nil
	}

	perms, err := s.ps.GetPermissionsByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, job, full, perms)
	if err != nil {
		return nil, err
	}

	permIds := make([]uint64, len(filtered))
	for i := range filtered {
		permIds[i] = filtered[i].Id
	}
	return permIds, nil
}

func (s *Server) GetRoles(ctx context.Context, req *pbsettings.GetRolesRequest) (*pbsettings.GetRolesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	var roles collections.Roles
	var err error

	if userInfo.Superuser && req.LowestRank != nil && *req.LowestRank {
		roles, err = s.ps.GetRoles(ctx, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		collectedRoles := map[string]*model.FivenetRbacRoles{}
		for _, role := range roles {
			if _, ok := collectedRoles[role.Job]; !ok {
				collectedRoles[role.Job] = role
				continue
			}
		}

		roles = collections.Roles{}
		for _, role := range collectedRoles {
			roles = append(roles, role)
		}
	} else {
		roles, err = s.ps.GetJobRolesUpTo(ctx, userInfo.Job, userInfo.JobGrade)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	resp := &pbsettings.GetRolesResponse{}
	for _, r := range roles {
		role := &permissions.Role{
			Id:          r.ID,
			CreatedAt:   timestamp.New(*r.CreatedAt),
			Job:         r.Job,
			Grade:       r.Grade,
			Permissions: []*permissions.Permission{},
		}

		s.enricher.EnrichJobInfo(role)

		resp.Roles = append(resp.Roles, role)
	}

	return resp, nil
}

func (s *Server) GetRole(ctx context.Context, req *pbsettings.GetRoleRequest) (*pbsettings.GetRoleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.Superuser {
		return nil, errorssettings.ErrNoPermission
	}

	fullFilter := !(userInfo.Superuser && req.Filtered != nil && !*req.Filtered)

	var perms []*permissions.Permission
	if fullFilter {
		perms, err = s.ps.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
		}
	} else {
		perms, err = s.ps.GetJobPermissions(ctx, role.Job)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
		}
	}

	resp := &pbsettings.GetRoleResponse{
		Role: &permissions.Role{
			Id:        role.ID,
			CreatedAt: timestamp.New(*role.CreatedAt),
			Job:       role.Job,
			Grade:     role.Grade,
		},
	}
	s.enricher.EnrichJobInfo(resp.Role)

	fPerms, err := s.filterPermissions(ctx, role.Job, fullFilter, perms)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	resp.Role.Permissions = make([]*permissions.Permission, len(fPerms))
	copy(resp.Role.Permissions, fPerms)

	resp.Role.Attributes, err = s.ps.GetRoleAttributes(role.Job, role.Grade)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return resp, nil
}

func (s *Server) CreateRole(ctx context.Context, req *pbsettings.CreateRoleRequest) (*pbsettings.CreateRoleResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int("fivenet.settings.Grade", int(req.Grade)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "CreateRole",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Make sure the user is from the job or is a super user
	if !userInfo.Superuser {
		if req.Job != userInfo.Job {
			return nil, errorssettings.ErrInvalidRequest
		}
		if req.Grade > userInfo.JobGrade {
			return nil, errorssettings.ErrInvalidRequest
		}
	}

	role, err := s.ps.GetRoleByJobAndGrade(ctx, req.Job, req.Grade)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}
	if role != nil {
		return nil, errorssettings.ErrRoleAlreadyExists
	}

	cr, err := s.ps.CreateRole(ctx, req.Job, req.Grade)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if cr == nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	r := permissions.ConvertFromRole(cr)
	s.enricher.EnrichJobInfo(r)

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	return &pbsettings.CreateRoleResponse{
		Role: r,
	}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *pbsettings.DeleteRoleRequest) (*pbsettings.DeleteRoleResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.settings.role_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "DeleteRole",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.Superuser {
		return nil, errorssettings.ErrNoPermission
	}

	roleCount, err := s.ps.CountRolesForJob(ctx, userInfo.Job)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	// Don't allow deleting the "last" role, one role should always remain
	if roleCount <= 1 {
		return nil, errorssettings.ErrInvalidRequest
	}

	// Don't allow deleting the own or higher role
	if role.Grade >= userInfo.JobGrade {
		return nil, errorssettings.ErrOwnRoleDeletion
	}

	if err := s.ps.DeleteRole(ctx, role.ID); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbsettings.DeleteRoleResponse{}, nil
}

func (s *Server) UpdateRolePerms(ctx context.Context, req *pbsettings.UpdateRolePermsRequest) (*pbsettings.UpdateRolePermsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.settings.role_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "UpdateRolePerms",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.Superuser {
		return nil, errswrap.NewError(err, errorssettings.ErrNoPermission)
	}

	if req.Perms != nil {
		if err := s.handlPermissionsUpdate(ctx, role, req.Perms); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidPerms)
		}
	}
	if req.Attrs != nil {
		if err := s.handleAttributeUpdate(ctx, userInfo, role, req.Attrs); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidAttrs)
		}
	}

	// Send event to every employee
	if _, err := s.js.PublishAsyncProto(ctx,
		fmt.Sprintf("%s.%s.%s.%d", notifi.BaseSubject, notifi.JobGradeTopic, role.Job, role.Grade),
		&notifications.JobGradeEvent{
			Data: &notifications.JobGradeEvent_RefreshToken{
				RefreshToken: true,
			},
		}); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbsettings.UpdateRolePermsResponse{}, nil
}

func (s *Server) handlPermissionsUpdate(ctx context.Context, role *model.FivenetRbacRoles, permsUpdate *pbsettings.PermsUpdate) error {
	updatePermIds := make([]uint64, len(permsUpdate.ToUpdate))
	for i := range permsUpdate.ToUpdate {
		updatePermIds[i] = permsUpdate.ToUpdate[i].Id
	}
	toUpdate, err := s.filterPermissionIDs(ctx, role.Job, false, updatePermIds)
	if err != nil {
		return err
	}

	toDelete, err := s.filterPermissionIDs(ctx, role.Job, false, permsUpdate.ToRemove)
	if err != nil {
		return err
	}

	permsToRemove := []uint64{}
	if len(toUpdate) > 0 {
		toUpdatePerms := make([]perms.AddPerm, len(permsUpdate.ToUpdate))
		for _, v := range toUpdate {
			for i := range permsUpdate.ToUpdate {
				if v == permsUpdate.ToUpdate[i].Id {
					toUpdatePerms[i] = perms.AddPerm{
						Id:  permsUpdate.ToUpdate[i].Id,
						Val: permsUpdate.ToUpdate[i].Val,
					}

					if !permsUpdate.ToUpdate[i].Val {
						permsToRemove = append(permsToRemove, permsUpdate.ToUpdate[i].Id)
					}
					break
				}
			}
		}

		if err := s.ps.UpdateRolePermissions(ctx, role.ID, toUpdatePerms...); err != nil {
			return err
		}
	}

	if len(toDelete) > 0 {
		if err := s.ps.RemovePermissionsFromRole(ctx, role.ID, toDelete...); err != nil {
			return err
		}
	}

	if len(permsToRemove) > 0 {
		for _, perm := range permsToRemove {
			if err := s.ps.RemoveAttributesFromRoleByPermission(ctx, role.ID, perm); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) handleAttributeUpdate(ctx context.Context, userInfo *userinfo.UserInfo, role *model.FivenetRbacRoles, attrUpdates *pbsettings.AttrsUpdate) error {
	if len(attrUpdates.ToUpdate) > 0 {
		if err := s.ps.UpdateRoleAttributes(ctx, userInfo.Job, role.ID, attrUpdates.ToUpdate...); err != nil {
			return err
		}
	}

	if len(attrUpdates.ToRemove) > 0 {
		if err := s.ps.RemoveAttributesFromRole(ctx, role.ID, attrUpdates.ToRemove...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetPermissions(ctx context.Context, req *pbsettings.GetPermissionsRequest) (*pbsettings.GetPermissionsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.settings.role_id", int64(req.RoleId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	perms, err := s.ps.GetAllPermissions(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	fullFilter := !(userInfo.Superuser && req.Filtered != nil && !*req.Filtered)

	filtered, err := s.filterPermissions(ctx, userInfo.Job, fullFilter, perms)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	resp := &pbsettings.GetPermissionsResponse{}
	resp.Permissions = filtered

	role, err := s.ps.GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	if role.Job != userInfo.Job && !userInfo.Superuser {
		return nil, errorssettings.ErrInvalidRequest
	}

	attrs, ok := s.ps.GetJobAttributes(role.Job)
	if !ok {
		return nil, errorssettings.ErrInvalidRequest
	}
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) GetEffectivePermissions(ctx context.Context, req *pbsettings.GetEffectivePermissionsRequest) (*pbsettings.GetEffectivePermissionsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.settings.role_id", int64(req.RoleId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.RoleId)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.Superuser {
		return nil, errswrap.NewError(err, errorssettings.ErrNoPermission)
	}

	perms, err := s.ps.GetEffectiveRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	attrs, err := s.ps.GetEffectiveRoleAttributes(role.Job, role.Grade)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	r := &permissions.Role{
		Id:    role.ID,
		Job:   role.Job,
		Grade: role.Grade,
	}

	s.enricher.EnrichJobInfo(r)

	resp := &pbsettings.GetEffectivePermissionsResponse{}
	resp.Role = r
	resp.Permissions = perms
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) GetAllPermissions(ctx context.Context, req *pbsettings.GetAllPermissionsRequest) (*pbsettings.GetAllPermissionsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "GetAllPermissions",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.Job)
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	perms, err := s.ps.GetAllPermissions(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	attrs, err := s.ps.GetAllAttributes(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	resp := &pbsettings.GetAllPermissionsResponse{}
	resp.Permissions = perms
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) GetJobLimits(ctx context.Context, req *pbsettings.GetJobLimitsRequest) (*pbsettings.GetJobLimitsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "GetJobLimits",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.Job)
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	resp := &pbsettings.GetJobLimitsResponse{}

	perms, err := s.ps.GetJobPermissions(ctx, job.Name)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	resp.Permissions = perms

	attrs, _ := s.ps.GetJobAttributes(job.Name)
	resp.Attributes = attrs

	resp.Job = job.Name
	resp.JobLabel = &job.Label

	return resp, nil
}

func (s *Server) UpdateJobLimits(ctx context.Context, req *pbsettings.UpdateJobLimitsRequest) (*pbsettings.UpdateJobLimitsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "UpdateJobLimits",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.Job)
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	for _, attr := range req.Attrs.ToUpdate {
		if err := s.ps.UpdateJobAttributes(ctx, job.Name, attr.AttrId, attr.MaxValues); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	for _, ps := range req.Perms.ToUpdate {
		if err := s.ps.UpdateJobPermissions(ctx, job.Name, ps.Id, ps.Val); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	for _, ps := range req.Perms.ToRemove {
		if err := s.ps.UpdateJobPermissions(ctx, job.Name, ps, false); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	if err := s.ps.ApplyJobPermissions(ctx, job.Name); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbsettings.UpdateJobLimitsResponse{}, nil
}
