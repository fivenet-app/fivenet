BEGIN;

CREATE TABLE `fivenet_user_notification_preferences` (
  `user_id` int(11) NOT NULL,
  `category` smallint(5) NOT NULL DEFAULT 0,
  `kind` smallint(5) NOT NULL DEFAULT 0,
  `inbox_enabled` tinyint(1) DEFAULT NULL,
  `toast_enabled` tinyint(1) DEFAULT NULL,
  `sound_enabled` tinyint(1) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `category`, `kind`),
  CONSTRAINT `fk_fivenet_user_notification_preferences_user_id` FOREIGN KEY (`user_id`) REFERENCES `fivenet_user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMIT;
