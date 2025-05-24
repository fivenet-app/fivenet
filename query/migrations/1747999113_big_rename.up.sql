BEGIN;

ALTER TABLE `fivenet_centrum_users` DROP FOREIGN KEY `fk_fivenet_centrum_users_identifier`;
ALTER TABLE `fivenet_centrum_users` DROP COLUMN `identifier`;

-- Rename tables
RENAME TABLE `fivenet_attrs` TO `fivenet_rbac_attrs`;
RENAME TABLE `fivenet_job_attrs` TO `fivenet_rbac_job_attrs`;
RENAME TABLE `fivenet_job_permissions` TO `fivenet_rbac_job_permissions`;
RENAME TABLE `fivenet_permissions` TO `fivenet_rbac_permissions`;
RENAME TABLE `fivenet_role_attrs` TO `fivenet_rbac_roles_attrs`;
RENAME TABLE `fivenet_role_permissions` TO `fivenet_rbac_roles_permissions`;
RENAME TABLE `fivenet_roles` TO `fivenet_rbac_roles`;
RENAME TABLE `fivenet_oauth2_accounts` TO `fivenet_accounts_oauth2`;
RENAME TABLE `fivenet_jobs_conduct` TO `fivenet_job_conduct`;
RENAME TABLE `fivenet_jobs_labels` TO `fivenet_job_labels`;
RENAME TABLE `fivenet_jobs_labels_users` TO `fivenet_job_colleague_labels`;
RENAME TABLE `fivenet_jobs_timeclock` TO `fivenet_job_timeclock`;
RENAME TABLE `fivenet_jobs_user_activity` TO `fivenet_job_colleague_activity`;
RENAME TABLE `fivenet_jobs_user_props` TO `fivenet_job_colleague_props`;
RENAME TABLE `fivenet_users` TO `fivenet_user`;
RENAME TABLE `fivenet_user_citizen_labels` TO `fivenet_user_labels`;
RENAME TABLE `fivenet_job_citizen_labels` TO `fivenet_user_labels_job`;
RENAME TABLE `fivenet_user_locations` TO `fivenet_centrum_user_locations`;
RENAME TABLE `fivenet_centrum_users` TO `fivenet_centrum_disponents`;
RENAME TABLE `fivenet_wiki_page_activity` TO `fivenet_wiki_pages_activity`;
RENAME TABLE `fivenet_job_grades` TO `fivenet_jobs_grades`;

-- Rename columns
ALTER TABLE `fivenet_user_labels` CHANGE `attribute_id` `label_id` bigint unsigned NOT NULL;

