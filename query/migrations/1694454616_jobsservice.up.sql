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
        KEY (`type`),
        KEY (`created_at`),
        KEY (`target_user_id`),
        CONSTRAINT `fk_fivenet_jobs_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_jobs_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

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
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

DROP TRIGGER IF EXISTS `fivenet_jobs_timeclock_spent_time_calc`;

-- Trigger: fivenet_jobs_timeclock_spent_time_calc
-- Requires `SUPER` privilege to be created...
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

-- Table: fivenet_jobs_requests_types
CREATE TABLE IF NOT EXISTS `fivenet_jobs_requests_types` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `job` varchar(50) NOT NULL,
  `name` varchar(32) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `weight` int(11) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_jobs_requests_job` (`job`),
  KEY `idx_fivenet_jobs_requests_weight` (`weight`),
  KEY `idx_fivenet_jobs_requests_types_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_jobs_requests
CREATE TABLE IF NOT EXISTS `fivenet_jobs_requests` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `job` varchar(50) NOT NULL,
  `type_id` bigint(20) unsigned DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `message` longtext NOT NULL,
  `status` varchar(24) DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `approved` tinyint(1) DEFAULT NULL,
  `approver_id` int(11) DEFAULT NULL,
  `begins_at` datetime(3) DEFAULT NULL,
  `ends_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_jobs_requests_created_at` (`created_at`),
  KEY `idx_fivenet_jobs_requests_deleted_at` (`deleted_at`),
  KEY `idx_fivenet_jobs_requests_type_id` (`type_id`),
  KEY `idx_fivenet_jobs_requests_creator_id` (`creator_id`),
  FULLTEXT KEY `idx_fivenet_jobs_requests_title` (`title`),
  FULLTEXT KEY `idx_fivenet_jobs_requests_message` (`message`),
  CONSTRAINT `fk_fivenet_jobs_requests_types` FOREIGN KEY (`type_id`) REFERENCES `fivenet_jobs_requests_types` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_fivenet_jobs_requests_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_jobs_requests_comments
CREATE TABLE IF NOT EXISTS `fivenet_jobs_requests_comments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `request_id` bigint(20) unsigned NOT NULL,
  `comment` longtext,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_jobs_requests_comments_request_id` (`request_id`),
  KEY `idx_fivenet_jobs_requests_comments_creator_id` (`creator_id`),
  CONSTRAINT `fk_fivenet_jobs_requests_comments_request_id` FOREIGN KEY (`request_id`) REFERENCES `fivenet_jobs_requests` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_jobs_requests_comments_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_jobs_qualifications
-- TODO

-- Table: fivenet_jobs_qualifications_results
-- TODO

-- Table: fivenet_jobs_qualifications_job_access
-- TODO

-- Table: fivenet_jobs_qualifications_reqs_access
-- TODO

COMMIT;
