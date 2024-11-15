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
  `job` varchar(40) DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `domain` varchar(80) NOT NULL,
  `label` varchar(128) NOT NULL,
  `internal` tinyint(1) DEFAULT 0,
  `signature` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_mailer_emails_domain` (`domain`),
  UNIQUE KEY `idx_fivenet_mailer_emails_domain_email` (`domain`, `email`),
  KEY `idx_fivenet_mailer_emails_job` (`job`)
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

-- Table: fivenet_mailer_templates
CREATE TABLE IF NOT EXISTS `fivenet_mailer_templates` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `email_id` bigint(20) unsigned DEFAULT NULL,
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
        `creator_job` varchar(40) DEFAULT NULL,
        `creator_id`int(11) DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_mailer_threads_deleted_at` (`deleted_at`)
    ) ENGINE = InnoDB;

-- Table: fivenet_mailer_threads_recipients_emails
CREATE TABLE IF NOT EXISTS `fivenet_mailer_threads_recipients_emails` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `thread_id` bigint(20) unsigned NOT NULL,
  `email_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_threads_recipients_emails` (`thread_id`, `email_id`),
  KEY `idx_fivenet_mailer_threads_recipients_emails_thread_id` (`thread_id`),
  KEY `idx_fivenet_mailer_threads_recipients_emails_email_id` (`email_id`),
  CONSTRAINT `fk_fivenet_mailer_threads_recipients_emails_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_threads_recipients_emails_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads_recipients_users
CREATE TABLE IF NOT EXISTS `fivenet_mailer_threads_recipients_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `thread_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_threads_recipients_users` (`thread_id`, `user_id`),
  KEY `idx_fivenet_mailer_threads_recipients_users_thread_id` (`thread_id`),
  KEY `idx_fivenet_mailer_threads_recipients_users_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_mailer_threads_recipients_users_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_threads_recipients_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads_state_email
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_threads_state_email` (
        `thread_id` bigint(20) unsigned NOT NULL,
        `email_id` bigint(20) unsigned NOT NULL,
        `last_read` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `unread` tinyint(1) DEFAULT 1,
        `important` tinyint(1) DEFAULT 0,
        `favorite` tinyint(1) DEFAULT 0,
        `muted` tinyint(1) DEFAULT 0,
        `archived` tinyint(1) DEFAULT 0,
        PRIMARY KEY (`thread_id`, `email_id`),
        CONSTRAINT `fk_fivenet_mailer_threads_state_email_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_threads_state_email_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads_state_user
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_threads_state_user` (
        `thread_id` bigint(20) unsigned NOT NULL,
        `user_id` int(11) NOT NULL,
        `last_read` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `unread` tinyint(1) DEFAULT 1,
        `important` tinyint(1) DEFAULT 0,
        `favorite` tinyint(1) DEFAULT 0,
        `muted` tinyint(1) DEFAULT 0,
        `archived` tinyint(1) DEFAULT 0,
        PRIMARY KEY (`thread_id`, `user_id`),
        CONSTRAINT `fk_fivenet_mailer_threads_state_user_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_threads_state_user_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_messages
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_messages` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `thread_id` bigint(20) unsigned DEFAULT NULL,
        `sender_email_id` bigint(20) unsigned DEFAULT NULL,
        `sender_user_id` int(11) DEFAULT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `title` varchar(255) NOT NULL,
        `content` longtext NOT NULL,
        `data` text DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_fivenet_mailer_messages_thread_id` (`thread_id`),
        KEY `idx_fivenet_mailer_messages_sender_email_id` (`sender_email_id`),
        KEY `idx_fivenet_mailer_messages_sender_user_id` (`sender_user_id`),
        KEY `idx_fivenet_mailer_messages_deleted_at` (`deleted_at`),
        FULLTEXT KEY `idx_fivenet_mailer_messages_title` (`title`),
        FULLTEXT KEY `idx_fivenet_mailer_messages_content` (`content`),
        CONSTRAINT `fk_fivenet_mailer_messages_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_messages_sender_email_id` FOREIGN KEY (`sender_email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_messages_sender_user_id` FOREIGN KEY (`sender_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_settings_blocked
CREATE TABLE IF NOT EXISTS `fivenet_mailer_settings_blocked` (
  `source_email_id` bigint(20) unsigned NOT NULL,
  `target_email_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`source_email_id`, `target_email_id`),
  UNIQUE KEY `idx_fivenet_mailer_settings_blocked` (`source_email_id`, `target_email_id`),
  KEY `idx_fivenet_mailer_settings_blocked_source_email_id` (`source_email_id`),
  KEY `idx_fivenet_mailer_settings_blocked_target_email_id` (`target_email_id`),
  CONSTRAINT `fk_fivenet_mailer_settings_blocked_source_email_id` FOREIGN KEY (`source_email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_settings_blocked_target_email_id` FOREIGN KEY (`target_email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
