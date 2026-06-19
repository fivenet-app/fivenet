BEGIN;

-- Table: fivenet_job_groups
CREATE TABLE IF NOT EXISTS `fivenet_job_groups` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `job` varchar(20) NOT NULL,
  `name` varchar(128) NOT NULL,
  `description` longtext DEFAULT NULL,
  `short_name` varchar(32) DEFAULT NULL,
  `logo_file_id` bigint(20) unsigned DEFAULT NULL,
  `color` varchar(32) DEFAULT NULL,
  `type` tinyint(2) NOT NULL DEFAULT 1,
  `state` tinyint(2) NOT NULL DEFAULT 1,
  `membership_mode` tinyint(2) NOT NULL DEFAULT 1,
  `sort_rank` varchar(32) NOT NULL DEFAULT '',
  `members_count` int(11) NOT NULL DEFAULT 0,
  `leaders_count` int(11) NOT NULL DEFAULT 0,
  `rules_count` int(11) NOT NULL DEFAULT 0,
  `exclusions_count` int(11) NOT NULL DEFAULT 0,
  `created_by_user_id` int(11) NOT NULL,
  `updated_by_user_id` int(11) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_job_groups_job_name` (`job`, `name`),
  KEY `idx_fivenet_job_groups_job_state_sort` (`job`, `state`, `sort_rank`, `name`),
  KEY `idx_fivenet_job_groups_deleted` (`job`, `deleted_at`),
  KEY `idx_fivenet_job_groups_job` (`job`),
  CONSTRAINT `fk_fivenet_job_groups_logo_file_id` FOREIGN KEY (`logo_file_id`) REFERENCES `fivenet_files` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_leaders
CREATE TABLE IF NOT EXISTS `fivenet_job_group_leaders` (
  `group_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `created_by_user_id` int(11) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`group_id`, `user_id`),
  KEY `idx_fivenet_job_group_leaders_user_id` (`user_id`, `group_id`),
  KEY `idx_fivenet_job_group_leaders_created_by_user_id` (`created_by_user_id`),
  CONSTRAINT `fk_fivenet_job_group_leaders_group_id` FOREIGN KEY (`group_id`) REFERENCES `fivenet_job_groups` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_group_leaders_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_manual_members
CREATE TABLE IF NOT EXISTS `fivenet_job_group_manual_members` (
  `group_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `created_by_user_id` int(11) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`group_id`, `user_id`),
  KEY `idx_fivenet_job_group_manual_members_user_id` (`user_id`, `group_id`),
  KEY `idx_fivenet_job_group_manual_members_created_by_user_id` (`created_by_user_id`),
  CONSTRAINT `fk_fivenet_job_group_manual_members_group_id` FOREIGN KEY (`group_id`) REFERENCES `fivenet_job_groups` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_group_manual_members_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_member_exclusions
CREATE TABLE IF NOT EXISTS `fivenet_job_group_member_exclusions` (
  `group_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `reason_type` tinyint(2) NOT NULL DEFAULT 1,
  `reason` varchar(255) DEFAULT NULL,
  `created_by_user_id` int(11) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`group_id`, `user_id`),
  KEY `idx_fivenet_job_group_member_exclusions_user_id` (`user_id`, `group_id`),
  KEY `idx_fivenet_job_group_member_exclusions_created_by_user_id` (`created_by_user_id`),
  KEY `idx_fivenet_job_group_member_exclusions_reason_type` (`reason_type`),
  CONSTRAINT `fk_fivenet_job_group_member_exclusions_group_id` FOREIGN KEY (`group_id`) REFERENCES `fivenet_job_groups` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_group_member_exclusions_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_rules
CREATE TABLE IF NOT EXISTS `fivenet_job_group_rules` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) unsigned NOT NULL,
  `rule_type` tinyint(2) NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT 1,
  `created_by_user_id` int(11) NOT NULL,
  `updated_by_user_id` int(11) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_job_group_rules_group_id_enabled` (`group_id`, `enabled`),
  KEY `idx_fivenet_job_group_rules_rule_type` (`rule_type`),
  CONSTRAINT `fk_fivenet_job_group_rules_group_id` FOREIGN KEY (`group_id`) REFERENCES `fivenet_job_groups` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_rule_grades
CREATE TABLE IF NOT EXISTS `fivenet_job_group_rule_grades` (
  `rule_id` bigint(20) unsigned NOT NULL,
  `grade_rule_type` tinyint(2) NOT NULL,
  `grade` int(11) DEFAULT NULL,
  `min_grade` int(11) DEFAULT NULL,
  `max_grade` int(11) DEFAULT NULL,
  PRIMARY KEY (`rule_id`),
  CONSTRAINT `fk_fivenet_job_group_rule_grades_rule_id` FOREIGN KEY (`rule_id`) REFERENCES `fivenet_job_group_rules` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_rule_qualifications
CREATE TABLE IF NOT EXISTS `fivenet_job_group_rule_qualifications` (
  `rule_id` bigint(20) unsigned NOT NULL,
  `qualification_rule_type` tinyint(2) NOT NULL,
  `require_completed` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`rule_id`),
  CONSTRAINT `fk_fivenet_job_group_rule_qualifications_rule_id` FOREIGN KEY (`rule_id`) REFERENCES `fivenet_job_group_rules` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_rule_qualification_items
CREATE TABLE IF NOT EXISTS `fivenet_job_group_rule_qualification_items` (
  `rule_id` bigint(20) unsigned NOT NULL,
  `qualification_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`rule_id`, `qualification_id`),
  KEY `idx_fivenet_job_group_rule_qualification_items_qualification_id` (`qualification_id`, `rule_id`),
  CONSTRAINT `fk_fivenet_job_group_rule_qualification_items_rule_id` FOREIGN KEY (`rule_id`) REFERENCES `fivenet_job_group_rules` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Table: fivenet_job_group_activity
CREATE TABLE IF NOT EXISTS `fivenet_job_group_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `job` varchar(20) NOT NULL,
  `group_id` bigint(20) unsigned NOT NULL,
  `activity_type` tinyint(2) NOT NULL,
  `actor_user_id` int(11) NOT NULL,
  `target_user_id` int(11) DEFAULT NULL,
  `rule_id` bigint(20) unsigned DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `data` longtext DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_job_group_activity_group_created` (`group_id`, `created_at`),
  KEY `idx_fivenet_job_group_activity_job_group_created` (`job`, `group_id`, `created_at`),
  KEY `idx_fivenet_job_group_activity_type` (`activity_type`),
  KEY `idx_fivenet_job_group_activity_actor_user_id` (`actor_user_id`),
  KEY `idx_fivenet_job_group_activity_target_user_id` (`target_user_id`),
  KEY `idx_fivenet_job_group_activity_rule_id` (`rule_id`),
  CONSTRAINT `fk_fivenet_job_group_activity_group_id` FOREIGN KEY (`group_id`) REFERENCES `fivenet_job_groups` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_group_activity_actor_user_id` FOREIGN KEY (`actor_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_group_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
