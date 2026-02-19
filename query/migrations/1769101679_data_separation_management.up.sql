BEGIN;

ALTER TABLE `fivenet_user` ADD KEY `idx_phone_number` (`phone_number`); -- Add index on phone_number

-- Table: `fivenet_user_licenses` - Drop dependent constraints and indexes
ALTER TABLE `fivenet_user_licenses` DROP CONSTRAINT `fk_fivenet_user_licenses_owner`;
ALTER TABLE `fivenet_user_licenses` DROP INDEX `fivenet_user_licenses_owner_IDX`;

ALTER TABLE `fivenet_user_licenses` DROP CONSTRAINT `fk_fivenet_user_licenses_type`;
ALTER TABLE `fivenet_user_licenses` DROP INDEX `PRIMARY`;

-- Table: Remove foreign keys and indexes depending on `fivenet_user`
ALTER TABLE `fivenet_calendar` DROP FOREIGN KEY `fk_fivenet_calendar_creator_id`;
ALTER TABLE `fivenet_calendar_access` DROP FOREIGN KEY `fk_fivenet_calendar_access_user_id`;
ALTER TABLE `fivenet_calendar_entries` DROP FOREIGN KEY `fk_fivenet_calendar_entries_creator_id`;
ALTER TABLE `fivenet_calendar_rsvp` DROP FOREIGN KEY `fk_fivenet_calendar_rsvp_user_id`;
ALTER TABLE `fivenet_calendar_subs` DROP FOREIGN KEY `fk_fivenet_calendar_subs_user_id`;
ALTER TABLE `fivenet_centrum_dispatchers` DROP FOREIGN KEY `fk_fivenet_centrum_dispatchers_user_id`;
ALTER TABLE `fivenet_centrum_dispatches` DROP FOREIGN KEY `fk_fivenet_centrum_dispatches_creator_id`;
ALTER TABLE `fivenet_centrum_dispatches_status` DROP FOREIGN KEY `fk_fivenet_centrum_dispatches_status_user_id`;
ALTER TABLE `fivenet_centrum_markers` DROP FOREIGN KEY `fk_fivenet_centrum_markers_creator_id`;
ALTER TABLE `fivenet_centrum_units_status` DROP FOREIGN KEY `fk_fivenet_centrum_units_status_creator_id`;
ALTER TABLE `fivenet_centrum_units_status` DROP FOREIGN KEY `fk_fivenet_centrum_units_status_user_id`;
ALTER TABLE `fivenet_centrum_units_users` DROP FOREIGN KEY `fk_fivenet_centrum_units_users_user_id`;
ALTER TABLE `fivenet_documents` DROP FOREIGN KEY `fk_fivenet_documents_creator_id`;
ALTER TABLE `fivenet_documents_access` DROP FOREIGN KEY `fk_fivenet_documents_access_user_id`;
ALTER TABLE `fivenet_documents_activity` DROP FOREIGN KEY `fk_fivenet_documents_activity_creator_id`;
ALTER TABLE `fivenet_documents_approval_tasks` DROP FOREIGN KEY `fk_fivenet_doc_apptsk_task_creator_id`;
ALTER TABLE `fivenet_documents_approval_tasks` DROP FOREIGN KEY `fk_fivenet_doc_apptsk_task_user_id`;
ALTER TABLE `fivenet_documents_approvals` DROP FOREIGN KEY `fk_fivenet_doc_approvals_user_id`;
ALTER TABLE `fivenet_documents_comments` DROP FOREIGN KEY `fk_fivenet_documents_comments_creator_id`;
ALTER TABLE `fivenet_documents_pins` DROP FOREIGN KEY `fk_fivenet_documents_pins_creator_id`;
ALTER TABLE `fivenet_documents_pins` DROP FOREIGN KEY `fk_fivenet_documents_pins_user_id`;
ALTER TABLE `fivenet_documents_references` DROP FOREIGN KEY `fk_fivenet_documents_references_creator_id`;
ALTER TABLE `fivenet_documents_relations` DROP FOREIGN KEY `fk_fivenet_documents_relations_source_user_id`;
ALTER TABLE `fivenet_documents_relations` DROP FOREIGN KEY `fk_fivenet_documents_relations_target_user_id`;
ALTER TABLE `fivenet_documents_requests` DROP FOREIGN KEY `fk_fivenet_documents_requests_creator_id`;
ALTER TABLE `fivenet_documents_stamps` DROP FOREIGN KEY `fk_fivenet_documents_signatures_stamp_user`;
ALTER TABLE `fivenet_documents_workflow_users` DROP FOREIGN KEY `fk_fivenet_documents_workflow_users_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` DROP FOREIGN KEY `fk_fivenet_job_colleague_activity_source_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` DROP FOREIGN KEY `fk_fivenet_job_colleague_activity_target_user_id`;
ALTER TABLE `fivenet_job_colleague_labels` DROP FOREIGN KEY `fk_fivenet_job_colleague_labels_user_id`;
ALTER TABLE `fivenet_job_colleague_props` DROP FOREIGN KEY `fk_fivenet_job_colleague_props_user_id`;
ALTER TABLE `fivenet_job_conduct` DROP FOREIGN KEY `fk_fivenet_job_conduct_creator_id`;
ALTER TABLE `fivenet_job_conduct` DROP FOREIGN KEY `fk_fivenet_job_conduct_target_user_id`;
ALTER TABLE `fivenet_job_timeclock` DROP FOREIGN KEY `fk_fivenet_job_timeclock_user_id`;
ALTER TABLE `fivenet_mailer_emails` DROP FOREIGN KEY `fk_fivenet_mailer_emails_user_id`;
ALTER TABLE `fivenet_mailer_emails_access` DROP FOREIGN KEY `fk_fivenet_mailer_emails_access_user_id`;
ALTER TABLE `fivenet_mailer_templates` DROP FOREIGN KEY `fk_fivenet_mailer_templates_creator_id`;
ALTER TABLE `fivenet_notifications` DROP FOREIGN KEY `fk_fivenet_notifications_user_id`;
ALTER TABLE `fivenet_qualifications` DROP FOREIGN KEY `fk_fivenet_qualifications_creator_id`;
ALTER TABLE `fivenet_qualifications_exam_responses` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_responses_user_id`;
ALTER TABLE `fivenet_qualifications_exam_users` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_users_user_id`;
ALTER TABLE `fivenet_qualifications_requests` DROP FOREIGN KEY `fk_fivenet_qualifications_requests_approver_id`;
ALTER TABLE `fivenet_qualifications_requests` DROP FOREIGN KEY `fk_fivenet_qualifications_requests_user_id`;
ALTER TABLE `fivenet_qualifications_results` DROP FOREIGN KEY `fk_fivenet_qualifications_results_creator_id`;
ALTER TABLE `fivenet_qualifications_results` DROP FOREIGN KEY `fk_fivenet_qualifications_results_user_id`;
ALTER TABLE `fivenet_user_activity` DROP FOREIGN KEY `fk_fivenet_user_activity_source_user_id`;
ALTER TABLE `fivenet_user_activity` DROP FOREIGN KEY `fk_fivenet_user_activity_target_user_id`;
ALTER TABLE `fivenet_user_labels` DROP FOREIGN KEY `fk_fivenet_user_labels_user_id`;
ALTER TABLE `fivenet_user_props` DROP FOREIGN KEY `fk_fivenet_user_props_user_id`;
ALTER TABLE `fivenet_wiki_pages` DROP FOREIGN KEY `fk_fivenet_wiki_pages_creator_id`;
ALTER TABLE `fivenet_wiki_pages_access` DROP FOREIGN KEY `fk_fivenet_wiki_pages_access_user_id`;
ALTER TABLE `fivenet_wiki_pages_activity` DROP FOREIGN KEY `fk_fivenet_wiki_pages_activity_creator_id`;

