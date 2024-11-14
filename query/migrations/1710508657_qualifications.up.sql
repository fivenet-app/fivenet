BEGIN;

-- Table: fivenet_qualifications
CREATE TABLE
    IF NOT EXISTS `fivenet_qualifications` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `job` varchar(20) NOT NULL,
        `weight` int(11) unsigned DEFAULT 0,
        `closed` tinyint(1) DEFAULT 0,
        `abbreviation` varchar(20) NOT NULL,
        `title` varchar(1024) NOT NULL,
        `description` varchar(512) DEFAULT NULL,
        `content` longtext,
        `creator_id` int(11) NULL DEFAULT NULL,
        `creator_job` varchar(50) NOT NULL,
        `discord_sync_enabled` tinyint(1) DEFAULT 0,
        `discord_settings` longtext,
        `exam_mode` smallint(2) DEFAULT 1,
        `exam_settings` longtext,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_qualifications_deleted_at` (`deleted_at`),
        KEY `idx_fivenet_qualifications_job` (`job`),
        KEY `idx_fivenet_qualifications_weight` (`weight`),
        KEY `idx_fivenet_qualifications_discord_sync_enabled` (`job`, `discord_sync_enabled`),
        CONSTRAINT `fk_fivenet_qualifications_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_qualifications_job_access
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `qualification_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_qualifications_job_access` (`qualification_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_qualifications_job_access_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_qualifications_job_access_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_qualifications_requirements
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_requirements` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `qualification_id` bigint(20) unsigned NOT NULL,
  `target_qualification_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_qualifications_requirements_qualification_id` (`qualification_id`),
  UNIQUE KEY `idx_fivenet_qualifications_requirements_qualification_ids` (`qualification_id`, `target_qualification_id`),
  CONSTRAINT `fk_fivenet_qualifications_requirements_quali_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_requirements_target_quali_id` FOREIGN KEY (`target_qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_qualifications_results
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_results` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `status` smallint(2) DEFAULT 0,
  `score` int(4) DEFAULT NULL,
  `summary` varchar(512) DEFAULT NULL,
  `creator_id` int(11) NULL DEFAULT NULL,
  `creator_job` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_qualifications_results_qualification_id_user_id` (`qualification_id`, `user_id`),
  KEY `idx_fivenet_qualifications_results_status` (`status`),
  CONSTRAINT `fk_fivenet_qualifications_results_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_results_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_results_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_qualifications_requests
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_requests` (
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `user_comment` varchar(512) DEFAULT NULL,
  `status` smallint(2) DEFAULT 0,
  `approved_at` datetime(3) DEFAULT NULL,
  `approver_comment` varchar(255) DEFAULT NULL,
  `approver_id` int(11) NULL DEFAULT NULL,
  `approver_job` varchar(50) DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_qualifications_requests_quali_id_user_id` (`qualification_id`, `user_id`),
  KEY `idx_fivenet_qualifications_requests_status` (`status`),
  KEY `idx_fivenet_qualifications_requests_approved_at` (`approved_at`),
  CONSTRAINT `fk_fivenet_qualifications_requests_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_requests_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_requests_approver_id` FOREIGN KEY (`approver_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
