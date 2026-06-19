BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_documents_stamps_visibility_creator` (
  `target_id` bigint unsigned NOT NULL,
  `creator_job` varchar(50) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  KEY `idx_fivenet_documents_stamps_visibility_creator_lookup` (`creator_job`, `target_id`),
  CONSTRAINT `fk_fivenet_documents_stamps_visibility_creator_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_stamps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_stamps_visibility_subject` (
  `target_id` bigint unsigned NOT NULL,
  `subject_id` bigint unsigned NOT NULL,
  `access` smallint NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_documents_stamps_visibility_subject_lookup` (`subject_id`, `access`, `target_id`, `effect`),
  KEY `idx_fivenet_documents_stamps_visibility_subject_target` (`target_id`, `access`, `subject_id`, `effect`),
  CONSTRAINT `fk_fivenet_documents_stamps_visibility_subject_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents_stamps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_stamps_visibility_subject_subject_id` FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_documents_stamps_visibility_subject_effect` CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

INSERT IGNORE INTO `fivenet_documents_stamps_visibility_creator` (`target_id`, `creator_job`)
SELECT s.`id`, s.`name`
FROM `fivenet_documents_stamps` s
WHERE s.`deleted_at` IS NULL
  AND s.`name` <> '';

INSERT IGNORE INTO `fivenet_documents_stamps_visibility_subject` (`target_id`, `subject_id`, `access`, `effect`)
SELECT sa.`target_id`, sa.`subject_id`, sa.`access`, sa.`effect`
FROM `fivenet_documents_stamps_access` sa
JOIN `fivenet_documents_stamps` s ON s.`id` = sa.`target_id`
WHERE s.`deleted_at` IS NULL
  AND sa.`effect` = 1;

COMMIT;
