BEGIN;

ALTER TABLE `fivenet_job_props` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `updated_at`,
  ADD KEY `idx_fivenet_calendar_deleted_at` (`deleted_at`);

COMMIT;
