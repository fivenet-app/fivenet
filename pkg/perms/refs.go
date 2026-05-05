package perms

import (
	"strings"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
)

type AttrKind interface {
	attrKind()
}

type (
	StringListAttr   struct{}
	JobListAttr      struct{}
	JobGradeListAttr struct{}
)

func (StringListAttr) attrKind()   {}
func (JobListAttr) attrKind()      {}
func (JobGradeListAttr) attrKind() {}

type PermissionRef struct {
	namespace Namespace
	service   Service
	name      Name
}

func NewPermissionRef(namespace Namespace, service Service, name Name) PermissionRef {
	return PermissionRef{
		namespace: namespace,
		service:   service,
		name:      name,
	}
}

func NewRawPermissionRef(namespace string, service string, name string) PermissionRef {
	return NewPermissionRef(Namespace(namespace), Service(service), Name(name))
}

func PermissionRefFromServiceMethod(serviceMethod string) (PermissionRef, bool) {
	service, method, ok := strings.Cut(serviceMethod, "/")
	if !ok || service == "" || method == "" {
		return PermissionRef{}, false
	}

	namespace, service, ok := splitNamespacedService(service)
	if !ok {
		return PermissionRef{}, false
	}

	return NewRawPermissionRef(namespace, service, method), true
}

func PermissionRefFromGRPCMethod(fullMethod string) (PermissionRef, bool) {
	serviceMethod, ok := strings.CutPrefix(fullMethod, "/services.")
	if !ok {
		return PermissionRef{}, false
	}

	return PermissionRefFromServiceMethod(serviceMethod)
}

func PermissionRefFromProto(perm *permissionspermissions.Permission) (PermissionRef, bool) {
	if perm == nil || perm.GetNamespace() == "" || perm.GetService() == "" || perm.GetName() == "" {
		return PermissionRef{}, false
	}

	return NewRawPermissionRef(perm.GetNamespace(), perm.GetService(), perm.GetName()), true
}

func (p PermissionRef) Namespace() Namespace {
	return p.namespace
}

func (p PermissionRef) Service() Service {
	return p.service
}

func (p PermissionRef) Name() Name {
	return p.name
}

type AttrRef[T AttrKind] struct {
	perm PermissionRef
	key  Key
}

func NewStringListAttrRef(perm PermissionRef, key Key) AttrRef[StringListAttr] {
	return AttrRef[StringListAttr]{perm: perm, key: key}
}

func NewJobListAttrRef(perm PermissionRef, key Key) AttrRef[JobListAttr] {
	return AttrRef[JobListAttr]{perm: perm, key: key}
}

func NewJobGradeListAttrRef(perm PermissionRef, key Key) AttrRef[JobGradeListAttr] {
	return AttrRef[JobGradeListAttr]{perm: perm, key: key}
}

func (a AttrRef[T]) Permission() PermissionRef {
	return a.perm
}

func (a AttrRef[T]) Key() Key {
	return a.key
}

func (ps *Perms) Can(userInfo *userinfo.UserInfo, perm PermissionRef) bool {
	return ps.can(userInfo, perm.namespace, perm.service, perm.name)
}

func (ps *Perms) CanRaw(
	userInfo *userinfo.UserInfo,
	namespace string,
	service string,
	name string,
) bool {
	return ps.Can(userInfo, NewRawPermissionRef(namespace, service, name))
}

func (ps *Perms) CanServiceMethod(userInfo *userinfo.UserInfo, serviceMethod string) bool {
	perm, ok := PermissionRefFromServiceMethod(serviceMethod)
	if !ok {
		return false
	}

	return ps.Can(userInfo, perm)
}

func (ps *Perms) CanProto(
	userInfo *userinfo.UserInfo,
	perm *permissionspermissions.Permission,
) bool {
	ref, ok := PermissionRefFromProto(perm)
	if !ok {
		return false
	}

	return ps.Can(userInfo, ref)
}

func (ps *Perms) AttrStringList(
	userInfo *userinfo.UserInfo,
	attr AttrRef[StringListAttr],
) (*permissionsattributes.StringList, error) {
	return ps.attrStringListRaw(
		userInfo,
		attr.perm.namespace,
		attr.perm.service,
		attr.perm.name,
		attr.key,
	)
}

func (ps *Perms) AttrJobList(
	userInfo *userinfo.UserInfo,
	attr AttrRef[JobListAttr],
) (*permissionsattributes.StringList, error) {
	return ps.attrJobListRaw(
		userInfo,
		attr.perm.namespace,
		attr.perm.service,
		attr.perm.name,
		attr.key,
	)
}

func (ps *Perms) AttrJobGradeList(
	userInfo *userinfo.UserInfo,
	attr AttrRef[JobGradeListAttr],
) (*permissionsattributes.JobGradeList, error) {
	return ps.attrJobGradeListRaw(
		userInfo,
		attr.perm.namespace,
		attr.perm.service,
		attr.perm.name,
		attr.key,
	)
}

func splitNamespacedService(service string) (string, string, bool) {
	idx := strings.LastIndex(service, ".")
	if idx <= 0 || idx == len(service)-1 {
		return "", "", false
	}

	return service[:idx], service[idx+1:], true
}
