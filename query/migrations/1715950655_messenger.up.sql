BEGIN;

-- Table: fivenet_msgs_threads
CREATE TABLE
    IF NOT EXISTS `fivenet_msgs_threads` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `title` varchar(255) NOT NULL,
        `closed` tinyint(1) DEFAULT 0,
        `creator_job` varchar(50) NOT NULL,
        `creator_id` int(11) NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_msgs_threads_deleted_at` (`deleted_at`),
        CONSTRAINT `fk_fivenet_msgs_threads_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
    ) ENGINE = InnoDB;

-- Table: fivenet_msgs_threads_user_state
CREATE TABLE
    IF NOT EXISTS `fivenet_msgs_threads_user_state` (
        `thread_id` bigint(20) unsigned NOT NULL,
        `user_id` int(11) NOT NULL,
        `last_read` bigint(20) unsigned DEFAULT NULL,
        `important` tinyint(1) DEFAULT 0,
        `favorite` tinyint(1) DEFAULT 0,
        `muted` tinyint(1) DEFAULT 0,
        PRIMARY KEY (`thread_id`, `user_id`),
        CONSTRAINT `fk_fivenet_msgs_threads_user_state_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_msgs_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_msgs_threads_user_state_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

-- Table: fivenet_msgs_threads_job_access
CREATE TABLE IF NOT EXISTS `fivenet_msgs_threads_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `thread_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 1,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_msgs_threads_job_access` (`thread_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_msgs_threads_job_access_thread_id` (`thread_id`),
  CONSTRAINT `fk_fivenet_msgs_threads_job_access_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_msgs_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_msgs_threads_user_access
CREATE TABLE IF NOT EXISTS `fivenet_msgs_threads_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `thread_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_msgs_threads_user_access` (`thread_id`, `user_id`),
  KEY `idx_fivenet_msgs_threads_user_access_thread_id` (`thread_id`),
  KEY `idx_fivenet_msgs_threads_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_msgs_threads_user_access_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_msgs_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_msgs_threads_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_msgs_messages
CREATE TABLE
    IF NOT EXISTS `fivenet_msgs_messages` (
        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
        `thread_id` bigint(20) unsigned NOT NULL,
        `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
        `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
        `deleted_at` datetime(3) DEFAULT NULL,
        `message` varchar(2048) NOT NULL,
        `data` text DEFAULT NULL,
        `creator_id` int(11) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `idx_fivenet_msgs_messages_deleted_at` (`deleted_at`),
        CONSTRAINT `fk_fivenet_msgs_messages_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_msgs_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
        CONSTRAINT `fk_fivenet_msgs_messages_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
    ) ENGINE=InnoDB;

COMMIT;
