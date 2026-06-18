BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_calendar_visibility_public` (
  `target_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  CONSTRAINT `fk_fivenet_calendar_visibility_public_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_calendar_visibility_creator` (
  `target_id` bigint unsigned NOT NULL,
  `creator_id` int NOT NULL,
  `creator_job` varchar(50) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  KEY `idx_fivenet_calendar_visibility_creator_lookup` (`creator_id`, `creator_job`, `target_id`),
  KEY `idx_fivenet_calendar_visibility_creator_calendar` (`target_id`, `creator_id`, `creator_job`),
  CONSTRAINT `fk_fivenet_calendar_visibility_creator_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_calendar_visibility_subject` (
  `target_id` bigint unsigned NOT NULL,
  `subject_id` bigint unsigned NOT NULL,
  `access` smallint NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_calendar_visibility_subject_lookup` (`subject_id`, `access`, `target_id`, `effect`),
  KEY `idx_fivenet_calendar_visibility_subject_calendar` (`target_id`, `access`, `subject_id`, `effect`),
  CONSTRAINT `fk_fivenet_calendar_visibility_subject_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_calendar` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_calendar_visibility_subject_subject_id` FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_calendar_visibility_subject_effect` CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

INSERT IGNORE INTO `fivenet_calendar_visibility_public` (`target_id`)
SELECT c.`id`
FROM `fivenet_calendar` c
WHERE c.`deleted_at` IS NULL
  AND c.`public` = 1;

INSERT IGNORE INTO `fivenet_calendar_visibility_creator` (`target_id`, `creator_id`, `creator_job`)
SELECT c.`id`, c.`creator_id`, c.`creator_job`
FROM `fivenet_calendar` c
WHERE c.`deleted_at` IS NULL
  AND c.`creator_id` IS NOT NULL;

INSERT IGNORE INTO `fivenet_calendar_visibility_subject` (`target_id`, `subject_id`, `access`, `effect`)
SELECT ca.`target_id`, ca.`subject_id`, ca.`access`, ca.`effect`
FROM `fivenet_calendar_access` ca
JOIN `fivenet_calendar` c ON c.`id` = ca.`target_id`
WHERE c.`deleted_at` IS NULL
  AND ca.`effect` = 1;

COMMIT;
