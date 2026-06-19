BEGIN;

ALTER TABLE `fivenet_documents_stamps`
  ADD COLUMN `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) AFTER `created_at`;

UPDATE `fivenet_documents_stamps`
SET `updated_at` = `created_at`
WHERE `updated_at` IS NULL;

COMMIT;