-- Rename indexes
ALTER TABLE `fivenet_rbac_attrs` RENAME INDEX `idx_fivenet_attrs_permission_id_key_unque` TO `idx_permission_id_key_unique`;
ALTER TABLE `fivenet_centrum_disponents` RENAME INDEX `idx_fivenet_centrum_users_job` TO `idx_job`;
ALTER TABLE `fivenet_rbac_job_attrs` RENAME INDEX `fk_fivenet_job_attrs_attr_id` TO `idx_attr_id`;
ALTER TABLE `fivenet_rbac_job_permissions` RENAME INDEX `fk_fivenet_job_permissions_permission` TO `idx_permission`;
ALTER TABLE `fivenet_job_conduct` RENAME INDEX `fivenet_jobs_conduct_type` TO `idx_conduct_type`;
ALTER TABLE `fivenet_job_conduct` RENAME INDEX `fivenet_jobs_conduct_created_at` TO `idx_conduct_created_at`;
ALTER TABLE `fivenet_job_conduct` RENAME INDEX `fivenet_jobs_conduct_target_user_id` TO `fk_fivenet_jobs_conduct_target_user_id`;
ALTER TABLE `fivenet_job_conduct` RENAME INDEX `fk_fivenet_jobs_conduct_creator_id` TO `fk_fivenet_jobs_conduct_creator_id`;
ALTER TABLE `fivenet_job_conduct` RENAME INDEX `idx_fivenet_jobs_conduct_deleted_at` TO `idx_deleted_at`;
ALTER TABLE `fivenet_user_labels` RENAME INDEX `idx_fivenet_user_citizen_attributes_unique` TO `idx_unique`;
ALTER TABLE `fivenet_user_labels` RENAME INDEX `fk_fivenet_user_citizen_attributes_attribute_id` TO `fk_fivenet_user_labels_label_id`;
ALTER TABLE `fivenet_user_labels_job` RENAME INDEX `idx_fivenet_job_citizen_attributes_unique` TO `idx_unique`;
ALTER TABLE `fivenet_user_labels_job` RENAME INDEX `idx_fivenet_job_citizen_attributes_name` TO `idx_name`;
ALTER TABLE `fivenet_user_labels_job` RENAME INDEX `idx_fivenet_job_citizen_labels_sort_key` TO `idx_sort_key`;
ALTER TABLE `fivenet_job_labels` RENAME INDEX `idx_fivenet_jobs_labels_unique` TO `idx_unique`;
ALTER TABLE `fivenet_job_labels` RENAME INDEX `idx_fivenet_jobs_labels_name` TO `idx_name`;
ALTER TABLE `fivenet_job_labels` RENAME INDEX `idx_fivenet_jobs_labels_order` TO `idx_order`;
ALTER TABLE `fivenet_job_labels` RENAME INDEX `idx_fivenet_jobs_labels_sort_key` TO `idx_sort_key`;
ALTER TABLE `fivenet_job_labels` RENAME INDEX `idx_fivenet_jobs_labels_deleted_at` TO `idx_deleted_at`;
ALTER TABLE `fivenet_job_colleague_labels` RENAME INDEX `idx_fivenet_jobs_labels_users_unique` TO `idx_unique`;
ALTER TABLE `fivenet_job_colleague_labels` RENAME INDEX `fk_fivenet_jobs_labels_users_label_id` TO `fk_fivenet_job_colleague_labels_label_id`;
ALTER TABLE `fivenet_job_timeclock` RENAME INDEX `idx_fivenet_jobs_timeclock_unique` TO `idx_unique`;
ALTER TABLE `fivenet_job_timeclock` RENAME INDEX `fk_fivenet_jobs_timeclock_user_id` TO `fk_fivenet_job_timeclock_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` RENAME INDEX `idx_fivenet_jobs_user_activity_job` TO `idx_job`;
ALTER TABLE `fivenet_job_colleague_activity` RENAME INDEX `idx_fivenet_jobs_user_activity_source_user_id` TO `fk_fivenet_jobs_user_activity_source_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` RENAME INDEX `idx_fivenet_jobs_user_activity_target_user_id` TO `fk_fivenet_jobs_user_activity_target_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` RENAME INDEX `idx_fivenet_jobs_user_activity_activity_type` TO `idx_activity_type`;
ALTER TABLE `fivenet_job_colleague_props` RENAME INDEX `idx_fivenet_jobs_user_props_unique` TO `idx_unique`;
ALTER TABLE `fivenet_job_colleague_props` RENAME INDEX `idx_fivenet_jobs_user_props_job_name_prefix` TO `idx_name_prefix`;
ALTER TABLE `fivenet_job_colleague_props` RENAME INDEX `idx_fivenet_jobs_user_props_job_name_suffix` TO `idx_name_suffix`;
ALTER TABLE `fivenet_job_colleague_props` RENAME INDEX `idx_fivenet_jobs_user_props_deleted_at` TO `idx_deleted_at`;
ALTER TABLE `fivenet_accounts_oauth2` RENAME INDEX `idx_fivenet_oauth2_accounts_unique` TO `idx_unique`;
ALTER TABLE `fivenet_accounts_oauth2` RENAME INDEX `idx_fivenet_oauth2_accounts_provider_external_id` TO `idx_provider_external_id`;
ALTER TABLE `fivenet_rbac_permissions` RENAME INDEX `idx_fivenet_permissions_category_name_unique` TO `idx_category_name_unique`;
ALTER TABLE `fivenet_rbac_permissions` RENAME INDEX `idx_fivenet_permissions_guard_name_unique` TO `idx_guard_name_unique`;
ALTER TABLE `fivenet_rbac_permissions` RENAME INDEX `idx_fivenet_permissions_category` TO `idx_category`;
ALTER TABLE `fivenet_rbac_permissions` RENAME INDEX `idx_fivenet_permissions_order` TO `idx_order`;
ALTER TABLE `fivenet_rbac_roles_attrs` RENAME INDEX `fk_fivenet_role_attrs_attr_id` TO `idx_attr_id`;
ALTER TABLE `fivenet_rbac_roles_permissions` RENAME INDEX `fk_fivenet_role_permissions_permission` TO `idx_permission`;
ALTER TABLE `fivenet_rbac_roles` RENAME INDEX `idx_fivenet_roles_job_grade_unique` TO `idx_job_grade_unique`;
ALTER TABLE `fivenet_centrum_user_locations` RENAME INDEX `idx_fivenet_user_locations_job` TO `idx_job`;
ALTER TABLE `fivenet_user_props` RENAME INDEX `idx_fivenet_user_props_unique` TO `idx_unique`;
ALTER TABLE `fivenet_user_props` RENAME INDEX `idx_fivenet_user_props_wanted` TO `idx_wanted`;
ALTER TABLE `fivenet_user_props` RENAME INDEX `idx_fivenet_user_props_avatar` TO `idx_avatar`;
ALTER TABLE `fivenet_user_props` RENAME INDEX `idx_fivenet_user_props_mug_shot` TO `idx_mug_shot`;
ALTER TABLE `fivenet_user` RENAME INDEX `id` TO `idx_id`;
ALTER TABLE `fivenet_user` RENAME INDEX `idx_fivenet_users_job` TO `idx_job`;
ALTER TABLE `fivenet_user` RENAME INDEX `idx_fivenet_users_dateofbirth` TO `idx_dateofbirth`;
ALTER TABLE `fivenet_user` RENAME INDEX `idx_fivenet_users_firstname_lastname_fulltext` TO `idx_firstname_lastname_fulltext`;

