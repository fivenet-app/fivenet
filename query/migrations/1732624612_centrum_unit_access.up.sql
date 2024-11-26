BEGIN;

-- Table: fivenet_centrum_units_job_access
CREATE TABLE IF NOT EXISTS `fivenet_centrum_units_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `unit_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_units_job_access` (`unit_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_centrum_units_job_access_unit_id` (`unit_id`),
  CONSTRAINT `fk_fivenet_centrum_units_job_access_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_centrum_units_qualifications_access
CREATE TABLE IF NOT EXISTS `fivenet_centrum_units_qualifications_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `unit_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_units_qualifications_access` (`unit_id`, `qualification_id`),
  KEY `idx_fivenet_centrum_units_qualifications_access_unit_id` (`unit_id`),
  KEY `idx_fivenet_centrum_units_qualifications_access_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_centrum_units_qualifications_access_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_centrum_units_qualifications_access_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

ALTER TABLE `fivenet_centrum_units` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL;
ALTER TABLE `fivenet_centrum_units` ADD KEY `idx_fivenet_centrum_units_deleted_at` (`deleted_at`);
ALTER TABLE `fivenet_centrum_units` CHANGE `deleted_at` `deleted_at` datetime(3) NULL AFTER `updated_at`;

COMMIT;
