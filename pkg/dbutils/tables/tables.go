package tables

import "sync"

var (
	jobs      = FivenetJobs
	jobGrades = FivenetJobsGrades
	licenses  = FivenetLicenses

	user          = FivenetUser
	userLicenses  = FivenetUserLicenses
	ownedVehicles = FivenetOwnedVehicles
)

var once sync.Once

var esxCompatEnabled = false

// EnableESXCompat called to enable ESX compat mode, overrides the `fivenet_` prefixed tables with the ESX names.
func EnableESXCompat() {
	once.Do(setESXTableNames)
}

func IsESXCompatEnabled() bool {
	return esxCompatEnabled
}

func setESXTableNames() {
	panic("ESX compat mode has been removed and is no longer supported.")
}

func Jobs() *FivenetJobsTable {
	return jobs
}

func JobsGrades() *FivenetJobsGradesTable {
	return jobGrades
}

func Licenses() *FivenetLicensesTable {
	return licenses
}

func User() *FivenetUserTable {
	return user
}

func UserLicenses() *FivenetUserLicensesTable {
	return userLicenses
}

func OwnedVehicles() *FivenetOwnedVehiclesTable {
	return ownedVehicles
}
