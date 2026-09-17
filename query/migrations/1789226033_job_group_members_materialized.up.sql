BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_job_group_members` (
  `group_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `is_leader` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`group_id`, `user_id`),
  KEY `idx_fivenet_job_group_members_user_id` (`user_id`, `group_id`),
  CONSTRAINT `fk_fivenet_job_group_members_group_id` FOREIGN KEY (`group_id`) REFERENCES `fivenet_job_groups` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_group_members_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
