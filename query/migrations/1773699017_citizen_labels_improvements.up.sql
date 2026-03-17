BEGIN;

ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) AFTER `id`;
ALTER TABLE `fivenet_user_labels_job` ADD COLUMN `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) AFTER `created_at`;

-- TODO

COMMIT;