ALTER TABLE `fivenet_centrum_user_locations` DROP FOREIGN KEY `fk_fivenet_centrum_user_locations_identifier`;

-- Table: `fivenet_user` - Modify primary key on fivenet_user to be `id` instead of `identifier`
ALTER TABLE `fivenet_user` MODIFY COLUMN `id` int NOT NULL;

-- Table: `fivenet_user` - `identifier` column was primary key
ALTER TABLE `fivenet_user` DROP INDEX `PRIMARY`;

-- Table: `fivenet_user` - Set `id` as primary key
ALTER TABLE `fivenet_user` ADD PRIMARY KEY (`id`);
-- Table: `fivenet_user` - `id` column was unique key
ALTER TABLE `fivenet_user` DROP INDEX `idx_id`;
ALTER TABLE  `fivenet_user` MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT;

-- Table: `fivenet_user` - Set `identifier` as unique key
ALTER TABLE `fivenet_user` ADD UNIQUE KEY `idx_identifier` (`identifier`);

-- Table: `fivenet_user` - Add `account_id` column
ALTER TABLE `fivenet_user` ADD COLUMN `account_id` bigint(20) unsigned AFTER `id`;
-- Table: `fivenet_user` - Add index on `account_id` (can't be unique as multiple users can belong to same account)
ALTER TABLE `fivenet_user` ADD KEY `idx_account_id` (`account_id`);
-- Table: `fivenet_user` - Add foreign key constraint to `fivenet_account`
ALTER TABLE `fivenet_user` ADD CONSTRAINT `fk_fivenet_user_account_id` FOREIGN KEY (`account_id`) REFERENCES `fivenet_accounts` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- Table: `fivenet_user` - Add `deleted_at` column for soft deletes
ALTER TABLE `fivenet_user` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `last_seen`;
ALTER TABLE `fivenet_user` ADD COLUMN `deleted_reason` VARCHAR(255) DEFAULT NULL AFTER `deleted_at`;

-- Table: `fivenet_user` - Migrate existing users from `users` to `fivenet_user` table
INSERT INTO `fivenet_user` (
    id,
    identifier,
    `group`,
    job,
    job_grade,
    firstname,
    lastname,
    dateofbirth,
    sex,
    height,
    phone_number,
    disabled,
    visum,
    playtime,
    created_at,
    last_seen
)
SELECT
    id,
    identifier,
    `group`,
    job,
    job_grade,
    firstname,
    lastname,
    dateofbirth,
    sex,
    height,
    phone_number,
    disabled,
    visum,
    playtime,
    created_at,
    last_seen
FROM `users`;

-- Table: `fivenet_user` - Fill in `account_id` for existing users with "matching" accounts
UPDATE `fivenet_user` u
JOIN `fivenet_accounts` a ON SUBSTRING_INDEX(u.`identifier`, ':', -1) = a.`license`
SET u.`account_id` = a.`id`;

-- Table: Recreate foreign keys and indexes depending on `fivenet_user`
ALTER TABLE `fivenet_calendar` ADD CONSTRAINT `fk_fivenet_calendar_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_access` ADD CONSTRAINT `fk_fivenet_calendar_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_entries` ADD CONSTRAINT `fk_fivenet_calendar_entries_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_rsvp` ADD CONSTRAINT `fk_fivenet_calendar_rsvp_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_calendar_subs` ADD CONSTRAINT `fk_fivenet_calendar_subs_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_dispatchers` ADD CONSTRAINT `fk_fivenet_centrum_dispatchers_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_dispatches` ADD CONSTRAINT `fk_fivenet_centrum_dispatches_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_dispatches_status` ADD CONSTRAINT `fk_fivenet_centrum_dispatches_status_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_markers` ADD CONSTRAINT `fk_fivenet_centrum_markers_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_units_status` ADD CONSTRAINT `fk_fivenet_centrum_units_status_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_units_status` ADD CONSTRAINT `fk_fivenet_centrum_units_status_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_units_users` ADD CONSTRAINT `fk_fivenet_centrum_units_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents` ADD CONSTRAINT `fk_fivenet_documents_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_access` ADD CONSTRAINT `fk_fivenet_documents_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_activity` ADD CONSTRAINT `fk_fivenet_documents_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_approval_tasks` ADD CONSTRAINT `fk_fivenet_doc_apptsk_task_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_approval_tasks` ADD CONSTRAINT `fk_fivenet_doc_apptsk_task_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_approvals` ADD CONSTRAINT `fk_fivenet_doc_approvals_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_comments` ADD CONSTRAINT `fk_fivenet_documents_comments_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_pins` ADD CONSTRAINT `fk_fivenet_documents_pins_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_pins` ADD CONSTRAINT `fk_fivenet_documents_pins_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_references` ADD CONSTRAINT `fk_fivenet_documents_references_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_relations` ADD CONSTRAINT `fk_fivenet_documents_relations_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_relations` ADD CONSTRAINT `fk_fivenet_documents_relations_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_requests` ADD CONSTRAINT `fk_fivenet_documents_requests_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_stamps` ADD CONSTRAINT `fk_fivenet_documents_signatures_stamp_user` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents_workflow_users` ADD CONSTRAINT `fk_fivenet_documents_workflow_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_colleague_activity` ADD CONSTRAINT `fk_fivenet_job_colleague_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_job_colleague_activity` ADD CONSTRAINT `fk_fivenet_job_colleague_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_colleague_labels` ADD CONSTRAINT `fk_fivenet_job_colleague_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_colleague_props` ADD CONSTRAINT `fk_fivenet_job_colleague_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_conduct` ADD CONSTRAINT `fk_fivenet_job_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_conduct` ADD CONSTRAINT `fk_fivenet_job_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_timeclock` ADD CONSTRAINT `fk_fivenet_job_timeclock_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_emails` ADD CONSTRAINT `fk_fivenet_mailer_emails_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_emails_access` ADD CONSTRAINT `fk_fivenet_mailer_emails_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_templates` ADD CONSTRAINT `fk_fivenet_mailer_templates_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_notifications` ADD CONSTRAINT `fk_fivenet_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications` ADD CONSTRAINT `fk_fivenet_qualifications_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_responses` ADD CONSTRAINT `fk_fivenet_qualifications_exam_responses_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_users` ADD CONSTRAINT `fk_fivenet_qualifications_exam_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_requests` ADD CONSTRAINT `fk_fivenet_qualifications_requests_approver_id` FOREIGN KEY (`approver_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_requests` ADD CONSTRAINT `fk_fivenet_qualifications_requests_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_results` ADD CONSTRAINT `fk_fivenet_qualifications_results_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_results` ADD CONSTRAINT `fk_fivenet_qualifications_results_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_activity` ADD CONSTRAINT `fk_fivenet_user_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_user_activity` ADD CONSTRAINT `fk_fivenet_user_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_labels` ADD CONSTRAINT `fk_fivenet_user_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_props` ADD CONSTRAINT `fk_fivenet_user_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_pages` ADD CONSTRAINT `fk_fivenet_wiki_pages_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_pages_access` ADD CONSTRAINT `fk_fivenet_wiki_pages_access_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_pages_activity` ADD CONSTRAINT `fk_fivenet_wiki_pages_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `fivenet_user` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

-- Table: `fivenet_user_licenses` - Create new `user_id` column to replace `owner` in fivenet_user_licenses
ALTER TABLE `fivenet_user_licenses` ADD COLUMN `user_id` int(11) NOT NULL FIRST;
ALTER TABLE `fivenet_user_licenses` ADD PRIMARY KEY (`type`, `user_id`);

INSERT
	INTO
	fivenet_user_licenses (`user_id`,
	`type`,
	`owner`)
SELECT
	u.id,
	ul.`type`,
	ul.`owner`
FROM
	user_licenses ul
INNER JOIN users u ON (u.identifier = ul.`owner`);

-- Table: `fivenet_user_licenses` - Migrate `owner` (identifier) to `user_id`
UPDATE `fivenet_user_licenses` ul
JOIN `fivenet_user` u ON ul.`owner` = u.`identifier`
SET ul.`user_id` = u.`id`;

-- Table: `fivenet_user_licenses` - Remove "broken" records where no matching user was found..
DELETE FROM `fivenet_user_licenses` WHERE `user_id` = 0;
ALTER TABLE `fivenet_user_licenses` ADD CONSTRAINT `fk_fivenet_user_licenses_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_user_licenses` - Recreate foreign key for `type`
ALTER TABLE `fivenet_user_licenses` ADD CONSTRAINT `fk_fivenet_user_licenses_type` FOREIGN KEY (`type`) REFERENCES `fivenet_licenses` (`type`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_owned_vehicles` - Drop dependent constraints and indexes
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_ownerplate`;
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_owner`;
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_owner_type`;
ALTER TABLE `fivenet_owned_vehicles` DROP INDEX `idx_fivenet_owned_vehicles_owner_model_type`;

ALTER TABLE `fivenet_owned_vehicles` ADD COLUMN `user_id` int(11) NULL AFTER `owner`;

INSERT INTO `fivenet_owned_vehicles` (`owner`, `user_id`, `job`, `plate`, `model`, `type`)
SELECT `owner`, NULL, NULL, `plate`, `model`, `type` FROM `owned_vehicles`;

-- Table: `fivenet_owned_vehicles` - Migrate `owner` (identifier) to `user_id`
UPDATE `fivenet_owned_vehicles` owv
JOIN `fivenet_user` u ON owv.`owner` = u.`identifier`
SET owv.`user_id` = u.`id`, owv.`owner` = NULL;

-- Table: `fivenet_owned_vehicles` - Remove "broken" records where no matching user was found..
DELETE FROM `fivenet_owned_vehicles` WHERE `user_id` = 0;

-- Table: `fivenet_owned_vehicles` - Recreate indexes with `user_id` instead of `owner`
ALTER TABLE `fivenet_owned_vehicles` ADD UNIQUE KEY `idx_fivenet_owned_vehicles_userplate` (`user_id`, `plate`);
ALTER TABLE `fivenet_owned_vehicles` ADD UNIQUE KEY `idx_fivenet_owned_vehicles_jobplate` (`job`, `plate`);
ALTER TABLE `fivenet_owned_vehicles` ADD KEY `idx_fivenet_owned_vehicles_user` (`user_id`);
ALTER TABLE `fivenet_owned_vehicles` ADD KEY `idx_fivenet_owned_vehicles_user_type` (`user_id`, `type`);
ALTER TABLE `fivenet_owned_vehicles` ADD KEY `idx_fivenet_owned_vehicles_user_model_type` (`user_id`, `model`, `type`);

ALTER TABLE `fivenet_owned_vehicles` ADD CONSTRAINT `fk_fivenet_owned_vehicles_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_user_phone_numbers`
CREATE TABLE IF NOT EXISTS `fivenet_user_phone_numbers` (
  `user_id` int(11) NOT NULL,
  `phone_number` varchar(15) NOT NULL,
  `is_primary` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `phone_number`),
  UNIQUE KEY `idx_phone_number` (`phone_number`),
  KEY `idx_is_primary` (`is_primary`),
  CONSTRAINT `fk_fivenet_user_phone_numbers_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: `fivenet_user_phone_numbers` - Add all current `phone_number`s as primary numbers in `fivenet_user_phone_numbers`
INSERT INTO `fivenet_user_phone_numbers` (`user_id`, `phone_number`, `is_primary`)
SELECT u.`id` AS `user_id`, u.`phone_number` AS `phone_number`, 1 AS `is_primary`
FROM `fivenet_user` u
WHERE u.`phone_number` IS NOT NULL AND u.`phone_number` != '';

-- Table: `fivenet_user_jobs`
CREATE TABLE IF NOT EXISTS `fivenet_user_jobs` (
  `user_id` int(11) NOT NULL,
  `job` varchar(50) NOT NULL,
  `grade` int(11) NOT NULL DEFAULT '0',
  `is_primary` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `job`),
  KEY `idx_job` (`job`),
  KEY `idx_is_primary` (`is_primary`),
  CONSTRAINT `fk_fivenet_user_jobs_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: `fivenet_user_jobs` - Add all current `job`s as primary job in `fivenet_user_jobs`
INSERT INTO `fivenet_user_jobs` (`user_id`, `job`, `grade`, `is_primary`)
SELECT u.`id` AS `user_id`, u.`job` AS `job`,  u.`job_grade` AS `grade`, 1 AS `is_primary`
FROM `fivenet_user` u
WHERE u.`job` IS NOT NULL AND u.`job` != '';

-- Table: `fivenet_user_licenses` - Drop old `owner` column from table
ALTER TABLE `fivenet_user_licenses` DROP COLUMN `owner`;
-- Table: `fivenet_owned_vehicles` - Drop old `owner` column from table
ALTER TABLE `fivenet_owned_vehicles` DROP COLUMN `owner`;

-- Table: `fivenet_accounts` - Drop old job override and `superuser` columns, job overrides and superuser state are determined by the token
ALTER TABLE `fivenet_accounts` DROP COLUMN `override_job`;
ALTER TABLE `fivenet_accounts` DROP COLUMN `override_job_grade`;
ALTER TABLE `fivenet_accounts` DROP COLUMN `superuser`;

-- Table: `fivenet_accounts` - Add `groups` column to store JSON array of group names (this allows for multiple groups per account instead of per user as before)
ALTER TABLE `fivenet_accounts` ADD COLUMN `groups` varchar(255) DEFAULT NULL AFTER `license`;

-- MariaDB doesn't support JSON multi-valued indexes, so we skip this part
-- It might land in MariaDB 12.2 or later, see https://jira.mariadb.org/browse/MDEV-25848
SET @ver := VERSION();

-- Build either the ALTER or a no-op SELECT
SET @ddl :=
  IF(
    LOCATE('MariaDB', @ver) = 0,
    'ALTER TABLE fivenet_accounts
       ADD INDEX idx_groups
         ((CAST(`groups` AS CHAR(255) ARRAY)));',
    CONCAT(
      'SELECT ''Skipping JSON multi-valued index on MariaDB: ',
      @ver,
      ''';'
    )
  );

PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Table: `fivenet_centrum_user_locations` - Add `user_id` column and foreign key constraint to `fivenet_user`
TRUNCATE `fivenet_centrum_user_locations`;

ALTER TABLE `fivenet_centrum_user_locations` DROP PRIMARY KEY;
ALTER TABLE `fivenet_centrum_user_locations` DROP COLUMN `identifier`;

ALTER TABLE `fivenet_centrum_user_locations` ADD COLUMN `user_id` int(11) NOT NULL FIRST;
ALTER TABLE `fivenet_centrum_user_locations` ADD PRIMARY KEY (`user_id`);
ALTER TABLE `fivenet_centrum_user_locations` ADD CONSTRAINT `fk_fivenet_centrum_user_locations_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- Table: `fivenet_user_props`
ALTER TABLE `fivenet_user_props` DROP COLUMN `attributes`;

-- Table: `fivenet_documents_stamps` - Drop `user_id` column as stamps are per job
ALTER TABLE `fivenet_documents_stamps` DROP FOREIGN KEY `fk_fivenet_documents_signatures_stamp_user`;
ALTER TABLE `fivenet_documents_stamps` DROP COLUMN `user_id`;

-- Table: `fivenet_documents_meta` - Add comment count
ALTER TABLE `fivenet_documents_meta` ADD COLUMN `comment_count` int(11) NOT NULL DEFAULT '0' AFTER `ap_policies_active`;

-- Table: `fivenet_jobs` - Add `created_at` column and `deleted_at` column for soft deletes
ALTER TABLE `fivenet_jobs` ADD COLUMN `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) AFTER `label`;
ALTER TABLE `fivenet_jobs` ADD COLUMN `deleted_at` datetime(3) DEFAULT NULL AFTER `created_at`;

-- Table: `fivenet_calendar` - Fix job inserted as empty string and not NULL in some cases.
UPDATE `fivenet_calendar` SET `job` = NULL WHERE `job` = '';

COMMIT;
