BEGIN;

-- Drop reason field from fivenet_user_labels
ALTER TABLE `fivenet_user_labels` DROP COLUMN `reason`;

-- Table: fivenet_user_labels_job_job_access
CREATE TABLE IF NOT EXISTS `fivenet_user_labels_job_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `label_id` bigint(20) unsigned NOT NULL,
  `job` varchar(40) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_user_labels_job_job_access` (`label_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_user_labels_job_job_access_label_id` (`label_id`),
  CONSTRAINT `fk_fivenet_user_labels_job_job_access_label_id` FOREIGN KEY (`label_id`) REFERENCES `fivenet_user_labels_job` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
