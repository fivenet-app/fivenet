BEGIN;

UPDATE `fivenet_job_group_rules` SET `created_by_user_id` = 0 WHERE `created_by_user_id` IS NULL;

ALTER TABLE `fivenet_job_group_rules`
  DROP FOREIGN KEY `fk_fivenet_job_group_rules_created_by_user_id`,
  DROP INDEX `idx_fivenet_job_group_rules_created_by_user_id`,
  DROP FOREIGN KEY `fk_fivenet_job_group_rules_updated_by_user_id`,
  DROP INDEX `idx_fivenet_job_group_rules_updated_by_user_id`,
  MODIFY COLUMN `created_by_user_id` int(11) NOT NULL AFTER `enabled`;

UPDATE `fivenet_job_group_member_exclusions` SET `created_by_user_id` = 0 WHERE `created_by_user_id` IS NULL;

ALTER TABLE `fivenet_job_group_member_exclusions`
  DROP FOREIGN KEY `fk_fivenet_job_group_member_exclusions_created_by_user_id`,
  MODIFY COLUMN `created_by_user_id` int(11) NOT NULL AFTER `reason`;

UPDATE `fivenet_job_group_manual_members` SET `created_by_user_id` = 0 WHERE `created_by_user_id` IS NULL;

ALTER TABLE `fivenet_job_group_manual_members`
  DROP FOREIGN KEY `fk_fivenet_job_group_manual_members_created_by_user_id`,
  MODIFY COLUMN `created_by_user_id` int(11) NOT NULL AFTER `reason`;

UPDATE `fivenet_job_group_leaders` SET `created_by_user_id` = 0 WHERE `created_by_user_id` IS NULL;

ALTER TABLE `fivenet_job_group_leaders`
  DROP FOREIGN KEY `fk_fivenet_job_group_leaders_created_by_user_id`,
  MODIFY COLUMN `created_by_user_id` int(11) NOT NULL AFTER `user_id`;

UPDATE `fivenet_job_groups` SET `created_by_user_id` = 0 WHERE `created_by_user_id` IS NULL;

ALTER TABLE `fivenet_job_groups`
  DROP FOREIGN KEY `fk_fivenet_job_groups_created_by_user_id`,
  DROP INDEX `idx_fivenet_job_groups_created_by_user_id`,
  DROP FOREIGN KEY `fk_fivenet_job_groups_updated_by_user_id`,
  DROP INDEX `idx_fivenet_job_groups_updated_by_user_id`,
  MODIFY COLUMN `created_by_user_id` int(11) NOT NULL AFTER `exclusions_count`;

COMMIT;
