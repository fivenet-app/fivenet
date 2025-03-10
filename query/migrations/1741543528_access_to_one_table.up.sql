BEGIN;

-- Table: fivenet_calendar_access
CREATE TABLE IF NOT EXISTS `fivenet_calendar_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_calendar_access_unique` (`target_id`, `user_id`, `job`, `minimum_grade`),
  UNIQUE KEY `idx_fivenet_calendar_access_unique_access` (`target_id`, `user_id`, `job`, `minimum_grade`, `access`),
  INDEX `fk_fivenet_calendar_access_user_id` (`user_id`),
  INDEX `fk_fivenet_calendar_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_calendar_access_access` (`access`),
  CONSTRAINT `fk_fivenet_calendar_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_calendar_access` (`target_id`, `user_id`, `access`)
SELECT `calendar_id`, `user_id`, `access` FROM `fivenet_calendar_user_access`;

INSERT INTO `fivenet_calendar_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `calendar_id`, `job`, `minimum_grade`, `access` FROM `fivenet_calendar_job_access`;

-- Table: fivenet_centrum_units_access
CREATE TABLE IF NOT EXISTS `fivenet_centrum_units_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `qualification_id` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_units_access_unique` (`target_id`, `job`, `minimum_grade`, `qualification_id`),
  UNIQUE KEY `idx_fivenet_centrum_units_access_unique_access` (`target_id`, `job`, `minimum_grade`, `qualification_id`, `access`),
  INDEX `fk_fivenet_centrum_units_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_centrum_units_access_qualification_id` (`qualification_id`),
  INDEX `fk_fivenet_centrum_units_access_access` (`access`),
  CONSTRAINT `fk_fivenet_centrum_units_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_centrum_units_access` (`target_id`, `qualification_id`, `access`)
SELECT `unit_id`, `qualification_id`, `access` FROM `fivenet_centrum_units_qualifications_access`;

INSERT INTO `fivenet_centrum_units_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `unit_id`, `job`, `minimum_grade`, `access` FROM `fivenet_centrum_units_job_access`;

-- Table: fivenet_documents_access
CREATE TABLE IF NOT EXISTS `fivenet_documents_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_access_unique` (`target_id`, `user_id`, `job`, `minimum_grade`),
  UNIQUE KEY `idx_fivenet_documents_access_unique_access` (`target_id`, `user_id`, `job`, `minimum_grade`, `access`),
  INDEX `fk_fivenet_documents_access_user_id` (`user_id`),
  INDEX `fk_fivenet_documents_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_documents_access_access` (`access`),
  CONSTRAINT `fk_fivenet_documents_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_documents_access` (`target_id`, `user_id`, `access`)
SELECT `document_id`, `user_id`, `access` FROM `fivenet_documents_user_access`;

INSERT INTO `fivenet_documents_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `document_id`, `job`, `minimum_grade`, `access` FROM `fivenet_documents_job_access`;

-- Table: fivenet_documents_templates_access
CREATE TABLE IF NOT EXISTS `fivenet_documents_templates_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_templates_access_unique` (`target_id`, `job`, `minimum_grade`),
  UNIQUE KEY `idx_fivenet_documents_templates_access_unique_access` (`target_id`, `job`, `minimum_grade`, `access`),
  INDEX `fk_fivenet_documents_templates_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_documents_templates_access_access` (`access`),
  CONSTRAINT `fk_fivenet_documents_templates_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_templates` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_documents_templates_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `template_id`, `job`, `minimum_grade`, `access` FROM `fivenet_documents_templates_job_access`;

-- Table: fivenet_internet_domains_access
CREATE TABLE IF NOT EXISTS `fivenet_internet_domains_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_internet_domain_access_unique` (`target_id`, `user_id`, `job`, `minimum_grade`),
  UNIQUE KEY `idx_fivenet_internet_domain_access_unique_access` (`target_id`, `user_id`, `job`, `minimum_grade`, `access`),
  INDEX `fk_fivenet_internet_domains_access_user_id` (`user_id`),
  INDEX `fk_fivenet_internet_domains_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_internet_domains_access_access` (`access`),
  CONSTRAINT `fk_fivenet_internet_domains_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_internet_domains` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_internet_domains_access` (`target_id`, `user_id`, `access`)
