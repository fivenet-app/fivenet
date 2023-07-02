BEGIN;

-- Table: fivenet_centrum_units
CREATE TABLE IF NOT EXISTS `fivenet_centrum_units` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `job` varchar(20) DEFAULT NULL,
  `name` varchar(128) NOT NULL,
  `initials` varchar(4) NOT NULL,
  `color` char(6) NOT NULL,
  `description` varchar(255) NULL DEFAULT NULL,
  `status` smallint(2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_centrum_units_users
CREATE TABLE IF NOT EXISTS `fivenet_centrum_units_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `unit_id` bigint(20) unsigned NOT NULL,
  `identifier` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_fivenet_centrum_units_users_unit_id` (`unit_id`),
  UNIQUE KEY `idx_fivenet_centrum_units_users_identifier` (`identifier`),
  CONSTRAINT `fk_fivenet_centrum_units_users_unit_id` FOREIGN KEY (`unit_id`) REFERENCES `fivenet_centrum_units` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_centrum_units_users_identifier` FOREIGN KEY (`identifier`) REFERENCES `fivenet_user_locations` (`identifier`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
