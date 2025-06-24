BEGIN;

ALTER TABLE `fivenet_centrum_dispatches` DROP INDEX `idx_fivenet_centrum_dispatches_job`;
ALTER TABLE `fivenet_centrum_dispatches` MODIFY COLUMN `jobs` json NOT NULL;
ALTER TABLE `fivenet_centrum_dispatches` ADD CONSTRAINT `chk_fivenet_centrum_dispatches_jobs` CHECK (json_valid(`jobs`));
ALTER TABLE `fivenet_centrum_dispatches` ADD INDEX `idx_jobs` ((CAST(`jobs` AS CHAR(255) ARRAY)));

COMMIT;
