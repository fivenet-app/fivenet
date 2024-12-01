BEGIN;

CREATE TABLE `fivenet_jobs_labels` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job` varchar(20) NOT NULL,
  `name` varchar(32) NOT NULL,
  `color` char(7) DEFAULT '#5c7aff',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_job_citizen_attributes_unique` (`job`, `name`),
  KEY `idx_fivenet_job_citizen_attributes_name` (`name`)
) ENGINE=InnoDB;

CREATE TABLE `fivenet_jobs_labels_users` (
  `user_id` int NOT NULL,
  `job` varchar(20) NOT NULL,
  `attribute_id` bigint unsigned NOT NULL,
  UNIQUE KEY `idx_fivenet_jobs_labels_users_unique` (`user_id`, `job`, `attribute_id`),
  KEY `fk_fivenet_jobs_labels_users_attribute_id` (`attribute_id`),
  CONSTRAINT `fk_fivenet_jobs_labels_users_attribute_id` FOREIGN KEY (`attribute_id`) REFERENCES `fivenet_job_citizen_attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_jobs_labels_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
