BEGIN;

-- Soft delete and notes columns for accounts
ALTER TABLE `fivenet_accounts` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `updated_at`;
ALTER TABLE `fivenet_accounts` ADD KEY `idx_deleted_at` (`deleted_at`);
ALTER TABLE `fivenet_accounts` ADD COLUMN `notes` varchar(512) DEFAULT NULL AFTER `last_char`;

-- Soft delete for lawbooks and laws
ALTER TABLE `fivenet_lawbooks` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `updated_at`;
ALTER TABLE `fivenet_lawbooks` ADD KEY `idx_deleted_at` (`deleted_at`);
ALTER TABLE `fivenet_lawbooks_laws` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `updated_at`;
ALTER TABLE `fivenet_lawbooks_laws` ADD KEY `idx_deleted_at` (`deleted_at`);

-- Rename label_id to be target_id column
ALTER TABLE `fivenet_user_labels_job_job_access` CHANGE `label_id` `target_id` bigint unsigned NOT NULL;

COMMIT;
