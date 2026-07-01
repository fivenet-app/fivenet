package userinfo

import (
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
)

// CanBeSuperuser determines whether an account currently qualifies for superuser mode.
func CanBeSuperuser(
	groups *accounts.AccountGroups,
	license string,
	superuserGroups []string,
	superuserUsers []string,
) bool {
	if groups != nil && groups.ContainsAnyGroup(superuserGroups) {
		return true
	}

	return slices.Contains(superuserUsers, license)
}
