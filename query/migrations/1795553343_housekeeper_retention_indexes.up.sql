BEGIN;

-- Support oldest-first housekeeper retention deletes.
ALTER TABLE `fivenet_job_props`
  RENAME INDEX `idx_fivenet_calendar_deleted_at` TO `idx_fivenet_job_props_deleted_at`;

ALTER TABLE `fivenet_wiki_pages`
  ADD INDEX `idx_fivenet_wiki_pages_deleted_at_id` (`deleted_at`, `id`);

ALTER TABLE `fivenet_property_definitions`
  ADD INDEX `idx_fivenet_property_definitions_deleted_at` (`deleted_at`);

ALTER TABLE `fivenet_documents_stats_metric`
  ADD INDEX `idx_fivenet_documents_stats_metric_created_at_id` (`created_at`, `id`);

ALTER TABLE `fivenet_job_timeclock`
  ADD INDEX `idx_fivenet_job_timeclock_date` (`date`);

ALTER TABLE `fivenet_job_groups`
  ADD INDEX `idx_fivenet_job_groups_deleted_at_id` (`deleted_at`, `id`);

ALTER TABLE `fivenet_notifications`
  ADD INDEX `idx_fivenet_notifications_created_at_id` (`created_at`, `id`);

ALTER TABLE `fivenet_mailer_emails`
  ADD INDEX `idx_fivenet_mailer_emails_deleted_at_id` (`deleted_at`, `id`);

ALTER TABLE `fivenet_centrum_units_status`
  ADD INDEX `idx_fivenet_centrum_units_status_created_at_id` (`created_at`, `id`);

ALTER TABLE `fivenet_centrum_dispatches_heatmaps`
  ADD INDEX `idx_fivenet_centrum_dispatches_heatmaps_generated_at` (`generated_at`);

COMMIT;
