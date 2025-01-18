BEGIN;

-- User Activity rework
ALTER TABLE `fivenet_user_activity` ADD COLUMN `data` longtext AFTER `reason`;

ALTER TABLE `fivenet_user_activity` DROP COLUMN `key`;
ALTER TABLE `fivenet_user_activity` DROP COLUMN `old_value`;
ALTER TABLE `fivenet_user_activity` DROP COLUMN `new_value`;

-- Internet Feature changes
ALTER TABLE `fivenet_internet_domains` ADD online tinyint(1) DEFAULT 1 NULL AFTER `deleted_at`;

COMMIT;
