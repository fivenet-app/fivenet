BEGIN;

ALTER TABLE `fivenet_centrum_dispatches` DROP INDEX `idx_fivenet_centrum_dispatches_job`;
ALTER TABLE `fivenet_centrum_dispatches` MODIFY COLUMN `jobs` json NOT NULL;
ALTER TABLE `fivenet_centrum_dispatches` ADD CONSTRAINT `chk_fivenet_centrum_dispatches_jobs` CHECK (json_valid(`jobs`));

-- MariaDB doesn't support JSON multi-valued indexes, so we skip this part
-- It might land in MariaDB 12.2 or later, see https://jira.mariadb.org/browse/MDEV-25848
SET @ver := VERSION();

-- Build either the ALTER or a no-op SELECT
SET @ddl :=
  IF(
    LOCATE('MariaDB', @ver) = 0,
    'ALTER TABLE fivenet_centrum_dispatches
       ADD INDEX idx_jobs_elements
         ((CAST(jobs AS CHAR(255) ARRAY)));',
    CONCAT(
      'SELECT ''Skipping JSON multi-valued index on MariaDB: ',
      @ver,
      ''';'
    )
  );

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

COMMIT;
