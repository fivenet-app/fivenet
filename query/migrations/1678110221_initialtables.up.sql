BEGIN;

-- Table: fivenet_accounts
CREATE TABLE IF NOT EXISTS `fivenet_accounts` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `enabled` tinyint(1) DEFAULT 0,
  `username` varchar(24) NULL,
  `password` varchar(60) NULL,
  `license` varchar(64) NOT NULL,
  `reg_token` char(6) DEFAULT NULL,
  `override_job` varchar(50) DEFAULT NULL,
  `override_job_grade` int(11) DEFAULT NULL,
  `superuser` tinyint(1) DEFAULT 0,
  `last_char` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_accounts_username` (`username`),
  UNIQUE KEY `idx_fivenet_accounts_license` (`license`),
  UNIQUE KEY `idx_fivenet_accounts_username_license` (`username`, `license`),
  UNIQUE KEY `idx_fivenet_accounts_reg_token` (`reg_token`)
) ENGINE=InnoDB;

-- Table: fivenet_documents_categories
CREATE TABLE IF NOT EXISTS `fivenet_documents_categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `description` longtext DEFAULT NULL,
  `job` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_documents_categories_job` (`job`)
) ENGINE=InnoDB;

-- Table: fivenet_documents_templates
CREATE TABLE IF NOT EXISTS `fivenet_documents_templates` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `weight` int(11) unsigned DEFAULT 0,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `title` longtext NOT NULL,
  `description` longtext NOT NULL,
  `content_title` longtext NOT NULL,
  `content` longtext NOT NULL,
  `state` varchar(24) NOT NULL,
  `access` longtext DEFAULT NULL,
  `schema` longtext DEFAULT NULL,
  `creator_job` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_documents_templates_deleted_at` (`deleted_at`),
  KEY `idx_fivenet_documents_templates_weight` (`weight`),
  KEY `idx_fivenet_documents_templates_category_id` (`category_id`),
  CONSTRAINT `fk_fivenet_documents_templates_categories` FOREIGN KEY (`category_id`) REFERENCES `fivenet_documents_categories` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB;

-- Table: fivenet_documents
CREATE TABLE IF NOT EXISTS `fivenet_documents` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `title` longtext NOT NULL,
  `summary` varchar(128) NOT NULL,
  `content_type` smallint(2) NOT NULL,
  `content` longtext NOT NULL,
  `data` longtext DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `creator_job` varchar(50) NOT NULL,
  `state` varchar(32) NOT NULL,
  `closed` tinyint(1) DEFAULT 0,
  `public` tinyint(1) NOT NULL DEFAULT 0,
  `template_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_documents_created_at` (`created_at`),
  KEY `idx_fivenet_documents_deleted_at` (`deleted_at`),
  KEY `idx_fivenet_documents_category_id` (`category_id`),
  KEY `idx_fivenet_documents_creator_id` (`creator_id`),
  FULLTEXT KEY `idx_fivenet_documents_title` (`title`),
  FULLTEXT KEY `idx_fivenet_documents_content` (`content`),
  CONSTRAINT `fk_fivenet_documents_categories` FOREIGN KEY (`category_id`) REFERENCES `fivenet_documents_categories` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_fivenet_documents_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB;

