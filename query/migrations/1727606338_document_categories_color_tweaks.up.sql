BEGIN;

ALTER TABLE `fivenet_documents_categories` MODIFY COLUMN `color` char(7) DEFAULT 'primary' NULL;

COMMIT;
