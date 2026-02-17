BEGIN;

ALTER TABLE `fivenet_user` MODIFY COLUMN `height` DECIMAL(5,2) NULL;

ALTER TABLE `fivenet_user` ADD INDEX `idx_height` (`height`);

COMMIT;
