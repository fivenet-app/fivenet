BEGIN;

RENAME TABLE `fivenet_centrum_disponents` TO `fivenet_centrum_dispatchers`;
ALTER TABLE `fivenet_centrum_dispatchers` DROP FOREIGN KEY `fk_fivenet_centrum_disponents_user_id`;
ALTER TABLE `fivenet_centrum_dispatchers` ADD CONSTRAINT `fk_fivenet_centrum_dispatchers_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_centrum_settings`
ALTER TABLE `fivenet_centrum_settings` ADD COLUMN `type` mediumint(2) DEFAULT 0 NULL AFTER `enabled`;
ALTER TABLE `fivenet_centrum_settings` ADD COLUMN `public` tinyint(1) DEFAULT 0 NOT NULL AFTER `type`;

-- Table: fivenet_centrum_job_access
CREATE TABLE IF NOT EXISTS `fivenet_centrum_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `unit_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_job_access` (`unit_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_centrum_job_access_unit_id` (`unit_id`),
  CONSTRAINT `fk_fivenet_centrum_job_access_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_centrum_qualifications_access
CREATE TABLE IF NOT EXISTS `fivenet_centrum_qualifications_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `unit_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_qualifications_access` (`unit_id`, `qualification_id`),
  KEY `idx_fivenet_centrum_qualifications_access_unit_id` (`unit_id`),
  KEY `idx_fivenet_centrum_qualifications_access_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_centrum_qualifications_access_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_centrum_qualifications_access_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_centrum_user_locations
ALTER TABLE `fivenet_centrum_user_locations` ADD COLUMN `job_grade` int(11) NULL AFTER `job`;
ALTER TABLE `fivenet_centrum_user_locations` MODIFY COLUMN `job` varchar(20) NULL;
ALTER TABLE `fivenet_centrum_user_locations` ADD COLUMN `on_duty` tinyint(1) DEFAULT 0 AFTER `hidden`;
ALTER TABLE `fivenet_centrum_user_locations` ADD COLUMN `data` text NULL AFTER `updated_at`;

-- Table: `fivenet_centrum_dispatches` - Generated coord
ALTER TABLE `fivenet_centrum_dispatches` MODIFY COLUMN `x` decimal(24,14) NOT NULL;
ALTER TABLE `fivenet_centrum_dispatches` MODIFY COLUMN `y` decimal(24,14) NOT NULL;

-- Table: `fivenet_centrum_dispatches_heatmaps`
CREATE TABLE IF NOT EXISTS `fivenet_centrum_dispatches_heatmaps` (
  `job` VARCHAR(20) PRIMARY KEY,
  `heat_json` JSON NOT NULL,
  `max` INT NOT NULL DEFAULT 0,
  `generated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

COMMIT;
