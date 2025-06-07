BEGIN;

ALTER TABLE `fivenet_centrum_disponents` DROP FOREIGN KEY `fk_fivenet_centrum_disponents_user_id`;
ALTER TABLE `fivenet_centrum_disponents` ADD CONSTRAINT `fk_fivenet_centrum_disponents_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_centrum_user_locations` DROP FOREIGN KEY `fk_fivenet_centrum_user_locations_identifier`;
ALTER TABLE `fivenet_centrum_user_locations` ADD CONSTRAINT `fk_fivenet_centrum_user_locations_identifier` FOREIGN KEY (`identifier`) REFERENCES `{{.UsersTableName}}` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_user_labels` DROP FOREIGN KEY `fk_fivenet_user_labels_user_id`;
ALTER TABLE `fivenet_user_labels` ADD CONSTRAINT `fk_fivenet_user_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_colleague_activity` DROP FOREIGN KEY `fk_fivenet_job_colleague_activity_source_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` ADD CONSTRAINT `fk_fivenet_job_colleague_activity_source_user_id` FOREIGN KEY (`source_user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

ALTER TABLE `fivenet_job_colleague_activity` DROP FOREIGN KEY `fk_fivenet_job_colleague_activity_target_user_id`;
ALTER TABLE `fivenet_job_colleague_activity` ADD CONSTRAINT `fk_fivenet_job_colleague_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_colleague_labels` DROP FOREIGN KEY `fk_fivenet_job_colleague_labels_user_id`;
ALTER TABLE `fivenet_job_colleague_labels` ADD CONSTRAINT `fk_fivenet_job_colleague_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_colleague_props` DROP FOREIGN KEY `fk_fivenet_job_colleague_props_user_id`;
ALTER TABLE `fivenet_job_colleague_props` ADD CONSTRAINT `fk_fivenet_job_colleague_props_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_conduct` DROP FOREIGN KEY `fk_fivenet_job_conduct_creator_id`;
ALTER TABLE `fivenet_job_conduct` ADD CONSTRAINT `fk_fivenet_job_conduct_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_conduct` DROP FOREIGN KEY `fk_fivenet_job_conduct_target_user_id`;
ALTER TABLE `fivenet_job_conduct` ADD CONSTRAINT `fk_fivenet_job_conduct_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_timeclock` DROP FOREIGN KEY `fk_fivenet_job_timeclock_user_id`;
ALTER TABLE `fivenet_job_timeclock` ADD CONSTRAINT `fk_fivenet_job_timeclock_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `fivenet_wiki_pages_activity` DROP FOREIGN KEY `fk_fivenet_wiki_pages_activity_creator_id`;
ALTER TABLE `fivenet_wiki_pages_activity` ADD CONSTRAINT `fk_fivenet_wiki_pages_activity_creator_id` FOREIGN KEY (`creator_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

COMMIT;
