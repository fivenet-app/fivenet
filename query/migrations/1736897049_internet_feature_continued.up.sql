BEGIN;

-- User Activity rework
ALTER TABLE `fivenet_user_activity` ADD COLUMN `data` longtext AFTER `reason`;
ALTER TABLE `fivenet_user_activity` MODIFY COLUMN `key` varchar(64) NULL;

-- Internet Feature changes
ALTER TABLE `fivenet_internet_domains` ADD online tinyint(1) DEFAULT 1 NULL AFTER `deleted_at`;

COMMIT;
