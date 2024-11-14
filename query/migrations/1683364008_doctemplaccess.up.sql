BEGIN;

-- Table: fivenet_documents_templates_job_access
CREATE TABLE IF NOT EXISTS `fivenet_documents_templates_job_access` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `template_id` bigint(20) unsigned NOT NULL,
  `job` varchar(20) NOT NULL,
  `minimum_grade` int(11) NOT NULL DEFAULT 0,
  `access` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_documents_templates_job_access` (`template_id`, `job`, `minimum_grade`),
  KEY `idx_fivenet_documents_templates_job_access_template_id` (`template_id`),
  CONSTRAINT `fk_fivenet_documents_templates_job_access_template_id` FOREIGN KEY (`template_id`) REFERENCES `fivenet_documents_templates` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
