BEGIN;

ALTER TABLE `fivenet_permissions` ADD COLUMN `order` mediumint(4) DEFAULT 0,
  ADD INDEX `idx_fivenet_permissions_order` (`order`);

COMMIT;
