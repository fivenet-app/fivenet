BEGIN;

-- Table: fivenet_job_groups
CREATE TABLE IF NOT EXISTS `fivenet_job_groups` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` bigint(20) unsigned NOT NULL,
  `name` varchar(128) NOT NULL,
  `description` longtext DEFAULT NULL,
  `short_name` varchar(32) DEFAULT NULL,
  `logo_file_id` varchar(128) DEFAULT NULL,
  `color` varchar(32) DEFAULT NULL,
  `type` tinyint(2) NOT NULL DEFAULT 1,
  `state` tinyint(2) NOT NULL DEFAULT 1,
  `membership_mode` tinyint(2) NOT NULL DEFAULT 1,
  `sort_order` int(11) NOT NULL DEFAULT 0,
  `members_count` int(11) NOT NULL DEFAULT 0,
  `leaders_count` int(11) NOT NULL DEFAULT 0,
  `rules_count` int(11) NOT NULL DEFAULT 0,
  `exclusions_count` int(11) NOT NULL DEFAULT 0,
  `created_by_user_id` int(11) NOT NULL,
  `updated_by_user_id` int(11) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  `archived_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_job_groups_job_name` (`job_id`, `name`),
  KEY `idx_fivenet_job_groups_job_state_sort` (`job_id`, `state`, `sort_order`, `name`),
  KEY `idx_fivenet_job_groups_archived` (`job_id`, `archived_at`),
  KEY `idx_fivenet_job_groups_job_id` (`job_id`)
) ENGINE=InnoDB;

COMMIT;
