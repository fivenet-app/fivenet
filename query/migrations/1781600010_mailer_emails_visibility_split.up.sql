BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_visibility_creator` (
  `target_id` bigint unsigned NOT NULL,
  `creator_id` int NOT NULL,
  `creator_job` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`),
  KEY `idx_fivenet_mailer_emails_visibility_creator_lookup` (`creator_id`, `creator_job`, `target_id`),
  KEY `idx_fivenet_mailer_emails_visibility_creator_target` (`target_id`, `creator_id`, `creator_job`),
  CONSTRAINT `fk_fivenet_mailer_emails_visibility_creator_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_mailer_emails_visibility_subject` (
  `target_id` bigint unsigned NOT NULL,
  `subject_id` bigint unsigned NOT NULL,
  `access` smallint NOT NULL,
  `effect` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`target_id`, `subject_id`, `access`, `effect`),
  KEY `idx_fivenet_mailer_emails_visibility_subject_lookup` (`subject_id`, `access`, `target_id`, `effect`),
  KEY `idx_fivenet_mailer_emails_visibility_subject_target` (`target_id`, `access`, `subject_id`, `effect`),
  CONSTRAINT `fk_fivenet_mailer_emails_visibility_subject_target_id` FOREIGN KEY (`target_id`) REFERENCES `fivenet_mailer_emails` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_mailer_emails_visibility_subject_subject_id` FOREIGN KEY (`subject_id`) REFERENCES `fivenet_acl_subjects` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_fivenet_mailer_emails_visibility_subject_effect` CHECK (`effect` IN (0, 1))
) ENGINE=InnoDB;

INSERT IGNORE INTO `fivenet_mailer_emails_visibility_creator` (`target_id`, `creator_id`, `creator_job`)
SELECT e.`id`, e.`user_id`, ''
FROM `fivenet_mailer_emails` e
WHERE e.`deleted_at` IS NULL
  AND e.`user_id` IS NOT NULL;

INSERT IGNORE INTO `fivenet_mailer_emails_visibility_subject` (`target_id`, `subject_id`, `access`, `effect`)
SELECT ea.`target_id`, ea.`subject_id`, ea.`access`, ea.`effect`
FROM `fivenet_mailer_emails_access` ea
JOIN `fivenet_mailer_emails` e ON e.`id` = ea.`target_id`
WHERE e.`deleted_at` IS NULL
  AND ea.`effect` = 1;

COMMIT;
