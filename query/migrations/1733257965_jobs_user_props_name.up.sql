BEGIN;

ALTER TABLE `fivenet_jobs_user_props` ADD COLUMN `name_prefix` varchar(24) NULL;
ALTER TABLE `fivenet_jobs_user_props` ADD COLUMN `name_suffix` varchar(24) NULL;

COMMIT;
