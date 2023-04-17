BEGIN;

-- Table: fivenet_audit_log
CREATE TABLE IF NOT EXISTS `fivenet_audit_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `user_id` int(11) NOT NULL,
  `user_job` varchar(20) NOT NULL,
  `target_job` varchar(20) DEFAULT NULL,
  `service` varchar(255) NOT NULL,
  `method` varchar(255) NOT NULL,
  `state` smallint(2) NOT NULL,
  `data` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_fivenet_audit_log_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
