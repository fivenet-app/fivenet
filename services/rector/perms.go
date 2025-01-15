package rector

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	rector "github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/perms/collections"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	errorsrector "github.com/fivenet-app/fivenet/services/rector/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"
)

var ignoredGuardPermissions = []string{}

func (s *Server) ensureUserCanAccessRole(ctx context.Context, roleId uint64) (*model.FivenetRoles, bool, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, err := s.ps.GetRole(ctx, roleId)
	if err != nil {
		return nil, false, err
	}

	if userInfo.SuperUser {
		return role, true, nil
	}

	// Make sure the user is from the job
	if role.Job != userInfo.Job {
		return nil, false, errorsrector.ErrInvalidRequest
	}

	if role.Grade > userInfo.JobGrade {
		return nil, false, errorsrector.ErrInvalidRequest
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
		for i := 0; i < len(ignoredGuardPermissions); i++ {
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
	for i := 0; i < len(filtered); i++ {
		permIds[i] = filtered[i].Id
	}
	return permIds, nil
}

func (s *Server) GetRoles(ctx context.Context, req *pbrector.GetRolesRequest) (*pbrector.GetRolesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	var roles collections.Roles
	var err error

	if userInfo.SuperUser && req.LowestRank != nil && *req.LowestRank {
		roles, err = s.ps.GetRoles(ctx, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}

		collectedRoles := map[string]*model.FivenetRoles{}
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
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	resp := &pbrector.GetRolesResponse{}
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

func (s *Server) GetRole(ctx context.Context, req *pbrector.GetRoleRequest) (*pbrector.GetRoleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsrector.ErrNoPermission
	}

	fullFilter := !(userInfo.SuperUser && req.Filtered != nil && !*req.Filtered)

	var perms []*permissions.Permission
	if fullFilter {
		perms, err = s.ps.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
		}
	} else {
		perms, err = s.ps.GetJobPermissions(ctx, role.Job)
		if err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
		}
	}

	resp := &pbrector.GetRoleResponse{
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
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	resp.Role.Permissions = make([]*permissions.Permission, len(fPerms))
	copy(resp.Role.Permissions, fPerms)

	resp.Role.Attributes, err = s.ps.GetRoleAttributes(role.Job, role.Grade)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	return resp, nil
}

func (s *Server) CreateRole(ctx context.Context, req *pbrector.CreateRoleRequest) (*pbrector.CreateRoleResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.rector.job", req.Job))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int("fivenet.rector.Grade", int(req.Grade)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "CreateRole",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Make sure the user is from the job or is a super user
	if !userInfo.SuperUser {
		if req.Job != userInfo.Job {
			return nil, errorsrector.ErrInvalidRequest
		}
		if req.Grade > userInfo.JobGrade {
			return nil, errorsrector.ErrInvalidRequest
		}
	}

	role, err := s.ps.GetRoleByJobAndGrade(ctx, req.Job, req.Grade)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}
	if role != nil {
		return nil, errorsrector.ErrRoleAlreadyExists
	}

	cr, err := s.ps.CreateRole(ctx, req.Job, req.Grade)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	if cr == nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	r := permissions.ConvertFromRole(cr)
	s.enricher.EnrichJobInfo(r)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	return &pbrector.CreateRoleResponse{
		Role: r,
	}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *pbrector.DeleteRoleRequest) (*pbrector.DeleteRoleResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.rector.role_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteRole",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsrector.ErrNoPermission
	}

	roleCount, err := s.ps.CountRolesForJob(ctx, userInfo.Job)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	// Don't allow deleting the "last" role, one role should always remain
	if roleCount <= 1 {
		return nil, errorsrector.ErrInvalidRequest
	}

	// Don't allow deleting the own or higher role
	if role.Grade >= userInfo.JobGrade {
		return nil, errorsrector.ErrOwnRoleDeletion
	}

	if err := s.ps.DeleteRole(ctx, role.ID); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbrector.DeleteRoleResponse{}, nil
}

