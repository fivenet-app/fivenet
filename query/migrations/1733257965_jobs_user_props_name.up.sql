BEGIN;

ALTER TABLE `fivenet_jobs_user_props` ADD COLUMN `name_prefix` varchar(24) NULL;
ALTER TABLE `fivenet_jobs_user_props` ADD KEY `idx_fivenet_jobs_user_props_job_name_prefix` (`job`, `name_prefix`);

ALTER TABLE `fivenet_jobs_user_props` ADD COLUMN `name_suffix` varchar(24) NULL;
ALTER TABLE `fivenet_jobs_user_props` ADD KEY `idx_fivenet_jobs_user_props_job_name_suffix` (`job`, `name_suffix`);

COMMIT;
