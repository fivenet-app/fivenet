BEGIN;

DROP TABLE `fivenet_msgs_threads_job_access`;
ALTER TABLE `fivenet_msgs_threads` DROP COLUMN `archived`;
ALTER TABLE `fivenet_msgs_threads_user_state` ADD `archived` tinyint(1) DEFAULT '0';

DELETE FROM `fivenet_permissions` WHERE `category` = 'MessengerService';

-- Table: fivenet_msgs_settings_blocks
CREATE TABLE IF NOT EXISTS `fivenet_msgs_settings_blocks` (
  `source_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`source_id`, `user_id`),
  UNIQUE KEY `idx_fivenet_msgs_settings_blocks` (`source_id`, `user_id`),
  KEY `idx_fivenet_msgs_settings_blocks_source_id` (`source_id`),
  KEY `idx_fivenet_msgs_settings_blocks_user_id` (`user_id`),
  CONSTRAINT `fk_fivenet_msgs_settings_blocks_source_id` FOREIGN KEY (`source_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_msgs_settings_blocks_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
