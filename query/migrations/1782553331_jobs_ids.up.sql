BEGIN;

-- Table: fivenet_jobs - Add id column
ALTER TABLE `fivenet_jobs`
  ADD COLUMN `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT FIRST,
  ADD UNIQUE KEY `idx_fivenet_jobs_id` (`id`);

-- Table: fivenet_jobs_grades - Add job_id column
ALTER TABLE `fivenet_jobs_grades`
  ADD COLUMN `job_id` bigint(20) unsigned DEFAULT NULL FIRST,
  ADD KEY `idx_fivenet_jobs_grades_job_id` (`job_id`);

UPDATE `fivenet_jobs_grades` AS `jg`
INNER JOIN `fivenet_jobs` AS `j`
  ON `jg`.`job_name` = `j`.`name`
SET `jg`.`job_id` = `j`.`id`
WHERE `jg`.`job_id` IS NULL;

ALTER TABLE `fivenet_jobs_grades`
  MODIFY COLUMN `job_id` bigint(20) unsigned NOT NULL,
  ADD CONSTRAINT `fk_fivenet_jobs_grades_job_id` FOREIGN KEY (`job_id`) REFERENCES `fivenet_jobs` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

COMMIT;
