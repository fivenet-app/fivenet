BEGIN;

-- Table: arpanet_permissions
CREATE TABLE IF NOT EXISTS `arpanet_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `name` varchar(255) NOT NULL,
  `guard_name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_arpanet_permissions_guard_name` (`guard_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_roles
CREATE TABLE IF NOT EXISTS `arpanet_roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `name` varchar(255) NOT NULL,
  `guard_name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_arpanet_roles_guard_name` (`guard_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_role_permissions
CREATE TABLE IF NOT EXISTS `arpanet_role_permissions` (
  `role_id` bigint(20) unsigned NOT NULL,
  `permission_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_arpanet_role_permissions_permission` (`permission_id`),
  CONSTRAINT `fk_arpanet_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `arpanet_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `arpanet_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_roles
CREATE TABLE IF NOT EXISTS `arpanet_user_roles` (
  `user_id` int(11) NOT NULL,
  `role_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  CONSTRAINT `fk_arpanet_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `arpanet_roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table: arpanet_user_permissions
CREATE TABLE IF NOT EXISTS `arpanet_user_permissions` (
  `user_id` int(11) NOT NULL,
  `permission_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`permission_id`),
  CONSTRAINT `fk_arpanet_user_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `arpanet_permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_arpanet_user_permissions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
