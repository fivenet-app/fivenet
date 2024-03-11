BEGIN;

-- Table: fivenet_jobs_conduct
CREATE TABLE
    IF NOT EXISTS `fivenet_jobs_conduct` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `job` varchar(20) NOT NULL,
        `type` smallint(2) NOT NULL,
        `message` longtext,
        `expires_at` datetime(3) DEFAULT NULL,
        `target_user_id` int(11) NULL DEFAULT NULL,
        `creator_id` int(11) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `fivenet_jobs_conduct_type` (`type`),
        KEY `fivenet_jobs_conduct_created_at` (`created_at`),
        KEY `fivenet_jobs_conduct_target_user_id` (`target_user_id`),
        CONSTRAINT `fk_fivenet_jobs_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_jobs_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_jobs_timeclock
CREATE TABLE
    IF NOT EXISTS `fivenet_jobs_timeclock` (
        `job` varchar(20) NOT NULL,
        `user_id` int(11) NOT NULL,
        `date` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
        `start_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `end_time` datetime(3) DEFAULT NULL,
        `spent_time` decimal(10,2) DEFAULT 0.0,
        PRIMARY KEY (`job`, `user_id`, `date`),
        CONSTRAINT `fk_fivenet_jobs_timeclock_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

DROP TRIGGER IF EXISTS `fivenet_jobs_timeclock_spent_time_calc`;

-- Trigger: fivenet_jobs_timeclock_spent_time_calc
-- Requires `SUPER` privilege to be created...
/*
CREATE TRIGGER `fivenet_jobs_timeclock_spent_time_calc` BEFORE UPDATE ON `fivenet_jobs_timeclock`
    FOR EACH ROW BEGIN
        DECLARE `duration` DECIMAL(10,2);

      IF (NEW.`start_time` IS NOT NULL AND NEW.`end_time` IS NOT NULL) THEN
          SELECT CAST((TIMESTAMPDIFF(SECOND, NEW.`start_time`, NEW.`end_time`) / 3600) AS DECIMAL(10,2))
              INTO `duration`;

        SET NEW.`spent_time` = (OLD.`spent_time` + `duration`);
        SET NEW.`start_time` = NULL;
        SET NEW.`end_time` = NULL;
    END IF;
END;
*/

-- Table: fivenet_jobs_user_props
CREATE TABLE IF NOT EXISTS `fivenet_jobs_user_props` (
  `user_id` int(11) NOT NULL,
  `absence_begin` date DEFAULT NULL,
  `absence_end` date DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_jobs_user_props_unique` (`user_id`),
  CONSTRAINT `fk_fivenet_jobs_user_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_jobs_user_activity
CREATE TABLE IF NOT EXISTS `fivenet_jobs_user_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `job` varchar(20) NOT NULL,
  `source_user_id` int(11) DEFAULT NULL,
  `target_user_id` int(11) NOT NULL,
  `activity_type` smallint(2) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `data` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_jobs_user_activity_job` (`job`),
  KEY `idx_fivenet_jobs_user_activity_source_user_id` (`source_user_id`),
  KEY `idx_fivenet_jobs_user_activity_target_user_id` (`target_user_id`),
  KEY `idx_fivenet_jobs_user_activity_activity_type` (`activity_type`),
  CONSTRAINT `fk_fivenet_jobs_user_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_fivenet_jobs_user_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_jobs_qualifications
CREATE TABLE
    IF NOT EXISTS `fivenet_jobs_qualifications` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `job` varchar(20) NOT NULL,
        `weight` int(11) unsigned DEFAULT 0,
        `closed` tinyint(1) DEFAULT 0,
        `abbreviation` varchar(20) NOT NULL,
        `title` longtext,
        `description` longtext,
        `creator_id` int(11) NULL DEFAULT NULL,
        `creator_job` varchar(50) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_jobs_qualifications_deleted_at` (`deleted_at`),
        KEY `idx_fivenet_jobs_qualifications_job` (`job`),
        KEY `idx_fivenet_jobs_qualifications_weight` (`weight`),
        CONSTRAINT `fk_fivenet_jobs_qualifications_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_jobs_qualifications_job_access
CREATE TABLE IF NOT EXISTS `fivenet_jobs_qualifications_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `qualification_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 1,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_jobs_qualifications_job_access` (`qualification_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_jobs_qualifications_job_access_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_jobs_qualifications_job_access_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_jobs_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_jobs_qualifications_reqs_access
CREATE TABLE IF NOT EXISTS `fivenet_jobs_qualifications_reqs_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `qualification_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_jobs_qualifications_reqs_access_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_jobs_qualifications_reqs_access_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_jobs_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_jobs_qualifications_results
CREATE TABLE IF NOT EXISTS `fivenet_jobs_qualifications_results` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `qualification_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `successful` tinyint(1) DEFAULT 0,
  `score` int(4) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `creator_id` int(11) NULL DEFAULT NULL,
  `creator_job` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_jobs_qualifications_results_qualification_id_user_id` (`qualification_id`, `user_id`),
  KEY `idx_fivenet_jobs_qualifications_results_successful` (`successful`),
  CONSTRAINT `fk_fivenet_jobs_qualifications_results_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_jobs_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_jobs_qualifications_results_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_jobs_qualifications_results_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
