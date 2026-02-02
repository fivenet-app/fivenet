---
title: GRPC Protobuf Documentation
description: Documentation for GRPC Protobuf files.
---



## codegen/audit/redacted.proto

 <!-- end messages -->

 <!-- end enums -->


### File-level Extensions

| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| `redacted` | bool | .google.protobuf.FieldOptions | 51006 |  |

 <!-- end HasExtensions -->

 <!-- end services -->



## codegen/dbscanner/dbscanner.proto


### codegen.dbscanner.MessageOptions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `not_json` | [bool](#bool) | optional |  |
| `partial` | [bool](#bool) | optional |  |




 <!-- end messages -->

 <!-- end enums -->


### File-level Extensions

| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| `dbscanner` | MessageOptions | .google.protobuf.MessageOptions | 51004 |  |

 <!-- end HasExtensions -->

 <!-- end services -->



## codegen/itemslen/itemslen.proto

 <!-- end messages -->

 <!-- end enums -->


### File-level Extensions

| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| `enabled` | bool | .google.protobuf.FieldOptions | 51001 |  |

 <!-- end HasExtensions -->

 <!-- end services -->



## codegen/sanitizer/sanitizer.proto


### codegen.sanitizer.FieldOptions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `method` | [string](#string) | optional |  |
| `strip_html_tags` | [bool](#bool) | optional |  |
| `tiptap_json` | [bool](#bool) | optional |  |
| `max_bytes` | [uint32](#uint32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->


### File-level Extensions

| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| `sanitizer` | FieldOptions | .google.protobuf.FieldOptions | 51003 |  |

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/timestamp/timestamp.proto


### resources.timestamp.Timestamp
Timestamp for storage messages. We've defined a new local type wrapper of google.protobuf.Timestamp so we can implement sql.Scanner and sql.Valuer interfaces. See: https://golang.org/pkg/database/sql/#Scanner https://golang.org/pkg/database/sql/driver/#Valuer



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timestamp` | [google.protobuf.Timestamp](https://protobuf.dev/reference/protobuf/google.protobuf/#timestamp) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/permissions/attributes/attributes.proto


### resources.permissions.attributes.AttributeValues


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `string_list` | [StringList](#resourcespermissionsattributesStringList) |  |  |
| `job_list` | [StringList](#resourcespermissionsattributesStringList) |  |  |
| `job_grade_list` | [JobGradeList](#resourcespermissionsattributesJobGradeList) |  |  |





### resources.permissions.attributes.JobGradeList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fine_grained` | [bool](#bool) |  |  |
| `jobs` | [JobGradeList.JobsEntry](#resourcespermissionsattributesJobGradeListJobsEntry) | repeated |  |
| `grades` | [JobGradeList.GradesEntry](#resourcespermissionsattributesJobGradeListGradesEntry) | repeated |  |





### resources.permissions.attributes.JobGradeList.GradesEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [JobGrades](#resourcespermissionsattributesJobGrades) |  |  |





### resources.permissions.attributes.JobGradeList.JobsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [int32](#int32) |  |  |





### resources.permissions.attributes.JobGrades


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grades` | [int32](#int32) | repeated |  |





### resources.permissions.attributes.RoleAttribute


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `attr_id` | [int64](#int64) |  |  |
| `permission_id` | [int64](#int64) |  |  |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `key` | [string](#string) |  |  |
| `type` | [string](#string) |  |  |
| `valid_values` | [AttributeValues](#resourcespermissionsattributesAttributeValues) |  |  |
| `value` | [AttributeValues](#resourcespermissionsattributesAttributeValues) |  |  |
| `max_values` | [AttributeValues](#resourcespermissionsattributesAttributeValues) | optional |  |





### resources.permissions.attributes.StringList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `strings` | [string](#string) | repeated |  |




 <!-- end messages -->


### resources.permissions.attributes.AttributeType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `ATTRIBUTE_TYPE_UNSPECIFIED` | 0 |  |
| `ATTRIBUTE_TYPE_STRING_LIST` | 1 |  |
| `ATTRIBUTE_TYPE_JOB_LIST` | 2 |  |
| `ATTRIBUTE_TYPE_JOB_GRADE_LIST` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## codegen/perms/perms.proto


### codegen.perms.Attr


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |
| `type` | [resources.permissions.attributes.AttributeType](#resourcespermissionsattributesAttributeType) |  |  |
| `valid_string_list` | [string](#string) | repeated |  |





### codegen.perms.PermsOptions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `service` | [string](#string) | optional |  |
| `name` | [string](#string) | optional |  |
| `names` | [string](#string) | repeated |  |
| `order` | [int32](#int32) |  |  |
| `attrs` | [Attr](#codegenpermsAttr) | repeated |  |





### codegen.perms.ServiceOptions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `order` | [int32](#int32) |  |  |
| `icon` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->


### File-level Extensions

| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| `perms` | PermsOptions | .google.protobuf.MethodOptions | 51002 |  |
| `perms_svc` | ServiceOptions | .google.protobuf.ServiceOptions | 51005 |  |

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/accounts/oauth2/oauth2.proto


### resources.accounts.oauth2.OAuth2Account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `provider_name` | [string](#string) |  |  |
| `provider` | [OAuth2Provider](#resourcesaccountsoauth2OAuth2Provider) |  |  |
| `external_id` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `avatar` | [string](#string) |  |  |





### resources.accounts.oauth2.OAuth2Provider


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



## resources/users/licenses/licenses.proto


### resources.users.licenses.CitizensLicenses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `licenses` | [License](#resourcesuserslicensesLicense) | repeated |  |





### resources.users.licenses.License


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `label` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/file/meta.proto


### resources.file.FileMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image` | [ImageMeta](#resourcesfileImageMeta) |  |  |





### resources.file.ImageMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `width` | [int64](#int64) |  |  |
| `height` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/file/file.proto


### resources.file.File


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [int64](#int64) | optional |  |
| `id` | [int64](#int64) |  |  |
| `file_path` | [string](#string) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `byte_size` | [int64](#int64) |  | Bytes stored |
| `content_type` | [string](#string) |  |  |
| `meta` | [FileMeta](#resourcesfileFileMeta) | optional |  |
| `is_dir` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/jobs.proto


### resources.jobs.Job


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `label` | [string](#string) |  |  |
| `grades` | [JobGrade](#resourcesjobsJobGrade) | repeated |  |





### resources.jobs.JobGrade


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_name` | [string](#string) | optional |  |
| `grade` | [int32](#int32) |  |  |
| `label` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/labels/labels.proto


### resources.users.labels.Label


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  |  |
| `color` | [string](#string) |  |  |





### resources.users.labels.Labels


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [Label](#resourcesuserslabelsLabel) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/props/props.proto


### resources.users.props.UserProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `wanted` | [bool](#bool) | optional |  |
| `job_name` | [string](#string) | optional |  |
| `job` | [resources.jobs.Job](#resourcesjobsJob) | optional |  |
| `job_grade_number` | [int32](#int32) | optional |  |
| `job_grade` | [resources.jobs.JobGrade](#resourcesjobsJobGrade) | optional |  |
| `traffic_infraction_points` | [uint32](#uint32) | optional |  |
| `traffic_infraction_points_updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `open_fines` | [int64](#int64) | optional |  |
| `blood_type` | [string](#string) | optional |  |
| `mugshot_file_id` | [int64](#int64) | optional |  |
| `mugshot` | [resources.file.File](#resourcesfileFile) | optional |  |
| `labels` | [resources.users.labels.Labels](#resourcesuserslabelsLabels) | optional |  |
| `email` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/user.proto


### resources.users.PhoneNumber


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `number` | [string](#string) |  |  |
| `is_primary` | [bool](#bool) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.users.User


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `identifier` | [string](#string) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `jobs` | [UserJob](#resourcesusersUserJob) | repeated |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `sex` | [string](#string) | optional |  |
| `height` | [string](#string) | optional |  |
| `phone_number` | [string](#string) | optional |  |
| `phone_numbers` | [PhoneNumber](#resourcesusersPhoneNumber) | repeated |  |
| `visum` | [int32](#int32) | optional |  |
| `playtime` | [int32](#int32) | optional |  |
| `props` | [props.UserProps](#resourcesuserspropsUserProps) |  |  |
| `licenses` | [licenses.License](#resourcesuserslicensesLicense) | repeated |  |
| `profile_picture_file_id` | [int64](#int64) | optional |  |
| `profile_picture` | [string](#string) | optional |  |





### resources.users.UserJob


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) |  |  |
| `grade_label` | [string](#string) | optional |  |
| `is_primary` | [bool](#bool) |  |  |





### resources.users.UserShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `phone_number` | [string](#string) | optional |  |
| `profile_picture_file_id` | [int64](#int64) | optional |  |
| `profile_picture` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/accounts/accounts.proto


### resources.accounts.Account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `enabled` | [bool](#bool) |  |  |
| `username` | [string](#string) |  |  |
| `license` | [string](#string) |  |  |
| `groups` | [AccountGroups](#resourcesaccountsAccountGroups) | optional |  |
| `last_char` | [int32](#int32) | optional |  |
| `oauth2_accounts` | [oauth2.OAuth2Account](#resourcesaccountsoauth2OAuth2Account) | repeated |  |





### resources.accounts.AccountGroups


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `groups` | [string](#string) | repeated |  |





### resources.accounts.Character


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `available` | [bool](#bool) |  |  |
| `group` | [string](#string) |  |  |
| `char` | [resources.users.User](#resourcesusersUser) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/short/user.proto


### resources.users.short.UserShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `phone_number` | [string](#string) | optional |  |
| `profile_picture_file_id` | [int64](#int64) | optional |  |
| `profile_picture` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/audit/audit.proto


### resources.audit.AuditEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `user_job` | [string](#string) |  |  |
| `target_user_id` | [int32](#int32) | optional |  |
| `target_user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `target_user_job` | [string](#string) | optional |  |
| `service` | [string](#string) |  | GRPC Service name |
| `method` | [string](#string) |  | GRPC Method name |
| `action` | [EventAction](#resourcesauditEventAction) |  |  |
| `result` | [EventResult](#resourcesauditEventResult) |  |  |
| `meta` | [AuditEntryMeta](#resourcesauditAuditEntryMeta) | optional |  |
| `data` | [string](#string) | optional |  |





### resources.audit.AuditEntryMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `meta` | [AuditEntryMeta.MetaEntry](#resourcesauditAuditEntryMetaMetaEntry) | repeated |  |





### resources.audit.AuditEntryMeta.MetaEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |




 <!-- end messages -->


### resources.audit.EventAction

| Name | Number | Description |
| ---- | ------ | ----------- |
| `EVENT_ACTION_UNSPECIFIED` | 0 |  |
| `EVENT_ACTION_VIEWED` | 2 |  |
| `EVENT_ACTION_CREATED` | 3 |  |
| `EVENT_ACTION_UPDATED` | 4 |  |
| `EVENT_ACTION_DELETED` | 5 |  |



### resources.audit.EventResult

| Name | Number | Description |
| ---- | ------ | ----------- |
| `EVENT_RESULT_UNSPECIFIED` | 0 |  |
| `EVENT_RESULT_SUCCEEDED` | 1 |  |
| `EVENT_RESULT_FAILED` | 2 |  |
| `EVENT_RESULT_ERRORED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/calendar/access/access.proto


### resources.calendar.access.CalendarAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [CalendarJobAccess](#resourcescalendaraccessCalendarJobAccess) | repeated |  |
| `users` | [CalendarUserAccess](#resourcescalendaraccessCalendarUserAccess) | repeated |  |





### resources.calendar.access.CalendarJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcescalendaraccessAccessLevel) |  |  |





### resources.calendar.access.CalendarUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `access` | [AccessLevel](#resourcescalendaraccessAccessLevel) |  |  |




 <!-- end messages -->


### resources.calendar.access.AccessLevel

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



## resources/calendar/calendar.proto


### resources.calendar.Calendar


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `public` | [bool](#bool) |  |  |
| `closed` | [bool](#bool) |  |  |
| `color` | [string](#string) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `subscription` | [CalendarSub](#resourcescalendarCalendarSub) | optional |  |
| `access` | [access.CalendarAccess](#resourcescalendaraccessCalendarAccess) |  |  |





### resources.calendar.CalendarShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `public` | [bool](#bool) |  |  |
| `closed` | [bool](#bool) |  |  |
| `color` | [string](#string) |  |  |
| `subscription` | [CalendarSub](#resourcescalendarCalendarSub) | optional |  |





### resources.calendar.CalendarSub


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `confirmed` | [bool](#bool) |  |  |
| `muted` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/content/content.proto


### resources.common.content.Content


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `content_type` | [ContentType](#resourcescommoncontentContentType) |  |  |
| `raw_html` | [string](#string) | optional | Deprecated: Legacy HTML (only for viewing old content) |
| `content` | [RichTextHtmlNode](#resourcescommoncontentRichTextHtmlNode) | optional | Deprecated: Legacy custom HTML to JSON AST (only for viewing old content) |
| `tiptap_json` | [google.protobuf.Struct](https://protobuf.dev/reference/protobuf/google.protobuf/#struct) |  | Tiptap JSON Document |





### resources.common.content.ExtractedContent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `text` | [string](#string) |  |  |
| `word_count` | [uint32](#uint32) |  |  |
| `first_heading` | [string](#string) |  |  |





### resources.common.content.RichTextHtmlNode


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [NodeType](#resourcescommoncontentNodeType) |  |  |
| `id` | [string](#string) | optional |  |
| `tag` | [string](#string) |  |  |
| `attrs` | [RichTextHtmlNode.AttrsEntry](#resourcescommoncontentRichTextHtmlNodeAttrsEntry) | repeated |  |
| `text` | [string](#string) | optional |  |
| `content` | [RichTextHtmlNode](#resourcescommoncontentRichTextHtmlNode) | repeated |  |





### resources.common.content.RichTextHtmlNode.AttrsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |




 <!-- end messages -->


### resources.common.content.ContentType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CONTENT_TYPE_UNSPECIFIED` | 0 |  |
| `CONTENT_TYPE_HTML` | 1 | Used for legacy HTML content |
| `CONTENT_TYPE_TIPTAP_JSON` | 2 | Used for Tiptap JSON content |



### resources.common.content.NodeType

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



## resources/calendar/entries/entries.proto


### resources.calendar.entries.CalendarEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `calendar_id` | [int64](#int64) |  |  |
| `calendar` | [resources.calendar.Calendar](#resourcescalendarCalendar) | optional |  |
| `job` | [string](#string) | optional |  |
| `start_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `end_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `closed` | [bool](#bool) |  |  |
| `rsvp_open` | [bool](#bool) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `recurring` | [CalendarEntryRecurring](#resourcescalendarentriesCalendarEntryRecurring) | optional |  |
| `rsvp` | [CalendarEntryRSVP](#resourcescalendarentriesCalendarEntryRSVP) | optional |  |





### resources.calendar.entries.CalendarEntryRSVP


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `response` | [RsvpResponses](#resourcescalendarentriesRsvpResponses) |  |  |





### resources.calendar.entries.CalendarEntryRecurring


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `every` | [string](#string) |  |  |
| `count` | [int32](#int32) |  |  |
| `until` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->


### resources.calendar.entries.RsvpResponses

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



## resources/centrum/access/access.proto


### resources.centrum.access.CentrumAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [CentrumJobAccess](#resourcescentrumaccessCentrumJobAccess) | repeated |  |





### resources.centrum.access.CentrumJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `source_job` | [string](#string) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [CentrumAccessLevel](#resourcescentrumaccessCentrumAccessLevel) |  |  |
| `accepted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.centrum.access.CentrumQualificationAccess
Dummy - DO NOT USE!






### resources.centrum.access.CentrumUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


### resources.centrum.access.CentrumAccessLevel

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `CENTRUM_ACCESS_LEVEL_BLOCKED` | 1 |  |
| `CENTRUM_ACCESS_LEVEL_VIEW` | 2 |  |
| `CENTRUM_ACCESS_LEVEL_PARTICIPATE` | 3 |  |
| `CENTRUM_ACCESS_LEVEL_DISPATCH` | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/labels/labels.proto


### resources.jobs.labels.Label


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `job` | [string](#string) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `color` | [string](#string) |  |  |
| `order` | [int32](#int32) |  |  |





### resources.jobs.labels.LabelCount


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `label` | [Label](#resourcesjobslabelsLabel) |  |  |
| `count` | [int64](#int64) |  |  |





### resources.jobs.labels.Labels


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [Label](#resourcesjobslabelsLabel) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/colleagues/colleagues.proto


### resources.jobs.colleagues.Colleague


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) |  |  |
| `dateofbirth` | [string](#string) |  |  |
| `phone_number` | [string](#string) | optional |  |
| `profile_picture_file_id` | [int64](#int64) | optional |  |
| `profile_picture` | [string](#string) | optional |  |
| `props` | [ColleagueProps](#resourcesjobscolleaguesColleagueProps) |  |  |
| `email` | [string](#string) | optional |  |





### resources.jobs.colleagues.ColleagueProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `absence_begin` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `absence_end` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `note` | [string](#string) | optional |  |
| `labels` | [resources.jobs.labels.Labels](#resourcesjobslabelsLabels) | optional |  |
| `name_prefix` | [string](#string) | optional |  |
| `name_suffix` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/dispatchers/dispatchers.proto


### resources.centrum.dispatchers.Dispatchers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `dispatchers` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | repeated |  |





### resources.centrum.dispatchers.JobDispatchers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [Dispatchers](#resourcescentrumdispatchersDispatchers) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/joblist.proto


### resources.centrum.JobList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [JobListEntry](#resourcescentrumJobListEntry) | repeated |  |





### resources.centrum.JobListEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `label` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/qualifications/access/access.proto


### resources.qualifications.access.QualificationAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [QualificationJobAccess](#resourcesqualificationsaccessQualificationJobAccess) | repeated |  |





### resources.qualifications.access.QualificationJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesqualificationsaccessAccessLevel) |  |  |





### resources.qualifications.access.QualificationUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


### resources.qualifications.access.AccessLevel

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



## resources/qualifications/exam/exam.proto


### resources.qualifications.exam.ExamGrading


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `responses` | [ExamGradingResponse](#resourcesqualificationsexamExamGradingResponse) | repeated |  |





### resources.qualifications.exam.ExamGradingResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `question_id` | [int64](#int64) |  |  |
| `points` | [float](#float) |  |  |
| `checked` | [bool](#bool) | optional |  |





### resources.qualifications.exam.ExamQuestion


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `qualification_id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `data` | [ExamQuestionData](#resourcesqualificationsexamExamQuestionData) |  |  |
| `answer` | [ExamQuestionAnswerData](#resourcesqualificationsexamExamQuestionAnswerData) | optional |  |
| `points` | [int32](#int32) | optional |  |
| `order` | [int32](#int32) |  |  |





### resources.qualifications.exam.ExamQuestionAnswerData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `answer_key` | [string](#string) |  |  |
| `yesno` | [ExamResponseYesNo](#resourcesqualificationsexamExamResponseYesNo) |  |  |
| `free_text` | [ExamResponseText](#resourcesqualificationsexamExamResponseText) |  |  |
| `single_choice` | [ExamResponseSingleChoice](#resourcesqualificationsexamExamResponseSingleChoice) |  |  |
| `multiple_choice` | [ExamResponseMultipleChoice](#resourcesqualificationsexamExamResponseMultipleChoice) |  |  |





### resources.qualifications.exam.ExamQuestionData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `separator` | [ExamQuestionSeparator](#resourcesqualificationsexamExamQuestionSeparator) |  |  |
| `image` | [ExamQuestionImage](#resourcesqualificationsexamExamQuestionImage) |  |  |
| `yesno` | [ExamQuestionYesNo](#resourcesqualificationsexamExamQuestionYesNo) |  |  |
| `free_text` | [ExamQuestionText](#resourcesqualificationsexamExamQuestionText) |  |  |
| `single_choice` | [ExamQuestionSingleChoice](#resourcesqualificationsexamExamQuestionSingleChoice) |  |  |
| `multiple_choice` | [ExamQuestionMultipleChoice](#resourcesqualificationsexamExamQuestionMultipleChoice) |  |  |





### resources.qualifications.exam.ExamQuestionImage


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image` | [resources.file.File](#resourcesfileFile) |  |  |
| `alt` | [string](#string) | optional |  |





### resources.qualifications.exam.ExamQuestionMultipleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated |  |
| `limit` | [int32](#int32) | optional |  |





### resources.qualifications.exam.ExamQuestionSeparator





### resources.qualifications.exam.ExamQuestionSingleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated |  |





### resources.qualifications.exam.ExamQuestionText


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min_length` | [int32](#int32) |  |  |
| `max_length` | [int32](#int32) |  |  |





### resources.qualifications.exam.ExamQuestionYesNo





### resources.qualifications.exam.ExamQuestions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `questions` | [ExamQuestion](#resourcesqualificationsexamExamQuestion) | repeated |  |





### resources.qualifications.exam.ExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `question_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `question` | [ExamQuestion](#resourcesqualificationsexamExamQuestion) |  |  |
| `response` | [ExamResponseData](#resourcesqualificationsexamExamResponseData) |  |  |





### resources.qualifications.exam.ExamResponseData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `separator` | [ExamResponseSeparator](#resourcesqualificationsexamExamResponseSeparator) |  |  |
| `yesno` | [ExamResponseYesNo](#resourcesqualificationsexamExamResponseYesNo) |  |  |
| `free_text` | [ExamResponseText](#resourcesqualificationsexamExamResponseText) |  |  |
| `single_choice` | [ExamResponseSingleChoice](#resourcesqualificationsexamExamResponseSingleChoice) |  |  |
| `multiple_choice` | [ExamResponseMultipleChoice](#resourcesqualificationsexamExamResponseMultipleChoice) |  |  |





### resources.qualifications.exam.ExamResponseMultipleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated |  |





### resources.qualifications.exam.ExamResponseSeparator





### resources.qualifications.exam.ExamResponseSingleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choice` | [string](#string) |  |  |





### resources.qualifications.exam.ExamResponseText


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `text` | [string](#string) |  | 0.5 Megabyte |





### resources.qualifications.exam.ExamResponseYesNo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [bool](#bool) |  |  |





### resources.qualifications.exam.ExamResponses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `responses` | [ExamResponse](#resourcesqualificationsexamExamResponse) | repeated |  |





### resources.qualifications.exam.ExamUser


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `started_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `ends_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `ended_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.qualifications.exam.QualificationExamSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `time` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `auto_grade` | [bool](#bool) |  |  |
| `auto_grade_mode` | [AutoGradeMode](#resourcesqualificationsexamAutoGradeMode) |  |  |
| `minimum_points` | [int32](#int32) |  |  |




 <!-- end messages -->


### resources.qualifications.exam.AutoGradeMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `AUTO_GRADE_MODE_UNSPECIFIED` | 0 |  |
| `AUTO_GRADE_MODE_STRICT` | 1 |  |
| `AUTO_GRADE_MODE_PARTIAL_CREDIT` | 2 |  |



### resources.qualifications.exam.QualificationExamMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `QUALIFICATION_EXAM_MODE_UNSPECIFIED` | 0 |  |
| `QUALIFICATION_EXAM_MODE_DISABLED` | 1 |  |
| `QUALIFICATION_EXAM_MODE_REQUEST_NEEDED` | 2 |  |
| `QUALIFICATION_EXAM_MODE_ENABLED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/qualifications/qualifications.proto


### resources.qualifications.Qualification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `abbreviation` | [string](#string) |  |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `access` | [access.QualificationAccess](#resourcesqualificationsaccessQualificationAccess) |  |  |
| `requirements` | [QualificationRequirement](#resourcesqualificationsQualificationRequirement) | repeated |  |
| `discord_sync_enabled` | [bool](#bool) |  |  |
| `discord_settings` | [QualificationDiscordSettings](#resourcesqualificationsQualificationDiscordSettings) | optional |  |
| `exam_mode` | [exam.QualificationExamMode](#resourcesqualificationsexamQualificationExamMode) |  |  |
| `exam_settings` | [exam.QualificationExamSettings](#resourcesqualificationsexamQualificationExamSettings) | optional |  |
| `exam` | [exam.ExamQuestions](#resourcesqualificationsexamExamQuestions) | optional |  |
| `result` | [QualificationResult](#resourcesqualificationsQualificationResult) | optional |  |
| `request` | [QualificationRequest](#resourcesqualificationsQualificationRequest) | optional |  |
| `label_sync_enabled` | [bool](#bool) |  |  |
| `label_sync_format` | [string](#string) | optional |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### resources.qualifications.QualificationDiscordSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_name` | [string](#string) | optional |  |
| `role_format` | [string](#string) | optional |  |





### resources.qualifications.QualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `qualification_id` | [int64](#int64) |  |  |
| `qualification` | [QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) |  |  |
| `user_comment` | [string](#string) | optional |  |
| `status` | [RequestStatus](#resourcesqualificationsRequestStatus) | optional |  |
| `approved_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `approver_comment` | [string](#string) | optional |  |
| `approver_id` | [int32](#int32) | optional |  |
| `approver` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `approver_job` | [string](#string) | optional |  |





### resources.qualifications.QualificationRequirement


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `qualification_id` | [int64](#int64) |  |  |
| `target_qualification_id` | [int64](#int64) |  |  |
| `target_qualification` | [QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |





### resources.qualifications.QualificationResult


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `qualification_id` | [int64](#int64) |  |  |
| `qualification` | [QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) |  |  |
| `status` | [ResultStatus](#resourcesqualificationsResultStatus) |  |  |
| `score` | [float](#float) | optional |  |
| `summary` | [string](#string) |  |  |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) |  |  |
| `creator_job` | [string](#string) |  |  |





### resources.qualifications.QualificationShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `abbreviation` | [string](#string) |  |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `requirements` | [QualificationRequirement](#resourcesqualificationsQualificationRequirement) | repeated |  |
| `exam_mode` | [exam.QualificationExamMode](#resourcesqualificationsexamQualificationExamMode) |  |  |
| `exam_settings` | [exam.QualificationExamSettings](#resourcesqualificationsexamQualificationExamSettings) | optional |  |
| `result` | [QualificationResult](#resourcesqualificationsQualificationResult) | optional |  |




 <!-- end messages -->


### resources.qualifications.RequestStatus

| Name | Number | Description |
| ---- | ------ | ----------- |
| `REQUEST_STATUS_UNSPECIFIED` | 0 |  |
| `REQUEST_STATUS_PENDING` | 1 |  |
| `REQUEST_STATUS_DENIED` | 2 |  |
| `REQUEST_STATUS_ACCEPTED` | 3 |  |
| `REQUEST_STATUS_EXAM_STARTED` | 4 |  |
| `REQUEST_STATUS_EXAM_GRADING` | 5 |  |
| `REQUEST_STATUS_COMPLETED` | 6 |  |



### resources.qualifications.ResultStatus

| Name | Number | Description |
| ---- | ------ | ----------- |
| `RESULT_STATUS_UNSPECIFIED` | 0 |  |
| `RESULT_STATUS_PENDING` | 1 |  |
| `RESULT_STATUS_FAILED` | 2 |  |
| `RESULT_STATUS_SUCCESSFUL` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/units/access/access.proto


### resources.centrum.units.access.UnitAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [UnitJobAccess](#resourcescentrumunitsaccessUnitJobAccess) | repeated |  |
| `qualifications` | [UnitQualificationAccess](#resourcescentrumunitsaccessUnitQualificationAccess) | repeated |  |





### resources.centrum.units.access.UnitJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [UnitAccessLevel](#resourcescentrumunitsaccessUnitAccessLevel) |  |  |





### resources.centrum.units.access.UnitQualificationAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `qualification_id` | [int64](#int64) |  |  |
| `qualification` | [resources.qualifications.QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `access` | [UnitAccessLevel](#resourcescentrumunitsaccessUnitAccessLevel) |  |  |





### resources.centrum.units.access.UnitUserAccess




 <!-- end messages -->


### resources.centrum.units.access.UnitAccessLevel

| Name | Number | Description |
| ---- | ------ | ----------- |
| `UNIT_ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `UNIT_ACCESS_LEVEL_BLOCKED` | 1 |  |
| `UNIT_ACCESS_LEVEL_JOIN` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/units/units.proto


### resources.centrum.units.Unit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `name` | [string](#string) |  |  |
| `initials` | [string](#string) |  |  |
| `color` | [string](#string) |  |  |
| `icon` | [string](#string) | optional |  |
| `description` | [string](#string) | optional |  |
| `status` | [UnitStatus](#resourcescentrumunitsUnitStatus) | optional |  |
| `users` | [UnitAssignment](#resourcescentrumunitsUnitAssignment) | repeated |  |
| `attributes` | [UnitAttributes](#resourcescentrumunitsUnitAttributes) | optional |  |
| `home_postal` | [string](#string) | optional |  |
| `access` | [access.UnitAccess](#resourcescentrumunitsaccessUnitAccess) |  |  |





### resources.centrum.units.UnitAssignment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |





### resources.centrum.units.UnitAssignments


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `users` | [UnitAssignment](#resourcescentrumunitsUnitAssignment) | repeated |  |





### resources.centrum.units.UnitAttributes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [UnitAttribute](#resourcescentrumunitsUnitAttribute) | repeated |  |





### resources.centrum.units.UnitStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `unit_id` | [int64](#int64) |  |  |
| `unit` | [Unit](#resourcescentrumunitsUnit) | optional |  |
| `status` | [StatusUnit](#resourcescentrumunitsStatusUnit) |  |  |
| `reason` | [string](#string) | optional |  |
| `code` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |
| `x` | [double](#double) | optional |  |
| `y` | [double](#double) | optional |  |
| `postal` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |
| `creator_job` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.centrum.units.StatusUnit

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



### resources.centrum.units.UnitAttribute

| Name | Number | Description |
| ---- | ------ | ----------- |
| `UNIT_ATTRIBUTE_UNSPECIFIED` | 0 |  |
| `UNIT_ATTRIBUTE_STATIC` | 1 |  |
| `UNIT_ATTRIBUTE_NO_DISPATCH_AUTO_ASSIGN` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/dispatches/dispatches.proto


### resources.centrum.dispatches.Dispatch


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  | **Deprecated.**  |
| `jobs` | [resources.centrum.JobList](#resourcescentrumJobList) |  |  |
| `status` | [DispatchStatus](#resourcescentrumdispatchesDispatchStatus) | optional |  |
| `message` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `attributes` | [DispatchAttributes](#resourcescentrumdispatchesDispatchAttributes) | optional |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `postal` | [string](#string) | optional |  |
| `anon` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.User](#resourcesusersUser) | optional |  |
| `units` | [DispatchAssignment](#resourcescentrumdispatchesDispatchAssignment) | repeated |  |
| `references` | [DispatchReferences](#resourcescentrumdispatchesDispatchReferences) | optional |  |





### resources.centrum.dispatches.DispatchAssignment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [int64](#int64) |  |  |
| `unit_id` | [int64](#int64) |  |  |
| `unit` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.centrum.dispatches.DispatchAssignments


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `units` | [DispatchAssignment](#resourcescentrumdispatchesDispatchAssignment) | repeated |  |





### resources.centrum.dispatches.DispatchAttributes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [DispatchAttribute](#resourcescentrumdispatchesDispatchAttribute) | repeated |  |





### resources.centrum.dispatches.DispatchReference


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_dispatch_id` | [int64](#int64) |  |  |
| `reference_type` | [DispatchReferenceType](#resourcescentrumdispatchesDispatchReferenceType) |  |  |





### resources.centrum.dispatches.DispatchReferences


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `references` | [DispatchReference](#resourcescentrumdispatchesDispatchReference) | repeated |  |





### resources.centrum.dispatches.DispatchStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `dispatch_id` | [int64](#int64) |  |  |
| `unit_id` | [int64](#int64) | optional |  |
| `unit` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) | optional |  |
| `status` | [StatusDispatch](#resourcescentrumdispatchesStatusDispatch) |  |  |
| `reason` | [string](#string) | optional |  |
| `code` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |
| `x` | [double](#double) | optional |  |
| `y` | [double](#double) | optional |  |
| `postal` | [string](#string) | optional |  |
| `creator_job` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.centrum.dispatches.DispatchAttribute

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISPATCH_ATTRIBUTE_UNSPECIFIED` | 0 |  |
| `DISPATCH_ATTRIBUTE_MULTIPLE` | 1 |  |
| `DISPATCH_ATTRIBUTE_DUPLICATE` | 2 |  |
| `DISPATCH_ATTRIBUTE_TOO_OLD` | 3 |  |
| `DISPATCH_ATTRIBUTE_AUTOMATIC` | 4 |  |



### resources.centrum.dispatches.DispatchReferenceType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISPATCH_REFERENCE_TYPE_UNSPECIFIED` | 0 |  |
| `DISPATCH_REFERENCE_TYPE_REFERENCED` | 1 |  |
| `DISPATCH_REFERENCE_TYPE_DUPLICATED_BY` | 2 |  |
| `DISPATCH_REFERENCE_TYPE_DUPLICATE_OF` | 3 |  |



### resources.centrum.dispatches.StatusDispatch

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
| `STATUS_DISPATCH_DELETED` | 14 |  |



### resources.centrum.dispatches.TakeDispatchResp

| Name | Number | Description |
| ---- | ------ | ----------- |
| `TAKE_DISPATCH_RESP_UNSPECIFIED` | 0 |  |
| `TAKE_DISPATCH_RESP_TIMEOUT` | 1 |  |
| `TAKE_DISPATCH_RESP_ACCEPTED` | 2 |  |
| `TAKE_DISPATCH_RESP_DECLINED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/settings/settings.proto


### resources.centrum.settings.Configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deduplication_enabled` | [bool](#bool) |  |  |
| `deduplication_radius` | [int64](#int64) |  |  |
| `deduplication_duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) | optional |  |





### resources.centrum.settings.EffectiveAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatches` | [EffectiveDispatchAccess](#resourcescentrumsettingsEffectiveDispatchAccess) |  |  |





### resources.centrum.settings.EffectiveDispatchAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [JobAccessEntry](#resourcescentrumsettingsJobAccessEntry) | repeated |  |





### resources.centrum.settings.JobAccessEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `access` | [resources.centrum.access.CentrumAccessLevel](#resourcescentrumaccessCentrumAccessLevel) |  |  |





### resources.centrum.settings.PredefinedStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_status` | [string](#string) | repeated |  |
| `dispatch_status` | [string](#string) | repeated |  |





### resources.centrum.settings.Settings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |
| `type` | [CentrumType](#resourcescentrumsettingsCentrumType) |  |  |
| `public` | [bool](#bool) |  |  |
| `mode` | [CentrumMode](#resourcescentrumsettingsCentrumMode) |  |  |
| `fallback_mode` | [CentrumMode](#resourcescentrumsettingsCentrumMode) |  |  |
| `predefined_status` | [PredefinedStatus](#resourcescentrumsettingsPredefinedStatus) | optional |  |
| `timings` | [Timings](#resourcescentrumsettingsTimings) |  |  |
| `configuration` | [Configuration](#resourcescentrumsettingsConfiguration) |  |  |
| `access` | [resources.centrum.access.CentrumAccess](#resourcescentrumaccessCentrumAccess) | optional |  |
| `offered_access` | [resources.centrum.access.CentrumAccess](#resourcescentrumaccessCentrumAccess) | optional |  |
| `effective_access` | [EffectiveAccess](#resourcescentrumsettingsEffectiveAccess) | optional |  |





### resources.centrum.settings.Timings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_max_wait` | [int64](#int64) |  |  |
| `require_unit` | [bool](#bool) |  |  |
| `require_unit_reminder_seconds` | [int64](#int64) |  |  |




 <!-- end messages -->


### resources.centrum.settings.CentrumMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_MODE_UNSPECIFIED` | 0 |  |
| `CENTRUM_MODE_MANUAL` | 1 |  |
| `CENTRUM_MODE_CENTRAL_COMMAND` | 2 |  |
| `CENTRUM_MODE_AUTO_ROUND_ROBIN` | 3 |  |
| `CENTRUM_MODE_SIMPLIFIED` | 4 |  |



### resources.centrum.settings.CentrumType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_TYPE_UNSPECIFIED` | 0 |  |
| `CENTRUM_TYPE_DISPATCH` | 1 |  |
| `CENTRUM_TYPE_DELIVERY` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/banner.proto


### resources.settings.BannerMessage


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `title` | [string](#string) |  |  |
| `icon` | [string](#string) | optional |  |
| `color` | [string](#string) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/data.proto


### resources.settings.Data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [DataMode](#resourcessettingsDataMode) |  |  |




 <!-- end messages -->


### resources.settings.DataMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DATA_MODE_UNSPECIFIED` | 0 |  |
| `DATA_MODE_UNAVAILABLE` | 1 |  |
| `DATA_MODE_READONLY` | 2 |  |
| `DATA_MODE_READWRITE` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/config.proto


### resources.settings.AppConfig


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) | optional |  |
| `default_locale` | [string](#string) |  |  |
| `auth` | [Auth](#resourcessettingsAuth) |  |  |
| `perms` | [Perms](#resourcessettingsPerms) |  |  |
| `website` | [Website](#resourcessettingsWebsite) |  |  |
| `job_info` | [JobInfo](#resourcessettingsJobInfo) |  |  |
| `user_tracker` | [UserTracker](#resourcessettingsUserTracker) |  |  |
| `discord` | [Discord](#resourcessettingsDiscord) |  |  |
| `system` | [System](#resourcessettingsSystem) |  |  |
| `display` | [Display](#resourcessettingsDisplay) |  |  |
| `quick_buttons` | [QuickButtons](#resourcessettingsQuickButtons) |  |  |
| `data` | [Data](#resourcessettingsData) |  |  |





### resources.settings.Auth


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signup_enabled` | [bool](#bool) |  |  |
| `last_char_lock` | [bool](#bool) |  |  |





### resources.settings.Discord


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `sync_interval` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `invite_url` | [string](#string) | optional |  |
| `ignored_jobs` | [string](#string) | repeated |  |
| `bot_presence` | [DiscordBotPresence](#resourcessettingsDiscordBotPresence) | optional |  |
| `bot_id` | [string](#string) | optional |  |
| `bot_permissions` | [int64](#int64) |  |  |





### resources.settings.DiscordBotPresence


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [DiscordBotPresenceType](#resourcessettingsDiscordBotPresenceType) |  |  |
| `status` | [string](#string) | optional |  |
| `url` | [string](#string) | optional |  |





### resources.settings.Display


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `intl_locale` | [string](#string) | optional | IETF BCP 47 language tag (e.g. "en-US", "de-DE") |
| `currency_name` | [string](#string) |  | ISO 4217 currency code (e.g. "USD", "EUR") |





### resources.settings.JobInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unemployed_job` | [UnemployedJob](#resourcessettingsUnemployedJob) |  |  |
| `public_jobs` | [string](#string) | repeated |  |
| `hidden_jobs` | [string](#string) | repeated |  |





### resources.settings.Links


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `privacy_policy` | [string](#string) | optional |  |
| `imprint` | [string](#string) | optional |  |





### resources.settings.PenaltyCalculator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_count` | [uint32](#uint32) | optional |  |
| `detention_time_unit` | [PenaltyCalculatorDetentionTimeUnit](#resourcessettingsPenaltyCalculatorDetentionTimeUnit) | optional |  |
| `warn_settings` | [PenaltyCalculatorWarn](#resourcessettingsPenaltyCalculatorWarn) | optional |  |
| `max_leeway` | [uint32](#uint32) | optional |  |





### resources.settings.PenaltyCalculatorDetentionTimeUnit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `singular` | [string](#string) | optional |  |
| `plural` | [string](#string) | optional |  |





### resources.settings.PenaltyCalculatorWarn


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `fine` | [uint32](#uint32) | optional |  |
| `detention_time` | [uint32](#uint32) | optional |  |
| `stvo_points` | [uint32](#uint32) | optional |  |
| `warn_message` | [string](#string) | optional |  |





### resources.settings.Perm


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |





### resources.settings.Perms


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default` | [Perm](#resourcessettingsPerm) | repeated |  |





### resources.settings.QuickButtons


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `penalty_calculator` | [PenaltyCalculator](#resourcessettingsPenaltyCalculator) |  |  |





### resources.settings.System


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `banner_message_enabled` | [bool](#bool) |  |  |
| `banner_message` | [BannerMessage](#resourcessettingsBannerMessage) |  |  |





### resources.settings.UnemployedJob


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |





### resources.settings.UserTracker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_time` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `db_refresh_time` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |





### resources.settings.Website


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `links` | [Links](#resourcessettingsLinks) |  |  |
| `stats_page` | [bool](#bool) |  |  |




 <!-- end messages -->


### resources.settings.DiscordBotPresenceType

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



## resources/clientconfig/clientconfig.proto


### resources.clientconfig.Auth


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signup_enabled` | [bool](#bool) |  |  |
| `last_char_lock` | [bool](#bool) |  |  |
| `providers` | [ProviderConfig](#resourcesclientconfigProviderConfig) | repeated |  |





### resources.clientconfig.ClientConfig


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `default_locale` | [string](#string) |  |  |
| `auth` | [Auth](#resourcesclientconfigAuth) |  |  |
| `discord` | [Discord](#resourcesclientconfigDiscord) |  |  |
| `website` | [Website](#resourcesclientconfigWebsite) |  |  |
| `feature_gates` | [FeatureGates](#resourcesclientconfigFeatureGates) |  |  |
| `game` | [Game](#resourcesclientconfigGame) |  |  |
| `system` | [System](#resourcesclientconfigSystem) |  |  |
| `display` | [Display](#resourcesclientconfigDisplay) |  |  |
| `quick_buttons` | [resources.settings.QuickButtons](#resourcessettingsQuickButtons) |  |  |
| `data` | [resources.settings.Data](#resourcessettingsData) |  |  |





### resources.clientconfig.Discord


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bot_enabled` | [bool](#bool) |  |  |





### resources.clientconfig.Display


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `intl_locale` | [string](#string) | optional | IETF BCP 47 language tag (e.g. "en-US", "de-DE") |
| `currency_name` | [string](#string) |  | ISO 4217 currency code (e.g. "USD", "EUR") |





### resources.clientconfig.FeatureGates





### resources.clientconfig.Game


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unemployed_job_name` | [string](#string) |  |  |
| `start_job_grade` | [int32](#int32) |  |  |





### resources.clientconfig.Links


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `imprint` | [string](#string) | optional |  |
| `privacy_policy` | [string](#string) | optional |  |





### resources.clientconfig.OTLPFrontend


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `url` | [string](#string) |  |  |
| `headers` | [OTLPFrontend.HeadersEntry](#resourcesclientconfigOTLPFrontendHeadersEntry) | repeated |  |





### resources.clientconfig.OTLPFrontend.HeadersEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |





### resources.clientconfig.ProviderConfig


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `label` | [string](#string) |  |  |
| `icon` | [string](#string) | optional |  |
| `homepage` | [string](#string) |  |  |





### resources.clientconfig.System


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `banner_message_enabled` | [bool](#bool) |  |  |
| `banner_message` | [resources.settings.BannerMessage](#resourcessettingsBannerMessage) | optional |  |
| `otlp` | [OTLPFrontend](#resourcesclientconfigOTLPFrontend) |  |  |





### resources.clientconfig.Website


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `links` | [Links](#resourcesclientconfigLinks) |  |  |
| `stats_page` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/collab/collab.proto


### resources.collab.AwarenessPing


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |





### resources.collab.ClientPacket


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hello` | [CollabInit](#resourcescollabCollabInit) |  | Must be the first message |
| `sync_step` | [SyncStep](#resourcescollabSyncStep) |  |  |
| `yjs_update` | [YjsUpdate](#resourcescollabYjsUpdate) |  |  |
| `awareness` | [AwarenessPing](#resourcescollabAwarenessPing) |  |  |





### resources.collab.ClientUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `joined` | [bool](#bool) |  |  |
| `client_id` | [uint64](#uint64) |  |  |
| `label` | [string](#string) | optional |  |





### resources.collab.CollabHandshake


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [uint64](#uint64) |  |  |





### resources.collab.CollabInit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_id` | [int64](#int64) |  |  |





### resources.collab.FirstPromote





### resources.collab.ServerPacket


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_id` | [uint64](#uint64) |  | Who generated this packet (same ID used in awareness) |
| `handshake` | [CollabHandshake](#resourcescollabCollabHandshake) |  |  |
| `sync_step` | [SyncStep](#resourcescollabSyncStep) |  |  |
| `yjs_update` | [YjsUpdate](#resourcescollabYjsUpdate) |  |  |
| `awareness` | [AwarenessPing](#resourcescollabAwarenessPing) |  |  |
| `target_saved` | [TargetSaved](#resourcescollabTargetSaved) |  |  |
| `promote` | [FirstPromote](#resourcescollabFirstPromote) |  |  |
| `client_update` | [ClientUpdate](#resourcescollabClientUpdate) |  |  |





### resources.collab.SyncStep


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `step` | [int32](#int32) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `receiver_id` | [uint64](#uint64) | optional |  |





### resources.collab.TargetSaved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_id` | [int64](#int64) |  |  |





### resources.collab.YjsUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |




 <!-- end messages -->


### resources.collab.ClientRole

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CLIENT_ROLE_UNSPECIFIED` | 0 |  |
| `CLIENT_ROLE_READER` | 1 |  |
| `CLIENT_ROLE_WRITER` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/content/diff_activity.proto


### resources.common.content.ContentDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stats` | [ContentDiffStats](#resourcescommoncontentContentDiffStats) |  | Quick summary for badge like "+12 / -3" |
| `ops` | [ContentDiffOp](#resourcescommoncontentContentDiffOp) | repeated | The diff itself for inline rendering |





### resources.common.content.ContentDiffOp
One diff operation, designed for inline client rendering.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `kind` | [Kind](#resourcescommoncontentKind) |  |  |
| `text` | [string](#string) |  | Plain text segment. Can contain whitespace and newlines. Client renders: - EQUAL: normal text - INSERT: highlighted - DELETE: strikethrough or hidden behind a toggle |





### resources.common.content.ContentDiffStats
Optional stats to quickly show "what changed" without parsing ops.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inserted_runes` | [uint32](#uint32) |  | Counts are measured in runes/codepoints on the server side. |
| `deleted_runes` | [uint32](#uint32) |  |  |
| `op_count` | [uint32](#uint32) |  | Optional: number of diff ops (after cleanup/coalescing). |




 <!-- end messages -->


### resources.common.content.Kind

| Name | Number | Description |
| ---- | ------ | ----------- |
| `KIND_UNSPECIFIED` | 0 |  |
| `KIND_EQUAL` | 1 |  |
| `KIND_INSERT` | 2 |  |
| `KIND_DELETE` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/database/database.proto


### resources.common.database.DateRange
DateRange represents a datetime range (uses Timestamp underneath) It depends on the API method if it will use date or date + time.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `start` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Start time |
| `end` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | End time |





### resources.common.database.PaginationRequest
Pagination for requests to the server



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `offset` | [int64](#int64) |  |  |
| `page_size` | [int64](#int64) | optional |  |





### resources.common.database.PaginationResponse
Server Pagination Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total_count` | [int64](#int64) |  |  |
| `offset` | [int64](#int64) |  |  |
| `end` | [int64](#int64) |  |  |
| `page_size` | [int64](#int64) |  |  |





### resources.common.database.Sort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `columns` | [SortByColumn](#resourcescommondatabaseSortByColumn) | repeated |  |





### resources.common.database.SortByColumn
SortByColumn sort by column and direction



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | ID is the column name. |
| `desc` | [bool](#bool) |  | Desc if true sorts descending, ascending otherwise. |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/i18n.proto


### resources.common.I18NItem
Wrapped translated message for the client



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `parameters` | [I18NItem.ParametersEntry](#resourcescommonI18NItemParametersEntry) | repeated |  |





### resources.common.I18NItem.ParametersEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/error.proto


### resources.common.Error


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [I18NItem](#resourcescommonI18NItem) | optional |  |
| `content` | [I18NItem](#resourcescommonI18NItem) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/id_mapping.proto


### resources.common.IDMapping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/tests/objects.proto


### resources.common.tests.SimpleObject
INTERNAL ONLY** SimpleObject is used as a test object where proto-based messages are required.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `field1` | [string](#string) |  |  |
| `field2` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/cron/cron.proto


### resources.cron.Cronjob


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Cronjob name |
| `schedule` | [string](#string) |  | Cron schedule expression For available valid expressions, see [adhocore/gronx - Cron Expressions Documentation](https://github.com/adhocore/gronx/blob/fea40e3e90e70476877cfb9b50fac10c7de41c5c/README.md#cron-expression).<br/><br/>To generate Cronjob schedule expressions, you can also use web tools like https://crontab.guru/. |
| `state` | [CronjobState](#resourcescronCronjobState) |  | Cronjob state |
| `next_schedule_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Next time the cronjob should be run |
| `last_attempt_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional | Last attempted start time of Cronjob |
| `started_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional | Time current cronjob was started |
| `timeout` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) | optional | Optional timeout for cronjob execution |
| `data` | [CronjobData](#resourcescronCronjobData) |  | Cronjob data |
| `last_completed_event` | [CronjobCompletedEvent](#resourcescronCronjobCompletedEvent) | optional | Last event info to ease debugging and tracking |





### resources.cron.CronjobCompletedEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Cronjob name |
| `success` | [bool](#bool) |  | Cronjob execution success status |
| `cancelled` | [bool](#bool) |  | Cronjob execution was cancelled |
| `end_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Cronjob end time |
| `elapsed` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  | Cronjob execution time/elapsed time |
| `data` | [CronjobData](#resourcescronCronjobData) | optional | Cronjob data (can be empty if not touched by the Cronjob handler) |
| `node_name` | [string](#string) |  | Name of the node where the cronjob was executed |
| `error_message` | [string](#string) | optional | Error message (if success = false) |





### resources.cron.CronjobData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `data` | [google.protobuf.Any](https://protobuf.dev/reference/protobuf/google.protobuf/#any) | optional |  |





### resources.cron.CronjobLockOwnerState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hostname` | [string](#string) |  | Hostname of the agent the cronjob is running on |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |





### resources.cron.CronjobSchedulerEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cronjob` | [Cronjob](#resourcescronCronjob) |  | Full Cronjob spec |





### resources.cron.GenericCronData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attributes` | [GenericCronData.AttributesEntry](#resourcescronGenericCronDataAttributesEntry) | repeated |  |





### resources.cron.GenericCronData.AttributesEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |




 <!-- end messages -->


### resources.cron.CronjobState
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



## resources/discord/discord.proto


### resources.discord.Channel


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `guild_id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `type` | [uint32](#uint32) |  |  |
| `position` | [int64](#int64) |  |  |





### resources.discord.Guild


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `icon` | [string](#string) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/access/access.proto


### resources.documents.access.DocumentAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [DocumentJobAccess](#resourcesdocumentsaccessDocumentJobAccess) | repeated |  |
| `users` | [DocumentUserAccess](#resourcesdocumentsaccessDocumentUserAccess) | repeated |  |





### resources.documents.access.DocumentJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesdocumentsaccessAccessLevel) |  |  |
| `required` | [bool](#bool) | optional |  |





### resources.documents.access.DocumentUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `access` | [AccessLevel](#resourcesdocumentsaccessAccessLevel) |  |  |
| `required` | [bool](#bool) | optional |  |




 <!-- end messages -->


### resources.documents.access.AccessLevel

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



## resources/documents/activity/activity.proto


### resources.documents.activity.DocAccessJobsDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [resources.documents.access.DocumentJobAccess](#resourcesdocumentsaccessDocumentJobAccess) | repeated |  |
| `to_update` | [resources.documents.access.DocumentJobAccess](#resourcesdocumentsaccessDocumentJobAccess) | repeated |  |
| `to_delete` | [resources.documents.access.DocumentJobAccess](#resourcesdocumentsaccessDocumentJobAccess) | repeated |  |





### resources.documents.activity.DocAccessRequested


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `level` | [resources.documents.access.AccessLevel](#resourcesdocumentsaccessAccessLevel) |  |  |





### resources.documents.activity.DocAccessUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [DocAccessJobsDiff](#resourcesdocumentsactivityDocAccessJobsDiff) |  |  |
| `users` | [DocAccessUsersDiff](#resourcesdocumentsactivityDocAccessUsersDiff) |  |  |





### resources.documents.activity.DocAccessUsersDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [resources.documents.access.DocumentUserAccess](#resourcesdocumentsaccessDocumentUserAccess) | repeated |  |
| `to_update` | [resources.documents.access.DocumentUserAccess](#resourcesdocumentsaccessDocumentUserAccess) | repeated |  |
| `to_delete` | [resources.documents.access.DocumentUserAccess](#resourcesdocumentsaccessDocumentUserAccess) | repeated |  |





### resources.documents.activity.DocActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `document_id` | [int64](#int64) |  |  |
| `activity_type` | [DocActivityType](#resourcesdocumentsactivityDocActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [DocActivityData](#resourcesdocumentsactivityDocActivityData) |  |  |





### resources.documents.activity.DocActivityData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [DocUpdated](#resourcesdocumentsactivityDocUpdated) |  |  |
| `owner_changed` | [DocOwnerChanged](#resourcesdocumentsactivityDocOwnerChanged) |  |  |
| `access_updated` | [DocAccessUpdated](#resourcesdocumentsactivityDocAccessUpdated) |  |  |
| `access_requested` | [DocAccessRequested](#resourcesdocumentsactivityDocAccessRequested) |  |  |
| `signing_requested` | [DocSigningRequested](#resourcesdocumentsactivityDocSigningRequested) |  |  |





### resources.documents.activity.DocFilesChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [int64](#int64) |  |  |
| `deleted` | [int64](#int64) |  |  |





### resources.documents.activity.DocOwnerChanged


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_owner_id` | [int32](#int32) |  |  |
| `new_owner` | [resources.users.short.UserShort](#resourcesusersshortUserShort) |  |  |





### resources.documents.activity.DocSigningRequested


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deadline` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `approvers` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | repeated |  |





### resources.documents.activity.DocUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title_diff` | [string](#string) | optional |  |
| `title_cdiff` | [resources.common.content.ContentDiff](#resourcescommoncontentContentDiff) | optional |  |
| `content_diff` | [string](#string) | optional |  |
| `content_cdiff` | [resources.common.content.ContentDiff](#resourcescommoncontentContentDiff) | optional |  |
| `state_diff` | [string](#string) | optional |  |
| `state_cdiff` | [resources.common.content.ContentDiff](#resourcescommoncontentContentDiff) | optional |  |
| `files_change` | [DocFilesChange](#resourcesdocumentsactivityDocFilesChange) | optional |  |




 <!-- end messages -->


### resources.documents.activity.DocActivityType

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
| `DOC_ACTIVITY_TYPE_REQUESTED_APPROVAL` | 20 |  |
| `DOC_ACTIVITY_TYPE_REQUESTED_SIGNING` | 21 |  |
| `DOC_ACTIVITY_TYPE_APPROVAL_ASSIGNED` | 40 | Approval |
| `DOC_ACTIVITY_TYPE_APPROVAL_APPROVED` | 41 |  |
| `DOC_ACTIVITY_TYPE_APPROVAL_REJECTED` | 42 |  |
| `DOC_ACTIVITY_TYPE_APPROVAL_REVOKED` | 43 |  |
| `DOC_ACTIVITY_TYPE_APPROVAL_REMOVED` | 44 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/category/category.proto


### resources.documents.category.Category


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `job` | [string](#string) | optional |  |
| `color` | [string](#string) | optional |  |
| `icon` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/pins/pins.proto


### resources.documents.pins.DocumentPin


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `job` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `state` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/workflow/workflow.proto


### resources.documents.workflow.AutoCloseSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `message` | [string](#string) |  |  |





### resources.documents.workflow.Reminder


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `message` | [string](#string) |  |  |





### resources.documents.workflow.ReminderSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reminders` | [Reminder](#resourcesdocumentsworkflowReminder) | repeated |  |
| `max_reminder_count` | [int32](#int32) |  |  |





### resources.documents.workflow.Workflow


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reminder` | [bool](#bool) |  |  |
| `reminder_settings` | [ReminderSettings](#resourcesdocumentsworkflowReminderSettings) |  |  |
| `auto_close` | [bool](#bool) |  |  |
| `auto_close_settings` | [AutoCloseSettings](#resourcesdocumentsworkflowAutoCloseSettings) |  |  |





### resources.documents.workflow.WorkflowCronData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_doc_id` | [int64](#int64) |  |  |





### resources.documents.workflow.WorkflowState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `next_reminder_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `next_reminder_count` | [int32](#int32) | optional |  |
| `reminder_count` | [int32](#int32) |  |  |
| `auto_close_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `workflow` | [Workflow](#resourcesdocumentsworkflowWorkflow) | optional |  |





### resources.documents.workflow.WorkflowUserState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `manual_reminder_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `manual_reminder_message` | [string](#string) | optional |  |
| `reminder_count` | [int32](#int32) |  |  |
| `max_reminder_count` | [int32](#int32) |  |  |
| `workflow` | [Workflow](#resourcesdocumentsworkflowWorkflow) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/documents.proto


### resources.documents.Document


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category_id` | [int64](#int64) | optional |  |
| `category` | [category.Category](#resourcesdocumentscategoryCategory) | optional |  |
| `title` | [string](#string) |  |  |
| `word_count` | [uint32](#uint32) | optional | Derived field (server authored) |
| `first_heading` | [string](#string) | optional | Derived field (server authored) |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `data` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `meta` | [DocumentMeta](#resourcesdocumentsDocumentMeta) |  |  |
| `template_id` | [int64](#int64) | optional |  |
| `pin` | [pins.DocumentPin](#resourcesdocumentspinsDocumentPin) | optional |  |
| `workflow_state` | [workflow.WorkflowState](#resourcesdocumentsworkflowWorkflowState) | optional |  |
| `workflow_user` | [workflow.WorkflowUserState](#resourcesdocumentsworkflowWorkflowUserState) | optional |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### resources.documents.DocumentMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `recomputed_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `state` | [string](#string) |  |  |
| `approved` | [bool](#bool) | optional | Overall aggregates - At least one approval policy fully satisfied |
| `ap_required_total` | [int32](#int32) | optional | Approval rollups Total approvals needed across policies |
| `ap_collected_approved` | [int32](#int32) | optional | Approvals collected |
| `ap_required_remaining` | [int32](#int32) | optional | How many left to satisfy |
| `ap_declined_count` | [int32](#int32) | optional | Number of declines |
| `ap_pending_count` | [int32](#int32) | optional | Tasks still pending (optional) |
| `ap_any_declined` | [bool](#bool) | optional | Quick flag if any declines |
| `ap_policies_active` | [int32](#int32) | optional | Number of active approval policies |





### resources.documents.DocumentShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category_id` | [int64](#int64) | optional |  |
| `category` | [category.Category](#resourcesdocumentscategoryCategory) | optional |  |
| `title` | [string](#string) |  |  |
| `word_count` | [uint32](#uint32) | optional | Derived field (server authored) |
| `first_heading` | [string](#string) | optional | Derived field (server authored) |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `meta` | [DocumentMeta](#resourcesdocumentsDocumentMeta) |  |  |
| `pin` | [pins.DocumentPin](#resourcesdocumentspinsDocumentPin) | optional |  |
| `workflow_state` | [workflow.WorkflowState](#resourcesdocumentsworkflowWorkflowState) | optional |  |
| `workflow_user` | [workflow.WorkflowUserState](#resourcesdocumentsworkflowWorkflowUserState) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/stamps/stamp.proto


### resources.documents.stamps.Stamp


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `name` | [string](#string) |  |  |
| `svg_template` | [string](#string) |  | Parameterized SVG with slots |
| `access` | [StampAccess](#resourcesdocumentsstampsStampAccess) |  |  |





### resources.documents.stamps.StampAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [StampJobAccess](#resourcesdocumentsstampsStampJobAccess) | repeated |  |





### resources.documents.stamps.StampJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [StampAccessLevel](#resourcesdocumentsstampsStampAccessLevel) |  |  |




 <!-- end messages -->


### resources.documents.stamps.StampAccessLevel

| Name | Number | Description |
| ---- | ------ | ----------- |
| `STAMP_ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `STAMP_ACCESS_LEVEL_BLOCKED` | 1 |  |
| `STAMP_ACCESS_LEVEL_USE` | 2 |  |
| `STAMP_ACCESS_LEVEL_MANAGE` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/approval/approval.proto


### resources.documents.approval.Approval


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `document_id` | [int64](#int64) |  |  |
| `snapshot_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `task_id` | [int64](#int64) | optional | Link to originating task (if any) |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `user_job` | [string](#string) | optional |  |
| `user_job_label` | [string](#string) | optional |  |
| `user_grade` | [int32](#int32) | optional |  |
| `user_grade_label` | [string](#string) | optional |  |
| `payload_svg` | [string](#string) | optional | SVG path, typed preview, stamp fill, etc. |
| `stamp_id` | [int64](#int64) | optional |  |
| `stamp` | [resources.documents.stamps.Stamp](#resourcesdocumentsstampsStamp) | optional |  |
| `status` | [ApprovalStatus](#resourcesdocumentsapprovalApprovalStatus) |  |  |
| `comment` | [string](#string) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `revoked_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.documents.approval.ApprovalPolicy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `snapshot_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `on_edit_behavior` | [OnEditBehavior](#resourcesdocumentsapprovalOnEditBehavior) |  |  |
| `rule_kind` | [ApprovalRuleKind](#resourcesdocumentsapprovalApprovalRuleKind) |  |  |
| `required_count` | [int32](#int32) | optional |  |
| `signature_required` | [bool](#bool) |  |  |
| `self_approve_allowed` | [bool](#bool) |  |  |
| `assigned_count` | [int32](#int32) |  |  |
| `approved_count` | [int32](#int32) |  |  |
| `declined_count` | [int32](#int32) |  |  |
| `pending_count` | [int32](#int32) |  |  |
| `any_declined` | [bool](#bool) |  |  |
| `started_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `completed_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.documents.approval.ApprovalTask


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `document_id` | [int64](#int64) |  |  |
| `snapshot_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `assignee_kind` | [ApprovalAssigneeKind](#resourcesdocumentsapprovalApprovalAssigneeKind) |  |  |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) | optional |  |
| `job_grade_label` | [string](#string) | optional |  |
| `label` | [string](#string) | optional | "Leadership", "Counterparty Rep" |
| `signature_required` | [bool](#bool) |  |  |
| `slot_no` | [int32](#int32) |  | >=1; meaningful only for Job tasks; always 1 for User |
| `status` | [ApprovalTaskStatus](#resourcesdocumentsapprovalApprovalTaskStatus) |  |  |
| `comment` | [string](#string) | optional | Optional comment on approve/decline |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `completed_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `due_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `decision_count` | [int32](#int32) |  |  |
| `approval_id` | [int64](#int64) | optional |  |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `document` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |




 <!-- end messages -->


### resources.documents.approval.ApprovalAssigneeKind

| Name | Number | Description |
| ---- | ------ | ----------- |
| `APPROVAL_ASSIGNEE_KIND_UNSPECIFIED` | 0 |  |
| `APPROVAL_ASSIGNEE_KIND_USER` | 1 |  |
| `APPROVAL_ASSIGNEE_KIND_JOB_GRADE` | 2 |  |



### resources.documents.approval.ApprovalRuleKind

| Name | Number | Description |
| ---- | ------ | ----------- |
| `APPROVAL_RULE_KIND_UNSPECIFIED` | 0 |  |
| `APPROVAL_RULE_KIND_REQUIRE_ALL` | 1 | Require all approvals (based on tasks created) |
| `APPROVAL_RULE_KIND_QUORUM_ANY` | 2 | Quorum of N approvals based on `required_count` |



### resources.documents.approval.ApprovalStatus

| Name | Number | Description |
| ---- | ------ | ----------- |
| `APPROVAL_STATUS_UNSPECIFIED` | 0 |  |
| `APPROVAL_STATUS_APPROVED` | 1 |  |
| `APPROVAL_STATUS_DECLINED` | 2 |  |
| `APPROVAL_STATUS_REVOKED` | 3 |  |



### resources.documents.approval.ApprovalTaskStatus

| Name | Number | Description |
| ---- | ------ | ----------- |
| `APPROVAL_TASK_STATUS_UNSPECIFIED` | 0 |  |
| `APPROVAL_TASK_STATUS_PENDING` | 1 |  |
| `APPROVAL_TASK_STATUS_APPROVED` | 2 |  |
| `APPROVAL_TASK_STATUS_DECLINED` | 3 |  |
| `APPROVAL_TASK_STATUS_EXPIRED` | 4 |  |
| `APPROVAL_TASK_STATUS_CANCELLED` | 5 |  |
| `APPROVAL_TASK_STATUS_COMPLETED` | 6 |  |



### resources.documents.approval.OnEditBehavior
Policy snapshot applied to a specific version


| Name | Number | Description |
| ---- | ------ | ----------- |
| `ON_EDIT_BEHAVIOR_UNSPECIFIED` | 0 |  |
| `ON_EDIT_BEHAVIOR_KEEP_PROGRESS` | 1 | Keep approvals where possible |
| `ON_EDIT_BEHAVIOR_RESET` | 2 | Reset review on content edits |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/comment/comment.proto


### resources.documents.comment.Comment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `document_id` | [int64](#int64) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/forms/forms.proto


### resources.documents.forms.Form


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/references/references.proto


### resources.documents.references.DocumentReference


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `source_document_id` | [int64](#int64) |  |  |
| `source_document` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |
| `reference` | [DocReference](#resourcesdocumentsreferencesDocReference) |  |  |
| `target_document_id` | [int64](#int64) |  |  |
| `target_document` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |




 <!-- end messages -->


### resources.documents.references.DocReference

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_REFERENCE_UNSPECIFIED` | 0 |  |
| `DOC_REFERENCE_LINKED` | 1 |  |
| `DOC_REFERENCE_SOLVES` | 2 |  |
| `DOC_REFERENCE_CLOSES` | 3 |  |
| `DOC_REFERENCE_DEPRECATES` | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/relations/relations.proto


### resources.documents.relations.DocumentRelation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `document_id` | [int64](#int64) |  |  |
| `document` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |
| `source_user_id` | [int32](#int32) |  |  |
| `source_user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `relation` | [DocRelation](#resourcesdocumentsrelationsDocRelation) |  |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |




 <!-- end messages -->


### resources.documents.relations.DocRelation

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_RELATION_UNSPECIFIED` | 0 |  |
| `DOC_RELATION_MENTIONED` | 1 |  |
| `DOC_RELATION_TARGETS` | 2 |  |
| `DOC_RELATION_CAUSED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/requests/requests.proto


### resources.documents.requests.DocRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `document_id` | [int64](#int64) |  |  |
| `request_type` | [resources.documents.activity.DocActivityType](#resourcesdocumentsactivityDocActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [resources.documents.activity.DocActivityData](#resourcesdocumentsactivityDocActivityData) |  |  |
| `accepted` | [bool](#bool) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/vehicles/props/props.proto


### resources.vehicles.props.VehicleProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plate` | [string](#string) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `wanted` | [bool](#bool) | optional |  |
| `wanted_reason` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/vehicles/vehicles.proto


### resources.vehicles.Vehicle


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plate` | [string](#string) |  |  |
| `model` | [string](#string) | optional |  |
| `type` | [string](#string) |  |  |
| `owner_id` | [int32](#int32) | optional |  |
| `owner_identifier` | [string](#string) | optional |  |
| `owner` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |
| `props` | [props.VehicleProps](#resourcesvehiclespropsVehicleProps) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/templates/templates.proto


### resources.documents.templates.ObjectSpecs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `required` | [bool](#bool) | optional |  |
| `min` | [int32](#int32) | optional |  |
| `max` | [int32](#int32) | optional |  |





### resources.documents.templates.Template


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category` | [resources.documents.category.Category](#resourcesdocumentscategoryCategory) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `color` | [string](#string) | optional |  |
| `icon` | [string](#string) | optional |  |
| `content_title` | [string](#string) |  |  |
| `content` | [string](#string) |  |  |
| `state` | [string](#string) |  |  |
| `schema` | [TemplateSchema](#resourcesdocumentstemplatesTemplateSchema) |  |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `job_access` | [TemplateJobAccess](#resourcesdocumentstemplatesTemplateJobAccess) | repeated |  |
| `content_access` | [resources.documents.access.DocumentAccess](#resourcesdocumentsaccessDocumentAccess) |  |  |
| `workflow` | [resources.documents.workflow.Workflow](#resourcesdocumentsworkflowWorkflow) | optional |  |
| `approval` | [TemplateApproval](#resourcesdocumentstemplatesTemplateApproval) | optional |  |





### resources.documents.templates.TemplateApproval


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `policy` | [TemplateApprovalPolicy](#resourcesdocumentstemplatesTemplateApprovalPolicy) | optional |  |
| `tasks` | [TemplateApprovalTaskSeed](#resourcesdocumentstemplatesTemplateApprovalTaskSeed) | repeated |  |





### resources.documents.templates.TemplateApprovalPolicy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rule_kind` | [resources.documents.approval.ApprovalRuleKind](#resourcesdocumentsapprovalApprovalRuleKind) |  |  |
| `on_edit_behavior` | [resources.documents.approval.OnEditBehavior](#resourcesdocumentsapprovalOnEditBehavior) |  |  |
| `required_count` | [int32](#int32) | optional |  |
| `signature_required` | [bool](#bool) |  |  |
| `self_approve_allowed` | [bool](#bool) |  |  |





### resources.documents.templates.TemplateApprovalTaskSeed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  | If user_id == 0 -> JOB task |
| `minimum_grade` | [int32](#int32) |  |  |
| `label` | [string](#string) | optional | Label of task |
| `signature_required` | [bool](#bool) |  |  |
| `slots` | [int32](#int32) |  | Only for JOB tasks; number of PENDING slots to ensure (>=1) |
| `due_in_days` | [int32](#int32) | optional | Optional default due date for created slots |
| `comment` | [string](#string) | optional | Optional note set on created tasks |





### resources.documents.templates.TemplateData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `active_char` | [resources.users.User](#resourcesusersUser) |  |  |
| `documents` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | repeated |  |
| `users` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | repeated |  |
| `vehicles` | [resources.vehicles.Vehicle](#resourcesvehiclesVehicle) | repeated |  |





### resources.documents.templates.TemplateJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [resources.documents.access.AccessLevel](#resourcesdocumentsaccessAccessLevel) |  |  |





### resources.documents.templates.TemplateRequirements


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `documents` | [ObjectSpecs](#resourcesdocumentstemplatesObjectSpecs) | optional |  |
| `users` | [ObjectSpecs](#resourcesdocumentstemplatesObjectSpecs) | optional |  |
| `vehicles` | [ObjectSpecs](#resourcesdocumentstemplatesObjectSpecs) | optional |  |





### resources.documents.templates.TemplateSchema


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `requirements` | [TemplateRequirements](#resourcesdocumentstemplatesTemplateRequirements) |  |  |





### resources.documents.templates.TemplateShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category` | [resources.documents.category.Category](#resourcesdocumentscategoryCategory) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `color` | [string](#string) | optional |  |
| `icon` | [string](#string) | optional |  |
| `schema` | [TemplateSchema](#resourcesdocumentstemplatesTemplateSchema) |  |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `workflow` | [resources.documents.workflow.Workflow](#resourcesdocumentsworkflowWorkflow) | optional |  |





### resources.documents.templates.TemplateUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/file/filestore.proto


### resources.file.DeleteFileRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [int64](#int64) |  |  |
| `file_id` | [int64](#int64) |  |  |





### resources.file.DeleteFileResponse





### resources.file.UploadFileRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `meta` | [UploadMeta](#resourcesfileUploadMeta) |  |  |
| `data` | [bytes](#bytes) |  | Raw bytes <= 128 KiB each, browsers should only read 64 KiB at a time, but this is a buffer just in case |





### resources.file.UploadFileResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  | Unique ID for the uploaded file |
| `url` | [string](#string) |  | URL to the uploaded file |
| `file` | [File](#resourcesfileFile) |  | File info |





### resources.file.UploadMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [int64](#int64) |  |  |
| `namespace` | [string](#string) |  | "documents", "wiki", … |
| `original_name` | [string](#string) |  |  |
| `content_type` | [string](#string) |  | optional - server re-validates |
| `size` | [int64](#int64) |  | Size in bytes |
| `reason` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/grpcws/grpcws.proto


### resources.grpcws.Body


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |
| `complete` | [bool](#bool) |  |  |





### resources.grpcws.Cancel





### resources.grpcws.Complete





### resources.grpcws.Failure


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `error_message` | [string](#string) |  |  |
| `error_status` | [string](#string) |  |  |
| `headers` | [Failure.HeadersEntry](#resourcesgrpcwsFailureHeadersEntry) | repeated |  |





### resources.grpcws.Failure.HeadersEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [HeaderValue](#resourcesgrpcwsHeaderValue) |  |  |





### resources.grpcws.GrpcFrame


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stream_id` | [uint32](#uint32) |  |  |
| `ping` | [Ping](#resourcesgrpcwsPing) |  |  |
| `header` | [Header](#resourcesgrpcwsHeader) |  |  |
| `body` | [Body](#resourcesgrpcwsBody) |  |  |
| `complete` | [Complete](#resourcesgrpcwsComplete) |  |  |
| `failure` | [Failure](#resourcesgrpcwsFailure) |  |  |
| `cancel` | [Cancel](#resourcesgrpcwsCancel) |  |  |





### resources.grpcws.Header


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operation` | [string](#string) |  |  |
| `headers` | [Header.HeadersEntry](#resourcesgrpcwsHeaderHeadersEntry) | repeated |  |
| `status` | [int32](#int32) |  |  |





### resources.grpcws.Header.HeadersEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [HeaderValue](#resourcesgrpcwsHeaderValue) |  |  |





### resources.grpcws.HeaderValue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [string](#string) | repeated |  |





### resources.grpcws.Ping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pong` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/colleagues/activity/activity.proto


### resources.jobs.colleagues.activity.AbsenceDateChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_begin` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `absence_end` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |





### resources.jobs.colleagues.activity.ColleagueActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `source_user_id` | [int32](#int32) | optional |  |
| `source_user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) |  |  |
| `activity_type` | [ColleagueActivityType](#resourcesjobscolleaguesactivityColleagueActivityType) |  |  |
| `reason` | [string](#string) |  |  |
| `data` | [ColleagueActivityData](#resourcesjobscolleaguesactivityColleagueActivityData) |  |  |





### resources.jobs.colleagues.activity.ColleagueActivityData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_date` | [AbsenceDateChange](#resourcesjobscolleaguesactivityAbsenceDateChange) |  |  |
| `grade_change` | [GradeChange](#resourcesjobscolleaguesactivityGradeChange) |  |  |
| `labels_change` | [LabelsChange](#resourcesjobscolleaguesactivityLabelsChange) |  |  |
| `name_change` | [NameChange](#resourcesjobscolleaguesactivityNameChange) |  |  |





### resources.jobs.colleagues.activity.GradeChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grade` | [int32](#int32) |  |  |
| `grade_label` | [string](#string) |  |  |





### resources.jobs.colleagues.activity.LabelsChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [resources.jobs.labels.Label](#resourcesjobslabelsLabel) | repeated |  |
| `removed` | [resources.jobs.labels.Label](#resourcesjobslabelsLabel) | repeated |  |





### resources.jobs.colleagues.activity.NameChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prefix` | [string](#string) | optional |  |
| `suffix` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.jobs.colleagues.activity.ColleagueActivityType

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



## resources/jobs/conduct/conduct.proto


### resources.jobs.conduct.ConductEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `type` | [ConductType](#resourcesjobsconductConductType) |  |  |
| `draft` | [bool](#bool) |  |  |
| `message` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |




 <!-- end messages -->


### resources.jobs.conduct.ConductType

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



## resources/jobs/settings/settings.proto


### resources.jobs.settings.DiscordSyncChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `plan` | [string](#string) |  |  |





### resources.jobs.settings.DiscordSyncChanges


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `changes` | [DiscordSyncChange](#resourcesjobssettingsDiscordSyncChange) | repeated |  |





### resources.jobs.settings.DiscordSyncSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dry_run` | [bool](#bool) |  |  |
| `user_info_sync` | [bool](#bool) |  |  |
| `user_info_sync_settings` | [UserInfoSyncSettings](#resourcesjobssettingsUserInfoSyncSettings) |  |  |
| `status_log` | [bool](#bool) |  |  |
| `status_log_settings` | [StatusLogSettings](#resourcesjobssettingsStatusLogSettings) |  |  |
| `jobs_absence` | [bool](#bool) |  |  |
| `jobs_absence_settings` | [JobsAbsenceSettings](#resourcesjobssettingsJobsAbsenceSettings) |  |  |
| `group_sync_settings` | [GroupSyncSettings](#resourcesjobssettingsGroupSyncSettings) |  |  |
| `qualifications_role_format` | [string](#string) |  |  |





### resources.jobs.settings.GroupMapping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `from_grade` | [int32](#int32) |  |  |
| `to_grade` | [int32](#int32) |  |  |





### resources.jobs.settings.GroupSyncSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ignored_role_ids` | [string](#string) | repeated |  |





### resources.jobs.settings.JobSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_past_days` | [int32](#int32) |  |  |
| `absence_future_days` | [int32](#int32) |  |  |





### resources.jobs.settings.JobsAbsenceSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_role` | [string](#string) |  |  |





### resources.jobs.settings.StatusLogSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel_id` | [string](#string) |  |  |





### resources.jobs.settings.UserInfoSyncSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `employee_role_enabled` | [bool](#bool) |  |  |
| `employee_role_format` | [string](#string) |  |  |
| `grade_role_format` | [string](#string) |  |  |
| `unemployed_enabled` | [bool](#bool) |  |  |
| `unemployed_mode` | [UserInfoSyncUnemployedMode](#resourcesjobssettingsUserInfoSyncUnemployedMode) |  |  |
| `unemployed_role_name` | [string](#string) |  |  |
| `sync_nicknames` | [bool](#bool) |  |  |
| `group_mapping` | [GroupMapping](#resourcesjobssettingsGroupMapping) | repeated |  |




 <!-- end messages -->


### resources.jobs.settings.UserInfoSyncUnemployedMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_UNSPECIFIED` | 0 |  |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE` | 1 |  |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_KICK` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/props/props.proto


### resources.jobs.props.JobProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `livemap_marker_color` | [string](#string) |  |  |
| `quick_buttons` | [QuickButtons](#resourcesjobspropsQuickButtons) |  |  |
| `radio_frequency` | [string](#string) | optional |  |
| `discord_guild_id` | [string](#string) | optional |  |
| `discord_last_sync` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `discord_sync_settings` | [resources.jobs.settings.DiscordSyncSettings](#resourcesjobssettingsDiscordSyncSettings) |  |  |
| `discord_sync_changes` | [resources.jobs.settings.DiscordSyncChanges](#resourcesjobssettingsDiscordSyncChanges) | optional |  |
| `motd` | [string](#string) | optional |  |
| `logo_file_id` | [int64](#int64) | optional |  |
| `logo_file` | [resources.file.File](#resourcesfileFile) | optional |  |
| `settings` | [resources.jobs.settings.JobSettings](#resourcesjobssettingsJobSettings) |  |  |





### resources.jobs.props.QuickButtons


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `penalty_calculator` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/timeclock/timeclock.proto


### resources.jobs.timeclock.TimeclockEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | optional |  |
| `start_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `end_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `spent_time` | [float](#float) |  |  |





### resources.jobs.timeclock.TimeclockStats


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `spent_time_sum` | [float](#float) |  |  |
| `spent_time_avg` | [float](#float) |  |  |
| `spent_time_max` | [float](#float) |  |  |





### resources.jobs.timeclock.TimeclockWeeklyStats


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `year` | [int32](#int32) |  |  |
| `calendar_week` | [int32](#int32) |  |  |
| `sum` | [float](#float) |  |  |
| `avg` | [float](#float) |  |  |
| `max` | [float](#float) |  |  |




 <!-- end messages -->


### resources.jobs.timeclock.TimeclockMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `TIMECLOCK_MODE_UNSPECIFIED` | 0 |  |
| `TIMECLOCK_MODE_DAILY` | 1 |  |
| `TIMECLOCK_MODE_WEEKLY` | 2 |  |
| `TIMECLOCK_MODE_RANGE` | 3 |  |
| `TIMECLOCK_MODE_TIMELINE` | 4 |  |



### resources.jobs.timeclock.TimeclockViewMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `TIMECLOCK_VIEW_MODE_UNSPECIFIED` | 0 |  |
| `TIMECLOCK_VIEW_MODE_SELF` | 1 |  |
| `TIMECLOCK_VIEW_MODE_ALL` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/laws/laws.proto


### resources.laws.Law


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `lawbook_id` | [int64](#int64) |  |  |
| `name` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `hint` | [string](#string) | optional |  |
| `fine` | [uint32](#uint32) | optional |  |
| `detention_time` | [uint32](#uint32) | optional |  |
| `stvo_points` | [uint32](#uint32) | optional |  |





### resources.laws.LawBook


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `laws` | [Law](#resourceslawsLaw) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/livemap/coords.proto


### resources.livemap.Coords


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/livemap/heatmap/heatmap.proto


### resources.livemap.heatmap.HeatmapEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `w` | [double](#double) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/livemap/markers/marker_marker.proto


### resources.livemap.markers.CircleMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `radius` | [int32](#int32) |  |  |
| `opacity` | [float](#float) | optional |  |





### resources.livemap.markers.IconMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `icon` | [string](#string) |  |  |





### resources.livemap.markers.MarkerData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `circle` | [CircleMarker](#resourceslivemapmarkersCircleMarker) |  |  |
| `icon` | [IconMarker](#resourceslivemapmarkersIconMarker) |  |  |





### resources.livemap.markers.MarkerMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `description` | [string](#string) | optional |  |
| `postal` | [string](#string) | optional |  |
| `color` | [string](#string) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) |  |  |
| `type` | [MarkerType](#resourceslivemapmarkersMarkerType) |  |  |
| `data` | [MarkerData](#resourceslivemapmarkersMarkerData) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |




 <!-- end messages -->


### resources.livemap.markers.MarkerType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `MARKER_TYPE_UNSPECIFIED` | 0 |  |
| `MARKER_TYPE_DOT` | 1 |  |
| `MARKER_TYPE_CIRCLE` | 2 |  |
| `MARKER_TYPE_ICON` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/livemap/markers/user_marker.proto


### resources.livemap.markers.UserMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `postal` | [string](#string) | optional |  |
| `color` | [string](#string) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) |  |  |
| `job_grade` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) |  |  |
| `unit_id` | [int64](#int64) | optional |  |
| `unit` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) | optional |  |
| `hidden` | [bool](#bool) |  |  |
| `data` | [UserMarkerData](#resourceslivemapmarkersUserMarkerData) | optional |  |





### resources.livemap.markers.UserMarkerData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `is_in_vehicle` | [bool](#bool) |  |  |
| `vehicle_plate` | [string](#string) | optional |  |
| `vehicle_updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/access/access.proto


### resources.mailer.access.Access


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [JobAccess](#resourcesmaileraccessJobAccess) | repeated |  |
| `users` | [UserAccess](#resourcesmaileraccessUserAccess) | repeated |  |
| `qualifications` | [QualificationAccess](#resourcesmaileraccessQualificationAccess) | repeated |  |





### resources.mailer.access.JobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesmaileraccessAccessLevel) |  |  |





### resources.mailer.access.QualificationAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `qualification_id` | [int64](#int64) |  |  |
| `qualification` | [resources.qualifications.QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `access` | [AccessLevel](#resourcesmaileraccessAccessLevel) |  |  |





### resources.mailer.access.UserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `access` | [AccessLevel](#resourcesmaileraccessAccessLevel) |  |  |




 <!-- end messages -->


### resources.mailer.access.AccessLevel

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



## resources/mailer/settings/settings.proto


### resources.mailer.settings.EmailSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `signature` | [resources.common.content.Content](#resourcescommoncontentContent) | optional |  |
| `blocked_emails` | [string](#string) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/emails/email.proto


### resources.mailer.emails.Email


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deactivated` | [bool](#bool) |  |  |
| `job` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `email` | [string](#string) |  |  |
| `email_changed` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `label` | [string](#string) | optional |  |
| `access` | [resources.mailer.access.Access](#resourcesmaileraccessAccess) |  |  |
| `settings` | [resources.mailer.settings.EmailSettings](#resourcesmailersettingsEmailSettings) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/messages/message.proto


### resources.mailer.messages.Message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `thread_id` | [int64](#int64) |  |  |
| `sender_id` | [int64](#int64) |  |  |
| `sender` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `data` | [MessageData](#resourcesmailermessagesMessageData) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |





### resources.mailer.messages.MessageAttachment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [MessageAttachmentDocument](#resourcesmailermessagesMessageAttachmentDocument) |  |  |





### resources.mailer.messages.MessageAttachmentDocument


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `title` | [string](#string) | optional |  |





### resources.mailer.messages.MessageData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attachments` | [MessageAttachment](#resourcesmailermessagesMessageAttachment) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/threads/thread.proto


### resources.mailer.threads.Thread


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `creator_email_id` | [int64](#int64) |  |  |
| `creator_email` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `title` | [string](#string) |  |  |
| `recipients` | [ThreadRecipientEmail](#resourcesmailerthreadsThreadRecipientEmail) | repeated |  |
| `state` | [ThreadState](#resourcesmailerthreadsThreadState) | optional |  |





### resources.mailer.threads.ThreadRecipientEmail


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `email_id` | [int64](#int64) |  |  |
| `email` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) | optional |  |





### resources.mailer.threads.ThreadState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread_id` | [int64](#int64) |  |  |
| `email_id` | [int64](#int64) |  |  |
| `last_read` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `unread` | [bool](#bool) | optional |  |
| `important` | [bool](#bool) | optional |  |
| `favorite` | [bool](#bool) | optional |  |
| `muted` | [bool](#bool) | optional |  |
| `archived` | [bool](#bool) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/events/events.proto


### resources.mailer.events.MailerEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_update` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) |  |  |
| `email_delete` | [int64](#int64) |  |  |
| `email_settings_updated` | [resources.mailer.settings.EmailSettings](#resourcesmailersettingsEmailSettings) |  |  |
| `thread_update` | [resources.mailer.threads.Thread](#resourcesmailerthreadsThread) |  |  |
| `thread_delete` | [int64](#int64) |  |  |
| `thread_state_update` | [resources.mailer.threads.ThreadState](#resourcesmailerthreadsThreadState) |  |  |
| `message_update` | [resources.mailer.messages.Message](#resourcesmailermessagesMessage) |  |  |
| `message_delete` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/templates/template.proto


### resources.mailer.templates.Template


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `email_id` | [int64](#int64) |  |  |
| `title` | [string](#string) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/notifications/clientview/clientview.proto


### resources.notifications.clientview.ClientView


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [ObjectType](#resourcesnotificationsclientviewObjectType) |  |  |
| `id` | [int64](#int64) | optional |  |





### resources.notifications.clientview.ObjectEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [ObjectType](#resourcesnotificationsclientviewObjectType) |  |  |
| `id` | [int64](#int64) | optional |  |
| `event_type` | [ObjectEventType](#resourcesnotificationsclientviewObjectEventType) |  |  |
| `user_id` | [int32](#int32) | optional |  |
| `job` | [string](#string) | optional |  |
| `data` | [google.protobuf.Struct](https://protobuf.dev/reference/protobuf/google.protobuf/#struct) | optional |  |




 <!-- end messages -->


### resources.notifications.clientview.ObjectEventType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `OBJECT_EVENT_TYPE_UNSPECIFIED` | 0 |  |
| `OBJECT_EVENT_TYPE_UPDATED` | 1 |  |
| `OBJECT_EVENT_TYPE_DELETED` | 2 |  |



### resources.notifications.clientview.ObjectType

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



## resources/notifications/notifications.proto


### resources.notifications.CalendarData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [int64](#int64) | optional |  |
| `calendar_entry_id` | [int64](#int64) | optional |  |





### resources.notifications.Data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `link` | [Link](#resourcesnotificationsLink) | optional |  |
| `caused_by` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `calendar` | [CalendarData](#resourcesnotificationsCalendarData) | optional |  |





### resources.notifications.Link


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to` | [string](#string) |  |  |
| `title` | [string](#string) | optional |  |
| `external` | [bool](#bool) | optional |  |





### resources.notifications.Notification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `read_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `title` | [resources.common.I18NItem](#resourcescommonI18NItem) |  |  |
| `type` | [NotificationType](#resourcesnotificationsNotificationType) |  |  |
| `content` | [resources.common.I18NItem](#resourcescommonI18NItem) |  |  |
| `category` | [NotificationCategory](#resourcesnotificationsNotificationCategory) |  |  |
| `data` | [Data](#resourcesnotificationsData) | optional |  |
| `starred` | [bool](#bool) | optional |  |




 <!-- end messages -->


### resources.notifications.NotificationCategory

| Name | Number | Description |
| ---- | ------ | ----------- |
| `NOTIFICATION_CATEGORY_UNSPECIFIED` | 0 |  |
| `NOTIFICATION_CATEGORY_GENERAL` | 1 |  |
| `NOTIFICATION_CATEGORY_DOCUMENT` | 2 |  |
| `NOTIFICATION_CATEGORY_CALENDAR` | 3 |  |



### resources.notifications.NotificationType

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



## resources/userinfo/userinfo.proto


### resources.userinfo.PollReq
PollReq: published to `userinfo.poll.request` when an active user connects or requests a refresh.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [int64](#int64) |  | The account the user belongs to |
| `user_id` | [int32](#int32) |  | The unique user identifier within the account |





### resources.userinfo.UserInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [int64](#int64) |  |  |
| `enabled` | [bool](#bool) |  |  |
| `license` | [string](#string) |  |  |
| `last_char` | [int32](#int32) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `job_grade` | [int32](#int32) |  |  |
| `groups` | [resources.accounts.AccountGroups](#resourcesaccountsAccountGroups) | optional |  |
| `can_be_superuser` | [bool](#bool) |  |  |
| `superuser` | [bool](#bool) |  |  |





### resources.userinfo.UserInfoChanged
UserInfoChanged used to signal Job or JobGrade changes.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [int64](#int64) |  | The account the user belongs to |
| `user_id` | [int32](#int32) |  | The unique user identifier within the account |
| `changed_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Timestamp of when the change was detected |
| `old_job` | [string](#string) |  | Previous job title |
| `new_job` | [string](#string) | optional | New job title |
| `new_job_label` | [string](#string) | optional |  |
| `old_job_grade` | [int32](#int32) |  | Previous job grade |
| `new_job_grade` | [int32](#int32) | optional | New job grade |
| `new_job_grade_label` | [string](#string) | optional | New job grade label |
| `can_be_superuser` | [bool](#bool) | optional | Can the user be superuser (by group or license) |
| `superuser` | [bool](#bool) | optional | Superuser state |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/notifications/events/events.proto


### resources.notifications.events.JobEvent
Job related events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.props.JobProps](#resourcesjobspropsJobProps) |  |  |





### resources.notifications.events.JobGradeEvent
Job grade events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_token` | [bool](#bool) |  |  |





### resources.notifications.events.SystemEvent
System related events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_config` | [resources.clientconfig.ClientConfig](#resourcesclientconfigClientConfig) |  | Client configuration update (e.g., feature gates, game settings, banner message) |





### resources.notifications.events.UserEvent
User related events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_token` | [bool](#bool) |  |  |
| `notification` | [resources.notifications.Notification](#resourcesnotificationsNotification) |  | Notifications |
| `notifications_read_count` | [int64](#int64) |  |  |
| `user_info_changed` | [resources.userinfo.UserInfoChanged](#resourcesuserinfoUserInfoChanged) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/permissions/events/events.proto


### resources.permissions.events.JobLimitsUpdatedEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### resources.permissions.events.RoleIDEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/permissions/permissions/permissions.proto


### resources.permissions.permissions.PermItem


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `val` | [bool](#bool) |  |  |





### resources.permissions.permissions.Permission


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `guard_name` | [string](#string) |  |  |
| `val` | [bool](#bool) |  |  |
| `order` | [int32](#int32) | optional |  |
| `icon` | [string](#string) | optional |  |





### resources.permissions.permissions.Role


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `permissions` | [Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/perms.proto


### resources.settings.AttrsUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_update` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |
| `to_remove` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |





### resources.settings.PermsUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_update` | [resources.permissions.permissions.PermItem](#resourcespermissionspermissionsPermItem) | repeated |  |
| `to_remove` | [resources.permissions.permissions.PermItem](#resourcespermissionspermissionsPermItem) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/status.proto


### resources.settings.DBSyncStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `last_synced_data` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `last_synced_activity` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `last_dbsync_version` | [string](#string) | optional |  |





### resources.settings.Database


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `connected` | [bool](#bool) |  |  |
| `migration_version` | [uint64](#uint64) |  |  |
| `migration_dirty` | [bool](#bool) |  |  |
| `db_charset` | [string](#string) |  |  |
| `db_collation` | [string](#string) |  |  |
| `tables_ok` | [bool](#bool) |  |  |





### resources.settings.Nats


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `connected` | [bool](#bool) |  |  |





### resources.settings.NewVersionInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `url` | [string](#string) |  |  |
| `release_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.settings.SystemStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `database` | [Database](#resourcessettingsDatabase) |  |  |
| `nats` | [Nats](#resourcessettingsNats) |  |  |
| `dbsync` | [DBSyncStatus](#resourcessettingsDBSyncStatus) |  |  |
| `version` | [VersionStatus](#resourcessettingsVersionStatus) |  |  |





### resources.settings.VersionStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current` | [string](#string) |  |  |
| `new_version` | [NewVersionInfo](#resourcessettingsNewVersionInfo) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/stats/stats.proto


### resources.stats.Stat


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [int32](#int32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/sync/activity/activity.proto


### resources.sync.activity.AccountUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `license` | [string](#string) |  |  |
| `group` | [string](#string) | optional |  |
| `groups` | [string](#string) | repeated |  |





### resources.sync.activity.ColleagueProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reason` | [string](#string) | optional |  |
| `props` | [resources.jobs.colleagues.ColleagueProps](#resourcesjobscolleaguesColleagueProps) |  |  |





### resources.sync.activity.TimeclockUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `start` | [bool](#bool) |  |  |





### resources.sync.activity.UserOAuth2Conn
Connect an identifier/license to the provider with the specified external id (e.g., auto discord social connect on server join)



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider_name` | [string](#string) |  |  |
| `identifier` | [string](#string) |  |  |
| `external_id` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |





### resources.sync.activity.UserProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reason` | [string](#string) | optional |  |
| `props` | [resources.users.props.UserProps](#resourcesuserspropsUserProps) |  |  |





### resources.sync.activity.UserUpdate


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



## resources/sync/data/data.proto


### resources.sync.data.CitizenLocations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `job_grade` | [int32](#int32) | optional |  |
| `coords` | [resources.livemap.Coords](#resourceslivemapCoords) |  |  |
| `hidden` | [bool](#bool) |  |  |
| `remove` | [bool](#bool) |  |  |





### resources.sync.data.DataAccounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_updates` | [resources.sync.activity.AccountUpdate](#resourcessyncactivityAccountUpdate) | repeated |  |





### resources.sync.data.DataJobs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |





### resources.sync.data.DataLicenses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `licenses` | [resources.users.licenses.License](#resourcesuserslicensesLicense) | repeated |  |





### resources.sync.data.DataStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `count` | [int64](#int64) |  |  |





### resources.sync.data.DataUser
Detailed user information for sync purposes Should be kept inline with `resources.users.User`.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `identifier` | [string](#string) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `job_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `jobs` | [resources.users.UserJob](#resourcesusersUserJob) | repeated |  |
| `firstname` | [string](#string) |  |  |
| `lastname` | [string](#string) | optional |  |
| `dateofbirth` | [string](#string) |  |  |
| `sex` | [string](#string) | optional |  |
| `height` | [string](#string) | optional |  |
| `phone_number` | [string](#string) | optional |  |
| `phone_numbers` | [resources.users.PhoneNumber](#resourcesusersPhoneNumber) | repeated |  |
| `visum` | [int32](#int32) | optional |  |
| `playtime` | [int32](#int32) | optional |  |
| `props` | [resources.users.props.UserProps](#resourcesuserspropsUserProps) | optional |  |
| `licenses` | [resources.users.licenses.License](#resourcesuserslicensesLicense) | repeated |  |
| `profile_picture_file_id` | [int64](#int64) | optional |  |
| `profile_picture` | [string](#string) | optional |  |





### resources.sync.data.DataUserLocations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [CitizenLocations](#resourcessyncdataCitizenLocations) | repeated |  |
| `clear_all` | [bool](#bool) | optional |  |





### resources.sync.data.DataUsers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [DataUser](#resourcessyncdataDataUser) | repeated |  |





### resources.sync.data.DataVehicles


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vehicles` | [resources.vehicles.Vehicle](#resourcesvehiclesVehicle) | repeated |  |





### resources.sync.data.DeleteUsers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_ids` | [int32](#int32) | repeated |  |





### resources.sync.data.DeleteVehicles


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plates` | [string](#string) | repeated |  |





### resources.sync.data.LastCharID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `last_char_id` | [int32](#int32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/tracker/mapping.proto


### resources.tracker.UserMapping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `unit_id` | [int64](#int64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `hidden` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/activity/activity.proto


### resources.users.activity.CitizenDocumentRelation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [bool](#bool) |  |  |
| `document_id` | [int64](#int64) |  |  |
| `relation` | [resources.documents.relations.DocRelation](#resourcesdocumentsrelationsDocRelation) |  |  |





### resources.users.activity.FineChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `removed` | [bool](#bool) |  |  |
| `amount` | [int64](#int64) |  |  |





### resources.users.activity.JailChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seconds` | [int32](#int32) |  |  |
| `admin` | [bool](#bool) |  |  |
| `location` | [string](#string) | optional |  |





### resources.users.activity.JobChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) | optional |  |
| `grade_label` | [string](#string) | optional |  |





### resources.users.activity.LabelsChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [resources.users.labels.Label](#resourcesuserslabelsLabel) | repeated |  |
| `removed` | [resources.users.labels.Label](#resourcesuserslabelsLabel) | repeated |  |





### resources.users.activity.LicenseChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [bool](#bool) |  | `true` if licenses were added, `false` if removed |
| `licenses` | [resources.users.licenses.License](#resourcesuserslicensesLicense) | repeated |  |





### resources.users.activity.MugshotChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new` | [string](#string) | optional |  |





### resources.users.activity.NameChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |





### resources.users.activity.TrafficInfractionPointsChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old` | [uint32](#uint32) |  |  |
| `new` | [uint32](#uint32) |  |  |





### resources.users.activity.UserActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `type` | [UserActivityType](#resourcesusersactivityUserActivityType) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `source_user_id` | [int32](#int32) | optional |  |
| `source_user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) |  |  |
| `key` | [string](#string) |  |  |
| `reason` | [string](#string) |  |  |
| `data` | [UserActivityData](#resourcesusersactivityUserActivityData) | optional |  |
| `old_value` | [string](#string) |  |  |
| `new_value` | [string](#string) |  |  |





### resources.users.activity.UserActivityData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name_change` | [NameChange](#resourcesusersactivityNameChange) |  |  |
| `licenses_change` | [LicenseChange](#resourcesusersactivityLicenseChange) |  |  |
| `wanted_change` | [WantedChange](#resourcesusersactivityWantedChange) |  | User Props |
| `traffic_infraction_points_change` | [TrafficInfractionPointsChange](#resourcesusersactivityTrafficInfractionPointsChange) |  |  |
| `mugshot_change` | [MugshotChange](#resourcesusersactivityMugshotChange) |  |  |
| `labels_change` | [LabelsChange](#resourcesusersactivityLabelsChange) |  |  |
| `job_change` | [JobChange](#resourcesusersactivityJobChange) |  |  |
| `document_relation` | [CitizenDocumentRelation](#resourcesusersactivityCitizenDocumentRelation) |  | Docstore related |
| `jail_change` | [JailChange](#resourcesusersactivityJailChange) |  | "Plugin" activities |
| `fine_change` | [FineChange](#resourcesusersactivityFineChange) |  |  |





### resources.users.activity.WantedChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `wanted` | [bool](#bool) |  |  |




 <!-- end messages -->


### resources.users.activity.UserActivityType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `USER_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `USER_ACTIVITY_TYPE_NAME` | 4 | Name change |
| `USER_ACTIVITY_TYPE_LICENSES` | 5 | License(s) change added/removed |
| `USER_ACTIVITY_TYPE_WANTED` | 6 | Wanted status change |
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



## resources/vehicles/activity/activity.proto


### resources.vehicles.activity.VehicleActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `plate` | [string](#string) |  |  |
| `activity_type` | [VehicleActivityType](#resourcesvehiclesactivityVehicleActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [VehicleActivityData](#resourcesvehiclesactivityVehicleActivityData) |  |  |





### resources.vehicles.activity.VehicleActivityData




 <!-- end messages -->


### resources.vehicles.activity.VehicleActivityType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `VEHICLE_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `VEHICLE_ACTIVITY_TYPE_WANTED` | 1 | Types for `VehicleActivityData` |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/wiki/access/access.proto


### resources.wiki.access.PageAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageJobAccess](#resourceswikiaccessPageJobAccess) | repeated |  |
| `users` | [PageUserAccess](#resourceswikiaccessPageUserAccess) | repeated |  |





### resources.wiki.access.PageJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourceswikiaccessAccessLevel) |  |  |





### resources.wiki.access.PageUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `access` | [AccessLevel](#resourceswikiaccessAccessLevel) |  |  |




 <!-- end messages -->


### resources.wiki.access.AccessLevel

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



## resources/wiki/activity/activity.proto


### resources.wiki.activity.PageAccessJobsDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [resources.wiki.access.PageJobAccess](#resourceswikiaccessPageJobAccess) | repeated |  |
| `to_update` | [resources.wiki.access.PageJobAccess](#resourceswikiaccessPageJobAccess) | repeated |  |
| `to_delete` | [resources.wiki.access.PageJobAccess](#resourceswikiaccessPageJobAccess) | repeated |  |





### resources.wiki.activity.PageAccessUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageAccessJobsDiff](#resourceswikiactivityPageAccessJobsDiff) |  |  |
| `users` | [PageAccessUsersDiff](#resourceswikiactivityPageAccessUsersDiff) |  |  |





### resources.wiki.activity.PageAccessUsersDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [resources.wiki.access.PageUserAccess](#resourceswikiaccessPageUserAccess) | repeated |  |
| `to_update` | [resources.wiki.access.PageUserAccess](#resourceswikiaccessPageUserAccess) | repeated |  |
| `to_delete` | [resources.wiki.access.PageUserAccess](#resourceswikiaccessPageUserAccess) | repeated |  |





### resources.wiki.activity.PageActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `page_id` | [int64](#int64) |  |  |
| `activity_type` | [PageActivityType](#resourceswikiactivityPageActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [PageActivityData](#resourceswikiactivityPageActivityData) |  |  |





### resources.wiki.activity.PageActivityData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [PageUpdated](#resourceswikiactivityPageUpdated) |  |  |
| `access_updated` | [PageAccessUpdated](#resourceswikiactivityPageAccessUpdated) |  |  |





### resources.wiki.activity.PageFilesChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [int64](#int64) |  |  |
| `deleted` | [int64](#int64) |  |  |





### resources.wiki.activity.PageUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title_diff` | [string](#string) | optional |  |
| `title_cdiff` | [resources.common.content.ContentDiff](#resourcescommoncontentContentDiff) | optional |  |
| `description_diff` | [string](#string) | optional |  |
| `description_cdiff` | [resources.common.content.ContentDiff](#resourcescommoncontentContentDiff) | optional |  |
| `content_diff` | [string](#string) | optional |  |
| `content_cdiff` | [resources.common.content.ContentDiff](#resourcescommoncontentContentDiff) | optional |  |
| `files_change` | [PageFilesChange](#resourceswikiactivityPageFilesChange) | optional |  |




 <!-- end messages -->


### resources.wiki.activity.PageActivityType

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



## resources/wiki/page.proto


### resources.wiki.Page


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `parent_id` | [int64](#int64) | optional |  |
| `meta` | [PageMeta](#resourceswikiPageMeta) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `access` | [access.PageAccess](#resourceswikiaccessPageAccess) |  |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### resources.wiki.PageMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `slug` | [string](#string) | optional |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional |  |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `tags` | [string](#string) | repeated |  |
| `toc` | [bool](#bool) | optional |  |
| `public` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `startpage` | [bool](#bool) |  |  |





### resources.wiki.PageRootInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `logo_file_id` | [int64](#int64) | optional |  |
| `logo` | [resources.file.File](#resourcesfileFile) | optional |  |





### resources.wiki.PageShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `parent_id` | [int64](#int64) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `slug` | [string](#string) | optional |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `children` | [PageShort](#resourceswikiPageShort) | repeated |  |
| `root_info` | [PageRootInfo](#resourceswikiPageRootInfo) | optional |  |
| `level` | [int32](#int32) | optional |  |
| `draft` | [bool](#bool) |  |  |
| `startpage` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## services/auth/auth.proto


### services.auth.ChangePasswordRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |





### services.auth.ChangePasswordResponse





### services.auth.ChangeUsernameRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |





### services.auth.ChangeUsernameResponse





### services.auth.ChooseCharacterRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `char_id` | [int32](#int32) |  |  |





### services.auth.ChooseCharacterResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `job_props` | [resources.jobs.props.JobProps](#resourcesjobspropsJobProps) |  |  |
| `char` | [resources.users.User](#resourcesusersUser) |  |  |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |





### services.auth.CreateAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `password` | [string](#string) |  |  |





### services.auth.CreateAccountResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [int64](#int64) |  |  |





### services.auth.DeleteSocialLoginRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider` | [string](#string) |  |  |





### services.auth.DeleteSocialLoginResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `success` | [bool](#bool) |  |  |





### services.auth.ForgotPasswordRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |





### services.auth.ForgotPasswordResponse





### services.auth.GetAccountInfoRequest





### services.auth.GetAccountInfoResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [resources.accounts.Account](#resourcesaccountsAccount) |  |  |
| `oauth2_providers` | [resources.accounts.oauth2.OAuth2Provider](#resourcesaccountsoauth2OAuth2Provider) | repeated |  |
| `oauth2_connections` | [resources.accounts.oauth2.OAuth2Account](#resourcesaccountsoauth2OAuth2Account) | repeated |  |





### services.auth.GetCharactersRequest





### services.auth.GetCharactersResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chars` | [resources.accounts.Character](#resourcesaccountsCharacter) | repeated |  |





### services.auth.ImpersonateJobRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_grade` | [int32](#int32) |  |  |





### services.auth.ImpersonateJobResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `char` | [resources.users.User](#resourcesusersUser) |  |  |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |
| `state` | [bool](#bool) |  |  |





### services.auth.LoginRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `username` | [string](#string) |  |  |
| `password` | [string](#string) |  |  |





### services.auth.LoginResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `account_id` | [int64](#int64) |  |  |
| `char` | [ChooseCharacterResponse](#servicesauthChooseCharacterResponse) | optional |  |





### services.auth.LogoutRequest





### services.auth.LogoutResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `success` | [bool](#bool) |  |  |





### services.auth.SetSuperuserModeRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `superuser` | [bool](#bool) |  |  |
| `job` | [string](#string) | optional |  |





### services.auth.SetSuperuserModeResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `job_props` | [resources.jobs.props.JobProps](#resourcesjobspropsJobProps) | optional |  |
| `char` | [resources.users.User](#resourcesusersUser) |  |  |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.auth.AuthService
Auth Service handles user authentication, character selection and oauth2 connections Some methods **must** be caled via HTTP-based GRPC web request to allow cookies to be set/unset.


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Login` | [LoginRequest](#servicesauthLoginRequest) | [LoginResponse](#servicesauthLoginResponse) | |
| `Logout` | [LogoutRequest](#servicesauthLogoutRequest) | [LogoutResponse](#servicesauthLogoutResponse) | |
| `CreateAccount` | [CreateAccountRequest](#servicesauthCreateAccountRequest) | [CreateAccountResponse](#servicesauthCreateAccountResponse) | |
| `ChangeUsername` | [ChangeUsernameRequest](#servicesauthChangeUsernameRequest) | [ChangeUsernameResponse](#servicesauthChangeUsernameResponse) | |
| `ChangePassword` | [ChangePasswordRequest](#servicesauthChangePasswordRequest) | [ChangePasswordResponse](#servicesauthChangePasswordResponse) | |
| `ForgotPassword` | [ForgotPasswordRequest](#servicesauthForgotPasswordRequest) | [ForgotPasswordResponse](#servicesauthForgotPasswordResponse) | |
| `GetCharacters` | [GetCharactersRequest](#servicesauthGetCharactersRequest) | [GetCharactersResponse](#servicesauthGetCharactersResponse) | |
| `ChooseCharacter` | [ChooseCharacterRequest](#servicesauthChooseCharacterRequest) | [ChooseCharacterResponse](#servicesauthChooseCharacterResponse) | |
| `ImpersonateJob` | [ImpersonateJobRequest](#servicesauthImpersonateJobRequest) | [ImpersonateJobResponse](#servicesauthImpersonateJobResponse) | |
| `GetAccountInfo` | [GetAccountInfoRequest](#servicesauthGetAccountInfoRequest) | [GetAccountInfoResponse](#servicesauthGetAccountInfoResponse) | |
| `DeleteSocialLogin` | [DeleteSocialLoginRequest](#servicesauthDeleteSocialLoginRequest) | [DeleteSocialLoginResponse](#servicesauthDeleteSocialLoginResponse) | |
| `SetSuperuserMode` | [SetSuperuserModeRequest](#servicesauthSetSuperuserModeRequest) | [SetSuperuserModeResponse](#servicesauthSetSuperuserModeResponse) | |

 <!-- end services -->



## services/calendar/calendar.proto


### services.calendar.CreateCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resourcescalendarCalendar) |  |  |





### services.calendar.CreateCalendarResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resourcescalendarCalendar) |  |  |





### services.calendar.CreateOrUpdateCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.entries.CalendarEntry](#resourcescalendarentriesCalendarEntry) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |





### services.calendar.CreateOrUpdateCalendarEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.entries.CalendarEntry](#resourcescalendarentriesCalendarEntry) |  |  |





### services.calendar.DeleteCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [int64](#int64) |  |  |





### services.calendar.DeleteCalendarEntryResponse





### services.calendar.DeleteCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [int64](#int64) |  |  |





### services.calendar.DeleteCalendarResponse





### services.calendar.GetCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [int64](#int64) |  |  |





### services.calendar.GetCalendarEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.entries.CalendarEntry](#resourcescalendarentriesCalendarEntry) |  |  |





### services.calendar.GetCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [int64](#int64) |  |  |





### services.calendar.GetCalendarResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resourcescalendarCalendar) |  |  |





### services.calendar.GetUpcomingEntriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seconds` | [int32](#int32) |  |  |





### services.calendar.GetUpcomingEntriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [resources.calendar.entries.CalendarEntry](#resourcescalendarentriesCalendarEntry) | repeated |  |





### services.calendar.ListCalendarEntriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `year` | [int32](#int32) |  |  |
| `month` | [int32](#int32) |  |  |
| `calendar_ids` | [int64](#int64) | repeated |  |
| `show_hidden` | [bool](#bool) | optional |  |
| `after` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### services.calendar.ListCalendarEntriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [resources.calendar.entries.CalendarEntry](#resourcescalendarentriesCalendarEntry) | repeated |  |





### services.calendar.ListCalendarEntryRSVPRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `entry_id` | [int64](#int64) |  |  |





### services.calendar.ListCalendarEntryRSVPResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `entries` | [resources.calendar.entries.CalendarEntryRSVP](#resourcescalendarentriesCalendarEntryRSVP) | repeated |  |





### services.calendar.ListCalendarsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `only_public` | [bool](#bool) |  | Search params |
| `min_access_level` | [resources.calendar.access.AccessLevel](#resourcescalendaraccessAccessLevel) | optional |  |
| `after` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `calendar_ids` | [int64](#int64) | repeated |  |





### services.calendar.ListCalendarsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `calendars` | [resources.calendar.Calendar](#resourcescalendarCalendar) | repeated |  |





### services.calendar.ListSubscriptionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |





### services.calendar.ListSubscriptionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `subs` | [resources.calendar.CalendarSub](#resourcescalendarCalendarSub) | repeated |  |





### services.calendar.RSVPCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.entries.CalendarEntryRSVP](#resourcescalendarentriesCalendarEntryRSVP) |  |  |
| `subscribe` | [bool](#bool) |  |  |
| `remove` | [bool](#bool) | optional |  |





### services.calendar.RSVPCalendarEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.entries.CalendarEntryRSVP](#resourcescalendarentriesCalendarEntryRSVP) | optional |  |





### services.calendar.ShareCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [int64](#int64) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |





### services.calendar.ShareCalendarEntryResponse





### services.calendar.SubscribeToCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sub` | [resources.calendar.CalendarSub](#resourcescalendarCalendarSub) |  |  |
| `delete` | [bool](#bool) |  |  |





### services.calendar.SubscribeToCalendarResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sub` | [resources.calendar.CalendarSub](#resourcescalendarCalendarSub) |  |  |





### services.calendar.UpdateCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resourcescalendarCalendar) |  |  |





### services.calendar.UpdateCalendarResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar` | [resources.calendar.Calendar](#resourcescalendarCalendar) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.calendar.CalendarService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCalendars` | [ListCalendarsRequest](#servicescalendarListCalendarsRequest) | [ListCalendarsResponse](#servicescalendarListCalendarsResponse) | |
| `GetCalendar` | [GetCalendarRequest](#servicescalendarGetCalendarRequest) | [GetCalendarResponse](#servicescalendarGetCalendarResponse) | |
| `CreateCalendar` | [CreateCalendarRequest](#servicescalendarCreateCalendarRequest) | [CreateCalendarResponse](#servicescalendarCreateCalendarResponse) | |
| `UpdateCalendar` | [UpdateCalendarRequest](#servicescalendarUpdateCalendarRequest) | [UpdateCalendarResponse](#servicescalendarUpdateCalendarResponse) | |
| `DeleteCalendar` | [DeleteCalendarRequest](#servicescalendarDeleteCalendarRequest) | [DeleteCalendarResponse](#servicescalendarDeleteCalendarResponse) | |
| `ListCalendarEntries` | [ListCalendarEntriesRequest](#servicescalendarListCalendarEntriesRequest) | [ListCalendarEntriesResponse](#servicescalendarListCalendarEntriesResponse) | |
| `GetUpcomingEntries` | [GetUpcomingEntriesRequest](#servicescalendarGetUpcomingEntriesRequest) | [GetUpcomingEntriesResponse](#servicescalendarGetUpcomingEntriesResponse) | |
| `GetCalendarEntry` | [GetCalendarEntryRequest](#servicescalendarGetCalendarEntryRequest) | [GetCalendarEntryResponse](#servicescalendarGetCalendarEntryResponse) | |
| `CreateOrUpdateCalendarEntry` | [CreateOrUpdateCalendarEntryRequest](#servicescalendarCreateOrUpdateCalendarEntryRequest) | [CreateOrUpdateCalendarEntryResponse](#servicescalendarCreateOrUpdateCalendarEntryResponse) | |
| `DeleteCalendarEntry` | [DeleteCalendarEntryRequest](#servicescalendarDeleteCalendarEntryRequest) | [DeleteCalendarEntryResponse](#servicescalendarDeleteCalendarEntryResponse) | |
| `ShareCalendarEntry` | [ShareCalendarEntryRequest](#servicescalendarShareCalendarEntryRequest) | [ShareCalendarEntryResponse](#servicescalendarShareCalendarEntryResponse) | |
| `ListCalendarEntryRSVP` | [ListCalendarEntryRSVPRequest](#servicescalendarListCalendarEntryRSVPRequest) | [ListCalendarEntryRSVPResponse](#servicescalendarListCalendarEntryRSVPResponse) | |
| `RSVPCalendarEntry` | [RSVPCalendarEntryRequest](#servicescalendarRSVPCalendarEntryRequest) | [RSVPCalendarEntryResponse](#servicescalendarRSVPCalendarEntryResponse) | |
| `ListSubscriptions` | [ListSubscriptionsRequest](#servicescalendarListSubscriptionsRequest) | [ListSubscriptionsResponse](#servicescalendarListSubscriptionsResponse) | |
| `SubscribeToCalendar` | [SubscribeToCalendarRequest](#servicescalendarSubscribeToCalendarRequest) | [SubscribeToCalendarResponse](#servicescalendarSubscribeToCalendarResponse) | |

 <!-- end services -->



## services/centrum/centrum.proto


### services.centrum.AssignDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [int64](#int64) |  |  |
| `to_add` | [int64](#int64) | repeated |  |
| `to_remove` | [int64](#int64) | repeated |  |
| `forced` | [bool](#bool) | optional |  |





### services.centrum.AssignDispatchResponse





### services.centrum.AssignUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [int64](#int64) |  |  |
| `to_add` | [int32](#int32) | repeated |  |
| `to_remove` | [int32](#int32) | repeated |  |





### services.centrum.AssignUnitResponse





### services.centrum.CreateDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |





### services.centrum.CreateDispatchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |





### services.centrum.CreateOrUpdateUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) |  |  |





### services.centrum.CreateOrUpdateUnitResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) |  |  |





### services.centrum.DeleteDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.centrum.DeleteDispatchResponse





### services.centrum.DeleteUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [int64](#int64) |  |  |





### services.centrum.DeleteUnitResponse





### services.centrum.GetDispatchHeatmapRequest





### services.centrum.GetDispatchHeatmapResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_entries` | [int32](#int32) |  |  |
| `entries` | [resources.livemap.heatmap.HeatmapEntry](#resourceslivemapheatmapHeatmapEntry) | repeated |  |





### services.centrum.GetDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.centrum.GetDispatchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |





### services.centrum.GetSettingsRequest





### services.centrum.GetSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.settings.Settings](#resourcescentrumsettingsSettings) |  |  |
| `effective_access` | [resources.centrum.settings.EffectiveAccess](#resourcescentrumsettingsEffectiveAccess) |  |  |





### services.centrum.JoinUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [int64](#int64) | optional |  |





### services.centrum.JoinUnitResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) |  |  |





### services.centrum.LatestState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [resources.centrum.dispatchers.JobDispatchers](#resourcescentrumdispatchersJobDispatchers) |  |  |
| `own_unit_id` | [int64](#int64) | optional |  |
| `units` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) | repeated | Send the current units and dispatches |
| `dispatches` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) | repeated |  |





### services.centrum.ListDispatchActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `id` | [int64](#int64) |  |  |





### services.centrum.ListDispatchActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.centrum.dispatches.DispatchStatus](#resourcescentrumdispatchesDispatchStatus) | repeated |  |





### services.centrum.ListDispatchTargetJobsRequest





### services.centrum.ListDispatchTargetJobsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |





### services.centrum.ListDispatchesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `status` | [resources.centrum.dispatches.StatusDispatch](#resourcescentrumdispatchesStatusDispatch) | repeated |  |
| `not_status` | [resources.centrum.dispatches.StatusDispatch](#resourcescentrumdispatchesStatusDispatch) | repeated |  |
| `ids` | [int64](#int64) | repeated |  |
| `postal` | [string](#string) | optional |  |





### services.centrum.ListDispatchesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `dispatches` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) | repeated |  |





### services.centrum.ListUnitActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `id` | [int64](#int64) |  |  |





### services.centrum.ListUnitActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.centrum.units.UnitStatus](#resourcescentrumunitsUnitStatus) | repeated |  |





### services.centrum.ListUnitsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [resources.centrum.units.StatusUnit](#resourcescentrumunitsStatusUnit) | repeated |  |





### services.centrum.ListUnitsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `units` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) | repeated |  |





### services.centrum.StreamHandshake


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `server_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `settings` | [resources.centrum.settings.Settings](#resourcescentrumsettingsSettings) |  |  |
| `access` | [resources.centrum.settings.EffectiveAccess](#resourcescentrumsettingsEffectiveAccess) |  |  |





### services.centrum.StreamRequest





### services.centrum.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `handshake` | [StreamHandshake](#servicescentrumStreamHandshake) |  |  |
| `latest_state` | [LatestState](#servicescentrumLatestState) |  |  |
| `settings` | [resources.centrum.settings.Settings](#resourcescentrumsettingsSettings) |  |  |
| `access` | [resources.centrum.settings.EffectiveAccess](#resourcescentrumsettingsEffectiveAccess) |  |  |
| `dispatchers` | [resources.centrum.dispatchers.Dispatchers](#resourcescentrumdispatchersDispatchers) |  |  |
| `unit_deleted` | [int64](#int64) |  |  |
| `unit_updated` | [resources.centrum.units.Unit](#resourcescentrumunitsUnit) |  |  |
| `unit_status` | [resources.centrum.units.UnitStatus](#resourcescentrumunitsUnitStatus) |  |  |
| `dispatch_deleted` | [int64](#int64) |  |  |
| `dispatch_updated` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |
| `dispatch_status` | [resources.centrum.dispatches.DispatchStatus](#resourcescentrumdispatchesDispatchStatus) |  |  |





### services.centrum.TakeControlRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signon` | [bool](#bool) |  |  |





### services.centrum.TakeControlResponse





### services.centrum.TakeDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_ids` | [int64](#int64) | repeated |  |
| `resp` | [resources.centrum.dispatches.TakeDispatchResp](#resourcescentrumdispatchesTakeDispatchResp) |  |  |
| `reason` | [string](#string) | optional |  |





### services.centrum.TakeDispatchResponse





### services.centrum.UpdateDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |





### services.centrum.UpdateDispatchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |





### services.centrum.UpdateDispatchStatusRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [int64](#int64) |  |  |
| `status` | [resources.centrum.dispatches.StatusDispatch](#resourcescentrumdispatchesStatusDispatch) |  |  |
| `reason` | [string](#string) | optional |  |
| `code` | [string](#string) | optional |  |





### services.centrum.UpdateDispatchStatusResponse





### services.centrum.UpdateDispatchersRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_remove` | [int32](#int32) | repeated |  |





### services.centrum.UpdateDispatchersResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [resources.centrum.dispatchers.Dispatchers](#resourcescentrumdispatchersDispatchers) |  |  |





### services.centrum.UpdateSettingsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.settings.Settings](#resourcescentrumsettingsSettings) |  |  |





### services.centrum.UpdateSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.settings.Settings](#resourcescentrumsettingsSettings) |  |  |





### services.centrum.UpdateUnitStatusRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [int64](#int64) |  |  |
| `status` | [resources.centrum.units.StatusUnit](#resourcescentrumunitsStatusUnit) |  |  |
| `reason` | [string](#string) | optional |  |
| `code` | [string](#string) | optional |  |





### services.centrum.UpdateUnitStatusResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.centrum.CentrumService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `UpdateSettings` | [UpdateSettingsRequest](#servicescentrumUpdateSettingsRequest) | [UpdateSettingsResponse](#servicescentrumUpdateSettingsResponse) | |
| `CreateDispatch` | [CreateDispatchRequest](#servicescentrumCreateDispatchRequest) | [CreateDispatchResponse](#servicescentrumCreateDispatchResponse) | |
| `UpdateDispatch` | [UpdateDispatchRequest](#servicescentrumUpdateDispatchRequest) | [UpdateDispatchResponse](#servicescentrumUpdateDispatchResponse) | |
| `DeleteDispatch` | [DeleteDispatchRequest](#servicescentrumDeleteDispatchRequest) | [DeleteDispatchResponse](#servicescentrumDeleteDispatchResponse) | |
| `ListDispatchTargetJobs` | [ListDispatchTargetJobsRequest](#servicescentrumListDispatchTargetJobsRequest) | [ListDispatchTargetJobsResponse](#servicescentrumListDispatchTargetJobsResponse) | |
| `TakeControl` | [TakeControlRequest](#servicescentrumTakeControlRequest) | [TakeControlResponse](#servicescentrumTakeControlResponse) | |
| `AssignDispatch` | [AssignDispatchRequest](#servicescentrumAssignDispatchRequest) | [AssignDispatchResponse](#servicescentrumAssignDispatchResponse) | |
| `AssignUnit` | [AssignUnitRequest](#servicescentrumAssignUnitRequest) | [AssignUnitResponse](#servicescentrumAssignUnitResponse) | |
| `GetDispatchHeatmap` | [GetDispatchHeatmapRequest](#servicescentrumGetDispatchHeatmapRequest) | [GetDispatchHeatmapResponse](#servicescentrumGetDispatchHeatmapResponse) | |
| `UpdateDispatchers` | [UpdateDispatchersRequest](#servicescentrumUpdateDispatchersRequest) | [UpdateDispatchersResponse](#servicescentrumUpdateDispatchersResponse) | |
| `Stream` | [StreamRequest](#servicescentrumStreamRequest) | [StreamResponse](#servicescentrumStreamResponse) stream | |
| `GetSettings` | [GetSettingsRequest](#servicescentrumGetSettingsRequest) | [GetSettingsResponse](#servicescentrumGetSettingsResponse) | |
| `JoinUnit` | [JoinUnitRequest](#servicescentrumJoinUnitRequest) | [JoinUnitResponse](#servicescentrumJoinUnitResponse) | |
| `ListUnits` | [ListUnitsRequest](#servicescentrumListUnitsRequest) | [ListUnitsResponse](#servicescentrumListUnitsResponse) | |
| `ListUnitActivity` | [ListUnitActivityRequest](#servicescentrumListUnitActivityRequest) | [ListUnitActivityResponse](#servicescentrumListUnitActivityResponse) | |
| `GetDispatch` | [GetDispatchRequest](#servicescentrumGetDispatchRequest) | [GetDispatchResponse](#servicescentrumGetDispatchResponse) | |
| `ListDispatches` | [ListDispatchesRequest](#servicescentrumListDispatchesRequest) | [ListDispatchesResponse](#servicescentrumListDispatchesResponse) | |
| `ListDispatchActivity` | [ListDispatchActivityRequest](#servicescentrumListDispatchActivityRequest) | [ListDispatchActivityResponse](#servicescentrumListDispatchActivityResponse) | |
| `CreateOrUpdateUnit` | [CreateOrUpdateUnitRequest](#servicescentrumCreateOrUpdateUnitRequest) | [CreateOrUpdateUnitResponse](#servicescentrumCreateOrUpdateUnitResponse) | |
| `DeleteUnit` | [DeleteUnitRequest](#servicescentrumDeleteUnitRequest) | [DeleteUnitResponse](#servicescentrumDeleteUnitResponse) | |
| `TakeDispatch` | [TakeDispatchRequest](#servicescentrumTakeDispatchRequest) | [TakeDispatchResponse](#servicescentrumTakeDispatchResponse) | |
| `UpdateUnitStatus` | [UpdateUnitStatusRequest](#servicescentrumUpdateUnitStatusRequest) | [UpdateUnitStatusResponse](#servicescentrumUpdateUnitStatusResponse) | |
| `UpdateDispatchStatus` | [UpdateDispatchStatusRequest](#servicescentrumUpdateDispatchStatusRequest) | [UpdateDispatchStatusResponse](#servicescentrumUpdateDispatchStatusResponse) | |

 <!-- end services -->



## services/citizens/citizens.proto


### services.citizens.DeleteAvatarRequest





### services.citizens.DeleteAvatarResponse





### services.citizens.DeleteMugshotRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `reason` | [string](#string) |  |  |





### services.citizens.DeleteMugshotResponse





### services.citizens.GetUserRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `info_only` | [bool](#bool) | optional |  |





### services.citizens.GetUserResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user` | [resources.users.User](#resourcesusersUser) |  |  |





### services.citizens.ListCitizensRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `search` | [string](#string) |  | Search params |
| `wanted` | [bool](#bool) | optional |  |
| `phone_number` | [string](#string) | optional |  |
| `traffic_infraction_points` | [uint32](#uint32) | optional |  |
| `dateofbirth` | [string](#string) | optional |  |
| `open_fines` | [int64](#int64) | optional |  |





### services.citizens.ListCitizensResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `users` | [resources.users.User](#resourcesusersUser) | repeated |  |





### services.citizens.ListUserActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_id` | [int32](#int32) |  | Search params |
| `types` | [resources.users.activity.UserActivityType](#resourcesusersactivityUserActivityType) | repeated |  |





### services.citizens.ListUserActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.users.activity.UserActivity](#resourcesusersactivityUserActivity) | repeated |  |





### services.citizens.ManageLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.labels.Label](#resourcesuserslabelsLabel) | repeated |  |





### services.citizens.ManageLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.labels.Label](#resourcesuserslabelsLabel) | repeated |  |





### services.citizens.SetUserPropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.users.props.UserProps](#resourcesuserspropsUserProps) |  |  |
| `reason` | [string](#string) |  |  |





### services.citizens.SetUserPropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.users.props.UserProps](#resourcesuserspropsUserProps) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.citizens.CitizensService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCitizens` | [ListCitizensRequest](#servicescitizensListCitizensRequest) | [ListCitizensResponse](#servicescitizensListCitizensResponse) | |
| `GetUser` | [GetUserRequest](#servicescitizensGetUserRequest) | [GetUserResponse](#servicescitizensGetUserResponse) | |
| `ListUserActivity` | [ListUserActivityRequest](#servicescitizensListUserActivityRequest) | [ListUserActivityResponse](#servicescitizensListUserActivityResponse) | |
| `SetUserProps` | [SetUserPropsRequest](#servicescitizensSetUserPropsRequest) | [SetUserPropsResponse](#servicescitizensSetUserPropsResponse) | |
| `UploadAvatar` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `DeleteAvatar` | [DeleteAvatarRequest](#servicescitizensDeleteAvatarRequest) | [DeleteAvatarResponse](#servicescitizensDeleteAvatarResponse) | |
| `UploadMugshot` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `DeleteMugshot` | [DeleteMugshotRequest](#servicescitizensDeleteMugshotRequest) | [DeleteMugshotResponse](#servicescitizensDeleteMugshotResponse) | |
| `ManageLabels` | [ManageLabelsRequest](#servicescitizensManageLabelsRequest) | [ManageLabelsResponse](#servicescitizensManageLabelsResponse) | |

 <!-- end services -->



## services/completor/completor.proto


### services.completor.CompleteCitizenLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |





### services.completor.CompleteCitizenLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.labels.Label](#resourcesuserslabelsLabel) | repeated |  |





### services.completor.CompleteCitizensRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |
| `current_job` | [bool](#bool) | optional |  |
| `on_duty` | [bool](#bool) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `user_ids_only` | [bool](#bool) | optional |  |





### services.completor.CompleteCitizensResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | repeated |  |





### services.completor.CompleteDocumentCategoriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |





### services.completor.CompleteDocumentCategoriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `categories` | [resources.documents.category.Category](#resourcesdocumentscategoryCategory) | repeated |  |





### services.completor.CompleteJobsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) | optional |  |
| `exact_match` | [bool](#bool) | optional |  |
| `current_job` | [bool](#bool) | optional |  |





### services.completor.CompleteJobsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |





### services.completor.ListLawBooksRequest





### services.completor.ListLawBooksResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `books` | [resources.laws.LawBook](#resourceslawsLawBook) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.completor.CompletorService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `CompleteCitizens` | [CompleteCitizensRequest](#servicescompletorCompleteCitizensRequest) | [CompleteCitizensResponse](#servicescompletorCompleteCitizensResponse) | |
| `CompleteJobs` | [CompleteJobsRequest](#servicescompletorCompleteJobsRequest) | [CompleteJobsResponse](#servicescompletorCompleteJobsResponse) | |
| `CompleteDocumentCategories` | [CompleteDocumentCategoriesRequest](#servicescompletorCompleteDocumentCategoriesRequest) | [CompleteDocumentCategoriesResponse](#servicescompletorCompleteDocumentCategoriesResponse) | |
| `ListLawBooks` | [ListLawBooksRequest](#servicescompletorListLawBooksRequest) | [ListLawBooksResponse](#servicescompletorListLawBooksResponse) | |
| `CompleteCitizenLabels` | [CompleteCitizenLabelsRequest](#servicescompletorCompleteCitizenLabelsRequest) | [CompleteCitizenLabelsResponse](#servicescompletorCompleteCitizenLabelsResponse) | |

 <!-- end services -->



## services/documents/approval.proto


### services.documents.ApprovalTaskSeed
A declarative "ensure" for tasks under one policy/snapshot. Exactly one target must be set: user_id OR (job + minimum_grade).



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  | If set -> USER task; slots is forced to 1 |
| `job` | [string](#string) |  | If user_id == 0 -> JOB task |
| `minimum_grade` | [int32](#int32) |  |  |
| `label` | [string](#string) | optional | Label of task |
| `signature_required` | [bool](#bool) |  |  |
| `slots` | [int32](#int32) |  | Only for JOB tasks; number of PENDING slots to ensure (>=1) |
| `due_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional | Optional default due date for created slots |
| `comment` | [string](#string) | optional | Optional note set on created tasks |





### services.documents.DecideApprovalRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `task_id` | [int64](#int64) | optional |  |
| `new_status` | [resources.documents.approval.ApprovalTaskStatus](#resourcesdocumentsapprovalApprovalTaskStatus) |  | APPROVED or DECLINED |
| `comment` | [string](#string) |  |  |
| `payload_svg` | [string](#string) | optional |  |
| `stamp_id` | [int64](#int64) | optional | When type=STAMP |





### services.documents.DecideApprovalResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approval` | [resources.documents.approval.Approval](#resourcesdocumentsapprovalApproval) |  |  |
| `task` | [resources.documents.approval.ApprovalTask](#resourcesdocumentsapprovalApprovalTask) |  |  |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  |  |





### services.documents.DeleteApprovalTasksRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `task_ids` | [int64](#int64) | repeated |  |
| `delete_all_pending` | [bool](#bool) |  | If true, ignore task_ids and delete all PENDING tasks under this policy |





### services.documents.DeleteApprovalTasksResponse





### services.documents.ListApprovalPoliciesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |





### services.documents.ListApprovalPoliciesResponse
Only one policy per document is supported currently.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  |  |
| `doc_meta` | [resources.documents.DocumentMeta](#resourcesdocumentsDocumentMeta) |  |  |





### services.documents.ListApprovalTasksInboxRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `statuses` | [resources.documents.approval.ApprovalTaskStatus](#resourcesdocumentsapprovalApprovalTaskStatus) | repeated |  |
| `only_drafts` | [bool](#bool) | optional | Controls inclusion of drafts in the result: - unset/null: include all documents (drafts and non-drafts) - false: only non-draft documents - true: only draft documents |





### services.documents.ListApprovalTasksInboxResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `tasks` | [resources.documents.approval.ApprovalTask](#resourcesdocumentsapprovalApprovalTask) | repeated |  |





### services.documents.ListApprovalTasksRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `statuses` | [resources.documents.approval.ApprovalTaskStatus](#resourcesdocumentsapprovalApprovalTaskStatus) | repeated |  |





### services.documents.ListApprovalTasksResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tasks` | [resources.documents.approval.ApprovalTask](#resourcesdocumentsapprovalApprovalTask) | repeated |  |





### services.documents.ListApprovalsRequest
List approvals (artifacts) for a policy/snapshot. If snapshot_date is unset, server defaults to policy.snapshot_date.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `task_id` | [int64](#int64) | optional |  |
| `snapshot_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `status` | [resources.documents.approval.ApprovalStatus](#resourcesdocumentsapprovalApprovalStatus) | optional | Optional filters |
| `user_id` | [int32](#int32) | optional | Filter by signer |





### services.documents.ListApprovalsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approvals` | [resources.documents.approval.Approval](#resourcesdocumentsapprovalApproval) | repeated |  |





### services.documents.RecomputeApprovalPolicyCountersRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |





### services.documents.RecomputeApprovalPolicyCountersResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  |  |





### services.documents.ReopenApprovalTaskRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `task_id` | [int64](#int64) |  |  |
| `comment` | [string](#string) |  |  |





### services.documents.ReopenApprovalTaskResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `task` | [resources.documents.approval.ApprovalTask](#resourcesdocumentsapprovalApprovalTask) |  |  |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  |  |





### services.documents.RevokeApprovalRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approval_id` | [int64](#int64) |  |  |
| `comment` | [string](#string) |  |  |





### services.documents.RevokeApprovalResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `approval` | [resources.documents.approval.Approval](#resourcesdocumentsapprovalApproval) |  |  |





### services.documents.UpsertApprovalPolicyRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  |  |





### services.documents.UpsertApprovalPolicyResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  |  |





### services.documents.UpsertApprovalTasksRequest
Upsert = insert missing PENDING tasks/slots; will NOT delete existing tasks. Identity rules (server-side): - USER task: unique by (document_id, snapshot_date, assignee_kind=USER, user_id) - JOB task: unique by (document_id, snapshot_date, assignee_kind=JOB, job, minimum_grade, slot_no) For JOB seeds with slots=N, the server ensures there are at least N PENDING slots (slot_no 1..N).



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `snapshot_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional | If empty, use policy.snapshot_date |
| `seeds` | [ApprovalTaskSeed](#servicesdocumentsApprovalTaskSeed) | repeated |  |





### services.documents.UpsertApprovalTasksResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tasks_created` | [int32](#int32) |  | Number of new task rows inserted |
| `tasks_ensured` | [int32](#int32) |  | Number of requested targets already satisfied (no-op) |
| `policy` | [resources.documents.approval.ApprovalPolicy](#resourcesdocumentsapprovalApprovalPolicy) |  | Echo (optional convenience) |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.ApprovalService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListApprovalTasksInbox` | [ListApprovalTasksInboxRequest](#servicesdocumentsListApprovalTasksInboxRequest) | [ListApprovalTasksInboxResponse](#servicesdocumentsListApprovalTasksInboxResponse) |Inbox (for tasks assigned to user) |
| `ListApprovalPolicies` | [ListApprovalPoliciesRequest](#servicesdocumentsListApprovalPoliciesRequest) | [ListApprovalPoliciesResponse](#servicesdocumentsListApprovalPoliciesResponse) |Policies |
| `UpsertApprovalPolicy` | [UpsertApprovalPolicyRequest](#servicesdocumentsUpsertApprovalPolicyRequest) | [UpsertApprovalPolicyResponse](#servicesdocumentsUpsertApprovalPolicyResponse) | |
| `ListApprovalTasks` | [ListApprovalTasksRequest](#servicesdocumentsListApprovalTasksRequest) | [ListApprovalTasksResponse](#servicesdocumentsListApprovalTasksResponse) |Tasks |
| `UpsertApprovalTasks` | [UpsertApprovalTasksRequest](#servicesdocumentsUpsertApprovalTasksRequest) | [UpsertApprovalTasksResponse](#servicesdocumentsUpsertApprovalTasksResponse) | |
| `DeleteApprovalTasks` | [DeleteApprovalTasksRequest](#servicesdocumentsDeleteApprovalTasksRequest) | [DeleteApprovalTasksResponse](#servicesdocumentsDeleteApprovalTasksResponse) | |
| `ListApprovals` | [ListApprovalsRequest](#servicesdocumentsListApprovalsRequest) | [ListApprovalsResponse](#servicesdocumentsListApprovalsResponse) |Approval |
| `RevokeApproval` | [RevokeApprovalRequest](#servicesdocumentsRevokeApprovalRequest) | [RevokeApprovalResponse](#servicesdocumentsRevokeApprovalResponse) | |
| `DecideApproval` | [DecideApprovalRequest](#servicesdocumentsDecideApprovalRequest) | [DecideApprovalResponse](#servicesdocumentsDecideApprovalResponse) | |
| `ReopenApprovalTask` | [ReopenApprovalTaskRequest](#servicesdocumentsReopenApprovalTaskRequest) | [ReopenApprovalTaskResponse](#servicesdocumentsReopenApprovalTaskResponse) | |
| `RecomputeApprovalPolicyCounters` | [RecomputeApprovalPolicyCountersRequest](#servicesdocumentsRecomputeApprovalPolicyCountersRequest) | [RecomputeApprovalPolicyCountersResponse](#servicesdocumentsRecomputeApprovalPolicyCountersResponse) |Helpers |

 <!-- end services -->



## services/documents/collab.proto

 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.CollabService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `JoinRoom` | [.resources.collab.ClientPacket](#resourcescollabClientPacket) stream | [.resources.collab.ServerPacket](#resourcescollabServerPacket) stream |buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |

 <!-- end services -->



## services/documents/documents.proto


### services.documents.AddDocumentReferenceRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reference` | [resources.documents.references.DocumentReference](#resourcesdocumentsreferencesDocumentReference) |  |  |





### services.documents.AddDocumentReferenceResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.AddDocumentRelationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relation` | [resources.documents.relations.DocumentRelation](#resourcesdocumentsrelationsDocumentRelation) |  |  |





### services.documents.AddDocumentRelationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.ChangeDocumentOwnerRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `new_user_id` | [int32](#int32) | optional |  |





### services.documents.ChangeDocumentOwnerResponse





### services.documents.CreateDocumentReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `request_type` | [resources.documents.activity.DocActivityType](#resourcesdocumentsactivityDocActivityType) |  |  |
| `reason` | [string](#string) | optional |  |
| `data` | [resources.documents.activity.DocActivityData](#resourcesdocumentsactivityDocActivityData) | optional |  |





### services.documents.CreateDocumentReqResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.documents.requests.DocRequest](#resourcesdocumentsrequestsDocRequest) |  |  |





### services.documents.CreateDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `template_id` | [int64](#int64) | optional |  |
| `template_data` | [resources.documents.templates.TemplateData](#resourcesdocumentstemplatesTemplateData) | optional |  |





### services.documents.CreateDocumentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.CreateOrUpdateCategoryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [resources.documents.category.Category](#resourcesdocumentscategoryCategory) |  |  |





### services.documents.CreateOrUpdateCategoryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [resources.documents.category.Category](#resourcesdocumentscategoryCategory) |  |  |





### services.documents.CreateTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.templates.Template](#resourcesdocumentstemplatesTemplate) |  |  |





### services.documents.CreateTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.DeleteCategoryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.DeleteCategoryResponse





### services.documents.DeleteCommentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment_id` | [int64](#int64) |  |  |





### services.documents.DeleteCommentResponse





### services.documents.DeleteDocumentReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request_id` | [int64](#int64) |  |  |





### services.documents.DeleteDocumentReqResponse





### services.documents.DeleteDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `reason` | [string](#string) | optional |  |





### services.documents.DeleteDocumentResponse





### services.documents.DeleteTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.DeleteTemplateResponse





### services.documents.EditCommentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.comment.Comment](#resourcesdocumentscommentComment) |  |  |





### services.documents.EditCommentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.comment.Comment](#resourcesdocumentscommentComment) |  |  |





### services.documents.GetCommentsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `document_id` | [int64](#int64) |  |  |





### services.documents.GetCommentsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `comments` | [resources.documents.comment.Comment](#resourcesdocumentscommentComment) | repeated |  |





### services.documents.GetDocumentAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |





### services.documents.GetDocumentAccessResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [resources.documents.access.DocumentAccess](#resourcesdocumentsaccessDocumentAccess) |  |  |





### services.documents.GetDocumentReferencesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |





### services.documents.GetDocumentReferencesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `references` | [resources.documents.references.DocumentReference](#resourcesdocumentsreferencesDocumentReference) | repeated |  |





### services.documents.GetDocumentRelationsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |





### services.documents.GetDocumentRelationsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relations` | [resources.documents.relations.DocumentRelation](#resourcesdocumentsrelationsDocumentRelation) | repeated |  |





### services.documents.GetDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `info_only` | [bool](#bool) | optional |  |





### services.documents.GetDocumentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [resources.documents.Document](#resourcesdocumentsDocument) |  |  |
| `access` | [resources.documents.access.DocumentAccess](#resourcesdocumentsaccessDocumentAccess) |  |  |





### services.documents.GetTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template_id` | [int64](#int64) |  |  |
| `data` | [resources.documents.templates.TemplateData](#resourcesdocumentstemplatesTemplateData) | optional |  |
| `render` | [bool](#bool) | optional |  |





### services.documents.GetTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.templates.Template](#resourcesdocumentstemplatesTemplate) |  |  |
| `rendered` | [bool](#bool) |  |  |





### services.documents.ListCategoriesRequest





### services.documents.ListCategoriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `categories` | [resources.documents.category.Category](#resourcesdocumentscategoryCategory) | repeated |  |





### services.documents.ListDocumentActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `document_id` | [int64](#int64) |  |  |
| `activity_types` | [resources.documents.activity.DocActivityType](#resourcesdocumentsactivityDocActivityType) | repeated | Search params |





### services.documents.ListDocumentActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.documents.activity.DocActivity](#resourcesdocumentsactivityDocActivity) | repeated |  |





### services.documents.ListDocumentPinsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `personal` | [bool](#bool) | optional | Search params If true, only personal pins are returned |





### services.documents.ListDocumentPinsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `documents` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | repeated |  |





### services.documents.ListDocumentReqsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `document_id` | [int64](#int64) |  |  |





### services.documents.ListDocumentReqsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `requests` | [resources.documents.requests.DocRequest](#resourcesdocumentsrequestsDocRequest) | repeated |  |





### services.documents.ListDocumentsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `search` | [string](#string) | optional | Search params |
| `category_ids` | [int64](#int64) | repeated |  |
| `creator_ids` | [int32](#int32) | repeated |  |
| `from` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `to` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `closed` | [bool](#bool) | optional |  |
| `document_ids` | [int64](#int64) | repeated |  |
| `only_drafts` | [bool](#bool) | optional | Controls inclusion of drafts in the result: - unset/null: include all documents (drafts and non-drafts) - false: only non-draft documents - true: only draft documents |





### services.documents.ListDocumentsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `documents` | [resources.documents.DocumentShort](#resourcesdocumentsDocumentShort) | repeated |  |





### services.documents.ListTemplatesRequest





### services.documents.ListTemplatesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `templates` | [resources.documents.templates.TemplateShort](#resourcesdocumentstemplatesTemplateShort) | repeated |  |





### services.documents.ListUserDocumentsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `relations` | [resources.documents.relations.DocRelation](#resourcesdocumentsrelationsDocRelation) | repeated |  |
| `closed` | [bool](#bool) | optional |  |





### services.documents.ListUserDocumentsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `relations` | [resources.documents.relations.DocumentRelation](#resourcesdocumentsrelationsDocumentRelation) | repeated |  |





### services.documents.PostCommentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.comment.Comment](#resourcesdocumentscommentComment) |  |  |





### services.documents.PostCommentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.comment.Comment](#resourcesdocumentscommentComment) |  |  |





### services.documents.RemoveDocumentReferenceRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.RemoveDocumentReferenceResponse





### services.documents.RemoveDocumentRelationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.documents.RemoveDocumentRelationResponse





### services.documents.SetDocumentAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `access` | [resources.documents.access.DocumentAccess](#resourcesdocumentsaccessDocumentAccess) |  |  |





### services.documents.SetDocumentAccessResponse





### services.documents.SetDocumentReminderRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `reminder_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `message` | [string](#string) | optional |  |
| `max_reminder_count` | [int32](#int32) |  |  |





### services.documents.SetDocumentReminderResponse





### services.documents.ToggleDocumentPinRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `state` | [bool](#bool) |  |  |
| `personal` | [bool](#bool) | optional | If true, the pin is personal and not shared with other job members |





### services.documents.ToggleDocumentPinResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pin` | [resources.documents.pins.DocumentPin](#resourcesdocumentspinsDocumentPin) | optional |  |





### services.documents.ToggleDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `closed` | [bool](#bool) |  |  |





### services.documents.ToggleDocumentResponse





### services.documents.UpdateDocumentReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `request_id` | [int64](#int64) |  |  |
| `reason` | [string](#string) | optional |  |
| `data` | [resources.documents.activity.DocActivityData](#resourcesdocumentsactivityDocActivityData) | optional |  |
| `accepted` | [bool](#bool) |  |  |





### services.documents.UpdateDocumentReqResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.documents.requests.DocRequest](#resourcesdocumentsrequestsDocRequest) |  |  |





### services.documents.UpdateDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [int64](#int64) |  |  |
| `category_id` | [int64](#int64) | optional |  |
| `title` | [string](#string) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `data` | [string](#string) | optional |  |
| `meta` | [resources.documents.DocumentMeta](#resourcesdocumentsDocumentMeta) |  |  |
| `access` | [resources.documents.access.DocumentAccess](#resourcesdocumentsaccessDocumentAccess) | optional |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### services.documents.UpdateDocumentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [resources.documents.Document](#resourcesdocumentsDocument) |  |  |





### services.documents.UpdateTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.templates.Template](#resourcesdocumentstemplatesTemplate) |  |  |





### services.documents.UpdateTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.templates.Template](#resourcesdocumentstemplatesTemplate) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.DocumentsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTemplates` | [ListTemplatesRequest](#servicesdocumentsListTemplatesRequest) | [ListTemplatesResponse](#servicesdocumentsListTemplatesResponse) | |
| `GetTemplate` | [GetTemplateRequest](#servicesdocumentsGetTemplateRequest) | [GetTemplateResponse](#servicesdocumentsGetTemplateResponse) | |
| `CreateTemplate` | [CreateTemplateRequest](#servicesdocumentsCreateTemplateRequest) | [CreateTemplateResponse](#servicesdocumentsCreateTemplateResponse) | |
| `UpdateTemplate` | [UpdateTemplateRequest](#servicesdocumentsUpdateTemplateRequest) | [UpdateTemplateResponse](#servicesdocumentsUpdateTemplateResponse) | |
| `DeleteTemplate` | [DeleteTemplateRequest](#servicesdocumentsDeleteTemplateRequest) | [DeleteTemplateResponse](#servicesdocumentsDeleteTemplateResponse) | |
| `ListDocuments` | [ListDocumentsRequest](#servicesdocumentsListDocumentsRequest) | [ListDocumentsResponse](#servicesdocumentsListDocumentsResponse) | |
| `GetDocument` | [GetDocumentRequest](#servicesdocumentsGetDocumentRequest) | [GetDocumentResponse](#servicesdocumentsGetDocumentResponse) | |
| `CreateDocument` | [CreateDocumentRequest](#servicesdocumentsCreateDocumentRequest) | [CreateDocumentResponse](#servicesdocumentsCreateDocumentResponse) | |
| `UpdateDocument` | [UpdateDocumentRequest](#servicesdocumentsUpdateDocumentRequest) | [UpdateDocumentResponse](#servicesdocumentsUpdateDocumentResponse) | |
| `DeleteDocument` | [DeleteDocumentRequest](#servicesdocumentsDeleteDocumentRequest) | [DeleteDocumentResponse](#servicesdocumentsDeleteDocumentResponse) | |
| `ToggleDocument` | [ToggleDocumentRequest](#servicesdocumentsToggleDocumentRequest) | [ToggleDocumentResponse](#servicesdocumentsToggleDocumentResponse) | |
| `ChangeDocumentOwner` | [ChangeDocumentOwnerRequest](#servicesdocumentsChangeDocumentOwnerRequest) | [ChangeDocumentOwnerResponse](#servicesdocumentsChangeDocumentOwnerResponse) | |
| `GetDocumentReferences` | [GetDocumentReferencesRequest](#servicesdocumentsGetDocumentReferencesRequest) | [GetDocumentReferencesResponse](#servicesdocumentsGetDocumentReferencesResponse) | |
| `GetDocumentRelations` | [GetDocumentRelationsRequest](#servicesdocumentsGetDocumentRelationsRequest) | [GetDocumentRelationsResponse](#servicesdocumentsGetDocumentRelationsResponse) | |
| `AddDocumentReference` | [AddDocumentReferenceRequest](#servicesdocumentsAddDocumentReferenceRequest) | [AddDocumentReferenceResponse](#servicesdocumentsAddDocumentReferenceResponse) | |
| `RemoveDocumentReference` | [RemoveDocumentReferenceRequest](#servicesdocumentsRemoveDocumentReferenceRequest) | [RemoveDocumentReferenceResponse](#servicesdocumentsRemoveDocumentReferenceResponse) | |
| `AddDocumentRelation` | [AddDocumentRelationRequest](#servicesdocumentsAddDocumentRelationRequest) | [AddDocumentRelationResponse](#servicesdocumentsAddDocumentRelationResponse) | |
| `RemoveDocumentRelation` | [RemoveDocumentRelationRequest](#servicesdocumentsRemoveDocumentRelationRequest) | [RemoveDocumentRelationResponse](#servicesdocumentsRemoveDocumentRelationResponse) | |
| `GetComments` | [GetCommentsRequest](#servicesdocumentsGetCommentsRequest) | [GetCommentsResponse](#servicesdocumentsGetCommentsResponse) | |
| `PostComment` | [PostCommentRequest](#servicesdocumentsPostCommentRequest) | [PostCommentResponse](#servicesdocumentsPostCommentResponse) | |
| `EditComment` | [EditCommentRequest](#servicesdocumentsEditCommentRequest) | [EditCommentResponse](#servicesdocumentsEditCommentResponse) | |
| `DeleteComment` | [DeleteCommentRequest](#servicesdocumentsDeleteCommentRequest) | [DeleteCommentResponse](#servicesdocumentsDeleteCommentResponse) | |
| `GetDocumentAccess` | [GetDocumentAccessRequest](#servicesdocumentsGetDocumentAccessRequest) | [GetDocumentAccessResponse](#servicesdocumentsGetDocumentAccessResponse) | |
| `SetDocumentAccess` | [SetDocumentAccessRequest](#servicesdocumentsSetDocumentAccessRequest) | [SetDocumentAccessResponse](#servicesdocumentsSetDocumentAccessResponse) | |
| `ListDocumentActivity` | [ListDocumentActivityRequest](#servicesdocumentsListDocumentActivityRequest) | [ListDocumentActivityResponse](#servicesdocumentsListDocumentActivityResponse) | |
| `ListDocumentReqs` | [ListDocumentReqsRequest](#servicesdocumentsListDocumentReqsRequest) | [ListDocumentReqsResponse](#servicesdocumentsListDocumentReqsResponse) | |
| `CreateDocumentReq` | [CreateDocumentReqRequest](#servicesdocumentsCreateDocumentReqRequest) | [CreateDocumentReqResponse](#servicesdocumentsCreateDocumentReqResponse) | |
| `UpdateDocumentReq` | [UpdateDocumentReqRequest](#servicesdocumentsUpdateDocumentReqRequest) | [UpdateDocumentReqResponse](#servicesdocumentsUpdateDocumentReqResponse) | |
| `DeleteDocumentReq` | [DeleteDocumentReqRequest](#servicesdocumentsDeleteDocumentReqRequest) | [DeleteDocumentReqResponse](#servicesdocumentsDeleteDocumentReqResponse) | |
| `ListUserDocuments` | [ListUserDocumentsRequest](#servicesdocumentsListUserDocumentsRequest) | [ListUserDocumentsResponse](#servicesdocumentsListUserDocumentsResponse) | |
| `ListCategories` | [ListCategoriesRequest](#servicesdocumentsListCategoriesRequest) | [ListCategoriesResponse](#servicesdocumentsListCategoriesResponse) | |
| `CreateOrUpdateCategory` | [CreateOrUpdateCategoryRequest](#servicesdocumentsCreateOrUpdateCategoryRequest) | [CreateOrUpdateCategoryResponse](#servicesdocumentsCreateOrUpdateCategoryResponse) | |
| `DeleteCategory` | [DeleteCategoryRequest](#servicesdocumentsDeleteCategoryRequest) | [DeleteCategoryResponse](#servicesdocumentsDeleteCategoryResponse) | |
| `ListDocumentPins` | [ListDocumentPinsRequest](#servicesdocumentsListDocumentPinsRequest) | [ListDocumentPinsResponse](#servicesdocumentsListDocumentPinsResponse) | |
| `ToggleDocumentPin` | [ToggleDocumentPinRequest](#servicesdocumentsToggleDocumentPinRequest) | [ToggleDocumentPinResponse](#servicesdocumentsToggleDocumentPinResponse) | |
| `SetDocumentReminder` | [SetDocumentReminderRequest](#servicesdocumentsSetDocumentReminderRequest) | [SetDocumentReminderResponse](#servicesdocumentsSetDocumentReminderResponse) | |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) | |

 <!-- end services -->



## services/documents/forms.proto


### services.documents.GetFormRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `form_id` | [int64](#int64) |  |  |





### services.documents.GetFormResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `form` | [resources.documents.forms.Form](#resourcesdocumentsformsForm) |  |  |





### services.documents.ListFormsRequest





### services.documents.ListFormsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `forms` | [resources.documents.forms.Form](#resourcesdocumentsformsForm) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.FormsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetForm` | [GetFormRequest](#servicesdocumentsGetFormRequest) | [GetFormResponse](#servicesdocumentsGetFormResponse) | |
| `ListForms` | [ListFormsRequest](#servicesdocumentsListFormsRequest) | [ListFormsResponse](#servicesdocumentsListFormsResponse) | |

 <!-- end services -->



## services/documents/stamps.proto


### services.documents.DeleteStampRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stamp_id` | [int64](#int64) |  |  |





### services.documents.DeleteStampResponse





### services.documents.ListUsableStampsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `document_id` | [int64](#int64) | optional | If set, only stamps usable for signing this document are returned |





### services.documents.ListUsableStampsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `stamps` | [resources.documents.stamps.Stamp](#resourcesdocumentsstampsStamp) | repeated |  |





### services.documents.UpsertStampRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stamp` | [resources.documents.stamps.Stamp](#resourcesdocumentsstampsStamp) |  |  |





### services.documents.UpsertStampResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stamp` | [resources.documents.stamps.Stamp](#resourcesdocumentsstampsStamp) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.StampsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListUsableStamps` | [ListUsableStampsRequest](#servicesdocumentsListUsableStampsRequest) | [ListUsableStampsResponse](#servicesdocumentsListUsableStampsResponse) | |
| `UpsertStamp` | [UpsertStampRequest](#servicesdocumentsUpsertStampRequest) | [UpsertStampResponse](#servicesdocumentsUpsertStampResponse) | |
| `DeleteStamp` | [DeleteStampRequest](#servicesdocumentsDeleteStampRequest) | [DeleteStampResponse](#servicesdocumentsDeleteStampResponse) | |

 <!-- end services -->



## services/filestore/filestore.proto


### services.filestore.DeleteFileByPathRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  |  |





### services.filestore.DeleteFileByPathResponse





### services.filestore.ListFilesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `path` | [string](#string) | optional |  |





### services.filestore.ListFilesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.filestore.FilestoreService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Upload` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `ListFiles` | [ListFilesRequest](#servicesfilestoreListFilesRequest) | [ListFilesResponse](#servicesfilestoreListFilesResponse) | |
| `DeleteFile` | [.resources.file.DeleteFileRequest](#resourcesfileDeleteFileRequest) | [.resources.file.DeleteFileResponse](#resourcesfileDeleteFileResponse) | |
| `DeleteFileByPath` | [DeleteFileByPathRequest](#servicesfilestoreDeleteFileByPathRequest) | [DeleteFileByPathResponse](#servicesfilestoreDeleteFileByPathResponse) | |

 <!-- end services -->



## services/jobs/conduct.proto


### services.jobs.CreateConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.conduct.ConductEntry](#resourcesjobsconductConductEntry) |  |  |





### services.jobs.CreateConductEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.conduct.ConductEntry](#resourcesjobsconductConductEntry) |  |  |





### services.jobs.DeleteConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.jobs.DeleteConductEntryResponse





### services.jobs.GetConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.jobs.GetConductEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.conduct.ConductEntry](#resourcesjobsconductConductEntry) |  |  |





### services.jobs.ListConductEntriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `types` | [resources.jobs.conduct.ConductType](#resourcesjobsconductConductType) | repeated | Search params |
| `show_expired` | [bool](#bool) | optional |  |
| `show_drafts` | [bool](#bool) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `ids` | [int64](#int64) | repeated |  |





### services.jobs.ListConductEntriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `entries` | [resources.jobs.conduct.ConductEntry](#resourcesjobsconductConductEntry) | repeated |  |





### services.jobs.UpdateConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.conduct.ConductEntry](#resourcesjobsconductConductEntry) |  |  |





### services.jobs.UpdateConductEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.conduct.ConductEntry](#resourcesjobsconductConductEntry) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.jobs.ConductService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListConductEntries` | [ListConductEntriesRequest](#servicesjobsListConductEntriesRequest) | [ListConductEntriesResponse](#servicesjobsListConductEntriesResponse) | |
| `GetConductEntry` | [GetConductEntryRequest](#servicesjobsGetConductEntryRequest) | [GetConductEntryResponse](#servicesjobsGetConductEntryResponse) | |
| `CreateConductEntry` | [CreateConductEntryRequest](#servicesjobsCreateConductEntryRequest) | [CreateConductEntryResponse](#servicesjobsCreateConductEntryResponse) | |
| `UpdateConductEntry` | [UpdateConductEntryRequest](#servicesjobsUpdateConductEntryRequest) | [UpdateConductEntryResponse](#servicesjobsUpdateConductEntryResponse) | |
| `DeleteConductEntry` | [DeleteConductEntryRequest](#servicesjobsDeleteConductEntryRequest) | [DeleteConductEntryResponse](#servicesjobsDeleteConductEntryResponse) | |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) | |

 <!-- end services -->



## services/jobs/jobs.proto


### services.jobs.GetColleagueLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) | optional |  |





### services.jobs.GetColleagueLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.labels.Label](#resourcesjobslabelsLabel) | repeated |  |





### services.jobs.GetColleagueLabelsStatsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `label_ids` | [int64](#int64) | repeated |  |





### services.jobs.GetColleagueLabelsStatsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `count` | [resources.jobs.labels.LabelCount](#resourcesjobslabelsLabelCount) | repeated |  |





### services.jobs.GetColleagueRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `info_only` | [bool](#bool) | optional |  |





### services.jobs.GetColleagueResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) |  |  |





### services.jobs.GetMOTDRequest





### services.jobs.GetMOTDResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  |  |





### services.jobs.GetSelfRequest





### services.jobs.GetSelfResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) |  |  |





### services.jobs.ListColleagueActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_ids` | [int32](#int32) | repeated | Search params |
| `activity_types` | [resources.jobs.colleagues.activity.ColleagueActivityType](#resourcesjobscolleaguesactivityColleagueActivityType) | repeated |  |





### services.jobs.ListColleagueActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.jobs.colleagues.activity.ColleagueActivity](#resourcesjobscolleaguesactivityColleagueActivity) | repeated |  |





### services.jobs.ListColleaguesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `search` | [string](#string) |  | Search params |
| `user_ids` | [int32](#int32) | repeated |  |
| `user_only` | [bool](#bool) | optional |  |
| `absent` | [bool](#bool) | optional |  |
| `label_ids` | [int64](#int64) | repeated |  |
| `name_prefix` | [string](#string) | optional |  |
| `name_suffix` | [string](#string) | optional |  |





### services.jobs.ListColleaguesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `colleagues` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | repeated |  |





### services.jobs.ManageLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.labels.Label](#resourcesjobslabelsLabel) | repeated |  |





### services.jobs.ManageLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.labels.Label](#resourcesjobslabelsLabel) | repeated |  |





### services.jobs.SetColleaguePropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.jobs.colleagues.ColleagueProps](#resourcesjobscolleaguesColleagueProps) |  |  |
| `reason` | [string](#string) |  |  |





### services.jobs.SetColleaguePropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.jobs.colleagues.ColleagueProps](#resourcesjobscolleaguesColleagueProps) |  |  |





### services.jobs.SetMOTDRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  |  |





### services.jobs.SetMOTDResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.jobs.JobsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListColleagues` | [ListColleaguesRequest](#servicesjobsListColleaguesRequest) | [ListColleaguesResponse](#servicesjobsListColleaguesResponse) | |
| `GetSelf` | [GetSelfRequest](#servicesjobsGetSelfRequest) | [GetSelfResponse](#servicesjobsGetSelfResponse) | |
| `GetColleague` | [GetColleagueRequest](#servicesjobsGetColleagueRequest) | [GetColleagueResponse](#servicesjobsGetColleagueResponse) | |
| `ListColleagueActivity` | [ListColleagueActivityRequest](#servicesjobsListColleagueActivityRequest) | [ListColleagueActivityResponse](#servicesjobsListColleagueActivityResponse) | |
| `SetColleagueProps` | [SetColleaguePropsRequest](#servicesjobsSetColleaguePropsRequest) | [SetColleaguePropsResponse](#servicesjobsSetColleaguePropsResponse) | |
| `GetColleagueLabels` | [GetColleagueLabelsRequest](#servicesjobsGetColleagueLabelsRequest) | [GetColleagueLabelsResponse](#servicesjobsGetColleagueLabelsResponse) | |
| `ManageLabels` | [ManageLabelsRequest](#servicesjobsManageLabelsRequest) | [ManageLabelsResponse](#servicesjobsManageLabelsResponse) | |
| `GetColleagueLabelsStats` | [GetColleagueLabelsStatsRequest](#servicesjobsGetColleagueLabelsStatsRequest) | [GetColleagueLabelsStatsResponse](#servicesjobsGetColleagueLabelsStatsResponse) | |
| `GetMOTD` | [GetMOTDRequest](#servicesjobsGetMOTDRequest) | [GetMOTDResponse](#servicesjobsGetMOTDResponse) | |
| `SetMOTD` | [SetMOTDRequest](#servicesjobsSetMOTDRequest) | [SetMOTDResponse](#servicesjobsSetMOTDResponse) | |

 <!-- end services -->



## services/jobs/timeclock.proto


### services.jobs.GetTimeclockStatsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) | optional |  |





### services.jobs.GetTimeclockStatsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stats` | [resources.jobs.timeclock.TimeclockStats](#resourcesjobstimeclockTimeclockStats) |  |  |
| `weekly` | [resources.jobs.timeclock.TimeclockWeeklyStats](#resourcesjobstimeclockTimeclockWeeklyStats) | repeated |  |





### services.jobs.ListInactiveEmployeesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `days` | [int32](#int32) |  | Search params |





### services.jobs.ListInactiveEmployeesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `colleagues` | [resources.jobs.colleagues.Colleague](#resourcesjobscolleaguesColleague) | repeated |  |





### services.jobs.ListTimeclockRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_mode` | [resources.jobs.timeclock.TimeclockViewMode](#resourcesjobstimeclockTimeclockViewMode) |  | Search params |
| `mode` | [resources.jobs.timeclock.TimeclockMode](#resourcesjobstimeclockTimeclockMode) |  |  |
| `date` | [resources.common.database.DateRange](#resourcescommondatabaseDateRange) | optional |  |
| `per_day` | [bool](#bool) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |





### services.jobs.ListTimeclockResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `stats` | [resources.jobs.timeclock.TimeclockStats](#resourcesjobstimeclockTimeclockStats) |  |  |
| `stats_weekly` | [resources.jobs.timeclock.TimeclockWeeklyStats](#resourcesjobstimeclockTimeclockWeeklyStats) | repeated |  |
| `daily` | [TimeclockDay](#servicesjobsTimeclockDay) |  |  |
| `weekly` | [TimeclockWeekly](#servicesjobsTimeclockWeekly) |  |  |
| `range` | [TimeclockRange](#servicesjobsTimeclockRange) |  |  |





### services.jobs.TimeclockDay


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `entries` | [resources.jobs.timeclock.TimeclockEntry](#resourcesjobstimeclockTimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |





### services.jobs.TimeclockRange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `entries` | [resources.jobs.timeclock.TimeclockEntry](#resourcesjobstimeclockTimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |





### services.jobs.TimeclockWeekly


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `entries` | [resources.jobs.timeclock.TimeclockEntry](#resourcesjobstimeclockTimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.jobs.TimeclockService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTimeclock` | [ListTimeclockRequest](#servicesjobsListTimeclockRequest) | [ListTimeclockResponse](#servicesjobsListTimeclockResponse) | |
| `GetTimeclockStats` | [GetTimeclockStatsRequest](#servicesjobsGetTimeclockStatsRequest) | [GetTimeclockStatsResponse](#servicesjobsGetTimeclockStatsResponse) | |
| `ListInactiveEmployees` | [ListInactiveEmployeesRequest](#servicesjobsListInactiveEmployeesRequest) | [ListInactiveEmployeesResponse](#servicesjobsListInactiveEmployeesResponse) | |

 <!-- end services -->



## services/livemap/livemap.proto


### services.livemap.CreateOrUpdateMarkerRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `marker` | [resources.livemap.markers.MarkerMarker](#resourceslivemapmarkersMarkerMarker) |  |  |





### services.livemap.CreateOrUpdateMarkerResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `marker` | [resources.livemap.markers.MarkerMarker](#resourceslivemapmarkersMarkerMarker) |  |  |





### services.livemap.DeleteMarkerRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.livemap.DeleteMarkerResponse





### services.livemap.JobsList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |
| `markers` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |





### services.livemap.MarkerMarkersUpdates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [resources.livemap.markers.MarkerMarker](#resourceslivemapmarkersMarkerMarker) | repeated |  |
| `deleted` | [int64](#int64) | repeated |  |
| `part` | [int32](#int32) |  |  |
| `partial` | [bool](#bool) |  |  |





### services.livemap.Snapshot
A roll-up of the entire USERLOC bucket. Published every N seconds on `$KV.user_locations._snapshot` with the headers: Nats-Rollup: all KV-Operation: ROLLUP



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `markers` | [resources.livemap.markers.UserMarker](#resourceslivemapmarkersUserMarker) | repeated | All currently-known user markers, already filtered for obsolete PURGE/DELETE events. |





### services.livemap.StreamRequest





### services.livemap.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_on_duty` | [bool](#bool) | optional |  |
| `jobs` | [JobsList](#serviceslivemapJobsList) |  |  |
| `markers` | [MarkerMarkersUpdates](#serviceslivemapMarkerMarkersUpdates) |  |  |
| `snapshot` | [Snapshot](#serviceslivemapSnapshot) |  |  |
| `user_updates` | [UserUpdates](#serviceslivemapUserUpdates) |  |  |
| `user_deletes` | [UserDeletes](#serviceslivemapUserDeletes) |  |  |





### services.livemap.UserDelete


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int32](#int32) |  | The user ID of an user marker that was deleted. |
| `job` | [string](#string) |  | The job of the user that was deleted. |





### services.livemap.UserDeletes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deletes` | [UserDelete](#serviceslivemapUserDelete) | repeated |  |





### services.livemap.UserUpdates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updates` | [resources.livemap.markers.UserMarker](#resourceslivemapmarkersUserMarker) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.livemap.LivemapService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Stream` | [StreamRequest](#serviceslivemapStreamRequest) | [StreamResponse](#serviceslivemapStreamResponse) stream | |
| `CreateOrUpdateMarker` | [CreateOrUpdateMarkerRequest](#serviceslivemapCreateOrUpdateMarkerRequest) | [CreateOrUpdateMarkerResponse](#serviceslivemapCreateOrUpdateMarkerResponse) | |
| `DeleteMarker` | [DeleteMarkerRequest](#serviceslivemapDeleteMarkerRequest) | [DeleteMarkerResponse](#serviceslivemapDeleteMarkerResponse) | |

 <!-- end services -->



## services/mailer/mailer.proto


### services.mailer.CreateOrUpdateEmailRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) |  |  |





### services.mailer.CreateOrUpdateEmailResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) |  |  |





### services.mailer.CreateOrUpdateTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.templates.Template](#resourcesmailertemplatesTemplate) |  |  |





### services.mailer.CreateOrUpdateTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.templates.Template](#resourcesmailertemplatesTemplate) |  |  |





### services.mailer.CreateThreadRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.threads.Thread](#resourcesmailerthreadsThread) |  |  |
| `message` | [resources.mailer.messages.Message](#resourcesmailermessagesMessage) |  |  |
| `recipients` | [string](#string) | repeated |  |





### services.mailer.CreateThreadResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.threads.Thread](#resourcesmailerthreadsThread) |  |  |





### services.mailer.DeleteEmailRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.mailer.DeleteEmailResponse





### services.mailer.DeleteMessageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `thread_id` | [int64](#int64) |  |  |
| `message_id` | [int64](#int64) |  |  |





### services.mailer.DeleteMessageResponse





### services.mailer.DeleteTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `id` | [int64](#int64) |  |  |





### services.mailer.DeleteTemplateResponse





### services.mailer.DeleteThreadRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `thread_id` | [int64](#int64) |  |  |





### services.mailer.DeleteThreadResponse





### services.mailer.GetEmailProposalsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `input` | [string](#string) |  |  |
| `job` | [bool](#bool) | optional |  |
| `user_id` | [int32](#int32) | optional |  |





### services.mailer.GetEmailProposalsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `emails` | [string](#string) | repeated |  |
| `domains` | [string](#string) | repeated |  |





### services.mailer.GetEmailRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.mailer.GetEmailResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) |  |  |





### services.mailer.GetEmailSettingsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |





### services.mailer.GetEmailSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.settings.EmailSettings](#resourcesmailersettingsEmailSettings) |  |  |





### services.mailer.GetTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `template_id` | [int64](#int64) |  |  |





### services.mailer.GetTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.templates.Template](#resourcesmailertemplatesTemplate) |  |  |





### services.mailer.GetThreadRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `thread_id` | [int64](#int64) |  |  |





### services.mailer.GetThreadResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.threads.Thread](#resourcesmailerthreadsThread) |  |  |





### services.mailer.GetThreadStateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |
| `thread_id` | [int64](#int64) |  |  |





### services.mailer.GetThreadStateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.threads.ThreadState](#resourcesmailerthreadsThreadState) |  |  |





### services.mailer.ListEmailsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `all` | [bool](#bool) | optional | Search params |





### services.mailer.ListEmailsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `emails` | [resources.mailer.emails.Email](#resourcesmaileremailsEmail) | repeated |  |





### services.mailer.ListTemplatesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [int64](#int64) |  |  |





### services.mailer.ListTemplatesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `templates` | [resources.mailer.templates.Template](#resourcesmailertemplatesTemplate) | repeated |  |





### services.mailer.ListThreadMessagesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `email_id` | [int64](#int64) |  |  |
| `thread_id` | [int64](#int64) |  |  |
| `after` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### services.mailer.ListThreadMessagesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `messages` | [resources.mailer.messages.Message](#resourcesmailermessagesMessage) | repeated |  |





### services.mailer.ListThreadsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `email_ids` | [int64](#int64) | repeated | Search params |
| `unread` | [bool](#bool) | optional |  |
| `archived` | [bool](#bool) | optional |  |





### services.mailer.ListThreadsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `threads` | [resources.mailer.threads.Thread](#resourcesmailerthreadsThread) | repeated |  |





### services.mailer.PostMessageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [resources.mailer.messages.Message](#resourcesmailermessagesMessage) |  |  |
| `recipients` | [string](#string) | repeated |  |





### services.mailer.PostMessageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [resources.mailer.messages.Message](#resourcesmailermessagesMessage) |  |  |





### services.mailer.SearchThreadsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `search` | [string](#string) |  | Search params |





### services.mailer.SearchThreadsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `messages` | [resources.mailer.messages.Message](#resourcesmailermessagesMessage) | repeated |  |





### services.mailer.SetEmailSettingsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.settings.EmailSettings](#resourcesmailersettingsEmailSettings) |  |  |





### services.mailer.SetEmailSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.settings.EmailSettings](#resourcesmailersettingsEmailSettings) |  |  |





### services.mailer.SetThreadStateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.threads.ThreadState](#resourcesmailerthreadsThreadState) |  |  |





### services.mailer.SetThreadStateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.threads.ThreadState](#resourcesmailerthreadsThreadState) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.mailer.MailerService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListEmails` | [ListEmailsRequest](#servicesmailerListEmailsRequest) | [ListEmailsResponse](#servicesmailerListEmailsResponse) | |
| `GetEmail` | [GetEmailRequest](#servicesmailerGetEmailRequest) | [GetEmailResponse](#servicesmailerGetEmailResponse) | |
| `CreateOrUpdateEmail` | [CreateOrUpdateEmailRequest](#servicesmailerCreateOrUpdateEmailRequest) | [CreateOrUpdateEmailResponse](#servicesmailerCreateOrUpdateEmailResponse) | |
| `DeleteEmail` | [DeleteEmailRequest](#servicesmailerDeleteEmailRequest) | [DeleteEmailResponse](#servicesmailerDeleteEmailResponse) | |
| `GetEmailProposals` | [GetEmailProposalsRequest](#servicesmailerGetEmailProposalsRequest) | [GetEmailProposalsResponse](#servicesmailerGetEmailProposalsResponse) | |
| `ListTemplates` | [ListTemplatesRequest](#servicesmailerListTemplatesRequest) | [ListTemplatesResponse](#servicesmailerListTemplatesResponse) | |
| `GetTemplate` | [GetTemplateRequest](#servicesmailerGetTemplateRequest) | [GetTemplateResponse](#servicesmailerGetTemplateResponse) | |
| `CreateOrUpdateTemplate` | [CreateOrUpdateTemplateRequest](#servicesmailerCreateOrUpdateTemplateRequest) | [CreateOrUpdateTemplateResponse](#servicesmailerCreateOrUpdateTemplateResponse) | |
| `DeleteTemplate` | [DeleteTemplateRequest](#servicesmailerDeleteTemplateRequest) | [DeleteTemplateResponse](#servicesmailerDeleteTemplateResponse) | |
| `ListThreads` | [ListThreadsRequest](#servicesmailerListThreadsRequest) | [ListThreadsResponse](#servicesmailerListThreadsResponse) | |
| `GetThread` | [GetThreadRequest](#servicesmailerGetThreadRequest) | [GetThreadResponse](#servicesmailerGetThreadResponse) | |
| `CreateThread` | [CreateThreadRequest](#servicesmailerCreateThreadRequest) | [CreateThreadResponse](#servicesmailerCreateThreadResponse) | |
| `DeleteThread` | [DeleteThreadRequest](#servicesmailerDeleteThreadRequest) | [DeleteThreadResponse](#servicesmailerDeleteThreadResponse) | |
| `GetThreadState` | [GetThreadStateRequest](#servicesmailerGetThreadStateRequest) | [GetThreadStateResponse](#servicesmailerGetThreadStateResponse) | |
| `SetThreadState` | [SetThreadStateRequest](#servicesmailerSetThreadStateRequest) | [SetThreadStateResponse](#servicesmailerSetThreadStateResponse) | |
| `SearchThreads` | [SearchThreadsRequest](#servicesmailerSearchThreadsRequest) | [SearchThreadsResponse](#servicesmailerSearchThreadsResponse) | |
| `ListThreadMessages` | [ListThreadMessagesRequest](#servicesmailerListThreadMessagesRequest) | [ListThreadMessagesResponse](#servicesmailerListThreadMessagesResponse) | |
| `PostMessage` | [PostMessageRequest](#servicesmailerPostMessageRequest) | [PostMessageResponse](#servicesmailerPostMessageResponse) | |
| `DeleteMessage` | [DeleteMessageRequest](#servicesmailerDeleteMessageRequest) | [DeleteMessageResponse](#servicesmailerDeleteMessageResponse) | |
| `GetEmailSettings` | [GetEmailSettingsRequest](#servicesmailerGetEmailSettingsRequest) | [GetEmailSettingsResponse](#servicesmailerGetEmailSettingsResponse) | |
| `SetEmailSettings` | [SetEmailSettingsRequest](#servicesmailerSetEmailSettingsRequest) | [SetEmailSettingsResponse](#servicesmailerSetEmailSettingsResponse) | |

 <!-- end services -->



## services/notifications/notifications.proto


### services.notifications.GetNotificationsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `include_read` | [bool](#bool) | optional |  |
| `categories` | [resources.notifications.NotificationCategory](#resourcesnotificationsNotificationCategory) | repeated |  |





### services.notifications.GetNotificationsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `notifications` | [resources.notifications.Notification](#resourcesnotificationsNotification) | repeated |  |





### services.notifications.MarkNotificationsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unread` | [bool](#bool) |  |  |
| `ids` | [int64](#int64) | repeated |  |
| `all` | [bool](#bool) | optional |  |





### services.notifications.MarkNotificationsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [int64](#int64) |  |  |





### services.notifications.StreamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `clientview` | [resources.notifications.clientview.ClientView](#resourcesnotificationsclientviewClientView) |  |  |





### services.notifications.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `notification_count` | [int64](#int64) |  |  |
| `restart` | [bool](#bool) | optional |  |
| `user_event` | [resources.notifications.events.UserEvent](#resourcesnotificationseventsUserEvent) |  |  |
| `job_event` | [resources.notifications.events.JobEvent](#resourcesnotificationseventsJobEvent) |  |  |
| `job_grade_event` | [resources.notifications.events.JobGradeEvent](#resourcesnotificationseventsJobGradeEvent) |  |  |
| `system_event` | [resources.notifications.events.SystemEvent](#resourcesnotificationseventsSystemEvent) |  |  |
| `mailer_event` | [resources.mailer.events.MailerEvent](#resourcesmailereventsMailerEvent) |  |  |
| `object_event` | [resources.notifications.clientview.ObjectEvent](#resourcesnotificationsclientviewObjectEvent) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.notifications.NotificationsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetNotifications` | [GetNotificationsRequest](#servicesnotificationsGetNotificationsRequest) | [GetNotificationsResponse](#servicesnotificationsGetNotificationsResponse) | |
| `MarkNotifications` | [MarkNotificationsRequest](#servicesnotificationsMarkNotificationsRequest) | [MarkNotificationsResponse](#servicesnotificationsMarkNotificationsResponse) | |
| `Stream` | [StreamRequest](#servicesnotificationsStreamRequest) stream | [StreamResponse](#servicesnotificationsStreamResponse) stream | |

 <!-- end services -->



## services/qualifications/qualifications.proto


### services.qualifications.CreateOrUpdateQualificationRequestRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.qualifications.QualificationRequest](#resourcesqualificationsQualificationRequest) |  |  |





### services.qualifications.CreateOrUpdateQualificationRequestResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.qualifications.QualificationRequest](#resourcesqualificationsQualificationRequest) |  |  |





### services.qualifications.CreateOrUpdateQualificationResultRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [resources.qualifications.QualificationResult](#resourcesqualificationsQualificationResult) |  |  |
| `grading` | [resources.qualifications.exam.ExamGrading](#resourcesqualificationsexamExamGrading) | optional |  |





### services.qualifications.CreateOrUpdateQualificationResultResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [resources.qualifications.QualificationResult](#resourcesqualificationsQualificationResult) |  |  |





### services.qualifications.CreateQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |





### services.qualifications.CreateQualificationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |





### services.qualifications.DeleteQualificationReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |





### services.qualifications.DeleteQualificationReqResponse





### services.qualifications.DeleteQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |





### services.qualifications.DeleteQualificationResponse





### services.qualifications.DeleteQualificationResultRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result_id` | [int64](#int64) |  |  |





### services.qualifications.DeleteQualificationResultResponse





### services.qualifications.GetExamInfoRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |





### services.qualifications.GetExamInfoResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.QualificationShort](#resourcesqualificationsQualificationShort) |  |  |
| `question_count` | [int64](#int64) |  |  |
| `exam_user` | [resources.qualifications.exam.ExamUser](#resourcesqualificationsexamExamUser) | optional |  |





### services.qualifications.GetQualificationAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |





### services.qualifications.GetQualificationAccessResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [resources.qualifications.access.QualificationAccess](#resourcesqualificationsaccessQualificationAccess) |  |  |





### services.qualifications.GetQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `with_exam` | [bool](#bool) | optional |  |





### services.qualifications.GetQualificationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.Qualification](#resourcesqualificationsQualification) |  |  |





### services.qualifications.GetUserExamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `user_id` | [int32](#int32) |  |  |





### services.qualifications.GetUserExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exam` | [resources.qualifications.exam.ExamQuestions](#resourcesqualificationsexamExamQuestions) |  |  |
| `exam_user` | [resources.qualifications.exam.ExamUser](#resourcesqualificationsexamExamUser) |  |  |
| `responses` | [resources.qualifications.exam.ExamResponses](#resourcesqualificationsexamExamResponses) |  |  |
| `grading` | [resources.qualifications.exam.ExamGrading](#resourcesqualificationsexamExamGrading) |  |  |





### services.qualifications.ListQualificationRequestsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `qualification_id` | [int64](#int64) | optional | Search params |
| `status` | [resources.qualifications.RequestStatus](#resourcesqualificationsRequestStatus) | repeated |  |
| `user_id` | [int32](#int32) | optional |  |





### services.qualifications.ListQualificationRequestsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `requests` | [resources.qualifications.QualificationRequest](#resourcesqualificationsQualificationRequest) | repeated |  |





### services.qualifications.ListQualificationsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `search` | [string](#string) | optional | Search params |
| `job` | [string](#string) | optional |  |





### services.qualifications.ListQualificationsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `qualifications` | [resources.qualifications.Qualification](#resourcesqualificationsQualification) | repeated |  |





### services.qualifications.ListQualificationsResultsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `qualification_id` | [int64](#int64) | optional | Search params |
| `status` | [resources.qualifications.ResultStatus](#resourcesqualificationsResultStatus) | repeated |  |
| `user_id` | [int32](#int32) | optional |  |





### services.qualifications.ListQualificationsResultsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `results` | [resources.qualifications.QualificationResult](#resourcesqualificationsQualificationResult) | repeated |  |





### services.qualifications.SetQualificationAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `access` | [resources.qualifications.access.QualificationAccess](#resourcesqualificationsaccessQualificationAccess) |  |  |





### services.qualifications.SetQualificationAccessResponse





### services.qualifications.SubmitExamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `responses` | [resources.qualifications.exam.ExamResponses](#resourcesqualificationsexamExamResponses) |  |  |





### services.qualifications.SubmitExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |





### services.qualifications.TakeExamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |
| `cancel` | [bool](#bool) | optional |  |





### services.qualifications.TakeExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exam` | [resources.qualifications.exam.ExamQuestions](#resourcesqualificationsexamExamQuestions) |  |  |
| `exam_user` | [resources.qualifications.exam.ExamUser](#resourcesqualificationsexamExamUser) |  |  |





### services.qualifications.UpdateQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.Qualification](#resourcesqualificationsQualification) |  |  |





### services.qualifications.UpdateQualificationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.qualifications.QualificationsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListQualifications` | [ListQualificationsRequest](#servicesqualificationsListQualificationsRequest) | [ListQualificationsResponse](#servicesqualificationsListQualificationsResponse) | |
| `GetQualification` | [GetQualificationRequest](#servicesqualificationsGetQualificationRequest) | [GetQualificationResponse](#servicesqualificationsGetQualificationResponse) | |
| `CreateQualification` | [CreateQualificationRequest](#servicesqualificationsCreateQualificationRequest) | [CreateQualificationResponse](#servicesqualificationsCreateQualificationResponse) | |
| `UpdateQualification` | [UpdateQualificationRequest](#servicesqualificationsUpdateQualificationRequest) | [UpdateQualificationResponse](#servicesqualificationsUpdateQualificationResponse) | |
| `DeleteQualification` | [DeleteQualificationRequest](#servicesqualificationsDeleteQualificationRequest) | [DeleteQualificationResponse](#servicesqualificationsDeleteQualificationResponse) | |
| `ListQualificationRequests` | [ListQualificationRequestsRequest](#servicesqualificationsListQualificationRequestsRequest) | [ListQualificationRequestsResponse](#servicesqualificationsListQualificationRequestsResponse) | |
| `CreateOrUpdateQualificationRequest` | [CreateOrUpdateQualificationRequestRequest](#servicesqualificationsCreateOrUpdateQualificationRequestRequest) | [CreateOrUpdateQualificationRequestResponse](#servicesqualificationsCreateOrUpdateQualificationRequestResponse) | |
| `DeleteQualificationReq` | [DeleteQualificationReqRequest](#servicesqualificationsDeleteQualificationReqRequest) | [DeleteQualificationReqResponse](#servicesqualificationsDeleteQualificationReqResponse) | |
| `ListQualificationsResults` | [ListQualificationsResultsRequest](#servicesqualificationsListQualificationsResultsRequest) | [ListQualificationsResultsResponse](#servicesqualificationsListQualificationsResultsResponse) | |
| `CreateOrUpdateQualificationResult` | [CreateOrUpdateQualificationResultRequest](#servicesqualificationsCreateOrUpdateQualificationResultRequest) | [CreateOrUpdateQualificationResultResponse](#servicesqualificationsCreateOrUpdateQualificationResultResponse) | |
| `DeleteQualificationResult` | [DeleteQualificationResultRequest](#servicesqualificationsDeleteQualificationResultRequest) | [DeleteQualificationResultResponse](#servicesqualificationsDeleteQualificationResultResponse) | |
| `GetExamInfo` | [GetExamInfoRequest](#servicesqualificationsGetExamInfoRequest) | [GetExamInfoResponse](#servicesqualificationsGetExamInfoResponse) | |
| `TakeExam` | [TakeExamRequest](#servicesqualificationsTakeExamRequest) | [TakeExamResponse](#servicesqualificationsTakeExamResponse) | |
| `SubmitExam` | [SubmitExamRequest](#servicesqualificationsSubmitExamRequest) | [SubmitExamResponse](#servicesqualificationsSubmitExamResponse) | |
| `GetUserExam` | [GetUserExamRequest](#servicesqualificationsGetUserExamRequest) | [GetUserExamResponse](#servicesqualificationsGetUserExamResponse) | |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) | |

 <!-- end services -->



## services/settings/accounts.proto


### services.settings.CreateAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `license` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `last_char` | [int32](#int32) | optional |  |
| `char` | [resources.users.short.UserShort](#resourcesusersshortUserShort) | optional | Allow creating a char at the same time (only when dbsync is used) |





### services.settings.CreateAccountResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) |  |  |





### services.settings.DeleteAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.settings.DeleteAccountResponse





### services.settings.DisconnectSocialLoginRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `provider_name` | [string](#string) |  |  |





### services.settings.DisconnectSocialLoginResponse





### services.settings.ListAccountsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `license` | [string](#string) | optional | Search params |
| `enabled` | [bool](#bool) | optional |  |
| `username` | [string](#string) | optional |  |
| `external_id` | [string](#string) | optional |  |
| `group` | [string](#string) | optional |  |





### services.settings.ListAccountsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `accounts` | [resources.accounts.Account](#resourcesaccountsAccount) | repeated |  |





### services.settings.UpdateAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `enabled` | [bool](#bool) | optional |  |
| `last_char` | [int32](#int32) | optional |  |





### services.settings.UpdateAccountResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [resources.accounts.Account](#resourcesaccountsAccount) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.AccountsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListAccounts` | [ListAccountsRequest](#servicessettingsListAccountsRequest) | [ListAccountsResponse](#servicessettingsListAccountsResponse) | |
| `CreateAccount` | [CreateAccountRequest](#servicessettingsCreateAccountRequest) | [CreateAccountResponse](#servicessettingsCreateAccountResponse) | |
| `UpdateAccount` | [UpdateAccountRequest](#servicessettingsUpdateAccountRequest) | [UpdateAccountResponse](#servicessettingsUpdateAccountResponse) | |
| `DisconnectSocialLogin` | [DisconnectSocialLoginRequest](#servicessettingsDisconnectSocialLoginRequest) | [DisconnectSocialLoginResponse](#servicessettingsDisconnectSocialLoginResponse) | |
| `DeleteAccount` | [DeleteAccountRequest](#servicessettingsDeleteAccountRequest) | [DeleteAccountResponse](#servicessettingsDeleteAccountResponse) | |

 <!-- end services -->



## services/settings/config.proto


### services.settings.GetAppConfigRequest





### services.settings.GetAppConfigResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [resources.settings.AppConfig](#resourcessettingsAppConfig) |  |  |





### services.settings.UpdateAppConfigRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [resources.settings.AppConfig](#resourcessettingsAppConfig) |  |  |





### services.settings.UpdateAppConfigResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [resources.settings.AppConfig](#resourcessettingsAppConfig) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.ConfigService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetAppConfig` | [GetAppConfigRequest](#servicessettingsGetAppConfigRequest) | [GetAppConfigResponse](#servicessettingsGetAppConfigResponse) | |
| `UpdateAppConfig` | [UpdateAppConfigRequest](#servicessettingsUpdateAppConfigRequest) | [UpdateAppConfigResponse](#servicessettingsUpdateAppConfigResponse) | |

 <!-- end services -->



## services/settings/cron.proto


### services.settings.ListCronjobsRequest





### services.settings.ListCronjobsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.cron.Cronjob](#resourcescronCronjob) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.CronService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCronjobs` | [ListCronjobsRequest](#servicessettingsListCronjobsRequest) | [ListCronjobsResponse](#servicessettingsListCronjobsResponse) | |

 <!-- end services -->



## services/settings/laws.proto


### services.settings.CreateOrUpdateLawBookRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `law_book` | [resources.laws.LawBook](#resourceslawsLawBook) |  |  |





### services.settings.CreateOrUpdateLawBookResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `law_book` | [resources.laws.LawBook](#resourceslawsLawBook) |  |  |





### services.settings.CreateOrUpdateLawRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `law` | [resources.laws.Law](#resourceslawsLaw) |  |  |





### services.settings.CreateOrUpdateLawResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `law` | [resources.laws.Law](#resourceslawsLaw) |  |  |





### services.settings.DeleteLawBookRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.settings.DeleteLawBookResponse





### services.settings.DeleteLawRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.settings.DeleteLawResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.LawsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `CreateOrUpdateLawBook` | [CreateOrUpdateLawBookRequest](#servicessettingsCreateOrUpdateLawBookRequest) | [CreateOrUpdateLawBookResponse](#servicessettingsCreateOrUpdateLawBookResponse) | |
| `DeleteLawBook` | [DeleteLawBookRequest](#servicessettingsDeleteLawBookRequest) | [DeleteLawBookResponse](#servicessettingsDeleteLawBookResponse) | |
| `CreateOrUpdateLaw` | [CreateOrUpdateLawRequest](#servicessettingsCreateOrUpdateLawRequest) | [CreateOrUpdateLawResponse](#servicessettingsCreateOrUpdateLawResponse) | |
| `DeleteLaw` | [DeleteLawRequest](#servicessettingsDeleteLawRequest) | [DeleteLawResponse](#servicessettingsDeleteLawResponse) | |

 <!-- end services -->



## services/settings/settings.proto


### services.settings.CreateRoleRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |





### services.settings.CreateRoleResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.permissions.Role](#resourcespermissionspermissionsRole) |  |  |





### services.settings.DeleteJobLogoRequest





### services.settings.DeleteJobLogoResponse





### services.settings.DeleteRoleRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.settings.DeleteRoleResponse





### services.settings.GetEffectivePermissionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [int64](#int64) |  |  |





### services.settings.GetEffectivePermissionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.permissions.Role](#resourcespermissionspermissionsRole) |  |  |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |





### services.settings.GetJobPropsRequest





### services.settings.GetJobPropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.props.JobProps](#resourcesjobspropsJobProps) |  |  |





### services.settings.GetPermissionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [int64](#int64) |  |  |





### services.settings.GetPermissionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |





### services.settings.GetRoleRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.settings.GetRoleResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.permissions.Role](#resourcespermissionspermissionsRole) |  |  |





### services.settings.GetRolesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `lowest_rank` | [bool](#bool) | optional |  |





### services.settings.GetRolesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `roles` | [resources.permissions.permissions.Role](#resourcespermissionspermissionsRole) | repeated |  |





### services.settings.ListDiscordChannelsRequest





### services.settings.ListDiscordChannelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [resources.discord.Channel](#resourcesdiscordChannel) | repeated |  |





### services.settings.ListUserGuildsRequest





### services.settings.ListUserGuildsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `guilds` | [resources.discord.Guild](#resourcesdiscordGuild) | repeated |  |





### services.settings.SetJobPropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.props.JobProps](#resourcesjobspropsJobProps) |  |  |





### services.settings.SetJobPropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.props.JobProps](#resourcesjobspropsJobProps) |  |  |





### services.settings.UpdateRolePermsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |
| `perms` | [resources.settings.PermsUpdate](#resourcessettingsPermsUpdate) | optional |  |
| `attrs` | [resources.settings.AttrsUpdate](#resourcessettingsAttrsUpdate) | optional |  |





### services.settings.UpdateRolePermsResponse





### services.settings.ViewAuditLogRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_ids` | [int32](#int32) | repeated | Search params |
| `from` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `to` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `services` | [string](#string) | repeated |  |
| `methods` | [string](#string) | repeated |  |
| `search` | [string](#string) | optional |  |
| `actions` | [resources.audit.EventAction](#resourcesauditEventAction) | repeated |  |
| `results` | [resources.audit.EventResult](#resourcesauditEventResult) | repeated |  |





### services.settings.ViewAuditLogResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `logs` | [resources.audit.AuditEntry](#resourcesauditAuditEntry) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.SettingsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetJobProps` | [GetJobPropsRequest](#servicessettingsGetJobPropsRequest) | [GetJobPropsResponse](#servicessettingsGetJobPropsResponse) | |
| `SetJobProps` | [SetJobPropsRequest](#servicessettingsSetJobPropsRequest) | [SetJobPropsResponse](#servicessettingsSetJobPropsResponse) | |
| `GetRoles` | [GetRolesRequest](#servicessettingsGetRolesRequest) | [GetRolesResponse](#servicessettingsGetRolesResponse) | |
| `GetRole` | [GetRoleRequest](#servicessettingsGetRoleRequest) | [GetRoleResponse](#servicessettingsGetRoleResponse) | |
| `CreateRole` | [CreateRoleRequest](#servicessettingsCreateRoleRequest) | [CreateRoleResponse](#servicessettingsCreateRoleResponse) | |
| `DeleteRole` | [DeleteRoleRequest](#servicessettingsDeleteRoleRequest) | [DeleteRoleResponse](#servicessettingsDeleteRoleResponse) | |
| `UpdateRolePerms` | [UpdateRolePermsRequest](#servicessettingsUpdateRolePermsRequest) | [UpdateRolePermsResponse](#servicessettingsUpdateRolePermsResponse) | |
| `GetPermissions` | [GetPermissionsRequest](#servicessettingsGetPermissionsRequest) | [GetPermissionsResponse](#servicessettingsGetPermissionsResponse) | |
| `GetEffectivePermissions` | [GetEffectivePermissionsRequest](#servicessettingsGetEffectivePermissionsRequest) | [GetEffectivePermissionsResponse](#servicessettingsGetEffectivePermissionsResponse) | |
| `ViewAuditLog` | [ViewAuditLogRequest](#servicessettingsViewAuditLogRequest) | [ViewAuditLogResponse](#servicessettingsViewAuditLogResponse) | |
| `ListDiscordChannels` | [ListDiscordChannelsRequest](#servicessettingsListDiscordChannelsRequest) | [ListDiscordChannelsResponse](#servicessettingsListDiscordChannelsResponse) | |
| `ListUserGuilds` | [ListUserGuildsRequest](#servicessettingsListUserGuildsRequest) | [ListUserGuildsResponse](#servicessettingsListUserGuildsResponse) | |
| `UploadJobLogo` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `DeleteJobLogo` | [DeleteJobLogoRequest](#servicessettingsDeleteJobLogoRequest) | [DeleteJobLogoResponse](#servicessettingsDeleteJobLogoResponse) | |

 <!-- end services -->



## services/settings/system.proto


### services.settings.DeleteFactionRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### services.settings.DeleteFactionResponse





### services.settings.GetAllPermissionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### services.settings.GetAllPermissionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |





### services.settings.GetJobLimitsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### services.settings.GetJobLimitsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `permissions` | [resources.permissions.permissions.Permission](#resourcespermissionspermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.attributes.RoleAttribute](#resourcespermissionsattributesRoleAttribute) | repeated |  |





### services.settings.GetStatusRequest





### services.settings.GetStatusResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [resources.settings.SystemStatus](#resourcessettingsSystemStatus) |  |  |





### services.settings.UpdateJobLimitsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `perms` | [resources.settings.PermsUpdate](#resourcessettingsPermsUpdate) | optional |  |
| `attrs` | [resources.settings.AttrsUpdate](#resourcessettingsAttrsUpdate) | optional |  |





### services.settings.UpdateJobLimitsResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.SystemService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetStatus` | [GetStatusRequest](#servicessettingsGetStatusRequest) | [GetStatusResponse](#servicessettingsGetStatusResponse) | |
| `GetAllPermissions` | [GetAllPermissionsRequest](#servicessettingsGetAllPermissionsRequest) | [GetAllPermissionsResponse](#servicessettingsGetAllPermissionsResponse) | |
| `GetJobLimits` | [GetJobLimitsRequest](#servicessettingsGetJobLimitsRequest) | [GetJobLimitsResponse](#servicessettingsGetJobLimitsResponse) | |
| `UpdateJobLimits` | [UpdateJobLimitsRequest](#servicessettingsUpdateJobLimitsRequest) | [UpdateJobLimitsResponse](#servicessettingsUpdateJobLimitsResponse) | |
| `DeleteFaction` | [DeleteFactionRequest](#servicessettingsDeleteFactionRequest) | [DeleteFactionResponse](#servicessettingsDeleteFactionResponse) | |

 <!-- end services -->



## services/stats/stats.proto


### services.stats.GetStatsRequest





### services.stats.GetStatsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stats` | [GetStatsResponse.StatsEntry](#servicesstatsGetStatsResponseStatsEntry) | repeated |  |





### services.stats.GetStatsResponse.StatsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [resources.stats.Stat](#resourcesstatsStat) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.stats.StatsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetStats` | [GetStatsRequest](#servicesstatsGetStatsRequest) | [GetStatsResponse](#servicesstatsGetStatsResponse) | |

 <!-- end services -->



## services/sync/sync.proto


### services.sync.AddActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_oauth2` | [resources.sync.activity.UserOAuth2Conn](#resourcessyncactivityUserOAuth2Conn) |  |  |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |
| `user_activity` | [resources.users.activity.UserActivity](#resourcesusersactivityUserActivity) |  | User activity |
| `user_props` | [resources.sync.activity.UserProps](#resourcessyncactivityUserProps) |  | Setting props will cause activity to be created automtically |
| `colleague_activity` | [resources.jobs.colleagues.activity.ColleagueActivity](#resourcesjobscolleaguesactivityColleagueActivity) |  | Jobs user activity |
| `colleague_props` | [resources.sync.activity.ColleagueProps](#resourcessyncactivityColleagueProps) |  | Setting props will cause activity to be created automtically |
| `job_timeclock` | [resources.sync.activity.TimeclockUpdate](#resourcessyncactivityTimeclockUpdate) |  | Timeclock user entry |
| `account_update` | [resources.sync.activity.AccountUpdate](#resourcessyncactivityAccountUpdate) |  | Account update for a signle account to update the group(s). |
| `user_update` | [resources.sync.activity.UserUpdate](#resourcessyncactivityUserUpdate) |  | User/Char info updates that aren't tracked by activity (yet) |





### services.sync.AddActivityResponse





### services.sync.DeleteDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.sync.data.DeleteUsers](#resourcessyncdataDeleteUsers) |  |  |
| `vehicles` | [resources.sync.data.DeleteVehicles](#resourcessyncdataDeleteVehicles) |  |  |





### services.sync.DeleteDataResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |





### services.sync.GetStatusRequest





### services.sync.GetStatusResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `licenses` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `accounts` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `users` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `vehicles` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |





### services.sync.RegisterAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `reset_token` | [bool](#bool) |  |  |
| `last_char_id` | [int32](#int32) | optional |  |





### services.sync.RegisterAccountResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) | optional |  |
| `account_id` | [int64](#int64) | optional |  |
| `username` | [string](#string) | optional |  |





### services.sync.SendDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.data.DataJobs](#resourcessyncdataDataJobs) |  |  |
| `licenses` | [resources.sync.data.DataLicenses](#resourcessyncdataDataLicenses) |  |  |
| `accounts` | [resources.sync.data.DataAccounts](#resourcessyncdataDataAccounts) |  |  |
| `users` | [resources.sync.data.DataUsers](#resourcessyncdataDataUsers) |  |  |
| `vehicles` | [resources.sync.data.DataVehicles](#resourcessyncdataDataVehicles) |  |  |
| `user_locations` | [resources.sync.data.DataUserLocations](#resourcessyncdataDataUserLocations) |  |  |
| `last_char_id` | [resources.sync.data.LastCharID](#resourcessyncdataLastCharID) |  |  |





### services.sync.SendDataResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |





### services.sync.StreamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) | optional |  |





### services.sync.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |





### services.sync.TransferAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old_license` | [string](#string) |  |  |
| `new_license` | [string](#string) |  |  |





### services.sync.TransferAccountResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.sync.SyncService
Sync Service handles the sync of data (e.g., users, jobs) to this FiveNet instance and API calls from the plugin (e.g., user activity, user props changes).


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetStatus` | [GetStatusRequest](#servicessyncGetStatusRequest) | [GetStatusResponse](#servicessyncGetStatusResponse) |Get basic "sync state" from server side (currently simply the count of records on the server side). |
| `AddActivity` | [AddActivityRequest](#servicessyncAddActivityRequest) | [AddActivityResponse](#servicessyncAddActivityResponse) |For "tracking" activity such as "user received traffic infraction points", timeclock entries, etc. |
| `RegisterAccount` | [RegisterAccountRequest](#servicessyncRegisterAccountRequest) | [RegisterAccountResponse](#servicessyncRegisterAccountResponse) |Get registration token for a new user account or return the account id and username, for a given identifier/license. |
| `TransferAccount` | [TransferAccountRequest](#servicessyncTransferAccountRequest) | [TransferAccountResponse](#servicessyncTransferAccountResponse) |Transfer account from one license to another |
| `SendData` | [SendDataRequest](#servicessyncSendDataRequest) | [SendDataResponse](#servicessyncSendDataResponse) |DBSync's method of sending (mass) data to the FiveNet server for storing. |
| `DeleteData` | [DeleteDataRequest](#servicessyncDeleteDataRequest) | [DeleteDataResponse](#servicessyncDeleteDataResponse) |Way for the gameserver to delete certain data as well |
| `Stream` | [StreamRequest](#servicessyncStreamRequest) | [StreamResponse](#servicessyncStreamResponse) stream |Used for the server to stream events to the dbsync (e.g., "refresh" of user/char data) |

 <!-- end services -->



## services/sync/v2/sync.proto


### services.sync.v2.AddAccountUpdateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_update` | [resources.sync.activity.AccountUpdate](#resourcessyncactivityAccountUpdate) |  |  |





### services.sync.v2.AddActivityResponse





### services.sync.v2.AddColleagueActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague_activity` | [resources.jobs.colleagues.activity.ColleagueActivity](#resourcesjobscolleaguesactivityColleagueActivity) |  |  |





### services.sync.v2.AddColleaguePropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague_props` | [resources.sync.activity.ColleagueProps](#resourcessyncactivityColleagueProps) |  |  |





### services.sync.v2.AddDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.dispatches.Dispatch](#resourcescentrumdispatchesDispatch) |  |  |





### services.sync.v2.AddJobTimeclockRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_timeclock` | [resources.sync.activity.TimeclockUpdate](#resourcessyncactivityTimeclockUpdate) |  |  |





### services.sync.v2.AddUserActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_activity` | [resources.users.activity.UserActivity](#resourcesusersactivityUserActivity) |  |  |





### services.sync.v2.AddUserOAuth2ConnRequest
Individual AddActivity request messages



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_oauth2` | [resources.sync.activity.UserOAuth2Conn](#resourcessyncactivityUserOAuth2Conn) |  |  |





### services.sync.v2.AddUserPropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_props` | [resources.sync.activity.UserProps](#resourcessyncactivityUserProps) |  |  |





### services.sync.v2.AddUserUpdateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_update` | [resources.sync.activity.UserUpdate](#resourcessyncactivityUserUpdate) |  |  |





### services.sync.v2.DeleteDataResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |





### services.sync.v2.DeleteUsersDataRequest
Individual DeleteData request messages



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.sync.data.DeleteUsers](#resourcessyncdataDeleteUsers) |  |  |





### services.sync.v2.DeleteVehiclesDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vehicles` | [resources.sync.data.DeleteVehicles](#resourcessyncdataDeleteVehicles) |  |  |





### services.sync.v2.GetStatusRequest





### services.sync.v2.GetStatusResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `licenses` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `accounts` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `users` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |
| `vehicles` | [resources.sync.data.DataStatus](#resourcessyncdataDataStatus) |  |  |





### services.sync.v2.RegisterAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `reset_token` | [bool](#bool) |  |  |
| `last_char_id` | [int32](#int32) | optional |  |





### services.sync.v2.RegisterAccountResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) | optional |  |
| `account_id` | [int64](#int64) | optional |  |
| `username` | [string](#string) | optional |  |





### services.sync.v2.SendAccountsDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [resources.sync.data.DataAccounts](#resourcessyncdataDataAccounts) |  |  |





### services.sync.v2.SendDataResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |





### services.sync.v2.SendJobsDataRequest
Individual SendData request messages



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.data.DataJobs](#resourcessyncdataDataJobs) |  |  |





### services.sync.v2.SendLastCharIDDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_char_id` | [resources.sync.data.LastCharID](#resourcessyncdataLastCharID) |  |  |





### services.sync.v2.SendLicensesDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `licenses` | [resources.sync.data.DataLicenses](#resourcessyncdataDataLicenses) |  |  |





### services.sync.v2.SendUserLocationsDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_locations` | [resources.sync.data.DataUserLocations](#resourcessyncdataDataUserLocations) |  |  |





### services.sync.v2.SendUsersDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.sync.data.DataUsers](#resourcessyncdataDataUsers) |  |  |





### services.sync.v2.SendVehiclesDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vehicles` | [resources.sync.data.DataVehicles](#resourcessyncdataDataVehicles) |  |  |





### services.sync.v2.StreamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) | optional |  |





### services.sync.v2.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |





### services.sync.v2.TransferAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old_license` | [string](#string) |  |  |
| `new_license` | [string](#string) |  |  |





### services.sync.v2.TransferAccountResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.sync.v2.SyncService
Sync Service handles the sync of data (e.g., users, jobs) to this FiveNet instance and API calls from the plugin (e.g., user activity, user props changes).


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetStatus` | [GetStatusRequest](#servicessyncv2GetStatusRequest) | [GetStatusResponse](#servicessyncv2GetStatusResponse) |Get basic "sync state" from server side (currently simply the count of records on the server side). |
| `RegisterAccount` | [RegisterAccountRequest](#servicessyncv2RegisterAccountRequest) | [RegisterAccountResponse](#servicessyncv2RegisterAccountResponse) |Get registration token for a new user account or return the account id and username, for a given identifier/license. |
| `TransferAccount` | [TransferAccountRequest](#servicessyncv2TransferAccountRequest) | [TransferAccountResponse](#servicessyncv2TransferAccountResponse) |Transfer account from one license to another |
| `AddUserOAuth2Conn` | [AddUserOAuth2ConnRequest](#servicessyncv2AddUserOAuth2ConnRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) |Individual AddActivity methods |
| `AddDispatch` | [AddDispatchRequest](#servicessyncv2AddDispatchRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddUserActivity` | [AddUserActivityRequest](#servicessyncv2AddUserActivityRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddUserProps` | [AddUserPropsRequest](#servicessyncv2AddUserPropsRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddColleagueActivity` | [AddColleagueActivityRequest](#servicessyncv2AddColleagueActivityRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddColleagueProps` | [AddColleaguePropsRequest](#servicessyncv2AddColleaguePropsRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddJobTimeclock` | [AddJobTimeclockRequest](#servicessyncv2AddJobTimeclockRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddAccountUpdate` | [AddAccountUpdateRequest](#servicessyncv2AddAccountUpdateRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `AddUserUpdate` | [AddUserUpdateRequest](#servicessyncv2AddUserUpdateRequest) | [AddActivityResponse](#servicessyncv2AddActivityResponse) | |
| `SendJobsData` | [SendJobsDataRequest](#servicessyncv2SendJobsDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) |Individual SendData methods |
| `SendLicensesData` | [SendLicensesDataRequest](#servicessyncv2SendLicensesDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) | |
| `SendAccountsData` | [SendAccountsDataRequest](#servicessyncv2SendAccountsDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) | |
| `SendUsersData` | [SendUsersDataRequest](#servicessyncv2SendUsersDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) | |
| `SendVehiclesData` | [SendVehiclesDataRequest](#servicessyncv2SendVehiclesDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) | |
| `SendUserLocationsData` | [SendUserLocationsDataRequest](#servicessyncv2SendUserLocationsDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) | |
| `SendLastCharIDData` | [SendLastCharIDDataRequest](#servicessyncv2SendLastCharIDDataRequest) | [SendDataResponse](#servicessyncv2SendDataResponse) | |
| `DeleteUsersData` | [DeleteUsersDataRequest](#servicessyncv2DeleteUsersDataRequest) | [DeleteDataResponse](#servicessyncv2DeleteDataResponse) |Individual DeleteData methods |
| `DeleteVehiclesData` | [DeleteVehiclesDataRequest](#servicessyncv2DeleteVehiclesDataRequest) | [DeleteDataResponse](#servicessyncv2DeleteDataResponse) | |
| `Stream` | [StreamRequest](#servicessyncv2StreamRequest) | [StreamResponse](#servicessyncv2StreamResponse) stream |Used for the server to stream events to the dbsync (e.g., "refresh" of user/char data) |

 <!-- end services -->



## services/vehicles/vehicles.proto


### services.vehicles.ListVehiclesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `license_plate` | [string](#string) | optional | Search params |
| `model` | [string](#string) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `job` | [string](#string) | optional |  |
| `wanted` | [bool](#bool) | optional |  |





### services.vehicles.ListVehiclesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `vehicles` | [resources.vehicles.Vehicle](#resourcesvehiclesVehicle) | repeated |  |





### services.vehicles.SetVehiclePropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.vehicles.props.VehicleProps](#resourcesvehiclespropsVehicleProps) |  |  |





### services.vehicles.SetVehiclePropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.vehicles.props.VehicleProps](#resourcesvehiclespropsVehicleProps) |  |  |
| `reason` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.vehicles.VehiclesService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListVehicles` | [ListVehiclesRequest](#servicesvehiclesListVehiclesRequest) | [ListVehiclesResponse](#servicesvehiclesListVehiclesResponse) | |
| `SetVehicleProps` | [SetVehiclePropsRequest](#servicesvehiclesSetVehiclePropsRequest) | [SetVehiclePropsResponse](#servicesvehiclesSetVehiclePropsResponse) | |

 <!-- end services -->



## services/wiki/collab.proto

 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.wiki.CollabService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `JoinRoom` | [.resources.collab.ClientPacket](#resourcescollabClientPacket) stream | [.resources.collab.ServerPacket](#resourcescollabServerPacket) stream |buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |

 <!-- end services -->



## services/wiki/wiki.proto


### services.wiki.CreatePageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [int64](#int64) | optional |  |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |





### services.wiki.CreatePageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `id` | [int64](#int64) |  |  |





### services.wiki.DeletePageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.wiki.DeletePageResponse





### services.wiki.GetPageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int64](#int64) |  |  |





### services.wiki.GetPageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resourceswikiPage) |  |  |





### services.wiki.ListPageActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `page_id` | [int64](#int64) |  |  |





### services.wiki.ListPageActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.wiki.activity.PageActivity](#resourceswikiactivityPageActivity) | repeated |  |





### services.wiki.ListPagesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `job` | [string](#string) | optional | Search params |
| `root_only` | [bool](#bool) | optional |  |
| `search` | [string](#string) | optional |  |





### services.wiki.ListPagesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `pages` | [resources.wiki.PageShort](#resourceswikiPageShort) | repeated |  |





### services.wiki.UpdatePageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resourceswikiPage) |  |  |





### services.wiki.UpdatePageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resourceswikiPage) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.wiki.WikiService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListPages` | [ListPagesRequest](#serviceswikiListPagesRequest) | [ListPagesResponse](#serviceswikiListPagesResponse) | |
| `GetPage` | [GetPageRequest](#serviceswikiGetPageRequest) | [GetPageResponse](#serviceswikiGetPageResponse) | |
| `CreatePage` | [CreatePageRequest](#serviceswikiCreatePageRequest) | [CreatePageResponse](#serviceswikiCreatePageResponse) | |
| `UpdatePage` | [UpdatePageRequest](#serviceswikiUpdatePageRequest) | [UpdatePageResponse](#serviceswikiUpdatePageResponse) | |
| `DeletePage` | [DeletePageRequest](#serviceswikiDeletePageRequest) | [DeletePageResponse](#serviceswikiDeletePageResponse) | |
| `ListPageActivity` | [ListPageActivityRequest](#serviceswikiListPageActivityRequest) | [ListPageActivityResponse](#serviceswikiListPageActivityResponse) | |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) | |

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