SELECT `domain_id`, `user_id`, `access` FROM `fivenet_internet_domains_user_access`;

INSERT INTO `fivenet_internet_domains_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `domain_id`, `job`, `minimum_grade`, `access` FROM `fivenet_internet_domains_job_access`;

-- Table: fivenet_mailer_emails_access
CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `qualification_id` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_access_unique` (`target_id`, `user_id`, `job`, `minimum_grade`, `qualification_id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_access_unique_access` (`target_id`, `user_id`, `job`, `minimum_grade`, `qualification_id`, `access`),
  INDEX `fk_fivenet_mailer_emails_access_user_id` (`user_id`),
  INDEX `fk_fivenet_mailer_emails_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_mailer_emails_access_access` (`access`),
  CONSTRAINT `fk_fivenet_mailer_emails_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_mailer_emails_access` (`target_id`, `user_id`, `access`)
SELECT `email_id`, `user_id`, `access` FROM `fivenet_mailer_emails_user_access`;

INSERT INTO `fivenet_mailer_emails_access` (`target_id`, `qualification_id`, `access`)
SELECT `email_id`, `qualification_id`, `access` FROM `fivenet_mailer_emails_qualifications_access`;

INSERT INTO `fivenet_mailer_emails_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `email_id`, `job`, `minimum_grade`, `access` FROM `fivenet_mailer_emails_job_access`;

-- Table: fivenet_qualifications_access
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_qualifications_access_unique` (`target_id`, `job`, `minimum_grade`),
  UNIQUE KEY `idx_fivenet_qualifications_access_unique_access` (`target_id`, `job`, `minimum_grade`, `access`),
  INDEX `fk_fivenet_qualifications_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_qualifications_access_access` (`access`),
  CONSTRAINT `fk_fivenet_qualifications_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_qualifications_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `qualification_id`, `job`, `minimum_grade`, `access` FROM `fivenet_qualifications_job_access`;

-- Table: fivenet_wiki_pages_access
CREATE TABLE IF NOT EXISTS `fivenet_wiki_pages_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NULL,
  `job` varchar(40) NULL,
  `minimum_grade` int(11) NULL,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_wiki_pages_access_unique` (`target_id`, `user_id`, `job`, `minimum_grade`),
  UNIQUE KEY `idx_fivenet_wiki_pages_access_unique_access` (`target_id`, `user_id`, `job`, `minimum_grade`, `access`),
  INDEX `fk_fivenet_wiki_pages_access_user_id` (`user_id`),
  INDEX `fk_fivenet_wiki_pages_access_job_grade` (`job`, `minimum_grade`),
  INDEX `fk_fivenet_wiki_pages_access_access` (`access`),
  CONSTRAINT `fk_fivenet_wiki_pages_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;

INSERT INTO `fivenet_wiki_pages_access` (`target_id`, `user_id`, `access`)
SELECT `page_id`, `user_id`, `access` FROM `fivenet_wiki_page_user_access`;

INSERT INTO `fivenet_wiki_pages_access` (`target_id`, `job`, `minimum_grade`, `access`)
SELECT `page_id`, `job`, `minimum_grade`, `access` FROM `fivenet_wiki_page_job_access`;

DROP TABLE `fivenet_calendar_job_access`;
DROP TABLE `fivenet_calendar_user_access`;
DROP TABLE `fivenet_centrum_units_job_access`;
DROP TABLE `fivenet_centrum_units_qualifications_access`;
DROP TABLE `fivenet_documents_job_access`;
DROP TABLE `fivenet_documents_user_access`;
DROP TABLE `fivenet_documents_templates_job_access`;
DROP TABLE `fivenet_internet_domains_job_access`;
DROP TABLE `fivenet_internet_domains_user_access`;
DROP TABLE `fivenet_mailer_emails_job_access`;
DROP TABLE `fivenet_mailer_emails_qualifications_access`;
DROP TABLE `fivenet_mailer_emails_user_access`;
DROP TABLE `fivenet_qualifications_job_access`;
DROP TABLE `fivenet_wiki_page_job_access`;
DROP TABLE `fivenet_wiki_page_user_access`;

COMMIT;
