BEGIN;

-- Table: fivenet_calendar
CREATE TABLE
    IF NOT EXISTS `fivenet_calendar` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `job` varchar(20) DEFAULT NULL,
        `name` varchar(255) NOT NULL,
        `description` varchar(512) DEFAULT NULL,
        `public` tinyint(1) DEFAULT 0,
        `closed` tinyint(1) DEFAULT 0,
        `creator_id` int(11) NULL DEFAULT NULL,
        `creator_job` varchar(50) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_calendar_deleted_at` (`deleted_at`),
        KEY `idx_fivenet_calendar_job` (`job`),
        CONSTRAINT `fk_fivenet_calendar_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_calendar_entries
CREATE TABLE
    IF NOT EXISTS `fivenet_calendar_entries` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `calendar_id` bigint(20) unsigned NOT NULL,
        `job` varchar(20) DEFAULT NULL,
        `start_time` datetime(3) NOT NULL,
        `end_time` datetime(3) DEFAULT NULL,
        `title` varchar(1024) NOT NULL,
        `content` longtext,
        `public` tinyint(1) DEFAULT 0,
        `creator_id` int(11) NULL DEFAULT NULL,
        `creator_job` varchar(50) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_calendar_entries_deleted_at` (`deleted_at`),
        KEY `idx_fivenet_calendar_entries_calendar_id` (`calendar_id`),
        KEY `idx_fivenet_calendar_entries_times` (`start_time`, `end_time`),
        KEY `idx_fivenet_calendar_entries_job` (`job`),
        CONSTRAINT `fk_fivenet_calendar_entries_calendar_id` FOREIGN KEY (`calendar_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_calendar_entries_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_calendar_job_access
CREATE TABLE IF NOT EXISTS `fivenet_calendar_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `calendar_id` bigint(20) unsigned NOT NULL,
  `entry_id` bigint(20) unsigned DEFAULT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 1,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_calendar_job_access` (`calendar_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_calendar_job_access_calendar_id` (`calendar_id`),
  CONSTRAINT `fk_fivenet_calendar_job_access_calendar_id` FOREIGN KEY (`calendar_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_calendar_user_access
CREATE TABLE IF NOT EXISTS `fivenet_calendar_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `calendar_id` bigint(20) unsigned DEFAULT NULL,
  `entry_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_calendar_user_access` (`calendar_id`, `user_id`),
  KEY `idx_fivenet_calendar_user_access_calendar_id` (`calendar_id`),
  KEY `idx_fivenet_calendar_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_calendar_user_access_calendar_id` FOREIGN KEY (`calendar_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_calendar_user_access_entry_id` FOREIGN KEY (`entry_id`) REFERENCES `fivenet_calendar_entries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_calendar_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
