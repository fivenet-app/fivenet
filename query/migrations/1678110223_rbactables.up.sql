BEGIN;

-- Table: fivenet_permissions
CREATE TABLE IF NOT EXISTS `fivenet_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `category` varchar(128) NOT NULL,
  `name` varchar(255) NOT NULL,
  `guard_name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_permissions_category_name_unique` (`category`, `name`),
  KEY `idx_fivenet_permissions_category` (`category`),
  UNIQUE KEY `idx_fivenet_permissions_guard_name_unique` (`guard_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_attrs
CREATE TABLE IF NOT EXISTS `fivenet_attrs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `permission_id` bigint(20) unsigned NOT NULL,
  `key` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `valid_values` text DEFAULT NULL,
  `default_values` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_attrs_permission_id_key_unque` (`permission_id`,`key`),
  CONSTRAINT `fk_fivenet_attrs_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `fivenet_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_roles
CREATE TABLE IF NOT EXISTS `fivenet_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `job` varchar(50) NOT NULL,
  `grade` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_fivenet_roles_job_grade_unique` (`job`,`grade`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_role_attrs
CREATE TABLE IF NOT EXISTS `fivenet_role_attrs` (
  `role_id` bigint(20) unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `attr_id` bigint(20) unsigned NOT NULL,
  `value` longtext NOT NULL,
  `max_values` text DEFAULT NULL,
  PRIMARY KEY (`role_id`,`attr_id`),
  CONSTRAINT `fk_fivenet_role_attrs_role_id` FOREIGN KEY (`role_id`) REFERENCES `fivenet_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_role_attrs_attr_id` FOREIGN KEY (`attr_id`) REFERENCES `fivenet_attrs` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: fivenet_role_permissions
CREATE TABLE IF NOT EXISTS `fivenet_role_permissions` (
  `role_id` bigint(20) unsigned NOT NULL,
  `permission_id` bigint(20) unsigned NOT NULL,
  `val` tinyint(1) NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  CONSTRAINT `fk_fivenet_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `fivenet_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_fivenet_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `fivenet_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
