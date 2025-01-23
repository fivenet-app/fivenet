# migrations

Uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate) to run migrations if needed.

Migration file naming:

```console
# Update
{timestamp}_{title}.up.sql
# Rollback
{timestamp}_{title}.down.sql
```

(Timestamp is generated using `date +%s` command)

For more information about the migration library, see https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md

## Create Ready-To-Edit Migration files

Tested with GNU Bash (version `5.1.16`):

```console
REASON="yourreasonhere"
FILENAME_PREFIX="$(date +%s)_$REASON"
touch "${FILENAME_PREFIX}."{up,down}.sql
```
Please note that when referencing the `users` table, please use `{{.UsersTableName}}` instead to allow the esx compat mode to work.


## Foreign Keys To ESX Framework Tables

### `users` Table

Might not be fully uptodate and is mainly here to be used to enable someone to switch from ESX Compat Mode `true` to `false`.

```sql
ALTER TABLE `fivenet_calendar` ADD CONSTRAINT `fk_fivenet_calendar_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_entries` ADD CONSTRAINT `fk_fivenet_calendar_entries_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_rsvp` ADD CONSTRAINT `fk_fivenet_calendar_rsvp_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_subs` ADD CONSTRAINT `fk_fivenet_calendar_subs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_user_access` ADD CONSTRAINT `fk_fivenet_calendar_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_dispatches_status` ADD CONSTRAINT `fk_fivenet_centrum_dispatches_status_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_markers` ADD CONSTRAINT `fk_fivenet_centrum_markers_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_units_status` ADD CONSTRAINT `fk_fivenet_centrum_units_status_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_units_status` ADD CONSTRAINT `fk_fivenet_centrum_units_status_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_units_users` ADD CONSTRAINT `fk_fivenet_centrum_units_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_users` ADD CONSTRAINT `fk_fivenet_centrum_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents` ADD CONSTRAINT `fk_fivenet_documents_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_activity` ADD CONSTRAINT `fk_fivenet_documents_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_comments` ADD CONSTRAINT `fk_fivenet_documents_comments_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_pins` ADD CONSTRAINT `fk_fivenet_documents_pins_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_references` ADD CONSTRAINT `fk_fivenet_documents_references_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_relations` ADD CONSTRAINT `fk_fivenet_documents_relations_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_relations` ADD CONSTRAINT `fk_fivenet_documents_relations_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_requests` ADD CONSTRAINT `fk_fivenet_documents_requests_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_user_access` ADD CONSTRAINT `fk_fivenet_documents_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_workflow_users` ADD CONSTRAINT `fk_fivenet_documents_workflow_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_internet_ads` ADD CONSTRAINT `fk_fivenet_internet_ads_approver_id` FOREIGN KEY (`approver_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_internet_ads` ADD CONSTRAINT `fk_fivenet_internet_ads_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_internet_domains` ADD CONSTRAINT `fk_fivenet_internet_domains_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_internet_pages` ADD CONSTRAINT `fk_fivenet_internet_pages_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_conduct` ADD CONSTRAINT `fk_fivenet_jobs_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_conduct` ADD CONSTRAINT `fk_fivenet_jobs_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_labels_users` ADD CONSTRAINT `fk_fivenet_jobs_labels_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_timeclock` ADD CONSTRAINT `fk_fivenet_jobs_timeclock_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_user_activity` ADD CONSTRAINT `fk_fivenet_jobs_user_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_jobs_user_activity` ADD CONSTRAINT `fk_fivenet_jobs_user_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_user_props` ADD CONSTRAINT `fk_fivenet_jobs_user_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_emails` ADD CONSTRAINT `fk_fivenet_mailer_emails_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_emails_user_access` ADD CONSTRAINT `fk_fivenet_mailer_emails_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_templates` ADD CONSTRAINT `fk_fivenet_mailer_templates_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_notifications` ADD CONSTRAINT `fk_fivenet_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications` ADD CONSTRAINT `fk_fivenet_qualifications_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_responses` ADD CONSTRAINT `fivenet_qualifications_exam_responses_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_users` ADD CONSTRAINT `fk_fivenet_qualifications_exam_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_requests` ADD CONSTRAINT `fk_fivenet_qualifications_requests_approver_id` FOREIGN KEY (`approver_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_requests` ADD CONSTRAINT `fk_fivenet_qualifications_requests_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_results` ADD CONSTRAINT `fk_fivenet_qualifications_results_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_results` ADD CONSTRAINT `fk_fivenet_qualifications_results_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_activity` ADD CONSTRAINT `fk_fivenet_user_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_user_activity` ADD CONSTRAINT `fk_fivenet_user_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_citizen_labels` ADD CONSTRAINT `fk_fivenet_user_citizen_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_locations` ADD CONSTRAINT `fk_fivenet_user_locations_identifier` FOREIGN KEY (`identifier`) REFERENCES `users` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_props` ADD CONSTRAINT `fk_fivenet_user_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_page_activity` ADD CONSTRAINT `fk_fivenet_wiki_page_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_wiki_page_user_access` ADD CONSTRAINT `fk_fivenet_wiki_page_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_pages` ADD CONSTRAINT `fk_fivenet_wiki_pages_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
```

