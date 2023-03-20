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
  UNIQUE KEY `idx_arpanet_accounts_username` (`username`),
  UNIQUE KEY `idx_arpanet_accounts_license` (`license`),
  UNIQUE KEY `idx_arpanet_accounts_username_license` (`username`, `license`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_categories
CREATE TABLE IF NOT EXISTS `arpanet_documents_categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `description` longtext DEFAULT NULL,
  `job` varchar(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_categories_job` (`job`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_templates
CREATE TABLE IF NOT EXISTS `arpanet_documents_templates` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE current_timestamp(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `job` varchar(20) NOT NULL,
  `job_grade` int(11) NOT NULL DEFAULT 0,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `title` longtext NOT NULL,
  `description` longtext NOT NULL,
  `content_title` longtext NOT NULL,
  `content` text NOT NULL,
  `additional_data` longtext DEFAULT NULL,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_templates_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_templates_category_id` (`category_id`),
  KEY `idx_arpanet_documents_templates_job_job_grade` (`job`, `job_grade`),
  CONSTRAINT `fk_arpanet_documents_templates_categories` FOREIGN KEY (`category_id`) REFERENCES `arpanet_documents_categories` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_arpanet_documents_templates_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents
CREATE TABLE IF NOT EXISTS `arpanet_documents` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `title` longtext NOT NULL,
  `content_type` smallint(2) NOT NULL,
  `content` longtext NOT NULL,
  `data` longtext DEFAULT NULL,
  `creator_id` int(11) NOT NULL,
  `creator_job` varchar(50) NOT NULL,
  `state` varchar(24) NOT NULL,
  `closed` tinyint(1) DEFAULT 0,
  `public` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_category_id` (`category_id`),
  KEY `idx_arpanet_documents_creator_id` (`creator_id`),
  FULLTEXT KEY `idx_arpanet_documents_title` (`title`),
  FULLTEXT KEY `idx_arpanet_documents_content` (`content`),
  CONSTRAINT `fk_arpanet_documents_categories` FOREIGN KEY (`category_id`) REFERENCES `arpanet_documents_categories` (`id`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fk_arpanet_documents_categories_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table:arpanet_documents_comments
CREATE TABLE IF NOT EXISTS `arpanet_documents_comments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `document_id` bigint(20) unsigned NOT NULL,
  `comment` longtext,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_documents_comments_document_id` (`document_id`),
  KEY `idx_arpanet_documents_comments_creator_id` (`creator_id`),
  CONSTRAINT `fk_arpanet_documents_comments_document_id` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_documents_comments_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_job_access
CREATE TABLE IF NOT EXISTS `arpanet_documents_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `document_id` bigint(20) unsigned NOT NULL,
  `job` varchar(20) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_arpanet_documents_job_access` (`document_id`, `job`),
  KEY `idx_arpanet_documents_job_access_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_job_access_document_id` (`document_id`),
  CONSTRAINT `fk_arpanet_documents_job_access_document_id` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_documents_job_access_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_references
CREATE TABLE IF NOT EXISTS `arpanet_documents_references` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `source_document_id` bigint(20) unsigned NOT NULL,
  `reference` smallint(2) NOT NULL,
  `target_document_id` bigint(20) unsigned NOT NULL,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_arpanet_documents_references_unique` (`source_document_id`, `reference`, `target_document_id`, `creator_id`),
  KEY `idx_arpanet_documents_references_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_references_source_document_id` (`source_document_id`),
  KEY `idx_arpanet_documents_references_target_document_id` (`target_document_id`),
  KEY `idx_arpanet_documents_references_creator_id` (`creator_id`),
  CONSTRAINT `fk_arpanet_documents_references_source_document_id` FOREIGN KEY (`source_document_id`) REFERENCES `arpanet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_documents_references_target_document_id` FOREIGN KEY (`target_document_id`) REFERENCES `arpanet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_documents_references_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_relations
CREATE TABLE IF NOT EXISTS `arpanet_documents_relations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `document_id` bigint(20) unsigned NOT NULL,
  `source_user_id` int(11) NOT NULL,
  `relation` smallint(2) NOT NULL,
  `target_user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_arpanet_documents_relations_unique` (`source_user_id`, `relation`, `target_user_id`),
  KEY `idx_arpanet_documents_relations_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_relations_document_id` (`document_id`),
  KEY `idx_arpanet_documents_relations_source_user_id` (`source_user_id`),
  KEY `idx_arpanet_documents_relations_target_user_id` (`target_user_id`),
  CONSTRAINT `fk_arpanet_documents_relations_document_id` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_documents_relations_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_arpanet_documents_relations_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_documents_user_access
CREATE TABLE IF NOT EXISTS `arpanet_documents_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `document_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  `creator_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_arpanet_documents_user_access` (`document_id`, `user_id`),
  KEY `idx_arpanet_documents_users_access_deleted_at` (`deleted_at`),
  KEY `idx_arpanet_documents_user_access_document_id` (`document_id`),
  KEY `idx_arpanet_documents_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_arpanet_documents_user_access_document_id` FOREIGN KEY (`document_id`) REFERENCES `arpanet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_documents_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_arpanet_documents_user_access_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_activity
CREATE TABLE IF NOT EXISTS `arpanet_user_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `source_user_id` int(11) NOT NULL,
  `target_user_id` int(11) NOT NULL,
  `type` smallint(2) NOT NULL,
  `key` varchar(64) NOT NULL,
  `old_value` varchar(255) DEFAULT NULL,
  `new_value` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_arpanet_user_activity_source_user_id` (`source_user_id`),
  KEY `idx_arpanet_user_activity_target_user_id` (`target_user_id`),
  CONSTRAINT `fk_arpanet_user_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_arpanet_user_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_locations
CREATE TABLE IF NOT EXISTS `arpanet_user_locations` (
  `user_id` int(11) NOT NULL,
  `job` varchar(20) DEFAULT NULL,
  `x` decimal(24,14) DEFAULT NULL,
  `y` decimal(24,14) DEFAULT NULL,
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
  UNIQUE KEY `idx_arpanet_user_props_unique` (`user_id`),
  KEY `idx_arpanet_user_props_wanted` (`wanted`),
  CONSTRAINT `fk_arpanet_user_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
