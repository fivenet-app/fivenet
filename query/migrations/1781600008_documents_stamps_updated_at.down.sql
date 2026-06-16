BEGIN;

ALTER TABLE `fivenet_documents_stamps`
  DROP COLUMN `updated_at`;

COMMIT;
