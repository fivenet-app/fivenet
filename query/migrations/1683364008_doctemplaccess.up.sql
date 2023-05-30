BEGIN;

-- Table: fivenet_documents_templates_job_access
CREATE TABLE IF NOT EXISTS `fivenet_documents_templates_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `template_id` bigint(20) unsigned NOT NULL,
  `job` varchar(20) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 1,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_templates_job_access` (`template_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_documents_templates_job_access_deleted_at` (`deleted_at`),
  KEY `idx_fivenet_documents_templates_job_access_template_id` (`template_id`),
  CONSTRAINT `fk_fivenet_documents_templates_job_access_template_id` FOREIGN KEY (`template_id`) REFERENCES `fivenet_documents_templates` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_templates_job_access_job` FOREIGN KEY (`job`) REFERENCES `jobs` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
