BEGIN;

ALTER TABLE `fivenet_rbac_permissions` ADD COLUMN `icon` VARCHAR(128) NULL DEFAULT NULL AFTER `order`;

COMMIT;
