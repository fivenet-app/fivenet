package mock

import (
	"strings"

	"github.com/galexrt/arpanet/pkg/perms/collections"
	"github.com/galexrt/arpanet/pkg/perms/helpers"
)

type PermsMock struct {
	Perms     map[string]PermCount
	UserPerms map[int32]map[string]PermCount
}

type PermCount struct {
	Count int
	Value bool
}

func NewMock() *PermsMock {
	return &PermsMock{
		Perms:     map[string]PermCount{},
		UserPerms: map[int32]map[string]PermCount{},
	}
}

func (p *PermsMock) GetAllPermissionsOfUser(userID int32) (collections.Permissions, error) {
	// TODO

	return nil, nil
}

func (p *PermsMock) GetAllPermissionsByPrefixOfUser(userID int32, prefix string) (collections.Permissions, error) {
	// TODO

	return nil, nil
}

func (p *PermsMock) GetSuffixOfPermissionsByPrefixOfUser(userID int32, prefix string) ([]string, error) {
	// TODO

	return nil, nil
}

func (p *PermsMock) CanID(userID int32, perm ...string) bool {
	guard := helpers.Guard(strings.Join(perm, "."))

	if pm, ok := p.Perms[guard]; ok {
		pm.Count++
	}

	return false
}
