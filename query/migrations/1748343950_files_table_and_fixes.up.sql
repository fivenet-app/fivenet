BEGIN;

-- Table: fivenet_files
CREATE TABLE IF NOT EXISTS `fivenet_files` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `file_path` varchar(512) NOT NULL,
    `byte_size` bigint unsigned NOT NULL,
    `content_type` varchar(255) NOT NULL,
    `meta` json DEFAULT NULL,
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_path` (`file_path`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB;

-- Table: fivenet_documents_files
CREATE TABLE IF NOT EXISTS `fivenet_documents_files` (
    `document_id` bigint unsigned NOT NULL,
    `file_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`document_id`, `file_id`),
    KEY `idx_file_id` (`file_id`),
    CONSTRAINT `fk_fivenet_documents_files_document_id` FOREIGN KEY (`document_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_fivenet_documents_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT
);

-- Table: fivenet_mailer_messages_files
CREATE TABLE IF NOT EXISTS `fivenet_mailer_messages_files` (
    `message_id` bigint unsigned NOT NULL,
    `file_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`message_id`, `file_id`),
    KEY `idx_file_id` (`file_id`),
    CONSTRAINT `fk_fivenet_mailer_messages_files_message_id` FOREIGN KEY (`message_id`) REFERENCES `fivenet_mailer_messages` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_fivenet_mailer_messages_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT
);

-- Table: fivenet_qualifications_files
CREATE TABLE IF NOT EXISTS `fivenet_qualifications_files` (
    `qualification_id` bigint unsigned NOT NULL,
    `file_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`qualification_id`, `file_id`),
    KEY `idx_file_id` (`file_id`),
    CONSTRAINT `fk_fivenet_qualifications_files_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_fivenet_qualifications_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT
);

-- Table: fivenet_wiki_pages_files
CREATE TABLE IF NOT EXISTS `fivenet_wiki_pages_files` (
    `page_id` bigint unsigned NOT NULL,
    `file_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`page_id`, `file_id`),
    KEY `idx_file_id` (`file_id`),
    CONSTRAINT `fk_fivenet_wiki_pages_files_page_id` FOREIGN KEY (`page_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_fivenet_wiki_pages_files_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT
);

-- Table: `fivenet_*_access` - Drop any non-existing user_id records from the access tables
DELETE FROM `fivenet_calendar_access` WHERE `user_id` NOT IN (SELECT `id` FROM `{{.UsersTableName}}`);
DELETE FROM `fivenet_documents_access` WHERE `user_id` NOT IN (SELECT `id` FROM `{{.UsersTableName}}`);
DELETE FROM `fivenet_internet_domains_access` WHERE `user_id` NOT IN (SELECT `id` FROM `{{.UsersTableName}}`);
DELETE FROM `fivenet_mailer_emails_access` WHERE `user_id` NOT IN (SELECT `id` FROM `{{.UsersTableName}}`);
DELETE FROM `fivenet_wiki_pages_access` WHERE `user_id` NOT IN (SELECT `id` FROM `{{.UsersTableName}}`);

-- Table: `fivenet_calendar_access` - Fix unique indexes not working with NULL values
ALTER TABLE `fivenet_calendar_access` DROP INDEX `idx_fivenet_calendar_access_unique_access`;
ALTER TABLE `fivenet_calendar_access` DROP INDEX `fk_fivenet_calendar_access_job_grade`;
ALTER TABLE `fivenet_calendar_access` DROP FOREIGN KEY `fk_fivenet_calendar_access_target_id`;
ALTER TABLE `fivenet_calendar_access` DROP INDEX `idx_fivenet_calendar_access_unique`;

ALTER TABLE `fivenet_calendar_access` ADD INDEX `idx_job_minimum_grade` (`job`, `minimum_grade`);
ALTER TABLE `fivenet_calendar_access` ADD UNIQUE KEY `idx_user_id_access_unique` (`target_id`,`user_id`);
ALTER TABLE `fivenet_calendar_access` ADD UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`,`job`,`minimum_grade`);
ALTER TABLE `fivenet_calendar_access` ADD CONSTRAINT `fk_fivenet_calendar_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_access` ADD CONSTRAINT `fk_fivenet_calendar_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_documents_access` - Fix unique indexes not working with NULL values
ALTER TABLE `fivenet_documents_access` DROP INDEX `idx_fivenet_documents_access_unique_access`;
ALTER TABLE `fivenet_documents_access` DROP INDEX `fk_fivenet_documents_access_job_grade`;
ALTER TABLE `fivenet_documents_access` DROP FOREIGN KEY `fk_fivenet_documents_access_target_id`;
ALTER TABLE `fivenet_documents_access` DROP INDEX `idx_fivenet_documents_access_unique`;

ALTER TABLE `fivenet_documents_access` ADD INDEX `idx_job_minimum_grade` (`job`, `minimum_grade`);
ALTER TABLE `fivenet_documents_access` ADD UNIQUE KEY `idx_user_id_access_unique` (`target_id`,`user_id`);
ALTER TABLE `fivenet_documents_access` ADD UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`,`job`,`minimum_grade`);
ALTER TABLE `fivenet_documents_access` ADD CONSTRAINT `fk_fivenet_documents_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_access` ADD CONSTRAINT `fk_fivenet_documents_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_internet_domains_access` - Fix unique indexes not working with NULL values
ALTER TABLE `fivenet_internet_domains_access` DROP INDEX `idx_fivenet_internet_domain_access_unique_access`;
ALTER TABLE `fivenet_internet_domains_access` DROP INDEX `fk_fivenet_internet_domains_access_job_grade`;
ALTER TABLE `fivenet_internet_domains_access` DROP FOREIGN KEY `fk_fivenet_internet_domains_access_target_id`;
ALTER TABLE `fivenet_internet_domains_access` DROP INDEX `idx_fivenet_internet_domains_access_unique`;

