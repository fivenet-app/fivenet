# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [resources/accounts/oauth2.proto](#resources_accounts_oauth2-proto)
    - [OAuth2Account](#resources-accounts-OAuth2Account)
    - [OAuth2Provider](#resources-accounts-OAuth2Provider)
  
- [resources/accounts/accounts.proto](#resources_accounts_accounts-proto)
    - [Account](#resources-accounts-Account)
    - [Character](#resources-accounts-Character)
  
- [resources/centrum/access.proto](#resources_centrum_access-proto)
    - [UnitAccess](#resources-centrum-UnitAccess)
    - [UnitJobAccess](#resources-centrum-UnitJobAccess)
    - [UnitQualificationAccess](#resources-centrum-UnitQualificationAccess)
    - [UnitUserAccess](#resources-centrum-UnitUserAccess)
  
    - [UnitAccessLevel](#resources-centrum-UnitAccessLevel)
  
- [resources/centrum/dispatches.proto](#resources_centrum_dispatches-proto)
    - [Dispatch](#resources-centrum-Dispatch)
    - [DispatchAssignment](#resources-centrum-DispatchAssignment)
    - [DispatchAssignments](#resources-centrum-DispatchAssignments)
    - [DispatchReference](#resources-centrum-DispatchReference)
    - [DispatchReferences](#resources-centrum-DispatchReferences)
    - [DispatchStatus](#resources-centrum-DispatchStatus)
  
    - [DispatchReferenceType](#resources-centrum-DispatchReferenceType)
    - [StatusDispatch](#resources-centrum-StatusDispatch)
    - [TakeDispatchResp](#resources-centrum-TakeDispatchResp)
  
- [resources/centrum/general.proto](#resources_centrum_general-proto)
    - [Attributes](#resources-centrum-Attributes)
    - [Disponents](#resources-centrum-Disponents)
    - [UserUnitMapping](#resources-centrum-UserUnitMapping)
  
- [resources/centrum/settings.proto](#resources_centrum_settings-proto)
    - [PredefinedStatus](#resources-centrum-PredefinedStatus)
    - [Settings](#resources-centrum-Settings)
    - [Timings](#resources-centrum-Timings)
  
    - [CentrumMode](#resources-centrum-CentrumMode)
  
- [resources/centrum/units.proto](#resources_centrum_units-proto)
    - [Unit](#resources-centrum-Unit)
    - [UnitAssignment](#resources-centrum-UnitAssignment)
    - [UnitAssignments](#resources-centrum-UnitAssignments)
    - [UnitStatus](#resources-centrum-UnitStatus)
  
    - [StatusUnit](#resources-centrum-StatusUnit)
  
- [resources/common/database/database.proto](#resources_common_database_database-proto)
    - [DateRange](#resources-common-database-DateRange)
    - [PaginationRequest](#resources-common-database-PaginationRequest)
    - [PaginationResponse](#resources-common-database-PaginationResponse)
    - [Sort](#resources-common-database-Sort)
  
- [resources/common/grpcws/grpcws.proto](#resources_common_grpcws_grpcws-proto)
    - [Body](#resources-common-grpcws-Body)
    - [Cancel](#resources-common-grpcws-Cancel)
    - [Complete](#resources-common-grpcws-Complete)
    - [Failure](#resources-common-grpcws-Failure)
    - [Failure.HeadersEntry](#resources-common-grpcws-Failure-HeadersEntry)
    - [GrpcFrame](#resources-common-grpcws-GrpcFrame)
    - [Header](#resources-common-grpcws-Header)
    - [Header.HeadersEntry](#resources-common-grpcws-Header-HeadersEntry)
    - [HeaderValue](#resources-common-grpcws-HeaderValue)
    - [Ping](#resources-common-grpcws-Ping)
  
- [resources/common/cron/cron.proto](#resources_common_cron_cron-proto)
    - [Cronjob](#resources-common-cron-Cronjob)
    - [CronjobCompletedEvent](#resources-common-cron-CronjobCompletedEvent)
    - [CronjobData](#resources-common-cron-CronjobData)
    - [CronjobLockOwnerState](#resources-common-cron-CronjobLockOwnerState)
    - [CronjobSchedulerEvent](#resources-common-cron-CronjobSchedulerEvent)
    - [GenericCronData](#resources-common-cron-GenericCronData)
    - [GenericCronData.AttributesEntry](#resources-common-cron-GenericCronData-AttributesEntry)
  
    - [CronjobState](#resources-common-cron-CronjobState)
  
- [resources/common/uuid.proto](#resources_common_uuid-proto)
    - [UUID](#resources-common-UUID)
  
- [resources/common/content/content.proto](#resources_common_content_content-proto)
    - [Content](#resources-common-content-Content)
    - [JSONNode](#resources-common-content-JSONNode)
    - [JSONNode.AttrsEntry](#resources-common-content-JSONNode-AttrsEntry)
  
    - [ContentType](#resources-common-content-ContentType)
  
- [resources/common/error.proto](#resources_common_error-proto)
    - [Error](#resources-common-Error)
  
- [resources/common/i18n.proto](#resources_common_i18n-proto)
    - [TranslateItem](#resources-common-TranslateItem)
    - [TranslateItem.ParametersEntry](#resources-common-TranslateItem-ParametersEntry)
  
- [resources/documents/category.proto](#resources_documents_category-proto)
    - [Category](#resources-documents-Category)
  
- [resources/documents/access.proto](#resources_documents_access-proto)
    - [DocumentAccess](#resources-documents-DocumentAccess)
    - [DocumentJobAccess](#resources-documents-DocumentJobAccess)
    - [DocumentUserAccess](#resources-documents-DocumentUserAccess)
  
    - [AccessLevel](#resources-documents-AccessLevel)
  
- [resources/documents/activity.proto](#resources_documents_activity-proto)
    - [DocAccessJobsDiff](#resources-documents-DocAccessJobsDiff)
    - [DocAccessRequested](#resources-documents-DocAccessRequested)
    - [DocAccessUpdated](#resources-documents-DocAccessUpdated)
    - [DocAccessUsersDiff](#resources-documents-DocAccessUsersDiff)
    - [DocActivity](#resources-documents-DocActivity)
    - [DocActivityData](#resources-documents-DocActivityData)
    - [DocOwnerChanged](#resources-documents-DocOwnerChanged)
    - [DocUpdated](#resources-documents-DocUpdated)
  
    - [DocActivityType](#resources-documents-DocActivityType)
  
- [resources/documents/comment.proto](#resources_documents_comment-proto)
    - [Comment](#resources-documents-Comment)
  
- [resources/documents/documents.proto](#resources_documents_documents-proto)
    - [Document](#resources-documents-Document)
    - [DocumentReference](#resources-documents-DocumentReference)
    - [DocumentRelation](#resources-documents-DocumentRelation)
    - [DocumentShort](#resources-documents-DocumentShort)
    - [WorkflowState](#resources-documents-WorkflowState)
    - [WorkflowUserState](#resources-documents-WorkflowUserState)
  
    - [DocReference](#resources-documents-DocReference)
    - [DocRelation](#resources-documents-DocRelation)
  
- [resources/documents/requests.proto](#resources_documents_requests-proto)
    - [DocRequest](#resources-documents-DocRequest)
  
- [resources/documents/templates.proto](#resources_documents_templates-proto)
    - [ObjectSpecs](#resources-documents-ObjectSpecs)
    - [Template](#resources-documents-Template)
    - [TemplateData](#resources-documents-TemplateData)
    - [TemplateJobAccess](#resources-documents-TemplateJobAccess)
    - [TemplateRequirements](#resources-documents-TemplateRequirements)
    - [TemplateSchema](#resources-documents-TemplateSchema)
    - [TemplateShort](#resources-documents-TemplateShort)
    - [TemplateUserAccess](#resources-documents-TemplateUserAccess)
  
- [resources/documents/workflow.proto](#resources_documents_workflow-proto)
    - [AutoCloseSettings](#resources-documents-AutoCloseSettings)
    - [Reminder](#resources-documents-Reminder)
    - [ReminderSettings](#resources-documents-ReminderSettings)
    - [Workflow](#resources-documents-Workflow)
    - [WorkflowCronData](#resources-documents-WorkflowCronData)
  
- [resources/filestore/file.proto](#resources_filestore_file-proto)
    - [File](#resources-filestore-File)
    - [FileInfo](#resources-filestore-FileInfo)
  
- [resources/jobs/conduct.proto](#resources_jobs_conduct-proto)
    - [ConductEntry](#resources-jobs-ConductEntry)
  
    - [ConductType](#resources-jobs-ConductType)
  
- [resources/jobs/colleagues.proto](#resources_jobs_colleagues-proto)
    - [Colleague](#resources-jobs-Colleague)
    - [ColleagueAbsenceDate](#resources-jobs-ColleagueAbsenceDate)
    - [ColleagueGradeChange](#resources-jobs-ColleagueGradeChange)
    - [ColleagueLabelsChange](#resources-jobs-ColleagueLabelsChange)
    - [ColleagueNameChange](#resources-jobs-ColleagueNameChange)
    - [JobsUserActivity](#resources-jobs-JobsUserActivity)
    - [JobsUserActivityData](#resources-jobs-JobsUserActivityData)
    - [JobsUserProps](#resources-jobs-JobsUserProps)
  
    - [JobsUserActivityType](#resources-jobs-JobsUserActivityType)
  
- [resources/jobs/labels.proto](#resources_jobs_labels-proto)
    - [Label](#resources-jobs-Label)
    - [LabelCount](#resources-jobs-LabelCount)
    - [Labels](#resources-jobs-Labels)
  
- [resources/jobs/timeclock.proto](#resources_jobs_timeclock-proto)
    - [TimeclockEntry](#resources-jobs-TimeclockEntry)
    - [TimeclockStats](#resources-jobs-TimeclockStats)
    - [TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats)
  
    - [TimeclockMode](#resources-jobs-TimeclockMode)
    - [TimeclockUserMode](#resources-jobs-TimeclockUserMode)
  
- [resources/laws/laws.proto](#resources_laws_laws-proto)
    - [Law](#resources-laws-Law)
    - [LawBook](#resources-laws-LawBook)
  
- [resources/livemap/tracker.proto](#resources_livemap_tracker-proto)
    - [UsersUpdateEvent](#resources-livemap-UsersUpdateEvent)
  
- [resources/livemap/livemap.proto](#resources_livemap_livemap-proto)
    - [CircleMarker](#resources-livemap-CircleMarker)
    - [Coords](#resources-livemap-Coords)
    - [IconMarker](#resources-livemap-IconMarker)
    - [MarkerData](#resources-livemap-MarkerData)
    - [MarkerInfo](#resources-livemap-MarkerInfo)
    - [MarkerMarker](#resources-livemap-MarkerMarker)
    - [UserMarker](#resources-livemap-UserMarker)
  
    - [MarkerType](#resources-livemap-MarkerType)
  
- [resources/notifications/events.proto](#resources_notifications_events-proto)
    - [JobEvent](#resources-notifications-JobEvent)
    - [JobGradeEvent](#resources-notifications-JobGradeEvent)
    - [SystemEvent](#resources-notifications-SystemEvent)
    - [UserEvent](#resources-notifications-UserEvent)
  
- [resources/notifications/notifications.proto](#resources_notifications_notifications-proto)
    - [CalendarData](#resources-notifications-CalendarData)
    - [Data](#resources-notifications-Data)
    - [Link](#resources-notifications-Link)
    - [Notification](#resources-notifications-Notification)
  
    - [NotificationCategory](#resources-notifications-NotificationCategory)
    - [NotificationType](#resources-notifications-NotificationType)
  
- [resources/permissions/permissions.proto](#resources_permissions_permissions-proto)
    - [AttributeValues](#resources-permissions-AttributeValues)
    - [JobGradeList](#resources-permissions-JobGradeList)
    - [JobGradeList.JobsEntry](#resources-permissions-JobGradeList-JobsEntry)
    - [Permission](#resources-permissions-Permission)
    - [RawRoleAttribute](#resources-permissions-RawRoleAttribute)
    - [Role](#resources-permissions-Role)
    - [RoleAttribute](#resources-permissions-RoleAttribute)
    - [StringList](#resources-permissions-StringList)
  
- [resources/qualifications/access.proto](#resources_qualifications_access-proto)
    - [QualificationAccess](#resources-qualifications-QualificationAccess)
    - [QualificationJobAccess](#resources-qualifications-QualificationJobAccess)
    - [QualificationUserAccess](#resources-qualifications-QualificationUserAccess)
  
    - [AccessLevel](#resources-qualifications-AccessLevel)
  
- [resources/qualifications/exam.proto](#resources_qualifications_exam-proto)
    - [ExamQuestion](#resources-qualifications-ExamQuestion)
    - [ExamQuestionAnswerData](#resources-qualifications-ExamQuestionAnswerData)
    - [ExamQuestionData](#resources-qualifications-ExamQuestionData)
    - [ExamQuestionImage](#resources-qualifications-ExamQuestionImage)
    - [ExamQuestionMultipleChoice](#resources-qualifications-ExamQuestionMultipleChoice)
    - [ExamQuestionSeparator](#resources-qualifications-ExamQuestionSeparator)
    - [ExamQuestionSingleChoice](#resources-qualifications-ExamQuestionSingleChoice)
    - [ExamQuestionText](#resources-qualifications-ExamQuestionText)
    - [ExamQuestionYesNo](#resources-qualifications-ExamQuestionYesNo)
    - [ExamQuestions](#resources-qualifications-ExamQuestions)
    - [ExamResponse](#resources-qualifications-ExamResponse)
    - [ExamResponseData](#resources-qualifications-ExamResponseData)
    - [ExamResponseMultipleChoice](#resources-qualifications-ExamResponseMultipleChoice)
    - [ExamResponseSeparator](#resources-qualifications-ExamResponseSeparator)
    - [ExamResponseSingleChoice](#resources-qualifications-ExamResponseSingleChoice)
    - [ExamResponseText](#resources-qualifications-ExamResponseText)
    - [ExamResponseYesNo](#resources-qualifications-ExamResponseYesNo)
    - [ExamResponses](#resources-qualifications-ExamResponses)
    - [ExamUser](#resources-qualifications-ExamUser)
  
- [resources/qualifications/qualifications.proto](#resources_qualifications_qualifications-proto)
    - [Qualification](#resources-qualifications-Qualification)
    - [QualificationDiscordSettings](#resources-qualifications-QualificationDiscordSettings)
    - [QualificationExamSettings](#resources-qualifications-QualificationExamSettings)
    - [QualificationRequest](#resources-qualifications-QualificationRequest)
    - [QualificationRequirement](#resources-qualifications-QualificationRequirement)
    - [QualificationResult](#resources-qualifications-QualificationResult)
    - [QualificationShort](#resources-qualifications-QualificationShort)
  
    - [QualificationExamMode](#resources-qualifications-QualificationExamMode)
    - [RequestStatus](#resources-qualifications-RequestStatus)
    - [ResultStatus](#resources-qualifications-ResultStatus)
  
- [resources/rector/audit.proto](#resources_rector_audit-proto)
    - [AuditEntry](#resources-rector-AuditEntry)
  
    - [EventType](#resources-rector-EventType)
  
- [resources/rector/config.proto](#resources_rector_config-proto)
    - [AppConfig](#resources-rector-AppConfig)
    - [Auth](#resources-rector-Auth)
    - [Discord](#resources-rector-Discord)
    - [DiscordBotPresence](#resources-rector-DiscordBotPresence)
    - [JobInfo](#resources-rector-JobInfo)
    - [Links](#resources-rector-Links)
    - [Perm](#resources-rector-Perm)
    - [Perms](#resources-rector-Perms)
    - [UnemployedJob](#resources-rector-UnemployedJob)
    - [UserTracker](#resources-rector-UserTracker)
    - [Website](#resources-rector-Website)
  
    - [DiscordBotPresenceType](#resources-rector-DiscordBotPresenceType)
  
- [resources/timestamp/timestamp.proto](#resources_timestamp_timestamp-proto)
    - [Timestamp](#resources-timestamp-Timestamp)
  
- [resources/users/job_props.proto](#resources_users_job_props-proto)
    - [DiscordSyncChange](#resources-users-DiscordSyncChange)
    - [DiscordSyncChanges](#resources-users-DiscordSyncChanges)
    - [DiscordSyncSettings](#resources-users-DiscordSyncSettings)
    - [GroupMapping](#resources-users-GroupMapping)
    - [GroupSyncSettings](#resources-users-GroupSyncSettings)
    - [JobProps](#resources-users-JobProps)
    - [JobSettings](#resources-users-JobSettings)
    - [JobsAbsenceSettings](#resources-users-JobsAbsenceSettings)
    - [QuickButtons](#resources-users-QuickButtons)
    - [StatusLogSettings](#resources-users-StatusLogSettings)
    - [UserInfoSyncSettings](#resources-users-UserInfoSyncSettings)
  
    - [UserInfoSyncUnemployedMode](#resources-users-UserInfoSyncUnemployedMode)
  
- [resources/users/jobs.proto](#resources_users_jobs-proto)
    - [Job](#resources-users-Job)
    - [JobGrade](#resources-users-JobGrade)
  
- [resources/users/users.proto](#resources_users_users-proto)
    - [CitizenAttribute](#resources-users-CitizenAttribute)
    - [CitizenAttributes](#resources-users-CitizenAttributes)
    - [License](#resources-users-License)
    - [User](#resources-users-User)
    - [UserActivity](#resources-users-UserActivity)
    - [UserProps](#resources-users-UserProps)
    - [UserShort](#resources-users-UserShort)
  
    - [UserActivityType](#resources-users-UserActivityType)
  
- [resources/vehicles/vehicles.proto](#resources_vehicles_vehicles-proto)
    - [Vehicle](#resources-vehicles-Vehicle)
  
- [resources/calendar/access.proto](#resources_calendar_access-proto)
    - [CalendarAccess](#resources-calendar-CalendarAccess)
    - [CalendarJobAccess](#resources-calendar-CalendarJobAccess)
    - [CalendarUserAccess](#resources-calendar-CalendarUserAccess)
  
    - [AccessLevel](#resources-calendar-AccessLevel)
  
- [resources/calendar/calendar.proto](#resources_calendar_calendar-proto)
    - [Calendar](#resources-calendar-Calendar)
    - [CalendarEntry](#resources-calendar-CalendarEntry)
    - [CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP)
    - [CalendarEntryRecurring](#resources-calendar-CalendarEntryRecurring)
    - [CalendarShort](#resources-calendar-CalendarShort)
    - [CalendarSub](#resources-calendar-CalendarSub)
  
    - [RsvpResponses](#resources-calendar-RsvpResponses)
  
- [resources/stats/stats.proto](#resources_stats_stats-proto)
    - [Stat](#resources-stats-Stat)
  
- [resources/internet/ads.proto](#resources_internet_ads-proto)
    - [Ad](#resources-internet-Ad)
  
    - [AdType](#resources-internet-AdType)
  
- [resources/internet/search.proto](#resources_internet_search-proto)
    - [SearchResult](#resources-internet-SearchResult)
  
- [resources/internet/domain.proto](#resources_internet_domain-proto)
    - [Domain](#resources-internet-Domain)
  
- [resources/internet/page.proto](#resources_internet_page-proto)
    - [Page](#resources-internet-Page)
    - [PageData](#resources-internet-PageData)
  
    - [PageLayoutType](#resources-internet-PageLayoutType)
  
- [resources/mailer/access.proto](#resources_mailer_access-proto)
    - [Access](#resources-mailer-Access)
    - [JobAccess](#resources-mailer-JobAccess)
    - [QualificationAccess](#resources-mailer-QualificationAccess)
    - [UserAccess](#resources-mailer-UserAccess)
  
    - [AccessLevel](#resources-mailer-AccessLevel)
  
- [resources/mailer/email.proto](#resources_mailer_email-proto)
    - [Email](#resources-mailer-Email)
  
- [resources/mailer/events.proto](#resources_mailer_events-proto)
    - [MailerEvent](#resources-mailer-MailerEvent)
  
- [resources/mailer/message.proto](#resources_mailer_message-proto)
    - [Message](#resources-mailer-Message)
    - [MessageData](#resources-mailer-MessageData)
    - [MessageDataEntry](#resources-mailer-MessageDataEntry)
  
- [resources/mailer/settings.proto](#resources_mailer_settings-proto)
    - [EmailSettings](#resources-mailer-EmailSettings)
  
- [resources/mailer/template.proto](#resources_mailer_template-proto)
    - [Template](#resources-mailer-Template)
  
- [resources/mailer/thread.proto](#resources_mailer_thread-proto)
    - [Thread](#resources-mailer-Thread)
    - [ThreadRecipientEmail](#resources-mailer-ThreadRecipientEmail)
    - [ThreadState](#resources-mailer-ThreadState)
  
- [resources/wiki/access.proto](#resources_wiki_access-proto)
    - [PageAccess](#resources-wiki-PageAccess)
    - [PageJobAccess](#resources-wiki-PageJobAccess)
    - [PageUserAccess](#resources-wiki-PageUserAccess)
  
    - [AccessLevel](#resources-wiki-AccessLevel)
  
- [resources/wiki/activity.proto](#resources_wiki_activity-proto)
    - [PageAccessJobsDiff](#resources-wiki-PageAccessJobsDiff)
    - [PageAccessUpdated](#resources-wiki-PageAccessUpdated)
    - [PageAccessUsersDiff](#resources-wiki-PageAccessUsersDiff)
    - [PageActivity](#resources-wiki-PageActivity)
    - [PageActivityData](#resources-wiki-PageActivityData)
    - [PageUpdated](#resources-wiki-PageUpdated)
  
    - [PageActivityType](#resources-wiki-PageActivityType)
  
- [resources/wiki/page.proto](#resources_wiki_page-proto)
    - [Page](#resources-wiki-Page)
    - [PageMeta](#resources-wiki-PageMeta)
    - [PageRootInfo](#resources-wiki-PageRootInfo)
    - [PageShort](#resources-wiki-PageShort)
  
- [services/auth/auth.proto](#services_auth_auth-proto)
    - [ChangePasswordRequest](#services-auth-ChangePasswordRequest)
    - [ChangePasswordResponse](#services-auth-ChangePasswordResponse)
    - [ChangeUsernameRequest](#services-auth-ChangeUsernameRequest)
    - [ChangeUsernameResponse](#services-auth-ChangeUsernameResponse)
    - [ChooseCharacterRequest](#services-auth-ChooseCharacterRequest)
    - [ChooseCharacterResponse](#services-auth-ChooseCharacterResponse)
    - [CreateAccountRequest](#services-auth-CreateAccountRequest)
    - [CreateAccountResponse](#services-auth-CreateAccountResponse)
    - [DeleteOAuth2ConnectionRequest](#services-auth-DeleteOAuth2ConnectionRequest)
    - [DeleteOAuth2ConnectionResponse](#services-auth-DeleteOAuth2ConnectionResponse)
    - [ForgotPasswordRequest](#services-auth-ForgotPasswordRequest)
    - [ForgotPasswordResponse](#services-auth-ForgotPasswordResponse)
    - [GetAccountInfoRequest](#services-auth-GetAccountInfoRequest)
    - [GetAccountInfoResponse](#services-auth-GetAccountInfoResponse)
    - [GetCharactersRequest](#services-auth-GetCharactersRequest)
    - [GetCharactersResponse](#services-auth-GetCharactersResponse)
    - [LoginRequest](#services-auth-LoginRequest)
    - [LoginResponse](#services-auth-LoginResponse)
    - [LogoutRequest](#services-auth-LogoutRequest)
    - [LogoutResponse](#services-auth-LogoutResponse)
    - [SetSuperUserModeRequest](#services-auth-SetSuperUserModeRequest)
    - [SetSuperUserModeResponse](#services-auth-SetSuperUserModeResponse)
  
    - [AuthService](#services-auth-AuthService)
  
- [services/centrum/centrum.proto](#services_centrum_centrum-proto)
    - [AssignDispatchRequest](#services-centrum-AssignDispatchRequest)
    - [AssignDispatchResponse](#services-centrum-AssignDispatchResponse)
    - [AssignUnitRequest](#services-centrum-AssignUnitRequest)
    - [AssignUnitResponse](#services-centrum-AssignUnitResponse)
    - [CreateDispatchRequest](#services-centrum-CreateDispatchRequest)
    - [CreateDispatchResponse](#services-centrum-CreateDispatchResponse)
    - [CreateOrUpdateUnitRequest](#services-centrum-CreateOrUpdateUnitRequest)
    - [CreateOrUpdateUnitResponse](#services-centrum-CreateOrUpdateUnitResponse)
    - [DeleteDispatchRequest](#services-centrum-DeleteDispatchRequest)
    - [DeleteDispatchResponse](#services-centrum-DeleteDispatchResponse)
    - [DeleteUnitRequest](#services-centrum-DeleteUnitRequest)
    - [DeleteUnitResponse](#services-centrum-DeleteUnitResponse)
    - [GetDispatchRequest](#services-centrum-GetDispatchRequest)
    - [GetDispatchResponse](#services-centrum-GetDispatchResponse)
    - [GetSettingsRequest](#services-centrum-GetSettingsRequest)
    - [GetSettingsResponse](#services-centrum-GetSettingsResponse)
    - [JoinUnitRequest](#services-centrum-JoinUnitRequest)
    - [JoinUnitResponse](#services-centrum-JoinUnitResponse)
    - [LatestState](#services-centrum-LatestState)
    - [ListDispatchActivityRequest](#services-centrum-ListDispatchActivityRequest)
    - [ListDispatchActivityResponse](#services-centrum-ListDispatchActivityResponse)
    - [ListDispatchesRequest](#services-centrum-ListDispatchesRequest)
    - [ListDispatchesResponse](#services-centrum-ListDispatchesResponse)
    - [ListUnitActivityRequest](#services-centrum-ListUnitActivityRequest)
    - [ListUnitActivityResponse](#services-centrum-ListUnitActivityResponse)
    - [ListUnitsRequest](#services-centrum-ListUnitsRequest)
    - [ListUnitsResponse](#services-centrum-ListUnitsResponse)
    - [StreamRequest](#services-centrum-StreamRequest)
    - [StreamResponse](#services-centrum-StreamResponse)
    - [TakeControlRequest](#services-centrum-TakeControlRequest)
    - [TakeControlResponse](#services-centrum-TakeControlResponse)
    - [TakeDispatchRequest](#services-centrum-TakeDispatchRequest)
    - [TakeDispatchResponse](#services-centrum-TakeDispatchResponse)
    - [UpdateDispatchRequest](#services-centrum-UpdateDispatchRequest)
    - [UpdateDispatchResponse](#services-centrum-UpdateDispatchResponse)
    - [UpdateDispatchStatusRequest](#services-centrum-UpdateDispatchStatusRequest)
    - [UpdateDispatchStatusResponse](#services-centrum-UpdateDispatchStatusResponse)
    - [UpdateSettingsRequest](#services-centrum-UpdateSettingsRequest)
    - [UpdateSettingsResponse](#services-centrum-UpdateSettingsResponse)
    - [UpdateUnitStatusRequest](#services-centrum-UpdateUnitStatusRequest)
    - [UpdateUnitStatusResponse](#services-centrum-UpdateUnitStatusResponse)
  
    - [CentrumService](#services-centrum-CentrumService)
  
- [services/citizenstore/citizenstore.proto](#services_citizenstore_citizenstore-proto)
    - [GetUserRequest](#services-citizenstore-GetUserRequest)
    - [GetUserResponse](#services-citizenstore-GetUserResponse)
    - [ListCitizensRequest](#services-citizenstore-ListCitizensRequest)
    - [ListCitizensResponse](#services-citizenstore-ListCitizensResponse)
    - [ListUserActivityRequest](#services-citizenstore-ListUserActivityRequest)
    - [ListUserActivityResponse](#services-citizenstore-ListUserActivityResponse)
    - [ManageCitizenAttributesRequest](#services-citizenstore-ManageCitizenAttributesRequest)
    - [ManageCitizenAttributesResponse](#services-citizenstore-ManageCitizenAttributesResponse)
    - [SetProfilePictureRequest](#services-citizenstore-SetProfilePictureRequest)
    - [SetProfilePictureResponse](#services-citizenstore-SetProfilePictureResponse)
    - [SetUserPropsRequest](#services-citizenstore-SetUserPropsRequest)
    - [SetUserPropsResponse](#services-citizenstore-SetUserPropsResponse)
  
    - [CitizenStoreService](#services-citizenstore-CitizenStoreService)
  
- [services/completor/completor.proto](#services_completor_completor-proto)
    - [CompleteCitizenAttributesRequest](#services-completor-CompleteCitizenAttributesRequest)
    - [CompleteCitizenAttributesResponse](#services-completor-CompleteCitizenAttributesResponse)
    - [CompleteCitizensRequest](#services-completor-CompleteCitizensRequest)
    - [CompleteCitizensRespoonse](#services-completor-CompleteCitizensRespoonse)
    - [CompleteDocumentCategoriesRequest](#services-completor-CompleteDocumentCategoriesRequest)
    - [CompleteDocumentCategoriesResponse](#services-completor-CompleteDocumentCategoriesResponse)
    - [CompleteJobsRequest](#services-completor-CompleteJobsRequest)
    - [CompleteJobsResponse](#services-completor-CompleteJobsResponse)
    - [ListLawBooksRequest](#services-completor-ListLawBooksRequest)
    - [ListLawBooksResponse](#services-completor-ListLawBooksResponse)
  
    - [CompletorService](#services-completor-CompletorService)
  
- [services/dmv/vehicles.proto](#services_dmv_vehicles-proto)
    - [ListVehiclesRequest](#services-dmv-ListVehiclesRequest)
    - [ListVehiclesResponse](#services-dmv-ListVehiclesResponse)
  
    - [DMVService](#services-dmv-DMVService)
  
- [services/docstore/docstore.proto](#services_docstore_docstore-proto)
    - [AddDocumentReferenceRequest](#services-docstore-AddDocumentReferenceRequest)
    - [AddDocumentReferenceResponse](#services-docstore-AddDocumentReferenceResponse)
    - [AddDocumentRelationRequest](#services-docstore-AddDocumentRelationRequest)
    - [AddDocumentRelationResponse](#services-docstore-AddDocumentRelationResponse)
    - [ChangeDocumentOwnerRequest](#services-docstore-ChangeDocumentOwnerRequest)
    - [ChangeDocumentOwnerResponse](#services-docstore-ChangeDocumentOwnerResponse)
    - [CreateCategoryRequest](#services-docstore-CreateCategoryRequest)
    - [CreateCategoryResponse](#services-docstore-CreateCategoryResponse)
    - [CreateDocumentReqRequest](#services-docstore-CreateDocumentReqRequest)
    - [CreateDocumentReqResponse](#services-docstore-CreateDocumentReqResponse)
    - [CreateDocumentRequest](#services-docstore-CreateDocumentRequest)
    - [CreateDocumentResponse](#services-docstore-CreateDocumentResponse)
    - [CreateTemplateRequest](#services-docstore-CreateTemplateRequest)
    - [CreateTemplateResponse](#services-docstore-CreateTemplateResponse)
    - [DeleteCategoryRequest](#services-docstore-DeleteCategoryRequest)
    - [DeleteCategoryResponse](#services-docstore-DeleteCategoryResponse)
    - [DeleteCommentRequest](#services-docstore-DeleteCommentRequest)
    - [DeleteCommentResponse](#services-docstore-DeleteCommentResponse)
    - [DeleteDocumentReqRequest](#services-docstore-DeleteDocumentReqRequest)
    - [DeleteDocumentReqResponse](#services-docstore-DeleteDocumentReqResponse)
    - [DeleteDocumentRequest](#services-docstore-DeleteDocumentRequest)
    - [DeleteDocumentResponse](#services-docstore-DeleteDocumentResponse)
    - [DeleteTemplateRequest](#services-docstore-DeleteTemplateRequest)
    - [DeleteTemplateResponse](#services-docstore-DeleteTemplateResponse)
    - [EditCommentRequest](#services-docstore-EditCommentRequest)
    - [EditCommentResponse](#services-docstore-EditCommentResponse)
    - [GetCommentsRequest](#services-docstore-GetCommentsRequest)
    - [GetCommentsResponse](#services-docstore-GetCommentsResponse)
    - [GetDocumentAccessRequest](#services-docstore-GetDocumentAccessRequest)
    - [GetDocumentAccessResponse](#services-docstore-GetDocumentAccessResponse)
    - [GetDocumentReferencesRequest](#services-docstore-GetDocumentReferencesRequest)
    - [GetDocumentReferencesResponse](#services-docstore-GetDocumentReferencesResponse)
    - [GetDocumentRelationsRequest](#services-docstore-GetDocumentRelationsRequest)
    - [GetDocumentRelationsResponse](#services-docstore-GetDocumentRelationsResponse)
    - [GetDocumentRequest](#services-docstore-GetDocumentRequest)
    - [GetDocumentResponse](#services-docstore-GetDocumentResponse)
    - [GetTemplateRequest](#services-docstore-GetTemplateRequest)
    - [GetTemplateResponse](#services-docstore-GetTemplateResponse)
    - [ListCategoriesRequest](#services-docstore-ListCategoriesRequest)
    - [ListCategoriesResponse](#services-docstore-ListCategoriesResponse)
    - [ListDocumentActivityRequest](#services-docstore-ListDocumentActivityRequest)
    - [ListDocumentActivityResponse](#services-docstore-ListDocumentActivityResponse)
    - [ListDocumentPinsRequest](#services-docstore-ListDocumentPinsRequest)
    - [ListDocumentPinsResponse](#services-docstore-ListDocumentPinsResponse)
    - [ListDocumentReqsRequest](#services-docstore-ListDocumentReqsRequest)
    - [ListDocumentReqsResponse](#services-docstore-ListDocumentReqsResponse)
    - [ListDocumentsRequest](#services-docstore-ListDocumentsRequest)
    - [ListDocumentsResponse](#services-docstore-ListDocumentsResponse)
    - [ListTemplatesRequest](#services-docstore-ListTemplatesRequest)
    - [ListTemplatesResponse](#services-docstore-ListTemplatesResponse)
    - [ListUserDocumentsRequest](#services-docstore-ListUserDocumentsRequest)
    - [ListUserDocumentsResponse](#services-docstore-ListUserDocumentsResponse)
    - [PostCommentRequest](#services-docstore-PostCommentRequest)
    - [PostCommentResponse](#services-docstore-PostCommentResponse)
    - [RemoveDocumentReferenceRequest](#services-docstore-RemoveDocumentReferenceRequest)
    - [RemoveDocumentReferenceResponse](#services-docstore-RemoveDocumentReferenceResponse)
    - [RemoveDocumentRelationRequest](#services-docstore-RemoveDocumentRelationRequest)
    - [RemoveDocumentRelationResponse](#services-docstore-RemoveDocumentRelationResponse)
    - [SetDocumentAccessRequest](#services-docstore-SetDocumentAccessRequest)
    - [SetDocumentAccessResponse](#services-docstore-SetDocumentAccessResponse)
    - [SetDocumentReminderRequest](#services-docstore-SetDocumentReminderRequest)
    - [SetDocumentReminderResponse](#services-docstore-SetDocumentReminderResponse)
    - [ToggleDocumentPinRequest](#services-docstore-ToggleDocumentPinRequest)
    - [ToggleDocumentPinResponse](#services-docstore-ToggleDocumentPinResponse)
    - [ToggleDocumentRequest](#services-docstore-ToggleDocumentRequest)
    - [ToggleDocumentResponse](#services-docstore-ToggleDocumentResponse)
    - [UpdateCategoryRequest](#services-docstore-UpdateCategoryRequest)
    - [UpdateCategoryResponse](#services-docstore-UpdateCategoryResponse)
    - [UpdateDocumentReqRequest](#services-docstore-UpdateDocumentReqRequest)
    - [UpdateDocumentReqResponse](#services-docstore-UpdateDocumentReqResponse)
    - [UpdateDocumentRequest](#services-docstore-UpdateDocumentRequest)
    - [UpdateDocumentResponse](#services-docstore-UpdateDocumentResponse)
    - [UpdateTemplateRequest](#services-docstore-UpdateTemplateRequest)
    - [UpdateTemplateResponse](#services-docstore-UpdateTemplateResponse)
  
    - [DocStoreService](#services-docstore-DocStoreService)
  
- [services/jobs/conduct.proto](#services_jobs_conduct-proto)
    - [CreateConductEntryRequest](#services-jobs-CreateConductEntryRequest)
    - [CreateConductEntryResponse](#services-jobs-CreateConductEntryResponse)
    - [DeleteConductEntryRequest](#services-jobs-DeleteConductEntryRequest)
    - [DeleteConductEntryResponse](#services-jobs-DeleteConductEntryResponse)
    - [ListConductEntriesRequest](#services-jobs-ListConductEntriesRequest)
    - [ListConductEntriesResponse](#services-jobs-ListConductEntriesResponse)
    - [UpdateConductEntryRequest](#services-jobs-UpdateConductEntryRequest)
    - [UpdateConductEntryResponse](#services-jobs-UpdateConductEntryResponse)
  
    - [JobsConductService](#services-jobs-JobsConductService)
  
- [services/jobs/jobs.proto](#services_jobs_jobs-proto)
    - [GetColleagueLabelsRequest](#services-jobs-GetColleagueLabelsRequest)
    - [GetColleagueLabelsResponse](#services-jobs-GetColleagueLabelsResponse)
    - [GetColleagueLabelsStatsRequest](#services-jobs-GetColleagueLabelsStatsRequest)
    - [GetColleagueLabelsStatsResponse](#services-jobs-GetColleagueLabelsStatsResponse)
    - [GetColleagueRequest](#services-jobs-GetColleagueRequest)
    - [GetColleagueResponse](#services-jobs-GetColleagueResponse)
    - [GetMOTDRequest](#services-jobs-GetMOTDRequest)
    - [GetMOTDResponse](#services-jobs-GetMOTDResponse)
    - [GetSelfRequest](#services-jobs-GetSelfRequest)
    - [GetSelfResponse](#services-jobs-GetSelfResponse)
    - [ListColleagueActivityRequest](#services-jobs-ListColleagueActivityRequest)
    - [ListColleagueActivityResponse](#services-jobs-ListColleagueActivityResponse)
    - [ListColleaguesRequest](#services-jobs-ListColleaguesRequest)
    - [ListColleaguesResponse](#services-jobs-ListColleaguesResponse)
    - [ManageColleagueLabelsRequest](#services-jobs-ManageColleagueLabelsRequest)
    - [ManageColleagueLabelsResponse](#services-jobs-ManageColleagueLabelsResponse)
    - [SetJobsUserPropsRequest](#services-jobs-SetJobsUserPropsRequest)
    - [SetJobsUserPropsResponse](#services-jobs-SetJobsUserPropsResponse)
    - [SetMOTDRequest](#services-jobs-SetMOTDRequest)
    - [SetMOTDResponse](#services-jobs-SetMOTDResponse)
  
    - [JobsService](#services-jobs-JobsService)
  
- [services/jobs/timeclock.proto](#services_jobs_timeclock-proto)
    - [GetTimeclockStatsRequest](#services-jobs-GetTimeclockStatsRequest)
    - [GetTimeclockStatsResponse](#services-jobs-GetTimeclockStatsResponse)
    - [ListInactiveEmployeesRequest](#services-jobs-ListInactiveEmployeesRequest)
    - [ListInactiveEmployeesResponse](#services-jobs-ListInactiveEmployeesResponse)
    - [ListTimeclockRequest](#services-jobs-ListTimeclockRequest)
    - [ListTimeclockResponse](#services-jobs-ListTimeclockResponse)
    - [TimeclockDay](#services-jobs-TimeclockDay)
    - [TimeclockRange](#services-jobs-TimeclockRange)
    - [TimeclockWeekly](#services-jobs-TimeclockWeekly)
  
    - [JobsTimeclockService](#services-jobs-JobsTimeclockService)
  
- [services/livemapper/livemap.proto](#services_livemapper_livemap-proto)
    - [CreateOrUpdateMarkerRequest](#services-livemapper-CreateOrUpdateMarkerRequest)
    - [CreateOrUpdateMarkerResponse](#services-livemapper-CreateOrUpdateMarkerResponse)
    - [DeleteMarkerRequest](#services-livemapper-DeleteMarkerRequest)
    - [DeleteMarkerResponse](#services-livemapper-DeleteMarkerResponse)
    - [JobsList](#services-livemapper-JobsList)
    - [MarkerMarkersUpdates](#services-livemapper-MarkerMarkersUpdates)
    - [StreamRequest](#services-livemapper-StreamRequest)
    - [StreamResponse](#services-livemapper-StreamResponse)
    - [UserMarkersUpdates](#services-livemapper-UserMarkersUpdates)
  
    - [LivemapperService](#services-livemapper-LivemapperService)
  
- [services/notificator/notificator.proto](#services_notificator_notificator-proto)
    - [GetNotificationsRequest](#services-notificator-GetNotificationsRequest)
    - [GetNotificationsResponse](#services-notificator-GetNotificationsResponse)
    - [MarkNotificationsRequest](#services-notificator-MarkNotificationsRequest)
    - [MarkNotificationsResponse](#services-notificator-MarkNotificationsResponse)
    - [StreamRequest](#services-notificator-StreamRequest)
    - [StreamResponse](#services-notificator-StreamResponse)
  
    - [NotificatorService](#services-notificator-NotificatorService)
  
- [services/qualifications/qualifications.proto](#services_qualifications_qualifications-proto)
    - [CreateOrUpdateQualificationRequestRequest](#services-qualifications-CreateOrUpdateQualificationRequestRequest)
    - [CreateOrUpdateQualificationRequestResponse](#services-qualifications-CreateOrUpdateQualificationRequestResponse)
    - [CreateOrUpdateQualificationResultRequest](#services-qualifications-CreateOrUpdateQualificationResultRequest)
    - [CreateOrUpdateQualificationResultResponse](#services-qualifications-CreateOrUpdateQualificationResultResponse)
    - [CreateQualificationRequest](#services-qualifications-CreateQualificationRequest)
    - [CreateQualificationResponse](#services-qualifications-CreateQualificationResponse)
    - [DeleteQualificationReqRequest](#services-qualifications-DeleteQualificationReqRequest)
    - [DeleteQualificationReqResponse](#services-qualifications-DeleteQualificationReqResponse)
    - [DeleteQualificationRequest](#services-qualifications-DeleteQualificationRequest)
    - [DeleteQualificationResponse](#services-qualifications-DeleteQualificationResponse)
    - [DeleteQualificationResultRequest](#services-qualifications-DeleteQualificationResultRequest)
    - [DeleteQualificationResultResponse](#services-qualifications-DeleteQualificationResultResponse)
    - [GetExamInfoRequest](#services-qualifications-GetExamInfoRequest)
    - [GetExamInfoResponse](#services-qualifications-GetExamInfoResponse)
    - [GetQualificationAccessRequest](#services-qualifications-GetQualificationAccessRequest)
    - [GetQualificationAccessResponse](#services-qualifications-GetQualificationAccessResponse)
    - [GetQualificationRequest](#services-qualifications-GetQualificationRequest)
    - [GetQualificationResponse](#services-qualifications-GetQualificationResponse)
    - [GetUserExamRequest](#services-qualifications-GetUserExamRequest)
    - [GetUserExamResponse](#services-qualifications-GetUserExamResponse)
    - [ListQualificationRequestsRequest](#services-qualifications-ListQualificationRequestsRequest)
    - [ListQualificationRequestsResponse](#services-qualifications-ListQualificationRequestsResponse)
    - [ListQualificationsRequest](#services-qualifications-ListQualificationsRequest)
    - [ListQualificationsResponse](#services-qualifications-ListQualificationsResponse)
    - [ListQualificationsResultsRequest](#services-qualifications-ListQualificationsResultsRequest)
    - [ListQualificationsResultsResponse](#services-qualifications-ListQualificationsResultsResponse)
    - [SetQualificationAccessRequest](#services-qualifications-SetQualificationAccessRequest)
    - [SetQualificationAccessResponse](#services-qualifications-SetQualificationAccessResponse)
    - [SubmitExamRequest](#services-qualifications-SubmitExamRequest)
    - [SubmitExamResponse](#services-qualifications-SubmitExamResponse)
    - [TakeExamRequest](#services-qualifications-TakeExamRequest)
    - [TakeExamResponse](#services-qualifications-TakeExamResponse)
    - [UpdateQualificationRequest](#services-qualifications-UpdateQualificationRequest)
    - [UpdateQualificationResponse](#services-qualifications-UpdateQualificationResponse)
  
    - [QualificationsService](#services-qualifications-QualificationsService)
  
- [services/rector/config.proto](#services_rector_config-proto)
    - [GetAppConfigRequest](#services-rector-GetAppConfigRequest)
    - [GetAppConfigResponse](#services-rector-GetAppConfigResponse)
    - [UpdateAppConfigRequest](#services-rector-UpdateAppConfigRequest)
    - [UpdateAppConfigResponse](#services-rector-UpdateAppConfigResponse)
  
    - [RectorConfigService](#services-rector-RectorConfigService)
  
- [services/rector/filestore.proto](#services_rector_filestore-proto)
    - [DeleteFileRequest](#services-rector-DeleteFileRequest)
    - [DeleteFileResponse](#services-rector-DeleteFileResponse)
    - [ListFilesRequest](#services-rector-ListFilesRequest)
    - [ListFilesResponse](#services-rector-ListFilesResponse)
    - [UploadFileRequest](#services-rector-UploadFileRequest)
    - [UploadFileResponse](#services-rector-UploadFileResponse)
  
    - [RectorFilestoreService](#services-rector-RectorFilestoreService)
  
- [services/rector/laws.proto](#services_rector_laws-proto)
    - [CreateOrUpdateLawBookRequest](#services-rector-CreateOrUpdateLawBookRequest)
    - [CreateOrUpdateLawBookResponse](#services-rector-CreateOrUpdateLawBookResponse)
    - [CreateOrUpdateLawRequest](#services-rector-CreateOrUpdateLawRequest)
    - [CreateOrUpdateLawResponse](#services-rector-CreateOrUpdateLawResponse)
    - [DeleteLawBookRequest](#services-rector-DeleteLawBookRequest)
    - [DeleteLawBookResponse](#services-rector-DeleteLawBookResponse)
    - [DeleteLawRequest](#services-rector-DeleteLawRequest)
    - [DeleteLawResponse](#services-rector-DeleteLawResponse)
  
    - [RectorLawsService](#services-rector-RectorLawsService)
  
- [services/rector/rector.proto](#services_rector_rector-proto)
    - [AttrsUpdate](#services-rector-AttrsUpdate)
    - [CreateRoleRequest](#services-rector-CreateRoleRequest)
    - [CreateRoleResponse](#services-rector-CreateRoleResponse)
    - [DeleteFactionRequest](#services-rector-DeleteFactionRequest)
    - [DeleteFactionResponse](#services-rector-DeleteFactionResponse)
    - [DeleteRoleRequest](#services-rector-DeleteRoleRequest)
    - [DeleteRoleResponse](#services-rector-DeleteRoleResponse)
    - [GetJobPropsRequest](#services-rector-GetJobPropsRequest)
    - [GetJobPropsResponse](#services-rector-GetJobPropsResponse)
    - [GetPermissionsRequest](#services-rector-GetPermissionsRequest)
    - [GetPermissionsResponse](#services-rector-GetPermissionsResponse)
    - [GetRoleRequest](#services-rector-GetRoleRequest)
    - [GetRoleResponse](#services-rector-GetRoleResponse)
    - [GetRolesRequest](#services-rector-GetRolesRequest)
    - [GetRolesResponse](#services-rector-GetRolesResponse)
    - [PermItem](#services-rector-PermItem)
    - [PermsUpdate](#services-rector-PermsUpdate)
    - [SetJobPropsRequest](#services-rector-SetJobPropsRequest)
    - [SetJobPropsResponse](#services-rector-SetJobPropsResponse)
    - [UpdateRoleLimitsRequest](#services-rector-UpdateRoleLimitsRequest)
    - [UpdateRoleLimitsResponse](#services-rector-UpdateRoleLimitsResponse)
    - [UpdateRolePermsRequest](#services-rector-UpdateRolePermsRequest)
    - [UpdateRolePermsResponse](#services-rector-UpdateRolePermsResponse)
    - [ViewAuditLogRequest](#services-rector-ViewAuditLogRequest)
    - [ViewAuditLogResponse](#services-rector-ViewAuditLogResponse)
  
    - [RectorService](#services-rector-RectorService)
  
- [services/calendar/calendar.proto](#services_calendar_calendar-proto)
    - [CreateOrUpdateCalendarEntryRequest](#services-calendar-CreateOrUpdateCalendarEntryRequest)
    - [CreateOrUpdateCalendarEntryResponse](#services-calendar-CreateOrUpdateCalendarEntryResponse)
    - [CreateOrUpdateCalendarRequest](#services-calendar-CreateOrUpdateCalendarRequest)
    - [CreateOrUpdateCalendarResponse](#services-calendar-CreateOrUpdateCalendarResponse)
    - [DeleteCalendarEntryRequest](#services-calendar-DeleteCalendarEntryRequest)
    - [DeleteCalendarEntryResponse](#services-calendar-DeleteCalendarEntryResponse)
    - [DeleteCalendarRequest](#services-calendar-DeleteCalendarRequest)
    - [DeleteCalendarResponse](#services-calendar-DeleteCalendarResponse)
    - [GetCalendarEntryRequest](#services-calendar-GetCalendarEntryRequest)
    - [GetCalendarEntryResponse](#services-calendar-GetCalendarEntryResponse)
    - [GetCalendarRequest](#services-calendar-GetCalendarRequest)
    - [GetCalendarResponse](#services-calendar-GetCalendarResponse)
    - [GetUpcomingEntriesRequest](#services-calendar-GetUpcomingEntriesRequest)
    - [GetUpcomingEntriesResponse](#services-calendar-GetUpcomingEntriesResponse)
    - [ListCalendarEntriesRequest](#services-calendar-ListCalendarEntriesRequest)
    - [ListCalendarEntriesResponse](#services-calendar-ListCalendarEntriesResponse)
    - [ListCalendarEntryRSVPRequest](#services-calendar-ListCalendarEntryRSVPRequest)
    - [ListCalendarEntryRSVPResponse](#services-calendar-ListCalendarEntryRSVPResponse)
    - [ListCalendarsRequest](#services-calendar-ListCalendarsRequest)
    - [ListCalendarsResponse](#services-calendar-ListCalendarsResponse)
    - [ListSubscriptionsRequest](#services-calendar-ListSubscriptionsRequest)
    - [ListSubscriptionsResponse](#services-calendar-ListSubscriptionsResponse)
    - [RSVPCalendarEntryRequest](#services-calendar-RSVPCalendarEntryRequest)
    - [RSVPCalendarEntryResponse](#services-calendar-RSVPCalendarEntryResponse)
    - [ShareCalendarEntryRequest](#services-calendar-ShareCalendarEntryRequest)
    - [ShareCalendarEntryResponse](#services-calendar-ShareCalendarEntryResponse)
    - [SubscribeToCalendarRequest](#services-calendar-SubscribeToCalendarRequest)
    - [SubscribeToCalendarResponse](#services-calendar-SubscribeToCalendarResponse)
  
    - [CalendarService](#services-calendar-CalendarService)
  
- [services/stats/stats.proto](#services_stats_stats-proto)
    - [GetStatsRequest](#services-stats-GetStatsRequest)
    - [GetStatsResponse](#services-stats-GetStatsResponse)
    - [GetStatsResponse.StatsEntry](#services-stats-GetStatsResponse-StatsEntry)
  
    - [StatsService](#services-stats-StatsService)
  
- [services/internet/ads.proto](#services_internet_ads-proto)
    - [GetAdsRequest](#services-internet-GetAdsRequest)
    - [GetAdsResponse](#services-internet-GetAdsResponse)
  
    - [AdsService](#services-internet-AdsService)
  
- [services/internet/internet.proto](#services_internet_internet-proto)
    - [GetPageRequest](#services-internet-GetPageRequest)
    - [GetPageResponse](#services-internet-GetPageResponse)
    - [SearchRequest](#services-internet-SearchRequest)
    - [SearchResponse](#services-internet-SearchResponse)
  
    - [InternetService](#services-internet-InternetService)
  
- [services/mailer/mailer.proto](#services_mailer_mailer-proto)
    - [CreateOrUpdateEmailRequest](#services-mailer-CreateOrUpdateEmailRequest)
    - [CreateOrUpdateEmailResponse](#services-mailer-CreateOrUpdateEmailResponse)
    - [CreateOrUpdateTemplateRequest](#services-mailer-CreateOrUpdateTemplateRequest)
    - [CreateOrUpdateTemplateResponse](#services-mailer-CreateOrUpdateTemplateResponse)
    - [CreateThreadRequest](#services-mailer-CreateThreadRequest)
    - [CreateThreadResponse](#services-mailer-CreateThreadResponse)
    - [DeleteEmailRequest](#services-mailer-DeleteEmailRequest)
    - [DeleteEmailResponse](#services-mailer-DeleteEmailResponse)
    - [DeleteMessageRequest](#services-mailer-DeleteMessageRequest)
    - [DeleteMessageResponse](#services-mailer-DeleteMessageResponse)
    - [DeleteTemplateRequest](#services-mailer-DeleteTemplateRequest)
    - [DeleteTemplateResponse](#services-mailer-DeleteTemplateResponse)
    - [DeleteThreadRequest](#services-mailer-DeleteThreadRequest)
    - [DeleteThreadResponse](#services-mailer-DeleteThreadResponse)
    - [GetEmailProposalsRequest](#services-mailer-GetEmailProposalsRequest)
    - [GetEmailProposalsResponse](#services-mailer-GetEmailProposalsResponse)
    - [GetEmailRequest](#services-mailer-GetEmailRequest)
    - [GetEmailResponse](#services-mailer-GetEmailResponse)
    - [GetEmailSettingsRequest](#services-mailer-GetEmailSettingsRequest)
    - [GetEmailSettingsResponse](#services-mailer-GetEmailSettingsResponse)
    - [GetTemplateRequest](#services-mailer-GetTemplateRequest)
    - [GetTemplateResponse](#services-mailer-GetTemplateResponse)
    - [GetThreadRequest](#services-mailer-GetThreadRequest)
    - [GetThreadResponse](#services-mailer-GetThreadResponse)
    - [GetThreadStateRequest](#services-mailer-GetThreadStateRequest)
    - [GetThreadStateResponse](#services-mailer-GetThreadStateResponse)
    - [ListEmailsRequest](#services-mailer-ListEmailsRequest)
    - [ListEmailsResponse](#services-mailer-ListEmailsResponse)
    - [ListTemplatesRequest](#services-mailer-ListTemplatesRequest)
    - [ListTemplatesResponse](#services-mailer-ListTemplatesResponse)
    - [ListThreadMessagesRequest](#services-mailer-ListThreadMessagesRequest)
    - [ListThreadMessagesResponse](#services-mailer-ListThreadMessagesResponse)
    - [ListThreadsRequest](#services-mailer-ListThreadsRequest)
    - [ListThreadsResponse](#services-mailer-ListThreadsResponse)
    - [PostMessageRequest](#services-mailer-PostMessageRequest)
    - [PostMessageResponse](#services-mailer-PostMessageResponse)
    - [SearchThreadsRequest](#services-mailer-SearchThreadsRequest)
    - [SearchThreadsResponse](#services-mailer-SearchThreadsResponse)
    - [SetEmailSettingsRequest](#services-mailer-SetEmailSettingsRequest)
    - [SetEmailSettingsResponse](#services-mailer-SetEmailSettingsResponse)
    - [SetThreadStateRequest](#services-mailer-SetThreadStateRequest)
    - [SetThreadStateResponse](#services-mailer-SetThreadStateResponse)
  
    - [MailerService](#services-mailer-MailerService)
  
- [services/wiki/wiki.proto](#services_wiki_wiki-proto)
    - [CreatePageRequest](#services-wiki-CreatePageRequest)
    - [CreatePageResponse](#services-wiki-CreatePageResponse)
    - [DeletePageRequest](#services-wiki-DeletePageRequest)
    - [DeletePageResponse](#services-wiki-DeletePageResponse)
    - [GetPageRequest](#services-wiki-GetPageRequest)
    - [GetPageResponse](#services-wiki-GetPageResponse)
    - [ListPageActivityRequest](#services-wiki-ListPageActivityRequest)
    - [ListPageActivityResponse](#services-wiki-ListPageActivityResponse)
    - [ListPagesRequest](#services-wiki-ListPagesRequest)
    - [ListPagesResponse](#services-wiki-ListPagesResponse)
    - [UpdatePageRequest](#services-wiki-UpdatePageRequest)
    - [UpdatePageResponse](#services-wiki-UpdatePageResponse)
  
    - [WikiService](#services-wiki-WikiService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="resources_accounts_oauth2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/accounts/oauth2.proto



<a name="resources-accounts-OAuth2Account"></a>

### OAuth2Account



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| provider_name | [string](#string) |  |  |
| provider | [OAuth2Provider](#resources-accounts-OAuth2Provider) |  |  |
| external_id | [uint64](#uint64) |  |  |
| username | [string](#string) |  |  |
| avatar | [string](#string) |  |  |






<a name="resources-accounts-OAuth2Provider"></a>

### OAuth2Provider



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| label | [string](#string) |  |  |
| homepage | [string](#string) |  |  |
| icon | [string](#string) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_accounts_accounts-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/accounts/accounts.proto



<a name="resources-accounts-Account"></a>

### Account



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| username | [string](#string) |  |  |
| license | [string](#string) |  |  |






<a name="resources-accounts-Character"></a>

### Character



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| available | [bool](#bool) |  |  |
| group | [string](#string) |  |  |
| char | [resources.users.User](#resources-users-User) |  | @gotags: alias:"user" |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/access.proto



<a name="resources-centrum-UnitAccess"></a>

### UnitAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [UnitJobAccess](#resources-centrum-UnitJobAccess) | repeated | @gotags: alias:"job_access" |
| qualifications | [UnitQualificationAccess](#resources-centrum-UnitQualificationAccess) | repeated | @gotags: alias:"qualification_access" |






<a name="resources-centrum-UnitJobAccess"></a>

### UnitJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"calendar_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [UnitAccessLevel](#resources-centrum-UnitAccessLevel) |  |  |






<a name="resources-centrum-UnitQualificationAccess"></a>

### UnitQualificationAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"thread_id" |
| qualification_id | [uint64](#uint64) |  |  |
| qualification | [resources.qualifications.QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| access | [UnitAccessLevel](#resources-centrum-UnitAccessLevel) |  |  |






<a name="resources-centrum-UnitUserAccess"></a>

### UnitUserAccess






 <!-- end messages -->


<a name="resources-centrum-UnitAccessLevel"></a>

### UnitAccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNIT_ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| UNIT_ACCESS_LEVEL_BLOCKED | 1 |  |
| UNIT_ACCESS_LEVEL_JOIN | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_dispatches-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/dispatches.proto



<a name="resources-centrum-Dispatch"></a>

### Dispatch



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| status | [DispatchStatus](#resources-centrum-DispatchStatus) | optional |  |
| message | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| attributes | [Attributes](#resources-centrum-Attributes) | optional |  |
| x | [double](#double) |  |  |
| y | [double](#double) |  |  |
| postal | [string](#string) | optional | @sanitize |
| anon | [bool](#bool) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.User](#resources-users-User) | optional |  |
| units | [DispatchAssignment](#resources-centrum-DispatchAssignment) | repeated |  |
| references | [DispatchReferences](#resources-centrum-DispatchReferences) | optional |  |






<a name="resources-centrum-DispatchAssignment"></a>

### DispatchAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch_id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"dispatch_id" |
| unit_id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"unit_id" |
| unit | [Unit](#resources-centrum-Unit) | optional |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| expires_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="resources-centrum-DispatchAssignments"></a>

### DispatchAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| units | [DispatchAssignment](#resources-centrum-DispatchAssignment) | repeated |  |






<a name="resources-centrum-DispatchReference"></a>

### DispatchReference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| target_dispatch_id | [uint64](#uint64) |  |  |
| reference_type | [DispatchReferenceType](#resources-centrum-DispatchReferenceType) |  |  |






<a name="resources-centrum-DispatchReferences"></a>

### DispatchReferences



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| references | [DispatchReference](#resources-centrum-DispatchReference) | repeated |  |






<a name="resources-centrum-DispatchStatus"></a>

### DispatchStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| dispatch_id | [uint64](#uint64) |  |  |
| unit_id | [uint64](#uint64) | optional |  |
| unit | [Unit](#resources-centrum-Unit) | optional |  |
| status | [StatusDispatch](#resources-centrum-StatusDispatch) |  |  |
| reason | [string](#string) | optional | @sanitize |
| code | [string](#string) | optional | @sanitize |
| user_id | [int32](#int32) | optional |  |
| user | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |
| x | [double](#double) | optional |  |
| y | [double](#double) | optional |  |
| postal | [string](#string) | optional | @sanitize |





 <!-- end messages -->


<a name="resources-centrum-DispatchReferenceType"></a>

### DispatchReferenceType


| Name | Number | Description |
| ---- | ------ | ----------- |
| DISPATCH_REFERENCE_TYPE_UNSPECIFIED | 0 |  |
| DISPATCH_REFERENCE_TYPE_REFERENCED | 1 |  |
| DISPATCH_REFERENCE_TYPE_DUPLICATED_BY | 2 |  |
| DISPATCH_REFERENCE_TYPE_DUPLICATE_OF | 3 |  |



<a name="resources-centrum-StatusDispatch"></a>

### StatusDispatch


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATUS_DISPATCH_UNSPECIFIED | 0 |  |
| STATUS_DISPATCH_NEW | 1 |  |
| STATUS_DISPATCH_UNASSIGNED | 2 |  |
| STATUS_DISPATCH_UPDATED | 3 |  |
| STATUS_DISPATCH_UNIT_ASSIGNED | 4 |  |
| STATUS_DISPATCH_UNIT_UNASSIGNED | 5 |  |
| STATUS_DISPATCH_UNIT_ACCEPTED | 6 |  |
| STATUS_DISPATCH_UNIT_DECLINED | 7 |  |
| STATUS_DISPATCH_EN_ROUTE | 8 |  |
| STATUS_DISPATCH_ON_SCENE | 9 |  |
| STATUS_DISPATCH_NEED_ASSISTANCE | 10 |  |
| STATUS_DISPATCH_COMPLETED | 11 |  |
| STATUS_DISPATCH_CANCELLED | 12 |  |
| STATUS_DISPATCH_ARCHIVED | 13 |  |



<a name="resources-centrum-TakeDispatchResp"></a>

### TakeDispatchResp


| Name | Number | Description |
| ---- | ------ | ----------- |
| TAKE_DISPATCH_RESP_UNSPECIFIED | 0 |  |
| TAKE_DISPATCH_RESP_TIMEOUT | 1 |  |
| TAKE_DISPATCH_RESP_ACCEPTED | 2 |  |
| TAKE_DISPATCH_RESP_DECLINED | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_general-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/general.proto



<a name="resources-centrum-Attributes"></a>

### Attributes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| list | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-centrum-Disponents"></a>

### Disponents



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| disponents | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |






<a name="resources-centrum-UserUnitMapping"></a>

### UserUnitMapping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| user_id | [int32](#int32) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_settings-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/settings.proto



<a name="resources-centrum-PredefinedStatus"></a>

### PredefinedStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_status | [string](#string) | repeated | @sanitize: method=StripTags |
| dispatch_status | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-centrum-Settings"></a>

### Settings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| enabled | [bool](#bool) |  |  |
| mode | [CentrumMode](#resources-centrum-CentrumMode) |  |  |
| fallback_mode | [CentrumMode](#resources-centrum-CentrumMode) |  |  |
| predefined_status | [PredefinedStatus](#resources-centrum-PredefinedStatus) | optional |  |
| timings | [Timings](#resources-centrum-Timings) |  |  |






<a name="resources-centrum-Timings"></a>

### Timings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch_max_wait | [int64](#int64) |  |  |
| require_unit | [bool](#bool) |  |  |
| require_unit_reminder_seconds | [int64](#int64) |  |  |





 <!-- end messages -->


<a name="resources-centrum-CentrumMode"></a>

### CentrumMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| CENTRUM_MODE_UNSPECIFIED | 0 |  |
| CENTRUM_MODE_MANUAL | 1 |  |
| CENTRUM_MODE_CENTRAL_COMMAND | 2 |  |
| CENTRUM_MODE_AUTO_ROUND_ROBIN | 3 |  |
| CENTRUM_MODE_SIMPLIFIED | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_units-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/units.proto



<a name="resources-centrum-Unit"></a>

### Unit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| name | [string](#string) |  | @sanitize |
| initials | [string](#string) |  | @sanitize |
| color | [string](#string) |  | @sanitize: method=StripTags |
| description | [string](#string) | optional | @sanitize |
| status | [UnitStatus](#resources-centrum-UnitStatus) | optional |  |
| users | [UnitAssignment](#resources-centrum-UnitAssignment) | repeated |  |
| attributes | [Attributes](#resources-centrum-Attributes) | optional |  |
| home_postal | [string](#string) | optional |  |
| access | [UnitAccess](#resources-centrum-UnitAccess) |  |  |






<a name="resources-centrum-UnitAssignment"></a>

### UnitAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"unit_id" |
| user_id | [int32](#int32) |  | @gotags: sql:"primary_key" alias:"user_id" |
| user | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |






<a name="resources-centrum-UnitAssignments"></a>

### UnitAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| users | [UnitAssignment](#resources-centrum-UnitAssignment) | repeated |  |






<a name="resources-centrum-UnitStatus"></a>

### UnitStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| unit_id | [uint64](#uint64) |  |  |
| unit | [Unit](#resources-centrum-Unit) | optional |  |
| status | [StatusUnit](#resources-centrum-StatusUnit) |  |  |
| reason | [string](#string) | optional | @sanitize |
| code | [string](#string) | optional | @sanitize |
| user_id | [int32](#int32) | optional |  |
| user | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |
| x | [double](#double) | optional |  |
| y | [double](#double) | optional |  |
| postal | [string](#string) | optional | @sanitize |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |





 <!-- end messages -->


<a name="resources-centrum-StatusUnit"></a>

### StatusUnit


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATUS_UNIT_UNSPECIFIED | 0 |  |
| STATUS_UNIT_UNKNOWN | 1 |  |
| STATUS_UNIT_USER_ADDED | 2 |  |
| STATUS_UNIT_USER_REMOVED | 3 |  |
| STATUS_UNIT_UNAVAILABLE | 4 |  |
| STATUS_UNIT_AVAILABLE | 5 |  |
| STATUS_UNIT_ON_BREAK | 6 |  |
| STATUS_UNIT_BUSY | 7 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_database_database-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/database/database.proto



<a name="resources-common-database-DateRange"></a>

### DateRange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| end | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="resources-common-database-PaginationRequest"></a>

### PaginationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| offset | [int64](#int64) |  |  |
| page_size | [int64](#int64) | optional |  |






<a name="resources-common-database-PaginationResponse"></a>

### PaginationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total_count | [int64](#int64) |  |  |
| offset | [int64](#int64) |  |  |
| end | [int64](#int64) |  |  |
| page_size | [int64](#int64) |  |  |






<a name="resources-common-database-Sort"></a>

### Sort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| column | [string](#string) |  |  |
| direction | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_grpcws_grpcws-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/grpcws/grpcws.proto



<a name="resources-common-grpcws-Body"></a>

### Body



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |
| complete | [bool](#bool) |  |  |






<a name="resources-common-grpcws-Cancel"></a>

### Cancel







<a name="resources-common-grpcws-Complete"></a>

### Complete







<a name="resources-common-grpcws-Failure"></a>

### Failure



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| error_message | [string](#string) |  |  |
| error_status | [string](#string) |  |  |
| headers | [Failure.HeadersEntry](#resources-common-grpcws-Failure-HeadersEntry) | repeated |  |






<a name="resources-common-grpcws-Failure-HeadersEntry"></a>

### Failure.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [HeaderValue](#resources-common-grpcws-HeaderValue) |  |  |






<a name="resources-common-grpcws-GrpcFrame"></a>

### GrpcFrame



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| streamId | [uint32](#uint32) |  |  |
| ping | [Ping](#resources-common-grpcws-Ping) |  |  |
| header | [Header](#resources-common-grpcws-Header) |  |  |
| body | [Body](#resources-common-grpcws-Body) |  |  |
| complete | [Complete](#resources-common-grpcws-Complete) |  |  |
| failure | [Failure](#resources-common-grpcws-Failure) |  |  |
| cancel | [Cancel](#resources-common-grpcws-Cancel) |  |  |






<a name="resources-common-grpcws-Header"></a>

### Header



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [string](#string) |  |  |
| headers | [Header.HeadersEntry](#resources-common-grpcws-Header-HeadersEntry) | repeated |  |
| status | [int32](#int32) |  |  |






<a name="resources-common-grpcws-Header-HeadersEntry"></a>

### Header.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [HeaderValue](#resources-common-grpcws-HeaderValue) |  |  |






<a name="resources-common-grpcws-HeaderValue"></a>

### HeaderValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) | repeated |  |






<a name="resources-common-grpcws-Ping"></a>

### Ping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pong | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_cron_cron-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/cron/cron.proto



<a name="resources-common-cron-Cronjob"></a>

### Cronjob



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| schedule | [string](#string) |  |  |
| state | [CronjobState](#resources-common-cron-CronjobState) |  |  |
| next_schedule_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| last_attempt_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| started_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| timeout | [google.protobuf.Duration](#google-protobuf-Duration) | optional |  |
| data | [CronjobData](#resources-common-cron-CronjobData) |  |  |






<a name="resources-common-cron-CronjobCompletedEvent"></a>

### CronjobCompletedEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| sucess | [bool](#bool) |  |  |
| endDate | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| elapsed | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| data | [CronjobData](#resources-common-cron-CronjobData) | optional |  |






<a name="resources-common-cron-CronjobData"></a>

### CronjobData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| data | [google.protobuf.Any](#google-protobuf-Any) | optional |  |






<a name="resources-common-cron-CronjobLockOwnerState"></a>

### CronjobLockOwnerState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hostname | [string](#string) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="resources-common-cron-CronjobSchedulerEvent"></a>

### CronjobSchedulerEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cronjob | [Cronjob](#resources-common-cron-Cronjob) |  |  |






<a name="resources-common-cron-GenericCronData"></a>

### GenericCronData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attributes | [GenericCronData.AttributesEntry](#resources-common-cron-GenericCronData-AttributesEntry) | repeated | @sanitize: method=StripTags |






<a name="resources-common-cron-GenericCronData-AttributesEntry"></a>

### GenericCronData.AttributesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 <!-- end messages -->


<a name="resources-common-cron-CronjobState"></a>

### CronjobState


| Name | Number | Description |
| ---- | ------ | ----------- |
| CRONJOB_STATE_UNSPECIFIED | 0 |  |
| CRONJOB_STATE_WAITING | 1 |  |
| CRONJOB_STATE_PENDING | 2 |  |
| CRONJOB_STATE_RUNNING | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_uuid-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/uuid.proto



<a name="resources-common-UUID"></a>

### UUID



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_content_content-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/content/content.proto



<a name="resources-common-content-Content"></a>

### Content



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) | optional |  |
| content | [JSONNode](#resources-common-content-JSONNode) | optional |  |
| raw_content | [string](#string) | optional | @sanitize |






<a name="resources-common-content-JSONNode"></a>

### JSONNode



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | @sanitize: method=StripTags |
| id | [string](#string) |  | @sanitize: method=StripTags |
| tag | [string](#string) |  | @sanitize: method=StripTags |
| attrs | [JSONNode.AttrsEntry](#resources-common-content-JSONNode-AttrsEntry) | repeated | @sanitize: method=StripTags |
| text | [string](#string) |  | @sanitize: method=StripTags |
| content | [JSONNode](#resources-common-content-JSONNode) | repeated |  |






<a name="resources-common-content-JSONNode-AttrsEntry"></a>

### JSONNode.AttrsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 <!-- end messages -->


<a name="resources-common-content-ContentType"></a>

### ContentType


| Name | Number | Description |
| ---- | ------ | ----------- |
| CONTENT_TYPE_UNSPECIFIED | 0 |  |
| CONTENT_TYPE_HTML | 1 |  |
| CONTENT_TYPE_PLAIN | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_error-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/error.proto



<a name="resources-common-Error"></a>

### Error



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [TranslateItem](#resources-common-TranslateItem) | optional |  |
| content | [TranslateItem](#resources-common-TranslateItem) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_i18n-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/i18n.proto



<a name="resources-common-TranslateItem"></a>

### TranslateItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | @sanitize: method=StripTags |
| parameters | [TranslateItem.ParametersEntry](#resources-common-TranslateItem-ParametersEntry) | repeated | @sanitize: method=StripTags |






<a name="resources-common-TranslateItem-ParametersEntry"></a>

### TranslateItem.ParametersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_category-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/category.proto



<a name="resources-documents-Category"></a>

### Category



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| name | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| job | [string](#string) | optional |  |
| color | [string](#string) | optional | @sanitize: method=StripTags |
| icon | [string](#string) | optional | @sanitize: method=StripTags |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/access.proto



<a name="resources-documents-DocumentAccess"></a>

### DocumentAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated | @gotags: alias:"job_access" |
| users | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-documents-DocumentJobAccess"></a>

### DocumentJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"document_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-documents-AccessLevel) |  |  |
| required | [bool](#bool) | optional | @gotags: alias:"required" |






<a name="resources-documents-DocumentUserAccess"></a>

### DocumentUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"document_id" |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-documents-AccessLevel) |  |  |
| required | [bool](#bool) | optional | @gotags: alias:"required" |





 <!-- end messages -->


<a name="resources-documents-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_BLOCKED | 1 |  |
| ACCESS_LEVEL_VIEW | 2 |  |
| ACCESS_LEVEL_COMMENT | 3 |  |
| ACCESS_LEVEL_STATUS | 4 |  |
| ACCESS_LEVEL_ACCESS | 5 |  |
| ACCESS_LEVEL_EDIT | 6 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_activity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/activity.proto



<a name="resources-documents-DocAccessJobsDiff"></a>

### DocAccessJobsDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to_create | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated |  |
| to_update | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated |  |
| to_delete | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated |  |






<a name="resources-documents-DocAccessRequested"></a>

### DocAccessRequested



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| level | [AccessLevel](#resources-documents-AccessLevel) |  |  |






<a name="resources-documents-DocAccessUpdated"></a>

### DocAccessUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [DocAccessJobsDiff](#resources-documents-DocAccessJobsDiff) |  |  |
| users | [DocAccessUsersDiff](#resources-documents-DocAccessUsersDiff) |  |  |






<a name="resources-documents-DocAccessUsersDiff"></a>

### DocAccessUsersDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to_create | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated |  |
| to_update | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated |  |
| to_delete | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated |  |






<a name="resources-documents-DocActivity"></a>

### DocActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| document_id | [uint64](#uint64) |  |  |
| activity_type | [DocActivityType](#resources-documents-DocActivityType) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| reason | [string](#string) | optional |  |
| data | [DocActivityData](#resources-documents-DocActivityData) |  |  |






<a name="resources-documents-DocActivityData"></a>

### DocActivityData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| updated | [DocUpdated](#resources-documents-DocUpdated) |  |  |
| owner_changed | [DocOwnerChanged](#resources-documents-DocOwnerChanged) |  |  |
| access_updated | [DocAccessUpdated](#resources-documents-DocAccessUpdated) |  |  |
| access_requested | [DocAccessRequested](#resources-documents-DocAccessRequested) |  |  |






<a name="resources-documents-DocOwnerChanged"></a>

### DocOwnerChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_owner_id | [int32](#int32) |  |  |
| new_owner | [resources.users.UserShort](#resources-users-UserShort) |  |  |






<a name="resources-documents-DocUpdated"></a>

### DocUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title_diff | [string](#string) | optional |  |
| content_diff | [string](#string) | optional |  |
| state_diff | [string](#string) | optional |  |





 <!-- end messages -->


<a name="resources-documents-DocActivityType"></a>

### DocActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOC_ACTIVITY_TYPE_UNSPECIFIED | 0 |  |
| DOC_ACTIVITY_TYPE_CREATED | 1 | Base |
| DOC_ACTIVITY_TYPE_STATUS_OPEN | 2 |  |
| DOC_ACTIVITY_TYPE_STATUS_CLOSED | 3 |  |
| DOC_ACTIVITY_TYPE_UPDATED | 4 |  |
| DOC_ACTIVITY_TYPE_RELATIONS_UPDATED | 5 |  |
| DOC_ACTIVITY_TYPE_REFERENCES_UPDATED | 6 |  |
| DOC_ACTIVITY_TYPE_ACCESS_UPDATED | 7 |  |
| DOC_ACTIVITY_TYPE_OWNER_CHANGED | 8 |  |
| DOC_ACTIVITY_TYPE_DELETED | 9 |  |
| DOC_ACTIVITY_TYPE_COMMENT_ADDED | 10 | Comments |
| DOC_ACTIVITY_TYPE_COMMENT_UPDATED | 11 |  |
| DOC_ACTIVITY_TYPE_COMMENT_DELETED | 12 |  |
| DOC_ACTIVITY_TYPE_REQUESTED_ACCESS | 13 | Requests |
| DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE | 14 |  |
| DOC_ACTIVITY_TYPE_REQUESTED_OPENING | 15 |  |
| DOC_ACTIVITY_TYPE_REQUESTED_UPDATE | 16 |  |
| DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE | 17 |  |
| DOC_ACTIVITY_TYPE_REQUESTED_DELETION | 18 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_comment-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/comment.proto



<a name="resources-documents-Comment"></a>

### Comment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| document_id | [uint64](#uint64) |  |  |
| content | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_documents-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/documents.proto



<a name="resources-documents-Document"></a>

### Document



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category_id | [uint64](#uint64) | optional |  |
| category | [Category](#resources-documents-Category) | optional | @gotags: alias:"category" |
| title | [string](#string) |  | @sanitize |
| content_type | [resources.common.content.ContentType](#resources-common-content-ContentType) |  | @gotags: alias:"content_type" |
| content | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| data | [string](#string) | optional | @sanitize

@gotags: alias:"data" |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| state | [string](#string) |  | @sanitize |
| closed | [bool](#bool) |  |  |
| public | [bool](#bool) |  |  |
| template_id | [uint64](#uint64) | optional |  |
| pinned | [bool](#bool) |  |  |
| workflow_state | [WorkflowState](#resources-documents-WorkflowState) | optional |  |
| workflow_user | [WorkflowUserState](#resources-documents-WorkflowUserState) | optional |  |






<a name="resources-documents-DocumentReference"></a>

### DocumentReference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) | optional |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| source_document_id | [uint64](#uint64) |  | @gotags: alias:"source_document_id" |
| source_document | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:"source_document" |
| reference | [DocReference](#resources-documents-DocReference) |  | @gotags: alias:"reference" |
| target_document_id | [uint64](#uint64) |  | @gotags: alias:"target_document_id" |
| target_document | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:"target_document" |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"ref_creator" |






<a name="resources-documents-DocumentRelation"></a>

### DocumentRelation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) | optional |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| document_id | [uint64](#uint64) |  |  |
| document | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:"document" |
| source_user_id | [int32](#int32) |  | @gotags: alias:"source_user_id" |
| source_user | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"source_user" |
| relation | [DocRelation](#resources-documents-DocRelation) |  | @gotags: alias:"relation" |
| target_user_id | [int32](#int32) |  | @gotags: alias:"target_user_id" |
| target_user | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"target_user" |






<a name="resources-documents-DocumentShort"></a>

### DocumentShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category_id | [uint64](#uint64) | optional |  |
| category | [Category](#resources-documents-Category) | optional | @gotags: alias:"category" |
| title | [string](#string) |  | @sanitize |
| content_type | [resources.common.content.ContentType](#resources-common-content-ContentType) |  | @gotags: alias:"content_type" |
| content | [resources.common.content.Content](#resources-common-content-Content) |  | @sanitize |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  | @gotags: alias:"creator_job" |
| creator_job_label | [string](#string) | optional |  |
| state | [string](#string) |  | @sanitize

@gotags: alias:"state" |
| closed | [bool](#bool) |  |  |
| public | [bool](#bool) |  |  |
| workflow_state | [WorkflowState](#resources-documents-WorkflowState) | optional |  |
| workflow_user | [WorkflowUserState](#resources-documents-WorkflowUserState) | optional |  |






<a name="resources-documents-WorkflowState"></a>

### WorkflowState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| next_reminder_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| next_reminder_count | [int32](#int32) | optional |  |
| auto_close_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| workflow | [Workflow](#resources-documents-Workflow) | optional | @gotags: alias:"workflow" |
| document | [DocumentShort](#resources-documents-DocumentShort) | optional |  |






<a name="resources-documents-WorkflowUserState"></a>

### WorkflowUserState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| manual_reminder_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| manual_reminder_message | [string](#string) | optional |  |
| workflow | [Workflow](#resources-documents-Workflow) | optional | @gotags: alias:"workflow" |
| document | [DocumentShort](#resources-documents-DocumentShort) | optional |  |





 <!-- end messages -->


<a name="resources-documents-DocReference"></a>

### DocReference


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOC_REFERENCE_UNSPECIFIED | 0 |  |
| DOC_REFERENCE_LINKED | 1 |  |
| DOC_REFERENCE_SOLVES | 2 |  |
| DOC_REFERENCE_CLOSES | 3 |  |
| DOC_REFERENCE_DEPRECATES | 4 |  |



<a name="resources-documents-DocRelation"></a>

### DocRelation


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOC_RELATION_UNSPECIFIED | 0 |  |
| DOC_RELATION_MENTIONED | 1 |  |
| DOC_RELATION_TARGETS | 2 |  |
| DOC_RELATION_CAUSED | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_requests-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/requests.proto



<a name="resources-documents-DocRequest"></a>

### DocRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| document_id | [uint64](#uint64) |  |  |
| request_type | [DocActivityType](#resources-documents-DocActivityType) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| reason | [string](#string) | optional |  |
| data | [DocActivityData](#resources-documents-DocActivityData) |  |  |
| accepted | [bool](#bool) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_templates-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/templates.proto



<a name="resources-documents-ObjectSpecs"></a>

### ObjectSpecs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| required | [bool](#bool) | optional |  |
| min | [int32](#int32) | optional |  |
| max | [int32](#int32) | optional |  |






<a name="resources-documents-Template"></a>

### Template



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category | [Category](#resources-documents-Category) |  | @gotags: alias:"category" |
| weight | [uint32](#uint32) |  |  |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) |  | @sanitize |
| color | [string](#string) | optional | @sanitize: method=StripTags |
| icon | [string](#string) | optional | @sanitize: method=StripTags |
| content_title | [string](#string) |  | @gotags: alias:"content_title" |
| content | [string](#string) |  | @gotags: alias:"content" |
| state | [string](#string) |  | @gotags: alias:"state" |
| schema | [TemplateSchema](#resources-documents-TemplateSchema) |  | @gotags: alias:"schema" |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| job_access | [TemplateJobAccess](#resources-documents-TemplateJobAccess) | repeated |  |
| content_access | [DocumentAccess](#resources-documents-DocumentAccess) |  | @gotags: alias:"access" |
| workflow | [Workflow](#resources-documents-Workflow) | optional |  |






<a name="resources-documents-TemplateData"></a>

### TemplateData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| activeChar | [resources.users.User](#resources-users-User) |  |  |
| documents | [DocumentShort](#resources-documents-DocumentShort) | repeated |  |
| users | [resources.users.UserShort](#resources-users-UserShort) | repeated |  |
| vehicles | [resources.vehicles.Vehicle](#resources-vehicles-Vehicle) | repeated |  |






<a name="resources-documents-TemplateJobAccess"></a>

### TemplateJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"template_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-documents-AccessLevel) |  |  |






<a name="resources-documents-TemplateRequirements"></a>

### TemplateRequirements



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| documents | [ObjectSpecs](#resources-documents-ObjectSpecs) | optional |  |
| users | [ObjectSpecs](#resources-documents-ObjectSpecs) | optional |  |
| vehicles | [ObjectSpecs](#resources-documents-ObjectSpecs) | optional |  |






<a name="resources-documents-TemplateSchema"></a>

### TemplateSchema



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requirements | [TemplateRequirements](#resources-documents-TemplateRequirements) |  |  |






<a name="resources-documents-TemplateShort"></a>

### TemplateShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category | [Category](#resources-documents-Category) |  | @gotags: alias:"category" |
| weight | [uint32](#uint32) |  |  |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) |  | @sanitize |
| color | [string](#string) | optional | @sanitize: method=StripTags |
| icon | [string](#string) | optional | @sanitize: method=StripTags |
| schema | [TemplateSchema](#resources-documents-TemplateSchema) |  | @gotags: alias:"schema" |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| workflow | [Workflow](#resources-documents-Workflow) | optional |  |






<a name="resources-documents-TemplateUserAccess"></a>

### TemplateUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_workflow-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/workflow.proto



<a name="resources-documents-AutoCloseSettings"></a>

### AutoCloseSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| duration | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| message | [string](#string) |  |  |






<a name="resources-documents-Reminder"></a>

### Reminder



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| duration | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| message | [string](#string) |  |  |






<a name="resources-documents-ReminderSettings"></a>

### ReminderSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reminders | [Reminder](#resources-documents-Reminder) | repeated |  |






<a name="resources-documents-Workflow"></a>

### Workflow



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reminder | [bool](#bool) |  |  |
| reminder_settings | [ReminderSettings](#resources-documents-ReminderSettings) |  |  |
| auto_close | [bool](#bool) |  |  |
| auto_close_settings | [AutoCloseSettings](#resources-documents-AutoCloseSettings) |  |  |






<a name="resources-documents-WorkflowCronData"></a>

### WorkflowCronData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| last_doc_id | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_filestore_file-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/filestore/file.proto



<a name="resources-filestore-File"></a>

### File



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url | [string](#string) | optional |  |
| data | [bytes](#bytes) |  |  |
| delete | [bool](#bool) | optional |  |
| content_type | [string](#string) | optional |  |
| extension | [string](#string) | optional |  |






<a name="resources-filestore-FileInfo"></a>

### FileInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| last_modified | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| size | [int64](#int64) |  |  |
| content_type | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_conduct-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/conduct.proto



<a name="resources-jobs-ConductEntry"></a>

### ConductEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| type | [ConductType](#resources-jobs-ConductType) |  |  |
| message | [string](#string) |  | @sanitize |
| expires_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_user_id | [int32](#int32) |  |  |
| target_user | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:"target_user" |
| creator_id | [int32](#int32) |  |  |
| creator | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:"creator" |





 <!-- end messages -->


<a name="resources-jobs-ConductType"></a>

### ConductType


| Name | Number | Description |
| ---- | ------ | ----------- |
| CONDUCT_TYPE_UNSPECIFIED | 0 |  |
| CONDUCT_TYPE_NEUTRAL | 1 |  |
| CONDUCT_TYPE_POSITIVE | 2 |  |
| CONDUCT_TYPE_NEGATIVE | 3 |  |
| CONDUCT_TYPE_WARNING | 4 |  |
| CONDUCT_TYPE_SUSPENSION | 5 |  |
| CONDUCT_TYPE_NOTE | 6 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_colleagues-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/colleagues.proto



<a name="resources-jobs-Colleague"></a>

### Colleague



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  | @gotags: alias:"id" |
| identifier | [string](#string) | optional |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| job_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| firstname | [string](#string) |  |  |
| lastname | [string](#string) |  |  |
| dateofbirth | [string](#string) |  |  |
| phone_number | [string](#string) | optional |  |
| avatar | [resources.filestore.File](#resources-filestore-File) | optional |  |
| props | [JobsUserProps](#resources-jobs-JobsUserProps) |  | @gotags: alias:"jobs_user_props" |
| email | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-jobs-ColleagueAbsenceDate"></a>

### ColleagueAbsenceDate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| absence_begin | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| absence_end | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="resources-jobs-ColleagueGradeChange"></a>

### ColleagueGradeChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| grade | [int32](#int32) |  |  |
| grade_label | [string](#string) |  |  |






<a name="resources-jobs-ColleagueLabelsChange"></a>

### ColleagueLabelsChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| added | [Label](#resources-jobs-Label) | repeated |  |
| removed | [Label](#resources-jobs-Label) | repeated |  |






<a name="resources-jobs-ColleagueNameChange"></a>

### ColleagueNameChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| prefix | [string](#string) | optional |  |
| suffix | [string](#string) | optional |  |






<a name="resources-jobs-JobsUserActivity"></a>

### JobsUserActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| source_user_id | [int32](#int32) |  |  |
| source_user | [Colleague](#resources-jobs-Colleague) |  | @gotags: alias:"source_user" |
| target_user_id | [int32](#int32) |  |  |
| target_user | [Colleague](#resources-jobs-Colleague) |  | @gotags: alias:"target_user" |
| activity_type | [JobsUserActivityType](#resources-jobs-JobsUserActivityType) |  |  |
| reason | [string](#string) |  | @sanitize |
| data | [JobsUserActivityData](#resources-jobs-JobsUserActivityData) |  |  |






<a name="resources-jobs-JobsUserActivityData"></a>

### JobsUserActivityData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| absence_date | [ColleagueAbsenceDate](#resources-jobs-ColleagueAbsenceDate) |  |  |
| grade_change | [ColleagueGradeChange](#resources-jobs-ColleagueGradeChange) |  |  |
| labels_change | [ColleagueLabelsChange](#resources-jobs-ColleagueLabelsChange) |  |  |
| name_change | [ColleagueNameChange](#resources-jobs-ColleagueNameChange) |  |  |






<a name="resources-jobs-JobsUserProps"></a>

### JobsUserProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| job | [string](#string) |  |  |
| absence_begin | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| absence_end | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| note | [string](#string) | optional | @sanitize: method=StripTags |
| labels | [Labels](#resources-jobs-Labels) | optional |  |
| name_prefix | [string](#string) | optional |  |
| name_suffix | [string](#string) | optional |  |





 <!-- end messages -->


<a name="resources-jobs-JobsUserActivityType"></a>

### JobsUserActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| JOBS_USER_ACTIVITY_TYPE_UNSPECIFIED | 0 |  |
| JOBS_USER_ACTIVITY_TYPE_HIRED | 1 |  |
| JOBS_USER_ACTIVITY_TYPE_FIRED | 2 |  |
| JOBS_USER_ACTIVITY_TYPE_PROMOTED | 3 |  |
| JOBS_USER_ACTIVITY_TYPE_DEMOTED | 4 |  |
| JOBS_USER_ACTIVITY_TYPE_ABSENCE_DATE | 5 |  |
| JOBS_USER_ACTIVITY_TYPE_NOTE | 6 |  |
| JOBS_USER_ACTIVITY_TYPE_LABELS | 7 |  |
| JOBS_USER_ACTIVITY_TYPE_NAME | 8 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_labels-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/labels.proto



<a name="resources-jobs-Label"></a>

### Label



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| job | [string](#string) | optional |  |
| name | [string](#string) |  |  |
| color | [string](#string) |  | @sanitize: method=StripTags |
| order | [int32](#int32) |  |  |






<a name="resources-jobs-LabelCount"></a>

### LabelCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| label | [Label](#resources-jobs-Label) |  |  |
| count | [int64](#int64) |  |  |






<a name="resources-jobs-Labels"></a>

### Labels



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| list | [Label](#resources-jobs-Label) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_timeclock-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/timeclock.proto



<a name="resources-jobs-TimeclockEntry"></a>

### TimeclockEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  | @gotags: sql:"primary_key" |
| date | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: sql:"primary_key" |
| user | [Colleague](#resources-jobs-Colleague) | optional |  |
| start_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| end_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| spent_time | [float](#float) |  |  |






<a name="resources-jobs-TimeclockStats"></a>

### TimeclockStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| spent_time_sum | [float](#float) |  |  |
| spent_time_avg | [float](#float) |  |  |
| spent_time_max | [float](#float) |  |  |






<a name="resources-jobs-TimeclockWeeklyStats"></a>

### TimeclockWeeklyStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| year | [int32](#int32) |  |  |
| calendar_week | [int32](#int32) |  |  |
| sum | [float](#float) |  |  |
| avg | [float](#float) |  |  |
| max | [float](#float) |  |  |





 <!-- end messages -->


<a name="resources-jobs-TimeclockMode"></a>

### TimeclockMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| TIMECLOCK_MODE_UNSPECIFIED | 0 |  |
| TIMECLOCK_MODE_DAILY | 1 |  |
| TIMECLOCK_MODE_WEEKLY | 2 |  |
| TIMECLOCK_MODE_RANGE | 3 |  |



<a name="resources-jobs-TimeclockUserMode"></a>

### TimeclockUserMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| TIMECLOCK_USER_MODE_UNSPECIFIED | 0 |  |
| TIMECLOCK_USER_MODE_SELF | 1 |  |
| TIMECLOCK_USER_MODE_ALL | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_laws_laws-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/laws/laws.proto



<a name="resources-laws-Law"></a>

### Law



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"law.id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| lawbook_id | [uint64](#uint64) |  |  |
| name | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| hint | [string](#string) | optional | @sanitize |
| fine | [uint32](#uint32) | optional |  |
| detention_time | [uint32](#uint32) | optional |  |
| stvo_points | [uint32](#uint32) | optional |  |






<a name="resources-laws-LawBook"></a>

### LawBook



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| name | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| laws | [Law](#resources-laws-Law) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_livemap_tracker-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/tracker.proto



<a name="resources-livemap-UsersUpdateEvent"></a>

### UsersUpdateEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| added | [UserMarker](#resources-livemap-UserMarker) | repeated |  |
| removed | [UserMarker](#resources-livemap-UserMarker) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_livemap_livemap-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/livemap.proto



<a name="resources-livemap-CircleMarker"></a>

### CircleMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| radius | [int32](#int32) |  |  |
| opacity | [float](#float) | optional |  |






<a name="resources-livemap-Coords"></a>

### Coords



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| x | [double](#double) |  |  |
| y | [double](#double) |  |  |






<a name="resources-livemap-IconMarker"></a>

### IconMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| icon | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-livemap-MarkerData"></a>

### MarkerData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| circle | [CircleMarker](#resources-livemap-CircleMarker) |  |  |
| icon | [IconMarker](#resources-livemap-IconMarker) |  |  |






<a name="resources-livemap-MarkerInfo"></a>

### MarkerInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) |  |  |
| name | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| x | [double](#double) |  |  |
| y | [double](#double) |  |  |
| postal | [string](#string) | optional | @sanitize |
| color | [string](#string) | optional | @sanitize: method=StripTags |
| icon | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-livemap-MarkerMarker"></a>

### MarkerMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| info | [MarkerInfo](#resources-livemap-MarkerInfo) |  |  |
| type | [MarkerType](#resources-livemap-MarkerType) |  | @gotags: alias:"markerType" |
| expires_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| data | [MarkerData](#resources-livemap-MarkerData) |  | @gotags: alias:"markerData" |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional |  |






<a name="resources-livemap-UserMarker"></a>

### UserMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| info | [MarkerInfo](#resources-livemap-MarkerInfo) |  |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.jobs.Colleague](#resources-jobs-Colleague) |  | @gotags: alias:"user" |
| unit_id | [uint64](#uint64) | optional |  |
| unit | [resources.centrum.Unit](#resources-centrum-Unit) | optional |  |
| hidden | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="resources-livemap-MarkerType"></a>

### MarkerType


| Name | Number | Description |
| ---- | ------ | ----------- |
| MARKER_TYPE_UNSPECIFIED | 0 |  |
| MARKER_TYPE_DOT | 1 |  |
| MARKER_TYPE_CIRCLE | 2 |  |
| MARKER_TYPE_ICON | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_notifications_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/notifications/events.proto



<a name="resources-notifications-JobEvent"></a>

### JobEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_props | [resources.users.JobProps](#resources-users-JobProps) |  |  |






<a name="resources-notifications-JobGradeEvent"></a>

### JobGradeEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [bool](#bool) |  |  |






<a name="resources-notifications-SystemEvent"></a>

### SystemEvent







<a name="resources-notifications-UserEvent"></a>

### UserEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [bool](#bool) |  |  |
| notification | [Notification](#resources-notifications-Notification) |  | Notifications |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_notifications_notifications-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/notifications/notifications.proto



<a name="resources-notifications-CalendarData"></a>

### CalendarData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar_id | [uint64](#uint64) | optional |  |
| calendar_entry_id | [uint64](#uint64) | optional |  |






<a name="resources-notifications-Data"></a>

### Data



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| link | [Link](#resources-notifications-Link) | optional |  |
| caused_by | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| calendar | [CalendarData](#resources-notifications-CalendarData) | optional |  |






<a name="resources-notifications-Link"></a>

### Link



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to | [string](#string) |  |  |
| title | [string](#string) | optional |  |
| external | [bool](#bool) | optional |  |






<a name="resources-notifications-Notification"></a>

### Notification



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| read_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| user_id | [int32](#int32) |  |  |
| title | [resources.common.TranslateItem](#resources-common-TranslateItem) |  | @sanitize |
| type | [NotificationType](#resources-notifications-NotificationType) |  |  |
| content | [resources.common.TranslateItem](#resources-common-TranslateItem) |  | @sanitize |
| category | [NotificationCategory](#resources-notifications-NotificationCategory) |  |  |
| data | [Data](#resources-notifications-Data) | optional |  |
| starred | [bool](#bool) | optional |  |





 <!-- end messages -->


<a name="resources-notifications-NotificationCategory"></a>

### NotificationCategory


| Name | Number | Description |
| ---- | ------ | ----------- |
| NOTIFICATION_CATEGORY_UNSPECIFIED | 0 |  |
| NOTIFICATION_CATEGORY_GENERAL | 1 |  |
| NOTIFICATION_CATEGORY_DOCUMENT | 2 |  |
| NOTIFICATION_CATEGORY_CALENDAR | 3 |  |



<a name="resources-notifications-NotificationType"></a>

### NotificationType


| Name | Number | Description |
| ---- | ------ | ----------- |
| NOTIFICATION_TYPE_UNSPECIFIED | 0 |  |
| NOTIFICATION_TYPE_ERROR | 1 |  |
| NOTIFICATION_TYPE_WARNING | 2 |  |
| NOTIFICATION_TYPE_INFO | 3 |  |
| NOTIFICATION_TYPE_SUCCESS | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_permissions_permissions-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/permissions/permissions.proto



<a name="resources-permissions-AttributeValues"></a>

### AttributeValues



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| string_list | [StringList](#resources-permissions-StringList) |  |  |
| job_list | [StringList](#resources-permissions-StringList) |  |  |
| job_grade_list | [JobGradeList](#resources-permissions-JobGradeList) |  |  |






<a name="resources-permissions-JobGradeList"></a>

### JobGradeList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [JobGradeList.JobsEntry](#resources-permissions-JobGradeList-JobsEntry) | repeated |  |






<a name="resources-permissions-JobGradeList-JobsEntry"></a>

### JobGradeList.JobsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [int32](#int32) |  |  |






<a name="resources-permissions-Permission"></a>

### Permission



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category | [string](#string) |  |  |
| name | [string](#string) |  |  |
| guard_name | [string](#string) |  |  |
| val | [bool](#bool) |  |  |






<a name="resources-permissions-RawRoleAttribute"></a>

### RawRoleAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| attr_id | [uint64](#uint64) |  |  |
| permission_id | [uint64](#uint64) |  |  |
| category | [string](#string) |  |  |
| name | [string](#string) |  |  |
| key | [string](#string) |  |  |
| type | [string](#string) |  |  |
| valid_values | [AttributeValues](#resources-permissions-AttributeValues) |  |  |
| value | [AttributeValues](#resources-permissions-AttributeValues) |  |  |






<a name="resources-permissions-Role"></a>

### Role



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| permissions | [Permission](#resources-permissions-Permission) | repeated |  |
| attributes | [RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="resources-permissions-RoleAttribute"></a>

### RoleAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| attr_id | [uint64](#uint64) |  |  |
| permission_id | [uint64](#uint64) |  |  |
| category | [string](#string) |  |  |
| name | [string](#string) |  |  |
| key | [string](#string) |  |  |
| type | [string](#string) |  |  |
| valid_values | [AttributeValues](#resources-permissions-AttributeValues) |  |  |
| value | [AttributeValues](#resources-permissions-AttributeValues) |  |  |
| max_values | [AttributeValues](#resources-permissions-AttributeValues) | optional |  |






<a name="resources-permissions-StringList"></a>

### StringList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| strings | [string](#string) | repeated | @sanitize: method=StripTags |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_qualifications_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/qualifications/access.proto



<a name="resources-qualifications-QualificationAccess"></a>

### QualificationAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [QualificationJobAccess](#resources-qualifications-QualificationJobAccess) | repeated |  |






<a name="resources-qualifications-QualificationJobAccess"></a>

### QualificationJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"qualification_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-qualifications-AccessLevel) |  |  |






<a name="resources-qualifications-QualificationUserAccess"></a>

### QualificationUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


<a name="resources-qualifications-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_BLOCKED | 1 |  |
| ACCESS_LEVEL_VIEW | 2 |  |
| ACCESS_LEVEL_REQUEST | 3 |  |
| ACCESS_LEVEL_TAKE | 4 |  |
| ACCESS_LEVEL_GRADE | 5 |  |
| ACCESS_LEVEL_MANAGE | 6 |  |
| ACCESS_LEVEL_EDIT | 7 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_qualifications_exam-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/qualifications/exam.proto



<a name="resources-qualifications-ExamQuestion"></a>

### ExamQuestion



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| qualification_id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| title | [string](#string) |  | @sanitize: method=StripTags |
| description | [string](#string) | optional | @sanitize: method=StripTags |
| data | [ExamQuestionData](#resources-qualifications-ExamQuestionData) |  |  |
| answer | [ExamQuestionAnswerData](#resources-qualifications-ExamQuestionAnswerData) | optional |  |
| points | [int32](#int32) | optional |  |






<a name="resources-qualifications-ExamQuestionAnswerData"></a>

### ExamQuestionAnswerData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| answer_key | [string](#string) |  |  |






<a name="resources-qualifications-ExamQuestionData"></a>

### ExamQuestionData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| separator | [ExamQuestionSeparator](#resources-qualifications-ExamQuestionSeparator) |  |  |
| image | [ExamQuestionImage](#resources-qualifications-ExamQuestionImage) |  |  |
| yesno | [ExamQuestionYesNo](#resources-qualifications-ExamQuestionYesNo) |  |  |
| free_text | [ExamQuestionText](#resources-qualifications-ExamQuestionText) |  |  |
| single_choice | [ExamQuestionSingleChoice](#resources-qualifications-ExamQuestionSingleChoice) |  |  |
| multiple_choice | [ExamQuestionMultipleChoice](#resources-qualifications-ExamQuestionMultipleChoice) |  |  |






<a name="resources-qualifications-ExamQuestionImage"></a>

### ExamQuestionImage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| image | [resources.filestore.File](#resources-filestore-File) |  |  |
| alt | [string](#string) | optional |  |






<a name="resources-qualifications-ExamQuestionMultipleChoice"></a>

### ExamQuestionMultipleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choices | [string](#string) | repeated | @sanitize: method=StripTags |
| limit | [int32](#int32) | optional |  |






<a name="resources-qualifications-ExamQuestionSeparator"></a>

### ExamQuestionSeparator







<a name="resources-qualifications-ExamQuestionSingleChoice"></a>

### ExamQuestionSingleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choices | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-qualifications-ExamQuestionText"></a>

### ExamQuestionText



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| min_length | [int32](#int32) |  |  |
| max_length | [int32](#int32) |  |  |






<a name="resources-qualifications-ExamQuestionYesNo"></a>

### ExamQuestionYesNo







<a name="resources-qualifications-ExamQuestions"></a>

### ExamQuestions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| questions | [ExamQuestion](#resources-qualifications-ExamQuestion) | repeated |  |






<a name="resources-qualifications-ExamResponse"></a>

### ExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| question_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| question | [ExamQuestion](#resources-qualifications-ExamQuestion) |  |  |
| response | [ExamResponseData](#resources-qualifications-ExamResponseData) |  |  |






<a name="resources-qualifications-ExamResponseData"></a>

### ExamResponseData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| separator | [ExamResponseSeparator](#resources-qualifications-ExamResponseSeparator) |  |  |
| yesno | [ExamResponseYesNo](#resources-qualifications-ExamResponseYesNo) |  |  |
| free_text | [ExamResponseText](#resources-qualifications-ExamResponseText) |  |  |
| single_choice | [ExamResponseSingleChoice](#resources-qualifications-ExamResponseSingleChoice) |  |  |
| multiple_choice | [ExamResponseMultipleChoice](#resources-qualifications-ExamResponseMultipleChoice) |  |  |






<a name="resources-qualifications-ExamResponseMultipleChoice"></a>

### ExamResponseMultipleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choices | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-qualifications-ExamResponseSeparator"></a>

### ExamResponseSeparator







<a name="resources-qualifications-ExamResponseSingleChoice"></a>

### ExamResponseSingleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choice | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-qualifications-ExamResponseText"></a>

### ExamResponseText



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  | @sanitize: method=StripTags

0.5 Megabyte |






<a name="resources-qualifications-ExamResponseYesNo"></a>

### ExamResponseYesNo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [bool](#bool) |  |  |






<a name="resources-qualifications-ExamResponses"></a>

### ExamResponses



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| responses | [ExamResponse](#resources-qualifications-ExamResponse) | repeated |  |






<a name="resources-qualifications-ExamUser"></a>

### ExamUser



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| started_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| ends_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| ended_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_qualifications_qualifications-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/qualifications/qualifications.proto



<a name="resources-qualifications-Qualification"></a>

### Qualification



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| weight | [uint32](#uint32) |  |  |
| closed | [bool](#bool) |  |  |
| abbreviation | [string](#string) |  | @sanitize: method=StripTags |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize: method=StripTags |
| content | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| access | [QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |
| requirements | [QualificationRequirement](#resources-qualifications-QualificationRequirement) | repeated |  |
| discord_sync_enabled | [bool](#bool) |  |  |
| discord_settings | [QualificationDiscordSettings](#resources-qualifications-QualificationDiscordSettings) | optional |  |
| exam_mode | [QualificationExamMode](#resources-qualifications-QualificationExamMode) |  |  |
| exam_settings | [QualificationExamSettings](#resources-qualifications-QualificationExamSettings) | optional |  |
| exam | [ExamQuestions](#resources-qualifications-ExamQuestions) | optional |  |
| result | [QualificationResult](#resources-qualifications-QualificationResult) | optional |  |
| request | [QualificationRequest](#resources-qualifications-QualificationRequest) | optional |  |
| label_sync_enabled | [bool](#bool) |  |  |
| label_sync_format | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-qualifications-QualificationDiscordSettings"></a>

### QualificationDiscordSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_name | [string](#string) | optional |  |
| role_format | [string](#string) | optional |  |






<a name="resources-qualifications-QualificationExamSettings"></a>

### QualificationExamSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |






<a name="resources-qualifications-QualificationRequest"></a>

### QualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| qualification_id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"qualification_id" |
| qualification | [QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| user_id | [int32](#int32) |  | @gotags: sql:"primary_key" |
| user | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:"user" |
| user_comment | [string](#string) | optional | @sanitize: method=StripTags |
| status | [RequestStatus](#resources-qualifications-RequestStatus) | optional |  |
| approved_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| approver_comment | [string](#string) | optional | @sanitize: method=StripTags |
| approver_id | [int32](#int32) | optional |  |
| approver | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"approver" |
| approver_job | [string](#string) | optional |  |






<a name="resources-qualifications-QualificationRequirement"></a>

### QualificationRequirement



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| qualification_id | [uint64](#uint64) |  |  |
| target_qualification_id | [uint64](#uint64) |  |  |
| target_qualification | [QualificationShort](#resources-qualifications-QualificationShort) | optional | @gotags: alias:"targetqualification.*" |






<a name="resources-qualifications-QualificationResult"></a>

### QualificationResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| qualification_id | [uint64](#uint64) |  |  |
| qualification | [QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:"user" |
| status | [ResultStatus](#resources-qualifications-ResultStatus) |  |  |
| score | [uint32](#uint32) | optional |  |
| summary | [string](#string) |  | @sanitize: method=StripTags |
| creator_id | [int32](#int32) |  |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |






<a name="resources-qualifications-QualificationShort"></a>

### QualificationShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| weight | [uint32](#uint32) |  |  |
| closed | [bool](#bool) |  |  |
| abbreviation | [string](#string) |  | @sanitize: method=StripTags |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize: method=StripTags |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| requirements | [QualificationRequirement](#resources-qualifications-QualificationRequirement) | repeated |  |
| exam_mode | [QualificationExamMode](#resources-qualifications-QualificationExamMode) |  |  |
| exam_settings | [QualificationExamSettings](#resources-qualifications-QualificationExamSettings) | optional |  |
| result | [QualificationResult](#resources-qualifications-QualificationResult) | optional |  |





 <!-- end messages -->


<a name="resources-qualifications-QualificationExamMode"></a>

### QualificationExamMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| QUALIFICATION_EXAM_MODE_UNSPECIFIED | 0 |  |
| QUALIFICATION_EXAM_MODE_DISABLED | 1 |  |
| QUALIFICATION_EXAM_MODE_REQUEST_NEEDED | 2 |  |
| QUALIFICATION_EXAM_MODE_ENABLED | 3 |  |



<a name="resources-qualifications-RequestStatus"></a>

### RequestStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| REQUEST_STATUS_UNSPECIFIED | 0 |  |
| REQUEST_STATUS_PENDING | 1 |  |
| REQUEST_STATUS_DENIED | 2 |  |
| REQUEST_STATUS_ACCEPTED | 3 |  |
| REQUEST_STATUS_EXAM_STARTED | 4 |  |
| REQUEST_STATUS_EXAM_GRADING | 5 |  |
| REQUEST_STATUS_COMPLETED | 6 |  |



<a name="resources-qualifications-ResultStatus"></a>

### ResultStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| RESULT_STATUS_UNSPECIFIED | 0 |  |
| RESULT_STATUS_PENDING | 1 |  |
| RESULT_STATUS_FAILED | 2 |  |
| RESULT_STATUS_SUCCESSFUL | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_rector_audit-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/rector/audit.proto



<a name="resources-rector-AuditEntry"></a>

### AuditEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| user_id | [uint64](#uint64) |  | @gotags: alias:"user_id" |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| user_job | [string](#string) |  | @gotags: alias:"user_job" |
| target_user_id | [int32](#int32) | optional | @gotags: alias:"target_user_id" |
| target_user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| target_user_job | [string](#string) |  | @gotags: alias:"target_user_job" |
| service | [string](#string) |  | @gotags: alias:"service" |
| method | [string](#string) |  | @gotags: alias:"method" |
| state | [EventType](#resources-rector-EventType) |  | @gotags: alias:"state" |
| data | [string](#string) | optional | @gotags: alias:"data" |





 <!-- end messages -->


<a name="resources-rector-EventType"></a>

### EventType


| Name | Number | Description |
| ---- | ------ | ----------- |
| EVENT_TYPE_UNSPECIFIED | 0 |  |
| EVENT_TYPE_ERRORED | 1 |  |
| EVENT_TYPE_VIEWED | 2 |  |
| EVENT_TYPE_CREATED | 3 |  |
| EVENT_TYPE_UPDATED | 4 |  |
| EVENT_TYPE_DELETED | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_rector_config-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/rector/config.proto



<a name="resources-rector-AppConfig"></a>

### AppConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) | optional |  |
| auth | [Auth](#resources-rector-Auth) |  |  |
| perms | [Perms](#resources-rector-Perms) |  |  |
| website | [Website](#resources-rector-Website) |  |  |
| job_info | [JobInfo](#resources-rector-JobInfo) |  |  |
| user_tracker | [UserTracker](#resources-rector-UserTracker) |  |  |
| discord | [Discord](#resources-rector-Discord) |  |  |






<a name="resources-rector-Auth"></a>

### Auth



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signup_enabled | [bool](#bool) |  |  |
| last_char_lock | [bool](#bool) |  |  |






<a name="resources-rector-Discord"></a>

### Discord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enabled | [bool](#bool) |  |  |
| sync_interval | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| invite_url | [string](#string) | optional | @sanitize: method=StripTags |
| ignored_jobs | [string](#string) | repeated | @sanitize: method=StripTags |
| bot_presence | [DiscordBotPresence](#resources-rector-DiscordBotPresence) | optional |  |






<a name="resources-rector-DiscordBotPresence"></a>

### DiscordBotPresence



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [DiscordBotPresenceType](#resources-rector-DiscordBotPresenceType) |  |  |
| status | [string](#string) | optional | @sanitize: method=StripTags |
| url | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-rector-JobInfo"></a>

### JobInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unemployed_job | [UnemployedJob](#resources-rector-UnemployedJob) |  |  |
| public_jobs | [string](#string) | repeated | @sanitize: method=StripTags |
| hidden_jobs | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-rector-Links"></a>

### Links



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| privacy_policy | [string](#string) | optional | @sanitize: method=StripTags |
| imprint | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-rector-Perm"></a>

### Perm



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [string](#string) |  | @sanitize: method=StripTags |
| name | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-rector-Perms"></a>

### Perms



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default | [Perm](#resources-rector-Perm) | repeated |  |






<a name="resources-rector-UnemployedJob"></a>

### UnemployedJob



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| grade | [int32](#int32) |  |  |






<a name="resources-rector-UserTracker"></a>

### UserTracker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_time | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| db_refresh_time | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| livemap_jobs | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-rector-Website"></a>

### Website



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| links | [Links](#resources-rector-Links) |  |  |
| stats_page | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="resources-rector-DiscordBotPresenceType"></a>

### DiscordBotPresenceType


| Name | Number | Description |
| ---- | ------ | ----------- |
| DISCORD_BOT_PRESENCE_TYPE_UNSPECIFIED | 0 |  |
| DISCORD_BOT_PRESENCE_TYPE_GAME | 1 |  |
| DISCORD_BOT_PRESENCE_TYPE_LISTENING | 2 |  |
| DISCORD_BOT_PRESENCE_TYPE_STREAMING | 3 |  |
| DISCORD_BOT_PRESENCE_TYPE_WATCH | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_timestamp_timestamp-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/timestamp/timestamp.proto



<a name="resources-timestamp-Timestamp"></a>

### Timestamp
Timestamp for storage messages. We've defined a new local type wrapper of google.protobuf.Timestamp so we can implement sql.Scanner and sql.Valuer interfaces. See: https://golang.org/pkg/database/sql/#Scanner https://golang.org/pkg/database/sql/driver/#Valuer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_job_props-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/job_props.proto



<a name="resources-users-DiscordSyncChange"></a>

### DiscordSyncChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| plan | [string](#string) |  |  |






<a name="resources-users-DiscordSyncChanges"></a>

### DiscordSyncChanges



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| changes | [DiscordSyncChange](#resources-users-DiscordSyncChange) | repeated |  |






<a name="resources-users-DiscordSyncSettings"></a>

### DiscordSyncSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dry_run | [bool](#bool) |  |  |
| user_info_sync | [bool](#bool) |  |  |
| user_info_sync_settings | [UserInfoSyncSettings](#resources-users-UserInfoSyncSettings) |  |  |
| status_log | [bool](#bool) |  |  |
| status_log_settings | [StatusLogSettings](#resources-users-StatusLogSettings) |  |  |
| jobs_absence | [bool](#bool) |  |  |
| jobs_absence_settings | [JobsAbsenceSettings](#resources-users-JobsAbsenceSettings) |  |  |
| group_sync_settings | [GroupSyncSettings](#resources-users-GroupSyncSettings) |  |  |
| qualifications_role_format | [string](#string) |  |  |






<a name="resources-users-GroupMapping"></a>

### GroupMapping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| from_grade | [int32](#int32) |  |  |
| to_grade | [int32](#int32) |  |  |






<a name="resources-users-GroupSyncSettings"></a>

### GroupSyncSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ignored_role_ids | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-users-JobProps"></a>

### JobProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| theme | [string](#string) |  |  |
| livemap_marker_color | [string](#string) |  |  |
| quick_buttons | [QuickButtons](#resources-users-QuickButtons) |  |  |
| radio_frequency | [string](#string) | optional |  |
| discord_guild_id | [string](#string) | optional |  |
| discord_last_sync | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| discord_sync_settings | [DiscordSyncSettings](#resources-users-DiscordSyncSettings) |  |  |
| discord_sync_changes | [DiscordSyncChanges](#resources-users-DiscordSyncChanges) | optional |  |
| motd | [string](#string) | optional |  |
| logo_url | [resources.filestore.File](#resources-filestore-File) | optional |  |
| settings | [JobSettings](#resources-users-JobSettings) |  |  |






<a name="resources-users-JobSettings"></a>

### JobSettings







<a name="resources-users-JobsAbsenceSettings"></a>

### JobsAbsenceSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| absence_role | [string](#string) |  |  |






<a name="resources-users-QuickButtons"></a>

### QuickButtons



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| penalty_calculator | [bool](#bool) |  |  |
| body_checkup | [bool](#bool) |  |  |
| math_calculator | [bool](#bool) |  |  |






<a name="resources-users-StatusLogSettings"></a>

### StatusLogSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel_id | [string](#string) |  |  |






<a name="resources-users-UserInfoSyncSettings"></a>

### UserInfoSyncSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_role_enabled | [bool](#bool) |  |  |
| employee_role_format | [string](#string) |  |  |
| grade_role_format | [string](#string) |  |  |
| unemployed_enabled | [bool](#bool) |  |  |
| unemployed_mode | [UserInfoSyncUnemployedMode](#resources-users-UserInfoSyncUnemployedMode) |  |  |
| unemployed_role_name | [string](#string) |  |  |
| sync_nicknames | [bool](#bool) |  |  |
| group_mapping | [GroupMapping](#resources-users-GroupMapping) | repeated |  |





 <!-- end messages -->


<a name="resources-users-UserInfoSyncUnemployedMode"></a>

### UserInfoSyncUnemployedMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_INFO_SYNC_UNEMPLOYED_MODE_UNSPECIFIED | 0 |  |
| USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE | 1 |  |
| USER_INFO_SYNC_UNEMPLOYED_MODE_KICK | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_jobs-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/jobs.proto



<a name="resources-users-Job"></a>

### Job



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | @gotags: sql:"primary_key" alias:"name" |
| label | [string](#string) |  |  |
| grades | [JobGrade](#resources-users-JobGrade) | repeated |  |






<a name="resources-users-JobGrade"></a>

### JobGrade



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_name | [string](#string) | optional |  |
| grade | [int32](#int32) |  |  |
| label | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_users-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/users.proto



<a name="resources-users-CitizenAttribute"></a>

### CitizenAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| job | [string](#string) | optional |  |
| name | [string](#string) |  |  |
| color | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-users-CitizenAttributes"></a>

### CitizenAttributes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| list | [CitizenAttribute](#resources-users-CitizenAttribute) | repeated |  |






<a name="resources-users-License"></a>

### License



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  |  |
| label | [string](#string) |  |  |






<a name="resources-users-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  | @gotags: alias:"id" |
| identifier | [string](#string) | optional |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| job_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| firstname | [string](#string) |  |  |
| lastname | [string](#string) |  |  |
| dateofbirth | [string](#string) |  |  |
| sex | [string](#string) | optional |  |
| height | [string](#string) | optional |  |
| phone_number | [string](#string) | optional |  |
| visum | [int32](#int32) | optional |  |
| playtime | [int32](#int32) | optional |  |
| props | [UserProps](#resources-users-UserProps) |  | @gotags: alias:"fivenet_user_props" |
| licenses | [License](#resources-users-License) | repeated | @gotags: alias:"user_licenses" |
| avatar | [resources.filestore.File](#resources-filestore-File) | optional |  |






<a name="resources-users-UserActivity"></a>

### UserActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:"fivenet_user_activity.id" |
| type | [UserActivityType](#resources-users-UserActivityType) |  | @gotags: alias:"fivenet_user_activity.type" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: alias:"fivenet_user_activity.created_at" |
| source_user | [UserShort](#resources-users-UserShort) |  | @gotags: alias:"source_user" |
| target_user | [UserShort](#resources-users-UserShort) |  | @gotags: alias:"target_user" |
| key | [string](#string) |  | @sanitize

@gotags: alias:"fivenet_user_activity.key" |
| old_value | [string](#string) |  | @gotags: alias:"fivenet_user_activity.old_value" |
| new_value | [string](#string) |  | @gotags: alias:"fivenet_user_activity.new_value" |
| reason | [string](#string) |  | @sanitize

@gotags: alias:"fivenet_user_activity.reason" |






<a name="resources-users-UserProps"></a>

### UserProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| wanted | [bool](#bool) | optional |  |
| job_name | [string](#string) | optional | @gotags: alias:"job" |
| job | [Job](#resources-users-Job) | optional |  |
| job_grade_number | [int32](#int32) | optional | @gotags: alias:"job_grade" |
| job_grade | [JobGrade](#resources-users-JobGrade) | optional |  |
| traffic_infraction_points | [uint32](#uint32) | optional |  |
| open_fines | [int64](#int64) | optional |  |
| blood_type | [string](#string) | optional |  |
| mug_shot | [resources.filestore.File](#resources-filestore-File) | optional |  |
| attributes | [CitizenAttributes](#resources-users-CitizenAttributes) | optional |  |
| email | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-users-UserShort"></a>

### UserShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  | @gotags: alias:"id" |
| identifier | [string](#string) | optional |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| job_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| firstname | [string](#string) |  |  |
| lastname | [string](#string) |  |  |
| dateofbirth | [string](#string) |  |  |
| phone_number | [string](#string) | optional |  |
| avatar | [resources.filestore.File](#resources-filestore-File) | optional |  |





 <!-- end messages -->


<a name="resources-users-UserActivityType"></a>

### UserActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_ACTIVITY_TYPE_UNSPECIFIED | 0 |  |
| USER_ACTIVITY_TYPE_CHANGED | 1 |  |
| USER_ACTIVITY_TYPE_MENTIONED | 2 |  |
| USER_ACTIVITY_TYPE_CREATED | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_vehicles_vehicles-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/vehicles/vehicles.proto



<a name="resources-vehicles-Vehicle"></a>

### Vehicle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plate | [string](#string) |  |  |
| model | [string](#string) | optional |  |
| type | [string](#string) |  |  |
| owner | [resources.users.UserShort](#resources-users-UserShort) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_calendar_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/calendar/access.proto



<a name="resources-calendar-CalendarAccess"></a>

### CalendarAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [CalendarJobAccess](#resources-calendar-CalendarJobAccess) | repeated | @gotags: alias:"job_access" |
| users | [CalendarUserAccess](#resources-calendar-CalendarUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-calendar-CalendarJobAccess"></a>

### CalendarJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"calendar_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-calendar-AccessLevel) |  |  |






<a name="resources-calendar-CalendarUserAccess"></a>

### CalendarUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"calendar_id" |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-calendar-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-calendar-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_BLOCKED | 1 |  |
| ACCESS_LEVEL_VIEW | 2 |  |
| ACCESS_LEVEL_SHARE | 3 |  |
| ACCESS_LEVEL_EDIT | 4 |  |
| ACCESS_LEVEL_MANAGE | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_calendar_calendar-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/calendar/calendar.proto



<a name="resources-calendar-Calendar"></a>

### Calendar



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) | optional |  |
| name | [string](#string) |  | @sanitize: method=StripTags |
| description | [string](#string) | optional | @sanitize: method=StripTags |
| public | [bool](#bool) |  |  |
| closed | [bool](#bool) |  |  |
| color | [string](#string) |  | @sanitize: method=StripTags |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| subscription | [CalendarSub](#resources-calendar-CalendarSub) | optional |  |
| access | [CalendarAccess](#resources-calendar-CalendarAccess) |  |  |






<a name="resources-calendar-CalendarEntry"></a>

### CalendarEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| calendar_id | [uint64](#uint64) |  |  |
| calendar | [Calendar](#resources-calendar-Calendar) | optional |  |
| job | [string](#string) | optional |  |
| start_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| end_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| title | [string](#string) |  | @sanitize: method=StripTags |
| content | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| closed | [bool](#bool) |  |  |
| rsvp_open | [bool](#bool) | optional |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| recurring | [CalendarEntryRecurring](#resources-calendar-CalendarEntryRecurring) | optional |  |
| rsvp | [CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) | optional |  |






<a name="resources-calendar-CalendarEntryRSVP"></a>

### CalendarEntryRSVP



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry_id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| response | [RsvpResponses](#resources-calendar-RsvpResponses) |  |  |






<a name="resources-calendar-CalendarEntryRecurring"></a>

### CalendarEntryRecurring



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| every | [string](#string) |  |  |
| count | [int32](#int32) |  |  |
| until | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="resources-calendar-CalendarShort"></a>

### CalendarShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| name | [string](#string) |  | @sanitize: method=StripTags |
| description | [string](#string) | optional | @sanitize: method=StripTags |
| public | [bool](#bool) |  |  |
| closed | [bool](#bool) |  |  |
| color | [string](#string) |  | @sanitize: method=StripTags |
| subscription | [CalendarSub](#resources-calendar-CalendarSub) | optional |  |






<a name="resources-calendar-CalendarSub"></a>

### CalendarSub



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| confirmed | [bool](#bool) |  |  |
| muted | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="resources-calendar-RsvpResponses"></a>

### RsvpResponses


| Name | Number | Description |
| ---- | ------ | ----------- |
| RSVP_RESPONSES_UNSPECIFIED | 0 |  |
| RSVP_RESPONSES_HIDDEN | 1 |  |
| RSVP_RESPONSES_INVITED | 2 |  |
| RSVP_RESPONSES_NO | 3 |  |
| RSVP_RESPONSES_MAYBE | 4 |  |
| RSVP_RESPONSES_YES | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_stats_stats-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/stats/stats.proto



<a name="resources-stats-Stat"></a>

### Stat



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [int32](#int32) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_internet_ads-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/internet/ads.proto



<a name="resources-internet-Ad"></a>

### Ad



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| disabled | [bool](#bool) |  |  |
| ad_type | [AdType](#resources-internet-AdType) |  |  |
| starts_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| ends_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| title | [string](#string) |  | @sanitize: method=StripTags |
| description | [string](#string) |  | @sanitize: method=StripTags |
| image | [resources.filestore.File](#resources-filestore-File) | optional |  |
| approver_id | [int32](#int32) | optional |  |
| approver_job | [string](#string) | optional |  |
| creator_id | [int32](#int32) | optional |  |
| creator_job | [string](#string) | optional |  |





 <!-- end messages -->


<a name="resources-internet-AdType"></a>

### AdType


| Name | Number | Description |
| ---- | ------ | ----------- |
| AD_TYPE_UNSPECIFIED | 0 |  |
| AD_TYPE_SPONSORED | 1 |  |
| AD_TYPE_SEARCH_RESULT | 2 |  |
| AD_TYPE_CONTENT_MAIN | 3 |  |
| AD_TYPE_CONTENT_ASIDE | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_internet_search-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/internet/search.proto



<a name="resources-internet-SearchResult"></a>

### SearchResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| url | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_internet_domain-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/internet/domain.proto



<a name="resources-internet-Domain"></a>

### Domain



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| name | [string](#string) |  |  |
| creator_job | [string](#string) | optional |  |
| creator_id | [int32](#int32) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_internet_page-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/internet/page.proto



<a name="resources-internet-Page"></a>

### Page



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| domain_id | [uint64](#uint64) |  |  |
| path | [string](#string) |  | @sanitize: method=StripTags |
| title | [string](#string) |  | @sanitize: method=StripTags |
| description | [string](#string) |  | @sanitize: method=StripTags |
| data | [PageData](#resources-internet-PageData) |  |  |
| creator_job | [string](#string) | optional |  |
| creator_id | [int32](#int32) | optional |  |






<a name="resources-internet-PageData"></a>

### PageData
TODO





 <!-- end messages -->


<a name="resources-internet-PageLayoutType"></a>

### PageLayoutType


| Name | Number | Description |
| ---- | ------ | ----------- |
| PAGE_LAYOUT_TYPE_UNSPECIFIED | 0 |  |
| PAGE_LAYOUT_TYPE_ | 1 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/access.proto



<a name="resources-mailer-Access"></a>

### Access



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [JobAccess](#resources-mailer-JobAccess) | repeated | @gotags: alias:"job_access" |
| users | [UserAccess](#resources-mailer-UserAccess) | repeated | @gotags: alias:"user_access" |
| qualifications | [QualificationAccess](#resources-mailer-QualificationAccess) | repeated | @gotags: alias:"qualification_access" |






<a name="resources-mailer-JobAccess"></a>

### JobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"email_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-mailer-AccessLevel) |  |  |






<a name="resources-mailer-QualificationAccess"></a>

### QualificationAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"thread_id" |
| qualification_id | [uint64](#uint64) |  |  |
| qualification | [resources.qualifications.QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| access | [AccessLevel](#resources-mailer-AccessLevel) |  |  |






<a name="resources-mailer-UserAccess"></a>

### UserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"thread_id" |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-mailer-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-mailer-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_BLOCKED | 1 |  |
| ACCESS_LEVEL_READ | 2 |  |
| ACCESS_LEVEL_WRITE | 3 |  |
| ACCESS_LEVEL_MANAGE | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_email-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/email.proto



<a name="resources-mailer-Email"></a>

### Email



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deactivated | [bool](#bool) |  |  |
| job | [string](#string) | optional |  |
| user_id | [int32](#int32) | optional |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| email | [string](#string) |  | @sanitize: method=StripTags |
| email_changed | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| label | [string](#string) | optional | @sanitize: method=StripTags |
| internal | [bool](#bool) |  |  |
| access | [Access](#resources-mailer-Access) |  |  |
| settings | [EmailSettings](#resources-mailer-EmailSettings) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/events.proto



<a name="resources-mailer-MailerEvent"></a>

### MailerEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_update | [Email](#resources-mailer-Email) |  |  |
| email_delete | [uint64](#uint64) |  |  |
| email_settings_updated | [EmailSettings](#resources-mailer-EmailSettings) |  |  |
| thread_update | [Thread](#resources-mailer-Thread) |  |  |
| thread_delete | [uint64](#uint64) |  |  |
| thread_state_update | [ThreadState](#resources-mailer-ThreadState) |  |  |
| message_update | [Message](#resources-mailer-Message) |  |  |
| message_delete | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_message-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/message.proto



<a name="resources-mailer-Message"></a>

### Message



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |
| sender_id | [uint64](#uint64) |  |  |
| sender | [Email](#resources-mailer-Email) | optional | @gotags: alias:"sender" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| title | [string](#string) |  | @sanitize: method=StripTags |
| content | [resources.common.content.Content](#resources-common-content-Content) |  | @sanitize |
| data | [MessageData](#resources-mailer-MessageData) | optional |  |
| creator_id | [int32](#int32) | optional |  |
| creator_job | [string](#string) | optional |  |






<a name="resources-mailer-MessageData"></a>

### MessageData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [MessageDataEntry](#resources-mailer-MessageDataEntry) | repeated |  |






<a name="resources-mailer-MessageDataEntry"></a>

### MessageDataEntry
TODO add way to link to, e.g., internal "objects" (citizens, documents, calendar entries, qualifications, etc.)





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_settings-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/settings.proto



<a name="resources-mailer-EmailSettings"></a>

### EmailSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| signature | [string](#string) | optional | @sanitize |
| blocked_emails | [string](#string) | repeated | @sanitize: method=StripTags |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_template-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/template.proto



<a name="resources-mailer-Template"></a>

### Template



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| email_id | [uint64](#uint64) |  |  |
| title | [string](#string) |  | @sanitize: method=StripTags |
| content | [string](#string) |  | @sanitize |
| creator_job | [string](#string) | optional |  |
| creator_id | [int32](#int32) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_mailer_thread-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/mailer/thread.proto



<a name="resources-mailer-Thread"></a>

### Thread



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| creator_email_id | [uint64](#uint64) |  |  |
| creator_email | [Email](#resources-mailer-Email) | optional |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| title | [string](#string) |  | @sanitize: method=StripTags |
| recipients | [ThreadRecipientEmail](#resources-mailer-ThreadRecipientEmail) | repeated |  |
| state | [ThreadState](#resources-mailer-ThreadState) | optional | @gotags: alias:"thread_state" |






<a name="resources-mailer-ThreadRecipientEmail"></a>

### ThreadRecipientEmail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"thread_id" |
| email_id | [uint64](#uint64) |  |  |
| email | [Email](#resources-mailer-Email) | optional |  |






<a name="resources-mailer-ThreadState"></a>

### ThreadState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |
| email_id | [uint64](#uint64) |  |  |
| last_read | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| unread | [bool](#bool) | optional |  |
| important | [bool](#bool) | optional |  |
| favorite | [bool](#bool) | optional |  |
| muted | [bool](#bool) | optional |  |
| archived | [bool](#bool) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_wiki_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/wiki/access.proto



<a name="resources-wiki-PageAccess"></a>

### PageAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated | @gotags: alias:"job_access" |
| users | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-wiki-PageJobAccess"></a>

### PageJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"page_id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-wiki-AccessLevel) |  |  |






<a name="resources-wiki-PageUserAccess"></a>

### PageUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_id | [uint64](#uint64) |  | @gotags: alias:"page_id" |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-wiki-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-wiki-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_BLOCKED | 1 |  |
| ACCESS_LEVEL_VIEW | 2 |  |
| ACCESS_LEVEL_ACCESS | 3 |  |
| ACCESS_LEVEL_EDIT | 4 |  |
| ACCESS_LEVEL_OWNER | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_wiki_activity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/wiki/activity.proto



<a name="resources-wiki-PageAccessJobsDiff"></a>

### PageAccessJobsDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to_create | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated |  |
| to_update | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated |  |
| to_delete | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated |  |






<a name="resources-wiki-PageAccessUpdated"></a>

### PageAccessUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [PageAccessJobsDiff](#resources-wiki-PageAccessJobsDiff) |  |  |
| users | [PageAccessUsersDiff](#resources-wiki-PageAccessUsersDiff) |  |  |






<a name="resources-wiki-PageAccessUsersDiff"></a>

### PageAccessUsersDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to_create | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated |  |
| to_update | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated |  |
| to_delete | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated |  |






<a name="resources-wiki-PageActivity"></a>

### PageActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| page_id | [uint64](#uint64) |  |  |
| activity_type | [PageActivityType](#resources-wiki-PageActivityType) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| reason | [string](#string) | optional |  |
| data | [PageActivityData](#resources-wiki-PageActivityData) |  |  |






<a name="resources-wiki-PageActivityData"></a>

### PageActivityData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| updated | [PageUpdated](#resources-wiki-PageUpdated) |  |  |
| access_updated | [PageAccessUpdated](#resources-wiki-PageAccessUpdated) |  |  |






<a name="resources-wiki-PageUpdated"></a>

### PageUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title_diff | [string](#string) | optional |  |
| description_diff | [string](#string) | optional |  |
| content_diff | [string](#string) | optional |  |





 <!-- end messages -->


<a name="resources-wiki-PageActivityType"></a>

### PageActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| PAGE_ACTIVITY_TYPE_UNSPECIFIED | 0 |  |
| PAGE_ACTIVITY_TYPE_CREATED | 1 | Base |
| PAGE_ACTIVITY_TYPE_UPDATED | 2 |  |
| PAGE_ACTIVITY_TYPE_ACCESS_UPDATED | 3 |  |
| PAGE_ACTIVITY_TYPE_OWNER_CHANGED | 4 |  |
| PAGE_ACTIVITY_TYPE_DELETED | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_wiki_page-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/wiki/page.proto



<a name="resources-wiki-Page"></a>

### Page



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| job | [string](#string) |  | @sanitize: method=StripTags |
| job_label | [string](#string) | optional |  |
| parent_id | [uint64](#uint64) | optional |  |
| meta | [PageMeta](#resources-wiki-PageMeta) |  |  |
| content | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| access | [PageAccess](#resources-wiki-PageAccess) |  |  |






<a name="resources-wiki-PageMeta"></a>

### PageMeta



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| slug | [string](#string) | optional | @sanitize: method=StripTags |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) |  | @sanitize: method=StripTags |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| content_type | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |
| tags | [string](#string) | repeated | @sanitize: method=StripTags |
| toc | [bool](#bool) | optional |  |
| public | [bool](#bool) |  |  |






<a name="resources-wiki-PageRootInfo"></a>

### PageRootInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| logo | [resources.filestore.File](#resources-filestore-File) | optional |  |






<a name="resources-wiki-PageShort"></a>

### PageShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| parent_id | [uint64](#uint64) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| slug | [string](#string) | optional | @sanitize: method=StripTags |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| children | [PageShort](#resources-wiki-PageShort) | repeated |  |
| root_info | [PageRootInfo](#resources-wiki-PageRootInfo) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="services_auth_auth-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/auth/auth.proto



<a name="services-auth-ChangePasswordRequest"></a>

### ChangePasswordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| current | [string](#string) |  |  |
| new | [string](#string) |  |  |






<a name="services-auth-ChangePasswordResponse"></a>

### ChangePasswordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| expires | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="services-auth-ChangeUsernameRequest"></a>

### ChangeUsernameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| current | [string](#string) |  |  |
| new | [string](#string) |  |  |






<a name="services-auth-ChangeUsernameResponse"></a>

### ChangeUsernameResponse







<a name="services-auth-ChooseCharacterRequest"></a>

### ChooseCharacterRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| char_id | [int32](#int32) |  |  |






<a name="services-auth-ChooseCharacterResponse"></a>

### ChooseCharacterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| expires | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| permissions | [string](#string) | repeated |  |
| job_props | [resources.users.JobProps](#resources-users-JobProps) |  |  |
| char | [resources.users.User](#resources-users-User) |  | @gotags: alias:"user" |
| username | [string](#string) |  |  |






<a name="services-auth-CreateAccountRequest"></a>

### CreateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reg_token | [string](#string) |  |  |
| username | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="services-auth-CreateAccountResponse"></a>

### CreateAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_id | [uint64](#uint64) |  |  |






<a name="services-auth-DeleteOAuth2ConnectionRequest"></a>

### DeleteOAuth2ConnectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| provider | [string](#string) |  |  |






<a name="services-auth-DeleteOAuth2ConnectionResponse"></a>

### DeleteOAuth2ConnectionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="services-auth-ForgotPasswordRequest"></a>

### ForgotPasswordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reg_token | [string](#string) |  |  |
| new | [string](#string) |  |  |






<a name="services-auth-ForgotPasswordResponse"></a>

### ForgotPasswordResponse







<a name="services-auth-GetAccountInfoRequest"></a>

### GetAccountInfoRequest







<a name="services-auth-GetAccountInfoResponse"></a>

### GetAccountInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [resources.accounts.Account](#resources-accounts-Account) |  |  |
| oauth2_providers | [resources.accounts.OAuth2Provider](#resources-accounts-OAuth2Provider) | repeated |  |
| oauth2_connections | [resources.accounts.OAuth2Account](#resources-accounts-OAuth2Account) | repeated |  |






<a name="services-auth-GetCharactersRequest"></a>

### GetCharactersRequest







<a name="services-auth-GetCharactersResponse"></a>

### GetCharactersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chars | [resources.accounts.Character](#resources-accounts-Character) | repeated |  |






<a name="services-auth-LoginRequest"></a>

### LoginRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="services-auth-LoginResponse"></a>

### LoginResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| expires | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| account_id | [uint64](#uint64) |  |  |
| char | [ChooseCharacterResponse](#services-auth-ChooseCharacterResponse) | optional |  |






<a name="services-auth-LogoutRequest"></a>

### LogoutRequest







<a name="services-auth-LogoutResponse"></a>

### LogoutResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="services-auth-SetSuperUserModeRequest"></a>

### SetSuperUserModeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| superuser | [bool](#bool) |  |  |
| job | [string](#string) | optional |  |






<a name="services-auth-SetSuperUserModeResponse"></a>

### SetSuperUserModeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| expires | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| job_props | [resources.users.JobProps](#resources-users-JobProps) | optional |  |
| char | [resources.users.User](#resources-users-User) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-auth-AuthService"></a>

### AuthService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Login | [LoginRequest](#services-auth-LoginRequest) | [LoginResponse](#services-auth-LoginResponse) |  |
| Logout | [LogoutRequest](#services-auth-LogoutRequest) | [LogoutResponse](#services-auth-LogoutResponse) |  |
| CreateAccount | [CreateAccountRequest](#services-auth-CreateAccountRequest) | [CreateAccountResponse](#services-auth-CreateAccountResponse) |  |
| ChangeUsername | [ChangeUsernameRequest](#services-auth-ChangeUsernameRequest) | [ChangeUsernameResponse](#services-auth-ChangeUsernameResponse) |  |
| ChangePassword | [ChangePasswordRequest](#services-auth-ChangePasswordRequest) | [ChangePasswordResponse](#services-auth-ChangePasswordResponse) |  |
| ForgotPassword | [ForgotPasswordRequest](#services-auth-ForgotPasswordRequest) | [ForgotPasswordResponse](#services-auth-ForgotPasswordResponse) |  |
| GetCharacters | [GetCharactersRequest](#services-auth-GetCharactersRequest) | [GetCharactersResponse](#services-auth-GetCharactersResponse) |  |
| ChooseCharacter | [ChooseCharacterRequest](#services-auth-ChooseCharacterRequest) | [ChooseCharacterResponse](#services-auth-ChooseCharacterResponse) | @perm |
| GetAccountInfo | [GetAccountInfoRequest](#services-auth-GetAccountInfoRequest) | [GetAccountInfoResponse](#services-auth-GetAccountInfoResponse) |  |
| DeleteOAuth2Connection | [DeleteOAuth2ConnectionRequest](#services-auth-DeleteOAuth2ConnectionRequest) | [DeleteOAuth2ConnectionResponse](#services-auth-DeleteOAuth2ConnectionResponse) |  |
| SetSuperUserMode | [SetSuperUserModeRequest](#services-auth-SetSuperUserModeRequest) | [SetSuperUserModeResponse](#services-auth-SetSuperUserModeResponse) |  |

 <!-- end services -->



<a name="services_centrum_centrum-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/centrum/centrum.proto



<a name="services-centrum-AssignDispatchRequest"></a>

### AssignDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch_id | [uint64](#uint64) |  |  |
| to_add | [uint64](#uint64) | repeated |  |
| to_remove | [uint64](#uint64) | repeated |  |
| forced | [bool](#bool) | optional |  |






<a name="services-centrum-AssignDispatchResponse"></a>

### AssignDispatchResponse







<a name="services-centrum-AssignUnitRequest"></a>

### AssignUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  |  |
| to_add | [int32](#int32) | repeated |  |
| to_remove | [int32](#int32) | repeated |  |






<a name="services-centrum-AssignUnitResponse"></a>

### AssignUnitResponse







<a name="services-centrum-CreateDispatchRequest"></a>

### CreateDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-CreateDispatchResponse"></a>

### CreateDispatchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-CreateOrUpdateUnitRequest"></a>

### CreateOrUpdateUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |






<a name="services-centrum-CreateOrUpdateUnitResponse"></a>

### CreateOrUpdateUnitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |






<a name="services-centrum-DeleteDispatchRequest"></a>

### DeleteDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-centrum-DeleteDispatchResponse"></a>

### DeleteDispatchResponse







<a name="services-centrum-DeleteUnitRequest"></a>

### DeleteUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  |  |






<a name="services-centrum-DeleteUnitResponse"></a>

### DeleteUnitResponse







<a name="services-centrum-GetDispatchRequest"></a>

### GetDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-centrum-GetDispatchResponse"></a>

### GetDispatchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-GetSettingsRequest"></a>

### GetSettingsRequest







<a name="services-centrum-GetSettingsResponse"></a>

### GetSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| settings | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |






<a name="services-centrum-JoinUnitRequest"></a>

### JoinUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) | optional |  |






<a name="services-centrum-JoinUnitResponse"></a>

### JoinUnitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |






<a name="services-centrum-LatestState"></a>

### LatestState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| settings | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |
| disponents | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |
| own_unit_id | [uint64](#uint64) | optional |  |
| units | [resources.centrum.Unit](#resources-centrum-Unit) | repeated | Send the current units and dispatches |
| dispatches | [resources.centrum.Dispatch](#resources-centrum-Dispatch) | repeated |  |






<a name="services-centrum-ListDispatchActivityRequest"></a>

### ListDispatchActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| id | [uint64](#uint64) |  |  |






<a name="services-centrum-ListDispatchActivityResponse"></a>

### ListDispatchActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.centrum.DispatchStatus](#resources-centrum-DispatchStatus) | repeated |  |






<a name="services-centrum-ListDispatchesRequest"></a>

### ListDispatchesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| status | [resources.centrum.StatusDispatch](#resources-centrum-StatusDispatch) | repeated |  |
| not_status | [resources.centrum.StatusDispatch](#resources-centrum-StatusDispatch) | repeated |  |
| ids | [uint64](#uint64) | repeated |  |
| postal | [string](#string) | optional |  |






<a name="services-centrum-ListDispatchesResponse"></a>

### ListDispatchesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| dispatches | [resources.centrum.Dispatch](#resources-centrum-Dispatch) | repeated |  |






<a name="services-centrum-ListUnitActivityRequest"></a>

### ListUnitActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| id | [uint64](#uint64) |  |  |






<a name="services-centrum-ListUnitActivityResponse"></a>

### ListUnitActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.centrum.UnitStatus](#resources-centrum-UnitStatus) | repeated |  |






<a name="services-centrum-ListUnitsRequest"></a>

### ListUnitsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [resources.centrum.StatusUnit](#resources-centrum-StatusUnit) | repeated |  |






<a name="services-centrum-ListUnitsResponse"></a>

### ListUnitsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| units | [resources.centrum.Unit](#resources-centrum-Unit) | repeated |  |






<a name="services-centrum-StreamRequest"></a>

### StreamRequest







<a name="services-centrum-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| latest_state | [LatestState](#services-centrum-LatestState) |  |  |
| settings | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |
| disponents | [resources.centrum.Disponents](#resources-centrum-Disponents) |  |  |
| unit_created | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |
| unit_deleted | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |
| unit_updated | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |
| unit_status | [resources.centrum.UnitStatus](#resources-centrum-UnitStatus) |  |  |
| dispatch_created | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |
| dispatch_deleted | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |
| dispatch_updated | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |
| dispatch_status | [resources.centrum.DispatchStatus](#resources-centrum-DispatchStatus) |  |  |






<a name="services-centrum-TakeControlRequest"></a>

### TakeControlRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signon | [bool](#bool) |  |  |






<a name="services-centrum-TakeControlResponse"></a>

### TakeControlResponse







<a name="services-centrum-TakeDispatchRequest"></a>

### TakeDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch_ids | [uint64](#uint64) | repeated |  |
| resp | [resources.centrum.TakeDispatchResp](#resources-centrum-TakeDispatchResp) |  |  |
| reason | [string](#string) | optional | @sanitize |






<a name="services-centrum-TakeDispatchResponse"></a>

### TakeDispatchResponse







<a name="services-centrum-UpdateDispatchRequest"></a>

### UpdateDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-UpdateDispatchResponse"></a>

### UpdateDispatchResponse







<a name="services-centrum-UpdateDispatchStatusRequest"></a>

### UpdateDispatchStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dispatch_id | [uint64](#uint64) |  |  |
| status | [resources.centrum.StatusDispatch](#resources-centrum-StatusDispatch) |  |  |
| reason | [string](#string) | optional | @sanitize |
| code | [string](#string) | optional | @sanitize |






<a name="services-centrum-UpdateDispatchStatusResponse"></a>

### UpdateDispatchStatusResponse







<a name="services-centrum-UpdateSettingsRequest"></a>

### UpdateSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| settings | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |






<a name="services-centrum-UpdateSettingsResponse"></a>

### UpdateSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| settings | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |






<a name="services-centrum-UpdateUnitStatusRequest"></a>

### UpdateUnitStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  |  |
| status | [resources.centrum.StatusUnit](#resources-centrum-StatusUnit) |  |  |
| reason | [string](#string) | optional | @sanitize |
| code | [string](#string) | optional | @sanitize |






<a name="services-centrum-UpdateUnitStatusResponse"></a>

### UpdateUnitStatusResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-centrum-CentrumService"></a>

### CentrumService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| UpdateSettings | [UpdateSettingsRequest](#services-centrum-UpdateSettingsRequest) | [UpdateSettingsResponse](#services-centrum-UpdateSettingsResponse) | @perm |
| CreateDispatch | [CreateDispatchRequest](#services-centrum-CreateDispatchRequest) | [CreateDispatchResponse](#services-centrum-CreateDispatchResponse) | @perm |
| UpdateDispatch | [UpdateDispatchRequest](#services-centrum-UpdateDispatchRequest) | [UpdateDispatchResponse](#services-centrum-UpdateDispatchResponse) | @perm |
| DeleteDispatch | [DeleteDispatchRequest](#services-centrum-DeleteDispatchRequest) | [DeleteDispatchResponse](#services-centrum-DeleteDispatchResponse) | @perm |
| TakeControl | [TakeControlRequest](#services-centrum-TakeControlRequest) | [TakeControlResponse](#services-centrum-TakeControlResponse) | @perm |
| AssignDispatch | [AssignDispatchRequest](#services-centrum-AssignDispatchRequest) | [AssignDispatchResponse](#services-centrum-AssignDispatchResponse) | @perm: Name=TakeControl |
| AssignUnit | [AssignUnitRequest](#services-centrum-AssignUnitRequest) | [AssignUnitResponse](#services-centrum-AssignUnitResponse) | @perm: Name=TakeControl |
| Stream | [StreamRequest](#services-centrum-StreamRequest) | [StreamResponse](#services-centrum-StreamResponse) stream | @perm |
| GetSettings | [GetSettingsRequest](#services-centrum-GetSettingsRequest) | [GetSettingsResponse](#services-centrum-GetSettingsResponse) | @perm: Name=Stream |
| JoinUnit | [JoinUnitRequest](#services-centrum-JoinUnitRequest) | [JoinUnitResponse](#services-centrum-JoinUnitResponse) | @perm: Name=Stream |
| ListUnits | [ListUnitsRequest](#services-centrum-ListUnitsRequest) | [ListUnitsResponse](#services-centrum-ListUnitsResponse) | @perm: Name=Stream |
| ListUnitActivity | [ListUnitActivityRequest](#services-centrum-ListUnitActivityRequest) | [ListUnitActivityResponse](#services-centrum-ListUnitActivityResponse) | @perm: Name=Stream |
| GetDispatch | [GetDispatchRequest](#services-centrum-GetDispatchRequest) | [GetDispatchResponse](#services-centrum-GetDispatchResponse) | @perm: Name=Stream |
| ListDispatches | [ListDispatchesRequest](#services-centrum-ListDispatchesRequest) | [ListDispatchesResponse](#services-centrum-ListDispatchesResponse) | @perm: Name=Stream |
| ListDispatchActivity | [ListDispatchActivityRequest](#services-centrum-ListDispatchActivityRequest) | [ListDispatchActivityResponse](#services-centrum-ListDispatchActivityResponse) | @perm: Name=Stream |
| CreateOrUpdateUnit | [CreateOrUpdateUnitRequest](#services-centrum-CreateOrUpdateUnitRequest) | [CreateOrUpdateUnitResponse](#services-centrum-CreateOrUpdateUnitResponse) | @perm |
| DeleteUnit | [DeleteUnitRequest](#services-centrum-DeleteUnitRequest) | [DeleteUnitResponse](#services-centrum-DeleteUnitResponse) | @perm |
| TakeDispatch | [TakeDispatchRequest](#services-centrum-TakeDispatchRequest) | [TakeDispatchResponse](#services-centrum-TakeDispatchResponse) | @perm |
| UpdateUnitStatus | [UpdateUnitStatusRequest](#services-centrum-UpdateUnitStatusRequest) | [UpdateUnitStatusResponse](#services-centrum-UpdateUnitStatusResponse) | @perm: Name=TakeDispatch |
| UpdateDispatchStatus | [UpdateDispatchStatusRequest](#services-centrum-UpdateDispatchStatusRequest) | [UpdateDispatchStatusResponse](#services-centrum-UpdateDispatchStatusResponse) | @perm: Name=TakeDispatch |

 <!-- end services -->



<a name="services_citizenstore_citizenstore-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/citizenstore/citizenstore.proto



<a name="services-citizenstore-GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| info_only | [bool](#bool) | optional |  |






<a name="services-citizenstore-GetUserResponse"></a>

### GetUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [resources.users.User](#resources-users-User) |  |  |






<a name="services-citizenstore-ListCitizensRequest"></a>

### ListCitizensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| search | [string](#string) |  | Search params |
| wanted | [bool](#bool) | optional |  |
| phone_number | [string](#string) | optional |  |
| traffic_infraction_points | [uint32](#uint32) | optional |  |
| dateofbirth | [string](#string) | optional |  |
| open_fines | [uint64](#uint64) | optional |  |






<a name="services-citizenstore-ListCitizensResponse"></a>

### ListCitizensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| users | [resources.users.User](#resources-users-User) | repeated |  |






<a name="services-citizenstore-ListUserActivityRequest"></a>

### ListUserActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| user_id | [int32](#int32) |  |  |






<a name="services-citizenstore-ListUserActivityResponse"></a>

### ListUserActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.users.UserActivity](#resources-users-UserActivity) | repeated |  |






<a name="services-citizenstore-ManageCitizenAttributesRequest"></a>

### ManageCitizenAttributesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attributes | [resources.users.CitizenAttribute](#resources-users-CitizenAttribute) | repeated |  |






<a name="services-citizenstore-ManageCitizenAttributesResponse"></a>

### ManageCitizenAttributesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attributes | [resources.users.CitizenAttribute](#resources-users-CitizenAttribute) | repeated |  |






<a name="services-citizenstore-SetProfilePictureRequest"></a>

### SetProfilePictureRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avatar | [resources.filestore.File](#resources-filestore-File) |  |  |






<a name="services-citizenstore-SetProfilePictureResponse"></a>

### SetProfilePictureResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avatar | [resources.filestore.File](#resources-filestore-File) |  |  |






<a name="services-citizenstore-SetUserPropsRequest"></a>

### SetUserPropsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| props | [resources.users.UserProps](#resources-users-UserProps) |  |  |
| reason | [string](#string) |  | @sanitize |






<a name="services-citizenstore-SetUserPropsResponse"></a>

### SetUserPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| props | [resources.users.UserProps](#resources-users-UserProps) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-citizenstore-CitizenStoreService"></a>

### CitizenStoreService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListCitizens | [ListCitizensRequest](#services-citizenstore-ListCitizensRequest) | [ListCitizensResponse](#services-citizenstore-ListCitizensResponse) | @perm: Attrs=Fields/StringList:[]string{"PhoneNumber", "Licenses", "UserProps.Wanted", "UserProps.Job", "UserProps.TrafficInfractionPoints", "UserProps.OpenFines", "UserProps.BloodType", "UserProps.MugShot", "UserProps.Attributes", "UserProps.Email"} |
| GetUser | [GetUserRequest](#services-citizenstore-GetUserRequest) | [GetUserResponse](#services-citizenstore-GetUserResponse) | @perm: Attrs=Jobs/JobGradeList |
| ListUserActivity | [ListUserActivityRequest](#services-citizenstore-ListUserActivityRequest) | [ListUserActivityResponse](#services-citizenstore-ListUserActivityResponse) | @perm: Attrs=Fields/StringList:[]string{"SourceUser", "Own"} |
| SetUserProps | [SetUserPropsRequest](#services-citizenstore-SetUserPropsRequest) | [SetUserPropsResponse](#services-citizenstore-SetUserPropsResponse) | @perm: Attrs=Fields/StringList:[]string{"Wanted", "Job", "TrafficInfractionPoints", "MugShot", "Attributes"} |
| SetProfilePicture | [SetProfilePictureRequest](#services-citizenstore-SetProfilePictureRequest) | [SetProfilePictureResponse](#services-citizenstore-SetProfilePictureResponse) | @perm: Name=Any |
| ManageCitizenAttributes | [ManageCitizenAttributesRequest](#services-citizenstore-ManageCitizenAttributesRequest) | [ManageCitizenAttributesResponse](#services-citizenstore-ManageCitizenAttributesResponse) | @perm |

 <!-- end services -->



<a name="services_completor_completor-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/completor/completor.proto



<a name="services-completor-CompleteCitizenAttributesRequest"></a>

### CompleteCitizenAttributesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| search | [string](#string) |  |  |






<a name="services-completor-CompleteCitizenAttributesResponse"></a>

### CompleteCitizenAttributesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attributes | [resources.users.CitizenAttribute](#resources-users-CitizenAttribute) | repeated |  |






<a name="services-completor-CompleteCitizensRequest"></a>

### CompleteCitizensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| search | [string](#string) |  |  |
| current_job | [bool](#bool) | optional |  |
| on_duty | [bool](#bool) | optional |  |
| user_id | [int32](#int32) | optional |  |






<a name="services-completor-CompleteCitizensRespoonse"></a>

### CompleteCitizensRespoonse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [resources.users.UserShort](#resources-users-UserShort) | repeated | @gotags: alias:"user" |






<a name="services-completor-CompleteDocumentCategoriesRequest"></a>

### CompleteDocumentCategoriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| search | [string](#string) |  |  |






<a name="services-completor-CompleteDocumentCategoriesResponse"></a>

### CompleteDocumentCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| categories | [resources.documents.Category](#resources-documents-Category) | repeated |  |






<a name="services-completor-CompleteJobsRequest"></a>

### CompleteJobsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| search | [string](#string) | optional |  |
| exact_match | [bool](#bool) | optional |  |
| current_job | [bool](#bool) | optional |  |






<a name="services-completor-CompleteJobsResponse"></a>

### CompleteJobsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [resources.users.Job](#resources-users-Job) | repeated |  |






<a name="services-completor-ListLawBooksRequest"></a>

### ListLawBooksRequest







<a name="services-completor-ListLawBooksResponse"></a>

### ListLawBooksResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| books | [resources.laws.LawBook](#resources-laws-LawBook) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-completor-CompletorService"></a>

### CompletorService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CompleteCitizens | [CompleteCitizensRequest](#services-completor-CompleteCitizensRequest) | [CompleteCitizensRespoonse](#services-completor-CompleteCitizensRespoonse) | @perm |
| CompleteJobs | [CompleteJobsRequest](#services-completor-CompleteJobsRequest) | [CompleteJobsResponse](#services-completor-CompleteJobsResponse) | @perm: Name=Any |
| CompleteDocumentCategories | [CompleteDocumentCategoriesRequest](#services-completor-CompleteDocumentCategoriesRequest) | [CompleteDocumentCategoriesResponse](#services-completor-CompleteDocumentCategoriesResponse) | @perm: Attrs=Jobs/JobList |
| ListLawBooks | [ListLawBooksRequest](#services-completor-ListLawBooksRequest) | [ListLawBooksResponse](#services-completor-ListLawBooksResponse) | @perm: Name=Any |
| CompleteCitizenAttributes | [CompleteCitizenAttributesRequest](#services-completor-CompleteCitizenAttributesRequest) | [CompleteCitizenAttributesResponse](#services-completor-CompleteCitizenAttributesResponse) | @perm: Attrs=Jobs/JobList |

 <!-- end services -->



<a name="services_dmv_vehicles-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/dmv/vehicles.proto



<a name="services-dmv-ListVehiclesRequest"></a>

### ListVehiclesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| license_plate | [string](#string) | optional | Search params |
| model | [string](#string) | optional |  |
| user_id | [int32](#int32) | optional |  |






<a name="services-dmv-ListVehiclesResponse"></a>

### ListVehiclesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| vehicles | [resources.vehicles.Vehicle](#resources-vehicles-Vehicle) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-dmv-DMVService"></a>

### DMVService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListVehicles | [ListVehiclesRequest](#services-dmv-ListVehiclesRequest) | [ListVehiclesResponse](#services-dmv-ListVehiclesResponse) | @perm |

 <!-- end services -->



<a name="services_docstore_docstore-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/docstore/docstore.proto



<a name="services-docstore-AddDocumentReferenceRequest"></a>

### AddDocumentReferenceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reference | [resources.documents.DocumentReference](#resources-documents-DocumentReference) |  |  |






<a name="services-docstore-AddDocumentReferenceResponse"></a>

### AddDocumentReferenceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-AddDocumentRelationRequest"></a>

### AddDocumentRelationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| relation | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) |  |  |






<a name="services-docstore-AddDocumentRelationResponse"></a>

### AddDocumentRelationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-ChangeDocumentOwnerRequest"></a>

### ChangeDocumentOwnerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| new_user_id | [int32](#int32) | optional |  |






<a name="services-docstore-ChangeDocumentOwnerResponse"></a>

### ChangeDocumentOwnerResponse







<a name="services-docstore-CreateCategoryRequest"></a>

### CreateCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [resources.documents.Category](#resources-documents-Category) |  |  |






<a name="services-docstore-CreateCategoryResponse"></a>

### CreateCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-CreateDocumentReqRequest"></a>

### CreateDocumentReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| request_type | [resources.documents.DocActivityType](#resources-documents-DocActivityType) |  |  |
| reason | [string](#string) | optional | @sanitize |
| data | [resources.documents.DocActivityData](#resources-documents-DocActivityData) | optional |  |






<a name="services-docstore-CreateDocumentReqResponse"></a>

### CreateDocumentReqResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [resources.documents.DocRequest](#resources-documents-DocRequest) |  |  |






<a name="services-docstore-CreateDocumentRequest"></a>

### CreateDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [uint64](#uint64) | optional | @gotags: alias:"category_id" |
| title | [string](#string) |  | @sanitize: method=StripTags

@gotags: alias:"title" |
| content | [resources.common.content.Content](#resources-common-content-Content) |  | @sanitize |
| content_type | [resources.common.content.ContentType](#resources-common-content-ContentType) |  | @gotags: alias:"content_type" |
| data | [string](#string) | optional | @gotags: alias:"data" |
| state | [string](#string) |  | @sanitize

@gotags: alias:"state" |
| closed | [bool](#bool) |  | @gotags: alias:"closed" |
| public | [bool](#bool) |  | @gotags: alias:"public" |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) | optional |  |
| template_id | [uint64](#uint64) | optional |  |






<a name="services-docstore-CreateDocumentResponse"></a>

### CreateDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  | @gotags: alias:"id" |






<a name="services-docstore-CreateTemplateRequest"></a>

### CreateTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.documents.Template](#resources-documents-Template) |  |  |






<a name="services-docstore-CreateTemplateResponse"></a>

### CreateTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-DeleteCategoryRequest"></a>

### DeleteCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ids | [uint64](#uint64) | repeated |  |






<a name="services-docstore-DeleteCategoryResponse"></a>

### DeleteCategoryResponse







<a name="services-docstore-DeleteCommentRequest"></a>

### DeleteCommentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| comment_id | [uint64](#uint64) |  |  |






<a name="services-docstore-DeleteCommentResponse"></a>

### DeleteCommentResponse







<a name="services-docstore-DeleteDocumentReqRequest"></a>

### DeleteDocumentReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [uint64](#uint64) |  |  |






<a name="services-docstore-DeleteDocumentReqResponse"></a>

### DeleteDocumentReqResponse







<a name="services-docstore-DeleteDocumentRequest"></a>

### DeleteDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  | @gotags: alias:"id" |






<a name="services-docstore-DeleteDocumentResponse"></a>

### DeleteDocumentResponse







<a name="services-docstore-DeleteTemplateRequest"></a>

### DeleteTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-DeleteTemplateResponse"></a>

### DeleteTemplateResponse







<a name="services-docstore-EditCommentRequest"></a>

### EditCommentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| comment | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-docstore-EditCommentResponse"></a>

### EditCommentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| comment | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-docstore-GetCommentsRequest"></a>

### GetCommentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| document_id | [uint64](#uint64) |  |  |






<a name="services-docstore-GetCommentsResponse"></a>

### GetCommentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| comments | [resources.documents.Comment](#resources-documents-Comment) | repeated |  |






<a name="services-docstore-GetDocumentAccessRequest"></a>

### GetDocumentAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |






<a name="services-docstore-GetDocumentAccessResponse"></a>

### GetDocumentAccessResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-docstore-GetDocumentReferencesRequest"></a>

### GetDocumentReferencesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |






<a name="services-docstore-GetDocumentReferencesResponse"></a>

### GetDocumentReferencesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| references | [resources.documents.DocumentReference](#resources-documents-DocumentReference) | repeated | @gotags: alias:"reference" |






<a name="services-docstore-GetDocumentRelationsRequest"></a>

### GetDocumentRelationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |






<a name="services-docstore-GetDocumentRelationsResponse"></a>

### GetDocumentRelationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| relations | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) | repeated | @gotags: alias:"relation" |






<a name="services-docstore-GetDocumentRequest"></a>

### GetDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| info_only | [bool](#bool) | optional |  |






<a name="services-docstore-GetDocumentResponse"></a>

### GetDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document | [resources.documents.Document](#resources-documents-Document) |  |  |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-docstore-GetTemplateRequest"></a>

### GetTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template_id | [uint64](#uint64) |  |  |
| data | [resources.documents.TemplateData](#resources-documents-TemplateData) | optional |  |
| render | [bool](#bool) | optional |  |






<a name="services-docstore-GetTemplateResponse"></a>

### GetTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.documents.Template](#resources-documents-Template) |  |  |
| rendered | [bool](#bool) |  |  |






<a name="services-docstore-ListCategoriesRequest"></a>

### ListCategoriesRequest







<a name="services-docstore-ListCategoriesResponse"></a>

### ListCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [resources.documents.Category](#resources-documents-Category) | repeated |  |






<a name="services-docstore-ListDocumentActivityRequest"></a>

### ListDocumentActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| document_id | [uint64](#uint64) |  |  |
| activity_types | [resources.documents.DocActivityType](#resources-documents-DocActivityType) | repeated | Search params |






<a name="services-docstore-ListDocumentActivityResponse"></a>

### ListDocumentActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.documents.DocActivity](#resources-documents-DocActivity) | repeated |  |






<a name="services-docstore-ListDocumentPinsRequest"></a>

### ListDocumentPinsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |






<a name="services-docstore-ListDocumentPinsResponse"></a>

### ListDocumentPinsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| documents | [resources.documents.DocumentShort](#resources-documents-DocumentShort) | repeated |  |






<a name="services-docstore-ListDocumentReqsRequest"></a>

### ListDocumentReqsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| document_id | [uint64](#uint64) |  |  |






<a name="services-docstore-ListDocumentReqsResponse"></a>

### ListDocumentReqsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| requests | [resources.documents.DocRequest](#resources-documents-DocRequest) | repeated |  |






<a name="services-docstore-ListDocumentsRequest"></a>

### ListDocumentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| search | [string](#string) | optional | Search params |
| category_ids | [uint64](#uint64) | repeated |  |
| creator_ids | [int32](#int32) | repeated |  |
| from | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| to | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| closed | [bool](#bool) | optional |  |
| document_ids | [uint64](#uint64) | repeated |  |






<a name="services-docstore-ListDocumentsResponse"></a>

### ListDocumentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| documents | [resources.documents.DocumentShort](#resources-documents-DocumentShort) | repeated |  |






<a name="services-docstore-ListTemplatesRequest"></a>

### ListTemplatesRequest







<a name="services-docstore-ListTemplatesResponse"></a>

### ListTemplatesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| templates | [resources.documents.TemplateShort](#resources-documents-TemplateShort) | repeated |  |






<a name="services-docstore-ListUserDocumentsRequest"></a>

### ListUserDocumentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| user_id | [int32](#int32) |  |  |
| relations | [resources.documents.DocRelation](#resources-documents-DocRelation) | repeated |  |
| closed | [bool](#bool) | optional |  |






<a name="services-docstore-ListUserDocumentsResponse"></a>

### ListUserDocumentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| relations | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) | repeated |  |






<a name="services-docstore-PostCommentRequest"></a>

### PostCommentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| comment | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-docstore-PostCommentResponse"></a>

### PostCommentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| comment | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-docstore-RemoveDocumentReferenceRequest"></a>

### RemoveDocumentReferenceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-RemoveDocumentReferenceResponse"></a>

### RemoveDocumentReferenceResponse







<a name="services-docstore-RemoveDocumentRelationRequest"></a>

### RemoveDocumentRelationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-docstore-RemoveDocumentRelationResponse"></a>

### RemoveDocumentRelationResponse







<a name="services-docstore-SetDocumentAccessRequest"></a>

### SetDocumentAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-docstore-SetDocumentAccessResponse"></a>

### SetDocumentAccessResponse







<a name="services-docstore-SetDocumentReminderRequest"></a>

### SetDocumentReminderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| reminder_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| message | [string](#string) | optional | @sanitize: method=StripTags |






<a name="services-docstore-SetDocumentReminderResponse"></a>

### SetDocumentReminderResponse







<a name="services-docstore-ToggleDocumentPinRequest"></a>

### ToggleDocumentPinRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| state | [bool](#bool) |  |  |






<a name="services-docstore-ToggleDocumentPinResponse"></a>

### ToggleDocumentPinResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [bool](#bool) |  |  |






<a name="services-docstore-ToggleDocumentRequest"></a>

### ToggleDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| closed | [bool](#bool) |  |  |






<a name="services-docstore-ToggleDocumentResponse"></a>

### ToggleDocumentResponse







<a name="services-docstore-UpdateCategoryRequest"></a>

### UpdateCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [resources.documents.Category](#resources-documents-Category) |  |  |






<a name="services-docstore-UpdateCategoryResponse"></a>

### UpdateCategoryResponse







<a name="services-docstore-UpdateDocumentReqRequest"></a>

### UpdateDocumentReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |
| request_id | [uint64](#uint64) |  |  |
| reason | [string](#string) | optional | @sanitize |
| data | [resources.documents.DocActivityData](#resources-documents-DocActivityData) | optional |  |
| accepted | [bool](#bool) |  |  |






<a name="services-docstore-UpdateDocumentReqResponse"></a>

### UpdateDocumentReqResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [resources.documents.DocRequest](#resources-documents-DocRequest) |  |  |






<a name="services-docstore-UpdateDocumentRequest"></a>

### UpdateDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  | @gotags: alias:"id" |
| category_id | [uint64](#uint64) | optional | @gotags: alias:"category_id" |
| title | [string](#string) |  | @sanitize: method=StripTags

@gotags: alias:"title" |
| content | [resources.common.content.Content](#resources-common-content-Content) |  | @sanitize |
| content_type | [resources.common.content.ContentType](#resources-common-content-ContentType) |  | @gotags: alias:"content_type" |
| data | [string](#string) | optional | @gotags: alias:"data" |
| state | [string](#string) |  | @sanitize

@gotags: alias:"state" |
| closed | [bool](#bool) |  | @gotags: alias:"closed" |
| public | [bool](#bool) |  | @gotags: alias:"public" |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) | optional |  |






<a name="services-docstore-UpdateDocumentResponse"></a>

### UpdateDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  | @gotags: alias:"id" |






<a name="services-docstore-UpdateTemplateRequest"></a>

### UpdateTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.documents.Template](#resources-documents-Template) |  |  |






<a name="services-docstore-UpdateTemplateResponse"></a>

### UpdateTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.documents.Template](#resources-documents-Template) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-docstore-DocStoreService"></a>

### DocStoreService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListTemplates | [ListTemplatesRequest](#services-docstore-ListTemplatesRequest) | [ListTemplatesResponse](#services-docstore-ListTemplatesResponse) | @perm |
| GetTemplate | [GetTemplateRequest](#services-docstore-GetTemplateRequest) | [GetTemplateResponse](#services-docstore-GetTemplateResponse) | @perm: Name=ListTemplates |
| CreateTemplate | [CreateTemplateRequest](#services-docstore-CreateTemplateRequest) | [CreateTemplateResponse](#services-docstore-CreateTemplateResponse) | @perm |
| UpdateTemplate | [UpdateTemplateRequest](#services-docstore-UpdateTemplateRequest) | [UpdateTemplateResponse](#services-docstore-UpdateTemplateResponse) | @perm: Name=CreateTemplate |
| DeleteTemplate | [DeleteTemplateRequest](#services-docstore-DeleteTemplateRequest) | [DeleteTemplateResponse](#services-docstore-DeleteTemplateResponse) | @perm |
| ListDocuments | [ListDocumentsRequest](#services-docstore-ListDocumentsRequest) | [ListDocumentsResponse](#services-docstore-ListDocumentsResponse) | @perm |
| GetDocument | [GetDocumentRequest](#services-docstore-GetDocumentRequest) | [GetDocumentResponse](#services-docstore-GetDocumentResponse) | @perm: Name=ListDocuments |
| CreateDocument | [CreateDocumentRequest](#services-docstore-CreateDocumentRequest) | [CreateDocumentResponse](#services-docstore-CreateDocumentResponse) | @perm |
| UpdateDocument | [UpdateDocumentRequest](#services-docstore-UpdateDocumentRequest) | [UpdateDocumentResponse](#services-docstore-UpdateDocumentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| DeleteDocument | [DeleteDocumentRequest](#services-docstore-DeleteDocumentRequest) | [DeleteDocumentResponse](#services-docstore-DeleteDocumentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| ToggleDocument | [ToggleDocumentRequest](#services-docstore-ToggleDocumentRequest) | [ToggleDocumentResponse](#services-docstore-ToggleDocumentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| ChangeDocumentOwner | [ChangeDocumentOwnerRequest](#services-docstore-ChangeDocumentOwnerRequest) | [ChangeDocumentOwnerResponse](#services-docstore-ChangeDocumentOwnerResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| GetDocumentReferences | [GetDocumentReferencesRequest](#services-docstore-GetDocumentReferencesRequest) | [GetDocumentReferencesResponse](#services-docstore-GetDocumentReferencesResponse) | @perm: Name=ListDocuments |
| GetDocumentRelations | [GetDocumentRelationsRequest](#services-docstore-GetDocumentRelationsRequest) | [GetDocumentRelationsResponse](#services-docstore-GetDocumentRelationsResponse) | @perm: Name=ListDocuments |
| AddDocumentReference | [AddDocumentReferenceRequest](#services-docstore-AddDocumentReferenceRequest) | [AddDocumentReferenceResponse](#services-docstore-AddDocumentReferenceResponse) | @perm |
| RemoveDocumentReference | [RemoveDocumentReferenceRequest](#services-docstore-RemoveDocumentReferenceRequest) | [RemoveDocumentReferenceResponse](#services-docstore-RemoveDocumentReferenceResponse) | @perm: Name=AddDocumentReference |
| AddDocumentRelation | [AddDocumentRelationRequest](#services-docstore-AddDocumentRelationRequest) | [AddDocumentRelationResponse](#services-docstore-AddDocumentRelationResponse) | @perm |
| RemoveDocumentRelation | [RemoveDocumentRelationRequest](#services-docstore-RemoveDocumentRelationRequest) | [RemoveDocumentRelationResponse](#services-docstore-RemoveDocumentRelationResponse) | @perm: Name=AddDocumentRelation |
| GetComments | [GetCommentsRequest](#services-docstore-GetCommentsRequest) | [GetCommentsResponse](#services-docstore-GetCommentsResponse) | @perm: Name=ListDocuments |
| PostComment | [PostCommentRequest](#services-docstore-PostCommentRequest) | [PostCommentResponse](#services-docstore-PostCommentResponse) | @perm |
| EditComment | [EditCommentRequest](#services-docstore-EditCommentRequest) | [EditCommentResponse](#services-docstore-EditCommentResponse) | @perm: Name=PostComment |
| DeleteComment | [DeleteCommentRequest](#services-docstore-DeleteCommentRequest) | [DeleteCommentResponse](#services-docstore-DeleteCommentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| GetDocumentAccess | [GetDocumentAccessRequest](#services-docstore-GetDocumentAccessRequest) | [GetDocumentAccessResponse](#services-docstore-GetDocumentAccessResponse) | @perm: Name=ListDocuments |
| SetDocumentAccess | [SetDocumentAccessRequest](#services-docstore-SetDocumentAccessRequest) | [SetDocumentAccessResponse](#services-docstore-SetDocumentAccessResponse) | @perm: Name=CreateDocument |
| ListDocumentActivity | [ListDocumentActivityRequest](#services-docstore-ListDocumentActivityRequest) | [ListDocumentActivityResponse](#services-docstore-ListDocumentActivityResponse) | @perm |
| ListDocumentReqs | [ListDocumentReqsRequest](#services-docstore-ListDocumentReqsRequest) | [ListDocumentReqsResponse](#services-docstore-ListDocumentReqsResponse) | @perm |
| CreateDocumentReq | [CreateDocumentReqRequest](#services-docstore-CreateDocumentReqRequest) | [CreateDocumentReqResponse](#services-docstore-CreateDocumentReqResponse) | @perm: Attrs=Types/StringList:[]string{"Access", "Closure", "Update", "Deletion", "OwnerChange"} |
| UpdateDocumentReq | [UpdateDocumentReqRequest](#services-docstore-UpdateDocumentReqRequest) | [UpdateDocumentReqResponse](#services-docstore-UpdateDocumentReqResponse) | @perm: Name=CreateDocumentReq |
| DeleteDocumentReq | [DeleteDocumentReqRequest](#services-docstore-DeleteDocumentReqRequest) | [DeleteDocumentReqResponse](#services-docstore-DeleteDocumentReqResponse) | @perm |
| ListUserDocuments | [ListUserDocumentsRequest](#services-docstore-ListUserDocumentsRequest) | [ListUserDocumentsResponse](#services-docstore-ListUserDocumentsResponse) | @perm |
| ListCategories | [ListCategoriesRequest](#services-docstore-ListCategoriesRequest) | [ListCategoriesResponse](#services-docstore-ListCategoriesResponse) | @perm |
| CreateCategory | [CreateCategoryRequest](#services-docstore-CreateCategoryRequest) | [CreateCategoryResponse](#services-docstore-CreateCategoryResponse) | @perm |
| UpdateCategory | [UpdateCategoryRequest](#services-docstore-UpdateCategoryRequest) | [UpdateCategoryResponse](#services-docstore-UpdateCategoryResponse) | @perm: Name=CreateCategory |
| DeleteCategory | [DeleteCategoryRequest](#services-docstore-DeleteCategoryRequest) | [DeleteCategoryResponse](#services-docstore-DeleteCategoryResponse) | @perm |
| ListDocumentPins | [ListDocumentPinsRequest](#services-docstore-ListDocumentPinsRequest) | [ListDocumentPinsResponse](#services-docstore-ListDocumentPinsResponse) | @perm: Name=ListDocuments |
| ToggleDocumentPin | [ToggleDocumentPinRequest](#services-docstore-ToggleDocumentPinRequest) | [ToggleDocumentPinResponse](#services-docstore-ToggleDocumentPinResponse) | @perm |
| SetDocumentReminder | [SetDocumentReminderRequest](#services-docstore-SetDocumentReminderRequest) | [SetDocumentReminderResponse](#services-docstore-SetDocumentReminderResponse) | @perm |

 <!-- end services -->



<a name="services_jobs_conduct-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/conduct.proto



<a name="services-jobs-CreateConductEntryRequest"></a>

### CreateConductEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |






<a name="services-jobs-CreateConductEntryResponse"></a>

### CreateConductEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |






<a name="services-jobs-DeleteConductEntryRequest"></a>

### DeleteConductEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-jobs-DeleteConductEntryResponse"></a>

### DeleteConductEntryResponse







<a name="services-jobs-ListConductEntriesRequest"></a>

### ListConductEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| types | [resources.jobs.ConductType](#resources-jobs-ConductType) | repeated | Search params |
| show_expired | [bool](#bool) | optional |  |
| user_ids | [int32](#int32) | repeated |  |
| ids | [uint64](#uint64) | repeated |  |






<a name="services-jobs-ListConductEntriesResponse"></a>

### ListConductEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| entries | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) | repeated |  |






<a name="services-jobs-UpdateConductEntryRequest"></a>

### UpdateConductEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |






<a name="services-jobs-UpdateConductEntryResponse"></a>

### UpdateConductEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-jobs-JobsConductService"></a>

### JobsConductService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListConductEntries | [ListConductEntriesRequest](#services-jobs-ListConductEntriesRequest) | [ListConductEntriesResponse](#services-jobs-ListConductEntriesResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "All"} |
| CreateConductEntry | [CreateConductEntryRequest](#services-jobs-CreateConductEntryRequest) | [CreateConductEntryResponse](#services-jobs-CreateConductEntryResponse) | @perm |
| UpdateConductEntry | [UpdateConductEntryRequest](#services-jobs-UpdateConductEntryRequest) | [UpdateConductEntryResponse](#services-jobs-UpdateConductEntryResponse) | @perm |
| DeleteConductEntry | [DeleteConductEntryRequest](#services-jobs-DeleteConductEntryRequest) | [DeleteConductEntryResponse](#services-jobs-DeleteConductEntryResponse) | @perm |

 <!-- end services -->



<a name="services_jobs_jobs-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/jobs.proto



<a name="services-jobs-GetColleagueLabelsRequest"></a>

### GetColleagueLabelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| search | [string](#string) | optional |  |






<a name="services-jobs-GetColleagueLabelsResponse"></a>

### GetColleagueLabelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| labels | [resources.jobs.Label](#resources-jobs-Label) | repeated |  |






<a name="services-jobs-GetColleagueLabelsStatsRequest"></a>

### GetColleagueLabelsStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| label_ids | [uint64](#uint64) | repeated |  |






<a name="services-jobs-GetColleagueLabelsStatsResponse"></a>

### GetColleagueLabelsStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [resources.jobs.LabelCount](#resources-jobs-LabelCount) | repeated |  |






<a name="services-jobs-GetColleagueRequest"></a>

### GetColleagueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| info_only | [bool](#bool) | optional |  |






<a name="services-jobs-GetColleagueResponse"></a>

### GetColleagueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| colleague | [resources.jobs.Colleague](#resources-jobs-Colleague) |  |  |






<a name="services-jobs-GetMOTDRequest"></a>

### GetMOTDRequest







<a name="services-jobs-GetMOTDResponse"></a>

### GetMOTDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| motd | [string](#string) |  |  |






<a name="services-jobs-GetSelfRequest"></a>

### GetSelfRequest







<a name="services-jobs-GetSelfResponse"></a>

### GetSelfResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| colleague | [resources.jobs.Colleague](#resources-jobs-Colleague) |  |  |






<a name="services-jobs-ListColleagueActivityRequest"></a>

### ListColleagueActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| user_ids | [int32](#int32) | repeated | Search params |
| activity_types | [resources.jobs.JobsUserActivityType](#resources-jobs-JobsUserActivityType) | repeated |  |






<a name="services-jobs-ListColleagueActivityResponse"></a>

### ListColleagueActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.jobs.JobsUserActivity](#resources-jobs-JobsUserActivity) | repeated |  |






<a name="services-jobs-ListColleaguesRequest"></a>

### ListColleaguesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| search | [string](#string) |  | Search params |
| user_id | [int32](#int32) | optional |  |
| absent | [bool](#bool) | optional |  |
| label_ids | [uint64](#uint64) | repeated |  |
| name_prefix | [string](#string) | optional |  |
| name_suffix | [string](#string) | optional |  |






<a name="services-jobs-ListColleaguesResponse"></a>

### ListColleaguesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| colleagues | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |






<a name="services-jobs-ManageColleagueLabelsRequest"></a>

### ManageColleagueLabelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| labels | [resources.jobs.Label](#resources-jobs-Label) | repeated |  |






<a name="services-jobs-ManageColleagueLabelsResponse"></a>

### ManageColleagueLabelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| labels | [resources.jobs.Label](#resources-jobs-Label) | repeated |  |






<a name="services-jobs-SetJobsUserPropsRequest"></a>

### SetJobsUserPropsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| props | [resources.jobs.JobsUserProps](#resources-jobs-JobsUserProps) |  |  |
| reason | [string](#string) |  | @sanitize |






<a name="services-jobs-SetJobsUserPropsResponse"></a>

### SetJobsUserPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| props | [resources.jobs.JobsUserProps](#resources-jobs-JobsUserProps) |  |  |






<a name="services-jobs-SetMOTDRequest"></a>

### SetMOTDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| motd | [string](#string) |  | @sanitize: method=StripTags |






<a name="services-jobs-SetMOTDResponse"></a>

### SetMOTDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| motd | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-jobs-JobsService"></a>

### JobsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListColleagues | [ListColleaguesRequest](#services-jobs-ListColleaguesRequest) | [ListColleaguesResponse](#services-jobs-ListColleaguesResponse) | @perm |
| GetSelf | [GetSelfRequest](#services-jobs-GetSelfRequest) | [GetSelfResponse](#services-jobs-GetSelfResponse) | @perm: Name=ListColleagues |
| GetColleague | [GetColleagueRequest](#services-jobs-GetColleagueRequest) | [GetColleagueResponse](#services-jobs-GetColleagueResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"Note", "Labels"} |
| ListColleagueActivity | [ListColleagueActivityRequest](#services-jobs-ListColleagueActivityRequest) | [ListColleagueActivityResponse](#services-jobs-ListColleagueActivityResponse) | @perm: Attrs=Types/StringList:[]string{"HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME"} |
| SetJobsUserProps | [SetJobsUserPropsRequest](#services-jobs-SetJobsUserPropsRequest) | [SetJobsUserPropsResponse](#services-jobs-SetJobsUserPropsResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"AbsenceDate", "Note", "Labels", "Name"} |
| GetColleagueLabels | [GetColleagueLabelsRequest](#services-jobs-GetColleagueLabelsRequest) | [GetColleagueLabelsResponse](#services-jobs-GetColleagueLabelsResponse) | @perm: Name=GetColleague |
| ManageColleagueLabels | [ManageColleagueLabelsRequest](#services-jobs-ManageColleagueLabelsRequest) | [ManageColleagueLabelsResponse](#services-jobs-ManageColleagueLabelsResponse) | @perm |
| GetColleagueLabelsStats | [GetColleagueLabelsStatsRequest](#services-jobs-GetColleagueLabelsStatsRequest) | [GetColleagueLabelsStatsResponse](#services-jobs-GetColleagueLabelsStatsResponse) | @perm: Name=GetColleague |
| GetMOTD | [GetMOTDRequest](#services-jobs-GetMOTDRequest) | [GetMOTDResponse](#services-jobs-GetMOTDResponse) | @perm: Name=Any |
| SetMOTD | [SetMOTDRequest](#services-jobs-SetMOTDRequest) | [SetMOTDResponse](#services-jobs-SetMOTDResponse) | @perm |

 <!-- end services -->



<a name="services_jobs_timeclock-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/timeclock.proto



<a name="services-jobs-GetTimeclockStatsRequest"></a>

### GetTimeclockStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) | optional |  |






<a name="services-jobs-GetTimeclockStatsResponse"></a>

### GetTimeclockStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stats | [resources.jobs.TimeclockStats](#resources-jobs-TimeclockStats) |  |  |
| weekly | [resources.jobs.TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats) | repeated |  |






<a name="services-jobs-ListInactiveEmployeesRequest"></a>

### ListInactiveEmployeesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| days | [int32](#int32) |  | Search params |






<a name="services-jobs-ListInactiveEmployeesResponse"></a>

### ListInactiveEmployeesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| colleagues | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |






<a name="services-jobs-ListTimeclockRequest"></a>

### ListTimeclockRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| user_mode | [resources.jobs.TimeclockUserMode](#resources-jobs-TimeclockUserMode) |  | Search params |
| mode | [resources.jobs.TimeclockMode](#resources-jobs-TimeclockMode) |  |  |
| date | [resources.common.database.DateRange](#resources-common-database-DateRange) | optional |  |
| per_day | [bool](#bool) |  |  |
| user_ids | [int32](#int32) | repeated |  |






<a name="services-jobs-ListTimeclockResponse"></a>

### ListTimeclockResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| stats | [resources.jobs.TimeclockStats](#resources-jobs-TimeclockStats) |  |  |
| stats_weekly | [resources.jobs.TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats) | repeated |  |
| daily | [TimeclockDay](#services-jobs-TimeclockDay) |  |  |
| weekly | [TimeclockWeekly](#services-jobs-TimeclockWeekly) |  |  |
| range | [TimeclockRange](#services-jobs-TimeclockRange) |  |  |






<a name="services-jobs-TimeclockDay"></a>

### TimeclockDay



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| date | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| entries | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| sum | [float](#float) |  |  |






<a name="services-jobs-TimeclockRange"></a>

### TimeclockRange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| date | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: sql:"primary_key" |
| entries | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| sum | [float](#float) |  |  |






<a name="services-jobs-TimeclockWeekly"></a>

### TimeclockWeekly



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| date | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: sql:"primary_key" |
| entries | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| sum | [float](#float) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-jobs-JobsTimeclockService"></a>

### JobsTimeclockService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListTimeclock | [ListTimeclockRequest](#services-jobs-ListTimeclockRequest) | [ListTimeclockResponse](#services-jobs-ListTimeclockResponse) | @perm: Attrs=Access/StringList:[]string{"All"} |
| GetTimeclockStats | [GetTimeclockStatsRequest](#services-jobs-GetTimeclockStatsRequest) | [GetTimeclockStatsResponse](#services-jobs-GetTimeclockStatsResponse) | @perm: Name=ListTimeclock |
| ListInactiveEmployees | [ListInactiveEmployeesRequest](#services-jobs-ListInactiveEmployeesRequest) | [ListInactiveEmployeesResponse](#services-jobs-ListInactiveEmployeesResponse) | @perm |

 <!-- end services -->



<a name="services_livemapper_livemap-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/livemapper/livemap.proto



<a name="services-livemapper-CreateOrUpdateMarkerRequest"></a>

### CreateOrUpdateMarkerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| marker | [resources.livemap.MarkerMarker](#resources-livemap-MarkerMarker) |  |  |






<a name="services-livemapper-CreateOrUpdateMarkerResponse"></a>

### CreateOrUpdateMarkerResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| marker | [resources.livemap.MarkerMarker](#resources-livemap-MarkerMarker) |  |  |






<a name="services-livemapper-DeleteMarkerRequest"></a>

### DeleteMarkerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-livemapper-DeleteMarkerResponse"></a>

### DeleteMarkerResponse







<a name="services-livemapper-JobsList"></a>

### JobsList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [resources.users.Job](#resources-users-Job) | repeated |  |
| markers | [resources.users.Job](#resources-users-Job) | repeated |  |






<a name="services-livemapper-MarkerMarkersUpdates"></a>

### MarkerMarkersUpdates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| markers | [resources.livemap.MarkerMarker](#resources-livemap-MarkerMarker) | repeated |  |






<a name="services-livemapper-StreamRequest"></a>

### StreamRequest







<a name="services-livemapper-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [JobsList](#services-livemapper-JobsList) |  |  |
| markers | [MarkerMarkersUpdates](#services-livemapper-MarkerMarkersUpdates) |  |  |
| users | [UserMarkersUpdates](#services-livemapper-UserMarkersUpdates) |  |  |
| user_on_duty | [bool](#bool) | optional |  |






<a name="services-livemapper-UserMarkersUpdates"></a>

### UserMarkersUpdates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [resources.livemap.UserMarker](#resources-livemap-UserMarker) | repeated |  |
| part | [int32](#int32) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-livemapper-LivemapperService"></a>

### LivemapperService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Stream | [StreamRequest](#services-livemapper-StreamRequest) | [StreamResponse](#services-livemapper-StreamResponse) stream | @perm: Attrs=Markers/JobList|Players/JobGradeList |
| CreateOrUpdateMarker | [CreateOrUpdateMarkerRequest](#services-livemapper-CreateOrUpdateMarkerRequest) | [CreateOrUpdateMarkerResponse](#services-livemapper-CreateOrUpdateMarkerResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| DeleteMarker | [DeleteMarkerRequest](#services-livemapper-DeleteMarkerRequest) | [DeleteMarkerResponse](#services-livemapper-DeleteMarkerResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |

 <!-- end services -->



<a name="services_notificator_notificator-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/notificator/notificator.proto



<a name="services-notificator-GetNotificationsRequest"></a>

### GetNotificationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| include_read | [bool](#bool) | optional |  |
| categories | [resources.notifications.NotificationCategory](#resources-notifications-NotificationCategory) | repeated |  |






<a name="services-notificator-GetNotificationsResponse"></a>

### GetNotificationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| notifications | [resources.notifications.Notification](#resources-notifications-Notification) | repeated |  |






<a name="services-notificator-MarkNotificationsRequest"></a>

### MarkNotificationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ids | [uint64](#uint64) | repeated |  |
| all | [bool](#bool) | optional |  |






<a name="services-notificator-MarkNotificationsResponse"></a>

### MarkNotificationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| updated | [uint64](#uint64) |  |  |






<a name="services-notificator-StreamRequest"></a>

### StreamRequest







<a name="services-notificator-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| notification_count | [int32](#int32) |  |  |
| restart | [bool](#bool) | optional |  |
| user_event | [resources.notifications.UserEvent](#resources-notifications-UserEvent) |  |  |
| job_event | [resources.notifications.JobEvent](#resources-notifications-JobEvent) |  |  |
| job_grade_event | [resources.notifications.JobGradeEvent](#resources-notifications-JobGradeEvent) |  |  |
| system_event | [resources.notifications.SystemEvent](#resources-notifications-SystemEvent) |  |  |
| mailer_event | [resources.mailer.MailerEvent](#resources-mailer-MailerEvent) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-notificator-NotificatorService"></a>

### NotificatorService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetNotifications | [GetNotificationsRequest](#services-notificator-GetNotificationsRequest) | [GetNotificationsResponse](#services-notificator-GetNotificationsResponse) | @perm: Name=Any |
| MarkNotifications | [MarkNotificationsRequest](#services-notificator-MarkNotificationsRequest) | [MarkNotificationsResponse](#services-notificator-MarkNotificationsResponse) | @perm: Name=Any |
| Stream | [StreamRequest](#services-notificator-StreamRequest) | [StreamResponse](#services-notificator-StreamResponse) stream | @perm: Name=Any |

 <!-- end services -->



<a name="services_qualifications_qualifications-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/qualifications/qualifications.proto



<a name="services-qualifications-CreateOrUpdateQualificationRequestRequest"></a>

### CreateOrUpdateQualificationRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [resources.qualifications.QualificationRequest](#resources-qualifications-QualificationRequest) |  |  |






<a name="services-qualifications-CreateOrUpdateQualificationRequestResponse"></a>

### CreateOrUpdateQualificationRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [resources.qualifications.QualificationRequest](#resources-qualifications-QualificationRequest) |  |  |






<a name="services-qualifications-CreateOrUpdateQualificationResultRequest"></a>

### CreateOrUpdateQualificationResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [resources.qualifications.QualificationResult](#resources-qualifications-QualificationResult) |  |  |






<a name="services-qualifications-CreateOrUpdateQualificationResultResponse"></a>

### CreateOrUpdateQualificationResultResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result | [resources.qualifications.QualificationResult](#resources-qualifications-QualificationResult) |  |  |






<a name="services-qualifications-CreateQualificationRequest"></a>

### CreateQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification | [resources.qualifications.Qualification](#resources-qualifications-Qualification) |  |  |






<a name="services-qualifications-CreateQualificationResponse"></a>

### CreateQualificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |






<a name="services-qualifications-DeleteQualificationReqRequest"></a>

### DeleteQualificationReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |






<a name="services-qualifications-DeleteQualificationReqResponse"></a>

### DeleteQualificationReqResponse







<a name="services-qualifications-DeleteQualificationRequest"></a>

### DeleteQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |






<a name="services-qualifications-DeleteQualificationResponse"></a>

### DeleteQualificationResponse







<a name="services-qualifications-DeleteQualificationResultRequest"></a>

### DeleteQualificationResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| result_id | [uint64](#uint64) |  |  |






<a name="services-qualifications-DeleteQualificationResultResponse"></a>

### DeleteQualificationResultResponse







<a name="services-qualifications-GetExamInfoRequest"></a>

### GetExamInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |






<a name="services-qualifications-GetExamInfoResponse"></a>

### GetExamInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification | [resources.qualifications.QualificationShort](#resources-qualifications-QualificationShort) |  |  |
| question_count | [int32](#int32) |  |  |
| exam_user | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) | optional |  |






<a name="services-qualifications-GetQualificationAccessRequest"></a>

### GetQualificationAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |






<a name="services-qualifications-GetQualificationAccessResponse"></a>

### GetQualificationAccessResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access | [resources.qualifications.QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |






<a name="services-qualifications-GetQualificationRequest"></a>

### GetQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| with_exam | [bool](#bool) | optional |  |






<a name="services-qualifications-GetQualificationResponse"></a>

### GetQualificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification | [resources.qualifications.Qualification](#resources-qualifications-Qualification) |  |  |






<a name="services-qualifications-GetUserExamRequest"></a>

### GetUserExamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |






<a name="services-qualifications-GetUserExamResponse"></a>

### GetUserExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| exam | [resources.qualifications.ExamQuestions](#resources-qualifications-ExamQuestions) |  |  |
| exam_user | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) |  |  |
| responses | [resources.qualifications.ExamResponses](#resources-qualifications-ExamResponses) |  |  |






<a name="services-qualifications-ListQualificationRequestsRequest"></a>

### ListQualificationRequestsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| qualification_id | [uint64](#uint64) | optional | Search params |
| status | [resources.qualifications.RequestStatus](#resources-qualifications-RequestStatus) | repeated |  |
| user_id | [int32](#int32) | optional |  |






<a name="services-qualifications-ListQualificationRequestsResponse"></a>

### ListQualificationRequestsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| requests | [resources.qualifications.QualificationRequest](#resources-qualifications-QualificationRequest) | repeated |  |






<a name="services-qualifications-ListQualificationsRequest"></a>

### ListQualificationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| search | [string](#string) | optional | Search params |
| job | [string](#string) | optional |  |






<a name="services-qualifications-ListQualificationsResponse"></a>

### ListQualificationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| qualifications | [resources.qualifications.Qualification](#resources-qualifications-Qualification) | repeated |  |






<a name="services-qualifications-ListQualificationsResultsRequest"></a>

### ListQualificationsResultsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| qualification_id | [uint64](#uint64) | optional | Search params |
| status | [resources.qualifications.ResultStatus](#resources-qualifications-ResultStatus) | repeated |  |
| user_id | [int32](#int32) | optional |  |






<a name="services-qualifications-ListQualificationsResultsResponse"></a>

### ListQualificationsResultsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| results | [resources.qualifications.QualificationResult](#resources-qualifications-QualificationResult) | repeated |  |






<a name="services-qualifications-SetQualificationAccessRequest"></a>

### SetQualificationAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| access | [resources.qualifications.QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |






<a name="services-qualifications-SetQualificationAccessResponse"></a>

### SetQualificationAccessResponse







<a name="services-qualifications-SubmitExamRequest"></a>

### SubmitExamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| responses | [resources.qualifications.ExamResponses](#resources-qualifications-ExamResponses) |  |  |






<a name="services-qualifications-SubmitExamResponse"></a>

### SubmitExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| duration | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |






<a name="services-qualifications-TakeExamRequest"></a>

### TakeExamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |
| cancel | [bool](#bool) | optional |  |






<a name="services-qualifications-TakeExamResponse"></a>

### TakeExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| exam | [resources.qualifications.ExamQuestions](#resources-qualifications-ExamQuestions) |  |  |
| exam_user | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) |  |  |






<a name="services-qualifications-UpdateQualificationRequest"></a>

### UpdateQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification | [resources.qualifications.Qualification](#resources-qualifications-Qualification) |  |  |






<a name="services-qualifications-UpdateQualificationResponse"></a>

### UpdateQualificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| qualification_id | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-qualifications-QualificationsService"></a>

### QualificationsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListQualifications | [ListQualificationsRequest](#services-qualifications-ListQualificationsRequest) | [ListQualificationsResponse](#services-qualifications-ListQualificationsResponse) | @perm |
| GetQualification | [GetQualificationRequest](#services-qualifications-GetQualificationRequest) | [GetQualificationResponse](#services-qualifications-GetQualificationResponse) | @perm: Name=ListQualifications |
| CreateQualification | [CreateQualificationRequest](#services-qualifications-CreateQualificationRequest) | [CreateQualificationResponse](#services-qualifications-CreateQualificationResponse) | @perm |
| UpdateQualification | [UpdateQualificationRequest](#services-qualifications-UpdateQualificationRequest) | [UpdateQualificationResponse](#services-qualifications-UpdateQualificationResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| DeleteQualification | [DeleteQualificationRequest](#services-qualifications-DeleteQualificationRequest) | [DeleteQualificationResponse](#services-qualifications-DeleteQualificationResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| ListQualificationRequests | [ListQualificationRequestsRequest](#services-qualifications-ListQualificationRequestsRequest) | [ListQualificationRequestsResponse](#services-qualifications-ListQualificationRequestsResponse) | @perm: Name=ListQualifications |
| CreateOrUpdateQualificationRequest | [CreateOrUpdateQualificationRequestRequest](#services-qualifications-CreateOrUpdateQualificationRequestRequest) | [CreateOrUpdateQualificationRequestResponse](#services-qualifications-CreateOrUpdateQualificationRequestResponse) | @perm: Name=ListQualifications |
| DeleteQualificationReq | [DeleteQualificationReqRequest](#services-qualifications-DeleteQualificationReqRequest) | [DeleteQualificationReqResponse](#services-qualifications-DeleteQualificationReqResponse) | @perm |
| ListQualificationsResults | [ListQualificationsResultsRequest](#services-qualifications-ListQualificationsResultsRequest) | [ListQualificationsResultsResponse](#services-qualifications-ListQualificationsResultsResponse) | @perm: Name=ListQualifications |
| CreateOrUpdateQualificationResult | [CreateOrUpdateQualificationResultRequest](#services-qualifications-CreateOrUpdateQualificationResultRequest) | [CreateOrUpdateQualificationResultResponse](#services-qualifications-CreateOrUpdateQualificationResultResponse) | @perm |
| DeleteQualificationResult | [DeleteQualificationResultRequest](#services-qualifications-DeleteQualificationResultRequest) | [DeleteQualificationResultResponse](#services-qualifications-DeleteQualificationResultResponse) | @perm |
| GetExamInfo | [GetExamInfoRequest](#services-qualifications-GetExamInfoRequest) | [GetExamInfoResponse](#services-qualifications-GetExamInfoResponse) | @perm: Name=ListQualifications |
| TakeExam | [TakeExamRequest](#services-qualifications-TakeExamRequest) | [TakeExamResponse](#services-qualifications-TakeExamResponse) | @perm: Name=ListQualifications |
| SubmitExam | [SubmitExamRequest](#services-qualifications-SubmitExamRequest) | [SubmitExamResponse](#services-qualifications-SubmitExamResponse) | @perm: Name=ListQualifications |
| GetUserExam | [GetUserExamRequest](#services-qualifications-GetUserExamRequest) | [GetUserExamResponse](#services-qualifications-GetUserExamResponse) | @perm: Name=CreateOrUpdateQualificationResult |

 <!-- end services -->



<a name="services_rector_config-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/rector/config.proto



<a name="services-rector-GetAppConfigRequest"></a>

### GetAppConfigRequest







<a name="services-rector-GetAppConfigResponse"></a>

### GetAppConfigResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [resources.rector.AppConfig](#resources-rector-AppConfig) |  |  |






<a name="services-rector-UpdateAppConfigRequest"></a>

### UpdateAppConfigRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [resources.rector.AppConfig](#resources-rector-AppConfig) |  |  |






<a name="services-rector-UpdateAppConfigResponse"></a>

### UpdateAppConfigResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [resources.rector.AppConfig](#resources-rector-AppConfig) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-rector-RectorConfigService"></a>

### RectorConfigService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAppConfig | [GetAppConfigRequest](#services-rector-GetAppConfigRequest) | [GetAppConfigResponse](#services-rector-GetAppConfigResponse) | @perm: Name=SuperUser |
| UpdateAppConfig | [UpdateAppConfigRequest](#services-rector-UpdateAppConfigRequest) | [UpdateAppConfigResponse](#services-rector-UpdateAppConfigResponse) | @perm: Name=SuperUser |

 <!-- end services -->



<a name="services_rector_filestore-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/rector/filestore.proto



<a name="services-rector-DeleteFileRequest"></a>

### DeleteFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |






<a name="services-rector-DeleteFileResponse"></a>

### DeleteFileResponse







<a name="services-rector-ListFilesRequest"></a>

### ListFilesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| path | [string](#string) | optional |  |






<a name="services-rector-ListFilesResponse"></a>

### ListFilesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| files | [resources.filestore.FileInfo](#resources-filestore-FileInfo) | repeated |  |






<a name="services-rector-UploadFileRequest"></a>

### UploadFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| prefix | [string](#string) |  |  |
| name | [string](#string) |  |  |
| file | [resources.filestore.File](#resources-filestore-File) |  |  |






<a name="services-rector-UploadFileResponse"></a>

### UploadFileResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file | [resources.filestore.FileInfo](#resources-filestore-FileInfo) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-rector-RectorFilestoreService"></a>

### RectorFilestoreService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListFiles | [ListFilesRequest](#services-rector-ListFilesRequest) | [ListFilesResponse](#services-rector-ListFilesResponse) | @perm: Name=SuperUser |
| UploadFile | [UploadFileRequest](#services-rector-UploadFileRequest) | [UploadFileResponse](#services-rector-UploadFileResponse) | @perm: Name=SuperUser |
| DeleteFile | [DeleteFileRequest](#services-rector-DeleteFileRequest) | [DeleteFileResponse](#services-rector-DeleteFileResponse) | @perm: Name=SuperUser |

 <!-- end services -->



<a name="services_rector_laws-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/rector/laws.proto



<a name="services-rector-CreateOrUpdateLawBookRequest"></a>

### CreateOrUpdateLawBookRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lawBook | [resources.laws.LawBook](#resources-laws-LawBook) |  |  |






<a name="services-rector-CreateOrUpdateLawBookResponse"></a>

### CreateOrUpdateLawBookResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lawBook | [resources.laws.LawBook](#resources-laws-LawBook) |  |  |






<a name="services-rector-CreateOrUpdateLawRequest"></a>

### CreateOrUpdateLawRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| law | [resources.laws.Law](#resources-laws-Law) |  |  |






<a name="services-rector-CreateOrUpdateLawResponse"></a>

### CreateOrUpdateLawResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| law | [resources.laws.Law](#resources-laws-Law) |  |  |






<a name="services-rector-DeleteLawBookRequest"></a>

### DeleteLawBookRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-rector-DeleteLawBookResponse"></a>

### DeleteLawBookResponse







<a name="services-rector-DeleteLawRequest"></a>

### DeleteLawRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-rector-DeleteLawResponse"></a>

### DeleteLawResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-rector-RectorLawsService"></a>

### RectorLawsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateOrUpdateLawBook | [CreateOrUpdateLawBookRequest](#services-rector-CreateOrUpdateLawBookRequest) | [CreateOrUpdateLawBookResponse](#services-rector-CreateOrUpdateLawBookResponse) | @perm: Name=SuperUser |
| DeleteLawBook | [DeleteLawBookRequest](#services-rector-DeleteLawBookRequest) | [DeleteLawBookResponse](#services-rector-DeleteLawBookResponse) | @perm: Name=SuperUser |
| CreateOrUpdateLaw | [CreateOrUpdateLawRequest](#services-rector-CreateOrUpdateLawRequest) | [CreateOrUpdateLawResponse](#services-rector-CreateOrUpdateLawResponse) | @perm: Name=SuperUser |
| DeleteLaw | [DeleteLawRequest](#services-rector-DeleteLawRequest) | [DeleteLawResponse](#services-rector-DeleteLawResponse) | @perm: Name=SuperUser |

 <!-- end services -->



<a name="services_rector_rector-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/rector/rector.proto



<a name="services-rector-AttrsUpdate"></a>

### AttrsUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to_update | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |
| to_remove | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-rector-CreateRoleRequest"></a>

### CreateRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| grade | [int32](#int32) |  |  |






<a name="services-rector-CreateRoleResponse"></a>

### CreateRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [resources.permissions.Role](#resources-permissions-Role) |  |  |






<a name="services-rector-DeleteFactionRequest"></a>

### DeleteFactionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint64](#uint64) |  |  |






<a name="services-rector-DeleteFactionResponse"></a>

### DeleteFactionResponse







<a name="services-rector-DeleteRoleRequest"></a>

### DeleteRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-rector-DeleteRoleResponse"></a>

### DeleteRoleResponse







<a name="services-rector-GetJobPropsRequest"></a>

### GetJobPropsRequest







<a name="services-rector-GetJobPropsResponse"></a>

### GetJobPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_props | [resources.users.JobProps](#resources-users-JobProps) |  |  |






<a name="services-rector-GetPermissionsRequest"></a>

### GetPermissionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint64](#uint64) |  |  |
| filtered | [bool](#bool) | optional |  |






<a name="services-rector-GetPermissionsResponse"></a>

### GetPermissionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| permissions | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| attributes | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-rector-GetRoleRequest"></a>

### GetRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| filtered | [bool](#bool) | optional |  |






<a name="services-rector-GetRoleResponse"></a>

### GetRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role | [resources.permissions.Role](#resources-permissions-Role) |  |  |






<a name="services-rector-GetRolesRequest"></a>

### GetRolesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lowest_rank | [bool](#bool) | optional |  |






<a name="services-rector-GetRolesResponse"></a>

### GetRolesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| roles | [resources.permissions.Role](#resources-permissions-Role) | repeated |  |






<a name="services-rector-PermItem"></a>

### PermItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| val | [bool](#bool) |  |  |






<a name="services-rector-PermsUpdate"></a>

### PermsUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| to_update | [PermItem](#services-rector-PermItem) | repeated |  |
| to_remove | [uint64](#uint64) | repeated |  |






<a name="services-rector-SetJobPropsRequest"></a>

### SetJobPropsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_props | [resources.users.JobProps](#resources-users-JobProps) |  |  |






<a name="services-rector-SetJobPropsResponse"></a>

### SetJobPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_props | [resources.users.JobProps](#resources-users-JobProps) |  |  |






<a name="services-rector-UpdateRoleLimitsRequest"></a>

### UpdateRoleLimitsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| role_id | [uint64](#uint64) |  |  |
| perms | [PermsUpdate](#services-rector-PermsUpdate) | optional |  |
| attrs | [AttrsUpdate](#services-rector-AttrsUpdate) | optional |  |






<a name="services-rector-UpdateRoleLimitsResponse"></a>

### UpdateRoleLimitsResponse







<a name="services-rector-UpdateRolePermsRequest"></a>

### UpdateRolePermsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| perms | [PermsUpdate](#services-rector-PermsUpdate) | optional |  |
| attrs | [AttrsUpdate](#services-rector-AttrsUpdate) | optional |  |






<a name="services-rector-UpdateRolePermsResponse"></a>

### UpdateRolePermsResponse







<a name="services-rector-ViewAuditLogRequest"></a>

### ViewAuditLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| user_ids | [int32](#int32) | repeated | Search params |
| from | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| to | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| services | [string](#string) | repeated | @sanitize: method=StripTags |
| methods | [string](#string) | repeated | @sanitize: method=StripTags |
| search | [string](#string) | optional |  |






<a name="services-rector-ViewAuditLogResponse"></a>

### ViewAuditLogResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| logs | [resources.rector.AuditEntry](#resources-rector-AuditEntry) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-rector-RectorService"></a>

### RectorService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetJobProps | [GetJobPropsRequest](#services-rector-GetJobPropsRequest) | [GetJobPropsResponse](#services-rector-GetJobPropsResponse) | @perm |
| SetJobProps | [SetJobPropsRequest](#services-rector-SetJobPropsRequest) | [SetJobPropsResponse](#services-rector-SetJobPropsResponse) | @perm |
| GetRoles | [GetRolesRequest](#services-rector-GetRolesRequest) | [GetRolesResponse](#services-rector-GetRolesResponse) | @perm |
| GetRole | [GetRoleRequest](#services-rector-GetRoleRequest) | [GetRoleResponse](#services-rector-GetRoleResponse) | @perm: Name=GetRoles |
| CreateRole | [CreateRoleRequest](#services-rector-CreateRoleRequest) | [CreateRoleResponse](#services-rector-CreateRoleResponse) | @perm |
| DeleteRole | [DeleteRoleRequest](#services-rector-DeleteRoleRequest) | [DeleteRoleResponse](#services-rector-DeleteRoleResponse) | @perm |
| UpdateRolePerms | [UpdateRolePermsRequest](#services-rector-UpdateRolePermsRequest) | [UpdateRolePermsResponse](#services-rector-UpdateRolePermsResponse) | @perm |
| GetPermissions | [GetPermissionsRequest](#services-rector-GetPermissionsRequest) | [GetPermissionsResponse](#services-rector-GetPermissionsResponse) | @perm: Name=GetRoles |
| ViewAuditLog | [ViewAuditLogRequest](#services-rector-ViewAuditLogRequest) | [ViewAuditLogResponse](#services-rector-ViewAuditLogResponse) | @perm |
| UpdateRoleLimits | [UpdateRoleLimitsRequest](#services-rector-UpdateRoleLimitsRequest) | [UpdateRoleLimitsResponse](#services-rector-UpdateRoleLimitsResponse) | @perm: Name=SuperUser |
| DeleteFaction | [DeleteFactionRequest](#services-rector-DeleteFactionRequest) | [DeleteFactionResponse](#services-rector-DeleteFactionResponse) | @perm: Name=SuperUser |

 <!-- end services -->



<a name="services_calendar_calendar-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/calendar/calendar.proto



<a name="services-calendar-CreateOrUpdateCalendarEntryRequest"></a>

### CreateOrUpdateCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) |  |  |
| user_ids | [int32](#int32) | repeated |  |






<a name="services-calendar-CreateOrUpdateCalendarEntryResponse"></a>

### CreateOrUpdateCalendarEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) |  |  |






<a name="services-calendar-CreateOrUpdateCalendarRequest"></a>

### CreateOrUpdateCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-CreateOrUpdateCalendarResponse"></a>

### CreateOrUpdateCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-DeleteCalendarEntryRequest"></a>

### DeleteCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry_id | [uint64](#uint64) |  |  |






<a name="services-calendar-DeleteCalendarEntryResponse"></a>

### DeleteCalendarEntryResponse







<a name="services-calendar-DeleteCalendarRequest"></a>

### DeleteCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar_id | [uint64](#uint64) |  |  |






<a name="services-calendar-DeleteCalendarResponse"></a>

### DeleteCalendarResponse







<a name="services-calendar-GetCalendarEntryRequest"></a>

### GetCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry_id | [uint64](#uint64) |  |  |






<a name="services-calendar-GetCalendarEntryResponse"></a>

### GetCalendarEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) |  |  |






<a name="services-calendar-GetCalendarRequest"></a>

### GetCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar_id | [uint64](#uint64) |  |  |






<a name="services-calendar-GetCalendarResponse"></a>

### GetCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| calendar | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-GetUpcomingEntriesRequest"></a>

### GetUpcomingEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| seconds | [int32](#int32) |  |  |






<a name="services-calendar-GetUpcomingEntriesResponse"></a>

### GetUpcomingEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entries | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) | repeated |  |






<a name="services-calendar-ListCalendarEntriesRequest"></a>

### ListCalendarEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| year | [int32](#int32) |  |  |
| month | [int32](#int32) |  |  |
| calendar_ids | [uint64](#uint64) | repeated |  |
| show_hidden | [bool](#bool) | optional |  |
| after | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-calendar-ListCalendarEntriesResponse"></a>

### ListCalendarEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entries | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) | repeated |  |






<a name="services-calendar-ListCalendarEntryRSVPRequest"></a>

### ListCalendarEntryRSVPRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| entry_id | [uint64](#uint64) |  |  |






<a name="services-calendar-ListCalendarEntryRSVPResponse"></a>

### ListCalendarEntryRSVPResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| entries | [resources.calendar.CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) | repeated |  |






<a name="services-calendar-ListCalendarsRequest"></a>

### ListCalendarsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| only_public | [bool](#bool) |  |  |
| min_access_level | [resources.calendar.AccessLevel](#resources-calendar-AccessLevel) | optional |  |
| after | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-calendar-ListCalendarsResponse"></a>

### ListCalendarsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| calendars | [resources.calendar.Calendar](#resources-calendar-Calendar) | repeated |  |






<a name="services-calendar-ListSubscriptionsRequest"></a>

### ListSubscriptionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |






<a name="services-calendar-ListSubscriptionsResponse"></a>

### ListSubscriptionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| subs | [resources.calendar.CalendarSub](#resources-calendar-CalendarSub) | repeated |  |






<a name="services-calendar-RSVPCalendarEntryRequest"></a>

### RSVPCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.calendar.CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) |  |  |
| subscribe | [bool](#bool) |  |  |
| remove | [bool](#bool) | optional |  |






<a name="services-calendar-RSVPCalendarEntryResponse"></a>

### RSVPCalendarEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [resources.calendar.CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) | optional |  |






<a name="services-calendar-ShareCalendarEntryRequest"></a>

### ShareCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry_id | [uint64](#uint64) |  |  |
| user_ids | [int32](#int32) | repeated |  |






<a name="services-calendar-ShareCalendarEntryResponse"></a>

### ShareCalendarEntryResponse







<a name="services-calendar-SubscribeToCalendarRequest"></a>

### SubscribeToCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sub | [resources.calendar.CalendarSub](#resources-calendar-CalendarSub) |  |  |
| delete | [bool](#bool) |  |  |






<a name="services-calendar-SubscribeToCalendarResponse"></a>

### SubscribeToCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sub | [resources.calendar.CalendarSub](#resources-calendar-CalendarSub) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-calendar-CalendarService"></a>

### CalendarService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListCalendars | [ListCalendarsRequest](#services-calendar-ListCalendarsRequest) | [ListCalendarsResponse](#services-calendar-ListCalendarsResponse) | @perm: Name=Any |
| GetCalendar | [GetCalendarRequest](#services-calendar-GetCalendarRequest) | [GetCalendarResponse](#services-calendar-GetCalendarResponse) | @perm: Name=Any |
| CreateOrUpdateCalendar | [CreateOrUpdateCalendarRequest](#services-calendar-CreateOrUpdateCalendarRequest) | [CreateOrUpdateCalendarResponse](#services-calendar-CreateOrUpdateCalendarResponse) | @perm: Attrs=Fields/StringList:[]string{"Job", "Public"} |
| DeleteCalendar | [DeleteCalendarRequest](#services-calendar-DeleteCalendarRequest) | [DeleteCalendarResponse](#services-calendar-DeleteCalendarResponse) | @perm |
| ListCalendarEntries | [ListCalendarEntriesRequest](#services-calendar-ListCalendarEntriesRequest) | [ListCalendarEntriesResponse](#services-calendar-ListCalendarEntriesResponse) | @perm: Name=Any |
| GetUpcomingEntries | [GetUpcomingEntriesRequest](#services-calendar-GetUpcomingEntriesRequest) | [GetUpcomingEntriesResponse](#services-calendar-GetUpcomingEntriesResponse) | @perm: Name=Any |
| GetCalendarEntry | [GetCalendarEntryRequest](#services-calendar-GetCalendarEntryRequest) | [GetCalendarEntryResponse](#services-calendar-GetCalendarEntryResponse) | @perm: Name=Any |
| CreateOrUpdateCalendarEntry | [CreateOrUpdateCalendarEntryRequest](#services-calendar-CreateOrUpdateCalendarEntryRequest) | [CreateOrUpdateCalendarEntryResponse](#services-calendar-CreateOrUpdateCalendarEntryResponse) | @perm |
| DeleteCalendarEntry | [DeleteCalendarEntryRequest](#services-calendar-DeleteCalendarEntryRequest) | [DeleteCalendarEntryResponse](#services-calendar-DeleteCalendarEntryResponse) | @perm |
| ShareCalendarEntry | [ShareCalendarEntryRequest](#services-calendar-ShareCalendarEntryRequest) | [ShareCalendarEntryResponse](#services-calendar-ShareCalendarEntryResponse) | @perm: Name=CreateOrUpdateCalendarEntry |
| ListCalendarEntryRSVP | [ListCalendarEntryRSVPRequest](#services-calendar-ListCalendarEntryRSVPRequest) | [ListCalendarEntryRSVPResponse](#services-calendar-ListCalendarEntryRSVPResponse) | @perm: Name=Any |
| RSVPCalendarEntry | [RSVPCalendarEntryRequest](#services-calendar-RSVPCalendarEntryRequest) | [RSVPCalendarEntryResponse](#services-calendar-RSVPCalendarEntryResponse) | @perm: Name=Any |
| ListSubscriptions | [ListSubscriptionsRequest](#services-calendar-ListSubscriptionsRequest) | [ListSubscriptionsResponse](#services-calendar-ListSubscriptionsResponse) | @perm: Name=Any |
| SubscribeToCalendar | [SubscribeToCalendarRequest](#services-calendar-SubscribeToCalendarRequest) | [SubscribeToCalendarResponse](#services-calendar-SubscribeToCalendarResponse) | @perm: Name=Any |

 <!-- end services -->



<a name="services_stats_stats-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/stats/stats.proto



<a name="services-stats-GetStatsRequest"></a>

### GetStatsRequest







<a name="services-stats-GetStatsResponse"></a>

### GetStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stats | [GetStatsResponse.StatsEntry](#services-stats-GetStatsResponse-StatsEntry) | repeated |  |






<a name="services-stats-GetStatsResponse-StatsEntry"></a>

### GetStatsResponse.StatsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [resources.stats.Stat](#resources-stats-Stat) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-stats-StatsService"></a>

### StatsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetStats | [GetStatsRequest](#services-stats-GetStatsRequest) | [GetStatsResponse](#services-stats-GetStatsResponse) |  |

 <!-- end services -->



<a name="services_internet_ads-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/internet/ads.proto



<a name="services-internet-GetAdsRequest"></a>

### GetAdsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ad_type | [resources.internet.AdType](#resources-internet-AdType) |  |  |
| count | [int32](#int32) |  |  |






<a name="services-internet-GetAdsResponse"></a>

### GetAdsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ads | [resources.internet.Ad](#resources-internet-Ad) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-internet-AdsService"></a>

### AdsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAds | [GetAdsRequest](#services-internet-GetAdsRequest) | [GetAdsResponse](#services-internet-GetAdsResponse) | @perm: Name=Any |

 <!-- end services -->



<a name="services_internet_internet-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/internet/internet.proto



<a name="services-internet-GetPageRequest"></a>

### GetPageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  |  |
| path | [string](#string) |  |  |






<a name="services-internet-GetPageResponse"></a>

### GetPageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [resources.internet.Page](#resources-internet-Page) | optional |  |






<a name="services-internet-SearchRequest"></a>

### SearchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| search | [string](#string) |  |  |






<a name="services-internet-SearchResponse"></a>

### SearchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| results | [resources.internet.SearchResult](#resources-internet-SearchResult) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-internet-InternetService"></a>

### InternetService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Search | [SearchRequest](#services-internet-SearchRequest) | [SearchResponse](#services-internet-SearchResponse) | @perm: Name=Any |
| GetPage | [GetPageRequest](#services-internet-GetPageRequest) | [GetPageResponse](#services-internet-GetPageResponse) | @perm: Name=Any |

 <!-- end services -->



<a name="services_mailer_mailer-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/mailer/mailer.proto



<a name="services-mailer-CreateOrUpdateEmailRequest"></a>

### CreateOrUpdateEmailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [resources.mailer.Email](#resources-mailer-Email) |  |  |






<a name="services-mailer-CreateOrUpdateEmailResponse"></a>

### CreateOrUpdateEmailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [resources.mailer.Email](#resources-mailer-Email) |  |  |






<a name="services-mailer-CreateOrUpdateTemplateRequest"></a>

### CreateOrUpdateTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.mailer.Template](#resources-mailer-Template) |  |  |






<a name="services-mailer-CreateOrUpdateTemplateResponse"></a>

### CreateOrUpdateTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.mailer.Template](#resources-mailer-Template) |  |  |






<a name="services-mailer-CreateThreadRequest"></a>

### CreateThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread | [resources.mailer.Thread](#resources-mailer-Thread) |  |  |
| message | [resources.mailer.Message](#resources-mailer-Message) |  |  |
| recipients | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="services-mailer-CreateThreadResponse"></a>

### CreateThreadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread | [resources.mailer.Thread](#resources-mailer-Thread) |  |  |






<a name="services-mailer-DeleteEmailRequest"></a>

### DeleteEmailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteEmailResponse"></a>

### DeleteEmailResponse







<a name="services-mailer-DeleteMessageRequest"></a>

### DeleteMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |
| message_id | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteMessageResponse"></a>

### DeleteMessageResponse







<a name="services-mailer-DeleteTemplateRequest"></a>

### DeleteTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| id | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteTemplateResponse"></a>

### DeleteTemplateResponse







<a name="services-mailer-DeleteThreadRequest"></a>

### DeleteThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteThreadResponse"></a>

### DeleteThreadResponse







<a name="services-mailer-GetEmailProposalsRequest"></a>

### GetEmailProposalsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| input | [string](#string) |  |  |
| job | [bool](#bool) | optional |  |
| user_id | [int32](#int32) | optional |  |






<a name="services-mailer-GetEmailProposalsResponse"></a>

### GetEmailProposalsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| emails | [string](#string) | repeated |  |
| domains | [string](#string) | repeated |  |






<a name="services-mailer-GetEmailRequest"></a>

### GetEmailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-mailer-GetEmailResponse"></a>

### GetEmailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [resources.mailer.Email](#resources-mailer-Email) |  |  |






<a name="services-mailer-GetEmailSettingsRequest"></a>

### GetEmailSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |






<a name="services-mailer-GetEmailSettingsResponse"></a>

### GetEmailSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| settings | [resources.mailer.EmailSettings](#resources-mailer-EmailSettings) |  |  |






<a name="services-mailer-GetTemplateRequest"></a>

### GetTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| template_id | [uint64](#uint64) |  |  |






<a name="services-mailer-GetTemplateResponse"></a>

### GetTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| template | [resources.mailer.Template](#resources-mailer-Template) |  |  |






<a name="services-mailer-GetThreadRequest"></a>

### GetThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |






<a name="services-mailer-GetThreadResponse"></a>

### GetThreadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread | [resources.mailer.Thread](#resources-mailer-Thread) |  |  |






<a name="services-mailer-GetThreadStateRequest"></a>

### GetThreadStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |






<a name="services-mailer-GetThreadStateResponse"></a>

### GetThreadStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [resources.mailer.ThreadState](#resources-mailer-ThreadState) |  |  |






<a name="services-mailer-ListEmailsRequest"></a>

### ListEmailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| all | [bool](#bool) | optional | Search params |






<a name="services-mailer-ListEmailsResponse"></a>

### ListEmailsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| emails | [resources.mailer.Email](#resources-mailer-Email) | repeated |  |






<a name="services-mailer-ListTemplatesRequest"></a>

### ListTemplatesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email_id | [uint64](#uint64) |  |  |






<a name="services-mailer-ListTemplatesResponse"></a>

### ListTemplatesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| templates | [resources.mailer.Template](#resources-mailer-Template) | repeated |  |






<a name="services-mailer-ListThreadMessagesRequest"></a>

### ListThreadMessagesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| email_id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |
| after | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-mailer-ListThreadMessagesResponse"></a>

### ListThreadMessagesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| messages | [resources.mailer.Message](#resources-mailer-Message) | repeated |  |






<a name="services-mailer-ListThreadsRequest"></a>

### ListThreadsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| email_ids | [uint64](#uint64) | repeated | Search params |
| unread | [bool](#bool) | optional |  |
| archived | [bool](#bool) | optional |  |






<a name="services-mailer-ListThreadsResponse"></a>

### ListThreadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| threads | [resources.mailer.Thread](#resources-mailer-Thread) | repeated |  |






<a name="services-mailer-PostMessageRequest"></a>

### PostMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [resources.mailer.Message](#resources-mailer-Message) |  |  |
| recipients | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="services-mailer-PostMessageResponse"></a>

### PostMessageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [resources.mailer.Message](#resources-mailer-Message) |  |  |






<a name="services-mailer-SearchThreadsRequest"></a>

### SearchThreadsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| search | [string](#string) |  | Search params |






<a name="services-mailer-SearchThreadsResponse"></a>

### SearchThreadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| messages | [resources.mailer.Message](#resources-mailer-Message) | repeated |  |






<a name="services-mailer-SetEmailSettingsRequest"></a>

### SetEmailSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| settings | [resources.mailer.EmailSettings](#resources-mailer-EmailSettings) |  |  |






<a name="services-mailer-SetEmailSettingsResponse"></a>

### SetEmailSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| settings | [resources.mailer.EmailSettings](#resources-mailer-EmailSettings) |  |  |






<a name="services-mailer-SetThreadStateRequest"></a>

### SetThreadStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [resources.mailer.ThreadState](#resources-mailer-ThreadState) |  |  |






<a name="services-mailer-SetThreadStateResponse"></a>

### SetThreadStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [resources.mailer.ThreadState](#resources-mailer-ThreadState) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-mailer-MailerService"></a>

### MailerService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListEmails | [ListEmailsRequest](#services-mailer-ListEmailsRequest) | [ListEmailsResponse](#services-mailer-ListEmailsResponse) | @perm |
| GetEmail | [GetEmailRequest](#services-mailer-GetEmailRequest) | [GetEmailResponse](#services-mailer-GetEmailResponse) | @perm: Name=ListEmails |
| CreateOrUpdateEmail | [CreateOrUpdateEmailRequest](#services-mailer-CreateOrUpdateEmailRequest) | [CreateOrUpdateEmailResponse](#services-mailer-CreateOrUpdateEmailResponse) | @perm: Attrs=Fields/StringList:[]string{"Job"} |
| DeleteEmail | [DeleteEmailRequest](#services-mailer-DeleteEmailRequest) | [DeleteEmailResponse](#services-mailer-DeleteEmailResponse) | @perm |
| GetEmailProposals | [GetEmailProposalsRequest](#services-mailer-GetEmailProposalsRequest) | [GetEmailProposalsResponse](#services-mailer-GetEmailProposalsResponse) | @perm: Name=ListEmails |
| ListTemplates | [ListTemplatesRequest](#services-mailer-ListTemplatesRequest) | [ListTemplatesResponse](#services-mailer-ListTemplatesResponse) | @perm: Name=ListEmails |
| GetTemplate | [GetTemplateRequest](#services-mailer-GetTemplateRequest) | [GetTemplateResponse](#services-mailer-GetTemplateResponse) | @perm: Name=ListEmails |
| CreateOrUpdateTemplate | [CreateOrUpdateTemplateRequest](#services-mailer-CreateOrUpdateTemplateRequest) | [CreateOrUpdateTemplateResponse](#services-mailer-CreateOrUpdateTemplateResponse) | @perm: Name=ListEmails |
| DeleteTemplate | [DeleteTemplateRequest](#services-mailer-DeleteTemplateRequest) | [DeleteTemplateResponse](#services-mailer-DeleteTemplateResponse) | @perm: Name=ListEmails |
| ListThreads | [ListThreadsRequest](#services-mailer-ListThreadsRequest) | [ListThreadsResponse](#services-mailer-ListThreadsResponse) | @perm: Name=ListEmails |
| GetThread | [GetThreadRequest](#services-mailer-GetThreadRequest) | [GetThreadResponse](#services-mailer-GetThreadResponse) | @perm: Name=ListEmails |
| CreateThread | [CreateThreadRequest](#services-mailer-CreateThreadRequest) | [CreateThreadResponse](#services-mailer-CreateThreadResponse) | @perm: Name=ListEmails |
| DeleteThread | [DeleteThreadRequest](#services-mailer-DeleteThreadRequest) | [DeleteThreadResponse](#services-mailer-DeleteThreadResponse) | @perm: Name=SuperUser |
| GetThreadState | [GetThreadStateRequest](#services-mailer-GetThreadStateRequest) | [GetThreadStateResponse](#services-mailer-GetThreadStateResponse) | @perm: Name=ListEmails |
| SetThreadState | [SetThreadStateRequest](#services-mailer-SetThreadStateRequest) | [SetThreadStateResponse](#services-mailer-SetThreadStateResponse) | @perm: Name=ListEmails |
| SearchThreads | [SearchThreadsRequest](#services-mailer-SearchThreadsRequest) | [SearchThreadsResponse](#services-mailer-SearchThreadsResponse) | @perm: Name=ListEmails |
| ListThreadMessages | [ListThreadMessagesRequest](#services-mailer-ListThreadMessagesRequest) | [ListThreadMessagesResponse](#services-mailer-ListThreadMessagesResponse) | @perm: Name=ListEmails |
| PostMessage | [PostMessageRequest](#services-mailer-PostMessageRequest) | [PostMessageResponse](#services-mailer-PostMessageResponse) | @perm: Name=ListEmails |
| DeleteMessage | [DeleteMessageRequest](#services-mailer-DeleteMessageRequest) | [DeleteMessageResponse](#services-mailer-DeleteMessageResponse) | @perm: Name=SuperUser |
| GetEmailSettings | [GetEmailSettingsRequest](#services-mailer-GetEmailSettingsRequest) | [GetEmailSettingsResponse](#services-mailer-GetEmailSettingsResponse) | @perm: Name=ListEmails |
| SetEmailSettings | [SetEmailSettingsRequest](#services-mailer-SetEmailSettingsRequest) | [SetEmailSettingsResponse](#services-mailer-SetEmailSettingsResponse) | @perm: Name=ListEmails |

 <!-- end services -->



<a name="services_wiki_wiki-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/wiki/wiki.proto



<a name="services-wiki-CreatePageRequest"></a>

### CreatePageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [resources.wiki.Page](#resources-wiki-Page) |  |  |






<a name="services-wiki-CreatePageResponse"></a>

### CreatePageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [resources.wiki.Page](#resources-wiki-Page) |  |  |






<a name="services-wiki-DeletePageRequest"></a>

### DeletePageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-wiki-DeletePageResponse"></a>

### DeletePageResponse







<a name="services-wiki-GetPageRequest"></a>

### GetPageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |






<a name="services-wiki-GetPageResponse"></a>

### GetPageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [resources.wiki.Page](#resources-wiki-Page) |  |  |






<a name="services-wiki-ListPageActivityRequest"></a>

### ListPageActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| page_id | [uint64](#uint64) |  |  |






<a name="services-wiki-ListPageActivityResponse"></a>

### ListPageActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.wiki.PageActivity](#resources-wiki-PageActivity) | repeated |  |






<a name="services-wiki-ListPagesRequest"></a>

### ListPagesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| sort | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| job | [string](#string) | optional | Search params |
| root_only | [bool](#bool) | optional |  |
| search | [string](#string) | optional |  |






<a name="services-wiki-ListPagesResponse"></a>

### ListPagesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| pages | [resources.wiki.PageShort](#resources-wiki-PageShort) | repeated |  |






<a name="services-wiki-UpdatePageRequest"></a>

### UpdatePageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [resources.wiki.Page](#resources-wiki-Page) |  |  |






<a name="services-wiki-UpdatePageResponse"></a>

### UpdatePageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [resources.wiki.Page](#resources-wiki-Page) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-wiki-WikiService"></a>

### WikiService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListPages | [ListPagesRequest](#services-wiki-ListPagesRequest) | [ListPagesResponse](#services-wiki-ListPagesResponse) | @perm |
| GetPage | [GetPageRequest](#services-wiki-GetPageRequest) | [GetPageResponse](#services-wiki-GetPageResponse) | @perm: Name=ListPages |
| CreatePage | [CreatePageRequest](#services-wiki-CreatePageRequest) | [CreatePageResponse](#services-wiki-CreatePageResponse) | @perm: Attrs=Fields/StringList:[]string{"Public"} |
| UpdatePage | [UpdatePageRequest](#services-wiki-UpdatePageRequest) | [UpdatePageResponse](#services-wiki-UpdatePageResponse) | @perm |
| DeletePage | [DeletePageRequest](#services-wiki-DeletePageRequest) | [DeletePageResponse](#services-wiki-DeletePageResponse) | @perm |
| ListPageActivity | [ListPageActivityRequest](#services-wiki-ListPageActivityRequest) | [ListPageActivityResponse](#services-wiki-ListPageActivityResponse) | @perm |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

