BEGIN;

CREATE TABLE `fivenet_jobs_labels` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job` varchar(20) NOT NULL,
  `name` varchar(32) NOT NULL,
  `color` char(7) DEFAULT '#5c7aff',
  `order` mediumint(4) DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_jobs_labels_unique` (`job`, `name`),
  KEY `idx_fivenet_jobs_labels_name` (`name`),
  KEY `idx_fivenet_jobs_labels_order` (`order`)
) ENGINE=InnoDB;

CREATE TABLE `fivenet_jobs_labels_users` (
  `user_id` int NOT NULL,
  `job` varchar(20) NOT NULL,
  `label_id` bigint unsigned NOT NULL,
  UNIQUE KEY `idx_fivenet_jobs_labels_users_unique` (`user_id`, `job`, `label_id`),
  KEY `fk_fivenet_jobs_labels_users_label_id` (`label_id`),
  CONSTRAINT `fk_fivenet_jobs_labels_users_label_id` FOREIGN KEY (`label_id`) REFERENCES `fivenet_jobs_labels` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_jobs_labels_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
