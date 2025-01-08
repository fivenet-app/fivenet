package tables

var (
	Jobs      = FivenetJobs
	JobGrades = FivenetJobGrades
	Licenses  = FivenetLicenses

	Users         = FivenetUsers
	UserLicenses  = FivenetUserLicenses
	OwnedVehicles = FivenetOwnedVehicles
)

// Called to enable ESX compat mode, overrides the `fivenet_` prefixed tables with the ESX names
func EnableESXCompat() {
	FivenetJobs = newFivenetJobsTable("", "jobs", "")
	FivenetJobGrades = newFivenetJobGradesTable("", "job_grades", "")
	FivenetLicenses = newFivenetLicensesTable("", "licenses", "")

	FivenetUsers = newFivenetUsersTable("", "users", "")
	FivenetUserLicenses = newFivenetUserLicensesTable("", "user_licenses", "")
	FivenetOwnedVehicles = newFivenetOwnedVehiclesTable("", "owned_vehicles", "")
}
