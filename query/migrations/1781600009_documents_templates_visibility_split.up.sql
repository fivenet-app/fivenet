BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_documents_templates_visibility_subject` (
  `target_id` bigint unsigned NOT NULL,
  `subject_id` bigint unsigned NOT NULL,
  `access` smallint NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_documents_templates_visibility_subject_lookup` (`subject_id`, `access`, `target_id`, `effect`),
  KEY `idx_fivenet_documents_templates_visibility_subject_target` (`target_id`, `access`, `subject_id`, `effect`),
  CONSTRAINT `fk_fivenet_documents_templates_visibility_subject_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_templates` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_templates_visibility_subject_subject_id` FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_documents_templates_visibility_subject_effect` CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

INSERT IGNORE INTO `fivenet_documents_templates_visibility_subject` (`target_id`, `subject_id`, `access`, `effect`)
SELECT ta.`target_id`, ta.`subject_id`, ta.`access`, ta.`effect`
FROM `fivenet_documents_templates_access` ta
JOIN `fivenet_documents_templates` t ON t.`id` = ta.`target_id`
WHERE t.`deleted_at` IS NULL
  AND ta.`effect` = 1;

COMMIT;
