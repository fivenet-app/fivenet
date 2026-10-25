BEGIN;

ALTER TABLE `fivenet_user_files`
  DROP INDEX `uq_fivenet_user_files_mugshot_user`,
  DROP COLUMN `mugshot_user_id`;

COMMIT;
