// Code generated by protoc-gen-fronthelper. DO NOT EDIT.
// source: resources/accounts/accounts.proto
// source: resources/accounts/oauth2.proto
// source: resources/centrum/dispatches.proto
// source: resources/centrum/general.proto
// source: resources/centrum/settings.proto
// source: resources/centrum/units.proto
// source: resources/common/i18n.proto
// source: resources/common/database/database.proto
// source: resources/documents/access.proto
// source: resources/documents/activity.proto
// source: resources/documents/category.proto
// source: resources/documents/comment.proto
// source: resources/documents/documents.proto
// source: resources/documents/requests.proto
// source: resources/documents/templates.proto
// source: resources/filestore/file.proto
// source: resources/jobs/colleagues.proto
// source: resources/jobs/conduct.proto
// source: resources/jobs/timeclock.proto
// source: resources/laws/laws.proto
// source: resources/livemap/livemap.proto
// source: resources/livemap/tracker.proto
// source: resources/notifications/notifications.proto
// source: resources/permissions/permissions.proto
// source: resources/qualifications/qualifications.proto
// source: resources/rector/audit.proto
// source: resources/rector/config.proto
// source: resources/timestamp/timestamp.proto
// source: resources/users/jobs.proto
// source: resources/users/users.proto
// source: resources/vehicles/vehicles.proto
// source: services/auth/auth.proto
// source: services/centrum/centrum.proto
// source: services/citizenstore/citizenstore.proto
// source: services/completor/completor.proto
// source: services/dmv/vehicles.proto
// source: services/docstore/docstore.proto
// source: services/jobs/conduct.proto
// source: services/jobs/jobs.proto
// source: services/jobs/timeclock.proto
// source: services/livemapper/livemap.proto
// source: services/notificator/notificator.proto
// source: services/qualifications/qualifications.proto
// source: services/rector/config.proto
// source: services/rector/filestore.proto
// source: services/rector/laws.proto
// source: services/rector/rector.proto

export type Perms =
    | 'CanBeSuper'
    | 'SuperUser'
    | 'TODOService.TODOMethod'
	| 'AuthService.ChooseCharacter'
	| 'CentrumService.CreateDispatch'
	| 'CentrumService.CreateOrUpdateUnit'
	| 'CentrumService.DeleteDispatch'
	| 'CentrumService.DeleteUnit'
	| 'CentrumService.Stream'
	| 'CentrumService.TakeControl'
	| 'CentrumService.TakeDispatch'
	| 'CentrumService.UpdateDispatch'
	| 'CentrumService.UpdateSettings'
	| 'CitizenStoreService.GetUser'
	| 'CitizenStoreService.ListCitizens'
	| 'CitizenStoreService.ListUserActivity'
	| 'CitizenStoreService.ManageCitizenAttributes'
	| 'CitizenStoreService.SetUserProps'
	| 'CompletorService.CompleteCitizenAttributes'
	| 'CompletorService.CompleteCitizens'
	| 'CompletorService.CompleteDocumentCategories'
	| 'CompletorService.CompleteJobs'
	| 'DMVService.ListVehicles'
	| 'DocStoreService.AddDocumentReference'
	| 'DocStoreService.AddDocumentRelation'
	| 'DocStoreService.ChangeDocumentOwner'
	| 'DocStoreService.CreateCategory'
	| 'DocStoreService.CreateDocument'
	| 'DocStoreService.CreateDocumentReq'
	| 'DocStoreService.CreateTemplate'
	| 'DocStoreService.DeleteCategory'
	| 'DocStoreService.DeleteComment'
	| 'DocStoreService.DeleteDocument'
	| 'DocStoreService.DeleteDocumentReq'
	| 'DocStoreService.DeleteTemplate'
	| 'DocStoreService.GetDocument'
	| 'DocStoreService.ListCategories'
	| 'DocStoreService.ListDocumentActivity'
	| 'DocStoreService.ListDocumentReqs'
	| 'DocStoreService.ListDocuments'
	| 'DocStoreService.ListTemplates'
	| 'DocStoreService.ListUserDocuments'
	| 'DocStoreService.PostComment'
	| 'DocStoreService.ToggleDocument'
	| 'DocStoreService.UpdateDocument'
	| 'JobsConductService.CreateConductEntry'
	| 'JobsConductService.DeleteConductEntry'
	| 'JobsConductService.ListConductEntries'
	| 'JobsConductService.UpdateConductEntry'
	| 'JobsService.GetColleague'
	| 'JobsService.ListColleagueActivity'
	| 'JobsService.ListColleagues'
	| 'JobsService.SetJobsUserProps'
	| 'JobsService.SetMOTD'
	| 'JobsTimeclockService.ListInactiveEmployees'
	| 'JobsTimeclockService.ListTimeclock'
	| 'LivemapperService.CreateOrUpdateMarker'
	| 'LivemapperService.DeleteMarker'
	| 'LivemapperService.Stream'
	| 'QualificationsService.CreateOrUpdateQualificationResult'
	| 'QualificationsService.CreateQualification'
	| 'QualificationsService.DeleteQualification'
	| 'QualificationsService.DeleteQualificationReq'
	| 'QualificationsService.DeleteQualificationResult'
	| 'QualificationsService.GetQualification'
	| 'QualificationsService.ListQualifications'
	| 'QualificationsService.UpdateQualification'
	| 'RectorService.CreateRole'
	| 'RectorService.DeleteRole'
	| 'RectorService.GetJobProps'
	| 'RectorService.GetRoles'
	| 'RectorService.SetJobProps'
	| 'RectorService.UpdateRolePerms'
	| 'RectorService.ViewAuditLog';
