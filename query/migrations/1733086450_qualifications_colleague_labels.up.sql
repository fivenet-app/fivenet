BEGIN;

ALTER TABLE `fivenet_jobs_labels` MODIFY COLUMN `name` varchar(128) NOT NULL;

ALTER TABLE `fivenet_qualifications` ADD COLUMN `label_sync_enabled` tinyint(1) DEFAULT '0';
ALTER TABLE `fivenet_qualifications` ADD COLUMN `label_sync_format` varchar(128) DEFAULT NULL;

COMMIT;
