BEGIN;

-- Remove old messenger tables
DROP TABLE IF EXISTS `fivenet_msgs_messages`;
DROP TABLE IF EXISTS `fivenet_msgs_settings_blocks`;
DROP TABLE IF EXISTS `fivenet_msgs_threads_user_access`;
DROP TABLE IF EXISTS `fivenet_msgs_threads_user_state`;
DROP TABLE IF EXISTS `fivenet_msgs_threads`;

-- Table: fivenet_mailer_emails
CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `disabled` tinyint(1) DEFAULT 0,
  `job` varchar(40) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `email` varchar(80) DEFAULT NULL,
  `label` varchar(128) DEFAULT NULL,
  `internal` tinyint(1) DEFAULT 0,
  `signature` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_email` (`email`),
  KEY `idx_fivenet_mailer_emails_job` (`job`),
  CONSTRAINT `fk_fivenet_mailer_emails_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_emails_job_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `email_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_job_access` (`email_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_mailer_emails_job_access_email_id` (`email_id`),
  CONSTRAINT `fk_fivenet_mailer_emails_job_access_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_emails_user_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `email_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_user_access` (`email_id`, `user_id`),
  KEY `idx_fivenet_mailer_emails_user_access_email_id` (`email_id`),
  KEY `idx_fivenet_mailer_emails_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_mailer_emails_user_access_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_emails_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_emails_qualifications_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_qualifications_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `email_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_qualifications_access` (`email_id`, `qualification_id`),
  KEY `idx_fivenet_mailer_emails_qualifications_access_email_id` (`email_id`),
  KEY `idx_fivenet_mailer_emails_qualifications_access_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_mailer_emails_qualifications_access_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_emails_qualifications_access_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_templates
CREATE TABLE IF NOT EXISTS `fivenet_mailer_templates` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `email_id` bigint(20) unsigned DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `creator_job` varchar(40) DEFAULT NULL,
  `creator_id`int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_mailer_templates_creator_id_creator_job` (`email_id`, `creator_id`, `creator_job`),
  CONSTRAINT `fk_fivenet_mailer_templates_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_templates_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_threads` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `creator_email_id` bigint(20) unsigned NOT NULL,
        `creator_id` int(11) DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_mailer_threads_deleted_at` (`deleted_at`),
        CONSTRAINT `fk_fivenet_mailer_threads_creator_email_id` FOREIGN KEY (`creator_email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_mailer_threads_recipients_emails
CREATE TABLE IF NOT EXISTS `fivenet_mailer_threads_recipients` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `thread_id` bigint(20) unsigned NOT NULL,
  `email_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_threads_recipients` (`thread_id`, `email_id`),
  KEY `idx_fivenet_mailer_threads_recipients_thread_id` (`thread_id`),
  KEY `idx_fivenet_mailer_threads_recipients_email_id` (`email_id`),
  CONSTRAINT `fk_fivenet_mailer_threads_recipients_emails_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_threads_recipients_emails_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads_state
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_threads_state` (
        `thread_id` bigint(20) unsigned NOT NULL,
        `email_id` bigint(20) unsigned NOT NULL,
        `last_read` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `unread` tinyint(1) DEFAULT 1,
        `important` tinyint(1) DEFAULT 0,
        `favorite` tinyint(1) DEFAULT 0,
        `muted` tinyint(1) DEFAULT 0,
        `archived` tinyint(1) DEFAULT 0,
        PRIMARY KEY (`thread_id`, `email_id`),
        CONSTRAINT `fk_fivenet_mailer_threads_state_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_threads_state_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_messages
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_messages` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `thread_id` bigint(20) unsigned DEFAULT NULL,
        `sender_id` bigint(20) unsigned DEFAULT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `title` varchar(255) NOT NULL,
        `content` longtext NOT NULL,
        `data` text DEFAULT NULL,
        `creator_id` int(11) DEFAULT NULL,
        `creator_job` varchar(40) DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_fivenet_mailer_messages_thread_id` (`thread_id`),
        KEY `idx_fivenet_mailer_messages_sender_id` (`sender_id`),
        KEY `idx_fivenet_mailer_messages_deleted_at` (`deleted_at`),
        KEY `idx_fivenet_mailer_messages_creator_id` (`creator_id`),
        FULLTEXT KEY `idx_fivenet_mailer_messages_title` (`title`),
        FULLTEXT KEY `idx_fivenet_mailer_messages_content` (`content`),
        CONSTRAINT `fk_fivenet_mailer_messages_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_messages_sender_id` FOREIGN KEY (`sender_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_settings_blocked
CREATE TABLE IF NOT EXISTS `fivenet_mailer_settings_blocked` (
  `email_id` bigint(20) unsigned NOT NULL,
  `target_email` varchar(80) NOT NULL,
  PRIMARY KEY (`email_id`, `target_email`),
  UNIQUE KEY `idx_fivenet_mailer_settings_blocked` (`email_id`, `target_email`),
  KEY `idx_fivenet_mailer_settings_blocked_email_id` (`email_id`),
  CONSTRAINT `fk_fivenet_mailer_settings_blocked_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