-- "Rename" foreign keys
SET FOREIGN_KEY_CHECKS=0;

ALTER TABLE `fivenet_accounts_oauth2` DROP FOREIGN KEY `fk_fivenet_oauth2_accounts_account_id`,
  ADD CONSTRAINT `fk_fivenet_accounts_oauth2_account_id` FOREIGN KEY (`account_id`) REFERENCES `fivenet_accounts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_disponents` DROP FOREIGN KEY `fk_fivenet_centrum_users_user_id`,
  ADD CONSTRAINT `fk_fivenet_centrum_disponents_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_centrum_user_locations` DROP FOREIGN KEY `fk_fivenet_user_locations_identifier`,
  ADD CONSTRAINT `fk_fivenet_centrum_user_locations_identifier` FOREIGN KEY (`identifier`) REFERENCES `users` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_labels` DROP FOREIGN KEY `fk_fivenet_user_citizen_attributes_attribute_id`,
  ADD CONSTRAINT `fk_fivenet_user_labels_label_id` FOREIGN KEY (`label_id`) REFERENCES `fivenet_user_labels_job` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_user_labels` DROP FOREIGN KEY `fk_fivenet_user_citizen_labels_user_id`,
  ADD CONSTRAINT `fk_fivenet_user_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_documents` DROP FOREIGN KEY `fk_fivenet_documents_categories`,
  ADD CONSTRAINT `fk_fivenet_documents_category_id` FOREIGN KEY (`category_id`) REFERENCES `fivenet_documents_categories` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_documents_templates` DROP FOREIGN KEY `fk_fivenet_documents_templates_categories`,
  ADD CONSTRAINT `fk_fivenet_documents_templates_category_id` FOREIGN KEY (`category_id`) REFERENCES `fivenet_documents_categories` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_job_colleague_activity` DROP FOREIGN KEY `fk_fivenet_jobs_user_activity_source_user_id`,
  ADD CONSTRAINT `fk_fivenet_job_colleague_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_job_colleague_activity` DROP FOREIGN KEY `fk_fivenet_jobs_user_activity_target_user_id`,
  ADD CONSTRAINT `fk_fivenet_job_colleague_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_colleague_labels` DROP FOREIGN KEY `fk_fivenet_jobs_labels_users_label_id`,
  ADD CONSTRAINT `fk_fivenet_job_colleague_labels_label_id` FOREIGN KEY (`label_id`) REFERENCES `fivenet_job_labels` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_colleague_labels` DROP FOREIGN KEY `fk_fivenet_jobs_labels_users_user_id`,
  ADD CONSTRAINT `fk_fivenet_job_colleague_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_colleague_props` DROP FOREIGN KEY `fk_fivenet_jobs_user_props_user_id`,
  ADD CONSTRAINT `fk_fivenet_job_colleague_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_conduct` DROP FOREIGN KEY `fk_fivenet_jobs_conduct_creator_id`,
  ADD CONSTRAINT `fk_fivenet_job_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_conduct` DROP FOREIGN KEY `fk_fivenet_jobs_conduct_target_user_id`,
  ADD CONSTRAINT `fk_fivenet_job_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_job_timeclock` DROP FOREIGN KEY `fk_fivenet_jobs_timeclock_user_id`,
  ADD CONSTRAINT `fk_fivenet_job_timeclock_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_threads_recipients` DROP FOREIGN KEY `fk_fivenet_mailer_threads_recipients_emails_email_id`,
  ADD CONSTRAINT `fk_fivenet_mailer_threads_recipients_email_id` FOREIGN KEY (`email_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_mailer_threads_recipients` DROP FOREIGN KEY `fk_fivenet_mailer_threads_recipients_emails_thread_id`,
  ADD CONSTRAINT `fk_fivenet_mailer_threads_recipients_thread_id` FOREIGN KEY (`thread_id`) REFERENCES `fivenet_mailer_threads` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_questions` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_questions_quali_id`,
  ADD CONSTRAINT `fk_fivenet_qualifications_exam_questions_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_responses` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_responses_quali_user_id`,
  ADD CONSTRAINT `fk_fivenet_qualifications_exam_responses_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_exam_users` DROP FOREIGN KEY `fk_fivenet_qualifications_exam_users_quali_id`,
  ADD CONSTRAINT `fk_fivenet_qualifications_exam_users_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_requirements` DROP FOREIGN KEY `fk_fivenet_qualifications_requirements_quali_id`,
  ADD CONSTRAINT `fk_fivenet_qualifications_requirements_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_qualifications_requirements` DROP FOREIGN KEY `fk_fivenet_qualifications_requirements_target_quali_id`,
  ADD CONSTRAINT `fk_fivenet_qualifications_requirements_target_qualification_id` FOREIGN KEY (`target_qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_attrs` DROP FOREIGN KEY `fk_fivenet_attrs_permissions_permission_id`,
  ADD CONSTRAINT `fk_fivenet_rbac_attrs_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `fivenet_rbac_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_job_attrs` DROP FOREIGN KEY `fk_fivenet_job_attrs_attr_id`,
  ADD CONSTRAINT `fk_fivenet_rbac_job_attrs_attr_id` FOREIGN KEY (`attr_id`) REFERENCES `fivenet_rbac_attrs` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_job_permissions` DROP FOREIGN KEY `fk_fivenet_job_permissions_permission`,
  ADD CONSTRAINT `fk_fivenet_rbac_job_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `fivenet_rbac_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_roles_attrs` DROP FOREIGN KEY `fk_fivenet_role_attrs_attr_id`,
  ADD CONSTRAINT `fk_fivenet_rbac_roles_attrs_attr_id` FOREIGN KEY (`attr_id`) REFERENCES `fivenet_rbac_attrs` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_roles_attrs` DROP FOREIGN KEY `fk_fivenet_role_attrs_role_id`,
  ADD CONSTRAINT `fk_fivenet_rbac_roles_attrs_role_id` FOREIGN KEY (`role_id`) REFERENCES `fivenet_rbac_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_roles_permissions` DROP FOREIGN KEY `fk_fivenet_role_permissions_permission`,
  ADD CONSTRAINT `fk_fivenet_rbac_roles_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `fivenet_rbac_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_rbac_roles_permissions` DROP FOREIGN KEY `fk_fivenet_role_permissions_role`,
  ADD CONSTRAINT `fk_fivenet_rbac_roles_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `fivenet_rbac_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_wiki_pages_activity` DROP FOREIGN KEY `fk_fivenet_wiki_page_activity_creator_id`,
  ADD CONSTRAINT `fk_fivenet_wiki_pages_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE `fivenet_wiki_pages_activity` DROP FOREIGN KEY `fk_fivenet_wiki_page_activity_page_id`,
  ADD CONSTRAINT `fk_fivenet_wiki_pages_activity_page_id` FOREIGN KEY (`page_id`) REFERENCES `fivenet_wiki_pages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `fivenet_jobs_grades` DROP FOREIGN KEY `fk_fivenet_job_grades_job_name`,
  ADD CONSTRAINT `fk_fivenet_jobs_grades_job_name` FOREIGN KEY (`job_name`) REFERENCES `fivenet_jobs` (`name`) ON DELETE CASCADE ON UPDATE CASCADE;

SET FOREIGN_KEY_CHECKS=1;

-- Rename permission categories + guard
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('auth-', REPLACE(guard_name, 'authservice', 'authservice')), `category` = 'auth.AuthService' WHERE `category` = 'AuthService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('calendar-', REPLACE(guard_name, 'calendarservice', 'calendarservice')), `category` = 'calendar.CalendarService' WHERE `category` = 'CalendarService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('centrum-', REPLACE(guard_name, 'centrumservice', 'centrumservice')), `category` = 'centrum.CentrumService' WHERE `category` = 'CentrumService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('citizens-', REPLACE(guard_name, 'citizenstoreservice', 'citizensservice')), `category` = 'citizens.CitizensService' WHERE `category` = 'CitizenStoreService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('completor-', REPLACE(guard_name, 'completorservice', 'completorservice')), `category` = 'completor.CompletorService' WHERE `category` = 'CompletorService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('vehicles-', REPLACE(guard_name, 'dmvservice', 'vehiclesservice')), `category` = 'vehicles.VehiclesService' WHERE `category` = 'DMVService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('documents-', REPLACE(guard_name, 'docstoreservice', 'documentsservice')), `category` = 'documents.DocumentsService' WHERE `category` = 'DocStoreService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('jobs-', REPLACE(guard_name, 'jobsconductservice', 'conductservice')), `category` = 'jobs.ConductService' WHERE `category` = 'JobsConductService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('jobs-', REPLACE(guard_name, 'jobsservice', 'jobsservice')), `category` = 'jobs.JobsService' WHERE `category` = 'JobsService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('jobs-', REPLACE(guard_name, 'jobstimeclockservice', 'timeclockservice')), `category` = 'jobs.TimeclockService' WHERE `category` = 'JobsTimeclockService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('livemap-', REPLACE(guard_name, 'livemapperservice', 'livemapservice')), `category` = 'livemap.LivemapService' WHERE `category` = 'LivemapperService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('mailer-', REPLACE(guard_name, 'mailerservice', 'mailerservice')), `category` = 'mailer.MailerService' WHERE `category` = 'MailerService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('qualifications-', REPLACE(guard_name, 'qualificationsservice', 'qualificationsservice')), `category` = 'qualifications.QualificationsService' WHERE `category` = 'QualificationsService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('settings-', REPLACE(guard_name, 'rectorlawsservice', 'lawsservice')), `category` = 'settings.LawsService' WHERE `category` = 'RectorLawsService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('settings-', REPLACE(guard_name, 'rectorservice', 'settingsservice')), `category` = 'settings.SettingsService' WHERE `category` = 'RectorService';
UPDATE `fivenet_rbac_permissions` SET `guard_name` = CONCAT('wiki-', REPLACE(guard_name, 'wikiservice', 'wikiservice')), `category` = 'wiki.WikiService' WHERE `category` = 'WikiService';

-- Rename permissions
UPDATE `fivenet_rbac_permissions`
	SET name='ManageLabels', guard_name='citizens-citizensservice-managelabels'
	WHERE category = 'citizens.CitizensService' AND name = 'ManageColleagueLabels';

UPDATE `fivenet_rbac_permissions`
	SET name='ManageLabels', guard_name='jobs-jobsservice-managelabels'
	WHERE category = 'jobs.JobsService' AND name = 'ManageColleagueLabels';

UPDATE `fivenet_rbac_permissions`
	SET name='SetColleagueProps', guard_name='jobs-jobsservice-setcolleagueprops'
	WHERE category = 'jobs.JobsService' AND name = 'SetJobsUserProps';

COMMIT;
