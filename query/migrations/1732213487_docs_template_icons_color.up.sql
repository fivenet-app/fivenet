BEGIN;

ALTER TABLE `fivenet_documents_templates` ADD COLUMN `color` char(7) DEFAULT 'blue' NULL AFTER `description`;
ALTER TABLE `fivenet_documents_templates` ADD COLUMN `icon` varchar(128) DEFAULT NULL AFTER `color`;

COMMIT;