-- Table: fivenet_documents_comments
CREATE TABLE IF NOT EXISTS `fivenet_documents_comments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `document_id` bigint(20) unsigned NOT NULL,
  `comment` longtext,
  `creator_id` int(11) NOT NULL,
  `creator_job` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_documents_comments_document_id` (`document_id`),
  KEY `idx_fivenet_documents_comments_creator_id` (`creator_id`),
  CONSTRAINT `fk_fivenet_documents_comments_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_comments_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_documents_job_access
CREATE TABLE IF NOT EXISTS `fivenet_documents_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `document_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 1,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_job_access` (`document_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_documents_job_access_document_id` (`document_id`),
  CONSTRAINT `fk_fivenet_documents_job_access_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_documents_references
CREATE TABLE IF NOT EXISTS `fivenet_documents_references` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `source_document_id` bigint(20) unsigned NOT NULL,
  `reference` smallint(2) NOT NULL,
  `target_document_id` bigint(20) unsigned NOT NULL,
  `creator_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_references_unique` (`source_document_id`, `reference`, `target_document_id`, `creator_id`),
  KEY `idx_fivenet_documents_references_deleted_at` (`deleted_at`),
  KEY `idx_fivenet_documents_references_source_document_id` (`source_document_id`),
  KEY `idx_fivenet_documents_references_target_document_id` (`target_document_id`),
  KEY `idx_fivenet_documents_references_creator_id` (`creator_id`),
  CONSTRAINT `fk_fivenet_documents_references_source_document_id` FOREIGN KEY (`source_document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_references_target_document_id` FOREIGN KEY (`target_document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_references_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB;

-- Table: fivenet_documents_relations
CREATE TABLE IF NOT EXISTS `fivenet_documents_relations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `document_id` bigint(20) unsigned NOT NULL,
  `source_user_id` int(11) DEFAULT NULL,
  `relation` smallint(2) NOT NULL,
  `target_user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_relations_unique` (`document_id`, `relation`, `target_user_id`),
  KEY `idx_fivenet_documents_relations_deleted_at` (`deleted_at`),
  KEY `idx_fivenet_documents_relations_document_id` (`document_id`),
  KEY `idx_fivenet_documents_relations_source_user_id` (`source_user_id`),
  KEY `idx_fivenet_documents_relations_target_user_id` (`target_user_id`),
  CONSTRAINT `fk_fivenet_documents_relations_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_relations_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_fivenet_documents_relations_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_documents_user_access
CREATE TABLE IF NOT EXISTS `fivenet_documents_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `document_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_user_access` (`document_id`, `user_id`),
  KEY `idx_fivenet_documents_user_access_document_id` (`document_id`),
  KEY `idx_fivenet_documents_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_documents_user_access_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_documents_activity
CREATE TABLE IF NOT EXISTS `fivenet_documents_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `document_id` bigint(20) unsigned NOT NULL,
  `activity_type` smallint(2) NOT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `creator_job` varchar(20) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `data` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_documents_activity_document_id` (`document_id`),
  KEY `idx_fivenet_documents_activity_creator_id` (`creator_id`),
  KEY `idx_fivenet_documents_activity_activity_type` (`activity_type`),
  CONSTRAINT `fk_fivenet_documents_activity_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB;

-- Table: fivenet_documents_requests
CREATE TABLE IF NOT EXISTS `fivenet_documents_requests` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `document_id` bigint(20) unsigned NOT NULL,
  `request_type` smallint(2) NOT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `creator_job` varchar(20) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `data` longtext DEFAULT NULL,
  `accepted` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_documents_requests_document_id` (`document_id`),
  KEY `idx_fivenet_documents_requests_creator_id` (`creator_id`),
  KEY `idx_fivenet_documents_requests_request_type` (`request_type`),
  KEY `idx_fivenet_documents_requests_accepted` (`accepted`),
  KEY `idx_fivenet_documents_requests_unique` (`document_id`, `request_type`),
  CONSTRAINT `fk_fivenet_documents_requests_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_requests_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB;

-- Table: fivenet_user_activity
CREATE TABLE IF NOT EXISTS `fivenet_user_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `source_user_id` int(11) DEFAULT NULL,
  `target_user_id` int(11) NOT NULL,
  `type` smallint(2) NOT NULL,
  `key` varchar(64) NOT NULL,
  `old_value` varchar(255) DEFAULT NULL,
  `new_value` varchar(255) DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_user_activity_source_user_id` (`source_user_id`),
  KEY `idx_fivenet_user_activity_target_user_id` (`target_user_id`),
  CONSTRAINT `fk_fivenet_user_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_fivenet_user_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_user_locations
CREATE TABLE IF NOT EXISTS `fivenet_user_locations` (
  `identifier` varchar(64) NOT NULL,
  `job` varchar(20) NOT NULL,
  `x` decimal(24,14) DEFAULT NULL,
  `y` decimal(24,14) DEFAULT NULL,
  `hidden` tinyint(1) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`identifier`),
  KEY `idx_fivenet_user_locations_job` (`job`),
  CONSTRAINT `fk_fivenet_user_locations_identifier` FOREIGN KEY (`identifier`) REFERENCES `users` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_user_props
CREATE TABLE IF NOT EXISTS `fivenet_user_props` (
  `user_id` int(11) NOT NULL,
  `wanted` tinyint(1) DEFAULT 0,
  `job` varchar(20) DEFAULT NULL,
  `job_grade` int(11) DEFAULT NULL,
  `traffic_infraction_points` mediumint(8) unsigned DEFAULT 0,
  `open_fines` bigint(20) DEFAULT 0,
  `blood_type` varchar(3) DEFAULT NULL,
  `avatar` varchar(128) DEFAULT NULL,
  `mug_shot` varchar(128) DEFAULT NULL,
  `attributes` text DEFAULT NULL,
  UNIQUE KEY `idx_fivenet_user_props_unique` (`user_id`),
  KEY `idx_fivenet_user_props_wanted` (`wanted`),
  KEY `idx_fivenet_user_props_avatar` (`avatar`),
  KEY `idx_fivenet_user_props_mug_shot` (`mug_shot`),
  CONSTRAINT `fk_fivenet_user_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