#### Drop Foreign Keys

Some foreign keys might not exist in your case due to newer database migrations.

```sql
ALTER TABLE `fivenet_calendar` DROP FOREIGN KEY `fk_fivenet_calendar_creator_id`;
ALTER TABLE `fivenet_calendar_entries` DROP FOREIGN KEY `fk_fivenet_calendar_entries_creator_id`;
ALTER TABLE `fivenet_calendar_rsvp` DROP FOREIGN KEY `fk_fivenet_calendar_rsvp_user_id`;
ALTER TABLE `fivenet_calendar_subs` DROP FOREIGN KEY `fk_fivenet_calendar_subs_user_id`;
ALTER TABLE `fivenet_calendar_user_access` DROP FOREIGN KEY `fk_fivenet_calendar_user_access_user_id`;
ALTER TABLE `fivenet_centrum_dispatches_status` DROP FOREIGN KEY `fk_fivenet_centrum_dispatches_status_user_id`;
ALTER TABLE `fivenet_centrum_markers` DROP FOREIGN KEY `fk_fivenet_centrum_markers_creator_id`;
ALTER TABLE `fivenet_centrum_units_status` DROP FOREIGN KEY `fk_fivenet_centrum_units_status_creator_id`;
ALTER TABLE `fivenet_centrum_units_status` DROP FOREIGN KEY `fk_fivenet_centrum_units_status_user_id`;
ALTER TABLE `fivenet_centrum_units_users` DROP FOREIGN KEY `fk_fivenet_centrum_units_users_user_id`;
ALTER TABLE `fivenet_centrum_users` DROP FOREIGN KEY `fk_fivenet_centrum_users_user_id`;
ALTER TABLE `fivenet_documents` DROP FOREIGN KEY `fk_fivenet_documents_creator_id`;
ALTER TABLE `fivenet_documents_activity` DROP FOREIGN KEY `fk_fivenet_documents_activity_creator_id`;
ALTER TABLE `fivenet_documents_comments` DROP FOREIGN KEY `fk_fivenet_documents_comments_creator_id`;
ALTER TABLE `fivenet_documents_pins` DROP FOREIGN KEY `fk_fivenet_documents_pins_creator_id`;
ALTER TABLE `fivenet_documents_references` DROP FOREIGN KEY `fk_fivenet_documents_references_creator_id`;
ALTER TABLE `fivenet_documents_relations` DROP FOREIGN KEY `fk_fivenet_documents_relations_source_user_id`;
ALTER TABLE `fivenet_documents_relations` DROP FOREIGN KEY `fk_fivenet_documents_relations_target_user_id`;
ALTER TABLE `fivenet_documents_requests` DROP FOREIGN KEY `fk_fivenet_documents_requests_creator_id`;
ALTER TABLE `fivenet_documents_user_access` DROP FOREIGN KEY `fk_fivenet_documents_user_access_user_id`;
ALTER TABLE `fivenet_documents_workflow_users` DROP FOREIGN KEY `fk_fivenet_documents_workflow_users_user_id`;
ALTER TABLE `fivenet_internet_ads` DROP FOREIGN KEY `fk_fivenet_internet_ads_approver_id`;
ALTER TABLE `fivenet_internet_ads` DROP FOREIGN KEY `fk_fivenet_internet_ads_creator_id`;
ALTER TABLE `fivenet_internet_domains` DROP FOREIGN KEY `fk_fivenet_internet_domains_creator_id`;
ALTER TABLE `fivenet_internet_pages` DROP FOREIGN KEY `fk_fivenet_internet_pages_creator_id`;
ALTER TABLE `fivenet_jobs_conduct` DROP FOREIGN KEY `fk_fivenet_jobs_conduct_creator_id`;
ALTER TABLE `fivenet_jobs_conduct` DROP FOREIGN KEY `fk_fivenet_jobs_conduct_target_user_id`;
ALTER TABLE `fivenet_jobs_labels_users` DROP FOREIGN KEY `fk_fivenet_jobs_labels_users_user_id`;
ALTER TABLE `fivenet_jobs_timeclock` DROP FOREIGN KEY `fk_fivenet_jobs_timeclock_user_id`;
ALTER TABLE `fivenet_jobs_user_activity` DROP FOREIGN KEY `fk_fivenet_jobs_user_activity_source_user_id`;
ALTER TABLE `fivenet_jobs_user_activity` DROP FOREIGN KEY `fk_fivenet_jobs_user_activity_target_user_id`;
ALTER TABLE `fivenet_jobs_user_props` DROP FOREIGN KEY `fk_fivenet_jobs_user_props_user_id`;
ALTER TABLE `fivenet_mailer_emails` DROP FOREIGN KEY `fk_fivenet_mailer_emails_user_id`;
ALTER TABLE `fivenet_mailer_emails_user_access` DROP FOREIGN KEY `fk_fivenet_mailer_emails_user_access_user_id`;
ALTER TABLE `fivenet_mailer_templates` DROP FOREIGN KEY `fk_fivenet_mailer_templates_creator_id`;
ALTER TABLE `fivenet_notifications` DROP FOREIGN KEY `fk_fivenet_notifications_user_id`;
ALTER TABLE `fivenet_qualifications` DROP FOREIGN KEY `fk_fivenet_qualifications_creator_id`;
ALTER TABLE `fivenet_qualifications_exam_responses` DROP FOREIGN KEY `fivenet_qualifications_exam_responses_user_id`;
ALTER TABLE `fivenet_qualifications_exam_users` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_users_user_id`;
ALTER TABLE `fivenet_qualifications_requests` DROP FOREIGN KEY `fk_fivenet_qualifications_requests_approver_id`;
ALTER TABLE `fivenet_qualifications_requests` DROP FOREIGN KEY `fk_fivenet_qualifications_requests_user_id`;
ALTER TABLE `fivenet_qualifications_results` DROP FOREIGN KEY `fk_fivenet_qualifications_results_creator_id`;
ALTER TABLE `fivenet_qualifications_results` DROP FOREIGN KEY `fk_fivenet_qualifications_results_user_id`;
ALTER TABLE `fivenet_user_activity` DROP FOREIGN KEY `fk_fivenet_user_activity_source_user_id`;
ALTER TABLE `fivenet_user_activity` DROP FOREIGN KEY `fk_fivenet_user_activity_target_user_id`;
ALTER TABLE `fivenet_user_citizen_labels` DROP FOREIGN KEY `fk_fivenet_user_citizen_labels_user_id`;
ALTER TABLE `fivenet_user_citizen_labels` DROP FOREIGN KEY `fk_fivenet_user_citizen_attributes_user_id`;
ALTER TABLE `fivenet_user_locations` DROP FOREIGN KEY `fk_fivenet_user_locations_identifier`;
ALTER TABLE `fivenet_user_props` DROP FOREIGN KEY `fk_fivenet_user_props_user_id`;
ALTER TABLE `fivenet_wiki_page_activity` DROP FOREIGN KEY `fk_fivenet_wiki_page_activity_creator_id`;
ALTER TABLE `fivenet_wiki_page_user_access` DROP FOREIGN KEY `fk_fivenet_wiki_page_user_access_user_id`;
ALTER TABLE `fivenet_wiki_pages` DROP FOREIGN KEY `fk_fivenet_wiki_pages_creator_id`;
```