func (s *Server) UpdateRolePerms(ctx context.Context, req *pbrector.UpdateRolePermsRequest) (*pbrector.UpdateRolePermsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.rector.role_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "UpdateRolePerms",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}
	if !check && !userInfo.SuperUser {
		return nil, errswrap.NewError(err, errorsrector.ErrNoPermission)
	}

	if req.Perms != nil {
		if err := s.handlPermissionsUpdate(ctx, role, req.Perms); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrInvalidPerms)
		}
	}
	if req.Attrs != nil {
		if err := s.handleAttributeUpdate(ctx, userInfo, role, req.Attrs); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrInvalidAttrs)
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
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbrector.UpdateRolePermsResponse{}, nil
}

func (s *Server) handlPermissionsUpdate(ctx context.Context, role *model.FivenetRoles, permsUpdate *pbrector.PermsUpdate) error {
	updatePermIds := make([]uint64, len(permsUpdate.ToUpdate))
	for i := 0; i < len(permsUpdate.ToUpdate); i++ {
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
		if err := s.ps.UpdateRolePermissions(ctx, role.ID, toUpdatePerms...); err != nil {
			return err
		}
	}
	if len(toDelete) > 0 {
		if err := s.ps.RemovePermissionsFromRole(ctx, role.ID, toDelete...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleAttributeUpdate(ctx context.Context, userInfo *userinfo.UserInfo, role *model.FivenetRoles, attrUpdates *pbrector.AttrsUpdate) error {
	if len(attrUpdates.ToUpdate) > 0 {
		if err := s.ps.AddOrUpdateAttributesToRole(ctx, userInfo.Job, role.ID, attrUpdates.ToUpdate...); err != nil {
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

func (s *Server) GetPermissions(ctx context.Context, req *pbrector.GetPermissionsRequest) (*pbrector.GetPermissionsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.rector.role_id", int64(req.RoleId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	perms, err := s.ps.GetAllPermissions(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	fullFilter := !(userInfo.SuperUser && req.Filtered != nil && !*req.Filtered)

	filtered, err := s.filterPermissions(ctx, userInfo.Job, fullFilter, perms)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	resp := &pbrector.GetPermissionsResponse{}
	resp.Permissions = filtered

	role, err := s.ps.GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}

	if role.Job != userInfo.Job && !userInfo.SuperUser {
		return nil, errorsrector.ErrInvalidRequest
	}

	attrs, err := s.ps.GetAllAttributes(ctx, role.Job, role.Grade)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) UpdateRoleLimits(ctx context.Context, req *pbrector.UpdateRoleLimitsRequest) (*pbrector.UpdateRoleLimitsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.rector.role_id", int64(req.RoleId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "UpdateRoleLimits",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	role, err := s.ps.GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrInvalidRequest)
	}
	for _, attr := range req.Attrs.ToUpdate {
		if err := s.ps.UpdateJobAttributeMaxValues(ctx, role.Job, attr.AttrId, attr.MaxValues); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	for _, ps := range req.Perms.ToUpdate {
		if err := s.ps.UpdateJobPermissions(ctx, role.Job, ps.Id, ps.Val); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	for _, ps := range req.Perms.ToRemove {
		if err := s.ps.UpdateJobPermissions(ctx, role.Job, ps, false); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	if err := s.ps.ApplyJobPermissions(ctx, role.Job); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbrector.UpdateRoleLimitsResponse{}, nil
}

func (s *Server) DeleteFaction(ctx context.Context, req *pbrector.DeleteFactionRequest) (*pbrector.DeleteFactionResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.rector.role_id", int64(req.RoleId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteFaction",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	role, err := s.ps.GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.rector.job", role.Job))

	roles, err := s.ps.GetJobRoles(ctx, role.Job)
	if err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	errs := multierr.Combine()
	for _, role := range roles {
		if err := s.ps.DeleteRole(ctx, role.ID); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}

	if err := s.ps.ClearJobAttributes(ctx, role.Job); err != nil {
		errs = multierr.Append(errs, err)
	}

	// Remove job props as last action to remove a faction from the data
	if err := s.deleteJobProps(ctx, role.Job); err != nil {
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		return nil, errswrap.NewError(errs, errorsrector.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbrector.DeleteFactionResponse{}, nil
}

func (s *Server) deleteJobProps(ctx context.Context, job string) error {
	stmt := tJobProps.
		DELETE().
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
