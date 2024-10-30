BEGIN;

ALTER TABLE `fivenet_lawbooks_laws` ADD `hint` varchar(512) NULL;
ALTER TABLE `fivenet_lawbooks_laws` MODIFY COLUMN `description` varchar(1024) NULL;
ALTER TABLE `fivenet_lawbooks_laws` CHANGE `hint` `hint` varchar(512) NULL AFTER description;

COMMIT;
