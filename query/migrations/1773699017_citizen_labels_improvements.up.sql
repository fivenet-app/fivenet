BEGIN;

ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) AFTER `id`;
ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) AFTER `created_at`;
ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `updated_at`;
ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `icon` varchar(128) DEFAULT NULL AFTER `color`;
ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `settings` varchar(255) DEFAULT NULL AFTER `icon`;

-- Moved from citizens.CitizensService to citizens.LabelsService
DELETE FROM `fivenet_rbac_permissions` WHERE `category` = 'citizens.CitizensService' AND `name` = 'ManageLabels' LIMIT 1;

-- Add new fields for label expiration to "mapping" table
ALTER TABLE `fivenet_user_labels` ADD COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) AFTER `label_id`;
ALTER TABLE `fivenet_user_labels` ADD COLUMN `expires_at` datetime(3) DEFAULT NULL AFTER `created_at`;
ALTER TABLE `fivenet_user_labels` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `expires_at`;
ALTER TABLE `fivenet_user_labels` ADD COLUMN `reason` varchar(255) DEFAULT NULL AFTER `deleted_at`;

-- Add icon field to (colleague) job labels
ALTER TABLE `fivenet_job_labels` ADD COLUMN `icon` varchar(128) DEFAULT NULL AFTER `color`;

COMMIT;
