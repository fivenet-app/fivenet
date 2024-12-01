BEGIN;

ALTER TABLE `fivenet_jobs_labels` MODIFY COLUMN `name` varchar(64) NOT NULL;
ALTER TABLE `fivenet_job_citizen_attributes` MODIFY COLUMN `name` varchar(64) NOT NULL;

ALTER TABLE `fivenet_qualifications` ADD COLUMN `label_sync_enabled` tinyint(1) DEFAULT '0';
ALTER TABLE `fivenet_qualifications` ADD COLUMN `label_sync_format` varchar(64) DEFAULT NULL;

COMMIT;
