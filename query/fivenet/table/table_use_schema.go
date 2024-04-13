//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	FivenetAccounts = FivenetAccounts.FromSchema(schema)
	FivenetAttrs = FivenetAttrs.FromSchema(schema)
	FivenetAuditLog = FivenetAuditLog.FromSchema(schema)
	FivenetCentrumDispatches = FivenetCentrumDispatches.FromSchema(schema)
	FivenetCentrumDispatchesAsgmts = FivenetCentrumDispatchesAsgmts.FromSchema(schema)
	FivenetCentrumDispatchesStatus = FivenetCentrumDispatchesStatus.FromSchema(schema)
	FivenetCentrumMarkers = FivenetCentrumMarkers.FromSchema(schema)
	FivenetCentrumSettings = FivenetCentrumSettings.FromSchema(schema)
	FivenetCentrumUnits = FivenetCentrumUnits.FromSchema(schema)
	FivenetCentrumUnitsStatus = FivenetCentrumUnitsStatus.FromSchema(schema)
	FivenetCentrumUnitsUsers = FivenetCentrumUnitsUsers.FromSchema(schema)
	FivenetCentrumUsers = FivenetCentrumUsers.FromSchema(schema)
	FivenetConfig = FivenetConfig.FromSchema(schema)
	FivenetDocuments = FivenetDocuments.FromSchema(schema)
	FivenetDocumentsActivity = FivenetDocumentsActivity.FromSchema(schema)
	FivenetDocumentsCategories = FivenetDocumentsCategories.FromSchema(schema)
	FivenetDocumentsComments = FivenetDocumentsComments.FromSchema(schema)
	FivenetDocumentsJobAccess = FivenetDocumentsJobAccess.FromSchema(schema)
	FivenetDocumentsReferences = FivenetDocumentsReferences.FromSchema(schema)
	FivenetDocumentsRelations = FivenetDocumentsRelations.FromSchema(schema)
	FivenetDocumentsRequests = FivenetDocumentsRequests.FromSchema(schema)
	FivenetDocumentsTemplates = FivenetDocumentsTemplates.FromSchema(schema)
	FivenetDocumentsTemplatesJobAccess = FivenetDocumentsTemplatesJobAccess.FromSchema(schema)
	FivenetDocumentsUserAccess = FivenetDocumentsUserAccess.FromSchema(schema)
	FivenetJobAttrs = FivenetJobAttrs.FromSchema(schema)
	FivenetJobCitizenAttributes = FivenetJobCitizenAttributes.FromSchema(schema)
	FivenetJobPermissions = FivenetJobPermissions.FromSchema(schema)
	FivenetJobProps = FivenetJobProps.FromSchema(schema)
	FivenetJobsConduct = FivenetJobsConduct.FromSchema(schema)
	FivenetJobsTimeclock = FivenetJobsTimeclock.FromSchema(schema)
	FivenetJobsUserActivity = FivenetJobsUserActivity.FromSchema(schema)
	FivenetJobsUserProps = FivenetJobsUserProps.FromSchema(schema)
	FivenetLawbooks = FivenetLawbooks.FromSchema(schema)
	FivenetLawbooksLaws = FivenetLawbooksLaws.FromSchema(schema)
	FivenetNotifications = FivenetNotifications.FromSchema(schema)
	FivenetOauth2Accounts = FivenetOauth2Accounts.FromSchema(schema)
	FivenetPermissions = FivenetPermissions.FromSchema(schema)
	FivenetQualifications = FivenetQualifications.FromSchema(schema)
	FivenetQualificationsJobAccess = FivenetQualificationsJobAccess.FromSchema(schema)
	FivenetQualificationsRequests = FivenetQualificationsRequests.FromSchema(schema)
	FivenetQualificationsRequirements = FivenetQualificationsRequirements.FromSchema(schema)
	FivenetQualificationsResults = FivenetQualificationsResults.FromSchema(schema)
	FivenetRoleAttrs = FivenetRoleAttrs.FromSchema(schema)
	FivenetRolePermissions = FivenetRolePermissions.FromSchema(schema)
	FivenetRoles = FivenetRoles.FromSchema(schema)
	FivenetUserActivity = FivenetUserActivity.FromSchema(schema)
	FivenetUserLocations = FivenetUserLocations.FromSchema(schema)
	FivenetUserProps = FivenetUserProps.FromSchema(schema)
	GksphoneJobMessage = GksphoneJobMessage.FromSchema(schema)
	GksphoneSettings = GksphoneSettings.FromSchema(schema)
	JobGrades = JobGrades.FromSchema(schema)
	Jobs = Jobs.FromSchema(schema)
	Licenses = Licenses.FromSchema(schema)
	OwnedVehicles = OwnedVehicles.FromSchema(schema)
	UserLicenses = UserLicenses.FromSchema(schema)
	Users = Users.FromSchema(schema)
}
