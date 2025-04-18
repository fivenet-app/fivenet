// Code generated by protoc-gen-fronthelper. DO NOT EDIT.
// source: resources/accounts/accounts.proto
// source: resources/accounts/oauth2.proto
// source: resources/calendar/access.proto
// source: resources/calendar/calendar.proto
// source: resources/centrum/access.proto
// source: resources/centrum/attributes.proto
// source: resources/centrum/dispatches.proto
// source: resources/centrum/disponents.proto
// source: resources/centrum/settings.proto
// source: resources/centrum/units.proto
// source: resources/centrum/user_unit.proto
// source: resources/common/error.proto
// source: resources/common/i18n.proto
// source: resources/common/uuid.proto
// source: resources/common/content/content.proto
// source: resources/common/cron/cron.proto
// source: resources/common/database/database.proto
// source: resources/common/grpcws/grpcws.proto
// source: resources/common/tests/objects.proto
// source: resources/documents/access.proto
// source: resources/documents/activity.proto
// source: resources/documents/category.proto
// source: resources/documents/comment.proto
// source: resources/documents/documents.proto
// source: resources/documents/requests.proto
// source: resources/documents/templates.proto
// source: resources/documents/workflow.proto
// source: resources/filestore/file.proto
// source: resources/internet/access.proto
// source: resources/internet/ads.proto
// source: resources/internet/domain.proto
// source: resources/internet/page.proto
// source: resources/internet/search.proto
// source: resources/jobs/activity.proto
// source: resources/jobs/colleagues.proto
// source: resources/jobs/conduct.proto
// source: resources/jobs/labels.proto
// source: resources/jobs/timeclock.proto
// source: resources/laws/laws.proto
// source: resources/livemap/livemap.proto
// source: resources/livemap/tracker.proto
// source: resources/mailer/access.proto
// source: resources/mailer/email.proto
// source: resources/mailer/events.proto
// source: resources/mailer/message.proto
// source: resources/mailer/settings.proto
// source: resources/mailer/template.proto
// source: resources/mailer/thread.proto
// source: resources/notifications/events.proto
// source: resources/notifications/notifications.proto
// source: resources/permissions/permissions.proto
// source: resources/qualifications/access.proto
// source: resources/qualifications/exam.proto
// source: resources/qualifications/qualifications.proto
// source: resources/rector/audit.proto
// source: resources/rector/banner.proto
// source: resources/rector/config.proto
// source: resources/stats/stats.proto
// source: resources/sync/activity.proto
// source: resources/sync/data.proto
// source: resources/timestamp/timestamp.proto
// source: resources/users/activity.proto
// source: resources/users/job_props.proto
// source: resources/users/job_settings.proto
// source: resources/users/jobs.proto
// source: resources/users/labels.proto
// source: resources/users/licenses.proto
// source: resources/users/props.proto
// source: resources/users/users.proto
// source: resources/vehicles/vehicles.proto
// source: resources/wiki/access.proto
// source: resources/wiki/activity.proto
// source: resources/wiki/page.proto
// source: services/auth/auth.proto
// source: services/calendar/calendar.proto
// source: services/centrum/centrum.proto
// source: services/citizenstore/citizenstore.proto
// source: services/completor/completor.proto
// source: services/dmv/vehicles.proto
// source: services/docstore/docstore.proto
// source: services/internet/ads.proto
// source: services/internet/domain.proto
// source: services/internet/internet.proto
// source: services/jobs/conduct.proto
// source: services/jobs/jobs.proto
// source: services/jobs/timeclock.proto
// source: services/livemapper/livemap.proto
// source: services/mailer/mailer.proto
// source: services/notificator/notificator.proto
// source: services/qualifications/qualifications.proto
// source: services/rector/config.proto
// source: services/rector/filestore.proto
// source: services/rector/laws.proto
// source: services/rector/rector.proto
// source: services/stats/stats.proto
// source: services/sync/sync.proto
// source: services/wiki/wiki.proto

export const grpcServices = [
    'AdsService',
    'AuthService',
    'CalendarService',
    'CentrumService',
    'CitizenStoreService',
    'CompletorService',
    'DMVService',
    'DocStoreService',
    'DomainService',
    'InternetService',
    'JobsConductService',
    'JobsService',
    'JobsTimeclockService',
    'LivemapperService',
    'MailerService',
    'NotificatorService',
    'QualificationsService',
    'RectorConfigService',
    'RectorFilestoreService',
    'RectorLawsService',
    'RectorService',
    'StatsService',
    'SyncService',
    'WikiService',
];

