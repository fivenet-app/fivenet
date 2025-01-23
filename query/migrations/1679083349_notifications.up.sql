BEGIN;

-- Table: fivenet_notifications
CREATE TABLE IF NOT EXISTS `fivenet_notifications` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `read_at` datetime(3) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `job` varchar(20) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `type` smallint(2) NOT NULL,
  `content` longtext DEFAULT NULL,
  `category` smallint(2) NOT NULL,
  `data` longtext DEFAULT NULL,
  `starred` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_notifications_read_at` (`read_at`),
  KEY `idx_fivenet_notifications_user_id` (`user_id`),
  KEY `idx_fivenet_notifications_type` (`type`),
  CONSTRAINT `fk_fivenet_notifications_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
