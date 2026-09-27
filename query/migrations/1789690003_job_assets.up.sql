BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_job_assets` (
  `job` varchar(20) NOT NULL,
  `file_id` bigint unsigned NOT NULL,
  `display_name` varchar(255) NOT NULL DEFAULT '',
  `created_by_user_id` int(11) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`job`, `file_id`),
  KEY `idx_fivenet_job_assets_file_id` (`file_id`),
  KEY `idx_fivenet_job_assets_created_by_user_id` (`created_by_user_id`),
  KEY `idx_fivenet_job_assets_job_created` (`job`, `created_at`),
  CONSTRAINT `fk_fivenet_job_assets_file_id` FOREIGN KEY (`file_id`) REFERENCES `fivenet_files` (`id`)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_job_assets_created_by_user_id` FOREIGN KEY (`created_by_user_id`) REFERENCES `fivenet_user` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
