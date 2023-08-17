BEGIN;

-- Table: fivenet_notifications
CREATE TABLE IF NOT EXISTS `fivenet_notifications` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `read_at` datetime(3) DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `title` varchar(255) NOT NULL,
  `type` varchar(128) NOT NULL,
  `content` longtext DEFAULT NULL,
  `category` smallint(2) NOT NULL,
  `data` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_notifications_read_at` (`read_at`),
  KEY `idx_fivenet_notifications_user_id` (`user_id`),
  KEY `idx_fivenet_notifications_type` (`type`),
  CONSTRAINT `fk_fivenet_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
