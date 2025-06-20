# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [resources/accounts/accounts.proto](#resources_accounts_accounts-proto)
    - [Account](#resources-accounts-Account)
    - [Character](#resources-accounts-Character)
  
- [resources/accounts/oauth2.proto](#resources_accounts_oauth2-proto)
    - [OAuth2Account](#resources-accounts-OAuth2Account)
    - [OAuth2Provider](#resources-accounts-OAuth2Provider)
  
- [resources/centrum/attributes.proto](#resources_centrum_attributes-proto)
    - [DispatchAttributes](#resources-centrum-DispatchAttributes)
    - [UnitAttributes](#resources-centrum-UnitAttributes)
  
    - [DispatchAttribute](#resources-centrum-DispatchAttribute)
    - [UnitAttribute](#resources-centrum-UnitAttribute)
  
- [resources/centrum/units.proto](#resources_centrum_units-proto)
    - [Unit](#resources-centrum-Unit)
    - [UnitAssignment](#resources-centrum-UnitAssignment)
    - [UnitAssignments](#resources-centrum-UnitAssignments)
    - [UnitStatus](#resources-centrum-UnitStatus)
  
    - [StatusUnit](#resources-centrum-StatusUnit)
  
- [resources/centrum/units_access.proto](#resources_centrum_units_access-proto)
    - [UnitAccess](#resources-centrum-UnitAccess)
    - [UnitJobAccess](#resources-centrum-UnitJobAccess)
    - [UnitQualificationAccess](#resources-centrum-UnitQualificationAccess)
    - [UnitUserAccess](#resources-centrum-UnitUserAccess)
  
    - [UnitAccessLevel](#resources-centrum-UnitAccessLevel)
  
- [resources/centrum/access.proto](#resources_centrum_access-proto)
    - [CentrumAccess](#resources-centrum-CentrumAccess)
    - [CentrumJobAccess](#resources-centrum-CentrumJobAccess)
    - [CentrumQualificationAccess](#resources-centrum-CentrumQualificationAccess)
    - [CentrumUserAccess](#resources-centrum-CentrumUserAccess)
  
    - [CentrumAccessLevel](#resources-centrum-CentrumAccessLevel)
  
- [resources/centrum/dispatchers.proto](#resources_centrum_dispatchers-proto)
    - [Dispatchers](#resources-centrum-Dispatchers)
  
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
  
- [resources/centrum/settings.proto](#resources_centrum_settings-proto)
    - [PredefinedStatus](#resources-centrum-PredefinedStatus)
    - [Settings](#resources-centrum-Settings)
    - [Timings](#resources-centrum-Timings)
  
    - [CentrumMode](#resources-centrum-CentrumMode)
    - [CentrumType](#resources-centrum-CentrumType)
  
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
  
- [resources/common/content/content.proto](#resources_common_content_content-proto)
    - [Content](#resources-common-content-Content)
    - [JSONNode](#resources-common-content-JSONNode)
    - [JSONNode.AttrsEntry](#resources-common-content-JSONNode-AttrsEntry)
  
    - [ContentType](#resources-common-content-ContentType)
    - [NodeType](#resources-common-content-NodeType)
  
- [resources/common/tests/objects.proto](#resources_common_tests_objects-proto)
    - [SimpleObject](#resources-common-tests-SimpleObject)
  
- [resources/common/error.proto](#resources_common_error-proto)
    - [Error](#resources-common-Error)
  
- [resources/common/i18n.proto](#resources_common_i18n-proto)
    - [I18NItem](#resources-common-I18NItem)
    - [I18NItem.ParametersEntry](#resources-common-I18NItem-ParametersEntry)
  
- [resources/common/uuid.proto](#resources_common_uuid-proto)
    - [UUID](#resources-common-UUID)
  
- [resources/documents/category.proto](#resources_documents_category-proto)
    - [Category](#resources-documents-Category)
  
- [resources/documents/comment.proto](#resources_documents_comment-proto)
    - [Comment](#resources-documents-Comment)
  
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
  
- [resources/documents/access.proto](#resources_documents_access-proto)
    - [DocumentAccess](#resources-documents-DocumentAccess)
    - [DocumentJobAccess](#resources-documents-DocumentJobAccess)
    - [DocumentUserAccess](#resources-documents-DocumentUserAccess)
  
    - [AccessLevel](#resources-documents-AccessLevel)
  
- [resources/documents/pins.proto](#resources_documents_pins-proto)
    - [DocumentPin](#resources-documents-DocumentPin)
  
- [resources/documents/activity.proto](#resources_documents_activity-proto)
    - [DocAccessJobsDiff](#resources-documents-DocAccessJobsDiff)
    - [DocAccessRequested](#resources-documents-DocAccessRequested)
    - [DocAccessUpdated](#resources-documents-DocAccessUpdated)
    - [DocAccessUsersDiff](#resources-documents-DocAccessUsersDiff)
    - [DocActivity](#resources-documents-DocActivity)
    - [DocActivityData](#resources-documents-DocActivityData)
    - [DocFilesChange](#resources-documents-DocFilesChange)
    - [DocOwnerChanged](#resources-documents-DocOwnerChanged)
    - [DocUpdated](#resources-documents-DocUpdated)
  
    - [DocActivityType](#resources-documents-DocActivityType)
  
- [resources/documents/documents.proto](#resources_documents_documents-proto)
    - [Document](#resources-documents-Document)
    - [DocumentReference](#resources-documents-DocumentReference)
    - [DocumentRelation](#resources-documents-DocumentRelation)
    - [DocumentShort](#resources-documents-DocumentShort)
    - [WorkflowState](#resources-documents-WorkflowState)
    - [WorkflowUserState](#resources-documents-WorkflowUserState)
  
    - [DocReference](#resources-documents-DocReference)
    - [DocRelation](#resources-documents-DocRelation)
  
- [resources/jobs/activity.proto](#resources_jobs_activity-proto)
    - [AbsenceDateChange](#resources-jobs-AbsenceDateChange)
    - [ColleagueActivity](#resources-jobs-ColleagueActivity)
    - [ColleagueActivityData](#resources-jobs-ColleagueActivityData)
    - [GradeChange](#resources-jobs-GradeChange)
    - [LabelsChange](#resources-jobs-LabelsChange)
    - [NameChange](#resources-jobs-NameChange)
  
    - [ColleagueActivityType](#resources-jobs-ColleagueActivityType)
  
- [resources/jobs/conduct.proto](#resources_jobs_conduct-proto)
    - [ConductEntry](#resources-jobs-ConductEntry)
  
    - [ConductType](#resources-jobs-ConductType)
  
- [resources/jobs/job_settings.proto](#resources_jobs_job_settings-proto)
    - [DiscordSyncChange](#resources-jobs-DiscordSyncChange)
    - [DiscordSyncChanges](#resources-jobs-DiscordSyncChanges)
    - [DiscordSyncSettings](#resources-jobs-DiscordSyncSettings)
    - [GroupMapping](#resources-jobs-GroupMapping)
    - [GroupSyncSettings](#resources-jobs-GroupSyncSettings)
    - [JobSettings](#resources-jobs-JobSettings)
    - [JobsAbsenceSettings](#resources-jobs-JobsAbsenceSettings)
    - [StatusLogSettings](#resources-jobs-StatusLogSettings)
    - [UserInfoSyncSettings](#resources-jobs-UserInfoSyncSettings)
  
    - [UserInfoSyncUnemployedMode](#resources-jobs-UserInfoSyncUnemployedMode)
  
- [resources/jobs/jobs.proto](#resources_jobs_jobs-proto)
    - [Job](#resources-jobs-Job)
    - [JobGrade](#resources-jobs-JobGrade)
  
- [resources/jobs/labels.proto](#resources_jobs_labels-proto)
    - [Label](#resources-jobs-Label)
    - [LabelCount](#resources-jobs-LabelCount)
    - [Labels](#resources-jobs-Labels)
  
- [resources/jobs/timeclock.proto](#resources_jobs_timeclock-proto)
    - [TimeclockEntry](#resources-jobs-TimeclockEntry)
    - [TimeclockStats](#resources-jobs-TimeclockStats)
    - [TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats)
  
    - [TimeclockMode](#resources-jobs-TimeclockMode)
    - [TimeclockViewMode](#resources-jobs-TimeclockViewMode)
  
- [resources/jobs/colleagues.proto](#resources_jobs_colleagues-proto)
    - [Colleague](#resources-jobs-Colleague)
    - [ColleagueProps](#resources-jobs-ColleagueProps)
  
- [resources/jobs/job_props.proto](#resources_jobs_job_props-proto)
    - [JobProps](#resources-jobs-JobProps)
    - [QuickButtons](#resources-jobs-QuickButtons)
  
- [resources/laws/laws.proto](#resources_laws_laws-proto)
    - [Law](#resources-laws-Law)
    - [LawBook](#resources-laws-LawBook)
  
- [resources/notifications/notifications.proto](#resources_notifications_notifications-proto)
    - [CalendarData](#resources-notifications-CalendarData)
    - [Data](#resources-notifications-Data)
    - [Link](#resources-notifications-Link)
    - [Notification](#resources-notifications-Notification)
  
    - [NotificationCategory](#resources-notifications-NotificationCategory)
    - [NotificationType](#resources-notifications-NotificationType)
  
- [resources/notifications/events.proto](#resources_notifications_events-proto)
    - [JobEvent](#resources-notifications-JobEvent)
    - [JobGradeEvent](#resources-notifications-JobGradeEvent)
    - [SystemEvent](#resources-notifications-SystemEvent)
    - [UserEvent](#resources-notifications-UserEvent)
  
- [resources/notifications/client_view.proto](#resources_notifications_client_view-proto)
    - [ClientView](#resources-notifications-ClientView)
    - [ObjectEvent](#resources-notifications-ObjectEvent)
  
    - [ObjectEventType](#resources-notifications-ObjectEventType)
    - [ObjectType](#resources-notifications-ObjectType)
  
- [resources/permissions/attributes.proto](#resources_permissions_attributes-proto)
    - [AttributeValues](#resources-permissions-AttributeValues)
    - [JobGradeList](#resources-permissions-JobGradeList)
    - [JobGradeList.GradesEntry](#resources-permissions-JobGradeList-GradesEntry)
    - [JobGradeList.JobsEntry](#resources-permissions-JobGradeList-JobsEntry)
    - [JobGrades](#resources-permissions-JobGrades)
    - [RoleAttribute](#resources-permissions-RoleAttribute)
    - [StringList](#resources-permissions-StringList)
  
- [resources/permissions/events.proto](#resources_permissions_events-proto)
    - [JobLimitsUpdatedEvent](#resources-permissions-JobLimitsUpdatedEvent)
    - [RoleIDEvent](#resources-permissions-RoleIDEvent)
  
- [resources/permissions/permissions.proto](#resources_permissions_permissions-proto)
    - [PermItem](#resources-permissions-PermItem)
    - [Permission](#resources-permissions-Permission)
    - [Role](#resources-permissions-Role)
  
- [resources/qualifications/access.proto](#resources_qualifications_access-proto)
    - [QualificationAccess](#resources-qualifications-QualificationAccess)
    - [QualificationJobAccess](#resources-qualifications-QualificationJobAccess)
    - [QualificationUserAccess](#resources-qualifications-QualificationUserAccess)
  
    - [AccessLevel](#resources-qualifications-AccessLevel)
  
- [resources/qualifications/exam.proto](#resources_qualifications_exam-proto)
    - [ExamGrading](#resources-qualifications-ExamGrading)
    - [ExamGradingResponse](#resources-qualifications-ExamGradingResponse)
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
  
    - [AutoGradeMode](#resources-qualifications-AutoGradeMode)
    - [QualificationExamMode](#resources-qualifications-QualificationExamMode)
    - [RequestStatus](#resources-qualifications-RequestStatus)
    - [ResultStatus](#resources-qualifications-ResultStatus)
  
- [resources/timestamp/timestamp.proto](#resources_timestamp_timestamp-proto)
    - [Timestamp](#resources-timestamp-Timestamp)
  
- [resources/vehicles/vehicles.proto](#resources_vehicles_vehicles-proto)
    - [Vehicle](#resources-vehicles-Vehicle)
  
- [resources/calendar/calendar.proto](#resources_calendar_calendar-proto)
    - [Calendar](#resources-calendar-Calendar)
    - [CalendarEntry](#resources-calendar-CalendarEntry)
    - [CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP)
    - [CalendarEntryRecurring](#resources-calendar-CalendarEntryRecurring)
    - [CalendarShort](#resources-calendar-CalendarShort)
    - [CalendarSub](#resources-calendar-CalendarSub)
  
    - [RsvpResponses](#resources-calendar-RsvpResponses)
  
- [resources/calendar/access.proto](#resources_calendar_access-proto)
    - [CalendarAccess](#resources-calendar-CalendarAccess)
    - [CalendarJobAccess](#resources-calendar-CalendarJobAccess)
    - [CalendarUserAccess](#resources-calendar-CalendarUserAccess)
  
    - [AccessLevel](#resources-calendar-AccessLevel)
  
- [resources/stats/stats.proto](#resources_stats_stats-proto)
    - [Stat](#resources-stats-Stat)
  
- [resources/internet/domain.proto](#resources_internet_domain-proto)
    - [Domain](#resources-internet-Domain)
    - [TLD](#resources-internet-TLD)
  
- [resources/internet/page.proto](#resources_internet_page-proto)
    - [ContentNode](#resources-internet-ContentNode)
    - [ContentNode.AttrsEntry](#resources-internet-ContentNode-AttrsEntry)
    - [Page](#resources-internet-Page)
    - [PageData](#resources-internet-PageData)
  
    - [PageLayoutType](#resources-internet-PageLayoutType)
  
- [resources/internet/search.proto](#resources_internet_search-proto)
    - [SearchResult](#resources-internet-SearchResult)
  
- [resources/internet/access.proto](#resources_internet_access-proto)
    - [PageAccess](#resources-internet-PageAccess)
    - [PageJobAccess](#resources-internet-PageJobAccess)
    - [PageUserAccess](#resources-internet-PageUserAccess)
  
    - [AccessLevel](#resources-internet-AccessLevel)
  
- [resources/internet/ads.proto](#resources_internet_ads-proto)
    - [Ad](#resources-internet-Ad)
  
    - [AdType](#resources-internet-AdType)
  
- [resources/mailer/events.proto](#resources_mailer_events-proto)
    - [MailerEvent](#resources-mailer-MailerEvent)
  
- [resources/mailer/message.proto](#resources_mailer_message-proto)
    - [Message](#resources-mailer-Message)
    - [MessageAttachment](#resources-mailer-MessageAttachment)
    - [MessageAttachmentDocument](#resources-mailer-MessageAttachmentDocument)
    - [MessageData](#resources-mailer-MessageData)
  
- [resources/mailer/settings.proto](#resources_mailer_settings-proto)
    - [EmailSettings](#resources-mailer-EmailSettings)
  
- [resources/mailer/template.proto](#resources_mailer_template-proto)
    - [Template](#resources-mailer-Template)
  
- [resources/mailer/thread.proto](#resources_mailer_thread-proto)
    - [Thread](#resources-mailer-Thread)
    - [ThreadRecipientEmail](#resources-mailer-ThreadRecipientEmail)
    - [ThreadState](#resources-mailer-ThreadState)
  
- [resources/mailer/email.proto](#resources_mailer_email-proto)
    - [Email](#resources-mailer-Email)
  
- [resources/mailer/access.proto](#resources_mailer_access-proto)
    - [Access](#resources-mailer-Access)
    - [JobAccess](#resources-mailer-JobAccess)
    - [QualificationAccess](#resources-mailer-QualificationAccess)
    - [UserAccess](#resources-mailer-UserAccess)
  
    - [AccessLevel](#resources-mailer-AccessLevel)
  
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
    - [PageFilesChange](#resources-wiki-PageFilesChange)
    - [PageUpdated](#resources-wiki-PageUpdated)
  
    - [PageActivityType](#resources-wiki-PageActivityType)
  
- [resources/wiki/page.proto](#resources_wiki_page-proto)
    - [Page](#resources-wiki-Page)
    - [PageMeta](#resources-wiki-PageMeta)
    - [PageRootInfo](#resources-wiki-PageRootInfo)
    - [PageShort](#resources-wiki-PageShort)
  
- [resources/sync/activity.proto](#resources_sync_activity-proto)
    - [ColleagueProps](#resources-sync-ColleagueProps)
    - [TimeclockUpdate](#resources-sync-TimeclockUpdate)
    - [UserOAuth2Conn](#resources-sync-UserOAuth2Conn)
    - [UserProps](#resources-sync-UserProps)
    - [UserUpdate](#resources-sync-UserUpdate)
  
- [resources/sync/data.proto](#resources_sync_data-proto)
    - [CitizenLocations](#resources-sync-CitizenLocations)
    - [DataJobs](#resources-sync-DataJobs)
    - [DataLicenses](#resources-sync-DataLicenses)
    - [DataStatus](#resources-sync-DataStatus)
    - [DataUserLocations](#resources-sync-DataUserLocations)
    - [DataUsers](#resources-sync-DataUsers)
    - [DataVehicles](#resources-sync-DataVehicles)
    - [DeleteUsers](#resources-sync-DeleteUsers)
    - [DeleteVehicles](#resources-sync-DeleteVehicles)
  
- [resources/users/activity.proto](#resources_users_activity-proto)
    - [CitizenDocumentRelation](#resources-users-CitizenDocumentRelation)
    - [FineChange](#resources-users-FineChange)
    - [JailChange](#resources-users-JailChange)
    - [JobChange](#resources-users-JobChange)
    - [LabelsChange](#resources-users-LabelsChange)
    - [LicenseChange](#resources-users-LicenseChange)
    - [MugshotChange](#resources-users-MugshotChange)
    - [NameChange](#resources-users-NameChange)
    - [TrafficInfractionPointsChange](#resources-users-TrafficInfractionPointsChange)
    - [UserActivity](#resources-users-UserActivity)
    - [UserActivityData](#resources-users-UserActivityData)
    - [WantedChange](#resources-users-WantedChange)
  
    - [UserActivityType](#resources-users-UserActivityType)
  
- [resources/users/licenses.proto](#resources_users_licenses-proto)
    - [CitizensLicenses](#resources-users-CitizensLicenses)
    - [License](#resources-users-License)
  
- [resources/users/labels.proto](#resources_users_labels-proto)
    - [Label](#resources-users-Label)
    - [Labels](#resources-users-Labels)
  
- [resources/users/props.proto](#resources_users_props-proto)
    - [UserProps](#resources-users-UserProps)
  
- [resources/users/users.proto](#resources_users_users-proto)
    - [User](#resources-users-User)
    - [UserShort](#resources-users-UserShort)
  
- [resources/audit/audit.proto](#resources_audit_audit-proto)
    - [AuditEntry](#resources-audit-AuditEntry)
  
    - [EventType](#resources-audit-EventType)
  
- [resources/settings/banner.proto](#resources_settings_banner-proto)
    - [BannerMessage](#resources-settings-BannerMessage)
  
- [resources/settings/config.proto](#resources_settings_config-proto)
    - [AppConfig](#resources-settings-AppConfig)
    - [Auth](#resources-settings-Auth)
    - [Discord](#resources-settings-Discord)
    - [DiscordBotPresence](#resources-settings-DiscordBotPresence)
    - [JobInfo](#resources-settings-JobInfo)
    - [Links](#resources-settings-Links)
    - [Perm](#resources-settings-Perm)
    - [Perms](#resources-settings-Perms)
    - [System](#resources-settings-System)
    - [UnemployedJob](#resources-settings-UnemployedJob)
    - [UserTracker](#resources-settings-UserTracker)
    - [Website](#resources-settings-Website)
  
    - [DiscordBotPresenceType](#resources-settings-DiscordBotPresenceType)
  
- [resources/collab/collab.proto](#resources_collab_collab-proto)
    - [AwarenessPing](#resources-collab-AwarenessPing)
    - [ClientPacket](#resources-collab-ClientPacket)
    - [CollabHandshake](#resources-collab-CollabHandshake)
    - [CollabInit](#resources-collab-CollabInit)
    - [ServerPacket](#resources-collab-ServerPacket)
    - [SyncStep](#resources-collab-SyncStep)
    - [TargetSaved](#resources-collab-TargetSaved)
    - [YjsUpdate](#resources-collab-YjsUpdate)
  
    - [ClientRole](#resources-collab-ClientRole)
  
- [resources/discord/discord.proto](#resources_discord_discord-proto)
    - [Channel](#resources-discord-Channel)
    - [Guild](#resources-discord-Guild)
  
- [resources/file/file.proto](#resources_file_file-proto)
    - [File](#resources-file-File)
  
- [resources/file/filestore.proto](#resources_file_filestore-proto)
    - [DeleteFileRequest](#resources-file-DeleteFileRequest)
    - [DeleteFileResponse](#resources-file-DeleteFileResponse)
    - [UploadMeta](#resources-file-UploadMeta)
    - [UploadPacket](#resources-file-UploadPacket)
    - [UploadResponse](#resources-file-UploadResponse)
  
- [resources/file/meta.proto](#resources_file_meta-proto)
    - [FileMeta](#resources-file-FileMeta)
    - [ImageMeta](#resources-file-ImageMeta)
  
- [resources/livemap/coords.proto](#resources_livemap_coords-proto)
    - [Coords](#resources-livemap-Coords)
  
- [resources/livemap/heatmap.proto](#resources_livemap_heatmap-proto)
    - [HeatmapEntry](#resources-livemap-HeatmapEntry)
  
- [resources/livemap/marker_marker.proto](#resources_livemap_marker_marker-proto)
    - [CircleMarker](#resources-livemap-CircleMarker)
    - [IconMarker](#resources-livemap-IconMarker)
    - [MarkerData](#resources-livemap-MarkerData)
    - [MarkerMarker](#resources-livemap-MarkerMarker)
  
    - [MarkerType](#resources-livemap-MarkerType)
  
- [resources/livemap/user_marker.proto](#resources_livemap_user_marker-proto)
    - [UserMarker](#resources-livemap-UserMarker)
  
- [resources/tracker/mapping.proto](#resources_tracker_mapping-proto)
    - [UserMapping](#resources-tracker-UserMapping)
  
- [resources/clientconfig/clientconfig.proto](#resources_clientconfig_clientconfig-proto)
    - [ClientConfig](#resources-clientconfig-ClientConfig)
    - [Discord](#resources-clientconfig-Discord)
    - [FeatureGates](#resources-clientconfig-FeatureGates)
    - [Game](#resources-clientconfig-Game)
    - [Links](#resources-clientconfig-Links)
    - [LoginConfig](#resources-clientconfig-LoginConfig)
    - [OTLPFrontend](#resources-clientconfig-OTLPFrontend)
    - [OTLPFrontend.HeadersEntry](#resources-clientconfig-OTLPFrontend-HeadersEntry)
    - [ProviderConfig](#resources-clientconfig-ProviderConfig)
    - [System](#resources-clientconfig-System)
    - [Website](#resources-clientconfig-Website)
  
- [resources/userinfo/user_info.proto](#resources_userinfo_user_info-proto)
    - [PollReq](#resources-userinfo-PollReq)
    - [UserInfo](#resources-userinfo-UserInfo)
    - [UserInfoChanged](#resources-userinfo-UserInfoChanged)
  
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
    - [SetSuperuserModeRequest](#services-auth-SetSuperuserModeRequest)
    - [SetSuperuserModeResponse](#services-auth-SetSuperuserModeResponse)
  
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
    - [Dispatchers](#services-centrum-Dispatchers)
    - [GetDispatchHeatmapRequest](#services-centrum-GetDispatchHeatmapRequest)
    - [GetDispatchHeatmapResponse](#services-centrum-GetDispatchHeatmapResponse)
    - [GetDispatchRequest](#services-centrum-GetDispatchRequest)
    - [GetDispatchResponse](#services-centrum-GetDispatchResponse)
    - [GetSettingsRequest](#services-centrum-GetSettingsRequest)
    - [GetSettingsResponse](#services-centrum-GetSettingsResponse)
    - [JobAccess](#services-centrum-JobAccess)
    - [JobAccessEntry](#services-centrum-JobAccessEntry)
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
    - [StreamHandshake](#services-centrum-StreamHandshake)
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
    - [UpdateDispatchersRequest](#services-centrum-UpdateDispatchersRequest)
    - [UpdateDispatchersResponse](#services-centrum-UpdateDispatchersResponse)
    - [UpdateSettingsRequest](#services-centrum-UpdateSettingsRequest)
    - [UpdateSettingsResponse](#services-centrum-UpdateSettingsResponse)
    - [UpdateUnitStatusRequest](#services-centrum-UpdateUnitStatusRequest)
    - [UpdateUnitStatusResponse](#services-centrum-UpdateUnitStatusResponse)
  
    - [CentrumService](#services-centrum-CentrumService)
  
- [services/completor/completor.proto](#services_completor_completor-proto)
    - [CompleteCitizenLabelsRequest](#services-completor-CompleteCitizenLabelsRequest)
    - [CompleteCitizenLabelsResponse](#services-completor-CompleteCitizenLabelsResponse)
    - [CompleteCitizensRequest](#services-completor-CompleteCitizensRequest)
    - [CompleteCitizensRespoonse](#services-completor-CompleteCitizensRespoonse)
    - [CompleteDocumentCategoriesRequest](#services-completor-CompleteDocumentCategoriesRequest)
    - [CompleteDocumentCategoriesResponse](#services-completor-CompleteDocumentCategoriesResponse)
    - [CompleteJobsRequest](#services-completor-CompleteJobsRequest)
    - [CompleteJobsResponse](#services-completor-CompleteJobsResponse)
    - [ListLawBooksRequest](#services-completor-ListLawBooksRequest)
    - [ListLawBooksResponse](#services-completor-ListLawBooksResponse)
  
    - [CompletorService](#services-completor-CompletorService)
  
- [services/jobs/conduct.proto](#services_jobs_conduct-proto)
    - [CreateConductEntryRequest](#services-jobs-CreateConductEntryRequest)
    - [CreateConductEntryResponse](#services-jobs-CreateConductEntryResponse)
    - [DeleteConductEntryRequest](#services-jobs-DeleteConductEntryRequest)
    - [DeleteConductEntryResponse](#services-jobs-DeleteConductEntryResponse)
    - [ListConductEntriesRequest](#services-jobs-ListConductEntriesRequest)
    - [ListConductEntriesResponse](#services-jobs-ListConductEntriesResponse)
    - [UpdateConductEntryRequest](#services-jobs-UpdateConductEntryRequest)
    - [UpdateConductEntryResponse](#services-jobs-UpdateConductEntryResponse)
  
    - [ConductService](#services-jobs-ConductService)
  
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
    - [ManageLabelsRequest](#services-jobs-ManageLabelsRequest)
    - [ManageLabelsResponse](#services-jobs-ManageLabelsResponse)
    - [SetColleaguePropsRequest](#services-jobs-SetColleaguePropsRequest)
    - [SetColleaguePropsResponse](#services-jobs-SetColleaguePropsResponse)
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
  
    - [TimeclockService](#services-jobs-TimeclockService)
  
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
  
- [services/calendar/calendar.proto](#services_calendar_calendar-proto)
    - [CreateCalendarRequest](#services-calendar-CreateCalendarRequest)
    - [CreateCalendarResponse](#services-calendar-CreateCalendarResponse)
    - [CreateOrUpdateCalendarEntryRequest](#services-calendar-CreateOrUpdateCalendarEntryRequest)
    - [CreateOrUpdateCalendarEntryResponse](#services-calendar-CreateOrUpdateCalendarEntryResponse)
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
    - [UpdateCalendarRequest](#services-calendar-UpdateCalendarRequest)
    - [UpdateCalendarResponse](#services-calendar-UpdateCalendarResponse)
  
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
  
- [services/internet/domain.proto](#services_internet_domain-proto)
    - [CheckDomainAvailabilityRequest](#services-internet-CheckDomainAvailabilityRequest)
    - [CheckDomainAvailabilityResponse](#services-internet-CheckDomainAvailabilityResponse)
    - [ListDomainsRequest](#services-internet-ListDomainsRequest)
    - [ListDomainsResponse](#services-internet-ListDomainsResponse)
    - [ListTLDsRequest](#services-internet-ListTLDsRequest)
    - [ListTLDsResponse](#services-internet-ListTLDsResponse)
    - [RegisterDomainRequest](#services-internet-RegisterDomainRequest)
    - [RegisterDomainResponse](#services-internet-RegisterDomainResponse)
    - [UpdateDomainRequest](#services-internet-UpdateDomainRequest)
    - [UpdateDomainResponse](#services-internet-UpdateDomainResponse)
  
    - [DomainService](#services-internet-DomainService)
  
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
  
- [services/wiki/collab.proto](#services_wiki_collab-proto)
    - [CollabService](#services-wiki-CollabService)
  
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
  
- [services/sync/sync.proto](#services_sync_sync-proto)
    - [AddActivityRequest](#services-sync-AddActivityRequest)
    - [AddActivityResponse](#services-sync-AddActivityResponse)
    - [DeleteDataRequest](#services-sync-DeleteDataRequest)
    - [DeleteDataResponse](#services-sync-DeleteDataResponse)
    - [GetStatusRequest](#services-sync-GetStatusRequest)
    - [GetStatusResponse](#services-sync-GetStatusResponse)
    - [RegisterAccountRequest](#services-sync-RegisterAccountRequest)
    - [RegisterAccountResponse](#services-sync-RegisterAccountResponse)
    - [SendDataRequest](#services-sync-SendDataRequest)
    - [SendDataResponse](#services-sync-SendDataResponse)
    - [StreamRequest](#services-sync-StreamRequest)
    - [StreamResponse](#services-sync-StreamResponse)
    - [TransferAccountRequest](#services-sync-TransferAccountRequest)
    - [TransferAccountResponse](#services-sync-TransferAccountResponse)
  
    - [SyncService](#services-sync-SyncService)
  
- [services/citizens/citizens.proto](#services_citizens_citizens-proto)
    - [DeleteAvatarRequest](#services-citizens-DeleteAvatarRequest)
    - [DeleteAvatarResponse](#services-citizens-DeleteAvatarResponse)
    - [DeleteMugshotRequest](#services-citizens-DeleteMugshotRequest)
    - [DeleteMugshotResponse](#services-citizens-DeleteMugshotResponse)
    - [GetUserRequest](#services-citizens-GetUserRequest)
    - [GetUserResponse](#services-citizens-GetUserResponse)
    - [ListCitizensRequest](#services-citizens-ListCitizensRequest)
    - [ListCitizensResponse](#services-citizens-ListCitizensResponse)
    - [ListUserActivityRequest](#services-citizens-ListUserActivityRequest)
    - [ListUserActivityResponse](#services-citizens-ListUserActivityResponse)
    - [ManageLabelsRequest](#services-citizens-ManageLabelsRequest)
    - [ManageLabelsResponse](#services-citizens-ManageLabelsResponse)
    - [SetUserPropsRequest](#services-citizens-SetUserPropsRequest)
    - [SetUserPropsResponse](#services-citizens-SetUserPropsResponse)
  
    - [CitizensService](#services-citizens-CitizensService)
  
- [services/documents/collab.proto](#services_documents_collab-proto)
    - [CollabService](#services-documents-CollabService)
  
- [services/documents/documents.proto](#services_documents_documents-proto)
    - [AddDocumentReferenceRequest](#services-documents-AddDocumentReferenceRequest)
    - [AddDocumentReferenceResponse](#services-documents-AddDocumentReferenceResponse)
    - [AddDocumentRelationRequest](#services-documents-AddDocumentRelationRequest)
    - [AddDocumentRelationResponse](#services-documents-AddDocumentRelationResponse)
    - [ChangeDocumentOwnerRequest](#services-documents-ChangeDocumentOwnerRequest)
    - [ChangeDocumentOwnerResponse](#services-documents-ChangeDocumentOwnerResponse)
    - [CreateDocumentReqRequest](#services-documents-CreateDocumentReqRequest)
    - [CreateDocumentReqResponse](#services-documents-CreateDocumentReqResponse)
    - [CreateDocumentRequest](#services-documents-CreateDocumentRequest)
    - [CreateDocumentResponse](#services-documents-CreateDocumentResponse)
    - [CreateOrUpdateCategoryRequest](#services-documents-CreateOrUpdateCategoryRequest)
    - [CreateOrUpdateCategoryResponse](#services-documents-CreateOrUpdateCategoryResponse)
    - [CreateTemplateRequest](#services-documents-CreateTemplateRequest)
    - [CreateTemplateResponse](#services-documents-CreateTemplateResponse)
    - [DeleteCategoryRequest](#services-documents-DeleteCategoryRequest)
    - [DeleteCategoryResponse](#services-documents-DeleteCategoryResponse)
    - [DeleteCommentRequest](#services-documents-DeleteCommentRequest)
    - [DeleteCommentResponse](#services-documents-DeleteCommentResponse)
    - [DeleteDocumentReqRequest](#services-documents-DeleteDocumentReqRequest)
    - [DeleteDocumentReqResponse](#services-documents-DeleteDocumentReqResponse)
    - [DeleteDocumentRequest](#services-documents-DeleteDocumentRequest)
    - [DeleteDocumentResponse](#services-documents-DeleteDocumentResponse)
    - [DeleteTemplateRequest](#services-documents-DeleteTemplateRequest)
    - [DeleteTemplateResponse](#services-documents-DeleteTemplateResponse)
    - [EditCommentRequest](#services-documents-EditCommentRequest)
    - [EditCommentResponse](#services-documents-EditCommentResponse)
    - [GetCommentsRequest](#services-documents-GetCommentsRequest)
    - [GetCommentsResponse](#services-documents-GetCommentsResponse)
    - [GetDocumentAccessRequest](#services-documents-GetDocumentAccessRequest)
    - [GetDocumentAccessResponse](#services-documents-GetDocumentAccessResponse)
    - [GetDocumentReferencesRequest](#services-documents-GetDocumentReferencesRequest)
    - [GetDocumentReferencesResponse](#services-documents-GetDocumentReferencesResponse)
    - [GetDocumentRelationsRequest](#services-documents-GetDocumentRelationsRequest)
    - [GetDocumentRelationsResponse](#services-documents-GetDocumentRelationsResponse)
    - [GetDocumentRequest](#services-documents-GetDocumentRequest)
    - [GetDocumentResponse](#services-documents-GetDocumentResponse)
    - [GetTemplateRequest](#services-documents-GetTemplateRequest)
    - [GetTemplateResponse](#services-documents-GetTemplateResponse)
    - [ListCategoriesRequest](#services-documents-ListCategoriesRequest)
    - [ListCategoriesResponse](#services-documents-ListCategoriesResponse)
    - [ListDocumentActivityRequest](#services-documents-ListDocumentActivityRequest)
    - [ListDocumentActivityResponse](#services-documents-ListDocumentActivityResponse)
    - [ListDocumentPinsRequest](#services-documents-ListDocumentPinsRequest)
    - [ListDocumentPinsResponse](#services-documents-ListDocumentPinsResponse)
    - [ListDocumentReqsRequest](#services-documents-ListDocumentReqsRequest)
    - [ListDocumentReqsResponse](#services-documents-ListDocumentReqsResponse)
    - [ListDocumentsRequest](#services-documents-ListDocumentsRequest)
    - [ListDocumentsResponse](#services-documents-ListDocumentsResponse)
    - [ListTemplatesRequest](#services-documents-ListTemplatesRequest)
    - [ListTemplatesResponse](#services-documents-ListTemplatesResponse)
    - [ListUserDocumentsRequest](#services-documents-ListUserDocumentsRequest)
    - [ListUserDocumentsResponse](#services-documents-ListUserDocumentsResponse)
    - [PostCommentRequest](#services-documents-PostCommentRequest)
    - [PostCommentResponse](#services-documents-PostCommentResponse)
    - [RemoveDocumentReferenceRequest](#services-documents-RemoveDocumentReferenceRequest)
    - [RemoveDocumentReferenceResponse](#services-documents-RemoveDocumentReferenceResponse)
    - [RemoveDocumentRelationRequest](#services-documents-RemoveDocumentRelationRequest)
    - [RemoveDocumentRelationResponse](#services-documents-RemoveDocumentRelationResponse)
    - [SetDocumentAccessRequest](#services-documents-SetDocumentAccessRequest)
    - [SetDocumentAccessResponse](#services-documents-SetDocumentAccessResponse)
    - [SetDocumentReminderRequest](#services-documents-SetDocumentReminderRequest)
    - [SetDocumentReminderResponse](#services-documents-SetDocumentReminderResponse)
    - [ToggleDocumentPinRequest](#services-documents-ToggleDocumentPinRequest)
    - [ToggleDocumentPinResponse](#services-documents-ToggleDocumentPinResponse)
    - [ToggleDocumentRequest](#services-documents-ToggleDocumentRequest)
    - [ToggleDocumentResponse](#services-documents-ToggleDocumentResponse)
    - [UpdateDocumentReqRequest](#services-documents-UpdateDocumentReqRequest)
    - [UpdateDocumentReqResponse](#services-documents-UpdateDocumentReqResponse)
    - [UpdateDocumentRequest](#services-documents-UpdateDocumentRequest)
    - [UpdateDocumentResponse](#services-documents-UpdateDocumentResponse)
    - [UpdateTemplateRequest](#services-documents-UpdateTemplateRequest)
    - [UpdateTemplateResponse](#services-documents-UpdateTemplateResponse)
  
    - [DocumentsService](#services-documents-DocumentsService)
  
- [services/livemap/livemap.proto](#services_livemap_livemap-proto)
    - [CreateOrUpdateMarkerRequest](#services-livemap-CreateOrUpdateMarkerRequest)
    - [CreateOrUpdateMarkerResponse](#services-livemap-CreateOrUpdateMarkerResponse)
    - [DeleteMarkerRequest](#services-livemap-DeleteMarkerRequest)
    - [DeleteMarkerResponse](#services-livemap-DeleteMarkerResponse)
    - [JobsList](#services-livemap-JobsList)
    - [MarkerMarkersUpdates](#services-livemap-MarkerMarkersUpdates)
    - [Snapshot](#services-livemap-Snapshot)
    - [StreamRequest](#services-livemap-StreamRequest)
    - [StreamResponse](#services-livemap-StreamResponse)
    - [UserDelete](#services-livemap-UserDelete)
  
    - [LivemapService](#services-livemap-LivemapService)
  
- [services/settings/accounts.proto](#services_settings_accounts-proto)
    - [DeleteAccountRequest](#services-settings-DeleteAccountRequest)
    - [DeleteAccountResponse](#services-settings-DeleteAccountResponse)
    - [DisconnectOAuth2ConnectionRequest](#services-settings-DisconnectOAuth2ConnectionRequest)
    - [DisconnectOAuth2ConnectionResponse](#services-settings-DisconnectOAuth2ConnectionResponse)
    - [ListAccountsRequest](#services-settings-ListAccountsRequest)
    - [ListAccountsResponse](#services-settings-ListAccountsResponse)
    - [UpdateAccountRequest](#services-settings-UpdateAccountRequest)
    - [UpdateAccountResponse](#services-settings-UpdateAccountResponse)
  
    - [AccountsService](#services-settings-AccountsService)
  
- [services/settings/config.proto](#services_settings_config-proto)
    - [GetAppConfigRequest](#services-settings-GetAppConfigRequest)
    - [GetAppConfigResponse](#services-settings-GetAppConfigResponse)
    - [UpdateAppConfigRequest](#services-settings-UpdateAppConfigRequest)
    - [UpdateAppConfigResponse](#services-settings-UpdateAppConfigResponse)
  
    - [ConfigService](#services-settings-ConfigService)
  
- [services/settings/cron.proto](#services_settings_cron-proto)
    - [ListCronjobsRequest](#services-settings-ListCronjobsRequest)
    - [ListCronjobsResponse](#services-settings-ListCronjobsResponse)
  
    - [CronService](#services-settings-CronService)
  
- [services/settings/laws.proto](#services_settings_laws-proto)
    - [CreateOrUpdateLawBookRequest](#services-settings-CreateOrUpdateLawBookRequest)
    - [CreateOrUpdateLawBookResponse](#services-settings-CreateOrUpdateLawBookResponse)
    - [CreateOrUpdateLawRequest](#services-settings-CreateOrUpdateLawRequest)
    - [CreateOrUpdateLawResponse](#services-settings-CreateOrUpdateLawResponse)
    - [DeleteLawBookRequest](#services-settings-DeleteLawBookRequest)
    - [DeleteLawBookResponse](#services-settings-DeleteLawBookResponse)
    - [DeleteLawRequest](#services-settings-DeleteLawRequest)
    - [DeleteLawResponse](#services-settings-DeleteLawResponse)
  
    - [LawsService](#services-settings-LawsService)
  
- [services/settings/settings.proto](#services_settings_settings-proto)
    - [AttrsUpdate](#services-settings-AttrsUpdate)
    - [CreateRoleRequest](#services-settings-CreateRoleRequest)
    - [CreateRoleResponse](#services-settings-CreateRoleResponse)
    - [DeleteFactionRequest](#services-settings-DeleteFactionRequest)
    - [DeleteFactionResponse](#services-settings-DeleteFactionResponse)
    - [DeleteJobLogoRequest](#services-settings-DeleteJobLogoRequest)
    - [DeleteJobLogoResponse](#services-settings-DeleteJobLogoResponse)
    - [DeleteRoleRequest](#services-settings-DeleteRoleRequest)
    - [DeleteRoleResponse](#services-settings-DeleteRoleResponse)
    - [GetAllPermissionsRequest](#services-settings-GetAllPermissionsRequest)
    - [GetAllPermissionsResponse](#services-settings-GetAllPermissionsResponse)
    - [GetEffectivePermissionsRequest](#services-settings-GetEffectivePermissionsRequest)
    - [GetEffectivePermissionsResponse](#services-settings-GetEffectivePermissionsResponse)
    - [GetJobLimitsRequest](#services-settings-GetJobLimitsRequest)
    - [GetJobLimitsResponse](#services-settings-GetJobLimitsResponse)
    - [GetJobPropsRequest](#services-settings-GetJobPropsRequest)
    - [GetJobPropsResponse](#services-settings-GetJobPropsResponse)
    - [GetPermissionsRequest](#services-settings-GetPermissionsRequest)
    - [GetPermissionsResponse](#services-settings-GetPermissionsResponse)
    - [GetRoleRequest](#services-settings-GetRoleRequest)
    - [GetRoleResponse](#services-settings-GetRoleResponse)
    - [GetRolesRequest](#services-settings-GetRolesRequest)
    - [GetRolesResponse](#services-settings-GetRolesResponse)
    - [ListDiscordChannelsRequest](#services-settings-ListDiscordChannelsRequest)
    - [ListDiscordChannelsResponse](#services-settings-ListDiscordChannelsResponse)
    - [ListUserGuildsRequest](#services-settings-ListUserGuildsRequest)
    - [ListUserGuildsResponse](#services-settings-ListUserGuildsResponse)
    - [PermsUpdate](#services-settings-PermsUpdate)
    - [SetJobPropsRequest](#services-settings-SetJobPropsRequest)
    - [SetJobPropsResponse](#services-settings-SetJobPropsResponse)
    - [UpdateJobLimitsRequest](#services-settings-UpdateJobLimitsRequest)
    - [UpdateJobLimitsResponse](#services-settings-UpdateJobLimitsResponse)
    - [UpdateRolePermsRequest](#services-settings-UpdateRolePermsRequest)
    - [UpdateRolePermsResponse](#services-settings-UpdateRolePermsResponse)
    - [ViewAuditLogRequest](#services-settings-ViewAuditLogRequest)
    - [ViewAuditLogResponse](#services-settings-ViewAuditLogResponse)
  
    - [SettingsService](#services-settings-SettingsService)
  
- [services/vehicles/vehicles.proto](#services_vehicles_vehicles-proto)
    - [ListVehiclesRequest](#services-vehicles-ListVehiclesRequest)
    - [ListVehiclesResponse](#services-vehicles-ListVehiclesResponse)
  
    - [VehiclesService](#services-vehicles-VehiclesService)
  
- [services/filestore/filestore.proto](#services_filestore_filestore-proto)
    - [DeleteFileByPathRequest](#services-filestore-DeleteFileByPathRequest)
    - [DeleteFileByPathResponse](#services-filestore-DeleteFileByPathResponse)
    - [ListFilesRequest](#services-filestore-ListFilesRequest)
    - [ListFilesResponse](#services-filestore-ListFilesResponse)
  
    - [FilestoreService](#services-filestore-FilestoreService)
  
- [services/notifications/notifications.proto](#services_notifications_notifications-proto)
    - [GetNotificationsRequest](#services-notifications-GetNotificationsRequest)
    - [GetNotificationsResponse](#services-notifications-GetNotificationsResponse)
    - [MarkNotificationsRequest](#services-notifications-MarkNotificationsRequest)
    - [MarkNotificationsResponse](#services-notifications-MarkNotificationsResponse)
    - [StreamMessage](#services-notifications-StreamMessage)
    - [StreamResponse](#services-notifications-StreamResponse)
  
    - [NotificationsService](#services-notifications-NotificationsService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="resources_accounts_accounts-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/accounts/accounts.proto



<a name="resources-accounts-Account"></a>

### Account



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `username` | [string](#string) |  |  |
| `license` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |
| `last_char` | [int32](#int32) | optional |  |
| `oauth2_accounts` | [OAuth2Account](#resources-accounts-OAuth2Account) | repeated | @gotags: alias:"oauth2_account" |






<a name="resources-accounts-Character"></a>

### Character



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `available` | [bool](#bool) |  |  |
| `group` | [string](#string) |  |  |
| `char` | [resources.users.User](#resources-users-User) |  | @gotags: alias:"user" |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_accounts_oauth2-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/accounts/oauth2.proto



<a name="resources-accounts-OAuth2Account"></a>

### OAuth2Account



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `provider_name` | [string](#string) |  | @gotags: sql:"primary_key" alias:"provider_name" |
| `provider` | [OAuth2Provider](#resources-accounts-OAuth2Provider) |  |  |
| `external_id` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `avatar` | [string](#string) |  |  |






<a name="resources-accounts-OAuth2Provider"></a>

### OAuth2Provider



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `label` | [string](#string) |  |  |
| `homepage` | [string](#string) |  |  |
| `icon` | [string](#string) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_attributes-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/attributes.proto



<a name="resources-centrum-DispatchAttributes"></a>

### DispatchAttributes
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [DispatchAttribute](#resources-centrum-DispatchAttribute) | repeated |  |






<a name="resources-centrum-UnitAttributes"></a>

### UnitAttributes
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [UnitAttribute](#resources-centrum-UnitAttribute) | repeated |  |





 <!-- end messages -->


<a name="resources-centrum-DispatchAttribute"></a>

### DispatchAttribute


| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISPATCH_ATTRIBUTE_UNSPECIFIED` | 0 |  |
| `DISPATCH_ATTRIBUTE_MULTIPLE` | 1 |  |
| `DISPATCH_ATTRIBUTE_DUPLICATE` | 2 |  |
| `DISPATCH_ATTRIBUTE_TOO_OLD` | 3 |  |
| `DISPATCH_ATTRIBUTE_AUTOMATIC` | 4 |  |



<a name="resources-centrum-UnitAttribute"></a>

### UnitAttribute


| Name | Number | Description |
| ---- | ------ | ----------- |
| `UNIT_ATTRIBUTE_UNSPECIFIED` | 0 |  |
| `UNIT_ATTRIBUTE_STATIC` | 1 |  |
| `UNIT_ATTRIBUTE_NO_DISPATCH_AUTO_ASSIGN` | 2 |  |


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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `name` | [string](#string) |  | @sanitize |
| `initials` | [string](#string) |  | @sanitize |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize |
| `status` | [UnitStatus](#resources-centrum-UnitStatus) | optional |  |
| `users` | [UnitAssignment](#resources-centrum-UnitAssignment) | repeated |  |
| `attributes` | [UnitAttributes](#resources-centrum-UnitAttributes) | optional |  |
| `home_postal` | [string](#string) | optional |  |
| `access` | [UnitAccess](#resources-centrum-UnitAccess) |  |  |






<a name="resources-centrum-UnitAssignment"></a>

### UnitAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"unit_id" |
| `user_id` | [int32](#int32) |  | @gotags: sql:"primary_key" alias:"user_id" |
| `user` | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |






<a name="resources-centrum-UnitAssignments"></a>

### UnitAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `users` | [UnitAssignment](#resources-centrum-UnitAssignment) | repeated |  |






<a name="resources-centrum-UnitStatus"></a>

### UnitStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `unit_id` | [uint64](#uint64) |  |  |
| `unit` | [Unit](#resources-centrum-Unit) | optional |  |
| `status` | [StatusUnit](#resources-centrum-StatusUnit) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |
| `x` | [double](#double) | optional |  |
| `y` | [double](#double) | optional |  |
| `postal` | [string](#string) | optional | @sanitize |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |





 <!-- end messages -->


<a name="resources-centrum-StatusUnit"></a>

### StatusUnit


| Name | Number | Description |
| ---- | ------ | ----------- |
| `STATUS_UNIT_UNSPECIFIED` | 0 |  |
| `STATUS_UNIT_UNKNOWN` | 1 |  |
| `STATUS_UNIT_USER_ADDED` | 2 |  |
| `STATUS_UNIT_USER_REMOVED` | 3 |  |
| `STATUS_UNIT_UNAVAILABLE` | 4 |  |
| `STATUS_UNIT_AVAILABLE` | 5 |  |
| `STATUS_UNIT_ON_BREAK` | 6 |  |
| `STATUS_UNIT_BUSY` | 7 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_units_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/units_access.proto



<a name="resources-centrum-UnitAccess"></a>

### UnitAccess
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [UnitJobAccess](#resources-centrum-UnitJobAccess) | repeated | @gotags: alias:"job_access" |
| `qualifications` | [UnitQualificationAccess](#resources-centrum-UnitQualificationAccess) | repeated | @gotags: alias:"qualification_access" |






<a name="resources-centrum-UnitJobAccess"></a>

### UnitJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [UnitAccessLevel](#resources-centrum-UnitAccessLevel) |  |  |






<a name="resources-centrum-UnitQualificationAccess"></a>

### UnitQualificationAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [resources.qualifications.QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| `access` | [UnitAccessLevel](#resources-centrum-UnitAccessLevel) |  |  |






<a name="resources-centrum-UnitUserAccess"></a>

### UnitUserAccess






 <!-- end messages -->


<a name="resources-centrum-UnitAccessLevel"></a>

### UnitAccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `UNIT_ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `UNIT_ACCESS_LEVEL_BLOCKED` | 1 |  |
| `UNIT_ACCESS_LEVEL_JOIN` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/access.proto



<a name="resources-centrum-CentrumAccess"></a>

### CentrumAccess
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [CentrumJobAccess](#resources-centrum-CentrumJobAccess) | repeated | @gotags: alias:"job_access" |






<a name="resources-centrum-CentrumJobAccess"></a>

### CentrumJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [CentrumAccessLevel](#resources-centrum-CentrumAccessLevel) |  |  |






<a name="resources-centrum-CentrumQualificationAccess"></a>

### CentrumQualificationAccess
Dummy - DO NOT USE!






<a name="resources-centrum-CentrumUserAccess"></a>

### CentrumUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


<a name="resources-centrum-CentrumAccessLevel"></a>

### CentrumAccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_VIEW` | 1 |  |
| `ACCESS_LEVEL_PARTICIPATE` | 2 |  |
| `ACCESS_LEVEL_DISPATCH` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_dispatchers-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/dispatchers.proto



<a name="resources-centrum-Dispatchers"></a>

### Dispatchers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `dispatchers` | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `jobs` | [string](#string) | repeated |  |
| `status` | [DispatchStatus](#resources-centrum-DispatchStatus) | optional |  |
| `message` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `attributes` | [DispatchAttributes](#resources-centrum-DispatchAttributes) | optional |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `postal` | [string](#string) | optional | @sanitize |
| `anon` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.User](#resources-users-User) | optional |  |
| `units` | [DispatchAssignment](#resources-centrum-DispatchAssignment) | repeated |  |
| `references` | [DispatchReferences](#resources-centrum-DispatchReferences) | optional |  |






<a name="resources-centrum-DispatchAssignment"></a>

### DispatchAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"dispatch_id" |
| `unit_id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"unit_id" |
| `unit` | [Unit](#resources-centrum-Unit) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `expires_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="resources-centrum-DispatchAssignments"></a>

### DispatchAssignments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `units` | [DispatchAssignment](#resources-centrum-DispatchAssignment) | repeated |  |






<a name="resources-centrum-DispatchReference"></a>

### DispatchReference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_dispatch_id` | [uint64](#uint64) |  |  |
| `reference_type` | [DispatchReferenceType](#resources-centrum-DispatchReferenceType) |  |  |






<a name="resources-centrum-DispatchReferences"></a>

### DispatchReferences
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `references` | [DispatchReference](#resources-centrum-DispatchReference) | repeated |  |






<a name="resources-centrum-DispatchStatus"></a>

### DispatchStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `unit_id` | [uint64](#uint64) | optional |  |
| `unit` | [Unit](#resources-centrum-Unit) | optional |  |
| `status` | [StatusDispatch](#resources-centrum-StatusDispatch) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.Colleague](#resources-jobs-Colleague) | optional |  |
| `x` | [double](#double) | optional |  |
| `y` | [double](#double) | optional |  |
| `postal` | [string](#string) | optional | @sanitize |





 <!-- end messages -->


<a name="resources-centrum-DispatchReferenceType"></a>

### DispatchReferenceType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISPATCH_REFERENCE_TYPE_UNSPECIFIED` | 0 |  |
| `DISPATCH_REFERENCE_TYPE_REFERENCED` | 1 |  |
| `DISPATCH_REFERENCE_TYPE_DUPLICATED_BY` | 2 |  |
| `DISPATCH_REFERENCE_TYPE_DUPLICATE_OF` | 3 |  |



<a name="resources-centrum-StatusDispatch"></a>

### StatusDispatch


| Name | Number | Description |
| ---- | ------ | ----------- |
| `STATUS_DISPATCH_UNSPECIFIED` | 0 |  |
| `STATUS_DISPATCH_NEW` | 1 |  |
| `STATUS_DISPATCH_UNASSIGNED` | 2 |  |
| `STATUS_DISPATCH_UPDATED` | 3 |  |
| `STATUS_DISPATCH_UNIT_ASSIGNED` | 4 |  |
| `STATUS_DISPATCH_UNIT_UNASSIGNED` | 5 |  |
| `STATUS_DISPATCH_UNIT_ACCEPTED` | 6 |  |
| `STATUS_DISPATCH_UNIT_DECLINED` | 7 |  |
| `STATUS_DISPATCH_EN_ROUTE` | 8 |  |
| `STATUS_DISPATCH_ON_SCENE` | 9 |  |
| `STATUS_DISPATCH_NEED_ASSISTANCE` | 10 |  |
| `STATUS_DISPATCH_COMPLETED` | 11 |  |
| `STATUS_DISPATCH_CANCELLED` | 12 |  |
| `STATUS_DISPATCH_ARCHIVED` | 13 |  |



<a name="resources-centrum-TakeDispatchResp"></a>

### TakeDispatchResp


| Name | Number | Description |
| ---- | ------ | ----------- |
| `TAKE_DISPATCH_RESP_UNSPECIFIED` | 0 |  |
| `TAKE_DISPATCH_RESP_TIMEOUT` | 1 |  |
| `TAKE_DISPATCH_RESP_ACCEPTED` | 2 |  |
| `TAKE_DISPATCH_RESP_DECLINED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_centrum_settings-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/centrum/settings.proto



<a name="resources-centrum-PredefinedStatus"></a>

### PredefinedStatus
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_status` | [string](#string) | repeated | @sanitize: method=StripTags |
| `dispatch_status` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-centrum-Settings"></a>

### Settings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |
| `type` | [CentrumType](#resources-centrum-CentrumType) |  |  |
| `public` | [bool](#bool) |  |  |
| `mode` | [CentrumMode](#resources-centrum-CentrumMode) |  |  |
| `fallback_mode` | [CentrumMode](#resources-centrum-CentrumMode) |  |  |
| `predefined_status` | [PredefinedStatus](#resources-centrum-PredefinedStatus) | optional |  |
| `timings` | [Timings](#resources-centrum-Timings) |  |  |
| `access` | [CentrumAccess](#resources-centrum-CentrumAccess) |  |  |






<a name="resources-centrum-Timings"></a>

### Timings
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_max_wait` | [int64](#int64) |  |  |
| `require_unit` | [bool](#bool) |  |  |
| `require_unit_reminder_seconds` | [int64](#int64) |  |  |





 <!-- end messages -->


<a name="resources-centrum-CentrumMode"></a>

### CentrumMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_MODE_UNSPECIFIED` | 0 |  |
| `CENTRUM_MODE_MANUAL` | 1 |  |
| `CENTRUM_MODE_CENTRAL_COMMAND` | 2 |  |
| `CENTRUM_MODE_AUTO_ROUND_ROBIN` | 3 |  |
| `CENTRUM_MODE_SIMPLIFIED` | 4 |  |



<a name="resources-centrum-CentrumType"></a>

### CentrumType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_TYPE_UNSPECIFIED` | 0 |  |
| `CENTRUM_TYPE_DISPATCH` | 1 |  |
| `CENTRUM_TYPE_DELIVERY` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_database_database-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/database/database.proto



<a name="resources-common-database-DateRange"></a>

### DateRange
Datetime range (uses Timestamp underneath) It depends on the API method if it will use date or date + time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `start` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | Start time |
| `end` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | End time |






<a name="resources-common-database-PaginationRequest"></a>

### PaginationRequest
Pagination for requests to the server


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `offset` | [int64](#int64) |  |  |
| `page_size` | [int64](#int64) | optional |  |






<a name="resources-common-database-PaginationResponse"></a>

### PaginationResponse
Server Pagination Response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total_count` | [int64](#int64) |  |  |
| `offset` | [int64](#int64) |  |  |
| `end` | [int64](#int64) |  |  |
| `page_size` | [int64](#int64) |  |  |






<a name="resources-common-database-Sort"></a>

### Sort
Sort by column


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `column` | [string](#string) |  | Column name |
| `direction` | [string](#string) |  | Sort direction, must be `asc` (ascending) or `desc` (descending) |





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
| `data` | [bytes](#bytes) |  |  |
| `complete` | [bool](#bool) |  |  |






<a name="resources-common-grpcws-Cancel"></a>

### Cancel







<a name="resources-common-grpcws-Complete"></a>

### Complete







<a name="resources-common-grpcws-Failure"></a>

### Failure



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `error_message` | [string](#string) |  |  |
| `error_status` | [string](#string) |  |  |
| `headers` | [Failure.HeadersEntry](#resources-common-grpcws-Failure-HeadersEntry) | repeated |  |






<a name="resources-common-grpcws-Failure-HeadersEntry"></a>

### Failure.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [HeaderValue](#resources-common-grpcws-HeaderValue) |  |  |






<a name="resources-common-grpcws-GrpcFrame"></a>

### GrpcFrame



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `streamId` | [uint32](#uint32) |  |  |
| `ping` | [Ping](#resources-common-grpcws-Ping) |  |  |
| `header` | [Header](#resources-common-grpcws-Header) |  |  |
| `body` | [Body](#resources-common-grpcws-Body) |  |  |
| `complete` | [Complete](#resources-common-grpcws-Complete) |  |  |
| `failure` | [Failure](#resources-common-grpcws-Failure) |  |  |
| `cancel` | [Cancel](#resources-common-grpcws-Cancel) |  |  |






<a name="resources-common-grpcws-Header"></a>

### Header



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operation` | [string](#string) |  |  |
| `headers` | [Header.HeadersEntry](#resources-common-grpcws-Header-HeadersEntry) | repeated |  |
| `status` | [int32](#int32) |  |  |






<a name="resources-common-grpcws-Header-HeadersEntry"></a>

### Header.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [HeaderValue](#resources-common-grpcws-HeaderValue) |  |  |






<a name="resources-common-grpcws-HeaderValue"></a>

### HeaderValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [string](#string) | repeated |  |






<a name="resources-common-grpcws-Ping"></a>

### Ping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pong` | [bool](#bool) |  |  |





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
| `name` | [string](#string) |  | Cronjob name |
| `schedule` | [string](#string) |  | Cron schedule expression For available valid expressions, see [adhocore/gronx - Cron Expressions Documentation](https://github.com/adhocore/gronx/blob/fea40e3e90e70476877cfb9b50fac10c7de41c5c/README.md#cron-expression).

To generate Cronjob schedule expressions, you can also use web tools like https://crontab.guru/. |
| `state` | [CronjobState](#resources-common-cron-CronjobState) |  | Cronjob state |
| `next_schedule_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | Next time the cronjob should be run |
| `last_attempt_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional | Last attempted start time of Cronjob |
| `started_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional | Time current cronjob was started |
| `timeout` | [google.protobuf.Duration](#google-protobuf-Duration) | optional | Optional timeout for cronjob execution |
| `data` | [CronjobData](#resources-common-cron-CronjobData) |  | Cronjob data |
| `last_completed_event` | [CronjobCompletedEvent](#resources-common-cron-CronjobCompletedEvent) | optional | Last event info to ease debugging and tracking |






<a name="resources-common-cron-CronjobCompletedEvent"></a>

### CronjobCompletedEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Cronjob name |
| `success` | [bool](#bool) |  | Cronjob execution success status |
| `cancelled` | [bool](#bool) |  | Cronjob execution was cancelled |
| `endDate` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | Cronjob end time |
| `elapsed` | [google.protobuf.Duration](#google-protobuf-Duration) |  | Cronjob execution time/elapsed time |
| `data` | [CronjobData](#resources-common-cron-CronjobData) | optional | Cronjob data (can be empty if not touched by the Cronjob handler) |
| `node_name` | [string](#string) |  | Name of the node where the cronjob was executed |
| `error_message` | [string](#string) | optional | Error message (if success = false) |






<a name="resources-common-cron-CronjobData"></a>

### CronjobData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `data` | [google.protobuf.Any](#google-protobuf-Any) | optional |  |






<a name="resources-common-cron-CronjobLockOwnerState"></a>

### CronjobLockOwnerState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hostname` | [string](#string) |  | Hostname of the agent the cronjob is running on |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="resources-common-cron-CronjobSchedulerEvent"></a>

### CronjobSchedulerEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cronjob` | [Cronjob](#resources-common-cron-Cronjob) |  | Full Cronjob spec |






<a name="resources-common-cron-GenericCronData"></a>

### GenericCronData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attributes` | [GenericCronData.AttributesEntry](#resources-common-cron-GenericCronData-AttributesEntry) | repeated | @sanitize: method=StripTags |






<a name="resources-common-cron-GenericCronData-AttributesEntry"></a>

### GenericCronData.AttributesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |





 <!-- end messages -->


<a name="resources-common-cron-CronjobState"></a>

### CronjobState
States of Cronjbo

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CRONJOB_STATE_UNSPECIFIED` | 0 |  |
| `CRONJOB_STATE_WAITING` | 1 |  |
| `CRONJOB_STATE_PENDING` | 2 |  |
| `CRONJOB_STATE_RUNNING` | 3 |  |


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
| `version` | [string](#string) | optional |  |
| `content` | [JSONNode](#resources-common-content-JSONNode) | optional |  |
| `raw_content` | [string](#string) | optional | @sanitize |






<a name="resources-common-content-JSONNode"></a>

### JSONNode



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [NodeType](#resources-common-content-NodeType) |  |  |
| `id` | [string](#string) | optional | @sanitize: method=StripTags |
| `tag` | [string](#string) |  | @sanitize: method=StripTags |
| `attrs` | [JSONNode.AttrsEntry](#resources-common-content-JSONNode-AttrsEntry) | repeated | @sanitize: method=StripTags |
| `text` | [string](#string) | optional | @sanitize: method=StripTags |
| `content` | [JSONNode](#resources-common-content-JSONNode) | repeated |  |






<a name="resources-common-content-JSONNode-AttrsEntry"></a>

### JSONNode.AttrsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |





 <!-- end messages -->


<a name="resources-common-content-ContentType"></a>

### ContentType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `CONTENT_TYPE_UNSPECIFIED` | 0 |  |
| `CONTENT_TYPE_HTML` | 1 |  |
| `CONTENT_TYPE_PLAIN` | 2 |  |
| `CONTENT_TYPE_TIPTAP_JSON` | 3 |  |



<a name="resources-common-content-NodeType"></a>

### NodeType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `NODE_TYPE_UNSPECIFIED` | 0 |  |
| `NODE_TYPE_DOC` | 1 |  |
| `NODE_TYPE_ELEMENT` | 2 |  |
| `NODE_TYPE_TEXT` | 3 |  |
| `NODE_TYPE_COMMENT` | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_tests_objects-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/tests/objects.proto



<a name="resources-common-tests-SimpleObject"></a>

### SimpleObject
INTERNAL ONLY** SimpleObject is used as a test object where proto-based messages are required.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `field1` | [string](#string) |  |  |
| `field2` | [bool](#bool) |  |  |





 <!-- end messages -->

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
| `title` | [I18NItem](#resources-common-I18NItem) | optional |  |
| `content` | [I18NItem](#resources-common-I18NItem) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_common_i18n-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/common/i18n.proto



<a name="resources-common-I18NItem"></a>

### I18NItem
Wrapped translated message for the client @dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  | @sanitize: method=StripTags |
| `parameters` | [I18NItem.ParametersEntry](#resources-common-I18NItem-ParametersEntry) | repeated | @sanitize: method=StripTags |






<a name="resources-common-I18NItem-ParametersEntry"></a>

### I18NItem.ParametersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |





 <!-- end messages -->

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
| `uuid` | [string](#string) |  |  |





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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `job` | [string](#string) | optional |  |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `document_id` | [uint64](#uint64) |  |  |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `request_type` | [DocActivityType](#resources-documents-DocActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [DocActivityData](#resources-documents-DocActivityData) |  |  |
| `accepted` | [bool](#bool) | optional |  |





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
| `required` | [bool](#bool) | optional |  |
| `min` | [int32](#int32) | optional |  |
| `max` | [int32](#int32) | optional |  |






<a name="resources-documents-Template"></a>

### Template



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `category` | [Category](#resources-documents-Category) |  | @gotags: alias:"category" |
| `weight` | [uint32](#uint32) |  |  |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) |  | @sanitize |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `content_title` | [string](#string) |  | @gotags: alias:"content_title" |
| `content` | [string](#string) |  | @gotags: alias:"content" |
| `state` | [string](#string) |  | @gotags: alias:"state" |
| `schema` | [TemplateSchema](#resources-documents-TemplateSchema) |  | @gotags: alias:"schema" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `job_access` | [TemplateJobAccess](#resources-documents-TemplateJobAccess) | repeated |  |
| `content_access` | [DocumentAccess](#resources-documents-DocumentAccess) |  | @gotags: alias:"access" |
| `workflow` | [Workflow](#resources-documents-Workflow) | optional |  |






<a name="resources-documents-TemplateData"></a>

### TemplateData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `activeChar` | [resources.users.User](#resources-users-User) |  |  |
| `documents` | [DocumentShort](#resources-documents-DocumentShort) | repeated |  |
| `users` | [resources.users.UserShort](#resources-users-UserShort) | repeated |  |
| `vehicles` | [resources.vehicles.Vehicle](#resources-vehicles-Vehicle) | repeated |  |






<a name="resources-documents-TemplateJobAccess"></a>

### TemplateJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  | @gotags: alias:"template_id" |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-documents-AccessLevel) |  |  |






<a name="resources-documents-TemplateRequirements"></a>

### TemplateRequirements



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `documents` | [ObjectSpecs](#resources-documents-ObjectSpecs) | optional |  |
| `users` | [ObjectSpecs](#resources-documents-ObjectSpecs) | optional |  |
| `vehicles` | [ObjectSpecs](#resources-documents-ObjectSpecs) | optional |  |






<a name="resources-documents-TemplateSchema"></a>

### TemplateSchema
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `requirements` | [TemplateRequirements](#resources-documents-TemplateRequirements) |  |  |






<a name="resources-documents-TemplateShort"></a>

### TemplateShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `category` | [Category](#resources-documents-Category) |  | @gotags: alias:"category" |
| `weight` | [uint32](#uint32) |  |  |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) |  | @sanitize |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `schema` | [TemplateSchema](#resources-documents-TemplateSchema) |  | @gotags: alias:"schema" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `workflow` | [Workflow](#resources-documents-Workflow) | optional |  |






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
| `duration` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| `message` | [string](#string) |  |  |






<a name="resources-documents-Reminder"></a>

### Reminder



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| `message` | [string](#string) |  |  |






<a name="resources-documents-ReminderSettings"></a>

### ReminderSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reminders` | [Reminder](#resources-documents-Reminder) | repeated |  |






<a name="resources-documents-Workflow"></a>

### Workflow
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reminder` | [bool](#bool) |  |  |
| `reminder_settings` | [ReminderSettings](#resources-documents-ReminderSettings) |  |  |
| `auto_close` | [bool](#bool) |  |  |
| `auto_close_settings` | [AutoCloseSettings](#resources-documents-AutoCloseSettings) |  |  |






<a name="resources-documents-WorkflowCronData"></a>

### WorkflowCronData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_doc_id` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/access.proto



<a name="resources-documents-DocumentAccess"></a>

### DocumentAccess
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated | @gotags: alias:"job_access" |
| `users` | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-documents-DocumentJobAccess"></a>

### DocumentJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-documents-AccessLevel) |  |  |
| `required` | [bool](#bool) | optional |  |






<a name="resources-documents-DocumentUserAccess"></a>

### DocumentUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `access` | [AccessLevel](#resources-documents-AccessLevel) |  |  |
| `required` | [bool](#bool) | optional |  |





 <!-- end messages -->


<a name="resources-documents-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_BLOCKED` | 1 |  |
| `ACCESS_LEVEL_VIEW` | 2 |  |
| `ACCESS_LEVEL_COMMENT` | 3 |  |
| `ACCESS_LEVEL_STATUS` | 4 |  |
| `ACCESS_LEVEL_ACCESS` | 5 |  |
| `ACCESS_LEVEL_EDIT` | 6 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_documents_pins-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/documents/pins.proto



<a name="resources-documents-DocumentPin"></a>

### DocumentPin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" |
| `job` | [string](#string) | optional | @gotags: sql:"primary_key" |
| `user_id` | [int32](#int32) | optional | @gotags: sql:"primary_key" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `state` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) |  |  |





 <!-- end messages -->

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
| `to_create` | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated |  |
| `to_update` | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated |  |
| `to_delete` | [DocumentJobAccess](#resources-documents-DocumentJobAccess) | repeated |  |






<a name="resources-documents-DocAccessRequested"></a>

### DocAccessRequested



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `level` | [AccessLevel](#resources-documents-AccessLevel) |  |  |






<a name="resources-documents-DocAccessUpdated"></a>

### DocAccessUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [DocAccessJobsDiff](#resources-documents-DocAccessJobsDiff) |  |  |
| `users` | [DocAccessUsersDiff](#resources-documents-DocAccessUsersDiff) |  |  |






<a name="resources-documents-DocAccessUsersDiff"></a>

### DocAccessUsersDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated |  |
| `to_update` | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated |  |
| `to_delete` | [DocumentUserAccess](#resources-documents-DocumentUserAccess) | repeated |  |






<a name="resources-documents-DocActivity"></a>

### DocActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `activity_type` | [DocActivityType](#resources-documents-DocActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [DocActivityData](#resources-documents-DocActivityData) |  |  |






<a name="resources-documents-DocActivityData"></a>

### DocActivityData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [DocUpdated](#resources-documents-DocUpdated) |  |  |
| `owner_changed` | [DocOwnerChanged](#resources-documents-DocOwnerChanged) |  |  |
| `access_updated` | [DocAccessUpdated](#resources-documents-DocAccessUpdated) |  |  |
| `access_requested` | [DocAccessRequested](#resources-documents-DocAccessRequested) |  |  |






<a name="resources-documents-DocFilesChange"></a>

### DocFilesChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [int64](#int64) |  |  |
| `deleted` | [int64](#int64) |  |  |






<a name="resources-documents-DocOwnerChanged"></a>

### DocOwnerChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_owner_id` | [int32](#int32) |  |  |
| `new_owner` | [resources.users.UserShort](#resources-users-UserShort) |  |  |






<a name="resources-documents-DocUpdated"></a>

### DocUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title_diff` | [string](#string) | optional |  |
| `content_diff` | [string](#string) | optional |  |
| `state_diff` | [string](#string) | optional |  |
| `files_change` | [DocFilesChange](#resources-documents-DocFilesChange) | optional |  |





 <!-- end messages -->


<a name="resources-documents-DocActivityType"></a>

### DocActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `DOC_ACTIVITY_TYPE_CREATED` | 1 | Base |
| `DOC_ACTIVITY_TYPE_STATUS_OPEN` | 2 |  |
| `DOC_ACTIVITY_TYPE_STATUS_CLOSED` | 3 |  |
| `DOC_ACTIVITY_TYPE_UPDATED` | 4 |  |
| `DOC_ACTIVITY_TYPE_RELATIONS_UPDATED` | 5 |  |
| `DOC_ACTIVITY_TYPE_REFERENCES_UPDATED` | 6 |  |
| `DOC_ACTIVITY_TYPE_ACCESS_UPDATED` | 7 |  |
| `DOC_ACTIVITY_TYPE_OWNER_CHANGED` | 8 |  |
| `DOC_ACTIVITY_TYPE_DELETED` | 9 |  |
| `DOC_ACTIVITY_TYPE_DRAFT_TOGGLED` | 19 |  |
| `DOC_ACTIVITY_TYPE_COMMENT_ADDED` | 10 | Comments |
| `DOC_ACTIVITY_TYPE_COMMENT_UPDATED` | 11 |  |
| `DOC_ACTIVITY_TYPE_COMMENT_DELETED` | 12 |  |
| `DOC_ACTIVITY_TYPE_REQUESTED_ACCESS` | 13 | Requests |
| `DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE` | 14 |  |
| `DOC_ACTIVITY_TYPE_REQUESTED_OPENING` | 15 |  |
| `DOC_ACTIVITY_TYPE_REQUESTED_UPDATE` | 16 |  |
| `DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE` | 17 |  |
| `DOC_ACTIVITY_TYPE_REQUESTED_DELETION` | 18 |  |


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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `category_id` | [uint64](#uint64) | optional |  |
| `category` | [Category](#resources-documents-Category) | optional | @gotags: alias:"category" |
| `title` | [string](#string) |  | @sanitize |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `data` | [string](#string) | optional | @sanitize

@gotags: alias:"data" |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `state` | [string](#string) |  | @sanitize |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `template_id` | [uint64](#uint64) | optional |  |
| `pin` | [DocumentPin](#resources-documents-DocumentPin) | optional | @gotags: alias:"pin" |
| `workflow_state` | [WorkflowState](#resources-documents-WorkflowState) | optional |  |
| `workflow_user` | [WorkflowUserState](#resources-documents-WorkflowUserState) | optional |  |
| `files` | [resources.file.File](#resources-file-File) | repeated | @gotags: alias:"files" |






<a name="resources-documents-DocumentReference"></a>

### DocumentReference



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `source_document_id` | [uint64](#uint64) |  | @gotags: alias:"source_document_id" |
| `source_document` | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:"source_document" |
| `reference` | [DocReference](#resources-documents-DocReference) |  | @gotags: alias:"reference" |
| `target_document_id` | [uint64](#uint64) |  | @gotags: alias:"target_document_id" |
| `target_document` | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:"target_document" |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"ref_creator" |






<a name="resources-documents-DocumentRelation"></a>

### DocumentRelation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `document_id` | [uint64](#uint64) |  |  |
| `document` | [DocumentShort](#resources-documents-DocumentShort) | optional | @gotags: alias:"document" |
| `source_user_id` | [int32](#int32) |  | @gotags: alias:"source_user_id" |
| `source_user` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"source_user" |
| `relation` | [DocRelation](#resources-documents-DocRelation) |  | @gotags: alias:"relation" |
| `target_user_id` | [int32](#int32) |  | @gotags: alias:"target_user_id" |
| `target_user` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"target_user" |






<a name="resources-documents-DocumentShort"></a>

### DocumentShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `category_id` | [uint64](#uint64) | optional |  |
| `category` | [Category](#resources-documents-Category) | optional | @gotags: alias:"category" |
| `title` | [string](#string) |  | @sanitize |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `state` | [string](#string) |  | @sanitize |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `pin` | [DocumentPin](#resources-documents-DocumentPin) | optional | @gotags: alias:"pin" |
| `workflow_state` | [WorkflowState](#resources-documents-WorkflowState) | optional |  |
| `workflow_user` | [WorkflowUserState](#resources-documents-WorkflowUserState) | optional |  |






<a name="resources-documents-WorkflowState"></a>

### WorkflowState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `next_reminder_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `next_reminder_count` | [int32](#int32) | optional |  |
| `auto_close_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `workflow` | [Workflow](#resources-documents-Workflow) | optional | @gotags: alias:"workflow" |
| `document` | [DocumentShort](#resources-documents-DocumentShort) | optional |  |






<a name="resources-documents-WorkflowUserState"></a>

### WorkflowUserState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `manual_reminder_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `manual_reminder_message` | [string](#string) | optional |  |
| `workflow` | [Workflow](#resources-documents-Workflow) | optional | @gotags: alias:"workflow" |
| `document` | [DocumentShort](#resources-documents-DocumentShort) | optional |  |





 <!-- end messages -->


<a name="resources-documents-DocReference"></a>

### DocReference


| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_REFERENCE_UNSPECIFIED` | 0 |  |
| `DOC_REFERENCE_LINKED` | 1 |  |
| `DOC_REFERENCE_SOLVES` | 2 |  |
| `DOC_REFERENCE_CLOSES` | 3 |  |
| `DOC_REFERENCE_DEPRECATES` | 4 |  |



<a name="resources-documents-DocRelation"></a>

### DocRelation


| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_RELATION_UNSPECIFIED` | 0 |  |
| `DOC_RELATION_MENTIONED` | 1 |  |
| `DOC_RELATION_TARGETS` | 2 |  |
| `DOC_RELATION_CAUSED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_activity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/activity.proto



<a name="resources-jobs-AbsenceDateChange"></a>

### AbsenceDateChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_begin` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `absence_end` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="resources-jobs-ColleagueActivity"></a>

### ColleagueActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `source_user_id` | [int32](#int32) | optional |  |
| `source_user` | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:"source_user" |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [Colleague](#resources-jobs-Colleague) |  | @gotags: alias:"target_user" |
| `activity_type` | [ColleagueActivityType](#resources-jobs-ColleagueActivityType) |  |  |
| `reason` | [string](#string) |  | @sanitize |
| `data` | [ColleagueActivityData](#resources-jobs-ColleagueActivityData) |  |  |






<a name="resources-jobs-ColleagueActivityData"></a>

### ColleagueActivityData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_date` | [AbsenceDateChange](#resources-jobs-AbsenceDateChange) |  |  |
| `grade_change` | [GradeChange](#resources-jobs-GradeChange) |  |  |
| `labels_change` | [LabelsChange](#resources-jobs-LabelsChange) |  |  |
| `name_change` | [NameChange](#resources-jobs-NameChange) |  |  |






<a name="resources-jobs-GradeChange"></a>

### GradeChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grade` | [int32](#int32) |  |  |
| `grade_label` | [string](#string) |  |  |






<a name="resources-jobs-LabelsChange"></a>

### LabelsChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [Label](#resources-jobs-Label) | repeated |  |
| `removed` | [Label](#resources-jobs-Label) | repeated |  |






<a name="resources-jobs-NameChange"></a>

### NameChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prefix` | [string](#string) | optional |  |
| `suffix` | [string](#string) | optional |  |





 <!-- end messages -->


<a name="resources-jobs-ColleagueActivityType"></a>

### ColleagueActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `COLLEAGUE_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `COLLEAGUE_ACTIVITY_TYPE_HIRED` | 1 |  |
| `COLLEAGUE_ACTIVITY_TYPE_FIRED` | 2 |  |
| `COLLEAGUE_ACTIVITY_TYPE_PROMOTED` | 3 |  |
| `COLLEAGUE_ACTIVITY_TYPE_DEMOTED` | 4 |  |
| `COLLEAGUE_ACTIVITY_TYPE_ABSENCE_DATE` | 5 |  |
| `COLLEAGUE_ACTIVITY_TYPE_NOTE` | 6 |  |
| `COLLEAGUE_ACTIVITY_TYPE_LABELS` | 7 |  |
| `COLLEAGUE_ACTIVITY_TYPE_NAME` | 8 |  |


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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `type` | [ConductType](#resources-jobs-ConductType) |  |  |
| `message` | [string](#string) |  | @sanitize |
| `expires_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:"target_user" |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [Colleague](#resources-jobs-Colleague) | optional | @gotags: alias:"creator" |





 <!-- end messages -->


<a name="resources-jobs-ConductType"></a>

### ConductType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `CONDUCT_TYPE_UNSPECIFIED` | 0 |  |
| `CONDUCT_TYPE_NEUTRAL` | 1 |  |
| `CONDUCT_TYPE_POSITIVE` | 2 |  |
| `CONDUCT_TYPE_NEGATIVE` | 3 |  |
| `CONDUCT_TYPE_WARNING` | 4 |  |
| `CONDUCT_TYPE_SUSPENSION` | 5 |  |
| `CONDUCT_TYPE_NOTE` | 6 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_job_settings-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/job_settings.proto



<a name="resources-jobs-DiscordSyncChange"></a>

### DiscordSyncChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `plan` | [string](#string) |  |  |






<a name="resources-jobs-DiscordSyncChanges"></a>

### DiscordSyncChanges
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `changes` | [DiscordSyncChange](#resources-jobs-DiscordSyncChange) | repeated |  |






<a name="resources-jobs-DiscordSyncSettings"></a>

### DiscordSyncSettings
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dry_run` | [bool](#bool) |  |  |
| `user_info_sync` | [bool](#bool) |  |  |
| `user_info_sync_settings` | [UserInfoSyncSettings](#resources-jobs-UserInfoSyncSettings) |  |  |
| `status_log` | [bool](#bool) |  |  |
| `status_log_settings` | [StatusLogSettings](#resources-jobs-StatusLogSettings) |  |  |
| `jobs_absence` | [bool](#bool) |  |  |
| `jobs_absence_settings` | [JobsAbsenceSettings](#resources-jobs-JobsAbsenceSettings) |  |  |
| `group_sync_settings` | [GroupSyncSettings](#resources-jobs-GroupSyncSettings) |  |  |
| `qualifications_role_format` | [string](#string) |  |  |






<a name="resources-jobs-GroupMapping"></a>

### GroupMapping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `from_grade` | [int32](#int32) |  |  |
| `to_grade` | [int32](#int32) |  |  |






<a name="resources-jobs-GroupSyncSettings"></a>

### GroupSyncSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ignored_role_ids` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-jobs-JobSettings"></a>

### JobSettings
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_past_days` | [int32](#int32) |  |  |
| `absence_future_days` | [int32](#int32) |  |  |






<a name="resources-jobs-JobsAbsenceSettings"></a>

### JobsAbsenceSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_role` | [string](#string) |  |  |






<a name="resources-jobs-StatusLogSettings"></a>

### StatusLogSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel_id` | [string](#string) |  |  |






<a name="resources-jobs-UserInfoSyncSettings"></a>

### UserInfoSyncSettings



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `employee_role_enabled` | [bool](#bool) |  |  |
| `employee_role_format` | [string](#string) |  |  |
| `grade_role_format` | [string](#string) |  |  |
| `unemployed_enabled` | [bool](#bool) |  |  |
| `unemployed_mode` | [UserInfoSyncUnemployedMode](#resources-jobs-UserInfoSyncUnemployedMode) |  |  |
| `unemployed_role_name` | [string](#string) |  |  |
| `sync_nicknames` | [bool](#bool) |  |  |
| `group_mapping` | [GroupMapping](#resources-jobs-GroupMapping) | repeated |  |





 <!-- end messages -->


<a name="resources-jobs-UserInfoSyncUnemployedMode"></a>

### UserInfoSyncUnemployedMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_UNSPECIFIED` | 0 |  |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE` | 1 |  |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_KICK` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_jobs-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/jobs.proto



<a name="resources-jobs-Job"></a>

### Job



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | @gotags: sql:"primary_key" alias:"name" |
| `label` | [string](#string) |  |  |
| `grades` | [JobGrade](#resources-jobs-JobGrade) | repeated |  |






<a name="resources-jobs-JobGrade"></a>

### JobGrade



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_name` | [string](#string) | optional |  |
| `grade` | [int32](#int32) |  |  |
| `label` | [string](#string) |  |  |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `job` | [string](#string) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `order` | [int32](#int32) |  |  |






<a name="resources-jobs-LabelCount"></a>

### LabelCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `label` | [Label](#resources-jobs-Label) |  |  |
| `count` | [int64](#int64) |  |  |






<a name="resources-jobs-Labels"></a>

### Labels



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [Label](#resources-jobs-Label) | repeated |  |





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
| `user_id` | [int32](#int32) |  | @gotags: sql:"primary_key" |
| `job` | [string](#string) |  |  |
| `date` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: sql:"primary_key" |
| `user` | [Colleague](#resources-jobs-Colleague) | optional |  |
| `start_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional | @gotags: sql:"primary_key" |
| `end_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `spent_time` | [float](#float) |  |  |






<a name="resources-jobs-TimeclockStats"></a>

### TimeclockStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `spent_time_sum` | [float](#float) |  |  |
| `spent_time_avg` | [float](#float) |  |  |
| `spent_time_max` | [float](#float) |  |  |






<a name="resources-jobs-TimeclockWeeklyStats"></a>

### TimeclockWeeklyStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `year` | [int32](#int32) |  |  |
| `calendar_week` | [int32](#int32) |  |  |
| `sum` | [float](#float) |  |  |
| `avg` | [float](#float) |  |  |
| `max` | [float](#float) |  |  |





 <!-- end messages -->


<a name="resources-jobs-TimeclockMode"></a>

### TimeclockMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| `TIMECLOCK_MODE_UNSPECIFIED` | 0 |  |
| `TIMECLOCK_MODE_DAILY` | 1 |  |
| `TIMECLOCK_MODE_WEEKLY` | 2 |  |
| `TIMECLOCK_MODE_RANGE` | 3 |  |
| `TIMECLOCK_MODE_TIMELINE` | 4 |  |



<a name="resources-jobs-TimeclockViewMode"></a>

### TimeclockViewMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| `TIMECLOCK_VIEW_MODE_UNSPECIFIED` | 0 |  |
| `TIMECLOCK_VIEW_MODE_SELF` | 1 |  |
| `TIMECLOCK_VIEW_MODE_ALL` | 2 |  |


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
| `user_id` | [int32](#int32) |  | @gotags: alias:"id" |
| `identifier` | [string](#string) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `phone_number` | [string](#string) | optional |  |
| `avatar_file_id` | [uint64](#uint64) | optional |  |
| `avatar` | [string](#string) | optional | @gotags: alias:"avatar" |
| `props` | [ColleagueProps](#resources-jobs-ColleagueProps) |  | @gotags: alias:"colleague_props" |
| `email` | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-jobs-ColleagueProps"></a>

### ColleagueProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `absence_begin` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `absence_end` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `note` | [string](#string) | optional | @sanitize: method=StripTags |
| `labels` | [Labels](#resources-jobs-Labels) | optional |  |
| `name_prefix` | [string](#string) | optional |  |
| `name_suffix` | [string](#string) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_jobs_job_props-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/jobs/job_props.proto



<a name="resources-jobs-JobProps"></a>

### JobProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `livemap_marker_color` | [string](#string) |  |  |
| `quick_buttons` | [QuickButtons](#resources-jobs-QuickButtons) |  |  |
| `radio_frequency` | [string](#string) | optional |  |
| `discord_guild_id` | [string](#string) | optional |  |
| `discord_last_sync` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `discord_sync_settings` | [DiscordSyncSettings](#resources-jobs-DiscordSyncSettings) |  |  |
| `discord_sync_changes` | [DiscordSyncChanges](#resources-jobs-DiscordSyncChanges) | optional |  |
| `motd` | [string](#string) | optional |  |
| `logo_file_id` | [uint64](#uint64) | optional |  |
| `logo_file` | [resources.file.File](#resources-file-File) | optional | @gotags: alias:"logo_file" |
| `settings` | [JobSettings](#resources-jobs-JobSettings) |  |  |






<a name="resources-jobs-QuickButtons"></a>

### QuickButtons
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `penalty_calculator` | [bool](#bool) |  |  |
| `math_calculator` | [bool](#bool) |  |  |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"law.id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `lawbook_id` | [uint64](#uint64) |  |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `hint` | [string](#string) | optional | @sanitize |
| `fine` | [uint32](#uint32) | optional |  |
| `detention_time` | [uint32](#uint32) | optional |  |
| `stvo_points` | [uint32](#uint32) | optional |  |






<a name="resources-laws-LawBook"></a>

### LawBook



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `laws` | [Law](#resources-laws-Law) | repeated |  |





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
| `calendar_id` | [uint64](#uint64) | optional |  |
| `calendar_entry_id` | [uint64](#uint64) | optional |  |






<a name="resources-notifications-Data"></a>

### Data
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `link` | [Link](#resources-notifications-Link) | optional |  |
| `caused_by` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `calendar` | [CalendarData](#resources-notifications-CalendarData) | optional |  |






<a name="resources-notifications-Link"></a>

### Link



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to` | [string](#string) |  |  |
| `title` | [string](#string) | optional |  |
| `external` | [bool](#bool) | optional |  |






<a name="resources-notifications-Notification"></a>

### Notification



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `read_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `title` | [resources.common.I18NItem](#resources-common-I18NItem) |  | @sanitize |
| `type` | [NotificationType](#resources-notifications-NotificationType) |  |  |
| `content` | [resources.common.I18NItem](#resources-common-I18NItem) |  | @sanitize |
| `category` | [NotificationCategory](#resources-notifications-NotificationCategory) |  |  |
| `data` | [Data](#resources-notifications-Data) | optional |  |
| `starred` | [bool](#bool) | optional |  |





 <!-- end messages -->


<a name="resources-notifications-NotificationCategory"></a>

### NotificationCategory


| Name | Number | Description |
| ---- | ------ | ----------- |
| `NOTIFICATION_CATEGORY_UNSPECIFIED` | 0 |  |
| `NOTIFICATION_CATEGORY_GENERAL` | 1 |  |
| `NOTIFICATION_CATEGORY_DOCUMENT` | 2 |  |
| `NOTIFICATION_CATEGORY_CALENDAR` | 3 |  |



<a name="resources-notifications-NotificationType"></a>

### NotificationType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `NOTIFICATION_TYPE_UNSPECIFIED` | 0 |  |
| `NOTIFICATION_TYPE_ERROR` | 1 |  |
| `NOTIFICATION_TYPE_WARNING` | 2 |  |
| `NOTIFICATION_TYPE_INFO` | 3 |  |
| `NOTIFICATION_TYPE_SUCCESS` | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_notifications_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/notifications/events.proto



<a name="resources-notifications-JobEvent"></a>

### JobEvent
Job related events


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resources-jobs-JobProps) |  |  |






<a name="resources-notifications-JobGradeEvent"></a>

### JobGradeEvent
Job grade events


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_token` | [bool](#bool) |  |  |






<a name="resources-notifications-SystemEvent"></a>

### SystemEvent
System related events


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_config` | [resources.clientconfig.ClientConfig](#resources-clientconfig-ClientConfig) |  | Client configuration update (e.g., feature gates, game settings, banner message) |






<a name="resources-notifications-UserEvent"></a>

### UserEvent
User related events


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_token` | [bool](#bool) |  |  |
| `notification` | [Notification](#resources-notifications-Notification) |  | Notifications |
| `notifications_read_count` | [int32](#int32) |  |  |
| `user_info_changed` | [resources.userinfo.UserInfoChanged](#resources-userinfo-UserInfoChanged) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_notifications_client_view-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/notifications/client_view.proto



<a name="resources-notifications-ClientView"></a>

### ClientView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [ObjectType](#resources-notifications-ObjectType) |  |  |
| `id` | [uint64](#uint64) | optional |  |






<a name="resources-notifications-ObjectEvent"></a>

### ObjectEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [ObjectType](#resources-notifications-ObjectType) |  |  |
| `id` | [uint64](#uint64) | optional |  |
| `event_type` | [ObjectEventType](#resources-notifications-ObjectEventType) |  |  |
| `user_id` | [int32](#int32) | optional |  |
| `job` | [string](#string) | optional |  |
| `data` | [google.protobuf.Struct](#google-protobuf-Struct) | optional |  |





 <!-- end messages -->


<a name="resources-notifications-ObjectEventType"></a>

### ObjectEventType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `OBJECT_EVENT_TYPE_UNSPECIFIED` | 0 |  |
| `OBJECT_EVENT_TYPE_UPDATED` | 1 |  |
| `OBJECT_EVENT_TYPE_DELETED` | 2 |  |



<a name="resources-notifications-ObjectType"></a>

### ObjectType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `OBJECT_TYPE_UNSPECIFIED` | 0 |  |
| `OBJECT_TYPE_CITIZEN` | 1 |  |
| `OBJECT_TYPE_DOCUMENT` | 2 |  |
| `OBJECT_TYPE_WIKI_PAGE` | 3 |  |
| `OBJECT_TYPE_JOBS_COLLEAGUE` | 4 |  |
| `OBJECT_TYPE_JOBS_CONDUCT` | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_permissions_attributes-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/permissions/attributes.proto



<a name="resources-permissions-AttributeValues"></a>

### AttributeValues
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `string_list` | [StringList](#resources-permissions-StringList) |  |  |
| `job_list` | [StringList](#resources-permissions-StringList) |  |  |
| `job_grade_list` | [JobGradeList](#resources-permissions-JobGradeList) |  |  |






<a name="resources-permissions-JobGradeList"></a>

### JobGradeList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fine_grained` | [bool](#bool) |  |  |
| `jobs` | [JobGradeList.JobsEntry](#resources-permissions-JobGradeList-JobsEntry) | repeated |  |
| `grades` | [JobGradeList.GradesEntry](#resources-permissions-JobGradeList-GradesEntry) | repeated |  |






<a name="resources-permissions-JobGradeList-GradesEntry"></a>

### JobGradeList.GradesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [JobGrades](#resources-permissions-JobGrades) |  |  |






<a name="resources-permissions-JobGradeList-JobsEntry"></a>

### JobGradeList.JobsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [int32](#int32) |  |  |






<a name="resources-permissions-JobGrades"></a>

### JobGrades



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grades` | [int32](#int32) | repeated |  |






<a name="resources-permissions-RoleAttribute"></a>

### RoleAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `attr_id` | [uint64](#uint64) |  |  |
| `permission_id` | [uint64](#uint64) |  |  |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `key` | [string](#string) |  |  |
| `type` | [string](#string) |  |  |
| `valid_values` | [AttributeValues](#resources-permissions-AttributeValues) |  |  |
| `value` | [AttributeValues](#resources-permissions-AttributeValues) |  |  |
| `max_values` | [AttributeValues](#resources-permissions-AttributeValues) | optional |  |






<a name="resources-permissions-StringList"></a>

### StringList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `strings` | [string](#string) | repeated | @sanitize: method=StripTags |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_permissions_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/permissions/events.proto



<a name="resources-permissions-JobLimitsUpdatedEvent"></a>

### JobLimitsUpdatedEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |






<a name="resources-permissions-RoleIDEvent"></a>

### RoleIDEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_permissions_permissions-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/permissions/permissions.proto



<a name="resources-permissions-PermItem"></a>

### PermItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `val` | [bool](#bool) |  |  |






<a name="resources-permissions-Permission"></a>

### Permission



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `guard_name` | [string](#string) |  |  |
| `val` | [bool](#bool) |  |  |
| `order` | [int32](#int32) | optional |  |






<a name="resources-permissions-Role"></a>

### Role



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `permissions` | [Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |





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
| `jobs` | [QualificationJobAccess](#resources-qualifications-QualificationJobAccess) | repeated |  |






<a name="resources-qualifications-QualificationJobAccess"></a>

### QualificationJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-qualifications-AccessLevel) |  |  |






<a name="resources-qualifications-QualificationUserAccess"></a>

### QualificationUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


<a name="resources-qualifications-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_BLOCKED` | 1 |  |
| `ACCESS_LEVEL_VIEW` | 2 |  |
| `ACCESS_LEVEL_REQUEST` | 3 |  |
| `ACCESS_LEVEL_TAKE` | 4 |  |
| `ACCESS_LEVEL_GRADE` | 5 |  |
| `ACCESS_LEVEL_EDIT` | 6 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_qualifications_exam-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/qualifications/exam.proto



<a name="resources-qualifications-ExamGrading"></a>

### ExamGrading
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `responses` | [ExamGradingResponse](#resources-qualifications-ExamGradingResponse) | repeated |  |






<a name="resources-qualifications-ExamGradingResponse"></a>

### ExamGradingResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `question_id` | [uint64](#uint64) |  |  |
| `points` | [float](#float) |  |  |
| `checked` | [bool](#bool) | optional |  |






<a name="resources-qualifications-ExamQuestion"></a>

### ExamQuestion



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `data` | [ExamQuestionData](#resources-qualifications-ExamQuestionData) |  |  |
| `answer` | [ExamQuestionAnswerData](#resources-qualifications-ExamQuestionAnswerData) | optional |  |
| `points` | [int32](#int32) | optional |  |






<a name="resources-qualifications-ExamQuestionAnswerData"></a>

### ExamQuestionAnswerData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `answer_key` | [string](#string) |  |  |
| `yesno` | [ExamResponseYesNo](#resources-qualifications-ExamResponseYesNo) |  |  |
| `free_text` | [ExamResponseText](#resources-qualifications-ExamResponseText) |  |  |
| `single_choice` | [ExamResponseSingleChoice](#resources-qualifications-ExamResponseSingleChoice) |  |  |
| `multiple_choice` | [ExamResponseMultipleChoice](#resources-qualifications-ExamResponseMultipleChoice) |  |  |






<a name="resources-qualifications-ExamQuestionData"></a>

### ExamQuestionData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `separator` | [ExamQuestionSeparator](#resources-qualifications-ExamQuestionSeparator) |  |  |
| `image` | [ExamQuestionImage](#resources-qualifications-ExamQuestionImage) |  |  |
| `yesno` | [ExamQuestionYesNo](#resources-qualifications-ExamQuestionYesNo) |  |  |
| `free_text` | [ExamQuestionText](#resources-qualifications-ExamQuestionText) |  |  |
| `single_choice` | [ExamQuestionSingleChoice](#resources-qualifications-ExamQuestionSingleChoice) |  |  |
| `multiple_choice` | [ExamQuestionMultipleChoice](#resources-qualifications-ExamQuestionMultipleChoice) |  |  |






<a name="resources-qualifications-ExamQuestionImage"></a>

### ExamQuestionImage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image` | [resources.file.File](#resources-file-File) |  |  |
| `alt` | [string](#string) | optional |  |






<a name="resources-qualifications-ExamQuestionMultipleChoice"></a>

### ExamQuestionMultipleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated | @sanitize: method=StripTags |
| `limit` | [int32](#int32) | optional |  |






<a name="resources-qualifications-ExamQuestionSeparator"></a>

### ExamQuestionSeparator







<a name="resources-qualifications-ExamQuestionSingleChoice"></a>

### ExamQuestionSingleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-qualifications-ExamQuestionText"></a>

### ExamQuestionText



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min_length` | [int32](#int32) |  |  |
| `max_length` | [int32](#int32) |  |  |






<a name="resources-qualifications-ExamQuestionYesNo"></a>

### ExamQuestionYesNo







<a name="resources-qualifications-ExamQuestions"></a>

### ExamQuestions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `questions` | [ExamQuestion](#resources-qualifications-ExamQuestion) | repeated |  |






<a name="resources-qualifications-ExamResponse"></a>

### ExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `question_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `question` | [ExamQuestion](#resources-qualifications-ExamQuestion) |  |  |
| `response` | [ExamResponseData](#resources-qualifications-ExamResponseData) |  |  |






<a name="resources-qualifications-ExamResponseData"></a>

### ExamResponseData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `separator` | [ExamResponseSeparator](#resources-qualifications-ExamResponseSeparator) |  |  |
| `yesno` | [ExamResponseYesNo](#resources-qualifications-ExamResponseYesNo) |  |  |
| `free_text` | [ExamResponseText](#resources-qualifications-ExamResponseText) |  |  |
| `single_choice` | [ExamResponseSingleChoice](#resources-qualifications-ExamResponseSingleChoice) |  |  |
| `multiple_choice` | [ExamResponseMultipleChoice](#resources-qualifications-ExamResponseMultipleChoice) |  |  |






<a name="resources-qualifications-ExamResponseMultipleChoice"></a>

### ExamResponseMultipleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-qualifications-ExamResponseSeparator"></a>

### ExamResponseSeparator







<a name="resources-qualifications-ExamResponseSingleChoice"></a>

### ExamResponseSingleChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choice` | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-qualifications-ExamResponseText"></a>

### ExamResponseText



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `text` | [string](#string) |  | @sanitize: method=StripTags

0.5 Megabyte |






<a name="resources-qualifications-ExamResponseYesNo"></a>

### ExamResponseYesNo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [bool](#bool) |  |  |






<a name="resources-qualifications-ExamResponses"></a>

### ExamResponses
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `responses` | [ExamResponse](#resources-qualifications-ExamResponse) | repeated |  |






<a name="resources-qualifications-ExamUser"></a>

### ExamUser



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `started_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `ends_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `ended_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |





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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `abbreviation` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `access` | [QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |
| `requirements` | [QualificationRequirement](#resources-qualifications-QualificationRequirement) | repeated |  |
| `discord_sync_enabled` | [bool](#bool) |  |  |
| `discord_settings` | [QualificationDiscordSettings](#resources-qualifications-QualificationDiscordSettings) | optional |  |
| `exam_mode` | [QualificationExamMode](#resources-qualifications-QualificationExamMode) |  |  |
| `exam_settings` | [QualificationExamSettings](#resources-qualifications-QualificationExamSettings) | optional |  |
| `exam` | [ExamQuestions](#resources-qualifications-ExamQuestions) | optional |  |
| `result` | [QualificationResult](#resources-qualifications-QualificationResult) | optional |  |
| `request` | [QualificationRequest](#resources-qualifications-QualificationRequest) | optional |  |
| `label_sync_enabled` | [bool](#bool) |  |  |
| `label_sync_format` | [string](#string) | optional | @sanitize: method=StripTags |
| `files` | [resources.file.File](#resources-file-File) | repeated | @gotags: alias:"files" |






<a name="resources-qualifications-QualificationDiscordSettings"></a>

### QualificationDiscordSettings
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_name` | [string](#string) | optional |  |
| `role_format` | [string](#string) | optional |  |






<a name="resources-qualifications-QualificationExamSettings"></a>

### QualificationExamSettings
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `time` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| `auto_grade` | [bool](#bool) |  |  |
| `auto_grade_mode` | [AutoGradeMode](#resources-qualifications-AutoGradeMode) |  |  |
| `minimum_points` | [int32](#int32) |  |  |






<a name="resources-qualifications-QualificationRequest"></a>

### QualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `qualification_id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"qualification_id" |
| `qualification` | [QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| `user_id` | [int32](#int32) |  | @gotags: sql:"primary_key" |
| `user` | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:"user" |
| `user_comment` | [string](#string) | optional | @sanitize: method=StripTags |
| `status` | [RequestStatus](#resources-qualifications-RequestStatus) | optional |  |
| `approved_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `approver_comment` | [string](#string) | optional | @sanitize: method=StripTags |
| `approver_id` | [int32](#int32) | optional |  |
| `approver` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"approver" |
| `approver_job` | [string](#string) | optional |  |






<a name="resources-qualifications-QualificationRequirement"></a>

### QualificationRequirement



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `target_qualification_id` | [uint64](#uint64) |  |  |
| `target_qualification` | [QualificationShort](#resources-qualifications-QualificationShort) | optional | @gotags: alias:"targetqualification.*" |






<a name="resources-qualifications-QualificationResult"></a>

### QualificationResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:"user" |
| `status` | [ResultStatus](#resources-qualifications-ResultStatus) |  |  |
| `score` | [float](#float) | optional |  |
| `summary` | [string](#string) |  | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) |  | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |






<a name="resources-qualifications-QualificationShort"></a>

### QualificationShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `abbreviation` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `requirements` | [QualificationRequirement](#resources-qualifications-QualificationRequirement) | repeated |  |
| `exam_mode` | [QualificationExamMode](#resources-qualifications-QualificationExamMode) |  |  |
| `exam_settings` | [QualificationExamSettings](#resources-qualifications-QualificationExamSettings) | optional |  |
| `result` | [QualificationResult](#resources-qualifications-QualificationResult) | optional |  |





 <!-- end messages -->


<a name="resources-qualifications-AutoGradeMode"></a>

### AutoGradeMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| `AUTO_GRADE_MODE_UNSPECIFIED` | 0 |  |
| `AUTO_GRADE_MODE_STRICT` | 1 |  |
| `AUTO_GRADE_MODE_PARTIAL_CREDIT` | 2 |  |



<a name="resources-qualifications-QualificationExamMode"></a>

### QualificationExamMode


| Name | Number | Description |
| ---- | ------ | ----------- |
| `QUALIFICATION_EXAM_MODE_UNSPECIFIED` | 0 |  |
| `QUALIFICATION_EXAM_MODE_DISABLED` | 1 |  |
| `QUALIFICATION_EXAM_MODE_REQUEST_NEEDED` | 2 |  |
| `QUALIFICATION_EXAM_MODE_ENABLED` | 3 |  |



<a name="resources-qualifications-RequestStatus"></a>

### RequestStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| `REQUEST_STATUS_UNSPECIFIED` | 0 |  |
| `REQUEST_STATUS_PENDING` | 1 |  |
| `REQUEST_STATUS_DENIED` | 2 |  |
| `REQUEST_STATUS_ACCEPTED` | 3 |  |
| `REQUEST_STATUS_EXAM_STARTED` | 4 |  |
| `REQUEST_STATUS_EXAM_GRADING` | 5 |  |
| `REQUEST_STATUS_COMPLETED` | 6 |  |



<a name="resources-qualifications-ResultStatus"></a>

### ResultStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| `RESULT_STATUS_UNSPECIFIED` | 0 |  |
| `RESULT_STATUS_PENDING` | 1 |  |
| `RESULT_STATUS_FAILED` | 2 |  |
| `RESULT_STATUS_SUCCESSFUL` | 3 |  |


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
| `timestamp` | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 <!-- end messages -->

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
| `plate` | [string](#string) |  |  |
| `model` | [string](#string) | optional |  |
| `type` | [string](#string) |  |  |
| `owner_id` | [int32](#int32) | optional |  |
| `owner_identifier` | [string](#string) | optional |  |
| `owner` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `public` | [bool](#bool) |  |  |
| `closed` | [bool](#bool) |  |  |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `subscription` | [CalendarSub](#resources-calendar-CalendarSub) | optional |  |
| `access` | [CalendarAccess](#resources-calendar-CalendarAccess) |  |  |






<a name="resources-calendar-CalendarEntry"></a>

### CalendarEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `calendar_id` | [uint64](#uint64) |  |  |
| `calendar` | [Calendar](#resources-calendar-Calendar) | optional |  |
| `job` | [string](#string) | optional |  |
| `start_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `end_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `closed` | [bool](#bool) |  |  |
| `rsvp_open` | [bool](#bool) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `recurring` | [CalendarEntryRecurring](#resources-calendar-CalendarEntryRecurring) | optional |  |
| `rsvp` | [CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) | optional |  |






<a name="resources-calendar-CalendarEntryRSVP"></a>

### CalendarEntryRSVP



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `response` | [RsvpResponses](#resources-calendar-RsvpResponses) |  |  |






<a name="resources-calendar-CalendarEntryRecurring"></a>

### CalendarEntryRecurring
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `every` | [string](#string) |  |  |
| `count` | [int32](#int32) |  |  |
| `until` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="resources-calendar-CalendarShort"></a>

### CalendarShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `public` | [bool](#bool) |  |  |
| `closed` | [bool](#bool) |  |  |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `subscription` | [CalendarSub](#resources-calendar-CalendarSub) | optional |  |






<a name="resources-calendar-CalendarSub"></a>

### CalendarSub



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `confirmed` | [bool](#bool) |  |  |
| `muted` | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="resources-calendar-RsvpResponses"></a>

### RsvpResponses


| Name | Number | Description |
| ---- | ------ | ----------- |
| `RSVP_RESPONSES_UNSPECIFIED` | 0 |  |
| `RSVP_RESPONSES_HIDDEN` | 1 |  |
| `RSVP_RESPONSES_INVITED` | 2 |  |
| `RSVP_RESPONSES_NO` | 3 |  |
| `RSVP_RESPONSES_MAYBE` | 4 |  |
| `RSVP_RESPONSES_YES` | 5 |  |


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
| `jobs` | [CalendarJobAccess](#resources-calendar-CalendarJobAccess) | repeated | @gotags: alias:"job_access" |
| `users` | [CalendarUserAccess](#resources-calendar-CalendarUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-calendar-CalendarJobAccess"></a>

### CalendarJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-calendar-AccessLevel) |  |  |






<a name="resources-calendar-CalendarUserAccess"></a>

### CalendarUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `access` | [AccessLevel](#resources-calendar-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-calendar-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_BLOCKED` | 1 |  |
| `ACCESS_LEVEL_VIEW` | 2 |  |
| `ACCESS_LEVEL_SHARE` | 3 |  |
| `ACCESS_LEVEL_EDIT` | 4 |  |
| `ACCESS_LEVEL_MANAGE` | 5 |  |


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
| `value` | [int32](#int32) | optional |  |





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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `tld_id` | [uint64](#uint64) |  |  |
| `tld` | [TLD](#resources-internet-TLD) | optional |  |
| `active` | [bool](#bool) |  |  |
| `name` | [string](#string) |  |  |
| `transfer_code` | [string](#string) | optional |  |
| `approver_job` | [string](#string) | optional |  |
| `approver_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |






<a name="resources-internet-TLD"></a>

### TLD



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `internal` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_internet_page-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/internet/page.proto



<a name="resources-internet-ContentNode"></a>

### ContentNode



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [resources.common.content.NodeType](#resources-common-content-NodeType) |  |  |
| `id` | [string](#string) | optional | @sanitize: method=StripTags |
| `tag` | [string](#string) |  | @sanitize: method=StripTags |
| `attrs` | [ContentNode.AttrsEntry](#resources-internet-ContentNode-AttrsEntry) | repeated | @sanitize: method=StripTags |
| `text` | [string](#string) | optional | @sanitize: method=StripTags |
| `content` | [ContentNode](#resources-internet-ContentNode) | repeated |  |
| `slots` | [ContentNode](#resources-internet-ContentNode) | repeated |  |






<a name="resources-internet-ContentNode-AttrsEntry"></a>

### ContentNode.AttrsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="resources-internet-Page"></a>

### Page



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `domain_id` | [uint64](#uint64) |  |  |
| `path` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) |  | @sanitize: method=StripTags |
| `data` | [PageData](#resources-internet-PageData) |  |  |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |






<a name="resources-internet-PageData"></a>

### PageData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `layout_type` | [PageLayoutType](#resources-internet-PageLayoutType) |  |  |
| `node` | [ContentNode](#resources-internet-ContentNode) | optional |  |





 <!-- end messages -->


<a name="resources-internet-PageLayoutType"></a>

### PageLayoutType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `PAGE_LAYOUT_TYPE_UNSPECIFIED` | 0 |  |
| `PAGE_LAYOUT_TYPE_BASIC_PAGE` | 1 |  |
| `PAGE_LAYOUT_TYPE_LANDING_PAGE` | 2 |  |


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
| `id` | [uint64](#uint64) |  |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `domain_id` | [uint64](#uint64) |  |  |
| `domain` | [Domain](#resources-internet-Domain) | optional |  |
| `path` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_internet_access-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/internet/access.proto



<a name="resources-internet-PageAccess"></a>

### PageAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageJobAccess](#resources-internet-PageJobAccess) | repeated | @gotags: alias:"job_access" |
| `users` | [PageUserAccess](#resources-internet-PageUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-internet-PageJobAccess"></a>

### PageJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-internet-AccessLevel) |  |  |






<a name="resources-internet-PageUserAccess"></a>

### PageUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `access` | [AccessLevel](#resources-internet-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-internet-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_BLOCKED` | 1 |  |
| `ACCESS_LEVEL_VIEW` | 2 |  |
| `ACCESS_LEVEL_EDIT` | 3 |  |
| `ACCESS_LEVEL_OWNER` | 4 |  |


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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `disabled` | [bool](#bool) |  |  |
| `ad_type` | [AdType](#resources-internet-AdType) |  |  |
| `starts_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `ends_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) |  | @sanitize: method=StripTags |
| `image` | [resources.file.File](#resources-file-File) | optional |  |
| `approver_id` | [int32](#int32) | optional |  |
| `approver_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |





 <!-- end messages -->


<a name="resources-internet-AdType"></a>

### AdType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `AD_TYPE_UNSPECIFIED` | 0 |  |
| `AD_TYPE_SPONSORED` | 1 |  |
| `AD_TYPE_SEARCH_RESULT` | 2 |  |
| `AD_TYPE_CONTENT_MAIN` | 3 |  |
| `AD_TYPE_CONTENT_ASIDE` | 4 |  |


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
| `email_update` | [Email](#resources-mailer-Email) |  |  |
| `email_delete` | [uint64](#uint64) |  |  |
| `email_settings_updated` | [EmailSettings](#resources-mailer-EmailSettings) |  |  |
| `thread_update` | [Thread](#resources-mailer-Thread) |  |  |
| `thread_delete` | [uint64](#uint64) |  |  |
| `thread_state_update` | [ThreadState](#resources-mailer-ThreadState) |  |  |
| `message_update` | [Message](#resources-mailer-Message) |  |  |
| `message_delete` | [uint64](#uint64) |  |  |





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
| `id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |
| `sender_id` | [uint64](#uint64) |  |  |
| `sender` | [Email](#resources-mailer-Email) | optional | @gotags: alias:"sender" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  | @sanitize |
| `data` | [MessageData](#resources-mailer-MessageData) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |






<a name="resources-mailer-MessageAttachment"></a>

### MessageAttachment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [MessageAttachmentDocument](#resources-mailer-MessageAttachmentDocument) |  |  |






<a name="resources-mailer-MessageAttachmentDocument"></a>

### MessageAttachmentDocument



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `title` | [string](#string) | optional |  |






<a name="resources-mailer-MessageData"></a>

### MessageData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attachments` | [MessageAttachment](#resources-mailer-MessageAttachment) | repeated |  |





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
| `email_id` | [uint64](#uint64) |  |  |
| `signature` | [string](#string) | optional | @sanitize |
| `blocked_emails` | [string](#string) | repeated | @sanitize: method=StripTags |





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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `email_id` | [uint64](#uint64) |  |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [string](#string) |  | @sanitize |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |





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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `creator_email_id` | [uint64](#uint64) |  |  |
| `creator_email` | [Email](#resources-mailer-Email) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `recipients` | [ThreadRecipientEmail](#resources-mailer-ThreadRecipientEmail) | repeated |  |
| `state` | [ThreadState](#resources-mailer-ThreadState) | optional | @gotags: alias:"thread_state" |






<a name="resources-mailer-ThreadRecipientEmail"></a>

### ThreadRecipientEmail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  | @gotags: alias:"thread_id" |
| `email_id` | [uint64](#uint64) |  |  |
| `email` | [Email](#resources-mailer-Email) | optional |  |






<a name="resources-mailer-ThreadState"></a>

### ThreadState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread_id` | [uint64](#uint64) |  |  |
| `email_id` | [uint64](#uint64) |  |  |
| `last_read` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `unread` | [bool](#bool) | optional |  |
| `important` | [bool](#bool) | optional |  |
| `favorite` | [bool](#bool) | optional |  |
| `muted` | [bool](#bool) | optional |  |
| `archived` | [bool](#bool) | optional |  |





 <!-- end messages -->

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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deactivated` | [bool](#bool) |  |  |
| `job` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `email` | [string](#string) |  | @sanitize: method=StripTags |
| `email_changed` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `label` | [string](#string) | optional | @sanitize: method=StripTags |
| `access` | [Access](#resources-mailer-Access) |  |  |
| `settings` | [EmailSettings](#resources-mailer-EmailSettings) | optional |  |





 <!-- end messages -->

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
| `jobs` | [JobAccess](#resources-mailer-JobAccess) | repeated | @gotags: alias:"job_access" |
| `users` | [UserAccess](#resources-mailer-UserAccess) | repeated | @gotags: alias:"user_access" |
| `qualifications` | [QualificationAccess](#resources-mailer-QualificationAccess) | repeated | @gotags: alias:"qualification_access" |






<a name="resources-mailer-JobAccess"></a>

### JobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-mailer-AccessLevel) |  |  |






<a name="resources-mailer-QualificationAccess"></a>

### QualificationAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [resources.qualifications.QualificationShort](#resources-qualifications-QualificationShort) | optional |  |
| `access` | [AccessLevel](#resources-mailer-AccessLevel) |  |  |






<a name="resources-mailer-UserAccess"></a>

### UserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `access` | [AccessLevel](#resources-mailer-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-mailer-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_BLOCKED` | 1 |  |
| `ACCESS_LEVEL_READ` | 2 |  |
| `ACCESS_LEVEL_WRITE` | 3 |  |
| `ACCESS_LEVEL_MANAGE` | 4 |  |


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
| `jobs` | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated | @gotags: alias:"job_access" |
| `users` | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated | @gotags: alias:"user_access" |






<a name="resources-wiki-PageJobAccess"></a>

### PageJobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resources-wiki-AccessLevel) |  |  |






<a name="resources-wiki-PageUserAccess"></a>

### PageUserAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `access` | [AccessLevel](#resources-wiki-AccessLevel) |  |  |





 <!-- end messages -->


<a name="resources-wiki-AccessLevel"></a>

### AccessLevel


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `ACCESS_LEVEL_BLOCKED` | 1 |  |
| `ACCESS_LEVEL_VIEW` | 2 |  |
| `ACCESS_LEVEL_ACCESS` | 3 |  |
| `ACCESS_LEVEL_EDIT` | 4 |  |


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
| `to_create` | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated |  |
| `to_update` | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated |  |
| `to_delete` | [PageJobAccess](#resources-wiki-PageJobAccess) | repeated |  |






<a name="resources-wiki-PageAccessUpdated"></a>

### PageAccessUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageAccessJobsDiff](#resources-wiki-PageAccessJobsDiff) |  |  |
| `users` | [PageAccessUsersDiff](#resources-wiki-PageAccessUsersDiff) |  |  |






<a name="resources-wiki-PageAccessUsersDiff"></a>

### PageAccessUsersDiff



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated |  |
| `to_update` | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated |  |
| `to_delete` | [PageUserAccess](#resources-wiki-PageUserAccess) | repeated |  |






<a name="resources-wiki-PageActivity"></a>

### PageActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `page_id` | [uint64](#uint64) |  |  |
| `activity_type` | [PageActivityType](#resources-wiki-PageActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [PageActivityData](#resources-wiki-PageActivityData) |  |  |






<a name="resources-wiki-PageActivityData"></a>

### PageActivityData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [PageUpdated](#resources-wiki-PageUpdated) |  |  |
| `access_updated` | [PageAccessUpdated](#resources-wiki-PageAccessUpdated) |  |  |






<a name="resources-wiki-PageFilesChange"></a>

### PageFilesChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [int64](#int64) |  |  |
| `deleted` | [int64](#int64) |  |  |






<a name="resources-wiki-PageUpdated"></a>

### PageUpdated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title_diff` | [string](#string) | optional |  |
| `description_diff` | [string](#string) | optional |  |
| `content_diff` | [string](#string) | optional |  |
| `files_change` | [PageFilesChange](#resources-wiki-PageFilesChange) | optional |  |





 <!-- end messages -->


<a name="resources-wiki-PageActivityType"></a>

### PageActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `PAGE_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `PAGE_ACTIVITY_TYPE_CREATED` | 1 | Base |
| `PAGE_ACTIVITY_TYPE_UPDATED` | 2 |  |
| `PAGE_ACTIVITY_TYPE_ACCESS_UPDATED` | 3 |  |
| `PAGE_ACTIVITY_TYPE_OWNER_CHANGED` | 4 |  |
| `PAGE_ACTIVITY_TYPE_DELETED` | 5 |  |
| `PAGE_ACTIVITY_TYPE_DRAFT_TOGGLED` | 6 |  |


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
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `job` | [string](#string) |  | @sanitize: method=StripTags |
| `job_label` | [string](#string) | optional |  |
| `parent_id` | [uint64](#uint64) | optional |  |
| `meta` | [PageMeta](#resources-wiki-PageMeta) |  |  |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `access` | [PageAccess](#resources-wiki-PageAccess) |  |  |
| `files` | [resources.file.File](#resources-file-File) | repeated | @gotags: alias:"files" |






<a name="resources-wiki-PageMeta"></a>

### PageMeta



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `slug` | [string](#string) | optional | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) |  | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional | @gotags: alias:"creator" |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |
| `tags` | [string](#string) | repeated | @sanitize: method=StripTags |
| `toc` | [bool](#bool) | optional |  |
| `public` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |






<a name="resources-wiki-PageRootInfo"></a>

### PageRootInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `logo_file_id` | [uint64](#uint64) | optional |  |
| `logo` | [resources.file.File](#resources-file-File) | optional | @gotags: alias:"logo" |






<a name="resources-wiki-PageShort"></a>

### PageShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `parent_id` | [uint64](#uint64) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `slug` | [string](#string) | optional | @sanitize: method=StripTags |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `children` | [PageShort](#resources-wiki-PageShort) | repeated |  |
| `root_info` | [PageRootInfo](#resources-wiki-PageRootInfo) | optional |  |
| `level` | [int32](#int32) | optional |  |
| `draft` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_sync_activity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/sync/activity.proto



<a name="resources-sync-ColleagueProps"></a>

### ColleagueProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reason` | [string](#string) | optional |  |
| `props` | [resources.jobs.ColleagueProps](#resources-jobs-ColleagueProps) |  |  |






<a name="resources-sync-TimeclockUpdate"></a>

### TimeclockUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `start` | [bool](#bool) |  |  |






<a name="resources-sync-UserOAuth2Conn"></a>

### UserOAuth2Conn
Connect an identifier/license to the provider with the specified external id (e.g., auto discord social connect on server join)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider_name` | [string](#string) |  |  |
| `identifier` | [string](#string) |  |  |
| `external_id` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |






<a name="resources-sync-UserProps"></a>

### UserProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reason` | [string](#string) | optional |  |
| `props` | [resources.users.UserProps](#resources-users-UserProps) |  |  |






<a name="resources-sync-UserUpdate"></a>

### UserUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `group` | [string](#string) | optional |  |
| `job` | [string](#string) | optional | Char details |
| `job_grade` | [int32](#int32) | optional |  |
| `firstname` | [string](#string) | optional |  |
| `lastname` | [string](#string) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_sync_data-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/sync/data.proto



<a name="resources-sync-CitizenLocations"></a>

### CitizenLocations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `job` | [string](#string) |  |  |
| `coords` | [resources.livemap.Coords](#resources-livemap-Coords) |  |  |
| `hidden` | [bool](#bool) |  |  |
| `remove` | [bool](#bool) |  |  |






<a name="resources-sync-DataJobs"></a>

### DataJobs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.jobs.Job](#resources-jobs-Job) | repeated |  |






<a name="resources-sync-DataLicenses"></a>

### DataLicenses



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `licenses` | [resources.users.License](#resources-users-License) | repeated |  |






<a name="resources-sync-DataStatus"></a>

### DataStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `count` | [int64](#int64) |  |  |






<a name="resources-sync-DataUserLocations"></a>

### DataUserLocations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [CitizenLocations](#resources-sync-CitizenLocations) | repeated |  |
| `clear_all` | [bool](#bool) | optional |  |






<a name="resources-sync-DataUsers"></a>

### DataUsers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.users.User](#resources-users-User) | repeated |  |






<a name="resources-sync-DataVehicles"></a>

### DataVehicles



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vehicles` | [resources.vehicles.Vehicle](#resources-vehicles-Vehicle) | repeated |  |






<a name="resources-sync-DeleteUsers"></a>

### DeleteUsers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_ids` | [int32](#int32) | repeated |  |






<a name="resources-sync-DeleteVehicles"></a>

### DeleteVehicles



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plates` | [string](#string) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_activity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/activity.proto



<a name="resources-users-CitizenDocumentRelation"></a>

### CitizenDocumentRelation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [bool](#bool) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `relation` | [int32](#int32) |  | resources.documents.DocRelation enum |






<a name="resources-users-FineChange"></a>

### FineChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `removed` | [bool](#bool) |  |  |
| `amount` | [int64](#int64) |  |  |






<a name="resources-users-JailChange"></a>

### JailChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seconds` | [int32](#int32) |  |  |
| `admin` | [bool](#bool) |  |  |
| `location` | [string](#string) | optional |  |






<a name="resources-users-JobChange"></a>

### JobChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) | optional |  |
| `grade_label` | [string](#string) | optional |  |






<a name="resources-users-LabelsChange"></a>

### LabelsChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [Label](#resources-users-Label) | repeated |  |
| `removed` | [Label](#resources-users-Label) | repeated |  |






<a name="resources-users-LicenseChange"></a>

### LicenseChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [bool](#bool) |  |  |
| `licenses` | [License](#resources-users-License) | repeated |  |






<a name="resources-users-MugshotChange"></a>

### MugshotChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new` | [string](#string) | optional |  |






<a name="resources-users-NameChange"></a>

### NameChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |






<a name="resources-users-TrafficInfractionPointsChange"></a>

### TrafficInfractionPointsChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old` | [uint32](#uint32) |  |  |
| `new` | [uint32](#uint32) |  |  |






<a name="resources-users-UserActivity"></a>

### UserActivity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: alias:"user_activity.id" |
| `type` | [UserActivityType](#resources-users-UserActivityType) |  | @gotags: alias:"user_activity.type" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: alias:"user_activity.created_at" |
| `source_user_id` | [int32](#int32) | optional | @gotags: alias:"source_user_id" |
| `source_user` | [UserShort](#resources-users-UserShort) | optional | @gotags: alias:"source_user" |
| `target_user_id` | [int32](#int32) |  | @gotags: alias:"target_user_id" |
| `target_user` | [UserShort](#resources-users-UserShort) |  | @gotags: alias:"target_user" |
| `key` | [string](#string) |  | @sanitize

@gotags: alias:"user_activity.key" |
| `reason` | [string](#string) |  | @sanitize

@gotags: alias:"user_activity.reason" |
| `data` | [UserActivityData](#resources-users-UserActivityData) | optional | @gotags: alias:"user_activity.data" |
| `old_value` | [string](#string) |  | @gotags: alias:"user_activity.old_value" |
| `new_value` | [string](#string) |  | @gotags: alias:"user_activity.new_value" |






<a name="resources-users-UserActivityData"></a>

### UserActivityData
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name_change` | [NameChange](#resources-users-NameChange) |  |  |
| `licenses_change` | [LicenseChange](#resources-users-LicenseChange) |  |  |
| `wanted_change` | [WantedChange](#resources-users-WantedChange) |  | User Props |
| `traffic_infraction_points_change` | [TrafficInfractionPointsChange](#resources-users-TrafficInfractionPointsChange) |  |  |
| `mugshot_change` | [MugshotChange](#resources-users-MugshotChange) |  |  |
| `labels_change` | [LabelsChange](#resources-users-LabelsChange) |  |  |
| `job_change` | [JobChange](#resources-users-JobChange) |  |  |
| `document_relation` | [CitizenDocumentRelation](#resources-users-CitizenDocumentRelation) |  | Docstore related |
| `jail_change` | [JailChange](#resources-users-JailChange) |  | "Plugin" activities |
| `fine_change` | [FineChange](#resources-users-FineChange) |  |  |






<a name="resources-users-WantedChange"></a>

### WantedChange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `wanted` | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="resources-users-UserActivityType"></a>

### UserActivityType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `USER_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `USER_ACTIVITY_TYPE_NAME` | 4 |  |
| `USER_ACTIVITY_TYPE_LICENSES` | 5 |  |
| `USER_ACTIVITY_TYPE_WANTED` | 6 |  |
| `USER_ACTIVITY_TYPE_TRAFFIC_INFRACTION_POINTS` | 7 |  |
| `USER_ACTIVITY_TYPE_MUGSHOT` | 8 |  |
| `USER_ACTIVITY_TYPE_LABELS` | 9 |  |
| `USER_ACTIVITY_TYPE_JOB` | 10 |  |
| `USER_ACTIVITY_TYPE_DOCUMENT` | 11 |  |
| `USER_ACTIVITY_TYPE_JAIL` | 12 |  |
| `USER_ACTIVITY_TYPE_FINE` | 13 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_licenses-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/licenses.proto



<a name="resources-users-CitizensLicenses"></a>

### CitizensLicenses



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `licenses` | [License](#resources-users-License) | repeated |  |






<a name="resources-users-License"></a>

### License



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `label` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_labels-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/labels.proto



<a name="resources-users-Label"></a>

### Label



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: sql:"primary_key" alias:"id" |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `color` | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-users-Labels"></a>

### Labels
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [Label](#resources-users-Label) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_props-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/props.proto



<a name="resources-users-UserProps"></a>

### UserProps



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `wanted` | [bool](#bool) | optional |  |
| `job_name` | [string](#string) | optional | @gotags: alias:"job" |
| `job` | [resources.jobs.Job](#resources-jobs-Job) | optional |  |
| `job_grade_number` | [int32](#int32) | optional | @gotags: alias:"job_grade" |
| `job_grade` | [resources.jobs.JobGrade](#resources-jobs-JobGrade) | optional |  |
| `traffic_infraction_points` | [uint32](#uint32) | optional |  |
| `traffic_infraction_points_updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `open_fines` | [int64](#int64) | optional |  |
| `blood_type` | [string](#string) | optional |  |
| `mugshot_file_id` | [uint64](#uint64) | optional |  |
| `mugshot` | [resources.file.File](#resources-file-File) | optional | @gotags: alias:"mugshot" |
| `labels` | [Labels](#resources-users-Labels) | optional |  |
| `email` | [string](#string) | optional | @sanitize: method=StripTags |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_users_users-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/users/users.proto



<a name="resources-users-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  | @gotags: alias:"id" |
| `identifier` | [string](#string) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `sex` | [string](#string) | optional |  |
| `height` | [string](#string) | optional |  |
| `phone_number` | [string](#string) | optional |  |
| `visum` | [int32](#int32) | optional |  |
| `playtime` | [int32](#int32) | optional |  |
| `props` | [UserProps](#resources-users-UserProps) |  | @gotags: alias:"fivenet_user_props" |
| `licenses` | [License](#resources-users-License) | repeated | @gotags: alias:"user_licenses" |
| `avatar_file_id` | [uint64](#uint64) | optional |  |
| `avatar` | [string](#string) | optional |  |
| `group` | [string](#string) | optional |  |






<a name="resources-users-UserShort"></a>

### UserShort



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  | @gotags: alias:"id" |
| `identifier` | [string](#string) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `phone_number` | [string](#string) | optional |  |
| `avatar_file_id` | [uint64](#uint64) | optional |  |
| `avatar` | [string](#string) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_audit_audit-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/audit/audit.proto



<a name="resources-audit-AuditEntry"></a>

### AuditEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `user_id` | [int32](#int32) |  | @gotags: alias:"user_id" |
| `user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `user_job` | [string](#string) |  | @gotags: alias:"user_job" |
| `target_user_id` | [int32](#int32) | optional | @gotags: alias:"target_user_id" |
| `target_user` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |
| `target_user_job` | [string](#string) | optional | @gotags: alias:"target_user_job" |
| `service` | [string](#string) |  | @gotags: alias:"service" |
| `method` | [string](#string) |  | @gotags: alias:"method" |
| `state` | [EventType](#resources-audit-EventType) |  | @gotags: alias:"state" |
| `data` | [string](#string) | optional | @gotags: alias:"data" |





 <!-- end messages -->


<a name="resources-audit-EventType"></a>

### EventType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `EVENT_TYPE_UNSPECIFIED` | 0 |  |
| `EVENT_TYPE_ERRORED` | 1 |  |
| `EVENT_TYPE_VIEWED` | 2 |  |
| `EVENT_TYPE_CREATED` | 3 |  |
| `EVENT_TYPE_UPDATED` | 4 |  |
| `EVENT_TYPE_DELETED` | 5 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_settings_banner-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/settings/banner.proto



<a name="resources-settings-BannerMessage"></a>

### BannerMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize: method |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `expires_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_settings_config-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/settings/config.proto



<a name="resources-settings-AppConfig"></a>

### AppConfig
@dbscanner: json,partial


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) | optional |  |
| `default_locale` | [string](#string) |  |  |
| `auth` | [Auth](#resources-settings-Auth) |  |  |
| `perms` | [Perms](#resources-settings-Perms) |  |  |
| `website` | [Website](#resources-settings-Website) |  |  |
| `job_info` | [JobInfo](#resources-settings-JobInfo) |  |  |
| `user_tracker` | [UserTracker](#resources-settings-UserTracker) |  |  |
| `discord` | [Discord](#resources-settings-Discord) |  |  |
| `system` | [System](#resources-settings-System) |  |  |






<a name="resources-settings-Auth"></a>

### Auth



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signup_enabled` | [bool](#bool) |  |  |
| `last_char_lock` | [bool](#bool) |  |  |






<a name="resources-settings-Discord"></a>

### Discord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `sync_interval` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| `invite_url` | [string](#string) | optional | @sanitize: method=StripTags |
| `ignored_jobs` | [string](#string) | repeated | @sanitize: method=StripTags |
| `bot_presence` | [DiscordBotPresence](#resources-settings-DiscordBotPresence) | optional |  |
| `bot_id` | [string](#string) | optional | @sanitize: method=StripTags |
| `bot_permissions` | [int64](#int64) |  |  |






<a name="resources-settings-DiscordBotPresence"></a>

### DiscordBotPresence



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [DiscordBotPresenceType](#resources-settings-DiscordBotPresenceType) |  |  |
| `status` | [string](#string) | optional | @sanitize: method=StripTags |
| `url` | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-settings-JobInfo"></a>

### JobInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unemployed_job` | [UnemployedJob](#resources-settings-UnemployedJob) |  |  |
| `public_jobs` | [string](#string) | repeated | @sanitize: method=StripTags |
| `hidden_jobs` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="resources-settings-Links"></a>

### Links



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `privacy_policy` | [string](#string) | optional | @sanitize: method=StripTags |
| `imprint` | [string](#string) | optional | @sanitize: method=StripTags |






<a name="resources-settings-Perm"></a>

### Perm



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [string](#string) |  | @sanitize: method=StripTags |
| `name` | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-settings-Perms"></a>

### Perms



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default` | [Perm](#resources-settings-Perm) | repeated |  |






<a name="resources-settings-System"></a>

### System



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `banner_message_enabled` | [bool](#bool) |  |  |
| `banner_message` | [BannerMessage](#resources-settings-BannerMessage) |  |  |






<a name="resources-settings-UnemployedJob"></a>

### UnemployedJob



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |






<a name="resources-settings-UserTracker"></a>

### UserTracker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_time` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |
| `db_refresh_time` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |






<a name="resources-settings-Website"></a>

### Website



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `links` | [Links](#resources-settings-Links) |  |  |
| `stats_page` | [bool](#bool) |  |  |





 <!-- end messages -->


<a name="resources-settings-DiscordBotPresenceType"></a>

### DiscordBotPresenceType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISCORD_BOT_PRESENCE_TYPE_UNSPECIFIED` | 0 |  |
| `DISCORD_BOT_PRESENCE_TYPE_GAME` | 1 |  |
| `DISCORD_BOT_PRESENCE_TYPE_LISTENING` | 2 |  |
| `DISCORD_BOT_PRESENCE_TYPE_STREAMING` | 3 |  |
| `DISCORD_BOT_PRESENCE_TYPE_WATCH` | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_collab_collab-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/collab/collab.proto



<a name="resources-collab-AwarenessPing"></a>

### AwarenessPing



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |






<a name="resources-collab-ClientPacket"></a>

### ClientPacket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hello` | [CollabInit](#resources-collab-CollabInit) |  | Must be the first message |
| `sync_step` | [SyncStep](#resources-collab-SyncStep) |  |  |
| `yjs_update` | [YjsUpdate](#resources-collab-YjsUpdate) |  |  |
| `awareness` | [AwarenessPing](#resources-collab-AwarenessPing) |  |  |






<a name="resources-collab-CollabHandshake"></a>

### CollabHandshake



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [uint64](#uint64) |  |  |
| `first` | [bool](#bool) |  |  |






<a name="resources-collab-CollabInit"></a>

### CollabInit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_id` | [uint64](#uint64) |  |  |






<a name="resources-collab-ServerPacket"></a>

### ServerPacket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_id` | [uint64](#uint64) |  | Who generated this packet (same ID used in awareness) |
| `handshake` | [CollabHandshake](#resources-collab-CollabHandshake) |  |  |
| `sync_step` | [SyncStep](#resources-collab-SyncStep) |  |  |
| `yjs_update` | [YjsUpdate](#resources-collab-YjsUpdate) |  |  |
| `awareness` | [AwarenessPing](#resources-collab-AwarenessPing) |  |  |
| `target_saved` | [TargetSaved](#resources-collab-TargetSaved) |  |  |






<a name="resources-collab-SyncStep"></a>

### SyncStep



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `step` | [int32](#int32) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `receiver_id` | [uint64](#uint64) | optional |  |






<a name="resources-collab-TargetSaved"></a>

### TargetSaved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_id` | [uint64](#uint64) |  |  |






<a name="resources-collab-YjsUpdate"></a>

### YjsUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |





 <!-- end messages -->


<a name="resources-collab-ClientRole"></a>

### ClientRole


| Name | Number | Description |
| ---- | ------ | ----------- |
| `CLIENT_ROLE_UNSPECIFIED` | 0 |  |
| `CLIENT_ROLE_READER` | 1 |  |
| `CLIENT_ROLE_WRITER` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_discord_discord-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/discord/discord.proto



<a name="resources-discord-Channel"></a>

### Channel



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `guild_id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `type` | [uint32](#uint32) |  |  |
| `position` | [int64](#int64) |  |  |






<a name="resources-discord-Guild"></a>

### Guild



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `icon` | [string](#string) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_file_file-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/file/file.proto



<a name="resources-file-File"></a>

### File



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) | optional |  |
| `id` | [uint64](#uint64) |  |  |
| `file_path` | [string](#string) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `byte_size` | [int64](#int64) |  | Bytes stored |
| `content_type` | [string](#string) |  |  |
| `meta` | [FileMeta](#resources-file-FileMeta) | optional |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_file_filestore-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/file/filestore.proto



<a name="resources-file-DeleteFileRequest"></a>

### DeleteFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) |  |  |
| `file_id` | [uint64](#uint64) |  |  |






<a name="resources-file-DeleteFileResponse"></a>

### DeleteFileResponse







<a name="resources-file-UploadMeta"></a>

### UploadMeta



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) |  |  |
| `namespace` | [string](#string) |  | "documents", "wiki", … |
| `original_name` | [string](#string) |  |  |
| `content_type` | [string](#string) |  | optional – server re-validates |
| `size` | [int64](#int64) |  | Size in bytes |
| `reason` | [string](#string) |  | @sanitize |






<a name="resources-file-UploadPacket"></a>

### UploadPacket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `meta` | [UploadMeta](#resources-file-UploadMeta) |  |  |
| `data` | [bytes](#bytes) |  | Raw bytes <= 128 KiB each, browsers should only read 64 KiB at a time, but this is a buffer just in case |






<a name="resources-file-UploadResponse"></a>

### UploadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | Unique ID for the uploaded file |
| `url` | [string](#string) |  | URL to the uploaded file |
| `file` | [File](#resources-file-File) |  | File info |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_file_meta-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/file/meta.proto



<a name="resources-file-FileMeta"></a>

### FileMeta
@dbscanner: json


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image` | [ImageMeta](#resources-file-ImageMeta) |  |  |






<a name="resources-file-ImageMeta"></a>

### ImageMeta



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `width` | [int64](#int64) |  |  |
| `height` | [int64](#int64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_livemap_coords-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/coords.proto



<a name="resources-livemap-Coords"></a>

### Coords



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_livemap_heatmap-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/heatmap.proto



<a name="resources-livemap-HeatmapEntry"></a>

### HeatmapEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `w` | [double](#double) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_livemap_marker_marker-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/marker_marker.proto



<a name="resources-livemap-CircleMarker"></a>

### CircleMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `radius` | [int32](#int32) |  |  |
| `opacity` | [float](#float) | optional |  |






<a name="resources-livemap-IconMarker"></a>

### IconMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `icon` | [string](#string) |  | @sanitize: method=StripTags |






<a name="resources-livemap-MarkerData"></a>

### MarkerData
@dbscanner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `circle` | [CircleMarker](#resources-livemap-CircleMarker) |  |  |
| `icon` | [IconMarker](#resources-livemap-IconMarker) |  |  |






<a name="resources-livemap-MarkerMarker"></a>

### MarkerMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `expires_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `postal` | [string](#string) | optional | @sanitize: method=StripTags |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) |  |  |
| `type` | [MarkerType](#resources-livemap-MarkerType) |  | @gotags: alias:"markerType" |
| `data` | [MarkerData](#resources-livemap-MarkerData) |  | @gotags: alias:"markerData" |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resources-users-UserShort) | optional |  |





 <!-- end messages -->


<a name="resources-livemap-MarkerType"></a>

### MarkerType


| Name | Number | Description |
| ---- | ------ | ----------- |
| `MARKER_TYPE_UNSPECIFIED` | 0 |  |
| `MARKER_TYPE_DOT` | 1 |  |
| `MARKER_TYPE_CIRCLE` | 2 |  |
| `MARKER_TYPE_ICON` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_livemap_user_marker-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/livemap/user_marker.proto



<a name="resources-livemap-UserMarker"></a>

### UserMarker



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `postal` | [string](#string) | optional | @sanitize: method=StripTags |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) |  |  |
| `job_grade` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.Colleague](#resources-jobs-Colleague) |  | @gotags: alias:"user" |
| `unit_id` | [uint64](#uint64) | optional |  |
| `unit` | [resources.centrum.Unit](#resources-centrum-Unit) | optional |  |
| `hidden` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_tracker_mapping-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/tracker/mapping.proto



<a name="resources-tracker-UserMapping"></a>

### UserMapping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `unit_id` | [uint64](#uint64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `hidden` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_clientconfig_clientconfig-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/clientconfig/clientconfig.proto



<a name="resources-clientconfig-ClientConfig"></a>

### ClientConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `default_locale` | [string](#string) |  |  |
| `login` | [LoginConfig](#resources-clientconfig-LoginConfig) |  |  |
| `discord` | [Discord](#resources-clientconfig-Discord) |  |  |
| `website` | [Website](#resources-clientconfig-Website) |  |  |
| `feature_gates` | [FeatureGates](#resources-clientconfig-FeatureGates) |  |  |
| `game` | [Game](#resources-clientconfig-Game) |  |  |
| `system` | [System](#resources-clientconfig-System) |  |  |






<a name="resources-clientconfig-Discord"></a>

### Discord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bot_enabled` | [bool](#bool) |  |  |






<a name="resources-clientconfig-FeatureGates"></a>

### FeatureGates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image_proxy` | [bool](#bool) |  |  |






<a name="resources-clientconfig-Game"></a>

### Game



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unemployed_job_name` | [string](#string) |  |  |
| `start_job_grade` | [int32](#int32) |  |  |






<a name="resources-clientconfig-Links"></a>

### Links



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `imprint` | [string](#string) | optional |  |
| `privacy_policy` | [string](#string) | optional |  |






<a name="resources-clientconfig-LoginConfig"></a>

### LoginConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signup_enabled` | [bool](#bool) |  |  |
| `last_char_lock` | [bool](#bool) |  |  |
| `providers` | [ProviderConfig](#resources-clientconfig-ProviderConfig) | repeated |  |






<a name="resources-clientconfig-OTLPFrontend"></a>

### OTLPFrontend



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `url` | [string](#string) |  |  |
| `headers` | [OTLPFrontend.HeadersEntry](#resources-clientconfig-OTLPFrontend-HeadersEntry) | repeated |  |






<a name="resources-clientconfig-OTLPFrontend-HeadersEntry"></a>

### OTLPFrontend.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="resources-clientconfig-ProviderConfig"></a>

### ProviderConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `label` | [string](#string) |  |  |
| `icon` | [string](#string) | optional |  |
| `homepage` | [string](#string) |  |  |






<a name="resources-clientconfig-System"></a>

### System



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `banner_message_enabled` | [bool](#bool) |  |  |
| `banner_message` | [resources.settings.BannerMessage](#resources-settings-BannerMessage) | optional |  |
| `otlp` | [OTLPFrontend](#resources-clientconfig-OTLPFrontend) |  |  |






<a name="resources-clientconfig-Website"></a>

### Website



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `links` | [Links](#resources-clientconfig-Links) |  |  |
| `stats_page` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="resources_userinfo_user_info-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## resources/userinfo/user_info.proto



<a name="resources-userinfo-PollReq"></a>

### PollReq
PollReq: published to `userinfo.poll.request` when an active user connects or requests a refresh.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  | The account the user belongs to |
| `user_id` | [int32](#int32) |  | The unique user identifier within the account |






<a name="resources-userinfo-UserInfo"></a>

### UserInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `account_id` | [uint64](#uint64) |  |  |
| `license` | [string](#string) |  |  |
| `last_char` | [int32](#int32) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `job_grade` | [int32](#int32) |  |  |
| `group` | [string](#string) |  |  |
| `can_be_superuser` | [bool](#bool) |  |  |
| `superuser` | [bool](#bool) |  |  |
| `override_job` | [string](#string) | optional |  |
| `override_job_grade` | [int32](#int32) | optional |  |






<a name="resources-userinfo-UserInfoChanged"></a>

### UserInfoChanged
UserInfoChanged used to signal Job or JobGrade changes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  | The account the user belongs to |
| `user_id` | [int32](#int32) |  | The unique user identifier within the account |
| `old_job` | [string](#string) |  | Previous job title |
| `new_job` | [string](#string) |  | New job title |
| `new_job_label` | [string](#string) | optional |  |
| `old_job_grade` | [int32](#int32) |  | Previous job grade |
| `new_job_grade` | [int32](#int32) |  | New job grade |
| `new_job_grade_label` | [string](#string) | optional | New job grade label |
| `changed_at` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | Timestamp of when the change was detected |





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
| `current` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |






<a name="services-auth-ChangePasswordResponse"></a>

### ChangePasswordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `expires` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |






<a name="services-auth-ChangeUsernameRequest"></a>

### ChangeUsernameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |






<a name="services-auth-ChangeUsernameResponse"></a>

### ChangeUsernameResponse







<a name="services-auth-ChooseCharacterRequest"></a>

### ChooseCharacterRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `char_id` | [int32](#int32) |  |  |






<a name="services-auth-ChooseCharacterResponse"></a>

### ChooseCharacterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `username` | [string](#string) |  |  |
| `expires` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `job_props` | [resources.jobs.JobProps](#resources-jobs-JobProps) |  |  |
| `char` | [resources.users.User](#resources-users-User) |  | @gotags: alias:"user" |
| `permissions` | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-auth-CreateAccountRequest"></a>

### CreateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `password` | [string](#string) |  |  |






<a name="services-auth-CreateAccountResponse"></a>

### CreateAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  |  |






<a name="services-auth-DeleteOAuth2ConnectionRequest"></a>

### DeleteOAuth2ConnectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider` | [string](#string) |  |  |






<a name="services-auth-DeleteOAuth2ConnectionResponse"></a>

### DeleteOAuth2ConnectionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `success` | [bool](#bool) |  |  |






<a name="services-auth-ForgotPasswordRequest"></a>

### ForgotPasswordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |






<a name="services-auth-ForgotPasswordResponse"></a>

### ForgotPasswordResponse







<a name="services-auth-GetAccountInfoRequest"></a>

### GetAccountInfoRequest







<a name="services-auth-GetAccountInfoResponse"></a>

### GetAccountInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [resources.accounts.Account](#resources-accounts-Account) |  |  |
| `oauth2_providers` | [resources.accounts.OAuth2Provider](#resources-accounts-OAuth2Provider) | repeated |  |
| `oauth2_connections` | [resources.accounts.OAuth2Account](#resources-accounts-OAuth2Account) | repeated |  |






<a name="services-auth-GetCharactersRequest"></a>

### GetCharactersRequest







<a name="services-auth-GetCharactersResponse"></a>

### GetCharactersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chars` | [resources.accounts.Character](#resources-accounts-Character) | repeated |  |






<a name="services-auth-LoginRequest"></a>

### LoginRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `username` | [string](#string) |  |  |
| `password` | [string](#string) |  |  |






<a name="services-auth-LoginResponse"></a>

### LoginResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `expires` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `account_id` | [uint64](#uint64) |  |  |
| `char` | [ChooseCharacterResponse](#services-auth-ChooseCharacterResponse) | optional |  |






<a name="services-auth-LogoutRequest"></a>

### LogoutRequest







<a name="services-auth-LogoutResponse"></a>

### LogoutResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `success` | [bool](#bool) |  |  |






<a name="services-auth-SetSuperuserModeRequest"></a>

### SetSuperuserModeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `superuser` | [bool](#bool) |  |  |
| `job` | [string](#string) | optional |  |






<a name="services-auth-SetSuperuserModeResponse"></a>

### SetSuperuserModeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `expires` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `job_props` | [resources.jobs.JobProps](#resources-jobs-JobProps) | optional |  |
| `char` | [resources.users.User](#resources-users-User) |  | @gotags: alias:"user" |
| `permissions` | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-auth-AuthService"></a>

### AuthService
Auth Service handles user authentication, character selection and oauth2 connections Some methods **must** be caled via HTTP-based GRPC web request to allow cookies to be set/unset.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Login` | [LoginRequest](#services-auth-LoginRequest) | [LoginResponse](#services-auth-LoginResponse) |  |
| `Logout` | [LogoutRequest](#services-auth-LogoutRequest) | [LogoutResponse](#services-auth-LogoutResponse) |  |
| `CreateAccount` | [CreateAccountRequest](#services-auth-CreateAccountRequest) | [CreateAccountResponse](#services-auth-CreateAccountResponse) |  |
| `ChangeUsername` | [ChangeUsernameRequest](#services-auth-ChangeUsernameRequest) | [ChangeUsernameResponse](#services-auth-ChangeUsernameResponse) |  |
| `ChangePassword` | [ChangePasswordRequest](#services-auth-ChangePasswordRequest) | [ChangePasswordResponse](#services-auth-ChangePasswordResponse) |  |
| `ForgotPassword` | [ForgotPasswordRequest](#services-auth-ForgotPasswordRequest) | [ForgotPasswordResponse](#services-auth-ForgotPasswordResponse) |  |
| `GetCharacters` | [GetCharactersRequest](#services-auth-GetCharactersRequest) | [GetCharactersResponse](#services-auth-GetCharactersResponse) |  |
| `ChooseCharacter` | [ChooseCharacterRequest](#services-auth-ChooseCharacterRequest) | [ChooseCharacterResponse](#services-auth-ChooseCharacterResponse) | @perm |
| `GetAccountInfo` | [GetAccountInfoRequest](#services-auth-GetAccountInfoRequest) | [GetAccountInfoResponse](#services-auth-GetAccountInfoResponse) |  |
| `DeleteOAuth2Connection` | [DeleteOAuth2ConnectionRequest](#services-auth-DeleteOAuth2ConnectionRequest) | [DeleteOAuth2ConnectionResponse](#services-auth-DeleteOAuth2ConnectionResponse) |  |
| `SetSuperuserMode` | [SetSuperuserModeRequest](#services-auth-SetSuperuserModeRequest) | [SetSuperuserModeResponse](#services-auth-SetSuperuserModeResponse) |  |

 <!-- end services -->



<a name="services_centrum_centrum-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/centrum/centrum.proto



<a name="services-centrum-AssignDispatchRequest"></a>

### AssignDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `to_add` | [uint64](#uint64) | repeated |  |
| `to_remove` | [uint64](#uint64) | repeated |  |
| `forced` | [bool](#bool) | optional |  |






<a name="services-centrum-AssignDispatchResponse"></a>

### AssignDispatchResponse







<a name="services-centrum-AssignUnitRequest"></a>

### AssignUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `to_add` | [int32](#int32) | repeated |  |
| `to_remove` | [int32](#int32) | repeated |  |






<a name="services-centrum-AssignUnitResponse"></a>

### AssignUnitResponse







<a name="services-centrum-CreateDispatchRequest"></a>

### CreateDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-CreateDispatchResponse"></a>

### CreateDispatchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-CreateOrUpdateUnitRequest"></a>

### CreateOrUpdateUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |






<a name="services-centrum-CreateOrUpdateUnitResponse"></a>

### CreateOrUpdateUnitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |






<a name="services-centrum-DeleteDispatchRequest"></a>

### DeleteDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-centrum-DeleteDispatchResponse"></a>

### DeleteDispatchResponse







<a name="services-centrum-DeleteUnitRequest"></a>

### DeleteUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |






<a name="services-centrum-DeleteUnitResponse"></a>

### DeleteUnitResponse







<a name="services-centrum-Dispatchers"></a>

### Dispatchers



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [resources.centrum.Dispatchers](#resources-centrum-Dispatchers) | repeated |  |






<a name="services-centrum-GetDispatchHeatmapRequest"></a>

### GetDispatchHeatmapRequest







<a name="services-centrum-GetDispatchHeatmapResponse"></a>

### GetDispatchHeatmapResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_entries` | [int32](#int32) |  |  |
| `entries` | [resources.livemap.HeatmapEntry](#resources-livemap-HeatmapEntry) | repeated |  |






<a name="services-centrum-GetDispatchRequest"></a>

### GetDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-centrum-GetDispatchResponse"></a>

### GetDispatchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-GetSettingsRequest"></a>

### GetSettingsRequest







<a name="services-centrum-GetSettingsResponse"></a>

### GetSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |






<a name="services-centrum-JobAccess"></a>

### JobAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatches` | [JobAccessEntry](#services-centrum-JobAccessEntry) | repeated |  |






<a name="services-centrum-JobAccessEntry"></a>

### JobAccessEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [resources.jobs.Job](#resources-jobs-Job) |  |  |
| `access` | [resources.centrum.CentrumAccessLevel](#resources-centrum-CentrumAccessLevel) |  |  |






<a name="services-centrum-JoinUnitRequest"></a>

### JoinUnitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) | optional |  |






<a name="services-centrum-JoinUnitResponse"></a>

### JoinUnitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |






<a name="services-centrum-LatestState"></a>

### LatestState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [Dispatchers](#services-centrum-Dispatchers) |  |  |
| `own_unit_id` | [uint64](#uint64) | optional |  |
| `units` | [resources.centrum.Unit](#resources-centrum-Unit) | repeated | Send the current units and dispatches |
| `dispatches` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) | repeated |  |






<a name="services-centrum-ListDispatchActivityRequest"></a>

### ListDispatchActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `id` | [uint64](#uint64) |  |  |






<a name="services-centrum-ListDispatchActivityResponse"></a>

### ListDispatchActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `activity` | [resources.centrum.DispatchStatus](#resources-centrum-DispatchStatus) | repeated |  |






<a name="services-centrum-ListDispatchesRequest"></a>

### ListDispatchesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `status` | [resources.centrum.StatusDispatch](#resources-centrum-StatusDispatch) | repeated |  |
| `not_status` | [resources.centrum.StatusDispatch](#resources-centrum-StatusDispatch) | repeated |  |
| `ids` | [uint64](#uint64) | repeated |  |
| `postal` | [string](#string) | optional |  |






<a name="services-centrum-ListDispatchesResponse"></a>

### ListDispatchesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `dispatches` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) | repeated |  |






<a name="services-centrum-ListUnitActivityRequest"></a>

### ListUnitActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `id` | [uint64](#uint64) |  |  |






<a name="services-centrum-ListUnitActivityResponse"></a>

### ListUnitActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `activity` | [resources.centrum.UnitStatus](#resources-centrum-UnitStatus) | repeated |  |






<a name="services-centrum-ListUnitsRequest"></a>

### ListUnitsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [resources.centrum.StatusUnit](#resources-centrum-StatusUnit) | repeated |  |






<a name="services-centrum-ListUnitsResponse"></a>

### ListUnitsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `units` | [resources.centrum.Unit](#resources-centrum-Unit) | repeated |  |






<a name="services-centrum-StreamHandshake"></a>

### StreamHandshake



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `server_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `settings` | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |
| `job_access` | [JobAccess](#services-centrum-JobAccess) |  |  |






<a name="services-centrum-StreamRequest"></a>

### StreamRequest







<a name="services-centrum-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `handshake` | [StreamHandshake](#services-centrum-StreamHandshake) |  |  |
| `latest_state` | [LatestState](#services-centrum-LatestState) |  |  |
| `settings` | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |
| `job_access` | [JobAccess](#services-centrum-JobAccess) |  |  |
| `dispatchers` | [resources.centrum.Dispatchers](#resources-centrum-Dispatchers) |  |  |
| `unit_deleted` | [uint64](#uint64) |  |  |
| `unit_updated` | [resources.centrum.Unit](#resources-centrum-Unit) |  |  |
| `unit_status` | [resources.centrum.UnitStatus](#resources-centrum-UnitStatus) |  |  |
| `dispatch_deleted` | [uint64](#uint64) |  |  |
| `dispatch_updated` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |
| `dispatch_status` | [resources.centrum.DispatchStatus](#resources-centrum-DispatchStatus) |  |  |






<a name="services-centrum-TakeControlRequest"></a>

### TakeControlRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signon` | [bool](#bool) |  |  |






<a name="services-centrum-TakeControlResponse"></a>

### TakeControlResponse







<a name="services-centrum-TakeDispatchRequest"></a>

### TakeDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_ids` | [uint64](#uint64) | repeated |  |
| `resp` | [resources.centrum.TakeDispatchResp](#resources-centrum-TakeDispatchResp) |  |  |
| `reason` | [string](#string) | optional | @sanitize |






<a name="services-centrum-TakeDispatchResponse"></a>

### TakeDispatchResponse







<a name="services-centrum-UpdateDispatchRequest"></a>

### UpdateDispatchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |






<a name="services-centrum-UpdateDispatchResponse"></a>

### UpdateDispatchResponse







<a name="services-centrum-UpdateDispatchStatusRequest"></a>

### UpdateDispatchStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `status` | [resources.centrum.StatusDispatch](#resources-centrum-StatusDispatch) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |






<a name="services-centrum-UpdateDispatchStatusResponse"></a>

### UpdateDispatchStatusResponse







<a name="services-centrum-UpdateDispatchersRequest"></a>

### UpdateDispatchersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_remove` | [int32](#int32) | repeated |  |






<a name="services-centrum-UpdateDispatchersResponse"></a>

### UpdateDispatchersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [resources.centrum.Dispatchers](#resources-centrum-Dispatchers) |  |  |






<a name="services-centrum-UpdateSettingsRequest"></a>

### UpdateSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |






<a name="services-centrum-UpdateSettingsResponse"></a>

### UpdateSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.Settings](#resources-centrum-Settings) |  |  |






<a name="services-centrum-UpdateUnitStatusRequest"></a>

### UpdateUnitStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `status` | [resources.centrum.StatusUnit](#resources-centrum-StatusUnit) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |






<a name="services-centrum-UpdateUnitStatusResponse"></a>

### UpdateUnitStatusResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-centrum-CentrumService"></a>

### CentrumService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `UpdateSettings` | [UpdateSettingsRequest](#services-centrum-UpdateSettingsRequest) | [UpdateSettingsResponse](#services-centrum-UpdateSettingsResponse) | @perm: Attrs=Access/StringList:[]string{"Shared"} |
| `CreateDispatch` | [CreateDispatchRequest](#services-centrum-CreateDispatchRequest) | [CreateDispatchResponse](#services-centrum-CreateDispatchResponse) | @perm |
| `UpdateDispatch` | [UpdateDispatchRequest](#services-centrum-UpdateDispatchRequest) | [UpdateDispatchResponse](#services-centrum-UpdateDispatchResponse) | @perm |
| `DeleteDispatch` | [DeleteDispatchRequest](#services-centrum-DeleteDispatchRequest) | [DeleteDispatchResponse](#services-centrum-DeleteDispatchResponse) | @perm |
| `TakeControl` | [TakeControlRequest](#services-centrum-TakeControlRequest) | [TakeControlResponse](#services-centrum-TakeControlResponse) | @perm |
| `AssignDispatch` | [AssignDispatchRequest](#services-centrum-AssignDispatchRequest) | [AssignDispatchResponse](#services-centrum-AssignDispatchResponse) | @perm: Name=TakeControl |
| `AssignUnit` | [AssignUnitRequest](#services-centrum-AssignUnitRequest) | [AssignUnitResponse](#services-centrum-AssignUnitResponse) | @perm: Name=TakeControl |
| `GetDispatchHeatmap` | [GetDispatchHeatmapRequest](#services-centrum-GetDispatchHeatmapRequest) | [GetDispatchHeatmapResponse](#services-centrum-GetDispatchHeatmapResponse) | @perm: Name=TakeControl |
| `UpdateDispatchers` | [UpdateDispatchersRequest](#services-centrum-UpdateDispatchersRequest) | [UpdateDispatchersResponse](#services-centrum-UpdateDispatchersResponse) | @perm |
| `Stream` | [StreamRequest](#services-centrum-StreamRequest) | [StreamResponse](#services-centrum-StreamResponse) stream | @perm |
| `GetSettings` | [GetSettingsRequest](#services-centrum-GetSettingsRequest) | [GetSettingsResponse](#services-centrum-GetSettingsResponse) | @perm: Name=Stream |
| `JoinUnit` | [JoinUnitRequest](#services-centrum-JoinUnitRequest) | [JoinUnitResponse](#services-centrum-JoinUnitResponse) | @perm: Name=Stream |
| `ListUnits` | [ListUnitsRequest](#services-centrum-ListUnitsRequest) | [ListUnitsResponse](#services-centrum-ListUnitsResponse) | @perm: Name=Stream |
| `ListUnitActivity` | [ListUnitActivityRequest](#services-centrum-ListUnitActivityRequest) | [ListUnitActivityResponse](#services-centrum-ListUnitActivityResponse) | @perm: Name=Stream |
| `GetDispatch` | [GetDispatchRequest](#services-centrum-GetDispatchRequest) | [GetDispatchResponse](#services-centrum-GetDispatchResponse) | @perm: Name=Stream |
| `ListDispatches` | [ListDispatchesRequest](#services-centrum-ListDispatchesRequest) | [ListDispatchesResponse](#services-centrum-ListDispatchesResponse) | @perm: Name=Stream |
| `ListDispatchActivity` | [ListDispatchActivityRequest](#services-centrum-ListDispatchActivityRequest) | [ListDispatchActivityResponse](#services-centrum-ListDispatchActivityResponse) | @perm: Name=Stream |
| `CreateOrUpdateUnit` | [CreateOrUpdateUnitRequest](#services-centrum-CreateOrUpdateUnitRequest) | [CreateOrUpdateUnitResponse](#services-centrum-CreateOrUpdateUnitResponse) | @perm |
| `DeleteUnit` | [DeleteUnitRequest](#services-centrum-DeleteUnitRequest) | [DeleteUnitResponse](#services-centrum-DeleteUnitResponse) | @perm |
| `TakeDispatch` | [TakeDispatchRequest](#services-centrum-TakeDispatchRequest) | [TakeDispatchResponse](#services-centrum-TakeDispatchResponse) | @perm |
| `UpdateUnitStatus` | [UpdateUnitStatusRequest](#services-centrum-UpdateUnitStatusRequest) | [UpdateUnitStatusResponse](#services-centrum-UpdateUnitStatusResponse) | @perm: Name=TakeDispatch |
| `UpdateDispatchStatus` | [UpdateDispatchStatusRequest](#services-centrum-UpdateDispatchStatusRequest) | [UpdateDispatchStatusResponse](#services-centrum-UpdateDispatchStatusResponse) | @perm: Name=TakeDispatch |

 <!-- end services -->



<a name="services_completor_completor-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/completor/completor.proto



<a name="services-completor-CompleteCitizenLabelsRequest"></a>

### CompleteCitizenLabelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |






<a name="services-completor-CompleteCitizenLabelsResponse"></a>

### CompleteCitizenLabelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.Label](#resources-users-Label) | repeated |  |






<a name="services-completor-CompleteCitizensRequest"></a>

### CompleteCitizensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |
| `current_job` | [bool](#bool) | optional |  |
| `on_duty` | [bool](#bool) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `user_ids_only` | [bool](#bool) | optional |  |






<a name="services-completor-CompleteCitizensRespoonse"></a>

### CompleteCitizensRespoonse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.users.UserShort](#resources-users-UserShort) | repeated | @gotags: alias:"user" |






<a name="services-completor-CompleteDocumentCategoriesRequest"></a>

### CompleteDocumentCategoriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |






<a name="services-completor-CompleteDocumentCategoriesResponse"></a>

### CompleteDocumentCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `categories` | [resources.documents.Category](#resources-documents-Category) | repeated |  |






<a name="services-completor-CompleteJobsRequest"></a>

### CompleteJobsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) | optional |  |
| `exact_match` | [bool](#bool) | optional |  |
| `current_job` | [bool](#bool) | optional |  |






<a name="services-completor-CompleteJobsResponse"></a>

### CompleteJobsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.jobs.Job](#resources-jobs-Job) | repeated |  |






<a name="services-completor-ListLawBooksRequest"></a>

### ListLawBooksRequest







<a name="services-completor-ListLawBooksResponse"></a>

### ListLawBooksResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `books` | [resources.laws.LawBook](#resources-laws-LawBook) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-completor-CompletorService"></a>

### CompletorService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `CompleteCitizens` | [CompleteCitizensRequest](#services-completor-CompleteCitizensRequest) | [CompleteCitizensRespoonse](#services-completor-CompleteCitizensRespoonse) | @perm |
| `CompleteJobs` | [CompleteJobsRequest](#services-completor-CompleteJobsRequest) | [CompleteJobsResponse](#services-completor-CompleteJobsResponse) | @perm: Name=Any |
| `CompleteDocumentCategories` | [CompleteDocumentCategoriesRequest](#services-completor-CompleteDocumentCategoriesRequest) | [CompleteDocumentCategoriesResponse](#services-completor-CompleteDocumentCategoriesResponse) | @perm: Attrs=Jobs/JobList |
| `ListLawBooks` | [ListLawBooksRequest](#services-completor-ListLawBooksRequest) | [ListLawBooksResponse](#services-completor-ListLawBooksResponse) | @perm: Name=Any |
| `CompleteCitizenLabels` | [CompleteCitizenLabelsRequest](#services-completor-CompleteCitizenLabelsRequest) | [CompleteCitizenLabelsResponse](#services-completor-CompleteCitizenLabelsResponse) | @perm: Attrs=Jobs/JobList |

 <!-- end services -->



<a name="services_jobs_conduct-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/conduct.proto



<a name="services-jobs-CreateConductEntryRequest"></a>

### CreateConductEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |






<a name="services-jobs-CreateConductEntryResponse"></a>

### CreateConductEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |






<a name="services-jobs-DeleteConductEntryRequest"></a>

### DeleteConductEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-jobs-DeleteConductEntryResponse"></a>

### DeleteConductEntryResponse







<a name="services-jobs-ListConductEntriesRequest"></a>

### ListConductEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `types` | [resources.jobs.ConductType](#resources-jobs-ConductType) | repeated | Search params |
| `show_expired` | [bool](#bool) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `ids` | [uint64](#uint64) | repeated |  |






<a name="services-jobs-ListConductEntriesResponse"></a>

### ListConductEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `entries` | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) | repeated |  |






<a name="services-jobs-UpdateConductEntryRequest"></a>

### UpdateConductEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |






<a name="services-jobs-UpdateConductEntryResponse"></a>

### UpdateConductEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resources-jobs-ConductEntry) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-jobs-ConductService"></a>

### ConductService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListConductEntries` | [ListConductEntriesRequest](#services-jobs-ListConductEntriesRequest) | [ListConductEntriesResponse](#services-jobs-ListConductEntriesResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "All"} |
| `CreateConductEntry` | [CreateConductEntryRequest](#services-jobs-CreateConductEntryRequest) | [CreateConductEntryResponse](#services-jobs-CreateConductEntryResponse) | @perm |
| `UpdateConductEntry` | [UpdateConductEntryRequest](#services-jobs-UpdateConductEntryRequest) | [UpdateConductEntryResponse](#services-jobs-UpdateConductEntryResponse) | @perm |
| `DeleteConductEntry` | [DeleteConductEntryRequest](#services-jobs-DeleteConductEntryRequest) | [DeleteConductEntryResponse](#services-jobs-DeleteConductEntryResponse) | @perm |

 <!-- end services -->



<a name="services_jobs_jobs-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/jobs.proto



<a name="services-jobs-GetColleagueLabelsRequest"></a>

### GetColleagueLabelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) | optional |  |






<a name="services-jobs-GetColleagueLabelsResponse"></a>

### GetColleagueLabelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.Label](#resources-jobs-Label) | repeated |  |






<a name="services-jobs-GetColleagueLabelsStatsRequest"></a>

### GetColleagueLabelsStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `label_ids` | [uint64](#uint64) | repeated |  |






<a name="services-jobs-GetColleagueLabelsStatsResponse"></a>

### GetColleagueLabelsStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `count` | [resources.jobs.LabelCount](#resources-jobs-LabelCount) | repeated |  |






<a name="services-jobs-GetColleagueRequest"></a>

### GetColleagueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `info_only` | [bool](#bool) | optional |  |






<a name="services-jobs-GetColleagueResponse"></a>

### GetColleagueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague` | [resources.jobs.Colleague](#resources-jobs-Colleague) |  |  |






<a name="services-jobs-GetMOTDRequest"></a>

### GetMOTDRequest







<a name="services-jobs-GetMOTDResponse"></a>

### GetMOTDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  |  |






<a name="services-jobs-GetSelfRequest"></a>

### GetSelfRequest







<a name="services-jobs-GetSelfResponse"></a>

### GetSelfResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague` | [resources.jobs.Colleague](#resources-jobs-Colleague) |  |  |






<a name="services-jobs-ListColleagueActivityRequest"></a>

### ListColleagueActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `user_ids` | [int32](#int32) | repeated | Search params |
| `activity_types` | [resources.jobs.ColleagueActivityType](#resources-jobs-ColleagueActivityType) | repeated |  |






<a name="services-jobs-ListColleagueActivityResponse"></a>

### ListColleagueActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `activity` | [resources.jobs.ColleagueActivity](#resources-jobs-ColleagueActivity) | repeated |  |






<a name="services-jobs-ListColleaguesRequest"></a>

### ListColleaguesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `search` | [string](#string) |  | Search params |
| `user_ids` | [int32](#int32) | repeated |  |
| `user_only` | [bool](#bool) | optional |  |
| `absent` | [bool](#bool) | optional |  |
| `label_ids` | [uint64](#uint64) | repeated |  |
| `name_prefix` | [string](#string) | optional |  |
| `name_suffix` | [string](#string) | optional |  |






<a name="services-jobs-ListColleaguesResponse"></a>

### ListColleaguesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `colleagues` | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |






<a name="services-jobs-ManageLabelsRequest"></a>

### ManageLabelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.Label](#resources-jobs-Label) | repeated |  |






<a name="services-jobs-ManageLabelsResponse"></a>

### ManageLabelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.Label](#resources-jobs-Label) | repeated |  |






<a name="services-jobs-SetColleaguePropsRequest"></a>

### SetColleaguePropsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.jobs.ColleagueProps](#resources-jobs-ColleagueProps) |  |  |
| `reason` | [string](#string) |  | @sanitize |






<a name="services-jobs-SetColleaguePropsResponse"></a>

### SetColleaguePropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.jobs.ColleagueProps](#resources-jobs-ColleagueProps) |  |  |






<a name="services-jobs-SetMOTDRequest"></a>

### SetMOTDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  | @sanitize: method=StripTags |






<a name="services-jobs-SetMOTDResponse"></a>

### SetMOTDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-jobs-JobsService"></a>

### JobsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListColleagues` | [ListColleaguesRequest](#services-jobs-ListColleaguesRequest) | [ListColleaguesResponse](#services-jobs-ListColleaguesResponse) | @perm |
| `GetSelf` | [GetSelfRequest](#services-jobs-GetSelfRequest) | [GetSelfResponse](#services-jobs-GetSelfResponse) | @perm: Name=ListColleagues |
| `GetColleague` | [GetColleagueRequest](#services-jobs-GetColleagueRequest) | [GetColleagueResponse](#services-jobs-GetColleagueResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"Note", "Labels"} |
| `ListColleagueActivity` | [ListColleagueActivityRequest](#services-jobs-ListColleagueActivityRequest) | [ListColleagueActivityResponse](#services-jobs-ListColleagueActivityResponse) | @perm: Attrs=Types/StringList:[]string{"HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME"} |
| `SetColleagueProps` | [SetColleaguePropsRequest](#services-jobs-SetColleaguePropsRequest) | [SetColleaguePropsResponse](#services-jobs-SetColleaguePropsResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"AbsenceDate", "Note", "Labels", "Name"} |
| `GetColleagueLabels` | [GetColleagueLabelsRequest](#services-jobs-GetColleagueLabelsRequest) | [GetColleagueLabelsResponse](#services-jobs-GetColleagueLabelsResponse) | @perm: Name=GetColleague |
| `ManageLabels` | [ManageLabelsRequest](#services-jobs-ManageLabelsRequest) | [ManageLabelsResponse](#services-jobs-ManageLabelsResponse) | @perm |
| `GetColleagueLabelsStats` | [GetColleagueLabelsStatsRequest](#services-jobs-GetColleagueLabelsStatsRequest) | [GetColleagueLabelsStatsResponse](#services-jobs-GetColleagueLabelsStatsResponse) | @perm: Name=GetColleague |
| `GetMOTD` | [GetMOTDRequest](#services-jobs-GetMOTDRequest) | [GetMOTDResponse](#services-jobs-GetMOTDResponse) | @perm: Name=Any |
| `SetMOTD` | [SetMOTDRequest](#services-jobs-SetMOTDRequest) | [SetMOTDResponse](#services-jobs-SetMOTDResponse) | @perm |

 <!-- end services -->



<a name="services_jobs_timeclock-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/jobs/timeclock.proto



<a name="services-jobs-GetTimeclockStatsRequest"></a>

### GetTimeclockStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) | optional |  |






<a name="services-jobs-GetTimeclockStatsResponse"></a>

### GetTimeclockStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stats` | [resources.jobs.TimeclockStats](#resources-jobs-TimeclockStats) |  |  |
| `weekly` | [resources.jobs.TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats) | repeated |  |






<a name="services-jobs-ListInactiveEmployeesRequest"></a>

### ListInactiveEmployeesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `days` | [int32](#int32) |  | Search params |






<a name="services-jobs-ListInactiveEmployeesResponse"></a>

### ListInactiveEmployeesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `colleagues` | [resources.jobs.Colleague](#resources-jobs-Colleague) | repeated |  |






<a name="services-jobs-ListTimeclockRequest"></a>

### ListTimeclockRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `user_mode` | [resources.jobs.TimeclockViewMode](#resources-jobs-TimeclockViewMode) |  | Search params |
| `mode` | [resources.jobs.TimeclockMode](#resources-jobs-TimeclockMode) |  |  |
| `date` | [resources.common.database.DateRange](#resources-common-database-DateRange) | optional |  |
| `per_day` | [bool](#bool) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |






<a name="services-jobs-ListTimeclockResponse"></a>

### ListTimeclockResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `stats` | [resources.jobs.TimeclockStats](#resources-jobs-TimeclockStats) |  |  |
| `stats_weekly` | [resources.jobs.TimeclockWeeklyStats](#resources-jobs-TimeclockWeeklyStats) | repeated |  |
| `daily` | [TimeclockDay](#services-jobs-TimeclockDay) |  |  |
| `weekly` | [TimeclockWeekly](#services-jobs-TimeclockWeekly) |  |  |
| `range` | [TimeclockRange](#services-jobs-TimeclockRange) |  |  |






<a name="services-jobs-TimeclockDay"></a>

### TimeclockDay



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  |  |
| `entries` | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |






<a name="services-jobs-TimeclockRange"></a>

### TimeclockRange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: sql:"primary_key" |
| `entries` | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |






<a name="services-jobs-TimeclockWeekly"></a>

### TimeclockWeekly



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) |  | @gotags: sql:"primary_key" |
| `entries` | [resources.jobs.TimeclockEntry](#resources-jobs-TimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-jobs-TimeclockService"></a>

### TimeclockService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTimeclock` | [ListTimeclockRequest](#services-jobs-ListTimeclockRequest) | [ListTimeclockResponse](#services-jobs-ListTimeclockResponse) | @perm: Attrs=Access/StringList:[]string{"All"} |
| `GetTimeclockStats` | [GetTimeclockStatsRequest](#services-jobs-GetTimeclockStatsRequest) | [GetTimeclockStatsResponse](#services-jobs-GetTimeclockStatsResponse) | @perm: Name=ListTimeclock |
| `ListInactiveEmployees` | [ListInactiveEmployeesRequest](#services-jobs-ListInactiveEmployeesRequest) | [ListInactiveEmployeesResponse](#services-jobs-ListInactiveEmployeesResponse) | @perm |

 <!-- end services -->



<a name="services_qualifications_qualifications-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/qualifications/qualifications.proto



<a name="services-qualifications-CreateOrUpdateQualificationRequestRequest"></a>

### CreateOrUpdateQualificationRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.qualifications.QualificationRequest](#resources-qualifications-QualificationRequest) |  |  |






<a name="services-qualifications-CreateOrUpdateQualificationRequestResponse"></a>

### CreateOrUpdateQualificationRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.qualifications.QualificationRequest](#resources-qualifications-QualificationRequest) |  |  |






<a name="services-qualifications-CreateOrUpdateQualificationResultRequest"></a>

### CreateOrUpdateQualificationResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [resources.qualifications.QualificationResult](#resources-qualifications-QualificationResult) |  |  |
| `grading` | [resources.qualifications.ExamGrading](#resources-qualifications-ExamGrading) | optional |  |






<a name="services-qualifications-CreateOrUpdateQualificationResultResponse"></a>

### CreateOrUpdateQualificationResultResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [resources.qualifications.QualificationResult](#resources-qualifications-QualificationResult) |  |  |






<a name="services-qualifications-CreateQualificationRequest"></a>

### CreateQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |






<a name="services-qualifications-CreateQualificationResponse"></a>

### CreateQualificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |






<a name="services-qualifications-DeleteQualificationReqRequest"></a>

### DeleteQualificationReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |






<a name="services-qualifications-DeleteQualificationReqResponse"></a>

### DeleteQualificationReqResponse







<a name="services-qualifications-DeleteQualificationRequest"></a>

### DeleteQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |






<a name="services-qualifications-DeleteQualificationResponse"></a>

### DeleteQualificationResponse







<a name="services-qualifications-DeleteQualificationResultRequest"></a>

### DeleteQualificationResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result_id` | [uint64](#uint64) |  |  |






<a name="services-qualifications-DeleteQualificationResultResponse"></a>

### DeleteQualificationResultResponse







<a name="services-qualifications-GetExamInfoRequest"></a>

### GetExamInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |






<a name="services-qualifications-GetExamInfoResponse"></a>

### GetExamInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.QualificationShort](#resources-qualifications-QualificationShort) |  |  |
| `question_count` | [int32](#int32) |  |  |
| `exam_user` | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) | optional |  |






<a name="services-qualifications-GetQualificationAccessRequest"></a>

### GetQualificationAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |






<a name="services-qualifications-GetQualificationAccessResponse"></a>

### GetQualificationAccessResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [resources.qualifications.QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |






<a name="services-qualifications-GetQualificationRequest"></a>

### GetQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `with_exam` | [bool](#bool) | optional |  |






<a name="services-qualifications-GetQualificationResponse"></a>

### GetQualificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.Qualification](#resources-qualifications-Qualification) |  |  |






<a name="services-qualifications-GetUserExamRequest"></a>

### GetUserExamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |






<a name="services-qualifications-GetUserExamResponse"></a>

### GetUserExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exam` | [resources.qualifications.ExamQuestions](#resources-qualifications-ExamQuestions) |  |  |
| `exam_user` | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) |  |  |
| `responses` | [resources.qualifications.ExamResponses](#resources-qualifications-ExamResponses) |  |  |
| `grading` | [resources.qualifications.ExamGrading](#resources-qualifications-ExamGrading) |  |  |






<a name="services-qualifications-ListQualificationRequestsRequest"></a>

### ListQualificationRequestsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `qualification_id` | [uint64](#uint64) | optional | Search params |
| `status` | [resources.qualifications.RequestStatus](#resources-qualifications-RequestStatus) | repeated |  |
| `user_id` | [int32](#int32) | optional |  |






<a name="services-qualifications-ListQualificationRequestsResponse"></a>

### ListQualificationRequestsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `requests` | [resources.qualifications.QualificationRequest](#resources-qualifications-QualificationRequest) | repeated |  |






<a name="services-qualifications-ListQualificationsRequest"></a>

### ListQualificationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `search` | [string](#string) | optional | Search params |
| `job` | [string](#string) | optional |  |






<a name="services-qualifications-ListQualificationsResponse"></a>

### ListQualificationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `qualifications` | [resources.qualifications.Qualification](#resources-qualifications-Qualification) | repeated |  |






<a name="services-qualifications-ListQualificationsResultsRequest"></a>

### ListQualificationsResultsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `qualification_id` | [uint64](#uint64) | optional | Search params |
| `status` | [resources.qualifications.ResultStatus](#resources-qualifications-ResultStatus) | repeated |  |
| `user_id` | [int32](#int32) | optional |  |






<a name="services-qualifications-ListQualificationsResultsResponse"></a>

### ListQualificationsResultsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `results` | [resources.qualifications.QualificationResult](#resources-qualifications-QualificationResult) | repeated |  |






<a name="services-qualifications-SetQualificationAccessRequest"></a>

### SetQualificationAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `access` | [resources.qualifications.QualificationAccess](#resources-qualifications-QualificationAccess) |  |  |






<a name="services-qualifications-SetQualificationAccessResponse"></a>

### SetQualificationAccessResponse







<a name="services-qualifications-SubmitExamRequest"></a>

### SubmitExamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `responses` | [resources.qualifications.ExamResponses](#resources-qualifications-ExamResponses) |  |  |






<a name="services-qualifications-SubmitExamResponse"></a>

### SubmitExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](#google-protobuf-Duration) |  |  |






<a name="services-qualifications-TakeExamRequest"></a>

### TakeExamRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `cancel` | [bool](#bool) | optional |  |






<a name="services-qualifications-TakeExamResponse"></a>

### TakeExamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exam` | [resources.qualifications.ExamQuestions](#resources-qualifications-ExamQuestions) |  |  |
| `exam_user` | [resources.qualifications.ExamUser](#resources-qualifications-ExamUser) |  |  |






<a name="services-qualifications-UpdateQualificationRequest"></a>

### UpdateQualificationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.Qualification](#resources-qualifications-Qualification) |  |  |






<a name="services-qualifications-UpdateQualificationResponse"></a>

### UpdateQualificationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-qualifications-QualificationsService"></a>

### QualificationsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListQualifications` | [ListQualificationsRequest](#services-qualifications-ListQualificationsRequest) | [ListQualificationsResponse](#services-qualifications-ListQualificationsResponse) | @perm |
| `GetQualification` | [GetQualificationRequest](#services-qualifications-GetQualificationRequest) | [GetQualificationResponse](#services-qualifications-GetQualificationResponse) | @perm: Name=ListQualifications |
| `CreateQualification` | [CreateQualificationRequest](#services-qualifications-CreateQualificationRequest) | [CreateQualificationResponse](#services-qualifications-CreateQualificationResponse) | @perm: Name=UpdateQualification |
| `UpdateQualification` | [UpdateQualificationRequest](#services-qualifications-UpdateQualificationRequest) | [UpdateQualificationResponse](#services-qualifications-UpdateQualificationResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Fields/StringList:[]string{"Public"} |
| `DeleteQualification` | [DeleteQualificationRequest](#services-qualifications-DeleteQualificationRequest) | [DeleteQualificationResponse](#services-qualifications-DeleteQualificationResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `ListQualificationRequests` | [ListQualificationRequestsRequest](#services-qualifications-ListQualificationRequestsRequest) | [ListQualificationRequestsResponse](#services-qualifications-ListQualificationRequestsResponse) | @perm: Name=ListQualifications |
| `CreateOrUpdateQualificationRequest` | [CreateOrUpdateQualificationRequestRequest](#services-qualifications-CreateOrUpdateQualificationRequestRequest) | [CreateOrUpdateQualificationRequestResponse](#services-qualifications-CreateOrUpdateQualificationRequestResponse) | @perm: Name=ListQualifications |
| `DeleteQualificationReq` | [DeleteQualificationReqRequest](#services-qualifications-DeleteQualificationReqRequest) | [DeleteQualificationReqResponse](#services-qualifications-DeleteQualificationReqResponse) | @perm: Name=ListQualifications |
| `ListQualificationsResults` | [ListQualificationsResultsRequest](#services-qualifications-ListQualificationsResultsRequest) | [ListQualificationsResultsResponse](#services-qualifications-ListQualificationsResultsResponse) | @perm: Name=ListQualifications |
| `CreateOrUpdateQualificationResult` | [CreateOrUpdateQualificationResultRequest](#services-qualifications-CreateOrUpdateQualificationResultRequest) | [CreateOrUpdateQualificationResultResponse](#services-qualifications-CreateOrUpdateQualificationResultResponse) | @perm: Name=ListQualifications |
| `DeleteQualificationResult` | [DeleteQualificationResultRequest](#services-qualifications-DeleteQualificationResultRequest) | [DeleteQualificationResultResponse](#services-qualifications-DeleteQualificationResultResponse) | @perm: Name=ListQualifications |
| `GetExamInfo` | [GetExamInfoRequest](#services-qualifications-GetExamInfoRequest) | [GetExamInfoResponse](#services-qualifications-GetExamInfoResponse) | @perm: Name=ListQualifications |
| `TakeExam` | [TakeExamRequest](#services-qualifications-TakeExamRequest) | [TakeExamResponse](#services-qualifications-TakeExamResponse) | @perm: Name=ListQualifications |
| `SubmitExam` | [SubmitExamRequest](#services-qualifications-SubmitExamRequest) | [SubmitExamResponse](#services-qualifications-SubmitExamResponse) | @perm: Name=ListQualifications |
| `GetUserExam` | [GetUserExamRequest](#services-qualifications-GetUserExamRequest) | [GetUserExamResponse](#services-qualifications-GetUserExamResponse) | @perm: Name=ListQualifications |
| `UploadFile` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=UpdateQualification |

 <!-- end services -->



<a name="services_calendar_calendar-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/calendar/calendar.proto



<a name="services-calendar-CreateCalendarRequest"></a>

### CreateCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-CreateCalendarResponse"></a>

### CreateCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-CreateOrUpdateCalendarEntryRequest"></a>

### CreateOrUpdateCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |






<a name="services-calendar-CreateOrUpdateCalendarEntryResponse"></a>

### CreateOrUpdateCalendarEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) |  |  |






<a name="services-calendar-DeleteCalendarEntryRequest"></a>

### DeleteCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |






<a name="services-calendar-DeleteCalendarEntryResponse"></a>

### DeleteCalendarEntryResponse







<a name="services-calendar-DeleteCalendarRequest"></a>

### DeleteCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [uint64](#uint64) |  |  |






<a name="services-calendar-DeleteCalendarResponse"></a>

### DeleteCalendarResponse







<a name="services-calendar-GetCalendarEntryRequest"></a>

### GetCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |






<a name="services-calendar-GetCalendarEntryResponse"></a>

### GetCalendarEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) |  |  |






<a name="services-calendar-GetCalendarRequest"></a>

### GetCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [uint64](#uint64) |  |  |






<a name="services-calendar-GetCalendarResponse"></a>

### GetCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-GetUpcomingEntriesRequest"></a>

### GetUpcomingEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seconds` | [int32](#int32) |  |  |






<a name="services-calendar-GetUpcomingEntriesResponse"></a>

### GetUpcomingEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) | repeated |  |






<a name="services-calendar-ListCalendarEntriesRequest"></a>

### ListCalendarEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `year` | [int32](#int32) |  |  |
| `month` | [int32](#int32) |  |  |
| `calendar_ids` | [uint64](#uint64) | repeated |  |
| `show_hidden` | [bool](#bool) | optional |  |
| `after` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-calendar-ListCalendarEntriesResponse"></a>

### ListCalendarEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [resources.calendar.CalendarEntry](#resources-calendar-CalendarEntry) | repeated |  |






<a name="services-calendar-ListCalendarEntryRSVPRequest"></a>

### ListCalendarEntryRSVPRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `entry_id` | [uint64](#uint64) |  |  |






<a name="services-calendar-ListCalendarEntryRSVPResponse"></a>

### ListCalendarEntryRSVPResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `entries` | [resources.calendar.CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) | repeated |  |






<a name="services-calendar-ListCalendarsRequest"></a>

### ListCalendarsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `only_public` | [bool](#bool) |  |  |
| `min_access_level` | [resources.calendar.AccessLevel](#resources-calendar-AccessLevel) | optional |  |
| `after` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-calendar-ListCalendarsResponse"></a>

### ListCalendarsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `calendars` | [resources.calendar.Calendar](#resources-calendar-Calendar) | repeated |  |






<a name="services-calendar-ListSubscriptionsRequest"></a>

### ListSubscriptionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |






<a name="services-calendar-ListSubscriptionsResponse"></a>

### ListSubscriptionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `subs` | [resources.calendar.CalendarSub](#resources-calendar-CalendarSub) | repeated |  |






<a name="services-calendar-RSVPCalendarEntryRequest"></a>

### RSVPCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) |  |  |
| `subscribe` | [bool](#bool) |  |  |
| `remove` | [bool](#bool) | optional |  |






<a name="services-calendar-RSVPCalendarEntryResponse"></a>

### RSVPCalendarEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntryRSVP](#resources-calendar-CalendarEntryRSVP) | optional |  |






<a name="services-calendar-ShareCalendarEntryRequest"></a>

### ShareCalendarEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |






<a name="services-calendar-ShareCalendarEntryResponse"></a>

### ShareCalendarEntryResponse







<a name="services-calendar-SubscribeToCalendarRequest"></a>

### SubscribeToCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sub` | [resources.calendar.CalendarSub](#resources-calendar-CalendarSub) |  |  |
| `delete` | [bool](#bool) |  |  |






<a name="services-calendar-SubscribeToCalendarResponse"></a>

### SubscribeToCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sub` | [resources.calendar.CalendarSub](#resources-calendar-CalendarSub) |  |  |






<a name="services-calendar-UpdateCalendarRequest"></a>

### UpdateCalendarRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |






<a name="services-calendar-UpdateCalendarResponse"></a>

### UpdateCalendarResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resources-calendar-Calendar) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-calendar-CalendarService"></a>

### CalendarService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCalendars` | [ListCalendarsRequest](#services-calendar-ListCalendarsRequest) | [ListCalendarsResponse](#services-calendar-ListCalendarsResponse) | @perm: Name=Any |
| `GetCalendar` | [GetCalendarRequest](#services-calendar-GetCalendarRequest) | [GetCalendarResponse](#services-calendar-GetCalendarResponse) | @perm: Name=Any |
| `CreateCalendar` | [CreateCalendarRequest](#services-calendar-CreateCalendarRequest) | [CreateCalendarResponse](#services-calendar-CreateCalendarResponse) | @perm: Attrs=Fields/StringList:[]string{"Job", "Public"} |
| `UpdateCalendar` | [UpdateCalendarRequest](#services-calendar-UpdateCalendarRequest) | [UpdateCalendarResponse](#services-calendar-UpdateCalendarResponse) | @perm: Name=Any |
| `DeleteCalendar` | [DeleteCalendarRequest](#services-calendar-DeleteCalendarRequest) | [DeleteCalendarResponse](#services-calendar-DeleteCalendarResponse) | @perm: Name=Any |
| `ListCalendarEntries` | [ListCalendarEntriesRequest](#services-calendar-ListCalendarEntriesRequest) | [ListCalendarEntriesResponse](#services-calendar-ListCalendarEntriesResponse) | @perm: Name=Any |
| `GetUpcomingEntries` | [GetUpcomingEntriesRequest](#services-calendar-GetUpcomingEntriesRequest) | [GetUpcomingEntriesResponse](#services-calendar-GetUpcomingEntriesResponse) | @perm: Name=Any |
| `GetCalendarEntry` | [GetCalendarEntryRequest](#services-calendar-GetCalendarEntryRequest) | [GetCalendarEntryResponse](#services-calendar-GetCalendarEntryResponse) | @perm: Name=Any |
| `CreateOrUpdateCalendarEntry` | [CreateOrUpdateCalendarEntryRequest](#services-calendar-CreateOrUpdateCalendarEntryRequest) | [CreateOrUpdateCalendarEntryResponse](#services-calendar-CreateOrUpdateCalendarEntryResponse) | @perm: Name=Any |
| `DeleteCalendarEntry` | [DeleteCalendarEntryRequest](#services-calendar-DeleteCalendarEntryRequest) | [DeleteCalendarEntryResponse](#services-calendar-DeleteCalendarEntryResponse) | @perm: Name=Any |
| `ShareCalendarEntry` | [ShareCalendarEntryRequest](#services-calendar-ShareCalendarEntryRequest) | [ShareCalendarEntryResponse](#services-calendar-ShareCalendarEntryResponse) | @perm: Name=Any |
| `ListCalendarEntryRSVP` | [ListCalendarEntryRSVPRequest](#services-calendar-ListCalendarEntryRSVPRequest) | [ListCalendarEntryRSVPResponse](#services-calendar-ListCalendarEntryRSVPResponse) | @perm: Name=Any |
| `RSVPCalendarEntry` | [RSVPCalendarEntryRequest](#services-calendar-RSVPCalendarEntryRequest) | [RSVPCalendarEntryResponse](#services-calendar-RSVPCalendarEntryResponse) | @perm: Name=Any |
| `ListSubscriptions` | [ListSubscriptionsRequest](#services-calendar-ListSubscriptionsRequest) | [ListSubscriptionsResponse](#services-calendar-ListSubscriptionsResponse) | @perm: Name=Any |
| `SubscribeToCalendar` | [SubscribeToCalendarRequest](#services-calendar-SubscribeToCalendarRequest) | [SubscribeToCalendarResponse](#services-calendar-SubscribeToCalendarResponse) | @perm: Name=Any |

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
| `stats` | [GetStatsResponse.StatsEntry](#services-stats-GetStatsResponse-StatsEntry) | repeated |  |






<a name="services-stats-GetStatsResponse-StatsEntry"></a>

### GetStatsResponse.StatsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [resources.stats.Stat](#resources-stats-Stat) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-stats-StatsService"></a>

### StatsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetStats` | [GetStatsRequest](#services-stats-GetStatsRequest) | [GetStatsResponse](#services-stats-GetStatsResponse) |  |

 <!-- end services -->



<a name="services_internet_ads-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/internet/ads.proto



<a name="services-internet-GetAdsRequest"></a>

### GetAdsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ad_type` | [resources.internet.AdType](#resources-internet-AdType) |  |  |
| `count` | [int32](#int32) |  |  |






<a name="services-internet-GetAdsResponse"></a>

### GetAdsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ads` | [resources.internet.Ad](#resources-internet-Ad) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-internet-AdsService"></a>

### AdsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetAds` | [GetAdsRequest](#services-internet-GetAdsRequest) | [GetAdsResponse](#services-internet-GetAdsResponse) | @perm: Name=Any |

 <!-- end services -->



<a name="services_internet_domain-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/internet/domain.proto



<a name="services-internet-CheckDomainAvailabilityRequest"></a>

### CheckDomainAvailabilityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tld_id` | [uint64](#uint64) |  |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |






<a name="services-internet-CheckDomainAvailabilityResponse"></a>

### CheckDomainAvailabilityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `available` | [bool](#bool) |  |  |
| `transferable` | [bool](#bool) | optional |  |






<a name="services-internet-ListDomainsRequest"></a>

### ListDomainsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |






<a name="services-internet-ListDomainsResponse"></a>

### ListDomainsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `domains` | [resources.internet.Domain](#resources-internet-Domain) | repeated |  |






<a name="services-internet-ListTLDsRequest"></a>

### ListTLDsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `internal` | [bool](#bool) | optional |  |






<a name="services-internet-ListTLDsResponse"></a>

### ListTLDsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tlds` | [resources.internet.TLD](#resources-internet-TLD) | repeated |  |






<a name="services-internet-RegisterDomainRequest"></a>

### RegisterDomainRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tld_id` | [uint64](#uint64) |  |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `transfer_code` | [string](#string) | optional | In case a domain will be transfered |






<a name="services-internet-RegisterDomainResponse"></a>

### RegisterDomainResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [resources.internet.Domain](#resources-internet-Domain) |  |  |






<a name="services-internet-UpdateDomainRequest"></a>

### UpdateDomainRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain_id` | [uint64](#uint64) |  |  |
| `transferable` | [bool](#bool) |  |  |






<a name="services-internet-UpdateDomainResponse"></a>

### UpdateDomainResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [resources.internet.Domain](#resources-internet-Domain) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-internet-DomainService"></a>

### DomainService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTLDs` | [ListTLDsRequest](#services-internet-ListTLDsRequest) | [ListTLDsResponse](#services-internet-ListTLDsResponse) | @perm: Name=Any |
| `CheckDomainAvailability` | [CheckDomainAvailabilityRequest](#services-internet-CheckDomainAvailabilityRequest) | [CheckDomainAvailabilityResponse](#services-internet-CheckDomainAvailabilityResponse) | @perm: Name=Any |
| `RegisterDomain` | [RegisterDomainRequest](#services-internet-RegisterDomainRequest) | [RegisterDomainResponse](#services-internet-RegisterDomainResponse) | @perm: Name=Any |
| `ListDomains` | [ListDomainsRequest](#services-internet-ListDomainsRequest) | [ListDomainsResponse](#services-internet-ListDomainsResponse) | @perm: Name=Any |
| `UpdateDomain` | [UpdateDomainRequest](#services-internet-UpdateDomainRequest) | [UpdateDomainResponse](#services-internet-UpdateDomainResponse) | @perm: Name=Any |

 <!-- end services -->



<a name="services_internet_internet-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/internet/internet.proto



<a name="services-internet-GetPageRequest"></a>

### GetPageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  |  |
| `path` | [string](#string) |  |  |






<a name="services-internet-GetPageResponse"></a>

### GetPageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.internet.Page](#resources-internet-Page) | optional |  |






<a name="services-internet-SearchRequest"></a>

### SearchRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |
| `domain_id` | [uint64](#uint64) | optional |  |






<a name="services-internet-SearchResponse"></a>

### SearchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `results` | [resources.internet.SearchResult](#resources-internet-SearchResult) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-internet-InternetService"></a>

### InternetService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Search` | [SearchRequest](#services-internet-SearchRequest) | [SearchResponse](#services-internet-SearchResponse) | @perm: Name=Any |
| `GetPage` | [GetPageRequest](#services-internet-GetPageRequest) | [GetPageResponse](#services-internet-GetPageResponse) | @perm: Name=Any |

 <!-- end services -->



<a name="services_mailer_mailer-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/mailer/mailer.proto



<a name="services-mailer-CreateOrUpdateEmailRequest"></a>

### CreateOrUpdateEmailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.Email](#resources-mailer-Email) |  |  |






<a name="services-mailer-CreateOrUpdateEmailResponse"></a>

### CreateOrUpdateEmailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.Email](#resources-mailer-Email) |  |  |






<a name="services-mailer-CreateOrUpdateTemplateRequest"></a>

### CreateOrUpdateTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.Template](#resources-mailer-Template) |  |  |






<a name="services-mailer-CreateOrUpdateTemplateResponse"></a>

### CreateOrUpdateTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.Template](#resources-mailer-Template) |  |  |






<a name="services-mailer-CreateThreadRequest"></a>

### CreateThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.Thread](#resources-mailer-Thread) |  |  |
| `message` | [resources.mailer.Message](#resources-mailer-Message) |  |  |
| `recipients` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="services-mailer-CreateThreadResponse"></a>

### CreateThreadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.Thread](#resources-mailer-Thread) |  |  |






<a name="services-mailer-DeleteEmailRequest"></a>

### DeleteEmailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteEmailResponse"></a>

### DeleteEmailResponse







<a name="services-mailer-DeleteMessageRequest"></a>

### DeleteMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |
| `message_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteMessageResponse"></a>

### DeleteMessageResponse







<a name="services-mailer-DeleteTemplateRequest"></a>

### DeleteTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `id` | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteTemplateResponse"></a>

### DeleteTemplateResponse







<a name="services-mailer-DeleteThreadRequest"></a>

### DeleteThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-DeleteThreadResponse"></a>

### DeleteThreadResponse







<a name="services-mailer-GetEmailProposalsRequest"></a>

### GetEmailProposalsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `input` | [string](#string) |  |  |
| `job` | [bool](#bool) | optional |  |
| `user_id` | [int32](#int32) | optional |  |






<a name="services-mailer-GetEmailProposalsResponse"></a>

### GetEmailProposalsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `emails` | [string](#string) | repeated |  |
| `domains` | [string](#string) | repeated |  |






<a name="services-mailer-GetEmailRequest"></a>

### GetEmailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-mailer-GetEmailResponse"></a>

### GetEmailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.Email](#resources-mailer-Email) |  |  |






<a name="services-mailer-GetEmailSettingsRequest"></a>

### GetEmailSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-GetEmailSettingsResponse"></a>

### GetEmailSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.EmailSettings](#resources-mailer-EmailSettings) |  |  |






<a name="services-mailer-GetTemplateRequest"></a>

### GetTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `template_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-GetTemplateResponse"></a>

### GetTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.Template](#resources-mailer-Template) |  |  |






<a name="services-mailer-GetThreadRequest"></a>

### GetThreadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-GetThreadResponse"></a>

### GetThreadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.Thread](#resources-mailer-Thread) |  |  |






<a name="services-mailer-GetThreadStateRequest"></a>

### GetThreadStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-GetThreadStateResponse"></a>

### GetThreadStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.ThreadState](#resources-mailer-ThreadState) |  |  |






<a name="services-mailer-ListEmailsRequest"></a>

### ListEmailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `all` | [bool](#bool) | optional | Search params |






<a name="services-mailer-ListEmailsResponse"></a>

### ListEmailsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `emails` | [resources.mailer.Email](#resources-mailer-Email) | repeated |  |






<a name="services-mailer-ListTemplatesRequest"></a>

### ListTemplatesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |






<a name="services-mailer-ListTemplatesResponse"></a>

### ListTemplatesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `templates` | [resources.mailer.Template](#resources-mailer-Template) | repeated |  |






<a name="services-mailer-ListThreadMessagesRequest"></a>

### ListThreadMessagesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |
| `after` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |






<a name="services-mailer-ListThreadMessagesResponse"></a>

### ListThreadMessagesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `messages` | [resources.mailer.Message](#resources-mailer-Message) | repeated |  |






<a name="services-mailer-ListThreadsRequest"></a>

### ListThreadsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `email_ids` | [uint64](#uint64) | repeated | Search params |
| `unread` | [bool](#bool) | optional |  |
| `archived` | [bool](#bool) | optional |  |






<a name="services-mailer-ListThreadsResponse"></a>

### ListThreadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `threads` | [resources.mailer.Thread](#resources-mailer-Thread) | repeated |  |






<a name="services-mailer-PostMessageRequest"></a>

### PostMessageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [resources.mailer.Message](#resources-mailer-Message) |  |  |
| `recipients` | [string](#string) | repeated | @sanitize: method=StripTags |






<a name="services-mailer-PostMessageResponse"></a>

### PostMessageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [resources.mailer.Message](#resources-mailer-Message) |  |  |






<a name="services-mailer-SearchThreadsRequest"></a>

### SearchThreadsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `search` | [string](#string) |  | Search params |






<a name="services-mailer-SearchThreadsResponse"></a>

### SearchThreadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `messages` | [resources.mailer.Message](#resources-mailer-Message) | repeated |  |






<a name="services-mailer-SetEmailSettingsRequest"></a>

### SetEmailSettingsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.EmailSettings](#resources-mailer-EmailSettings) |  |  |






<a name="services-mailer-SetEmailSettingsResponse"></a>

### SetEmailSettingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.EmailSettings](#resources-mailer-EmailSettings) |  |  |






<a name="services-mailer-SetThreadStateRequest"></a>

### SetThreadStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.ThreadState](#resources-mailer-ThreadState) |  |  |






<a name="services-mailer-SetThreadStateResponse"></a>

### SetThreadStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.ThreadState](#resources-mailer-ThreadState) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-mailer-MailerService"></a>

### MailerService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListEmails` | [ListEmailsRequest](#services-mailer-ListEmailsRequest) | [ListEmailsResponse](#services-mailer-ListEmailsResponse) | @perm |
| `GetEmail` | [GetEmailRequest](#services-mailer-GetEmailRequest) | [GetEmailResponse](#services-mailer-GetEmailResponse) | @perm: Name=ListEmails |
| `CreateOrUpdateEmail` | [CreateOrUpdateEmailRequest](#services-mailer-CreateOrUpdateEmailRequest) | [CreateOrUpdateEmailResponse](#services-mailer-CreateOrUpdateEmailResponse) | @perm: Attrs=Fields/StringList:[]string{"Job"} |
| `DeleteEmail` | [DeleteEmailRequest](#services-mailer-DeleteEmailRequest) | [DeleteEmailResponse](#services-mailer-DeleteEmailResponse) | @perm |
| `GetEmailProposals` | [GetEmailProposalsRequest](#services-mailer-GetEmailProposalsRequest) | [GetEmailProposalsResponse](#services-mailer-GetEmailProposalsResponse) | @perm: Name=ListEmails |
| `ListTemplates` | [ListTemplatesRequest](#services-mailer-ListTemplatesRequest) | [ListTemplatesResponse](#services-mailer-ListTemplatesResponse) | @perm: Name=ListEmails |
| `GetTemplate` | [GetTemplateRequest](#services-mailer-GetTemplateRequest) | [GetTemplateResponse](#services-mailer-GetTemplateResponse) | @perm: Name=ListEmails |
| `CreateOrUpdateTemplate` | [CreateOrUpdateTemplateRequest](#services-mailer-CreateOrUpdateTemplateRequest) | [CreateOrUpdateTemplateResponse](#services-mailer-CreateOrUpdateTemplateResponse) | @perm: Name=ListEmails |
| `DeleteTemplate` | [DeleteTemplateRequest](#services-mailer-DeleteTemplateRequest) | [DeleteTemplateResponse](#services-mailer-DeleteTemplateResponse) | @perm: Name=ListEmails |
| `ListThreads` | [ListThreadsRequest](#services-mailer-ListThreadsRequest) | [ListThreadsResponse](#services-mailer-ListThreadsResponse) | @perm: Name=ListEmails |
| `GetThread` | [GetThreadRequest](#services-mailer-GetThreadRequest) | [GetThreadResponse](#services-mailer-GetThreadResponse) | @perm: Name=ListEmails |
| `CreateThread` | [CreateThreadRequest](#services-mailer-CreateThreadRequest) | [CreateThreadResponse](#services-mailer-CreateThreadResponse) | @perm: Name=ListEmails |
| `DeleteThread` | [DeleteThreadRequest](#services-mailer-DeleteThreadRequest) | [DeleteThreadResponse](#services-mailer-DeleteThreadResponse) | @perm: Name=Superuser |
| `GetThreadState` | [GetThreadStateRequest](#services-mailer-GetThreadStateRequest) | [GetThreadStateResponse](#services-mailer-GetThreadStateResponse) | @perm: Name=ListEmails |
| `SetThreadState` | [SetThreadStateRequest](#services-mailer-SetThreadStateRequest) | [SetThreadStateResponse](#services-mailer-SetThreadStateResponse) | @perm: Name=ListEmails |
| `SearchThreads` | [SearchThreadsRequest](#services-mailer-SearchThreadsRequest) | [SearchThreadsResponse](#services-mailer-SearchThreadsResponse) | @perm: Name=ListEmails |
| `ListThreadMessages` | [ListThreadMessagesRequest](#services-mailer-ListThreadMessagesRequest) | [ListThreadMessagesResponse](#services-mailer-ListThreadMessagesResponse) | @perm: Name=ListEmails |
| `PostMessage` | [PostMessageRequest](#services-mailer-PostMessageRequest) | [PostMessageResponse](#services-mailer-PostMessageResponse) | @perm: Name=ListEmails |
| `DeleteMessage` | [DeleteMessageRequest](#services-mailer-DeleteMessageRequest) | [DeleteMessageResponse](#services-mailer-DeleteMessageResponse) | @perm: Name=Superuser |
| `GetEmailSettings` | [GetEmailSettingsRequest](#services-mailer-GetEmailSettingsRequest) | [GetEmailSettingsResponse](#services-mailer-GetEmailSettingsResponse) | @perm: Name=ListEmails |
| `SetEmailSettings` | [SetEmailSettingsRequest](#services-mailer-SetEmailSettingsRequest) | [SetEmailSettingsResponse](#services-mailer-SetEmailSettingsResponse) | @perm: Name=ListEmails |

 <!-- end services -->



<a name="services_wiki_collab-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/wiki/collab.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-wiki-CollabService"></a>

### CollabService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `JoinRoom` | [.resources.collab.ClientPacket](#resources-collab-ClientPacket) stream | [.resources.collab.ServerPacket](#resources-collab-ServerPacket) stream | @perm: Name=wiki.WikiService/ListPages |

 <!-- end services -->



<a name="services_wiki_wiki-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/wiki/wiki.proto



<a name="services-wiki-CreatePageRequest"></a>

### CreatePageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) | optional |  |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |






<a name="services-wiki-CreatePageResponse"></a>

### CreatePageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `id` | [uint64](#uint64) |  |  |






<a name="services-wiki-DeletePageRequest"></a>

### DeletePageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-wiki-DeletePageResponse"></a>

### DeletePageResponse







<a name="services-wiki-GetPageRequest"></a>

### GetPageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-wiki-GetPageResponse"></a>

### GetPageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resources-wiki-Page) |  |  |






<a name="services-wiki-ListPageActivityRequest"></a>

### ListPageActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `page_id` | [uint64](#uint64) |  |  |






<a name="services-wiki-ListPageActivityResponse"></a>

### ListPageActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `activity` | [resources.wiki.PageActivity](#resources-wiki-PageActivity) | repeated |  |






<a name="services-wiki-ListPagesRequest"></a>

### ListPagesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `job` | [string](#string) | optional | Search params |
| `root_only` | [bool](#bool) | optional |  |
| `search` | [string](#string) | optional |  |






<a name="services-wiki-ListPagesResponse"></a>

### ListPagesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `pages` | [resources.wiki.PageShort](#resources-wiki-PageShort) | repeated |  |






<a name="services-wiki-UpdatePageRequest"></a>

### UpdatePageRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resources-wiki-Page) |  |  |






<a name="services-wiki-UpdatePageResponse"></a>

### UpdatePageResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resources-wiki-Page) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-wiki-WikiService"></a>

### WikiService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListPages` | [ListPagesRequest](#services-wiki-ListPagesRequest) | [ListPagesResponse](#services-wiki-ListPagesResponse) | @perm |
| `GetPage` | [GetPageRequest](#services-wiki-GetPageRequest) | [GetPageResponse](#services-wiki-GetPageResponse) | @perm: Name=ListPages |
| `CreatePage` | [CreatePageRequest](#services-wiki-CreatePageRequest) | [CreatePageResponse](#services-wiki-CreatePageResponse) | @perm: Name=UpdatePage |
| `UpdatePage` | [UpdatePageRequest](#services-wiki-UpdatePageRequest) | [UpdatePageResponse](#services-wiki-UpdatePageResponse) | @perm: Attrs=Fields/StringList:[]string{"Public"} |
| `DeletePage` | [DeletePageRequest](#services-wiki-DeletePageRequest) | [DeletePageResponse](#services-wiki-DeletePageResponse) | @perm |
| `ListPageActivity` | [ListPageActivityRequest](#services-wiki-ListPageActivityRequest) | [ListPageActivityResponse](#services-wiki-ListPageActivityResponse) | @perm |
| `UploadFile` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=UpdatePage |

 <!-- end services -->



<a name="services_sync_sync-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/sync/sync.proto



<a name="services-sync-AddActivityRequest"></a>

### AddActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_oauth2` | [resources.sync.UserOAuth2Conn](#resources-sync-UserOAuth2Conn) |  |  |
| `dispatch` | [resources.centrum.Dispatch](#resources-centrum-Dispatch) |  |  |
| `user_activity` | [resources.users.UserActivity](#resources-users-UserActivity) |  | User activity |
| `user_props` | [resources.sync.UserProps](#resources-sync-UserProps) |  | Setting props will cause activity to be created automtically |
| `colleague_activity` | [resources.jobs.ColleagueActivity](#resources-jobs-ColleagueActivity) |  | Jobs user activity |
| `colleague_props` | [resources.sync.ColleagueProps](#resources-sync-ColleagueProps) |  | Setting props will cause activity to be created automtically |
| `job_timeclock` | [resources.sync.TimeclockUpdate](#resources-sync-TimeclockUpdate) |  | Timeclock user entry |
| `user_update` | [resources.sync.UserUpdate](#resources-sync-UserUpdate) |  | User/Char info updates that aren't tracked by activity (yet) |






<a name="services-sync-AddActivityResponse"></a>

### AddActivityResponse







<a name="services-sync-DeleteDataRequest"></a>

### DeleteDataRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.sync.DeleteUsers](#resources-sync-DeleteUsers) |  |  |
| `vehicles` | [resources.sync.DeleteVehicles](#resources-sync-DeleteVehicles) |  |  |






<a name="services-sync-DeleteDataResponse"></a>

### DeleteDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |






<a name="services-sync-GetStatusRequest"></a>

### GetStatusRequest







<a name="services-sync-GetStatusResponse"></a>

### GetStatusResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.DataStatus](#resources-sync-DataStatus) |  |  |
| `licenses` | [resources.sync.DataStatus](#resources-sync-DataStatus) |  |  |
| `users` | [resources.sync.DataStatus](#resources-sync-DataStatus) |  |  |
| `vehicles` | [resources.sync.DataStatus](#resources-sync-DataStatus) |  |  |






<a name="services-sync-RegisterAccountRequest"></a>

### RegisterAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `reset_token` | [bool](#bool) |  |  |
| `last_char_id` | [int32](#int32) | optional |  |






<a name="services-sync-RegisterAccountResponse"></a>

### RegisterAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) | optional |  |
| `account_id` | [uint64](#uint64) | optional |  |
| `username` | [string](#string) | optional |  |






<a name="services-sync-SendDataRequest"></a>

### SendDataRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.DataJobs](#resources-sync-DataJobs) |  |  |
| `licenses` | [resources.sync.DataLicenses](#resources-sync-DataLicenses) |  |  |
| `users` | [resources.sync.DataUsers](#resources-sync-DataUsers) |  |  |
| `vehicles` | [resources.sync.DataVehicles](#resources-sync-DataVehicles) |  |  |
| `user_locations` | [resources.sync.DataUserLocations](#resources-sync-DataUserLocations) |  |  |






<a name="services-sync-SendDataResponse"></a>

### SendDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |






<a name="services-sync-StreamRequest"></a>

### StreamRequest







<a name="services-sync-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |






<a name="services-sync-TransferAccountRequest"></a>

### TransferAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old_license` | [string](#string) |  |  |
| `new_license` | [string](#string) |  |  |






<a name="services-sync-TransferAccountResponse"></a>

### TransferAccountResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-sync-SyncService"></a>

### SyncService
Sync Service handles the sync of data (e.g., users, jobs) to this FiveNet instance and API calls from the plugin (e.g., user activity, user props changes).

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetStatus` | [GetStatusRequest](#services-sync-GetStatusRequest) | [GetStatusResponse](#services-sync-GetStatusResponse) | Get basic "sync state" from server side (currently simply the count of records on the server side). |
| `AddActivity` | [AddActivityRequest](#services-sync-AddActivityRequest) | [AddActivityResponse](#services-sync-AddActivityResponse) | For "tracking" activity such as "user received traffic infraction points", timeclock entries, etc. |
| `RegisterAccount` | [RegisterAccountRequest](#services-sync-RegisterAccountRequest) | [RegisterAccountResponse](#services-sync-RegisterAccountResponse) | Get registration token for a new user account or return the account id and username, for a given identifier/license. |
| `TransferAccount` | [TransferAccountRequest](#services-sync-TransferAccountRequest) | [TransferAccountResponse](#services-sync-TransferAccountResponse) | Transfer account from one license to another |
| `SendData` | [SendDataRequest](#services-sync-SendDataRequest) | [SendDataResponse](#services-sync-SendDataResponse) | DBSync's method of sending (mass) data to the FiveNet server for storing. |
| `DeleteData` | [DeleteDataRequest](#services-sync-DeleteDataRequest) | [DeleteDataResponse](#services-sync-DeleteDataResponse) | Way for the gameserver to delete certain data as well |
| `Stream` | [StreamRequest](#services-sync-StreamRequest) | [StreamResponse](#services-sync-StreamResponse) stream | Used for the server to stream events to the dbsync (e.g., "refresh" of user/char data) |

 <!-- end services -->



<a name="services_citizens_citizens-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/citizens/citizens.proto



<a name="services-citizens-DeleteAvatarRequest"></a>

### DeleteAvatarRequest







<a name="services-citizens-DeleteAvatarResponse"></a>

### DeleteAvatarResponse







<a name="services-citizens-DeleteMugshotRequest"></a>

### DeleteMugshotRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `reason` | [string](#string) |  | @sanitize |






<a name="services-citizens-DeleteMugshotResponse"></a>

### DeleteMugshotResponse







<a name="services-citizens-GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `info_only` | [bool](#bool) | optional |  |






<a name="services-citizens-GetUserResponse"></a>

### GetUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user` | [resources.users.User](#resources-users-User) |  |  |






<a name="services-citizens-ListCitizensRequest"></a>

### ListCitizensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `search` | [string](#string) |  | Search params |
| `wanted` | [bool](#bool) | optional |  |
| `phone_number` | [string](#string) | optional |  |
| `traffic_infraction_points` | [uint32](#uint32) | optional |  |
| `dateofbirth` | [string](#string) | optional |  |
| `open_fines` | [uint64](#uint64) | optional |  |






<a name="services-citizens-ListCitizensResponse"></a>

### ListCitizensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `users` | [resources.users.User](#resources-users-User) | repeated |  |






<a name="services-citizens-ListUserActivityRequest"></a>

### ListUserActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `user_id` | [int32](#int32) |  | Search params |
| `types` | [resources.users.UserActivityType](#resources-users-UserActivityType) | repeated |  |






<a name="services-citizens-ListUserActivityResponse"></a>

### ListUserActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `activity` | [resources.users.UserActivity](#resources-users-UserActivity) | repeated |  |






<a name="services-citizens-ManageLabelsRequest"></a>

### ManageLabelsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.Label](#resources-users-Label) | repeated |  |






<a name="services-citizens-ManageLabelsResponse"></a>

### ManageLabelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.Label](#resources-users-Label) | repeated |  |






<a name="services-citizens-SetUserPropsRequest"></a>

### SetUserPropsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.users.UserProps](#resources-users-UserProps) |  |  |
| `reason` | [string](#string) |  | @sanitize |






<a name="services-citizens-SetUserPropsResponse"></a>

### SetUserPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.users.UserProps](#resources-users-UserProps) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-citizens-CitizensService"></a>

### CitizensService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCitizens` | [ListCitizensRequest](#services-citizens-ListCitizensRequest) | [ListCitizensResponse](#services-citizens-ListCitizensResponse) | @perm: Attrs=Fields/StringList:[]string{"PhoneNumber", "Licenses", "UserProps.Wanted", "UserProps.Job", "UserProps.TrafficInfractionPoints", "UserProps.OpenFines", "UserProps.BloodType", "UserProps.Mugshot", "UserProps.Labels", "UserProps.Email"} |
| `GetUser` | [GetUserRequest](#services-citizens-GetUserRequest) | [GetUserResponse](#services-citizens-GetUserResponse) | @perm: Attrs=Jobs/JobGradeList |
| `ListUserActivity` | [ListUserActivityRequest](#services-citizens-ListUserActivityRequest) | [ListUserActivityResponse](#services-citizens-ListUserActivityResponse) | @perm: Attrs=Fields/StringList:[]string{"SourceUser", "Own"} |
| `SetUserProps` | [SetUserPropsRequest](#services-citizens-SetUserPropsRequest) | [SetUserPropsResponse](#services-citizens-SetUserPropsResponse) | @perm: Attrs=Fields/StringList:[]string{"Wanted", "Job", "TrafficInfractionPoints", "Mugshot", "Labels"} |
| `UploadAvatar` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=Any |
| `DeleteAvatar` | [DeleteAvatarRequest](#services-citizens-DeleteAvatarRequest) | [DeleteAvatarResponse](#services-citizens-DeleteAvatarResponse) | @perm: Name=Any |
| `UploadMugshot` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=SetUserProps |
| `DeleteMugshot` | [DeleteMugshotRequest](#services-citizens-DeleteMugshotRequest) | [DeleteMugshotResponse](#services-citizens-DeleteMugshotResponse) | @perm: Name=SetUserProps |
| `ManageLabels` | [ManageLabelsRequest](#services-citizens-ManageLabelsRequest) | [ManageLabelsResponse](#services-citizens-ManageLabelsResponse) | @perm |

 <!-- end services -->



<a name="services_documents_collab-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/documents/collab.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-documents-CollabService"></a>

### CollabService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `JoinRoom` | [.resources.collab.ClientPacket](#resources-collab-ClientPacket) stream | [.resources.collab.ServerPacket](#resources-collab-ServerPacket) stream | @perm: Name=documents.DocumentsService/ListDocuments |

 <!-- end services -->



<a name="services_documents_documents-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/documents/documents.proto



<a name="services-documents-AddDocumentReferenceRequest"></a>

### AddDocumentReferenceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reference` | [resources.documents.DocumentReference](#resources-documents-DocumentReference) |  |  |






<a name="services-documents-AddDocumentReferenceResponse"></a>

### AddDocumentReferenceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-AddDocumentRelationRequest"></a>

### AddDocumentRelationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relation` | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) |  |  |






<a name="services-documents-AddDocumentRelationResponse"></a>

### AddDocumentRelationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-ChangeDocumentOwnerRequest"></a>

### ChangeDocumentOwnerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `new_user_id` | [int32](#int32) | optional |  |






<a name="services-documents-ChangeDocumentOwnerResponse"></a>

### ChangeDocumentOwnerResponse







<a name="services-documents-CreateDocumentReqRequest"></a>

### CreateDocumentReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `request_type` | [resources.documents.DocActivityType](#resources-documents-DocActivityType) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `data` | [resources.documents.DocActivityData](#resources-documents-DocActivityData) | optional |  |






<a name="services-documents-CreateDocumentReqResponse"></a>

### CreateDocumentReqResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.documents.DocRequest](#resources-documents-DocRequest) |  |  |






<a name="services-documents-CreateDocumentRequest"></a>

### CreateDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |
| `template_id` | [uint64](#uint64) | optional |  |
| `template_data` | [resources.documents.TemplateData](#resources-documents-TemplateData) | optional |  |






<a name="services-documents-CreateDocumentResponse"></a>

### CreateDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-CreateOrUpdateCategoryRequest"></a>

### CreateOrUpdateCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [resources.documents.Category](#resources-documents-Category) |  |  |






<a name="services-documents-CreateOrUpdateCategoryResponse"></a>

### CreateOrUpdateCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [resources.documents.Category](#resources-documents-Category) |  |  |






<a name="services-documents-CreateTemplateRequest"></a>

### CreateTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resources-documents-Template) |  |  |






<a name="services-documents-CreateTemplateResponse"></a>

### CreateTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-DeleteCategoryRequest"></a>

### DeleteCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-DeleteCategoryResponse"></a>

### DeleteCategoryResponse







<a name="services-documents-DeleteCommentRequest"></a>

### DeleteCommentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment_id` | [uint64](#uint64) |  |  |






<a name="services-documents-DeleteCommentResponse"></a>

### DeleteCommentResponse







<a name="services-documents-DeleteDocumentReqRequest"></a>

### DeleteDocumentReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request_id` | [uint64](#uint64) |  |  |






<a name="services-documents-DeleteDocumentReqResponse"></a>

### DeleteDocumentReqResponse







<a name="services-documents-DeleteDocumentRequest"></a>

### DeleteDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `reason` | [string](#string) | optional | @sanitize: method=StripTags |






<a name="services-documents-DeleteDocumentResponse"></a>

### DeleteDocumentResponse







<a name="services-documents-DeleteTemplateRequest"></a>

### DeleteTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-DeleteTemplateResponse"></a>

### DeleteTemplateResponse







<a name="services-documents-EditCommentRequest"></a>

### EditCommentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-documents-EditCommentResponse"></a>

### EditCommentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-documents-GetCommentsRequest"></a>

### GetCommentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `document_id` | [uint64](#uint64) |  |  |






<a name="services-documents-GetCommentsResponse"></a>

### GetCommentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `comments` | [resources.documents.Comment](#resources-documents-Comment) | repeated |  |






<a name="services-documents-GetDocumentAccessRequest"></a>

### GetDocumentAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |






<a name="services-documents-GetDocumentAccessResponse"></a>

### GetDocumentAccessResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-documents-GetDocumentReferencesRequest"></a>

### GetDocumentReferencesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |






<a name="services-documents-GetDocumentReferencesResponse"></a>

### GetDocumentReferencesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `references` | [resources.documents.DocumentReference](#resources-documents-DocumentReference) | repeated | @gotags: alias:"reference" |






<a name="services-documents-GetDocumentRelationsRequest"></a>

### GetDocumentRelationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |






<a name="services-documents-GetDocumentRelationsResponse"></a>

### GetDocumentRelationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relations` | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) | repeated | @gotags: alias:"relation" |






<a name="services-documents-GetDocumentRequest"></a>

### GetDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `info_only` | [bool](#bool) | optional |  |






<a name="services-documents-GetDocumentResponse"></a>

### GetDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [resources.documents.Document](#resources-documents-Document) |  |  |
| `access` | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-documents-GetTemplateRequest"></a>

### GetTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template_id` | [uint64](#uint64) |  |  |
| `data` | [resources.documents.TemplateData](#resources-documents-TemplateData) | optional |  |
| `render` | [bool](#bool) | optional |  |






<a name="services-documents-GetTemplateResponse"></a>

### GetTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resources-documents-Template) |  |  |
| `rendered` | [bool](#bool) |  |  |






<a name="services-documents-ListCategoriesRequest"></a>

### ListCategoriesRequest







<a name="services-documents-ListCategoriesResponse"></a>

### ListCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `categories` | [resources.documents.Category](#resources-documents-Category) | repeated |  |






<a name="services-documents-ListDocumentActivityRequest"></a>

### ListDocumentActivityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `activity_types` | [resources.documents.DocActivityType](#resources-documents-DocActivityType) | repeated | Search params |






<a name="services-documents-ListDocumentActivityResponse"></a>

### ListDocumentActivityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `activity` | [resources.documents.DocActivity](#resources-documents-DocActivity) | repeated |  |






<a name="services-documents-ListDocumentPinsRequest"></a>

### ListDocumentPinsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `personal` | [bool](#bool) | optional | Search params If true, only personal pins are returned |






<a name="services-documents-ListDocumentPinsResponse"></a>

### ListDocumentPinsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `documents` | [resources.documents.DocumentShort](#resources-documents-DocumentShort) | repeated |  |






<a name="services-documents-ListDocumentReqsRequest"></a>

### ListDocumentReqsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `document_id` | [uint64](#uint64) |  |  |






<a name="services-documents-ListDocumentReqsResponse"></a>

### ListDocumentReqsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `requests` | [resources.documents.DocRequest](#resources-documents-DocRequest) | repeated |  |






<a name="services-documents-ListDocumentsRequest"></a>

### ListDocumentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `search` | [string](#string) | optional | Search params |
| `category_ids` | [uint64](#uint64) | repeated |  |
| `creator_ids` | [int32](#int32) | repeated |  |
| `from` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `to` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `closed` | [bool](#bool) | optional |  |
| `document_ids` | [uint64](#uint64) | repeated |  |
| `only_drafts` | [bool](#bool) | optional | Controls inclusion of drafts in the result: - unset/null: include all documents (drafts and non-drafts) - false: only non-draft documents - true: only draft documents |






<a name="services-documents-ListDocumentsResponse"></a>

### ListDocumentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `documents` | [resources.documents.DocumentShort](#resources-documents-DocumentShort) | repeated |  |






<a name="services-documents-ListTemplatesRequest"></a>

### ListTemplatesRequest







<a name="services-documents-ListTemplatesResponse"></a>

### ListTemplatesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `templates` | [resources.documents.TemplateShort](#resources-documents-TemplateShort) | repeated |  |






<a name="services-documents-ListUserDocumentsRequest"></a>

### ListUserDocumentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `relations` | [resources.documents.DocRelation](#resources-documents-DocRelation) | repeated |  |
| `closed` | [bool](#bool) | optional |  |






<a name="services-documents-ListUserDocumentsResponse"></a>

### ListUserDocumentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `relations` | [resources.documents.DocumentRelation](#resources-documents-DocumentRelation) | repeated |  |






<a name="services-documents-PostCommentRequest"></a>

### PostCommentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-documents-PostCommentResponse"></a>

### PostCommentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resources-documents-Comment) |  |  |






<a name="services-documents-RemoveDocumentReferenceRequest"></a>

### RemoveDocumentReferenceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-RemoveDocumentReferenceResponse"></a>

### RemoveDocumentReferenceResponse







<a name="services-documents-RemoveDocumentRelationRequest"></a>

### RemoveDocumentRelationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-documents-RemoveDocumentRelationResponse"></a>

### RemoveDocumentRelationResponse







<a name="services-documents-SetDocumentAccessRequest"></a>

### SetDocumentAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `access` | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) |  |  |






<a name="services-documents-SetDocumentAccessResponse"></a>

### SetDocumentAccessResponse







<a name="services-documents-SetDocumentReminderRequest"></a>

### SetDocumentReminderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `reminder_time` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `message` | [string](#string) | optional | @sanitize: method=StripTags |






<a name="services-documents-SetDocumentReminderResponse"></a>

### SetDocumentReminderResponse







<a name="services-documents-ToggleDocumentPinRequest"></a>

### ToggleDocumentPinRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `state` | [bool](#bool) |  |  |
| `personal` | [bool](#bool) | optional | If true, the pin is personal and not shared with other job members |






<a name="services-documents-ToggleDocumentPinResponse"></a>

### ToggleDocumentPinResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pin` | [resources.documents.DocumentPin](#resources-documents-DocumentPin) | optional | @gotags: alias:"pin" |






<a name="services-documents-ToggleDocumentRequest"></a>

### ToggleDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `closed` | [bool](#bool) |  |  |






<a name="services-documents-ToggleDocumentResponse"></a>

### ToggleDocumentResponse







<a name="services-documents-UpdateDocumentReqRequest"></a>

### UpdateDocumentReqRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `request_id` | [uint64](#uint64) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `data` | [resources.documents.DocActivityData](#resources-documents-DocActivityData) | optional |  |
| `accepted` | [bool](#bool) |  |  |






<a name="services-documents-UpdateDocumentReqResponse"></a>

### UpdateDocumentReqResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.documents.DocRequest](#resources-documents-DocRequest) |  |  |






<a name="services-documents-UpdateDocumentRequest"></a>

### UpdateDocumentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  | @gotags: alias:"id" |
| `category_id` | [uint64](#uint64) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags

@gotags: alias:"title" |
| `content` | [resources.common.content.Content](#resources-common-content-Content) |  |  |
| `content_type` | [resources.common.content.ContentType](#resources-common-content-ContentType) |  |  |
| `data` | [string](#string) | optional |  |
| `state` | [string](#string) |  | @sanitize |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `access` | [resources.documents.DocumentAccess](#resources-documents-DocumentAccess) | optional |  |
| `files` | [resources.file.File](#resources-file-File) | repeated | @gotags: alias:"files" |






<a name="services-documents-UpdateDocumentResponse"></a>

### UpdateDocumentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [resources.documents.Document](#resources-documents-Document) |  |  |






<a name="services-documents-UpdateTemplateRequest"></a>

### UpdateTemplateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resources-documents-Template) |  |  |






<a name="services-documents-UpdateTemplateResponse"></a>

### UpdateTemplateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resources-documents-Template) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-documents-DocumentsService"></a>

### DocumentsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTemplates` | [ListTemplatesRequest](#services-documents-ListTemplatesRequest) | [ListTemplatesResponse](#services-documents-ListTemplatesResponse) | @perm |
| `GetTemplate` | [GetTemplateRequest](#services-documents-GetTemplateRequest) | [GetTemplateResponse](#services-documents-GetTemplateResponse) | @perm: Name=ListTemplates |
| `CreateTemplate` | [CreateTemplateRequest](#services-documents-CreateTemplateRequest) | [CreateTemplateResponse](#services-documents-CreateTemplateResponse) | @perm |
| `UpdateTemplate` | [UpdateTemplateRequest](#services-documents-UpdateTemplateRequest) | [UpdateTemplateResponse](#services-documents-UpdateTemplateResponse) | @perm: Name=CreateTemplate |
| `DeleteTemplate` | [DeleteTemplateRequest](#services-documents-DeleteTemplateRequest) | [DeleteTemplateResponse](#services-documents-DeleteTemplateResponse) | @perm |
| `ListDocuments` | [ListDocumentsRequest](#services-documents-ListDocumentsRequest) | [ListDocumentsResponse](#services-documents-ListDocumentsResponse) | @perm |
| `GetDocument` | [GetDocumentRequest](#services-documents-GetDocumentRequest) | [GetDocumentResponse](#services-documents-GetDocumentResponse) | @perm: Name=ListDocuments |
| `CreateDocument` | [CreateDocumentRequest](#services-documents-CreateDocumentRequest) | [CreateDocumentResponse](#services-documents-CreateDocumentResponse) | @perm: Name=UpdateDocument |
| `UpdateDocument` | [UpdateDocumentRequest](#services-documents-UpdateDocumentRequest) | [UpdateDocumentResponse](#services-documents-UpdateDocumentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `DeleteDocument` | [DeleteDocumentRequest](#services-documents-DeleteDocumentRequest) | [DeleteDocumentResponse](#services-documents-DeleteDocumentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `ToggleDocument` | [ToggleDocumentRequest](#services-documents-ToggleDocumentRequest) | [ToggleDocumentResponse](#services-documents-ToggleDocumentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `ChangeDocumentOwner` | [ChangeDocumentOwnerRequest](#services-documents-ChangeDocumentOwnerRequest) | [ChangeDocumentOwnerResponse](#services-documents-ChangeDocumentOwnerResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `GetDocumentReferences` | [GetDocumentReferencesRequest](#services-documents-GetDocumentReferencesRequest) | [GetDocumentReferencesResponse](#services-documents-GetDocumentReferencesResponse) | @perm: Name=ListDocuments |
| `GetDocumentRelations` | [GetDocumentRelationsRequest](#services-documents-GetDocumentRelationsRequest) | [GetDocumentRelationsResponse](#services-documents-GetDocumentRelationsResponse) | @perm: Name=ListDocuments |
| `AddDocumentReference` | [AddDocumentReferenceRequest](#services-documents-AddDocumentReferenceRequest) | [AddDocumentReferenceResponse](#services-documents-AddDocumentReferenceResponse) | @perm |
| `RemoveDocumentReference` | [RemoveDocumentReferenceRequest](#services-documents-RemoveDocumentReferenceRequest) | [RemoveDocumentReferenceResponse](#services-documents-RemoveDocumentReferenceResponse) | @perm: Name=AddDocumentReference |
| `AddDocumentRelation` | [AddDocumentRelationRequest](#services-documents-AddDocumentRelationRequest) | [AddDocumentRelationResponse](#services-documents-AddDocumentRelationResponse) | @perm |
| `RemoveDocumentRelation` | [RemoveDocumentRelationRequest](#services-documents-RemoveDocumentRelationRequest) | [RemoveDocumentRelationResponse](#services-documents-RemoveDocumentRelationResponse) | @perm: Name=AddDocumentRelation |
| `GetComments` | [GetCommentsRequest](#services-documents-GetCommentsRequest) | [GetCommentsResponse](#services-documents-GetCommentsResponse) | @perm: Name=ListDocuments |
| `PostComment` | [PostCommentRequest](#services-documents-PostCommentRequest) | [PostCommentResponse](#services-documents-PostCommentResponse) | @perm: Name=ListDocuments |
| `EditComment` | [EditCommentRequest](#services-documents-EditCommentRequest) | [EditCommentResponse](#services-documents-EditCommentResponse) | @perm: Name=ListDocuments |
| `DeleteComment` | [DeleteCommentRequest](#services-documents-DeleteCommentRequest) | [DeleteCommentResponse](#services-documents-DeleteCommentResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `GetDocumentAccess` | [GetDocumentAccessRequest](#services-documents-GetDocumentAccessRequest) | [GetDocumentAccessResponse](#services-documents-GetDocumentAccessResponse) | @perm: Name=ListDocuments |
| `SetDocumentAccess` | [SetDocumentAccessRequest](#services-documents-SetDocumentAccessRequest) | [SetDocumentAccessResponse](#services-documents-SetDocumentAccessResponse) | @perm: Name=UpdateDocument |
| `ListDocumentActivity` | [ListDocumentActivityRequest](#services-documents-ListDocumentActivityRequest) | [ListDocumentActivityResponse](#services-documents-ListDocumentActivityResponse) | @perm |
| `ListDocumentReqs` | [ListDocumentReqsRequest](#services-documents-ListDocumentReqsRequest) | [ListDocumentReqsResponse](#services-documents-ListDocumentReqsResponse) | @perm |
| `CreateDocumentReq` | [CreateDocumentReqRequest](#services-documents-CreateDocumentReqRequest) | [CreateDocumentReqResponse](#services-documents-CreateDocumentReqResponse) | @perm: Attrs=Types/StringList:[]string{"Access", "Closure", "Update", "Deletion", "OwnerChange"} |
| `UpdateDocumentReq` | [UpdateDocumentReqRequest](#services-documents-UpdateDocumentReqRequest) | [UpdateDocumentReqResponse](#services-documents-UpdateDocumentReqResponse) | @perm: Name=CreateDocumentReq |
| `DeleteDocumentReq` | [DeleteDocumentReqRequest](#services-documents-DeleteDocumentReqRequest) | [DeleteDocumentReqResponse](#services-documents-DeleteDocumentReqResponse) | @perm |
| `ListUserDocuments` | [ListUserDocumentsRequest](#services-documents-ListUserDocumentsRequest) | [ListUserDocumentsResponse](#services-documents-ListUserDocumentsResponse) | @perm |
| `ListCategories` | [ListCategoriesRequest](#services-documents-ListCategoriesRequest) | [ListCategoriesResponse](#services-documents-ListCategoriesResponse) | @perm |
| `CreateOrUpdateCategory` | [CreateOrUpdateCategoryRequest](#services-documents-CreateOrUpdateCategoryRequest) | [CreateOrUpdateCategoryResponse](#services-documents-CreateOrUpdateCategoryResponse) | @perm |
| `DeleteCategory` | [DeleteCategoryRequest](#services-documents-DeleteCategoryRequest) | [DeleteCategoryResponse](#services-documents-DeleteCategoryResponse) | @perm |
| `ListDocumentPins` | [ListDocumentPinsRequest](#services-documents-ListDocumentPinsRequest) | [ListDocumentPinsResponse](#services-documents-ListDocumentPinsResponse) | @perm: Name=ListDocuments |
| `ToggleDocumentPin` | [ToggleDocumentPinRequest](#services-documents-ToggleDocumentPinRequest) | [ToggleDocumentPinResponse](#services-documents-ToggleDocumentPinResponse) | @perm: Attrs=Types/StringList:[]string{"JobWide"} |
| `SetDocumentReminder` | [SetDocumentReminderRequest](#services-documents-SetDocumentReminderRequest) | [SetDocumentReminderResponse](#services-documents-SetDocumentReminderResponse) | @perm |
| `UploadFile` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=UpdateDocument |

 <!-- end services -->



<a name="services_livemap_livemap-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/livemap/livemap.proto



<a name="services-livemap-CreateOrUpdateMarkerRequest"></a>

### CreateOrUpdateMarkerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `marker` | [resources.livemap.MarkerMarker](#resources-livemap-MarkerMarker) |  |  |






<a name="services-livemap-CreateOrUpdateMarkerResponse"></a>

### CreateOrUpdateMarkerResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `marker` | [resources.livemap.MarkerMarker](#resources-livemap-MarkerMarker) |  |  |






<a name="services-livemap-DeleteMarkerRequest"></a>

### DeleteMarkerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-livemap-DeleteMarkerResponse"></a>

### DeleteMarkerResponse







<a name="services-livemap-JobsList"></a>

### JobsList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.jobs.Job](#resources-jobs-Job) | repeated |  |
| `markers` | [resources.jobs.Job](#resources-jobs-Job) | repeated |  |






<a name="services-livemap-MarkerMarkersUpdates"></a>

### MarkerMarkersUpdates



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [resources.livemap.MarkerMarker](#resources-livemap-MarkerMarker) | repeated |  |
| `deleted` | [uint64](#uint64) | repeated |  |
| `part` | [int32](#int32) |  |  |
| `partial` | [bool](#bool) |  |  |






<a name="services-livemap-Snapshot"></a>

### Snapshot
A roll-up of the entire USERLOC bucket. Published every N seconds on `$KV.user_locations._snapshot` with the headers: Nats-Rollup: all KV-Operation: ROLLUP


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `markers` | [resources.livemap.UserMarker](#resources-livemap-UserMarker) | repeated | All currently-known user markers, already filtered for obsolete PURGE/DELETE events. |
| `generated_at` | [int64](#int64) |  | When the snapshot was generated (Unix epoch millis). |
| `snapshot_seq` | [uint64](#uint64) |  | Optional monotonic counter so a client can ignore older roll-ups that arrive out-of-order. |
| `schema_version` | [uint32](#uint32) |  | Version in case we extend the definition later (e.g. add units). |






<a name="services-livemap-StreamRequest"></a>

### StreamRequest







<a name="services-livemap-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_on_duty` | [bool](#bool) | optional |  |
| `jobs` | [JobsList](#services-livemap-JobsList) |  |  |
| `markers` | [MarkerMarkersUpdates](#services-livemap-MarkerMarkersUpdates) |  |  |
| `snapshot` | [Snapshot](#services-livemap-Snapshot) |  |  |
| `user_update` | [resources.livemap.UserMarker](#resources-livemap-UserMarker) |  |  |
| `user_delete` | [UserDelete](#services-livemap-UserDelete) |  |  |






<a name="services-livemap-UserDelete"></a>

### UserDelete



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int32](#int32) |  | The user ID of the user that was deleted. |
| `job` | [string](#string) |  | The job of the user that was deleted. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-livemap-LivemapService"></a>

### LivemapService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Stream` | [StreamRequest](#services-livemap-StreamRequest) | [StreamResponse](#services-livemap-StreamResponse) stream | @perm: Attrs=Markers/JobList|Players/JobGradeList |
| `CreateOrUpdateMarker` | [CreateOrUpdateMarkerRequest](#services-livemap-CreateOrUpdateMarkerRequest) | [CreateOrUpdateMarkerResponse](#services-livemap-CreateOrUpdateMarkerResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `DeleteMarker` | [DeleteMarkerRequest](#services-livemap-DeleteMarkerRequest) | [DeleteMarkerResponse](#services-livemap-DeleteMarkerResponse) | @perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |

 <!-- end services -->



<a name="services_settings_accounts-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/settings/accounts.proto



<a name="services-settings-DeleteAccountRequest"></a>

### DeleteAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-settings-DeleteAccountResponse"></a>

### DeleteAccountResponse







<a name="services-settings-DisconnectOAuth2ConnectionRequest"></a>

### DisconnectOAuth2ConnectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `providerName` | [string](#string) |  |  |






<a name="services-settings-DisconnectOAuth2ConnectionResponse"></a>

### DisconnectOAuth2ConnectionResponse







<a name="services-settings-ListAccountsRequest"></a>

### ListAccountsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `license` | [string](#string) | optional | Search params |
| `enabled` | [bool](#bool) | optional |  |






<a name="services-settings-ListAccountsResponse"></a>

### ListAccountsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `accounts` | [resources.accounts.Account](#resources-accounts-Account) | repeated |  |






<a name="services-settings-UpdateAccountRequest"></a>

### UpdateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `enabled` | [bool](#bool) | optional |  |
| `last_char` | [int32](#int32) | optional |  |






<a name="services-settings-UpdateAccountResponse"></a>

### UpdateAccountResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [resources.accounts.Account](#resources-accounts-Account) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-settings-AccountsService"></a>

### AccountsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListAccounts` | [ListAccountsRequest](#services-settings-ListAccountsRequest) | [ListAccountsResponse](#services-settings-ListAccountsResponse) | @perm: Name=Superuser |
| `UpdateAccount` | [UpdateAccountRequest](#services-settings-UpdateAccountRequest) | [UpdateAccountResponse](#services-settings-UpdateAccountResponse) | @perm: Name=Superuser |
| `DisconnectOAuth2Connection` | [DisconnectOAuth2ConnectionRequest](#services-settings-DisconnectOAuth2ConnectionRequest) | [DisconnectOAuth2ConnectionResponse](#services-settings-DisconnectOAuth2ConnectionResponse) | @perm: Name=Superuser |
| `DeleteAccount` | [DeleteAccountRequest](#services-settings-DeleteAccountRequest) | [DeleteAccountResponse](#services-settings-DeleteAccountResponse) | @perm: Name=Superuser |

 <!-- end services -->



<a name="services_settings_config-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/settings/config.proto



<a name="services-settings-GetAppConfigRequest"></a>

### GetAppConfigRequest







<a name="services-settings-GetAppConfigResponse"></a>

### GetAppConfigResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [resources.settings.AppConfig](#resources-settings-AppConfig) |  |  |






<a name="services-settings-UpdateAppConfigRequest"></a>

### UpdateAppConfigRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [resources.settings.AppConfig](#resources-settings-AppConfig) |  |  |






<a name="services-settings-UpdateAppConfigResponse"></a>

### UpdateAppConfigResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [resources.settings.AppConfig](#resources-settings-AppConfig) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-settings-ConfigService"></a>

### ConfigService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetAppConfig` | [GetAppConfigRequest](#services-settings-GetAppConfigRequest) | [GetAppConfigResponse](#services-settings-GetAppConfigResponse) | @perm: Name=Superuser |
| `UpdateAppConfig` | [UpdateAppConfigRequest](#services-settings-UpdateAppConfigRequest) | [UpdateAppConfigResponse](#services-settings-UpdateAppConfigResponse) | @perm: Name=Superuser |

 <!-- end services -->



<a name="services_settings_cron-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/settings/cron.proto



<a name="services-settings-ListCronjobsRequest"></a>

### ListCronjobsRequest







<a name="services-settings-ListCronjobsResponse"></a>

### ListCronjobsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.common.cron.Cronjob](#resources-common-cron-Cronjob) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-settings-CronService"></a>

### CronService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCronjobs` | [ListCronjobsRequest](#services-settings-ListCronjobsRequest) | [ListCronjobsResponse](#services-settings-ListCronjobsResponse) | @perm: Name=Superuser |

 <!-- end services -->



<a name="services_settings_laws-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/settings/laws.proto



<a name="services-settings-CreateOrUpdateLawBookRequest"></a>

### CreateOrUpdateLawBookRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `lawBook` | [resources.laws.LawBook](#resources-laws-LawBook) |  |  |






<a name="services-settings-CreateOrUpdateLawBookResponse"></a>

### CreateOrUpdateLawBookResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `lawBook` | [resources.laws.LawBook](#resources-laws-LawBook) |  |  |






<a name="services-settings-CreateOrUpdateLawRequest"></a>

### CreateOrUpdateLawRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `law` | [resources.laws.Law](#resources-laws-Law) |  |  |






<a name="services-settings-CreateOrUpdateLawResponse"></a>

### CreateOrUpdateLawResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `law` | [resources.laws.Law](#resources-laws-Law) |  |  |






<a name="services-settings-DeleteLawBookRequest"></a>

### DeleteLawBookRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-settings-DeleteLawBookResponse"></a>

### DeleteLawBookResponse







<a name="services-settings-DeleteLawRequest"></a>

### DeleteLawRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-settings-DeleteLawResponse"></a>

### DeleteLawResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-settings-LawsService"></a>

### LawsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `CreateOrUpdateLawBook` | [CreateOrUpdateLawBookRequest](#services-settings-CreateOrUpdateLawBookRequest) | [CreateOrUpdateLawBookResponse](#services-settings-CreateOrUpdateLawBookResponse) | @perm |
| `DeleteLawBook` | [DeleteLawBookRequest](#services-settings-DeleteLawBookRequest) | [DeleteLawBookResponse](#services-settings-DeleteLawBookResponse) | @perm |
| `CreateOrUpdateLaw` | [CreateOrUpdateLawRequest](#services-settings-CreateOrUpdateLawRequest) | [CreateOrUpdateLawResponse](#services-settings-CreateOrUpdateLawResponse) | @perm: Name=CreateOrUpdateLawBook |
| `DeleteLaw` | [DeleteLawRequest](#services-settings-DeleteLawRequest) | [DeleteLawResponse](#services-settings-DeleteLawResponse) | @perm: Name=DeleteLawBook |

 <!-- end services -->



<a name="services_settings_settings-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/settings/settings.proto



<a name="services-settings-AttrsUpdate"></a>

### AttrsUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_update` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |
| `to_remove` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-settings-CreateRoleRequest"></a>

### CreateRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |






<a name="services-settings-CreateRoleResponse"></a>

### CreateRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.Role](#resources-permissions-Role) |  |  |






<a name="services-settings-DeleteFactionRequest"></a>

### DeleteFactionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |






<a name="services-settings-DeleteFactionResponse"></a>

### DeleteFactionResponse







<a name="services-settings-DeleteJobLogoRequest"></a>

### DeleteJobLogoRequest







<a name="services-settings-DeleteJobLogoResponse"></a>

### DeleteJobLogoResponse







<a name="services-settings-DeleteRoleRequest"></a>

### DeleteRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-settings-DeleteRoleResponse"></a>

### DeleteRoleResponse







<a name="services-settings-GetAllPermissionsRequest"></a>

### GetAllPermissionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |






<a name="services-settings-GetAllPermissionsResponse"></a>

### GetAllPermissionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permissions` | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-settings-GetEffectivePermissionsRequest"></a>

### GetEffectivePermissionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |






<a name="services-settings-GetEffectivePermissionsResponse"></a>

### GetEffectivePermissionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.Role](#resources-permissions-Role) |  |  |
| `permissions` | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-settings-GetJobLimitsRequest"></a>

### GetJobLimitsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |






<a name="services-settings-GetJobLimitsResponse"></a>

### GetJobLimitsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `permissions` | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-settings-GetJobPropsRequest"></a>

### GetJobPropsRequest







<a name="services-settings-GetJobPropsResponse"></a>

### GetJobPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resources-jobs-JobProps) |  |  |






<a name="services-settings-GetPermissionsRequest"></a>

### GetPermissionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |






<a name="services-settings-GetPermissionsResponse"></a>

### GetPermissionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permissions` | [resources.permissions.Permission](#resources-permissions-Permission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resources-permissions-RoleAttribute) | repeated |  |






<a name="services-settings-GetRoleRequest"></a>

### GetRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |






<a name="services-settings-GetRoleResponse"></a>

### GetRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.Role](#resources-permissions-Role) |  |  |






<a name="services-settings-GetRolesRequest"></a>

### GetRolesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `lowest_rank` | [bool](#bool) | optional |  |






<a name="services-settings-GetRolesResponse"></a>

### GetRolesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `roles` | [resources.permissions.Role](#resources-permissions-Role) | repeated |  |






<a name="services-settings-ListDiscordChannelsRequest"></a>

### ListDiscordChannelsRequest







<a name="services-settings-ListDiscordChannelsResponse"></a>

### ListDiscordChannelsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [resources.discord.Channel](#resources-discord-Channel) | repeated |  |






<a name="services-settings-ListUserGuildsRequest"></a>

### ListUserGuildsRequest







<a name="services-settings-ListUserGuildsResponse"></a>

### ListUserGuildsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `guilds` | [resources.discord.Guild](#resources-discord-Guild) | repeated |  |






<a name="services-settings-PermsUpdate"></a>

### PermsUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_update` | [resources.permissions.PermItem](#resources-permissions-PermItem) | repeated |  |
| `to_remove` | [resources.permissions.PermItem](#resources-permissions-PermItem) | repeated |  |






<a name="services-settings-SetJobPropsRequest"></a>

### SetJobPropsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resources-jobs-JobProps) |  |  |






<a name="services-settings-SetJobPropsResponse"></a>

### SetJobPropsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resources-jobs-JobProps) |  |  |






<a name="services-settings-UpdateJobLimitsRequest"></a>

### UpdateJobLimitsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `perms` | [PermsUpdate](#services-settings-PermsUpdate) | optional |  |
| `attrs` | [AttrsUpdate](#services-settings-AttrsUpdate) | optional |  |






<a name="services-settings-UpdateJobLimitsResponse"></a>

### UpdateJobLimitsResponse







<a name="services-settings-UpdateRolePermsRequest"></a>

### UpdateRolePermsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `perms` | [PermsUpdate](#services-settings-PermsUpdate) | optional |  |
| `attrs` | [AttrsUpdate](#services-settings-AttrsUpdate) | optional |  |






<a name="services-settings-UpdateRolePermsResponse"></a>

### UpdateRolePermsResponse







<a name="services-settings-ViewAuditLogRequest"></a>

### ViewAuditLogRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `user_ids` | [int32](#int32) | repeated | Search params |
| `from` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `to` | [resources.timestamp.Timestamp](#resources-timestamp-Timestamp) | optional |  |
| `services` | [string](#string) | repeated | @sanitize: method=StripTags |
| `methods` | [string](#string) | repeated | @sanitize: method=StripTags |
| `search` | [string](#string) | optional |  |






<a name="services-settings-ViewAuditLogResponse"></a>

### ViewAuditLogResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `logs` | [resources.audit.AuditEntry](#resources-audit-AuditEntry) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-settings-SettingsService"></a>

### SettingsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetJobProps` | [GetJobPropsRequest](#services-settings-GetJobPropsRequest) | [GetJobPropsResponse](#services-settings-GetJobPropsResponse) | @perm |
| `SetJobProps` | [SetJobPropsRequest](#services-settings-SetJobPropsRequest) | [SetJobPropsResponse](#services-settings-SetJobPropsResponse) | @perm |
| `GetRoles` | [GetRolesRequest](#services-settings-GetRolesRequest) | [GetRolesResponse](#services-settings-GetRolesResponse) | @perm |
| `GetRole` | [GetRoleRequest](#services-settings-GetRoleRequest) | [GetRoleResponse](#services-settings-GetRoleResponse) | @perm: Name=GetRoles |
| `CreateRole` | [CreateRoleRequest](#services-settings-CreateRoleRequest) | [CreateRoleResponse](#services-settings-CreateRoleResponse) | @perm |
| `DeleteRole` | [DeleteRoleRequest](#services-settings-DeleteRoleRequest) | [DeleteRoleResponse](#services-settings-DeleteRoleResponse) | @perm |
| `UpdateRolePerms` | [UpdateRolePermsRequest](#services-settings-UpdateRolePermsRequest) | [UpdateRolePermsResponse](#services-settings-UpdateRolePermsResponse) | @perm |
| `GetPermissions` | [GetPermissionsRequest](#services-settings-GetPermissionsRequest) | [GetPermissionsResponse](#services-settings-GetPermissionsResponse) | @perm: Name=GetRoles |
| `GetEffectivePermissions` | [GetEffectivePermissionsRequest](#services-settings-GetEffectivePermissionsRequest) | [GetEffectivePermissionsResponse](#services-settings-GetEffectivePermissionsResponse) | @perm: Name=GetRoles |
| `ViewAuditLog` | [ViewAuditLogRequest](#services-settings-ViewAuditLogRequest) | [ViewAuditLogResponse](#services-settings-ViewAuditLogResponse) | @perm |
| `GetAllPermissions` | [GetAllPermissionsRequest](#services-settings-GetAllPermissionsRequest) | [GetAllPermissionsResponse](#services-settings-GetAllPermissionsResponse) | @perm: Name=Superuser |
| `GetJobLimits` | [GetJobLimitsRequest](#services-settings-GetJobLimitsRequest) | [GetJobLimitsResponse](#services-settings-GetJobLimitsResponse) | @perm: Name=Superuser |
| `UpdateJobLimits` | [UpdateJobLimitsRequest](#services-settings-UpdateJobLimitsRequest) | [UpdateJobLimitsResponse](#services-settings-UpdateJobLimitsResponse) | @perm: Name=Superuser |
| `DeleteFaction` | [DeleteFactionRequest](#services-settings-DeleteFactionRequest) | [DeleteFactionResponse](#services-settings-DeleteFactionResponse) | @perm: Name=Superuser |
| `ListDiscordChannels` | [ListDiscordChannelsRequest](#services-settings-ListDiscordChannelsRequest) | [ListDiscordChannelsResponse](#services-settings-ListDiscordChannelsResponse) | @perm: Name=SetJobProps |
| `ListUserGuilds` | [ListUserGuildsRequest](#services-settings-ListUserGuildsRequest) | [ListUserGuildsResponse](#services-settings-ListUserGuildsResponse) | @perm: Name=SetJobProps |
| `UploadJobLogo` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=SetJobProps |
| `DeleteJobLogo` | [DeleteJobLogoRequest](#services-settings-DeleteJobLogoRequest) | [DeleteJobLogoResponse](#services-settings-DeleteJobLogoResponse) | @perm: Name=SetJobProps |

 <!-- end services -->



<a name="services_vehicles_vehicles-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/vehicles/vehicles.proto



<a name="services-vehicles-ListVehiclesRequest"></a>

### ListVehiclesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resources-common-database-Sort) | optional |  |
| `license_plate` | [string](#string) | optional | Search params |
| `model` | [string](#string) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `job` | [string](#string) | optional |  |






<a name="services-vehicles-ListVehiclesResponse"></a>

### ListVehiclesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `vehicles` | [resources.vehicles.Vehicle](#resources-vehicles-Vehicle) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-vehicles-VehiclesService"></a>

### VehiclesService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListVehicles` | [ListVehiclesRequest](#services-vehicles-ListVehiclesRequest) | [ListVehiclesResponse](#services-vehicles-ListVehiclesResponse) | @perm |

 <!-- end services -->



<a name="services_filestore_filestore-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/filestore/filestore.proto



<a name="services-filestore-DeleteFileByPathRequest"></a>

### DeleteFileByPathRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  |  |






<a name="services-filestore-DeleteFileByPathResponse"></a>

### DeleteFileByPathResponse







<a name="services-filestore-ListFilesRequest"></a>

### ListFilesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `path` | [string](#string) | optional |  |






<a name="services-filestore-ListFilesResponse"></a>

### ListFilesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `files` | [resources.file.File](#resources-file-File) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-filestore-FilestoreService"></a>

### FilestoreService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Upload` | [.resources.file.UploadPacket](#resources-file-UploadPacket) stream | [.resources.file.UploadResponse](#resources-file-UploadResponse) | @perm: Name=Superuser |
| `ListFiles` | [ListFilesRequest](#services-filestore-ListFilesRequest) | [ListFilesResponse](#services-filestore-ListFilesResponse) | @perm: Name=Superuser |
| `DeleteFile` | [.resources.file.DeleteFileRequest](#resources-file-DeleteFileRequest) | [.resources.file.DeleteFileResponse](#resources-file-DeleteFileResponse) | @perm: Name=Superuser |
| `DeleteFileByPath` | [DeleteFileByPathRequest](#services-filestore-DeleteFileByPathRequest) | [DeleteFileByPathResponse](#services-filestore-DeleteFileByPathResponse) | @perm: Name=Superuser |

 <!-- end services -->



<a name="services_notifications_notifications-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## services/notifications/notifications.proto



<a name="services-notifications-GetNotificationsRequest"></a>

### GetNotificationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resources-common-database-PaginationRequest) |  |  |
| `include_read` | [bool](#bool) | optional |  |
| `categories` | [resources.notifications.NotificationCategory](#resources-notifications-NotificationCategory) | repeated |  |






<a name="services-notifications-GetNotificationsResponse"></a>

### GetNotificationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resources-common-database-PaginationResponse) |  |  |
| `notifications` | [resources.notifications.Notification](#resources-notifications-Notification) | repeated |  |






<a name="services-notifications-MarkNotificationsRequest"></a>

### MarkNotificationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unread` | [bool](#bool) |  |  |
| `ids` | [uint64](#uint64) | repeated |  |
| `all` | [bool](#bool) | optional |  |






<a name="services-notifications-MarkNotificationsResponse"></a>

### MarkNotificationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [uint64](#uint64) |  |  |






<a name="services-notifications-StreamMessage"></a>

### StreamMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_view` | [resources.notifications.ClientView](#resources-notifications-ClientView) |  |  |






<a name="services-notifications-StreamResponse"></a>

### StreamResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `notification_count` | [int32](#int32) |  |  |
| `restart` | [bool](#bool) | optional |  |
| `user_event` | [resources.notifications.UserEvent](#resources-notifications-UserEvent) |  |  |
| `job_event` | [resources.notifications.JobEvent](#resources-notifications-JobEvent) |  |  |
| `job_grade_event` | [resources.notifications.JobGradeEvent](#resources-notifications-JobGradeEvent) |  |  |
| `system_event` | [resources.notifications.SystemEvent](#resources-notifications-SystemEvent) |  |  |
| `mailer_event` | [resources.mailer.MailerEvent](#resources-mailer-MailerEvent) |  |  |
| `object_event` | [resources.notifications.ObjectEvent](#resources-notifications-ObjectEvent) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="services-notifications-NotificationsService"></a>

### NotificationsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetNotifications` | [GetNotificationsRequest](#services-notifications-GetNotificationsRequest) | [GetNotificationsResponse](#services-notifications-GetNotificationsResponse) | @perm: Name=Any |
| `MarkNotifications` | [MarkNotificationsRequest](#services-notifications-MarkNotificationsRequest) | [MarkNotificationsResponse](#services-notifications-MarkNotificationsResponse) | @perm: Name=Any |
| `Stream` | [StreamMessage](#services-notifications-StreamMessage) stream | [StreamResponse](#services-notifications-StreamResponse) stream | @perm: Name=Any |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> `double` |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> `float` |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> `int32` | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> `int64` | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> `uint32` | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> `uint64` | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> `sint32` | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> `sint64` | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> `fixed32` | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> `fixed64` | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> `sfixed32` | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> `sfixed64` | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> `bool` |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> `string` | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> `bytes` | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

