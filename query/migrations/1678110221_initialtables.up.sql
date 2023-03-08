BEGIN;
-- Table: arpanet_accounts
CREATE TABLE IF NOT EXISTS `arpanet_accounts` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `enabled` tinyint(1) DEFAULT 0,
  `username` varchar(24) NOT NULL,
  `password` varchar(64) NOT NULL,
  `license` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_accounts_username` (`username`),
  KEY `idx_arpanet_accounts_license` (`license`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents
CREATE TABLE IF NOT EXISTS `arpanet_documents` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` longtext NOT NULL,
  `content` longtext NOT NULL,
  `content_type` varchar(24) NOT NULL,
  `closed` tinyint(1) DEFAULT 0,
  `state` varchar(24) NOT NULL,
  `creator_id` int(11) NOT NULL,
  `creator_job` varchar(20) NOT NULL,
  `public` tinyint(1) NOT NULL DEFAULT 0,
  `response_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_creator_id` (`creator_id`),
  KEY `idx_arpanet_documents_creator_job` (`creator_job`),
  KEY `idx_arpanet_documents_response_id` (`response_id`),
  CONSTRAINT `fk_arpanet_documents_responses` FOREIGN KEY (`response_id`) REFERENCES `arpanet_documents` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_job_access
CREATE TABLE IF NOT EXISTS `arpanet_documents_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `document_id` bigint(20) unsigned NOT NULL,
  `name` varchar(20) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` varchar(12) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_job_access_document_id` (`document_id`),
  CONSTRAINT `fk_arpanet_documents_job_access` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`),
  CONSTRAINT `fk_arpanet_documents_jobs` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- arpanet_documents_mentions
CREATE TABLE IF NOT EXISTS `arpanet_documents_mentions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `document_id` bigint(20) unsigned DEFAULT NULL,
  `identifier` varchar(64) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_mentions_document_id` (`document_id`),
  KEY `idx_arpanet_documents_mentions_identifier` (`identifier`),
  KEY `idx_arpanet_documents_mentions_user_id` (`user_id`),
  CONSTRAINT `fk_arpanet_documents_mentions` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_templates
CREATE TABLE `arpanet_documents_templates` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE current_timestamp(3),
  `job` varchar(20) NOT NULL,
  `job_grade` int(11) NOT NULL DEFAULT 1,
  `title` longtext NOT NULL,
  `description` varchar(255) NOT NULL,
  `content_title` longtext NOT NULL,
  `content` text NOT NULL,
  `additional_data` longtext DEFAULT NULL,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_user_access
CREATE TABLE IF NOT EXISTS `arpanet_documents_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `document_id` bigint(20) unsigned DEFAULT NULL,
  `identifier` varchar(64) DEFAULT NULL,
  `access` varchar(12) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_user_access_identifier` (`identifier`),
  KEY `idx_arpanet_documents_user_access_document_id` (`document_id`),
  KEY `idx_arpanet_documents_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_arpanet_documents_user_access` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`),
  CONSTRAINT `fk_arpanet_documents_users` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_activity
CREATE TABLE IF NOT EXISTS `arpanet_user_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `target_user_id` int(11) NOT NULL,
  `cause_user_id` int(11) NOT NULL,
  `type` longtext DEFAULT NULL,
  `key` varchar(64) DEFAULT NULL,
  `old_value` varchar(256) DEFAULT NULL,
  `new_value` varchar(256) DEFAULT NULL,
  `reason` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_user_activity_cause_user_id` (`cause_user_id`),
  KEY `idx_arpanet_user_activity_target_user_id` (`target_user_id`),
  CONSTRAINT `fk_users_cause_activity` FOREIGN KEY (`cause_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_users_target_activity` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_locations
CREATE TABLE IF NOT EXISTS `arpanet_user_locations` (
  `user_id` int(11) NOT NULL,
  `job` varchar(20) DEFAULT NULL,
  `x` float DEFAULT NULL,
  `y` float DEFAULT NULL,
  `hidden` tinyint(1) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  KEY `idx_arpanet_user_locations_job` (`job`),
  CONSTRAINT `fk_arpanet_user_locations_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_props
CREATE TABLE IF NOT EXISTS `arpanet_user_props` (
  `user_id` int(11) NOT NULL,
  `wanted` tinyint(1) NOT NULL DEFAULT 0,
  UNIQUE KEY `arpanet_user_props_UN` (`user_id`),
  KEY `idx_arpanet_user_props_wanted` (`wanted`),
  CONSTRAINT `arpanet_user_props_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- # Table: job_grades - Should already exist
-- CREATE TABLE IF NOT EXISTS `job_grades` (
--   `job_name` varchar(50) NOT NULL,
--   `grade` int(11) NOT NULL,
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) NOT NULL,
--   `salary` int(11) NOT NULL,
--   `skin_male` longtext NOT NULL,
--   `skin_female` longtext NOT NULL,
--   PRIMARY KEY (`job_name`,`grade`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- # Table: jobs - Should already exist
-- CREATE TABLE IF NOT EXISTS `jobs` (
--   `name` varchar(50) NOT NULL,
--   `label` varchar(50) DEFAULT NULL,
--   PRIMARY KEY (`name`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- # Table: owned_vehicles -- Should already exist
-- CREATE TABLE IF NOT EXISTS `owned_vehicles` (
--   `owner` varchar(64) DEFAULT NULL,
--   `plate` varchar(12) NOT NULL,
--   `model` varchar(60) NOT NULL,
--   `vehicle` longtext DEFAULT NULL,
--   `type` varchar(20) NOT NULL,
--   `stored` tinyint(1) NOT NULL DEFAULT 0,
--   `carseller` int(11) DEFAULT 0,
--   `owners` longtext DEFAULT NULL,
--   `trunk` longtext DEFAULT NULL,
--   PRIMARY KEY (`plate`),
--   UNIQUE KEY `IDX_OWNED_VEHICLES_OWNERPLATE` (`owner`,`plate`) USING BTREE,
--   KEY `IDX_OWNED_VEHICLES_OWNER` (`owner`),
--   KEY `IDX_OWNED_VEHICLES_OWNERTYPE` (`owner`,`type`),
--   KEY `IDX_OWNED_VEHICLES_OWNERRMODELTYPE` (`owner`,`model`,`type`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- # Table: user_licenses - Should already exist
-- CREATE TABLE IF NOT EXISTS `user_licenses` (
--   `type` varchar(60) NOT NULL,
--   `owner` varchar(64) NOT NULL,
--   PRIMARY KEY (`type`,`owner`),
--   KEY `idx_user_licenses_owner` (`owner`) USING BTREE
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- # Table: users - Should already exist
COMMIT;
