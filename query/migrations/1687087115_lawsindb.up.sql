BEGIN;

-- Table: fivenet_lawbooks
CREATE TABLE IF NOT EXISTS `fivenet_lawbooks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `name` varchar(128) NOT NULL,
  `description` varchar(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_lawbooks_laws
CREATE TABLE IF NOT EXISTS `fivenet_lawbooks_laws` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `lawbook_id` bigint(20) unsigned NOT NULL,
  `name` varchar(128) NOT NULL,
  `description` varchar(255) NULL,
  `fine` bigint(20) unsigned DEFAULT 0,
  `detention_time` bigint(20) unsigned DEFAULT 0,
  `stvo_points` bigint(20) unsigned DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`lawbook_id`, `name`),
  CONSTRAINT `fk_fivenet_lawbooks_laws_lawbook_id` FOREIGN KEY (`lawbook_id`) REFERENCES `fivenet_lawbooks` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
