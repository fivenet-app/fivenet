BEGIN;

ALTER TABLE `fivenet_user` ADD COLUMN `license` VARCHAR(64) NULL AFTER `account_id`;
ALTER TABLE `fivenet_user` ADD INDEX `idx_license` (`id`, `license`);

COMMIT;
