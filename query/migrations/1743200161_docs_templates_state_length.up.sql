BEGIN;

ALTER TABLE `fivenet_documents_templates` MODIFY COLUMN `state` text NOT NULL;

COMMIT;
