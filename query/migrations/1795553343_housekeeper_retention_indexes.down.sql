BEGIN;

ALTER TABLE `fivenet_job_props`
  RENAME INDEX `idx_fivenet_job_props_deleted_at` TO `idx_fivenet_calendar_deleted_at`;

ALTER TABLE `fivenet_centrum_dispatches_heatmaps`
  DROP INDEX `idx_fivenet_centrum_dispatches_heatmaps_generated_at`;

ALTER TABLE `fivenet_centrum_units_status`
  DROP INDEX `idx_fivenet_centrum_units_status_created_at_id`;

ALTER TABLE `fivenet_mailer_emails`
  DROP INDEX `idx_fivenet_mailer_emails_deleted_at_id`;

ALTER TABLE `fivenet_notifications`
  DROP INDEX `idx_fivenet_notifications_created_at_id`;

ALTER TABLE `fivenet_job_groups`
  DROP INDEX `idx_fivenet_job_groups_deleted_at_id`;

ALTER TABLE `fivenet_job_timeclock`
  DROP INDEX `idx_fivenet_job_timeclock_date`;

ALTER TABLE `fivenet_documents_stats_metric`
  DROP INDEX `idx_fivenet_documents_stats_metric_created_at_id`;

ALTER TABLE `fivenet_property_definitions`
  DROP INDEX `idx_fivenet_property_definitions_deleted_at`;

ALTER TABLE `fivenet_wiki_pages`
  DROP INDEX `idx_fivenet_wiki_pages_deleted_at_id`;

COMMIT;
