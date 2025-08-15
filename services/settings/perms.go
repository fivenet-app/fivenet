package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms/collections"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/multierr"
)

var ignoredGuardPermissions = []string{}

func (s *Server) ensureUserCanAccessRole(
	ctx context.Context,
	roleId int64,
) (*permissions.Role, bool, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, err := s.ps.GetRole(ctx, roleId)
	if err != nil {
		return nil, false, fmt.Errorf(
			"failed to get role %d for user access check. %w",
			roleId,
			err,
		)
	}

	if userInfo.GetSuperuser() {
		return role, true, nil
	}

	// Make sure the user is from the job
	if role.GetJob() != userInfo.GetJob() {
		return nil, false, errorssettings.ErrInvalidRequest
	}

	if role.GetGrade() > userInfo.GetJobGrade() {
		return nil, false, errorssettings.ErrInvalidRequest
	}

	return role, true, nil
}

func (s *Server) filterPermissions(
	ctx context.Context,
	job string,
	ps []*permissions.Permission,
) ([]*permissions.Permission, error) {
	filtered := []*permissions.Permission{}

	filters, err := s.ps.GetJobPermissions(ctx, job)
	if err != nil {
		return nil, err
	}

outer:
	for _, p := range ps {
		for i := range ignoredGuardPermissions {
			if p.GetGuardName() == ignoredGuardPermissions[i] {
				continue outer
			}
		}

		found := false
		for _, filter := range filters {
			if p.GetId() == filter.GetId() {
				if !filter.GetVal() {
					continue outer
				}
				found = true
			}
		}
		if !found {
			continue
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}

func (s *Server) filterPermissionIDs(
	ctx context.Context,
	job string,
	ids []int64,
) ([]int64, error) {
	if len(ids) == 0 {
		return ids, nil
	}

	perms, err := s.ps.GetPermissionsByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}

	filtered, err := s.filterPermissions(ctx, job, perms)
	if err != nil {
		return nil, err
	}

	permIds := make([]int64, len(filtered))
	for i := range filtered {
		permIds[i] = filtered[i].GetId()
	}
	return permIds, nil
}

func (s *Server) GetRoles(
	ctx context.Context,
	req *pbsettings.GetRolesRequest,
) (*pbsettings.GetRolesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	var roles collections.Roles
	var err error

	if userInfo.GetSuperuser() && req.LowestRank != nil && req.GetLowestRank() {
		roles, err = s.ps.GetRoles(ctx, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}

		collectedRoles := map[string]*permissions.Role{}
		for _, role := range roles {
			if _, ok := collectedRoles[role.GetJob()]; !ok {
				collectedRoles[role.GetJob()] = role
				continue
			}
		}

		roles = collections.Roles{}
		for _, role := range collectedRoles {
			roles = append(roles, role)
		}
	} else {
		roles, err = s.ps.GetJobRolesUpTo(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	resp := &pbsettings.GetRolesResponse{}
	for _, r := range roles {
		role := &permissions.Role{
			Id:          r.GetId(),
			CreatedAt:   r.GetCreatedAt(),
			Job:         r.GetJob(),
			Grade:       r.GetGrade(),
			Permissions: []*permissions.Permission{},
		}

		s.enricher.EnrichJobInfoNoFallback(role)

		resp.Roles = append(resp.Roles, role)
	}

	return resp, nil
}

func (s *Server) GetRole(
	ctx context.Context,
	req *pbsettings.GetRoleRequest,
) (*pbsettings.GetRoleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorssettings.ErrNoPermission
	}

	perms, err := s.ps.GetRolePermissions(ctx, role.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	fPerms, err := s.filterPermissions(ctx, role.GetJob(), perms)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	resp := &pbsettings.GetRoleResponse{
		Role: &permissions.Role{
			Id:        role.GetId(),
			CreatedAt: role.GetCreatedAt(),
			Job:       role.GetJob(),
			Grade:     role.GetGrade(),

			Permissions: fPerms,
		},
	}
	resp.Role.Attributes, err = s.ps.GetRoleAttributes(ctx, role.GetJob(), role.GetGrade())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	s.enricher.EnrichJobInfoNoFallback(resp.GetRole())

	return resp, nil
}

func (s *Server) CreateRole(
	ctx context.Context,
	req *pbsettings.CreateRoleRequest,
) (*pbsettings.CreateRoleResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.settings.job", req.GetJob(),
		"fivenet.settings.grade", req.GetGrade(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "CreateRole",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Make sure the user is from the job or is a super user
	if !userInfo.GetSuperuser() {
		if req.GetJob() != userInfo.GetJob() {
			return nil, errorssettings.ErrInvalidRequest
		}
		if req.GetGrade() > userInfo.GetJobGrade() {
			return nil, errorssettings.ErrInvalidRequest
		}
	}

	role, err := s.ps.GetRoleByJobAndGrade(ctx, req.GetJob(), req.GetGrade())
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}
	if role != nil {
		return nil, errorssettings.ErrRoleAlreadyExists
	}

	r, err := s.ps.CreateRole(ctx, req.GetJob(), req.GetGrade())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if r == nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	s.enricher.EnrichJobInfoNoFallback(r)

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	return &pbsettings.CreateRoleResponse{
		Role: r,
	}, nil
}

func (s *Server) DeleteRole(
	ctx context.Context,
	req *pbsettings.DeleteRoleRequest,
) (*pbsettings.DeleteRoleResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.role_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "DeleteRole",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorssettings.ErrNoPermission
	}

	roleCount, err := s.ps.CountRolesForJob(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	// Don't allow deleting the "last" role, one role should always remain
	if roleCount <= 1 {
		return nil, errorssettings.ErrInvalidRequest
	}

	// Don't allow deleting the own or higher role
	if role.GetGrade() >= userInfo.GetJobGrade() {
		return nil, errorssettings.ErrOwnRoleDeletion
	}

	if err := s.ps.DeleteRole(ctx, role.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbsettings.DeleteRoleResponse{}, nil
}

func (s *Server) UpdateRolePerms(
	ctx context.Context,
	req *pbsettings.UpdateRolePermsRequest,
) (*pbsettings.UpdateRolePermsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.role_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "UpdateRolePerms",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errswrap.NewError(err, errorssettings.ErrNoPermission)
	}

	if req.GetPerms() != nil {
		if err := s.handlPermissionsUpdate(ctx, role, req.GetPerms()); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidPerms)
		}
	}
	if req.GetAttrs() != nil {
		if err := s.handleAttributeUpdate(ctx, userInfo, role, req.GetAttrs()); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidAttrs)
		}
	}

	// Send event to job grade employees
	if _, err := s.js.PublishAsyncProto(ctx,
		fmt.Sprintf("%s.%s.%s.%d", notifi.BaseSubject, notifi.JobGradeTopic, role.GetJob(), role.GetGrade()),
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

func (s *Server) handlPermissionsUpdate(
	ctx context.Context,
	role *permissions.Role,
	permsUpdate *settings.PermsUpdate,
) error {
	updatePermIds := make([]int64, len(permsUpdate.GetToUpdate()))
	for i := range permsUpdate.GetToUpdate() {
		updatePermIds[i] = permsUpdate.GetToUpdate()[i].GetId()
	}
	toUpdate, err := s.filterPermissionIDs(ctx, role.GetJob(), updatePermIds)
	if err != nil {
		return err
	}

	removePermIds := make([]int64, len(permsUpdate.GetToRemove()))
	for i := range permsUpdate.GetToRemove() {
		removePermIds[i] = permsUpdate.GetToUpdate()[i].GetId()
	}
	toDelete, err := s.filterPermissionIDs(ctx, role.GetJob(), removePermIds)
	if err != nil {
		return err
	}

	permsToRemove := []int64{}
	if len(toUpdate) > 0 {
		toUpdatePerms := make([]perms.AddPerm, len(permsUpdate.GetToUpdate()))
		for _, v := range toUpdate {
			for i := range permsUpdate.GetToUpdate() {
				if v == permsUpdate.GetToUpdate()[i].GetId() {
					toUpdatePerms[i] = perms.AddPerm{
						Id:  permsUpdate.GetToUpdate()[i].GetId(),
						Val: permsUpdate.GetToUpdate()[i].GetVal(),
					}

					if !permsUpdate.GetToUpdate()[i].GetVal() {
						permsToRemove = append(permsToRemove, permsUpdate.GetToUpdate()[i].GetId())
					}
					break
				}
			}
		}

		if err := s.ps.UpdateRolePermissions(ctx, role.GetId(), toUpdatePerms...); err != nil {
			return err
		}
	}

	if len(toDelete) > 0 {
		if err := s.ps.RemovePermissionsFromRole(ctx, role.GetId(), toDelete...); err != nil {
			return err
		}
	}

	if len(permsToRemove) > 0 {
		for _, perm := range permsToRemove {
			if err := s.ps.RemoveAttributesFromRoleByPermission(ctx, role.GetId(), perm); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) handleAttributeUpdate(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	role *permissions.Role,
	attrUpdates *settings.AttrsUpdate,
) error {
	if len(attrUpdates.GetToUpdate()) > 0 {
		if err := s.ps.UpdateRoleAttributes(ctx, userInfo.GetJob(), role.GetId(), attrUpdates.GetToUpdate()...); err != nil {
			return err
		}
	}

	if len(attrUpdates.GetToRemove()) > 0 {
		if err := s.ps.RemoveAttributesFromRole(ctx, role.GetId(), attrUpdates.GetToRemove()...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetPermissions(
	ctx context.Context,
	req *pbsettings.GetPermissionsRequest,
) (*pbsettings.GetPermissionsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.role_id", req.GetRoleId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	perms, err := s.ps.GetAllPermissions(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	filtered, err := s.filterPermissions(ctx, userInfo.GetJob(), perms)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	resp := &pbsettings.GetPermissionsResponse{}
	resp.Permissions = filtered

	role, err := s.ps.GetRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}

	if role.GetJob() != userInfo.GetJob() && !userInfo.GetSuperuser() {
		return nil, errorssettings.ErrInvalidRequest
	}

	attrs, err := s.ps.GetJobAttributes(ctx, role.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) GetEffectivePermissions(
	ctx context.Context,
	req *pbsettings.GetEffectivePermissionsRequest,
) (*pbsettings.GetEffectivePermissionsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.role_id", req.GetRoleId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	role, check, err := s.ensureUserCanAccessRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrInvalidRequest)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errswrap.NewError(err, errorssettings.ErrNoPermission)
	}

	perms, err := s.ps.GetEffectiveRolePermissions(ctx, role.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	attrs, err := s.ps.GetEffectiveRoleAttributes(ctx, role.GetJob(), role.GetGrade())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	r := &permissions.Role{
		Id:    role.GetId(),
		Job:   role.GetJob(),
		Grade: role.GetGrade(),
	}

	s.enricher.EnrichJobInfoNoFallback(r)

	resp := &pbsettings.GetEffectivePermissionsResponse{}
	resp.Role = r
	resp.Permissions = perms
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) DeleteFaction(
	ctx context.Context,
	req *pbsettings.DeleteFactionRequest,
) (*pbsettings.DeleteFactionResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.job", req.GetJob()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "DeleteFaction",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.job", req.GetJob()})

	roles, err := s.ps.GetJobRoles(ctx, req.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	errs := multierr.Combine()
	for _, role := range roles {
		if err := s.ps.DeleteRole(ctx, role.GetId()); err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}

	if err := s.ps.ClearJobAttributes(ctx, req.GetJob()); err != nil {
		errs = multierr.Append(errs, err)
		return nil, errswrap.NewError(errs, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.ClearJobPermissions(ctx, req.GetJob()); err != nil {
		errs = multierr.Append(errs, err)
		return nil, errswrap.NewError(errs, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.ApplyJobPermissions(ctx, req.GetJob()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	// Set job props to be deleted as last action to start the removal of a faction and it's data from the database
	if err := s.deleteJobProps(ctx, s.db, req.GetJob()); err != nil {
		errs = multierr.Append(errs, err)
	}

	if errs != nil {
		return nil, errswrap.NewError(errs, errorssettings.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbsettings.DeleteFactionResponse{}, nil
}

func (s *Server) deleteJobProps(ctx context.Context, tx qrm.DB, job string) error {
	stmt := tJobProps.
		UPDATE().
		SET(
			tJobProps.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tJobProps.Job.EQ(jet.String(job)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
