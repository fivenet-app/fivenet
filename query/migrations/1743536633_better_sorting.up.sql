BEGIN;

-- Cries in no natural sorting in MySQL

ALTER TABLE `fivenet_documents_categories` ADD COLUMN `sort_key` VARCHAR(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`name`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`name`, '[0-9]+', LPAD(REGEXP_SUBSTR(`name`, '[0-9]+'), 8, '0'))
      ELSE `name`
    END
  )) STORED AFTER `name`,
  ADD INDEX `idx_fivenet_documents_categories_sort_key` (`sort_key`);

ALTER TABLE `fivenet_job_citizen_labels` ADD COLUMN `sort_key` VARCHAR(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`name`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`name`, '[0-9]+', LPAD(REGEXP_SUBSTR(`name`, '[0-9]+'), 8, '0'))
      ELSE `name`
    END
  )) STORED AFTER `name`,
  ADD INDEX `idx_fivenet_job_citizen_labels_sort_key` (`sort_key`);

ALTER TABLE `fivenet_jobs_labels` ADD COLUMN `sort_key` VARCHAR(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`name`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`name`, '[0-9]+', LPAD(REGEXP_SUBSTR(`name`, '[0-9]+'), 8, '0'))
      ELSE `name`
    END
  )) STORED AFTER `name`,
  ADD INDEX `idx_fivenet_jobs_labels_sort_key` (`sort_key`);

ALTER TABLE `fivenet_lawbooks` ADD COLUMN `sort_key` VARCHAR(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`name`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`name`, '[0-9]+', LPAD(REGEXP_SUBSTR(`name`, '[0-9]+'), 8, '0'))
      ELSE `name`
    END
  )) STORED AFTER `name`,
  ADD INDEX `idx_fivenet_lawbooks_sort_key` (`sort_key`);

ALTER TABLE `fivenet_lawbooks_laws` ADD COLUMN `sort_key` VARCHAR(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`name`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`name`, '[0-9]+', LPAD(REGEXP_SUBSTR(`name`, '[0-9]+'), 8, '0'))
      ELSE `name`
    END
  )) STORED AFTER `name`,
  ADD INDEX `idx_fivenet_lawbooks_laws_sort_key` (`sort_key`);

ALTER TABLE `fivenet_wiki_pages` ADD COLUMN `sort_key` VARCHAR(255) GENERATED ALWAYS AS ((
    CASE
      WHEN (REGEXP_SUBSTR(`title`, '[0-9]+') IS NOT NULL) THEN
        REGEXP_REPLACE(`title`, '[0-9]+', LPAD(REGEXP_SUBSTR(`title`, '[0-9]+'), 8, '0'))
      ELSE `title`
    END
  )) STORED AFTER `title`,
  ADD INDEX `idx_fivenet_wiki_pages_sort_key` (`sort_key`);

COMMIT;