ALTER TABLE `fivenet_internet_domains_access` ADD INDEX `idx_job_minimum_grade` (`job`, `minimum_grade`);
ALTER TABLE `fivenet_internet_domains_access` ADD UNIQUE KEY `idx_user_id_access_unique` (`target_id`,`user_id`);
ALTER TABLE `fivenet_internet_domains_access` ADD UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`,`job`,`minimum_grade`);
ALTER TABLE `fivenet_internet_domains_access` ADD CONSTRAINT `fk_fivenet_internet_domains_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_internet_domains_access` ADD CONSTRAINT `fk_fivenet_internet_domains_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_internet_domains` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_mailer_emails_access` - Fix unique indexes not working with NULL values
ALTER TABLE `fivenet_mailer_emails_access` DROP INDEX `idx_fivenet_mailer_emails_access_unique_access`;
ALTER TABLE `fivenet_mailer_emails_access` DROP INDEX `fk_fivenet_mailer_emails_access_job_grade`;
ALTER TABLE `fivenet_mailer_emails_access` DROP FOREIGN KEY `fk_fivenet_mailer_emails_access_target_id`;
ALTER TABLE `fivenet_mailer_emails_access` DROP INDEX `idx_fivenet_mailer_emails_access_unique`;

ALTER TABLE `fivenet_mailer_emails_access` ADD INDEX `idx_job_minimum_grade` (`job`, `minimum_grade`);
ALTER TABLE `fivenet_mailer_emails_access` ADD UNIQUE KEY `idx_user_id_access_unique` (`target_id`,`user_id`);
ALTER TABLE `fivenet_mailer_emails_access` ADD UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`,`job`,`minimum_grade`);
ALTER TABLE `fivenet_mailer_emails_access` ADD CONSTRAINT `fk_fivenet_mailer_emails_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_emails_access` ADD CONSTRAINT `fk_fivenet_mailer_emails_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_wiki_pages_access` - Fix unique indexes not working with NULL values
ALTER TABLE `fivenet_wiki_pages_access` DROP INDEX `idx_fivenet_wiki_pages_access_unique_access`;
ALTER TABLE `fivenet_wiki_pages_access` DROP INDEX `fk_fivenet_wiki_pages_access_job_grade`;
ALTER TABLE `fivenet_wiki_pages_access` DROP FOREIGN KEY `fk_fivenet_wiki_pages_access_target_id`;
ALTER TABLE `fivenet_wiki_pages_access` DROP INDEX `idx_fivenet_wiki_pages_access_unique`;

ALTER TABLE `fivenet_wiki_pages_access` ADD INDEX `idx_job_minimum_grade` (`job`, `minimum_grade`);
ALTER TABLE `fivenet_wiki_pages_access` ADD UNIQUE KEY `idx_user_id_access_unique` (`target_id`,`user_id`);
ALTER TABLE `fivenet_wiki_pages_access` ADD UNIQUE KEY `idx_job_minimum_grade_access_unique` (`target_id`,`job`,`minimum_grade`);
ALTER TABLE `fivenet_wiki_pages_access` ADD CONSTRAINT `fk_fivenet_wiki_pages_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_pages_access` ADD CONSTRAINT `fk_fivenet_wiki_pages_access_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_wiki_pages` - Add draft field
ALTER TABLE `fivenet_wiki_pages` ADD COLUMN `draft` tinyint(1) DEFAULT '0';
ALTER TABLE `fivenet_wiki_pages` CHANGE `draft` `draft` tinyint(1) DEFAULT '0' AFTER `toc`;

ALTER TABLE `fivenet_wiki_pages` ADD INDEX `idx_draft` (`draft`);

-- Table `fivenet_job_props` - Add `file_id` column
ALTER TABLE `fivenet_job_props` ADD COLUMN `logo_file_id` bigint unsigned DEFAULT NULL AFTER `logo_url`;
ALTER TABLE `fivenet_job_props` ADD CONSTRAINT `fk_fivenet_job_props_logo_file_id` FOREIGN KEY (`logo_file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- Table: `fivenet_user_props` - Add file id columns for avatar and mug shot
ALTER TABLE `fivenet_user_props` ADD COLUMN `avatar_file_id` bigint unsigned DEFAULT NULL AFTER `avatar`;
ALTER TABLE `fivenet_user_props` ADD COLUMN `mugshot_file_id` bigint unsigned DEFAULT NULL AFTER `mug_shot`;
ALTER TABLE `fivenet_user_props` ADD CONSTRAINT `fk_fivenet_user_props_avatar_file_id` FOREIGN KEY (`avatar_file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_props` ADD CONSTRAINT `fk_fivenet_user_props_mugshot_file_id` FOREIGN KEY (`mugshot_file_id`) REFERENCES `fivenet_files` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

UPDATE `fivenet_rbac_job_attrs` SET `max_values` = REPLACE(`max_values`, 'MugShot', 'Mugshot') WHERE `max_values` LIKE '%MugShot%';
UPDATE `fivenet_rbac_roles_attrs` SET `value` = REPLACE(`value`, 'MugShot', 'Mugshot') WHERE `value` LIKE '%MugShot%';

UPDATE `fivenet_rbac_job_attrs` SET `max_values` = REPLACE(`max_values`, 'Attributes', 'Labels') WHERE `max_values` LIKE '%Attributes%';
UPDATE `fivenet_rbac_roles_attrs` SET `value` = REPLACE(`value`, 'Attributes', 'Labels') WHERE `value` LIKE '%Attributes%';

COMMIT;
