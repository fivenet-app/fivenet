BEGIN;

ALTER TABLE `fivenet_documents_templates` ADD COLUMN `approval` mediumtext AFTER `workflow`;

COMMIT;
