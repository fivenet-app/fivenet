BEGIN;

-- Table: fivenet_audit_log
CREATE TABLE IF NOT EXISTS `fivenet_audit_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `user_id` int(11) NOT NULL,
  `user_job` varchar(20) NOT NULL,
  `target_user_id` int(11) DEFAULT NULL,
  `service` varchar(255) NOT NULL,
  `method` varchar(255) NOT NULL,
  `state` smallint(2) NOT NULL,
  `data` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_audit_log_user_id` (`user_id`),
  KEY `idx_fivenet_audit_log_user_job` (`user_job`),
  KEY `idx_fivenet_audit_log_created_at` (`created_at`),
  KEY `idx_fivenet_audit_log_service` (`service`),
  KEY `idx_fivenet_audit_log_method` (`method`),
  FULLTEXT KEY `idx_fivenet_audit_log_data` (`data`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
