package tables

import "sync"

var (
	jobs      = FivenetJobs
	jobGrades = FivenetJobGrades
	licenses  = FivenetLicenses

	users         = FivenetUsers
	userLicenses  = FivenetUserLicenses
	ownedVehicles = FivenetOwnedVehicles
)

var once sync.Once

// Called to enable ESX compat mode, overrides the `fivenet_` prefixed tables with the ESX names
func EnableESXCompat() {
	once.Do(setESXTableNames)
}

func setESXTableNames() {
	jobs = newFivenetJobsTable("", "jobs", "")
	jobGrades = newFivenetJobGradesTable("", "job_grades", "")
	licenses = newFivenetLicensesTable("", "licenses", "")

	users = newFivenetUsersTable("", "users", "")
	userLicenses = newFivenetUserLicensesTable("", "user_licenses", "")
	ownedVehicles = newFivenetOwnedVehiclesTable("", "owned_vehicles", "")
}

func Jobs() *FivenetJobsTable {
	return jobs
}

func JobGrades() *FivenetJobGradesTable {
	return jobGrades
}

func Licenses() *FivenetLicensesTable {
	return licenses
}

func Users() *FivenetUsersTable {
	return users
}

func UserLicenses() *FivenetUserLicensesTable {
	return userLicenses
}

func OwnedVehicles() *FivenetOwnedVehiclesTable {
	return ownedVehicles
}
