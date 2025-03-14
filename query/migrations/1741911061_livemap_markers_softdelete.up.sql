BEGIN;

ALTER TABLE `fivenet_centrum_markers` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `created_at`;
ALTER TABLE `fivenet_centrum_markers` ADD COLUMN `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) DEFAULT NULL AFTER `created_at`;
ALTER TABLE `fivenet_centrum_markers` ADD KEY `idx_fivenet_centrum_markers_deleted_at` (`deleted_at`);

COMMIT;
