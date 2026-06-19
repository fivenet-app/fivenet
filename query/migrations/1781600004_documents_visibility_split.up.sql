BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_documents_visibility_public` (
  `target_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  CONSTRAINT `fk_fivenet_documents_visibility_public_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_visibility_creator` (
  `target_id` bigint unsigned NOT NULL,
  `creator_id` int NOT NULL,
  `creator_job` varchar(50) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  KEY `idx_fivenet_documents_visibility_creator_lookup` (`creator_id`, `creator_job`, `target_id`),
  KEY `idx_fivenet_documents_visibility_creator_document` (`target_id`, `creator_id`, `creator_job`),
  CONSTRAINT `fk_fivenet_documents_visibility_creator_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_documents_visibility_subject` (
  `target_id` bigint unsigned NOT NULL,
  `subject_id` bigint unsigned NOT NULL,
  `access` smallint NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_documents_visibility_subject_lookup` (`subject_id`, `access`, `target_id`, `effect`),
  KEY `idx_fivenet_documents_visibility_subject_document` (`target_id`, `access`, `subject_id`, `effect`),
  CONSTRAINT `fk_fivenet_documents_visibility_subject_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_documents` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_documents_visibility_subject_subject_id` FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_documents_visibility_subject_effect` CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

INSERT IGNORE INTO `fivenet_documents_visibility_public` (`target_id`)
SELECT d.`id`
FROM `fivenet_documents` d
WHERE d.`deleted_at` IS NULL
  AND d.`public` = 1;

INSERT IGNORE INTO `fivenet_documents_visibility_creator` (`target_id`, `creator_id`, `creator_job`)
SELECT d.`id`, d.`creator_id`, d.`creator_job`
FROM `fivenet_documents` d
WHERE d.`deleted_at` IS NULL
  AND d.`creator_id` IS NOT NULL;

INSERT IGNORE INTO `fivenet_documents_visibility_subject` (`target_id`, `subject_id`, `access`, `effect`)
SELECT da.`target_id`, da.`subject_id`, da.`access`, da.`effect`
FROM `fivenet_documents_access` da
JOIN `fivenet_documents` d ON d.`id` = da.`target_id`
WHERE d.`deleted_at` IS NULL
  AND da.`effect` = 1;

ALTER TABLE `fivenet_documents` DROP INDEX `idx_fivenet_documents_deleted_at`;
ALTER TABLE `fivenet_documents` DROP INDEX `idx_fivenet_documents_created_at`;
ALTER TABLE `fivenet_documents` ADD INDEX `idx_fivenet_documents_deleted_created_updated_at` (`deleted_at`, `created_at` DESC, `updated_at` DESC);

COMMIT;
