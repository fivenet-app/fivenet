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
        `expires_at` date DEFAULT NULL,
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
        `date` date DEFAULT (CURRENT_DATE) NOT NULL,
        `start_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `end_time` datetime(3) DEFAULT NULL,
        `spent_time` decimal(10,2) DEFAULT 0.0,
        PRIMARY KEY (`job`, `user_id`, `date`),
        CONSTRAINT `fk_fivenet_jobs_timeclock_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Trigger: `fivenet_jobs_timeclock_spent_time_calc` isn't used anymore in the future
DROP TRIGGER IF EXISTS `fivenet_jobs_timeclock_spent_time_calc`;

-- Table: fivenet_jobs_user_props
CREATE TABLE IF NOT EXISTS `fivenet_jobs_user_props` (
  `user_id` int(11) NOT NULL,
  `job` varchar(20) NOT NULL,
  `absence_begin` date DEFAULT NULL,
  `absence_end` date DEFAULT NULL,
  `note` text DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_jobs_user_props_unique` (`user_id`, `job`),
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

COMMIT;
