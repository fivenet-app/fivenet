BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_qualifications_visibility_public` (
  `target_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  CONSTRAINT `fk_fivenet_qualifications_visibility_public_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_qualifications_visibility_subject` (
  `target_id` bigint unsigned NOT NULL,
  `subject_id` bigint unsigned NOT NULL,
  `access` smallint NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_qualifications_visibility_subject_lookup` (`subject_id`, `access`, `target_id`, `effect`),
  KEY `idx_fivenet_qualifications_visibility_subject_qualification` (`target_id`, `access`, `subject_id`, `effect`),
  CONSTRAINT `fk_fivenet_qualifications_visibility_subject_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_visibility_subject_subject_id` FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_qualifications_visibility_subject_effect` CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

INSERT IGNORE INTO `fivenet_qualifications_visibility_public` (`target_id`)
SELECT q.`id`
FROM `fivenet_qualifications` q
WHERE q.`deleted_at` IS NULL
  AND q.`public` = 1;

INSERT IGNORE INTO `fivenet_qualifications_visibility_subject` (`target_id`, `subject_id`, `access`, `effect`)
SELECT qa.`target_id`, qa.`subject_id`, qa.`access`, qa.`effect`
FROM `fivenet_qualifications_access` qa
JOIN `fivenet_qualifications` q ON q.`id` = qa.`target_id`
WHERE q.`deleted_at` IS NULL
  AND qa.`effect` = 1;

COMMIT;
