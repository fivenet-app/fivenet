BEGIN;

-- Table: `fivenet_centrum_settings` add index for public option
ALTER TABLE `fivenet_centrum_settings` ADD INDEX `idx_fivenet_centrum_settings_public` (`public`);

-- Table: `fivenet_centrum_qualifications_access` drop unused centrum qualifications access table
DROP TABLE IF EXISTS `fivenet_centrum_qualifications_access`;

-- Table: `fivenet_centrum_job_access` replace table with new structure
DROP TABLE IF EXISTS `fivenet_centrum_job_access`;
CREATE TABLE IF NOT EXISTS `fivenet_centrum_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `source_job` varchar(40) NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  `accepted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_job_access` (`source_job`, `job`, `minimum_grade`),
  KEY `idx_fivenet_centrum_job_access_unit_id` (`source_job`),
  KEY `idx_fivenet_centrum_job_access_job` (`job`),
  KEY `idx_fivenet_centrum_job_access_accepted_at` (`accepted_at`)
) ENGINE=InnoDB;

ALTER TABLE `fivenet_centrum_settings` DROP COLUMN `access`;

COMMIT;
