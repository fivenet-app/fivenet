# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [resources/accounts/accounts.proto](#resources_accounts_accounts-proto)
    - [Account](#resources-accounts-Account)
    - [Character](#resources-accounts-Character)
  
- [resources/accounts/oauth2.proto](#resources_accounts_oauth2-proto)
    - [OAuth2Account](#resources-accounts-OAuth2Account)
    - [OAuth2Provider](#resources-accounts-OAuth2Provider)
  
- [resources/calendar/access.proto](#resources_calendar_access-proto)
    - [CalendarAccess](#resources-calendar-CalendarAccess)
    - [CalendarJobAccess](#resources-calendar-CalendarJobAccess)
    - [CalendarUserAccess](#resources-calendar-CalendarUserAccess)
  
    - [AccessLevel](#resources-calendar-AccessLevel)
    - [AccessLevelUpdateMode](#resources-calendar-AccessLevelUpdateMode)
  
- [resources/calendar/calendar.proto](#resources_calendar_calendar-proto)
    - [Calendar](#resources-calendar-Calendar)
    - [CalendarEntry](#resources-calendar-CalendarEntry)
    - [CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP)
    - [CalendarEntryRecurring](#resources-calendar-CalendarEntryRecurring)
    - [CalendarShort](#resources-calendar-CalendarShort)
    - [CalendarSub](#resources-calendar-CalendarSub)
  
    - [RsvpResponses](#resources-calendar-RsvpResponses)
  
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
    - [OrderBy](#resources-common-database-OrderBy)
    - [PaginationRequest](#resources-common-database-PaginationRequest)
    - [PaginationResponse](#resources-common-database-PaginationResponse)
  
- [resources/common/i18n.proto](#resources_common_i18n-proto)
    - [TranslateItem](#resources-common-TranslateItem)
    - [TranslateItem.ParametersEntry](#resources-common-TranslateItem-ParametersEntry)
  
- [resources/documents/access.proto](#resources_documents_access-proto)
    - [DocumentAccess](#resources-documents-DocumentAccess)
    - [DocumentJobAccess](#resources-documents-DocumentJobAccess)
    - [DocumentUserAccess](#resources-documents-DocumentUserAccess)
  
    - [AccessLevel](#resources-documents-AccessLevel)
    - [AccessLevelUpdateMode](#resources-documents-AccessLevelUpdateMode)
  
- [resources/documents/activity.proto](#resources_documents_activity-proto)
    - [DocAccessRequested](#resources-documents-DocAccessRequested)
    - [DocActivity](#resources-documents-DocActivity)
    - [DocActivityData](#resources-documents-DocActivityData)
    - [DocOwnerChanged](#resources-documents-DocOwnerChanged)
    - [DocUpdated](#resources-documents-DocUpdated)
  
    - [DocActivityType](#resources-documents-DocActivityType)
  
- [resources/documents/category.proto](#resources_documents_category-proto)
    - [Category](#resources-documents-Category)
  
- [resources/documents/comment.proto](#resources_documents_comment-proto)
    - [Comment](#resources-documents-Comment)
  
- [resources/documents/documents.proto](#resources_documents_documents-proto)
    - [Document](#resources-documents-Document)
    - [DocumentReference](#resources-documents-DocumentReference)
    - [DocumentRelation](#resources-documents-DocumentRelation)
    - [DocumentShort](#resources-documents-DocumentShort)
  
    - [DocContentType](#resources-documents-DocContentType)
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
  
- [resources/filestore/file.proto](#resources_filestore_file-proto)
    - [File](#resources-filestore-File)
    - [FileInfo](#resources-filestore-FileInfo)
  
- [resources/jobs/colleagues.proto](#resources_jobs_colleagues-proto)
    - [Colleague](#resources-jobs-Colleague)
    - [ColleagueAbsenceDate](#resources-jobs-ColleagueAbsenceDate)
    - [ColleagueGradeChange](#resources-jobs-ColleagueGradeChange)
    - [JobsUserActivity](#resources-jobs-JobsUserActivity)
    - [JobsUserActivityData](#resources-jobs-JobsUserActivityData)
    - [JobsUserProps](#resources-jobs-JobsUserProps)
  
    - [JobsUserActivityType](#resources-jobs-JobsUserActivityType)
  
- [resources/jobs/conduct.proto](#resources_jobs_conduct-proto)
    - [ConductEntry](#resources-jobs-ConductEntry)
  
    - [ConductType](#resources-jobs-ConductType)
  
- [resources/jobs/timeclock.proto](#resources_jobs_timeclock-proto)
    - [TimeclockEntry](#resources-jobs-TimeclockEntry)
    - [TimeclockStats](#resources-jobs-TimeclockStats)
    - [TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats)
  
- [resources/laws/laws.proto](#resources_laws_laws-proto)
    - [Law](#resources-laws-Law)
    - [LawBook](#resources-laws-LawBook)
  
- [resources/livemap/livemap.proto](#resources_livemap_livemap-proto)
    - [CircleMarker](#resources-livemap-CircleMarker)
    - [Coords](#resources-livemap-Coords)
    - [IconMarker](#resources-livemap-IconMarker)
    - [MarkerData](#resources-livemap-MarkerData)
    - [MarkerInfo](#resources-livemap-MarkerInfo)
    - [MarkerMarker](#resources-livemap-MarkerMarker)
    - [UserMarker](#resources-livemap-UserMarker)
  
    - [MarkerType](#resources-livemap-MarkerType)
  
- [resources/livemap/tracker.proto](#resources_livemap_tracker-proto)
    - [UsersUpdateEvent](#resources-livemap-UsersUpdateEvent)
  
- [resources/notifications/notifications.proto](#resources_notifications_notifications-proto)
    - [CalendarData](#resources-notifications-CalendarData)
    - [Data](#resources-notifications-Data)
    - [Link](#resources-notifications-Link)
    - [Notification](#resources-notifications-Notification)
  
    - [NotificationCategory](#resources-notifications-NotificationCategory)
    - [NotificationType](#resources-notifications-NotificationType)
  
- [resources/notifications/events.proto](#resources_notifications_events-proto)
    - [JobEvent](#resources-notifications-JobEvent)
    - [SystemEvent](#resources-notifications-SystemEvent)
    - [UserEvent](#resources-notifications-UserEvent)
  
- [resources/permissions/permissions.proto](#resources_permissions_permissions-proto)
    - [AttributeValues](#resources-permissions-AttributeValues)
    - [JobGradeList](#resources-permissions-JobGradeList)
    - [JobGradeList.JobsEntry](#resources-permissions-JobGradeList-JobsEntry)
    - [Permission](#resources-permissions-Permission)
    - [RawRoleAttribute](#resources-permissions-RawRoleAttribute)
    - [Role](#resources-permissions-Role)
    - [RoleAttribute](#resources-permissions-RoleAttribute)
    - [StringList](#resources-permissions-StringList)
  
- [resources/qualifications/exam.proto](#resources_qualifications_exam-proto)
    - [ExamQuestion](#resources-qualifications-ExamQuestion)
    - [ExamQuestionAnswerData](#resources-qualifications-ExamQuestionAnswerData)
    - [ExamQuestionData](#resources-qualifications-ExamQuestionData)
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
  
- [resources/qualifications/access.proto](#resources_qualifications_access-proto)
    - [QualificationAccess](#resources-qualifications-QualificationAccess)
    - [QualificationJobAccess](#resources-qualifications-QualificationJobAccess)
  
    - [AccessLevel](#resources-qualifications-AccessLevel)
    - [AccessLevelUpdateMode](#resources-qualifications-AccessLevelUpdateMode)
  
- [resources/rector/audit.proto](#resources_rector_audit-proto)
    - [AuditEntry](#resources-rector-AuditEntry)
  
    - [EventType](#resources-rector-EventType)
  
- [resources/rector/config.proto](#resources_rector_config-proto)
    - [AppConfig](#resources-rector-AppConfig)
    - [Auth](#resources-rector-Auth)
    - [Discord](#resources-rector-Discord)
    - [JobInfo](#resources-rector-JobInfo)
    - [Links](#resources-rector-Links)
    - [Perm](#resources-rector-Perm)
    - [Perms](#resources-rector-Perms)
    - [UnemployedJob](#resources-rector-UnemployedJob)
    - [UserTracker](#resources-rector-UserTracker)
    - [Website](#resources-rector-Website)
  
- [resources/timestamp/timestamp.proto](#resources_timestamp_timestamp-proto)
    - [Timestamp](#resources-timestamp-Timestamp)
  
- [resources/users/users.proto](#resources_users_users-proto)
    - [CitizenAttribute](#resources-users-CitizenAttribute)
    - [CitizenAttributes](#resources-users-CitizenAttributes)
    - [License](#resources-users-License)
    - [User](#resources-users-User)
    - [UserActivity](#resources-users-UserActivity)
    - [UserProps](#resources-users-UserProps)
    - [UserShort](#resources-users-UserShort)
  
    - [UserActivityType](#resources-users-UserActivityType)
  
- [resources/users/jobs.proto](#resources_users_jobs-proto)
    - [DiscordSyncChange](#resources-users-DiscordSyncChange)
    - [DiscordSyncChanges](#resources-users-DiscordSyncChanges)
    - [DiscordSyncSettings](#resources-users-DiscordSyncSettings)
    - [GroupMapping](#resources-users-GroupMapping)
    - [GroupSyncSettings](#resources-users-GroupSyncSettings)
    - [Job](#resources-users-Job)
    - [JobGrade](#resources-users-JobGrade)
    - [JobProps](#resources-users-JobProps)
    - [JobSettings](#resources-users-JobSettings)
    - [JobsAbsenceSettings](#resources-users-JobsAbsenceSettings)
    - [QuickButtons](#resources-users-QuickButtons)
    - [StatusLogSettings](#resources-users-StatusLogSettings)
    - [UserInfoSyncSettings](#resources-users-UserInfoSyncSettings)
  
    - [UserInfoSyncUnemployedMode](#resources-users-UserInfoSyncUnemployedMode)
  
- [resources/vehicles/vehicles.proto](#resources_vehicles_vehicles-proto)
    - [Vehicle](#resources-vehicles-Vehicle)
  
- [resources/messenger/message.proto](#resources_messenger_message-proto)
    - [Message](#resources-messenger-Message)
    - [MessageData](#resources-messenger-MessageData)
  
- [resources/messenger/access.proto](#resources_messenger_access-proto)
    - [ThreadAccess](#resources-messenger-ThreadAccess)
    - [ThreadJobAccess](#resources-messenger-ThreadJobAccess)
    - [ThreadUserAccess](#resources-messenger-ThreadUserAccess)
  
    - [AccessLevel](#resources-messenger-AccessLevel)
    - [AccessLevelUpdateMode](#resources-messenger-AccessLevelUpdateMode)
  
- [resources/messenger/events.proto](#resources_messenger_events-proto)
    - [MessengerEvent](#resources-messenger-MessengerEvent)
  
- [resources/messenger/thread.proto](#resources_messenger_thread-proto)
    - [Thread](#resources-messenger-Thread)
    - [ThreadUserState](#resources-messenger-ThreadUserState)
  
- [resources/messenger/user.proto](#resources_messenger_user-proto)
    - [UserStatus](#resources-messenger-UserStatus)
  
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
  
- [services/jobs/timeclock.proto](#services_jobs_timeclock-proto)
    - [GetTimeclockStatsRequest](#services-jobs-GetTimeclockStatsRequest)
    - [GetTimeclockStatsResponse](#services-jobs-GetTimeclockStatsResponse)
    - [ListInactiveEmployeesRequest](#services-jobs-ListInactiveEmployeesRequest)
    - [ListInactiveEmployeesResponse](#services-jobs-ListInactiveEmployeesResponse)
    - [ListTimeclockRequest](#services-jobs-ListTimeclockRequest)
    - [ListTimeclockResponse](#services-jobs-ListTimeclockResponse)
  
    - [JobsTimeclockService](#services-jobs-JobsTimeclockService)
  
- [services/jobs/jobs.proto](#services_jobs_jobs-proto)
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
    - [SetJobsUserPropsRequest](#services-jobs-SetJobsUserPropsRequest)
    - [SetJobsUserPropsResponse](#services-jobs-SetJobsUserPropsResponse)
    - [SetMOTDRequest](#services-jobs-SetMOTDRequest)
    - [SetMOTDResponse](#services-jobs-SetMOTDResponse)
  
    - [JobsService](#services-jobs-JobsService)
  
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
  
- [services/messenger/messenger.proto](#services_messenger_messenger-proto)
    - [CreateOrUpdateThreadRequest](#services-messenger-CreateOrUpdateThreadRequest)
    - [CreateOrUpdateThreadResponse](#services-messenger-CreateOrUpdateThreadResponse)
    - [DeleteMessageRequest](#services-messenger-DeleteMessageRequest)
    - [DeleteMessageResponse](#services-messenger-DeleteMessageResponse)
    - [DeleteThreadRequest](#services-messenger-DeleteThreadRequest)
    - [DeleteThreadResponse](#services-messenger-DeleteThreadResponse)
    - [GetThreadMessagesRequest](#services-messenger-GetThreadMessagesRequest)
    - [GetThreadMessagesResponse](#services-messenger-GetThreadMessagesResponse)
    - [GetThreadRequest](#services-messenger-GetThreadRequest)
    - [GetThreadResponse](#services-messenger-GetThreadResponse)
    - [LeaveThreadRequest](#services-messenger-LeaveThreadRequest)
    - [LeaveThreadResponse](#services-messenger-LeaveThreadResponse)
    - [ListThreadsRequest](#services-messenger-ListThreadsRequest)
    - [ListThreadsResponse](#services-messenger-ListThreadsResponse)
    - [PostMessageRequest](#services-messenger-PostMessageRequest)
    - [PostMessageResponse](#services-messenger-PostMessageResponse)
    - [SetThreadUserStateRequest](#services-messenger-SetThreadUserStateRequest)
    - [SetThreadUserStateResponse](#services-messenger-SetThreadUserStateResponse)
  
    - [MessengerService](#services-messenger-MessengerService)
  
- [Scalar Value Types](#scalar-value-types)



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
| char | [resources.users.User](#resources-users-User) |  |  |





 

 

 

 



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





 

 

 

 



<a name="resources_calendar_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/calendar/access.proto



<a name="resources-calendar-CalendarAccess"></a>

### CalendarAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [CalendarJobAccess](#resources-calendar-CalendarJobAccess) | repeated | @gotags: alias:&#34;job_access&#34; |
| users | [CalendarUserAccess](#resources-calendar-CalendarUserAccess) | repeated | @gotags: alias:&#34;user_access&#34; |






<a name="resources-calendar-CalendarJobAccess"></a>

### CalendarJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| calendar_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional | @gotags: alias:&#34;job_label&#34; |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional | @gotags: alias:&#34;job_grade_label&#34; |
| access | [AccessLevel](#resources-calendar-AccessLevel) |  |  |






<a name="resources-calendar-CalendarUserAccess"></a>

### CalendarUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| calendar_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-calendar-AccessLevel) |  |  |





 


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



<a name="resources-calendar-AccessLevelUpdateMode"></a>

### AccessLevelUpdateMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_UPDATE_MODE_UPDATE | 1 |  |
| ACCESS_LEVEL_UPDATE_MODE_DELETE | 2 |  |
| ACCESS_LEVEL_UPDATE_MODE_CLEAR | 3 |  |


 

 

 



<a name="resources_calendar_calendar-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/calendar/calendar.proto



<a name="resources-calendar-Calendar"></a>

### Calendar



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
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
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |
| subscription | [CalendarSub](#resources-calendar-CalendarSub) | optional |  |
| access | [CalendarAccess](#resources-calendar-CalendarAccess) |  |  |






<a name="resources-calendar-CalendarEntry"></a>

### CalendarEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| calendar_id | [uint64](#uint64) |  |  |
| calendar | [Calendar](#resources-calendar-Calendar) | optional |  |
| job | [string](#string) | optional |  |
| start_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| end_time | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| title | [string](#string) |  | @sanitize: method=StripTags |
| content | [string](#string) |  | @sanitize |
| closed | [bool](#bool) |  |  |
| rsvp_open | [bool](#bool) | optional |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
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
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
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


 

 

 



<a name="resources_centrum_dispatches-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/dispatches.proto



<a name="resources-centrum-Dispatch"></a>

### Dispatch



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
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
| dispatch_id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;dispatch_id&#34; |
| unit_id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;unit_id&#34; |
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
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| dispatch_id | [uint64](#uint64) |  |  |
| unit_id | [uint64](#uint64) | optional |  |
| unit | [Unit](#resources-centrum-Unit) | optional |  |
| status | [StatusDispatch](#resources-centrum-StatusDispatch) |  |  |
| reason | [string](#string) | optional | @sanitize |
| code | [string](#string) | optional | @sanitize |
| user_id | [int32](#int32) | optional |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| x | [double](#double) | optional |  |
| y | [double](#double) | optional |  |
| postal | [string](#string) | optional | @sanitize |





 


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


 

 

 



<a name="resources_centrum_general-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/general.proto



<a name="resources-centrum-Attributes"></a>

### Attributes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| list | [string](#string) | repeated |  |






<a name="resources-centrum-Disponents"></a>

### Disponents



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| disponents | [resources.users.UserShort](#resources-users-UserShort) | repeated |  |






<a name="resources-centrum-UserUnitMapping"></a>

### UserUnitMapping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| user_id | [int32](#int32) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |





 

 

 

 



<a name="resources_centrum_settings-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/settings.proto



<a name="resources-centrum-PredefinedStatus"></a>

### PredefinedStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_status | [string](#string) | repeated |  |
| dispatch_status | [string](#string) | repeated |  |






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





 


<a name="resources-centrum-CentrumMode"></a>

### CentrumMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| CENTRUM_MODE_UNSPECIFIED | 0 |  |
| CENTRUM_MODE_MANUAL | 1 |  |
| CENTRUM_MODE_CENTRAL_COMMAND | 2 |  |
| CENTRUM_MODE_AUTO_ROUND_ROBIN | 3 |  |
| CENTRUM_MODE_SIMPLIFIED | 4 |  |


 

 

 



<a name="resources_centrum_units-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/units.proto



<a name="resources-centrum-Unit"></a>

### Unit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| name | [string](#string) |  | @sanitize |
| initials | [string](#string) |  | @sanitize |
| color | [string](#string) |  |  |
| description | [string](#string) | optional | @sanitize |
| status | [UnitStatus](#resources-centrum-UnitStatus) | optional |  |
| users | [UnitAssignment](#resources-centrum-UnitAssignment) | repeated |  |
| attributes | [Attributes](#resources-centrum-Attributes) | optional |  |
| home_postal | [string](#string) | optional |  |






<a name="resources-centrum-UnitAssignment"></a>

### UnitAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unit_id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;unit_id&#34; |
| user_id | [int32](#int32) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;user_id&#34; |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |






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
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| unit_id | [uint64](#uint64) |  |  |
| unit | [Unit](#resources-centrum-Unit) | optional |  |
| status | [StatusUnit](#resources-centrum-StatusUnit) |  |  |
| reason | [string](#string) | optional | @sanitize |
| code | [string](#string) | optional | @sanitize |
| user_id | [int32](#int32) | optional |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| x | [double](#double) | optional |  |
| y | [double](#double) | optional |  |
| postal | [string](#string) | optional | @sanitize |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional |  |





 


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


 

 

 



<a name="resources_common_database_database-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/database/database.proto



<a name="resources-common-database-OrderBy"></a>

### OrderBy



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| column | [string](#string) |  |  |
| desc | [bool](#bool) |  |  |






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





 

 

 

 



<a name="resources_common_i18n-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/i18n.proto



<a name="resources-common-TranslateItem"></a>

### TranslateItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| parameters | [TranslateItem.ParametersEntry](#resources-common-TranslateItem-ParametersEntry) | repeated |  |






<a name="resources-common-TranslateItem-ParametersEntry"></a>

### TranslateItem.ParametersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 

 

 

 



<a name="resources_documents_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/access.proto



<a name="resources-documents-DocumentAccess"></a>

### DocumentAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated | @gotags: alias:&#34;job_access&#34; |
| users | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated | @gotags: alias:&#34;user_access&#34; |






<a name="resources-documents-DocumentJobAccess"></a>

### DocumentJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| document_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional | @gotags: alias:&#34;job_label&#34; |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional | @gotags: alias:&#34;job_grade_label&#34; |
| access | [AccessLevel](#resources-documents-AccessLevel) |  |  |
| required | [bool](#bool) | optional | @gotags: alias:&#34;required&#34; |






<a name="resources-documents-DocumentUserAccess"></a>

### DocumentUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| document_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-documents-AccessLevel) |  |  |
| required | [bool](#bool) | optional | @gotags: alias:&#34;required&#34; |





 


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



<a name="resources-documents-AccessLevelUpdateMode"></a>

### AccessLevelUpdateMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_UPDATE_MODE_UPDATE | 1 |  |
| ACCESS_LEVEL_UPDATE_MODE_DELETE | 2 |  |
| ACCESS_LEVEL_UPDATE_MODE_CLEAR | 3 |  |


 

 

 



<a name="resources_documents_activity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/activity.proto



<a name="resources-documents-DocAccessRequested"></a>

### DocAccessRequested



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| level | [AccessLevel](#resources-documents-AccessLevel) |  |  |






<a name="resources-documents-DocActivity"></a>

### DocActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| document_id | [uint64](#uint64) |  |  |
| activity_type | [DocActivityType](#resources-documents-DocActivityType) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
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
| access_updated | [DocumentAccess](#resources-documents-DocumentAccess) |  |  |
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





 

 

 

 



<a name="resources_documents_comment-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/comment.proto



<a name="resources-documents-Comment"></a>

### Comment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| document_id | [uint64](#uint64) |  |  |
| comment | [string](#string) |  | @sanitize: method=StripTags |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |





 

 

 

 



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
| category | [Category](#resources-documents-Category) | optional | @gotags: alias:&#34;category&#34; |
| title | [string](#string) |  | @sanitize |
| content_type | [DocContentType](#resources-documents-DocContentType) |  | @gotags: alias:&#34;content_type&#34; |
| content | [string](#string) |  | @sanitize |
| data | [string](#string) | optional | @sanitize

@gotags: alias:&#34;data&#34; |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |
| state | [string](#string) |  | @sanitize |
| closed | [bool](#bool) |  |  |
| public | [bool](#bool) |  |  |
| template_id | [uint64](#uint64) | optional |  |






<a name="resources-documents-DocumentReference"></a>

### DocumentReference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) | optional |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| source_document_id | [uint64](#uint64) |  | @gotags: alias:&#34;source_document_id&#34; |
| source_document | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:&#34;source_document&#34; |
| reference | [DocReference](#resources-documents-DocReference) |  | @gotags: alias:&#34;reference&#34; |
| target_document_id | [uint64](#uint64) |  | @gotags: alias:&#34;target_document_id&#34; |
| target_document | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:&#34;target_document&#34; |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;ref_creator&#34; |






<a name="resources-documents-DocumentRelation"></a>

### DocumentRelation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) | optional |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| document_id | [uint64](#uint64) |  |  |
| document | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:&#34;document&#34; |
| source_user_id | [int32](#int32) |  | @gotags: alias:&#34;source_user_id&#34; |
| source_user | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;source_user&#34; |
| relation | [DocRelation](#resources-documents-DocRelation) |  | @gotags: alias:&#34;relation&#34; |
| target_user_id | [int32](#int32) |  | @gotags: alias:&#34;target_user_id&#34; |
| target_user | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;target_user&#34; |






<a name="resources-documents-DocumentShort"></a>

### DocumentShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category_id | [uint64](#uint64) | optional |  |
| category | [Category](#resources-documents-Category) | optional | @gotags: alias:&#34;category&#34; |
| title | [string](#string) |  | @sanitize |
| content_type | [DocContentType](#resources-documents-DocContentType) |  | @gotags: alias:&#34;content_type&#34; |
| content | [string](#string) |  | @sanitize |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  | @gotags: alias:&#34;creator_job&#34; |
| state | [string](#string) |  | @sanitize

@gotags: alias:&#34;state&#34; |
| closed | [bool](#bool) |  |  |
| public | [bool](#bool) |  |  |





 


<a name="resources-documents-DocContentType"></a>

### DocContentType


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOC_CONTENT_TYPE_UNSPECIFIED | 0 |  |
| DOC_CONTENT_TYPE_HTML | 1 |  |
| DOC_CONTENT_TYPE_PLAIN | 2 |  |



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
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| reason | [string](#string) | optional |  |
| data | [DocActivityData](#resources-documents-DocActivityData) |  |  |
| accepted | [bool](#bool) | optional |  |





 

 

 

 



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
| id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category | [Category](#resources-documents-Category) |  | @gotags: alias:&#34;category&#34; |
| weight | [uint32](#uint32) |  |  |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) |  | @sanitize |
| content_title | [string](#string) |  | @gotags: alias:&#34;content_title&#34; |
| content | [string](#string) |  | @gotags: alias:&#34;content&#34; |
| state | [string](#string) |  | @gotags: alias:&#34;state&#34; |
| schema | [TemplateSchema](#resources-documents-TemplateSchema) |  | @gotags: alias:&#34;schema&#34; |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |
| job_access | [TemplateJobAccess](#resources-documents-TemplateJobAccess) | repeated |  |
| content_access | [DocumentAccess](#resources-documents-DocumentAccess) |  | @gotags: alias:&#34;access&#34; |






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
| id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| template_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional | @gotags: alias:&#34;job_grade_label&#34; |
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
| id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| category | [Category](#resources-documents-Category) |  | @gotags: alias:&#34;category&#34; |
| weight | [uint32](#uint32) |  |  |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) |  | @sanitize |
| schema | [TemplateSchema](#resources-documents-TemplateSchema) |  | @gotags: alias:&#34;schema&#34; |
| creator_job | [string](#string) |  |  |
| creator_job_label | [string](#string) | optional |  |





 

 

 

 



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





 

 

 

 



<a name="resources_jobs_colleagues-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/colleagues.proto



<a name="resources-jobs-Colleague"></a>

### Colleague



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  | @gotags: alias:&#34;id&#34; |
| identifier | [string](#string) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| job_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| firstname | [string](#string) |  |  |
| lastname | [string](#string) |  |  |
| dateofbirth | [string](#string) |  |  |
| phone_number | [string](#string) | optional |  |
| avatar | [resources.filestore.File](#resources-filestore-File) | optional |  |
| props | [JobsUserProps](#resources-jobs-JobsUserProps) |  | @gotags: alias:&#34;fivenet_jobs_user_props&#34; |






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






<a name="resources-jobs-JobsUserActivity"></a>

### JobsUserActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| source_user_id | [int32](#int32) |  |  |
| source_user | [Colleague](#resources-jobs-Colleague) |  | @gotags: alias:&#34;source_user&#34; |
| target_user_id | [int32](#int32) |  |  |
| target_user | [Colleague](#resources-jobs-Colleague) |  | @gotags: alias:&#34;target_user&#34; |
| activity_type | [JobsUserActivityType](#resources-jobs-JobsUserActivityType) |  |  |
| reason | [string](#string) |  | @sanitize |
| data | [JobsUserActivityData](#resources-jobs-JobsUserActivityData) |  |  |






<a name="resources-jobs-JobsUserActivityData"></a>

### JobsUserActivityData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| absence_date | [ColleagueAbsenceDate](#resources-jobs-ColleagueAbsenceDate) |  |  |
| grade_change | [ColleagueGradeChange](#resources-jobs-ColleagueGradeChange) |  |  |






<a name="resources-jobs-JobsUserProps"></a>

### JobsUserProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| job | [string](#string) |  |  |
| absence_begin | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| absence_end | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| note | [string](#string) | optional | @sanitize: method=StripTags |





 


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


 

 

 



<a name="resources_jobs_conduct-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/conduct.proto



<a name="resources-jobs-ConductEntry"></a>

### ConductEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| type | [ConductType](#resources-jobs-ConductType) |  |  |
| message | [string](#string) |  | @sanitize |
| expires_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| target_user_id | [int32](#int32) |  |  |
| target_user | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:&#34;target_user&#34; |
| creator_id | [int32](#int32) |  |  |
| creator | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:&#34;creator&#34; |





 


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


 

 

 



<a name="resources_jobs_timeclock-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/timeclock.proto



<a name="resources-jobs-TimeclockEntry"></a>

### TimeclockEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job | [string](#string) |  |  |
| date | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| user_id | [int32](#int32) |  |  |
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





 

 

 

 



<a name="resources_laws_laws-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/laws/laws.proto



<a name="resources-laws-Law"></a>

### Law



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;law.id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| lawbook_id | [uint64](#uint64) |  |  |
| name | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| fine | [uint32](#uint32) | optional |  |
| detention_time | [uint32](#uint32) | optional |  |
| stvo_points | [uint32](#uint32) | optional |  |






<a name="resources-laws-LawBook"></a>

### LawBook



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| name | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize |
| laws | [Law](#resources-laws-Law) | repeated |  |





 

 

 

 



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
| icon | [string](#string) |  |  |






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
| color | [string](#string) | optional |  |
| icon | [string](#string) | optional |  |






<a name="resources-livemap-MarkerMarker"></a>

### MarkerMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| info | [MarkerInfo](#resources-livemap-MarkerInfo) |  |  |
| type | [MarkerType](#resources-livemap-MarkerType) |  | @gotags: alias:&#34;markerType&#34; |
| expires_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| data | [MarkerData](#resources-livemap-MarkerData) |  | @gotags: alias:&#34;markerData&#34; |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional |  |






<a name="resources-livemap-UserMarker"></a>

### UserMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| info | [MarkerInfo](#resources-livemap-MarkerInfo) |  |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:&#34;user&#34; |
| unit_id | [uint64](#uint64) | optional |  |
| unit | [resources.centrum.Unit](#resources-centrum-Unit) | optional |  |





 


<a name="resources-livemap-MarkerType"></a>

### MarkerType


| Name | Number | Description |
| ---- | ------ | ----------- |
| MARKER_TYPE_UNSPECIFIED | 0 |  |
| MARKER_TYPE_DOT | 1 |  |
| MARKER_TYPE_CIRCLE | 2 |  |
| MARKER_TYPE_ICON | 3 |  |


 

 

 



<a name="resources_livemap_tracker-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/tracker.proto



<a name="resources-livemap-UsersUpdateEvent"></a>

### UsersUpdateEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| added | [UserMarker](#resources-livemap-UserMarker) | repeated |  |
| removed | [UserMarker](#resources-livemap-UserMarker) | repeated |  |





 

 

 

 



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


 

 

 



<a name="resources_notifications_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/notifications/events.proto



<a name="resources-notifications-JobEvent"></a>

### JobEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_props | [resources.users.JobProps](#resources-users-JobProps) |  |  |






<a name="resources-notifications-SystemEvent"></a>

### SystemEvent







<a name="resources-notifications-UserEvent"></a>

### UserEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [bool](#bool) |  |  |
| notification | [Notification](#resources-notifications-Notification) |  | Notifications |
| messenger | [resources.messenger.MessengerEvent](#resources-messenger-MessengerEvent) |  | Messenger |





 

 

 

 



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
| strings | [string](#string) | repeated |  |





 

 

 

 



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
| title | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| data | [ExamQuestionData](#resources-qualifications-ExamQuestionData) |  |  |
| answer | [ExamQuestionAnswerData](#resources-qualifications-ExamQuestionAnswerData) | optional |  |






<a name="resources-qualifications-ExamQuestionAnswerData"></a>

### ExamQuestionAnswerData







<a name="resources-qualifications-ExamQuestionData"></a>

### ExamQuestionData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| separator | [ExamQuestionSeparator](#resources-qualifications-ExamQuestionSeparator) |  |  |
| yesno | [ExamQuestionYesNo](#resources-qualifications-ExamQuestionYesNo) |  |  |
| free_text | [ExamQuestionText](#resources-qualifications-ExamQuestionText) |  |  |
| single_choice | [ExamQuestionSingleChoice](#resources-qualifications-ExamQuestionSingleChoice) |  |  |
| multiple_choice | [ExamQuestionMultipleChoice](#resources-qualifications-ExamQuestionMultipleChoice) |  |  |






<a name="resources-qualifications-ExamQuestionMultipleChoice"></a>

### ExamQuestionMultipleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choices | [string](#string) | repeated |  |
| limit | [int32](#int32) | optional |  |






<a name="resources-qualifications-ExamQuestionSeparator"></a>

### ExamQuestionSeparator







<a name="resources-qualifications-ExamQuestionSingleChoice"></a>

### ExamQuestionSingleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choices | [string](#string) | repeated |  |






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
| choices | [string](#string) | repeated |  |






<a name="resources-qualifications-ExamResponseSeparator"></a>

### ExamResponseSeparator







<a name="resources-qualifications-ExamResponseSingleChoice"></a>

### ExamResponseSingleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| choice | [string](#string) |  |  |






<a name="resources-qualifications-ExamResponseText"></a>

### ExamResponseText



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  | 0.5 Megabyte |






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





 

 

 

 



<a name="resources_qualifications_qualifications-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/qualifications/qualifications.proto



<a name="resources-qualifications-Qualification"></a>

### Qualification



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| job | [string](#string) |  |  |
| weight | [uint32](#uint32) |  |  |
| closed | [bool](#bool) |  |  |
| abbreviation | [string](#string) |  | @sanitize: method=StripTags |
| title | [string](#string) |  | @sanitize |
| description | [string](#string) | optional | @sanitize: method=StripTags |
| content | [string](#string) |  | @sanitize |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |
| access | [QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |
| requirements | [QualificationRequirement](#resources-qualifications-QualificationRequirement) | repeated |  |
| discord_settings | [QualificationDiscordSettings](#resources-qualifications-QualificationDiscordSettings) | optional |  |
| exam_mode | [QualificationExamMode](#resources-qualifications-QualificationExamMode) |  |  |
| exam_settings | [QualificationExamSettings](#resources-qualifications-QualificationExamSettings) | optional |  |
| exam | [ExamQuestions](#resources-qualifications-ExamQuestions) | optional |  |
| result | [QualificationResult](#resources-qualifications-QualificationResult) | optional |  |
| request | [QualificationRequest](#resources-qualifications-QualificationRequest) | optional |  |






<a name="resources-qualifications-QualificationDiscordSettings"></a>

### QualificationDiscordSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sync_enabled | [bool](#bool) |  |  |
| role_name | [string](#string) | optional |  |






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
| qualification_id | [uint64](#uint64) |  |  |
| qualification | [QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:&#34;user&#34; |
| user_comment | [string](#string) | optional | @sanitize: method=StripTags |
| status | [RequestStatus](#resources-qualifications-RequestStatus) | optional |  |
| approved_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| approver_comment | [string](#string) | optional | @sanitize: method=StripTags |
| approver_id | [int32](#int32) | optional |  |
| approver | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;approver&#34; |
| approver_job | [string](#string) | optional |  |






<a name="resources-qualifications-QualificationRequirement"></a>

### QualificationRequirement



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| qualification_id | [uint64](#uint64) |  |  |
| target_qualification_id | [uint64](#uint64) |  |  |
| target_qualification | [QualificationShort](#resources-qualifications-QualificationShort) | optional | @gotags: alias:&#34;targetqualification.*&#34; |






<a name="resources-qualifications-QualificationResult"></a>

### QualificationResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| qualification_id | [uint64](#uint64) |  |  |
| qualification | [QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:&#34;user&#34; |
| status | [ResultStatus](#resources-qualifications-ResultStatus) |  |  |
| score | [uint32](#uint32) | optional |  |
| summary | [string](#string) |  | @sanitize: method=StripTags |
| creator_id | [int32](#int32) |  |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |






<a name="resources-qualifications-QualificationShort"></a>

### QualificationShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
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
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| creator_job | [string](#string) |  |  |
| requirements | [QualificationRequirement](#resources-qualifications-QualificationRequirement) | repeated |  |
| exam_mode | [QualificationExamMode](#resources-qualifications-QualificationExamMode) |  |  |
| exam_settings | [QualificationExamSettings](#resources-qualifications-QualificationExamSettings) | optional |  |
| result | [QualificationResult](#resources-qualifications-QualificationResult) | optional |  |





 


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
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| qualification_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| access | [AccessLevel](#resources-qualifications-AccessLevel) |  |  |





 


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



<a name="resources-qualifications-AccessLevelUpdateMode"></a>

### AccessLevelUpdateMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_UPDATE_MODE_UPDATE | 1 |  |
| ACCESS_LEVEL_UPDATE_MODE_DELETE | 2 |  |
| ACCESS_LEVEL_UPDATE_MODE_CLEAR | 3 |  |


 

 

 



<a name="resources_rector_audit-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/rector/audit.proto



<a name="resources-rector-AuditEntry"></a>

### AuditEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| user_id | [uint64](#uint64) |  | @gotags: alias:&#34;user_id&#34; |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| user_job | [string](#string) |  | @gotags: alias:&#34;user_job&#34; |
| target_user_id | [int32](#int32) | optional | @gotags: alias:&#34;target_user_id&#34; |
| target_user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| target_user_job | [string](#string) |  | @gotags: alias:&#34;target_user_job&#34; |
| service | [string](#string) |  | @gotags: alias:&#34;service&#34; |
| method | [string](#string) |  | @gotags: alias:&#34;method&#34; |
| state | [EventType](#resources-rector-EventType) |  | @gotags: alias:&#34;state&#34; |
| data | [string](#string) | optional | @gotags: alias:&#34;data&#34; |





 


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


 

 

 



<a name="resources_rector_config-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/rector/config.proto



<a name="resources-rector-AppConfig"></a>

### AppConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
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
| invite_url | [string](#string) | optional |  |
| ignored_jobs | [string](#string) | repeated |  |






<a name="resources-rector-JobInfo"></a>

### JobInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unemployed_job | [UnemployedJob](#resources-rector-UnemployedJob) |  |  |
| public_jobs | [string](#string) | repeated |  |
| hidden_jobs | [string](#string) | repeated |  |






<a name="resources-rector-Links"></a>

### Links



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| privacy_policy | [string](#string) | optional |  |
| imprint | [string](#string) | optional |  |






<a name="resources-rector-Perm"></a>

### Perm



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [string](#string) |  |  |
| name | [string](#string) |  |  |






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
| livemap_jobs | [string](#string) | repeated |  |






<a name="resources-rector-Website"></a>

### Website



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| links | [Links](#resources-rector-Links) |  |  |





 

 

 

 



<a name="resources_timestamp_timestamp-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/timestamp/timestamp.proto



<a name="resources-timestamp-Timestamp"></a>

### Timestamp
Timestamp for storage messages.  We&#39;ve defined a new local type wrapper
of google.protobuf.Timestamp so we can implement sql.Scanner and sql.Valuer
interfaces.  See:
https://golang.org/pkg/database/sql/#Scanner
https://golang.org/pkg/database/sql/driver/#Valuer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| timestamp | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 

 

 

 



<a name="resources_users_users-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/users.proto



<a name="resources-users-CitizenAttribute"></a>

### CitizenAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;id&#34; |
| job | [string](#string) | optional |  |
| name | [string](#string) |  |  |
| color | [string](#string) |  |  |






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
| user_id | [int32](#int32) |  | @gotags: alias:&#34;id&#34; |
| identifier | [string](#string) |  |  |
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
| props | [UserProps](#resources-users-UserProps) |  | @gotags: alias:&#34;fivenet_user_props&#34; |
| licenses | [License](#resources-users-License) | repeated | @gotags: alias:&#34;user_licenses&#34; |
| avatar | [resources.filestore.File](#resources-filestore-File) | optional |  |






<a name="resources-users-UserActivity"></a>

### UserActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | @gotags: alias:&#34;fivenet_user_activity.id&#34; |
| type | [UserActivityType](#resources-users-UserActivityType) |  | @gotags: alias:&#34;fivenet_user_activity.type&#34; |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: alias:&#34;fivenet_user_activity.created_at&#34; |
| source_user | [UserShort](#resources-users-UserShort) |  | @gotags: alias:&#34;source_user&#34; |
| target_user | [UserShort](#resources-users-UserShort) |  | @gotags: alias:&#34;target_user&#34; |
| key | [string](#string) |  | @sanitize

@gotags: alias:&#34;fivenet_user_activity.key&#34; |
| old_value | [string](#string) |  | @gotags: alias:&#34;fivenet_user_activity.old_value&#34; |
| new_value | [string](#string) |  | @gotags: alias:&#34;fivenet_user_activity.new_value&#34; |
| reason | [string](#string) |  | @sanitize

@gotags: alias:&#34;fivenet_user_activity.reason&#34; |






<a name="resources-users-UserProps"></a>

### UserProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| wanted | [bool](#bool) | optional |  |
| job_name | [string](#string) | optional | @gotags: alias:&#34;job&#34; |
| job | [Job](#resources-users-Job) | optional |  |
| job_grade_number | [int32](#int32) | optional | @gotags: alias:&#34;job_grade&#34; |
| job_grade | [JobGrade](#resources-users-JobGrade) | optional |  |
| traffic_infraction_points | [uint32](#uint32) | optional |  |
| open_fines | [int64](#int64) | optional |  |
| blood_type | [string](#string) | optional |  |
| mug_shot | [resources.filestore.File](#resources-filestore-File) | optional |  |
| attributes | [CitizenAttributes](#resources-users-CitizenAttributes) | optional |  |






<a name="resources-users-UserShort"></a>

### UserShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  | @gotags: alias:&#34;id&#34; |
| identifier | [string](#string) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional |  |
| job_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional |  |
| firstname | [string](#string) |  |  |
| lastname | [string](#string) |  |  |
| dateofbirth | [string](#string) |  |  |
| phone_number | [string](#string) | optional |  |
| avatar | [resources.filestore.File](#resources-filestore-File) | optional |  |





 


<a name="resources-users-UserActivityType"></a>

### UserActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_ACTIVITY_TYPE_UNSPECIFIED | 0 |  |
| USER_ACTIVITY_TYPE_CHANGED | 1 |  |
| USER_ACTIVITY_TYPE_MENTIONED | 2 |  |
| USER_ACTIVITY_TYPE_CREATED | 3 |  |


 

 

 



<a name="resources_users_jobs-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/jobs.proto



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
| ignored_role_ids | [string](#string) | repeated |  |






<a name="resources-users-Job"></a>

### Job



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | @gotags: sql:&#34;primary_key&#34; alias:&#34;name&#34; |
| label | [string](#string) |  |  |
| grades | [JobGrade](#resources-users-JobGrade) | repeated |  |






<a name="resources-users-JobGrade"></a>

### JobGrade



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| job_name | [string](#string) | optional |  |
| grade | [int32](#int32) |  |  |
| label | [string](#string) |  |  |






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





 


<a name="resources-users-UserInfoSyncUnemployedMode"></a>

### UserInfoSyncUnemployedMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_INFO_SYNC_UNEMPLOYED_MODE_UNSPECIFIED | 0 |  |
| USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE | 1 |  |
| USER_INFO_SYNC_UNEMPLOYED_MODE_KICK | 2 |  |


 

 

 



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





 

 

 

 



<a name="resources_messenger_message-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/messenger/message.proto



<a name="resources-messenger-Message"></a>

### Message



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| thread_id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| message | [string](#string) |  | @sanitize: method=StripTags |
| data | [MessageData](#resources-messenger-MessageData) | optional |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |






<a name="resources-messenger-MessageData"></a>

### MessageData
TODO allow links to internal





 

 

 

 



<a name="resources_messenger_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/messenger/access.proto



<a name="resources-messenger-ThreadAccess"></a>

### ThreadAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [ThreadJobAccess](#resources-messenger-ThreadJobAccess) | repeated | @gotags: alias:&#34;job_access&#34; |
| users | [ThreadUserAccess](#resources-messenger-ThreadUserAccess) | repeated | @gotags: alias:&#34;user_access&#34; |






<a name="resources-messenger-ThreadJobAccess"></a>

### ThreadJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| thread_id | [uint64](#uint64) |  |  |
| job | [string](#string) |  |  |
| job_label | [string](#string) | optional | @gotags: alias:&#34;job_label&#34; |
| minimum_grade | [int32](#int32) |  |  |
| job_grade_label | [string](#string) | optional | @gotags: alias:&#34;job_grade_label&#34; |
| access | [AccessLevel](#resources-messenger-AccessLevel) |  |  |






<a name="resources-messenger-ThreadUserAccess"></a>

### ThreadUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| thread_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| user | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| access | [AccessLevel](#resources-messenger-AccessLevel) |  |  |





 


<a name="resources-messenger-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_BLOCKED | 1 |  |
| ACCESS_LEVEL_VIEW | 2 |  |
| ACCESS_LEVEL_MESSAGE | 3 |  |
| ACCESS_LEVEL_MANAGE | 4 |  |
| ACCESS_LEVEL_ADMIN | 5 |  |



<a name="resources-messenger-AccessLevelUpdateMode"></a>

### AccessLevelUpdateMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED | 0 |  |
| ACCESS_LEVEL_UPDATE_MODE_UPDATE | 1 |  |
| ACCESS_LEVEL_UPDATE_MODE_DELETE | 2 |  |
| ACCESS_LEVEL_UPDATE_MODE_CLEAR | 3 |  |


 

 

 



<a name="resources_messenger_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/messenger/events.proto



<a name="resources-messenger-MessengerEvent"></a>

### MessengerEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_update | [Thread](#resources-messenger-Thread) |  |  |
| thread_delete | [uint64](#uint64) |  |  |
| message_update | [Message](#resources-messenger-Message) |  |  |
| message_delete | [uint64](#uint64) |  |  |





 

 

 

 



<a name="resources_messenger_thread-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/messenger/thread.proto



<a name="resources-messenger-Thread"></a>

### Thread



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  |  |
| created_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| updated_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| deleted_at | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| title | [string](#string) |  | @sanitize |
| archived | [bool](#bool) |  |  |
| last_message | [Message](#resources-messenger-Message) | optional |  |
| user_state | [ThreadUserState](#resources-messenger-ThreadUserState) |  |  |
| creator_job | [string](#string) |  |  |
| creator_id | [int32](#int32) | optional |  |
| creator | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:&#34;creator&#34; |
| access | [ThreadAccess](#resources-messenger-ThreadAccess) |  |  |






<a name="resources-messenger-ThreadUserState"></a>

### ThreadUserState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |
| user_id | [int32](#int32) |  |  |
| unread | [bool](#bool) |  |  |
| last_read | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| important | [bool](#bool) |  |  |
| favorite | [bool](#bool) |  |  |
| muted | [bool](#bool) |  |  |





 

 

 

 



<a name="resources_messenger_user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/messenger/user.proto



<a name="resources-messenger-UserStatus"></a>

### UserStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| last_seen | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| status | [string](#string) | optional | @sanitize: method=StripTags |





 

 

 

 



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
| char | [resources.users.User](#resources-users-User) |  |  |
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





 

 

 


<a name="services-calendar-CalendarService"></a>

### CalendarService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListCalendars | [ListCalendarsRequest](#services-calendar-ListCalendarsRequest) | [ListCalendarsResponse](#services-calendar-ListCalendarsResponse) | @perm: Name=Any |
| GetCalendar | [GetCalendarRequest](#services-calendar-GetCalendarRequest) | [GetCalendarResponse](#services-calendar-GetCalendarResponse) | @perm: Name=Any |
| CreateOrUpdateCalendar | [CreateOrUpdateCalendarRequest](#services-calendar-CreateOrUpdateCalendarRequest) | [CreateOrUpdateCalendarResponse](#services-calendar-CreateOrUpdateCalendarResponse) | @perm: Attrs=Fields/StringList:[]string{&#34;Job&#34;, &#34;Public&#34;} |
| DeleteCalendar | [DeleteCalendarRequest](#services-calendar-DeleteCalendarRequest) | [DeleteCalendarResponse](#services-calendar-DeleteCalendarResponse) | @perm |
| ListCalendarEntries | [ListCalendarEntriesRequest](#services-calendar-ListCalendarEntriesRequest) | [ListCalendarEntriesResponse](#services-calendar-ListCalendarEntriesResponse) | @perm: Name=Any |
| GetCalendarEntry | [GetCalendarEntryRequest](#services-calendar-GetCalendarEntryRequest) | [GetCalendarEntryResponse](#services-calendar-GetCalendarEntryResponse) | @perm: Name=Any |
| CreateOrUpdateCalendarEntry | [CreateOrUpdateCalendarEntryRequest](#services-calendar-CreateOrUpdateCalendarEntryRequest) | [CreateOrUpdateCalendarEntryResponse](#services-calendar-CreateOrUpdateCalendarEntryResponse) | @perm |
| DeleteCalendarEntry | [DeleteCalendarEntryRequest](#services-calendar-DeleteCalendarEntryRequest) | [DeleteCalendarEntryResponse](#services-calendar-DeleteCalendarEntryResponse) | @perm |
| ShareCalendarEntry | [ShareCalendarEntryRequest](#services-calendar-ShareCalendarEntryRequest) | [ShareCalendarEntryResponse](#services-calendar-ShareCalendarEntryResponse) | @perm: Name=CreateOrUpdateCalendarEntry |
| ListCalendarEntryRSVP | [ListCalendarEntryRSVPRequest](#services-calendar-ListCalendarEntryRSVPRequest) | [ListCalendarEntryRSVPResponse](#services-calendar-ListCalendarEntryRSVPResponse) | @perm: Name=Any |
| RSVPCalendarEntry | [RSVPCalendarEntryRequest](#services-calendar-RSVPCalendarEntryRequest) | [RSVPCalendarEntryResponse](#services-calendar-RSVPCalendarEntryResponse) | @perm: Name=Any |
| ListSubscriptions | [ListSubscriptionsRequest](#services-calendar-ListSubscriptionsRequest) | [ListSubscriptionsResponse](#services-calendar-ListSubscriptionsResponse) | @perm: Name=Any |
| SubscribeToCalendar | [SubscribeToCalendarRequest](#services-calendar-SubscribeToCalendarRequest) | [SubscribeToCalendarResponse](#services-calendar-SubscribeToCalendarResponse) | @perm: Name=Any |

 



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
| disponents | [resources.users.UserShort](#resources-users-UserShort) | repeated |  |
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





 

 

 


<a name="services-citizenstore-CitizenStoreService"></a>

### CitizenStoreService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListCitizens | [ListCitizensRequest](#services-citizenstore-ListCitizensRequest) | [ListCitizensResponse](#services-citizenstore-ListCitizensResponse) | @perm: Attrs=Fields/StringList:[]string{&#34;PhoneNumber&#34;, &#34;Licenses&#34;, &#34;UserProps.Wanted&#34;, &#34;UserProps.Job&#34;, &#34;UserProps.TrafficInfractionPoints&#34;, &#34;UserProps.OpenFines&#34;, &#34;UserProps.BloodType&#34;, &#34;UserProps.MugShot&#34;, &#34;UserProps.Attributes&#34;} |
| GetUser | [GetUserRequest](#services-citizenstore-GetUserRequest) | [GetUserResponse](#services-citizenstore-GetUserResponse) | @perm: Attrs=Jobs/JobGradeList |
| ListUserActivity | [ListUserActivityRequest](#services-citizenstore-ListUserActivityRequest) | [ListUserActivityResponse](#services-citizenstore-ListUserActivityResponse) | @perm: Attrs=Fields/StringList:[]string{&#34;SourceUser&#34;, &#34;Own&#34;} |
| SetUserProps | [SetUserPropsRequest](#services-citizenstore-SetUserPropsRequest) | [SetUserPropsResponse](#services-citizenstore-SetUserPropsResponse) | @perm: Attrs=Fields/StringList:[]string{&#34;Wanted&#34;, &#34;Job&#34;, &#34;TrafficInfractionPoints&#34;, &#34;MugShot&#34;, &#34;Attributes&#34;} |
| SetProfilePicture | [SetProfilePictureRequest](#services-citizenstore-SetProfilePictureRequest) | [SetProfilePictureResponse](#services-citizenstore-SetProfilePictureResponse) | @perm: Name=Any |
| ManageCitizenAttributes | [ManageCitizenAttributesRequest](#services-citizenstore-ManageCitizenAttributesRequest) | [ManageCitizenAttributesResponse](#services-citizenstore-ManageCitizenAttributesResponse) | @perm |

 



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
| users | [resources.users.UserShort](#resources-users-UserShort) | repeated | @gotags: alias:&#34;user&#34; |






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





 

 

 


<a name="services-completor-CompletorService"></a>

### CompletorService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CompleteCitizens | [CompleteCitizensRequest](#services-completor-CompleteCitizensRequest) | [CompleteCitizensRespoonse](#services-completor-CompleteCitizensRespoonse) | @perm |
| CompleteJobs | [CompleteJobsRequest](#services-completor-CompleteJobsRequest) | [CompleteJobsResponse](#services-completor-CompleteJobsResponse) | @perm |
| CompleteDocumentCategories | [CompleteDocumentCategoriesRequest](#services-completor-CompleteDocumentCategoriesRequest) | [CompleteDocumentCategoriesResponse](#services-completor-CompleteDocumentCategoriesResponse) | @perm: Attrs=Jobs/JobList |
| ListLawBooks | [ListLawBooksRequest](#services-completor-ListLawBooksRequest) | [ListLawBooksResponse](#services-completor-ListLawBooksResponse) | @perm: Name=Any |
| CompleteCitizenAttributes | [CompleteCitizenAttributesRequest](#services-completor-CompleteCitizenAttributesRequest) | [CompleteCitizenAttributesResponse](#services-completor-CompleteCitizenAttributesResponse) | @perm: Attrs=Jobs/JobList |

 



<a name="services_dmv_vehicles-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/dmv/vehicles.proto



<a name="services-dmv-ListVehiclesRequest"></a>

### ListVehiclesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| orderBy | [resources.common.database.OrderBy](#resources-common-database-OrderBy) | repeated |  |
| license_plate | [string](#string) | optional | Search params |
| model | [string](#string) | optional |  |
| user_id | [int32](#int32) | optional |  |






<a name="services-dmv-ListVehiclesResponse"></a>

### ListVehiclesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| vehicles | [resources.vehicles.Vehicle](#resources-vehicles-Vehicle) | repeated |  |





 

 

 


<a name="services-dmv-DMVService"></a>

### DMVService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListVehicles | [ListVehiclesRequest](#services-dmv-ListVehiclesRequest) | [ListVehiclesResponse](#services-dmv-ListVehiclesResponse) | @perm |

 



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
| category_id | [uint64](#uint64) | optional | @gotags: alias:&#34;category_id&#34; |
| title | [string](#string) |  | @sanitize: method=StripTags

@gotags: alias:&#34;title&#34; |
| content | [string](#string) |  | @sanitize

@gotags: alias:&#34;content&#34; |
| content_type | [resources.documents.DocContentType](#resources-documents-DocContentType) |  | @gotags: alias:&#34;content_type&#34; |
| data | [string](#string) | optional | @gotags: alias:&#34;data&#34; |
| state | [string](#string) |  | @sanitize

@gotags: alias:&#34;state&#34; |
| closed | [bool](#bool) |  | @gotags: alias:&#34;closed&#34; |
| public | [bool](#bool) |  | @gotags: alias:&#34;public&#34; |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) | optional |  |
| template_id | [uint64](#uint64) | optional |  |






<a name="services-docstore-CreateDocumentResponse"></a>

### CreateDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |






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
| document_id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |






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
Comments ===============================================================


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
Access =====================================================================


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
| references | [resources.documents.DocumentReference](#resources-documents-DocumentReference) | repeated | @gotags: alias:&#34;reference&#34; |






<a name="services-docstore-GetDocumentRelationsRequest"></a>

### GetDocumentRelationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  |  |






<a name="services-docstore-GetDocumentRelationsResponse"></a>

### GetDocumentRelationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| relations | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) | repeated | @gotags: alias:&#34;relation&#34; |






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
Categories






<a name="services-docstore-ListCategoriesResponse"></a>

### ListCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [resources.documents.Category](#resources-documents-Category) | repeated |  |






<a name="services-docstore-ListDocumentActivityRequest"></a>

### ListDocumentActivityRequest
Document Activity and Requests =============================================


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| document_id | [uint64](#uint64) |  |  |
| activity_types | [resources.documents.DocActivityType](#resources-documents-DocActivityType) | repeated | Search |






<a name="services-docstore-ListDocumentActivityResponse"></a>

### ListDocumentActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| activity | [resources.documents.DocActivity](#resources-documents-DocActivity) | repeated |  |






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
Documents ==================================================================


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| orderBy | [resources.common.database.OrderBy](#resources-common-database-OrderBy) | repeated |  |
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
Templates ==================================================================






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
| mode | [resources.documents.AccessLevelUpdateMode](#resources-documents-AccessLevelUpdateMode) |  |  |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-docstore-SetDocumentAccessResponse"></a>

### SetDocumentAccessResponse







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
| document_id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |
| category_id | [uint64](#uint64) | optional | @gotags: alias:&#34;category_id&#34; |
| title | [string](#string) |  | @sanitize: method=StripTags

@gotags: alias:&#34;title&#34; |
| content | [string](#string) |  | @sanitize

@gotags: alias:&#34;content&#34; |
| content_type | [resources.documents.DocContentType](#resources-documents-DocContentType) |  | @gotags: alias:&#34;content_type&#34; |
| data | [string](#string) | optional | @gotags: alias:&#34;data&#34; |
| state | [string](#string) |  | @sanitize

@gotags: alias:&#34;state&#34; |
| closed | [bool](#bool) |  | @gotags: alias:&#34;closed&#34; |
| public | [bool](#bool) |  | @gotags: alias:&#34;public&#34; |
| access | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) | optional |  |






<a name="services-docstore-UpdateDocumentResponse"></a>

### UpdateDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| document_id | [uint64](#uint64) |  | @gotags: alias:&#34;id&#34; |






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
| GetDocument | [GetDocumentRequest](#services-docstore-GetDocumentRequest) | [GetDocumentResponse](#services-docstore-GetDocumentResponse) | @perm |
| CreateDocument | [CreateDocumentRequest](#services-docstore-CreateDocumentRequest) | [CreateDocumentResponse](#services-docstore-CreateDocumentResponse) | @perm |
| UpdateDocument | [UpdateDocumentRequest](#services-docstore-UpdateDocumentRequest) | [UpdateDocumentResponse](#services-docstore-UpdateDocumentResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| DeleteDocument | [DeleteDocumentRequest](#services-docstore-DeleteDocumentRequest) | [DeleteDocumentResponse](#services-docstore-DeleteDocumentResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| ToggleDocument | [ToggleDocumentRequest](#services-docstore-ToggleDocumentRequest) | [ToggleDocumentResponse](#services-docstore-ToggleDocumentResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| ChangeDocumentOwner | [ChangeDocumentOwnerRequest](#services-docstore-ChangeDocumentOwnerRequest) | [ChangeDocumentOwnerResponse](#services-docstore-ChangeDocumentOwnerResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| GetDocumentReferences | [GetDocumentReferencesRequest](#services-docstore-GetDocumentReferencesRequest) | [GetDocumentReferencesResponse](#services-docstore-GetDocumentReferencesResponse) | @perm: Name=GetDocument |
| GetDocumentRelations | [GetDocumentRelationsRequest](#services-docstore-GetDocumentRelationsRequest) | [GetDocumentRelationsResponse](#services-docstore-GetDocumentRelationsResponse) | @perm: Name=GetDocument |
| AddDocumentReference | [AddDocumentReferenceRequest](#services-docstore-AddDocumentReferenceRequest) | [AddDocumentReferenceResponse](#services-docstore-AddDocumentReferenceResponse) | @perm |
| RemoveDocumentReference | [RemoveDocumentReferenceRequest](#services-docstore-RemoveDocumentReferenceRequest) | [RemoveDocumentReferenceResponse](#services-docstore-RemoveDocumentReferenceResponse) | @perm: Name=AddDocumentReference |
| AddDocumentRelation | [AddDocumentRelationRequest](#services-docstore-AddDocumentRelationRequest) | [AddDocumentRelationResponse](#services-docstore-AddDocumentRelationResponse) | @perm |
| RemoveDocumentRelation | [RemoveDocumentRelationRequest](#services-docstore-RemoveDocumentRelationRequest) | [RemoveDocumentRelationResponse](#services-docstore-RemoveDocumentRelationResponse) | @perm: Name=AddDocumentRelation |
| GetComments | [GetCommentsRequest](#services-docstore-GetCommentsRequest) | [GetCommentsResponse](#services-docstore-GetCommentsResponse) | @perm: Name=GetDocument |
| PostComment | [PostCommentRequest](#services-docstore-PostCommentRequest) | [PostCommentResponse](#services-docstore-PostCommentResponse) | @perm |
| EditComment | [EditCommentRequest](#services-docstore-EditCommentRequest) | [EditCommentResponse](#services-docstore-EditCommentResponse) | @perm: Name=PostComment |
| DeleteComment | [DeleteCommentRequest](#services-docstore-DeleteCommentRequest) | [DeleteCommentResponse](#services-docstore-DeleteCommentResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| GetDocumentAccess | [GetDocumentAccessRequest](#services-docstore-GetDocumentAccessRequest) | [GetDocumentAccessResponse](#services-docstore-GetDocumentAccessResponse) | @perm: Name=GetDocument |
| SetDocumentAccess | [SetDocumentAccessRequest](#services-docstore-SetDocumentAccessRequest) | [SetDocumentAccessResponse](#services-docstore-SetDocumentAccessResponse) | @perm: Name=CreateDocument |
| ListDocumentActivity | [ListDocumentActivityRequest](#services-docstore-ListDocumentActivityRequest) | [ListDocumentActivityResponse](#services-docstore-ListDocumentActivityResponse) | @perm |
| ListDocumentReqs | [ListDocumentReqsRequest](#services-docstore-ListDocumentReqsRequest) | [ListDocumentReqsResponse](#services-docstore-ListDocumentReqsResponse) | @perm |
| CreateDocumentReq | [CreateDocumentReqRequest](#services-docstore-CreateDocumentReqRequest) | [CreateDocumentReqResponse](#services-docstore-CreateDocumentReqResponse) | @perm: Attrs=Types/StringList:[]string{&#34;Access&#34;, &#34;Closure&#34;, &#34;Update&#34;, &#34;Deletion&#34;, &#34;OwnerChange&#34;} |
| UpdateDocumentReq | [UpdateDocumentReqRequest](#services-docstore-UpdateDocumentReqRequest) | [UpdateDocumentReqResponse](#services-docstore-UpdateDocumentReqResponse) | @perm: Name=CreateDocumentReq |
| DeleteDocumentReq | [DeleteDocumentReqRequest](#services-docstore-DeleteDocumentReqRequest) | [DeleteDocumentReqResponse](#services-docstore-DeleteDocumentReqResponse) | @perm |
| ListUserDocuments | [ListUserDocumentsRequest](#services-docstore-ListUserDocumentsRequest) | [ListUserDocumentsResponse](#services-docstore-ListUserDocumentsResponse) | @perm |
| ListCategories | [ListCategoriesRequest](#services-docstore-ListCategoriesRequest) | [ListCategoriesResponse](#services-docstore-ListCategoriesResponse) | @perm |
| CreateCategory | [CreateCategoryRequest](#services-docstore-CreateCategoryRequest) | [CreateCategoryResponse](#services-docstore-CreateCategoryResponse) | @perm |
| UpdateCategory | [UpdateCategoryRequest](#services-docstore-UpdateCategoryRequest) | [UpdateCategoryResponse](#services-docstore-UpdateCategoryResponse) | @perm: Name=CreateCategory |
| DeleteCategory | [DeleteCategoryRequest](#services-docstore-DeleteCategoryRequest) | [DeleteCategoryResponse](#services-docstore-DeleteCategoryResponse) | @perm |

 



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





 

 

 


<a name="services-jobs-JobsConductService"></a>

### JobsConductService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListConductEntries | [ListConductEntriesRequest](#services-jobs-ListConductEntriesRequest) | [ListConductEntriesResponse](#services-jobs-ListConductEntriesResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;All&#34;} |
| CreateConductEntry | [CreateConductEntryRequest](#services-jobs-CreateConductEntryRequest) | [CreateConductEntryResponse](#services-jobs-CreateConductEntryResponse) | @perm |
| UpdateConductEntry | [UpdateConductEntryRequest](#services-jobs-UpdateConductEntryRequest) | [UpdateConductEntryResponse](#services-jobs-UpdateConductEntryResponse) | @perm |
| DeleteConductEntry | [DeleteConductEntryRequest](#services-jobs-DeleteConductEntryRequest) | [DeleteConductEntryResponse](#services-jobs-DeleteConductEntryResponse) | @perm |

 



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
| days | [int32](#int32) |  |  |






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
| user_ids | [int32](#int32) | repeated | Search |
| from | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| to | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| per_day | [bool](#bool) | optional |  |






<a name="services-jobs-ListTimeclockResponse"></a>

### ListTimeclockResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| entries | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| stats | [resources.jobs.TimeclockStats](#resources-jobs-TimeclockStats) |  |  |
| weekly | [resources.jobs.TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats) | repeated |  |





 

 

 


<a name="services-jobs-JobsTimeclockService"></a>

### JobsTimeclockService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListTimeclock | [ListTimeclockRequest](#services-jobs-ListTimeclockRequest) | [ListTimeclockResponse](#services-jobs-ListTimeclockResponse) | @perm: Attrs=Access/StringList:[]string{&#34;All&#34;} |
| GetTimeclockStats | [GetTimeclockStatsRequest](#services-jobs-GetTimeclockStatsRequest) | [GetTimeclockStatsResponse](#services-jobs-GetTimeclockStatsResponse) | @perm: Name=ListTimeclock |
| ListInactiveEmployees | [ListInactiveEmployeesRequest](#services-jobs-ListInactiveEmployeesRequest) | [ListInactiveEmployeesResponse](#services-jobs-ListInactiveEmployeesResponse) | @perm |

 



<a name="services_jobs_jobs-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/jobs.proto



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
| user_ids | [int32](#int32) | repeated |  |
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
| search | [string](#string) |  | Search params |
| user_id | [int32](#int32) | optional |  |
| absent | [bool](#bool) | optional |  |






<a name="services-jobs-ListColleaguesResponse"></a>

### ListColleaguesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| colleagues | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |






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





 

 

 


<a name="services-jobs-JobsService"></a>

### JobsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListColleagues | [ListColleaguesRequest](#services-jobs-ListColleaguesRequest) | [ListColleaguesResponse](#services-jobs-ListColleaguesResponse) | @perm |
| GetSelf | [GetSelfRequest](#services-jobs-GetSelfRequest) | [GetSelfResponse](#services-jobs-GetSelfResponse) | @perm: Name=ListColleagues |
| GetColleague | [GetColleagueRequest](#services-jobs-GetColleagueRequest) | [GetColleagueResponse](#services-jobs-GetColleagueResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;}|Types/StringList:[]string{&#34;Note&#34;} |
| ListColleagueActivity | [ListColleagueActivityRequest](#services-jobs-ListColleagueActivityRequest) | [ListColleagueActivityResponse](#services-jobs-ListColleagueActivityResponse) | @perm: Attrs=Types/StringList:[]string{&#34;HIRED&#34;, &#34;FIRED&#34;, &#34;PROMOTED&#34;, &#34;DEMOTED&#34;, &#34;ABSENCE_DATE&#34;, &#34;NOTE&#34;} |
| SetJobsUserProps | [SetJobsUserPropsRequest](#services-jobs-SetJobsUserPropsRequest) | [SetJobsUserPropsResponse](#services-jobs-SetJobsUserPropsResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;}|Types/StringList:[]string{&#34;AbsenceDate&#34;,&#34;Note&#34;} |
| GetMOTD | [GetMOTDRequest](#services-jobs-GetMOTDRequest) | [GetMOTDResponse](#services-jobs-GetMOTDResponse) | @perm: Name=Any |
| SetMOTD | [SetMOTDRequest](#services-jobs-SetMOTDRequest) | [SetMOTDResponse](#services-jobs-SetMOTDResponse) | @perm |

 



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






<a name="services-livemapper-UserMarkersUpdates"></a>

### UserMarkersUpdates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [resources.livemap.UserMarker](#resources-livemap-UserMarker) | repeated |  |
| part | [int32](#int32) |  |  |





 

 

 


<a name="services-livemapper-LivemapperService"></a>

### LivemapperService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Stream | [StreamRequest](#services-livemapper-StreamRequest) | [StreamResponse](#services-livemapper-StreamResponse) stream | @perm: Attrs=Markers/JobList|Players/JobGradeList |
| CreateOrUpdateMarker | [CreateOrUpdateMarkerRequest](#services-livemapper-CreateOrUpdateMarkerRequest) | [CreateOrUpdateMarkerResponse](#services-livemapper-CreateOrUpdateMarkerResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| DeleteMarker | [DeleteMarkerRequest](#services-livemapper-DeleteMarkerRequest) | [DeleteMarkerResponse](#services-livemapper-DeleteMarkerResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |

 



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
| system_event | [resources.notifications.SystemEvent](#resources-notifications-SystemEvent) |  |  |





 

 

 


<a name="services-notificator-NotificatorService"></a>

### NotificatorService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetNotifications | [GetNotificationsRequest](#services-notificator-GetNotificationsRequest) | [GetNotificationsResponse](#services-notificator-GetNotificationsResponse) | @perm: Name=Any |
| MarkNotifications | [MarkNotificationsRequest](#services-notificator-MarkNotificationsRequest) | [MarkNotificationsResponse](#services-notificator-MarkNotificationsResponse) | @perm: Name=Any |
| Stream | [StreamRequest](#services-notificator-StreamRequest) | [StreamResponse](#services-notificator-StreamResponse) stream | @perm: Name=Any |

 



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
| user | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) | optional |  |






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
| search | [string](#string) | optional | Search params |






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
| mode | [resources.qualifications.AccessLevelUpdateMode](#resources-qualifications-AccessLevelUpdateMode) |  |  |
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





 

 

 


<a name="services-qualifications-QualificationsService"></a>

### QualificationsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListQualifications | [ListQualificationsRequest](#services-qualifications-ListQualificationsRequest) | [ListQualificationsResponse](#services-qualifications-ListQualificationsResponse) | @perm |
| GetQualification | [GetQualificationRequest](#services-qualifications-GetQualificationRequest) | [GetQualificationResponse](#services-qualifications-GetQualificationResponse) | @perm |
| CreateQualification | [CreateQualificationRequest](#services-qualifications-CreateQualificationRequest) | [CreateQualificationResponse](#services-qualifications-CreateQualificationResponse) | @perm |
| UpdateQualification | [UpdateQualificationRequest](#services-qualifications-UpdateQualificationRequest) | [UpdateQualificationResponse](#services-qualifications-UpdateQualificationResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| DeleteQualification | [DeleteQualificationRequest](#services-qualifications-DeleteQualificationRequest) | [DeleteQualificationResponse](#services-qualifications-DeleteQualificationResponse) | @perm: Attrs=Access/StringList:[]string{&#34;Own&#34;, &#34;Lower_Rank&#34;, &#34;Same_Rank&#34;, &#34;Any&#34;} |
| ListQualificationRequests | [ListQualificationRequestsRequest](#services-qualifications-ListQualificationRequestsRequest) | [ListQualificationRequestsResponse](#services-qualifications-ListQualificationRequestsResponse) | @perm: Name=GetQualification |
| CreateOrUpdateQualificationRequest | [CreateOrUpdateQualificationRequestRequest](#services-qualifications-CreateOrUpdateQualificationRequestRequest) | [CreateOrUpdateQualificationRequestResponse](#services-qualifications-CreateOrUpdateQualificationRequestResponse) | @perm: Name=GetQualification |
| DeleteQualificationReq | [DeleteQualificationReqRequest](#services-qualifications-DeleteQualificationReqRequest) | [DeleteQualificationReqResponse](#services-qualifications-DeleteQualificationReqResponse) | @perm |
| ListQualificationsResults | [ListQualificationsResultsRequest](#services-qualifications-ListQualificationsResultsRequest) | [ListQualificationsResultsResponse](#services-qualifications-ListQualificationsResultsResponse) | @perm: Name=GetQualification |
| CreateOrUpdateQualificationResult | [CreateOrUpdateQualificationResultRequest](#services-qualifications-CreateOrUpdateQualificationResultRequest) | [CreateOrUpdateQualificationResultResponse](#services-qualifications-CreateOrUpdateQualificationResultResponse) | @perm |
| DeleteQualificationResult | [DeleteQualificationResultRequest](#services-qualifications-DeleteQualificationResultRequest) | [DeleteQualificationResultResponse](#services-qualifications-DeleteQualificationResultResponse) | @perm |
| GetExamInfo | [GetExamInfoRequest](#services-qualifications-GetExamInfoRequest) | [GetExamInfoResponse](#services-qualifications-GetExamInfoResponse) | @perm: Name=GetQualification |
| TakeExam | [TakeExamRequest](#services-qualifications-TakeExamRequest) | [TakeExamResponse](#services-qualifications-TakeExamResponse) | @perm: Name=GetQualification |
| SubmitExam | [SubmitExamRequest](#services-qualifications-SubmitExamRequest) | [SubmitExamResponse](#services-qualifications-SubmitExamResponse) | @perm: Name=GetQualification |
| GetUserExam | [GetUserExamRequest](#services-qualifications-GetUserExamRequest) | [GetUserExamResponse](#services-qualifications-GetUserExamResponse) | @perm: Name=CreateOrUpdateQualificationResult |

 



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





 

 

 


<a name="services-rector-RectorConfigService"></a>

### RectorConfigService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAppConfig | [GetAppConfigRequest](#services-rector-GetAppConfigRequest) | [GetAppConfigResponse](#services-rector-GetAppConfigResponse) | @perm: Name=SuperUser |
| UpdateAppConfig | [UpdateAppConfigRequest](#services-rector-UpdateAppConfigRequest) | [UpdateAppConfigResponse](#services-rector-UpdateAppConfigResponse) | @perm: Name=SuperUser |

 



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





 

 

 


<a name="services-rector-RectorFilestoreService"></a>

### RectorFilestoreService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListFiles | [ListFilesRequest](#services-rector-ListFilesRequest) | [ListFilesResponse](#services-rector-ListFilesResponse) | @perm: Name=SuperUser |
| UploadFile | [UploadFileRequest](#services-rector-UploadFileRequest) | [UploadFileResponse](#services-rector-UploadFileResponse) | @perm: Name=SuperUser |
| DeleteFile | [DeleteFileRequest](#services-rector-DeleteFileRequest) | [DeleteFileResponse](#services-rector-DeleteFileResponse) | @perm: Name=SuperUser |

 



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






 

 

 


<a name="services-rector-RectorLawsService"></a>

### RectorLawsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateOrUpdateLawBook | [CreateOrUpdateLawBookRequest](#services-rector-CreateOrUpdateLawBookRequest) | [CreateOrUpdateLawBookResponse](#services-rector-CreateOrUpdateLawBookResponse) | @perm: Name=SuperUser |
| DeleteLawBook | [DeleteLawBookRequest](#services-rector-DeleteLawBookRequest) | [DeleteLawBookResponse](#services-rector-DeleteLawBookResponse) | @perm: Name=SuperUser |
| CreateOrUpdateLaw | [CreateOrUpdateLawRequest](#services-rector-CreateOrUpdateLawRequest) | [CreateOrUpdateLawResponse](#services-rector-CreateOrUpdateLawResponse) | @perm: Name=SuperUser |
| DeleteLaw | [DeleteLawRequest](#services-rector-DeleteLawRequest) | [DeleteLawResponse](#services-rector-DeleteLawResponse) | @perm: Name=SuperUser |

 



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
| user_ids | [int32](#int32) | repeated |  |
| from | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| to | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| service | [string](#string) | optional |  |
| method | [string](#string) | optional |  |
| search | [string](#string) | optional |  |






<a name="services-rector-ViewAuditLogResponse"></a>

### ViewAuditLogResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| logs | [resources.rector.AuditEntry](#resources-rector-AuditEntry) | repeated |  |





 

 

 


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

 



<a name="services_messenger_messenger-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/messenger/messenger.proto



<a name="services-messenger-CreateOrUpdateThreadRequest"></a>

### CreateOrUpdateThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread | [resources.messenger.Thread](#resources-messenger-Thread) |  |  |






<a name="services-messenger-CreateOrUpdateThreadResponse"></a>

### CreateOrUpdateThreadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread | [resources.messenger.Thread](#resources-messenger-Thread) |  |  |






<a name="services-messenger-DeleteMessageRequest"></a>

### DeleteMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |
| message_id | [uint64](#uint64) |  |  |






<a name="services-messenger-DeleteMessageResponse"></a>

### DeleteMessageResponse







<a name="services-messenger-DeleteThreadRequest"></a>

### DeleteThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |






<a name="services-messenger-DeleteThreadResponse"></a>

### DeleteThreadResponse







<a name="services-messenger-GetThreadMessagesRequest"></a>

### GetThreadMessagesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |
| after | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="services-messenger-GetThreadMessagesResponse"></a>

### GetThreadMessagesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| messages | [resources.messenger.Message](#resources-messenger-Message) | repeated |  |






<a name="services-messenger-GetThreadRequest"></a>

### GetThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |






<a name="services-messenger-GetThreadResponse"></a>

### GetThreadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread | [resources.messenger.Thread](#resources-messenger-Thread) |  |  |






<a name="services-messenger-LeaveThreadRequest"></a>

### LeaveThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thread_id | [uint64](#uint64) |  |  |






<a name="services-messenger-LeaveThreadResponse"></a>

### LeaveThreadResponse







<a name="services-messenger-ListThreadsRequest"></a>

### ListThreadsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| after | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-messenger-ListThreadsResponse"></a>

### ListThreadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| threads | [resources.messenger.Thread](#resources-messenger-Thread) | repeated |  |






<a name="services-messenger-PostMessageRequest"></a>

### PostMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [resources.messenger.Message](#resources-messenger-Message) |  |  |






<a name="services-messenger-PostMessageResponse"></a>

### PostMessageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [resources.messenger.Message](#resources-messenger-Message) |  |  |






<a name="services-messenger-SetThreadUserStateRequest"></a>

### SetThreadUserStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [resources.messenger.ThreadUserState](#resources-messenger-ThreadUserState) |  |  |






<a name="services-messenger-SetThreadUserStateResponse"></a>

### SetThreadUserStateResponse






 

 

 


<a name="services-messenger-MessengerService"></a>

### MessengerService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListThreads | [ListThreadsRequest](#services-messenger-ListThreadsRequest) | [ListThreadsResponse](#services-messenger-ListThreadsResponse) | @perm |
| GetThread | [GetThreadRequest](#services-messenger-GetThreadRequest) | [GetThreadResponse](#services-messenger-GetThreadResponse) | @perm: Name=ListThreads |
| CreateOrUpdateThread | [CreateOrUpdateThreadRequest](#services-messenger-CreateOrUpdateThreadRequest) | [CreateOrUpdateThreadResponse](#services-messenger-CreateOrUpdateThreadResponse) | @perm |
| DeleteThread | [DeleteThreadRequest](#services-messenger-DeleteThreadRequest) | [DeleteThreadResponse](#services-messenger-DeleteThreadResponse) | @perm |
| SetThreadUserState | [SetThreadUserStateRequest](#services-messenger-SetThreadUserStateRequest) | [SetThreadUserStateResponse](#services-messenger-SetThreadUserStateResponse) | @perm: Name=ListThreads |
| LeaveThread | [LeaveThreadRequest](#services-messenger-LeaveThreadRequest) | [LeaveThreadResponse](#services-messenger-LeaveThreadResponse) | @perm: Name=ListThreads |
| GetThreadMessages | [GetThreadMessagesRequest](#services-messenger-GetThreadMessagesRequest) | [GetThreadMessagesResponse](#services-messenger-GetThreadMessagesResponse) | @perm: Name=ListThreads |
| PostMessage | [PostMessageRequest](#services-messenger-PostMessageRequest) | [PostMessageResponse](#services-messenger-PostMessageResponse) | @perm |
| DeleteMessage | [DeleteMessageRequest](#services-messenger-DeleteMessageRequest) | [DeleteMessageResponse](#services-messenger-DeleteMessageResponse) | @perm: Name=SuperUser |

 



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

