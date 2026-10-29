BEGIN;

CREATE TABLE IF NOT EXISTS `fivenet_property_definitions` (
  `scope` varchar(32) NOT NULL,
  `property_key` varchar(64) NOT NULL,
  `value_type` tinyint unsigned NOT NULL,
  `label` varchar(128) NOT NULL DEFAULT '',
  `description` varchar(512) DEFAULT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT 1,
  `searchable` tinyint(1) NOT NULL DEFAULT 0,
  `value_format` tinyint unsigned NOT NULL DEFAULT 0,
  `behavior` tinyint unsigned NOT NULL DEFAULT 1,
  `group_key` varchar(64) NOT NULL DEFAULT '',
  `sort_order` int NOT NULL DEFAULT 0,
  `pinned` tinyint(1) NOT NULL DEFAULT 0,
  `renderer` tinyint unsigned NOT NULL DEFAULT 1,
  `renderer_config` text DEFAULT NULL,
  `display_mode` tinyint unsigned NOT NULL DEFAULT 1,
  `implicit` tinyint(1) NOT NULL DEFAULT 0,
  `unit` varchar(32) DEFAULT NULL,
  `decimal_places` tinyint unsigned DEFAULT NULL,
  `currency` char(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`scope`, `property_key`),
  KEY `idx_property_definition_order` (`scope`, `group_key`, `sort_order`, `property_key`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_property_system_overrides` (
  `scope` varchar(32) NOT NULL,
  `property_key` varchar(64) NOT NULL,
  `enabled` tinyint(1) NOT NULL,
  `sort_order` int DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`scope`, `property_key`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_property_group_overrides` (
  `scope` varchar(32) NOT NULL,
  `group_key` varchar(64) NOT NULL,
  `sort_order` int NOT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`scope`, `group_key`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_property_groups` (
  `scope` varchar(32) NOT NULL,
  `group_key` varchar(64) NOT NULL,
  `label` varchar(128) NOT NULL,
  `sort_order` int NOT NULL DEFAULT 0,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`scope`, `group_key`),
  KEY `idx_property_group_order` (`scope`, `sort_order`, `group_key`),
  KEY `idx_property_group_deleted` (`scope`, `deleted_at`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `fivenet_user_property_values` (
  `user_id` int NOT NULL,
  `property_key` varchar(64) NOT NULL,
  `value_type` tinyint unsigned NOT NULL,
  `value_string` varchar(2048) DEFAULT NULL,
  `value_integer` bigint DEFAULT NULL,
  `value_number` double DEFAULT NULL,
  `value_bool` tinyint(1) DEFAULT NULL,
  `value_timestamp` datetime(3) DEFAULT NULL,
  `value_date` date DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`, `property_key`),
  KEY `idx_property_string` (`property_key`, `value_string`(255)),
  KEY `idx_property_integer` (`property_key`, `value_integer`),
  KEY `idx_property_number` (`property_key`, `value_number`),
  KEY `idx_property_bool` (`property_key`, `value_bool`),
  KEY `idx_property_timestamp` (`property_key`, `value_timestamp`),
  CONSTRAINT `fk_fivenet_user_property_values_user_id` FOREIGN KEY (`user_id`) REFERENCES `{{.UsersTableName}}` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Typed scalar values remain indexable and searchable. File attachments will
-- use a separate relation table when the generic attachment system is added.
CREATE TABLE IF NOT EXISTS `fivenet_vehicle_property_values` (
  `plate` varchar(12) NOT NULL,
  `property_key` varchar(64) NOT NULL,
  `value_type` tinyint unsigned NOT NULL,
  `value_string` varchar(2048) DEFAULT NULL,
  `value_integer` bigint DEFAULT NULL,
  `value_number` double DEFAULT NULL,
  `value_bool` tinyint(1) DEFAULT NULL,
  `value_timestamp` datetime(3) DEFAULT NULL,
  `value_date` date DEFAULT NULL,
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`plate`, `property_key`),
  KEY `idx_property_string` (`property_key`, `value_string`(255)),
  KEY `idx_property_integer` (`property_key`, `value_integer`),
  KEY `idx_property_number` (`property_key`, `value_number`),
  KEY `idx_property_bool` (`property_key`, `value_bool`),
  KEY `idx_property_timestamp` (`property_key`, `value_timestamp`),
  CONSTRAINT `fk_fivenet_vehicle_property_values_plate` FOREIGN KEY (`plate`) REFERENCES `fivenet_vehicle` (`plate`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

COMMIT;
