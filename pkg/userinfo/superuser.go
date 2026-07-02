package userinfo

import (
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
)

func canBeWithLists(
	groups *accounts.AccountGroups,
	license string,
	allowedGroups []string,
	allowedUsers []string,
) bool {
	if groups != nil && groups.ContainsAnyGroup(allowedGroups) {
		return true
	}

	return slices.Contains(allowedUsers, license)
}

// CanBeSuperuser determines whether an account currently qualifies for job-admin superuser mode.
func CanBeSuperuser(
	groups *accounts.AccountGroups,
	license string,
	superuserGroups []string,
	superuserUsers []string,
) bool {
	return canBeWithLists(groups, license, superuserGroups, superuserUsers)
}

// CanBeConfigAdmin determines whether an account can access config-admin gated screens and RPCs.
func CanBeConfigAdmin(
	groups *accounts.AccountGroups,
	license string,
	configAdminGroups []string,
	configAdminUsers []string,
) bool {
	return canBeWithLists(groups, license, configAdminGroups, configAdminUsers)
}

// EffectiveAdminLists returns the merged job-admin and config-admin lists after applying app-config overrides.
func EffectiveAdminLists(
	baseJobAdminGroups []string,
	baseJobAdminUsers []string,
	baseConfigAdminGroups []string,
	baseConfigAdminUsers []string,
	appCfg appconfig.IConfig,
) (
	jobAdminGroups []string,
	jobAdminUsers []string,
	configAdminGroups []string,
	configAdminUsers []string,
) {
	jobAdminGroups = utils.MergeUniqueStrings(baseJobAdminGroups)
	jobAdminUsers = utils.MergeUniqueStrings(baseJobAdminUsers)
	configAdminGroups = utils.MergeUniqueStrings(baseConfigAdminGroups)
	configAdminUsers = utils.MergeUniqueStrings(baseConfigAdminUsers)

	if appCfg == nil || appCfg.Get() == nil || appCfg.Get().GetAuth() == nil {
		return jobAdminGroups, jobAdminUsers, configAdminGroups, configAdminUsers
	}

	appAuth := appCfg.Get().GetAuth()
	jobAdminGroups = utils.MergeUniqueStrings(jobAdminGroups, appAuth.GetJobAdminGroups())
	jobAdminUsers = utils.MergeUniqueStrings(jobAdminUsers, appAuth.GetJobAdminUsers())
	configAdminGroups = utils.MergeUniqueStrings(configAdminGroups, appAuth.GetConfigAdminGroups())
	configAdminUsers = utils.MergeUniqueStrings(configAdminUsers, appAuth.GetConfigAdminUsers())

	return jobAdminGroups, jobAdminUsers, configAdminGroups, configAdminUsers
}

// EffectiveJobAdminLists returns the merged admin lists relevant for job-admin eligibility.
func EffectiveJobAdminLists(
	baseJobAdminGroups []string,
	baseJobAdminUsers []string,
	baseConfigAdminGroups []string,
	baseConfigAdminUsers []string,
	appCfg appconfig.IConfig,
) (jobAdminGroups []string, jobAdminUsers []string) {
	jobAdminGroups, jobAdminUsers, configAdminGroups, configAdminUsers := EffectiveAdminLists(
		baseJobAdminGroups,
		baseJobAdminUsers,
		baseConfigAdminGroups,
		baseConfigAdminUsers,
		appCfg,
	)

	return utils.MergeUniqueStrings(jobAdminGroups, configAdminGroups),
		utils.MergeUniqueStrings(jobAdminUsers, configAdminUsers)
}
