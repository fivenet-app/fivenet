BEGIN;

-- Remove old messenger tables
DROP TABLE IF EXISTS `fivenet_msgs_messages`;
DROP TABLE IF EXISTS `fivenet_msgs_settings_blocks`;
DROP TABLE IF EXISTS `fivenet_msgs_threads`;
DROP TABLE IF EXISTS `fivenet_msgs_threads_user_access`;
DROP TABLE IF EXISTS `fivenet_msgs_threads_user_state`;

-- Table: fivenet_mailer_addresses
CREATE TABLE IF NOT EXISTS `fivenet_mailer_addresses` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `domain` varchar(80) NOT NULL,
  `email` varchar(50) DEFAULT NULL,
  `signature` varchar(1024) DEFAULT NULL,
  `job` varchar(40) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_mailer_addresses_domain` (`domain`),
  UNIQUE KEY `idx_fivenet_mailer_addresses_domain_email` (`domain`, `email`),
  KEY `idx_fivenet_mailer_addresses_job` (`job`),
  KEY `idx_fivenet_mailer_addresses_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_mailer_addresses_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_addresses_job_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_addresses_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `address_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_addresses_job_access` (`address_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_mailer_addresses_job_access_address_id` (`address_id`),
  CONSTRAINT `fk_fivenet_mailer_addresses_job_access_address_id` FOREIGN KEY (`address_id`) REFERENCES `fivenet_mailer_addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_addresses_user_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_addresses_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `address_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_addresses_user_access` (`address_id`, `user_id`),
  KEY `idx_fivenet_mailer_addresses_user_access_address_id` (`address_id`),
  KEY `idx_fivenet_mailer_addresses_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_mailer_addresses_user_access_address_id` FOREIGN KEY (`address_id`) REFERENCES `fivenet_mailer_addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_addresses_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_threads` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `title` varchar(255) NOT NULL,
        `address_id` bigint(20) unsigned,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_mailer_threads_deleted_at` (`deleted_at`),
        CONSTRAINT `fk_fivenet_mailer_threads_address_id` FOREIGN KEY (`address_id`) REFERENCES `fivenet_mailer_addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_mailer_threads_job_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_threads_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `address_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_threads_job_access` (`address_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_mailer_threads_job_access_address_id` (`address_id`),
  CONSTRAINT `fk_fivenet_mailer_threads_job_access_address_id` FOREIGN KEY (`address_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads_user_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_threads_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `address_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_threads_user_access` (`address_id`, `user_id`),
  KEY `idx_fivenet_mailer_threads_user_access_address_id` (`address_id`),
  KEY `idx_fivenet_mailer_threads_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_mailer_threads_user_access_address_id` FOREIGN KEY (`address_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_threads_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_mailer_threads_user_state
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_threads_user_state` (
        `thread_id` bigint(20) unsigned NOT NULL,
        `user_id` int(11) NOT NULL,
        `last_read` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `unread` tinyint(1) DEFAULT 0,
        `important` tinyint(1) DEFAULT 0,
        `favorite` tinyint(1) DEFAULT 0,
        `muted` tinyint(1) DEFAULT 0,
        `archived` tinyint(1) DEFAULT '0',
        PRIMARY KEY (`thread_id`, `user_id`),
        CONSTRAINT `fk_fivenet_mailer_threads_user_state_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_threads_user_state_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_messages
CREATE TABLE
    IF NOT EXISTS `fivenet_mailer_messages` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `thread_id` bigint(20) unsigned DEFAULT NULL,
        `sender_id` bigint(20) unsigned NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `title` varchar(255) NOT NULL,
        `content` longtext NOT NULL,
        `data` text DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_mailer_messages_deleted_at` (`deleted_at`),
        FULLTEXT KEY `idx_fivenet_mailer_messages_title` (`title`),
        FULLTEXT KEY `idx_fivenet_mailer_messages_content` (`content`),
        UNIQUE KEY `idx_fivenet_mailer_messages_thread_id` (`thread_id`),
        UNIQUE KEY `idx_fivenet_mailer_messages_sender_id` (`sender_id`),
        CONSTRAINT `fk_fivenet_mailer_messages_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_mailer_messages_sender_id` FOREIGN KEY (`sender_id`) REFERENCES `fivenet_mailer_addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_mailer_settings_blocked
CREATE TABLE IF NOT EXISTS `fivenet_mailer_settings_blocked` (
  `source_address_id` bigint(20) unsigned NOT NULL,
  `target_address_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`source_address_id`, `target_address_id`),
  UNIQUE KEY `idx_fivenet_mailer_settings_blocked` (`source_address_id`, `target_address_id`),
  KEY `idx_fivenet_mailer_settings_blocked_source_address_id` (`source_address_id`),
  KEY `idx_fivenet_mailer_settings_blocked_target_address_id` (`target_address_id`),
  CONSTRAINT `fk_fivenet_mailer_settings_blocked_source_address_id` FOREIGN KEY (`source_address_id`) REFERENCES `fivenet_mailer_addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_settings_blocked_target_address_id` FOREIGN KEY (`target_address_id`) REFERENCES `fivenet_mailer_addresses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
