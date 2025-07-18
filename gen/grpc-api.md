---
title: GRPC Protobuf Documentation
description: Documentation for GRPC Protobuf files.
---



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



## resources/accounts/oauth2.proto


### resources.accounts.OAuth2Account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `provider_name` | [string](#string) |  |  |
| `provider` | [OAuth2Provider](#resourcesaccountsOAuth2Provider) |  |  |
| `external_id` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `avatar` | [string](#string) |  |  |





### resources.accounts.OAuth2Provider


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



## resources/users/licenses.proto


### resources.users.CitizensLicenses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `licenses` | [License](#resourcesusersLicense) | repeated |  |





### resources.users.License


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
@dbscanner: json



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
| `parent_id` | [uint64](#uint64) | optional |  |
| `id` | [uint64](#uint64) |  |  |
| `file_path` | [string](#string) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `byte_size` | [int64](#int64) |  | Bytes stored |
| `content_type` | [string](#string) |  |  |
| `meta` | [FileMeta](#resourcesfileFileMeta) | optional |  |




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



## resources/users/labels.proto


### resources.users.Label


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `color` | [string](#string) |  | @sanitize: method=StripTags |





### resources.users.Labels
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [Label](#resourcesusersLabel) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/props.proto


### resources.users.UserProps


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
| `mugshot_file_id` | [uint64](#uint64) | optional |  |
| `mugshot` | [resources.file.File](#resourcesfileFile) | optional |  |
| `labels` | [Labels](#resourcesusersLabels) | optional |  |
| `email` | [string](#string) | optional | @sanitize: method=StripTags |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/users.proto


### resources.users.User


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
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
| `props` | [UserProps](#resourcesusersUserProps) |  |  |
| `licenses` | [License](#resourcesusersLicense) | repeated |  |
| `avatar_file_id` | [uint64](#uint64) | optional |  |
| `avatar` | [string](#string) | optional |  |
| `group` | [string](#string) | optional |  |





### resources.users.UserShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
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



## resources/accounts/accounts.proto


### resources.accounts.Account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `username` | [string](#string) |  |  |
| `license` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |
| `last_char` | [int32](#int32) | optional |  |
| `oauth2_accounts` | [OAuth2Account](#resourcesaccountsOAuth2Account) | repeated |  |





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



## resources/audit/audit.proto


### resources.audit.AuditEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `user_job` | [string](#string) |  |  |
| `target_user_id` | [int32](#int32) | optional |  |
| `target_user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `target_user_job` | [string](#string) | optional |  |
| `service` | [string](#string) |  |  |
| `method` | [string](#string) |  |  |
| `state` | [EventType](#resourcesauditEventType) |  |  |
| `data` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.audit.EventType

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



## resources/calendar/access.proto


### resources.calendar.CalendarAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [CalendarJobAccess](#resourcescalendarCalendarJobAccess) | repeated |  |
| `users` | [CalendarUserAccess](#resourcescalendarCalendarUserAccess) | repeated |  |





### resources.calendar.CalendarJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcescalendarAccessLevel) |  |  |





### resources.calendar.CalendarUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `access` | [AccessLevel](#resourcescalendarAccessLevel) |  |  |




 <!-- end messages -->


### resources.calendar.AccessLevel

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



## resources/common/content/content.proto


### resources.common.content.Content


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) | optional |  |
| `content` | [JSONNode](#resourcescommoncontentJSONNode) | optional |  |
| `raw_content` | [string](#string) | optional | @sanitize |





### resources.common.content.JSONNode


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [NodeType](#resourcescommoncontentNodeType) |  |  |
| `id` | [string](#string) | optional | @sanitize: method=StripTags |
| `tag` | [string](#string) |  | @sanitize: method=StripTags |
| `attrs` | [JSONNode.AttrsEntry](#resourcescommoncontentJSONNodeAttrsEntry) | repeated | @sanitize: method=StripTags |
| `text` | [string](#string) | optional | @sanitize: method=StripTags |
| `content` | [JSONNode](#resourcescommoncontentJSONNode) | repeated |  |





### resources.common.content.JSONNode.AttrsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |




 <!-- end messages -->


### resources.common.content.ContentType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CONTENT_TYPE_UNSPECIFIED` | 0 |  |
| `CONTENT_TYPE_HTML` | 1 |  |
| `CONTENT_TYPE_PLAIN` | 2 |  |
| `CONTENT_TYPE_TIPTAP_JSON` | 3 |  |



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



## resources/calendar/calendar.proto


### resources.calendar.Calendar


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) | optional |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `public` | [bool](#bool) |  |  |
| `closed` | [bool](#bool) |  |  |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `subscription` | [CalendarSub](#resourcescalendarCalendarSub) | optional |  |
| `access` | [CalendarAccess](#resourcescalendarCalendarAccess) |  |  |





### resources.calendar.CalendarEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `calendar_id` | [uint64](#uint64) |  |  |
| `calendar` | [Calendar](#resourcescalendarCalendar) | optional |  |
| `job` | [string](#string) | optional |  |
| `start_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `end_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `closed` | [bool](#bool) |  |  |
| `rsvp_open` | [bool](#bool) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `recurring` | [CalendarEntryRecurring](#resourcescalendarCalendarEntryRecurring) | optional |  |
| `rsvp` | [CalendarEntryRSVP](#resourcescalendarCalendarEntryRSVP) | optional |  |





### resources.calendar.CalendarEntryRSVP


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `response` | [RsvpResponses](#resourcescalendarRsvpResponses) |  |  |





### resources.calendar.CalendarEntryRecurring
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `every` | [string](#string) |  |  |
| `count` | [int32](#int32) |  |  |
| `until` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.calendar.CalendarShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `public` | [bool](#bool) |  |  |
| `closed` | [bool](#bool) |  |  |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `subscription` | [CalendarSub](#resourcescalendarCalendarSub) | optional |  |





### resources.calendar.CalendarSub


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `confirmed` | [bool](#bool) |  |  |
| `muted` | [bool](#bool) |  |  |




 <!-- end messages -->


### resources.calendar.RsvpResponses

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



## resources/centrum/access.proto


### resources.centrum.CentrumAccess
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [CentrumJobAccess](#resourcescentrumCentrumJobAccess) | repeated |  |





### resources.centrum.CentrumJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [CentrumAccessLevel](#resourcescentrumCentrumAccessLevel) |  |  |





### resources.centrum.CentrumQualificationAccess
Dummy - DO NOT USE!






### resources.centrum.CentrumUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


### resources.centrum.CentrumAccessLevel

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



## resources/centrum/attributes.proto


### resources.centrum.DispatchAttributes
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [DispatchAttribute](#resourcescentrumDispatchAttribute) | repeated |  |





### resources.centrum.UnitAttributes
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [UnitAttribute](#resourcescentrumUnitAttribute) | repeated |  |




 <!-- end messages -->


### resources.centrum.DispatchAttribute

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISPATCH_ATTRIBUTE_UNSPECIFIED` | 0 |  |
| `DISPATCH_ATTRIBUTE_MULTIPLE` | 1 |  |
| `DISPATCH_ATTRIBUTE_DUPLICATE` | 2 |  |
| `DISPATCH_ATTRIBUTE_TOO_OLD` | 3 |  |
| `DISPATCH_ATTRIBUTE_AUTOMATIC` | 4 |  |



### resources.centrum.UnitAttribute

| Name | Number | Description |
| ---- | ------ | ----------- |
| `UNIT_ATTRIBUTE_UNSPECIFIED` | 0 |  |
| `UNIT_ATTRIBUTE_STATIC` | 1 |  |
| `UNIT_ATTRIBUTE_NO_DISPATCH_AUTO_ASSIGN` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/labels.proto


### resources.jobs.Label


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `order` | [int32](#int32) |  |  |





### resources.jobs.LabelCount


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `label` | [Label](#resourcesjobsLabel) |  |  |
| `count` | [int64](#int64) |  |  |





### resources.jobs.Labels


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `list` | [Label](#resourcesjobsLabel) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/colleagues.proto


### resources.jobs.Colleague


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
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
| `props` | [ColleagueProps](#resourcesjobsColleagueProps) |  |  |
| `email` | [string](#string) | optional | @sanitize: method=StripTags |





### resources.jobs.ColleagueProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `absence_begin` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `absence_end` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `note` | [string](#string) | optional | @sanitize: method=StripTags |
| `labels` | [Labels](#resourcesjobsLabels) | optional |  |
| `name_prefix` | [string](#string) | optional |  |
| `name_suffix` | [string](#string) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/dispatchers.proto


### resources.centrum.Dispatchers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `dispatchers` | [resources.jobs.Colleague](#resourcesjobsColleague) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/settings.proto


### resources.centrum.Configuration
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deduplication_enabled` | [bool](#bool) |  |  |
| `deduplication_radius` | [int64](#int64) |  |  |
| `deduplication_duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) | optional |  |





### resources.centrum.Job


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `label` | [string](#string) | optional |  |





### resources.centrum.JobList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [Job](#resourcescentrumJob) | repeated |  |





### resources.centrum.PredefinedStatus
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_status` | [string](#string) | repeated | @sanitize: method=StripTags |
| `dispatch_status` | [string](#string) | repeated | @sanitize: method=StripTags |





### resources.centrum.Settings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |
| `type` | [CentrumType](#resourcescentrumCentrumType) |  |  |
| `public` | [bool](#bool) |  |  |
| `mode` | [CentrumMode](#resourcescentrumCentrumMode) |  |  |
| `fallback_mode` | [CentrumMode](#resourcescentrumCentrumMode) |  |  |
| `predefined_status` | [PredefinedStatus](#resourcescentrumPredefinedStatus) | optional |  |
| `timings` | [Timings](#resourcescentrumTimings) |  |  |
| `access` | [CentrumAccess](#resourcescentrumCentrumAccess) |  |  |
| `configuration` | [Configuration](#resourcescentrumConfiguration) |  |  |





### resources.centrum.Timings
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_max_wait` | [int64](#int64) |  |  |
| `require_unit` | [bool](#bool) |  |  |
| `require_unit_reminder_seconds` | [int64](#int64) |  |  |




 <!-- end messages -->


### resources.centrum.CentrumMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_MODE_UNSPECIFIED` | 0 |  |
| `CENTRUM_MODE_MANUAL` | 1 |  |
| `CENTRUM_MODE_CENTRAL_COMMAND` | 2 |  |
| `CENTRUM_MODE_AUTO_ROUND_ROBIN` | 3 |  |
| `CENTRUM_MODE_SIMPLIFIED` | 4 |  |



### resources.centrum.CentrumType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `CENTRUM_TYPE_UNSPECIFIED` | 0 |  |
| `CENTRUM_TYPE_DISPATCH` | 1 |  |
| `CENTRUM_TYPE_DELIVERY` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/qualifications/access.proto


### resources.qualifications.QualificationAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [QualificationJobAccess](#resourcesqualificationsQualificationJobAccess) | repeated |  |





### resources.qualifications.QualificationJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesqualificationsAccessLevel) |  |  |





### resources.qualifications.QualificationUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->


### resources.qualifications.AccessLevel

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



## resources/qualifications/exam.proto


### resources.qualifications.ExamGrading
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `responses` | [ExamGradingResponse](#resourcesqualificationsExamGradingResponse) | repeated |  |





### resources.qualifications.ExamGradingResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `question_id` | [uint64](#uint64) |  |  |
| `points` | [float](#float) |  |  |
| `checked` | [bool](#bool) | optional |  |





### resources.qualifications.ExamQuestion


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `data` | [ExamQuestionData](#resourcesqualificationsExamQuestionData) |  |  |
| `answer` | [ExamQuestionAnswerData](#resourcesqualificationsExamQuestionAnswerData) | optional |  |
| `points` | [int32](#int32) | optional |  |
| `order` | [int32](#int32) |  |  |





### resources.qualifications.ExamQuestionAnswerData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `answer_key` | [string](#string) |  |  |
| `yesno` | [ExamResponseYesNo](#resourcesqualificationsExamResponseYesNo) |  |  |
| `free_text` | [ExamResponseText](#resourcesqualificationsExamResponseText) |  |  |
| `single_choice` | [ExamResponseSingleChoice](#resourcesqualificationsExamResponseSingleChoice) |  |  |
| `multiple_choice` | [ExamResponseMultipleChoice](#resourcesqualificationsExamResponseMultipleChoice) |  |  |





### resources.qualifications.ExamQuestionData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `separator` | [ExamQuestionSeparator](#resourcesqualificationsExamQuestionSeparator) |  |  |
| `image` | [ExamQuestionImage](#resourcesqualificationsExamQuestionImage) |  |  |
| `yesno` | [ExamQuestionYesNo](#resourcesqualificationsExamQuestionYesNo) |  |  |
| `free_text` | [ExamQuestionText](#resourcesqualificationsExamQuestionText) |  |  |
| `single_choice` | [ExamQuestionSingleChoice](#resourcesqualificationsExamQuestionSingleChoice) |  |  |
| `multiple_choice` | [ExamQuestionMultipleChoice](#resourcesqualificationsExamQuestionMultipleChoice) |  |  |





### resources.qualifications.ExamQuestionImage


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image` | [resources.file.File](#resourcesfileFile) |  |  |
| `alt` | [string](#string) | optional |  |





### resources.qualifications.ExamQuestionMultipleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated | @sanitize: method=StripTags |
| `limit` | [int32](#int32) | optional |  |





### resources.qualifications.ExamQuestionSeparator





### resources.qualifications.ExamQuestionSingleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated | @sanitize: method=StripTags |





### resources.qualifications.ExamQuestionText


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min_length` | [int32](#int32) |  |  |
| `max_length` | [int32](#int32) |  |  |





### resources.qualifications.ExamQuestionYesNo





### resources.qualifications.ExamQuestions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `questions` | [ExamQuestion](#resourcesqualificationsExamQuestion) | repeated |  |





### resources.qualifications.ExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `question_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `question` | [ExamQuestion](#resourcesqualificationsExamQuestion) |  |  |
| `response` | [ExamResponseData](#resourcesqualificationsExamResponseData) |  |  |





### resources.qualifications.ExamResponseData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `separator` | [ExamResponseSeparator](#resourcesqualificationsExamResponseSeparator) |  |  |
| `yesno` | [ExamResponseYesNo](#resourcesqualificationsExamResponseYesNo) |  |  |
| `free_text` | [ExamResponseText](#resourcesqualificationsExamResponseText) |  |  |
| `single_choice` | [ExamResponseSingleChoice](#resourcesqualificationsExamResponseSingleChoice) |  |  |
| `multiple_choice` | [ExamResponseMultipleChoice](#resourcesqualificationsExamResponseMultipleChoice) |  |  |





### resources.qualifications.ExamResponseMultipleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choices` | [string](#string) | repeated | @sanitize: method=StripTags |





### resources.qualifications.ExamResponseSeparator





### resources.qualifications.ExamResponseSingleChoice


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `choice` | [string](#string) |  | @sanitize: method=StripTags |





### resources.qualifications.ExamResponseText


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `text` | [string](#string) |  | @sanitize: method=StripTags<br/><br/>0.5 Megabyte |





### resources.qualifications.ExamResponseYesNo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [bool](#bool) |  |  |





### resources.qualifications.ExamResponses
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `responses` | [ExamResponse](#resourcesqualificationsExamResponse) | repeated |  |





### resources.qualifications.ExamUser


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `started_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `ends_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `ended_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/qualifications/qualifications.proto


### resources.qualifications.Qualification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `abbreviation` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `access` | [QualificationAccess](#resourcesqualificationsQualificationAccess) |  |  |
| `requirements` | [QualificationRequirement](#resourcesqualificationsQualificationRequirement) | repeated |  |
| `discord_sync_enabled` | [bool](#bool) |  |  |
| `discord_settings` | [QualificationDiscordSettings](#resourcesqualificationsQualificationDiscordSettings) | optional |  |
| `exam_mode` | [QualificationExamMode](#resourcesqualificationsQualificationExamMode) |  |  |
| `exam_settings` | [QualificationExamSettings](#resourcesqualificationsQualificationExamSettings) | optional |  |
| `exam` | [ExamQuestions](#resourcesqualificationsExamQuestions) | optional |  |
| `result` | [QualificationResult](#resourcesqualificationsQualificationResult) | optional |  |
| `request` | [QualificationRequest](#resourcesqualificationsQualificationRequest) | optional |  |
| `label_sync_enabled` | [bool](#bool) |  |  |
| `label_sync_format` | [string](#string) | optional | @sanitize: method=StripTags |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### resources.qualifications.QualificationDiscordSettings
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_name` | [string](#string) | optional |  |
| `role_format` | [string](#string) | optional |  |





### resources.qualifications.QualificationExamSettings
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `time` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `auto_grade` | [bool](#bool) |  |  |
| `auto_grade_mode` | [AutoGradeMode](#resourcesqualificationsAutoGradeMode) |  |  |
| `minimum_points` | [int32](#int32) |  |  |





### resources.qualifications.QualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) |  |  |
| `user_comment` | [string](#string) | optional | @sanitize: method=StripTags |
| `status` | [RequestStatus](#resourcesqualificationsRequestStatus) | optional |  |
| `approved_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `approver_comment` | [string](#string) | optional | @sanitize: method=StripTags |
| `approver_id` | [int32](#int32) | optional |  |
| `approver` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `approver_job` | [string](#string) | optional |  |





### resources.qualifications.QualificationRequirement


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `target_qualification_id` | [uint64](#uint64) |  |  |
| `target_qualification` | [QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |





### resources.qualifications.QualificationResult


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) |  |  |
| `status` | [ResultStatus](#resourcesqualificationsResultStatus) |  |  |
| `score` | [float](#float) | optional |  |
| `summary` | [string](#string) |  | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) |  |  |
| `creator_job` | [string](#string) |  |  |





### resources.qualifications.QualificationShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `abbreviation` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `requirements` | [QualificationRequirement](#resourcesqualificationsQualificationRequirement) | repeated |  |
| `exam_mode` | [QualificationExamMode](#resourcesqualificationsQualificationExamMode) |  |  |
| `exam_settings` | [QualificationExamSettings](#resourcesqualificationsQualificationExamSettings) | optional |  |
| `result` | [QualificationResult](#resourcesqualificationsQualificationResult) | optional |  |




 <!-- end messages -->


### resources.qualifications.AutoGradeMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `AUTO_GRADE_MODE_UNSPECIFIED` | 0 |  |
| `AUTO_GRADE_MODE_STRICT` | 1 |  |
| `AUTO_GRADE_MODE_PARTIAL_CREDIT` | 2 |  |



### resources.qualifications.QualificationExamMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `QUALIFICATION_EXAM_MODE_UNSPECIFIED` | 0 |  |
| `QUALIFICATION_EXAM_MODE_DISABLED` | 1 |  |
| `QUALIFICATION_EXAM_MODE_REQUEST_NEEDED` | 2 |  |
| `QUALIFICATION_EXAM_MODE_ENABLED` | 3 |  |



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



## resources/centrum/units_access.proto


### resources.centrum.UnitAccess
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [UnitJobAccess](#resourcescentrumUnitJobAccess) | repeated |  |
| `qualifications` | [UnitQualificationAccess](#resourcescentrumUnitQualificationAccess) | repeated |  |





### resources.centrum.UnitJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [UnitAccessLevel](#resourcescentrumUnitAccessLevel) |  |  |





### resources.centrum.UnitQualificationAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [resources.qualifications.QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `access` | [UnitAccessLevel](#resourcescentrumUnitAccessLevel) |  |  |





### resources.centrum.UnitUserAccess




 <!-- end messages -->


### resources.centrum.UnitAccessLevel

| Name | Number | Description |
| ---- | ------ | ----------- |
| `UNIT_ACCESS_LEVEL_UNSPECIFIED` | 0 |  |
| `UNIT_ACCESS_LEVEL_BLOCKED` | 1 |  |
| `UNIT_ACCESS_LEVEL_JOIN` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/centrum/units.proto


### resources.centrum.Unit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `initials` | [string](#string) |  | @sanitize |
| `color` | [string](#string) |  | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `description` | [string](#string) | optional | @sanitize |
| `status` | [UnitStatus](#resourcescentrumUnitStatus) | optional |  |
| `users` | [UnitAssignment](#resourcescentrumUnitAssignment) | repeated |  |
| `attributes` | [UnitAttributes](#resourcescentrumUnitAttributes) | optional |  |
| `home_postal` | [string](#string) | optional |  |
| `access` | [UnitAccess](#resourcescentrumUnitAccess) |  |  |





### resources.centrum.UnitAssignment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.jobs.Colleague](#resourcesjobsColleague) | optional |  |





### resources.centrum.UnitAssignments


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `users` | [UnitAssignment](#resourcescentrumUnitAssignment) | repeated |  |





### resources.centrum.UnitStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `unit_id` | [uint64](#uint64) |  |  |
| `unit` | [Unit](#resourcescentrumUnit) | optional |  |
| `status` | [StatusUnit](#resourcescentrumStatusUnit) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.Colleague](#resourcesjobsColleague) | optional |  |
| `x` | [double](#double) | optional |  |
| `y` | [double](#double) | optional |  |
| `postal` | [string](#string) | optional | @sanitize |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.jobs.Colleague](#resourcesjobsColleague) | optional |  |
| `creator_job` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.centrum.StatusUnit

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



## resources/centrum/dispatches.proto


### resources.centrum.Dispatch


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  | **Deprecated.**  |
| `jobs` | [JobList](#resourcescentrumJobList) |  |  |
| `status` | [DispatchStatus](#resourcescentrumDispatchStatus) | optional |  |
| `message` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `attributes` | [DispatchAttributes](#resourcescentrumDispatchAttributes) | optional |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `postal` | [string](#string) | optional | @sanitize |
| `anon` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.User](#resourcesusersUser) | optional |  |
| `units` | [DispatchAssignment](#resourcescentrumDispatchAssignment) | repeated |  |
| `references` | [DispatchReferences](#resourcescentrumDispatchReferences) | optional |  |





### resources.centrum.DispatchAssignment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `unit_id` | [uint64](#uint64) |  |  |
| `unit` | [Unit](#resourcescentrumUnit) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### resources.centrum.DispatchAssignments


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `units` | [DispatchAssignment](#resourcescentrumDispatchAssignment) | repeated |  |





### resources.centrum.DispatchReference


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_dispatch_id` | [uint64](#uint64) |  |  |
| `reference_type` | [DispatchReferenceType](#resourcescentrumDispatchReferenceType) |  |  |





### resources.centrum.DispatchReferences
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `references` | [DispatchReference](#resourcescentrumDispatchReference) | repeated |  |





### resources.centrum.DispatchStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `unit_id` | [uint64](#uint64) | optional |  |
| `unit` | [Unit](#resourcescentrumUnit) | optional |  |
| `status` | [StatusDispatch](#resourcescentrumStatusDispatch) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.Colleague](#resourcesjobsColleague) | optional |  |
| `x` | [double](#double) | optional |  |
| `y` | [double](#double) | optional |  |
| `postal` | [string](#string) | optional | @sanitize |
| `creator_job` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.centrum.DispatchReferenceType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DISPATCH_REFERENCE_TYPE_UNSPECIFIED` | 0 |  |
| `DISPATCH_REFERENCE_TYPE_REFERENCED` | 1 |  |
| `DISPATCH_REFERENCE_TYPE_DUPLICATED_BY` | 2 |  |
| `DISPATCH_REFERENCE_TYPE_DUPLICATE_OF` | 3 |  |



### resources.centrum.StatusDispatch

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



### resources.centrum.TakeDispatchResp

| Name | Number | Description |
| ---- | ------ | ----------- |
| `TAKE_DISPATCH_RESP_UNSPECIFIED` | 0 |  |
| `TAKE_DISPATCH_RESP_TIMEOUT` | 1 |  |
| `TAKE_DISPATCH_RESP_ACCEPTED` | 2 |  |
| `TAKE_DISPATCH_RESP_DECLINED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/banner.proto


### resources.settings.BannerMessage


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize: method |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/clientconfig/clientconfig.proto


### resources.clientconfig.ClientConfig


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  |  |
| `default_locale` | [string](#string) |  |  |
| `login` | [LoginConfig](#resourcesclientconfigLoginConfig) |  |  |
| `discord` | [Discord](#resourcesclientconfigDiscord) |  |  |
| `website` | [Website](#resourcesclientconfigWebsite) |  |  |
| `feature_gates` | [FeatureGates](#resourcesclientconfigFeatureGates) |  |  |
| `game` | [Game](#resourcesclientconfigGame) |  |  |
| `system` | [System](#resourcesclientconfigSystem) |  |  |





### resources.clientconfig.Discord


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bot_enabled` | [bool](#bool) |  |  |





### resources.clientconfig.FeatureGates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `image_proxy` | [bool](#bool) |  |  |





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





### resources.clientconfig.LoginConfig


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signup_enabled` | [bool](#bool) |  |  |
| `last_char_lock` | [bool](#bool) |  |  |
| `providers` | [ProviderConfig](#resourcesclientconfigProviderConfig) | repeated |  |





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
| `first` | [bool](#bool) |  |  |





### resources.collab.CollabInit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `target_id` | [uint64](#uint64) |  |  |





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
| `target_id` | [uint64](#uint64) |  |  |





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



## resources/common/cron/cron.proto


### resources.common.cron.Cronjob


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Cronjob name |
| `schedule` | [string](#string) |  | Cron schedule expression For available valid expressions, see [adhocore/gronx - Cron Expressions Documentation](https://github.com/adhocore/gronx/blob/fea40e3e90e70476877cfb9b50fac10c7de41c5c/README.md#cron-expression).<br/><br/>To generate Cronjob schedule expressions, you can also use web tools like https://crontab.guru/. |
| `state` | [CronjobState](#resourcescommoncronCronjobState) |  | Cronjob state |
| `next_schedule_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Next time the cronjob should be run |
| `last_attempt_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional | Last attempted start time of Cronjob |
| `started_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional | Time current cronjob was started |
| `timeout` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) | optional | Optional timeout for cronjob execution |
| `data` | [CronjobData](#resourcescommoncronCronjobData) |  | Cronjob data |
| `last_completed_event` | [CronjobCompletedEvent](#resourcescommoncronCronjobCompletedEvent) | optional | Last event info to ease debugging and tracking |





### resources.common.cron.CronjobCompletedEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Cronjob name |
| `success` | [bool](#bool) |  | Cronjob execution success status |
| `cancelled` | [bool](#bool) |  | Cronjob execution was cancelled |
| `end_date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Cronjob end time |
| `elapsed` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  | Cronjob execution time/elapsed time |
| `data` | [CronjobData](#resourcescommoncronCronjobData) | optional | Cronjob data (can be empty if not touched by the Cronjob handler) |
| `node_name` | [string](#string) |  | Name of the node where the cronjob was executed |
| `error_message` | [string](#string) | optional | Error message (if success = false) |





### resources.common.cron.CronjobData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `data` | [google.protobuf.Any](https://protobuf.dev/reference/protobuf/google.protobuf/#any) | optional |  |





### resources.common.cron.CronjobLockOwnerState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hostname` | [string](#string) |  | Hostname of the agent the cronjob is running on |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |





### resources.common.cron.CronjobSchedulerEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cronjob` | [Cronjob](#resourcescommoncronCronjob) |  | Full Cronjob spec |





### resources.common.cron.GenericCronData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attributes` | [GenericCronData.AttributesEntry](#resourcescommoncronGenericCronDataAttributesEntry) | repeated | @sanitize: method=StripTags |





### resources.common.cron.GenericCronData.AttributesEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |




 <!-- end messages -->


### resources.common.cron.CronjobState
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



## resources/common/database/database.proto


### resources.common.database.DateRange
Datetime range (uses Timestamp underneath) It depends on the API method if it will use date or date + time.



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
Sort by column



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `column` | [string](#string) |  | Column name |
| `direction` | [string](#string) |  | Sort direction, must be `asc` (ascending) or `desc` (descending) |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/i18n.proto


### resources.common.I18NItem
Wrapped translated message for the client @dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  | @sanitize: method=StripTags |
| `parameters` | [I18NItem.ParametersEntry](#resourcescommonI18NItemParametersEntry) | repeated | @sanitize: method=StripTags |





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



## resources/common/grpcws/grpcws.proto


### resources.common.grpcws.Body


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |
| `complete` | [bool](#bool) |  |  |





### resources.common.grpcws.Cancel





### resources.common.grpcws.Complete





### resources.common.grpcws.Failure


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `error_message` | [string](#string) |  |  |
| `error_status` | [string](#string) |  |  |
| `headers` | [Failure.HeadersEntry](#resourcescommongrpcwsFailureHeadersEntry) | repeated |  |





### resources.common.grpcws.Failure.HeadersEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [HeaderValue](#resourcescommongrpcwsHeaderValue) |  |  |





### resources.common.grpcws.GrpcFrame


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stream_id` | [uint32](#uint32) |  |  |
| `ping` | [Ping](#resourcescommongrpcwsPing) |  |  |
| `header` | [Header](#resourcescommongrpcwsHeader) |  |  |
| `body` | [Body](#resourcescommongrpcwsBody) |  |  |
| `complete` | [Complete](#resourcescommongrpcwsComplete) |  |  |
| `failure` | [Failure](#resourcescommongrpcwsFailure) |  |  |
| `cancel` | [Cancel](#resourcescommongrpcwsCancel) |  |  |





### resources.common.grpcws.Header


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operation` | [string](#string) |  |  |
| `headers` | [Header.HeadersEntry](#resourcescommongrpcwsHeaderHeadersEntry) | repeated |  |
| `status` | [int32](#int32) |  |  |





### resources.common.grpcws.Header.HeadersEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [HeaderValue](#resourcescommongrpcwsHeaderValue) |  |  |





### resources.common.grpcws.HeaderValue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [string](#string) | repeated |  |





### resources.common.grpcws.Ping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pong` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/common/id_mapping.proto


### resources.common.IDMapping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |




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



## resources/common/uuid.proto


### resources.common.UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `uuid` | [string](#string) |  |  |




 <!-- end messages -->

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



## resources/documents/access.proto


### resources.documents.DocumentAccess
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [DocumentJobAccess](#resourcesdocumentsDocumentJobAccess) | repeated |  |
| `users` | [DocumentUserAccess](#resourcesdocumentsDocumentUserAccess) | repeated |  |





### resources.documents.DocumentJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesdocumentsAccessLevel) |  |  |
| `required` | [bool](#bool) | optional |  |





### resources.documents.DocumentUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `access` | [AccessLevel](#resourcesdocumentsAccessLevel) |  |  |
| `required` | [bool](#bool) | optional |  |




 <!-- end messages -->


### resources.documents.AccessLevel

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



## resources/documents/activity.proto


### resources.documents.DocAccessJobsDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [DocumentJobAccess](#resourcesdocumentsDocumentJobAccess) | repeated |  |
| `to_update` | [DocumentJobAccess](#resourcesdocumentsDocumentJobAccess) | repeated |  |
| `to_delete` | [DocumentJobAccess](#resourcesdocumentsDocumentJobAccess) | repeated |  |





### resources.documents.DocAccessRequested


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `level` | [AccessLevel](#resourcesdocumentsAccessLevel) |  |  |





### resources.documents.DocAccessUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [DocAccessJobsDiff](#resourcesdocumentsDocAccessJobsDiff) |  |  |
| `users` | [DocAccessUsersDiff](#resourcesdocumentsDocAccessUsersDiff) |  |  |





### resources.documents.DocAccessUsersDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [DocumentUserAccess](#resourcesdocumentsDocumentUserAccess) | repeated |  |
| `to_update` | [DocumentUserAccess](#resourcesdocumentsDocumentUserAccess) | repeated |  |
| `to_delete` | [DocumentUserAccess](#resourcesdocumentsDocumentUserAccess) | repeated |  |





### resources.documents.DocActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `activity_type` | [DocActivityType](#resourcesdocumentsDocActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional | @sanitize |
| `data` | [DocActivityData](#resourcesdocumentsDocActivityData) |  |  |





### resources.documents.DocActivityData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [DocUpdated](#resourcesdocumentsDocUpdated) |  |  |
| `owner_changed` | [DocOwnerChanged](#resourcesdocumentsDocOwnerChanged) |  |  |
| `access_updated` | [DocAccessUpdated](#resourcesdocumentsDocAccessUpdated) |  |  |
| `access_requested` | [DocAccessRequested](#resourcesdocumentsDocAccessRequested) |  |  |





### resources.documents.DocFilesChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [int64](#int64) |  |  |
| `deleted` | [int64](#int64) |  |  |





### resources.documents.DocOwnerChanged


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_owner_id` | [int32](#int32) |  |  |
| `new_owner` | [resources.users.UserShort](#resourcesusersUserShort) |  |  |





### resources.documents.DocUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title_diff` | [string](#string) | optional |  |
| `content_diff` | [string](#string) | optional |  |
| `state_diff` | [string](#string) | optional |  |
| `files_change` | [DocFilesChange](#resourcesdocumentsDocFilesChange) | optional |  |




 <!-- end messages -->


### resources.documents.DocActivityType

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



## resources/documents/category.proto


### resources.documents.Category


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `job` | [string](#string) | optional |  |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/comment.proto


### resources.documents.Comment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `document_id` | [uint64](#uint64) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/pins.proto


### resources.documents.DocumentPin


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `state` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/workflow.proto


### resources.documents.AutoCloseSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `message` | [string](#string) |  |  |





### resources.documents.Reminder


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |
| `message` | [string](#string) |  |  |





### resources.documents.ReminderSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reminders` | [Reminder](#resourcesdocumentsReminder) | repeated |  |





### resources.documents.Workflow
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reminder` | [bool](#bool) |  |  |
| `reminder_settings` | [ReminderSettings](#resourcesdocumentsReminderSettings) |  |  |
| `auto_close` | [bool](#bool) |  |  |
| `auto_close_settings` | [AutoCloseSettings](#resourcesdocumentsAutoCloseSettings) |  |  |





### resources.documents.WorkflowCronData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_doc_id` | [uint64](#uint64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/documents.proto


### resources.documents.Document


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category_id` | [uint64](#uint64) | optional |  |
| `category` | [Category](#resourcesdocumentsCategory) | optional |  |
| `title` | [string](#string) |  | @sanitize |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `data` | [string](#string) | optional | @sanitize |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `state` | [string](#string) |  | @sanitize |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `template_id` | [uint64](#uint64) | optional |  |
| `pin` | [DocumentPin](#resourcesdocumentsDocumentPin) | optional |  |
| `workflow_state` | [WorkflowState](#resourcesdocumentsWorkflowState) | optional |  |
| `workflow_user` | [WorkflowUserState](#resourcesdocumentsWorkflowUserState) | optional |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### resources.documents.DocumentReference


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `source_document_id` | [uint64](#uint64) |  |  |
| `source_document` | [DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |
| `reference` | [DocReference](#resourcesdocumentsDocReference) |  |  |
| `target_document_id` | [uint64](#uint64) |  |  |
| `target_document` | [DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |





### resources.documents.DocumentRelation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `document_id` | [uint64](#uint64) |  |  |
| `document` | [DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |
| `source_user_id` | [int32](#int32) |  |  |
| `source_user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `relation` | [DocRelation](#resourcesdocumentsDocRelation) |  |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |





### resources.documents.DocumentShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category_id` | [uint64](#uint64) | optional |  |
| `category` | [Category](#resourcesdocumentsCategory) | optional |  |
| `title` | [string](#string) |  | @sanitize |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `state` | [string](#string) |  | @sanitize |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `pin` | [DocumentPin](#resourcesdocumentsDocumentPin) | optional |  |
| `workflow_state` | [WorkflowState](#resourcesdocumentsWorkflowState) | optional |  |
| `workflow_user` | [WorkflowUserState](#resourcesdocumentsWorkflowUserState) | optional |  |





### resources.documents.WorkflowState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `next_reminder_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `next_reminder_count` | [int32](#int32) | optional |  |
| `auto_close_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `workflow` | [Workflow](#resourcesdocumentsWorkflow) | optional |  |
| `document` | [DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |





### resources.documents.WorkflowUserState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `manual_reminder_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `manual_reminder_message` | [string](#string) | optional |  |
| `workflow` | [Workflow](#resourcesdocumentsWorkflow) | optional |  |
| `document` | [DocumentShort](#resourcesdocumentsDocumentShort) | optional |  |




 <!-- end messages -->


### resources.documents.DocReference

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_REFERENCE_UNSPECIFIED` | 0 |  |
| `DOC_REFERENCE_LINKED` | 1 |  |
| `DOC_REFERENCE_SOLVES` | 2 |  |
| `DOC_REFERENCE_CLOSES` | 3 |  |
| `DOC_REFERENCE_DEPRECATES` | 4 |  |



### resources.documents.DocRelation

| Name | Number | Description |
| ---- | ------ | ----------- |
| `DOC_RELATION_UNSPECIFIED` | 0 |  |
| `DOC_RELATION_MENTIONED` | 1 |  |
| `DOC_RELATION_TARGETS` | 2 |  |
| `DOC_RELATION_CAUSED` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/requests.proto


### resources.documents.DocRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `request_type` | [DocActivityType](#resourcesdocumentsDocActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [DocActivityData](#resourcesdocumentsDocActivityData) |  |  |
| `accepted` | [bool](#bool) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/vehicles/props.proto


### resources.vehicles.VehicleProps


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
| `owner` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |
| `props` | [VehicleProps](#resourcesvehiclesVehicleProps) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/documents/templates.proto


### resources.documents.ObjectSpecs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `required` | [bool](#bool) | optional |  |
| `min` | [int32](#int32) | optional |  |
| `max` | [int32](#int32) | optional |  |





### resources.documents.Template


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category` | [Category](#resourcesdocumentsCategory) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) |  | @sanitize |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `content_title` | [string](#string) |  |  |
| `content` | [string](#string) |  |  |
| `state` | [string](#string) |  |  |
| `schema` | [TemplateSchema](#resourcesdocumentsTemplateSchema) |  |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `job_access` | [TemplateJobAccess](#resourcesdocumentsTemplateJobAccess) | repeated |  |
| `content_access` | [DocumentAccess](#resourcesdocumentsDocumentAccess) |  |  |
| `workflow` | [Workflow](#resourcesdocumentsWorkflow) | optional |  |





### resources.documents.TemplateData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `active_char` | [resources.users.User](#resourcesusersUser) |  |  |
| `documents` | [DocumentShort](#resourcesdocumentsDocumentShort) | repeated |  |
| `users` | [resources.users.UserShort](#resourcesusersUserShort) | repeated |  |
| `vehicles` | [resources.vehicles.Vehicle](#resourcesvehiclesVehicle) | repeated |  |





### resources.documents.TemplateJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesdocumentsAccessLevel) |  |  |





### resources.documents.TemplateRequirements


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `documents` | [ObjectSpecs](#resourcesdocumentsObjectSpecs) | optional |  |
| `users` | [ObjectSpecs](#resourcesdocumentsObjectSpecs) | optional |  |
| `vehicles` | [ObjectSpecs](#resourcesdocumentsObjectSpecs) | optional |  |





### resources.documents.TemplateSchema
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `requirements` | [TemplateRequirements](#resourcesdocumentsTemplateRequirements) |  |  |





### resources.documents.TemplateShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category` | [Category](#resourcesdocumentsCategory) |  |  |
| `weight` | [uint32](#uint32) |  |  |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) |  | @sanitize |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `icon` | [string](#string) | optional | @sanitize: method=StripTags |
| `schema` | [TemplateSchema](#resourcesdocumentsTemplateSchema) |  |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `workflow` | [Workflow](#resourcesdocumentsWorkflow) | optional |  |





### resources.documents.TemplateUserAccess
Dummy - DO NOT USE!





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/file/filestore.proto


### resources.file.DeleteFileRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) |  |  |
| `file_id` | [uint64](#uint64) |  |  |





### resources.file.DeleteFileResponse





### resources.file.UploadFileRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `meta` | [UploadMeta](#resourcesfileUploadMeta) |  |  |
| `data` | [bytes](#bytes) |  | Raw bytes <= 128 KiB each, browsers should only read 64 KiB at a time, but this is a buffer just in case |





### resources.file.UploadFileResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  | Unique ID for the uploaded file |
| `url` | [string](#string) |  | URL to the uploaded file |
| `file` | [File](#resourcesfileFile) |  | File info |





### resources.file.UploadMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) |  |  |
| `namespace` | [string](#string) |  | "documents", "wiki", … |
| `original_name` | [string](#string) |  |  |
| `content_type` | [string](#string) |  | optional – server re-validates |
| `size` | [int64](#int64) |  | Size in bytes |
| `reason` | [string](#string) |  | @sanitize |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/internet/access.proto


### resources.internet.PageAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageJobAccess](#resourcesinternetPageJobAccess) | repeated |  |
| `users` | [PageUserAccess](#resourcesinternetPageUserAccess) | repeated |  |





### resources.internet.PageJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesinternetAccessLevel) |  |  |





### resources.internet.PageUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `access` | [AccessLevel](#resourcesinternetAccessLevel) |  |  |




 <!-- end messages -->


### resources.internet.AccessLevel

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



## resources/internet/ads.proto


### resources.internet.Ad


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `disabled` | [bool](#bool) |  |  |
| `ad_type` | [AdType](#resourcesinternetAdType) |  |  |
| `starts_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `ends_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) |  | @sanitize: method=StripTags |
| `image` | [resources.file.File](#resourcesfileFile) | optional |  |
| `approver_id` | [int32](#int32) | optional |  |
| `approver_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.internet.AdType

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



## resources/internet/domain.proto


### resources.internet.Domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `tld_id` | [uint64](#uint64) |  |  |
| `tld` | [TLD](#resourcesinternetTLD) | optional |  |
| `active` | [bool](#bool) |  |  |
| `name` | [string](#string) |  |  |
| `transfer_code` | [string](#string) | optional |  |
| `approver_job` | [string](#string) | optional |  |
| `approver_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |





### resources.internet.TLD


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  |  |
| `internal` | [bool](#bool) |  |  |
| `creator_id` | [int32](#int32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/internet/page.proto


### resources.internet.ContentNode


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [resources.common.content.NodeType](#resourcescommoncontentNodeType) |  |  |
| `id` | [string](#string) | optional | @sanitize: method=StripTags |
| `tag` | [string](#string) |  | @sanitize: method=StripTags |
| `attrs` | [ContentNode.AttrsEntry](#resourcesinternetContentNodeAttrsEntry) | repeated | @sanitize: method=StripTags |
| `text` | [string](#string) | optional | @sanitize: method=StripTags |
| `content` | [ContentNode](#resourcesinternetContentNode) | repeated |  |
| `slots` | [ContentNode](#resourcesinternetContentNode) | repeated |  |





### resources.internet.ContentNode.AttrsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |





### resources.internet.Page


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `domain_id` | [uint64](#uint64) |  |  |
| `path` | [string](#string) |  | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `description` | [string](#string) |  | @sanitize: method=StripTags |
| `data` | [PageData](#resourcesinternetPageData) |  |  |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |





### resources.internet.PageData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `layout_type` | [PageLayoutType](#resourcesinternetPageLayoutType) |  |  |
| `node` | [ContentNode](#resourcesinternetContentNode) | optional |  |




 <!-- end messages -->


### resources.internet.PageLayoutType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `PAGE_LAYOUT_TYPE_UNSPECIFIED` | 0 |  |
| `PAGE_LAYOUT_TYPE_BASIC_PAGE` | 1 |  |
| `PAGE_LAYOUT_TYPE_LANDING_PAGE` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/internet/search.proto


### resources.internet.SearchResult


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `domain_id` | [uint64](#uint64) |  |  |
| `domain` | [Domain](#resourcesinternetDomain) | optional |  |
| `path` | [string](#string) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/activity.proto


### resources.jobs.AbsenceDateChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_begin` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `absence_end` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |





### resources.jobs.ColleagueActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `source_user_id` | [int32](#int32) | optional |  |
| `source_user` | [Colleague](#resourcesjobsColleague) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [Colleague](#resourcesjobsColleague) |  |  |
| `activity_type` | [ColleagueActivityType](#resourcesjobsColleagueActivityType) |  |  |
| `reason` | [string](#string) |  | @sanitize |
| `data` | [ColleagueActivityData](#resourcesjobsColleagueActivityData) |  |  |





### resources.jobs.ColleagueActivityData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_date` | [AbsenceDateChange](#resourcesjobsAbsenceDateChange) |  |  |
| `grade_change` | [GradeChange](#resourcesjobsGradeChange) |  |  |
| `labels_change` | [LabelsChange](#resourcesjobsLabelsChange) |  |  |
| `name_change` | [NameChange](#resourcesjobsNameChange) |  |  |





### resources.jobs.GradeChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grade` | [int32](#int32) |  |  |
| `grade_label` | [string](#string) |  |  |





### resources.jobs.LabelsChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [Label](#resourcesjobsLabel) | repeated |  |
| `removed` | [Label](#resourcesjobsLabel) | repeated |  |





### resources.jobs.NameChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prefix` | [string](#string) | optional |  |
| `suffix` | [string](#string) | optional |  |




 <!-- end messages -->


### resources.jobs.ColleagueActivityType

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



## resources/jobs/conduct.proto


### resources.jobs.ConductEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `type` | [ConductType](#resourcesjobsConductType) |  |  |
| `message` | [string](#string) |  | @sanitize |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [Colleague](#resourcesjobsColleague) | optional |  |
| `creator_id` | [int32](#int32) |  |  |
| `creator` | [Colleague](#resourcesjobsColleague) | optional |  |




 <!-- end messages -->


### resources.jobs.ConductType

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



## resources/jobs/job_settings.proto


### resources.jobs.DiscordSyncChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `plan` | [string](#string) |  |  |





### resources.jobs.DiscordSyncChanges
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `changes` | [DiscordSyncChange](#resourcesjobsDiscordSyncChange) | repeated |  |





### resources.jobs.DiscordSyncSettings
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dry_run` | [bool](#bool) |  |  |
| `user_info_sync` | [bool](#bool) |  |  |
| `user_info_sync_settings` | [UserInfoSyncSettings](#resourcesjobsUserInfoSyncSettings) |  |  |
| `status_log` | [bool](#bool) |  |  |
| `status_log_settings` | [StatusLogSettings](#resourcesjobsStatusLogSettings) |  |  |
| `jobs_absence` | [bool](#bool) |  |  |
| `jobs_absence_settings` | [JobsAbsenceSettings](#resourcesjobsJobsAbsenceSettings) |  |  |
| `group_sync_settings` | [GroupSyncSettings](#resourcesjobsGroupSyncSettings) |  |  |
| `qualifications_role_format` | [string](#string) |  |  |





### resources.jobs.GroupMapping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `from_grade` | [int32](#int32) |  |  |
| `to_grade` | [int32](#int32) |  |  |





### resources.jobs.GroupSyncSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ignored_role_ids` | [string](#string) | repeated | @sanitize: method=StripTags |





### resources.jobs.JobSettings
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_past_days` | [int32](#int32) |  |  |
| `absence_future_days` | [int32](#int32) |  |  |





### resources.jobs.JobsAbsenceSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `absence_role` | [string](#string) |  |  |





### resources.jobs.StatusLogSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel_id` | [string](#string) |  |  |





### resources.jobs.UserInfoSyncSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `employee_role_enabled` | [bool](#bool) |  |  |
| `employee_role_format` | [string](#string) |  |  |
| `grade_role_format` | [string](#string) |  |  |
| `unemployed_enabled` | [bool](#bool) |  |  |
| `unemployed_mode` | [UserInfoSyncUnemployedMode](#resourcesjobsUserInfoSyncUnemployedMode) |  |  |
| `unemployed_role_name` | [string](#string) |  |  |
| `sync_nicknames` | [bool](#bool) |  |  |
| `group_mapping` | [GroupMapping](#resourcesjobsGroupMapping) | repeated |  |




 <!-- end messages -->


### resources.jobs.UserInfoSyncUnemployedMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_UNSPECIFIED` | 0 |  |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_GIVE_ROLE` | 1 |  |
| `USER_INFO_SYNC_UNEMPLOYED_MODE_KICK` | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/job_props.proto


### resources.jobs.JobProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `livemap_marker_color` | [string](#string) |  |  |
| `quick_buttons` | [QuickButtons](#resourcesjobsQuickButtons) |  |  |
| `radio_frequency` | [string](#string) | optional |  |
| `discord_guild_id` | [string](#string) | optional |  |
| `discord_last_sync` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `discord_sync_settings` | [DiscordSyncSettings](#resourcesjobsDiscordSyncSettings) |  |  |
| `discord_sync_changes` | [DiscordSyncChanges](#resourcesjobsDiscordSyncChanges) | optional |  |
| `motd` | [string](#string) | optional |  |
| `logo_file_id` | [uint64](#uint64) | optional |  |
| `logo_file` | [resources.file.File](#resourcesfileFile) | optional |  |
| `settings` | [JobSettings](#resourcesjobsJobSettings) |  |  |





### resources.jobs.QuickButtons
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `penalty_calculator` | [bool](#bool) |  |  |
| `math_calculator` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/jobs/timeclock.proto


### resources.jobs.TimeclockEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `job` | [string](#string) |  |  |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `user` | [Colleague](#resourcesjobsColleague) | optional |  |
| `start_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `end_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `spent_time` | [float](#float) |  |  |





### resources.jobs.TimeclockStats


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `spent_time_sum` | [float](#float) |  |  |
| `spent_time_avg` | [float](#float) |  |  |
| `spent_time_max` | [float](#float) |  |  |





### resources.jobs.TimeclockWeeklyStats


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `year` | [int32](#int32) |  |  |
| `calendar_week` | [int32](#int32) |  |  |
| `sum` | [float](#float) |  |  |
| `avg` | [float](#float) |  |  |
| `max` | [float](#float) |  |  |




 <!-- end messages -->


### resources.jobs.TimeclockMode

| Name | Number | Description |
| ---- | ------ | ----------- |
| `TIMECLOCK_MODE_UNSPECIFIED` | 0 |  |
| `TIMECLOCK_MODE_DAILY` | 1 |  |
| `TIMECLOCK_MODE_WEEKLY` | 2 |  |
| `TIMECLOCK_MODE_RANGE` | 3 |  |
| `TIMECLOCK_MODE_TIMELINE` | 4 |  |



### resources.jobs.TimeclockViewMode

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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `lawbook_id` | [uint64](#uint64) |  |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `hint` | [string](#string) | optional | @sanitize |
| `fine` | [uint32](#uint32) | optional |  |
| `detention_time` | [uint32](#uint32) | optional |  |
| `stvo_points` | [uint32](#uint32) | optional |  |





### resources.laws.LawBook


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
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



## resources/livemap/heatmap.proto


### resources.livemap.HeatmapEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `w` | [double](#double) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/livemap/marker_marker.proto


### resources.livemap.CircleMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `radius` | [int32](#int32) |  |  |
| `opacity` | [float](#float) | optional |  |





### resources.livemap.IconMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `icon` | [string](#string) |  | @sanitize: method=StripTags |





### resources.livemap.MarkerData
@dbscanner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `circle` | [CircleMarker](#resourceslivemapCircleMarker) |  |  |
| `icon` | [IconMarker](#resourceslivemapIconMarker) |  |  |





### resources.livemap.MarkerMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `expires_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `name` | [string](#string) |  | @sanitize |
| `description` | [string](#string) | optional | @sanitize |
| `postal` | [string](#string) | optional | @sanitize: method=StripTags |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) |  |  |
| `type` | [MarkerType](#resourceslivemapMarkerType) |  |  |
| `data` | [MarkerData](#resourceslivemapMarkerData) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |




 <!-- end messages -->


### resources.livemap.MarkerType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `MARKER_TYPE_UNSPECIFIED` | 0 |  |
| `MARKER_TYPE_DOT` | 1 |  |
| `MARKER_TYPE_CIRCLE` | 2 |  |
| `MARKER_TYPE_ICON` | 3 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/livemap/user_marker.proto


### resources.livemap.UserMarker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `x` | [double](#double) |  |  |
| `y` | [double](#double) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `postal` | [string](#string) | optional | @sanitize: method=StripTags |
| `color` | [string](#string) | optional | @sanitize: method=StripTags |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) |  |  |
| `job_grade` | [int32](#int32) | optional |  |
| `user` | [resources.jobs.Colleague](#resourcesjobsColleague) |  |  |
| `unit_id` | [uint64](#uint64) | optional |  |
| `unit` | [resources.centrum.Unit](#resourcescentrumUnit) | optional |  |
| `hidden` | [bool](#bool) |  |  |
| `data` | [UserMarkerData](#resourceslivemapUserMarkerData) | optional |  |





### resources.livemap.UserMarkerData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `is_in_vehicle` | [bool](#bool) |  |  |
| `vehicle_plate` | [string](#string) | optional |  |
| `vehicle_updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/access.proto


### resources.mailer.Access


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [JobAccess](#resourcesmailerJobAccess) | repeated |  |
| `users` | [UserAccess](#resourcesmailerUserAccess) | repeated |  |
| `qualifications` | [QualificationAccess](#resourcesmailerQualificationAccess) | repeated |  |





### resources.mailer.JobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourcesmailerAccessLevel) |  |  |





### resources.mailer.QualificationAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `qualification_id` | [uint64](#uint64) |  |  |
| `qualification` | [resources.qualifications.QualificationShort](#resourcesqualificationsQualificationShort) | optional |  |
| `access` | [AccessLevel](#resourcesmailerAccessLevel) |  |  |





### resources.mailer.UserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `access` | [AccessLevel](#resourcesmailerAccessLevel) |  |  |




 <!-- end messages -->


### resources.mailer.AccessLevel

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



## resources/mailer/settings.proto


### resources.mailer.EmailSettings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `signature` | [string](#string) | optional | @sanitize |
| `blocked_emails` | [string](#string) | repeated | @sanitize: method=StripTags |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/email.proto


### resources.mailer.Email


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deactivated` | [bool](#bool) |  |  |
| `job` | [string](#string) | optional |  |
| `user_id` | [int32](#int32) | optional |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `email` | [string](#string) |  | @sanitize: method=StripTags |
| `email_changed` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `label` | [string](#string) | optional | @sanitize: method=StripTags |
| `access` | [Access](#resourcesmailerAccess) |  |  |
| `settings` | [EmailSettings](#resourcesmailerEmailSettings) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/message.proto


### resources.mailer.Message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |
| `sender_id` | [uint64](#uint64) |  |  |
| `sender` | [Email](#resourcesmailerEmail) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  | @sanitize |
| `data` | [MessageData](#resourcesmailerMessageData) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator_job` | [string](#string) | optional |  |





### resources.mailer.MessageAttachment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [MessageAttachmentDocument](#resourcesmailerMessageAttachmentDocument) |  |  |





### resources.mailer.MessageAttachmentDocument


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `title` | [string](#string) | optional |  |





### resources.mailer.MessageData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `attachments` | [MessageAttachment](#resourcesmailerMessageAttachment) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/thread.proto


### resources.mailer.Thread


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `creator_email_id` | [uint64](#uint64) |  |  |
| `creator_email` | [Email](#resourcesmailerEmail) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `recipients` | [ThreadRecipientEmail](#resourcesmailerThreadRecipientEmail) | repeated |  |
| `state` | [ThreadState](#resourcesmailerThreadState) | optional |  |





### resources.mailer.ThreadRecipientEmail


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `email_id` | [uint64](#uint64) |  |  |
| `email` | [Email](#resourcesmailerEmail) | optional |  |





### resources.mailer.ThreadState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread_id` | [uint64](#uint64) |  |  |
| `email_id` | [uint64](#uint64) |  |  |
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



## resources/mailer/events.proto


### resources.mailer.MailerEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_update` | [Email](#resourcesmailerEmail) |  |  |
| `email_delete` | [uint64](#uint64) |  |  |
| `email_settings_updated` | [EmailSettings](#resourcesmailerEmailSettings) |  |  |
| `thread_update` | [Thread](#resourcesmailerThread) |  |  |
| `thread_delete` | [uint64](#uint64) |  |  |
| `thread_state_update` | [ThreadState](#resourcesmailerThreadState) |  |  |
| `message_update` | [Message](#resourcesmailerMessage) |  |  |
| `message_delete` | [uint64](#uint64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/mailer/template.proto


### resources.mailer.Template


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `email_id` | [uint64](#uint64) |  |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [string](#string) |  | @sanitize |
| `creator_job` | [string](#string) | optional |  |
| `creator_id` | [int32](#int32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/notifications/client_view.proto


### resources.notifications.ClientView


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [ObjectType](#resourcesnotificationsObjectType) |  |  |
| `id` | [uint64](#uint64) | optional |  |





### resources.notifications.ObjectEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [ObjectType](#resourcesnotificationsObjectType) |  |  |
| `id` | [uint64](#uint64) | optional |  |
| `event_type` | [ObjectEventType](#resourcesnotificationsObjectEventType) |  |  |
| `user_id` | [int32](#int32) | optional |  |
| `job` | [string](#string) | optional |  |
| `data` | [google.protobuf.Struct](https://protobuf.dev/reference/protobuf/google.protobuf/#struct) | optional |  |




 <!-- end messages -->


### resources.notifications.ObjectEventType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `OBJECT_EVENT_TYPE_UNSPECIFIED` | 0 |  |
| `OBJECT_EVENT_TYPE_UPDATED` | 1 |  |
| `OBJECT_EVENT_TYPE_DELETED` | 2 |  |



### resources.notifications.ObjectType

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
| `calendar_id` | [uint64](#uint64) | optional |  |
| `calendar_entry_id` | [uint64](#uint64) | optional |  |





### resources.notifications.Data
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `link` | [Link](#resourcesnotificationsLink) | optional |  |
| `caused_by` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
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
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `read_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `title` | [resources.common.I18NItem](#resourcescommonI18NItem) |  | @sanitize |
| `type` | [NotificationType](#resourcesnotificationsNotificationType) |  |  |
| `content` | [resources.common.I18NItem](#resourcescommonI18NItem) |  | @sanitize |
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



## resources/userinfo/user_info.proto


### resources.userinfo.PollReq
PollReq: published to `userinfo.poll.request` when an active user connects or requests a refresh.



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  | The account the user belongs to |
| `user_id` | [int32](#int32) |  | The unique user identifier within the account |





### resources.userinfo.UserInfo


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





### resources.userinfo.UserInfoChanged
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
| `changed_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  | Timestamp of when the change was detected |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/notifications/events.proto


### resources.notifications.JobEvent
Job related events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resourcesjobsJobProps) |  |  |





### resources.notifications.JobGradeEvent
Job grade events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_token` | [bool](#bool) |  |  |





### resources.notifications.SystemEvent
System related events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_config` | [resources.clientconfig.ClientConfig](#resourcesclientconfigClientConfig) |  | Client configuration update (e.g., feature gates, game settings, banner message) |





### resources.notifications.UserEvent
User related events



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `refresh_token` | [bool](#bool) |  |  |
| `notification` | [Notification](#resourcesnotificationsNotification) |  | Notifications |
| `notifications_read_count` | [int32](#int32) |  |  |
| `user_info_changed` | [resources.userinfo.UserInfoChanged](#resourcesuserinfoUserInfoChanged) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/permissions/attributes.proto


### resources.permissions.AttributeValues
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `string_list` | [StringList](#resourcespermissionsStringList) |  |  |
| `job_list` | [StringList](#resourcespermissionsStringList) |  |  |
| `job_grade_list` | [JobGradeList](#resourcespermissionsJobGradeList) |  |  |





### resources.permissions.JobGradeList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fine_grained` | [bool](#bool) |  |  |
| `jobs` | [JobGradeList.JobsEntry](#resourcespermissionsJobGradeListJobsEntry) | repeated |  |
| `grades` | [JobGradeList.GradesEntry](#resourcespermissionsJobGradeListGradesEntry) | repeated |  |





### resources.permissions.JobGradeList.GradesEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [JobGrades](#resourcespermissionsJobGrades) |  |  |





### resources.permissions.JobGradeList.JobsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [int32](#int32) |  |  |





### resources.permissions.JobGrades


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grades` | [int32](#int32) | repeated |  |





### resources.permissions.RoleAttribute


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `attr_id` | [uint64](#uint64) |  |  |
| `permission_id` | [uint64](#uint64) |  |  |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `key` | [string](#string) |  |  |
| `type` | [string](#string) |  |  |
| `valid_values` | [AttributeValues](#resourcespermissionsAttributeValues) |  |  |
| `value` | [AttributeValues](#resourcespermissionsAttributeValues) |  |  |
| `max_values` | [AttributeValues](#resourcespermissionsAttributeValues) | optional |  |





### resources.permissions.StringList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `strings` | [string](#string) | repeated | @sanitize: method=StripTags |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/permissions/events.proto


### resources.permissions.JobLimitsUpdatedEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### resources.permissions.RoleIDEvent


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/permissions/permissions.proto


### resources.permissions.PermItem


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `val` | [bool](#bool) |  |  |





### resources.permissions.Permission


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `category` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `guard_name` | [string](#string) |  |  |
| `val` | [bool](#bool) |  |  |
| `order` | [int32](#int32) | optional |  |





### resources.permissions.Role


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `permissions` | [Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/settings/config.proto


### resources.settings.AppConfig
@dbscanner: json,partial



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
| `invite_url` | [string](#string) | optional | @sanitize: method=StripTags |
| `ignored_jobs` | [string](#string) | repeated | @sanitize: method=StripTags |
| `bot_presence` | [DiscordBotPresence](#resourcessettingsDiscordBotPresence) | optional |  |
| `bot_id` | [string](#string) | optional | @sanitize: method=StripTags |
| `bot_permissions` | [int64](#int64) |  |  |





### resources.settings.DiscordBotPresence


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [DiscordBotPresenceType](#resourcessettingsDiscordBotPresenceType) |  |  |
| `status` | [string](#string) | optional | @sanitize: method=StripTags |
| `url` | [string](#string) | optional | @sanitize: method=StripTags |





### resources.settings.JobInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unemployed_job` | [UnemployedJob](#resourcessettingsUnemployedJob) |  |  |
| `public_jobs` | [string](#string) | repeated | @sanitize: method=StripTags |
| `hidden_jobs` | [string](#string) | repeated | @sanitize: method=StripTags |





### resources.settings.Links


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `privacy_policy` | [string](#string) | optional | @sanitize: method=StripTags |
| `imprint` | [string](#string) | optional | @sanitize: method=StripTags |





### resources.settings.Perm


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [string](#string) |  | @sanitize: method=StripTags |
| `name` | [string](#string) |  | @sanitize: method=StripTags |





### resources.settings.Perms


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default` | [Perm](#resourcessettingsPerm) | repeated |  |





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



## resources/stats/stats.proto


### resources.stats.Stat


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [int32](#int32) | optional |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/sync/activity.proto


### resources.sync.ColleagueProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reason` | [string](#string) | optional |  |
| `props` | [resources.jobs.ColleagueProps](#resourcesjobsColleagueProps) |  |  |





### resources.sync.TimeclockUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `start` | [bool](#bool) |  |  |





### resources.sync.UserOAuth2Conn
Connect an identifier/license to the provider with the specified external id (e.g., auto discord social connect on server join)



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider_name` | [string](#string) |  |  |
| `identifier` | [string](#string) |  |  |
| `external_id` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |





### resources.sync.UserProps


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reason` | [string](#string) | optional |  |
| `props` | [resources.users.UserProps](#resourcesusersUserProps) |  |  |





### resources.sync.UserUpdate


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



## resources/sync/data.proto


### resources.sync.CitizenLocations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  |  |
| `job` | [string](#string) |  |  |
| `coords` | [resources.livemap.Coords](#resourceslivemapCoords) |  |  |
| `hidden` | [bool](#bool) |  |  |
| `remove` | [bool](#bool) |  |  |





### resources.sync.DataJobs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |





### resources.sync.DataLicenses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `licenses` | [resources.users.License](#resourcesusersLicense) | repeated |  |





### resources.sync.DataStatus


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `count` | [int64](#int64) |  |  |





### resources.sync.DataUserLocations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [CitizenLocations](#resourcessyncCitizenLocations) | repeated |  |
| `clear_all` | [bool](#bool) | optional |  |





### resources.sync.DataUsers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.users.User](#resourcesusersUser) | repeated |  |





### resources.sync.DataVehicles


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vehicles` | [resources.vehicles.Vehicle](#resourcesvehiclesVehicle) | repeated |  |





### resources.sync.DeleteUsers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_ids` | [int32](#int32) | repeated |  |





### resources.sync.DeleteVehicles


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plates` | [string](#string) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/tracker/mapping.proto


### resources.tracker.UserMapping


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `unit_id` | [uint64](#uint64) | optional |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `hidden` | [bool](#bool) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/users/activity.proto


### resources.users.CitizenDocumentRelation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [bool](#bool) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `relation` | [int32](#int32) |  | resources.documents.DocRelation enum |





### resources.users.FineChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `removed` | [bool](#bool) |  |  |
| `amount` | [int64](#int64) |  |  |





### resources.users.JailChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seconds` | [int32](#int32) |  |  |
| `admin` | [bool](#bool) |  |  |
| `location` | [string](#string) | optional |  |





### resources.users.JobChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) | optional |  |
| `job_label` | [string](#string) | optional |  |
| `grade` | [int32](#int32) | optional |  |
| `grade_label` | [string](#string) | optional |  |





### resources.users.LabelsChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [Label](#resourcesusersLabel) | repeated |  |
| `removed` | [Label](#resourcesusersLabel) | repeated |  |





### resources.users.LicenseChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [bool](#bool) |  |  |
| `licenses` | [License](#resourcesusersLicense) | repeated |  |





### resources.users.MugshotChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new` | [string](#string) | optional |  |





### resources.users.NameChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old` | [string](#string) |  |  |
| `new` | [string](#string) |  |  |





### resources.users.TrafficInfractionPointsChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `old` | [uint32](#uint32) |  |  |
| `new` | [uint32](#uint32) |  |  |





### resources.users.UserActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `type` | [UserActivityType](#resourcesusersUserActivityType) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `source_user_id` | [int32](#int32) | optional |  |
| `source_user` | [UserShort](#resourcesusersUserShort) | optional |  |
| `target_user_id` | [int32](#int32) |  |  |
| `target_user` | [UserShort](#resourcesusersUserShort) |  |  |
| `key` | [string](#string) |  | @sanitize |
| `reason` | [string](#string) |  | @sanitize |
| `data` | [UserActivityData](#resourcesusersUserActivityData) | optional |  |
| `old_value` | [string](#string) |  |  |
| `new_value` | [string](#string) |  |  |





### resources.users.UserActivityData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name_change` | [NameChange](#resourcesusersNameChange) |  |  |
| `licenses_change` | [LicenseChange](#resourcesusersLicenseChange) |  |  |
| `wanted_change` | [WantedChange](#resourcesusersWantedChange) |  | User Props |
| `traffic_infraction_points_change` | [TrafficInfractionPointsChange](#resourcesusersTrafficInfractionPointsChange) |  |  |
| `mugshot_change` | [MugshotChange](#resourcesusersMugshotChange) |  |  |
| `labels_change` | [LabelsChange](#resourcesusersLabelsChange) |  |  |
| `job_change` | [JobChange](#resourcesusersJobChange) |  |  |
| `document_relation` | [CitizenDocumentRelation](#resourcesusersCitizenDocumentRelation) |  | Docstore related |
| `jail_change` | [JailChange](#resourcesusersJailChange) |  | "Plugin" activities |
| `fine_change` | [FineChange](#resourcesusersFineChange) |  |  |





### resources.users.WantedChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `wanted` | [bool](#bool) |  |  |




 <!-- end messages -->


### resources.users.UserActivityType

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



## resources/vehicles/activity.proto


### resources.vehicles.VehicleActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `plate` | [string](#string) |  |  |
| `activity_type` | [VehicleActivityType](#resourcesvehiclesVehicleActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional | @sanitize |
| `data` | [VehicleActivityData](#resourcesvehiclesVehicleActivityData) |  |  |





### resources.vehicles.VehicleActivityData
@dbscanner: json





 <!-- end messages -->


### resources.vehicles.VehicleActivityType

| Name | Number | Description |
| ---- | ------ | ----------- |
| `VEHICLE_ACTIVITY_TYPE_UNSPECIFIED` | 0 |  |
| `VEHICLE_ACTIVITY_TYPE_WANTED` | 1 | Types for `VehicleActivityData` |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## resources/wiki/access.proto


### resources.wiki.PageAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageJobAccess](#resourceswikiPageJobAccess) | repeated |  |
| `users` | [PageUserAccess](#resourceswikiPageUserAccess) | repeated |  |





### resources.wiki.PageJobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `minimum_grade` | [int32](#int32) |  |  |
| `job_grade_label` | [string](#string) | optional |  |
| `access` | [AccessLevel](#resourceswikiAccessLevel) |  |  |





### resources.wiki.PageUserAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `target_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |
| `user` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `access` | [AccessLevel](#resourceswikiAccessLevel) |  |  |




 <!-- end messages -->


### resources.wiki.AccessLevel

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



## resources/wiki/activity.proto


### resources.wiki.PageAccessJobsDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [PageJobAccess](#resourceswikiPageJobAccess) | repeated |  |
| `to_update` | [PageJobAccess](#resourceswikiPageJobAccess) | repeated |  |
| `to_delete` | [PageJobAccess](#resourceswikiPageJobAccess) | repeated |  |





### resources.wiki.PageAccessUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [PageAccessJobsDiff](#resourceswikiPageAccessJobsDiff) |  |  |
| `users` | [PageAccessUsersDiff](#resourceswikiPageAccessUsersDiff) |  |  |





### resources.wiki.PageAccessUsersDiff


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_create` | [PageUserAccess](#resourceswikiPageUserAccess) | repeated |  |
| `to_update` | [PageUserAccess](#resourceswikiPageUserAccess) | repeated |  |
| `to_delete` | [PageUserAccess](#resourceswikiPageUserAccess) | repeated |  |





### resources.wiki.PageActivity


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `page_id` | [uint64](#uint64) |  |  |
| `activity_type` | [PageActivityType](#resourceswikiPageActivityType) |  |  |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `creator_job` | [string](#string) |  |  |
| `creator_job_label` | [string](#string) | optional |  |
| `reason` | [string](#string) | optional |  |
| `data` | [PageActivityData](#resourceswikiPageActivityData) |  |  |





### resources.wiki.PageActivityData
@dbscanner: json



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [PageUpdated](#resourceswikiPageUpdated) |  |  |
| `access_updated` | [PageAccessUpdated](#resourceswikiPageAccessUpdated) |  |  |





### resources.wiki.PageFilesChange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `added` | [int64](#int64) |  |  |
| `deleted` | [int64](#int64) |  |  |





### resources.wiki.PageUpdated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title_diff` | [string](#string) | optional |  |
| `description_diff` | [string](#string) | optional |  |
| `content_diff` | [string](#string) | optional |  |
| `files_change` | [PageFilesChange](#resourceswikiPageFilesChange) | optional |  |




 <!-- end messages -->


### resources.wiki.PageActivityType

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
| `id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  | @sanitize: method=StripTags |
| `job_label` | [string](#string) | optional |  |
| `parent_id` | [uint64](#uint64) | optional |  |
| `meta` | [PageMeta](#resourceswikiPageMeta) |  |  |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `access` | [PageAccess](#resourceswikiPageAccess) |  |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### resources.wiki.PageMeta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `created_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `updated_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `slug` | [string](#string) | optional | @sanitize: method=StripTags |
| `title` | [string](#string) |  | @sanitize |
| `description` | [string](#string) |  | @sanitize: method=StripTags |
| `creator_id` | [int32](#int32) | optional |  |
| `creator` | [resources.users.UserShort](#resourcesusersUserShort) | optional |  |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `tags` | [string](#string) | repeated | @sanitize: method=StripTags |
| `toc` | [bool](#bool) | optional |  |
| `public` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |





### resources.wiki.PageRootInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `logo_file_id` | [uint64](#uint64) | optional |  |
| `logo` | [resources.file.File](#resourcesfileFile) | optional |  |





### resources.wiki.PageShort


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `parent_id` | [uint64](#uint64) | optional |  |
| `deleted_at` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `slug` | [string](#string) | optional | @sanitize: method=StripTags |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `children` | [PageShort](#resourceswikiPageShort) | repeated |  |
| `root_info` | [PageRootInfo](#resourceswikiPageRootInfo) | optional |  |
| `level` | [int32](#int32) | optional |  |
| `draft` | [bool](#bool) |  |  |




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


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |





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
| `username` | [string](#string) |  |  |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `job_props` | [resources.jobs.JobProps](#resourcesjobsJobProps) |  |  |
| `char` | [resources.users.User](#resourcesusersUser) |  |  |
| `permissions` | [resources.permissions.Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |





### services.auth.CreateAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reg_token` | [string](#string) |  |  |
| `username` | [string](#string) |  |  |
| `password` | [string](#string) |  |  |





### services.auth.CreateAccountResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_id` | [uint64](#uint64) |  |  |





### services.auth.DeleteOAuth2ConnectionRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider` | [string](#string) |  |  |





### services.auth.DeleteOAuth2ConnectionResponse


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
| `oauth2_providers` | [resources.accounts.OAuth2Provider](#resourcesaccountsOAuth2Provider) | repeated |  |
| `oauth2_connections` | [resources.accounts.OAuth2Account](#resourcesaccountsOAuth2Account) | repeated |  |





### services.auth.GetCharactersRequest





### services.auth.GetCharactersResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chars` | [resources.accounts.Character](#resourcesaccountsCharacter) | repeated |  |





### services.auth.LoginRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `username` | [string](#string) |  |  |
| `password` | [string](#string) |  |  |





### services.auth.LoginResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `account_id` | [uint64](#uint64) |  |  |
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
| `expires` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `job_props` | [resources.jobs.JobProps](#resourcesjobsJobProps) | optional |  |
| `char` | [resources.users.User](#resourcesusersUser) |  |  |
| `permissions` | [resources.permissions.Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |




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
| `ChooseCharacter` | [ChooseCharacterRequest](#servicesauthChooseCharacterRequest) | [ChooseCharacterResponse](#servicesauthChooseCharacterResponse) |@perm |
| `GetAccountInfo` | [GetAccountInfoRequest](#servicesauthGetAccountInfoRequest) | [GetAccountInfoResponse](#servicesauthGetAccountInfoResponse) | |
| `DeleteOAuth2Connection` | [DeleteOAuth2ConnectionRequest](#servicesauthDeleteOAuth2ConnectionRequest) | [DeleteOAuth2ConnectionResponse](#servicesauthDeleteOAuth2ConnectionResponse) | |
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
| `entry` | [resources.calendar.CalendarEntry](#resourcescalendarCalendarEntry) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |





### services.calendar.CreateOrUpdateCalendarEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntry](#resourcescalendarCalendarEntry) |  |  |





### services.calendar.DeleteCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |





### services.calendar.DeleteCalendarEntryResponse





### services.calendar.DeleteCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [uint64](#uint64) |  |  |





### services.calendar.DeleteCalendarResponse





### services.calendar.GetCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |





### services.calendar.GetCalendarEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntry](#resourcescalendarCalendarEntry) |  |  |





### services.calendar.GetCalendarRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `calendar_id` | [uint64](#uint64) |  |  |





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
| `entries` | [resources.calendar.CalendarEntry](#resourcescalendarCalendarEntry) | repeated |  |





### services.calendar.ListCalendarEntriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `year` | [int32](#int32) |  |  |
| `month` | [int32](#int32) |  |  |
| `calendar_ids` | [uint64](#uint64) | repeated |  |
| `show_hidden` | [bool](#bool) | optional |  |
| `after` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### services.calendar.ListCalendarEntriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [resources.calendar.CalendarEntry](#resourcescalendarCalendarEntry) | repeated |  |





### services.calendar.ListCalendarEntryRSVPRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `entry_id` | [uint64](#uint64) |  |  |





### services.calendar.ListCalendarEntryRSVPResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `entries` | [resources.calendar.CalendarEntryRSVP](#resourcescalendarCalendarEntryRSVP) | repeated |  |





### services.calendar.ListCalendarsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `only_public` | [bool](#bool) |  |  |
| `min_access_level` | [resources.calendar.AccessLevel](#resourcescalendarAccessLevel) | optional |  |
| `after` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





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
| `entry` | [resources.calendar.CalendarEntryRSVP](#resourcescalendarCalendarEntryRSVP) |  |  |
| `subscribe` | [bool](#bool) |  |  |
| `remove` | [bool](#bool) | optional |  |





### services.calendar.RSVPCalendarEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.calendar.CalendarEntryRSVP](#resourcescalendarCalendarEntryRSVP) | optional |  |





### services.calendar.ShareCalendarEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry_id` | [uint64](#uint64) |  |  |
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
| `ListCalendars` | [ListCalendarsRequest](#servicescalendarListCalendarsRequest) | [ListCalendarsResponse](#servicescalendarListCalendarsResponse) |@perm: Name=Any |
| `GetCalendar` | [GetCalendarRequest](#servicescalendarGetCalendarRequest) | [GetCalendarResponse](#servicescalendarGetCalendarResponse) |@perm: Name=Any |
| `CreateCalendar` | [CreateCalendarRequest](#servicescalendarCreateCalendarRequest) | [CreateCalendarResponse](#servicescalendarCreateCalendarResponse) |@perm: Attrs=Fields/StringList:[]string{"Job", "Public"} |
| `UpdateCalendar` | [UpdateCalendarRequest](#servicescalendarUpdateCalendarRequest) | [UpdateCalendarResponse](#servicescalendarUpdateCalendarResponse) |@perm: Name=Any |
| `DeleteCalendar` | [DeleteCalendarRequest](#servicescalendarDeleteCalendarRequest) | [DeleteCalendarResponse](#servicescalendarDeleteCalendarResponse) |@perm: Name=Any |
| `ListCalendarEntries` | [ListCalendarEntriesRequest](#servicescalendarListCalendarEntriesRequest) | [ListCalendarEntriesResponse](#servicescalendarListCalendarEntriesResponse) |@perm: Name=Any |
| `GetUpcomingEntries` | [GetUpcomingEntriesRequest](#servicescalendarGetUpcomingEntriesRequest) | [GetUpcomingEntriesResponse](#servicescalendarGetUpcomingEntriesResponse) |@perm: Name=Any |
| `GetCalendarEntry` | [GetCalendarEntryRequest](#servicescalendarGetCalendarEntryRequest) | [GetCalendarEntryResponse](#servicescalendarGetCalendarEntryResponse) |@perm: Name=Any |
| `CreateOrUpdateCalendarEntry` | [CreateOrUpdateCalendarEntryRequest](#servicescalendarCreateOrUpdateCalendarEntryRequest) | [CreateOrUpdateCalendarEntryResponse](#servicescalendarCreateOrUpdateCalendarEntryResponse) |@perm: Name=Any |
| `DeleteCalendarEntry` | [DeleteCalendarEntryRequest](#servicescalendarDeleteCalendarEntryRequest) | [DeleteCalendarEntryResponse](#servicescalendarDeleteCalendarEntryResponse) |@perm: Name=Any |
| `ShareCalendarEntry` | [ShareCalendarEntryRequest](#servicescalendarShareCalendarEntryRequest) | [ShareCalendarEntryResponse](#servicescalendarShareCalendarEntryResponse) |@perm: Name=Any |
| `ListCalendarEntryRSVP` | [ListCalendarEntryRSVPRequest](#servicescalendarListCalendarEntryRSVPRequest) | [ListCalendarEntryRSVPResponse](#servicescalendarListCalendarEntryRSVPResponse) |@perm: Name=Any |
| `RSVPCalendarEntry` | [RSVPCalendarEntryRequest](#servicescalendarRSVPCalendarEntryRequest) | [RSVPCalendarEntryResponse](#servicescalendarRSVPCalendarEntryResponse) |@perm: Name=Any |
| `ListSubscriptions` | [ListSubscriptionsRequest](#servicescalendarListSubscriptionsRequest) | [ListSubscriptionsResponse](#servicescalendarListSubscriptionsResponse) |@perm: Name=Any |
| `SubscribeToCalendar` | [SubscribeToCalendarRequest](#servicescalendarSubscribeToCalendarRequest) | [SubscribeToCalendarResponse](#servicescalendarSubscribeToCalendarResponse) |@perm: Name=Any |

 <!-- end services -->



## services/centrum/centrum.proto


### services.centrum.AssignDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `to_add` | [uint64](#uint64) | repeated |  |
| `to_remove` | [uint64](#uint64) | repeated |  |
| `forced` | [bool](#bool) | optional |  |





### services.centrum.AssignDispatchResponse





### services.centrum.AssignUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `to_add` | [int32](#int32) | repeated |  |
| `to_remove` | [int32](#int32) | repeated |  |





### services.centrum.AssignUnitResponse





### services.centrum.CreateDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |





### services.centrum.CreateDispatchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |





### services.centrum.CreateOrUpdateUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.Unit](#resourcescentrumUnit) |  |  |





### services.centrum.CreateOrUpdateUnitResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.Unit](#resourcescentrumUnit) |  |  |





### services.centrum.DeleteDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.centrum.DeleteDispatchResponse





### services.centrum.DeleteUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |





### services.centrum.DeleteUnitResponse





### services.centrum.Dispatchers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [resources.centrum.Dispatchers](#resourcescentrumDispatchers) | repeated |  |





### services.centrum.GetDispatchHeatmapRequest





### services.centrum.GetDispatchHeatmapResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_entries` | [int32](#int32) |  |  |
| `entries` | [resources.livemap.HeatmapEntry](#resourceslivemapHeatmapEntry) | repeated |  |





### services.centrum.GetDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.centrum.GetDispatchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |





### services.centrum.GetSettingsRequest





### services.centrum.GetSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.Settings](#resourcescentrumSettings) |  |  |





### services.centrum.JobAccess


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatches` | [JobAccessEntry](#servicescentrumJobAccessEntry) | repeated |  |





### services.centrum.JobAccessEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `access` | [resources.centrum.CentrumAccessLevel](#resourcescentrumCentrumAccessLevel) |  |  |





### services.centrum.JoinUnitRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) | optional |  |





### services.centrum.JoinUnitResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit` | [resources.centrum.Unit](#resourcescentrumUnit) |  |  |





### services.centrum.LatestState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [Dispatchers](#servicescentrumDispatchers) |  |  |
| `own_unit_id` | [uint64](#uint64) | optional |  |
| `units` | [resources.centrum.Unit](#resourcescentrumUnit) | repeated | Send the current units and dispatches |
| `dispatches` | [resources.centrum.Dispatch](#resourcescentrumDispatch) | repeated |  |





### services.centrum.ListDispatchActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `id` | [uint64](#uint64) |  |  |





### services.centrum.ListDispatchActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.centrum.DispatchStatus](#resourcescentrumDispatchStatus) | repeated |  |





### services.centrum.ListDispatchesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `status` | [resources.centrum.StatusDispatch](#resourcescentrumStatusDispatch) | repeated |  |
| `not_status` | [resources.centrum.StatusDispatch](#resourcescentrumStatusDispatch) | repeated |  |
| `ids` | [uint64](#uint64) | repeated |  |
| `postal` | [string](#string) | optional |  |





### services.centrum.ListDispatchesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `dispatches` | [resources.centrum.Dispatch](#resourcescentrumDispatch) | repeated |  |





### services.centrum.ListUnitActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `id` | [uint64](#uint64) |  |  |





### services.centrum.ListUnitActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.centrum.UnitStatus](#resourcescentrumUnitStatus) | repeated |  |





### services.centrum.ListUnitsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [resources.centrum.StatusUnit](#resourcescentrumStatusUnit) | repeated |  |





### services.centrum.ListUnitsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `units` | [resources.centrum.Unit](#resourcescentrumUnit) | repeated |  |





### services.centrum.StreamHandshake


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `server_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `settings` | [resources.centrum.Settings](#resourcescentrumSettings) |  |  |
| `job_access` | [JobAccess](#servicescentrumJobAccess) |  |  |





### services.centrum.StreamRequest





### services.centrum.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `handshake` | [StreamHandshake](#servicescentrumStreamHandshake) |  |  |
| `latest_state` | [LatestState](#servicescentrumLatestState) |  |  |
| `settings` | [resources.centrum.Settings](#resourcescentrumSettings) |  |  |
| `job_access` | [JobAccess](#servicescentrumJobAccess) |  |  |
| `dispatchers` | [resources.centrum.Dispatchers](#resourcescentrumDispatchers) |  |  |
| `unit_deleted` | [uint64](#uint64) |  |  |
| `unit_updated` | [resources.centrum.Unit](#resourcescentrumUnit) |  |  |
| `unit_status` | [resources.centrum.UnitStatus](#resourcescentrumUnitStatus) |  |  |
| `dispatch_deleted` | [uint64](#uint64) |  |  |
| `dispatch_updated` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |
| `dispatch_status` | [resources.centrum.DispatchStatus](#resourcescentrumDispatchStatus) |  |  |





### services.centrum.TakeControlRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signon` | [bool](#bool) |  |  |





### services.centrum.TakeControlResponse





### services.centrum.TakeDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_ids` | [uint64](#uint64) | repeated |  |
| `resp` | [resources.centrum.TakeDispatchResp](#resourcescentrumTakeDispatchResp) |  |  |
| `reason` | [string](#string) | optional | @sanitize |





### services.centrum.TakeDispatchResponse





### services.centrum.UpdateDispatchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |





### services.centrum.UpdateDispatchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |





### services.centrum.UpdateDispatchStatusRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatch_id` | [uint64](#uint64) |  |  |
| `status` | [resources.centrum.StatusDispatch](#resourcescentrumStatusDispatch) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |





### services.centrum.UpdateDispatchStatusResponse





### services.centrum.UpdateDispatchersRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_remove` | [int32](#int32) | repeated |  |





### services.centrum.UpdateDispatchersResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dispatchers` | [resources.centrum.Dispatchers](#resourcescentrumDispatchers) |  |  |





### services.centrum.UpdateSettingsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.Settings](#resourcescentrumSettings) |  |  |





### services.centrum.UpdateSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.centrum.Settings](#resourcescentrumSettings) |  |  |





### services.centrum.UpdateUnitStatusRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unit_id` | [uint64](#uint64) |  |  |
| `status` | [resources.centrum.StatusUnit](#resourcescentrumStatusUnit) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `code` | [string](#string) | optional | @sanitize |





### services.centrum.UpdateUnitStatusResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.centrum.CentrumService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `UpdateSettings` | [UpdateSettingsRequest](#servicescentrumUpdateSettingsRequest) | [UpdateSettingsResponse](#servicescentrumUpdateSettingsResponse) |@perm: Attrs=Access/StringList:[]string{"Shared"} |
| `CreateDispatch` | [CreateDispatchRequest](#servicescentrumCreateDispatchRequest) | [CreateDispatchResponse](#servicescentrumCreateDispatchResponse) |@perm |
| `UpdateDispatch` | [UpdateDispatchRequest](#servicescentrumUpdateDispatchRequest) | [UpdateDispatchResponse](#servicescentrumUpdateDispatchResponse) |@perm |
| `DeleteDispatch` | [DeleteDispatchRequest](#servicescentrumDeleteDispatchRequest) | [DeleteDispatchResponse](#servicescentrumDeleteDispatchResponse) |@perm |
| `TakeControl` | [TakeControlRequest](#servicescentrumTakeControlRequest) | [TakeControlResponse](#servicescentrumTakeControlResponse) |@perm |
| `AssignDispatch` | [AssignDispatchRequest](#servicescentrumAssignDispatchRequest) | [AssignDispatchResponse](#servicescentrumAssignDispatchResponse) |@perm: Name=TakeControl |
| `AssignUnit` | [AssignUnitRequest](#servicescentrumAssignUnitRequest) | [AssignUnitResponse](#servicescentrumAssignUnitResponse) |@perm: Name=TakeControl |
| `GetDispatchHeatmap` | [GetDispatchHeatmapRequest](#servicescentrumGetDispatchHeatmapRequest) | [GetDispatchHeatmapResponse](#servicescentrumGetDispatchHeatmapResponse) |@perm: Name=TakeControl |
| `UpdateDispatchers` | [UpdateDispatchersRequest](#servicescentrumUpdateDispatchersRequest) | [UpdateDispatchersResponse](#servicescentrumUpdateDispatchersResponse) |@perm |
| `Stream` | [StreamRequest](#servicescentrumStreamRequest) | [StreamResponse](#servicescentrumStreamResponse) stream |@perm |
| `GetSettings` | [GetSettingsRequest](#servicescentrumGetSettingsRequest) | [GetSettingsResponse](#servicescentrumGetSettingsResponse) |@perm: Name=Stream |
| `JoinUnit` | [JoinUnitRequest](#servicescentrumJoinUnitRequest) | [JoinUnitResponse](#servicescentrumJoinUnitResponse) |@perm: Name=Stream |
| `ListUnits` | [ListUnitsRequest](#servicescentrumListUnitsRequest) | [ListUnitsResponse](#servicescentrumListUnitsResponse) |@perm: Name=Stream |
| `ListUnitActivity` | [ListUnitActivityRequest](#servicescentrumListUnitActivityRequest) | [ListUnitActivityResponse](#servicescentrumListUnitActivityResponse) |@perm: Name=Stream |
| `GetDispatch` | [GetDispatchRequest](#servicescentrumGetDispatchRequest) | [GetDispatchResponse](#servicescentrumGetDispatchResponse) |@perm: Name=Stream |
| `ListDispatches` | [ListDispatchesRequest](#servicescentrumListDispatchesRequest) | [ListDispatchesResponse](#servicescentrumListDispatchesResponse) |@perm: Name=Stream |
| `ListDispatchActivity` | [ListDispatchActivityRequest](#servicescentrumListDispatchActivityRequest) | [ListDispatchActivityResponse](#servicescentrumListDispatchActivityResponse) |@perm: Name=Stream |
| `CreateOrUpdateUnit` | [CreateOrUpdateUnitRequest](#servicescentrumCreateOrUpdateUnitRequest) | [CreateOrUpdateUnitResponse](#servicescentrumCreateOrUpdateUnitResponse) |@perm |
| `DeleteUnit` | [DeleteUnitRequest](#servicescentrumDeleteUnitRequest) | [DeleteUnitResponse](#servicescentrumDeleteUnitResponse) |@perm |
| `TakeDispatch` | [TakeDispatchRequest](#servicescentrumTakeDispatchRequest) | [TakeDispatchResponse](#servicescentrumTakeDispatchResponse) |@perm |
| `UpdateUnitStatus` | [UpdateUnitStatusRequest](#servicescentrumUpdateUnitStatusRequest) | [UpdateUnitStatusResponse](#servicescentrumUpdateUnitStatusResponse) |@perm: Name=TakeDispatch |
| `UpdateDispatchStatus` | [UpdateDispatchStatusRequest](#servicescentrumUpdateDispatchStatusRequest) | [UpdateDispatchStatusResponse](#servicescentrumUpdateDispatchStatusResponse) |@perm: Name=TakeDispatch |

 <!-- end services -->



## services/citizens/citizens.proto


### services.citizens.DeleteAvatarRequest





### services.citizens.DeleteAvatarResponse





### services.citizens.DeleteMugshotRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `reason` | [string](#string) |  | @sanitize |





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
| `open_fines` | [uint64](#uint64) | optional |  |





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
| `types` | [resources.users.UserActivityType](#resourcesusersUserActivityType) | repeated |  |





### services.citizens.ListUserActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.users.UserActivity](#resourcesusersUserActivity) | repeated |  |





### services.citizens.ManageLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.Label](#resourcesusersLabel) | repeated |  |





### services.citizens.ManageLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.Label](#resourcesusersLabel) | repeated |  |





### services.citizens.SetUserPropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.users.UserProps](#resourcesusersUserProps) |  |  |
| `reason` | [string](#string) |  | @sanitize |





### services.citizens.SetUserPropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.users.UserProps](#resourcesusersUserProps) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.citizens.CitizensService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCitizens` | [ListCitizensRequest](#servicescitizensListCitizensRequest) | [ListCitizensResponse](#servicescitizensListCitizensResponse) |@perm: Attrs=Fields/StringList:[]string{"PhoneNumber", "Licenses", "UserProps.Wanted", "UserProps.Job", "UserProps.TrafficInfractionPoints", "UserProps.OpenFines", "UserProps.BloodType", "UserProps.Mugshot", "UserProps.Labels", "UserProps.Email"} |
| `GetUser` | [GetUserRequest](#servicescitizensGetUserRequest) | [GetUserResponse](#servicescitizensGetUserResponse) |@perm: Attrs=Jobs/JobGradeList |
| `ListUserActivity` | [ListUserActivityRequest](#servicescitizensListUserActivityRequest) | [ListUserActivityResponse](#servicescitizensListUserActivityResponse) |@perm: Attrs=Fields/StringList:[]string{"SourceUser", "Own"} |
| `SetUserProps` | [SetUserPropsRequest](#servicescitizensSetUserPropsRequest) | [SetUserPropsResponse](#servicescitizensSetUserPropsResponse) |@perm: Attrs=Fields/StringList:[]string{"Wanted", "Job", "TrafficInfractionPoints", "Mugshot", "Labels"} |
| `UploadAvatar` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=Any buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `DeleteAvatar` | [DeleteAvatarRequest](#servicescitizensDeleteAvatarRequest) | [DeleteAvatarResponse](#servicescitizensDeleteAvatarResponse) |@perm: Name=Any |
| `UploadMugshot` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=SetUserProps buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `DeleteMugshot` | [DeleteMugshotRequest](#servicescitizensDeleteMugshotRequest) | [DeleteMugshotResponse](#servicescitizensDeleteMugshotResponse) |@perm: Name=SetUserProps |
| `ManageLabels` | [ManageLabelsRequest](#servicescitizensManageLabelsRequest) | [ManageLabelsResponse](#servicescitizensManageLabelsResponse) |@perm |

 <!-- end services -->



## services/completor/completor.proto


### services.completor.CompleteCitizenLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |





### services.completor.CompleteCitizenLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.users.Label](#resourcesusersLabel) | repeated |  |





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
| `users` | [resources.users.UserShort](#resourcesusersUserShort) | repeated |  |





### services.completor.CompleteDocumentCategoriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |





### services.completor.CompleteDocumentCategoriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `categories` | [resources.documents.Category](#resourcesdocumentsCategory) | repeated |  |





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
| `CompleteCitizens` | [CompleteCitizensRequest](#servicescompletorCompleteCitizensRequest) | [CompleteCitizensResponse](#servicescompletorCompleteCitizensResponse) |@perm |
| `CompleteJobs` | [CompleteJobsRequest](#servicescompletorCompleteJobsRequest) | [CompleteJobsResponse](#servicescompletorCompleteJobsResponse) |@perm: Name=Any |
| `CompleteDocumentCategories` | [CompleteDocumentCategoriesRequest](#servicescompletorCompleteDocumentCategoriesRequest) | [CompleteDocumentCategoriesResponse](#servicescompletorCompleteDocumentCategoriesResponse) |@perm: Attrs=Jobs/JobList |
| `ListLawBooks` | [ListLawBooksRequest](#servicescompletorListLawBooksRequest) | [ListLawBooksResponse](#servicescompletorListLawBooksResponse) |@perm: Name=Any |
| `CompleteCitizenLabels` | [CompleteCitizenLabelsRequest](#servicescompletorCompleteCitizenLabelsRequest) | [CompleteCitizenLabelsResponse](#servicescompletorCompleteCitizenLabelsResponse) |@perm: Attrs=Jobs/JobList |

 <!-- end services -->



## services/documents/collab.proto

 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.CollabService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `JoinRoom` | [.resources.collab.ClientPacket](#resourcescollabClientPacket) stream | [.resources.collab.ServerPacket](#resourcescollabServerPacket) stream |@perm: Name=documents.DocumentsService/UpdateDocument buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |

 <!-- end services -->



## services/documents/documents.proto


### services.documents.AddDocumentReferenceRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reference` | [resources.documents.DocumentReference](#resourcesdocumentsDocumentReference) |  |  |





### services.documents.AddDocumentReferenceResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.AddDocumentRelationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relation` | [resources.documents.DocumentRelation](#resourcesdocumentsDocumentRelation) |  |  |





### services.documents.AddDocumentRelationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.ChangeDocumentOwnerRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `new_user_id` | [int32](#int32) | optional |  |





### services.documents.ChangeDocumentOwnerResponse





### services.documents.CreateDocumentReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `request_type` | [resources.documents.DocActivityType](#resourcesdocumentsDocActivityType) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `data` | [resources.documents.DocActivityData](#resourcesdocumentsDocActivityData) | optional |  |





### services.documents.CreateDocumentReqResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.documents.DocRequest](#resourcesdocumentsDocRequest) |  |  |





### services.documents.CreateDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `template_id` | [uint64](#uint64) | optional |  |
| `template_data` | [resources.documents.TemplateData](#resourcesdocumentsTemplateData) | optional |  |





### services.documents.CreateDocumentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.CreateOrUpdateCategoryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [resources.documents.Category](#resourcesdocumentsCategory) |  |  |





### services.documents.CreateOrUpdateCategoryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `category` | [resources.documents.Category](#resourcesdocumentsCategory) |  |  |





### services.documents.CreateTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resourcesdocumentsTemplate) |  |  |





### services.documents.CreateTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.DeleteCategoryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.DeleteCategoryResponse





### services.documents.DeleteCommentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment_id` | [uint64](#uint64) |  |  |





### services.documents.DeleteCommentResponse





### services.documents.DeleteDocumentReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request_id` | [uint64](#uint64) |  |  |





### services.documents.DeleteDocumentReqResponse





### services.documents.DeleteDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `reason` | [string](#string) | optional | @sanitize: method=StripTags |





### services.documents.DeleteDocumentResponse





### services.documents.DeleteTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.DeleteTemplateResponse





### services.documents.EditCommentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resourcesdocumentsComment) |  |  |





### services.documents.EditCommentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resourcesdocumentsComment) |  |  |





### services.documents.GetCommentsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `document_id` | [uint64](#uint64) |  |  |





### services.documents.GetCommentsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `comments` | [resources.documents.Comment](#resourcesdocumentsComment) | repeated |  |





### services.documents.GetDocumentAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |





### services.documents.GetDocumentAccessResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [resources.documents.DocumentAccess](#resourcesdocumentsDocumentAccess) |  |  |





### services.documents.GetDocumentReferencesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |





### services.documents.GetDocumentReferencesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `references` | [resources.documents.DocumentReference](#resourcesdocumentsDocumentReference) | repeated |  |





### services.documents.GetDocumentRelationsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |





### services.documents.GetDocumentRelationsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relations` | [resources.documents.DocumentRelation](#resourcesdocumentsDocumentRelation) | repeated |  |





### services.documents.GetDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `info_only` | [bool](#bool) | optional |  |





### services.documents.GetDocumentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [resources.documents.Document](#resourcesdocumentsDocument) |  |  |
| `access` | [resources.documents.DocumentAccess](#resourcesdocumentsDocumentAccess) |  |  |





### services.documents.GetTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template_id` | [uint64](#uint64) |  |  |
| `data` | [resources.documents.TemplateData](#resourcesdocumentsTemplateData) | optional |  |
| `render` | [bool](#bool) | optional |  |





### services.documents.GetTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resourcesdocumentsTemplate) |  |  |
| `rendered` | [bool](#bool) |  |  |





### services.documents.ListCategoriesRequest





### services.documents.ListCategoriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `categories` | [resources.documents.Category](#resourcesdocumentsCategory) | repeated |  |





### services.documents.ListDocumentActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `document_id` | [uint64](#uint64) |  |  |
| `activity_types` | [resources.documents.DocActivityType](#resourcesdocumentsDocActivityType) | repeated | Search params |





### services.documents.ListDocumentActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.documents.DocActivity](#resourcesdocumentsDocActivity) | repeated |  |





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
| `document_id` | [uint64](#uint64) |  |  |





### services.documents.ListDocumentReqsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `requests` | [resources.documents.DocRequest](#resourcesdocumentsDocRequest) | repeated |  |





### services.documents.ListDocumentsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `search` | [string](#string) | optional | Search params |
| `category_ids` | [uint64](#uint64) | repeated |  |
| `creator_ids` | [int32](#int32) | repeated |  |
| `from` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `to` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `closed` | [bool](#bool) | optional |  |
| `document_ids` | [uint64](#uint64) | repeated |  |
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
| `templates` | [resources.documents.TemplateShort](#resourcesdocumentsTemplateShort) | repeated |  |





### services.documents.ListUserDocumentsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_id` | [int32](#int32) |  |  |
| `relations` | [resources.documents.DocRelation](#resourcesdocumentsDocRelation) | repeated |  |
| `closed` | [bool](#bool) | optional |  |





### services.documents.ListUserDocumentsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `relations` | [resources.documents.DocumentRelation](#resourcesdocumentsDocumentRelation) | repeated |  |





### services.documents.PostCommentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resourcesdocumentsComment) |  |  |





### services.documents.PostCommentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `comment` | [resources.documents.Comment](#resourcesdocumentsComment) |  |  |





### services.documents.RemoveDocumentReferenceRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.RemoveDocumentReferenceResponse





### services.documents.RemoveDocumentRelationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.documents.RemoveDocumentRelationResponse





### services.documents.SetDocumentAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `access` | [resources.documents.DocumentAccess](#resourcesdocumentsDocumentAccess) |  |  |





### services.documents.SetDocumentAccessResponse





### services.documents.SetDocumentReminderRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `reminder_time` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `message` | [string](#string) | optional | @sanitize: method=StripTags |





### services.documents.SetDocumentReminderResponse





### services.documents.ToggleDocumentPinRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `state` | [bool](#bool) |  |  |
| `personal` | [bool](#bool) | optional | If true, the pin is personal and not shared with other job members |





### services.documents.ToggleDocumentPinResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pin` | [resources.documents.DocumentPin](#resourcesdocumentsDocumentPin) | optional |  |





### services.documents.ToggleDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `closed` | [bool](#bool) |  |  |





### services.documents.ToggleDocumentResponse





### services.documents.UpdateDocumentReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `request_id` | [uint64](#uint64) |  |  |
| `reason` | [string](#string) | optional | @sanitize |
| `data` | [resources.documents.DocActivityData](#resourcesdocumentsDocActivityData) | optional |  |
| `accepted` | [bool](#bool) |  |  |





### services.documents.UpdateDocumentReqResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `request` | [resources.documents.DocRequest](#resourcesdocumentsDocRequest) |  |  |





### services.documents.UpdateDocumentRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document_id` | [uint64](#uint64) |  |  |
| `category_id` | [uint64](#uint64) | optional |  |
| `title` | [string](#string) |  | @sanitize: method=StripTags |
| `content` | [resources.common.content.Content](#resourcescommoncontentContent) |  |  |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |
| `data` | [string](#string) | optional |  |
| `state` | [string](#string) |  | @sanitize |
| `closed` | [bool](#bool) |  |  |
| `draft` | [bool](#bool) |  |  |
| `public` | [bool](#bool) |  |  |
| `access` | [resources.documents.DocumentAccess](#resourcesdocumentsDocumentAccess) | optional |  |
| `files` | [resources.file.File](#resourcesfileFile) | repeated |  |





### services.documents.UpdateDocumentResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `document` | [resources.documents.Document](#resourcesdocumentsDocument) |  |  |





### services.documents.UpdateTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resourcesdocumentsTemplate) |  |  |





### services.documents.UpdateTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.documents.Template](#resourcesdocumentsTemplate) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.documents.DocumentsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTemplates` | [ListTemplatesRequest](#servicesdocumentsListTemplatesRequest) | [ListTemplatesResponse](#servicesdocumentsListTemplatesResponse) |@perm |
| `GetTemplate` | [GetTemplateRequest](#servicesdocumentsGetTemplateRequest) | [GetTemplateResponse](#servicesdocumentsGetTemplateResponse) |@perm: Name=ListTemplates |
| `CreateTemplate` | [CreateTemplateRequest](#servicesdocumentsCreateTemplateRequest) | [CreateTemplateResponse](#servicesdocumentsCreateTemplateResponse) |@perm |
| `UpdateTemplate` | [UpdateTemplateRequest](#servicesdocumentsUpdateTemplateRequest) | [UpdateTemplateResponse](#servicesdocumentsUpdateTemplateResponse) |@perm: Name=CreateTemplate |
| `DeleteTemplate` | [DeleteTemplateRequest](#servicesdocumentsDeleteTemplateRequest) | [DeleteTemplateResponse](#servicesdocumentsDeleteTemplateResponse) |@perm |
| `ListDocuments` | [ListDocumentsRequest](#servicesdocumentsListDocumentsRequest) | [ListDocumentsResponse](#servicesdocumentsListDocumentsResponse) |@perm |
| `GetDocument` | [GetDocumentRequest](#servicesdocumentsGetDocumentRequest) | [GetDocumentResponse](#servicesdocumentsGetDocumentResponse) |@perm: Name=ListDocuments |
| `CreateDocument` | [CreateDocumentRequest](#servicesdocumentsCreateDocumentRequest) | [CreateDocumentResponse](#servicesdocumentsCreateDocumentResponse) |@perm: Name=UpdateDocument |
| `UpdateDocument` | [UpdateDocumentRequest](#servicesdocumentsUpdateDocumentRequest) | [UpdateDocumentResponse](#servicesdocumentsUpdateDocumentResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `DeleteDocument` | [DeleteDocumentRequest](#servicesdocumentsDeleteDocumentRequest) | [DeleteDocumentResponse](#servicesdocumentsDeleteDocumentResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `ToggleDocument` | [ToggleDocumentRequest](#servicesdocumentsToggleDocumentRequest) | [ToggleDocumentResponse](#servicesdocumentsToggleDocumentResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `ChangeDocumentOwner` | [ChangeDocumentOwnerRequest](#servicesdocumentsChangeDocumentOwnerRequest) | [ChangeDocumentOwnerResponse](#servicesdocumentsChangeDocumentOwnerResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `GetDocumentReferences` | [GetDocumentReferencesRequest](#servicesdocumentsGetDocumentReferencesRequest) | [GetDocumentReferencesResponse](#servicesdocumentsGetDocumentReferencesResponse) |@perm: Name=ListDocuments |
| `GetDocumentRelations` | [GetDocumentRelationsRequest](#servicesdocumentsGetDocumentRelationsRequest) | [GetDocumentRelationsResponse](#servicesdocumentsGetDocumentRelationsResponse) |@perm: Name=ListDocuments |
| `AddDocumentReference` | [AddDocumentReferenceRequest](#servicesdocumentsAddDocumentReferenceRequest) | [AddDocumentReferenceResponse](#servicesdocumentsAddDocumentReferenceResponse) |@perm |
| `RemoveDocumentReference` | [RemoveDocumentReferenceRequest](#servicesdocumentsRemoveDocumentReferenceRequest) | [RemoveDocumentReferenceResponse](#servicesdocumentsRemoveDocumentReferenceResponse) |@perm: Name=AddDocumentReference |
| `AddDocumentRelation` | [AddDocumentRelationRequest](#servicesdocumentsAddDocumentRelationRequest) | [AddDocumentRelationResponse](#servicesdocumentsAddDocumentRelationResponse) |@perm |
| `RemoveDocumentRelation` | [RemoveDocumentRelationRequest](#servicesdocumentsRemoveDocumentRelationRequest) | [RemoveDocumentRelationResponse](#servicesdocumentsRemoveDocumentRelationResponse) |@perm: Name=AddDocumentRelation |
| `GetComments` | [GetCommentsRequest](#servicesdocumentsGetCommentsRequest) | [GetCommentsResponse](#servicesdocumentsGetCommentsResponse) |@perm: Name=ListDocuments |
| `PostComment` | [PostCommentRequest](#servicesdocumentsPostCommentRequest) | [PostCommentResponse](#servicesdocumentsPostCommentResponse) |@perm: Name=ListDocuments |
| `EditComment` | [EditCommentRequest](#servicesdocumentsEditCommentRequest) | [EditCommentResponse](#servicesdocumentsEditCommentResponse) |@perm: Name=ListDocuments |
| `DeleteComment` | [DeleteCommentRequest](#servicesdocumentsDeleteCommentRequest) | [DeleteCommentResponse](#servicesdocumentsDeleteCommentResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `GetDocumentAccess` | [GetDocumentAccessRequest](#servicesdocumentsGetDocumentAccessRequest) | [GetDocumentAccessResponse](#servicesdocumentsGetDocumentAccessResponse) |@perm: Name=ListDocuments |
| `SetDocumentAccess` | [SetDocumentAccessRequest](#servicesdocumentsSetDocumentAccessRequest) | [SetDocumentAccessResponse](#servicesdocumentsSetDocumentAccessResponse) |@perm: Name=UpdateDocument |
| `ListDocumentActivity` | [ListDocumentActivityRequest](#servicesdocumentsListDocumentActivityRequest) | [ListDocumentActivityResponse](#servicesdocumentsListDocumentActivityResponse) |@perm |
| `ListDocumentReqs` | [ListDocumentReqsRequest](#servicesdocumentsListDocumentReqsRequest) | [ListDocumentReqsResponse](#servicesdocumentsListDocumentReqsResponse) |@perm |
| `CreateDocumentReq` | [CreateDocumentReqRequest](#servicesdocumentsCreateDocumentReqRequest) | [CreateDocumentReqResponse](#servicesdocumentsCreateDocumentReqResponse) |@perm: Attrs=Types/StringList:[]string{"Access", "Closure", "Update", "Deletion", "OwnerChange"} |
| `UpdateDocumentReq` | [UpdateDocumentReqRequest](#servicesdocumentsUpdateDocumentReqRequest) | [UpdateDocumentReqResponse](#servicesdocumentsUpdateDocumentReqResponse) |@perm: Name=CreateDocumentReq |
| `DeleteDocumentReq` | [DeleteDocumentReqRequest](#servicesdocumentsDeleteDocumentReqRequest) | [DeleteDocumentReqResponse](#servicesdocumentsDeleteDocumentReqResponse) |@perm |
| `ListUserDocuments` | [ListUserDocumentsRequest](#servicesdocumentsListUserDocumentsRequest) | [ListUserDocumentsResponse](#servicesdocumentsListUserDocumentsResponse) |@perm |
| `ListCategories` | [ListCategoriesRequest](#servicesdocumentsListCategoriesRequest) | [ListCategoriesResponse](#servicesdocumentsListCategoriesResponse) |@perm |
| `CreateOrUpdateCategory` | [CreateOrUpdateCategoryRequest](#servicesdocumentsCreateOrUpdateCategoryRequest) | [CreateOrUpdateCategoryResponse](#servicesdocumentsCreateOrUpdateCategoryResponse) |@perm |
| `DeleteCategory` | [DeleteCategoryRequest](#servicesdocumentsDeleteCategoryRequest) | [DeleteCategoryResponse](#servicesdocumentsDeleteCategoryResponse) |@perm |
| `ListDocumentPins` | [ListDocumentPinsRequest](#servicesdocumentsListDocumentPinsRequest) | [ListDocumentPinsResponse](#servicesdocumentsListDocumentPinsResponse) |@perm: Name=ListDocuments |
| `ToggleDocumentPin` | [ToggleDocumentPinRequest](#servicesdocumentsToggleDocumentPinRequest) | [ToggleDocumentPinResponse](#servicesdocumentsToggleDocumentPinResponse) |@perm: Attrs=Types/StringList:[]string{"JobWide"} |
| `SetDocumentReminder` | [SetDocumentReminderRequest](#servicesdocumentsSetDocumentReminderRequest) | [SetDocumentReminderResponse](#servicesdocumentsSetDocumentReminderResponse) |@perm |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=UpdateDocument |

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
| `Upload` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=Superuser buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `ListFiles` | [ListFilesRequest](#servicesfilestoreListFilesRequest) | [ListFilesResponse](#servicesfilestoreListFilesResponse) |@perm: Name=Superuser |
| `DeleteFile` | [.resources.file.DeleteFileRequest](#resourcesfileDeleteFileRequest) | [.resources.file.DeleteFileResponse](#resourcesfileDeleteFileResponse) |@perm: Name=Superuser |
| `DeleteFileByPath` | [DeleteFileByPathRequest](#servicesfilestoreDeleteFileByPathRequest) | [DeleteFileByPathResponse](#servicesfilestoreDeleteFileByPathResponse) |@perm: Name=Superuser |

 <!-- end services -->



## services/internet/ads.proto


### services.internet.GetAdsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ad_type` | [resources.internet.AdType](#resourcesinternetAdType) |  |  |
| `count` | [int32](#int32) |  |  |





### services.internet.GetAdsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ads` | [resources.internet.Ad](#resourcesinternetAd) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.internet.AdsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetAds` | [GetAdsRequest](#servicesinternetGetAdsRequest) | [GetAdsResponse](#servicesinternetGetAdsResponse) |@perm: Name=Any |

 <!-- end services -->



## services/internet/domain.proto


### services.internet.CheckDomainAvailabilityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tld_id` | [uint64](#uint64) |  |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |





### services.internet.CheckDomainAvailabilityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `available` | [bool](#bool) |  |  |
| `transferable` | [bool](#bool) | optional |  |





### services.internet.ListDomainsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |





### services.internet.ListDomainsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `domains` | [resources.internet.Domain](#resourcesinternetDomain) | repeated |  |





### services.internet.ListTLDsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `internal` | [bool](#bool) | optional |  |





### services.internet.ListTLDsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tlds` | [resources.internet.TLD](#resourcesinternetTLD) | repeated |  |





### services.internet.RegisterDomainRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tld_id` | [uint64](#uint64) |  |  |
| `name` | [string](#string) |  | @sanitize: method=StripTags |
| `transfer_code` | [string](#string) | optional | In case a domain will be transfered |





### services.internet.RegisterDomainResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [resources.internet.Domain](#resourcesinternetDomain) |  |  |





### services.internet.UpdateDomainRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain_id` | [uint64](#uint64) |  |  |
| `transferable` | [bool](#bool) |  |  |





### services.internet.UpdateDomainResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [resources.internet.Domain](#resourcesinternetDomain) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.internet.DomainService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTLDs` | [ListTLDsRequest](#servicesinternetListTLDsRequest) | [ListTLDsResponse](#servicesinternetListTLDsResponse) |@perm: Name=Any |
| `CheckDomainAvailability` | [CheckDomainAvailabilityRequest](#servicesinternetCheckDomainAvailabilityRequest) | [CheckDomainAvailabilityResponse](#servicesinternetCheckDomainAvailabilityResponse) |@perm: Name=Any |
| `RegisterDomain` | [RegisterDomainRequest](#servicesinternetRegisterDomainRequest) | [RegisterDomainResponse](#servicesinternetRegisterDomainResponse) |@perm: Name=Any |
| `ListDomains` | [ListDomainsRequest](#servicesinternetListDomainsRequest) | [ListDomainsResponse](#servicesinternetListDomainsResponse) |@perm: Name=Any |
| `UpdateDomain` | [UpdateDomainRequest](#servicesinternetUpdateDomainRequest) | [UpdateDomainResponse](#servicesinternetUpdateDomainResponse) |@perm: Name=Any |

 <!-- end services -->



## services/internet/internet.proto


### services.internet.GetPageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  |  |
| `path` | [string](#string) |  |  |





### services.internet.GetPageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.internet.Page](#resourcesinternetPage) | optional |  |





### services.internet.SearchRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) |  |  |
| `domain_id` | [uint64](#uint64) | optional |  |





### services.internet.SearchResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `results` | [resources.internet.SearchResult](#resourcesinternetSearchResult) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.internet.InternetService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Search` | [SearchRequest](#servicesinternetSearchRequest) | [SearchResponse](#servicesinternetSearchResponse) |@perm: Name=Any |
| `GetPage` | [GetPageRequest](#servicesinternetGetPageRequest) | [GetPageResponse](#servicesinternetGetPageResponse) |@perm: Name=Any |

 <!-- end services -->



## services/jobs/conduct.proto


### services.jobs.CreateConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resourcesjobsConductEntry) |  |  |





### services.jobs.CreateConductEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resourcesjobsConductEntry) |  |  |





### services.jobs.DeleteConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.jobs.DeleteConductEntryResponse





### services.jobs.ListConductEntriesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `types` | [resources.jobs.ConductType](#resourcesjobsConductType) | repeated | Search params |
| `show_expired` | [bool](#bool) | optional |  |
| `user_ids` | [int32](#int32) | repeated |  |
| `ids` | [uint64](#uint64) | repeated |  |





### services.jobs.ListConductEntriesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `entries` | [resources.jobs.ConductEntry](#resourcesjobsConductEntry) | repeated |  |





### services.jobs.UpdateConductEntryRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resourcesjobsConductEntry) |  |  |





### services.jobs.UpdateConductEntryResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entry` | [resources.jobs.ConductEntry](#resourcesjobsConductEntry) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.jobs.ConductService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListConductEntries` | [ListConductEntriesRequest](#servicesjobsListConductEntriesRequest) | [ListConductEntriesResponse](#servicesjobsListConductEntriesResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "All"} |
| `CreateConductEntry` | [CreateConductEntryRequest](#servicesjobsCreateConductEntryRequest) | [CreateConductEntryResponse](#servicesjobsCreateConductEntryResponse) |@perm |
| `UpdateConductEntry` | [UpdateConductEntryRequest](#servicesjobsUpdateConductEntryRequest) | [UpdateConductEntryResponse](#servicesjobsUpdateConductEntryResponse) |@perm |
| `DeleteConductEntry` | [DeleteConductEntryRequest](#servicesjobsDeleteConductEntryRequest) | [DeleteConductEntryResponse](#servicesjobsDeleteConductEntryResponse) |@perm |

 <!-- end services -->



## services/jobs/jobs.proto


### services.jobs.GetColleagueLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `search` | [string](#string) | optional |  |





### services.jobs.GetColleagueLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.Label](#resourcesjobsLabel) | repeated |  |





### services.jobs.GetColleagueLabelsStatsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `label_ids` | [uint64](#uint64) | repeated |  |





### services.jobs.GetColleagueLabelsStatsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `count` | [resources.jobs.LabelCount](#resourcesjobsLabelCount) | repeated |  |





### services.jobs.GetColleagueRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) |  |  |
| `info_only` | [bool](#bool) | optional |  |





### services.jobs.GetColleagueResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague` | [resources.jobs.Colleague](#resourcesjobsColleague) |  |  |





### services.jobs.GetMOTDRequest





### services.jobs.GetMOTDResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  |  |





### services.jobs.GetSelfRequest





### services.jobs.GetSelfResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `colleague` | [resources.jobs.Colleague](#resourcesjobsColleague) |  |  |





### services.jobs.ListColleagueActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_ids` | [int32](#int32) | repeated | Search params |
| `activity_types` | [resources.jobs.ColleagueActivityType](#resourcesjobsColleagueActivityType) | repeated |  |





### services.jobs.ListColleagueActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.jobs.ColleagueActivity](#resourcesjobsColleagueActivity) | repeated |  |





### services.jobs.ListColleaguesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `search` | [string](#string) |  | Search params |
| `user_ids` | [int32](#int32) | repeated |  |
| `user_only` | [bool](#bool) | optional |  |
| `absent` | [bool](#bool) | optional |  |
| `label_ids` | [uint64](#uint64) | repeated |  |
| `name_prefix` | [string](#string) | optional |  |
| `name_suffix` | [string](#string) | optional |  |





### services.jobs.ListColleaguesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `colleagues` | [resources.jobs.Colleague](#resourcesjobsColleague) | repeated |  |





### services.jobs.ManageLabelsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.Label](#resourcesjobsLabel) | repeated |  |





### services.jobs.ManageLabelsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `labels` | [resources.jobs.Label](#resourcesjobsLabel) | repeated |  |





### services.jobs.SetColleaguePropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.jobs.ColleagueProps](#resourcesjobsColleagueProps) |  |  |
| `reason` | [string](#string) |  | @sanitize |





### services.jobs.SetColleaguePropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.jobs.ColleagueProps](#resourcesjobsColleagueProps) |  |  |





### services.jobs.SetMOTDRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `motd` | [string](#string) |  | @sanitize: method=StripTags |





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
| `ListColleagues` | [ListColleaguesRequest](#servicesjobsListColleaguesRequest) | [ListColleaguesResponse](#servicesjobsListColleaguesResponse) |@perm |
| `GetSelf` | [GetSelfRequest](#servicesjobsGetSelfRequest) | [GetSelfResponse](#servicesjobsGetSelfResponse) |@perm: Name=ListColleagues |
| `GetColleague` | [GetColleagueRequest](#servicesjobsGetColleagueRequest) | [GetColleagueResponse](#servicesjobsGetColleagueResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"Note", "Labels"} |
| `ListColleagueActivity` | [ListColleagueActivityRequest](#servicesjobsListColleagueActivityRequest) | [ListColleagueActivityResponse](#servicesjobsListColleagueActivityResponse) |@perm: Attrs=Types/StringList:[]string{"HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME"} |
| `SetColleagueProps` | [SetColleaguePropsRequest](#servicesjobsSetColleaguePropsRequest) | [SetColleaguePropsResponse](#servicesjobsSetColleaguePropsResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Types/StringList:[]string{"AbsenceDate", "Note", "Labels", "Name"} |
| `GetColleagueLabels` | [GetColleagueLabelsRequest](#servicesjobsGetColleagueLabelsRequest) | [GetColleagueLabelsResponse](#servicesjobsGetColleagueLabelsResponse) |@perm: Name=GetColleague |
| `ManageLabels` | [ManageLabelsRequest](#servicesjobsManageLabelsRequest) | [ManageLabelsResponse](#servicesjobsManageLabelsResponse) |@perm |
| `GetColleagueLabelsStats` | [GetColleagueLabelsStatsRequest](#servicesjobsGetColleagueLabelsStatsRequest) | [GetColleagueLabelsStatsResponse](#servicesjobsGetColleagueLabelsStatsResponse) |@perm: Name=GetColleague |
| `GetMOTD` | [GetMOTDRequest](#servicesjobsGetMOTDRequest) | [GetMOTDResponse](#servicesjobsGetMOTDResponse) |@perm: Name=Any |
| `SetMOTD` | [SetMOTDRequest](#servicesjobsSetMOTDRequest) | [SetMOTDResponse](#servicesjobsSetMOTDResponse) |@perm |

 <!-- end services -->



## services/jobs/timeclock.proto


### services.jobs.GetTimeclockStatsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_id` | [int32](#int32) | optional |  |





### services.jobs.GetTimeclockStatsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stats` | [resources.jobs.TimeclockStats](#resourcesjobsTimeclockStats) |  |  |
| `weekly` | [resources.jobs.TimeclockWeeklyStats](#resourcesjobsTimeclockWeeklyStats) | repeated |  |





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
| `colleagues` | [resources.jobs.Colleague](#resourcesjobsColleague) | repeated |  |





### services.jobs.ListTimeclockRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_mode` | [resources.jobs.TimeclockViewMode](#resourcesjobsTimeclockViewMode) |  | Search params |
| `mode` | [resources.jobs.TimeclockMode](#resourcesjobsTimeclockMode) |  |  |
| `date` | [resources.common.database.DateRange](#resourcescommondatabaseDateRange) | optional |  |
| `per_day` | [bool](#bool) |  |  |
| `user_ids` | [int32](#int32) | repeated |  |





### services.jobs.ListTimeclockResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `stats` | [resources.jobs.TimeclockStats](#resourcesjobsTimeclockStats) |  |  |
| `stats_weekly` | [resources.jobs.TimeclockWeeklyStats](#resourcesjobsTimeclockWeeklyStats) | repeated |  |
| `daily` | [TimeclockDay](#servicesjobsTimeclockDay) |  |  |
| `weekly` | [TimeclockWeekly](#servicesjobsTimeclockWeekly) |  |  |
| `range` | [TimeclockRange](#servicesjobsTimeclockRange) |  |  |





### services.jobs.TimeclockDay


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `entries` | [resources.jobs.TimeclockEntry](#resourcesjobsTimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |





### services.jobs.TimeclockRange


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `entries` | [resources.jobs.TimeclockEntry](#resourcesjobsTimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |





### services.jobs.TimeclockWeekly


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `date` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) |  |  |
| `entries` | [resources.jobs.TimeclockEntry](#resourcesjobsTimeclockEntry) | repeated |  |
| `sum` | [int64](#int64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.jobs.TimeclockService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListTimeclock` | [ListTimeclockRequest](#servicesjobsListTimeclockRequest) | [ListTimeclockResponse](#servicesjobsListTimeclockResponse) |@perm: Attrs=Access/StringList:[]string{"All"} |
| `GetTimeclockStats` | [GetTimeclockStatsRequest](#servicesjobsGetTimeclockStatsRequest) | [GetTimeclockStatsResponse](#servicesjobsGetTimeclockStatsResponse) |@perm: Name=ListTimeclock |
| `ListInactiveEmployees` | [ListInactiveEmployeesRequest](#servicesjobsListInactiveEmployeesRequest) | [ListInactiveEmployeesResponse](#servicesjobsListInactiveEmployeesResponse) |@perm |

 <!-- end services -->



## services/livemap/livemap.proto


### services.livemap.CreateOrUpdateMarkerRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `marker` | [resources.livemap.MarkerMarker](#resourceslivemapMarkerMarker) |  |  |





### services.livemap.CreateOrUpdateMarkerResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `marker` | [resources.livemap.MarkerMarker](#resourceslivemapMarkerMarker) |  |  |





### services.livemap.DeleteMarkerRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.livemap.DeleteMarkerResponse





### services.livemap.JobsList


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |
| `markers` | [resources.jobs.Job](#resourcesjobsJob) | repeated |  |





### services.livemap.MarkerMarkersUpdates


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [resources.livemap.MarkerMarker](#resourceslivemapMarkerMarker) | repeated |  |
| `deleted` | [uint64](#uint64) | repeated |  |
| `part` | [int32](#int32) |  |  |
| `partial` | [bool](#bool) |  |  |





### services.livemap.Snapshot
A roll-up of the entire USERLOC bucket. Published every N seconds on `$KV.user_locations._snapshot` with the headers: Nats-Rollup: all KV-Operation: ROLLUP



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `markers` | [resources.livemap.UserMarker](#resourceslivemapUserMarker) | repeated | All currently-known user markers, already filtered for obsolete PURGE/DELETE events. |





### services.livemap.StreamRequest





### services.livemap.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `user_on_duty` | [bool](#bool) | optional |  |
| `jobs` | [JobsList](#serviceslivemapJobsList) |  |  |
| `markers` | [MarkerMarkersUpdates](#serviceslivemapMarkerMarkersUpdates) |  |  |
| `snapshot` | [Snapshot](#serviceslivemapSnapshot) |  |  |
| `user_update` | [resources.livemap.UserMarker](#resourceslivemapUserMarker) |  |  |
| `user_delete` | [UserDelete](#serviceslivemapUserDelete) |  |  |





### services.livemap.UserDelete


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [int32](#int32) |  | The user ID of the user that was deleted. |
| `job` | [string](#string) |  | The job of the user that was deleted. |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.livemap.LivemapService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `Stream` | [StreamRequest](#serviceslivemapStreamRequest) | [StreamResponse](#serviceslivemapStreamResponse) stream |@perm: Attrs=Markers/JobList|Players/JobGradeList |
| `CreateOrUpdateMarker` | [CreateOrUpdateMarkerRequest](#serviceslivemapCreateOrUpdateMarkerRequest) | [CreateOrUpdateMarkerResponse](#serviceslivemapCreateOrUpdateMarkerResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `DeleteMarker` | [DeleteMarkerRequest](#serviceslivemapDeleteMarkerRequest) | [DeleteMarkerResponse](#serviceslivemapDeleteMarkerResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |

 <!-- end services -->



## services/mailer/mailer.proto


### services.mailer.CreateOrUpdateEmailRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.Email](#resourcesmailerEmail) |  |  |





### services.mailer.CreateOrUpdateEmailResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.Email](#resourcesmailerEmail) |  |  |





### services.mailer.CreateOrUpdateTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.Template](#resourcesmailerTemplate) |  |  |





### services.mailer.CreateOrUpdateTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.Template](#resourcesmailerTemplate) |  |  |





### services.mailer.CreateThreadRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.Thread](#resourcesmailerThread) |  |  |
| `message` | [resources.mailer.Message](#resourcesmailerMessage) |  |  |
| `recipients` | [string](#string) | repeated | @sanitize: method=StripTags |





### services.mailer.CreateThreadResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.Thread](#resourcesmailerThread) |  |  |





### services.mailer.DeleteEmailRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.mailer.DeleteEmailResponse





### services.mailer.DeleteMessageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |
| `message_id` | [uint64](#uint64) |  |  |





### services.mailer.DeleteMessageResponse





### services.mailer.DeleteTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `id` | [uint64](#uint64) |  |  |





### services.mailer.DeleteTemplateResponse





### services.mailer.DeleteThreadRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |





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
| `id` | [uint64](#uint64) |  |  |





### services.mailer.GetEmailResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email` | [resources.mailer.Email](#resourcesmailerEmail) |  |  |





### services.mailer.GetEmailSettingsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |





### services.mailer.GetEmailSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.EmailSettings](#resourcesmailerEmailSettings) |  |  |





### services.mailer.GetTemplateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `template_id` | [uint64](#uint64) |  |  |





### services.mailer.GetTemplateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `template` | [resources.mailer.Template](#resourcesmailerTemplate) |  |  |





### services.mailer.GetThreadRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |





### services.mailer.GetThreadResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `thread` | [resources.mailer.Thread](#resourcesmailerThread) |  |  |





### services.mailer.GetThreadStateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |





### services.mailer.GetThreadStateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.ThreadState](#resourcesmailerThreadState) |  |  |





### services.mailer.ListEmailsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `all` | [bool](#bool) | optional | Search params |





### services.mailer.ListEmailsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `emails` | [resources.mailer.Email](#resourcesmailerEmail) | repeated |  |





### services.mailer.ListTemplatesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `email_id` | [uint64](#uint64) |  |  |





### services.mailer.ListTemplatesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `templates` | [resources.mailer.Template](#resourcesmailerTemplate) | repeated |  |





### services.mailer.ListThreadMessagesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `email_id` | [uint64](#uint64) |  |  |
| `thread_id` | [uint64](#uint64) |  |  |
| `after` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |





### services.mailer.ListThreadMessagesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `messages` | [resources.mailer.Message](#resourcesmailerMessage) | repeated |  |





### services.mailer.ListThreadsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `email_ids` | [uint64](#uint64) | repeated | Search params |
| `unread` | [bool](#bool) | optional |  |
| `archived` | [bool](#bool) | optional |  |





### services.mailer.ListThreadsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `threads` | [resources.mailer.Thread](#resourcesmailerThread) | repeated |  |





### services.mailer.PostMessageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [resources.mailer.Message](#resourcesmailerMessage) |  |  |
| `recipients` | [string](#string) | repeated | @sanitize: method=StripTags |





### services.mailer.PostMessageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `message` | [resources.mailer.Message](#resourcesmailerMessage) |  |  |





### services.mailer.SearchThreadsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `search` | [string](#string) |  | Search params |





### services.mailer.SearchThreadsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `messages` | [resources.mailer.Message](#resourcesmailerMessage) | repeated |  |





### services.mailer.SetEmailSettingsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.EmailSettings](#resourcesmailerEmailSettings) |  |  |





### services.mailer.SetEmailSettingsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `settings` | [resources.mailer.EmailSettings](#resourcesmailerEmailSettings) |  |  |





### services.mailer.SetThreadStateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.ThreadState](#resourcesmailerThreadState) |  |  |





### services.mailer.SetThreadStateResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [resources.mailer.ThreadState](#resourcesmailerThreadState) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.mailer.MailerService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListEmails` | [ListEmailsRequest](#servicesmailerListEmailsRequest) | [ListEmailsResponse](#servicesmailerListEmailsResponse) |@perm |
| `GetEmail` | [GetEmailRequest](#servicesmailerGetEmailRequest) | [GetEmailResponse](#servicesmailerGetEmailResponse) |@perm: Name=ListEmails |
| `CreateOrUpdateEmail` | [CreateOrUpdateEmailRequest](#servicesmailerCreateOrUpdateEmailRequest) | [CreateOrUpdateEmailResponse](#servicesmailerCreateOrUpdateEmailResponse) |@perm: Attrs=Fields/StringList:[]string{"Job"} |
| `DeleteEmail` | [DeleteEmailRequest](#servicesmailerDeleteEmailRequest) | [DeleteEmailResponse](#servicesmailerDeleteEmailResponse) |@perm |
| `GetEmailProposals` | [GetEmailProposalsRequest](#servicesmailerGetEmailProposalsRequest) | [GetEmailProposalsResponse](#servicesmailerGetEmailProposalsResponse) |@perm: Name=ListEmails |
| `ListTemplates` | [ListTemplatesRequest](#servicesmailerListTemplatesRequest) | [ListTemplatesResponse](#servicesmailerListTemplatesResponse) |@perm: Name=ListEmails |
| `GetTemplate` | [GetTemplateRequest](#servicesmailerGetTemplateRequest) | [GetTemplateResponse](#servicesmailerGetTemplateResponse) |@perm: Name=ListEmails |
| `CreateOrUpdateTemplate` | [CreateOrUpdateTemplateRequest](#servicesmailerCreateOrUpdateTemplateRequest) | [CreateOrUpdateTemplateResponse](#servicesmailerCreateOrUpdateTemplateResponse) |@perm: Name=ListEmails |
| `DeleteTemplate` | [DeleteTemplateRequest](#servicesmailerDeleteTemplateRequest) | [DeleteTemplateResponse](#servicesmailerDeleteTemplateResponse) |@perm: Name=ListEmails |
| `ListThreads` | [ListThreadsRequest](#servicesmailerListThreadsRequest) | [ListThreadsResponse](#servicesmailerListThreadsResponse) |@perm: Name=ListEmails |
| `GetThread` | [GetThreadRequest](#servicesmailerGetThreadRequest) | [GetThreadResponse](#servicesmailerGetThreadResponse) |@perm: Name=ListEmails |
| `CreateThread` | [CreateThreadRequest](#servicesmailerCreateThreadRequest) | [CreateThreadResponse](#servicesmailerCreateThreadResponse) |@perm: Name=ListEmails |
| `DeleteThread` | [DeleteThreadRequest](#servicesmailerDeleteThreadRequest) | [DeleteThreadResponse](#servicesmailerDeleteThreadResponse) |@perm: Name=Superuser |
| `GetThreadState` | [GetThreadStateRequest](#servicesmailerGetThreadStateRequest) | [GetThreadStateResponse](#servicesmailerGetThreadStateResponse) |@perm: Name=ListEmails |
| `SetThreadState` | [SetThreadStateRequest](#servicesmailerSetThreadStateRequest) | [SetThreadStateResponse](#servicesmailerSetThreadStateResponse) |@perm: Name=ListEmails |
| `SearchThreads` | [SearchThreadsRequest](#servicesmailerSearchThreadsRequest) | [SearchThreadsResponse](#servicesmailerSearchThreadsResponse) |@perm: Name=ListEmails |
| `ListThreadMessages` | [ListThreadMessagesRequest](#servicesmailerListThreadMessagesRequest) | [ListThreadMessagesResponse](#servicesmailerListThreadMessagesResponse) |@perm: Name=ListEmails |
| `PostMessage` | [PostMessageRequest](#servicesmailerPostMessageRequest) | [PostMessageResponse](#servicesmailerPostMessageResponse) |@perm: Name=ListEmails |
| `DeleteMessage` | [DeleteMessageRequest](#servicesmailerDeleteMessageRequest) | [DeleteMessageResponse](#servicesmailerDeleteMessageResponse) |@perm: Name=Superuser |
| `GetEmailSettings` | [GetEmailSettingsRequest](#servicesmailerGetEmailSettingsRequest) | [GetEmailSettingsResponse](#servicesmailerGetEmailSettingsResponse) |@perm: Name=ListEmails |
| `SetEmailSettings` | [SetEmailSettingsRequest](#servicesmailerSetEmailSettingsRequest) | [SetEmailSettingsResponse](#servicesmailerSetEmailSettingsResponse) |@perm: Name=ListEmails |

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
| `ids` | [uint64](#uint64) | repeated |  |
| `all` | [bool](#bool) | optional |  |





### services.notifications.MarkNotificationsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `updated` | [uint64](#uint64) |  |  |





### services.notifications.StreamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_view` | [resources.notifications.ClientView](#resourcesnotificationsClientView) |  |  |





### services.notifications.StreamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `notification_count` | [int32](#int32) |  |  |
| `restart` | [bool](#bool) | optional |  |
| `user_event` | [resources.notifications.UserEvent](#resourcesnotificationsUserEvent) |  |  |
| `job_event` | [resources.notifications.JobEvent](#resourcesnotificationsJobEvent) |  |  |
| `job_grade_event` | [resources.notifications.JobGradeEvent](#resourcesnotificationsJobGradeEvent) |  |  |
| `system_event` | [resources.notifications.SystemEvent](#resourcesnotificationsSystemEvent) |  |  |
| `mailer_event` | [resources.mailer.MailerEvent](#resourcesmailerMailerEvent) |  |  |
| `object_event` | [resources.notifications.ObjectEvent](#resourcesnotificationsObjectEvent) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.notifications.NotificationsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `GetNotifications` | [GetNotificationsRequest](#servicesnotificationsGetNotificationsRequest) | [GetNotificationsResponse](#servicesnotificationsGetNotificationsResponse) |@perm: Name=Any |
| `MarkNotifications` | [MarkNotificationsRequest](#servicesnotificationsMarkNotificationsRequest) | [MarkNotificationsResponse](#servicesnotificationsMarkNotificationsResponse) |@perm: Name=Any |
| `Stream` | [StreamRequest](#servicesnotificationsStreamRequest) stream | [StreamResponse](#servicesnotificationsStreamResponse) stream |@perm: Name=Any |

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
| `grading` | [resources.qualifications.ExamGrading](#resourcesqualificationsExamGrading) | optional |  |





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
| `qualification_id` | [uint64](#uint64) |  |  |





### services.qualifications.DeleteQualificationReqRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |





### services.qualifications.DeleteQualificationReqResponse





### services.qualifications.DeleteQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |





### services.qualifications.DeleteQualificationResponse





### services.qualifications.DeleteQualificationResultRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result_id` | [uint64](#uint64) |  |  |





### services.qualifications.DeleteQualificationResultResponse





### services.qualifications.GetExamInfoRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |





### services.qualifications.GetExamInfoResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.QualificationShort](#resourcesqualificationsQualificationShort) |  |  |
| `question_count` | [int32](#int32) |  |  |
| `exam_user` | [resources.qualifications.ExamUser](#resourcesqualificationsExamUser) | optional |  |





### services.qualifications.GetQualificationAccessRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |





### services.qualifications.GetQualificationAccessResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [resources.qualifications.QualificationAccess](#resourcesqualificationsQualificationAccess) |  |  |





### services.qualifications.GetQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `with_exam` | [bool](#bool) | optional |  |





### services.qualifications.GetQualificationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.Qualification](#resourcesqualificationsQualification) |  |  |





### services.qualifications.GetUserExamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `user_id` | [int32](#int32) |  |  |





### services.qualifications.GetUserExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exam` | [resources.qualifications.ExamQuestions](#resourcesqualificationsExamQuestions) |  |  |
| `exam_user` | [resources.qualifications.ExamUser](#resourcesqualificationsExamUser) |  |  |
| `responses` | [resources.qualifications.ExamResponses](#resourcesqualificationsExamResponses) |  |  |
| `grading` | [resources.qualifications.ExamGrading](#resourcesqualificationsExamGrading) |  |  |





### services.qualifications.ListQualificationRequestsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `qualification_id` | [uint64](#uint64) | optional | Search params |
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
| `qualification_id` | [uint64](#uint64) | optional | Search params |
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
| `qualification_id` | [uint64](#uint64) |  |  |
| `access` | [resources.qualifications.QualificationAccess](#resourcesqualificationsQualificationAccess) |  |  |





### services.qualifications.SetQualificationAccessResponse





### services.qualifications.SubmitExamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `responses` | [resources.qualifications.ExamResponses](#resourcesqualificationsExamResponses) |  |  |





### services.qualifications.SubmitExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `duration` | [google.protobuf.Duration](https://protobuf.dev/reference/protobuf/google.protobuf/#duration) |  |  |





### services.qualifications.TakeExamRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |
| `cancel` | [bool](#bool) | optional |  |





### services.qualifications.TakeExamResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exam` | [resources.qualifications.ExamQuestions](#resourcesqualificationsExamQuestions) |  |  |
| `exam_user` | [resources.qualifications.ExamUser](#resourcesqualificationsExamUser) |  |  |





### services.qualifications.UpdateQualificationRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification` | [resources.qualifications.Qualification](#resourcesqualificationsQualification) |  |  |





### services.qualifications.UpdateQualificationResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `qualification_id` | [uint64](#uint64) |  |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.qualifications.QualificationsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListQualifications` | [ListQualificationsRequest](#servicesqualificationsListQualificationsRequest) | [ListQualificationsResponse](#servicesqualificationsListQualificationsResponse) |@perm |
| `GetQualification` | [GetQualificationRequest](#servicesqualificationsGetQualificationRequest) | [GetQualificationResponse](#servicesqualificationsGetQualificationResponse) |@perm: Name=ListQualifications |
| `CreateQualification` | [CreateQualificationRequest](#servicesqualificationsCreateQualificationRequest) | [CreateQualificationResponse](#servicesqualificationsCreateQualificationResponse) |@perm: Name=UpdateQualification |
| `UpdateQualification` | [UpdateQualificationRequest](#servicesqualificationsUpdateQualificationRequest) | [UpdateQualificationResponse](#servicesqualificationsUpdateQualificationResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"}|Fields/StringList:[]string{"Public"} |
| `DeleteQualification` | [DeleteQualificationRequest](#servicesqualificationsDeleteQualificationRequest) | [DeleteQualificationResponse](#servicesqualificationsDeleteQualificationResponse) |@perm: Attrs=Access/StringList:[]string{"Own", "Lower_Rank", "Same_Rank", "Any"} |
| `ListQualificationRequests` | [ListQualificationRequestsRequest](#servicesqualificationsListQualificationRequestsRequest) | [ListQualificationRequestsResponse](#servicesqualificationsListQualificationRequestsResponse) |@perm: Name=ListQualifications |
| `CreateOrUpdateQualificationRequest` | [CreateOrUpdateQualificationRequestRequest](#servicesqualificationsCreateOrUpdateQualificationRequestRequest) | [CreateOrUpdateQualificationRequestResponse](#servicesqualificationsCreateOrUpdateQualificationRequestResponse) |@perm: Name=ListQualifications |
| `DeleteQualificationReq` | [DeleteQualificationReqRequest](#servicesqualificationsDeleteQualificationReqRequest) | [DeleteQualificationReqResponse](#servicesqualificationsDeleteQualificationReqResponse) |@perm: Name=ListQualifications |
| `ListQualificationsResults` | [ListQualificationsResultsRequest](#servicesqualificationsListQualificationsResultsRequest) | [ListQualificationsResultsResponse](#servicesqualificationsListQualificationsResultsResponse) |@perm: Name=ListQualifications |
| `CreateOrUpdateQualificationResult` | [CreateOrUpdateQualificationResultRequest](#servicesqualificationsCreateOrUpdateQualificationResultRequest) | [CreateOrUpdateQualificationResultResponse](#servicesqualificationsCreateOrUpdateQualificationResultResponse) |@perm: Name=ListQualifications |
| `DeleteQualificationResult` | [DeleteQualificationResultRequest](#servicesqualificationsDeleteQualificationResultRequest) | [DeleteQualificationResultResponse](#servicesqualificationsDeleteQualificationResultResponse) |@perm: Name=ListQualifications |
| `GetExamInfo` | [GetExamInfoRequest](#servicesqualificationsGetExamInfoRequest) | [GetExamInfoResponse](#servicesqualificationsGetExamInfoResponse) |@perm: Name=ListQualifications |
| `TakeExam` | [TakeExamRequest](#servicesqualificationsTakeExamRequest) | [TakeExamResponse](#servicesqualificationsTakeExamResponse) |@perm: Name=ListQualifications |
| `SubmitExam` | [SubmitExamRequest](#servicesqualificationsSubmitExamRequest) | [SubmitExamResponse](#servicesqualificationsSubmitExamResponse) |@perm: Name=ListQualifications |
| `GetUserExam` | [GetUserExamRequest](#servicesqualificationsGetUserExamRequest) | [GetUserExamResponse](#servicesqualificationsGetUserExamResponse) |@perm: Name=ListQualifications |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=UpdateQualification |

 <!-- end services -->



## services/settings/accounts.proto


### services.settings.DeleteAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.settings.DeleteAccountResponse





### services.settings.DisconnectOAuth2ConnectionRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `provider_name` | [string](#string) |  |  |





### services.settings.DisconnectOAuth2ConnectionResponse





### services.settings.ListAccountsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `license` | [string](#string) | optional | Search params |
| `enabled` | [bool](#bool) | optional |  |
| `username` | [string](#string) | optional |  |
| `external_id` | [string](#string) | optional |  |





### services.settings.ListAccountsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `accounts` | [resources.accounts.Account](#resourcesaccountsAccount) | repeated |  |





### services.settings.UpdateAccountRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
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
| `ListAccounts` | [ListAccountsRequest](#servicessettingsListAccountsRequest) | [ListAccountsResponse](#servicessettingsListAccountsResponse) |@perm: Name=Superuser |
| `UpdateAccount` | [UpdateAccountRequest](#servicessettingsUpdateAccountRequest) | [UpdateAccountResponse](#servicessettingsUpdateAccountResponse) |@perm: Name=Superuser |
| `DisconnectOAuth2Connection` | [DisconnectOAuth2ConnectionRequest](#servicessettingsDisconnectOAuth2ConnectionRequest) | [DisconnectOAuth2ConnectionResponse](#servicessettingsDisconnectOAuth2ConnectionResponse) |@perm: Name=Superuser |
| `DeleteAccount` | [DeleteAccountRequest](#servicessettingsDeleteAccountRequest) | [DeleteAccountResponse](#servicessettingsDeleteAccountResponse) |@perm: Name=Superuser |

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
| `GetAppConfig` | [GetAppConfigRequest](#servicessettingsGetAppConfigRequest) | [GetAppConfigResponse](#servicessettingsGetAppConfigResponse) |@perm: Name=Superuser |
| `UpdateAppConfig` | [UpdateAppConfigRequest](#servicessettingsUpdateAppConfigRequest) | [UpdateAppConfigResponse](#servicessettingsUpdateAppConfigResponse) |@perm: Name=Superuser |

 <!-- end services -->



## services/settings/cron.proto


### services.settings.ListCronjobsRequest





### services.settings.ListCronjobsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.common.cron.Cronjob](#resourcescommoncronCronjob) | repeated |  |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.CronService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListCronjobs` | [ListCronjobsRequest](#servicessettingsListCronjobsRequest) | [ListCronjobsResponse](#servicessettingsListCronjobsResponse) |@perm: Name=Superuser |

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
| `id` | [uint64](#uint64) |  |  |





### services.settings.DeleteLawBookResponse





### services.settings.DeleteLawRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.settings.DeleteLawResponse




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.settings.LawsService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `CreateOrUpdateLawBook` | [CreateOrUpdateLawBookRequest](#servicessettingsCreateOrUpdateLawBookRequest) | [CreateOrUpdateLawBookResponse](#servicessettingsCreateOrUpdateLawBookResponse) |@perm |
| `DeleteLawBook` | [DeleteLawBookRequest](#servicessettingsDeleteLawBookRequest) | [DeleteLawBookResponse](#servicessettingsDeleteLawBookResponse) |@perm |
| `CreateOrUpdateLaw` | [CreateOrUpdateLawRequest](#servicessettingsCreateOrUpdateLawRequest) | [CreateOrUpdateLawResponse](#servicessettingsCreateOrUpdateLawResponse) |@perm: Name=CreateOrUpdateLawBook |
| `DeleteLaw` | [DeleteLawRequest](#servicessettingsDeleteLawRequest) | [DeleteLawResponse](#servicessettingsDeleteLawResponse) |@perm: Name=DeleteLawBook |

 <!-- end services -->



## services/settings/settings.proto


### services.settings.AttrsUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_update` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |
| `to_remove` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |





### services.settings.CreateRoleRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `grade` | [int32](#int32) |  |  |





### services.settings.CreateRoleResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.Role](#resourcespermissionsRole) |  |  |





### services.settings.DeleteFactionRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### services.settings.DeleteFactionResponse





### services.settings.DeleteJobLogoRequest





### services.settings.DeleteJobLogoResponse





### services.settings.DeleteRoleRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.settings.DeleteRoleResponse





### services.settings.GetAllPermissionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### services.settings.GetAllPermissionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permissions` | [resources.permissions.Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |





### services.settings.GetEffectivePermissionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |





### services.settings.GetEffectivePermissionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.Role](#resourcespermissionsRole) |  |  |
| `permissions` | [resources.permissions.Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |





### services.settings.GetJobLimitsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |





### services.settings.GetJobLimitsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `job_label` | [string](#string) | optional |  |
| `permissions` | [resources.permissions.Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |





### services.settings.GetJobPropsRequest





### services.settings.GetJobPropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resourcesjobsJobProps) |  |  |





### services.settings.GetPermissionsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role_id` | [uint64](#uint64) |  |  |





### services.settings.GetPermissionsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `permissions` | [resources.permissions.Permission](#resourcespermissionsPermission) | repeated |  |
| `attributes` | [resources.permissions.RoleAttribute](#resourcespermissionsRoleAttribute) | repeated |  |





### services.settings.GetRoleRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.settings.GetRoleResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `role` | [resources.permissions.Role](#resourcespermissionsRole) |  |  |





### services.settings.GetRolesRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `lowest_rank` | [bool](#bool) | optional |  |





### services.settings.GetRolesResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `roles` | [resources.permissions.Role](#resourcespermissionsRole) | repeated |  |





### services.settings.GetStatusRequest





### services.settings.GetStatusResponse
TODO






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





### services.settings.PermsUpdate


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to_update` | [resources.permissions.PermItem](#resourcespermissionsPermItem) | repeated |  |
| `to_remove` | [resources.permissions.PermItem](#resourcespermissionsPermItem) | repeated |  |





### services.settings.SetJobPropsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resourcesjobsJobProps) |  |  |





### services.settings.SetJobPropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job_props` | [resources.jobs.JobProps](#resourcesjobsJobProps) |  |  |





### services.settings.UpdateJobLimitsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `perms` | [PermsUpdate](#servicessettingsPermsUpdate) | optional |  |
| `attrs` | [AttrsUpdate](#servicessettingsAttrsUpdate) | optional |  |





### services.settings.UpdateJobLimitsResponse





### services.settings.UpdateRolePermsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `perms` | [PermsUpdate](#servicessettingsPermsUpdate) | optional |  |
| `attrs` | [AttrsUpdate](#servicessettingsAttrsUpdate) | optional |  |





### services.settings.UpdateRolePermsResponse





### services.settings.ViewAuditLogRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `sort` | [resources.common.database.Sort](#resourcescommondatabaseSort) | optional |  |
| `user_ids` | [int32](#int32) | repeated | Search params |
| `from` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `to` | [resources.timestamp.Timestamp](#resourcestimestampTimestamp) | optional |  |
| `services` | [string](#string) | repeated | @sanitize: method=StripTags |
| `methods` | [string](#string) | repeated | @sanitize: method=StripTags |
| `search` | [string](#string) | optional |  |
| `states` | [resources.audit.EventType](#resourcesauditEventType) | repeated |  |





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
| `GetJobProps` | [GetJobPropsRequest](#servicessettingsGetJobPropsRequest) | [GetJobPropsResponse](#servicessettingsGetJobPropsResponse) |@perm |
| `SetJobProps` | [SetJobPropsRequest](#servicessettingsSetJobPropsRequest) | [SetJobPropsResponse](#servicessettingsSetJobPropsResponse) |@perm |
| `GetRoles` | [GetRolesRequest](#servicessettingsGetRolesRequest) | [GetRolesResponse](#servicessettingsGetRolesResponse) |@perm |
| `GetRole` | [GetRoleRequest](#servicessettingsGetRoleRequest) | [GetRoleResponse](#servicessettingsGetRoleResponse) |@perm: Name=GetRoles |
| `CreateRole` | [CreateRoleRequest](#servicessettingsCreateRoleRequest) | [CreateRoleResponse](#servicessettingsCreateRoleResponse) |@perm |
| `DeleteRole` | [DeleteRoleRequest](#servicessettingsDeleteRoleRequest) | [DeleteRoleResponse](#servicessettingsDeleteRoleResponse) |@perm |
| `UpdateRolePerms` | [UpdateRolePermsRequest](#servicessettingsUpdateRolePermsRequest) | [UpdateRolePermsResponse](#servicessettingsUpdateRolePermsResponse) |@perm |
| `GetPermissions` | [GetPermissionsRequest](#servicessettingsGetPermissionsRequest) | [GetPermissionsResponse](#servicessettingsGetPermissionsResponse) |@perm: Name=GetRoles |
| `GetEffectivePermissions` | [GetEffectivePermissionsRequest](#servicessettingsGetEffectivePermissionsRequest) | [GetEffectivePermissionsResponse](#servicessettingsGetEffectivePermissionsResponse) |@perm: Name=GetRoles |
| `ViewAuditLog` | [ViewAuditLogRequest](#servicessettingsViewAuditLogRequest) | [ViewAuditLogResponse](#servicessettingsViewAuditLogResponse) |@perm |
| `GetAllPermissions` | [GetAllPermissionsRequest](#servicessettingsGetAllPermissionsRequest) | [GetAllPermissionsResponse](#servicessettingsGetAllPermissionsResponse) |@perm: Name=Superuser |
| `GetJobLimits` | [GetJobLimitsRequest](#servicessettingsGetJobLimitsRequest) | [GetJobLimitsResponse](#servicessettingsGetJobLimitsResponse) |@perm: Name=Superuser |
| `UpdateJobLimits` | [UpdateJobLimitsRequest](#servicessettingsUpdateJobLimitsRequest) | [UpdateJobLimitsResponse](#servicessettingsUpdateJobLimitsResponse) |@perm: Name=Superuser |
| `DeleteFaction` | [DeleteFactionRequest](#servicessettingsDeleteFactionRequest) | [DeleteFactionResponse](#servicessettingsDeleteFactionResponse) |@perm: Name=Superuser |
| `ListDiscordChannels` | [ListDiscordChannelsRequest](#servicessettingsListDiscordChannelsRequest) | [ListDiscordChannelsResponse](#servicessettingsListDiscordChannelsResponse) |@perm: Name=SetJobProps |
| `ListUserGuilds` | [ListUserGuildsRequest](#servicessettingsListUserGuildsRequest) | [ListUserGuildsResponse](#servicessettingsListUserGuildsResponse) |@perm: Name=SetJobProps |
| `UploadJobLogo` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=SetJobProps buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |
| `DeleteJobLogo` | [DeleteJobLogoRequest](#servicessettingsDeleteJobLogoRequest) | [DeleteJobLogoResponse](#servicessettingsDeleteJobLogoResponse) |@perm: Name=SetJobProps |
| `GetStatus` | [GetStatusRequest](#servicessettingsGetStatusRequest) | [GetStatusResponse](#servicessettingsGetStatusResponse) |@perm: Name=Superuser |

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
| `user_oauth2` | [resources.sync.UserOAuth2Conn](#resourcessyncUserOAuth2Conn) |  |  |
| `dispatch` | [resources.centrum.Dispatch](#resourcescentrumDispatch) |  |  |
| `user_activity` | [resources.users.UserActivity](#resourcesusersUserActivity) |  | User activity |
| `user_props` | [resources.sync.UserProps](#resourcessyncUserProps) |  | Setting props will cause activity to be created automtically |
| `colleague_activity` | [resources.jobs.ColleagueActivity](#resourcesjobsColleagueActivity) |  | Jobs user activity |
| `colleague_props` | [resources.sync.ColleagueProps](#resourcessyncColleagueProps) |  | Setting props will cause activity to be created automtically |
| `job_timeclock` | [resources.sync.TimeclockUpdate](#resourcessyncTimeclockUpdate) |  | Timeclock user entry |
| `user_update` | [resources.sync.UserUpdate](#resourcessyncUserUpdate) |  | User/Char info updates that aren't tracked by activity (yet) |





### services.sync.AddActivityResponse





### services.sync.DeleteDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `users` | [resources.sync.DeleteUsers](#resourcessyncDeleteUsers) |  |  |
| `vehicles` | [resources.sync.DeleteVehicles](#resourcessyncDeleteVehicles) |  |  |





### services.sync.DeleteDataResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |





### services.sync.GetStatusRequest





### services.sync.GetStatusResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.DataStatus](#resourcessyncDataStatus) |  |  |
| `licenses` | [resources.sync.DataStatus](#resourcessyncDataStatus) |  |  |
| `users` | [resources.sync.DataStatus](#resourcessyncDataStatus) |  |  |
| `vehicles` | [resources.sync.DataStatus](#resourcessyncDataStatus) |  |  |





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
| `account_id` | [uint64](#uint64) | optional |  |
| `username` | [string](#string) | optional |  |





### services.sync.SendDataRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `jobs` | [resources.sync.DataJobs](#resourcessyncDataJobs) |  |  |
| `licenses` | [resources.sync.DataLicenses](#resourcessyncDataLicenses) |  |  |
| `users` | [resources.sync.DataUsers](#resourcessyncDataUsers) |  |  |
| `vehicles` | [resources.sync.DataVehicles](#resourcessyncDataVehicles) |  |  |
| `user_locations` | [resources.sync.DataUserLocations](#resourcessyncDataUserLocations) |  |  |





### services.sync.SendDataResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `affected_rows` | [int64](#int64) |  |  |





### services.sync.StreamRequest





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
| `props` | [resources.vehicles.VehicleProps](#resourcesvehiclesVehicleProps) |  |  |





### services.vehicles.SetVehiclePropsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `props` | [resources.vehicles.VehicleProps](#resourcesvehiclesVehicleProps) |  |  |
| `reason` | [string](#string) |  | @sanitize |




 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.vehicles.VehiclesService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `ListVehicles` | [ListVehiclesRequest](#servicesvehiclesListVehiclesRequest) | [ListVehiclesResponse](#servicesvehiclesListVehiclesResponse) |@perm: Attrs=Fields/StringList:[]string{"Wanted"} |
| `SetVehicleProps` | [SetVehiclePropsRequest](#servicesvehiclesSetVehiclePropsRequest) | [SetVehiclePropsResponse](#servicesvehiclesSetVehiclePropsResponse) |@perm: Attrs=Fields/StringList:[]string{"Wanted"} |

 <!-- end services -->



## services/wiki/collab.proto

 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


### services.wiki.CollabService

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| `JoinRoom` | [.resources.collab.ClientPacket](#resourcescollabClientPacket) stream | [.resources.collab.ServerPacket](#resourcescollabServerPacket) stream |@perm: Name=wiki.WikiService/UpdatePage buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE buf:lint:ignore RPC_REQUEST_STANDARD_NAME buf:lint:ignore RPC_RESPONSE_STANDARD_NAME |

 <!-- end services -->



## services/wiki/wiki.proto


### services.wiki.CreatePageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `parent_id` | [uint64](#uint64) | optional |  |
| `content_type` | [resources.common.content.ContentType](#resourcescommoncontentContentType) |  |  |





### services.wiki.CreatePageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `job` | [string](#string) |  |  |
| `id` | [uint64](#uint64) |  |  |





### services.wiki.DeletePageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.wiki.DeletePageResponse





### services.wiki.GetPageRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |





### services.wiki.GetPageResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `page` | [resources.wiki.Page](#resourceswikiPage) |  |  |





### services.wiki.ListPageActivityRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationRequest](#resourcescommondatabasePaginationRequest) |  |  |
| `page_id` | [uint64](#uint64) |  |  |





### services.wiki.ListPageActivityResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [resources.common.database.PaginationResponse](#resourcescommondatabasePaginationResponse) |  |  |
| `activity` | [resources.wiki.PageActivity](#resourceswikiPageActivity) | repeated |  |





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
| `ListPages` | [ListPagesRequest](#serviceswikiListPagesRequest) | [ListPagesResponse](#serviceswikiListPagesResponse) |@perm |
| `GetPage` | [GetPageRequest](#serviceswikiGetPageRequest) | [GetPageResponse](#serviceswikiGetPageResponse) |@perm: Name=ListPages |
| `CreatePage` | [CreatePageRequest](#serviceswikiCreatePageRequest) | [CreatePageResponse](#serviceswikiCreatePageResponse) |@perm: Name=UpdatePage |
| `UpdatePage` | [UpdatePageRequest](#serviceswikiUpdatePageRequest) | [UpdatePageResponse](#serviceswikiUpdatePageResponse) |@perm: Attrs=Fields/StringList:[]string{"Public"} |
| `DeletePage` | [DeletePageRequest](#serviceswikiDeletePageRequest) | [DeletePageResponse](#serviceswikiDeletePageResponse) |@perm |
| `ListPageActivity` | [ListPageActivityRequest](#serviceswikiListPageActivityRequest) | [ListPageActivityResponse](#serviceswikiListPageActivityResponse) |@perm |
| `UploadFile` | [.resources.file.UploadFileRequest](#resourcesfileUploadFileRequest) stream | [.resources.file.UploadFileResponse](#resourcesfileUploadFileResponse) |@perm: Name=UpdatePage |

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

