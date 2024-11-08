BEGIN;

-- Table: fivenet_wiki_pages
CREATE TABLE IF NOT EXISTS `fivenet_wiki_pages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `job` varchar(50) NOT NULL,
  `parent_id` bigint(20) unsigned DEFAULT NULL,
  `content_type` smallint(2) NOT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `toc` tinyint(1) DEFAULT 1,
  `public` tinyint(1) NOT NULL DEFAULT 0,
  `slug` varchar(100) NOT NULL,
  `title` longtext NOT NULL,
  `description` varchar(128) NOT NULL,
  `content` longtext NOT NULL,
  `data` longtext DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_wiki_pages_id_job` (`id`, `job`),
  KEY `idx_fivenet_wiki_pages_parent_id` (`parent_id`),
  FULLTEXT KEY `idx_fivenet_wiki_pages_title` (`title`),
  FULLTEXT KEY `idx_fivenet_wiki_pages_content` (`content`),
  KEY `idx_fivenet_wiki_pages_creator_id` (`creator_id`),
  CONSTRAINT `fk_fivenet_wiki_pages_parent_id` FOREIGN KEY (`parent_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_wiki_pages_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_wiki_page_user_access
CREATE TABLE IF NOT EXISTS `fivenet_wiki_page_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `page_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 1,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_wiki_page_job_access` (`page_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_wiki_page_job_access_page_id` (`page_id`),
  CONSTRAINT `fk_fivenet_wiki_page_job_access_page_id` FOREIGN KEY (`page_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_wiki_page_user_access
CREATE TABLE IF NOT EXISTS `fivenet_wiki_page_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `page_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_wiki_page_user_access` (`page_id`, `user_id`),
  KEY `idx_fivenet_wiki_page_user_access_page_id` (`page_id`),
  KEY `idx_fivenet_wiki_page_user_access_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_wiki_page_user_access_page_id` FOREIGN KEY (`page_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_wiki_page_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_wiki_page_activity
CREATE TABLE IF NOT EXISTS `fivenet_wiki_page_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `page_id` bigint(20) unsigned NOT NULL,
  `activity_type` smallint(2) NOT NULL,
  `creator_id` int(11) DEFAULT NULL,
  `creator_job` varchar(20) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `data` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_wiki_page_activity_page_id` (`page_id`),
  KEY `idx_fivenet_wiki_page_activity_creator_id` (`creator_id`),
  KEY `idx_fivenet_wiki_page_activity_activity_type` (`activity_type`),
  CONSTRAINT `fk_fivenet_wiki_page_activity_page_id` FOREIGN KEY (`page_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_wiki_page_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB;

COMMIT;
