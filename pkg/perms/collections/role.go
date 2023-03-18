package collections

import (
	"github.com/galexrt/arpanet/query/arpanet/model"
)

// Roles provides methods for you to manage array data more easily.
type Roles []*model.ArpanetRoles

// Origin convert the collection to role array.
// @return []models.ArpanetRoles
func (u Roles) Origin() []*model.ArpanetRoles {
	return []*model.ArpanetRoles(u)
}

// Len returns the number of elements of the array.
// @return int64
func (u Roles) Len() (length int64) {
	return int64(len(u))
}

// IDs returns an array of the role array's ids.
// @return []uint
func (u Roles) IDs() (IDs []uint64) {
	for _, role := range u {
		IDs = append(IDs, role.ID)
	}
	return IDs
}

// Names returns an array of the role array's names.
// @return []string
func (u Roles) Names() (names []string) {
	for _, role := range u {
		names = append(names, role.Name)
	}
	return names
}

// GuardNames returns an array of the permission array's guard names.
// @return []string
func (u Roles) GuardNames() (guards []string) {
	for _, role := range u {
		guards = append(guards, role.GuardName)
	}
	return guards
}
