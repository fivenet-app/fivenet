// Code generated by protoc-gen-fronthelper. DO NOT EDIT.
// source: resources/centrum/dispatches.proto
// source: resources/centrum/general.proto
// source: resources/centrum/settings.proto
// source: resources/centrum/units.proto
// source: services/dmv/vehicles.proto
// source: resources/common/database/database.proto
// source: resources/laws/laws.proto
// source: resources/common/i18n.proto
// source: services/rector/rector.proto
// source: services/notificator/notificator.proto
// source: services/docstore/docstore.proto
// source: services/filestore/filestore.proto
// source: resources/jobs/conduct.proto
// source: resources/jobs/qualifications.proto
// source: resources/jobs/requests.proto
// source: resources/jobs/timeclock.proto
// source: resources/permissions/permissions.proto
// source: resources/documents/access.proto
// source: resources/documents/activity.proto
// source: resources/documents/category.proto
// source: resources/documents/comment.proto
// source: resources/documents/documents.proto
// source: resources/documents/requests.proto
// source: resources/documents/templates.proto
// source: resources/vehicles/vehicles.proto
// source: resources/livemap/livemap.proto
// source: resources/rector/audit.proto
// source: services/jobs/conduct.proto
// source: services/jobs/jobs.proto
// source: services/jobs/qualifications.proto
// source: services/jobs/requests.proto
// source: services/jobs/timeclock.proto
// source: services/livemapper/livemap.proto
// source: services/auth/auth.proto
// source: services/completor/completor.proto
// source: services/centrum/centrum.proto
// source: resources/accounts/accounts.proto
// source: resources/accounts/oauth2.proto
// source: resources/notifications/notifications.proto
// source: services/citizenstore/citizenstore.proto
// source: resources/users/jobs.proto
// source: resources/users/users.proto
// source: resources/timestamp/timestamp.proto

export type Perms =
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
	| 'CitizenStoreService.SetUserProps'
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
	| 'JobsRequestsService.CreateOrUpdateRequestType'
	| 'JobsRequestsService.CreateRequest'
	| 'JobsRequestsService.DeleteRequest'
	| 'JobsRequestsService.DeleteRequestComment'
	| 'JobsRequestsService.DeleteRequestType'
	| 'JobsRequestsService.ListRequests'
	| 'JobsRequestsService.UpdateRequest'
	| 'JobsService.ListColleagues'
	| 'JobsService.SetMOTD'
	| 'JobsTimeclockService.ListTimeclock'
	| 'LivemapperService.CreateOrUpdateMarker'
	| 'LivemapperService.DeleteMarker'
	| 'LivemapperService.Stream'
	| 'RectorService.CreateRole'
	| 'RectorService.DeleteRole'
	| 'RectorService.GetJobProps'
	| 'RectorService.GetRoles'
	| 'RectorService.SetJobProps'
	| 'RectorService.UpdateRolePerms'
	| 'RectorService.ViewAuditLog';
