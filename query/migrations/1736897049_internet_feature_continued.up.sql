BEGIN;

ALTER TABLE `fivenet_jobs_timeclock` DROP KEY `idx_fivenet_jobs_timeclock_unique`;
ALTER TABLE `fivenet_jobs_timeclock` ADD UNIQUE KEY `idx_fivenet_jobs_timeclock_unique` (`job`, `user_id`, `date`, `start_time`);

-- User Activity rework
ALTER TABLE `fivenet_user_activity` ADD COLUMN `data` longtext AFTER `reason`;

ALTER TABLE `fivenet_user_activity` DROP COLUMN `key`;
ALTER TABLE `fivenet_user_activity` DROP COLUMN `old_value`;
ALTER TABLE `fivenet_user_activity` DROP COLUMN `new_value`;

COMMIT;
