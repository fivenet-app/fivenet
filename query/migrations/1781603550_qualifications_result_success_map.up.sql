BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_qualifications_result_success_map` (
  `qualification_id` bigint(20) unsigned NOT NULL,
  `user_id` int(11) NOT NULL,
  `result_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`qualification_id`, `user_id`),
  UNIQUE KEY `idx_fivenet_qualifications_result_success_map_result_id` (`result_id`),
  KEY `idx_fivenet_qualifications_result_success_map_qualification_id` (`qualification_id`),
  KEY `idx_fivenet_qualifications_result_success_map_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_qualifications_result_success_map_qualification_id` FOREIGN KEY (`qualification_id`) REFERENCES `fivenet_qualifications` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_result_success_map_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_qualifications_result_success_map_result_id` FOREIGN KEY (`result_id`) REFERENCES `fivenet_qualifications_results` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

INSERT INTO `fivenet_qualifications_result_success_map` (`qualification_id`, `user_id`, `result_id`)
SELECT
  `qualification_id`,
  `user_id`,
  MAX(`id`) AS `result_id`
FROM `fivenet_qualifications_results`
WHERE `deleted_at` IS NULL
  AND `status` = 3
GROUP BY `qualification_id`, `user_id`
ON DUPLICATE KEY UPDATE `result_id` = VALUES(`result_id`);

COMMIT;
