BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_qualifications_activity` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `qualification_id` bigint(20) unsigned NOT NULL,
  `activity_type` tinyint(2) NOT NULL,
  `actor_user_id` int(11) DEFAULT NULL,
  `target_user_id` int(11) DEFAULT NULL,
  `data` longtext DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_qualifications_activity_qualification_created` (`qualification_id`, `created_at`),
  KEY `idx_fivenet_qualifications_activity_type` (`activity_type`),
  KEY `idx_fivenet_qualifications_activity_actor_user_id` (`actor_user_id`),
  KEY `idx_fivenet_qualifications_activity_target_user_id` (`target_user_id`),
  CONSTRAINT `fk_fivenet_qualifications_activity_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_activity_actor_user_id` FOREIGN KEY (`actor_user_id`) REFERENCES `{{.UsersTableName}}` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_activity_target_user_id` FOREIGN KEY (`target_user_id`) REFERENCES `{{.UsersTableName}}` (`id`)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
