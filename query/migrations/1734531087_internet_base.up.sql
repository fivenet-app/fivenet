-- Table: fivenet_internet_tlds
CREATE TABLE IF NOT EXISTS `fivenet_internet_tlds` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `internal` tinyint(1) DEFAULT 0 NULL,
  `creator_id`int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_internet_tlds_name` (`name`),
  CONSTRAINT `fk_fivenet_internet_tlds_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_internet_domains
CREATE TABLE IF NOT EXISTS `fivenet_internet_domains` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `tld_id` bigint(20) unsigned NOT NULL,
  `name` varchar(64) NOT NULL,
  `active` tinyint(1) DEFAULT 0 NULL,
  `transfer_code` varchar(10) DEFAULT NULL,
  `approver_job` varchar(40) DEFAULT NULL,
  `approver_id` int DEFAULT NULL,
  `creator_job` varchar(40) DEFAULT NULL,
  `creator_id`int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_internet_domains_name` (`name`),
  CONSTRAINT `fk_fivenet_internet_domains_tld_id` FOREIGN KEY (`tld_id`) REFERENCES `fivenet_internet_tlds` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_internet_domains_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_internet_domains_job_access
CREATE TABLE IF NOT EXISTS `fivenet_internet_domains_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `domain_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_internet_domains_job_access` (`domain_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_internet_domains_job_access_domain_id` (`domain_id`),
  CONSTRAINT `fk_fivenet_internet_domains_job_access_domain_id` FOREIGN KEY (`domain_id`) REFERENCES `fivenet_internet_domains` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_internet_domains_user_access
CREATE TABLE IF NOT EXISTS `fivenet_internet_domains_user_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `domain_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_internet_domains_user_access` (`domain_id`, `user_id`),
  KEY `idx_fivenet_internet_domains_user_access_domain_id` (`domain_id`),
  CONSTRAINT `fk_fivenet_internet_domains_user_access_domain_id` FOREIGN KEY (`domain_id`) REFERENCES `fivenet_internet_domains` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_internet_domains_user_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_internet_pages
CREATE TABLE IF NOT EXISTS `fivenet_internet_pages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `domain_id` bigint(20) unsigned NOT NULL,
  `path` varchar(128) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` varchar(512) NOT NULL,
  `data` text DEFAULT NULL,
  `creator_job` varchar(40) DEFAULT NULL,
  `creator_id`int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_internet_pages_domain_id_path` (`domain_id`, `path`),
  FULLTEXT KEY `idx_fivenet_internet_pages_title` (`title`),
  FULLTEXT KEY `idx_fivenet_internet_pages_description` (`description`),
  CONSTRAINT `fk_fivenet_internet_pages_domain_id` FOREIGN KEY (`domain_id`) REFERENCES `fivenet_internet_domains` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_internet_pages_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_internet_ads
CREATE TABLE IF NOT EXISTS `fivenet_internet_ads` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `disabled` tinyint(1) DEFAULT 1,
  `ad_type` smallint(2) NOT NULL,
  `starts_at` datetime(3) DEFAULT NULL,
  `ends_at` datetime(3) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `description` varchar(1024) NOT NULL,
  `image` varchar(128) DEFAULT NULL,
  `approver_job` varchar(40) DEFAULT NULL,
  `approver_id`int(11) DEFAULT NULL,
  `creator_job` varchar(40) DEFAULT NULL,
  `creator_id`int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_internet_ads_starts_at_ends_at` (`starts_at`, `ends_at`),
  CONSTRAINT `fk_fivenet_internet_ads_approver_id` FOREIGN KEY (`approver_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_internet_ads_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

BEGIN;

INSERT INTO `fivenet_internet_tlds` (`name`, `internal`) VALUES('ls', 0);
INSERT INTO `fivenet_internet_tlds` (`name`, `internal`) VALUES('search', 1);

INSERT INTO `fivenet_internet_domains` (`tld_id`, `name`, `active`) VALUES((SELECT `id` FROM `fivenet_internet_tlds` WHERE `name` = 'ls' LIMIT 1), 'nic', 1);
INSERT INTO `fivenet_internet_domains` (`tld_id`, `name`, `active`) VALUES((SELECT `id` FROM `fivenet_internet_tlds` WHERE `name` = 'ls' LIMIT 1), 'ads', 1);
INSERT INTO `fivenet_internet_domains` (`tld_id`, `name`, `active`) VALUES((SELECT `id` FROM `fivenet_internet_tlds` WHERE `name` = 'search' LIMIT 1), 'internet', 1);
INSERT INTO `fivenet_internet_domains` (`tld_id`, `name`, `active`) VALUES((SELECT `id` FROM `fivenet_internet_tlds` WHERE `name` = 'ls' LIMIT 1), 'example', 1);

-- Table: fivenet_internet_pages - Insert `example.ls` page entry
INSERT INTO fivenet_internet_pages (id, domain_id, `path`, title, description, `data`, creator_job, creator_id) VALUES(1, 4, '/', 'Example Domain', 'Example Domain for demonstration purposes.', '{"layoutType":"PAGE_LAYOUT_TYPE_LANDING_PAGE","node":{"type":"NODE_TYPE_ELEMENT","tag":"ULandingSection","attrs":{},"content":[{"type":"NODE_TYPE_ELEMENT","tag":"ULandingCard","attrs":{"title":"Example Domain","description":"Example Domain for demonstration purposes."},"content":[],"slots":[]}],"slots":[]}}', NULL, NULL);

COMMIT;
