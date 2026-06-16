BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_acl_subjects` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `subject_type` smallint(2) unsigned NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_acl_subjects_subject_type` (`subject_type`),
  CONSTRAINT `chk_fivenet_acl_subjects_subject_type`
    CHECK (`subject_type` IN (1, 2, 3))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_acl_subject_users` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_fivenet_acl_subject_users_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_acl_subject_users_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_acl_subject_users_user_id`
    FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_acl_subject_qualifications` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_fivenet_acl_subject_qualifications_qualification_id` (`qualification_id`),
  CONSTRAINT `fk_fivenet_acl_subject_qualifications_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_acl_subject_qualifications_qualification_id`
    FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_acl_subject_job_grade_scopes` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_fivenet_acl_subject_job_grade_scope` (`job`, `minimum_grade`),
  KEY `idx_fivenet_acl_subject_job_grade_match` (`job`, `minimum_grade`, `subject_id`),
  CONSTRAINT `fk_fivenet_acl_subject_job_grade_scopes_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_documents_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_documents_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_documents_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_documents_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_calendar_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_calendar_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_calendar_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_calendar_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_calendar_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_calendar` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_calendar_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_calendar_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_wiki_pages_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_wiki_pages_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_wiki_pages_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_wiki_pages_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_wiki_pages_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_wiki_pages` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_wiki_pages_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_wiki_pages_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_templates_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_templates_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_documents_templates_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_documents_templates_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_documents_templates_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_templates` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_templates_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_documents_templates_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_stamps_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_stamps_access_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_documents_stamps_access_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_documents_stamps_access_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_documents_stamps_access_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_stamps` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_stamps_access_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_documents_stamps_access_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_user_labels_job_job_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_user_labels_job_job_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_user_labels_job_job_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_user_labels_job_job_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_user_labels_job_job_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_user_labels_job` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_user_labels_job_job_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_user_labels_job_job_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_qualifications_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_qualifications_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_qualifications_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_qualifications_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_qualifications_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_qualifications` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_qualifications_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_mailer_emails_access_v2_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_mailer_emails_access_v2_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_mailer_emails_access_v2_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_mailer_emails_access_v2_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_mailer_emails` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_emails_access_v2_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_mailer_emails_access_v2_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_centrum_units_access_v2` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `target_id` bigint(20) unsigned NOT NULL,
  `subject_id` bigint(20) unsigned NOT NULL,
  `access` smallint(2) NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_centrum_units_access_unique` (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_centrum_units_access_target_access` (`target_id`, `access`, `subject_id`, `effect`),
  KEY `idx_fivenet_centrum_units_access_subject_access` (`subject_id`, `access`, `target_id`, `effect`),
  CONSTRAINT `fk_fivenet_centrum_units_access_subject_target_id`
    FOREIGN KEY (`target_id`) REFERENCES `fivenet_centrum_units` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_centrum_units_access_subject_subject_id`
    FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_centrum_units_access_effect`
    CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

-- Backfill ACL subjects from the legacy documents access table.
CREATE TEMPORARY TABLE `tmp_fivenet_acl_subject_users_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_subject_users_user_id` (`user_id`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_subject_users_backfill` (`subject_id`, `user_id`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_user_access`.`user_id`
FROM (
  SELECT DISTINCT `da`.`user_id`
  FROM `fivenet_documents_access` `da`
  LEFT JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `da`.`user_id`
  WHERE `da`.`user_id` IS NOT NULL
    AND `su`.`subject_id` IS NULL
  ORDER BY `da`.`user_id`
) `legacy_user_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 1 FROM `tmp_fivenet_acl_subject_users_backfill`;

INSERT INTO `fivenet_acl_subject_users` (`subject_id`, `user_id`)
SELECT `subject_id`, `user_id` FROM `tmp_fivenet_acl_subject_users_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_subject_users_backfill`;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `da`.`job`, `da`.`minimum_grade`
  FROM `fivenet_documents_access` `da`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `da`.`job` AND `sj`.`minimum_grade` = `da`.`minimum_grade`
  WHERE `da`.`job` IS NOT NULL
    AND `da`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `da`.`job`, `da`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_subject_job_grades_backfill`;

-- Legacy BLOCKED access (1) represented a subject-level block. In v2 this is
-- stored as deny rows for each concrete document access level.
INSERT IGNORE INTO `fivenet_documents_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `da`.`target_id`, `su`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_documents_access` `da`
INNER JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `da`.`user_id`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5 UNION ALL
  SELECT 6
) `denied_levels`
WHERE `da`.`user_id` IS NOT NULL
  AND `da`.`access` = 1;

INSERT IGNORE INTO `fivenet_documents_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `da`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_documents_access` `da`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `da`.`job` AND `sj`.`minimum_grade` = `da`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5 UNION ALL
  SELECT 6
) `denied_levels`
WHERE `da`.`job` IS NOT NULL
  AND `da`.`minimum_grade` IS NOT NULL
  AND `da`.`access` = 1;

INSERT IGNORE INTO `fivenet_documents_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `da`.`target_id`, `su`.`subject_id`, `da`.`access`, 1
FROM `fivenet_documents_access` `da`
INNER JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `da`.`user_id`
WHERE `da`.`user_id` IS NOT NULL
  AND `da`.`access` > 1;

INSERT IGNORE INTO `fivenet_documents_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `da`.`target_id`, `sj`.`subject_id`, `da`.`access`, 1
FROM `fivenet_documents_access` `da`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `da`.`job` AND `sj`.`minimum_grade` = `da`.`minimum_grade`
WHERE `da`.`job` IS NOT NULL
  AND `da`.`minimum_grade` IS NOT NULL
  AND `da`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_calendar_subject_users_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_calendar_subject_users_user_id` (`user_id`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_calendar_subject_users_backfill` (`subject_id`, `user_id`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_user_access`.`user_id`
FROM (
  SELECT DISTINCT `ca`.`user_id`
  FROM `fivenet_calendar_access` `ca`
  LEFT JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `ca`.`user_id`
  WHERE `ca`.`user_id` IS NOT NULL
    AND `su`.`subject_id` IS NULL
  ORDER BY `ca`.`user_id`
) `legacy_user_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 1 FROM `tmp_fivenet_acl_calendar_subject_users_backfill`;

INSERT INTO `fivenet_acl_subject_users` (`subject_id`, `user_id`)
SELECT `subject_id`, `user_id` FROM `tmp_fivenet_acl_calendar_subject_users_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_calendar_subject_users_backfill`;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_calendar_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_calendar_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_calendar_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `ca`.`job`, `ca`.`minimum_grade`
  FROM `fivenet_calendar_access` `ca`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `ca`.`job` AND `sj`.`minimum_grade` = `ca`.`minimum_grade`
  WHERE `ca`.`job` IS NOT NULL
    AND `ca`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `ca`.`job`, `ca`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_calendar_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_calendar_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_calendar_subject_job_grades_backfill`;

INSERT IGNORE INTO `fivenet_calendar_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ca`.`target_id`, `su`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_calendar_access` `ca`
INNER JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `ca`.`user_id`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5
) `denied_levels`
WHERE `ca`.`user_id` IS NOT NULL
  AND `ca`.`access` = 1;

INSERT IGNORE INTO `fivenet_calendar_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ca`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_calendar_access` `ca`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ca`.`job` AND `sj`.`minimum_grade` = `ca`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5
) `denied_levels`
WHERE `ca`.`job` IS NOT NULL
  AND `ca`.`minimum_grade` IS NOT NULL
  AND `ca`.`access` = 1;

INSERT IGNORE INTO `fivenet_calendar_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ca`.`target_id`, `su`.`subject_id`, `ca`.`access`, 1
FROM `fivenet_calendar_access` `ca`
INNER JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `ca`.`user_id`
WHERE `ca`.`user_id` IS NOT NULL
  AND `ca`.`access` > 1;

INSERT IGNORE INTO `fivenet_calendar_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ca`.`target_id`, `sj`.`subject_id`, `ca`.`access`, 1
FROM `fivenet_calendar_access` `ca`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ca`.`job` AND `sj`.`minimum_grade` = `ca`.`minimum_grade`
WHERE `ca`.`job` IS NOT NULL
  AND `ca`.`minimum_grade` IS NOT NULL
  AND `ca`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_wiki_subject_users_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_wiki_subject_users_user_id` (`user_id`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_wiki_subject_users_backfill` (`subject_id`, `user_id`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_user_access`.`user_id`
FROM (
  SELECT DISTINCT `pa`.`user_id`
  FROM `fivenet_wiki_pages_access` `pa`
  LEFT JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `pa`.`user_id`
  WHERE `pa`.`user_id` IS NOT NULL
    AND `su`.`subject_id` IS NULL
  ORDER BY `pa`.`user_id`
) `legacy_user_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 1 FROM `tmp_fivenet_acl_wiki_subject_users_backfill`;

INSERT INTO `fivenet_acl_subject_users` (`subject_id`, `user_id`)
SELECT `subject_id`, `user_id` FROM `tmp_fivenet_acl_wiki_subject_users_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_wiki_subject_users_backfill`;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_wiki_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_wiki_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_wiki_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `pa`.`job`, `pa`.`minimum_grade`
  FROM `fivenet_wiki_pages_access` `pa`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `pa`.`job` AND `sj`.`minimum_grade` = `pa`.`minimum_grade`
  WHERE `pa`.`job` IS NOT NULL
    AND `pa`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `pa`.`job`, `pa`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_wiki_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_wiki_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_wiki_subject_job_grades_backfill`;

INSERT IGNORE INTO `fivenet_wiki_pages_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `pa`.`target_id`, `su`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_wiki_pages_access` `pa`
INNER JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `pa`.`user_id`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4
) `denied_levels`
WHERE `pa`.`user_id` IS NOT NULL
  AND `pa`.`access` = 1;

INSERT IGNORE INTO `fivenet_wiki_pages_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `pa`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_wiki_pages_access` `pa`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `pa`.`job` AND `sj`.`minimum_grade` = `pa`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4
) `denied_levels`
WHERE `pa`.`job` IS NOT NULL
  AND `pa`.`minimum_grade` IS NOT NULL
  AND `pa`.`access` = 1;

INSERT IGNORE INTO `fivenet_wiki_pages_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `pa`.`target_id`, `su`.`subject_id`, `pa`.`access`, 1
FROM `fivenet_wiki_pages_access` `pa`
INNER JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `pa`.`user_id`
WHERE `pa`.`user_id` IS NOT NULL
  AND `pa`.`access` > 1;

INSERT IGNORE INTO `fivenet_wiki_pages_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `pa`.`target_id`, `sj`.`subject_id`, `pa`.`access`, 1
FROM `fivenet_wiki_pages_access` `pa`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `pa`.`job` AND `sj`.`minimum_grade` = `pa`.`minimum_grade`
WHERE `pa`.`job` IS NOT NULL
  AND `pa`.`minimum_grade` IS NOT NULL
  AND `pa`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_template_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_template_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_template_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `ta`.`job`, `ta`.`minimum_grade`
  FROM `fivenet_documents_templates_access` `ta`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `ta`.`job` AND `sj`.`minimum_grade` = `ta`.`minimum_grade`
  WHERE `ta`.`job` IS NOT NULL
    AND `ta`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `ta`.`job`, `ta`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_template_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_template_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_template_subject_job_grades_backfill`;

INSERT IGNORE INTO `fivenet_documents_templates_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ta`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_documents_templates_access` `ta`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ta`.`job` AND `sj`.`minimum_grade` = `ta`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5 UNION ALL
  SELECT 6
) `denied_levels`
WHERE `ta`.`job` IS NOT NULL
  AND `ta`.`minimum_grade` IS NOT NULL
  AND `ta`.`access` = 1;

INSERT IGNORE INTO `fivenet_documents_templates_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ta`.`target_id`, `sj`.`subject_id`, `ta`.`access`, 1
FROM `fivenet_documents_templates_access` `ta`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ta`.`job` AND `sj`.`minimum_grade` = `ta`.`minimum_grade`
WHERE `ta`.`job` IS NOT NULL
  AND `ta`.`minimum_grade` IS NOT NULL
  AND `ta`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_stamp_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_stamp_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_stamp_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `sa`.`job`, `sa`.`minimum_grade`
  FROM `fivenet_documents_stamps_access` `sa`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `sa`.`job` AND `sj`.`minimum_grade` = `sa`.`minimum_grade`
  WHERE `sa`.`job` IS NOT NULL
    AND `sa`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `sa`.`job`, `sa`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_stamp_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_stamp_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_stamp_subject_job_grades_backfill`;

INSERT IGNORE INTO `fivenet_documents_stamps_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `sa`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_documents_stamps_access` `sa`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `sa`.`job` AND `sj`.`minimum_grade` = `sa`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3
) `denied_levels`
WHERE `sa`.`job` IS NOT NULL
  AND `sa`.`minimum_grade` IS NOT NULL
  AND `sa`.`access` = 1;

INSERT IGNORE INTO `fivenet_documents_stamps_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `sa`.`target_id`, `sj`.`subject_id`, `sa`.`access`, 1
FROM `fivenet_documents_stamps_access` `sa`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `sa`.`job` AND `sj`.`minimum_grade` = `sa`.`minimum_grade`
WHERE `sa`.`job` IS NOT NULL
  AND `sa`.`minimum_grade` IS NOT NULL
  AND `sa`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_user_label_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_user_label_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_user_label_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `la`.`job`, `la`.`minimum_grade`
  FROM `fivenet_user_labels_job_job_access` `la`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `la`.`job` AND `sj`.`minimum_grade` = `la`.`minimum_grade`
  WHERE `la`.`job` IS NOT NULL
    AND `la`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `la`.`job`, `la`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_user_label_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_user_label_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_user_label_subject_job_grades_backfill`;

INSERT IGNORE INTO `fivenet_user_labels_job_job_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `la`.`target_id`, `sj`.`subject_id`, `la`.`access`, 1
FROM `fivenet_user_labels_job_job_access` `la`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `la`.`job` AND `sj`.`minimum_grade` = `la`.`minimum_grade`
WHERE `la`.`job` IS NOT NULL
  AND `la`.`minimum_grade` IS NOT NULL
  AND `la`.`access` > 0;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_qualification_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_qualification_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_qualification_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `qa`.`job`, `qa`.`minimum_grade`
  FROM `fivenet_qualifications_access` `qa`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `qa`.`job` AND `sj`.`minimum_grade` = `qa`.`minimum_grade`
  WHERE `qa`.`job` IS NOT NULL
    AND `qa`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `qa`.`job`, `qa`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_qualification_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_qualification_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_qualification_subject_job_grades_backfill`;

INSERT IGNORE INTO `fivenet_qualifications_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `qa`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_qualifications_access` `qa`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `qa`.`job` AND `sj`.`minimum_grade` = `qa`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5 UNION ALL
  SELECT 6
) `denied_levels`
WHERE `qa`.`job` IS NOT NULL
  AND `qa`.`minimum_grade` IS NOT NULL
  AND `qa`.`access` = 1;

INSERT IGNORE INTO `fivenet_qualifications_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `qa`.`target_id`, `sj`.`subject_id`, `qa`.`access`, 1
FROM `fivenet_qualifications_access` `qa`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `qa`.`job` AND `sj`.`minimum_grade` = `qa`.`minimum_grade`
WHERE `qa`.`job` IS NOT NULL
  AND `qa`.`minimum_grade` IS NOT NULL
  AND `qa`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_mailer_subject_users_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_mailer_subject_users_user_id` (`user_id`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_mailer_subject_users_backfill` (`subject_id`, `user_id`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_user_access`.`user_id`
FROM (
  SELECT DISTINCT `ma`.`user_id`
  FROM `fivenet_mailer_emails_access` `ma`
  LEFT JOIN `fivenet_acl_subject_users` `su` ON `su`.`user_id` = `ma`.`user_id`
  WHERE `ma`.`user_id` IS NOT NULL
    AND `su`.`subject_id` IS NULL
  ORDER BY `ma`.`user_id`
) `legacy_user_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 1 FROM `tmp_fivenet_acl_mailer_subject_users_backfill`;

INSERT INTO `fivenet_acl_subject_users` (`subject_id`, `user_id`)
SELECT `subject_id`, `user_id` FROM `tmp_fivenet_acl_mailer_subject_users_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_mailer_subject_users_backfill`;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_mailer_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_mailer_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_mailer_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `ma`.`job`, `ma`.`minimum_grade`
  FROM `fivenet_mailer_emails_access` `ma`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `ma`.`job` AND `sj`.`minimum_grade` = `ma`.`minimum_grade`
  WHERE `ma`.`job` IS NOT NULL
    AND `ma`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `ma`.`job`, `ma`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_mailer_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_mailer_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_mailer_subject_job_grades_backfill`;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_mailer_subject_qualifications_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_mailer_subject_qualification` (`qualification_id`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_mailer_subject_qualifications_backfill` (`subject_id`, `qualification_id`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_qualification_access`.`qualification_id`
FROM (
  SELECT DISTINCT `ma`.`qualification_id`
  FROM `fivenet_mailer_emails_access` `ma`
  LEFT JOIN `fivenet_acl_subject_qualifications` `sq` ON `sq`.`qualification_id` = `ma`.`qualification_id`
  WHERE `ma`.`qualification_id` IS NOT NULL
    AND `sq`.`subject_id` IS NULL
  ORDER BY `ma`.`qualification_id`
) `legacy_qualification_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 2 FROM `tmp_fivenet_acl_mailer_subject_qualifications_backfill`;

INSERT INTO `fivenet_acl_subject_qualifications` (`subject_id`, `qualification_id`)
SELECT `subject_id`, `qualification_id` FROM `tmp_fivenet_acl_mailer_subject_qualifications_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_mailer_subject_qualifications_backfill`;

INSERT IGNORE INTO `fivenet_mailer_emails_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ma`.`target_id`, `subject_map`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_mailer_emails_access` `ma`
INNER JOIN (
  SELECT `user_id` AS `source_id`, `subject_id`, 1 AS `subject_type` FROM `fivenet_acl_subject_users`
  UNION ALL
  SELECT `qualification_id` AS `source_id`, `subject_id`, 2 AS `subject_type` FROM `fivenet_acl_subject_qualifications`
) `subject_map` ON (
  (`subject_map`.`subject_type` = 1 AND `subject_map`.`source_id` = `ma`.`user_id`) OR
  (`subject_map`.`subject_type` = 2 AND `subject_map`.`source_id` = `ma`.`qualification_id`)
)
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4
) `denied_levels`
WHERE (`ma`.`user_id` IS NOT NULL OR `ma`.`qualification_id` IS NOT NULL)
  AND `ma`.`access` = 1;

INSERT IGNORE INTO `fivenet_mailer_emails_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ma`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_mailer_emails_access` `ma`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ma`.`job` AND `sj`.`minimum_grade` = `ma`.`minimum_grade`
INNER JOIN (
  SELECT 2 AS `access` UNION ALL
  SELECT 3 UNION ALL
  SELECT 4
) `denied_levels`
WHERE `ma`.`job` IS NOT NULL
  AND `ma`.`minimum_grade` IS NOT NULL
  AND `ma`.`access` = 1;

INSERT IGNORE INTO `fivenet_mailer_emails_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ma`.`target_id`, `subject_map`.`subject_id`, `ma`.`access`, 1
FROM `fivenet_mailer_emails_access` `ma`
INNER JOIN (
  SELECT `user_id` AS `source_id`, `subject_id`, 1 AS `subject_type` FROM `fivenet_acl_subject_users`
  UNION ALL
  SELECT `qualification_id` AS `source_id`, `subject_id`, 2 AS `subject_type` FROM `fivenet_acl_subject_qualifications`
) `subject_map` ON (
  (`subject_map`.`subject_type` = 1 AND `subject_map`.`source_id` = `ma`.`user_id`) OR
  (`subject_map`.`subject_type` = 2 AND `subject_map`.`source_id` = `ma`.`qualification_id`)
)
WHERE (`ma`.`user_id` IS NOT NULL OR `ma`.`qualification_id` IS NOT NULL)
  AND `ma`.`access` > 1;

INSERT IGNORE INTO `fivenet_mailer_emails_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ma`.`target_id`, `sj`.`subject_id`, `ma`.`access`, 1
FROM `fivenet_mailer_emails_access` `ma`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ma`.`job` AND `sj`.`minimum_grade` = `ma`.`minimum_grade`
WHERE `ma`.`job` IS NOT NULL
  AND `ma`.`minimum_grade` IS NOT NULL
  AND `ma`.`access` > 1;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_unit_subject_job_grades_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `job` varchar(50) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_unit_subject_job_grade` (`job`, `minimum_grade`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_unit_subject_job_grades_backfill` (`subject_id`, `job`, `minimum_grade`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_job_access`.`job`, `legacy_job_access`.`minimum_grade`
FROM (
  SELECT DISTINCT `ua`.`job`, `ua`.`minimum_grade`
  FROM `fivenet_centrum_units_access` `ua`
  LEFT JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
    ON `sj`.`job` = `ua`.`job` AND `sj`.`minimum_grade` = `ua`.`minimum_grade`
  WHERE `ua`.`job` IS NOT NULL
    AND `ua`.`minimum_grade` IS NOT NULL
    AND `sj`.`subject_id` IS NULL
  ORDER BY `ua`.`job`, `ua`.`minimum_grade`
) `legacy_job_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 3 FROM `tmp_fivenet_acl_unit_subject_job_grades_backfill`;

INSERT INTO `fivenet_acl_subject_job_grade_scopes` (`subject_id`, `job`, `minimum_grade`)
SELECT `subject_id`, `job`, `minimum_grade` FROM `tmp_fivenet_acl_unit_subject_job_grades_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_unit_subject_job_grades_backfill`;

CREATE TEMPORARY TABLE `tmp_fivenet_acl_unit_subject_qualifications_backfill` (
  `subject_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`subject_id`),
  UNIQUE KEY `idx_tmp_fivenet_acl_unit_subject_qualification` (`qualification_id`)
) ENGINE=MEMORY;

SET @next_acl_subject_id := (SELECT COALESCE(MAX(`id`), 0) FROM `fivenet_acl_subjects`);

INSERT INTO `tmp_fivenet_acl_unit_subject_qualifications_backfill` (`subject_id`, `qualification_id`)
SELECT (@next_acl_subject_id := @next_acl_subject_id + 1), `legacy_qualification_access`.`qualification_id`
FROM (
  SELECT DISTINCT `ua`.`qualification_id`
  FROM `fivenet_centrum_units_access` `ua`
  LEFT JOIN `fivenet_acl_subject_qualifications` `sq` ON `sq`.`qualification_id` = `ua`.`qualification_id`
  WHERE `ua`.`qualification_id` IS NOT NULL
    AND `sq`.`subject_id` IS NULL
  ORDER BY `ua`.`qualification_id`
) `legacy_qualification_access`;

INSERT INTO `fivenet_acl_subjects` (`id`, `subject_type`)
SELECT `subject_id`, 2 FROM `tmp_fivenet_acl_unit_subject_qualifications_backfill`;

INSERT INTO `fivenet_acl_subject_qualifications` (`subject_id`, `qualification_id`)
SELECT `subject_id`, `qualification_id` FROM `tmp_fivenet_acl_unit_subject_qualifications_backfill`;

DROP TEMPORARY TABLE `tmp_fivenet_acl_unit_subject_qualifications_backfill`;

INSERT IGNORE INTO `fivenet_centrum_units_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ua`.`target_id`, `subject_map`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_centrum_units_access` `ua`
INNER JOIN (
  SELECT `qualification_id` AS `source_id`, `subject_id` FROM `fivenet_acl_subject_qualifications`
) `subject_map` ON `subject_map`.`source_id` = `ua`.`qualification_id`
INNER JOIN (SELECT 2 AS `access`) `denied_levels`
WHERE `ua`.`qualification_id` IS NOT NULL
  AND `ua`.`access` = 1;

INSERT IGNORE INTO `fivenet_centrum_units_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ua`.`target_id`, `sj`.`subject_id`, `denied_levels`.`access`, 0
FROM `fivenet_centrum_units_access` `ua`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ua`.`job` AND `sj`.`minimum_grade` = `ua`.`minimum_grade`
INNER JOIN (SELECT 2 AS `access`) `denied_levels`
WHERE `ua`.`job` IS NOT NULL
  AND `ua`.`minimum_grade` IS NOT NULL
  AND `ua`.`access` = 1;

INSERT IGNORE INTO `fivenet_centrum_units_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ua`.`target_id`, `subject_map`.`subject_id`, `ua`.`access`, 1
FROM `fivenet_centrum_units_access` `ua`
INNER JOIN (
  SELECT `qualification_id` AS `source_id`, `subject_id` FROM `fivenet_acl_subject_qualifications`
) `subject_map` ON `subject_map`.`source_id` = `ua`.`qualification_id`
WHERE `ua`.`qualification_id` IS NOT NULL
  AND `ua`.`access` > 1;

INSERT IGNORE INTO `fivenet_centrum_units_access_v2` (`target_id`, `subject_id`, `access`, `effect`)
SELECT `ua`.`target_id`, `sj`.`subject_id`, `ua`.`access`, 1
FROM `fivenet_centrum_units_access` `ua`
INNER JOIN `fivenet_acl_subject_job_grade_scopes` `sj`
  ON `sj`.`job` = `ua`.`job` AND `sj`.`minimum_grade` = `ua`.`minimum_grade`
WHERE `ua`.`job` IS NOT NULL
  AND `ua`.`minimum_grade` IS NOT NULL
  AND `ua`.`access` > 1;

-- Drop old access tables and rename the new ones
DROP TABLE IF EXISTS `fivenet_calendar_access`;
DROP TABLE IF EXISTS `fivenet_centrum_units_access`;
DROP TABLE IF EXISTS `fivenet_documents_access`;
DROP TABLE IF EXISTS `fivenet_documents_stamps_access`;
DROP TABLE IF EXISTS `fivenet_documents_templates_access`;
DROP TABLE IF EXISTS `fivenet_mailer_emails_access`;
DROP TABLE IF EXISTS `fivenet_qualifications_access`;
DROP TABLE IF EXISTS `fivenet_user_labels_job_job_access`;
DROP TABLE IF EXISTS `fivenet_wiki_pages_access`;

RENAME TABLE `fivenet_calendar_access_v2` TO `fivenet_calendar_access`;
RENAME TABLE `fivenet_centrum_units_access_v2` TO `fivenet_centrum_units_access`;
RENAME TABLE `fivenet_documents_access_v2` TO `fivenet_documents_access`;
RENAME TABLE `fivenet_documents_stamps_access_v2` TO `fivenet_documents_stamps_access`;
RENAME TABLE `fivenet_documents_templates_access_v2` TO `fivenet_documents_templates_access`;
RENAME TABLE `fivenet_mailer_emails_access_v2` TO `fivenet_mailer_emails_access`;
RENAME TABLE `fivenet_qualifications_access_v2` TO `fivenet_qualifications_access`;
RENAME TABLE `fivenet_user_labels_job_job_access_v2` TO `fivenet_user_labels_job_job_access`;
RENAME TABLE `fivenet_wiki_pages_access_v2` TO `fivenet_wiki_pages_access`;

-- Table: fivenet_documents_templates - Convert access enum to int32
UPDATE `fivenet_documents_templates`
SET `access` = REPLACE(
  REPLACE(
    REPLACE(
      REPLACE(
        REPLACE(
          REPLACE(
            REPLACE(
              REPLACE(
                REPLACE(
                  REPLACE(
                    REPLACE(
                      REPLACE(
                        REPLACE(
                          REPLACE(`access`, '"access":"ACCESS_LEVEL_UNSPECIFIED"', '"access":0'),
                          '"access": "ACCESS_LEVEL_UNSPECIFIED"', '"access":0'
                        ),
                        '"access":"ACCESS_LEVEL_BLOCKED"', '"access":1'
                      ),
                      '"access": "ACCESS_LEVEL_BLOCKED"', '"access":1'
                    ),
                    '"access":"ACCESS_LEVEL_VIEW"', '"access":2'
                  ),
                  '"access": "ACCESS_LEVEL_VIEW"', '"access":2'
                ),
                '"access":"ACCESS_LEVEL_COMMENT"', '"access":3'
              ),
              '"access": "ACCESS_LEVEL_COMMENT"', '"access":3'
            ),
            '"access":"ACCESS_LEVEL_STATUS"', '"access":4'
          ),
          '"access": "ACCESS_LEVEL_STATUS"', '"access":4'
        ),
        '"access":"ACCESS_LEVEL_ACCESS"', '"access":5'
      ),
      '"access": "ACCESS_LEVEL_ACCESS"', '"access":5'
    ),
    '"access":"ACCESS_LEVEL_EDIT"', '"access":6'
  ),
  '"access": "ACCESS_LEVEL_EDIT"', '"access":6'
)
WHERE `access` LIKE '%ACCESS_LEVEL_%';

COMMIT;
