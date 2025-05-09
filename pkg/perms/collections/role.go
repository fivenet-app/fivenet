package collections

import (
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

// Roles provides methods for you to manage array data more easily.
type Roles []*model.FivenetRoles

// Origin convert the collection to role array.
// @return []models.ArpanetRoles
func (u Roles) Origin() []*model.FivenetRoles {
	return []*model.FivenetRoles(u)
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
