BEGIN;

ALTER TABLE `fivenet_jobs_grades`
  DROP FOREIGN KEY `fk_fivenet_jobs_grades_job_id`;

ALTER TABLE `fivenet_jobs_grades`
  DROP INDEX `idx_fivenet_jobs_grades_job_id`,
  DROP COLUMN `job_id`;

ALTER TABLE `fivenet_jobs`
  DROP INDEX `idx_fivenet_jobs_id`,
  DROP COLUMN `id`;

COMMIT;