export const grpcMethods = [
	'AdsService/GetAds',
	'AuthService/Login',
	'AuthService/Logout',
	'AuthService/CreateAccount',
	'AuthService/ChangeUsername',
	'AuthService/ChangePassword',
	'AuthService/ForgotPassword',
	'AuthService/GetCharacters',
	'AuthService/ChooseCharacter',
	'AuthService/GetAccountInfo',
	'AuthService/DeleteOAuth2Connection',
	'AuthService/SetSuperUserMode',
	'CalendarService/ListCalendars',
	'CalendarService/GetCalendar',
	'CalendarService/CreateCalendar',
	'CalendarService/UpdateCalendar',
	'CalendarService/DeleteCalendar',
	'CalendarService/ListCalendarEntries',
	'CalendarService/GetUpcomingEntries',
	'CalendarService/GetCalendarEntry',
	'CalendarService/CreateOrUpdateCalendarEntry',
	'CalendarService/DeleteCalendarEntry',
	'CalendarService/ShareCalendarEntry',
	'CalendarService/ListCalendarEntryRSVP',
	'CalendarService/RSVPCalendarEntry',
	'CalendarService/ListSubscriptions',
	'CalendarService/SubscribeToCalendar',
	'CentrumService/UpdateSettings',
	'CentrumService/CreateDispatch',
	'CentrumService/UpdateDispatch',
	'CentrumService/DeleteDispatch',
	'CentrumService/TakeControl',
	'CentrumService/AssignDispatch',
	'CentrumService/AssignUnit',
	'CentrumService/Stream',
	'CentrumService/GetSettings',
	'CentrumService/JoinUnit',
	'CentrumService/ListUnits',
	'CentrumService/ListUnitActivity',
	'CentrumService/GetDispatch',
	'CentrumService/ListDispatches',
	'CentrumService/ListDispatchActivity',
	'CentrumService/CreateOrUpdateUnit',
	'CentrumService/DeleteUnit',
	'CentrumService/TakeDispatch',
	'CentrumService/UpdateUnitStatus',
	'CentrumService/UpdateDispatchStatus',
	'CitizenStoreService/ListCitizens',
	'CitizenStoreService/GetUser',
	'CitizenStoreService/ListUserActivity',
	'CitizenStoreService/SetUserProps',
	'CitizenStoreService/SetProfilePicture',
	'CitizenStoreService/ManageCitizenLabels',
	'CompletorService/CompleteCitizens',
	'CompletorService/CompleteJobs',
	'CompletorService/CompleteDocumentCategories',
	'CompletorService/ListLawBooks',
	'CompletorService/CompleteCitizenLabels',
	'DMVService/ListVehicles',
	'DocStoreService/ListTemplates',
	'DocStoreService/GetTemplate',
	'DocStoreService/CreateTemplate',
	'DocStoreService/UpdateTemplate',
	'DocStoreService/DeleteTemplate',
	'DocStoreService/ListDocuments',
	'DocStoreService/GetDocument',
	'DocStoreService/CreateDocument',
	'DocStoreService/UpdateDocument',
	'DocStoreService/DeleteDocument',
	'DocStoreService/ToggleDocument',
	'DocStoreService/ChangeDocumentOwner',
	'DocStoreService/GetDocumentReferences',
	'DocStoreService/GetDocumentRelations',
	'DocStoreService/AddDocumentReference',
	'DocStoreService/RemoveDocumentReference',
	'DocStoreService/AddDocumentRelation',
	'DocStoreService/RemoveDocumentRelation',
	'DocStoreService/GetComments',
	'DocStoreService/PostComment',
	'DocStoreService/EditComment',
	'DocStoreService/DeleteComment',
	'DocStoreService/GetDocumentAccess',
	'DocStoreService/SetDocumentAccess',
	'DocStoreService/ListDocumentActivity',
	'DocStoreService/ListDocumentReqs',
	'DocStoreService/CreateDocumentReq',
	'DocStoreService/UpdateDocumentReq',
	'DocStoreService/DeleteDocumentReq',
	'DocStoreService/ListUserDocuments',
	'DocStoreService/ListCategories',
	'DocStoreService/CreateCategory',
	'DocStoreService/UpdateCategory',
	'DocStoreService/DeleteCategory',
	'DocStoreService/ListDocumentPins',
	'DocStoreService/ToggleDocumentPin',
	'DocStoreService/SetDocumentReminder',
	'DomainService/ListTLDs',
	'DomainService/CheckDomainAvailability',
	'DomainService/RegisterDomain',
	'DomainService/ListDomains',
	'DomainService/UpdateDomain',
	'InternetService/Search',
	'InternetService/GetPage',
	'JobsConductService/ListConductEntries',
	'JobsConductService/CreateConductEntry',
	'JobsConductService/UpdateConductEntry',
	'JobsConductService/DeleteConductEntry',
	'JobsService/ListColleagues',
	'JobsService/GetSelf',
	'JobsService/GetColleague',
	'JobsService/ListColleagueActivity',
	'JobsService/SetJobsUserProps',
	'JobsService/GetColleagueLabels',
	'JobsService/ManageColleagueLabels',
	'JobsService/GetColleagueLabelsStats',
	'JobsService/GetMOTD',
	'JobsService/SetMOTD',
	'JobsTimeclockService/ListTimeclock',
	'JobsTimeclockService/GetTimeclockStats',
	'JobsTimeclockService/ListInactiveEmployees',
	'LivemapperService/Stream',
	'LivemapperService/CreateOrUpdateMarker',
	'LivemapperService/DeleteMarker',
	'MailerService/ListEmails',
	'MailerService/GetEmail',
	'MailerService/CreateOrUpdateEmail',
	'MailerService/DeleteEmail',
	'MailerService/GetEmailProposals',
	'MailerService/ListTemplates',
	'MailerService/GetTemplate',
	'MailerService/CreateOrUpdateTemplate',
	'MailerService/DeleteTemplate',
	'MailerService/ListThreads',
	'MailerService/GetThread',
	'MailerService/CreateThread',
	'MailerService/DeleteThread',
	'MailerService/GetThreadState',
	'MailerService/SetThreadState',
	'MailerService/SearchThreads',
	'MailerService/ListThreadMessages',
	'MailerService/PostMessage',
	'MailerService/DeleteMessage',
	'MailerService/GetEmailSettings',
	'MailerService/SetEmailSettings',
	'NotificatorService/GetNotifications',
	'NotificatorService/MarkNotifications',
	'NotificatorService/Stream',
	'QualificationsService/ListQualifications',
	'QualificationsService/GetQualification',
	'QualificationsService/CreateQualification',
	'QualificationsService/UpdateQualification',
	'QualificationsService/DeleteQualification',
	'QualificationsService/ListQualificationRequests',
	'QualificationsService/CreateOrUpdateQualificationRequest',
	'QualificationsService/DeleteQualificationReq',
	'QualificationsService/ListQualificationsResults',
	'QualificationsService/CreateOrUpdateQualificationResult',
	'QualificationsService/DeleteQualificationResult',
	'QualificationsService/GetExamInfo',
	'QualificationsService/TakeExam',
	'QualificationsService/SubmitExam',
	'QualificationsService/GetUserExam',
	'RectorConfigService/GetAppConfig',
	'RectorConfigService/UpdateAppConfig',
	'RectorFilestoreService/ListFiles',
	'RectorFilestoreService/UploadFile',
	'RectorFilestoreService/DeleteFile',
	'RectorLawsService/CreateOrUpdateLawBook',
	'RectorLawsService/DeleteLawBook',
	'RectorLawsService/CreateOrUpdateLaw',
	'RectorLawsService/DeleteLaw',
	'RectorService/GetJobProps',
	'RectorService/SetJobProps',
	'RectorService/GetRoles',
	'RectorService/GetRole',
	'RectorService/CreateRole',
	'RectorService/DeleteRole',
	'RectorService/UpdateRolePerms',
	'RectorService/GetPermissions',
	'RectorService/ViewAuditLog',
	'RectorService/UpdateRoleLimits',
	'RectorService/DeleteFaction',
	'StatsService/GetStats',
	'SyncService/GetStatus',
	'SyncService/AddActivity',
	'SyncService/RegisterAccount',
	'SyncService/TransferAccount',
	'SyncService/SendData',
	'SyncService/DeleteData',
	'SyncService/Stream',
	'WikiService/ListPages',
	'WikiService/GetPage',
	'WikiService/CreatePage',
	'WikiService/UpdatePage',
	'WikiService/DeletePage',
	'WikiService/ListPageActivity',
];
