package permsstub

import (
	"context"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
)

type Permissions struct {
	CanFunc            func(*pbuserinfo.UserInfo, perms.PermissionRef) bool
	CanServiceMethodFn func(*pbuserinfo.UserInfo, string) bool
}

var _ perms.Permissions = (*Permissions)(nil)

func (p *Permissions) SetDefaultRolePerms(context.Context, []string) error { return nil }

func (p *Permissions) GetAllPermissions(context.Context) ([]*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) GetPermissionsByIDs(
	context.Context,
	...int64,
) ([]*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) GetPermission(
	context.Context,
	perms.Namespace,
	perms.Service,
	perms.Name,
) (*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) GetPermissionsOfUser(*pbuserinfo.UserInfo) ([]*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) GetAllAttributes(context.Context) ([]*permissionsattributes.RoleAttribute, error) {
	return nil, nil
}

func (p *Permissions) GetRoles(context.Context, bool) ([]*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) GetRole(context.Context, int64) (*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) GetRoleByJobAndGrade(
	context.Context,
	string,
	int32,
) (*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) GetJobRoles(context.Context, string) ([]*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) GetJobRolesUpTo(
	context.Context,
	string,
	int32,
) ([]*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) GetClosestJobRole(
	context.Context,
	string,
	int32,
) (*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) CountRolesForJob(context.Context, string) (int64, error) { return 0, nil }

func (p *Permissions) CreateRole(context.Context, string, int32) (*permissionspermissions.Role, error) {
	return nil, nil
}

func (p *Permissions) DeleteRole(context.Context, int64) error { return nil }

func (p *Permissions) GetRolePermissions(
	context.Context,
	int64,
) ([]*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) GetEffectiveRolePermissions(
	context.Context,
	int64,
) ([]*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) UpdateRolePermissions(context.Context, int64, ...perms.AddPerm) error {
	return nil
}

func (p *Permissions) RemovePermissionsFromRole(context.Context, int64, ...int64) error { return nil }

func (p *Permissions) GetRoleAttributes(
	context.Context,
	string,
	int32,
) ([]*permissionsattributes.RoleAttribute, error) {
	return nil, nil
}

func (p *Permissions) GetEffectiveRoleAttributes(
	context.Context,
	string,
	int32,
) ([]*permissionsattributes.RoleAttribute, error) {
	return nil, nil
}

func (p *Permissions) UpdateRoleAttributes(
	context.Context,
	string,
	int64,
	...*permissionsattributes.RoleAttribute,
) error {
	return nil
}

func (p *Permissions) RemoveAttributesFromRole(
	context.Context,
	int64,
	...*permissionsattributes.RoleAttribute,
) error {
	return nil
}

func (p *Permissions) RemoveAttributesFromRoleByPermission(context.Context, int64, int64) error {
	return nil
}

func (p *Permissions) GetJobPermissions(context.Context, string) ([]*permissionspermissions.Permission, error) {
	return nil, nil
}

func (p *Permissions) UpdateJobPermissions(
	context.Context,
	string,
	...*permissionspermissions.PermItem,
) error {
	return nil
}

func (p *Permissions) ApplyJobPermissions(context.Context, string) error { return nil }

func (p *Permissions) ClearJobPermissions(context.Context, string) error { return nil }

func (p *Permissions) GetJobAttributes(
	context.Context,
	string,
) ([]*permissionsattributes.RoleAttribute, error) {
	return nil, nil
}

func (p *Permissions) UpdateJobAttributes(
	context.Context,
	string,
	...*permissionsattributes.RoleAttribute,
) error {
	return nil
}

func (p *Permissions) ClearJobAttributes(context.Context, string) error { return nil }

func (p *Permissions) Can(userInfo *pbuserinfo.UserInfo, perm perms.PermissionRef) bool {
	if p.CanFunc != nil {
		return p.CanFunc(userInfo, perm)
	}
	return false
}

func (p *Permissions) CanRaw(
	userInfo *pbuserinfo.UserInfo,
	namespace string,
	service string,
	name string,
) bool {
	return p.Can(userInfo, perms.NewRawPermissionRef(namespace, service, name))
}

func (p *Permissions) CanServiceMethod(userInfo *pbuserinfo.UserInfo, serviceMethod string) bool {
	if p.CanServiceMethodFn != nil {
		return p.CanServiceMethodFn(userInfo, serviceMethod)
	}

	perm, ok := perms.PermissionRefFromServiceMethod(serviceMethod)
	if !ok {
		return false
	}

	return p.Can(userInfo, perm)
}

func (p *Permissions) CanProto(
	userInfo *pbuserinfo.UserInfo,
	perm *permissionspermissions.Permission,
) bool {
	ref, ok := perms.PermissionRefFromProto(perm)
	if !ok {
		return false
	}

	return p.Can(userInfo, ref)
}

func (p *Permissions) Attr(
	*pbuserinfo.UserInfo,
	perms.Namespace,
	perms.Service,
	perms.Name,
	perms.Key,
) (*permissionsattributes.AttributeValues, error) {
	return nil, nil
}

func (p *Permissions) AttrStringList(*pbuserinfo.UserInfo, perms.AttrRef[perms.StringListAttr]) (*permissionsattributes.StringList, error) {
	return nil, nil
}

func (p *Permissions) AttrJobList(*pbuserinfo.UserInfo, perms.AttrRef[perms.JobListAttr]) (*permissionsattributes.StringList, error) {
	return nil, nil
}

func (p *Permissions) AttrJobGradeList(*pbuserinfo.UserInfo, perms.AttrRef[perms.JobGradeListAttr]) (*permissionsattributes.JobGradeList, error) {
	return nil, nil
}
