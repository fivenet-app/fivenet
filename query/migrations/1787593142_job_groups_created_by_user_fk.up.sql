BEGIN;

ALTER TABLE `fivenet_job_groups`
  MODIFY COLUMN `created_by_user_id` int(11) DEFAULT NULL AFTER `exclusions_count`,
  ADD KEY `idx_fivenet_job_groups_created_by_user_id` (`created_by_user_id`),
  ADD KEY `idx_fivenet_job_groups_updated_by_user_id` (`updated_by_user_id`),
  ADD CONSTRAINT `fk_fivenet_job_groups_created_by_user_id` FOREIGN KEY (`created_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_fivenet_job_groups_updated_by_user_id` FOREIGN KEY (`updated_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_group_leaders`
  MODIFY COLUMN `created_by_user_id` int(11) DEFAULT NULL AFTER `user_id`,
  ADD CONSTRAINT `fk_fivenet_job_group_leaders_created_by_user_id` FOREIGN KEY (`created_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_group_manual_members`
  MODIFY COLUMN `created_by_user_id` int(11) DEFAULT NULL AFTER `reason`,
  ADD CONSTRAINT `fk_fivenet_job_group_manual_members_created_by_user_id` FOREIGN KEY (`created_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_group_member_exclusions`
  MODIFY COLUMN `created_by_user_id` int(11) DEFAULT NULL AFTER `reason`,
  ADD CONSTRAINT `fk_fivenet_job_group_member_exclusions_created_by_user_id` FOREIGN KEY (`created_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `fivenet_job_group_rules`
  MODIFY COLUMN `created_by_user_id` int(11) DEFAULT NULL AFTER `enabled`,
  ADD KEY `idx_fivenet_job_group_rules_created_by_user_id` (`created_by_user_id`),
  ADD KEY `idx_fivenet_job_group_rules_updated_by_user_id` (`updated_by_user_id`),
  ADD CONSTRAINT `fk_fivenet_job_group_rules_created_by_user_id` FOREIGN KEY (`created_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_fivenet_job_group_rules_updated_by_user_id` FOREIGN KEY (`updated_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE;

COMMIT;
