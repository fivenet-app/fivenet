BEGIN;

ALTER TABLE `fivenet_documents_categories` ADD COLUMN `color` char(7) DEFAULT NULL;
ALTER TABLE `fivenet_documents_categories` ADD COLUMN `icon` varchar(128) DEFAULT NULL;

COMMIT;
