package collections

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
)

// Roles provides methods for you to manage array data more easily.
type Roles []*permissions.Role

// Origin convert the collection to role array.
//
//	@return	[]models.ArpanetRoles
func (r Roles) Origin() []*permissions.Role {
	return []*permissions.Role(r)
}

// Len returns the number of elements of the array.
//
//	@return	int64
func (u Roles) Len() int64 {
	return int64(len(u))
}

// IDs returns an array of the role array's ids.
//
//	@return	[]uint
func (r Roles) IDs() []int64 {
	ids := make([]int64, 0, len(r))
	for _, role := range r {
		ids = append(ids, role.GetId())
	}
	return ids
}
