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
	esxCompatEnabled = true

	jobs = newFivenetJobsTable("", "jobs", "")
	jobGrades = newFivenetJobsGradesTable("", "job_grades", "")
	licenses = newFivenetLicensesTable("", "licenses", "")

	user = newFivenetUserTable("", "users", "")
	userLicenses = newFivenetUserLicensesTable("", "user_licenses", "")
	ownedVehicles = newFivenetOwnedVehiclesTable("", "owned_vehicles", "")
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
