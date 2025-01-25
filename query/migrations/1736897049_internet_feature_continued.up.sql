BEGIN;

-- User Activity rework
ALTER TABLE `fivenet_user_activity` ADD COLUMN `data` longtext AFTER `reason`;

ALTER TABLE `fivenet_user_activity` DROP COLUMN `key`;
ALTER TABLE `fivenet_user_activity` DROP COLUMN `old_value`;
ALTER TABLE `fivenet_user_activity` DROP COLUMN `new_value`;

COMMIT;
