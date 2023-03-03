package options

import (
	"github.com/galexrt/arpanet/pkg/permify/utils"
)

// RoleOption represents options when fetching roles.
type RoleOption struct {
	WithPermissions bool
	Pagination      *utils.Pagination
}
